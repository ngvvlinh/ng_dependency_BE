package query

import (
	"context"
	"strconv"
	"time"

	"etop.vn/api/main/location"
	"etop.vn/api/summary"
	"etop.vn/backend/com/summary/sqlstore"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/redis"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/capi/dot"
)

const currentVersion = "2"
const dateLayout = "2016-01-02"
const redisTime = 2 * 60

var _ summary.QueryService = &DashboardQuery{}

func (q *DashboardQuery) MessageBus() summary.QueryBus {
	b := bus.New()
	return summary.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (d *DashboardQuery) SummaryTopShip(ctx context.Context, req *summary.SummaryTopShipRequest) (*summary.SummaryTopShipResponse, error) {
	var summaryTableRedis summary.SummaryTopShipResponse
	keyRedis := buildKey(req.ShopID, req.DateFrom, req.DateTo, "topship")
	isReturn, err := d.checkRedis(keyRedis, &summaryTableRedis)
	if isReturn {
		return &summaryTableRedis, err
	}

	tables, err := d.store.SummarizeTopship(ctx, req)
	var totalAmountFFm *summary.SummaryTable
	var moneyTransaction *summary.SummaryTable
	var keyTotalAmountFFm int
	var keyMoneyTransaction int
	for key, value := range tables {
		if value.Tags[0] == "fulfillment06" {
			keyTotalAmountFFm = key
			totalAmountFFm = value
		}
		if value.Tags[0] == "fulfillment07" {
			keyMoneyTransaction = key
			moneyTransaction = value
		}
	}
	if tables[keyTotalAmountFFm] != nil && tables[keyMoneyTransaction] != nil {
		tables[keyTotalAmountFFm], tables[keyMoneyTransaction] = buildTotalMoneyTransaction(totalAmountFFm, moneyTransaction)
	}

	resulFfmByArea, err := d.factory(ctx).GetCountFfmByArea(req.ShopID, req.DateFrom, req.DateTo)
	if err != nil {
		return nil, err
	}
	summaryTableFfmByArea, err := d.buildFfmByAreaTable(ctx, resulFfmByArea)
	if err != nil {
		return nil, err
	}
	tables = append(tables, summaryTableFfmByArea)

	result := &summary.SummaryTopShipResponse{
		ListTable: tables,
	}
	err = d.redisStore.SetWithTTL(keyRedis, result, redisTime)
	if err != nil && err != redis.ErrNil {
		return nil, err
	}

	return result, err
}

type DashboardQuery struct {
	store       sqlstore.SummaryStore
	factory     sqlstore.SummaryStoreFactory
	redisStore  redis.Store
	locationBus location.QueryBus
}

func NewDashboardQuery(db *cmsql.Database, resdis redis.Store, locationBus location.QueryBus) *DashboardQuery {
	return &DashboardQuery{
		store:       sqlstore.NewSummaryStore(db),
		factory:     sqlstore.NewSummaryStoreFactory(db),
		redisStore:  resdis,
		locationBus: locationBus,
	}
}

func buildKey(shopID dot.ID, dateFrom, dateTo time.Time, keycode string) string {
	key := "summary/" + keycode + ":version=" + currentVersion +
		",shop=" + shopID.String() +
		",from=" + dateFrom.Format(dateLayout) +
		",to=" + dateTo.Format(dateLayout)
	return key
}

func (q *DashboardQuery) SummaryPOS(ctx context.Context, req *summary.SummaryPOSRequest) (*summary.SummaryPOSResponse, error) {
	var summaryTableRedis summary.SummaryPOSResponse
	keyRedis := buildKey(req.ShopID, req.DateFrom, req.DateTo, "pos")
	isReturn, err := q.checkRedis(keyRedis, &summaryTableRedis)
	if isReturn {
		return &summaryTableRedis, err
	}

	startDayTime := getStartDayTime()
	resultSummaryOrderToday, err := q.factory(ctx).GetOrderSummary(req.ShopID, startDayTime, startDayTime.Add(24*time.Hour))
	if err != nil {
		return nil, err
	}
	resultSummaryOrderYesterday, err := q.factory(ctx).GetOrderSummary(req.ShopID, startDayTime.Add(-24*time.Hour), startDayTime)
	if err != nil {
		return nil, err
	}
	// Lay cach 30 ngay
	timeLastMonthStart := getStartDayTime().Add(-30 * 24 * time.Hour)
	timeLastMonthEnd := timeLastMonthStart.Add(24 * time.Hour)
	resultSummaryADayLastMonth, err := q.factory(ctx).GetOrderSummary(req.ShopID, timeLastMonthStart, timeLastMonthEnd)
	if err != nil {
		return nil, err
	}
	summaryTableRedis.ListTable = append(summaryTableRedis.ListTable, buildSummaryToday(resultSummaryOrderToday, resultSummaryOrderYesterday, resultSummaryADayLastMonth))

	resultTotal, err := q.factory(ctx).GetOrderSummary(req.ShopID, req.DateFrom, req.DateTo)
	if err != nil {
		return nil, err
	}
	summaryTableRedis.ListTable = append(summaryTableRedis.ListTable, buildTotalTable(resultTotal))

	resultTotalPerDay, err := q.factory(ctx).GetAmoumtPerDay(req.ShopID, req.DateFrom, req.DateTo)
	if err != nil {
		return nil, err
	}
	summaryTableRedis.ListTable = append(summaryTableRedis.ListTable, buildDiagramOrderTable(resultTotalPerDay, req.DateFrom, req.DateTo))

	resultTopSellItem, err := q.factory(ctx).GetTopSellItem(req.ShopID, req.DateFrom, req.DateTo)
	if err != nil {
		return nil, err
	}
	summaryTableRedis.ListTable = append(summaryTableRedis.ListTable, buildTopSellTable(resultTopSellItem))

	resultStaffOrder, err := q.factory(ctx).GetListStaffOrder(req.ShopID, req.DateFrom, req.DateTo)
	if err != nil {
		return nil, err
	}
	summaryTableRedis.ListTable = append(summaryTableRedis.ListTable, buildStaffOrderTable(resultStaffOrder))
	summaryTableRedis.ListTable = append(summaryTableRedis.ListTable, buildOtherSourceTable(resultStaffOrder))

	err = q.redisStore.SetWithTTL(keyRedis, summaryTableRedis, redisTime)
	if err != nil && err != redis.ErrNil {
		return nil, err
	}
	return &summaryTableRedis, nil
}

func buildTotalMoneyTransaction(totalAmountFfm *summary.SummaryTable, moneyTransaction *summary.SummaryTable) (*summary.SummaryTable, *summary.SummaryTable) {
	var summaryTableTotal summary.SummaryTable
	summaryTableTotal.Tags = totalAmountFfm.Tags
	summaryTableTotal.Rows = totalAmountFfm.Rows
	summaryTableTotal.Cols = totalAmountFfm.Cols

	summaryTableTotal.Cols = append(summaryTableTotal.Cols, summary.SummaryColRow{
		Spec:  "sum(cod,actual_comppensation_amount)",
		Label: "Tổng cộng",
	})

	summaryTableTotal.Rows = append(summaryTableTotal.Rows, summary.SummaryColRow{
		Spec:  "sum(cod,actual_comppensation_amount)-sum(shipping_fee)",
		Label: "Tiền trả shop",
	})

	for i := 0; i < len(totalAmountFfm.Data); i++ {
		if i == 0 {
			totalCod := totalAmountFfm.Data[i]
			totalCod.Value = totalCod.Value + totalAmountFfm.Data[4].Value
			summaryTableTotal.Data = append(summaryTableTotal.Data, totalCod)
			continue
		}
		summaryTableTotal.Data = append(summaryTableTotal.Data, totalAmountFfm.Data[i])
		if i%2 == 1 {
			// nếu là phần tử cuối cùng thì tính môt ô tổng
			summaryTableTotal.Data = append(summaryTableTotal.Data, summary.SummaryItem{
				Label: totalAmountFfm.Data[i-1].Label + "summary",
				Value: totalAmountFfm.Data[i].Value + totalAmountFfm.Data[i-1].Value,
			})
		}
	}

	summaryTableTotal.Data = append(summaryTableTotal.Data, summary.SummaryItem{
		Label: "sum(cod,actual_comppensation_amount)-sum(shipping_fee):delivered|undeliverable",
		Value: summaryTableTotal.Data[0].Value - summaryTableTotal.Data[3].Value,
	})

	summaryTableTotal.Data = append(summaryTableTotal.Data, summary.SummaryItem{
		Label: "-sum(shipping_fee):returning|returned",
		Value: 0 - summaryTableTotal.Data[4].Value,
	})
	summaryTableTotal.Data = append(summaryTableTotal.Data, summary.SummaryItem{
		Label: "sum(cod,actual_comppensation_amount)-sum(shipping_fee)&delivered|undeliverable|returning|returned",
		Value: summaryTableTotal.Data[2].Value - summaryTableTotal.Data[5].Value,
	})

	// money transaction
	var summaryTableMoneyTransaction summary.SummaryTable
	summaryTableMoneyTransaction.Tags = moneyTransaction.Tags
	summaryTableMoneyTransaction.Rows = totalAmountFfm.Rows
	summaryTableMoneyTransaction.Cols = totalAmountFfm.Rows

	summaryTableMoneyTransaction.Rows = append(summaryTableMoneyTransaction.Rows, summary.SummaryColRow{
		Spec:  "sum(cod,actual_comppensation_amount)-sum(shipping_fee)",
		Label: "Tiền trả shop",
	})

	for i := 0; i < len(moneyTransaction.Data); i++ {
		if i < 3 {
			// tất cả totol_cod = total_cd(delivered) +
			totalCod := moneyTransaction.Data[i]
			totalCod.Value = totalCod.Value + moneyTransaction.Data[i+6].Value
			summaryTableMoneyTransaction.Data = append(summaryTableMoneyTransaction.Data, totalCod)
			continue
		}
		summaryTableMoneyTransaction.Data = append(summaryTableMoneyTransaction.Data, moneyTransaction.Data[i])
	}
	summaryTableMoneyTransaction.Data = append(summaryTableMoneyTransaction.Data, summary.SummaryItem{
		Label: "sum(cod,actual_comppensation_amount)-sum(shipping_fee):money_transaction_id!=nil,cod_etop_transfered_at!=nil",
		Value: summaryTableMoneyTransaction.Data[0].Value - summaryTableMoneyTransaction.Data[3].Value,
	})

	summaryTableMoneyTransaction.Data = append(summaryTableMoneyTransaction.Data, summary.SummaryItem{
		Label: "sum(cod,actual_comppensation_amount)-sum(shipping_fee):money_transaction_id!=nil,cod_etop_transfered_at=nil",
		Value: summaryTableMoneyTransaction.Data[1].Value - summaryTableMoneyTransaction.Data[4].Value,
	})
	summaryTableMoneyTransaction.Data = append(summaryTableMoneyTransaction.Data, summary.SummaryItem{
		Label: "sum(cod,actual_comppensation_amount)-sum(shipping_fee):money_transaction_id=nil,cod_etop_transfered_at=nil",
		Value: summaryTableMoneyTransaction.Data[2].Value - summaryTableMoneyTransaction.Data[5].Value,
	})

	return &summaryTableTotal, &summaryTableMoneyTransaction
}

func (d *DashboardQuery) buildFfmByAreaTable(ctx context.Context, args []*sqlstore.FfmByArea) (*summary.SummaryTable, error) {
	var summaryTable summary.SummaryTable
	summaryTable.Label = "Thống kê giao hàng theo khu vực"

	// cols
	summaryTable.Cols = append(summaryTable.Cols, summary.SummaryColRow{
		Spec:  "count(id):",
		Label: "Tổng số đơn",
	})
	summaryTable.Cols = append(summaryTable.Cols, summary.SummaryColRow{
		Spec:  "province_name:",
		Label: "Tên Tỉnh/Thành phố",
	})
	summaryTable.Cols = append(summaryTable.Cols, summary.SummaryColRow{
		Spec:  "district_name:",
		Label: "Tên Huyện/Quận",
	})
	summaryTable.Cols = append(summaryTable.Cols, summary.SummaryColRow{
		Spec:  "province_code:",
		Label: "Mã Tỉnh/Thành phố",
	})
	summaryTable.Cols = append(summaryTable.Cols, summary.SummaryColRow{
		Spec:  "district_code:",
		Label: "Mã Huyện/Quận",
	})

	key := -1
	for _, value := range args {
		if value.ProvinceCode == "0" || value.DistrictCode == "0" {
			continue
		}
		key++
		query := &location.GetLocationQuery{
			ProvinceCode: value.ProvinceCode,
			DistrictCode: value.DistrictCode,
			Result:       nil,
		}
		err := d.locationBus.Dispatch(ctx, query)
		if err != nil {
			return nil, err
		}
		summaryTable.Rows = append(summaryTable.Rows, summary.SummaryColRow{
			Spec:  "groupby(product_id),orderby(sum(quantity)),row_number(" + strconv.Itoa(key) + ")",
			Label: "value",
		})
		summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
			Spec:  summaryTable.Cols[0].Spec + summaryTable.Rows[key].Spec,
			Value: value.Count,
		})
		summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
			Spec:  summaryTable.Cols[1].Spec + summaryTable.Rows[key].Spec,
			Label: query.Result.Province.Name,
		})
		summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
			Spec:  summaryTable.Cols[2].Spec + summaryTable.Rows[key].Spec,
			Label: query.Result.District.Name,
		})
		summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
			Spec:  summaryTable.Cols[3].Spec + summaryTable.Rows[key].Spec,
			Label: value.ProvinceCode,
		})
		summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
			Spec:  summaryTable.Cols[4].Spec + summaryTable.Rows[key].Spec,
			Label: value.DistrictCode,
		})
	}
	return &summaryTable, nil
}

func getStartDayTime() time.Time {
	timeToday := time.Now()
	timeToday = timeToday.Add(-time.Duration(timeToday.Second()) * time.Second)
	timeToday = timeToday.Add(-time.Duration(timeToday.Hour()) * time.Hour)
	timeToday = timeToday.Add(-time.Duration(timeToday.Minute()) * time.Minute)
	return timeToday
}

func getLastMonthTime() (time.Time, time.Time) {
	timeNow := time.Now()
	lastDayLastMonth := timeNow.Add(-time.Duration(timeNow.Day()+1) * time.Hour * 24)
	if timeNow.Day() > lastDayLastMonth.Day() {
		return startOfDay(lastDayLastMonth), startOfDay(lastDayLastMonth).Add(24 * time.Hour)
	}
	dayOfLastMonth := lastDayLastMonth.Add((time.Duration(lastDayLastMonth.Day() - timeNow.Day())) * time.Hour * -24)
	return startOfDay(dayOfLastMonth), startOfDay(dayOfLastMonth).Add(24 * time.Hour)
}

func startOfDay(args time.Time) time.Time {
	args = args.Add(-time.Duration(args.Second()) * time.Second)
	args = args.Add(-time.Duration(args.Hour()) * time.Hour)
	args = args.Add(-time.Duration(args.Minute()) * time.Minute)
	return args
}

func (q *DashboardQuery) checkRedis(key string, obj interface{}) (bool, error) {
	err := q.redisStore.Get(key, obj)
	switch err {
	case nil:
		return true, nil
	case redis.ErrNil:
		return false, nil
	default:
		return true, err
	}
}

func buildSummaryToday(summaryToday *sqlstore.Total, summaryYesterday *sqlstore.Total, summaryDayLastMonth *sqlstore.Total) *summary.SummaryTable {
	var summaryTable summary.SummaryTable
	summaryTable.Label = "Tổng hợp thống kê ngày hôm nay"
	// Today
	summaryTable.Cols = append(summaryTable.Cols, summary.SummaryColRow{
		Spec:  "sum(total_amount):",
		Label: "Tổng doanh thu",
	})
	summaryTable.Cols = append(summaryTable.Cols, summary.SummaryColRow{
		Spec:  "count(order_id):",
		Label: "Số lượng đơn hàng",
	})
	summaryTable.Cols = append(summaryTable.Cols, summary.SummaryColRow{
		Spec:  "avg(total_amount):",
		Label: "Giá trị trung bình đơn hàng",
	})

	summaryTable.Rows = append(summaryTable.Rows, summary.SummaryColRow{
		Spec:  "today(created_at),returning|returned",
		Label: "Return Today",
	})
	summaryTable.Rows = append(summaryTable.Rows, summary.SummaryColRow{
		Spec:  "today(created_at)",
		Label: "Today",
	})
	summaryTable.Rows = append(summaryTable.Rows, summary.SummaryColRow{
		Spec:  "(today-1)(created_at)",
		Label: "Yesterday",
	})
	summaryTable.Rows = append(summaryTable.Rows, summary.SummaryColRow{
		Spec:  "(today-1month)(created_at)",
		Label: "Day Of Last Month",
	})

	// Return
	summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
		Spec:  "sum(total_amount):today(created_at),returning|returned",
		Value: 0,
	})
	summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
		Spec:  "count(order_id):today(created_at),returning|returned",
		Value: 0,
	})
	summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
		Spec:  "avg(total_amount):today(created_at),returning|returned",
		Value: 0,
	})

	// Today
	summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
		Spec:  "sum(total_amount):today(created_at)",
		Value: summaryToday.TotalAmount,
	})
	summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
		Spec:  "count(order_id):today(created_at)",
		Value: summaryToday.TotalOrder,
	})
	summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
		Spec:  "avg(total_amount):today(created_at)",
		Value: int(summaryToday.AverageOrder),
	})

	// Yesterday
	summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
		Spec:  "sum(total_amount):(today-1)(created_at)",
		Value: summaryYesterday.TotalAmount,
	})
	summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
		Spec:  "count(order_id):(today-1)(created_at)",
		Value: summaryYesterday.TotalOrder,
	})
	summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
		Spec:  "avg(total_amount):(today-1)(created_at)",
		Value: int(summaryYesterday.AverageOrder),
	})

	// Day of last month
	summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
		Spec:  "sum(total_amount):(today-1month)(created_at)",
		Value: summaryDayLastMonth.TotalAmount,
	})
	summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
		Spec:  "count(order_id):(today-1month)(created_at)",
		Value: summaryDayLastMonth.TotalOrder,
	})
	summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
		Spec:  "avg(total_amount):(today-1month)(created_at)",
		Value: int(summaryDayLastMonth.AverageOrder),
	})

	return &summaryTable
}

func buildTotalTable(args *sqlstore.Total) *summary.SummaryTable {
	var summaryTable summary.SummaryTable
	summaryTable.Label = "Tổng hợp thống kê theo thời gian"
	summaryTable.Cols = append(summaryTable.Cols, summary.SummaryColRow{
		Spec:  "sum(total_amount):",
		Label: "Tổng doanh thu",
	})
	summaryTable.Cols = append(summaryTable.Cols, summary.SummaryColRow{
		Spec:  "count(order_id):",
		Label: "Số lượng đơn hàng",
	})
	summaryTable.Cols = append(summaryTable.Cols, summary.SummaryColRow{
		Spec:  "avg(total_amount):",
		Label: "Giá trị trung bình đơn hàng",
	})

	summaryTable.Rows = append(summaryTable.Rows, summary.SummaryColRow{
		Spec:  "positive",
		Label: "Positive",
	})
	summaryTable.Rows = append(summaryTable.Rows, summary.SummaryColRow{
		Spec:  "negative",
		Label: "Negative",
	})

	// Positive
	summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
		Spec:  "sum(total_amount):positive",
		Value: args.TotalAmount,
	})
	summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
		Spec:  "count(order_id):positive",
		Value: args.TotalOrder,
	})
	summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
		Spec:  "avg(total_amount):positive",
		Value: int(args.AverageOrder),
	})

	// Negative
	summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
		Spec:  "sum(total_amount):Negative",
		Value: 0,
	})
	summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
		Spec:  "count(order_id):Negative",
		Value: 0,
	})
	summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
		Spec:  "avg(total_amount):Negative",
		Value: 0,
	})

	return &summaryTable
}

func buildDiagramOrderTable(args []*sqlstore.TotalPerDate, dateFrom time.Time, dateTo time.Time) *summary.SummaryTable {
	var summaryTable summary.SummaryTable
	summaryTable.Label = "Thống kê"
	summaryTable.Cols = append(summaryTable.Cols, summary.SummaryColRow{
		Spec:  "sum(total_count):",
		Label: "Doanh Thu",
	})
	summaryTable.Cols = append(summaryTable.Cols, summary.SummaryColRow{
		Spec:  "count(order_id):",
		Label: "Số Lượng Đơn",
	})

	for index, value := range args {
		summaryRow := summary.SummaryColRow{
			Spec:  "date(" + value.Day.Format("2006-01-02") + ")" + ",status!=-1",
			Label: value.Day.Format("2006-01-02"),
		}
		summaryTable.Rows = append(summaryTable.Rows, summaryRow)
		summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
			Label: "Tổng giá trị đơn ngày " + summaryRow.Spec,
			Spec:  summaryTable.Cols[0].Spec + ":" + summaryTable.Rows[index].Spec,
			Value: int(value.TotalAmount),
			Unit:  "",
		})
		summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
			Label: "Tổng số lượng đơn ngày " + summaryRow.Spec,
			Spec:  summaryTable.Cols[1].Spec + ":" + summaryTable.Rows[index].Spec,
			Value: int(value.Count),
			Unit:  "",
		})
	}

	return &summaryTable
}

func buildTopSellTable(args []*sqlstore.TopSellItem) *summary.SummaryTable {
	var summaryTable summary.SummaryTable
	summaryTable.Label = "Sản phẩm bán chạy"
	summaryTable.Cols = append(summaryTable.Cols, summary.SummaryColRow{
		Spec:  "product_code:",
		Label: "Mã sản phẩm",
	})
	summaryTable.Cols = append(summaryTable.Cols, summary.SummaryColRow{
		Spec:  "product_id:",
		Label: "ID sản phẩm",
	})
	summaryTable.Cols = append(summaryTable.Cols, summary.SummaryColRow{
		Spec:  "img_urls:",
		Label: "hình ảnh sản phẩm",
	})
	summaryTable.Cols = append(summaryTable.Cols, summary.SummaryColRow{
		Spec:  "name_product:",
		Label: "Tên Sản Phẩm",
	})
	summaryTable.Cols = append(summaryTable.Cols, summary.SummaryColRow{
		Spec:  "sum(quantity)",
		Label: "Số Lượng Sản Phẩm",
	})
	for index := range args {
		summaryTable.Rows = append(summaryTable.Rows, summary.SummaryColRow{
			Spec:  "groupby(product_id),orderby(sum(quantity)),row_number(" + strconv.Itoa(index) + ")" + "status!=-1",
			Label: strconv.Itoa(index),
		})
	}

	for index, value := range args {
		summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
			Spec:  summaryTable.Cols[0].Spec + summaryTable.Rows[index].Spec,
			Label: value.ProductCode,
		})
		summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
			Spec:  summaryTable.Cols[1].Spec + summaryTable.Rows[index].Spec,
			Label: value.ProductId.String(),
		})
		summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
			Spec:      summaryTable.Cols[2].Spec + summaryTable.Rows[index].Spec,
			ImageUrls: value.ImageUrls,
		})
		summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
			Spec:  summaryTable.Cols[3].Spec + summaryTable.Rows[index].Spec,
			Label: value.Name,
		})
		summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
			Spec:  summaryTable.Cols[4].Spec + summaryTable.Rows[index].Spec,
			Value: value.Count,
		})
	}

	return &summaryTable
}

func buildStaffOrderTable(args []*sqlstore.StaffOrder) *summary.SummaryTable {
	var summaryTable summary.SummaryTable
	summaryTable.Label = "Doanh số bán hàng theo nhân viên"
	summaryTable.Cols = append(summaryTable.Cols, summary.SummaryColRow{
		Spec:  "staff_id:",
		Label: "Mã Nhân viên",
	})
	summaryTable.Cols = append(summaryTable.Cols, summary.SummaryColRow{
		Spec:  "staff_name:",
		Label: "Tên nhân viên",
	})
	summaryTable.Cols = append(summaryTable.Cols, summary.SummaryColRow{
		Spec:  "count(order_id):",
		Label: "Tổng đơn",
	})
	summaryTable.Cols = append(summaryTable.Cols, summary.SummaryColRow{
		Spec:  "sum(total_amount):",
		Label: "Tổng doanh thu",
	})

	for index, value := range args {
		if value.UserID == 0 {
			continue
		}
		summaryTable.Rows = append(summaryTable.Rows, summary.SummaryColRow{
			Spec:  "STT-" + strconv.Itoa(index) + "status!=-1",
			Label: strconv.Itoa(index),
		})
		summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
			Spec:  summaryTable.Cols[0].Spec + summaryTable.Rows[index].Spec,
			Label: value.UserID.String(),
		})
		summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
			Spec:  summaryTable.Cols[1].Spec + summaryTable.Rows[index].Spec,
			Label: value.UserName,
		})
		summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
			Spec:  summaryTable.Cols[2].Spec + summaryTable.Rows[index].Spec,
			Value: int(value.TotalCount),
		})
		summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
			Spec:  summaryTable.Cols[3].Spec + summaryTable.Rows[index].Spec,
			Value: int(value.TotalAmount),
		})
	}
	return &summaryTable
}

func buildOtherSourceTable(args []*sqlstore.StaffOrder) *summary.SummaryTable {
	var summaryTable summary.SummaryTable
	summaryTable.Label = "Doanh số bán hàng từ nguồn khác"
	summaryTable.Cols = append(summaryTable.Cols, summary.SummaryColRow{
		Spec:  "other_source_id:",
		Label: "ID nguồn khác",
	})
	summaryTable.Cols = append(summaryTable.Cols, summary.SummaryColRow{
		Spec:  "other_soure_type:",
		Label: "Loại nguồn",
	})
	summaryTable.Cols = append(summaryTable.Cols, summary.SummaryColRow{
		Spec:  "count(order_id):",
		Label: "Tổng đơn",
	})
	summaryTable.Cols = append(summaryTable.Cols, summary.SummaryColRow{
		Spec:  "sum(total_amount):",
		Label: "Tổng doanh thu",
	})

	for index, value := range args {
		if value.UserID != 0 {
			continue
		}
		summaryTable.Rows = append(summaryTable.Rows, summary.SummaryColRow{
			Spec:  "groupby(other_source_id),orderby(sum(total_amount)),row_number(" + strconv.Itoa(index) + ")" + "status!=-1",
			Label: strconv.Itoa(index),
		})
		summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
			Spec:  summaryTable.Cols[0].Spec + summaryTable.Rows[index].Spec,
			Label: value.OrderSourceID.String(),
		})
		summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
			Spec:  summaryTable.Cols[1].Spec + summaryTable.Rows[index].Spec,
			Label: value.OrderSourceType,
		})
		summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
			Spec:  summaryTable.Cols[2].Spec + summaryTable.Rows[index].Spec,
			Value: int(value.TotalCount),
		})
		summaryTable.Data = append(summaryTable.Data, summary.SummaryItem{
			Spec:  summaryTable.Cols[3].Spec + summaryTable.Rows[index].Spec,
			Value: int(value.TotalAmount),
		})
	}
	return &summaryTable
}
