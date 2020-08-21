package ws

import (
	"context"

	"o.o/api/main/catalog"
	"o.o/api/main/inventory"
	"o.o/api/top/int/shop"
	api "o.o/api/top/int/shop"
	"o.o/api/webserver"
	"o.o/backend/pkg/common/apifw/cmapi"
	shop2 "o.o/backend/pkg/etop/api/shop"
	"o.o/backend/pkg/etop/api/shop/product"
	"o.o/capi/dot"
)

func (s *WebServerService) CreateOrUpdateWsProduct(ctx context.Context, r *api.CreateOrUpdateWsProductRequest) (*api.WsProduct, error) {
	shopID := s.SS.Shop().ID
	cmd := &webserver.CreateOrUpdateWsProductCommand{
		ID:           r.ProductID,
		ShopID:       shopID,
		SEOConfig:    shop2.ConvertSEOConfig(r.SEOConfig),
		Slug:         r.Slug,
		Appear:       r.Appear,
		ComparePrice: shop2.ConvertComparePrice(r.ComparePrices),
		DescHTML:     r.DescHTML,
	}
	err := s.WebserverAggr.Dispatch(ctx, cmd)
	if err != nil {
		return nil, err
	}
	result := PbWsProduct(cmd.Result)
	return result, nil
}

func (s *WebServerService) GetWsProduct(ctx context.Context, r *api.GetWsProductRequest) (*api.WsProduct, error) {
	shopID := s.SS.Shop().ID
	query := &webserver.GetWsProductByIDQuery{
		ID:     r.ID,
		ShopID: shopID,
	}
	err := s.WebserverQuery.Dispatch(ctx, query)
	if err != nil {
		return nil, err
	}
	resp, err := product.GetProductQuantity(ctx, s.InventoryQuery, shopID, query.Result.Product)
	if err != nil {
		return nil, err
	}
	result := PbWsProduct(query.Result)
	result.Product = resp
	return result, nil
}

func (s *WebServerService) GetWsProducts(ctx context.Context, r *api.GetWsProductsRequest) (*api.GetWsProductsResponse, error) {
	shopID := s.SS.Shop().ID
	paging := cmapi.CMPaging(r.Paging)
	query := &webserver.ListWsProductsQuery{
		ShopID:  shopID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(r.Filters),
	}
	err := s.WebserverQuery.Dispatch(ctx, query)
	if err != nil {
		return nil, err
	}
	mapProductWithQuantity, err := mapListProductsWithQuantity(ctx, s.InventoryQuery, query.Result.WsProducts)
	if err != nil {
		return nil, err
	}
	resultWsProducts := PbWsProducts(query.Result.WsProducts)
	for k, v := range resultWsProducts {
		resultWsProducts[k].Product = mapProductWithQuantity[v.ID]
	}
	result := &api.GetWsProductsResponse{
		WsProducts: resultWsProducts,
		Paging:     cmapi.PbPaging(query.Paging),
	}
	return result, nil
}

func mapListProductsWithQuantity(ctx context.Context, inventoryQuery inventory.QueryBus, args []*webserver.WsProduct) (map[dot.ID]*shop.ShopProduct, error) {
	if len(args) == 0 {
		return nil, nil
	}
	var products []*catalog.ShopProductWithVariants
	for _, v := range args {
		products = append(products, v.Product)
	}
	var mapProduct = make(map[dot.ID]*shop.ShopProduct)
	result, err := product.GetProductsQuantity(ctx, inventoryQuery, args[0].ShopID, products)
	if err != nil {
		return nil, err
	}
	for _, v := range result {
		mapProduct[v.Id] = v
	}
	return mapProduct, nil
}

func (s *WebServerService) GetWsProductsByIDs(ctx context.Context, r *api.GetWsProductsByIDsRequest) (*api.GetWsProductsByIDsResponse, error) {
	shopID := s.SS.Shop().ID
	query := &webserver.ListWsProductsByIDsQuery{
		ShopID: shopID,
		IDs:    r.IDs,
	}
	err := s.WebserverQuery.Dispatch(ctx, query)
	if err != nil {
		return nil, err
	}
	result := &api.GetWsProductsByIDsResponse{
		WsProducts: PbWsProducts(query.Result),
	}
	return result, nil
}
