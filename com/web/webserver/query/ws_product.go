package query

import (
	"context"
	"sort"

	"o.o/api/main/catalog"
	"o.o/api/meta"
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
	var mapVariants = make(map[dot.ID]*catalog.ShopVariant)
	for _, v := range queryProduct.Result.Variants {
		mapVariants[v.VariantID] = v
	}
	for _, v := range result.ComparePrice {
		if mapVariants[v.VariantID].RetailPrice > v.ComparePrice {
			result.IsSale = true
			break
		}
	}
	result.Product = sortVariantByRetailPrice(result.Product)
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

func (w WebserverQueryService) ListWsProductsByIDsWithPaging(ctx context.Context, shopID dot.ID, IDs []dot.ID, paging meta.Paging) (*webserver.ListWsProductsResponse, error) {
	if shopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Mising shop_id")
	}
	queryProduct := &catalog.ListShopProductWithVariantByIDsWithPagingQuery{
		ShopID: shopID,
		IDs:    IDs,
		Paging: paging,
	}
	if err := w.cataglogQuery.Dispatch(ctx, queryProduct); err != nil {
		return nil, err
	}
	result, err := w.getWsProducts(ctx, queryProduct.Result.Products)
	if err != nil {
		return nil, err
	}
	return &webserver.ListWsProductsResponse{
		ShopID:     shopID,
		PageInfo:   queryProduct.Result.Paging,
		WsProducts: result,
	}, nil
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
		} else {
			wsProductsResult = append(wsProductsResult, &webserver.WsProduct{
				ID:      shopProduct.ProductID,
				ShopID:  shopProduct.ShopID,
				Product: shopProduct,
				Appear:  true,
			})
		}
	}
	for key, product := range wsProductsResult {
		var mapVariants = make(map[dot.ID]*catalog.ShopVariant)
		for _, variant := range product.Product.Variants {
			mapVariants[variant.VariantID] = variant
		}
		for _, v := range product.ComparePrice {
			if mapVariants[v.VariantID].RetailPrice > v.ComparePrice {
				wsProductsResult[key].IsSale = true
				break
			}
		}
	}
	for k, v := range wsProductsResult {
		wsProductsResult[k].Product = sortVariantByRetailPrice(v.Product)
	}
	return wsProductsResult, nil
}

func (w WebserverQueryService) SearchProductByName(ctx context.Context, shopID dot.ID, name string) (*webserver.ListWsProductsResponse, error) {
	if shopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Mising shop_id")
	}
	queryProduct := &catalog.SearchProductByNameQuery{
		ShopID: shopID,
		Name:   name,
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

func sortVariantByRetailPrice(product *catalog.ShopProductWithVariants) *catalog.ShopProductWithVariants {
	sort.Slice(product.Variants, func(i, j int) bool {
		return product.Variants[i].RetailPrice < product.Variants[j].RetailPrice
	})
	return product
}
