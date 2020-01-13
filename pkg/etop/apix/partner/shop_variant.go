package partner

import (
	"context"

	"etop.vn/api/main/catalog"
	"etop.vn/api/top/external/types"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/cmapi"
	"etop.vn/backend/pkg/etop/apix/convertpb"
	"etop.vn/capi/filter"
)

func (s *VariantService) GetVariant(ctx context.Context, r *GetVariantEndpoint) error {
	query := &catalog.GetShopVariantQuery{
		ExternalID: r.ExternalId,
		VariantID:  r.Id,
		ShopID:     r.Context.Shop.ID,
		Code:       r.Code,
		Result:     nil,
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = convertpb.PbShopVariant(query.Result)
	return nil
}

func (s *VariantService) ListVariants(ctx context.Context, r *ListVariantsEndpoint) error {
	paging, err := cmapi.CMCursorPaging(r.Paging)
	if err != nil {
		return err
	}
	var IDs filter.IDs
	if len(r.Filter.ID) != 0 {
		IDs = r.Filter.ID
	}
	query := &catalog.ListShopVariantsByIDsQuery{
		IDs:    IDs,
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

func (s *VariantService) CreateVariant(ctx context.Context, r *CreateVariantEndpoint) error {
	cmd := &catalog.CreateShopVariantCommand{
		ExternalID:   r.ExternalId,
		ExternalCode: r.ExternalCode,
		PartnerID:    r.Context.AuthPartnerID,
		ShopID:       r.Context.Shop.ID,
		ProductID:    r.ProductId,
		Code:         r.Code,
		Name:         r.Name,
		ImageURLs:    r.ImageUrls,
		Note:         r.Note,
		Attributes:   r.Attributes,
		DescriptionInfo: catalog.DescriptionInfo{
			ShortDesc:   r.ShortDesc,
			Description: r.Description,
		},
		PriceInfo: catalog.PriceInfo{
			CostPrice:   r.CostPrice,
			ListPrice:   r.ListPrice,
			RetailPrice: r.RetailPrice,
		},
	}
	if err := catalogAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = convertpb.PbShopVariant(cmd.Result)
	return nil
}

func (s *VariantService) UpdateVariant(ctx context.Context, r *UpdateVariantEndpoint) error {
	shopID := r.Context.Shop.ID
	cmd := &catalog.UpdateShopVariantInfoCommand{
		ShopID:    shopID,
		VariantID: r.Id,
		Name:      r.Name,
		Code:      r.Code,
		Note:      r.Note,

		ShortDesc:    r.ShortDesc,
		Descripttion: r.Description,

		CostPrice:   r.CostPrice,
		ListPrice:   r.ListPrice,
		RetailPrice: r.RetailPrice,
	}
	if err := catalogAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = convertpb.PbShopVariant(cmd.Result)
	return nil
}

func (s *VariantService) DeleteVariant(ctx context.Context, r *DeleteVariantEndpoint) error {
	return cm.ErrTODO
}
