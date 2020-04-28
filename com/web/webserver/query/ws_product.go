package query

import (
	"context"

	"o.o/api/main/catalog"
	"o.o/api/webserver"
	cm "o.o/backend/pkg/common"
	"o.o/capi/dot"
)

func (w WebserverQueryService) GetWsProductByID(ctx context.Context, shopID dot.ID, ID dot.ID) (*webserver.WsProduct, error) {
	if shopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Mising shop_id")
	}
	queryProduct := &catalog.GetShopProductWithVariantsByIDQuery{
		ProductID: ID,
		ShopID:    shopID,
	}
	if err := w.cataglogQuery.Dispatch(ctx, queryProduct); err != nil {
		return nil, err
	}
	result, err := w.wsProductStore(ctx).ShopID(shopID).ID(ID).GetWsProduct()
	switch cm.ErrorCode(err) {
	case cm.NoError:
		result.Product = queryProduct.Result
	case cm.NotFound:
		result = &webserver.WsProduct{
			ID:        ID,
			ShopID:    shopID,
			Appear:    true,
			CreatedAt: queryProduct.Result.CreatedAt,
			UpdatedAt: queryProduct.Result.UpdatedAt,
			Product:   queryProduct.Result,
		}
	default:
		return nil, err
	}
	return result, nil
}

func (w WebserverQueryService) ListWsProductsByIDs(ctx context.Context, shopID dot.ID, IDs []dot.ID) ([]*webserver.WsProduct, error) {
	if shopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Mising shop_id")
	}
	queryProduct := &catalog.ListShopProductsWithVariantsByIDsQuery{
		IDs:    IDs,
		ShopID: shopID,
	}
	if err := w.cataglogQuery.Dispatch(ctx, queryProduct); err != nil {
		return nil, err
	}
	return w.getWsProducts(ctx, queryProduct.Result.Products)
}

func (w WebserverQueryService) ListWsProducts(ctx context.Context, args webserver.ListWsProductsArgs) (*webserver.ListWsProductsResponse, error) {
	if args.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Mising shop_id")
	}
	queryProduct := &catalog.ListShopProductsWithVariantsQuery{
		ShopID:  args.ShopID,
		Paging:  args.Paging,
		Filters: args.Filters,
	}
	if err := w.cataglogQuery.Dispatch(ctx, queryProduct); err != nil {
		return nil, err
	}
	wsProducts, err := w.getWsProducts(ctx, queryProduct.Result.Products)
	if err != nil {
		return nil, err
	}
	return &webserver.ListWsProductsResponse{
		PageInfo:   queryProduct.Result.Paging,
		WsProducts: wsProducts,
	}, nil
}

func (w WebserverQueryService) getWsProducts(ctx context.Context, shopProducts []*catalog.ShopProductWithVariants) ([]*webserver.WsProduct, error) {
	if len(shopProducts) < 1 {
		return []*webserver.WsProduct{}, nil
	}
	var productIDs []dot.ID
	for _, v := range shopProducts {
		productIDs = append(productIDs, v.ProductID)
	}
	wsProducts, err := w.wsProductStore(ctx).ShopID(shopProducts[0].ShopID).IDs(productIDs...).ListWsProducts()
	if err != nil {
		return nil, err
	}
	mapShopProducts := make(map[dot.ID]*webserver.WsProduct)
	for _, wsProduct := range wsProducts {
		mapShopProducts[wsProduct.ID] = wsProduct
	}
	var wsProductsResult []*webserver.WsProduct
	for _, shopProduct := range shopProducts {
		if mapShopProducts[shopProduct.ProductID] != nil {
			mapShopProducts[shopProduct.ProductID].Product = shopProduct
			wsProductsResult = append(wsProductsResult, mapShopProducts[shopProduct.ProductID])
			continue
		}
		wsProductsResult = append(wsProductsResult, &webserver.WsProduct{
			ID:      shopProduct.ProductID,
			ShopID:  shopProduct.ShopID,
			Product: shopProduct,
		})
	}
	return wsProductsResult, nil
}
