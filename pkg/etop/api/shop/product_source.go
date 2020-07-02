package shop

import (
	"context"

	api "o.o/api/top/int/shop"
	pbcm "o.o/api/top/types/common"
	catalogmodelx "o.o/backend/com/main/catalog/modelx"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
)

type ProductSourceService struct {
	session.Session
}

func (s *ProductSourceService) Clone() api.ProductSourceService { res := *s; return &res }

func (s *ProductSourceService) CreateVariant(ctx context.Context, q *api.DeprecatedCreateVariantRequest) (*api.ShopProduct, error) {
	cmd := &catalogmodelx.DeprecatedCreateVariantCommand{
		ShopID:      s.SS.Shop().ID,
		ProductID:   q.ProductId,
		ProductName: q.ProductName,
		Name:        q.Name,
		Description: q.Description,
		ShortDesc:   q.ShortDesc,
		ImageURLs:   q.ImageUrls,
		Tags:        q.Tags,
		Status:      q.Status,

		CostPrice:   q.CostPrice,
		ListPrice:   q.ListPrice,
		RetailPrice: q.RetailPrice,

		ProductCode:       q.Code,
		VariantCode:       q.Code,
		QuantityAvailable: q.QuantityAvailable,
		QuantityOnHand:    q.QuantityOnHand,
		QuantityReserved:  q.QuantityReserved,

		Attributes: q.Attributes,
		DescHTML:   q.DescHtml,
	}

	if err := bus.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	result := PbShopProductWithVariants(cmd.Result)
	return result, nil
}

func (s *ProductSourceService) CreateProductSourceCategory(ctx context.Context, q *api.CreatePSCategoryRequest) (*api.Category, error) {
	cmd := &catalogmodelx.CreateShopCategoryCommand{
		ShopID:   s.SS.Shop().ID,
		Name:     q.Name,
		ParentID: q.ParentId,
	}

	if err := bus.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := convertpb.PbCategory(cmd.Result)
	return result, nil
}

func (s *ProductSourceService) UpdateProductsPSCategory(ctx context.Context, q *api.UpdateProductsPSCategoryRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &catalogmodelx.UpdateProductsShopCategoryCommand{
		CategoryID: q.CategoryId,
		ProductIDs: q.ProductIds,
		ShopID:     s.SS.Shop().ID,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{
		Updated: cmd.Result.Updated,
	}
	return result, nil
}

func (s *ProductSourceService) GetProductSourceCategory(ctx context.Context, q *pbcm.IDRequest) (*api.Category, error) {
	cmd := &catalogmodelx.GetShopCategoryQuery{
		ShopID:     s.SS.Shop().ID,
		CategoryID: q.Id,
	}

	if err := bus.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	result := convertpb.PbCategory(cmd.Result)
	return result, nil
}

func (s *ProductSourceService) GetProductSourceCategories(ctx context.Context, q *api.GetProductSourceCategoriesRequest) (*api.CategoriesResponse, error) {
	cmd := &catalogmodelx.GetProductSourceCategoriesQuery{
		ShopID: s.SS.Shop().ID,
	}

	if err := bus.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	result := &api.CategoriesResponse{
		Categories: convertpb.PbCategories(cmd.Result.Categories),
	}
	return result, nil
}

func (s *ProductSourceService) UpdateProductSourceCategory(ctx context.Context, q *api.UpdateProductSourceCategoryRequest) (*api.Category, error) {
	cmd := &catalogmodelx.UpdateShopCategoryCommand{
		ID:       q.Id,
		ShopID:   s.SS.Shop().ID,
		ParentID: q.ParentId,
		Name:     q.Name,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := convertpb.PbCategory(cmd.Result)
	return result, nil
}

func (s *ProductSourceService) RemoveProductSourceCategory(ctx context.Context, q *pbcm.IDRequest) (*pbcm.RemovedResponse, error) {
	cmd := &catalogmodelx.RemoveShopCategoryCommand{
		ID:     q.Id,
		ShopID: s.SS.Shop().ID,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.RemovedResponse{
		Removed: cmd.Result.Removed,
	}
	return result, nil
}
