package shop

import (
	"context"

	"etop.vn/api/shopping/customering"
	pbcm "etop.vn/backend/pb/common"
	pbshop "etop.vn/backend/pb/etop/shop"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	wrapshop "etop.vn/backend/wrapper/etop/shop"
	. "etop.vn/capi/dot"
)

func init() {
	bus.AddHandlers("api",
		CreateCustomer,
		UpdateCustomer,
		DeleteCustomer,
		GetCustomer,
		GetCustomers,
		GetCustomersByIDs,
		GetCustomerDetails,

		BatchSetCustomersStatus,
	)
}

func CreateCustomer(ctx context.Context, r *wrapshop.CreateCustomerEndpoint) error {
	cmd := &customering.CreateCustomerCommand{
		ShopID:   r.Context.Shop.ID,
		Code:     r.Code,
		FullName: r.FullName,
		Gender:   r.Gender,
		Type:     r.Type,
		Birthday: r.Birthday,
		Note:     r.Note,
		Phone:    r.Phone,
		Email:    r.Email,
	}
	if err := customerAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = pbshop.PbCustomer(cmd.Result)
	return nil
}

func UpdateCustomer(ctx context.Context, r *wrapshop.UpdateCustomerEndpoint) error {
	cmd := &customering.UpdateCustomerCommand{
		ID:       r.Id,
		ShopID:   r.Context.Shop.ID,
		Code:     PString(r.Code),
		FullName: PString(r.FullName),
		Gender:   PString(r.Gender),
		Type:     PString(r.Type),
		Birthday: PString(r.Birthday),
		Note:     PString(r.Note),
		Phone:    PString(r.Phone),
		Email:    PString(r.Email),
	}
	if err := customerAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = pbshop.PbCustomer(cmd.Result)
	return nil
}

func BatchSetCustomersStatus(ctx context.Context, r *wrapshop.BatchSetCustomersStatusEndpoint) error {
	cmd := &customering.BatchSetCustomersStatusCommand{
		IDs:    r.Ids,
		ShopID: r.Context.Shop.ID,
		Status: int32(r.Status),
	}
	if err := customerAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.UpdatedResponse{Updated: cmd.Result.Updated}
	return nil
}

func DeleteCustomer(ctx context.Context, r *wrapshop.DeleteCustomerEndpoint) error {
	cmd := &customering.DeleteCustomerCommand{
		ID:     r.Id,
		ShopID: r.Context.Shop.ID,
	}
	if err := customerAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.DeletedResponse{Deleted: int32(cmd.Result)}
	return nil
}

func GetCustomer(ctx context.Context, r *wrapshop.GetCustomerEndpoint) error {
	query := &customering.GetCustomerByIDQuery{
		ID:     r.Id,
		ShopID: r.Context.Shop.ID,
	}
	if err := customerQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = pbshop.PbCustomer(query.Result)
	return nil
}

func GetCustomers(ctx context.Context, r *wrapshop.GetCustomersEndpoint) error {
	paging := r.Paging.CMPaging()
	query := &customering.ListCustomersQuery{
		ShopID:  r.Context.Shop.ID,
		Paging:  *paging,
		Filters: pbcm.ToFilters(r.Filters),
	}
	if err := customerQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = &pbshop.CustomersResponse{
		Customers: pbshop.PbCustomers(query.Result.Customers),
		Paging:    pbcm.PbPageInfo(paging, query.Result.Count),
	}
	return nil
}

func GetCustomersByIDs(ctx context.Context, r *wrapshop.GetCustomersByIDsEndpoint) error {
	query := &customering.ListCustomersByIDsQuery{
		IDs:    r.Ids,
		ShopID: r.Context.Shop.ID,
	}
	if err := customerQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = &pbshop.CustomersResponse{
		Customers: pbshop.PbCustomers(query.Result.Customers),
	}
	return nil
}

func GetCustomerDetails(ctx context.Context, r *wrapshop.GetCustomerDetailsEndpoint) error {
	return cm.ErrTODO
}
