package shop

import (
	"context"

	"o.o/api/main/catalog"
	"o.o/api/main/inventory"
	"o.o/api/top/int/shop"
	"o.o/api/webserver"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/capi/dot"
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
	err := s.WebserverAggr.Dispatch(ctx, cmd)
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
	err := s.WebserverQuery.Dispatch(ctx, query)
	if err != nil {
		return err
	}
	result, err := getProductQuantity(ctx, s.InventoryQuery, shopID, query.Result.Product)
	if err != nil {
		return err
	}
	r.Result = PbWsProduct(query.Result)
	r.Result.Product = result
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
	err := s.WebserverQuery.Dispatch(ctx, query)
	if err != nil {
		return err
	}
	mapProductWithQuantity, err := mapListProductsWithQuantity(ctx, s.InventoryQuery, query.Result.WsProducts)
	if err != nil {
		return err
	}
	resultWsProducts := PbWsProducts(query.Result.WsProducts)
	for k, v := range resultWsProducts {
		resultWsProducts[k].Product = mapProductWithQuantity[v.ID]
	}
	r.Result = &shop.GetWsProductsResponse{
		WsProducts: resultWsProducts,
		Paging:     cmapi.PbPaging(query.Paging),
	}
	return nil
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
	result, err := getProductsQuantity(ctx, inventoryQuery, args[0].ShopID, products)
	if err != nil {
		return nil, err
	}
	for _, v := range result {
		mapProduct[v.Id] = v
	}
	return mapProduct, nil
}

func (s *WebServerService) GetWsProductsByIDs(ctx context.Context, r *GetWsProductsByIDsEndpoint) error {
	shopID := r.Context.Shop.ID
	query := &webserver.ListWsProductsByIDsQuery{
		ShopID: shopID,
		IDs:    r.IDs,
	}
	err := s.WebserverQuery.Dispatch(ctx, query)
	if err != nil {
		return err
	}
	r.Result = &shop.GetWsProductsByIDsResponse{
		WsProducts: PbWsProducts(query.Result),
	}
	return nil
}
