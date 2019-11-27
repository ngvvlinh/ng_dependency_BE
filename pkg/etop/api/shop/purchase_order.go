package shop

import (
	"context"

	"etop.vn/api/main/inventory"
	"etop.vn/api/main/purchaseorder"
	pbcm "etop.vn/api/pb/common"
	pbshop "etop.vn/api/pb/etop/shop"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmapi"
	"etop.vn/backend/pkg/etop/api/convertpb"
)

func init() {
	bus.AddHandlers("api",
		purchaseOrderService.GetPurchaseOrder,
		purchaseOrderService.GetPurchaseOrders,
		purchaseOrderService.GetPurchaseOrdersByIDs,
		purchaseOrderService.GetPurchaseOrdersByReceiptID,
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
	r.Result = convertpb.PbPurchaseOrder(query.Result)
	r.Result.InventoryVoucher = PbShopInventoryVoucher(query.Result.InventoryVoucher)
	return nil
}

func (s *PurchaseOrderService) GetPurchaseOrders(ctx context.Context, r *GetPurchaseOrdersEndpoint) error {
	paging := cmapi.CMPaging(r.Paging)
	query := &purchaseorder.ListPurchaseOrdersQuery{
		ShopID:  r.Context.Shop.ID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(r.Filters),
	}
	if err := purchaseOrderQuery.Dispatch(ctx, query); err != nil {
		return err
	}

	var purchaseOrders []*pbshop.PurchaseOrder
	for _, purchaseOrder := range query.Result.PurchaseOrders {
		purchaseOrderTemp := convertpb.PbPurchaseOrder(purchaseOrder)
		purchaseOrderTemp.InventoryVoucher = PbShopInventoryVoucher(purchaseOrder.InventoryVoucher)
		purchaseOrders = append(purchaseOrders, purchaseOrderTemp)
	}

	r.Result = &pbshop.PurchaseOrdersResponse{
		PurchaseOrders: purchaseOrders,
		Paging:         cmapi.PbPageInfo(paging, query.Result.Count),
	}
	return nil
}

func (s *PurchaseOrderService) GetPurchaseOrdersByIDs(ctx context.Context, r *GetPurchaseOrdersByIDsEndpoint) error {
	query := &purchaseorder.GetPurchaseOrdersByIDsQuery{
		IDs:    r.Ids,
		ShopID: r.Context.Shop.ID,
	}
	if err := purchaseOrderQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = &pbshop.PurchaseOrdersResponse{PurchaseOrders: convertpb.PbPurchaseOrders(query.Result.PurchaseOrders)}
	return nil
}

func (s *PurchaseOrderService) GetPurchaseOrdersByReceiptID(ctx context.Context, r *GetPurchaseOrdersByReceiptIDEndpoint) error {
	query := &purchaseorder.ListPurchaseOrdersByReceiptIDQuery{
		ReceiptID: r.Id,
		ShopID:    r.Context.Shop.ID,
	}
	if err := purchaseOrderQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = &pbshop.PurchaseOrdersResponse{PurchaseOrders: convertpb.PbPurchaseOrders(query.Result.PurchaseOrders)}
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
		Lines:         convertpb.Convert_api_PurchaseOrderLines_To_core_PurchaseOrderLines(r.Lines),
		CreatedBy:     r.Context.UserID,
	}
	if err := purchaseOrderAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = convertpb.PbPurchaseOrder(cmd.Result)
	return nil
}

func (s *PurchaseOrderService) UpdatePurchaseOrder(ctx context.Context, r *UpdatePurchaseOrderEndpoint) error {
	cmd := &purchaseorder.UpdatePurchaseOrderCommand{
		ID:            r.Id,
		ShopID:        r.Context.Shop.ID,
		SupplierID:    r.SupplierId,
		BasketValue:   r.BasketValue,
		TotalDiscount: r.TotalDiscount,
		TotalAmount:   r.TotalAmount,
		Note:          r.Note,
		Lines:         convertpb.Convert_api_PurchaseOrderLines_To_core_PurchaseOrderLines(r.Lines),
	}
	if err := purchaseOrderAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = convertpb.PbPurchaseOrder(cmd.Result)
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
		AutoInventoryVoucher: inventory.AutoInventoryVoucher(r.AutoInventoryVoucher),
		ShopID:               r.Context.Shop.ID,
	}
	if err := purchaseOrderAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.UpdatedResponse{Updated: int(cmd.Result)}
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
	r.Result = &pbcm.UpdatedResponse{Updated: int(cmd.Result)}
	return nil
}
