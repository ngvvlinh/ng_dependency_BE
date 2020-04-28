package query

import (
	"context"

	"o.o/api/main/catalog"
	"o.o/api/webserver"
	cm "o.o/backend/pkg/common"
	"o.o/capi/dot"
)

func (w WebserverQueryService) GetWsCategoryByID(ctx context.Context, shopID dot.ID, ID dot.ID) (*webserver.WsCategory, error) {
	if shopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Mising shop_id")
	}
	queryCategory := &catalog.GetShopCategoryQuery{
		ID:     ID,
		ShopID: shopID,
	}
	if err := w.cataglogQuery.Dispatch(ctx, queryCategory); err != nil {
		return nil, err
	}
	result, err := w.wsCategoryStore(ctx).ShopID(shopID).ID(ID).GetWsCategory()
	switch cm.ErrorCode(err) {
	case cm.NoError:
		result.Category = queryCategory.Result
	case cm.NotFound:
		result = &webserver.WsCategory{
			ID:        ID,
			ShopID:    shopID,
			Appear:    true,
			CreatedAt: queryCategory.Result.CreatedAt,
			UpdatedAt: queryCategory.Result.UpdatedAt,
			Category:  queryCategory.Result,
		}
	default:
		return nil, err
	}
	return result, nil
}

func (w WebserverQueryService) ListWsCategoriesByIDs(ctx context.Context, shopID dot.ID, IDs []dot.ID) ([]*webserver.WsCategory, error) {
	if shopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Mising shop_id")
	}
	queryCategory := &catalog.ListShopCategoriesByIDsQuery{
		ShopID: shopID,
		Ids:    IDs,
	}
	if err := w.cataglogQuery.Dispatch(ctx, queryCategory); err != nil {
		return nil, err
	}
	return w.getWsCategories(ctx, queryCategory.Result.Categories)
}

func (w WebserverQueryService) ListWsCategories(ctx context.Context, args webserver.ListWsCategoriesArgs) (*webserver.ListWsCategoriesResponse, error) {
	if args.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Mising shop_id")
	}
	queryCategory := &catalog.ListShopCategoriesQuery{
		ShopID:  args.ShopID,
		Paging:  args.Paging,
		Filters: args.Filters,
	}
	if err := w.cataglogQuery.Dispatch(ctx, queryCategory); err != nil {
		return nil, err
	}
	wsCategories, err := w.getWsCategories(ctx, queryCategory.Result.Categories)
	if err != nil {
		return nil, err
	}
	return &webserver.ListWsCategoriesResponse{
		PageInfo:     queryCategory.Result.Paging,
		WsCategories: wsCategories,
	}, nil
}

func (w WebserverQueryService) getWsCategories(ctx context.Context, shopCategories []*catalog.ShopCategory) ([]*webserver.WsCategory, error) {
	if len(shopCategories) < 1 {
		return []*webserver.WsCategory{}, nil
	}
	var productIDs []dot.ID
	for _, v := range shopCategories {
		productIDs = append(productIDs, v.ID)
	}
	wsCategories, err := w.wsCategoryStore(ctx).ShopID(shopCategories[0].ShopID).IDs(productIDs...).ListWsCategories()
	if err != nil {
		return nil, err
	}
	mapShopCategories := make(map[dot.ID]*webserver.WsCategory)
	for _, wsCategory := range wsCategories {
		mapShopCategories[wsCategory.ID] = wsCategory
	}
	var wsCategoriesResult []*webserver.WsCategory
	for _, shopCategory := range shopCategories {
		if mapShopCategories[shopCategory.ID] != nil {
			mapShopCategories[shopCategory.ID].Category = shopCategory
			wsCategoriesResult = append(wsCategoriesResult, mapShopCategories[shopCategory.ID])
			continue
		}
		wsCategoriesResult = append(wsCategoriesResult, &webserver.WsCategory{
			ID:       shopCategory.ID,
			ShopID:   shopCategory.ShopID,
			Category: shopCategory,
		})
	}
	return wsCategoriesResult, nil
}
