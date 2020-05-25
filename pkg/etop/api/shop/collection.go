package shop

import (
	"context"

	"o.o/api/main/catalog"
	"o.o/api/top/int/shop"
	"o.o/backend/pkg/common/apifw/cmapi"
)

type CollectionService struct {
	CatalogQuery catalog.QueryBus
	CatalogAggr  catalog.CommandBus
}

func (s *CollectionService) Clone() *CollectionService { res := *s; return &res }

func (s *CollectionService) GetCollection(ctx context.Context, q *GetCollectionEndpoint) error {
	query := &catalog.GetShopCollectionQuery{
		ID:     q.Id,
		ShopID: q.Context.Shop.ID,
	}
	if err := s.CatalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = PbShopCollection(query.Result)
	return nil
}

func (s *CollectionService) GetCollections(ctx context.Context, q *GetCollectionsEndpoint) error {
	paging := cmapi.CMPaging(q.Paging)
	query := &catalog.ListShopCollectionsQuery{
		ShopID:  q.Context.Shop.ID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(q.Filters),
	}
	if err := s.CatalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &shop.ShopCollectionsResponse{
		Paging:      cmapi.PbPageInfo(paging),
		Collections: PbShopCollections(query.Result.Collections),
	}
	return nil
}

func (s *CollectionService) UpdateCollection(ctx context.Context, q *UpdateCollectionEndpoint) error {
	shopID := q.Context.Shop.ID
	cmd := &catalog.UpdateShopCollectionCommand{
		ID:          q.Id,
		ShopID:      shopID,
		Name:        q.Name,
		Description: q.Description,
		DescHTML:    q.DescHtml,
		ShortDesc:   q.ShortDesc,
	}
	if err := s.CatalogAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = PbShopCollection(cmd.Result)
	return nil
}

func (s *CollectionService) CreateCollection(ctx context.Context, q *CreateCollectionEndpoint) error {
	cmd := &catalog.CreateShopCollectionCommand{
		ShopID:      q.Context.Shop.ID,
		Name:        q.Name,
		DescHTML:    q.DescHtml,
		Description: q.Description,
		ShortDesc:   q.ShortDesc,
	}
	if err := s.CatalogAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = PbShopCollection(cmd.Result)
	return nil
}

func (s *CollectionService) GetCollectionsByProductID(ctx context.Context, q *GetCollectionsByProductIDEndpoint) error {
	query := &catalog.ListShopCollectionsByProductIDQuery{
		ShopID:    q.Context.Shop.ID,
		ProductID: q.ProductId,
	}
	if err := s.CatalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &shop.CollectionsResponse{
		Collections: PbShopCollections(query.Result),
	}
	return nil
}
