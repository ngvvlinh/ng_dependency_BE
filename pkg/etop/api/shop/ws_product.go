package shop

import (
	"context"

	"etop.vn/api/top/int/shop"
	"etop.vn/api/webserver"
	"etop.vn/backend/pkg/common/apifw/cmapi"
)

func (s *WebServerService) CreateOrUpdateWsProduct(ctx context.Context, r *CreateOrUpdateWsProductEndpoint) error {
	shopID := r.Context.Shop.ID
	cmd := &webserver.CreateOrUpdateWsProductCommand{
		ID:           r.ProductID,
		ShopID:       shopID,
		SEOConfig:    ConvertSEOConfig(r.SEOConfig),
		Slug:         r.Slug,
		Appear:       r.Appear,
		ComparePrice: ConvertComparePrice(r.ComparePrices),
		DescHTML:     r.DescHTML,
	}
	err := webserverAggr.Dispatch(ctx, cmd)
	if err != nil {
		return err
	}
	r.Result = PbWsProduct(cmd.Result)
	return nil
}

func (s *WebServerService) GetWsProduct(ctx context.Context, r *GetWsProductEndpoint) error {
	shopID := r.Context.Shop.ID
	query := &webserver.GetWsProductByIDQuery{
		ID:     r.ID,
		ShopID: shopID,
	}
	err := webserverQuery.Dispatch(ctx, query)
	if err != nil {
		return err
	}
	r.Result = PbWsProduct(query.Result)
	return nil
}

func (s *WebServerService) GetWsProducts(ctx context.Context, r *GetWsProductsEndpoint) error {
	shopID := r.Context.Shop.ID
	paging := cmapi.CMPaging(r.Paging)
	query := &webserver.ListWsProductsQuery{
		ShopID:  shopID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(r.Filters),
	}
	err := webserverQuery.Dispatch(ctx, query)
	if err != nil {
		return err
	}
	r.Result = &shop.GetWsProductsResponse{
		WsProducts: PbWsProducts(query.Result.WsProducts),
		Paging:     cmapi.PbPaging(query.Paging),
	}
	return nil
}

func (s *WebServerService) GetWsProductsByIDs(ctx context.Context, r *GetWsProductsByIDsEndpoint) error {
	shopID := r.Context.Shop.ID
	query := &webserver.ListWsProductsByIDsQuery{
		ShopID: shopID,
		IDs:    r.IDs,
	}
	err := webserverQuery.Dispatch(ctx, query)
	if err != nil {
		return err
	}
	r.Result = &shop.GetWsProductsByIDsResponse{
		WsProducts: PbWsProducts(query.Result),
	}
	return nil
}
