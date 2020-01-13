package partner

import (
	"context"

	"etop.vn/api/main/catalog"
	externaltypes "etop.vn/api/top/external/types"
	"etop.vn/backend/pkg/common/apifw/cmapi"
	"etop.vn/backend/pkg/etop/apix/convertpb"
	"etop.vn/capi/filter"
)

func (s *ProductService) GetProduct(ctx context.Context, r *GetProductEndpoint) error {
	query := &catalog.GetShopProductQuery{
		ExternalID: r.ExternalId,
		Code:       r.Code,
		ProductID:  r.Id,
		ShopID:     r.Context.Shop.ID,
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = convertpb.PbShopProduct(query.Result)
	return nil
}

func (s *ProductService) ListProducts(ctx context.Context, r *ListProductsEndpoint) error {
	paging, err := cmapi.CMCursorPaging(r.Paging)
	if err != nil {
		return err
	}
	var IDs filter.IDs
	if len(r.Filter.ID) != 0 {
		IDs = r.Filter.ID
	}
	query := &catalog.ListShopProductsByIDsQuery{
		IDs:    IDs,
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

func (s *ProductService) CreateProduct(ctx context.Context, r *CreateProductEndpoint) error {
	cmd := &catalog.CreateShopProductCommand{
		ExternalID:   r.ExternalId,
		ExternalCode: r.ExternalCode,
		PartnerID:    r.Context.AuthPartnerID,
		ShopID:       r.Context.Shop.ID,
		Code:         r.Code,
		Name:         r.Name,
		Unit:         r.Unit,
		ImageURLs:    r.ImageUrls,
		Note:         r.Note,
		DescriptionInfo: catalog.DescriptionInfo{
			ShortDesc:   r.ShortDesc,
			Description: r.Description,
		},
		PriceInfo: catalog.PriceInfo{
			CostPrice:   r.CostPrice,
			ListPrice:   r.ListPrice,
			RetailPrice: r.RetailPrice,
		},
		BrandID: r.BrandId,
	}
	if err := catalogAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = convertpb.ConvertProductWithVariantsToPbProduct(cmd.Result)
	return nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, r *UpdateProductEndpoint) error {
	cmd := &catalog.UpdateShopProductInfoCommand{
		ShopID:    r.Context.Shop.ID,
		ProductID: r.Id,
		Name:      r.Name,
		Unit:      r.Unit,
		Note:      r.Note,
		BrandID:   r.BrandId,

		ShortDesc:   r.ShortDesc,
		Description: r.Description,

		CostPrice:   r.CostPrice,
		ListPrice:   r.ListPrice,
		RetailPrice: r.RetailPrice,
	}
	if err := catalogAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = convertpb.ConvertProductWithVariantsToPbProduct(cmd.Result)
	return nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, r *DeleteProductEndpoint) error {
	// TODO:
	return nil
}
