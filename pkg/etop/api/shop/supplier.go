package shop

import (
	"context"

	"etop.vn/api/main/etop"
	"etop.vn/api/main/purchaseorder"
	"etop.vn/api/main/receipting"
	"etop.vn/api/shopping/suppliering"
	pbcm "etop.vn/backend/pb/common"
	pbshop "etop.vn/backend/pb/etop/shop"
	"etop.vn/backend/pkg/common/bus"
	. "etop.vn/capi/dot"
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
	r.Result = pbshop.PbSupplier(query.Result)

	if err := s.listLiabilities(ctx, r.Context.Shop.ID, []*pbshop.Supplier{r.Result}); err != nil {
		return err
	}
	return nil
}

func (s *SupplierService) GetSuppliers(ctx context.Context, r *GetSuppliersEndpoint) error {
	paging := r.Paging.CMPaging()
	query := &suppliering.ListSuppliersQuery{
		ShopID:  r.Context.Shop.ID,
		Paging:  *paging,
		Filters: pbcm.ToFilters(r.Filters),
	}
	if err := supplierQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = &pbshop.SuppliersResponse{
		Suppliers: pbshop.PbSuppliers(query.Result.Suppliers),
		Paging:    pbcm.PbPageInfo(paging, query.Result.Count),
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
	r.Result = &pbshop.SuppliersResponse{Suppliers: pbshop.PbSuppliers(query.Result.Suppliers)}

	if err := s.listLiabilities(ctx, r.Context.Shop.ID, r.Result.Suppliers); err != nil {
		return err
	}
	return nil
}

func (s *SupplierService) CreateSupplier(ctx context.Context, r *CreateSupplierEndpoint) error {
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
		return err
	}
	r.Result = pbshop.PbSupplier(cmd.Result)
	return nil
}

func (s *SupplierService) UpdateSupplier(ctx context.Context, r *UpdateSupplierEndpoint) error {
	cmd := &suppliering.UpdateSupplierCommand{
		ID:                r.Id,
		ShopID:            r.Context.Shop.ID,
		FullName:          PString(r.FullName),
		Phone:             PString(r.Phone),
		Email:             PString(r.Email),
		CompanyName:       PString(r.CompanyName),
		TaxNumber:         PString(r.TaxNumber),
		HeadquaterAddress: PString(r.HeadquaterAddress),
	}
	if err := supplierAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = pbshop.PbSupplier(cmd.Result)
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
	r.Result = &pbcm.DeletedResponse{Deleted: int32(cmd.Result)}
	return nil
}

func (s *SupplierService) listLiabilities(ctx context.Context, shopID int64, suppliers []*pbshop.Supplier) error {
	var supplierIDs []int64
	mapSupplierIDAndNumberOfPurchaseOrders := make(map[int64]int)
	mapSupplierIDAndTotalAmountPurchaseOrders := make(map[int64]int64)
	mapSupplierIDAndTotalAmountReceipts := make(map[int64]int64)

	for _, supplier := range suppliers {
		supplierIDs = append(supplierIDs, supplier.Id)
	}
	listPurchaseOrdersBySuppliersQuery := &purchaseorder.ListPurchaseOrdersBySupplierIDsAndStatusesQuery{
		SupplierIDs: supplierIDs,
		ShopID:      shopID,
		Statuses:    []etop.Status3{etop.S3Zero, etop.S3Positive},
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
		Statuses:  []etop.Status3{etop.S3Positive},
	}
	if err := receiptQuery.Dispatch(ctx, listReceiptsBySupplierIDs); err != nil {
		return err
	}
	receipts := listReceiptsBySupplierIDs.Result.Receipts
	for _, receipt := range receipts {
		mapSupplierIDAndTotalAmountReceipts[receipt.TraderID] += int64(receipt.Amount)
	}

	for _, supplier := range suppliers {
		supplier.Liability = &pbshop.SupplierLiability{
			TotalPurchaseOrders: int32(mapSupplierIDAndNumberOfPurchaseOrders[supplier.Id]),
			TotalAmount:         mapSupplierIDAndTotalAmountPurchaseOrders[supplier.Id],
			PaidAmount:          mapSupplierIDAndTotalAmountReceipts[supplier.Id],
			Liability:           mapSupplierIDAndTotalAmountPurchaseOrders[supplier.Id] - mapSupplierIDAndTotalAmountReceipts[supplier.Id],
		}
	}
	return nil
}
