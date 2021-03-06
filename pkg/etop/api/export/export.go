package export

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"math"
	"path/filepath"
	"strconv"
	"time"
	"unicode/utf8"

	apishop "o.o/api/top/int/shop"
	pbcm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/status4"
	identitymodel "o.o/backend/com/main/identity/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/common/sql/sq/core"
	"o.o/backend/pkg/common/storage"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/eventstream"
	"o.o/backend/pkg/etop/model"
	"o.o/common/jsonx"
	"o.o/common/l"
)

var ll = l.New()
var ErrAborted = cm.Error(cm.Aborted, "abort", nil)

const PathShopFulfillments = "shop/fulfillments"
const PathShopOrders = "shop/orders"
const PathShopCallLogs = "shop/calllogs"
const BaseRowsErrors = 10

type ConfigDirs struct {
	Export storage.DirConfig `yaml:"export"`
}

type ExportOption struct {
	Delimiter rune
	ExcelMode bool
}

type RowsInterface interface {
	Err() error
	Next() bool
	Scan(args ...interface{}) error
	Close() error
}

type ExportFunction func(ctx context.Context, id string, shop *identitymodel.Shop,
	exportOpts ExportOption, output io.Writer,
	result chan<- *apishop.ExportStatusItem,
	total int, rows RowsInterface, opts core.Opts) (_err error)

func (s *Service) exportAndReportProgress(
	cleanup func(),
	shop *identitymodel.Shop,
	exportResult *model.ExportAttempt, bareFilename string, exportOpts ExportOption,
	total int, rows RowsInterface, opts core.Opts,
	exportFunction ExportFunction) (_err error) {

	exportResult.StartedAt = time.Now()

	var errs []*pbcm.Error
	defer cleanup()
	defer func() {
		exportResult.DoneAt = time.Now()
		exportResult.ExpiresAt = exportResult.DoneAt.Add(7 * 24 * time.Hour).Truncate(24 * time.Hour)

		if err := s.exportAttemptStore(context.Background()).
			UpdateByID(exportResult.ID, exportResult); err != nil {
			ll.Error("error updating import attempt", l.Error(err))
		}
	}()
	defer func() {
		// record the last aborted error
		if _err == nil {
			if len(errs) > 0 && errs[len(errs)-1].Code == "aborted" {
				_err = errs[len(errs)-1]
			}
		}

		if _err != nil {
			exportResult.Status = status4.N
			exportResult.Error = cmapi.ErrorToModel(cmapi.PbError(_err))

			event := eventstream.Event{
				Type:      "export/error",
				AccountID: exportResult.AccountID,
				UserID:    exportResult.UserID,
				Payload:   cmapi.PbError(_err),
			}
			s.publisher.Publish(event)
			return
		}

		exportResult.Status = status4.P

		event := eventstream.Event{
			Type:      "export/ok",
			AccountID: exportResult.AccountID,
			UserID:    exportResult.UserID,
			Payload:   convertpb.PbExportAttempt(exportResult),
		}
		s.publisher.Publish(event)
	}()
	defer rows.Close()

	// prepare output file as zip
	exportID := exportResult.ID
	midPath, zipFilename := filepath.Split(exportResult.StoredFile)

	// create .zip file on disk
	zipFilePath := filepath.Join(s.config.Export.Path, midPath, zipFilename)
	zipFile, err := s.storageBucket.OpenForWrite(context.Background(), zipFilePath)
	if err != nil {
		return err
	}

	// prepare zip stream
	zipWriter := zip.NewWriter(zipFile)
	defer func() {
		// close the zip file, record the error
		if err2 := zipWriter.Close(); _err == nil && err2 != nil {
			_err = err2
		}
		if err2 := zipFile.Close(); _err == nil && err2 != nil {
			_err = err2
		}

		// update response
		exportResult.ID = exportID
		if _err == nil {
			exportResult.DownloadURL = s.config.Export.URLPrefix + "/" + exportResult.StoredFile

		} else {
			exportResult.FileName = ""
			exportResult.StoredFile = ""

			// TODO(vu): remove the file if there is any error
		}
	}()

	fileWriter, err := zipWriter.Create(bareFilename + ".csv")
	if err != nil {
		return err
	}

	exportCtx, exportCancel := context.WithCancel(context.Background())
	defer exportCancel()

	// perform export in another goroutine
	result := make(chan *apishop.ExportStatusItem, BaseRowsErrors)
	go exportFunction(
		exportCtx, exportID, shop, exportOpts, fileWriter, result,
		total, rows, opts)

	// send progress to client
	var statusItem *apishop.ExportStatusItem
	for statusItem = range result {
		buf := &bytes.Buffer{}
		if err2 := jsonx.MarshalTo(buf, statusItem); err2 != nil {
			panic(err2)
		}
		event := eventstream.Event{
			Type:      "export/progress",
			Global:    false,
			AccountID: exportResult.AccountID,
			UserID:    exportResult.UserID,
			Payload:   buf.Bytes(),
		}
		s.publisher.Publish(event)
	}

	// store the last status item as export result
	if statusItem != nil {
		// progress_max may not always be correct, but still good enough
		exportResult.NTotal = statusItem.ProgressMax
		exportResult.NExported = statusItem.ProgressValue
		exportResult.NError = len(errs)
	}
	return nil
}

func WriteBOM(w io.Writer) (n int, err error) {
	utf8bom := make([]byte, 3)
	utf8.EncodeRune(utf8bom, '\uFEFF')
	return w.Write(utf8bom)
}

func Noempty(s string) string {
	if s == "" {
		return " "
	}
	return s
}

// Excel and GoogleSheet auto convert fields look like number to text. We have
// to append a "\t" to work around. See https://stackoverflow.com/a/15107122
func FormatAsTextForExcel(s string) string {
	if s == "" {
		return ""
	}
	return fmt.Sprintf("=%q", s)
}

func FormatDateShort(time time.Time) string {
	year, month, day := time.Date()
	return fmt.Sprintf("%04d%02d%02d", year, month, day)
}

func FormatDateTimeShort(time time.Time) string {
	year, month, day := time.Date()
	hour, min, _ := time.Clock()
	return fmt.Sprintf("%04d%02d%02d.%02d%02d", year, month, day, hour, min)
}

func FormatDate(time time.Time) string {
	year, month, day := time.Date()
	if year <= 1990 {
		return ""
	}
	return fmt.Sprintf("%04d-%02d-%02d", year, month, day)
}

func FormatDateTime(t time.Time) string {
	year, month, day := t.Date()
	hour, min, second := t.Clock()
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", year, month, day, hour, min, second)
}

func FormatDuration(t time.Duration) string {
	// hh:mm:ss
	_t := t.Seconds()
	hours := math.Floor(_t / 3600)
	_t -= hours * 3600
	minutes := math.Floor(_t / 60)
	_t -= minutes * 60
	seconds := _t
	return fmt.Sprintf("%02d:%02d:%02d", int(hours), int(minutes), int(seconds))
}

func FirstLine(line *int, v string) string {
	if *line == 0 {
		return v
	}
	if v == "" {
		return ""
	}
	return "(nt)"
}

func LabelCODTransferAt(paymentStatus status4.Status, transferedAt time.Time) string {
	switch {
	case !transferedAt.IsZero():
		return "???? chuy???n ti???n thu h???"
	case paymentStatus == status4.S:
		return "Ch??a chuy???n ti???n thu h???"
	default:
		return ""
	}
}

func LabelShippingFeeTransferAt(paymentStatus status4.Status, transferedAt time.Time) string {
	switch {
	case !transferedAt.IsZero():
		return "???? thanh to??n"
	case paymentStatus == status4.S:
		return "Ch??a thanh to??n"
	default:
		return ""
	}
}

func LabelBool(b bool) string {
	if b {
		return "C??"
	}
	return "Kh??ng"
}

func FormatInt(i int) string {
	return strconv.Itoa(i)
}

func FormatIntBeforeVAT(i int) string {
	v := i * 10 / 11
	return strconv.Itoa(v)
}

func FormatIntVAT(i int) string {
	v := i / 11
	return strconv.Itoa(v)
}
