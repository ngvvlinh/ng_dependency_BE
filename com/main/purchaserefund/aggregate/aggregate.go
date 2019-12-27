package aggregate

import (
	"context"
	"time"

	"etop.vn/api/main/purchaseorder"
	"etop.vn/api/main/purchaserefund"
	"etop.vn/api/top/types/etc/status3"
	catalogconvert "etop.vn/backend/com/main/catalog/convert"
	"etop.vn/backend/com/main/purchaserefund/convert"
	"etop.vn/backend/com/main/purchaserefund/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/conversion"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/capi"
	"etop.vn/capi/dot"
	"etop.vn/common/l"
)

var ll = l.New()
var _ purchaserefund.Aggregate = &PurchaseRefundAggregate{}
var scheme = conversion.Build(convert.RegisterConversions, catalogconvert.RegisterConversions)

type PurchaseRefundAggregate struct {
	PurchaseOrderQuery purchaseorder.QueryBus
	db                 *cmsql.Database
	store              sqlstore.PurchaseRefundStoreFactory
	eventBus           capi.EventBus
}

func NewPurchaseRefundAggregate(
	database *cmsql.Database,
	eventBus capi.EventBus,
	purchaseOrderQuery purchaseorder.QueryBus,
) *PurchaseRefundAggregate {
	return &PurchaseRefundAggregate{
		db:                 database,
		eventBus:           eventBus,
		store:              sqlstore.NewPurchaseRefundStore(database),
		PurchaseOrderQuery: purchaseOrderQuery,
	}
}

func (a *PurchaseRefundAggregate) MessageBus() purchaserefund.CommandBus {
	b := bus.New()
	return purchaserefund.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *PurchaseRefundAggregate) CreatePurchaseRefund(ctx context.Context, args *purchaserefund.CreatePurchaseRefundArgs) (*purchaserefund.PurchaseRefund, error) {
	var purchaseRefundResult purchaserefund.PurchaseRefund
	err := scheme.Convert(args, &purchaseRefundResult)
	if err != nil {
		return nil, err
	}

	preLine, err := a.checkLinePurchaseOrder(ctx, args.ShopID, purchaseRefundResult.PurchaseOrderID, purchaseRefundResult.ID, purchaseRefundResult.Lines)
	if err != nil {
		return nil, err
	}
	purchaseRefundResult.SupplierID = preLine.SupplierID
	purchaseRefundResult.Lines = preLine.PurchaseRefundLine
	if preLine.BasketValue < purchaseRefundResult.Discount {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Giảm giá không được lớn hơn giá trị hàng trả")
	}
	purchaseRefundResult.TotalAmount = preLine.BasketValue - purchaseRefundResult.Discount
	purchaseRefundResult.BasketValue = preLine.BasketValue
	var maxCodeNorm = 1
	purchaserefundCode, err := a.store(ctx).ShopID(args.ShopID).GetPurchaseRefundByMaximumCodeNorm()
	switch cm.ErrorCode(err) {
	case cm.NoError:
		maxCodeNorm = purchaserefundCode.CodeNorm + 1
	case cm.NotFound:
		// no-op
	default:
		return nil, err
	}
	purchaseRefundResult.SupplierID = preLine.SupplierID
	purchaseRefundResult.Code = convert.GenerateCode(maxCodeNorm)
	purchaseRefundResult.CodeNorm = maxCodeNorm
	err = a.store(ctx).Create(&purchaseRefundResult)
	return &purchaseRefundResult, err
}

func (a *PurchaseRefundAggregate) UpdatePurchaseRefund(ctx context.Context, args *purchaserefund.UpdatePurchaseRefundArgs) (*purchaserefund.PurchaseRefund, error) {
	purchaserefundDB, err := a.store(ctx).ID(args.ID).ShopID(args.ShopID).GetPurchaseRefund()
	if err != nil {
		return nil, err
	}
	if purchaserefundDB.Status != status3.Z {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Phiếu trả hàng %v đã hủy hoặc đã xác nhận không thể cập nhập", purchaserefundDB.Code)
	}
	err = scheme.Convert(args, purchaserefundDB)
	if err != nil {
		return nil, err
	}
	preLine, err := a.checkLinePurchaseOrder(ctx, args.ShopID, purchaserefundDB.PurchaseOrderID, purchaserefundDB.ID, purchaserefundDB.Lines)
	if err != nil {
		return nil, err
	}
	purchaserefundDB.SupplierID = preLine.SupplierID
	purchaserefundDB.Lines = preLine.PurchaseRefundLine
	purchaserefundDB.BasketValue = preLine.BasketValue

	if preLine.BasketValue < purchaserefundDB.Discount {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Giảm giá không được lớn hơn giá trị hàng trả")
	}
	purchaserefundDB.TotalAmount = preLine.BasketValue - purchaserefundDB.Discount
	err = a.store(ctx).ShopID(args.ShopID).ID(args.ID).UpdatePurchaseRefundAll(purchaserefundDB)
	return purchaserefundDB, err
}

func (a *PurchaseRefundAggregate) checkLinePurchaseOrder(ctx context.Context, shopID dot.ID, purchaseOrderID dot.ID, purchaserefundID dot.ID, lines []*purchaserefund.PurchaseRefundLine) (*purchaserefund.CheckReceiptLinesResponse, error) {
	basketValue := 0
	queryPurchaseOrder := &purchaseorder.GetPurchaseOrderByIDQuery{
		ID:     purchaseOrderID,
		ShopID: shopID,
	}
	err := a.PurchaseOrderQuery.Dispatch(ctx, queryPurchaseOrder)
	if err != nil {
		return nil, err
	}
	purchaseOrder := queryPurchaseOrder.Result

	var linesVariant = make(map[dot.ID]*purchaseorder.PurchaseOrderLine, len(purchaseOrder.Lines))
	for _, value := range purchaseOrder.Lines {
		linesVariant[value.VariantID] = value
	}
	purchaserefundByOrder, err := a.store(ctx).ShopID(shopID).PurchaseOrderID(purchaseOrderID).ListPurchaseRefunds()
	for _, value := range purchaserefundByOrder {
		if value.ID == purchaserefundID || value.Status == status3.N {
			continue
		}
		for _, line := range value.Lines {
			linesVariant[line.VariantID].Quantity = linesVariant[line.VariantID].Quantity - line.Quantity
		}
	}
	for key, value := range lines {
		if linesVariant[value.VariantID] == nil {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Sản phẩm không tồn tại trong đơn hàng %v", queryPurchaseOrder.Result.Code)
		}
		if linesVariant[value.VariantID].Quantity < value.Quantity {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Số lượng sản phẩm trong đơn trả hàng lớn hơn đơn hàng")
		}
		lines[key].Code = linesVariant[value.VariantID].Code
		lines[key].ImageURL = linesVariant[value.VariantID].ImageUrl
		lines[key].ProductName = linesVariant[value.VariantID].ProductName
		lines[key].PaymentPrice = linesVariant[value.VariantID].PaymentPrice
		lines[key].ProductID = linesVariant[value.VariantID].ProductID
		lines[key].Attributes = linesVariant[value.VariantID].Attributes
		basketValue = basketValue + lines[key].PaymentPrice*lines[key].Quantity
	}
	return &purchaserefund.CheckReceiptLinesResponse{
		SupplierID:         queryPurchaseOrder.Result.SupplierID,
		PurchaseRefundLine: lines,
		BasketValue:        basketValue,
	}, nil
}

func (a *PurchaseRefundAggregate) CancelPurchaseRefund(ctx context.Context, args *purchaserefund.CancelPurchaseRefundArgs) (*purchaserefund.PurchaseRefund, error) {
	purchaserefundDB, err := a.store(ctx).ID(args.ID).ShopID(args.ShopID).GetPurchaseRefund()
	if err != nil {
		return nil, err
	}
	if purchaserefundDB.Status != status3.N {
		purchaserefundDB.CancelledAt = time.Now()
		purchaserefundDB.Status = status3.N
		purchaserefundDB.CancelReason = args.CancelReason
		purchaserefundDB.UpdatedAt = time.Now()
		err = a.store(ctx).ID(args.ID).ShopID(args.ShopID).UpdatePurchaseRefundAll(purchaserefundDB)
		if err != nil {
			return nil, err
		}
	}
	event := &purchaserefund.PurchaseRefundCancelledEvent{
		PurchaseRefundID:     args.ID,
		ShopID:               args.ShopID,
		UpdatedBy:            args.UpdatedBy,
		AutoInventoryVoucher: args.AutoInventoryVoucher,
		InventoryOverStock:   args.InventoryOverStock,
	}
	err = a.eventBus.Publish(ctx, event)
	if err != nil {
		return nil, err
	}
	return purchaserefundDB, nil
}

func (a *PurchaseRefundAggregate) ConfirmPurchaseRefund(ctx context.Context, args *purchaserefund.ConfirmPurchaseRefundArgs) (*purchaserefund.PurchaseRefund, error) {
	purchaserefundDB, err := a.store(ctx).ID(args.ID).ShopID(args.ShopID).GetPurchaseRefund()
	if err != nil {
		return nil, err
	}
	if purchaserefundDB.Status != status3.Z {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Không thể xác nhận đơn trả hàng nhập (%v)", purchaserefundDB.Code)
	}
	purchaserefundDB.ConfirmedAt = time.Now()
	purchaserefundDB.Status = status3.P
	purchaserefundDB.UpdatedBy = args.UpdatedBy
	err = a.store(ctx).ID(args.ID).ShopID(args.ShopID).UpdatePurchaseRefundAll(purchaserefundDB)
	if err != nil {
		return nil, err
	}
	event := &purchaserefund.ConfirmedPurchaseRefundEvent{
		ShopID:               args.ShopID,
		PurchaseRefundID:     args.ID,
		UpdatedBy:            args.UpdatedBy,
		AutoInventoryVoucher: args.AutoInventoryVoucher,
		InventoryOverStock:   args.InventoryOverStock,
	}
	err = a.eventBus.Publish(ctx, event)
	if err != nil {
		return nil, err
	}
	return purchaserefundDB, nil
}
