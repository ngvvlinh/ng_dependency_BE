package sqlstore

import (
	"context"
	"fmt"
	"strings"
	"time"

	"o.o/api/fabo/summary"
	commonsummary "o.o/api/summary"
	"o.o/backend/com/main/identity/model"
	"o.o/backend/com/summary/util"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sq/core"
	smry "o.o/backend/pkg/etop/logic/summary"
	"o.o/capi/dot"
)

const (
	minimumTablesReturn = 2
)

type SummaryStore struct {
	db    *cmsql.Database
	query cmsql.QueryFactory
	preds []interface{}
}

func NewSummaryStore(db *cmsql.Database) SummaryStore {
	return SummaryStore{db: db}
}

func (s *SummaryStore) SummarizeShop(ctx context.Context, req *summary.SummaryShopArgs) (_ []*commonsummary.SummaryTable, err error) {
	dateFrom, dateTo, shopID := req.DateFrom, req.DateTo, req.ShopID

	//
	// fulfillment and customer
	//
	tablesFfm := buildTablesFulfillment(dateFrom, dateTo)
	if err := s.execQuery(ctx, tablesFfm, shopID, "fulfillment"); err != nil {
		return nil, err
	}

	tableCustomersTotal := buildTableCustomersTotal(dateFrom, dateTo)
	if err := s.execQuery(ctx, []*smry.Table{tableCustomersTotal}, shopID, "shop_customer"); err != nil {
		return nil, err
	}

	var tableFfmPerDay, tableFfmTotal *smry.Table
	{
		if len(tablesFfm) < minimumTablesReturn {
			return nil, cm.Errorf(cm.Internal, nil, "number of tablesFfm is invalid")
		}

		tableFfmTotal = tablesFfm[0]
		tableFfmTotal, err = util.MergeTables([]*smry.Table{tableFfmTotal, tableCustomersTotal}, util.SameCols)
		if err != nil {
			return nil, cm.Errorf(cm.Internal, err, "Can't merge tables")
		}

		tableFfmPerDay = tablesFfm[1]

		// merge data of more tables per-day
		// reason: postgres just return 1664 entires each request
		// then separate time into more time series (from-to)
		for i := 2; i < len(tablesFfm); i++ {
			table := tablesFfm[i]
			tableFfmPerDay.Data = append(tableFfmPerDay.Data, table.Data...)
		}
	}

	//
	// dashboard staff
	//
	users, userIDs, err := s.getAllStaffs(shopID)
	if err != nil {
		return nil, err
	}

	tableFfmStaffs := buildTableFfmByStaffs(dateFrom, dateTo, userIDs)
	if err := s.execQuery(ctx, []*smry.Table{tableFfmStaffs}, shopID, "fulfillment"); err != nil {
		return nil, err
	}

	tableCustomerStaffs := buildTableNewCustomerByStaffs(dateFrom, dateTo, userIDs)
	if err := s.execQuery(ctx, []*smry.Table{tableCustomerStaffs}, shopID, "shop_customer"); err != nil {
		return nil, err
	}

	tableCustomerNameStaffs := buildTableCustomerNameByStaffs(dateFrom, dateTo, userIDs, users)

	externalPageIDs, err := s.getAllFbExternalPageIDs(shopID)
	if err != nil {
		return nil, err
	}

	tableFbUsersWereAdvised, err := buildTableFbUsersWereAdvised(s.db, dateFrom, dateTo, userIDs, externalPageIDs)
	if err != nil {
		return nil, err
	}

	tableMessagesStaffs, err := buildTableMessagesByStaffs(s.db, dateFrom, dateTo, userIDs, externalPageIDs)
	if err != nil {
		return nil, err
	}

	tableCommentsStaffs, err := buildTableCommentsByStaffs(s.db, dateFrom, dateTo, userIDs, externalPageIDs)
	if err != nil {
		return nil, err
	}

	newTableStaffs, err := util.MergeTables([]*smry.Table{
		tableFfmStaffs,
		tableCustomerStaffs,
		tableCustomerNameStaffs,
		tableMessagesStaffs,
		tableCommentsStaffs,
		tableFbUsersWereAdvised,
	}, util.SameCols)
	if err != nil {
		return nil, cm.Errorf(cm.Internal, err, "Can't merge tables")
	}
	newTableStaffs.Label = "K???t qu??? t???ng k???t t???ng nh??n vi??n"

	tables := []*smry.Table{tableFfmPerDay, tableFfmTotal, newTableStaffs}

	resTablesFfm := util.BuildResponse(tables)
	return resTablesFfm, nil
}

func (s *SummaryStore) SummarizeShopToday(ctx context.Context, req *summary.SummaryShopArgs) (_ []*commonsummary.SummaryTable, err error) {
	shopID := req.ShopID

	//
	// fulfillment and customer
	//
	tableFfmToday := buildTablesFulfillmentToday()
	if err := s.execQuery(ctx, []*smry.Table{tableFfmToday}, shopID, "fulfillment"); err != nil {
		return nil, err
	}

	tableCustomersToday := buildTableCustomersToday()
	if err := s.execQuery(ctx, []*smry.Table{tableCustomersToday}, shopID, "shop_customer"); err != nil {
		return nil, err
	}

	tableFfmToday, err = util.MergeTables([]*smry.Table{tableFfmToday, tableCustomersToday}, util.SameCols)
	if err != nil {
		return nil, cm.Errorf(cm.Internal, err, "Can't merge tables")
	}

	resTablesFfm := util.BuildResponse([]*smry.Table{tableFfmToday})
	return resTablesFfm, nil
}

func (s SummaryStore) getAllStaffs(shopID dot.ID) ([]*model.User, []dot.ID, error) {
	var users model.Users

	err := s.db.
		Where("id in (SELECT au.user_id FROM account_user au WHERE au.deleted_at IS NULL AND au.account_id = ?)", shopID).
		Find(&users)
	if err != nil {
		return nil, nil, err
	}

	var userIDs []dot.ID
	for _, user := range users {
		userIDs = append(userIDs, user.ID)
	}

	return users, userIDs, nil
}

func (s SummaryStore) getAllFbExternalPageIDs(shopID dot.ID) (externalPageIDs []string, _ error) {
	dbRows, err := s.db.
		SQL("SELECT DISTINCT external_id "+
			"FROM fb_external_page ").
		Where("shop_id = ? AND status = 1", shopID).
		Clone().
		Query()
	if err != nil {
		return nil, err
	}

	var externalPageID string
	for dbRows.Next() {
		err := dbRows.Scan(&externalPageID)
		if err != nil {
			return nil, err
		}

		externalPageIDs = append(externalPageIDs, externalPageID)
	}

	return externalPageIDs, nil
}

func (s SummaryStore) execQuery(ctx context.Context, tables []*smry.Table, shopID dot.ID, tableName string) error {
	for _, table := range tables {
		builder := smry.NewSummaryQueryBuilder(tableName)
		for i := range table.Data {
			// must always use [i] because we want to take the address
			builder.AddCell(&table.Data[i].Subject, (*core.Int64)(&table.Data[i].Value))
		}
		q := s.db.SQL(builder).WithContext(ctx)
		if shopID != 0 {
			q = q.Where("shop_id = ?", shopID)
		}
		q = q.Clone()
		if err := q.Scan(builder.ScanArgs...); err != nil {
			return err
		}
	}
	return nil
}

// result contains [tableFfmTotal, tableFfmPerDay]
func buildTablesFulfillment(dateFrom, dateTo time.Time) (tablesFfm []*smry.Table) {
	pred_doanh_thu_s???n_ph???m := smry.Predicate{
		Label: "Doanh thu s???n ph???m",
		Spec:  "shop_confirm = 1 and shipping_code != null and shipping_state not in (default, cancelled) and status not in (0, -1)",
		Expr:  sq.NewExpr("shop_confirm = 1 and shipping_code is not null and shipping_state not in ('default', 'cancelled') and status not in (0, -1)"),
	}

	row_doanh_thu_s???n_ph???m := smry.NewSubject("T???ng danh thu s???n ph???m", "", "SUM(basket_value)", "SUM(basket_value)", nil)

	pred_doanh_thu_COD := smry.Predicate{
		Label: "Doanh thu COD",
		Spec:  "shop_confirm = 1 and shipping_code != null and shipping_state not in (default, cancelled) and status not in (0, -1)",
		Expr:  sq.NewExpr("shop_confirm = 1 and shipping_code is not null and shipping_state not in ('default', 'cancelled') and status not in (0, -1)"),
	}

	row_doanh_thu_COD := smry.NewSubject("T???ng doanh thu COD", "", "SUM(total_cod_amount)", "SUM(total_cod_amount)", nil)

	pred_????_ch???t_????n := smry.Predicate{
		Label: "???? ch???t ????n",
		Spec:  "shop_confirm = 1 and shipping_code != null and shipping_state not in (default, cancelled) and status not in (0, -1)",
		Expr:  sq.NewExpr("shop_confirm = 1 and shipping_code is not null and shipping_state not in ('default', 'cancelled') and status not in (0, -1)"),
	}

	pred_????_b??n_giao_NVC := smry.Predicate{
		Label: "???? b??n giao NVC",
		Spec:  "shipping_state in (holding, delivering, undeliverable)",
		Expr:  sq.NewExpr("shipping_state in ('holding', 'delivering', 'undeliverable')"),
	}

	pred_????_giao_th??nh_c??ng := smry.Predicate{
		Label: "???? giao th??nh c??ng",
		Spec:  "shipping_state in (delivered)",
		Expr:  sq.NewExpr("shipping_state = 'delivered'"),
	}

	pred_tr???_h??ng_v??? := smry.Predicate{
		Label: "Tr??? h??ng v???",
		Spec:  "shipping_state in (returning, returned)",
		Expr:  sq.NewExpr("shipping_state in ('returning', 'returned')"),
	}

	pred_hu??? := smry.Predicate{
		Label: "Hu???",
		Spec:  "shipping_state = cancelled",
		Expr:  sq.NewExpr("shipping_state = 'cancelled'"),
	}

	row_t???ng_????n := smry.NewSubject("T???ng ????n", "", "count", "COUNT(*)", nil)

	//
	// today
	//
	rows := []smry.Subject{
		row_doanh_thu_s???n_ph???m.Combine("Doanh thu s???n ph???m", pred_doanh_thu_s???n_ph???m),
		row_doanh_thu_COD.Combine("Doanh thu COD", pred_doanh_thu_COD),
		row_t???ng_????n.Combine("???? ch???t ????n", pred_????_ch???t_????n),
		row_t???ng_????n.Combine("???? b??n giao NVC", pred_????_b??n_giao_NVC),
		row_t???ng_????n.Combine("???? giao th??nh c??ng", pred_????_giao_th??nh_c??ng),
		row_t???ng_????n.Combine("Tr??? h??ng v???", pred_tr???_h??ng_v???),
		row_t???ng_????n.Combine("Hu???", pred_hu???),
	}

	//
	// total with column customer
	//
	colsFfmTotal := []smry.Predicator{
		smry.Predicate{
			Label: "T???ng qu??t",
			Spec:  "total",
			Expr:  sq.NewExpr("created_at >= ? and ? > created_at", dateFrom, dateTo),
		},
	}
	tableFfmTotal := smry.BuildTable(rows, colsFfmTotal, "K???t qu??? t???ng qu??t", "fulfillments", "datefrom-dateto", "total")
	tablesFfm = append(tablesFfm, tableFfmTotal)

	// per-day
	isContinue := true

	timeStart, timeEnd := dateFrom, dateFrom
	for isContinue {
		timeEnd = timeStart.Add(200 * 24 * time.Hour)
		if timeEnd.After(dateTo) {
			timeEnd = dateTo
			isContinue = false
		}

		colsPerDay := util.BuildRowPerDate(timeStart, timeEnd)
		tablePerDay := smry.BuildTable(rows, colsPerDay, "K???t qu??? theo ng??y", "fulfillments", "datefrom-dateto", "per-day")
		tablesFfm = append(tablesFfm, tablePerDay)

		timeStart = timeEnd.Add(24 * time.Hour)
	}

	return
}

func buildTablesFulfillmentToday() (tablesFfm *smry.Table) {
	pred_doanh_thu_s???n_ph???m := smry.Predicate{
		Label: "Doanh thu s???n ph???m",
		Spec:  "shop_confirm = 1 and shipping_code != null and shipping_state not in (default, cancelled) and status not in (0, -1)",
		Expr:  sq.NewExpr("shop_confirm = 1 and shipping_code is not null and shipping_state not in ('default', 'cancelled') and status not in (0, -1)"),
	}

	row_doanh_thu_s???n_ph???m := smry.NewSubject("T???ng danh thu s???n ph???m", "", "SUM(basket_value)", "SUM(basket_value)", nil)

	pred_doanh_thu_COD := smry.Predicate{
		Label: "Doanh thu COD",
		Spec:  "shop_confirm = 1 and shipping_code != null and shipping_state not in (default, cancelled) and status not in (0, -1)",
		Expr:  sq.NewExpr("shop_confirm = 1 and shipping_code is not null and shipping_state not in ('default', 'cancelled') and status not in (0, -1)"),
	}

	row_doanh_thu_COD := smry.NewSubject("T???ng doanh thu COD", "", "SUM(total_cod_amount)", "SUM(total_cod_amount)", nil)

	pred_????_ch???t_????n := smry.Predicate{
		Label: "???? ch???t ????n",
		Spec:  "shop_confirm = 1 and shipping_code != null and shipping_state not in (default, cancelled) and status not in (0, -1)",
		Expr:  sq.NewExpr("shop_confirm = 1 and shipping_code is not null and shipping_state not in ('default', 'cancelled') and status not in (0, -1)"),
	}

	pred_????_b??n_giao_NVC := smry.Predicate{
		Label: "???? b??n giao NVC",
		Spec:  "shipping_state in (holding, delivering, undeliverable)",
		Expr:  sq.NewExpr("shipping_state in ('holding', 'delivering', 'undeliverable')"),
	}

	pred_????_giao_th??nh_c??ng := smry.Predicate{
		Label: "???? giao th??nh c??ng",
		Spec:  "shipping_state in (delivered)",
		Expr:  sq.NewExpr("shipping_state = 'delivered'"),
	}

	pred_tr???_h??ng_v??? := smry.Predicate{
		Label: "Tr??? h??ng v???",
		Spec:  "shipping_state in (returning, returned)",
		Expr:  sq.NewExpr("shipping_state in ('returning', 'returned')"),
	}

	pred_hu??? := smry.Predicate{
		Label: "Hu???",
		Spec:  "shipping_state = cancelled",
		Expr:  sq.NewExpr("shipping_state = 'cancelled'"),
	}

	row_t???ng_????n := smry.NewSubject("T???ng ????n", "", "count", "COUNT(*)", nil)

	//
	// today
	//
	rows := []smry.Subject{
		row_doanh_thu_s???n_ph???m.Combine("Doanh thu s???n ph???m", pred_doanh_thu_s???n_ph???m),
		row_doanh_thu_COD.Combine("Doanh thu COD", pred_doanh_thu_COD),
		row_t???ng_????n.Combine("???? ch???t ????n", pred_????_ch???t_????n),
		row_t???ng_????n.Combine("???? b??n giao NVC", pred_????_b??n_giao_NVC),
		row_t???ng_????n.Combine("???? giao th??nh c??ng", pred_????_giao_th??nh_c??ng),
		row_t???ng_????n.Combine("Tr??? h??ng v???", pred_tr???_h??ng_v???),
		row_t???ng_????n.Combine("Hu???", pred_hu???),
	}

	dateFromNow, dateToNow := getToday(time.Now())
	colsFfmToday := []smry.Predicator{
		smry.Predicate{
			Label: "H??m nay",
			Spec:  "today",
			Expr:  sq.NewExpr("created_at >= ? and ? > created_at", dateFromNow, dateToNow),
		},
	}
	return smry.BuildTable(rows, colsFfmToday, "K???t qu??? h??m nay", "fulfillments", "datefrom-dateto", "today")
}

func buildTableCustomersToday() (tableCustomers *smry.Table) {
	pred_kh??ch_h??ng_m???i := smry.Predicate{
		Label: "Kh??ch h??ng m???i",
		Spec:  "deleted_at is null",
		Expr:  sq.NewExpr("deleted_at is null"),
	}
	row_t???ng_kh??ch_h??ng := smry.NewSubject("T???ng kh??ch h??ng m???i", "", "count", "COUNT(*)", nil)

	dateFromNow, dateToNow := getToday(time.Now())
	rowsToday := []smry.Subject{
		row_t???ng_kh??ch_h??ng.Combine("Kh??ch h??ng m???i", pred_kh??ch_h??ng_m???i),
	}
	colsToday := []smry.Predicator{
		smry.Predicate{
			Label: "H??m nay",
			Spec:  "today",
			Expr:  sq.NewExpr("created_at >= ? and ? > created_at", dateFromNow, dateToNow),
		},
	}

	return smry.BuildTable(rowsToday, colsToday, "K???t qu??? c???a ng??y h??m nay", "customers", "datefrom-dateto", "today")
}

func buildTableCustomersTotal(dateFrom, dateTo time.Time) (tableCustomers *smry.Table) {
	pred_kh??ch_h??ng_m???i := smry.Predicate{
		Label: "Kh??ch h??ng m???i",
		Spec:  "deleted_at is null",
		Expr:  sq.NewExpr("deleted_at is null"),
	}

	//
	// total
	//
	row_t???ng_kh??ch_h??ng := smry.NewSubject("T???ng kh??ch h??ng m???i", "", "count", "COUNT(*)", nil)
	rowsTotal := []smry.Subject{
		row_t???ng_kh??ch_h??ng.Combine("Kh??ch h??ng m???i", pred_kh??ch_h??ng_m???i),
	}
	colsTotal := []smry.Predicator{
		smry.Predicate{
			Label: "T???ng qu??t",
			Spec:  "total",
			Expr:  sq.NewExpr("created_at >= ? and ? > created_at", dateFrom, dateTo),
		},
	}

	return smry.BuildTable(rowsTotal, colsTotal, "K???t qu??? t???ng qu??t", "customers", "datefrom-dateto", "total")
}

//
// build tables for staffs
//

func buildTableFfmByStaffs(
	dateFrom, dateTo time.Time, userIDs []dot.ID,
) (tableFfmStaffs *smry.Table) {
	pred_doanh_thu_s???n_ph???m := smry.Predicate{
		Label: "Doanh thu s???n ph???m",
		Spec:  "shop_confirm = 1 and shipping_code != null and shipping_state not in (default, cancelled) and status not in (0, -1)",
		Expr:  sq.NewExpr("shop_confirm = 1 and shipping_code is not null and shipping_state not in ('default', 'cancelled') and status not in (0, -1)"),
	}

	row_doanh_thu_s???n_ph???m := smry.NewSubject("doanh_thu_san_pham", "", "SUM(basket_value)", "SUM(basket_value)", nil)

	pred_doanh_thu_COD := smry.Predicate{
		Label: "Doanh thu COD",
		Spec:  "shop_confirm = 1 and shipping_code != null and shipping_state not in (default, cancelled) and status not in (0, -1)",
		Expr:  sq.NewExpr("shop_confirm = 1 and shipping_code is not null and shipping_state not in ('default', 'cancelled') and status not in (0, -1)"),
	}

	row_doanh_thu_COD := smry.NewSubject("T???ng doanh thu COD", "", "SUM(total_cod_amount)", "SUM(total_cod_amount)", nil)

	pred_t???ng_????n := smry.Predicate{
		Label: "T???ng ????n",
		Spec:  "status != -1 or shipping_state != cancelled",
		Expr:  sq.NewExpr("status != -1 OR shipping_state != 'cancelled'"),
	}

	row_t???ng_????n := smry.NewSubject("T???ng ????n", "", "count", "COUNT(*)", nil)

	rows := []smry.Subject{
		row_doanh_thu_s???n_ph???m.Combine("Doanh thu s???n ph???m", pred_doanh_thu_s???n_ph???m),
		row_doanh_thu_COD.Combine("Doanh thu COD", pred_doanh_thu_COD),
		row_t???ng_????n.Combine("T???ng ????n", pred_t???ng_????n),
	}

	var cols []smry.Predicator
	for _, userID := range userIDs {
		cols = append(cols, smry.Predicate{
			Label: fmt.Sprintf("user_id = %d, from(%s) - to(%s)", userID, dateFrom.Format("2006-01-02"), dateTo.Format("2006-01-02")),
			Spec:  fmt.Sprintf("user_id = %d, from(%s) - to(%s)", userID, dateFrom.Format("2006-01-02"), dateTo.Format("2006-01-02")),
			Expr:  sq.NewExpr("created_at >= ? and ? > created_at and created_by = ?", dateFrom, dateTo, userID),
		})
	}

	return smry.BuildTable(rows, cols, "K???t qu??? t???ng qu??t", "fulfillments", "datefrom-dateto", "total")
}

func buildTableNewCustomerByStaffs(
	dateFrom, dateTo time.Time, userIDs []dot.ID,
) (tableCustomers *smry.Table) {
	pred_kh??ch_m???i := smry.Predicate{
		Label: "Kh??ch h??ng m???i",
		Spec:  "deleted_at is null",
		Expr:  sq.NewExpr("deleted_at is null"),
	}

	row_t???ng_kh??ch_m???i := smry.NewSubject("T???ng Kh??ch h??ng m???i", "", "count", "COUNT(*)", nil)

	rows := []smry.Subject{
		row_t???ng_kh??ch_m???i.Combine("Kh??ch h??ng m???i", pred_kh??ch_m???i),
	}

	var cols []smry.Predicator
	for _, userID := range userIDs {
		cols = append(cols, smry.Predicate{
			Label: fmt.Sprintf("user_id = %d, from(%s) - to(%s)", userID, dateFrom.Format("2006-01-02"), dateTo.Format("2006-01-02")),
			Spec:  fmt.Sprintf("user_id = %d, from(%s) - to(%s)", userID, dateFrom.Format("2006-01-02"), dateTo.Format("2006-01-02")),
			Expr:  sq.NewExpr("created_at >= ? and ? > created_at and created_by = ?", dateFrom, dateTo, userID),
		})
	}
	return smry.BuildTable(rows, cols, "K???t qu??? t???ng qu??t", "shop_customer", "datefrom-dateto", "total")
}

func buildTableCustomerNameByStaffs(
	dateFrom, dateTo time.Time,
	userIDs []dot.ID, users []*model.User,
) (tableCustomerName *smry.Table) {
	pred_t??n_nh??n_vi??n := smry.Predicate{
		Label: "T??n nh??n vi??n",
		Spec:  "t??n nh??n vi??n",
	}

	row_t??n_nh??n_vi??n := smry.NewSubject("T??n nh??n vi??n", "", "", "", nil)

	rows := []smry.Subject{
		row_t??n_nh??n_vi??n.Combine("T??n nh??n vi??n", pred_t??n_nh??n_vi??n),
	}

	var cols []smry.Predicator
	for _, userID := range userIDs {
		cols = append(cols, smry.Predicate{
			Label: fmt.Sprintf("user_id = %d, from(%s) - to(%s)", userID, dateFrom.Format("2006-01-02"), dateTo.Format("2006-01-02")),
			Spec:  fmt.Sprintf("user_id = %d, from(%s) - to(%s)", userID, dateFrom.Format("2006-01-02"), dateTo.Format("2006-01-02")),
		})
	}

	tableCustomerName = smry.BuildTable(rows, cols, "K???t qu??? t???ng qu??t", "shop_customer", "datefrom-dateto", "total")
	for idx, user := range users {
		tableCustomerName.Data[idx].Value = user.ID.Int64()
	}

	return
}

func buildTableMessagesByStaffs(
	db *cmsql.Database, dateFrom, dateTo time.Time, userIDs []dot.ID, externalPageIDs []string,
) (tableMessages *smry.Table, _ error) {
	pred_tin_nh???n_????_g???i := smry.Predicate{
		Label: "Tin nh???n ???? g???i",
		Spec:  "tin_nhan_da_gui",
		Expr:  sq.NewExpr("deleted_at is null"),
	}

	row_t???ng := smry.NewSubject("T???ng kh??ch tin nh???n ???? g???i", "", "count", "COUNT(external_id)", nil)

	rows := []smry.Subject{
		row_t???ng.Combine("Tin nh???n ???? g???i", pred_tin_nh???n_????_g???i),
	}

	var cols []smry.Predicator
	for _, userID := range userIDs {
		cols = append(cols, smry.Predicate{
			Label: fmt.Sprintf("user_id = %d, from(%s) - to(%s)", userID, dateFrom.Format("2006-01-02"), dateTo.Format("2006-01-02")),
			Spec:  fmt.Sprintf("user_id = %d, from(%s) - to(%s)", userID, dateFrom.Format("2006-01-02"), dateTo.Format("2006-01-02")),
			Expr:  sq.NewExpr("external_created_time >= ? and ? > external_created_time and created_by = ?", dateFrom, dateTo, userID),
		})
	}
	table := smry.BuildTable(rows, cols, "K???t qu??? t???ng qu??t", "fb_external_message", "datefrom-dateto", "total")

	dbRows, err := db.
		SQL("SELECT COUNT(external_id), created_by "+
			"FROM fb_external_message ").
		Where("external_created_time >= ? AND external_created_time < ? and deleted_at is null", dateFrom, dateTo).
		In("created_by", userIDs).
		In("external_page_id", externalPageIDs).
		GroupBy("created_by").
		Clone().
		Query()
	if err != nil {
		return nil, err
	}

	mapCreatedByAndCountMessages := make(map[dot.ID]int64)
	var (
		createdBy     dot.ID
		countMessages int64
	)
	for dbRows.Next() {
		err := dbRows.Scan(&countMessages, &createdBy)
		if err != nil {
			return nil, err
		}
		mapCreatedByAndCountMessages[createdBy] = countMessages
	}

	for i, col := range table.Cols {
		label := col.GetLabel()
		userIDStr := label[len("user_id = "):strings.Index(label, ",")]
		userID, err := dot.ParseID(userIDStr)
		if err != nil {
			return nil, err
		}

		if countMessages, ok := mapCreatedByAndCountMessages[userID]; ok {
			table.Data[i].Value = countMessages
		}
	}

	return table, nil
}

func buildTableCommentsByStaffs(
	db *cmsql.Database, dateFrom, dateTo time.Time, userIDs []dot.ID, externalPageIDs []string,
) (tableMessages *smry.Table, _ error) {
	pred_tin_nh???n_????_g???i := smry.Predicate{
		Label: "Comment ???? g???i",
		Spec:  "comment_da_gui",
		Expr:  sq.NewExpr("deleted_at is null"),
	}

	row_t???ng := smry.NewSubject("T???ng kh??ch comment ???? g???i", "", "count", "COUNT(external_id)", nil)

	rows := []smry.Subject{
		row_t???ng.Combine("Comment ???? g???i", pred_tin_nh???n_????_g???i),
	}

	var cols []smry.Predicator
	for _, userID := range userIDs {
		cols = append(cols, smry.Predicate{
			Label: fmt.Sprintf("user_id = %d, from(%s) - to(%s)", userID, dateFrom.Format("2006-01-02"), dateTo.Format("2006-01-02")),
			Spec:  fmt.Sprintf("user_id = %d, from(%s) - to(%s)", userID, dateFrom.Format("2006-01-02"), dateTo.Format("2006-01-02")),
			Expr:  sq.NewExpr("external_created_time >= ? and ? > external_created_time and created_by = ?", dateFrom, dateTo, userID),
		})
	}
	table := smry.BuildTable(rows, cols, "K???t qu??? t???ng qu??t", "fb_external_comment", "datefrom-dateto", "total")

	dbRows, err := db.
		SQL("SELECT COUNT(external_id), created_by "+
			"FROM fb_external_comment ").
		Where("external_created_time >= ? AND external_created_time < ? and deleted_at is null", dateFrom, dateTo).
		In("created_by", userIDs).
		In("external_page_id", externalPageIDs).
		GroupBy("created_by").
		Clone().
		Query()
	if err != nil {
		return nil, err
	}

	mapCreatedByAndCountComments := make(map[dot.ID]int64)
	var (
		createdBy     dot.ID
		countComments int64
	)
	for dbRows.Next() {
		err := dbRows.Scan(&countComments, &createdBy)
		if err != nil {
			return nil, err
		}
		mapCreatedByAndCountComments[createdBy] = countComments
	}

	for i, col := range table.Cols {
		label := col.GetLabel()
		userIDStr := label[len("user_id = "):strings.Index(label, ",")]
		userID, err := dot.ParseID(userIDStr)
		if err != nil {
			return nil, err
		}

		if countComments, ok := mapCreatedByAndCountComments[userID]; ok {
			table.Data[i].Value = countComments
		}
	}

	return table, nil
}

func buildTableFbUsersWereAdvised(
	db *cmsql.Database, dateFrom, dateTo time.Time,
	userIDs []dot.ID, externalPageIDs []string,
) (table *smry.Table, _ error) {

	pred_kh??ch_h??ng_????_t??_v???n := smry.Predicate{
		Label: "Kh??ch h??ng ???? t?? v???n",
		Spec:  "deleted is null",
		Expr:  sq.NewExpr("fb_external_message.deleted_at is null and fb_external_comment.deleted_at is null"),
	}

	row_t???ng := smry.NewSubject("T???ng kh??ch h??ng ???? t?? v???n", "", "count(khach_hang_da_tu_van)", "count(khach_hang_da_tu_van)", nil)

	rows := []smry.Subject{
		row_t???ng.Combine("Kh??ch h??ng ???? t?? v???n", pred_kh??ch_h??ng_????_t??_v???n),
	}

	var cols []smry.Predicator
	for _, userID := range userIDs {
		cols = append(cols, smry.Predicate{
			Label: fmt.Sprintf("user_id = %d, from(%s) - to(%s)", userID, dateFrom.Format("2006-01-02"), dateTo.Format("2006-01-02")),
			Spec:  fmt.Sprintf("user_id = %d, from(%s) - to(%s)", userID, dateFrom.Format("2006-01-02"), dateTo.Format("2006-01-02")),
			Expr:  sq.NewExpr("user_id = ?, from(?) - to(?)", userID, dateFrom.Format("2006-01-02"), dateTo.Format("2006-01-02")),
		})
	}

	tableFbUsersWereAdvised := smry.BuildTable(rows, cols, "K???t qu??? t???ng qu??t", "khach_hang_da_tu_van", "datefrom-dateto", "total")

	mapCreatedByAndMapFbUserID := make(map[dot.ID]map[dot.ID]bool)
	{
		dbRows, err := db.
			SQL("SELECT DISTINCT fec.external_parent_user_id, fec.created_by "+
				"FROM fb_external_comment fec ").
			Where("fec.external_created_time >= ? AND fec.external_created_time < ? and fec.deleted_at is null", dateFrom, dateTo).
			In("fec.created_by", userIDs).
			In("fec.external_page_id", externalPageIDs).
			Clone().
			Query()
		if err != nil {
			return nil, err
		}

		var fbUserID, createdBy dot.ID
		for dbRows.Next() {
			err := dbRows.Scan(&fbUserID, &createdBy)
			if err != nil {
				return nil, err
			}

			if _, ok := mapCreatedByAndMapFbUserID[createdBy]; !ok {
				mapCreatedByAndMapFbUserID[createdBy] = make(map[dot.ID]bool)
			}
			mapCreatedByAndMapFbUserID[createdBy][fbUserID] = true
		}
	}

	{
		dbRows, err := db.
			SQL("SELECT DISTINCT fec.external_user_id, fem.created_by "+
				"FROM fb_external_message fem "+
				"JOIN fb_external_conversation fec on fem.external_conversation_id = fec.external_id ").
			Where("fem.external_created_time >= ? AND fem.external_created_time < ? and fem.deleted_at is null", dateFrom, dateTo).
			In("fem.created_by", userIDs).
			In("fem.external_page_id", externalPageIDs).
			Clone().
			Query()
		if err != nil {
			return nil, err
		}

		var fbUserID, createdBy dot.ID
		for dbRows.Next() {
			err := dbRows.Scan(&fbUserID, &createdBy)
			if err != nil {
				return nil, err
			}

			if _, ok := mapCreatedByAndMapFbUserID[createdBy]; !ok {
				mapCreatedByAndMapFbUserID[createdBy] = make(map[dot.ID]bool)
			}
			mapCreatedByAndMapFbUserID[createdBy][fbUserID] = true
		}
	}

	for i, userID := range userIDs {
		if mapFbUserID, ok := mapCreatedByAndMapFbUserID[userID]; ok {
			tableFbUsersWereAdvised.Data[i].Value = int64(len(mapFbUserID))
		}
	}

	return tableFbUsersWereAdvised, nil
}

func getToday(now time.Time) (from, to time.Time) {
	dateFrom, dateTo := now, now
	dateFrom = dateFrom.Add(-time.Duration(dateFrom.Second()) * time.Second)
	dateFrom = dateFrom.Add(-time.Duration(dateFrom.Minute()) * time.Minute)
	dateFrom = dateFrom.Add(-time.Duration(dateFrom.Hour()) * time.Hour)

	return dateFrom, dateTo
}
