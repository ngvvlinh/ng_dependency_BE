package shopping

import (
	"context"

	"etop.vn/api/main/catalog"
	"etop.vn/api/top/external/types"
	externaltypes "etop.vn/api/top/external/types"
	"etop.vn/api/top/types/common"
	cm "etop.vn/api/top/types/common"
	"etop.vn/backend/pkg/common/apifw/cmapi"
	"etop.vn/backend/pkg/etop/apix/convertpb"
	"etop.vn/capi/dot"
	"etop.vn/capi/filter"
)

func GetVariant(ctx context.Context, shopID dot.ID, request *externaltypes.GetVariantRequest) (*externaltypes.ShopVariant, error) {
	query := &catalog.GetShopVariantQuery{
		ExternalID: request.ExternalId,
		VariantID:  request.Id,
		ShopID:     shopID,
		Code:       request.Code,
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	return convertpb.PbShopVariant(query.Result), nil
}

func ListVariants(ctx context.Context, shopID dot.ID, request *externaltypes.ListVariantsRequest) (*externaltypes.ShopVariantsResponse, error) {
	paging, err := cmapi.CMCursorPaging(request.Paging)
	if err != nil {
		return nil, err
	}
	var IDs filter.IDs
	if len(request.Filter.ID) != 0 {
		IDs = request.Filter.ID
	}
	query := &catalog.ListShopVariantsByIDsQuery{
		IDs:    IDs,
		ShopID: shopID,
		Paging: *paging,
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	return &types.ShopVariantsResponse{
		ShopVariants: convertpb.PbShopVariants(query.Result.Variants),
		Paging:       convertpb.PbPageInfo(paging, &query.Result.Paging),
	}, nil
}

func CreateVariant(ctx context.Context, shopID dot.ID, partnerID dot.ID, request *externaltypes.CreateVariantRequest) (*externaltypes.ShopVariant, error) {
	cmd := &catalog.CreateShopVariantCommand{
		ExternalID:   request.ExternalId,
		ExternalCode: request.ExternalCode,
		PartnerID:    partnerID,
		ShopID:       shopID,
		ProductID:    request.ProductId,
		Code:         request.Code,
		Name:         request.Name,
		ImageURLs:    request.ImageUrls,
		Note:         request.Note,
		Attributes:   request.Attributes,
		DescriptionInfo: catalog.DescriptionInfo{
			ShortDesc:   request.ShortDesc,
			Description: request.Description,
		},
		PriceInfo: catalog.PriceInfo{
			CostPrice:   request.CostPrice,
			ListPrice:   request.ListPrice,
			RetailPrice: request.RetailPrice,
		},
	}
	if err := catalogAggregate.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return convertpb.PbShopVariant(cmd.Result), nil
}

func UpdateVariant(ctx context.Context, shopID dot.ID, request *externaltypes.UpdateVariantRequest) (*externaltypes.ShopVariant, error) {
	cmd := &catalog.UpdateShopVariantInfoCommand{
		ShopID:    shopID,
		VariantID: request.Id,
		Name:      request.Name,
		Code:      request.Code,
		Note:      request.Note,

		ShortDesc:    request.ShortDesc,
		Descripttion: request.Description,

		CostPrice:   request.CostPrice,
		ListPrice:   request.ListPrice,
		RetailPrice: request.RetailPrice,
	}
	if err := catalogAggregate.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return convertpb.PbShopVariant(cmd.Result), nil
}

func DeleteVariant(ctx context.Context, shopID dot.ID, request *externaltypes.GetVariantRequest) (*cm.Empty, error) {
	var IDs []dot.ID
	IDs = append(IDs, request.Id)
	cmd := &catalog.DeleteShopVariantsCommand{
		IDs:    IDs,
		ShopID: shopID,
	}
	if err := catalogAggregate.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return &common.Empty{}, nil
}
