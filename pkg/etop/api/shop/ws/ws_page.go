package ws

import (
	"context"

	api "o.o/api/top/int/shop"
	"o.o/api/webserver"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/api/shop"
)

func (s *WebServerService) CreateWsPage(ctx context.Context, r *api.CreateWsPageRequest) (*api.WsPage, error) {
	shopID := s.SS.Shop().ID
	cmd := &webserver.CreateWsPageCommand{
		ShopID:    shopID,
		SEOConfig: shop.ConvertSEOConfig(r.SEOConfig),
		Name:      r.Name,
		Slug:      r.Slug,
		DescHTML:  r.DescHTML,
		Image:     r.Image,
		Appear:    r.Appear,
	}
	err := s.WebserverAggr.Dispatch(ctx, cmd)
	if err != nil {
		return nil, err
	}
	result := PbWsPage(cmd.Result)
	return result, nil
}

func (s *WebServerService) UpdateWsPage(ctx context.Context, r *api.UpdateWsPageRequest) (*api.WsPage, error) {
	shopID := s.SS.Shop().ID
	cmd := &webserver.UpdateWsPageCommand{
		ShopID:    shopID,
		ID:        r.ID,
		SEOConfig: shop.ConvertSEOConfig(r.SEOConfig),
		Name:      r.Name,
		Slug:      r.Slug,
		DescHTML:  r.DescHTML,
		Image:     r.Image,
		Appear:    r.Appear,
	}
	err := s.WebserverAggr.Dispatch(ctx, cmd)
	if err != nil {
		return nil, err
	}
	result := PbWsPage(cmd.Result)
	return result, nil
}

func (s *WebServerService) DeleteWsPage(ctx context.Context, r *api.DeteleWsPageRequest) (*api.DeteleWsPageResponse, error) {
	shopID := s.SS.Shop().ID
	cmd := &webserver.DeleteWsPageCommand{
		ID:     r.ID,
		ShopID: shopID,
	}
	err := s.WebserverAggr.Dispatch(ctx, cmd)
	if err != nil {
		return nil, err
	}
	result := &api.DeteleWsPageResponse{
		Count: cmd.Result,
	}
	return result, nil
}

func (s *WebServerService) GetWsPage(ctx context.Context, r *api.GetWsPageRequest) (*api.WsPage, error) {
	shopID := s.SS.Shop().ID
	query := &webserver.GetWsPageByIDQuery{
		ID:     r.ID,
		ShopID: shopID,
	}
	err := s.WebserverQuery.Dispatch(ctx, query)
	if err != nil {
		return nil, err
	}
	result := PbWsPage(query.Result)
	return result, nil
}

func (s *WebServerService) GetWsPages(ctx context.Context, r *api.GetWsPagesRequest) (*api.GetWsPagesResponse, error) {
	shopID := s.SS.Shop().ID
	paging := cmapi.CMPaging(r.Paging)
	query := &webserver.ListWsPagesQuery{
		ShopID:  shopID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(r.Filters),
	}
	err := s.WebserverQuery.Dispatch(ctx, query)
	if err != nil {
		return nil, err
	}
	result := &api.GetWsPagesResponse{
		WsPages: PbWsPages(query.Result.WsPages),
		Paging:  cmapi.PbPaging(query.Paging),
	}
	return result, nil
}

func (s *WebServerService) GetWsPagesByIDs(ctx context.Context, r *api.GetWsPagesByIDsRequest) (*api.GetWsPagesByIDsResponse, error) {
	shopID := s.SS.Shop().ID
	query := &webserver.ListWsPagesByIDsQuery{
		ShopID: shopID,
		IDs:    r.IDs,
	}
	err := s.WebserverQuery.Dispatch(ctx, query)
	if err != nil {
		return nil, err
	}
	result := &api.GetWsPagesByIDsResponse{
		WsPages: PbWsPages(query.Result),
	}
	return result, nil
}
