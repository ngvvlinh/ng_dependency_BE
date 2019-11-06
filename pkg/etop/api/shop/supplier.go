package shop

import (
	"context"

	"etop.vn/api/shopping/suppliering"
	pbcm "etop.vn/backend/pb/common"
	pbshop "etop.vn/backend/pb/etop/shop"
	"etop.vn/backend/pkg/common/bus"
	wrapshop "etop.vn/backend/wrapper/etop/shop"
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

func (s *SupplierService) GetSupplier(ctx context.Context, r *wrapshop.GetSupplierEndpoint) error {
	query := &suppliering.GetSupplierByIDQuery{
		ID:     r.Id,
		ShopID: r.Context.Shop.ID,
	}
	if err := supplierQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = pbshop.PbSupplier(query.Result)
	return nil
}

func (s *SupplierService) GetSuppliers(ctx context.Context, r *wrapshop.GetSuppliersEndpoint) error {
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
	return nil
}

func (s *SupplierService) GetSuppliersByIDs(ctx context.Context, r *wrapshop.GetSuppliersByIDsEndpoint) error {
	query := &suppliering.ListSuppliersByIDsQuery{
		IDs:    r.Ids,
		ShopID: r.Context.Shop.ID,
	}
	if err := supplierQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = &pbshop.SuppliersResponse{
		Suppliers: pbshop.PbSuppliers(query.Result.Suppliers),
	}
	return nil
}

func (s *SupplierService) CreateSupplier(ctx context.Context, r *wrapshop.CreateSupplierEndpoint) error {
	cmd := &suppliering.CreateSupplierCommand{
		ShopID:   r.Context.Shop.ID,
		FullName: r.FullName,
		Note:     r.Note,
	}
	if err := supplierAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = pbshop.PbSupplier(cmd.Result)
	return nil
}

func (s *SupplierService) UpdateSupplier(ctx context.Context, r *wrapshop.UpdateSupplierEndpoint) error {
	cmd := &suppliering.UpdateSupplierCommand{
		ID:       r.Id,
		ShopID:   r.Context.Shop.ID,
		FullName: PString(r.FullName),
		Note:     PString(r.Note),
	}
	if err := supplierAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = pbshop.PbSupplier(cmd.Result)
	return nil
}

func (s *SupplierService) DeleteSupplier(ctx context.Context, r *wrapshop.DeleteSupplierEndpoint) error {
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
