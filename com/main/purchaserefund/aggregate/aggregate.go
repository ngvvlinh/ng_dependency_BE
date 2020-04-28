package aggregate

import (
	"context"
	"time"

	"o.o/api/main/purchaseorder"
	"o.o/api/main/purchaserefund"
	"o.o/api/top/types/etc/status3"
	catalogconvert "o.o/backend/com/main/catalog/convert"
	"o.o/backend/com/main/purchaserefund/convert"
	"o.o/backend/com/main/purchaserefund/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi"
	"o.o/capi/dot"
	"o.o/common/l"
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
	err = checkPurchaseRefund(&purchaseRefundResult)
	if err != nil {
		return nil, err
	}
	preLine, err := a.checkLinePurchaseOrder(ctx, args.ShopID, purchaseRefundResult.PurchaseOrderID, purchaseRefundResult.ID, purchaseRefundResult.Lines)
	if err != nil {
		return nil, err
	}
	purchaseRefundResult.SupplierID = preLine.SupplierID
	purchaseRefundResult.Lines = preLine.PurchaseRefundLine
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
	err = checkPurchaseRefund(purchaserefundDB)
	if err != nil {
		return nil, err
	}
	err = a.store(ctx).ShopID(args.ShopID).ID(args.ID).UpdatePurchaseRefundAll(purchaserefundDB)
	return purchaserefundDB, err
}

func checkPurchaseRefund(args *purchaserefund.PurchaseRefund) error {
	if args.BasketValue <= 0 {
		return cm.Errorf(cm.NotFound, nil, "basket_value không thể nhỏ hơn hoặc bằng 0")
	}
	basketValueLine := 0
	for _, value := range args.Lines {
		basketValueLine += (value.Adjustment + value.PaymentPrice) * value.Quantity
	}
	if basketValueLine != args.BasketValue {
		return cm.Errorf(cm.NotFound, nil, "basket_value không không hợp lệ.")
	}

	totalAdjustmentLine := 0
	for _, value := range args.AdjustmentLines {
		totalAdjustmentLine += value.Amount
	}
	if totalAdjustmentLine != args.TotalAdjustment {
		return cm.Errorf(cm.NotFound, nil, "total_adjustment không không hợp lệ.")
	}
	if args.BasketValue+args.TotalAdjustment != args.TotalAmount {
		return cm.Errorf(cm.NotFound, nil, "total_amount không không hợp lệ. total_amount = basket_value + total_adjustment")
	}
	return nil
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
		if lines[key].PaymentPrice != linesVariant[value.VariantID].PaymentPrice {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Giá sản phẩm %v trong request %v khác giá sản phẩm trong Purchase Order %v", lines[key].Code, lines[key].PaymentPrice, linesVariant[value.VariantID].PaymentPrice)
		}
		lines[key].Code = linesVariant[value.VariantID].Code
		lines[key].ImageURL = linesVariant[value.VariantID].ImageUrl
		lines[key].ProductName = linesVariant[value.VariantID].ProductName
		lines[key].ProductID = linesVariant[value.VariantID].ProductID
		lines[key].Attributes = linesVariant[value.VariantID].Attributes
		basketValue = basketValue + (lines[key].PaymentPrice+lines[key].Adjustment)*lines[key].Quantity
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
	err = a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		err = a.store(ctx).ID(args.ID).ShopID(args.ShopID).UpdatePurchaseRefundAll(purchaserefundDB)
		if err != nil {
			return err
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
			return err
		}
		return nil
	})

	return purchaserefundDB, nil
}
