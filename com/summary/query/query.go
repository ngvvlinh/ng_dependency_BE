package query

import (
	"context"
	"strconv"
	"time"

	"etop.vn/api/summary"
	"etop.vn/backend/com/summary/sqlstore"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/redis"
)

const currentVersion = "1"
const dateLayout = "2016-01-02"

var _ summary.QueryService = &DashboardQuery{}

type DashboardQuery struct {
	store      sqlstore.OrderStoreFactory
	redisStore redis.Store
}

func NewDashboardQuery(db cmsql.Database, resdis redis.Store) *DashboardQuery {
	return &DashboardQuery{
		store:      sqlstore.NewOrderStore(db),
		redisStore: resdis,
	}
}

func (q *DashboardQuery) MessageBus() summary.QueryBus {
	b := bus.New()
	return summary.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func buildKey(shopID int64, dateFrom, dateTo time.Time) string {
	key := "summary/pos:version=" + currentVersion +
		",shop=" + strconv.Itoa(int(shopID)) +
		",from=" + dateFrom.Format(dateLayout) +
		",to=" + dateTo.Format(dateLayout)
	return key
}

func (q *DashboardQuery) SummaryPOS(ctx context.Context, req *summary.SummaryPOSRequest) (*summary.SummaryPOSResponse, error) {
	var summaryTableRedis summary.SummaryPOSResponse
	keyRedis := buildKey(req.ShopID, req.DateFrom, req.DateTo)
	err := q.redisStore.Get(keyRedis, &summaryTableRedis)
	if err != nil && err != redis.ErrNil {
		return nil, err
	}
	if err == nil {
		return &summaryTableRedis, nil
	}

	var summaryTable summary.SummaryPOSResponse
	resultTotal, err := q.store(ctx).ShopID(req.ShopID).GetOrderSummary(req.DateFrom, req.DateTo)
	if err != nil {
		return nil, err
	}
	summaryTable.ListTable = append(summaryTable.ListTable, BuildTotalTable(resultTotal))
	resultTotolPerDay, err := q.store(ctx).ShopID(req.ShopID).GetAmoumtPerDay(req.DateFrom, req.DateTo)
	if err != nil {
		return nil, err
	}
	summaryTable.ListTable = append(summaryTable.ListTable, BuildDiagramTable(resultTotolPerDay))
	resultTopSellItem, err := q.store(ctx).GetTopSellItem(req.ShopID, req.DateFrom, req.DateTo)
	if err != nil {
		return nil, err
	}
	summaryTable.ListTable = append(summaryTable.ListTable, BuildTopSellTable(resultTopSellItem))

	err = q.redisStore.SetWithTTL(keyRedis, summaryTable, 30)
	if err != nil && err != redis.ErrNil {
		return nil, err
	}
	return &summaryTable, nil
}

func BuildTotalTable(args *sqlstore.Total) *summary.SummaryTable {
	var summaryTable summary.SummaryTable
	summaryTable.Label = "Tổng hợp thống kê"
	summaryTable.Cols = append(summaryTable.Cols, summary.SummaryColRow{
		Spec:  "sum(total_count)",
		Label: "Tổng doanh thu",
	})
	summaryTable.Cols = append(summaryTable.Cols, summary.SummaryColRow{
		Spec:  "count(order_id)",
		Label: "Số lượng đơn hàng",
	})
	summaryTable.Cols = append(summaryTable.Cols, summary.SummaryColRow{
		Spec:  "avg(total_count)",
		Label: "Giá trị trung bình đơn hàng",
	})

	summaryTable.Rows = append(summaryTable.Rows, summary.SummaryColRow{
		Spec:  "value",
		Label: "Value",
	})

	summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
		Spec:  "sum(total_count):value",
		Value: int(args.TotalAmount.Int64),
	})
	summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
		Spec:  "count(order_id):value",
		Value: int(args.TotalOrder.Int64),
	})
	summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
		Spec:  "avg(total_count):value",
		Value: int(args.AverageOrder.Float64),
	})

	return &summaryTable
}

func BuildDiagramTable(args []*sqlstore.TotalPerDate) *summary.SummaryTable {
	var summaryTable summary.SummaryTable
	summaryTable.Label = "Thống kê"
	summaryTable.Rows = append(summaryTable.Rows, summary.SummaryColRow{
		Spec:  "sum(total_count)&&perday=true",
		Label: "Doanh Thu",
	})
	summaryTable.Rows = append(summaryTable.Rows, summary.SummaryColRow{
		Spec:  "count(order_id)&&perday=true",
		Label: "Số Lượng Đơn",
	})

	for _, value := range args {
		summaryTable.Cols = append(summaryTable.Cols, summary.SummaryColRow{
			Spec:  "date(" + value.Day.Format("2006-01-02") + ")",
			Label: value.Day.Format("2006-01-02"),
		})
	}

	for index, value := range args {
		summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
			Spec:  summaryTable.Cols[index].Spec + ":" + summaryTable.Rows[0].Spec,
			Value: int(value.TotalAmount),
			Unit:  "",
		})
	}
	for index, value := range args {
		summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
			Spec:  summaryTable.Cols[index].Spec + ":" + summaryTable.Rows[1].Spec,
			Value: int(value.Count),
			Unit:  "",
		})
	}

	return &summaryTable
}

func BuildTopSellTable(args []*sqlstore.TopSellItem) *summary.SummaryTable {
	var summaryTable summary.SummaryTable
	summaryTable.Label = "Sản phẩm bán chạy"
	summaryTable.Cols = append(summaryTable.Cols, summary.SummaryColRow{
		Spec:  "img_urls",
		Label: "hình ảnh sản phẩm",
	})
	summaryTable.Cols = append(summaryTable.Cols, summary.SummaryColRow{
		Spec:  "name_product",
		Label: "Tên Sản Phẩm",
	})
	summaryTable.Cols = append(summaryTable.Cols, summary.SummaryColRow{
		Spec:  "sum(quantity)",
		Label: "Số Lượng Sản Phẩm",
	})
	for index, _ := range args {
		summaryTable.Rows = append(summaryTable.Rows, summary.SummaryColRow{
			Spec:  "STT-" + strconv.Itoa(index),
			Label: strconv.Itoa(index),
		})
	}

	for index, value := range args {
		summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
			Spec:      summaryTable.Cols[0].Spec + ":" + summaryTable.Rows[index].Spec,
			ImageUrls: value.ImageUrls,
		})
		summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
			Spec:  summaryTable.Cols[1].Spec + ":" + summaryTable.Rows[index].Spec,
			Label: value.Name,
		})
		summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
			Spec:  summaryTable.Cols[2].Spec + ":" + summaryTable.Rows[index].Spec,
			Value: int(value.Count),
			Unit:  "",
		})
	}

	return &summaryTable
}
