package shop

import (
	"context"

	"o.o/api/top/int/shop"
	"o.o/api/webserver"
	"o.o/backend/pkg/common/apifw/cmapi"
)

func (s *WebServerService) CreateOrUpdateWsCategory(ctx context.Context, r *CreateOrUpdateWsCategoryEndpoint) error {
	shopID := r.Context.Shop.ID
	cmd := &webserver.CreateOrUpdateWsCategoryCommand{
		ID:        r.CategoryID,
		ShopID:    shopID,
		SEOConfig: ConvertSEOConfig(r.SEOConfig),
		Slug:      r.Slug,
		Appear:    r.Appear,
		Image:     r.Image,
	}
	err := webserverAggr.Dispatch(ctx, cmd)
	if err != nil {
		return err
	}
	r.Result = PbWsCategory(cmd.Result)
	return nil
}

func (s *WebServerService) GetWsCategory(ctx context.Context, r *GetWsCategoryEndpoint) error {
	shopID := r.Context.Shop.ID
	query := &webserver.GetWsCategoryByIDQuery{
		ID:     r.ID,
		ShopID: shopID,
	}
	err := webserverQuery.Dispatch(ctx, query)
	if err != nil {
		return err
	}
	r.Result = PbWsCategory(query.Result)
	return nil
}

func (s *WebServerService) GetWsCategories(ctx context.Context, r *GetWsCategoriesEndpoint) error {
	shopID := r.Context.Shop.ID
	paging := cmapi.CMPaging(r.Paging)
	query := &webserver.ListWsCategoriesQuery{
		ShopID:  shopID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(r.Filters),
	}
	err := webserverQuery.Dispatch(ctx, query)
	if err != nil {
		return err
	}
	r.Result = &shop.GetWsCategoriesResponse{
		WsCategories: PbWsCategories(query.Result.WsCategories),
		Paging:       cmapi.PbPaging(query.Paging),
	}
	return nil
}

func (s *WebServerService) GetWsCategoriesByIDs(ctx context.Context, r *GetWsCategoriesByIDsEndpoint) error {
	shopID := r.Context.Shop.ID
	query := &webserver.ListWsCategoriesByIDsQuery{
		ShopID: shopID,
		IDs:    r.IDs,
	}
	err := webserverQuery.Dispatch(ctx, query)
	if err != nil {
		return err
	}
	r.Result = &shop.GetWsCategoriesByIDsResponse{
		WsCategories: PbWsCategories(query.Result),
	}
	return nil
}
