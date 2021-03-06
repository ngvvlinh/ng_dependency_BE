package aggregate

import (
	"context"
	"time"

	"o.o/api/main/refund"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/status5"
	com "o.o/backend/com/main"
	catalogconvert "o.o/backend/com/main/catalog/convert"
	"o.o/backend/com/main/ordering/model"
	ordermodel "o.o/backend/com/main/ordering/model"
	ordermodelx "o.o/backend/com/main/ordering/modelx"
	"o.o/backend/com/main/refund/convert"
	refundstore "o.o/backend/com/main/refund/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/capi"
	"o.o/capi/dot"
	"o.o/common/l"
)

var ll = l.New()
var _ refund.Aggregate = &RefundAggregate{}
var scheme = conversion.Build(convert.RegisterConversions, catalogconvert.RegisterConversions)

type RefundAggregate struct {
	db       *cmsql.Database
	store    refundstore.RefundStoreFactory
	eventBus capi.EventBus
	bus      bus.Bus

	OrderStore sqlstore.OrderStoreInterface
}

func NewRefundAggregate(
	database com.MainDB, eventBus capi.EventBus,
	OrderStore sqlstore.OrderStoreInterface,

) *RefundAggregate {
	return &RefundAggregate{
		db:         database,
		eventBus:   eventBus,
		store:      refundstore.NewRefundStore(database),
		OrderStore: OrderStore,
	}
}

func RefundAggregateMessageBus(a *RefundAggregate) refund.CommandBus {
	b := bus.New()
	return refund.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *RefundAggregate) CreateRefund(ctx context.Context, args *refund.CreateRefundArgs) (*refund.Refund, error) {
	var refundResult refund.Refund
	err := scheme.Convert(args, &refundResult)
	if err != nil {
		return nil, err
	}
	preLine, err := a.checkLineOrder(ctx, args.ShopID, refundResult.OrderID, refundResult.ID, refundResult.Lines, true)
	if err != nil {
		return nil, err
	}
	refundResult.CustomerID = preLine.CustomerID
	refundResult.Lines = preLine.RefundLine
	err = checkRefundBeforeCreateOrUpdate(&refundResult)
	if err != nil {
		return nil, err
	}
	var maxCodeNorm = 1
	refundCode, err := a.store(ctx).ShopID(args.ShopID).GetRefundByMaximumCodeNorm()
	switch cm.ErrorCode(err) {
	case cm.NoError:
		maxCodeNorm = refundCode.CodeNorm + 1
	case cm.NotFound:
		// no-op
	default:
		return nil, err
	}
	refundResult.Code = convert.GenerateCode(maxCodeNorm)
	refundResult.CodeNorm = maxCodeNorm
	err = a.store(ctx).Create(&refundResult)
	return &refundResult, err
}

func (a *RefundAggregate) UpdateRefund(ctx context.Context, args *refund.UpdateRefundArgs) (*refund.Refund, error) {
	refundDB, err := a.store(ctx).ID(args.ID).ShopID(args.ShopID).GetRefund()
	if err != nil {
		return nil, err
	}
	if refundDB.Status != status3.Z {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Phi???u tr??? h??ng %v ???? h???y ho???c ???? x??c nh???n kh??ng th??? c???p nh???p", refundDB.Code)
	}
	err = scheme.Convert(args, refundDB)
	if err != nil {
		return nil, err
	}
	preLine, err := a.checkLineOrder(ctx, args.ShopID, refundDB.OrderID, refundDB.ID, refundDB.Lines, false)
	if err != nil {
		return nil, err
	}
	refundDB.CustomerID = preLine.CustomerID
	refundDB.Lines = preLine.RefundLine
	refundDB.BasketValue = preLine.BasketValue
	err = checkRefundBeforeCreateOrUpdate(refundDB)
	if err != nil {
		return nil, err
	}
	err = a.store(ctx).ShopID(args.ShopID).ID(args.ID).UpdateRefundAll(refundDB)
	return refundDB, err
}

func (a *RefundAggregate) checkOrder(arg *ordermodel.Order) error {
	if len(arg.FulfillmentIDs) == 0 {
		// ????n kh??ng c?? giao h??ng
		if arg.Status == status5.NS || arg.Status == status5.Z || arg.Status == status5.P {
			return nil
		}
	} else {
		//  ????n c?? giao h??ng v?? b??? tr??? h??ng
		if arg.FulfillmentShippingStatus == status5.NS {
			for _, v := range arg.FulfillmentShippingStates {
				if v == "returned" {
					return nil
				}
			}
		}
		//  ????n c?? giao h??ng v?? ????n giao th??nh c??ng
		if arg.FulfillmentShippingStatus == status5.P {
			return nil
		}
	}

	return cm.Errorf(cm.FailedPrecondition, nil, "Kh??ng th??? t???o ????n tr??? h??ng cho ????n h??ng %v", arg.ID)
}

func (a *RefundAggregate) checkLineOrder(ctx context.Context, shopID dot.ID, orderID dot.ID, refundID dot.ID, lines []*refund.RefundLine, isCreate bool) (*refund.CheckReceiptLinesResponse, error) {
	queryOrder := &ordermodelx.GetOrderQuery{
		OrderID: orderID,
		ShopID:  shopID,
	}
	err := a.OrderStore.GetOrder(ctx, queryOrder)
	if err != nil {
		return nil, err
	}
	order := queryOrder.Result.Order
	if isCreate {
		err = a.checkOrder(order)
		if err != nil {
			return nil, err
		}
	}
	var linesByCode = make(map[string]*model.OrderLine, len(order.Lines))
	for _, value := range order.Lines {
		linesByCode[value.Code] = value
	}
	refundByOrder, err := a.store(ctx).ShopID(shopID).OrderID(orderID).ListRefunds()
	for _, value := range refundByOrder {
		if value.ID == refundID || value.Status == status3.N {
			continue
		}
		for _, _value := range value.Lines {
			linesByCode[_value.Code].Quantity = linesByCode[_value.Code].Quantity - _value.Quantity
		}
	}
	for key, value := range lines {
		if linesByCode[value.Code] == nil {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "S???n ph???m kh??ng t???n t???i trong ????n h??ng %v", order.Code)
		}
		if linesByCode[value.Code].Quantity < value.Quantity {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "S??? l?????ng s???n ph???m trong ????n tr??? h??ng l???n h??n ????n h??ng")
		}
		lines[key].Code = linesByCode[value.Code].Code
		lines[key].ImageURL = linesByCode[value.Code].ImageURL
		lines[key].ProductName = linesByCode[value.Code].ProductName
		lines[key].RetailPrice = linesByCode[value.Code].RetailPrice
		lines[key].ProductID = linesByCode[value.Code].ProductID
		lines[key].VariantID = linesByCode[value.Code].VariantID
		err = scheme.Convert(linesByCode[value.Code].Attributes, &lines[key].Attributes)
		if err != nil {
			return nil, err
		}
	}
	return &refund.CheckReceiptLinesResponse{
		CustomerID: order.CustomerID,
		RefundLine: lines,
	}, nil
}

func (a *RefundAggregate) CancelRefund(ctx context.Context, args *refund.CancelRefundArgs) (*refund.Refund, error) {
	refundDB, err := a.store(ctx).ID(args.ID).ShopID(args.ShopID).GetRefund()
	if err != nil {
		return nil, err
	}
	refundDB.CancelledAt = time.Now()
	refundDB.Status = status3.N
	refundDB.CancelReason = args.CancelReason
	refundDB.UpdatedAt = time.Now()
	refundDB.UpdatedBy = args.UpdatedBy
	err = a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		err = a.store(ctx).ID(args.ID).ShopID(args.ShopID).UpdateRefundAll(refundDB)
		if err != nil {
			return err
		}
		event := &refund.RefundCancelledEvent{
			ShopID:               args.ShopID,
			RefundID:             args.ID,
			UpdatedBy:            args.UpdatedBy,
			AutoInventoryVoucher: args.AutoInventoryVoucher,
		}
		err = a.eventBus.Publish(ctx, event)
		if err != nil {
			return err
		}
		return nil
	})

	return refundDB, err
}

func (a *RefundAggregate) ConfirmRefund(ctx context.Context, args *refund.ConfirmRefundArgs) (*refund.Refund, error) {
	refundDB, err := a.store(ctx).ID(args.ID).ShopID(args.ShopID).GetRefund()
	if err != nil {
		return nil, err
	}
	if refundDB.Status != status3.Z {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Phi???u tr??? h??ng %v ???? h???y ho???c ???? x??c nh???n kh??ng th??? c???p nh???p tr???ng th??i")
	}
	refundDB.ConfirmedAt = time.Now()
	refundDB.Status = status3.P
	refundDB.UpdatedAt = time.Now()
	refundDB.UpdatedBy = args.UpdatedBy
	err = a.store(ctx).ID(args.ID).ShopID(args.ShopID).UpdateRefundAll(refundDB)
	if err != nil {
		return nil, err
	}
	event := &refund.RefundConfirmedEvent{
		ShopID:               args.ShopID,
		RefundID:             args.ID,
		AutoInventoryVoucher: args.AutoInventoryVoucher,
	}
	err = a.eventBus.Publish(ctx, event)
	if err != nil {
		return nil, err
	}
	return refundDB, nil
}

func checkRefundBeforeCreateOrUpdate(args *refund.Refund) error {
	if args.BasketValue < 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "basket_value kh??ng ???????c nh??? h??n 0")
	}
	if args.TotalAmount < 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "total_amount kh??ng ???????c nh??? h??n 0")
	}
	if args.BasketValue+args.TotalAdjustment != args.TotalAmount {
		return cm.Errorf(cm.InvalidArgument, nil, "total_amount kh??ng ????ng, basket_value + total_adjustment == total_amount")
	}
	basketValue := 0
	for _, line := range args.Lines {
		basketValue += line.Quantity * (line.RetailPrice + line.Adjustment)
	}
	if basketValue != args.BasketValue {
		return cm.Errorf(cm.InvalidArgument, nil, "basket_value kh??ng ????ng, basket value ??ang l?? %, gi?? tr??? mong ?????i l?? %v", args.BasketValue, basketValue)
	}
	totalAdjustment := 0
	for _, line := range args.AdjustmentLines {
		totalAdjustment += line.Amount
	}
	if totalAdjustment != args.TotalAdjustment {
		return cm.Errorf(cm.InvalidArgument, nil, "total_adjustment kh??ng ????ng, basket value ??ang l?? %, gi?? tr??? mong ?????i l?? %v", args.TotalAdjustment, totalAdjustment)
	}
	return nil
}
