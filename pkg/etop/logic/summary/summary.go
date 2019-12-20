package summary

import (
	"context"
	"fmt"
	"math"
	"time"

	txmodel "etop.vn/backend/com/main/moneytx/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/idemp"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/pkg/common/sql/sq/core"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
	"etop.vn/common/l"
)

var (
	x  *cmsql.Database
	ll = l.New()

	idempgroup = idemp.NewGroup()
)

func init() {
	bus.AddHandlers("sql", SummarizeFulfillments)
}

func Init(db *cmsql.Database) {
	if x != nil && (*x).DB() != nil {
		ll.Panic("Already initialized")
	}
	x = db
}

func SummarizeFulfillments(ctx context.Context, query *model.SummarizeFulfillmentsRequest) error {
	var timeout time.Duration
	if cm.IsProd() {
		timeout = 5 * time.Second
	}

	key := fmt.Sprintf("SummarizeFulfillments %v", query.ShopID)
	resp, err, _ := idempgroup.Do(key, timeout, func() (interface{}, error) {
		return summarizeFulfillments(ctx, query)
	})
	query.Result = resp.(*model.SummarizeFulfillmentsRequest).Result
	return err
}

func summarizeFulfillments(ctx context.Context, query *model.SummarizeFulfillmentsRequest) (*model.SummarizeFulfillmentsRequest, error) {
	if query.ShopID == 0 {
		return query, cm.Error(cm.InvalidArgument, "missing shop_id", nil)
	}
	from, to, err := cm.ParseDateFromTo(query.DateFrom, query.DateTo)
	if err != nil {
		return nil, err
	}

	var moneyTransactions txmodel.MoneyTransactionShippings
	if err := x.Table("money_transaction_shipping").
		Where("shop_id = ?", query.ShopID).
		Where("status = 0").
		Where("total_orders > 0").
		Find(&moneyTransactions); err != nil {
		return query, err
	}

	moneyTransactionIDs := make([]dot.ID, len(moneyTransactions))
	for i, mt := range moneyTransactions {
		moneyTransactionIDs[i] = mt.ID
	}

	includeNewTables := !from.IsZero() && !to.IsZero()
	tables := buildTables(moneyTransactionIDs)
	if includeNewTables {
		tables = append(tables, buildTables2(from, to)...)
	}
	if err := execQuery(ctx, tables, query.ShopID); err != nil {
		return query, cm.Error(cm.Internal, "can not execute query", err)
	}
	if includeNewTables {
		tables = append(tables, buildTableAverage(tables[len(tables)-1]))
	}
	query.Result.Tables = buildResponse(tables)
	return query, nil
}

func execQuery(ctx context.Context, tables []*Table, shopID dot.ID) error {
	builder := NewSummaryQueryBuilder("fulfillment")
	for _, table := range tables {
		for i := range table.Data {
			// must always use [i] because we want to take the address
			builder.AddCell(&table.Data[i].Subject, (*core.Int)(&table.Data[i].Value))
		}
	}
	return x.SQL(builder).WithContext(ctx).
		Where("shop_id = ?", shopID).Scan(builder.ScanArgs...)
}

func buildResponse(tables []*Table) []*model.SummaryTable {
	res := make([]*model.SummaryTable, len(tables))
	for i, table := range tables {
		stable := &model.SummaryTable{
			Label: table.Label,
			Tags:  table.Tags,
			Cols:  buildCols(table.Cols),
			Rows:  buildRows(table.Rows),
			Data:  buildData(table.Data),
		}
		res[i] = stable
	}
	return res
}

func buildCols(cols []Predicator) []model.SummaryColRow {
	res := make([]model.SummaryColRow, len(cols))
	for i, col := range cols {
		res[i] = model.SummaryColRow{
			Label:  col.GetLabel(),
			Spec:   col.GetSpec(),
			Unit:   "",
			Indent: 0,
		}
	}
	return res
}

func buildRows(rows []Subject) []model.SummaryColRow {
	res := make([]model.SummaryColRow, len(rows))
	for i, row := range rows {
		res[i] = model.SummaryColRow{
			Label:  row.GetLabel(),
			Spec:   row.GetSpec(),
			Unit:   row.Unit,
			Indent: row.Ident,
		}
	}
	return res
}

func buildData(data []Cell) []model.SummaryItem {
	res := make([]model.SummaryItem, len(data))
	for i, item := range data {
		res[i] = model.SummaryItem{
			Spec:  item.Subject.GetSpec(),
			Value: item.Value,
			Unit:  item.Subject.Unit,
		}
	}
	return res
}

func buildTableAverage(input *Table) *Table {
	if input.NRow != 16 {
		panic("invalid input table")
	}

	nrows := 4
	table := NewTable(nrows, input.NCol, "Nội dung các giá trị trung bình", "average_created_at_after_cross_check", "after_cross_check", "average", "created_at")
	table.Cols = input.Cols
	table.Rows = make([]Subject, nrows)
	for col := 0; col < input.NCol; col++ {
		setRow(table, 0, col, calcDiv(input, 5, 0, col),
			"Trung bình phí đơn giao hàng thành công [Tổng (Phí giao + Trả hàng) / Tổng số lượng đơn giao hàng thành công]",
			"₫", "sum(shipping_fee_shop)/(count:shipping_status=P)")

		setRow(table, 1, col, calcDiv(input, 5, 1, col),
			"Trung bình phí tất cả đơn [Tổng phí giao hàng / Tổng số lượng đơn giao hàng]",
			"₫", "sum(shipping_fee_shop)/count")

		setRow(table, 2, col, calcDiv(input, 3, 0, col),
			"Giá trị trung bình đơn giao hàng thành công [Tổng giá trị đơn giao hàng thành công / Tổng số lượng đơn giao hàng thành công]",
			"₫", "sum(shipping_fee_shop)/(count:shipping_status=P)")

		setRow(table, 3, col, calcDiv(input, 3, 1, col),
			"Giá trị trung bình tất cả đơn [Tổng giá trị đơn giao hàng đã tạo / Tổng số lượng đơn giao hàng]",
			"₫", "sum(shipping_fee_shop)/count")
	}
	return table
}

func setRow(table *Table, row, col int, value int, label, unit, spec string) {
	subj := Subject{
		Label: label,
		Unit:  unit,
		Spec:  spec,
	}
	table.Rows[row] = subj

	cell := table.Cell(row, col)
	cell.Subject = subj.Combine("", table.Cols[col])
	cell.Value = value
}

func calcDiv(input *Table, above, below, col int) int {
	n := input.Cell(below, col).Value
	if n == 0 {
		return 0
	}
	return int(math.Trunc(float64(input.Cell(above, col).Value) / float64(n)))
}
