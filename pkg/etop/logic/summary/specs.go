package summary

import (
	"fmt"
	"strings"
	"time"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/sq"
)

func startOfDay(dayFromToday int) time.Time {
	return time.Now().Truncate(24 * time.Hour).
		Add(time.Duration(dayFromToday) * 24 * time.Hour).
		// Truncate() truncates time by UTC, so we have to convert to ICT
		Add(-7 * time.Hour)
}

func endOfDay(dayFromToday int) time.Time {
	return time.Now().Truncate(24 * time.Hour).
		Add(time.Duration(dayFromToday+1) * 24 * time.Hour).
		// convert to ICT
		Add(-7 * time.Hour)
}

func getWeek(t time.Time) int {
	_, week := t.ISOWeek()

	// our week starts on Monday
	if t.Weekday() == time.Sunday {
		week--
	}
	return week
}

func dayOfWeek(weekday time.Weekday) time.Time {
	today := time.Now()

	// our week starts on Monday
	if today.Weekday() == time.Sunday {
		today = today.Add(-7 * 24 * time.Hour)
	}

	d := today.Add(time.Duration(weekday-today.Weekday()) * 24 * time.Hour)
	if weekday == time.Sunday {
		d = d.Add(7 * 24 * time.Hour)
	}
	return d
}

func pred_giao_sau(labelFormat string, dayFromToday int) Predicate {
	dayStr := "today"
	if dayFromToday > 0 {
		dayStr = fmt.Sprintf("(today+%v)", dayFromToday)
	}
	start := startOfDay(dayFromToday)
	end := endOfDay(dayFromToday)
	_, month, day := start.Date()

	return Predicate{
		Label: fmt.Sprintf(labelFormat, fmt.Sprintf("%2d/%d", day, month)),
		Spec:  fmt.Sprintf("%v(expected_delivery_at)", dayStr),
		Expr:  sq.NewExpr("expected_delivery_at BETWEEN ? AND ?", start, end),
	}
}

func pred_tạo_thứ(labelFormat string, weekday time.Weekday) Predicate {
	start := dayOfWeek(weekday)
	end := start.Add(24 * time.Hour)
	_, month, day := start.Date()
	return Predicate{
		Label: fmt.Sprintf(labelFormat, fmt.Sprintf("%2d/%d", day, month)),
		Spec:  fmt.Sprintf("%v(created_at)", strings.ToLower(weekday.String())),
		Expr:  sq.NewExpr("created_at BETWEEN ? AND ?", start, end),
	}
}

func pred_tuần(column string, weekFromToday int) Predicate {
	start := dayOfWeek(time.Monday).Add(time.Duration(weekFromToday) * 7 * 24 * time.Hour)
	end := start.Add(7 * 24 * time.Hour)
	week := getWeek(start)

	p := Predicate{
		Expr: sq.NewExpr(column+" BETWEEN ? AND ?", start, end),
	}
	if weekFromToday == 0 {
		p.Label = fmt.Sprintf("Tuần này [Tuần %v]", week)
		p.Spec = fmt.Sprintf("thisweek(" + column + ")")
	} else {
		p.Label = fmt.Sprintf("Tuần %v", week)
		p.Spec = fmt.Sprintf("(thisweek%+d)("+column+")", weekFromToday)
	}
	return p
}

func pred_tạo_tuần(weekFromToday int) Predicate {
	return pred_tuần("created_at", weekFromToday)
}

var (
	row_tổng_đơn          = NewSubject("Tổng đơn", "", "count", "COUNT(*)", nil)
	row_tổng_giá_trị      = NewSubject("Tổng giá trị giao hàng", "₫", "sum(total_amount)", "SUM(total_amount)", nil)
	row_giá_trị_cod       = NewSubject("Tiền COD", "₫", "sum(total_cod_amount)", "SUM(total_cod_amount)", nil)
	row_tổng_cước_phí     = NewSubject("Tổng chi phí", "₫", "sum(shipping_fee_shop)", "SUM(shipping_fee_shop)", nil)
	row_phí_giao_hàng     = NewSubject("Phí giao hàng [gồm các đơn trả hàng]", "₫", "sum(shipping_fee_main)", "SUM(shipping_fee_main)", nil)
	row_phí_trả_hàng      = NewSubject("Phí trả hàng [Tính riêng đơn thất bại]", "₫", "sum(shipping_fee_return)", "SUM(shipping_fee_return)", nil)
	row_phí_bảo_hiểm      = NewSubject("Phí bảo hiểm [Theo nhà vận chuyển]", "₫", "sum(shipping_fee_insurance)", "SUM(shipping_fee_insurance)", nil)
	row_phí_vượt_cân      = NewSubject("Phí vượt cân [Theo khối lượng vượt cân]", "₫", "sum(shipping_fee_adjustment)", "SUM(shipping_fee_adjustment)", nil)
	row_phí_thu_hộ        = NewSubject("Phí thu hộ (COD) [Viettel Post tính phí với hàng trên 3.000.000 ₫]", "₫", "sum(shipping_fee_cods)", "SUM(shipping_fee_cods)", nil)
	row_phí_đổi_thông_tin = NewSubject("Phí đổi thông tin [GHTK có dịch vụ này]", "₫", "sum(shipping_fee_info_change)", "SUM(shipping_fee_info_change)", nil)
	row_phí_khác          = NewSubject("Các phụ phí khác [Có liệt kê trong đối soát, chi tiết đơn giao hàng]", "₫", "sum(shipping_fee_other)", "SUM(shipping_fee_other)", nil)
)

func buildTables(moneyTransactionIDs []int64) []*Table {
	inArgs := make([]interface{}, len(moneyTransactionIDs))
	for i, id := range moneyTransactionIDs {
		inArgs[i] = id
	}

	// Bao gồm các dơn mới, chưa giao, trả hàng, ... nhưng chưa đối soát.
	// Không bao gồm đơn huỷ và đơn tạo không thành công.
	pred_chưa_đối_soát := Predicate{
		Spec: "status=S",
		Expr: sq.NewExpr("status = 2"),
	}

	pred_sắp_chuyển := Predicate{
		Spec: "next_transaction",
		Expr: sq.In("money_transaction_id", inArgs...),
	}
	// không bao gồm các đơn đã huỷ (-1) hoặc tạo không thành công (0)
	pred_đã_tạo_không_huỷ := Predicate{
		Spec: "status!=Z|N",
		Expr: sq.NewExpr("status != -1 AND status != 0"),
	}
	pred_mới := Predicate{
		Spec: "shipping_status=Z",
		Expr: sq.NewExpr("shipping_status = 0"),
	}
	pred_không_mới := Predicate{
		Spec: "shipping_status=!Z",
		Expr: sq.NewExpr("shipping_status != 0"),
	}
	pred_đang_lấy := Predicate{
		Spec: "shipping_state=picking",
		Expr: sq.NewExpr("shipping_state = 'picking'"),
	}
	pred_lấy_hàng_trễ := Predicate{
		Spec: "late(expected_pick_at)",
		Expr: sq.NewExpr("expected_pick_at < NOW()"),
	}
	pred_đang_giao := Predicate{
		Spec: "shipping_state=holding|delivering|undeliverable",
		Expr: sq.NewExpr("shipping_state IN ('holding','delivering','undeliverable')"),
	}
	pred_đã_giao := Predicate{
		Spec: "shipping_status=P",
		Expr: sq.NewExpr("shipping_status = 1"),
	}
	pred_đã_giao_hoặc_đang_giao := Predicate{
		Spec: "shipping_status=P|S",
		Expr: sq.NewExpr("shipping_status IN (1,2)"),
	}
	pred_đã_giao_hoặc_đã_trả_hàng := Predicate{
		Spec: "shipping_status=P|NS",
		Expr: sq.NewExpr("shipping_status IN (1,-2)"),
	}
	pred_trả_hàng := Predicate{
		Spec: "shipping_state=returning|returned",
		Expr: sq.NewExpr("shipping_state IN ('returning','returned')"),
	}
	pred_không_trả_hàng := Predicate{
		Spec: "shipping_status!=returning|returned",
		Expr: sq.NewExpr("shipping_state NOT IN ('returning','returned')"),
	}
	pred_giao_trễ := Predicate{
		Spec: "late(expected_delivery_at)",
		Expr: sq.NewExpr("expected_delivery_at < NOW()"),
	}
	pred_giao_7_ngày_trở_lên := Predicate{
		Spec: "after_6_days(expected_delivery_at)",
		Expr: sq.NewExpr("expected_delivery_at >= ?", endOfDay(6)),
	}
	pred_cod := Predicate{
		Spec: "cod",
		Expr: sq.NewExpr("total_cod_amount != 0"),
	}
	pred_not_cod := Predicate{
		Spec: "!cod",
		Expr: sq.NewExpr("total_cod_amount = 0"),
	}

	// chưa đối soát
	var table00, table01, table02, table03 *Table
	{
		col_sắp_chuyển := Compose("Sắp chuyển [Trong phiên kế tiếp]", pred_chưa_đối_soát, pred_sắp_chuyển)
		col_đã_giao := Compose("Đã giao", pred_chưa_đối_soát, pred_đã_giao)
		col_đang_giao := Compose("Đang giao", pred_chưa_đối_soát, pred_đang_giao)

		col_mới := Compose("Mới", pred_đã_tạo_không_huỷ, pred_mới)
		col_đang_lấy := Compose("Đang lấy", pred_chưa_đối_soát, pred_đang_lấy)
		col_lấy_hàng_trễ := Compose("Lấy trễ", pred_chưa_đối_soát, pred_đang_lấy, pred_lấy_hàng_trễ)
		col_trả_hàng := Compose("Trả hàng", pred_chưa_đối_soát, pred_trả_hàng)
		col_tất_cả := Compose("Tất cả", pred_chưa_đối_soát)

		col_giao_trễ := Compose("Giao trễ [So với dự kiến]", pred_giao_trễ, pred_chưa_đối_soát, pred_đang_giao)
		col_giao_hôm_nay := Compose("", pred_giao_sau("Hôm nay [(%v)]", 0), pred_chưa_đối_soát, pred_đang_giao)
		col_giao_ngày_mai := Compose("", pred_giao_sau("Ngày mai [(%v)]", 1), pred_chưa_đối_soát, pred_đang_giao)
		col_giao_ngày_mốt := Compose("", pred_giao_sau("Ngày mốt [(%v)]", 2), pred_chưa_đối_soát, pred_đang_giao)
		col_giao_3_ngày_sau := Compose("", pred_giao_sau("3 ngày sau [(%v)]", 3), pred_chưa_đối_soát, pred_đang_giao)
		col_giao_4_ngày_sau := Compose("", pred_giao_sau("4 ngày sau [(%v)]", 4), pred_chưa_đối_soát, pred_đang_giao)
		col_giao_5_ngày_sau := Compose("", pred_giao_sau("5 ngày sau [(%v)]", 5), pred_chưa_đối_soát, pred_đang_giao)
		col_giao_6_ngày_sau := Compose("", pred_giao_sau("6 ngày sau [(%v)]", 6), pred_chưa_đối_soát, pred_đang_giao)
		col_giao_7_ngày_trở_lên := Compose("Khác [Các ngày xa hơn]", pred_giao_7_ngày_trở_lên, pred_chưa_đối_soát, pred_đang_giao)

		{
			rows := []Subject{
				row_tổng_đơn,
				row_tổng_đơn.Combine("Đơn COD [Đã giao, có thu hộ]", pred_cod, pred_đã_giao_hoặc_đang_giao).WithIdent(1),
				row_tổng_đơn.Combine("Đơn không COD [Đã giao, không thu hộ]", pred_not_cod, pred_đã_giao_hoặc_đang_giao).WithIdent(1),
				row_tổng_đơn.Combine("Đơn trả hàng [Giao thất bại]", pred_trả_hàng).WithIdent(1),
				row_giá_trị_cod.Combine("Tiền COD [Đã giao thành công]", pred_đã_giao),
				row_phí_giao_hàng,
				row_phí_trả_hàng.Combine("Phí trả hàng [Tính riêng đơn thất bại]", pred_trả_hàng),
				row_phí_bảo_hiểm,
				row_phí_vượt_cân,
				row_phí_thu_hộ,
				row_phí_đổi_thông_tin,
				row_phí_khác,
			}
			cols := []Predicator{
				col_sắp_chuyển,
				Compose("Đã giao hoặc trả hàng [Gồm đơn sắp chuyển]", pred_chưa_đối_soát, pred_đã_giao_hoặc_đã_trả_hàng),
				Compose("Đang giao [Chưa hoàn tất]", pred_chưa_đối_soát, pred_đang_giao),
			}
			table00 = buildTable(rows, cols, "Thống kê đơn hàng chưa đối soát", "overview_before_cross_check", "before_cross_check")
		}
		{
			rows := []Subject{
				row_tổng_đơn,
				row_tổng_đơn.Combine("Đơn COD", pred_cod),
				row_tổng_đơn.Combine("Đơn không COD", pred_not_cod),
				row_tổng_giá_trị,
				row_giá_trị_cod,
			}
			cols := []Predicator{
				col_mới,
				col_đang_lấy,
				col_lấy_hàng_trễ,
				col_đang_giao,
				col_giao_hôm_nay.Clone("Giao hôm nay"),
				col_đã_giao,
				col_trả_hàng,
				col_tất_cả,
			}
			table01 = buildTable(rows, cols, "Trạng thái giao hàng", "shipping_status_before_cross_check", "before_cross_check", "shipping_status")
		}
		{
			// Không bao gồm các đơn mới vì shop có thể huỷ.
			// Điều kiện chưa đối soát (không huỷ, không lỗi) đã bao gồm trong danh sách cột.
			rows := []Subject{
				row_tổng_đơn.Combine("Tổng đơn", pred_không_mới),
				row_tổng_đơn.Combine("Đơn COD", pred_không_mới, pred_cod),
				row_tổng_đơn.Combine("Đơn không COD", pred_không_mới, pred_not_cod),
				row_tổng_giá_trị.Combine("Giá trị", pred_không_mới),
				row_giá_trị_cod.Combine("Tiền COD", pred_không_mới, pred_cod),
			}
			cols := []Predicator{
				col_giao_trễ,
				col_giao_hôm_nay,
				col_giao_ngày_mai,
				col_giao_ngày_mốt,
				col_giao_3_ngày_sau,
				col_giao_4_ngày_sau,
				col_giao_5_ngày_sau,
				col_giao_6_ngày_sau,
				col_giao_7_ngày_trở_lên,
				col_đang_giao.Clone("Tất cả"),
			}
			table02 = buildTable(rows, cols, "Ngày giao hàng dự kiến [Không thống kê các đơn hàng MỚI vì có thể hủy]", "expected_delivery_at_before_cross_check", "before_cross_check", "expected_delivery_at")
		}
	}
	// cả đối soát và chưa đối soát
	{
		col_thứ_hai := Compose("", pred_tạo_thứ("Thứ 2 [Ngày %v]", time.Monday), pred_đã_tạo_không_huỷ)
		col_thứ_ba_ := Compose("", pred_tạo_thứ("Thứ 3 [Ngày %v]", time.Tuesday), pred_đã_tạo_không_huỷ)
		col_thứ_tư_ := Compose("", pred_tạo_thứ("Thứ 4 [Ngày %v]", time.Wednesday), pred_đã_tạo_không_huỷ)
		col_thứ_năm := Compose("", pred_tạo_thứ("Thứ 5 [Ngày %v]", time.Thursday), pred_đã_tạo_không_huỷ)
		col_thứ_sáu := Compose("", pred_tạo_thứ("Thứ 6 [Ngày %v]", time.Friday), pred_đã_tạo_không_huỷ)
		col_thứ_bảy := Compose("", pred_tạo_thứ("Thứ 7 [Ngày %v]", time.Saturday), pred_đã_tạo_không_huỷ)
		col_chủ_nhật := Compose("", pred_tạo_thứ("Chủ Nhật [Ngày %v]", time.Sunday), pred_đã_tạo_không_huỷ)
		col_tuần_này := Compose("", pred_tạo_tuần(0), pred_đã_tạo_không_huỷ)
		col_tuần_trước := Compose("", pred_tạo_tuần(-1), pred_đã_tạo_không_huỷ)
		col_tuần_trước_nữa := Compose("", pred_tạo_tuần(-2), pred_đã_tạo_không_huỷ)

		{
			// col đã bao gồm điều kiện không huỷ, nên không lặp lại ở đây
			rows := []Subject{
				row_tổng_đơn,
				row_tổng_đơn.Combine("Đơn giao hàng", pred_không_trả_hàng),
				row_tổng_đơn.Combine("Đơn COD", pred_cod),
				row_tổng_đơn.Combine("Đơn trả hàng", pred_trả_hàng),
				row_tổng_giá_trị.Combine("Giá trị giao hàng", pred_không_trả_hàng),
				row_tổng_giá_trị.Combine("Giá trị trả hàng", pred_trả_hàng),
				row_giá_trị_cod.Combine("Tiền COD [Không bao gồm đơn trả hàng]", pred_không_trả_hàng),
			}
			cols := []Predicator{
				col_thứ_hai,
				col_thứ_ba_,
				col_thứ_tư_,
				col_thứ_năm,
				col_thứ_sáu,
				col_thứ_bảy,
				col_chủ_nhật,
				col_tuần_này,
				col_tuần_trước,
				col_tuần_trước_nữa,
			}
			table03 = buildTable(rows, cols, "Số lượng đơn giao hàng đã được tạo [Tính theo đơn giao hàng được tạo ra và không HUỶ mỗi ngày.\nBao gồm cả đơn đã đối soát và đơn đang giao]", "created_at_before_cross_check", "before_cross_check", "created_at")
		}
	}
	return []*Table{table00, table01, table02, table03}
}

func formatDate(t time.Time) string {
	return t.Format(cm.DateLayout)
}

func buildTables2(dateFrom, dateTo time.Time) []*Table {
	pred_time_between := Predicate{
		Spec: fmt.Sprintf("created_at[%v|%v)", formatDate(dateFrom), formatDate(dateTo)),
		Expr: sq.NewExpr("created_at BETWEEN ? AND ?", dateFrom, dateTo),
	}

	// bao gồm các đơn đã đối soát
	pred_đã_đối_soát := Predicate{
		Spec: "money_transaction_id",
		Expr: sq.NewExpr("money_transaction_id IS NOT NULL"),
	}
	pred_ghn := Predicate{
		Spec: "shipping_provider=ghn",
		Expr: sq.NewExpr("shipping_provider = 'ghn'"),
	}
	pred_ghtk := Predicate{
		Spec: "shipping_provider=ghtk",
		Expr: sq.NewExpr("shipping_provider = 'ghtk'"),
	}
	pred_vtp := Predicate{
		Spec: "shipping_provider=vtpost",
		Expr: sq.NewExpr("shipping_provider = 'vtpost'"),
	}
	pred_đã_giao := Predicate{
		Spec: "shipping_status=P",
		Expr: sq.NewExpr("status = 1"),
	}
	pred_đã_trả_hàng := Predicate{
		Spec: "shipping_status=NS",
		Expr: sq.NewExpr("status = -2"),
	}

	var table00, table01 *Table
	{
		rows := []Subject{
			row_tổng_đơn.Combine("Tổng quan giao hàng"),
			row_tổng_đơn.Combine("Thành công", pred_đã_giao),
			row_tổng_đơn.Combine("Trả hàng", pred_đã_trả_hàng),
			row_tổng_cước_phí.Combine("Tổng chi phí [Bao gồm tất cả phí giao hàng, trả hàng, vượt cân, ...]"),
		}
		cols := []Predicator{
			Compose("GHN", pred_đã_đối_soát, pred_time_between, pred_ghn),
			Compose("GHTK", pred_đã_đối_soát, pred_time_between, pred_ghtk),
			Compose("VTP", pred_đã_đối_soát, pred_time_between, pred_vtp),
			Compose("Tổng", pred_đã_đối_soát, pred_time_between),
		}
		table00 = buildTable(rows, cols, "Tổng quan giao hàng", "shipping_provider_after_cross_check", "after_cross_check", "shipping_provider")
	}
	{
		pred_đã_chuyển_tiền := Predicate{
			Spec: "cod_etop_transfered_at",
			Expr: sq.NewExpr("cod_etop_transfered_at IS NOT NULL"),
		}
		row_doanh_thu_còn_lại := NewSubject("Doanh thu còn lại [Đã trừ các chi phí trên]", "₫", "sum(total_amount-shipping_fee_shop)", "SUM(total_amount - shipping_fee_shop)", nil)
		row_cod_còn_lại := NewSubject("COD còn lại [Đã trừ các chi phí trên]", "₫", "sum(total_cod_amount-shipping_fee_shop)", "SUM(total_cod_amount - shipping_fee_shop)", nil)
		row_cod_đã_nhận := NewSubject("COD đã nhận [TOPSHIP đã chuyển & Cấn trừ các chi phí trên]", "₫", "sum(total_cod_amount-shipping_fee_shop)", "SUM(total_cod_amount - shipping_fee_shop)", pred_đã_chuyển_tiền)
		rows := []Subject{
			row_tổng_đơn.Combine("Giao thành công", pred_đã_giao),
			row_tổng_đơn.Combine("Tổng đơn giao hàng").WithIdent(1),
			row_tổng_đơn.Combine("Đơn trả hàng", pred_đã_trả_hàng).WithIdent(1),
			row_tổng_giá_trị.Combine("Doanh thu"),
			row_giá_trị_cod.Combine("Tiền thu hộ (COD)", pred_đã_giao),
			row_tổng_cước_phí.Combine("Tổng chi phí"),
			row_phí_giao_hàng.WithIdent(1),
			row_phí_trả_hàng.Combine("Phí trả hàng [Tính riêng đơn thất bại]", pred_đã_trả_hàng).WithIdent(1),
			row_phí_bảo_hiểm.WithIdent(1),
			row_phí_vượt_cân.WithIdent(1),
			row_phí_thu_hộ.WithIdent(1),
			row_phí_đổi_thông_tin.WithIdent(1),
			row_phí_khác.WithIdent(1),
			row_doanh_thu_còn_lại.WithIdent(-1),
			row_cod_còn_lại.WithIdent(-1),
			row_cod_đã_nhận.WithIdent(-1),
		}
		cols := []Predicator{
			pred_đã_đối_soát.Clone("Toàn thời gian"),
			Compose("", pred_tạo_tuần(0), pred_đã_đối_soát),
			Compose("", pred_tạo_tuần(-1), pred_đã_đối_soát),
			Compose("", pred_tạo_tuần(-2), pred_đã_đối_soát),
			Compose("", pred_tạo_tuần(-3), pred_đã_đối_soát),
			Compose("", pred_tạo_tuần(-4), pred_đã_đối_soát),
			Compose("", pred_tạo_tuần(-5), pred_đã_đối_soát),
			Compose("", pred_tạo_tuần(-6), pred_đã_đối_soát),
			Compose("", pred_tạo_tuần(-7), pred_đã_đối_soát),
			Compose("", pred_tạo_tuần(-8), pred_đã_đối_soát),
			Compose("", pred_tạo_tuần(-9), pred_đã_đối_soát),
		}
		table01 = buildTable(rows, cols, "Số lượng đơn giao hàng đã được đối soát trong 10 tuần gần đây [Tính theo mốc thời gian đơn giao hàng tạo ra. Ghi nhận khi đơn giao hàng đã đối soát hoàn tất. Nghĩa là đã thanh toán phí vận chuyển và nhận tiền COD.\nĐơn COD hoặc giá trị đơn giao COD được tính trên các đơn giao hàng thành công. Các đơn trả hàng xem như bằng 0.]", "created_at_after_cross_check", "after_cross_check", "created_at")
	}
	return []*Table{table00, table01}
}

func buildTable(rows []Subject, cols []Predicator, label string, tags ...string) *Table {
	table := NewTable(len(rows), len(cols), label, tags...)
	table.Rows = rows
	table.Cols = cols
	for r, row := range rows {
		for c, col := range cols {
			table.SetCell(r, c, row.Combine("", col))
		}
	}
	return table
}
