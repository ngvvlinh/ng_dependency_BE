package query

import (
	"context"
	"time"

	"o.o/api/main/catalog"
	"o.o/api/main/identity"
	"o.o/api/main/ordering"
	"o.o/api/main/receipting"
	"o.o/api/main/reporting"
	"o.o/api/main/stocktaking"
	"o.o/api/top/types/etc/receipt_ref"
	"o.o/api/top/types/etc/receipt_type"
	"o.o/api/top/types/etc/report_time_filter"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/stocktake_type"
	"o.o/backend/pkg/common/bus"
	"o.o/capi/dot"
)

var _ reporting.QueryService = &ReportQuery{}

type ReportQuery struct {
	orderQuery       ordering.QueryBus
	identityQuery    identity.QueryBus
	receiptQuery     receipting.QueryBus
	catalogQuery     catalog.QueryBus
	stocktakingQuery stocktaking.QueryBus
}

func NewReportQuery(
	orderQuery ordering.QueryBus,
	identityQuery identity.QueryBus,
	receiptQuery receipting.QueryBus,
	catalogQuery catalog.QueryBus,
	stocktakingQuery stocktaking.QueryBus,
) *ReportQuery {
	return &ReportQuery{
		orderQuery:       orderQuery,
		identityQuery:    identityQuery,
		receiptQuery:     receiptQuery,
		catalogQuery:     catalogQuery,
		stocktakingQuery: stocktakingQuery,
	}
}

func ReportQueryMessageBus(q *ReportQuery) reporting.QueryBus {
	b := bus.New()
	return reporting.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (r *ReportQuery) ReportOrders(
	ctx context.Context, args *reporting.ReportOrdersArgs,
) (result []*reporting.ReportOrder, _ error) {
	listOrdersQuery := &ordering.ListOrdersConfirmedQuery{
		ShopID:        args.ShopID,
		CreatedAtFrom: args.CreatedAtFrom,
		CreatedAtTo:   args.CreatedAtTo,
		CreatedBy:     args.CreatedBy,
	}
	if err := r.orderQuery.Dispatch(ctx, listOrdersQuery); err != nil {
		return nil, err
	}

	orders := listOrdersQuery.Result
	if len(orders) == 0 {
		return nil, nil
	}
	var orderIDs []dot.ID
	mReportOrderEndOfDay := make(map[dot.ID]*reporting.ReportOrder)
	{
		for _, order := range orders {
			orderIDs = append(orderIDs, order.ID)
			report := &reporting.ReportOrder{
				OrderCode:     order.Code,
				CreatedAt:     order.CreatedAt,
				TotalItems:    order.TotalItems,
				TotalFee:      order.TotalFee,
				TotalDiscount: order.TotalDiscount,
				TotalAmount:   order.TotalAmount,
			}
			result = append(result, report)
			mReportOrderEndOfDay[order.ID] = report
		}
	}

	listReceiptsQuery := &receipting.ListReceiptsByRefsAndStatusAndTypeQuery{
		ShopID:      args.ShopID,
		RefIDs:      orderIDs,
		RefType:     receipt_ref.Order,
		ReceiptType: receipt_type.Receipt,
		Status:      status3.P,
	}
	if err := r.receiptQuery.Dispatch(ctx, listReceiptsQuery); err != nil {
		return nil, err
	}
	receipts := listReceiptsQuery.Result
	for _, receipt := range receipts {
		for _, line := range receipt.Lines {
			mReportOrderEndOfDay[line.RefID].Revenue += line.Amount
		}
	}

	return result, nil
}

func (r *ReportQuery) ReportIncomeStatement(
	ctx context.Context, args *reporting.ReportIncomeStatementArgs,
) (result map[int]*reporting.ReportIncomeStatement, _ error) {
	var (
		orders     []*ordering.Order
		variants   []*catalog.ShopVariant
		receipts   []*receipting.Receipt
		stocktakes []*stocktaking.ShopStocktake
	)
	timeFilter := args.TimeFilter

	result = make(map[int]*reporting.ReportIncomeStatement)

	createdAtFrom := time.Date(args.Year, 1, 1, 0, 0, 0, 0, time.Local)
	createdAtTo := time.Date(args.Year, 12, 31, 0, 0, 0, 0, time.Local)

	listOrdersQuery := &ordering.ListOrdersConfirmedQuery{
		ShopID: args.ShopID,
	}
	if args.TimeFilter != report_time_filter.Year {
		listOrdersQuery.CreatedAtFrom = createdAtFrom
		listOrdersQuery.CreatedAtTo = createdAtTo
	}
	if err := r.orderQuery.Dispatch(ctx, listOrdersQuery); err != nil {
		return nil, err
	}
	orders = listOrdersQuery.Result

	var variantIDs []dot.ID
	{
		mVariant := make(map[dot.ID]bool)
		for _, order := range orders {
			for _, variantID := range order.VariantIDs {
				if _, ok := mVariant[variantID]; ok {
					continue
				}
				mVariant[variantID] = true
				variantIDs = append(variantIDs, variantID)
			}
		}
	}

	listVariantsQuery := &catalog.ListShopVariantsByIDsQuery{
		IDs:    variantIDs,
		ShopID: args.ShopID,
	}
	if err := r.catalogQuery.Dispatch(ctx, listVariantsQuery); err != nil {
		return nil, err
	}
	variants = listVariantsQuery.Result.Variants
	mVariant := make(map[dot.ID]*catalog.ShopVariant)
	{
		for _, variant := range variants {
			mVariant[variant.VariantID] = variant
		}
	}

	listStocktakeQuery := &stocktaking.ListStocktakeQuery{
		ShopID: args.ShopID,
		Type:   stocktake_type.WrapStocktakeType(stocktake_type.Discard),
	}
	if args.TimeFilter != report_time_filter.Year {
		listStocktakeQuery.CreatedAtFrom = createdAtFrom
		listStocktakeQuery.CreatedAtTo = createdAtTo
	}
	if err := r.stocktakingQuery.Dispatch(ctx, listStocktakeQuery); err != nil {
		return nil, err
	}
	stocktakes = listStocktakeQuery.Result.Stocktakes

	listReceiptsQuery := &receipting.ListReceiptsQuery{
		ShopID: args.ShopID,
	}
	if args.TimeFilter != report_time_filter.Year {
		listReceiptsQuery.CreatedAtFrom = createdAtFrom
		listReceiptsQuery.CreatedAtTo = createdAtTo
	}
	if err := r.receiptQuery.Dispatch(ctx, listReceiptsQuery); err != nil {
		return nil, err
	}
	receipts = listReceiptsQuery.Result.Receipts

	for _, order := range orders {
		key := getTimeKey(order.CreatedAt, timeFilter)
		if _, ok := result[key]; !ok {
			result[key] = &reporting.ReportIncomeStatement{}
		}
		// Doanh thu b??n h??ng (1) = Gi?? tr??? h??ng ho?? + Ph??? thu
		result[key].Revenue += order.BasketValue + order.TotalFee
		// Gi???m tr??? Doanh thu (2) = T???ng (c??c kho???ng gi???m gi??)
		result[key].Discounts += order.TotalDiscount
		for _, line := range order.Lines {
			variant, ok := mVariant[line.VariantID]
			if !ok {
				continue
			}
			// Gi?? v???n h??ng b??n (4) = T???ng(Gi?? v???n * S??? l?????ng item b??n trong c??c ho?? ????n)
			result[key].CostPrice += variant.CostPrice * line.Quantity
		}
	}

	for _, stockTake := range stocktakes {
		key := getTimeKey(stockTake.CreatedAt, timeFilter)
		if _, ok := result[key]; !ok {
			result[key] = &reporting.ReportIncomeStatement{}
		}
		for _, line := range stockTake.Lines {
			// Xu???t h???y h??ng h??a  (6.2) = T???ng gi?? tr???  xu???t h???y h??ng h??a (Gi?? v???n * s??? l?????ng item xu???t hu???)
			result[key].Discards += line.CostPrice * line.NewQuantity
		}
	}

	for _, receipt := range receipts {
		key := getTimeKey(receipt.CreatedAt, timeFilter)
		if _, ok := result[key]; !ok {
			result[key] = &reporting.ReportIncomeStatement{}
		}
		if receipt.Type == receipt_type.Payment {
			if receipt.RefType == receipt_ref.Fulfillment {
				// Ph?? giao h??ng (6.1) = T???ng (Phi???u chi cho giao h??ng tr??n FFM)
				result[key].ShippingFee += receipt.Amount
			}

			// Chi ph?? kh??c (9) = T???ng (Phi???u chi kh??ng ???????c g???n v???i Ho?? ????n)
			if receipt.RefType != receipt_ref.Order {
				result[key].OtherExpenses += receipt.Amount
			} else {
				// Kh??c (6.3) = T???ng (Phi???u chi kh??c trong c??c Ho?? ????n)
				result[key].Others += receipt.Amount
			}
		} else {
			// Thu nh???p kh??c (8) = T???ng (Phi???u thu kh??ng ???????c g???n v???i Ho?? ????n)
			if receipt.RefType != receipt_ref.Order {
				result[key].OtherIncomes += receipt.Amount
			}
		}
	}

	for key := range result {
		// Doanh thu thu???n (3=1-2)
		result[key].NetRevenue = result[key].Revenue - result[key].Discounts
		// L???i nhu???n g???p v??? b??n h??ng (5=3-4)
		result[key].GrossProfit = result[key].NetRevenue - result[key].CostPrice
		// Chi ph?? (6 = 6.1 + 6.2 + 6.3)
		result[key].Expenses = result[key].ShippingFee + result[key].Discards + result[key].Others
		// L???i nhu???n t??? ho???t ?????ng kinh doanh (7=5-6)
		result[key].ProfitStatement = result[key].GrossProfit - result[key].Expenses
		// L???i nhu???n thu???n (10=(7+8)-9)
		result[key].NetProfit = (result[key].ProfitStatement + result[key].OtherIncomes) - result[key].OtherExpenses
	}

	return result, nil
}

func getTimeKey(_time time.Time, timeFilter report_time_filter.TimeFilter) int {
	switch timeFilter {
	case report_time_filter.Month:
		return int(_time.Month())
	case report_time_filter.Quater:
		return getQuater(_time)
	case report_time_filter.Year:
		return _time.Year()
	}
	return -1
}

func getQuater(_time time.Time) int {
	return 1 + ((int(_time.Month()) - 1) / 3)
}
