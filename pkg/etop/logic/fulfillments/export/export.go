package export

import (
	"archive/zip"
	"bytes"
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"unicode/utf8"

	"etop.vn/backend/pkg/services/shipping/modely"

	"github.com/golang/protobuf/jsonpb"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/idemp"
	"etop.vn/backend/pkg/common/l"
	"etop.vn/backend/pkg/common/redis"
	cmService "etop.vn/backend/pkg/common/service"
	"etop.vn/backend/pkg/common/sql/core"
	"etop.vn/backend/pkg/etop/eventstream"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/etop/sqlstore"

	pbcm "etop.vn/backend/pb/common"
	pbshop "etop.vn/backend/pb/etop/shop"
)

var ll = l.New()
var idempgroup *idemp.RedisGroup
var errAborted = cm.Error(cm.Aborted, "abort", nil)
var config Config
var publisher eventstream.Publisher

const PathShopFulfillments = "shop/fulfillments"
const BaseRowsErrors = 10

type Config struct {
	UrlPrefix string
	DirExport string
}

type ExportOption struct {
	Delimiter rune
	ExcelMode bool
}

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

var marshaler = jsonpb.Marshaler{OrigName: true, EmitDefaults: true}

func exportFulfillmentsAndReportProgress(
	cleanup func(),
	exportResult *model.ExportAttempt, bareFilename string, exportOpts ExportOption,
	total int, rows *sql.Rows, opts core.Opts,
) (_err error) {

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
			exportResult.Status = model.S4Negative
			exportResult.Error = pbcm.ErrorToModel(pbcm.PbError(_err))

			event := eventstream.Event{
				Type:      "export/error",
				AccountID: exportResult.AccountID,
				UserID:    exportResult.UserID,
				Payload:   pbcm.PbError(_err),
			}
			publisher.Publish(event)
			return
		}

		exportResult.Status = model.S4Positive

		event := eventstream.Event{
			Type:      "export/ok",
			AccountID: exportResult.AccountID,
			UserID:    exportResult.UserID,
			Payload:   pbshop.PbExportAttempt(exportResult),
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
	result := make(chan *pbshop.ExportStatusItem, BaseRowsErrors)
	go exportFulfillments(
		exportCtx, exportID, exportOpts, fileWriter, result,
		total, rows, opts)

	// send progress to client
	var statusItem *pbshop.ExportStatusItem
	for statusItem = range result {
		buf := &bytes.Buffer{}
		if err := marshaler.Marshal(buf, statusItem); err != nil {
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
		exportResult.NTotal = int(statusItem.ProgressMax)
		exportResult.NExported = int(statusItem.ProgressValue)
		exportResult.NError = len(errs)
	}
	return nil
}

const exportFulfillmentLines = false

func exportFulfillments(
	ctx context.Context, id string, exportOpts ExportOption, output io.Writer,
	result chan<- *pbshop.ExportStatusItem,
	total int, rows *sql.Rows, opts core.Opts,
) (_err error) {

	var count, countError, line int
	maxErrors := total / 100
	if maxErrors > BaseRowsErrors {
		maxErrors = BaseRowsErrors
	}

	makeProgress := func() *pbshop.ExportStatusItem {
		return &pbshop.ExportStatusItem{
			Id:            id,
			ProgressMax:   int32(total),
			ProgressValue: int32(count),
			ProgressError: int32(countError),
		}
	}
	handleError := func(err error) bool {
		switch err {
		case nil:
			return true
		case errAborted:
			return false
		default:
			countError++
			statusItem := makeProgress()
			statusItem.Error = pbcm.PbError(err)
			if countError >= maxErrors {
				err = cm.Errorf(cm.Aborted, nil, "Quá nhiều lỗi xảy ra")
				statusItem.Error = pbcm.PbError(err)
				result <- statusItem
			}
			return false
		}
	}
	defer func() {
		handleError(_err)
		if _err == nil {
			result <- makeProgress()
		}
		close(result)
	}()

	// BOM is required for Excel to recognize UTF-8 file
	if _, err := writeBOM(output); err != nil {
		return err
	}

	csvWriter := csv.NewWriter(output)
	csvWriter.Comma = exportOpts.Delimiter
	defer func() {
		// calling Flush() is necessary
		csvWriter.Flush()
		if err := csvWriter.Error(); _err == nil && err != nil {
			_err = err
		}
	}()

	tableWriter := NewTableWriter(csvWriter, noempty)
	var ffm modely.FulfillmentExtended
	if err := buildTableFulfillment(csvWriter, tableWriter, exportOpts, &count, &line, &ffm); err != nil {
		return err
	}

	result <- makeProgress()
	lastTime := time.Now()
	for {
		select {
		case <-ctx.Done():
			return errAborted

		default:
			if !rows.Next() {
				return rows.Err()
			}
			args := ffm.SQLScanArgs(opts)
			if err := rows.Scan(args...); err != nil {
				if !handleError(err) {
					return errAborted
				}
			}

			// increase count and line for each line written to csv file
			count++
			if exportFulfillmentLines {
				for line = 0; line < len(ffm.Lines); line++ {
					if err := tableWriter.WriteRow(); err != nil {
						if !handleError(err) {
							return errAborted
						}
					}
				}
			} else {
				if err := tableWriter.WriteRow(); err != nil {
					if !handleError(err) {
						return errAborted
					}
				}
			}
			if now := time.Now(); now.Sub(lastTime) > 100*time.Millisecond {
				result <- makeProgress()
				lastTime = now
			}
		}
	}
}

func writeBOM(w io.Writer) (n int, err error) {
	utf8bom := make([]byte, 3)
	utf8.EncodeRune(utf8bom, '\uFEFF')
	return w.Write(utf8bom)
}

func buildTableFulfillment(csvWriter *csv.Writer, w *TableWriter, exportOpts ExportOption, count *int, line *int, ffm *modely.FulfillmentExtended) error {
	var formatAsText = func(s string) string { return s }
	if exportOpts.ExcelMode {
		formatAsText = formatAsTextForExcel
	}

	var firstHeaders []string
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("#", func() string {
		return strconv.Itoa(*count)
	})

	firstHeaders = append(firstHeaders, "GIAO HÀNG")
	w.AddColumn("Mã đơn hàng", func() string { return formatAsText(ffm.Order.Code) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Mã nội bộ", func() string { return formatAsText(cm.Coalesce(ffm.Order.EdCode, ffm.Order.ExternalOrderID)) })
	if exportFulfillmentLines {
		firstHeaders = append(firstHeaders, "Ghi chú: (nt) = như trên")
	} else {
		firstHeaders = append(firstHeaders, "")
	}
	w.AddColumn("Mã giao hàng", func() string { return formatAsText(firstLine(line, ffm.ShippingCode)) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Nhà vận chuyển", func() string { return firstLine(line, ffm.ShippingProvider.Label()) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Ngày tạo", func() string { return firstLine(line, formatDate(ffm.CreatedAt)) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Bảo hiểm", func() string { return firstLine(line, labelBool(ffm.IncludeInsurance)) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Lý do huỷ", func() string { return firstLine(line, ffm.CancelReason) })

	firstHeaders = append(firstHeaders, "TRẠNG THÁI")
	w.AddColumn("Giao hàng", func() string { return firstLine(line, ffm.ShippingState.Text()) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Ngày lấy", func() string { return firstLine(line, formatDate(ffm.ShippingPickingAt)) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Ngày giao", func() string { return firstLine(line, formatDate(ffm.ShippingDeliveredAt)) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Ngày trả hàng", func() string { return firstLine(line, formatDate(ffm.ShippingReturnedAt)) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Thu hộ", func() string {
		return firstLine(line, labelCODTransferAt(ffm.EtopPaymentStatus, ffm.CODEtopTransferedAt))
	})
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Ngày thanh toán tiền thu hộ", func() string {
		return firstLine(line, formatDate(ffm.CODEtopTransferedAt))
	})
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Phí giao hàng", func() string {
		return firstLine(line, labelShippingFeeTransferAt(ffm.EtopPaymentStatus, ffm.ShippingFeeShopTransferedAt))
	})
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Ngày thanh toán phí giao hàng", func() string {
		return firstLine(line, formatDate(ffm.ShippingFeeShopTransferedAt))
	})
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Mã phiên thanh toán", func() string {
		if ffm.MoneyTransactionShipping == nil {
			return ""
		}
		return formatAsText(firstLine(line, ffm.MoneyTransactionShipping.Code))
	})

	firstHeaders = append(firstHeaders, "THÔNG TIN LẤY HÀNG")
	w.AddColumn("Tên người gửi", func() string { return firstLine(line, ffm.AddressFrom.GetFullName()) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Số điện thoại", func() string { return formatAsText(firstLine(line, ffm.AddressFrom.GetPhone())) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Tỉnh thành", func() string { return firstLine(line, ffm.AddressFrom.GetProvince()) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Quận huyện", func() string { return firstLine(line, ffm.AddressFrom.GetDistrict()) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Phường xã", func() string { return firstLine(line, ffm.AddressFrom.GetWard()) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Địa chỉ", func() string { return firstLine(line, ffm.AddressFrom.GetShortAddress()) })

	firstHeaders = append(firstHeaders, "THÔNG TIN NHẬN HÀNG")
	w.AddColumn("Tên người nhận", func() string { return firstLine(line, ffm.AddressTo.GetFullName()) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Số điện thoại", func() string { return formatAsText(firstLine(line, ffm.AddressTo.GetPhone())) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Tỉnh thành", func() string { return firstLine(line, ffm.AddressTo.GetProvince()) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Quận huyện", func() string { return firstLine(line, ffm.AddressTo.GetDistrict()) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Phường xã", func() string { return firstLine(line, ffm.AddressTo.GetWard()) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Địa chỉ", func() string { return firstLine(line, ffm.AddressTo.GetShortAddress()) })

	if exportFulfillmentLines {
		firstHeaders = append(firstHeaders, "THÔNG TIN SẢN PHẨM")
		w.AddColumn("Mã phiên bản sản phẩm", func() string { return formatAsText(ffm.Order.Lines[*line].Code) })
		firstHeaders = append(firstHeaders, "")
		w.AddColumn("Tên sản phẩm", func() string { return ffm.Order.Lines[*line].ProductName })
		firstHeaders = append(firstHeaders, "")
		w.AddColumn("Số lượng", func() string { return formatInt(ffm.Order.Lines[*line].Quantity) })
		firstHeaders = append(firstHeaders, "")
		w.AddColumn("Đơn giá (₫)", func() string { return formatInt(ffm.Order.Lines[*line].RetailPrice) })
		firstHeaders = append(firstHeaders, "")
		w.AddColumn("Thành tiền (trước giảm giá)", func() string { return formatInt(ffm.Order.Lines[*line].GetRetailAmount()) })
		firstHeaders = append(firstHeaders, "")
		w.AddColumn("Giảm giá (₫)", func() string { return formatInt(ffm.Order.Lines[*line].GetTotalDiscount()) })
		firstHeaders = append(firstHeaders, "")
		w.AddColumn("Thành tiền (sau giảm giá)", func() string { return formatInt(ffm.Order.Lines[*line].GetPaymentAmount()) })
	} else {
		firstHeaders = append(firstHeaders, "THÔNG TIN SẢN PHẨM")
		w.AddColumn("Sản phẩm", func() string { return ffm.Order.Lines.GetSummary() })
	}

	firstHeaders = append(firstHeaders, "GIÁ TRỊ ĐƠN HÀNG")
	w.AddColumn("Tổng số lượng sản phẩm", func() string { return firstLine(line, formatInt(ffm.Order.Lines.GetTotalItems())) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Tổng khối lượng (g)", func() string { return firstLine(line, formatInt(ffm.Order.Lines.GetTotalWeight())) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Tổng tiền hàng (trước giảm giá)", func() string { return firstLine(line, formatInt(ffm.Order.Lines.GetTotalRetailAmount())) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Tổng tiền hàng (sau giảm giá)", func() string { return firstLine(line, formatInt(ffm.Order.Lines.GetTotalPaymentAmount())) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Giảm giá đơn hàng (₫)", func() string { return firstLine(line, formatInt(ffm.Order.OrderDiscount)) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Phí giao hàng tính cho khách", func() string { return firstLine(line, formatInt(ffm.Order.FeeLines.GetShippingFee())) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Thuế tính cho khách", func() string { return firstLine(line, formatInt(ffm.Order.FeeLines.GetTax())) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Phí khác tính cho khách", func() string { return firstLine(line, formatInt(ffm.Order.FeeLines.GetOtherFee())) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Tổng phí tính cho khách", func() string { return firstLine(line, formatInt(ffm.Order.FeeLines.GetTotalFee())) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Tổng tiền thanh toán (bao gồm phí và giảm giá)", func() string {
		total := ffm.Order.Lines.GetTotalPaymentAmount() - ffm.Order.OrderDiscount + ffm.Order.FeeLines.GetTotalFee()
		return firstLine(line, formatInt(total))
	})
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Giá trị bảo hiểm", func() string {
		if !ffm.IncludeInsurance {
			return ""
		}
		return firstLine(line, formatInt(ffm.BasketValue))
	})

	firstHeaders = append(firstHeaders, "TIỀN THU HỘ VÀ CHI PHÍ")
	w.AddColumn("Tiền thu hộ", func() string { return firstLine(line, formatInt(ffm.TotalCODAmount)) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Phí giao hàng", func() string { return firstLine(line, formatIntBeforeVAT(ffm.ShippingFeeMain)) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Phí trả hàng", func() string { return firstLine(line, formatIntBeforeVAT(ffm.ShippingFeeReturn)) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Phí vượt cân", func() string { return firstLine(line, formatIntBeforeVAT(ffm.ShippingFeeAdjustment)) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Phí bảo hiểm", func() string { return firstLine(line, formatIntBeforeVAT(ffm.ShippingFeeInsurance)) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Phí thu hộ", func() string { return firstLine(line, formatIntBeforeVAT(ffm.ShippingFeeCODS)) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Phí đổi thông tin", func() string { return firstLine(line, formatIntBeforeVAT(ffm.ShippingFeeInfoChange)) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Phí khác", func() string { return firstLine(line, formatIntBeforeVAT(ffm.ShippingFeeOther)) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Tổng phí (trước VAT)", func() string { return firstLine(line, formatIntBeforeVAT(ffm.ShippingFeeShop)) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("VAT (10%)", func() string { return firstLine(line, formatIntVAT(ffm.ShippingFeeShop)) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Tổng phí (sau VAT)", func() string { return firstLine(line, formatInt(ffm.ShippingFeeShop)) })

	if err := csvWriter.Write(firstHeaders); err != nil {
		return err
	}
	return w.WriteHeader()
}

func noempty(s string) string {
	if s == "" {
		return " "
	}
	return s
}

// Excel and GoogleSheet auto convert fields look like number to text. We have
// to append a "\t" to work around. See https://stackoverflow.com/a/15107122
func formatAsTextForExcel(s string) string {
	if s == "" {
		return ""
	}
	return fmt.Sprintf("=%q", s)
}

func formatDateShort(time time.Time) string {
	year, month, day := time.Date()
	return fmt.Sprintf("%04d%02d%02d", year, month, day)
}

func formatDateTimeShort(time time.Time) string {
	year, month, day := time.Date()
	hour, min, _ := time.Clock()
	return fmt.Sprintf("%04d%02d%02d.%02d%02d", year, month, day, hour, min)
}

func formatDate(time time.Time) string {
	year, month, day := time.Date()
	if year <= 1990 {
		return ""
	}
	return fmt.Sprintf("%04d-%02d-%02d", year, month, day)
}

func firstLine(line *int, v string) string {
	if *line == 0 {
		return v
	}
	if v == "" {
		return ""
	}
	return "(nt)"
}

func labelCODTransferAt(paymentStatus model.Status4, transferedAt time.Time) string {
	switch {
	case !transferedAt.IsZero():
		return "Đã chuyển tiền thu hộ"
	case paymentStatus == model.S4SuperPos:
		return "Chưa chuyển tiền thu hộ"
	default:
		return ""
	}
}

func labelShippingFeeTransferAt(paymentStatus model.Status4, transferedAt time.Time) string {
	switch {
	case !transferedAt.IsZero():
		return "Đã thanh toán"
	case paymentStatus == model.S4SuperPos:
		return "Chưa thanh toán"
	default:
		return ""
	}
}

func labelBool(b bool) string {
	if b {
		return "Có"
	}
	return "Không"
}

func formatInt(i int) string {
	return strconv.Itoa(i)
}

func formatIntBeforeVAT(i int) string {
	v := i * 10 / 11
	return strconv.Itoa(v)
}

func formatIntVAT(i int) string {
	v := i / 11
	return strconv.Itoa(v)
}
