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
	pred_t???o_h??m_nay := smry.Pred_sau_ng??y("created_at", "H??m nay [%v]", 0)

	pred_????_h???y := smry.Predicate{
		Spec: "status=N",
		Expr: sq.NewExpr("status=?", status5.N),
	}

	pred_ch??a_h???y := smry.Predicate{
		Spec: "status!=-1",
		Expr: sq.NewExpr("status!=?", status5.N),
	}
	pre_t??nh_cod := smry.Predicate{
		Spec: "shipping_state!=returned,returning|status!=-1",
		Expr: sq.NewExpr("shipping_state not in ('returned','returning') and status!=?", status5.N),
	}

	pred_????_l???y := smry.Predicate{
		Spec: "shipping_state=holding,delivering,delivered,underliverable,returning,returned",
		Expr: sq.NewExpr("shipping_state IN ('holding','delivering','delivered','undeliverable','returning','returned')"),
	}
	pred_ch??a_l???y := smry.Predicate{
		Spec: "shipping_state=defaul,created,picking",
		Expr: sq.NewExpr("shipping_state IN ('default','created','picking')"),
	}

	row_t???ng_????n := smry.NewSubject("T???ng ????n", "", "count", "COUNT(*)", nil)

	//table t???ng k???t ng??y
	rows00 := []smry.Subject{
		row_t???ng_????n.Combine(""),
		row_t???ng_????n.Combine("", pred_????_h???y),
		row_t???ng_????n.Combine("", pred_????_l???y, pred_ch??a_h???y),
		row_t???ng_????n.Combine("", pred_ch??a_l???y, pred_ch??a_h???y),
	}
	cols00 := []smry.Predicator{
		pred_t???o_h??m_nay,
	}

	table00 := smry.BuildTable(rows00, cols00, "H??m nay", "fulfillment", "shipnow_fulfillment", "today")

	pred_kho???ng_th???i_gian := smry.Predicate{
		Spec: "datefrom-dateto",
		Expr: sq.NewExpr("created_at >= ? and ? > created_at ", dateFrom, dateTo),
	}

	pred_??ang_x???_l?? := smry.Predicate{
		Spec: "shipping_state=holding,delivering",
		Expr: sq.NewExpr("shipping_state in ('holding', 'delivering')"),
	}
	pred_giao_th??nh_c??ng := smry.Predicate{
		Spec: "shipping_state=delivered",
		Expr: sq.NewExpr("shipping_state in ('delivered')"),
	}
	pred_tr???_h??ng := smry.Predicate{
		Spec: "shipping_state=returned,returning",
		Expr: sq.NewExpr("shipping_state in ('returned', 'returning')"),
	}

	// T???ng k???t theo th???i gian dateFrom - dateTo
	rows01 := []smry.Subject{
		row_t???ng_????n.Combine(""),
		row_t???ng_????n.Combine("", pred_ch??a_l???y, pred_ch??a_h???y),
		row_t???ng_????n.Combine("", pred_??ang_x???_l??, pred_ch??a_h???y),
		row_t???ng_????n.Combine("", pred_giao_th??nh_c??ng, pred_ch??a_h???y),
		row_t???ng_????n.Combine("", pred_tr???_h??ng, pred_ch??a_h???y),
		row_t???ng_????n.Combine("", pred_????_h???y),
	}
	cols01 := []smry.Predicator{
		pred_kho???ng_th???i_gian,
	}
	table01 := smry.BuildTable(rows01, cols01, "T???ng ????n trong th??ng", "fulfillment01", "shipnow_fulfillment", "datefrom-dateto")

	// T???ng ????n t???ng ng??y
	rows02 := []smry.Subject{
		row_t???ng_????n,
	}
	rowsPerDay := util.BuildRowPerDate(dateFrom, dateTo)
	table02 := smry.BuildTable(rows02, rowsPerDay, "S??? l?????ng ????n theo ng??y", "fulfillment02", "shipnow_fulfillment", "datefrom-dateto")

	pred_giao_h??ng_nhanh := smry.Predicate{
		Spec: "shipping_provider=ghn",
		Expr: sq.NewExpr("shipping_provider = 'ghn'"),
	}

	pred_giao_h??ng_ti???t_ki???m := smry.Predicate{
		Spec: "shipping_provider=ghtk",
		Expr: sq.NewExpr("shipping_provider = 'ghtk'"),
	}

	pred_viettel_post := smry.Predicate{
		Spec: "shipping_provider=vtpost",
		Expr: sq.NewExpr("shipping_provider = 'vtpost'"),
	}

	rows03 := []smry.Subject{
		row_t???ng_????n.Combine("", pred_ch??a_h???y),
		row_t???ng_????n.Combine("", pred_giao_h??ng_nhanh, pred_ch??a_h???y),
		row_t???ng_????n.Combine("", pred_giao_h??ng_ti???t_ki???m, pred_ch??a_h???y),
		row_t???ng_????n.Combine("", pred_viettel_post, pred_ch??a_h???y),
	}
	cols03 := []smry.Predicator{
		pred_kho???ng_th???i_gian,
	}
	table03 := smry.BuildTable(rows03, cols03, "S??? l?????ng ????n theo nh?? v???n chuy???n", "fulfillment03", "datefrom-dateto&&provider")

	rows10 := []smry.Subject{
		row_t???ng_????n.Combine(""),
	}
	cols10 := []smry.Predicator{
		pred_kho???ng_th???i_gian,
	}
	table10 := smry.BuildTable(rows10, cols10, "S??? l?????ng ????n theo nh?? v???n chuy???n shipnow", "fulfillment03", "shipnow_fullfillment", "datefrom-dateto&&provider")

	pred_n???i_t???nh := smry.Predicate{
		Spec: "delivery_route=noi_tinh",
		Expr: sq.NewExpr("delivery_route=?", etopmodel.RouteSameProvince),
	}

	pred_n???i_mi???n := smry.Predicate{
		Spec: "delivery_route=noi_mien",
		Expr: sq.NewExpr("delivery_route=?", etopmodel.RouteSameRegion),
	}

	pred_ngo???i_mi???n := smry.Predicate{
		Spec: "delivery_route=toan_quoc",
		Expr: sq.NewExpr("delivery_route=?", etopmodel.RouteNationWide),
	}

	rows04 := []smry.Subject{
		row_t???ng_????n.Combine("", pred_ch??a_h???y),
		row_t???ng_????n.Combine("", pred_n???i_t???nh, pred_ch??a_h???y),
		row_t???ng_????n.Combine("", pred_n???i_mi???n, pred_ch??a_h???y),
		row_t???ng_????n.Combine("", pred_ngo???i_mi???n, pred_ch??a_h???y),
	}
	cols04 := []smry.Predicator{
		pred_kho???ng_th???i_gian,
	}
	table04 := smry.BuildTable(rows04, cols04, "S??? l?????ng ????n theo tuy???n giao h??ng", "fulfillment04", "datefrom-dateto&&delivery_route")

	row_t???ng_thu_h??? := smry.NewSubject("T???ng ti???n thu h???", "", "SUM(total_cod_amount)", "SUM(total_cod_amount)", nil)

	row_t???ng_ph??_v???n_chuy???n := smry.NewSubject("T???ng ph?? v???n chuy???n", "", "SUM(shipping_fee_shop)", "SUM(shipping_fee_shop)", nil)

	row_t???ng_ti???n_b???i_ho??n := smry.NewSubject("T???ng ti???n b???i ho??n", "", "SUM(actual_compensation_amount)", "SUM(actual_compensation_amount)", nil)

	rows05 := []smry.Subject{
		row_t???ng_thu_h???.Combine("", pre_t??nh_cod),
		row_t???ng_ph??_v???n_chuy???n.Combine("", pre_t??nh_cod),
	}
	cols05 := util.BuildRowPerDate(dateFrom, dateTo)
	table05 := smry.BuildTable(rows05, cols05, "Gi?? tr??? giao h??ng theo ng??y", "fulfillment05", "datefrom-dateto&&total_amount")

	row_ti???n_thu_h??? := smry.NewSubject("T???ng ti???n thu h???", "", "sum", "SUM(cod_amount)", nil)

	row_t???ng_ph?? := smry.NewSubject("T???ng ph?? v???n chuy???n", "", "sum", "SUM(total_fee)", nil)

	rows08 := []smry.Subject{
		row_ti???n_thu_h???.Combine("", pred_ch??a_h???y),
		row_t???ng_ph??.Combine("", pred_ch??a_h???y),
	}
	cols08 := util.BuildRowPerDate(dateFrom, dateTo)
	table08 := smry.BuildTable(rows08, cols08, "Gi?? tr??? giao h??ng", "fulfillment05", "shipnow_fulfillment", "datefrom-dateto&&total_amount")

	pred_b???i_ho??n := smry.Predicate{
		Spec: "shipping_state=undeliverable",
		Expr: sq.NewExpr("shipping_state='undeliverable'"),
	}
	pred_????_giao := smry.Predicate{
		Spec: "shipping_state=delivered",
		Expr: sq.NewExpr("shipping_state='delivered'"),
	}
	pred_th??nh_c??ng := smry.Predicate{
		Spec: "shipping_state=delivered|undeliverable",
		Expr: sq.NewExpr("shipping_state in ('delivered', 'undeliverable')"),
	}

	rows06 := []smry.Subject{
		row_t???ng_thu_h???.Combine("", pred_????_giao, pred_kho???ng_th???i_gian, pred_ch??a_h???y),
		row_t???ng_ph??_v???n_chuy???n.Combine("", pred_kho???ng_th???i_gian, pred_ch??a_h???y),
		row_t???ng_ti???n_b???i_ho??n.Combine("", pred_b???i_ho??n, pred_kho???ng_th???i_gian, pred_ch??a_h???y),
	}
	cols06 := []smry.Predicator{
		pred_th??nh_c??ng,
		pred_tr???_h??ng,
	}
	table06 := smry.BuildTable(rows06, cols06, "Th???ng k?? gi?? tr??? ?????i so??t", "fulfillment06", "money&&shipping_state")

	pred_????_?????i_so??t := smry.Predicate{
		Label: "ch??a ?????i so??t",
		Spec:  "money_transaction_id!=nil,cod_etop_transfered_at!=nil",
		Expr:  sq.NewExpr("money_transaction_id is not null and cod_etop_transfered_at is not null"),
	}
	pred_????_l??n_phi??n := smry.Predicate{
		Label: "???? l??n phi??n",
		Spec:  "money_transaction_id!=nil,cod_etop_transfered_at=nil",
		Expr:  sq.NewExpr("money_transaction_id is not null and cod_etop_transfered_at is null"),
	}
	pred_ch??a_l??n_phi??n := smry.Predicate{
		Label: "ch??a l??n phi??n",
		Spec:  "money_transaction_id=nil,cod_etop_transfered_at=nil",
		Expr:  sq.NewExpr("money_transaction_id is null and cod_etop_transfered_at is null and shipping_state in ('delivered','returning','returned', 'undeliverable')"),
	}
	pred_kh??ng_b???i_ho??n := smry.Predicate{
		Label: "Kh??ng b???i ho??n",
		Spec:  "shipping_state!=undeliverable",
		Expr:  sq.NewExpr("shipping_state in ('delivered','returning','returned')"),
	}
	rows07 := []smry.Subject{
		row_t???ng_thu_h???.Combine("", pred_kh??ng_b???i_ho??n, pred_kho???ng_th???i_gian, pred_ch??a_h???y),
		row_t???ng_ph??_v???n_chuy???n.Combine("", pred_kho???ng_th???i_gian, pred_ch??a_h???y),
		row_t???ng_ti???n_b???i_ho??n.Combine("", pred_b???i_ho??n, pred_kho???ng_th???i_gian, pred_ch??a_h???y),
	}
	cols07 := []smry.Predicator{
		pred_????_?????i_so??t,
		pred_????_l??n_phi??n,
		pred_ch??a_l??n_phi??n,
	}
	table07 := smry.BuildTable(rows07, cols07, "?????i so??t", "fulfillment07", "money_transaction")
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
