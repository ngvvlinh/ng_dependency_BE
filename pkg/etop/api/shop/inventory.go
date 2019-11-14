package shop

import (
	"context"
	"fmt"

	"etop.vn/api/main/etop"
	"etop.vn/api/main/inventory"
	"etop.vn/api/meta"
	ordermodelx "etop.vn/backend/com/main/ordering/modelx"
	pbs4 "etop.vn/backend/pb/etop/etc/status4"
	pbshop "etop.vn/backend/pb/etop/shop"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/model"
	. "etop.vn/capi/dot"
)

func (s *InventoryService) CreateInventoryVoucher(ctx context.Context, q *CreateInventoryVoucherEndpoint) error {
	cmd, err := PreCreateInventoryVoucher(ctx, q)
	if err != nil {
		return err
	}

	if err := inventoryAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbshop.CreateInventoryVoucherResponse{
		Inventory: PbShopInventoryVoucher(cmd.Result),
	}
	return nil
}

func (s *InventoryService) ConfirmInventoryVoucher(ctx context.Context, q *ConfirmInventoryVoucherEndpoint) error {
	shopID := q.Context.Shop.ID
	userID := q.Context.UserID

	cmd := &inventory.ConfirmInventoryVoucherCommand{
		ShopID:    shopID,
		ID:        q.Id,
		UpdatedBy: userID,
	}
	if err := inventoryAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbshop.ConfirmInventoryVoucherResponse{
		InventoryVoucher: PbShopInventoryVoucher(cmd.Result),
	}
	return nil
}

func (s *InventoryService) CancelInventoryVoucher(ctx context.Context, q *CancelInventoryVoucherEndpoint) error {
	shopID := q.Context.Shop.ID
	userID := q.Context.UserID

	cmd := &inventory.CancelInventoryVoucherCommand{
		ShopID:    shopID,
		ID:        q.Id,
		UpdatedBy: userID,
		Reason:    q.Reason,
	}
	if err := inventoryAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbshop.CancelInventoryVoucherResponse{
		Inventory: PbShopInventoryVoucher(cmd.Result),
	}
	return nil
}

func (s *InventoryService) UpdateInventoryVoucher(ctx context.Context, q *UpdateInventoryVoucherEndpoint) error {
	shopID := q.Context.Shop.ID
	userID := q.Context.UserID
	var items []*inventory.InventoryVoucherItem
	for _, value := range q.Lines {
		items = append(items, &inventory.InventoryVoucherItem{
			VariantID: value.VariantId,
			Price:     value.Price,
			Quantity:  value.Quantity,
		})
	}
	cmd := &inventory.UpdateInventoryVoucherCommand{
		Title:       PString(q.Title),
		ID:          q.Id,
		ShopID:      shopID,
		TotalAmount: q.TotalAmount,
		UpdatedBy:   userID,
		TraderID:    PInt64(q.TraderId),
		Note:        PString(q.Note),
		Lines:       items,
	}
	if err := inventoryAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbshop.UpdateInventoryVoucherResponse{
		InventoryVoucher: PbShopInventoryVoucher(cmd.Result),
	}
	return nil
}

func (s *InventoryService) AdjustInventoryQuantity(ctx context.Context, q *AdjustInventoryQuantityEndpoint) error {
	shopID := q.Context.Shop.ID
	userID := q.Context.UserID
	inventoryOverstock := q.Context.Shop.InventoryOverstock
	var items []*inventory.InventoryVariant
	for _, value := range q.InventoryVariants {
		items = append(items, &inventory.InventoryVariant{
			ShopID:          shopID,
			VariantID:       value.VariantId,
			QuantityOnHand:  value.QuantityOnHand,
			QuantitySummary: value.Quantity,
			QuantityPicked:  value.QuantityPicked,
		})
	}
	cmd := &inventory.AdjustInventoryQuantityCommand{
		Overstock: cm.BoolDefault(inventoryOverstock, true),
		ShopID:    shopID,
		Lines:     items,
		UserID:    userID,
		Note:      q.Note,
	}
	if err := inventoryAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbshop.AdjustInventoryQuantityResponse{
		InventoryVariants: PbInventoryVariants(cmd.Result.InventoryVariants),
		InventoryVouchers: PbShopInventoryVouchers(cmd.Result.InventoryVouchers),
	}
	return nil
}

func (s *InventoryService) GetInventoryVariants(ctx context.Context, q *GetInventoryVariantsEndpoint) error {
	shopID := q.Context.Shop.ID
	query := &inventory.GetInventoryVariantsQuery{
		ShopID: shopID,
		Paging: &meta.Paging{
			Offset: q.Paging.Offset,
			Limit:  q.Paging.Limit,
		},
	}
	if err := inventoryQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pbshop.GetInventoryVariantsResponse{
		InventoryVariants: PbInventoryVariants(query.Result.InventoryVariants),
	}
	return nil
}

func (s *InventoryService) GetInventoryVouchersByReference(ctx context.Context, q *GetInventoryVouchersByReferenceEndpoint) error {
	shopID := q.Context.Shop.ID
	query := &inventory.GetInventoryVoucherByReferenceQuery{
		ShopID:  shopID,
		RefID:   q.RefId,
		RefType: inventory.InventoryRefType(q.RefType),
	}
	if err := inventoryQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pbshop.GetInventoryVouchersByReferenceResponse{
		InventoryVoucher: PbShopInventoryVouchers(query.Result.InventoryVouchers),
		Status:           pbs4.Pb(model.Status4(query.Result.Status)),
	}
	return nil
}

func (s *InventoryService) GetInventoryVariant(ctx context.Context, q *GetInventoryVariantEndpoint) error {
	shopID := q.Context.Shop.ID
	query := &inventory.GetInventoryVariantQuery{
		ShopID:    shopID,
		VariantID: q.VariantId,
	}
	if err := inventoryQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = PbInventory(query.Result)
	return nil
}

func (s *InventoryService) GetInventoryVariantsByVariantIDs(ctx context.Context, q *GetInventoryVariantsByVariantIDsEndpoint) error {
	shopID := q.Context.Shop.ID
	query := &inventory.GetInventoryVariantsByVariantIDsQuery{
		ShopID:     shopID,
		VariantIDs: q.VariantIds,
	}
	if err := inventoryQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pbshop.GetInventoryVariantsResponse{
		InventoryVariants: PbInventoryVariants(query.Result.InventoryVariants),
	}
	return nil
}

func (s *InventoryService) GetInventoryVoucher(ctx context.Context, q *GetInventoryVoucherEndpoint) error {
	shopID := q.Context.Shop.ID
	query := &inventory.GetInventoryVoucherQuery{
		ShopID: shopID,
		ID:     q.Id,
	}
	if err := inventoryQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = PbShopInventoryVoucher(query.Result)
	return nil
}

func (s *InventoryService) GetInventoryVouchers(ctx context.Context, q *GetInventoryVouchersEndpoint) error {
	shopID := q.Context.Shop.ID
	query := &inventory.GetInventoryVouchersQuery{
		ShopID: shopID,
		Paging: &meta.Paging{
			Offset: q.Paging.Offset,
			Limit:  q.Paging.Limit,
		},
	}
	if err := inventoryQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pbshop.GetInventoryVouchersResponse{
		InventoryVouchers: PbShopInventoryVouchers(query.Result.InventoryVoucher),
	}
	return nil
}

func (s *InventoryService) GetInventoryVouchersByIDs(ctx context.Context, q *GetInventoryVouchersByIDsEndpoint) error {
	shopID := q.Context.Shop.ID
	query := &inventory.GetInventoryVouchersByIDsQuery{
		ShopID: shopID,
		IDs:    q.Ids,
	}
	if err := inventoryQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pbshop.GetInventoryVouchersResponse{
		InventoryVouchers: PbShopInventoryVouchers(query.Result.InventoryVoucher),
	}
	return nil
}

func PreCreateInventoryVoucherRefOrder(ctx context.Context, cmd *inventory.CreateInventoryVoucherCommand) error {
	var items []*inventory.InventoryVoucherItem

	// check order_id exit
	queryOrder := &ordermodelx.GetOrderQuery{
		OrderID: cmd.RefID,
		ShopID:  cmd.ShopID,
	}
	if err := bus.Dispatch(ctx, queryOrder); err != nil {
		return err
	}
	// check voucher already create
	queryInventoryVoucher := &inventory.GetInventoryVoucherByReferenceQuery{
		ShopID:  cmd.ShopID,
		RefID:   cmd.RefID,
		RefType: inventory.RefTypeOrder,
		Result:  nil,
	}
	err := inventoryQuery.Dispatch(ctx, queryInventoryVoucher)
	if err != nil {
		return err
	}
	for _, value := range queryInventoryVoucher.Result.InventoryVouchers {
		if value.Status == etop.S3Positive || value.Status == etop.S3Zero {
			return cm.Errorf(cm.InvalidArgument, nil, "Order đã có phiếu xuất kho inventory_voucher_id = %v, Vui lòng kiểm tra lại ", value.ID)
		}
	}

	// GET info and put it to cmd
	for _, value := range queryOrder.Result.Order.Lines {
		items = append(items, &inventory.InventoryVoucherItem{
			VariantID: value.VariantID,
			Quantity:  int32(value.Quantity),
		})
	}
	cmd.Title = "Xuất kho khi bán hàng"
	cmd.Lines = items
	cmd.Type = "out"
	cmd.TraderID = queryOrder.Result.Order.CustomerID
	cmd.Note = fmt.Sprintf("Tạo phiếu nhập kho theo đơn đặt hàng mã %v", queryOrder.Result.Order.ID)
	return nil
}

func PreCreateInventoryVoucher(ctx context.Context, q *CreateInventoryVoucherEndpoint) (*inventory.CreateInventoryVoucherCommand, error) {
	shopID := q.Context.Shop.ID
	inventoryOverstock := q.Context.Shop.InventoryOverstock

	// default when not have ref_id, ref_type
	var items []*inventory.InventoryVoucherItem
	for _, value := range q.Lines {
		items = append(items, &inventory.InventoryVoucherItem{
			VariantID: value.VariantId,
			Price:     value.Price,
			Quantity:  value.Quantity,
		})
	}

	cmd := &inventory.CreateInventoryVoucherCommand{
		Title:       q.Title,
		ShopID:      shopID,
		Overstock:   cm.BoolDefault(inventoryOverstock, true),
		RefType:     inventory.InventoryRefType(q.RefType),
		RefID:       q.RefId,
		TotalAmount: q.TotalAmount,
		CreatedBy:   q.Context.UserID,
		TraderID:    q.TraderId,
		Type:        inventory.InventoryVoucherType(q.Type),
		Note:        q.Note,
		Lines:       items,
	}
	// Check ref ID here "order", "purchaseorder", "stocktake", "return"
	// modify value flow reftype
	switch cmd.RefType {
	case inventory.RefTypeOrder:
		if err := PreCreateInventoryVoucherRefOrder(ctx, cmd); err != nil {
			return nil, err
		}
	case inventory.RefTypeReturns:
	case inventory.RefTypePurchaseOrder:
	case inventory.RefTypeStockTake:
	default:

	}
	return cmd, nil
}
