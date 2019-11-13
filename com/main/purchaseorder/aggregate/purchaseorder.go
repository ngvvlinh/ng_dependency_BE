package aggregate

import (
	"context"

	"etop.vn/api/main/catalog"
	"etop.vn/api/main/etop"
	"etop.vn/api/main/inventory"
	"etop.vn/api/main/purchaseorder"
	"etop.vn/api/shopping/suppliering"
	"etop.vn/backend/com/main/purchaseorder/convert"
	"etop.vn/backend/com/main/purchaseorder/model"
	"etop.vn/backend/com/main/purchaseorder/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/conversion"
	"etop.vn/capi"
	"etop.vn/common/l"
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
	database *cmsql.Database, eventBus capi.EventBus,
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

func (a *PurchaseOrderAggregate) MessageBus() purchaseorder.CommandBus {
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
			Wrap(cm.NotFound, "Không tìm thấy nhà phân phối").
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

	var maxCodeNorm int32
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
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Vui lòng nhập mã")
	}
	codeNorm := maxCodeNorm + 1
	purchaseOrder.Code = convert.GenerateCode(int(codeNorm))
	purchaseOrder.CodeNorm = codeNorm

	if err := a.store(ctx).CreatePurchaseOrder(purchaseOrder); err != nil {
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
			Wrap(cm.NotFound, "Không tìm thấy đơn nhập hàng.").
			Throw()
	}

	if purchaseOrder.Status != etop.S3Zero {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Không thể chỉnh sửa đơn nhập hàng, kiểm tra trạng thái đơn.")
	}

	if args.SupplierID.Valid && args.SupplierID.Int64 != purchaseOrder.SupplierID {
		// check supplier_id
		getSupplier := &suppliering.GetSupplierByIDQuery{
			ID:     args.SupplierID.Int64,
			ShopID: args.ShopID,
		}
		if err := a.supplierQuery.Dispatch(ctx, getSupplier); err != nil {
			return nil, cm.MapError(err).
				Wrap(cm.NotFound, "Không tìm thấy nhà phân phối").
				Throw()
		}

		purchaseOrder.Supplier = &purchaseorder.PurchaseOrderSupplier{
			FullName:           getSupplier.Result.FullName,
			Phone:              getSupplier.Result.Phone,
			Email:              getSupplier.Result.Email,
			CompanyName:        getSupplier.Result.CompanyName,
			TaxNumber:          getSupplier.Result.TaxNumber,
			HeadquarterAddress: getSupplier.Result.HeadquaterAddress,
		}
	}

	if err := scheme.Convert(args, purchaseOrder); err != nil {
		return nil, err
	}

	if err := a.checkPurchaseOrder(ctx, purchaseOrder); err != nil {
		return nil, err
	}

	purchaseOrderDB := new(model.PurchaseOrder)
	if err := scheme.Convert(purchaseOrder, purchaseOrderDB); err != nil {
		return nil, err
	}

	err = a.store(ctx).UpdatePurchaseOrderDB(purchaseOrderDB)
	return purchaseOrder, err
}

func (a *PurchaseOrderAggregate) checkPurchaseOrder(
	ctx context.Context, args *purchaseorder.PurchaseOrder,
) error {
	if args.BasketValue <= 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Tiền hàng phải lớn hơn 0")
	}
	if args.TotalDiscount < 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Giảm giá phải lớn hơn hoặc bằng 0")
	}
	if args.TotalAmount < 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Tiền thanh toán phải lớn hơn hoặc bằng 0")
	}
	if (args.BasketValue-args.TotalDiscount > 0 && args.BasketValue-args.TotalDiscount != args.TotalAmount) ||
		(args.BasketValue-args.TotalDiscount <= 0 && args.TotalAmount != 0) {
		return cm.Errorf(cm.InvalidArgument, nil, "Tiền thanh toán không hợp lệ")
	}

	var variantIDs []int64
	var totalPrice int64
	mapVariant := make(map[int64]*catalog.ShopVariant)
	for _, line := range args.Lines {
		if line.VariantID == 0 {
			return cm.Errorf(cm.NotFound, nil, "variant_id không thể bằng 0")
		}
		if line.Price <= 0 {
			return cm.Errorf(cm.InvalidArgument, nil, "gía của phiên bản sản phẩm không hợp lệ")
		}
		if line.Quantity <= 0 {
			return cm.Errorf(cm.InvalidArgument, nil, "số lượng của phiên bản sản phẩm không hợp lệ")
		}
		variantIDs = append(variantIDs, line.VariantID)
		totalPrice += line.Quantity * line.Price
	}
	if totalPrice != args.BasketValue {
		return cm.Errorf(cm.NotFound, nil, "Tiền hàng không hợp lệ")
	}

	query := &catalog.ListShopVariantsByIDsQuery{
		IDs:    variantIDs,
		ShopID: args.ShopID,
		Result: nil,
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
				return cm.Errorf(cm.FailedPrecondition, nil, "variant_id %v không tồn tại", variantID)
			}
		}
	}

	return nil
}

func (a *PurchaseOrderAggregate) CancelPurchaseOrder(
	ctx context.Context, args *purchaseorder.CancelPurchaseOrderArgs,
) (updated int, err error) {
	purchaseOrder, err := a.store(ctx).ID(args.ID).ShopID(args.ShopID).GetPurchaseOrder()
	if err != nil {
		return 0, cm.MapError(err).
			Wrap(cm.NotFound, "Không tìm thấy đơn nhập hàng.").
			Throw()
	}
	if purchaseOrder.Status != etop.S3Zero {
		return 0, cm.Errorf(cm.FailedPrecondition, nil, "Không thể huỷ đơn nhập hàng.")
	}

	updated, err = a.store(ctx).ID(args.ID).ShopID(args.ShopID).CancelPurchaseOrder(args.Reason)
	return updated, err
}

func (a *PurchaseOrderAggregate) ConfirmPurchaseOrder(
	ctx context.Context, args *purchaseorder.ConfirmPurchaseOrderArgs,
) (updated int, err error) {
	if args.AutoInventoryVoucher != "" &&
		args.AutoInventoryVoucher != purchaseorder.AutoInventoryVoucherCreate &&
		args.AutoInventoryVoucher != purchaseorder.AutoInventoryVoucherConfirm {
		return 0, cm.Errorf(cm.InvalidArgument, nil, "AutoInventoryVoucher không hợp lệ")
	}

	purchaseOrder, err := a.store(ctx).ID(args.ID).ShopID(args.ShopID).GetPurchaseOrder()
	if err != nil {
		return 0, cm.MapError(err).
			Wrap(cm.NotFound, "Không tìm thấy đơn nhập hàng.").
			Throw()
	}
	if purchaseOrder.Status != etop.S3Zero {
		return 0, cm.Errorf(cm.FailedPrecondition, nil, "Không thể xác nhận đơn nhập hàng.")
	}

	var lines []*inventory.InventoryVoucherItem
	for _, line := range purchaseOrder.Lines {
		lines = append(lines, &inventory.InventoryVoucherItem{
			VariantID: line.VariantID,
			Price:     int32(line.Price),
			Quantity:  int32(line.Quantity),
		})
	}

	err = a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		updated, err = a.store(ctx).ID(args.ID).ShopID(args.ShopID).ConfirmPurchaseOrder()
		if err != nil {
			return err
		}

		event := &purchaseorder.PurchaseOrderConfirmedEvent{
			ShopID:               args.ShopID,
			UserID:               purchaseOrder.CreatedBy,
			PurchaseOrderID:      args.ID,
			TraderID:             purchaseOrder.SupplierID,
			TotalAmount:          purchaseOrder.BasketValue,
			AutoInventoryVoucher: args.AutoInventoryVoucher,
			Lines:                lines,
		}
		if err := a.eventBus.Publish(ctx, event); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return updated, nil
}

func (a *PurchaseOrderAggregate) DeletePurchaseOrder(ctx context.Context, ID, shopID int64) (deleted int, _ error) {
	if _, err := a.store(ctx).ID(ID).ShopID(shopID).GetPurchaseOrder(); err != nil {
		return 0, cm.MapError(err).
			Wrap(cm.NotFound, "Không tìm thấy đơn nhập hàng.").
			Throw()
	}

	deleted, err := a.store(ctx).ID(ID).ShopID(shopID).SoftDelete()
	return deleted, err
}
