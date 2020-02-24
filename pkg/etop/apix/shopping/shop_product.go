package shopping

import (
	"context"

	"etop.vn/api/main/catalog"
	externaltypes "etop.vn/api/top/external/types"
	"etop.vn/api/top/types/common"
	cm "etop.vn/api/top/types/common"
	"etop.vn/backend/pkg/common/apifw/cmapi"
	"etop.vn/backend/pkg/etop/apix/convertpb"
	"etop.vn/capi/dot"
	"etop.vn/capi/filter"
)

func GetProduct(ctx context.Context, shopID dot.ID, request *externaltypes.GetProductRequest) (*externaltypes.ShopProduct, error) {
	query := &catalog.GetShopProductQuery{
		ExternalID: request.ExternalId,
		Code:       request.Code,
		ProductID:  request.Id,
		ShopID:     shopID,
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	return convertpb.PbShopProduct(query.Result), nil
}

func ListProducts(ctx context.Context, shopID dot.ID, request *externaltypes.ListProductsRequest) (*externaltypes.ShopProductsResponse, error) {
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
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	return &externaltypes.ShopProductsResponse{
		Products: convertpb.PbShopProducts(query.Result.Products),
		Paging:   convertpb.PbPageInfo(paging, &query.Result.Paging),
	}, nil
}

func CreateProduct(ctx context.Context, shopID dot.ID, partnerID dot.ID, request *externaltypes.CreateProductRequest) (*externaltypes.ShopProduct, error) {
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
	if err := catalogAggregate.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return convertpb.ConvertProductWithVariantsToPbProduct(cmd.Result), nil
}

func UpdateProduct(ctx context.Context, shopID dot.ID, request *externaltypes.UpdateProductRequest) (*externaltypes.ShopProduct, error) {
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
	if err := catalogAggregate.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return convertpb.ConvertProductWithVariantsToPbProduct(cmd.Result), nil
}

func DeleteProduct(ctx context.Context, shopID dot.ID, request *externaltypes.GetProductRequest) (*cm.Empty, error) {
	var IDs []dot.ID
	IDs = append(IDs, request.Id)
	cmd := &catalog.DeleteShopProductsCommand{
		IDs:    IDs,
		ShopID: shopID,
	}
	if err := catalogAggregate.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return &common.Empty{}, nil
}
