package aggregate

import (
	"context"

	"o.o/api/main/catalog"
	"o.o/api/main/inventory"
	"o.o/api/main/purchaseorder"
	"o.o/api/shopping/suppliering"
	"o.o/api/top/types/etc/status3"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/purchaseorder/convert"
	"o.o/backend/com/main/purchaseorder/model"
	"o.o/backend/com/main/purchaseorder/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi"
	"o.o/capi/dot"
	"o.o/common/l"
)

var ll = l.New()
var _ purchaseorder.Aggregate = &PurchaseOrderAggregate{}
var scheme = conversion.Build(convert.RegisterConversions)

type PurchaseOrderAggregate struct {
	db             *cmsql.Database
	store          sqlstore.PurchaseOrderStoreFactory
	eventBus       capi.EventBus
	catalogQuery   catalog.QueryBus
	supplierQuery  suppliering.QueryBus
	inventoryQuery inventory.QueryBus
}

func NewPurchaseOrderAggregate(
	database com.MainDB, eventBus capi.EventBus,
	catalogQ catalog.QueryBus, supplierQ suppliering.QueryBus,
	inventoryQ inventory.QueryBus,
) *PurchaseOrderAggregate {
	return &PurchaseOrderAggregate{
		db:             database,
		store:          sqlstore.NewPurchaseOrderStore(database),
		eventBus:       eventBus,
		catalogQuery:   catalogQ,
		supplierQuery:  supplierQ,
		inventoryQuery: inventoryQ,
	}
}

func PurchaseOrderAggregateMessageBus(a *PurchaseOrderAggregate) purchaseorder.CommandBus {
	b := bus.New()
	return purchaseorder.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *PurchaseOrderAggregate) CreatePurchaseOrder(
	ctx context.Context, args *purchaseorder.CreatePurchaseOrderArgs,
) (*purchaseorder.PurchaseOrder, error) {
	purchaseOrderNeedValidate := &purchaseorder.PurchaseOrder{
		ShopID:        args.ShopID,
		BasketValue:   args.BasketValue,
		TotalDiscount: args.TotalDiscount,
		DiscountLines: args.DiscountLines,
		TotalFee:      args.TotalFee,
		FeeLines:      args.FeeLines,
		TotalAmount:   args.TotalAmount,
		Lines:         args.Lines,
	}
	if err := a.checkPurchaseOrder(ctx, purchaseOrderNeedValidate); err != nil {
		return nil, err
	}
	// check supplier_id
	getSupplier := &suppliering.GetSupplierByIDQuery{
		ID:     args.SupplierID,
		ShopID: args.ShopID,
	}
	if err := a.supplierQuery.Dispatch(ctx, getSupplier); err != nil {
		return nil, cm.MapError(err).
			Wrap(cm.NotFound, "Kh??ng t??m th???y nh?? ph??n ph???i").
			Throw()
	}

	purchaseOrder := new(purchaseorder.PurchaseOrder)
	if err := scheme.Convert(args, purchaseOrder); err != nil {
		return nil, err
	}
	purchaseOrder.Supplier = &purchaseorder.PurchaseOrderSupplier{
		FullName:           getSupplier.Result.FullName,
		Phone:              getSupplier.Result.Phone,
		Email:              getSupplier.Result.Email,
		CompanyName:        getSupplier.Result.CompanyName,
		TaxNumber:          getSupplier.Result.TaxNumber,
		HeadquarterAddress: getSupplier.Result.HeadquaterAddress,
	}

	var maxCodeNorm int
	purchaseOrderTemp, err := a.store(ctx).ShopID(args.ShopID).IncludeDeleted().GetReceiptByMaximumCodeNorm()
	switch cm.ErrorCode(err) {
	case cm.NoError:
		maxCodeNorm = purchaseOrderTemp.CodeNorm
	case cm.NotFound:
	// no-op
	default:
		return nil, err
	}

	if maxCodeNorm >= convert.MaxCodeNorm {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Vui l??ng nh???p m??")
	}
	codeNorm := maxCodeNorm + 1
	purchaseOrder.Code = convert.GenerateCode(codeNorm)
	purchaseOrder.CodeNorm = codeNorm
	if err = a.getLinesInPurchaseOrder(ctx, purchaseOrder.Lines, args.ShopID); err != nil {
		return nil, err
	}
	if err = a.store(ctx).CreatePurchaseOrder(purchaseOrder); err != nil {
		return nil, err
	}
	return purchaseOrder, nil
}

func (a *PurchaseOrderAggregate) UpdatePurchaseOrder(
	ctx context.Context, args *purchaseorder.UpdatePurchaseOrderArgs,
) (*purchaseorder.PurchaseOrder, error) {
	purchaseOrder, err := a.store(ctx).ID(args.ID).ShopID(args.ShopID).GetPurchaseOrder()
	if err != nil {
		return nil, cm.MapError(err).
			Wrap(cm.NotFound, "Kh??ng t??m th???y ????n nh???p h??ng.").
			Throw()
	}
	if purchaseOrder.Status != status3.Z {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Kh??ng th??? ch???nh s???a ????n nh???p h??ng, ki???m tra tr???ng th??i ????n.")
	}

	if err = scheme.Convert(args, purchaseOrder); err != nil {
		return nil, err
	}
	if err = a.checkPurchaseOrder(ctx, purchaseOrder); err != nil {
		return nil, err
	}
	if err = a.getLinesInPurchaseOrder(ctx, purchaseOrder.Lines, args.ShopID); err != nil {
		return nil, err
	}
	purchaseOrderDB := new(model.PurchaseOrder)
	if err = scheme.Convert(purchaseOrder, purchaseOrderDB); err != nil {
		return nil, err
	}
	err = a.store(ctx).UpdatePurchaseOrderDB(purchaseOrderDB)
	return purchaseOrder, err
}

func (a *PurchaseOrderAggregate) getLinesInPurchaseOrder(ctx context.Context, lines []*purchaseorder.PurchaseOrderLine, ShopID dot.ID) error {
	var variantIDs []dot.ID
	var productIDs []dot.ID
	mapVariantShopVariant := make(map[dot.ID]*catalog.ShopVariant)
	mapProductShopProduct := make(map[dot.ID]*catalog.ShopProduct)
	for _, ln := range lines {
		variantIDs = append(variantIDs, ln.VariantID)
	}
	query := &catalog.ListShopVariantsByIDsQuery{
		IDs:    variantIDs,
		ShopID: ShopID,
	}
	if err := a.catalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	for _, variant := range query.Result.Variants {
		mapVariantShopVariant[variant.VariantID] = variant
		productIDs = append(productIDs, variant.ProductID)
	}
	q := &catalog.ListShopProductsByIDsQuery{
		IDs:    productIDs,
		ShopID: ShopID,
	}
	if err := a.catalogQuery.Dispatch(ctx, q); err != nil {
		return err
	}
	for _, product := range q.Result.Products {
		mapProductShopProduct[product.ProductID] = product
	}
	for _, ln := range lines {
		variant := mapVariantShopVariant[ln.VariantID]
		product := mapProductShopProduct[variant.ProductID]
		fillPOLineInfo(ln, variant, product)
	}
	return nil
}

func fillPOLineInfo(line *purchaseorder.PurchaseOrderLine, variant *catalog.ShopVariant, product *catalog.ShopProduct) {
	line.VariantID = variant.VariantID
	line.ProductID = product.ProductID
	if line.ProductName == "" {
		line.ProductName = product.Name
	}
	line.Code = variant.Code
	if line.ImageUrl == "" {
		line.ImageUrl = getImage(variant, product)
	}
	line.Attributes = variant.Attributes
}

func getImage(variant *catalog.ShopVariant, product *catalog.ShopProduct) string {
	if len(variant.ImageURLs) > 0 {
		return variant.ImageURLs[0]
	}
	if len(product.ImageURLs) > 0 {
		return product.ImageURLs[0]
	}
	return ""
}

func (a *PurchaseOrderAggregate) checkPurchaseOrder(
	ctx context.Context, args *purchaseorder.PurchaseOrder,
) error {
	if args.BasketValue < 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Ti???n h??ng ph???i l???n h??n ho???c b???ng 0")
	}
	if args.TotalDiscount < 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Gi???m gi?? ph???i l???n h??n ho???c b???ng 0")
	}
	if args.TotalAmount < 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Ti???n thanh to??n ph???i l???n h??n ho???c b???ng 0")
	}
	if args.TotalFee < 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Ti???n ph?? ph???i l???n h??n ho???c b???ng 0")
	}
	if args.BasketValue-args.TotalDiscount+args.TotalFee != args.TotalAmount {
		return cm.Errorf(cm.InvalidArgument, nil, "Ti???n thanh to??n kh??ng h???p l???")
	}

	var variantIDs []dot.ID
	var totalLinesValue int
	mapVariant := make(map[dot.ID]*catalog.ShopVariant)
	for _, line := range args.Lines {
		if line.VariantID == 0 {
			return cm.Errorf(cm.NotFound, nil, "variant_id kh??ng th??? b???ng 0")
		}
		if line.PaymentPrice < 0 {
			return cm.Errorf(cm.InvalidArgument, nil, "g??a c???a phi??n b???n s???n ph???m kh??ng h???p l???")
		}
		if line.Quantity <= 0 {
			return cm.Errorf(cm.InvalidArgument, nil, "s??? l?????ng c???a phi??n b???n s???n ph???m kh??ng h???p l???")
		}
		variantIDs = append(variantIDs, line.VariantID)
		totalLinesValue += line.Quantity * (line.PaymentPrice - line.Discount)
	}
	if totalLinesValue != args.BasketValue {
		return cm.Errorf(cm.NotFound, nil, "Ti???n h??ng kh??ng h???p l???")
	}

	totalFee := 0
	for _, value := range args.FeeLines {
		totalFee += value.Amount
	}
	if totalFee != args.TotalFee {
		return cm.Errorf(cm.NotFound, nil, "T???ng gi?? tr??? ph?? kh??ng h???p l???")
	}

	totalDiscount := 0
	for _, value := range args.DiscountLines {
		totalDiscount += value.Amount
	}
	if totalDiscount != args.TotalDiscount {
		return cm.Errorf(cm.NotFound, nil, "T???ng gi?? tr??? gi???m gi?? kh??ng h???p l???")
	}

	query := &catalog.ListShopVariantsByIDsQuery{
		IDs:    variantIDs,
		ShopID: args.ShopID,
	}
	if err := a.catalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	for _, variant := range query.Result.Variants {
		mapVariant[variant.VariantID] = variant
	}
	if len(variantIDs) != len(query.Result.Variants) {
		for _, variantID := range variantIDs {
			if _, ok := mapVariant[variantID]; !ok {
				return cm.Errorf(cm.FailedPrecondition, nil, "Phi??n b???n c???a s???n ph???m kh??ng c??n t???n t???i.")
			}
		}
	}

	return nil
}

func (a *PurchaseOrderAggregate) CancelPurchaseOrder(
	ctx context.Context, args *purchaseorder.CancelPurchaseOrderArgs,
) (updated int, err error) {
	_, err = a.store(ctx).ID(args.ID).ShopID(args.ShopID).GetPurchaseOrder()
	if err != nil {
		return 0, cm.MapError(err).
			Wrap(cm.NotFound, "Kh??ng t??m th???y ????n nh???p h??ng.").
			Throw()
	}
	err = a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		updated, err = a.store(ctx).ID(args.ID).ShopID(args.ShopID).CancelPurchaseOrder(args.CancelReason)
		event := &purchaseorder.PurchaseOrderCancelledEvent{
			PurchaseOrderID:      args.ID,
			ShopID:               args.ShopID,
			UpdatedBy:            args.UpdatedBy,
			AutoInventoryVoucher: args.AutoInventoryVoucher,
			InventoryOverStock:   args.InventoryOverStock,
		}
		err = a.eventBus.Publish(ctx, event)
		return err
	})
	return updated, err
}

func (a *PurchaseOrderAggregate) ConfirmPurchaseOrder(
	ctx context.Context, args *purchaseorder.ConfirmPurchaseOrderArgs,
) (updated int, err error) {

	purchaseOrder, err := a.store(ctx).ID(args.ID).ShopID(args.ShopID).GetPurchaseOrder()
	if err != nil {
		return 0, cm.MapError(err).
			Wrap(cm.NotFound, "Kh??ng t??m th???y ????n nh???p h??ng.").
			Throw()
	}
	if purchaseOrder.Status != status3.Z {
		return 0, cm.Errorf(cm.FailedPrecondition, nil, "Kh??ng th??? x??c nh???n ????n nh???p h??ng.")
	}

	var lines []*inventory.InventoryVoucherItem
	for _, line := range purchaseOrder.Lines {
		lines = append(lines, &inventory.InventoryVoucherItem{
			VariantID: line.VariantID,
			Price:     line.PaymentPrice,
			Quantity:  line.Quantity,
		})
	}

	err = a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		updated, err = a.store(ctx).ID(args.ID).ShopID(args.ShopID).ConfirmPurchaseOrder()
		if err != nil {
			return err
		}
		event := &purchaseorder.PurchaseOrderConfirmedEvent{
			ShopID:               args.ShopID,
			PurchaseOrderCode:    purchaseOrder.Code,
			UserID:               purchaseOrder.CreatedBy,
			PurchaseOrderID:      args.ID,
			TraderID:             purchaseOrder.SupplierID,
			TotalAmount:          purchaseOrder.BasketValue,
			AutoInventoryVoucher: args.AutoInventoryVoucher,
			Lines:                lines,
		}
		if err = a.eventBus.Publish(ctx, event); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	return updated, nil
}

func (a *PurchaseOrderAggregate) DeletePurchaseOrder(ctx context.Context, ID, shopID dot.ID) (deleted int, _ error) {
	if _, err := a.store(ctx).ID(ID).ShopID(shopID).GetPurchaseOrder(); err != nil {
		return 0, cm.MapError(err).
			Wrap(cm.NotFound, "Kh??ng t??m th???y ????n nh???p h??ng.").
			Throw()
	}

	deleted, err := a.store(ctx).ID(ID).ShopID(shopID).SoftDelete()
	return deleted, err
}
