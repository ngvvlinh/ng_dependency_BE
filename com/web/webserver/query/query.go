package query

import (
	"context"

	"o.o/api/main/catalog"
	"o.o/api/webserver"
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

func New(eventBus capi.EventBus, db *cmsql.Database, cataglogQ catalog.QueryBus) *WebserverQueryService {
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

func (q *WebserverQueryService) MessageBus() webserver.QueryBus {
	b := bus.New()
	return webserver.NewQueryServiceHandler(q).RegisterHandlers(b)
}

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
	return w.wsWebsiteStore(ctx).ShopID(shopID).ID(ID).GetWsWebsite()
}

func (w WebserverQueryService) ListWsWebsitesByIDs(ctx context.Context, shopID dot.ID, IDs []dot.ID) ([]*webserver.WsWebsite, error) {
	if shopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Mising shop_id")
	}
	return w.wsWebsiteStore(ctx).ShopID(shopID).IDs(IDs...).ListWsWebsites()
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
	paging := q.GetPaging()
	return &webserver.ListWsWebsitesResponse{
		WsWebsites: result,
		PageInfo:   paging,
	}, nil
}
