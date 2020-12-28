package purchase_order

import (
	"context"

	"o.o/api/main/purchaseorder"
	"o.o/api/top/int/shop"
	api "o.o/api/top/int/shop"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	convertpball "o.o/backend/pkg/etop/api/convertpb/_all"
	"o.o/backend/pkg/etop/api/shop/inventory"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/capi/util"
)

type PurchaseOrderService struct {
	session.Session

	PurchaseOrderAggr  purchaseorder.CommandBus
	PurchaseOrderQuery purchaseorder.QueryBus
}

func (s *PurchaseOrderService) Clone() api.PurchaseOrderService { res := *s; return &res }

func (s *PurchaseOrderService) GetPurchaseOrder(ctx context.Context, r *pbcm.IDRequest) (*api.PurchaseOrder, error) {
	query := &purchaseorder.GetPurchaseOrderByIDQuery{
		ID:     r.Id,
		ShopID: s.SS.Shop().ID,
	}
	if err := s.PurchaseOrderQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := convertpball.PbPurchaseOrder(query.Result)
	result.InventoryVoucher = inventory.PbShopInventoryVoucher(query.Result.InventoryVoucher)
	return result, nil
}

func (s *PurchaseOrderService) GetPurchaseOrders(ctx context.Context, r *api.GetPurchaseOrdersRequest) (*api.PurchaseOrdersResponse, error) {
	paging := cmapi.CMPaging(r.Paging)
	query := &purchaseorder.ListPurchaseOrdersQuery{
		ShopID:  s.SS.Shop().ID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(r.Filters),
	}
	if err := s.PurchaseOrderQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}

	var purchaseOrders []*shop.PurchaseOrder
	for _, purchaseOrder := range query.Result.PurchaseOrders {
		purchaseOrderTemp := convertpball.PbPurchaseOrder(purchaseOrder)
		purchaseOrderTemp.InventoryVoucher = inventory.PbShopInventoryVoucher(purchaseOrder.InventoryVoucher)
		purchaseOrders = append(purchaseOrders, purchaseOrderTemp)
	}

	result := &api.PurchaseOrdersResponse{
		PurchaseOrders: purchaseOrders,
		Paging:         cmapi.PbPageInfo(paging),
	}
	return result, nil
}

func (s *PurchaseOrderService) GetPurchaseOrdersByIDs(ctx context.Context, r *pbcm.IDsRequest) (*api.PurchaseOrdersResponse, error) {
	query := &purchaseorder.GetPurchaseOrdersByIDsQuery{
		IDs:    r.Ids,
		ShopID: s.SS.Shop().ID,
	}
	if err := s.PurchaseOrderQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &api.PurchaseOrdersResponse{PurchaseOrders: convertpball.PbPurchaseOrders(query.Result.PurchaseOrders)}
	return result, nil
}

func (s *PurchaseOrderService) GetPurchaseOrdersByReceiptID(ctx context.Context, r *pbcm.IDRequest) (*api.PurchaseOrdersResponse, error) {
	query := &purchaseorder.ListPurchaseOrdersByReceiptIDQuery{
		ReceiptID: r.Id,
		ShopID:    s.SS.Shop().ID,
	}
	if err := s.PurchaseOrderQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &api.PurchaseOrdersResponse{PurchaseOrders: convertpball.PbPurchaseOrders(query.Result.PurchaseOrders)}
	return result, nil
}

func (s *PurchaseOrderService) CreatePurchaseOrder(ctx context.Context, r *api.CreatePurchaseOrderRequest) (*api.PurchaseOrder, error) {
	cmd := &purchaseorder.CreatePurchaseOrderCommand{
		ShopID:        s.SS.Shop().ID,
		SupplierID:    r.SupplierId,
		BasketValue:   r.BasketValue,
		TotalDiscount: r.TotalDiscount,
		DiscountLines: r.DiscountLines,
		TotalFee:      r.TotalFee,
		FeeLines:      r.FeeLines,
		TotalAmount:   r.TotalAmount,
		Note:          r.Note,
		Lines:         convertpball.Convert_api_PurchaseOrderLines_To_core_PurchaseOrderLines(r.Lines),
		CreatedBy:     s.SS.Claim().UserID,
	}
	if err := s.PurchaseOrderAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := convertpball.PbPurchaseOrder(cmd.Result)
	return result, nil
}

func (s *PurchaseOrderService) UpdatePurchaseOrder(ctx context.Context, r *api.UpdatePurchaseOrderRequest) (*api.PurchaseOrder, error) {
	cmd := &purchaseorder.UpdatePurchaseOrderCommand{
		ID:            r.Id,
		ShopID:        s.SS.Shop().ID,
		BasketValue:   r.BasketValue,
		TotalDiscount: r.TotalDiscount,
		DiscountLines: r.DiscountLines,
		TotalFee:      r.TotalFee,
		FeeLines:      r.FeeLines,
		TotalAmount:   r.TotalAmount,
		Note:          r.Note,
		Lines:         convertpball.Convert_api_PurchaseOrderLines_To_core_PurchaseOrderLines(r.Lines),
	}
	if err := s.PurchaseOrderAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := convertpball.PbPurchaseOrder(cmd.Result)
	return result, nil
}

func (s *PurchaseOrderService) DeletePurchaseOrder(ctx context.Context, r *pbcm.IDRequest) (*pbcm.DeletedResponse, error) {
	cmd := &purchaseorder.DeletePurchaseOrderCommand{
		ID:     r.Id,
		ShopID: s.SS.Shop().ID,
	}
	if err := s.PurchaseOrderAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return &pbcm.DeletedResponse{Deleted: cmd.Result}, nil
}

func (s *PurchaseOrderService) ConfirmPurchaseOrder(ctx context.Context, r *api.ConfirmPurchaseOrderRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &purchaseorder.ConfirmPurchaseOrderCommand{
		ID:                   r.Id,
		AutoInventoryVoucher: inventory.CheckRoleAutoInventoryVoucher(s.SS.CheckRoles, r.AutoInventoryVoucher),
		ShopID:               s.SS.Shop().ID,
	}
	if err := s.PurchaseOrderAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{Updated: cmd.Result}
	return result, nil
}

func (s *PurchaseOrderService) CancelPurchaseOrder(ctx context.Context, r *api.CancelPurchaseOrderRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &purchaseorder.CancelPurchaseOrderCommand{
		ID:                   r.Id,
		ShopID:               s.SS.Shop().ID,
		CancelReason:         util.CoalesceString(r.CancelReason, r.Reason),
		UpdatedBy:            s.SS.Claim().UserID,
		InventoryOverStock:   s.SS.Shop().InventoryOverstock.Apply(true),
		AutoInventoryVoucher: inventory.CheckRoleAutoInventoryVoucher(s.SS.CheckRoles, r.AutoInventoryVoucher),
	}
	if err := s.PurchaseOrderAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{Updated: cmd.Result}
	return result, nil
}
