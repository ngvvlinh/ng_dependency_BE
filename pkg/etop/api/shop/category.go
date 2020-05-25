package shop

import (
	"context"

	"o.o/api/main/catalog"
	"o.o/api/top/int/shop"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/common/apifw/cmapi"
)

type CategoryService struct {
	CatalogQuery catalog.QueryBus
	CatalogAggr  catalog.CommandBus
}

func (s *CategoryService) Clone() *CategoryService { res := *s; return &res }

func (s *CategoryService) CreateCategory(ctx context.Context, q *CreateCategoryEndpoint) error {
	cmd := &catalog.CreateShopCategoryCommand{
		ShopID:   q.Context.Shop.ID,
		Name:     q.Name,
		ParentID: q.ParentId,
	}
	if err := s.CatalogAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = PbShopCategory(cmd.Result)
	return nil
}

func (s *CategoryService) GetCategory(ctx context.Context, q *GetCategoryEndpoint) error {
	query := &catalog.GetShopCategoryQuery{
		ID:     q.Id,
		ShopID: q.Context.Shop.ID,
	}
	if err := s.CatalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = PbShopCategory(query.Result)
	return nil
}

func (s *CategoryService) GetCategories(ctx context.Context, q *GetCategoriesEndpoint) error {
	paging := cmapi.CMPaging(q.Paging)
	query := &catalog.ListShopCategoriesQuery{
		ShopID:  q.Context.Shop.ID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(q.Filters),
	}
	if err := s.CatalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = &shop.ShopCategoriesResponse{
		Paging:     cmapi.PbPageInfo(paging),
		Categories: PbShopCategories(query.Result.Categories),
	}
	return nil
}

func (s *CategoryService) UpdateCategory(ctx context.Context, q *UpdateCategoryEndpoint) error {
	shopID := q.Context.Shop.ID
	cmd := &catalog.UpdateShopCategoryCommand{
		ID:       q.Id,
		ShopID:   shopID,
		Name:     q.Name,
		ParentID: q.ParentId,
	}
	if err := s.CatalogAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = PbShopCategory(cmd.Result)
	return nil
}

func (s *CategoryService) DeleteCategory(ctx context.Context, r *DeleteCategoryEndpoint) error {
	cmd := &catalog.DeleteShopCategoryCommand{
		ID:     r.Id,
		ShopID: r.Context.Shop.ID,
	}
	if err := s.CatalogAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.DeletedResponse{Deleted: cmd.Result}
	return nil
}
