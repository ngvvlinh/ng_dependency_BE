package shop

import (
	"context"

	"etop.vn/api/main/purchaseorder"
	pbcm "etop.vn/backend/pb/common"
	pbshop "etop.vn/backend/pb/etop/shop"
	"etop.vn/backend/pkg/common/bus"
	. "etop.vn/capi/dot"
)

func init() {
	bus.AddHandlers("api",
		purchaseOrderService.GetPurchaseOrder,
		purchaseOrderService.GetPurchaseOrders,
		purchaseOrderService.CreatePurchaseOrder,
		purchaseOrderService.UpdatePurchaseOrder,
		purchaseOrderService.DeletePurchaseOrder,
		purchaseOrderService.ConfirmPurchaseOrder,
		purchaseOrderService.CancelPurchaseOrder)
}

func (s *PurchaseOrderService) GetPurchaseOrder(ctx context.Context, r *GetPurchaseOrderEndpoint) error {
	query := &purchaseorder.GetPurchaseOrderByIDQuery{
		ID:     r.Id,
		ShopID: r.Context.Shop.ID,
	}
	if err := purchaseOrderQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = pbshop.PbPurchaseOrder(query.Result)
	r.Result.InventoryVoucher = PbShopInventoryVoucher(query.Result.InventoryVoucher)
	return nil
}

func (s *PurchaseOrderService) GetPurchaseOrders(ctx context.Context, r *GetPurchaseOrdersEndpoint) error {
	paging := r.Paging.CMPaging()
	query := &purchaseorder.ListPurchaseOrdersQuery{
		ShopID:  r.Context.Shop.ID,
		Paging:  *paging,
		Filters: pbcm.ToFilters(r.Filters),
	}
	if err := purchaseOrderQuery.Dispatch(ctx, query); err != nil {
		return err
	}

	var purchaseOrders []*pbshop.PurchaseOrder
	for _, purchaseOrder := range query.Result.PurchaseOrders {
		purchaseOrderTemp := pbshop.PbPurchaseOrder(purchaseOrder)
		purchaseOrderTemp.InventoryVoucher = PbShopInventoryVoucher(purchaseOrder.InventoryVoucher)
		purchaseOrders = append(purchaseOrders, purchaseOrderTemp)
	}

	r.Result = &pbshop.PurchaseOrdersResponse{
		PurchaseOrders: purchaseOrders,
		Paging:         pbcm.PbPageInfo(paging, query.Result.Count),
	}
	return nil
}

func (s *PurchaseOrderService) CreatePurchaseOrder(ctx context.Context, r *CreatePurchaseOrderEndpoint) error {
	cmd := &purchaseorder.CreatePurchaseOrderCommand{
		ShopID:        r.Context.Shop.ID,
		SupplierID:    r.SupplierId,
		BasketValue:   r.BasketValue,
		TotalDiscount: r.TotalDiscount,
		TotalAmount:   r.TotalAmount,
		Note:          r.Note,
		Lines:         pbshop.Convert_api_PurchaseOrderLines_To_core_PurchaseOrderLines(r.Lines),
		CreatedBy:     r.Context.UserID,
	}
	if err := purchaseOrderAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = pbshop.PbPurchaseOrder(cmd.Result)
	return nil
}

func (s *PurchaseOrderService) UpdatePurchaseOrder(ctx context.Context, r *UpdatePurchaseOrderEndpoint) error {
	cmd := &purchaseorder.UpdatePurchaseOrderCommand{
		ID:            r.Id,
		ShopID:        r.Context.Shop.ID,
		SupplierID:    PInt64(r.SupplierId),
		BasketValue:   PInt64(r.BasketValue),
		TotalDiscount: PInt64(r.TotalDiscount),
		TotalAmount:   PInt64(r.TotalAmount),
		Note:          PString(r.Note),
		Lines:         pbshop.Convert_api_PurchaseOrderLines_To_core_PurchaseOrderLines(r.Lines),
	}
	if err := purchaseOrderAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = pbshop.PbPurchaseOrder(cmd.Result)
	return nil
}

func (s *PurchaseOrderService) DeletePurchaseOrder(ctx context.Context, r *DeletePurchaseOrderEndpoint) error {
	cmd := &purchaseorder.DeletePurchaseOrderCommand{
		ID:     r.Id,
		ShopID: r.Context.Shop.ID,
	}
	if err := purchaseOrderAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	return nil
}

func (s *PurchaseOrderService) ConfirmPurchaseOrder(ctx context.Context, r *ConfirmPurchaseOrderEndpoint) error {
	cmd := &purchaseorder.ConfirmPurchaseOrderCommand{
		ID:                   r.Id,
		AutoInventoryVoucher: purchaseorder.PurchaseOrderAutoInventoryVoucher(r.AutoInventoryVoucher),
		ShopID:               r.Context.Shop.ID,
	}
	if err := purchaseOrderAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.UpdatedResponse{Updated: int32(cmd.Result)}
	return nil
}

func (s *PurchaseOrderService) CancelPurchaseOrder(ctx context.Context, r *CancelPurchaseOrderEndpoint) error {
	cmd := &purchaseorder.CancelPurchaseOrderCommand{
		ID:     r.Id,
		ShopID: r.Context.Shop.ID,
		Reason: r.Reason,
	}
	if err := purchaseOrderAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.UpdatedResponse{Updated: int32(cmd.Result)}
	return nil
}
