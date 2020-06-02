package query

import (
	"context"

	"o.o/api/main/catalog"
	"o.o/api/webserver"
	"o.o/backend/com/web"
	"o.o/backend/com/web/webserver/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi"
	"o.o/capi/dot"
)

var _ webserver.QueryService = &WebserverQueryService{}

type WebserverQueryService struct {
	db              *cmsql.Database
	wsCategoryStore sqlstore.WsCategoryStoreFactory
	wsProductStore  sqlstore.WsProductStoreFactory
	wsPageStore     sqlstore.WsPageStoreFactory
	wsWebsiteStore  sqlstore.WsWebsiteStoreFactory
	cataglogQuery   catalog.QueryBus
	eventBus        capi.EventBus
	bus             bus.Bus
}

func New(eventBus capi.EventBus, db web.WebServerDB, cataglogQ catalog.QueryBus) *WebserverQueryService {
	return &WebserverQueryService{
		db:              db,
		wsCategoryStore: sqlstore.NewWsCategoryStore(db),
		wsProductStore:  sqlstore.NewWsProductStore(db),
		wsPageStore:     sqlstore.NewWsPageStore(db),
		wsWebsiteStore:  sqlstore.NewWsWebsiteStore(db),
		cataglogQuery:   cataglogQ,
		eventBus:        eventBus,
	}
}

func WebserverQueryServiceMessageBus(q *WebserverQueryService) webserver.QueryBus {
	b := bus.New()
	return webserver.NewQueryServiceHandler(q).RegisterHandlers(b)
}

// TODO decide after release
// func (w WebserverQueryService) GetWsCategoryByID(ctx context.Context, shopID dot.ID, ID dot.ID) (*webserver.WsCategory, error) {
// 	if shopID == 0 {
// 		return nil, cm.Errorf(cm.InvalidArgument, nil, "Mising shop_id")
// 	}
// 	return w.wsCategoryStore(ctx).ShopID(shopID).ID(ID).GetWsCategory()
// }
//
// func (w WebserverQueryService) ListWsCategoriesByIDs(ctx context.Context, shopID dot.ID, IDs []dot.ID) ([]*webserver.WsCategory, error) {
// 	if shopID == 0 {
// 		return nil, cm.Errorf(cm.InvalidArgument, nil, "Mising shop_id")
// 	}
// 	return w.wsCategoryStore(ctx).ShopID(shopID).IDs(IDs...).ListWsCategories()
// }
//
// func (w WebserverQueryService) ListWsCategories(ctx context.Context, args webserver.ListWsCategoriesArgs) (*webserver.ListWsCategoriesResponse, error) {
// 	if args.ShopID == 0 {
// 		return nil, cm.Errorf(cm.InvalidArgument, nil, "Mising shop_id")
// 	}
// 	q := w.wsCategoryStore(ctx).ShopID(args.ShopID).Filters(args.Filters).WithPaging(args.Paging)
// 	result, err := q.ListWsCategories()
// 	if err != nil {
// 		return nil, err
// 	}
// 	paging := q.GetPaging()
// 	return &webserver.ListWsCategoriesResponse{
// 		PageInfo:     paging,
// 		WsCategories: result,
// 	}, nil
// }

func (w WebserverQueryService) GetWsPageByID(ctx context.Context, shopID dot.ID, ID dot.ID) (*webserver.WsPage, error) {
	if shopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Mising shop_id")
	}
	return w.wsPageStore(ctx).ShopID(shopID).ID(ID).GetWsPage()
}

func (w WebserverQueryService) ListWsPagesByIDs(ctx context.Context, shopID dot.ID, IDs []dot.ID) ([]*webserver.WsPage, error) {
	if shopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Mising shop_id")
	}
	return w.wsPageStore(ctx).ShopID(shopID).IDs(IDs...).ListWsPages()
}

func (w WebserverQueryService) ListWsPages(ctx context.Context, args webserver.ListWsPagesArgs) (*webserver.ListWsPagesResponse, error) {
	if args.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Mising shop_id")
	}
	q := w.wsPageStore(ctx).ShopID(args.ShopID).Filters(args.Filters).WithPaging(args.Paging)
	result, err := q.ListWsPages()
	if err != nil {
		return nil, err
	}
	paging := q.GetPaging()
	return &webserver.ListWsPagesResponse{
		PageInfo: paging,
		WsPages:  result,
	}, nil
}

func (w WebserverQueryService) GetWsWebsiteByID(ctx context.Context, shopID dot.ID, ID dot.ID) (*webserver.WsWebsite, error) {
	if shopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Mising shop_id")
	}
	result, err := w.wsWebsiteStore(ctx).ShopID(shopID).ID(ID).GetWsWebsite()
	if err != nil {
		return nil, err
	}
	return w.addProductInfo(ctx, result)
}

func (w WebserverQueryService) ListWsWebsitesByIDs(ctx context.Context, shopID dot.ID, IDs []dot.ID) ([]*webserver.WsWebsite, error) {
	if shopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Mising shop_id")
	}
	result, err := w.wsWebsiteStore(ctx).ShopID(shopID).IDs(IDs...).ListWsWebsites()
	if err != nil {
		return nil, err
	}
	return w.addProductsInfo(ctx, result)
}

func (w WebserverQueryService) ListWsWebsites(ctx context.Context, args webserver.ListWsWebsitesArgs) (*webserver.ListWsWebsitesResponse, error) {
	if args.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Mising shop_id")
	}
	q := w.wsWebsiteStore(ctx).ShopID(args.ShopID).Filters(args.Filters).WithPaging(args.Paging)
	result, err := q.ListWsWebsites()
	if err != nil {
		return nil, err
	}
	result, err = w.addProductsInfo(ctx, result)
	if err != nil {
		return nil, err
	}
	paging := q.GetPaging()
	return &webserver.ListWsWebsitesResponse{
		WsWebsites: result,
		PageInfo:   paging,
	}, nil
}

func (w WebserverQueryService) addProductInfo(ctx context.Context, args *webserver.WsWebsite) (*webserver.WsWebsite, error) {
	if args != nil {
		var productIDs []dot.ID
		if args.OutstandingProduct != nil {
			productIDs = append(productIDs, args.OutstandingProduct.ProductIDs...)
		}
		if args.NewProduct != nil {
			productIDs = append(productIDs, args.NewProduct.ProductIDs...)
		}
		listWsProduct, err := w.ListWsProductsByIDs(ctx, args.ShopID, productIDs)
		if err != nil {
			return nil, err
		}
		var mapProducts = make(map[dot.ID]*webserver.WsProduct)
		for _, v := range listWsProduct {
			mapProducts[v.ID] = v
		}
		if args.OutstandingProduct != nil {
			for _, productID := range args.OutstandingProduct.ProductIDs {
				if mapProducts[productID] != nil && mapProducts[productID].Appear {
					args.OutstandingProduct.Products = append(args.OutstandingProduct.Products, mapProducts[productID])
				}
			}
		}
		if args.NewProduct != nil {
			for _, productID := range args.NewProduct.ProductIDs {
				if mapProducts[productID] != nil && mapProducts[productID].Appear {
					args.NewProduct.Products = append(args.NewProduct.Products, mapProducts[productID])
				}
			}
		}
	}
	return args, nil
}

func (w WebserverQueryService) addProductsInfo(ctx context.Context, args []*webserver.WsWebsite) ([]*webserver.WsWebsite, error) {
	if len(args) > 0 {
		var productIDs []dot.ID
		var mapProducts = make(map[dot.ID]*webserver.WsProduct)
		for _, v := range args {
			if v.OutstandingProduct != nil {
				productIDs = append(productIDs, v.OutstandingProduct.ProductIDs...)
			}
			if v.NewProduct != nil {
				productIDs = append(productIDs, v.NewProduct.ProductIDs...)
			}
		}
		if len(productIDs) > 0 {
			listWsProduct, err := w.ListWsProductsByIDs(ctx, args[0].ShopID, productIDs)
			if err != nil {
				return nil, err
			}
			for _, v := range listWsProduct {
				mapProducts[v.ID] = v
			}
		}
		for _, wsWebsite := range args {
			//update info outstanding product
			var productOutStanding []dot.ID
			if wsWebsite.OutstandingProduct != nil {
				for _, productID := range wsWebsite.OutstandingProduct.ProductIDs {
					if mapProducts[productID] != nil && mapProducts[productID].Appear {
						productOutStanding = append(productOutStanding, mapProducts[productID].ID)
						wsWebsite.OutstandingProduct.Products = append(wsWebsite.OutstandingProduct.Products, mapProducts[productID])
					}
				}
			}
			wsWebsite.OutstandingProduct.ProductIDs = productOutStanding

			//update info new product
			var productNew []dot.ID
			if wsWebsite.NewProduct != nil {
				for _, productID := range wsWebsite.NewProduct.ProductIDs {
					if mapProducts[productID] != nil && mapProducts[productID].Appear {
						productNew = append(productNew, mapProducts[productID].ID)
						wsWebsite.NewProduct.Products = append(wsWebsite.NewProduct.Products, mapProducts[productID])
					}
				}
			}
			wsWebsite.NewProduct.ProductIDs = productNew
		}
	}
	return args, nil
}

func (w WebserverQueryService) GetShopIDBySiteSubdomain(ctx context.Context, siteSubDoimain string) (dot.ID, error) {
	if siteSubDoimain == "" {
		return 0, cm.Errorf(cm.InvalidArgument, nil, "Site SubDomain không thể rổng")
	}
	wsWesite, err := w.wsWebsiteStore(ctx).SiteSubdomain(siteSubDoimain).GetWsWebsite()
	if err != nil {
		return 0, err
	}
	return wsWesite.ShopID, nil
}
