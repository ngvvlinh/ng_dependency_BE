package aggregate

import (
	"context"
	"time"

	"o.o/api/main/refund"
	"o.o/api/top/types/etc/status3"
	com "o.o/backend/com/main"
	catalogconvert "o.o/backend/com/main/catalog/convert"
	"o.o/backend/com/main/ordering/model"
	ordermodelx "o.o/backend/com/main/ordering/modelx"
	"o.o/backend/com/main/refund/convert"
	"o.o/backend/com/main/refund/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi"
	"o.o/capi/dot"
	"o.o/common/l"
)

var ll = l.New()
var _ refund.Aggregate = &RefundAggregate{}
var scheme = conversion.Build(convert.RegisterConversions, catalogconvert.RegisterConversions)

type RefundAggregate struct {
	db       *cmsql.Database
	store    sqlstore.RefundStoreFactory
	eventBus capi.EventBus
	bus      bus.Bus
}

func NewRefundAggregate(
	database com.MainDB, eventBus capi.EventBus,
) *RefundAggregate {
	return &RefundAggregate{
		db:       database,
		eventBus: eventBus,
		store:    sqlstore.NewRefundStore(database),
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
	preLine, err := a.checkLineOrder(ctx, args.ShopID, refundResult.OrderID, refundResult.ID, refundResult.Lines)
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
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Phiếu trả hàng %v đã hủy hoặc đã xác nhận không thể cập nhập", refundDB.Code)
	}
	err = scheme.Convert(args, refundDB)
	if err != nil {
		return nil, err
	}
	preLine, err := a.checkLineOrder(ctx, args.ShopID, refundDB.OrderID, refundDB.ID, refundDB.Lines)
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

func (a *RefundAggregate) checkLineOrder(ctx context.Context, shopID dot.ID, orderID dot.ID, refundID dot.ID, lines []*refund.RefundLine) (*refund.CheckReceiptLinesResponse, error) {

	queryOrder := &ordermodelx.GetOrderQuery{
		OrderID: orderID,
		ShopID:  shopID,
	}
	err := bus.Dispatch(ctx, queryOrder)
	if err != nil {
		return nil, err
	}
	order := queryOrder.Result.Order

	var linesVariant = make(map[dot.ID]*model.OrderLine, len(order.Lines))
	for _, value := range order.Lines {
		linesVariant[value.VariantID] = value
	}
	refundByOrder, err := a.store(ctx).ShopID(shopID).OrderID(orderID).ListRefunds()
	for _, value := range refundByOrder {
		if value.ID == refundID || value.Status == status3.N {
			continue
		}
		for _, _value := range value.Lines {
			linesVariant[_value.VariantID].Quantity = linesVariant[_value.VariantID].Quantity - _value.Quantity
		}
	}
	for key, value := range lines {
		if linesVariant[value.VariantID] == nil {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Sản phẩm không tồn tại trong đơn hàng %v", queryOrder.Result.Order.Code)
		}
		if linesVariant[value.VariantID].Quantity < value.Quantity {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Số lượng sản phẩm trong đơn trả hàng lớn hơn đơn hàng")
		}
		lines[key].Code = linesVariant[value.VariantID].Code
		lines[key].ImageURL = linesVariant[value.VariantID].ImageURL
		lines[key].ProductName = linesVariant[value.VariantID].ProductName
		lines[key].RetailPrice = linesVariant[value.VariantID].RetailPrice
		lines[key].ProductID = linesVariant[value.VariantID].ProductID
		err = scheme.Convert(linesVariant[value.VariantID].Attributes, &lines[key].Attributes)
		if err != nil {
			return nil, err
		}
	}
	return &refund.CheckReceiptLinesResponse{
		CustomerID: queryOrder.Result.Order.CustomerID,
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
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Phiếu trả hàng %v đã hủy hoặc đã xác nhận không thể cập nhập trạng thái")
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
		return cm.Errorf(cm.InvalidArgument, nil, "basket_value không được nhỏ hơn 0")
	}
	if args.TotalAmount < 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "total_amount không được nhỏ hơn 0")
	}
	if args.BasketValue+args.TotalAdjustment != args.TotalAmount {
		return cm.Errorf(cm.InvalidArgument, nil, "total_amount không đúng, basket_value + total_adjustment == total_amount")
	}
	basketValue := 0
	for _, line := range args.Lines {
		basketValue += line.Quantity * (line.RetailPrice + line.Adjustment)
	}
	if basketValue != args.BasketValue {
		return cm.Errorf(cm.InvalidArgument, nil, "basket_value không đúng, basket value đang là %, giá trị mong đợi là %v", args.BasketValue, basketValue)
	}
	totalAdjustment := 0
	for _, line := range args.AdjustmentLines {
		totalAdjustment += line.Amount
	}
	if totalAdjustment != args.TotalAdjustment {
		return cm.Errorf(cm.InvalidArgument, nil, "total_adjustment không đúng, basket value đang là %, giá trị mong đợi là %v", args.TotalAdjustment, totalAdjustment)
	}
	return nil
}
