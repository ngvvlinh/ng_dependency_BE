package sqlstore

import (
	"context"
	"fmt"
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
	minimumTablesReturn = 3
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

	tableCustomersToday := buildTableCustomersToday()
	if err := s.execQuery(ctx, []*smry.Table{tableCustomersToday}, shopID, "shop_customer"); err != nil {
		return nil, err
	}

	tableCustomersTotal := buildTableCustomersTotal(dateFrom, dateTo)
	if err := s.execQuery(ctx, []*smry.Table{tableCustomersTotal}, shopID, "shop_customer"); err != nil {
		return nil, err
	}

	var tableFfmPerDay, tableFfmToday, tableFfmTotal *smry.Table
	{
		if len(tablesFfm) < minimumTablesReturn {
			return nil, cm.Errorf(cm.Internal, nil, "number of tablesFfm is invalid")
		}
		tableFfmToday = tablesFfm[0]
		tableFfmToday, err = util.MergeTables([]*smry.Table{tableFfmToday, tableCustomersToday}, util.SameCols)
		if err != nil {
			return nil, cm.Errorf(cm.Internal, err, "Can't merge tables")
		}

		tableFfmTotal = tablesFfm[1]
		tableFfmTotal, err = util.MergeTables([]*smry.Table{tableFfmTotal, tableCustomersTotal}, util.SameCols)
		if err != nil {
			return nil, cm.Errorf(cm.Internal, err, "Can't merge tables")
		}

		tableFfmPerDay = tablesFfm[2]

		// merge data of more tables per-day
		// reason: postgres just return 1664 entires each request
		// then separate time into more time series (from-to)
		for i := 3; i < len(tablesFfm); i++ {
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

	tableFbUsersWereAdvised, err := buildTableFbUsersWereAdvised(s.db, dateFrom, dateTo, userIDs)
	if err != nil {
		return nil, err
	}

	tableCustomerStaffs := buildTableNewCustomerByStaffs(dateFrom, dateTo, userIDs)
	if err := s.execQuery(ctx, []*smry.Table{tableCustomerStaffs}, shopID, "shop_customer"); err != nil {
		return nil, err
	}

	tableCustomerNameStaffs := buildTableCustomerNameByStaffs(dateFrom, dateTo, userIDs, users)

	tableMessagesStaffs := buildTableMessagesByStaffs(dateFrom, dateTo, userIDs)
	if err := s.execQuery(ctx, []*smry.Table{tableMessagesStaffs}, 0, "fb_external_message"); err != nil {
		return nil, err
	}

	tableCommentsStaffs := buildTableCommentsByStaffs(dateFrom, dateTo, userIDs)
	if err := s.execQuery(ctx, []*smry.Table{tableCommentsStaffs}, 0, "fb_external_comment"); err != nil {
		return nil, err
	}

	newTableStaffs, err := util.MergeTables([]*smry.Table{
		tableFfmStaffs,
		tableFbUsersWereAdvised,
		tableCustomerStaffs,
		tableCustomerNameStaffs,
		tableMessagesStaffs,
		tableCommentsStaffs,
	}, util.SameCols)
	if err != nil {
		return nil, cm.Errorf(cm.Internal, err, "Can't merge tables")
	}
	newTableStaffs.Label = "Kết quả tổng kết từng nhân viên"

	tables := []*smry.Table{tableFfmToday, tableFfmPerDay, tableFfmTotal, newTableStaffs}

	resTablesFfm := util.BuildResponse(tables)
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

// result contains [tableFfmToday, tableFfmTotal, tableFfmPerDay]
func buildTablesFulfillment(dateFrom, dateTo time.Time) (tablesFfm []*smry.Table) {
	pred_doanh_thu_sản_phẩm := smry.Predicate{
		Label: "Doanh thu sản phẩm",
		Spec:  "shop_confirm = 1 and shipping_code != null and shipping_state not in (default, cancelled) and status not in (0, -1)",
		Expr:  sq.NewExpr("shop_confirm = 1 and shipping_code is not null and shipping_state not in ('default', 'cancelled') and status not in (0, -1)"),
	}

	row_doanh_thu_sản_phẩm := smry.NewSubject("Tổng danh thu sản phẩm", "", "SUM(basket_value)", "SUM(basket_value)", nil)

	pred_doanh_thu_COD := smry.Predicate{
		Label: "Doanh thu COD",
		Spec:  "shop_confirm = 1 and shipping_code != null and shipping_state not in (default, cancelled) and status not in (0, -1)",
		Expr:  sq.NewExpr("shop_confirm = 1 and shipping_code is not null and shipping_state not in ('default', 'cancelled') and status not in (0, -1)"),
	}

	row_doanh_thu_COD := smry.NewSubject("Tổng doanh thu COD", "", "SUM(total_cod_amount)", "SUM(total_cod_amount)", nil)

	pred_đã_chốt_đơn := smry.Predicate{
		Label: "Đã chốt đơn",
		Spec:  "shop_confirm = 1 and shipping_code != null and shipping_state not in (default, cancelled) and status not in (0, -1)",
		Expr:  sq.NewExpr("shop_confirm = 1 and shipping_code is not null and shipping_state not in ('default', 'cancelled') and status not in (0, -1)"),
	}

	pred_đã_bàn_giao_NVC := smry.Predicate{
		Label: "Đã bàn giao NVC",
		Spec:  "shipping_state in (holding, delivering, undeliverable)",
		Expr:  sq.NewExpr("shipping_state in ('holding', 'delivering', 'undeliverable')"),
	}

	pred_đã_giao_thành_công := smry.Predicate{
		Label: "Đã giao thành công",
		Spec:  "shipping_state in (delivered)",
		Expr:  sq.NewExpr("shipping_state = 'delivered'"),
	}

	pred_trả_hàng_về := smry.Predicate{
		Label: "Trả hàng về",
		Spec:  "shipping_state in (returning, returned)",
		Expr:  sq.NewExpr("shipping_state in ('returning', 'returned')"),
	}

	pred_huỷ := smry.Predicate{
		Label: "Huỷ",
		Spec:  "shipping_state = cancelled",
		Expr:  sq.NewExpr("shipping_state = 'cancelled'"),
	}

	row_tổng_đơn := smry.NewSubject("Tổng đơn", "", "count", "COUNT(*)", nil)

	//
	// today
	//
	rows := []smry.Subject{
		row_doanh_thu_sản_phẩm.Combine("Doanh thu sản phẩm", pred_doanh_thu_sản_phẩm),
		row_doanh_thu_COD.Combine("Doanh thu COD", pred_doanh_thu_COD),
		row_tổng_đơn.Combine("Đã chốt đơn", pred_đã_chốt_đơn),
		row_tổng_đơn.Combine("Đã bàn giao NVC", pred_đã_bàn_giao_NVC),
		row_tổng_đơn.Combine("Đã giao thành công", pred_đã_giao_thành_công),
		row_tổng_đơn.Combine("Trả hàng về", pred_trả_hàng_về),
		row_tổng_đơn.Combine("Huỷ", pred_huỷ),
	}

	dateFromNow, dateToNow := getToday(time.Now())
	colsFfmToday := []smry.Predicator{
		smry.Predicate{
			Label: "Hôm nay",
			Spec:  "today",
			Expr:  sq.NewExpr("created_at >= ? and ? > created_at", dateFromNow, dateToNow),
		},
	}
	tableFfmToday := smry.BuildTable(rows, colsFfmToday, "Kết quả hôm nay", "fulfillments", "datefrom-dateto", "today")
	tablesFfm = append(tablesFfm, tableFfmToday)

	//
	// total with column customer
	//
	colsFfmTotal := []smry.Predicator{
		smry.Predicate{
			Label: "Tổng quát",
			Spec:  "total",
			Expr:  sq.NewExpr("created_at >= ? and ? > created_at", dateFrom, dateTo),
		},
	}
	tableFfmTotal := smry.BuildTable(rows, colsFfmTotal, "Kết quả tổng quát", "fulfillments", "datefrom-dateto", "total")
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
		tablePerDay := smry.BuildTable(rows, colsPerDay, "Kết quả theo ngày", "fulfillments", "datefrom-dateto", "per-day")
		tablesFfm = append(tablesFfm, tablePerDay)

		timeStart = timeEnd.Add(24 * time.Hour)
	}

	return
}

func buildTableCustomersToday() (tableCustomers *smry.Table) {
	pred_khách_hàng_mới := smry.Predicate{
		Label: "Khách hàng mới",
		Spec:  "deleted_at is null",
		Expr:  sq.NewExpr("deleted_at is null"),
	}
	row_tổng_khách_hàng := smry.NewSubject("Tổng khách hàng mới", "", "count", "COUNT(*)", nil)

	dateFromNow, dateToNow := getToday(time.Now())
	rowsToday := []smry.Subject{
		row_tổng_khách_hàng.Combine("Khách hàng mới", pred_khách_hàng_mới),
	}
	colsToday := []smry.Predicator{
		smry.Predicate{
			Label: "Hôm nay",
			Spec:  "today",
			Expr:  sq.NewExpr("created_at >= ? and ? > created_at", dateFromNow, dateToNow),
		},
	}

	return smry.BuildTable(rowsToday, colsToday, "Kết quả của ngày hôm nay", "customers", "datefrom-dateto", "today")
}

func buildTableCustomersTotal(dateFrom, dateTo time.Time) (tableCustomers *smry.Table) {
	pred_khách_hàng_mới := smry.Predicate{
		Label: "Khách hàng mới",
		Spec:  "deleted_at is null",
		Expr:  sq.NewExpr("deleted_at is null"),
	}

	//
	// total
	//
	row_tổng_khách_hàng := smry.NewSubject("Tổng khách hàng mới", "", "count", "COUNT(*)", nil)
	rowsTotal := []smry.Subject{
		row_tổng_khách_hàng.Combine("Khách hàng mới", pred_khách_hàng_mới),
	}
	colsTotal := []smry.Predicator{
		smry.Predicate{
			Label: "Tổng quát",
			Spec:  "total",
			Expr:  sq.NewExpr("created_at >= ? and ? > created_at", dateFrom, dateTo),
		},
	}

	return smry.BuildTable(rowsTotal, colsTotal, "Kết quả tổng quát", "customers", "datefrom-dateto", "total")
}

//
// build tables for staffs
//

func buildTableFfmByStaffs(
	dateFrom, dateTo time.Time, userIDs []dot.ID,
) (tableFfmStaffs *smry.Table) {
	pred_doanh_thu_sản_phẩm := smry.Predicate{
		Label: "Doanh thu sản phẩm",
		Spec:  "shop_confirm = 1 and shipping_code != null and shipping_state not in (default, cancelled) and status not in (0, -1)",
		Expr:  sq.NewExpr("shop_confirm = 1 and shipping_code is not null and shipping_state not in ('default', 'cancelled') and status not in (0, -1)"),
	}

	row_doanh_thu_sản_phẩm := smry.NewSubject("doanh_thu_san_pham", "", "SUM(basket_value)", "SUM(basket_value)", nil)

	pred_doanh_thu_COD := smry.Predicate{
		Label: "Doanh thu COD",
		Spec:  "shop_confirm = 1 and shipping_code != null and shipping_state not in (default, cancelled) and status not in (0, -1)",
		Expr:  sq.NewExpr("shop_confirm = 1 and shipping_code is not null and shipping_state not in ('default', 'cancelled') and status not in (0, -1)"),
	}

	row_doanh_thu_COD := smry.NewSubject("Tổng doanh thu COD", "", "SUM(total_cod_amount)", "SUM(total_cod_amount)", nil)

	pred_tổng_đơn := smry.Predicate{
		Label: "Tổng đơn",
		Spec:  "status != -1 or shipping_state != cancelled",
		Expr:  sq.NewExpr("status != -1 OR shipping_state != 'cancelled'"),
	}

	row_tổng_đơn := smry.NewSubject("Tổng đơn", "", "count", "COUNT(*)", nil)

	rows := []smry.Subject{
		row_doanh_thu_sản_phẩm.Combine("Doanh thu sản phẩm", pred_doanh_thu_sản_phẩm),
		row_doanh_thu_COD.Combine("Doanh thu COD", pred_doanh_thu_COD),
		row_tổng_đơn.Combine("Tổng đơn", pred_tổng_đơn),
	}

	var cols []smry.Predicator
	for _, userID := range userIDs {
		cols = append(cols, smry.Predicate{
			Label: fmt.Sprintf("user_id = %d, from(%s) - to(%s)", userID, dateFrom.Format("2006-01-02"), dateTo.Format("2006-01-02")),
			Spec:  fmt.Sprintf("user_id = %d, from(%s) - to(%s)", userID, dateFrom.Format("2006-01-02"), dateTo.Format("2006-01-02")),
			Expr:  sq.NewExpr("created_at >= ? and ? > created_at and created_by = ?", dateFrom, dateTo, userID),
		})
	}

	return smry.BuildTable(rows, cols, "Kết quả tổng quát", "fulfillments", "datefrom-dateto", "total")
}

func buildTableNewCustomerByStaffs(
	dateFrom, dateTo time.Time, userIDs []dot.ID,
) (tableCustomers *smry.Table) {
	pred_khách_mới := smry.Predicate{
		Label: "Khách hàng mới",
		Spec:  "deleted_at is null",
		Expr:  sq.NewExpr("deleted_at is null"),
	}

	row_tổng_khách_mới := smry.NewSubject("Tổng Khách hàng mới", "", "count", "COUNT(*)", nil)

	rows := []smry.Subject{
		row_tổng_khách_mới.Combine("Khách hàng mới", pred_khách_mới),
	}

	var cols []smry.Predicator
	for _, userID := range userIDs {
		cols = append(cols, smry.Predicate{
			Label: fmt.Sprintf("user_id = %d, from(%s) - to(%s)", userID, dateFrom.Format("2006-01-02"), dateTo.Format("2006-01-02")),
			Spec:  fmt.Sprintf("user_id = %d, from(%s) - to(%s)", userID, dateFrom.Format("2006-01-02"), dateTo.Format("2006-01-02")),
			Expr:  sq.NewExpr("created_at >= ? and ? > created_at and created_by = ?", dateFrom, dateTo, userID),
		})
	}
	return smry.BuildTable(rows, cols, "Kết quả tổng quát", "shop_customer", "datefrom-dateto", "total")
}

func buildTableCustomerNameByStaffs(
	dateFrom, dateTo time.Time,
	userIDs []dot.ID, users []*model.User,
) (tableCustomerName *smry.Table) {
	pred_tên_nhân_viên := smry.Predicate{
		Label: "Tên nhân viên",
		Spec:  "tên nhân viên",
	}

	row_tên_nhân_viên := smry.NewSubject("Tên nhân viên", "", "", "", nil)

	rows := []smry.Subject{
		row_tên_nhân_viên.Combine("Tên nhân viên", pred_tên_nhân_viên),
	}

	var cols []smry.Predicator
	for _, userID := range userIDs {
		cols = append(cols, smry.Predicate{
			Label: fmt.Sprintf("user_id = %d, from(%s) - to(%s)", userID, dateFrom.Format("2006-01-02"), dateTo.Format("2006-01-02")),
			Spec:  fmt.Sprintf("user_id = %d, from(%s) - to(%s)", userID, dateFrom.Format("2006-01-02"), dateTo.Format("2006-01-02")),
		})
	}

	tableCustomerName = smry.BuildTable(rows, cols, "Kết quả tổng quát", "shop_customer", "datefrom-dateto", "total")
	for idx, user := range users {
		tableCustomerName.Data[idx].Value = user.ID.Int64()
	}

	return
}

func buildTableMessagesByStaffs(
	dateFrom, dateTo time.Time, userIDs []dot.ID,
) (tableMessages *smry.Table) {
	pred_tin_nhắn_đã_gửi := smry.Predicate{
		Label: "Tin nhắn đã gửi",
		Spec:  "tin_nhan_da_gui",
		Expr:  sq.NewExpr("deleted_at is null"),
	}

	row_tổng := smry.NewSubject("Tổng khách tin nhắn đã gửi", "", "count", "COUNT(*)", nil)

	rows := []smry.Subject{
		row_tổng.Combine("Tin nhắn đã gửi", pred_tin_nhắn_đã_gửi),
	}

	var cols []smry.Predicator
	for _, userID := range userIDs {
		cols = append(cols, smry.Predicate{
			Label: fmt.Sprintf("user_id = %d, from(%s) - to(%s)", userID, dateFrom.Format("2006-01-02"), dateTo.Format("2006-01-02")),
			Spec:  fmt.Sprintf("user_id = %d, from(%s) - to(%s)", userID, dateFrom.Format("2006-01-02"), dateTo.Format("2006-01-02")),
			Expr:  sq.NewExpr("external_created_time >= ? and ? > external_created_time and created_by = ?", dateFrom, dateTo, userID),
		})
	}
	return smry.BuildTable(rows, cols, "Kết quả tổng quát", "fb_external_message", "datefrom-dateto", "total")
	return nil
}

func buildTableCommentsByStaffs(
	dateFrom, dateTo time.Time, userIDs []dot.ID,
) (tableMessages *smry.Table) {
	pred_tin_nhắn_đã_gửi := smry.Predicate{
		Label: "Comment đã gửi",
		Spec:  "comment_da_gui",
		Expr:  sq.NewExpr("deleted_at is null"),
	}

	row_tổng := smry.NewSubject("Tổng khách comment đã gửi", "", "count", "COUNT(*)", nil)

	rows := []smry.Subject{
		row_tổng.Combine("Comment đã gửi", pred_tin_nhắn_đã_gửi),
	}

	var cols []smry.Predicator
	for _, userID := range userIDs {
		cols = append(cols, smry.Predicate{
			Label: fmt.Sprintf("user_id = %d, from(%s) - to(%s)", userID, dateFrom.Format("2006-01-02"), dateTo.Format("2006-01-02")),
			Spec:  fmt.Sprintf("user_id = %d, from(%s) - to(%s)", userID, dateFrom.Format("2006-01-02"), dateTo.Format("2006-01-02")),
			Expr:  sq.NewExpr("external_created_time >= ? and ? > external_created_time and created_by = ?", dateFrom, dateTo, userID),
		})
	}
	return smry.BuildTable(rows, cols, "Kết quả tổng quát", "fb_external_comment", "datefrom-dateto", "total")
}

func buildTableFbUsersWereAdvisedByMessage(
	dateFrom, dateTo time.Time, userIDs []dot.ID,
) (tableFbUsers *smry.Table) {
	pred_khách_hàng_tư_vấn_qua_message := smry.Predicate{
		Label: "Khách hàng tư vấn qua message",
		Spec:  "deleted is null",
		Expr:  sq.NewExpr("deleted_at is null"),
	}

	row_tổng := smry.NewSubject("Tổng khách hàng tư vấn qua message", "", "count(distinct(external_conversation_id))", "COUNT(DISTINCT(external_conversation_id))", nil)

	rows := []smry.Subject{
		row_tổng.Combine("Khách hàng tư vấn qua message", pred_khách_hàng_tư_vấn_qua_message),
	}

	var cols []smry.Predicator
	for _, userID := range userIDs {
		cols = append(cols, smry.Predicate{
			Label: fmt.Sprintf("user_id = %d, from(%s) - to(%s)", userID, dateFrom.Format("2006-01-02"), dateTo.Format("2006-01-02")),
			Spec:  fmt.Sprintf("user_id = %d, from(%s) - to(%s)", userID, dateFrom.Format("2006-01-02"), dateTo.Format("2006-01-02")),
			Expr:  sq.NewExpr("external_created_time >= ? and ? > external_created_time and created_by = ?", dateFrom, dateTo, userID),
		})
	}

	return smry.BuildTable(rows, cols, "Kết quả tổng quát", "fb_external_message", "datefrom-dateto", "total")
}

func buildTableFbUsersWereAdvisedByComment(
	dateFrom, dateTo time.Time, userIDs []dot.ID,
) (tableFbUsers *smry.Table) {
	pred_khách_hàng_tư_vấn_qua_comment := smry.Predicate{
		Label: "Khách hàng tư vấn qua comment",
		Spec:  "deleted is null",
		Expr:  sq.NewExpr("deleted_at is null"),
	}

	row_tổng := smry.NewSubject("Tổng khách hàng tư vấn qua comment", "", "count(distinct(external_conversation_id))", "COUNT(DISTINCT(external_parent_user_id))", nil)

	rows := []smry.Subject{
		row_tổng.Combine("Khách hàng tư vấn qua ", pred_khách_hàng_tư_vấn_qua_comment),
	}

	var cols []smry.Predicator
	for _, userID := range userIDs {
		cols = append(cols, smry.Predicate{
			Label: fmt.Sprintf("user_id = %d, from(%s) - to(%s)", userID, dateFrom.Format("2006-01-02"), dateTo.Format("2006-01-02")),
			Spec:  fmt.Sprintf("user_id = %d, from(%s) - to(%s)", userID, dateFrom.Format("2006-01-02"), dateTo.Format("2006-01-02")),
			Expr:  sq.NewExpr("external_created_time >= ? and ? > external_created_time and created_by = ?", dateFrom, dateTo, userID),
		})
	}

	return smry.BuildTable(rows, cols, "Kết quả tổng quát", "fb_external_comment", "datefrom-dateto", "total")
}

func buildTableFbUsersWereAdvised(
	db *cmsql.Database, dateFrom, dateTo time.Time, userIDs []dot.ID,
) (*smry.Table, error) {
	pred_khách_hàng_đã_tư_vấn := smry.Predicate{
		Label: "Khách hàng đã tư vấn",
		Spec:  "deleted is null",
		Expr:  sq.NewExpr("fb_external_message.deleted_at is null and fb_external_comment.deleted_at is null"),
	}

	row_tổng := smry.NewSubject("Tổng khách hàng đã tư vấn", "", "count(khach_hang_da_tu_van)", "count(khach_hang_da_tu_van)", nil)

	rows := []smry.Subject{
		row_tổng.Combine("Khách hàng đã tư vấn", pred_khách_hàng_đã_tư_vấn),
	}

	var cols []smry.Predicator
	for _, userID := range userIDs {
		cols = append(cols, smry.Predicate{
			Label: fmt.Sprintf("user_id = %d, from(%s) - to(%s)", userID, dateFrom.Format("2006-01-02"), dateTo.Format("2006-01-02")),
			Spec:  fmt.Sprintf("user_id = %d, from(%s) - to(%s)", userID, dateFrom.Format("2006-01-02"), dateTo.Format("2006-01-02")),
			Expr:  sq.NewExpr("user_id = ?, from(?) - to(?)", userID, dateFrom.Format("2006-01-02"), dateTo.Format("2006-01-02")),
		})
	}

	table := smry.BuildTable(rows, cols, "Kết quả tổng quát", "khach_hang_da_tu_van", "datefrom-dateto", "total")

	mapCreatedByAndMapFbUserID := make(map[dot.ID]map[dot.ID]bool)
	{
		dbRows, err := db.
			SQL("SELECT DISTINCT fec.external_parent_user_id, fec.created_by "+
				"FROM fb_external_comment fec ").
			Where("fec.external_created_time >= ? AND fec.external_created_time < ? ", dateFrom, dateTo).
			In("fec.created_by", userIDs).
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
			Where("fem.external_created_time >= ? AND fem.external_created_time < ?", dateFrom, dateTo).
			In("fem.created_by", userIDs).
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
			table.Data[i].Value = int64(len(mapFbUserID))
		}
	}

	return table, nil
}

func getToday(now time.Time) (from, to time.Time) {
	dateFrom, dateTo := now, now
	dateFrom = dateFrom.Add(-time.Duration(dateFrom.Second()) * time.Second)
	dateFrom = dateFrom.Add(-time.Duration(dateFrom.Minute()) * time.Minute)
	dateFrom = dateFrom.Add(-time.Duration(dateFrom.Hour()) * time.Hour)

	return dateFrom, dateTo
}

func convertIDsToInt64s(ids []dot.ID) []int64 {
	res := make([]int64, 0, len(ids))
	for _, id := range ids {
		res = append(res, id.Int64())
	}
	return res
}
