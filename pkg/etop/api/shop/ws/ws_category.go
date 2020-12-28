package ws

import (
	"context"

	"o.o/api/main/catalog"
	"o.o/api/top/int/shop"
	api "o.o/api/top/int/shop"
	"o.o/api/webserver"
	"o.o/backend/pkg/common/apifw/cmapi"
	convertpball "o.o/backend/pkg/etop/api/convertpb/_all"
	"o.o/capi/dot"
)

func (s *WebServerService) CreateOrUpdateWsCategory(ctx context.Context, r *api.CreateOrUpdateWsCategoryRequest) (*api.WsCategory, error) {
	shopID := s.SS.Shop().ID
	cmd := &webserver.CreateOrUpdateWsCategoryCommand{
		ID:        r.CategoryID,
		ShopID:    shopID,
		SEOConfig: convertpball.ConvertSEOConfig(r.SEOConfig),
		Slug:      r.Slug,
		Appear:    r.Appear,
		Image:     r.Image,
	}
	err := s.WebserverAggr.Dispatch(ctx, cmd)
	if err != nil {
		return nil, err
	}
	result := PbWsCategory(cmd.Result)
	result, err = s.populateWsCategoryWithProductCount(ctx, result)
	return result, err
}

func (s *WebServerService) GetWsCategory(ctx context.Context, r *api.GetWsCategoryRequest) (*api.WsCategory, error) {
	shopID := s.SS.Shop().ID
	query := &webserver.GetWsCategoryByIDQuery{
		ID:     r.ID,
		ShopID: shopID,
	}
	err := s.WebserverQuery.Dispatch(ctx, query)
	if err != nil {
		return nil, err
	}
	result := PbWsCategory(query.Result)
	result, err = s.populateWsCategoryWithProductCount(ctx, result)
	return nil, err
}

func (s *WebServerService) GetWsCategories(ctx context.Context, r *api.GetWsCategoriesRequest) (*api.GetWsCategoriesResponse, error) {
	shopID := s.SS.Shop().ID
	paging := cmapi.CMPaging(r.Paging)
	query := &webserver.ListWsCategoriesQuery{
		ShopID:  shopID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(r.Filters),
	}
	err := s.WebserverQuery.Dispatch(ctx, query)
	if err != nil {
		return nil, err
	}
	result := &api.GetWsCategoriesResponse{
		WsCategories: PbWsCategories(query.Result.WsCategories),
		Paging:       cmapi.PbPaging(query.Paging),
	}
	result.WsCategories, err = s.populateWsCategoriesWithProductCount(ctx, result.WsCategories)
	return nil, err
}

func (s *WebServerService) GetWsCategoriesByIDs(ctx context.Context, r *api.GetWsCategoriesByIDsRequest) (*api.GetWsCategoriesByIDsResponse, error) {
	shopID := s.SS.Shop().ID
	query := &webserver.ListWsCategoriesByIDsQuery{
		ShopID: shopID,
		IDs:    r.IDs,
	}
	err := s.WebserverQuery.Dispatch(ctx, query)
	if err != nil {
		return nil, err
	}
	result := &api.GetWsCategoriesByIDsResponse{
		WsCategories: PbWsCategories(query.Result),
	}
	result.WsCategories, err = s.populateWsCategoriesWithProductCount(ctx, result.WsCategories)
	return nil, err
}

func (s *WebServerService) populateWsCategoriesWithProductCount(ctx context.Context, args []*shop.WsCategory) ([]*shop.WsCategory, error) {
	if len(args) == 0 {
		return []*shop.WsCategory{}, nil
	}
	var categoriesIDs []dot.ID
	for _, v := range args {
		categoriesIDs = append(categoriesIDs, v.ID)
	}
	query := &catalog.ListShopProductWithVariantByCategoriesIDsQuery{
		ShopID:        args[0].ShopID,
		CategoriesIds: categoriesIDs,
	}
	err := s.CatalogQuery.Dispatch(ctx, query)
	if err != nil {
		return nil, err
	}
	var mapCategoriesCountProduct = make(map[dot.ID]int)
	for _, v := range query.Result.Products {
		mapCategoriesCountProduct[v.CategoryID]++
	}
	for k, v := range args {
		args[k].ProductCount = mapCategoriesCountProduct[v.ID]
	}
	return args, nil
}

func (s *WebServerService) populateWsCategoryWithProductCount(ctx context.Context, args *shop.WsCategory) (*shop.WsCategory, error) {
	if args == nil {
		return nil, nil
	}
	query := &catalog.ListShopProductWithVariantByCategoriesIDsQuery{
		ShopID:        args.ShopID,
		CategoriesIds: []dot.ID{args.ID},
	}
	err := s.CatalogQuery.Dispatch(ctx, query)
	if err != nil {
		return nil, err
	}
	args.ProductCount = len(query.Result.Products)
	return args, nil
}
