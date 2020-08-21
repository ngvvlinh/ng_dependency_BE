package supplier

import (
	"context"
	"fmt"
	"time"

	"o.o/api/main/catalog"
	"o.o/api/main/purchaseorder"
	"o.o/api/main/receipting"
	"o.o/api/shopping/suppliering"
	"o.o/api/top/int/shop"
	api "o.o/api/top/int/shop"
	pbcm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/status3"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/api/convertpb"
	shop2 "o.o/backend/pkg/etop/api/shop"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/capi/dot"
)

type SupplierService struct {
	session.Session

	CatalogQuery       catalog.QueryBus
	PurchaseOrderQuery purchaseorder.QueryBus
	ReceiptQuery       receipting.QueryBus
	SupplierAggr       suppliering.CommandBus
	SupplierQuery      suppliering.QueryBus
}

func (s *SupplierService) Clone() api.SupplierService { res := *s; return &res }

func (s *SupplierService) GetSupplier(ctx context.Context, r *pbcm.IDRequest) (*api.Supplier, error) {
	query := &suppliering.GetSupplierByIDQuery{
		ID:     r.Id,
		ShopID: s.SS.Shop().ID,
	}
	if err := s.SupplierQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := convertpb.PbSupplier(query.Result)

	if err := s.listLiabilities(ctx, s.SS.Shop().ID, []*shop.Supplier{result}); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *SupplierService) GetSuppliers(ctx context.Context, r *api.GetSuppliersRequest) (*api.SuppliersResponse, error) {
	paging := cmapi.CMPaging(r.Paging)
	query := &suppliering.ListSuppliersQuery{
		ShopID:  s.SS.Shop().ID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(r.Filters),
	}
	if err := s.SupplierQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &api.SuppliersResponse{
		Suppliers: convertpb.PbSuppliers(query.Result.Suppliers),
		Paging:    cmapi.PbPageInfo(paging),
	}

	if err := s.listLiabilities(ctx, s.SS.Shop().ID, result.Suppliers); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *SupplierService) GetSuppliersByIDs(ctx context.Context, r *pbcm.IDsRequest) (*api.SuppliersResponse, error) {
	query := &suppliering.ListSuppliersByIDsQuery{
		IDs:    r.Ids,
		ShopID: s.SS.Shop().ID,
	}
	if err := s.SupplierQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &api.SuppliersResponse{Suppliers: convertpb.PbSuppliers(query.Result.Suppliers)}

	if err := s.listLiabilities(ctx, s.SS.Shop().ID, result.Suppliers); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *SupplierService) CreateSupplier(ctx context.Context, r *api.CreateSupplierRequest) (*api.Supplier, error) {
	key := fmt.Sprintf("CreateOrder %v-%v-%v-%v-%v",
		s.SS.Shop().ID, s.SS.Claim().UserID, r.FullName, r.Phone, r.Email)
	res, _, err := shop2.Idempgroup.DoAndWrap(
		ctx, key, 15*time.Second, "tạo nhà cung cấp",
		func() (interface{}, error) {
			cmd := &suppliering.CreateSupplierCommand{
				ShopID:            s.SS.Shop().ID,
				FullName:          r.FullName,
				Note:              r.Note,
				Phone:             r.Phone,
				Email:             r.Email,
				CompanyName:       r.CompanyName,
				TaxNumber:         r.TaxNumber,
				HeadquaterAddress: r.HeadquaterAddress,
			}
			if err := s.SupplierAggr.Dispatch(ctx, cmd); err != nil {
				return nil, err
			}
			result := convertpb.PbSupplier(cmd.Result)
			return result, nil
		})

	if err != nil {
		return nil, err
	}
	result := res.(*api.Supplier)
	return result, nil
}

func (s *SupplierService) UpdateSupplier(ctx context.Context, r *api.UpdateSupplierRequest) (*api.Supplier, error) {
	cmd := &suppliering.UpdateSupplierCommand{
		ID:                r.Id,
		ShopID:            s.SS.Shop().ID,
		FullName:          r.FullName,
		Phone:             r.Phone,
		Email:             r.Email,
		CompanyName:       r.CompanyName,
		TaxNumber:         r.TaxNumber,
		HeadquaterAddress: r.HeadquaterAddress,
		Note:              r.Note,
	}
	if err := s.SupplierAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := convertpb.PbSupplier(cmd.Result)
	return result, nil
}

func (s *SupplierService) DeleteSupplier(ctx context.Context, r *pbcm.IDRequest) (*pbcm.DeletedResponse, error) {
	cmd := &suppliering.DeleteSupplierCommand{
		ID:     r.Id,
		ShopID: s.SS.Shop().ID,
	}
	if err := s.SupplierAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.DeletedResponse{Deleted: cmd.Result}
	return result, nil
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
	if err := s.PurchaseOrderQuery.Dispatch(ctx, listPurchaseOrdersBySuppliersQuery); err != nil {
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
	if err := s.ReceiptQuery.Dispatch(ctx, listReceiptsBySupplierIDs); err != nil {
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

func (s *SupplierService) GetSuppliersByVariantID(ctx context.Context, r *api.GetSuppliersByVariantIDRequest) (*api.SuppliersResponse, error) {
	query := &catalog.GetSupplierIDsByVariantIDQuery{
		VariantID: r.VariantId,
		ShopID:    s.SS.Shop().ID,
	}
	if err := s.CatalogQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	querySuppplies := &suppliering.ListSuppliersByIDsQuery{
		IDs:    query.Result,
		ShopID: s.SS.Shop().ID,
	}
	if err := s.SupplierQuery.Dispatch(ctx, querySuppplies); err != nil {
		return nil, err
	}
	result := &api.SuppliersResponse{Suppliers: convertpb.PbSuppliers(querySuppplies.Result.Suppliers)}

	if err := s.listLiabilities(ctx, s.SS.Shop().ID, result.Suppliers); err != nil {
		return nil, err
	}
	return result, nil
}
