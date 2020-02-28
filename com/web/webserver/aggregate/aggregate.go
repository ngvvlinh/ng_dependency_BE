package aggregate

import (
	"context"

	"etop.vn/api/main/catalog"
	"etop.vn/api/webserver"
	"etop.vn/backend/com/web/webserver/convert"
	"etop.vn/backend/com/web/webserver/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/conversion"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/capi"
	"etop.vn/capi/dot"
)

var _ webserver.Aggregate = &WebserverAggregate{}
var scheme = conversion.Build(convert.RegisterConversions)

type WebserverAggregate struct {
	db              *cmsql.Database
	wsCategoryStore sqlstore.WsCategoryStoreFactory
	wsProductStore  sqlstore.WsProductStoreFactory
	wsPageStore     sqlstore.WsPageStoreFactory
	wsWebsiteStore  sqlstore.WsWebsiteStoreFactory
	eventBus        capi.EventBus
	bus             bus.Bus
	categoryQuery   catalog.QueryBus
}

func New(eventBus capi.EventBus, db *cmsql.Database, categoryQ catalog.QueryBus) *WebserverAggregate {
	return &WebserverAggregate{
		db:              db,
		wsCategoryStore: sqlstore.NewWsCategoryStore(db),
		wsProductStore:  sqlstore.NewWsProductStore(db),
		wsPageStore:     sqlstore.NewWsPageStore(db),
		wsWebsiteStore:  sqlstore.NewWsWebsiteStore(db),
		eventBus:        eventBus,
		categoryQuery:   categoryQ,
	}
}
func (q *WebserverAggregate) MessageBus() webserver.CommandBus {
	b := bus.New()
	return webserver.NewAggregateHandler(q).RegisterHandlers(b)
}

func (w WebserverAggregate) CreateOrUpdateWsCategory(ctx context.Context, args *webserver.CreateOrUpdateWsCategoryArgs) (*webserver.WsCategory, error) {
	if args.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Mising shop_id")
	}
	query := &catalog.GetShopCategoryQuery{
		ID:     args.ID,
		ShopID: args.ShopID,
	}
	err := w.categoryQuery.Dispatch(ctx, query)
	if err != nil {
		return nil, err
	}
	wsCategory, err := w.wsCategoryStore(ctx).ID(args.ID).GetWsCategory()
	switch cm.ErrorCode(err) {
	case cm.NoError:
		err = scheme.Convert(args, wsCategory)
		if err != nil {
			return nil, err
		}
		err = w.wsCategoryStore(ctx).UpdateWsCategoryAll(wsCategory)
		if err != nil {
			return nil, err
		}
	case cm.NotFound:
		wsCategoryCreate := &webserver.WsCategory{}
		err = scheme.Convert(args, wsCategoryCreate)
		if err != nil {
			return nil, err
		}
		err = w.wsCategoryStore(ctx).Create(wsCategoryCreate)
		if err != nil {
			return nil, err
		}
		wsCategory = wsCategoryCreate
	default:
		return nil, err
	}
	wsCategory.Category = query.Result
	return wsCategory, nil
}

func (w WebserverAggregate) CreateOrUpdateWsProduct(ctx context.Context, args *webserver.CreateOrUpdateWsProductArgs) (*webserver.WsProduct, error) {
	if args.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Mising shop_id")
	}
	query := &catalog.GetShopProductWithVariantsByIDQuery{
		ProductID: args.ID,
		ShopID:    args.ShopID,
	}
	err := w.categoryQuery.Dispatch(ctx, query)
	if err != nil {
		return nil, err
	}
	var mapVariant = make(map[dot.ID]*catalog.ShopVariant)
	for _, variant := range query.Result.Variants {
		mapVariant[variant.VariantID] = variant
	}
	for _, comparePrice := range args.ComparePrice {
		if mapVariant[comparePrice.VariantID] == nil {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "variant_id in compare price doesnt exist in %v product.", query.Result.Code)
		}
	}
	wsProduct, err := w.wsProductStore(ctx).ShopID(args.ShopID).ID(args.ID).GetWsProduct()
	switch cm.ErrorCode(err) {
	case cm.NoError:
		err = scheme.Convert(args, wsProduct)
		if err != nil {
			return nil, err
		}
		err = w.wsProductStore(ctx).ShopID(args.ShopID).ID(args.ID).UpdateWsProductAll(wsProduct)
		if err != nil {
			return nil, err
		}
	case cm.NotFound:
		wsProductCreate := &webserver.WsProduct{}
		err = scheme.Convert(args, wsProductCreate)
		if err != nil {
			return nil, err
		}
		err = w.wsProductStore(ctx).Create(wsProductCreate)
		if err != nil {
			return nil, err
		}
		wsProduct = wsProductCreate
	default:
		return nil, err
	}
	wsProduct.Product = query.Result
	return wsProduct, nil
}

func (w WebserverAggregate) CreateWsPage(ctx context.Context, args *webserver.CreateWsPageArgs) (*webserver.WsPage, error) {
	if args.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Mising shop_id")
	}
	var wsPage = &webserver.WsPage{}
	err := scheme.Convert(args, wsPage)
	if err != nil {
		return nil, err
	}
	err = w.wsPageStore(ctx).Create(wsPage)
	if err != nil {
		return nil, err
	}
	return wsPage, nil
}

func (w WebserverAggregate) UpdateWsPage(ctx context.Context, args *webserver.UpdateWsPageArgs) (*webserver.WsPage, error) {
	if args.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Mising shop_id")
	}
	wsPage, err := w.wsPageStore(ctx).ShopID(args.ShopID).ID(args.ID).GetWsPage()
	if err != nil {
		return nil, err
	}
	err = scheme.Convert(args, wsPage)
	if err != nil {
		return nil, err
	}
	err = w.wsPageStore(ctx).ShopID(args.ShopID).ID(args.ID).UpdateWsPageAll(wsPage)
	if err != nil {
		return nil, err
	}
	return wsPage, nil
}

func (w WebserverAggregate) DeleteWsPage(ctx context.Context, shopID dot.ID, ID dot.ID) (int, error) {
	if shopID == 0 {
		return 0, cm.Errorf(cm.InvalidArgument, nil, "Mising shop_id")
	}
	_, err := w.wsPageStore(ctx).ShopID(shopID).ID(ID).GetWsPage()
	if err != nil {
		return 0, err
	}
	deleted, err := w.wsPageStore(ctx).ShopID(shopID).ID(ID).SoftDelete()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

func (w WebserverAggregate) CreateOrUpdateWsWebsite(ctx context.Context, args *webserver.CreateWsPageArgs) (*webserver.WsWebsite, error) {
	if args.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Mising shop_id")
	}
	var wsWebsite = &webserver.WsWebsite{}
	wsWebsite, err := w.wsWebsiteStore(ctx).ShopID(args.ShopID).GetWsWebsite()
	switch cm.ErrorCode(err) {
	case cm.NoError:
		err = scheme.Convert(args, wsWebsite)
		if err != nil {
			return nil, err
		}
		err = w.wsWebsiteStore(ctx).UpdateWsWebsiteAll(wsWebsite)
		if err != nil {
			return nil, err
		}
	case cm.NotFound:
		err := scheme.Convert(args, wsWebsite)
		if err != nil {
			return nil, err
		}
		err = w.wsWebsiteStore(ctx).Create(wsWebsite)
		if err != nil {
			return nil, err
		}
	default:
		return nil, err
	}
	return wsWebsite, nil
}

func (w WebserverAggregate) CreateWsWebsite(ctx context.Context, args *webserver.CreateWsWebsiteArgs) (*webserver.WsWebsite, error) {
	if args.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Mising shop_id")
	}
	var wsWebsite = &webserver.WsWebsite{}
	err := scheme.Convert(args, wsWebsite)
	if err != nil {
		return nil, err
	}
	var productIDs = []dot.ID{}
	if wsWebsite.OutstandingProduct != nil {
		productIDs = makeUnduplicatedIDList(productIDs, wsWebsite.OutstandingProduct.ProductIDs)
	}
	if wsWebsite.NewProduct != nil {
		productIDs = makeUnduplicatedIDList(productIDs, wsWebsite.NewProduct.ProductIDs)
	}
	err = w.checkListproduct(ctx, args.ShopID, productIDs)
	if err != nil {
		return nil, err
	}
	err = w.wsWebsiteStore(ctx).Create(wsWebsite)
	if err != nil {
		return nil, err
	}
	return wsWebsite, nil
}

func (w WebserverAggregate) UpdateWsWebsite(ctx context.Context, args *webserver.UpdateWsWebsiteArgs) (*webserver.WsWebsite, error) {
	if args.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Mising shop_id")
	}
	wsWebsite, err := w.wsWebsiteStore(ctx).ShopID(args.ShopID).ID(args.ID).GetWsWebsite()
	if err != nil {
		return nil, err
	}
	err = scheme.Convert(args, wsWebsite)
	if err != nil {
		return nil, err
	}
	var productIDs = []dot.ID{}
	if wsWebsite.OutstandingProduct != nil {
		productIDs = makeUnduplicatedIDList(productIDs, wsWebsite.OutstandingProduct.ProductIDs)
	}
	if wsWebsite.NewProduct != nil {
		productIDs = makeUnduplicatedIDList(productIDs, wsWebsite.NewProduct.ProductIDs)
	}
	err = w.checkListproduct(ctx, args.ShopID, productIDs)
	if err != nil {
		return nil, err
	}
	err = w.wsWebsiteStore(ctx).Create(wsWebsite)
	if err != nil {
		return nil, err
	}
	return wsWebsite, nil
}

func (w WebserverAggregate) checkListproduct(ctx context.Context, shopID dot.ID, ids []dot.ID) error {
	query := &catalog.ListShopProductsByIDsQuery{
		IDs:    ids,
		ShopID: shopID,
	}
	err := w.categoryQuery.Dispatch(ctx, query)
	if err != nil {
		return err
	}
	var productIDs []dot.ID
	for _, product := range query.Result.Products {
		productIDs = append(productIDs, product.ProductID)
	}

	for _, productIDArg := range ids {
		check := true
		for _, productID := range productIDs {
			if productID == productIDArg {
				check = false
				break
			}
		}
		if check {
			return cm.Errorf(cm.InvalidArgument, nil, "Sản phẩm có id = %v không tồn tại", productIDArg)
		}
	}
	return nil
}

func makeUnduplicatedIDList(in []dot.ID, out []dot.ID) []dot.ID {
	for _, value := range out {
		isExisted := false
		for _, productID := range in {
			if value == productID {
				isExisted = true
				break
			}
		}
		if !isExisted {
			in = append(in, value)
		}
	}
	return in
}
