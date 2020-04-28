package shop

import (
	"context"
	"fmt"
	"time"

	"o.o/api/main/catalog"
	"o.o/api/main/purchaseorder"
	"o.o/api/main/receipting"
	"o.o/api/shopping/suppliering"
	"o.o/api/top/int/shop"
	pbcm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/status3"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/capi/dot"
)

func init() {
	bus.AddHandlers("api",
		supplierService.GetSupplier,
		supplierService.GetSuppliers,
		supplierService.GetSuppliersByIDs,
		supplierService.CreateSupplier,
		supplierService.UpdateSupplier,
		supplierService.DeleteSupplier)
}

func (s *SupplierService) GetSupplier(ctx context.Context, r *GetSupplierEndpoint) error {
	query := &suppliering.GetSupplierByIDQuery{
		ID:     r.Id,
		ShopID: r.Context.Shop.ID,
	}
	if err := supplierQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = convertpb.PbSupplier(query.Result)

	if err := s.listLiabilities(ctx, r.Context.Shop.ID, []*shop.Supplier{r.Result}); err != nil {
		return err
	}
	return nil
}

func (s *SupplierService) GetSuppliers(ctx context.Context, r *GetSuppliersEndpoint) error {
	paging := cmapi.CMPaging(r.Paging)
	query := &suppliering.ListSuppliersQuery{
		ShopID:  r.Context.Shop.ID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(r.Filters),
	}
	if err := supplierQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = &shop.SuppliersResponse{
		Suppliers: convertpb.PbSuppliers(query.Result.Suppliers),
		Paging:    cmapi.PbPageInfo(paging),
	}

	if err := s.listLiabilities(ctx, r.Context.Shop.ID, r.Result.Suppliers); err != nil {
		return err
	}
	return nil
}

func (s *SupplierService) GetSuppliersByIDs(ctx context.Context, r *GetSuppliersByIDsEndpoint) error {
	query := &suppliering.ListSuppliersByIDsQuery{
		IDs:    r.Ids,
		ShopID: r.Context.Shop.ID,
	}
	if err := supplierQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = &shop.SuppliersResponse{Suppliers: convertpb.PbSuppliers(query.Result.Suppliers)}

	if err := s.listLiabilities(ctx, r.Context.Shop.ID, r.Result.Suppliers); err != nil {
		return err
	}
	return nil
}

func (s *SupplierService) CreateSupplier(ctx context.Context, r *CreateSupplierEndpoint) error {
	key := fmt.Sprintf("CreateOrder %v-%v-%v-%v-%v",
		r.Context.Shop.ID, r.Context.UserID, r.FullName, r.Phone, r.Email)
	res, err := idempgroup.DoAndWrap(ctx, key, 15*time.Second,
		func() (interface{}, error) {
			cmd := &suppliering.CreateSupplierCommand{
				ShopID:            r.Context.Shop.ID,
				FullName:          r.FullName,
				Note:              r.Note,
				Phone:             r.Phone,
				Email:             r.Email,
				CompanyName:       r.CompanyName,
				TaxNumber:         r.TaxNumber,
				HeadquaterAddress: r.HeadquaterAddress,
			}
			if err := supplierAggr.Dispatch(ctx, cmd); err != nil {
				return nil, err
			}
			r.Result = convertpb.PbSupplier(cmd.Result)
			return r, nil
		}, "tạo nhà cung cấp")

	if err != nil {
		return err
	}
	r.Result = res.(*CreateSupplierEndpoint).Result
	return nil
}

func (s *SupplierService) UpdateSupplier(ctx context.Context, r *UpdateSupplierEndpoint) error {
	cmd := &suppliering.UpdateSupplierCommand{
		ID:                r.Id,
		ShopID:            r.Context.Shop.ID,
		FullName:          r.FullName,
		Phone:             r.Phone,
		Email:             r.Email,
		CompanyName:       r.CompanyName,
		TaxNumber:         r.TaxNumber,
		HeadquaterAddress: r.HeadquaterAddress,
		Note:              r.Note,
	}
	if err := supplierAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = convertpb.PbSupplier(cmd.Result)
	return nil
}

func (s *SupplierService) DeleteSupplier(ctx context.Context, r *DeleteSupplierEndpoint) error {
	cmd := &suppliering.DeleteSupplierCommand{
		ID:     r.Id,
		ShopID: r.Context.Shop.ID,
	}
	if err := supplierAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.DeletedResponse{Deleted: cmd.Result}
	return nil
}

func (s *SupplierService) listLiabilities(ctx context.Context, shopID dot.ID, suppliers []*shop.Supplier) error {
	var supplierIDs []dot.ID
	mapSupplierIDAndNumberOfPurchaseOrders := make(map[dot.ID]int)
	mapSupplierIDAndTotalAmountPurchaseOrders := make(map[dot.ID]int)
	mapSupplierIDAndTotalAmountReceipts := make(map[dot.ID]int)

	for _, supplier := range suppliers {
		supplierIDs = append(supplierIDs, supplier.Id)
	}
	listPurchaseOrdersBySuppliersQuery := &purchaseorder.ListPurchaseOrdersBySupplierIDsAndStatusesQuery{
		SupplierIDs: supplierIDs,
		ShopID:      shopID,
		Statuses:    []status3.Status{status3.Z, status3.P},
	}
	if err := purchaseOrderQuery.Dispatch(ctx, listPurchaseOrdersBySuppliersQuery); err != nil {
		return err
	}
	purchaseOrders := listPurchaseOrdersBySuppliersQuery.Result.PurchaseOrders
	for _, purchaseOrder := range purchaseOrders {
		mapSupplierIDAndNumberOfPurchaseOrders[purchaseOrder.SupplierID] += 1
		mapSupplierIDAndTotalAmountPurchaseOrders[purchaseOrder.SupplierID] += purchaseOrder.TotalAmount
	}

	listReceiptsBySupplierIDs := &receipting.ListReceiptsByTraderIDsAndStatusesQuery{
		ShopID:    shopID,
		TraderIDs: supplierIDs,
		Statuses:  []status3.Status{status3.P},
	}
	if err := receiptQuery.Dispatch(ctx, listReceiptsBySupplierIDs); err != nil {
		return err
	}
	receipts := listReceiptsBySupplierIDs.Result.Receipts
	for _, receipt := range receipts {
		mapSupplierIDAndTotalAmountReceipts[receipt.TraderID] += receipt.Amount
	}

	for _, supplier := range suppliers {
		supplier.Liability = &shop.SupplierLiability{
			TotalPurchaseOrders: mapSupplierIDAndNumberOfPurchaseOrders[supplier.Id],
			TotalAmount:         mapSupplierIDAndTotalAmountPurchaseOrders[supplier.Id],
			PaidAmount:          mapSupplierIDAndTotalAmountReceipts[supplier.Id],
			Liability:           mapSupplierIDAndTotalAmountPurchaseOrders[supplier.Id] - mapSupplierIDAndTotalAmountReceipts[supplier.Id],
		}
	}
	return nil
}

func (s *SupplierService) GetSuppliersByVariantID(ctx context.Context, r *GetSuppliersByVariantIDEndpoint) error {
	query := &catalog.GetSupplierIDsByVariantIDQuery{
		VariantID: r.VariantId,
		ShopID:    r.Context.Shop.ID,
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	querySuppplies := &suppliering.ListSuppliersByIDsQuery{
		IDs:    query.Result,
		ShopID: r.Context.Shop.ID,
	}
	if err := supplierQuery.Dispatch(ctx, querySuppplies); err != nil {
		return err
	}
	r.Result = &shop.SuppliersResponse{Suppliers: convertpb.PbSuppliers(querySuppplies.Result.Suppliers)}

	if err := s.listLiabilities(ctx, r.Context.Shop.ID, r.Result.Suppliers); err != nil {
		return err
	}
	return nil
}
