package partner

import (
	"context"

	"etop.vn/api/main/catalog"
	externaltypes "etop.vn/api/top/external/types"
	"etop.vn/backend/pkg/common/apifw/cmapi"
	"etop.vn/backend/pkg/etop/apix/convertpb"
)

func (s *ProductService) GetProducts(ctx context.Context, r *GetProductsEndpoint) error {
	paging, err := cmapi.CMCursorPaging(r.Paging)
	if err != nil {
		return err
	}

	query := &catalog.ListShopProductsByIDsQuery{
		IDs:    r.Ids,
		ShopID: r.Context.Shop.ID,
		Paging: *paging,
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = &externaltypes.ShopProductsResponse{
		Products: convertpb.PbShopProducts(query.Result.Products),
		Paging:   convertpb.PbPageInfo(r.Paging, &query.Result.Paging),
	}
	return nil
}