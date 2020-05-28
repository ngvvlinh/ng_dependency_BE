package shopping

import (
	"context"

	"o.o/api/main/catalog"
	"o.o/api/top/external/types"
	externaltypes "o.o/api/top/external/types"
	"o.o/api/top/types/common"
	cm "o.o/api/top/types/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/apix/convertpb"
	"o.o/capi/dot"
	"o.o/capi/filter"
)

func (s *Shopping) GetVariant(ctx context.Context, shopID dot.ID, request *externaltypes.GetVariantRequest) (*externaltypes.ShopVariant, error) {
	query := &catalog.GetShopVariantQuery{
		ExternalID: request.ExternalId,
		VariantID:  request.Id,
		ShopID:     shopID,
		Code:       request.Code,
	}
	if err := s.CatalogQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	return convertpb.PbShopVariant(query.Result), nil
}

func (s *Shopping) ListVariants(ctx context.Context, shopID dot.ID, request *externaltypes.ListVariantsRequest) (*externaltypes.ShopVariantsResponse, error) {
	paging, err := cmapi.CMCursorPaging(request.Paging)
	if err != nil {
		return nil, err
	}
	var IDs filter.IDs
	if len(request.Filter.ID) != 0 {
		IDs = request.Filter.ID
	}
	query := &catalog.ListShopVariantsByIDsQuery{
		IDs:            IDs,
		ShopID:         shopID,
		Paging:         *paging,
		IncludeDeleted: request.IncludeDeleted,
	}
	if err := s.CatalogQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	return &types.ShopVariantsResponse{
		ShopVariants: convertpb.PbShopVariants(query.Result.Variants),
		Paging:       convertpb.PbPageInfo(paging, &query.Result.Paging),
	}, nil
}

func (s *Shopping) CreateVariant(ctx context.Context, shopID dot.ID, partnerID dot.ID, request *externaltypes.CreateVariantRequest) (*externaltypes.ShopVariant, error) {
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
	if err := s.CatalogAggregate.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return convertpb.PbShopVariant(cmd.Result), nil
}

func (s *Shopping) UpdateVariant(ctx context.Context, shopID dot.ID, request *externaltypes.UpdateVariantRequest) (*externaltypes.ShopVariant, error) {
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
	if err := s.CatalogAggregate.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return convertpb.PbShopVariant(cmd.Result), nil
}

func (s *Shopping) DeleteVariant(ctx context.Context, shopID dot.ID, request *externaltypes.GetVariantRequest) (*cm.Empty, error) {
	var IDs []dot.ID
	IDs = append(IDs, request.Id)
	cmd := &catalog.DeleteShopVariantsCommand{
		IDs:    IDs,
		ShopID: shopID,
	}
	if err := s.CatalogAggregate.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return &common.Empty{}, nil
}
