package export

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"unicode/utf8"

	"etop.vn/api/top/int/shop"
	pbcm "etop.vn/api/top/types/common"
	"etop.vn/api/top/types/etc/status4"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/cmapi"
	"etop.vn/backend/pkg/common/apifw/idemp"
	cmService "etop.vn/backend/pkg/common/apifw/service"
	"etop.vn/backend/pkg/common/redis"
	"etop.vn/backend/pkg/common/sql/sq/core"
	"etop.vn/backend/pkg/etop/api/convertpb"
	"etop.vn/backend/pkg/etop/eventstream"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/etop/sqlstore"
	"etop.vn/common/jsonx"
	"etop.vn/common/l"
)

var ll = l.New()
var idempgroup *idemp.RedisGroup
var ErrAborted = cm.Error(cm.Aborted, "abort", nil)
var config Config
var publisher eventstream.Publisher

const PathShopFulfillments = "shop/fulfillments"
const PathShopOrders = "shop/orders"
const BaseRowsErrors = 10

type Config struct {
	UrlPrefix string
	DirExport string
}

type ExportOption struct {
	Delimiter rune
	ExcelMode bool
}

type rowsInterface interface {
	Err() error
	Next() bool
	Scan(args ...interface{}) error
	Close() error
}

type ExportFunction func(ctx context.Context, id string, exportOpts ExportOption, output io.Writer,
	result chan<- *shop.ExportStatusItem,
	total int, rows rowsInterface, opts core.Opts) (_err error)

func Init(sd cmService.Shutdowner, rd redis.Store, p eventstream.Publisher, cfg Config) {
	idempgroup = idemp.NewRedisGroup(rd, "export", 60)
	sd.Register(idempgroup.Shutdown)
	publisher = p

	if cfg.DirExport == "" || cfg.UrlPrefix == "" {
		panic("must provide export dir and url prefix ")
	}

	var err error
	config = cfg
	config.DirExport, err = verifyDir(config.DirExport)
	if err != nil {
		ll.Panic("invalid export config", l.Error(err))
	}
}

func verifyDir(dir string) (absPath string, err error) {
	if absPath, err = filepath.Abs(dir); err != nil {
		ll.S.Panicf("invalid directory: %v", err)
	}

	info, err := os.Stat(absPath)
	if err != nil {
		return
	}
	if !info.IsDir() {
		err = cm.Errorf(cm.InvalidArgument, nil, "must be a directory: %v", absPath)
	}
	return
}

func ensureDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}

func exportAndReportProgress(
	cleanup func(),
	exportResult *model.ExportAttempt, bareFilename string, exportOpts ExportOption,
	total int, rows rowsInterface, opts core.Opts,
	exportFunction ExportFunction) (_err error) {

	exportResult.StartedAt = time.Now()

	var errs []*pbcm.Error
	defer cleanup()
	defer func() {
		exportResult.DoneAt = time.Now()
		exportResult.ExpiresAt = exportResult.DoneAt.Add(7 * 24 * time.Hour).Truncate(24 * time.Hour)

		if err := sqlstore.ExportAttempt(context.Background()).
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
			publisher.Publish(event)
			return
		}

		exportResult.Status = status4.P

		event := eventstream.Event{
			Type:      "export/ok",
			AccountID: exportResult.AccountID,
			UserID:    exportResult.UserID,
			Payload:   convertpb.PbExportAttempt(exportResult),
		}
		publisher.Publish(event)
	}()
	defer rows.Close()

	// prepare output file as zip
	exportID := exportResult.ID
	midPath, zipFilename := filepath.Split(exportResult.StoredFile)

	// create .zip file on disk
	dirPath := filepath.Join(config.DirExport, midPath)
	if err := ensureDir(dirPath); err != nil {
		return err
	}
	zipFilePath := filepath.Join(config.DirExport, midPath, zipFilename)
	zipFile, err := os.OpenFile(zipFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	// prepare zip stream
	zipWriter := zip.NewWriter(zipFile)
	defer func() {
		// close the zip file, record the error
		if err := zipWriter.Close(); _err == nil && err != nil {
			_err = err
		}
		if err := zipFile.Close(); _err == nil && err != nil {
			_err = err
		}

		// update response
		exportResult.ID = exportID
		if _err == nil {
			exportResult.DownloadURL = config.UrlPrefix + "/" + exportResult.StoredFile

		} else {
			exportResult.FileName = ""
			exportResult.StoredFile = ""

			// remove the file if there is any error
			removeError := os.Remove(zipFilePath)
			if removeError != nil {
				ll.Error("error removing file", l.String("path", zipFilePath), l.Error(err))
			}
		}
	}()

	fileWriter, err := zipWriter.Create(bareFilename + ".csv")
	if err != nil {
		return err
	}

	exportCtx, exportCancel := context.WithCancel(context.Background())
	defer exportCancel()

	// perform export in another goroutine
	result := make(chan *shop.ExportStatusItem, BaseRowsErrors)
	go exportFunction(
		exportCtx, exportID, exportOpts, fileWriter, result,
		total, rows, opts)

	// send progress to client
	var statusItem *shop.ExportStatusItem
	for statusItem = range result {
		buf := &bytes.Buffer{}
		if err := jsonx.MarshalTo(buf, statusItem); err != nil {
			panic(err)
		}
		event := eventstream.Event{
			Type:      "export/progress",
			Global:    false,
			AccountID: exportResult.AccountID,
			UserID:    exportResult.UserID,
			Payload:   buf.Bytes(),
		}
		publisher.Publish(event)
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
		return "Đã chuyển tiền thu hộ"
	case paymentStatus == status4.S:
		return "Chưa chuyển tiền thu hộ"
	default:
		return ""
	}
}

func LabelShippingFeeTransferAt(paymentStatus status4.Status, transferedAt time.Time) string {
	switch {
	case !transferedAt.IsZero():
		return "Đã thanh toán"
	case paymentStatus == status4.S:
		return "Chưa thanh toán"
	default:
		return ""
	}
}

func LabelBool(b bool) string {
	if b {
		return "Có"
	}
	return "Không"
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
