package shop

import (
	"context"

	"o.o/api/top/int/shop"
	"o.o/api/webserver"
	"o.o/backend/pkg/common/apifw/cmapi"
)

func (s *WebServerService) CreateWsPage(ctx context.Context, r *CreateWsPageEndpoint) error {
	shopID := r.Context.Shop.ID
	cmd := &webserver.CreateWsPageCommand{
		ShopID:    shopID,
		SEOConfig: ConvertSEOConfig(r.SEOConfig),
		Name:      r.Name,
		Slug:      r.Slug,
		DescHTML:  r.DescHTML,
		Image:     r.Image,
		Appear:    r.Appear,
	}
	err := s.WebserverAggr.Dispatch(ctx, cmd)
	if err != nil {
		return err
	}
	r.Result = PbWsPage(cmd.Result)
	return nil
}

func (s *WebServerService) UpdateWsPage(ctx context.Context, r *UpdateWsPageEndpoint) error {
	shopID := r.Context.Shop.ID
	cmd := &webserver.UpdateWsPageCommand{
		ShopID:    shopID,
		ID:        r.ID,
		SEOConfig: ConvertSEOConfig(r.SEOConfig),
		Name:      r.Name,
		Slug:      r.Slug,
		DescHTML:  r.DescHTML,
		Image:     r.Image,
		Appear:    r.Appear,
	}
	err := s.WebserverAggr.Dispatch(ctx, cmd)
	if err != nil {
		return err
	}
	r.Result = PbWsPage(cmd.Result)
	return nil
}

func (s *WebServerService) DeleteWsPage(ctx context.Context, r *DeleteWsPageEndpoint) error {
	shopID := r.Context.Shop.ID
	cmd := &webserver.DeleteWsPageCommand{
		ID:     r.ID,
		ShopID: shopID,
	}
	err := s.WebserverAggr.Dispatch(ctx, cmd)
	if err != nil {
		return err
	}
	r.Result = &shop.DeteleWsPageResponse{
		Count: cmd.Result,
	}
	return nil
}

func (s *WebServerService) GetWsPage(ctx context.Context, r *GetWsPageEndpoint) error {
	shopID := r.Context.Shop.ID
	query := &webserver.GetWsPageByIDQuery{
		ID:     r.ID,
		ShopID: shopID,
	}
	err := s.WebserverQuery.Dispatch(ctx, query)
	if err != nil {
		return err
	}
	r.Result = PbWsPage(query.Result)
	return nil
}

func (s *WebServerService) GetWsPages(ctx context.Context, r *GetWsPagesEndpoint) error {
	shopID := r.Context.Shop.ID
	paging := cmapi.CMPaging(r.Paging)
	query := &webserver.ListWsPagesQuery{
		ShopID:  shopID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(r.Filters),
	}
	err := s.WebserverQuery.Dispatch(ctx, query)
	if err != nil {
		return err
	}
	r.Result = &shop.GetWsPagesResponse{
		WsPages: PbWsPages(query.Result.WsPages),
		Paging:  cmapi.PbPaging(query.Paging),
	}
	return nil
}

func (s *WebServerService) GetWsPagesByIDs(ctx context.Context, r *GetWsPagesByIDsEndpoint) error {
	shopID := r.Context.Shop.ID
	query := &webserver.ListWsPagesByIDsQuery{
		ShopID: shopID,
		IDs:    r.IDs,
	}
	err := s.WebserverQuery.Dispatch(ctx, query)
	if err != nil {
		return err
	}
	r.Result = &shop.GetWsPagesByIDsResponse{
		WsPages: PbWsPages(query.Result),
	}
	return nil
}
