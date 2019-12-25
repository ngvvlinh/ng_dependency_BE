package export

import (
	"context"
	"database/sql"
	"encoding/csv"
	"io"
	"strconv"
	"time"

	"etop.vn/api/top/int/shop"
	"etop.vn/backend/com/main/shipping/modely"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/cmapi"
	"etop.vn/backend/pkg/common/sql/sq/core"
)

const exportFulfillmentLines = false

func ExportFulfillments(
	ctx context.Context, id string, exportOpts ExportOption, output io.Writer,
	result chan<- *shop.ExportStatusItem,
	total int, rows *sql.Rows, opts core.Opts,
) (_err error) {

	var count, countError, line int
	maxErrors := total / 100
	if maxErrors > BaseRowsErrors {
		maxErrors = BaseRowsErrors
	}

	makeProgress := func() *shop.ExportStatusItem {
		return &shop.ExportStatusItem{
			Id:            id,
			ProgressMax:   total,
			ProgressValue: count,
			ProgressError: countError,
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
	var ffm modely.FulfillmentExtended
	if err := buildTableFulfillment(csvWriter, tableWriter, exportOpts, &count, &line, &ffm); err != nil {
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
			if exportFulfillmentLines {
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

func buildTableFulfillment(csvWriter *csv.Writer,
	w *TableWriter, exportOpts ExportOption,
	count *int, line *int, ffm *modely.FulfillmentExtended) error {
	var formatAsText = func(s string) string { return s }
	if exportOpts.ExcelMode {
		formatAsText = FormatAsTextForExcel
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
	w.AddColumn("Mã giao hàng", func() string { return formatAsText(FirstLine(line, ffm.ShippingCode)) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Nhà vận chuyển", func() string { return FirstLine(line, ffm.ShippingProvider.Label()) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Ngày tạo", func() string { return FirstLine(line, FormatDate(ffm.CreatedAt)) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Bảo hiểm", func() string { return FirstLine(line, LabelBool(ffm.IncludeInsurance)) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Lý do huỷ", func() string { return FirstLine(line, ffm.CancelReason) })

	firstHeaders = append(firstHeaders, "TRẠNG THÁI")
	w.AddColumn("Giao hàng", func() string { return FirstLine(line, ffm.ShippingState.Text()) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Ngày lấy", func() string { return FirstLine(line, FormatDate(ffm.ShippingPickingAt)) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Ngày giao", func() string { return FirstLine(line, FormatDate(ffm.ShippingDeliveredAt)) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Ngày trả hàng", func() string { return FirstLine(line, FormatDate(ffm.ShippingReturnedAt)) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Thu hộ", func() string {
		return FirstLine(line, LabelCODTransferAt(ffm.EtopPaymentStatus, ffm.CODEtopTransferedAt))
	})
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Ngày thanh toán tiền thu hộ", func() string {
		return FirstLine(line, FormatDate(ffm.CODEtopTransferedAt))
	})
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Phí giao hàng", func() string {
		return FirstLine(line, LabelShippingFeeTransferAt(ffm.EtopPaymentStatus, ffm.ShippingFeeShopTransferedAt))
	})
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Ngày thanh toán phí giao hàng", func() string {
		return FirstLine(line, FormatDate(ffm.ShippingFeeShopTransferedAt))
	})
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Mã phiên thanh toán", func() string {
		if ffm.MoneyTransactionShipping == nil {
			return ""
		}
		return formatAsText(FirstLine(line, ffm.MoneyTransactionShipping.Code))
	})

	firstHeaders = append(firstHeaders, "THÔNG TIN LẤY HÀNG")
	w.AddColumn("Tên người gửi", func() string { return FirstLine(line, ffm.AddressFrom.GetFullName()) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Số điện thoại", func() string { return formatAsText(FirstLine(line, ffm.AddressFrom.GetPhone())) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Tỉnh thành", func() string { return FirstLine(line, ffm.AddressFrom.GetProvince()) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Quận huyện", func() string { return FirstLine(line, ffm.AddressFrom.GetDistrict()) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Phường xã", func() string { return FirstLine(line, ffm.AddressFrom.GetWard()) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Địa chỉ", func() string { return FirstLine(line, ffm.AddressFrom.GetShortAddress()) })

	firstHeaders = append(firstHeaders, "THÔNG TIN NHẬN HÀNG")
	w.AddColumn("Tên người nhận", func() string { return FirstLine(line, ffm.AddressTo.GetFullName()) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Số điện thoại", func() string { return formatAsText(FirstLine(line, ffm.AddressTo.GetPhone())) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Tỉnh thành", func() string { return FirstLine(line, ffm.AddressTo.GetProvince()) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Quận huyện", func() string { return FirstLine(line, ffm.AddressTo.GetDistrict()) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Phường xã", func() string { return FirstLine(line, ffm.AddressTo.GetWard()) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Địa chỉ", func() string { return FirstLine(line, ffm.AddressTo.GetShortAddress()) })

	if exportFulfillmentLines {
		firstHeaders = append(firstHeaders, "THÔNG TIN SẢN PHẨM")
		w.AddColumn("Mã phiên bản sản phẩm", func() string { return formatAsText(ffm.Order.Lines[*line].Code) })
		firstHeaders = append(firstHeaders, "")
		w.AddColumn("Tên sản phẩm", func() string { return ffm.Order.Lines[*line].ProductName })
		firstHeaders = append(firstHeaders, "")
		w.AddColumn("Số lượng", func() string { return FormatInt(ffm.Order.Lines[*line].Quantity) })
		firstHeaders = append(firstHeaders, "")
		w.AddColumn("Đơn giá (₫)", func() string { return FormatInt(ffm.Order.Lines[*line].RetailPrice) })
		firstHeaders = append(firstHeaders, "")
		w.AddColumn("Thành tiền (trước giảm giá)", func() string { return FormatInt(ffm.Order.Lines[*line].GetRetailAmount()) })
		firstHeaders = append(firstHeaders, "")
		w.AddColumn("Giảm giá (₫)", func() string { return FormatInt(ffm.Order.Lines[*line].GetTotalDiscount()) })
		firstHeaders = append(firstHeaders, "")
		w.AddColumn("Thành tiền (sau giảm giá)", func() string { return FormatInt(ffm.Order.Lines[*line].GetPaymentAmount()) })
	} else {
		firstHeaders = append(firstHeaders, "THÔNG TIN SẢN PHẨM")
		w.AddColumn("Sản phẩm", func() string { return ffm.Order.Lines.GetSummary() })
	}

	firstHeaders = append(firstHeaders, "GIÁ TRỊ ĐƠN HÀNG")
	w.AddColumn("Tổng số lượng sản phẩm", func() string { return FirstLine(line, FormatInt(ffm.Order.Lines.GetTotalItems())) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Tổng khối lượng (g)", func() string { return FirstLine(line, FormatInt(ffm.Order.Lines.GetTotalWeight())) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Tổng tiền hàng (trước giảm giá)", func() string { return FirstLine(line, FormatInt(ffm.Order.Lines.GetTotalRetailAmount())) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Tổng tiền hàng (sau giảm giá)", func() string { return FirstLine(line, FormatInt(ffm.Order.Lines.GetTotalPaymentAmount())) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Giảm giá đơn hàng (₫)", func() string { return FirstLine(line, FormatInt(ffm.Order.OrderDiscount)) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Phí giao hàng tính cho khách", func() string { return FirstLine(line, FormatInt(ffm.Order.FeeLines.GetShippingFee())) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Thuế tính cho khách", func() string { return FirstLine(line, FormatInt(ffm.Order.FeeLines.GetTax())) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Phí khác tính cho khách", func() string { return FirstLine(line, FormatInt(ffm.Order.FeeLines.GetOtherFee())) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Tổng phí tính cho khách", func() string { return FirstLine(line, FormatInt(ffm.Order.FeeLines.GetTotalFee())) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Tổng tiền thanh toán (bao gồm phí và giảm giá)", func() string {
		total := ffm.Order.Lines.GetTotalPaymentAmount() - ffm.Order.OrderDiscount + ffm.Order.FeeLines.GetTotalFee()
		return FirstLine(line, FormatInt(total))
	})
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Giá trị bảo hiểm", func() string {
		if !ffm.IncludeInsurance {
			return ""
		}
		return FirstLine(line, FormatInt(ffm.BasketValue))
	})

	firstHeaders = append(firstHeaders, "TIỀN THU HỘ VÀ CHI PHÍ")
	w.AddColumn("Tiền thu hộ", func() string { return FirstLine(line, FormatInt(ffm.TotalCODAmount)) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Phí giao hàng", func() string { return FirstLine(line, FormatIntBeforeVAT(ffm.ShippingFeeMain)) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Phí trả hàng", func() string { return FirstLine(line, FormatIntBeforeVAT(ffm.ShippingFeeReturn)) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Phí vượt cân", func() string { return FirstLine(line, FormatIntBeforeVAT(ffm.ShippingFeeAdjustment)) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Phí bảo hiểm", func() string { return FirstLine(line, FormatIntBeforeVAT(ffm.ShippingFeeInsurance)) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Phí thu hộ", func() string { return FirstLine(line, FormatIntBeforeVAT(ffm.ShippingFeeCODS)) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Phí đổi thông tin", func() string { return FirstLine(line, FormatIntBeforeVAT(ffm.ShippingFeeInfoChange)) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Phí khác", func() string { return FirstLine(line, FormatIntBeforeVAT(ffm.ShippingFeeOther)) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Tổng phí (trước VAT)", func() string { return FirstLine(line, FormatIntBeforeVAT(ffm.ShippingFeeShop)) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("VAT (10%)", func() string { return FirstLine(line, FormatIntVAT(ffm.ShippingFeeShop)) })
	firstHeaders = append(firstHeaders, "")
	w.AddColumn("Tổng phí (sau VAT)", func() string { return FirstLine(line, FormatInt(ffm.ShippingFeeShop)) })

	if err := csvWriter.Write(firstHeaders); err != nil {
		return err
	}
	return w.WriteHeader()
}
