package sqlstore

import (
	"context"
	"fmt"
	"time"

	"o.o/api/etelecom/call_log_direction"
	"o.o/api/etelecom/summary"
	commonsummary "o.o/api/summary"
	"o.o/backend/com/etelecom/model"
	"o.o/backend/com/summary/util"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sq/core"
	smry "o.o/backend/pkg/etop/logic/summary"
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

func (s SummaryStore) Summary(
	ctx context.Context, req *summary.SummaryArgs) ([]*commonsummary.SummaryTable, error) {
	dateFrom, dateTo, shopID := req.DateFrom, req.DateTo, req.ShopID
	_, _, _ = dateFrom, dateTo, shopID

	var tables []*smry.Table

	tableCallLogTotal := buildTableCallLogTotal(dateFrom, dateTo)
	if err := s.execQuery(ctx, []*smry.Table{tableCallLogTotal}, shopID, "call_log"); err != nil {
		return nil, err
	}
	tables = append(tables, tableCallLogTotal)

	extensionIDs, err := s.getExtensionIDs(shopID)
	if err != nil {
		return nil, err
	}
	if len(extensionIDs) != 0 {
		tableCallLogPerExtension := buildTableCallLogPerExtension(dateFrom, dateTo, extensionIDs)
		if err := s.execQuery(ctx, []*smry.Table{tableCallLogPerExtension}, shopID, "call_log"); err != nil {
			return nil, err
		}
		tables = append(tables, tableCallLogPerExtension)
	}

	return util.BuildResponse(tables), nil
}

func (s SummaryStore) execQuery(ctx context.Context, tables []*smry.Table, accountID dot.ID, tableName string) error {
	for _, table := range tables {
		builder := smry.NewSummaryQueryBuilder(tableName)
		for i := range table.Data {
			// must always use [i] because we want to take the address
			builder.AddCell(&table.Data[i].Subject, (*core.Int64)(&table.Data[i].Value))
		}
		q := s.db.SQL(builder).WithContext(ctx).Where("account_id = ?", accountID).Clone()
		if err := q.Scan(builder.ScanArgs...); err != nil {
			return err
		}
	}
	return nil
}

func (s SummaryStore) getExtensionIDs(shopID dot.ID) ([]dot.ID, error) {
	var extensions model.Extensions

	err := s.db.
		Where("account_id = ?", shopID).
		Find(&extensions)
	if err != nil {
		return nil, err
	}

	var extensionIDs []dot.ID
	for _, extension := range extensions {
		extensionIDs = append(extensionIDs, extension.ID)
	}

	return extensionIDs, nil
}

func buildTableCallLogTotal(dateFrom, dateTo time.Time) *smry.Table {
	pred_tổng_số_cuộc_gọi := smry.Predicate{
		Label: "Tổng số cuộc gọi",
		Spec:  "Tổng số cuộc gọi",
		Expr:  sq.NewExpr("true"),
	}

	pred_tổng_thời_lượng_gọi := smry.Predicate{
		Label: "Tổng thời lượng gọi",
		Spec:  "Tổng thời lượng gọi",
		Expr:  sq.NewExpr("duration != 0 and direction = ?", call_log_direction.Out),
	}

	pred_tổng_thời_lượng_gọi_làm_tròn := smry.Predicate{
		Label: "Tổng thời lượng gọi làm tròn",
		Spec:  "Tổng thời lượng gọi làm tròn",
		Expr:  sq.NewExpr("duration != 0 and direction = ?", call_log_direction.Out),
	}

	pred_tổng_thời_lượng_nghe := smry.Predicate{
		Label: "Tổng thời lượng nghe",
		Spec:  "Tổng thời lượng nghe",
		Expr:  sq.NewExpr("duration != 0 and direction = ?", call_log_direction.In),
	}

	pred_tổng_số_cuộc_gọi_thành_công := smry.Predicate{
		Label: "Tổng số cuộc gọi thành công",
		Spec:  "call_status = 1",
		Expr:  sq.NewExpr("call_status = 1"),
	}

	pred_tổng_số_cuộc_gọi_thất_bại := smry.Predicate{
		Label: "Tổng số cuộc gọi thất bại",
		Spec:  "call_status != 1",
		Expr:  sq.NewExpr("call_status != 1"),
	}

	row_tổng_số_cuộc_gọi := smry.NewSubject("Tổng số cuộc gọi", "", "count", "COUNT(*)", nil)
	row_tổng_thời_lượng_làm_tròn := smry.NewSubject("Tổng thời lượng làm tròn", "", "sum(duration_postage)", "SUM(duration_postage)", nil)
	row_tổng_thời_lượng := smry.NewSubject("Tổng thời lượng", "", "sum(duration)", "SUM(duration)", nil)
	row_tổng_phí_gọi := smry.NewSubject("Tổng phí gọi", "", "sum(postage)", "SUM(postage)", nil)

	rowsTotal := []smry.Subject{
		row_tổng_số_cuộc_gọi.Combine("Tổng số cuộc gọi", pred_tổng_số_cuộc_gọi),
		row_tổng_thời_lượng.Combine("Tổng thời lượng gọi", pred_tổng_thời_lượng_gọi),
		row_tổng_thời_lượng_làm_tròn.Combine("Tổng thời lượng gọi làm tròn", pred_tổng_thời_lượng_gọi_làm_tròn),
		row_tổng_thời_lượng.Combine("Tổng thời lượng nghe", pred_tổng_thời_lượng_nghe),
		row_tổng_phí_gọi.Combine("Tổng phí gọi", pred_tổng_số_cuộc_gọi_thành_công),
		row_tổng_số_cuộc_gọi.Combine("Tổng số cuộc gọi thành công", pred_tổng_số_cuộc_gọi_thành_công),
		row_tổng_số_cuộc_gọi.Combine("Tổng số cuộc gọi thất baị", pred_tổng_số_cuộc_gọi_thất_bại),
	}

	colsTotal := []smry.Predicator{
		smry.Predicate{
			Label: "Tổng quát",
			Spec:  "total",
			Expr:  sq.NewExpr("created_at >= ? and ? > created_at", dateFrom, dateTo),
		},
	}

	return smry.BuildTable(rowsTotal, colsTotal, "Kết quả tổng quát", "call_log", "datefrom-dateto", "total")
}

func buildTableCallLogPerExtension(dateFrom, dateTo time.Time, extensionIDs []dot.ID) *smry.Table {
	pred_tổng_số_cuộc_gọi := smry.Predicate{
		Label: "Số cuộc gọi",
		Spec:  "Số cuộc gọi",
		Expr:  sq.NewExpr("true"),
	}

	pred_tổng_thời_lượng_gọi := smry.Predicate{
		Label: "Thời lượng gọi",
		Spec:  "Thời lượng gọi",
		Expr:  sq.NewExpr("duration != 0 and direction = ?", call_log_direction.Out),
	}

	pred_tổng_thời_lượng_gọi_làm_tròn := smry.Predicate{
		Label: "Tổng thời lượng gọi làm tròn",
		Spec:  "Tổng thời lượng gọi làm tròn",
		Expr:  sq.NewExpr("duration != 0 and direction = ?", call_log_direction.Out),
	}

	pred_tổng_thời_lượng_nghe := smry.Predicate{
		Label: "Thời lượng nghe",
		Spec:  "Thời lượng nghe",
		Expr:  sq.NewExpr("duration != 0 and direction = ?", call_log_direction.In),
	}

	pred_tổng_số_cuộc_gọi_thành_công := smry.Predicate{
		Label: "ExtensionID",
		Spec:  "call_status = 1",
		Expr:  sq.NewExpr("call_status = 1"),
	}

	pred_tổng_số_cuộc_gọi_thất_bại := smry.Predicate{
		Label: "Tổng số cuộc gọi thất bại",
		Spec:  "call_status != 1",
		Expr:  sq.NewExpr("call_status != 1"),
	}

	row_tổng_số_cuộc_gọi := smry.NewSubject("Tổng số cuộc gọi", "", "count", "COUNT(*)", nil)
	row_tổng_thời_lượng_làm_tròn := smry.NewSubject("Tổng thời lượng làm tròn", "", "sum(duration_postage)", "SUM(duration_postage)", nil)
	row_tổng_thời_lượng := smry.NewSubject("Tổng thời lượng", "", "sum(duration)", "SUM(duration)", nil)
	row_tổng_phí_gọi := smry.NewSubject("Tổng phí gọi", "", "sum(postage)", "SUM(postage)", nil)

	rowsTotal := []smry.Subject{
		row_tổng_số_cuộc_gọi.Combine("Tổng số cuộc gọi", pred_tổng_số_cuộc_gọi),
		row_tổng_thời_lượng.Combine("Tổng thời lượng gọi", pred_tổng_thời_lượng_gọi),
		row_tổng_thời_lượng_làm_tròn.Combine("Tổng thời lượng gọi làm tròn", pred_tổng_thời_lượng_gọi_làm_tròn),
		row_tổng_thời_lượng.Combine("Tổng thời lượng", pred_tổng_thời_lượng_nghe),
		row_tổng_phí_gọi.Combine("Tổng phí gọi", pred_tổng_số_cuộc_gọi_thành_công),
		row_tổng_số_cuộc_gọi.Combine("Tổng số cuộc gọi thành công", pred_tổng_số_cuộc_gọi_thành_công),
		row_tổng_số_cuộc_gọi.Combine("Tổng số cuộc gọi thất baị", pred_tổng_số_cuộc_gọi_thất_bại),
	}

	var colsPerExtension []smry.Predicator

	for _, extensionID := range extensionIDs {
		colsPerExtension = append(colsPerExtension, smry.Predicate{
			Label: fmt.Sprintf("extensionID(%d)", extensionID),
			Spec:  fmt.Sprintf("extension_id = %d", extensionID),
			Expr:  sq.NewExpr("extension_id = ? and created_at >= ? and created_at < ?", extensionID, dateFrom, dateTo),
		})
	}

	return smry.BuildTable(rowsTotal, colsPerExtension, "Kết quả theo nhân viên", "call_log", "datefrom-dateto", "per-extension")
}
