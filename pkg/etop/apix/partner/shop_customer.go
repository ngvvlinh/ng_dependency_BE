package partner

import (
	"context"

	"etop.vn/api/shopping/customering"
	externaltypes "etop.vn/api/top/external/types"
	"etop.vn/backend/pkg/common/apifw/cmapi"
	"etop.vn/backend/pkg/etop/apix/convertpb"
)

func (s *CustomerService) GetCustomers(ctx context.Context, r *GetCustomersEndpoint) error {
	paging, err := cmapi.CMCursorPaging(r.Paging)
	if err != nil {
		return err
	}

	query := &customering.ListCustomersByIDsQuery{
		ShopID: r.Context.Shop.ID,
		IDs:    r.Ids,
		Paging: *paging,
	}
	if err := customerQuery.Dispatch(ctx, query); err != nil {
		return err
	}

	r.Result = &externaltypes.CustomersResponse{
		Customers: convertpb.PbShopCustomers(query.Result.Customers),
		Paging:    convertpb.PbPageInfo(r.Paging, &query.Result.Paging),
	}
	return nil
}
