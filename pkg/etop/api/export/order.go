package export

import (
	"context"
	"database/sql"
	"encoding/csv"
	"io"
	"strconv"
	"time"

	pbshop "etop.vn/api/pb/etop/shop"
	orderingmodely "etop.vn/backend/com/main/ordering/modely"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmapi"
	"etop.vn/backend/pkg/common/sq/core"
	"etop.vn/backend/pkg/etop/model"
)

const exportOrderLines = false

func ExportOrders(
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
			ProgressMax:   int(total),
			ProgressValue: int(count),
			ProgressError: int(countError),
		}
	}
	handleError := func(err error) bool {
		switch err {
		case nil:
			return true
		case ErrAborted:
			return false
		default:
			countError++
			statusItem := makeProgress()
			statusItem.Error = cmapi.PbError(err)
			if countError >= maxErrors {
				err = cm.Errorf(cm.Aborted, nil, "Quá nhiều lỗi xảy ra")
				statusItem.Error = cmapi.PbError(err)
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
	if _, err := WriteBOM(output); err != nil {
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

	tableWriter := NewTableWriter(csvWriter, Noempty)
	var ffm orderingmodely.OrderExtended
	if err := buildOrder(csvWriter, tableWriter, exportOpts, &count, &line, &ffm); err != nil {
		return err
	}

	result <- makeProgress()
	lastTime := time.Now()
	for {
		select {
		case <-ctx.Done():
			return ErrAborted

		default:
			if !rows.Next() {
				return rows.Err()
			}
			args := ffm.SQLScanArgs(opts)
			if err := rows.Scan(args...); err != nil {
				if !handleError(err) {
					return ErrAborted
				}
			}

			// increase count and line for each line written to csv file
			count++
			if exportOrderLines {
				for line = 0; line < len(ffm.Lines); line++ {
					if err := tableWriter.WriteRow(); err != nil {
						if !handleError(err) {
							return ErrAborted
						}
					}
				}
			} else {
				if err := tableWriter.WriteRow(); err != nil {
					if !handleError(err) {
						return ErrAborted
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

func buildOrder(csvWriter *csv.Writer, w *TableWriter, exportOpts ExportOption, count *int, line *int, ffm *orderingmodely.OrderExtended) error {
	var formatAsText = func(s string) string { return s }
	if exportOpts.ExcelMode {
		formatAsText = FormatAsTextForExcel
	}

	var firstHeaders []string
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("#", func() string {
		return strconv.Itoa(*count)
	})

	firstHeaders = append(firstHeaders, "ĐƠN HÀNG")
	w.AddColumn("Mã đơn hàng", func() string { return formatAsText(ffm.Code) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Mã nội bộ", func() string { return formatAsText(cm.Coalesce(ffm.EdCode, ffm.ExternalOrderID)) })
	if exportOrderLines {
		firstHeaders = append(firstHeaders, "Ghi chú: (nt) = như trên")
	} else {
		firstHeaders = append(firstHeaders, "")
	}
	w.AddColumn("Ngày tạo", func() string { return FirstLine(line, FormatDate(ffm.CreatedAt)) })

	if exportOrderLines {
		firstHeaders = append(firstHeaders, "")
		w.AddColumn("Mã phiên bản sản phẩm", func() string { return formatAsText(ffm.Lines[*line].Code) })
		firstHeaders = append(firstHeaders, "")
		w.AddColumn("Tên sản phẩm", func() string { return ffm.Lines[*line].ProductName })
		firstHeaders = append(firstHeaders, "")
		w.AddColumn("Số lượng", func() string { return FormatInt(ffm.Lines[*line].Quantity) })
		firstHeaders = append(firstHeaders, "")
		w.AddColumn("Đơn giá (₫)", func() string { return FormatInt(ffm.Lines[*line].RetailPrice) })
		firstHeaders = append(firstHeaders, "")
		w.AddColumn("Thành tiền (trước giảm giá)", func() string { return FormatInt(ffm.Lines[*line].GetRetailAmount()) })
		firstHeaders = append(firstHeaders, "")
		w.AddColumn("Giảm giá (₫)", func() string { return FormatInt(ffm.Lines[*line].GetTotalDiscount()) })
		firstHeaders = append(firstHeaders, "")
		w.AddColumn("Thành tiền (sau giảm giá)", func() string { return FormatInt(ffm.Lines[*line].GetPaymentAmount()) })
	} else {
		firstHeaders = append(firstHeaders, "")
		w.AddColumn("Sản phẩm", func() string { return ffm.Lines.GetSummary() })
	}

	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Tổng số lượng sản phẩm", func() string { return FirstLine(line, FormatInt(ffm.Lines.GetTotalItems())) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Tổng tiền hàng (trước giảm giá)", func() string { return FirstLine(line, FormatInt(ffm.Lines.GetTotalRetailAmount())) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Giảm giá đơn hàng (₫)", func() string { return FirstLine(line, FormatInt(ffm.OrderDiscount)) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Tổng tiền hàng (sau giảm giá)", func() string { return FirstLine(line, FormatInt(ffm.Lines.GetTotalPaymentAmount())) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Tổng phí tính cho khách", func() string { return FirstLine(line, FormatInt(ffm.FeeLines.GetTotalFee())) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Tổng tiền thanh toán (bao gồm phí và giảm giá)", func() string {
		total := ffm.Lines.GetTotalPaymentAmount() - ffm.OrderDiscount + ffm.FeeLines.GetTotalFee()
		return FirstLine(line, FormatInt(total))
	})
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Trạng thái đơn hàng", func() string { return FirstLine(line, model.OrderStatusLabel(ffm.Status)) })

	firstHeaders = append(firstHeaders, "KHÁCH HÀNG")
	// TODO: add value customer code
	w.AddColumn("Mã KH", func() string { return FirstLine(line, "") })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Tên KH", func() string {
		if ffm.Customer == nil {
			return FirstLine(line, "")
		}
		return FirstLine(line, ffm.Customer.FullName)
	})
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Số điện thoại KH", func() string {
		if ffm.Customer == nil {
			return FirstLine(line, "")
		}
		return FirstLine(line, ffm.Customer.Phone)
	})

	firstHeaders = append(firstHeaders, "GIAO HÀNG")
	w.AddColumn("Mã giao hàng", func() string {
		if ffm.Fulfillment == nil {
			return FirstLine(line, "")
		}
		return FirstLine(line, ffm.Fulfillment.ShippingCode)
	})
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Nhà vận chuyển", func() string {
		if ffm.Fulfillment == nil {
			return FirstLine(line, "")
		}
		return FirstLine(line, ffm.Fulfillment.ShippingProvider.Label())
	})
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Trạng thái giao hàng", func() string {
		if ffm.Fulfillment == nil {
			return FirstLine(line, "")
		}
		return FirstLine(line, ffm.Fulfillment.ShippingState.Text())
	})

	firstHeaders = append(firstHeaders, "THANH TOÁN")
	w.AddColumn("Trạng thái thanh toán", func() string { return FirstLine(line, model.EtopPaymentStatusLabel(ffm.PaymentStatus)) })

	if err := csvWriter.Write(firstHeaders); err != nil {
		return err
	}
	return w.WriteHeader()
}
