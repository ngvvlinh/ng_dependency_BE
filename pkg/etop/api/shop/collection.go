package shop

import (
	"context"

	"o.o/api/main/catalog"
	api "o.o/api/top/int/shop"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/authorize/session"
)

type CollectionService struct {
	session.Session

	CatalogQuery catalog.QueryBus
	CatalogAggr  catalog.CommandBus
}

func (s *CollectionService) Clone() api.CollectionService { res := *s; return &res }

func (s *CollectionService) GetCollection(ctx context.Context, q *pbcm.IDRequest) (*api.ShopCollection, error) {
	query := &catalog.GetShopCollectionQuery{
		ID:     q.Id,
		ShopID: s.SS.Shop().ID,
	}
	if err := s.CatalogQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := PbShopCollection(query.Result)
	return result, nil
}

func (s *CollectionService) GetCollections(ctx context.Context, q *api.GetCollectionsRequest) (*api.ShopCollectionsResponse, error) {
	paging := cmapi.CMPaging(q.Paging)
	query := &catalog.ListShopCollectionsQuery{
		ShopID:  s.SS.Shop().ID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(q.Filters),
	}
	if err := s.CatalogQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &api.ShopCollectionsResponse{
		Paging:      cmapi.PbPageInfo(paging),
		Collections: PbShopCollections(query.Result.Collections),
	}
	return result, nil
}

func (s *CollectionService) UpdateCollection(ctx context.Context, q *api.UpdateCollectionRequest) (*api.ShopCollection, error) {
	shopID := s.SS.Shop().ID
	cmd := &catalog.UpdateShopCollectionCommand{
		ID:          q.Id,
		ShopID:      shopID,
		Name:        q.Name,
		Description: q.Description,
		DescHTML:    q.DescHtml,
		ShortDesc:   q.ShortDesc,
	}
	if err := s.CatalogAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := PbShopCollection(cmd.Result)
	return result, nil
}

func (s *CollectionService) CreateCollection(ctx context.Context, q *api.CreateCollectionRequest) (*api.ShopCollection, error) {
	cmd := &catalog.CreateShopCollectionCommand{
		ShopID:      s.SS.Shop().ID,
		Name:        q.Name,
		DescHTML:    q.DescHtml,
		Description: q.Description,
		ShortDesc:   q.ShortDesc,
	}
	if err := s.CatalogAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := PbShopCollection(cmd.Result)
	return result, nil
}

func (s *CollectionService) GetCollectionsByProductID(ctx context.Context, q *api.GetShopCollectionsByProductIDRequest) (*api.CollectionsResponse, error) {
	query := &catalog.ListShopCollectionsByProductIDQuery{
		ShopID:    s.SS.Shop().ID,
		ProductID: q.ProductId,
	}
	if err := s.CatalogQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &api.CollectionsResponse{
		Collections: PbShopCollections(query.Result),
	}
	return result, nil
}
