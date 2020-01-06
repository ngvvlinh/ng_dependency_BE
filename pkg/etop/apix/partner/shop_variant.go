package partner

import (
	"context"

	"etop.vn/api/main/catalog"
	"etop.vn/api/top/external/types"
	"etop.vn/backend/pkg/common/apifw/cmapi"
	"etop.vn/backend/pkg/etop/apix/convertpb"
)

func (s *VariantService) GetVariants(ctx context.Context, r *GetVariantsEndpoint) error {
	paging, err := cmapi.CMCursorPaging(r.Paging)
	if err != nil {
		return err
	}

	query := &catalog.ListShopVariantsByIDsQuery{
		IDs:    r.Ids,
		ShopID: r.Context.Shop.ID,
		Paging: *paging,
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = &types.ShopVariantsResponse{
		ShopVariants: convertpb.PbShopVariants(query.Result.Variants),
		Paging:       convertpb.PbPageInfo(r.Paging, &query.Result.Paging),
	}
	return nil
}
