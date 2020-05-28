package shopping

import (
	"context"

	"o.o/api/main/catalog"
	externaltypes "o.o/api/top/external/types"
	"o.o/api/top/types/common"
	cm "o.o/api/top/types/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/apix/convertpb"
	"o.o/capi/dot"
	"o.o/capi/filter"
)

func (s *Shopping) GetProduct(ctx context.Context, shopID dot.ID, request *externaltypes.GetProductRequest) (*externaltypes.ShopProduct, error) {
	query := &catalog.GetShopProductQuery{
		ExternalID: request.ExternalId,
		Code:       request.Code,
		ProductID:  request.Id,
		ShopID:     shopID,
	}
	if err := s.CatalogQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	return convertpb.PbShopProduct(query.Result), nil
}

func (s *Shopping) ListProducts(ctx context.Context, shopID dot.ID, request *externaltypes.ListProductsRequest) (*externaltypes.ShopProductsResponse, error) {
	paging, err := cmapi.CMCursorPaging(request.Paging)
	if err != nil {
		return nil, err
	}
	var IDs filter.IDs
	if len(request.Filter.ID) != 0 {
		IDs = request.Filter.ID
	}
	query := &catalog.ListShopProductsByIDsQuery{
		IDs:            IDs,
		ShopID:         shopID,
		Paging:         *paging,
		IncludeDeleted: request.IncludeDeleted,
	}
	if err := s.CatalogQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	return &externaltypes.ShopProductsResponse{
		Products: convertpb.PbShopProducts(query.Result.Products),
		Paging:   convertpb.PbPageInfo(paging, &query.Result.Paging),
	}, nil
}

func (s *Shopping) CreateProduct(ctx context.Context, shopID dot.ID, partnerID dot.ID, request *externaltypes.CreateProductRequest) (*externaltypes.ShopProduct, error) {
	cmd := &catalog.CreateShopProductCommand{
		ExternalID:   request.ExternalId,
		ExternalCode: request.ExternalCode,
		PartnerID:    partnerID,
		ShopID:       shopID,
		Code:         request.Code,
		Name:         request.Name,
		Unit:         request.Unit,
		ImageURLs:    request.ImageUrls,
		Note:         request.Note,
		DescriptionInfo: catalog.DescriptionInfo{
			ShortDesc:   request.ShortDesc,
			Description: request.Description,
		},
		PriceInfo: catalog.PriceInfo{
			CostPrice:   request.CostPrice,
			ListPrice:   request.ListPrice,
			RetailPrice: request.RetailPrice,
		},
		BrandID: request.BrandId,
	}
	if err := s.CatalogAggregate.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return convertpb.ConvertProductWithVariantsToPbProduct(cmd.Result), nil
}

func (s *Shopping) UpdateProduct(ctx context.Context, shopID dot.ID, request *externaltypes.UpdateProductRequest) (*externaltypes.ShopProduct, error) {
	cmd := &catalog.UpdateShopProductInfoCommand{
		ShopID:    shopID,
		ProductID: request.Id,
		Name:      request.Name,
		Unit:      request.Unit,
		Note:      request.Note,
		BrandID:   request.BrandId,

		ShortDesc:   request.ShortDesc,
		Description: request.Description,

		CostPrice:   request.CostPrice,
		ListPrice:   request.ListPrice,
		RetailPrice: request.RetailPrice,
	}
	if err := s.CatalogAggregate.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return convertpb.ConvertProductWithVariantsToPbProduct(cmd.Result), nil
}

func (s *Shopping) DeleteProduct(ctx context.Context, shopID dot.ID, request *externaltypes.GetProductRequest) (*cm.Empty, error) {
	var IDs []dot.ID
	IDs = append(IDs, request.Id)
	cmd := &catalog.DeleteShopProductsCommand{
		IDs:    IDs,
		ShopID: shopID,
	}
	if err := s.CatalogAggregate.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return &common.Empty{}, nil
}
