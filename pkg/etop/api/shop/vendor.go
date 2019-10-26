package shop

import (
	"context"

	"etop.vn/api/shopping/vendoring"
	pbcm "etop.vn/backend/pb/common"
	pbshop "etop.vn/backend/pb/etop/shop"
	"etop.vn/backend/pkg/common/bus"
	wrapshop "etop.vn/backend/wrapper/etop/shop"
	. "etop.vn/capi/dot"
)

func init() {
	bus.AddHandlers("api",
		vendorService.GetVendor,
		vendorService.GetVendors,
		vendorService.GetVendorsByIDs,
		vendorService.CreateVendor,
		vendorService.UpdateVendor,
		vendorService.DeleteVendor)
}

func (s *VendorService) GetVendor(ctx context.Context, r *wrapshop.GetVendorEndpoint) error {
	query := &vendoring.GetVendorByIDQuery{
		ID:     r.Id,
		ShopID: r.Context.Shop.ID,
	}
	if err := vendorQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = pbshop.PbVendor(query.Result)
	return nil
}

func (s *VendorService) GetVendors(ctx context.Context, r *wrapshop.GetVendorsEndpoint) error {
	paging := r.Paging.CMPaging()
	query := &vendoring.ListVendorsQuery{
		ShopID:  r.Context.Shop.ID,
		Paging:  *paging,
		Filters: pbcm.ToFilters(r.Filters),
	}
	if err := vendorQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = &pbshop.VendorsResponse{
		Vendors: pbshop.PbVendors(query.Result.Vendors),
		Paging:  pbcm.PbPageInfo(paging, query.Result.Count),
	}
	return nil
}

func (s *VendorService) GetVendorsByIDs(ctx context.Context, r *wrapshop.GetVendorsByIDsEndpoint) error {
	query := &vendoring.ListVendorsByIDsQuery{
		IDs:    r.Ids,
		ShopID: r.Context.Shop.ID,
	}
	if err := vendorQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = &pbshop.VendorsResponse{
		Vendors: pbshop.PbVendors(query.Result.Vendors),
	}
	return nil
}

func (s *VendorService) CreateVendor(ctx context.Context, r *wrapshop.CreateVendorEndpoint) error {
	cmd := &vendoring.CreateVendorCommand{
		ShopID:   r.Context.Shop.ID,
		FullName: r.FullName,
		Note:     r.Note,
	}
	if err := vendorAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = pbshop.PbVendor(cmd.Result)
	return nil
}

func (s *VendorService) UpdateVendor(ctx context.Context, r *wrapshop.UpdateVendorEndpoint) error {
	cmd := &vendoring.UpdateVendorCommand{
		ID:       r.Id,
		ShopID:   r.Context.Shop.ID,
		FullName: PString(r.FullName),
		Note:     PString(r.Note),
	}
	if err := vendorAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = pbshop.PbVendor(cmd.Result)
	return nil
}

func (s *VendorService) DeleteVendor(ctx context.Context, r *wrapshop.DeleteVendorEndpoint) error {
	cmd := &vendoring.DeleteVendorCommand{
		ID:     r.Id,
		ShopID: r.Context.Shop.ID,
	}
	if err := vendorAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.DeletedResponse{Deleted: int32(cmd.Result)}
	return nil
}
