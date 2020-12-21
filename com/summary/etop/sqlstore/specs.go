package sqlstore

import (
	"context"
	"time"

	"o.o/api/summary"
	"o.o/api/top/types/etc/status5"
	"o.o/backend/com/summary/util"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sq/core"
	smry "o.o/backend/pkg/etop/logic/summary"
	etopmodel "o.o/backend/pkg/etop/model"
	"o.o/capi/dot"
)

type SummaryStore struct {
	db    *cmsql.Database
	query cmsql.QueryFactory
	preds []interface{}
}

func NewSummaryStore(db *cmsql.Database) SummaryStore {
	return SummaryStore{db: db}
}

func (s SummaryStore) SummarizeTopship(ctx context.Context, req *summary.SummaryTopShipRequest) ([]*summary.SummaryTable, error) {
	dateFrom, dateTo, shopID := req.DateFrom, req.DateTo, req.ShopID
	_, _ = dateFrom, dateTo

	tablesFfm, tablesShipnowFfm := buildTablesTopShipFulfillment(dateFrom, dateTo)
	if err := s.execQuery(ctx, tablesFfm, shopID, "fulfillment"); err != nil {
		return nil, err
	}
	if err := s.execQuery(ctx, tablesShipnowFfm, shopID, "shipnow_fulfillment"); err != nil {
		return nil, err
	}
	resTablesFfm := util.BuildResponse(tablesFfm)
	resTablesShipnowFfm := util.BuildResponse(tablesShipnowFfm)
	for key, value := range resTablesFfm {
		for _, _value := range resTablesShipnowFfm {
			if Contain(value.Tags, "fulfillment03") && Contain(_value.Tags, "fulfillment03") {
				resTablesFfm[key] = addCarrierShipNow(value, _value)
				break
			}
			if value.Tags[0] == _value.Tags[0] {
				resTablesFfm[key] = mergeData(value, _value)
				break
			}
		}

	}
	countShipnowFfm := int64(0)
	for _, value := range resTablesShipnowFfm {
		if Contain(value.Tags, "shipnow_fullfillment") && Contain(value.Tags, "fulfillment03") {
			countShipnowFfm = value.Data[0].Value
			break
		}
	}
	for key, value := range resTablesFfm {
		if value.Tags[0] == "fulfillment04" {
			resTablesFfm[key].Data[0].Value = resTablesFfm[key].Data[0].Value + countShipnowFfm
			break
		}
	}
	return resTablesFfm, nil
}

func addCarrierShipNow(table1 *summary.SummaryTable, table2 *summary.SummaryTable) *summary.SummaryTable {
	// them 1 dong nha van chuyen shipnow
	table1.Rows = append(table1.Rows, summary.SummaryColRow{
		Label:  "Shipnow",
		Spec:   "shipping_provider=shipnow",
		Unit:   "",
		Indent: 0,
	})
	// them 1 dong so don hang nha van chuyen shipnow
	table1.Data = append(table1.Data, summary.SummaryItem{
		ImageUrls: nil,
		Label:     "Shipnow",
		Spec:      "count:shipping_provider=shipnow",
		Value:     table2.Data[0].Value,
		Unit:      "",
	})
	// tinh lai tong so don hang
	table1.Data[0].Value = table1.Data[0].Value + table2.Data[0].Value
	return table1
}

func mergeData(table1 *summary.SummaryTable, table2 *summary.SummaryTable) *summary.SummaryTable {
	for key, _ := range table1.Data {
		table1.Data[key].Value = table1.Data[key].Value + table2.Data[key].Value
	}
	return table1
}

func Contain(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

func buildData(data []smry.Cell) []summary.SummaryItem {
	res := make([]summary.SummaryItem, len(data))
	for i, item := range data {
		res[i] = summary.SummaryItem{
			Spec:  item.Subject.GetSpec(),
			Value: item.Value,
			Unit:  item.Subject.Unit,
		}
	}
	return res
}

func (s SummaryStore) execQuery(ctx context.Context, tables []*smry.Table, shopID dot.ID, tableName string) error {
	builder := smry.NewSummaryQueryBuilder(tableName)
	count := 1
	for _, table := range tables {
		for i := range table.Data {
			count++
			// must always use [i] because we want to take the address
			builder.AddCell(&table.Data[i].Subject, (*core.Int64)(&table.Data[i].Value))
		}
	}
	q := s.db.SQL(builder).WithContext(ctx).
		Where("shop_id = ?", shopID).Clone()
	err := q.Scan(builder.ScanArgs...)
	q.Clone()

	return err
}

func buildTablesTopShipFulfillment(dateFrom time.Time, dateTo time.Time) (ffm []*smry.Table, shipnowFfm []*smry.Table) {
	pred_tạo_hôm_nay := smry.Pred_sau_ngày("created_at", "Hôm nay [%v]", 0)

	pred_đã_hủy := smry.Predicate{
		Spec: "status=N",
		Expr: sq.NewExpr("status=?", status5.N),
	}

	pred_chưa_hủy := smry.Predicate{
		Spec: "status!=-1",
		Expr: sq.NewExpr("status!=?", status5.N),
	}
	pre_tính_cod := smry.Predicate{
		Spec: "shipping_state!=returned,returning|status!=-1",
		Expr: sq.NewExpr("shipping_state not in ('returned','returning') and status!=?", status5.N),
	}

	pred_đã_lấy := smry.Predicate{
		Spec: "shipping_state=holding,delivering,delivered,underliverable,returning,returned",
		Expr: sq.NewExpr("shipping_state IN ('holding','delivering','delivered','undeliverable','returning','returned')"),
	}
	pred_chưa_lấy := smry.Predicate{
		Spec: "shipping_state=defaul,created,picking",
		Expr: sq.NewExpr("shipping_state IN ('default','created','picking')"),
	}

	row_tổng_đơn := smry.NewSubject("Tổng đơn", "", "count", "COUNT(*)", nil)

	//table tổng kết ngày
	rows00 := []smry.Subject{
		row_tổng_đơn.Combine(""),
		row_tổng_đơn.Combine("", pred_đã_hủy),
		row_tổng_đơn.Combine("", pred_đã_lấy, pred_chưa_hủy),
		row_tổng_đơn.Combine("", pred_chưa_lấy, pred_chưa_hủy),
	}
	cols00 := []smry.Predicator{
		pred_tạo_hôm_nay,
	}

	table00 := smry.BuildTable(rows00, cols00, "Hôm nay", "fulfillment", "shipnow_fulfillment", "today")

	pred_khoảng_thời_gian := smry.Predicate{
		Spec: "datefrom-dateto",
		Expr: sq.NewExpr("created_at >= ? and ? > created_at ", dateFrom, dateTo),
	}

	pred_đang_xử_lý := smry.Predicate{
		Spec: "shipping_state=holding,delivering",
		Expr: sq.NewExpr("shipping_state in ('holding', 'delivering')"),
	}
	pred_giao_thành_công := smry.Predicate{
		Spec: "shipping_state=delivered",
		Expr: sq.NewExpr("shipping_state in ('delivered')"),
	}
	pred_trả_hàng := smry.Predicate{
		Spec: "shipping_state=returned,returning",
		Expr: sq.NewExpr("shipping_state in ('returned', 'returning')"),
	}

	// Tổng kết theo thời gian dateFrom - dateTo
	rows01 := []smry.Subject{
		row_tổng_đơn.Combine(""),
		row_tổng_đơn.Combine("", pred_chưa_lấy, pred_chưa_hủy),
		row_tổng_đơn.Combine("", pred_đang_xử_lý, pred_chưa_hủy),
		row_tổng_đơn.Combine("", pred_giao_thành_công, pred_chưa_hủy),
		row_tổng_đơn.Combine("", pred_trả_hàng, pred_chưa_hủy),
		row_tổng_đơn.Combine("", pred_đã_hủy),
	}
	cols01 := []smry.Predicator{
		pred_khoảng_thời_gian,
	}
	table01 := smry.BuildTable(rows01, cols01, "Tổng đơn trong tháng", "fulfillment01", "shipnow_fulfillment", "datefrom-dateto")

	// Tổng đơn từng ngày
	rows02 := []smry.Subject{
		row_tổng_đơn,
	}
	rowsPerDay := buildRowPerDate(dateFrom, dateTo)
	table02 := smry.BuildTable(rows02, rowsPerDay, "Số lượng đơn theo ngày", "fulfillment02", "shipnow_fulfillment", "datefrom-dateto")

	pred_giao_hàng_nhanh := smry.Predicate{
		Spec: "shipping_provider=ghn",
		Expr: sq.NewExpr("shipping_provider = 'ghn'"),
	}

	pred_giao_hàng_tiết_kiệm := smry.Predicate{
		Spec: "shipping_provider=ghtk",
		Expr: sq.NewExpr("shipping_provider = 'ghtk'"),
	}

	pred_viettel_post := smry.Predicate{
		Spec: "shipping_provider=vtpost",
		Expr: sq.NewExpr("shipping_provider = 'vtpost'"),
	}

	rows03 := []smry.Subject{
		row_tổng_đơn.Combine("", pred_chưa_hủy),
		row_tổng_đơn.Combine("", pred_giao_hàng_nhanh, pred_chưa_hủy),
		row_tổng_đơn.Combine("", pred_giao_hàng_tiết_kiệm, pred_chưa_hủy),
		row_tổng_đơn.Combine("", pred_viettel_post, pred_chưa_hủy),
	}
	cols03 := []smry.Predicator{
		pred_khoảng_thời_gian,
	}
	table03 := smry.BuildTable(rows03, cols03, "Số lượng đơn theo nhà vận chuyển", "fulfillment03", "datefrom-dateto&&provider")

	rows10 := []smry.Subject{
		row_tổng_đơn.Combine(""),
	}
	cols10 := []smry.Predicator{
		pred_khoảng_thời_gian,
	}
	table10 := smry.BuildTable(rows10, cols10, "Số lượng đơn theo nhà vận chuyển shipnow", "fulfillment03", "shipnow_fullfillment", "datefrom-dateto&&provider")

	pred_nội_tỉnh := smry.Predicate{
		Spec: "delivery_route=noi_tinh",
		Expr: sq.NewExpr("delivery_route=?", etopmodel.RouteSameProvince),
	}

	pred_nội_miền := smry.Predicate{
		Spec: "delivery_route=noi_mien",
		Expr: sq.NewExpr("delivery_route=?", etopmodel.RouteSameRegion),
	}

	pred_ngoại_miền := smry.Predicate{
		Spec: "delivery_route=toan_quoc",
		Expr: sq.NewExpr("delivery_route=?", etopmodel.RouteNationWide),
	}

	rows04 := []smry.Subject{
		row_tổng_đơn.Combine("", pred_chưa_hủy),
		row_tổng_đơn.Combine("", pred_nội_tỉnh, pred_chưa_hủy),
		row_tổng_đơn.Combine("", pred_nội_miền, pred_chưa_hủy),
		row_tổng_đơn.Combine("", pred_ngoại_miền, pred_chưa_hủy),
	}
	cols04 := []smry.Predicator{
		pred_khoảng_thời_gian,
	}
	table04 := smry.BuildTable(rows04, cols04, "Số lượng đơn theo tuyến giao hàng", "fulfillment04", "datefrom-dateto&&delivery_route")

	row_tổng_thu_hộ := smry.NewSubject("Tổng tiền thu hộ", "", "SUM(total_cod_amount)", "SUM(total_cod_amount)", nil)

	row_tổng_phí_vận_chuyển := smry.NewSubject("Tổng phí vận chuyển", "", "SUM(shipping_fee_shop)", "SUM(shipping_fee_shop)", nil)

	row_tổng_tiền_bồi_hoàn := smry.NewSubject("Tổng tiền bồi hoàn", "", "SUM(actual_compensation_amount)", "SUM(actual_compensation_amount)", nil)

	rows05 := []smry.Subject{
		row_tổng_thu_hộ.Combine("", pre_tính_cod),
		row_tổng_phí_vận_chuyển.Combine("", pre_tính_cod),
	}
	cols05 := buildRowPerDate(dateFrom, dateTo)
	table05 := smry.BuildTable(rows05, cols05, "Giá trị giao hàng theo ngày", "fulfillment05", "datefrom-dateto&&total_amount")

	row_tiền_thu_hộ := smry.NewSubject("Tổng tiền thu hộ", "", "sum", "SUM(cod_amount)", nil)

	row_tổng_phí := smry.NewSubject("Tổng phí vận chuyển", "", "sum", "SUM(total_fee)", nil)

	rows08 := []smry.Subject{
		row_tiền_thu_hộ.Combine("", pred_chưa_hủy),
		row_tổng_phí.Combine("", pred_chưa_hủy),
	}
	cols08 := buildRowPerDate(dateFrom, dateTo)
	table08 := smry.BuildTable(rows08, cols08, "Giá trị giao hàng", "fulfillment05", "shipnow_fulfillment", "datefrom-dateto&&total_amount")

	pred_bồi_hoàn := smry.Predicate{
		Spec: "shipping_state=undeliverable",
		Expr: sq.NewExpr("shipping_state='undeliverable'"),
	}
	pred_đã_giao := smry.Predicate{
		Spec: "shipping_state=delivered",
		Expr: sq.NewExpr("shipping_state='delivered'"),
	}
	pred_thành_công := smry.Predicate{
		Spec: "shipping_state=delivered|undeliverable",
		Expr: sq.NewExpr("shipping_state in ('delivered', 'undeliverable')"),
	}

	rows06 := []smry.Subject{
		row_tổng_thu_hộ.Combine("", pred_đã_giao, pred_khoảng_thời_gian, pred_chưa_hủy),
		row_tổng_phí_vận_chuyển.Combine("", pred_khoảng_thời_gian, pred_chưa_hủy),
		row_tổng_tiền_bồi_hoàn.Combine("", pred_bồi_hoàn, pred_khoảng_thời_gian, pred_chưa_hủy),
	}
	cols06 := []smry.Predicator{
		pred_thành_công,
		pred_trả_hàng,
	}
	table06 := smry.BuildTable(rows06, cols06, "Thống kê giá trị đối soát", "fulfillment06", "money&&shipping_state")

	pred_đã_đối_soát := smry.Predicate{
		Label: "chưa đối soát",
		Spec:  "money_transaction_id!=nil,cod_etop_transfered_at!=nil",
		Expr:  sq.NewExpr("money_transaction_id is not null and cod_etop_transfered_at is not null"),
	}
	pred_đã_lên_phiên := smry.Predicate{
		Label: "Đã lên phiên",
		Spec:  "money_transaction_id!=nil,cod_etop_transfered_at=nil",
		Expr:  sq.NewExpr("money_transaction_id is not null and cod_etop_transfered_at is null"),
	}
	pred_chưa_lên_phiên := smry.Predicate{
		Label: "chưa lên phiên",
		Spec:  "money_transaction_id=nil,cod_etop_transfered_at=nil",
		Expr:  sq.NewExpr("money_transaction_id is null and cod_etop_transfered_at is null and shipping_state in ('delivered','returning','returned', 'undeliverable')"),
	}
	pred_không_bồi_hoàn := smry.Predicate{
		Label: "Không bồi hoàn",
		Spec:  "shipping_state!=undeliverable",
		Expr:  sq.NewExpr("shipping_state in ('delivered','returning','returned')"),
	}
	rows07 := []smry.Subject{
		row_tổng_thu_hộ.Combine("", pred_không_bồi_hoàn, pred_khoảng_thời_gian, pred_chưa_hủy),
		row_tổng_phí_vận_chuyển.Combine("", pred_khoảng_thời_gian, pred_chưa_hủy),
		row_tổng_tiền_bồi_hoàn.Combine("", pred_bồi_hoàn, pred_khoảng_thời_gian, pred_chưa_hủy),
	}
	cols07 := []smry.Predicator{
		pred_đã_đối_soát,
		pred_đã_lên_phiên,
		pred_chưa_lên_phiên,
	}
	table07 := smry.BuildTable(rows07, cols07, "Đối soát", "fulfillment07", "money_transaction")
	suryFfm := []*smry.Table{table00, table01, table02, table03, table04, table05, table06, table07}
	var tablesShipnowFfm []*smry.Table
	for _, value := range suryFfm {
		if Contain(value.Tags, "shipnow_fulfillment") {
			var tablesSn = smry.BuildTable(value.Rows, value.Cols, value.Label, value.Tags...)
			tablesShipnowFfm = append(tablesShipnowFfm, tablesSn)
		}
	}
	tablesShipnowFfm = append(tablesShipnowFfm, table08)
	tablesShipnowFfm = append(tablesShipnowFfm, table10)

	return suryFfm, tablesShipnowFfm

}

func buildRowPerDate(dateFrom time.Time, dateTo time.Time) []smry.Predicator {
	var result []smry.Predicator
	var timeStart = dateFrom
	for timeStart.Before(dateTo) {
		result = append(result, smry.Predicate{
			Spec:  timeStart.Format("2006-01-02"),
			Label: timeStart.Format("2006-01-02"),
			Expr:  sq.NewExpr("created_at >= ? and ? > created_at", timeStart, timeStart.Add(24*time.Hour)),
		})
		timeStart = timeStart.Add(24 * time.Hour)
	}
	return result
}
