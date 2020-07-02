package shop

import (
	"context"

	"o.o/api/main/catalog"
	api "o.o/api/top/int/shop"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/authorize/session"
)

type CategoryService struct {
	session.Session

	CatalogQuery catalog.QueryBus
	CatalogAggr  catalog.CommandBus
}

func (s *CategoryService) Clone() api.CategoryService { res := *s; return &res }

func (s *CategoryService) CreateCategory(ctx context.Context, q *api.CreateCategoryRequest) (*api.ShopCategory, error) {
	cmd := &catalog.CreateShopCategoryCommand{
		ShopID:   s.SS.Shop().ID,
		Name:     q.Name,
		ParentID: q.ParentId,
	}
	if err := s.CatalogAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := PbShopCategory(cmd.Result)
	return result, nil
}

func (s *CategoryService) GetCategory(ctx context.Context, q *pbcm.IDRequest) (*api.ShopCategory, error) {
	query := &catalog.GetShopCategoryQuery{
		ID:     q.Id,
		ShopID: s.SS.Shop().ID,
	}
	if err := s.CatalogQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := PbShopCategory(query.Result)
	return result, nil
}

func (s *CategoryService) GetCategories(ctx context.Context, q *api.GetCategoriesRequest) (*api.ShopCategoriesResponse, error) {
	paging := cmapi.CMPaging(q.Paging)
	query := &catalog.ListShopCategoriesQuery{
		ShopID:  s.SS.Shop().ID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(q.Filters),
	}
	if err := s.CatalogQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}

	result := &api.ShopCategoriesResponse{
		Paging:     cmapi.PbPageInfo(paging),
		Categories: PbShopCategories(query.Result.Categories),
	}
	return result, nil
}

func (s *CategoryService) UpdateCategory(ctx context.Context, q *api.UpdateCategoryRequest) (*api.ShopCategory, error) {
	shopID := s.SS.Shop().ID
	cmd := &catalog.UpdateShopCategoryCommand{
		ID:       q.Id,
		ShopID:   shopID,
		Name:     q.Name,
		ParentID: q.ParentId,
	}
	if err := s.CatalogAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := PbShopCategory(cmd.Result)
	return result, nil
}

func (s *CategoryService) DeleteCategory(ctx context.Context, r *pbcm.IDRequest) (*pbcm.DeletedResponse, error) {
	cmd := &catalog.DeleteShopCategoryCommand{
		ID:     r.Id,
		ShopID: s.SS.Shop().ID,
	}
	if err := s.CatalogAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.DeletedResponse{Deleted: cmd.Result}
	return result, nil
}
