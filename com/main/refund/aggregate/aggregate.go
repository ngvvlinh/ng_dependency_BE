package aggregate

import (
	"context"
	"time"

	"etop.vn/api/main/refund"
	"etop.vn/api/top/types/etc/status3"
	catalogconvert "etop.vn/backend/com/main/catalog/convert"
	"etop.vn/backend/com/main/ordering/model"
	ordermodelx "etop.vn/backend/com/main/ordering/modelx"
	"etop.vn/backend/com/main/refund/convert"
	"etop.vn/backend/com/main/refund/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/conversion"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/capi"
	"etop.vn/capi/dot"
	"etop.vn/common/l"
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

func NewReceiptAggregate(
	database *cmsql.Database, eventBus capi.EventBus,
) *RefundAggregate {
	return &RefundAggregate{
		db:       database,
		eventBus: eventBus,
		store:    sqlstore.NewRefundStore(database),
	}
}

func (a *RefundAggregate) MessageBus() refund.CommandBus {
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
	if preLine.BasketValue < refundResult.Discount {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Giảm giá không được lớn hơn giá trị hàng trả")
	}
	refundResult.TotalAmount = preLine.BasketValue - refundResult.Discount
	refundResult.BasketValue = preLine.BasketValue
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
	refundResult.CustomerID = preLine.CustomerID
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

	if preLine.BasketValue < refundDB.Discount {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Giảm giá không được lớn hơn giá trị hàng trả")
	}
	refundDB.TotalAmount = preLine.BasketValue - refundDB.Discount
	err = a.store(ctx).ShopID(args.ShopID).ID(args.ID).UpdateRefundAll(refundDB)
	return refundDB, err
}

func (a *RefundAggregate) checkLineOrder(ctx context.Context, shopID dot.ID, orderID dot.ID, refundID dot.ID, lines []*refund.RefundLine) (*refund.CheckReceiptLinesResponse, error) {
	basketValue := 0
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
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Số lượng sản phẩm trong đơn trả hàng lơn hơn đơn hàng")
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
		basketValue = basketValue + lines[key].RetailPrice*lines[key].Quantity
	}
	return &refund.CheckReceiptLinesResponse{
		CustomerID:  queryOrder.Result.Order.CustomerID,
		RefundLine:  lines,
		BasketValue: basketValue,
	}, nil
}

func (a *RefundAggregate) CancelRefund(ctx context.Context, args *refund.CancelRefundArgs) (*refund.Refund, error) {
	refundDB, err := a.store(ctx).ID(args.ID).ShopID(args.ShopID).GetRefund()
	if err != nil {
		return nil, err
	}
	if refundDB.Status != status3.Z {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Phiếu trả hàng %v đã hủy hoặc đã xác nhận không thể cập nhập trạng thái")
	}
	refundDB.CancelledAt = time.Now()
	refundDB.Status = status3.N
	refundDB.CancelReason = args.CancelReason
	err = a.store(ctx).ID(args.ID).ShopID(args.ShopID).UpdateRefundAll(refundDB)
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
	refundDB.CancelledAt = time.Now()
	refundDB.Status = status3.P
	err = a.store(ctx).ID(args.ID).ShopID(args.ShopID).UpdateRefundAll(refundDB)
	if err != nil {
		return nil, err
	}
	event := &refund.ConfirmedRefundEvent{
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
