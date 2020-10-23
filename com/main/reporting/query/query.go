package query

import (
	"context"

	"o.o/api/main/identity"
	"o.o/api/main/ordering"
	"o.o/api/main/receipting"
	"o.o/api/main/reporting"
	"o.o/api/top/types/etc/receipt_ref"
	"o.o/api/top/types/etc/receipt_type"
	"o.o/api/top/types/etc/status3"
	"o.o/backend/pkg/common/bus"
	"o.o/capi/dot"
)

var _ reporting.QueryService = &ReportQuery{}

type ReportQuery struct {
	orderQuery    ordering.QueryBus
	identityQuery identity.QueryBus
	receiptQuery  receipting.QueryBus
}

func NewReportQuery(
	orderQuery ordering.QueryBus,
	identityQuery identity.QueryBus,
	receiptQuery receipting.QueryBus,
) *ReportQuery {
	return &ReportQuery{
		orderQuery:    orderQuery,
		identityQuery: identityQuery,
		receiptQuery:  receiptQuery,
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
		IsContains:  true,
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
