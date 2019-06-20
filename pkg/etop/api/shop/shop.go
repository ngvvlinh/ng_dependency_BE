package shop

import (
	"context"

	"etop.vn/backend/pkg/etop/api"

	"etop.vn/api/main/address"

	"etop.vn/backend/pkg/etop/logic/shipping_provider"

	"etop.vn/api/main/identity"

	"etop.vn/api/main/shipping/v1/types"

	"etop.vn/api/main/shipnow/carrier"

	"etop.vn/api/main/catalog"

	"etop.vn/api/main/shipnow"
	pbcm "etop.vn/backend/pb/common"
	pbetop "etop.vn/backend/pb/etop"
	pborder "etop.vn/backend/pb/etop/order"
	pbshop "etop.vn/backend/pb/etop/shop"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/idemp"
	"etop.vn/backend/pkg/common/l"
	"etop.vn/backend/pkg/common/redis"
	cmservice "etop.vn/backend/pkg/common/service"
	"etop.vn/backend/pkg/etop/api/convertpb"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/etop/sqlstore"
	notimodel "etop.vn/backend/pkg/notifier/model"
	catalogmodel "etop.vn/backend/pkg/services/catalog/model"
	catalogmodelx "etop.vn/backend/pkg/services/catalog/modelx"
	moneymodelx "etop.vn/backend/pkg/services/moneytx/modelx"
	wrapshop "etop.vn/backend/wrapper/etop/shop"
)

var ll = l.New()

func init() {
	bus.AddHandler("api", CreateCollection)
	bus.AddHandler("api", DeleteCollection)
	bus.AddHandler("api", GetCollection)
	bus.AddHandler("api", GetCollections)
	bus.AddHandler("api", GetCollectionsByIDs)
	bus.AddHandler("api", RemoveVariants)
	bus.AddHandler("api", RemoveProductsCollection)
	bus.AddHandler("api", UpdateCollection)
	bus.AddHandler("api", UpdateVariant)
	bus.AddHandler("api", UpdateProducts)
	bus.AddHandler("api", UpdateProductsCollection)
	bus.AddHandler("api", VersionInfo)

	bus.AddHandler("api", AddProducts)
	bus.AddHandler("api", GetProduct)
	bus.AddHandler("api", GetProducts)
	bus.AddHandler("api", GetProductsByIDs)
	bus.AddHandler("api", UpdateProduct)
	bus.AddHandler("api", UpdateProductsTags)
	bus.AddHandler("api", RemoveProducts)

	bus.AddHandler("api", CreateProductSource)
	bus.AddHandler("api", CreateVariant)
	bus.AddHandler("api", GetShopProductSources)
	bus.AddHandler("api", CreateProductSourceCategory)
	bus.AddHandler("api", UpdateProductsPSCategory)
	bus.AddHandler("api", UpdateProductImages)
	bus.AddHandler("api", GetProductSourceCategory)
	bus.AddHandler("api", GetProductSourceCategories)
	bus.AddHandler("api", UpdateProductSourceCategory)
	bus.AddHandler("api", RemoveProductSourceCategory)
	bus.AddHandler("api", UpdateVariantImages)

	bus.AddHandler("api", GetMoneyTransaction)
	bus.AddHandler("api", GetMoneyTransactions)

	bus.AddHandler("api", SummarizeFulfillments)
	bus.AddHandler("api", CalcBalance)
	bus.AddHandler("api", CreateDevice)
	bus.AddHandler("api", DeleteDevice)
	bus.AddHandler("api", GetNotifications)
	bus.AddHandler("api", GetNotification)
	bus.AddHandler("api", UpdateNotifications)
}

const PrefixIdemp = "IdempOrder"

var idempgroup *idemp.RedisGroup
var shipnowAggr shipnow.CommandBus
var shipnowQuery shipnow.QueryBus
var identityAggr identity.CommandBus
var identityQuery identity.QueryBus
var addressQuery address.QueryBus
var shippingCtrl *shipping_provider.ProviderManager
var catalogQuery catalog.QueryBus

func Init(catalogQueryBus catalog.QueryBus, shipnow shipnow.CommandBus, shipnowQS shipnow.QueryBus, identity identity.CommandBus, identityQS identity.QueryBus, addressQS address.QueryBus, shippingCtrl *shipping_provider.ProviderManager, sd cmservice.Shutdowner, rd redis.Store) {
	idempgroup = idemp.NewRedisGroup(rd, PrefixIdemp, 5*60)
	catalogQuery = catalogQueryBus
	shippingCtrl = shippingCtrl
	shipnowAggr = shipnow
	shipnowQuery = shipnowQS
	identityAggr = identity
	identityQuery = identityQS
	addressQuery = addressQS
	sd.Register(idempgroup.Shutdown)
}

func VersionInfo(ctx context.Context, q *wrapshop.VersionInfoEndpoint) error {
	q.Result = &pbcm.VersionInfoResponse{
		Service: "etop.Shop",
		Version: "0.1",
	}
	return nil
}

func GetCollection(ctx context.Context, q *wrapshop.GetCollectionEndpoint) error {
	query := &catalogmodelx.GetShopCollectionQuery{
		ShopID:       q.Context.Shop.ID,
		CollectionID: q.Id,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = pbshop.PbCollection(query.Result)
	return nil
}

func GetCollectionsByIDs(ctx context.Context, q *wrapshop.GetCollectionsByIDsEndpoint) error {
	query := &catalogmodelx.GetShopCollectionsQuery{
		ShopID:        q.Context.Shop.ID,
		CollectionIDs: q.Ids,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pbshop.CollectionsResponse{
		Collections: pbshop.PbCollections(query.Result.Collections),
	}
	return nil
}

func GetCollections(ctx context.Context, q *wrapshop.GetCollectionsEndpoint) error {
	query := &catalogmodelx.GetShopCollectionsQuery{
		ShopID: q.Context.Shop.ID,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pbshop.CollectionsResponse{
		Collections: pbshop.PbCollections(query.Result.Collections),
	}
	return nil
}

func UpdateVariant(ctx context.Context, q *wrapshop.UpdateVariantEndpoint) error {
	shopID := q.Context.Shop.ID
	productSourceID := q.Context.Shop.ProductSourceID
	cmd := &catalogmodelx.UpdateShopVariantCommand{
		ShopID:          shopID,
		Variant:         pbshop.PbUpdateVariantToModel(shopID, q.UpdateVariantRequest),
		CostPrice:       q.CostPrice,
		Code:            q.Sku,
		Attributes:      convertpb.AttributesTomodel(q.Attributes),
		ProductSourceID: productSourceID,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = PbShopVariant(cmd.Result)
	return nil
}

func UpdateProducts(ctx context.Context, q *wrapshop.UpdateVariantsEndpoint) error {
	return cm.ErrTODO
}

func CreateCollection(ctx context.Context, q *wrapshop.CreateCollectionEndpoint) error {
	cmd := &catalogmodelx.CreateShopCollectionCommand{
		Collection: pbshop.PbCreateCollection(q.Context.Shop.ID, q.CreateCollectionRequest),
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = pbshop.PbCollection(cmd.Result)
	return nil
}

func DeleteCollection(ctx context.Context, q *wrapshop.DeleteCollectionEndpoint) error {
	cmd := &catalogmodelx.RemoveShopCollectionCommand{
		ShopID:       q.Context.Shop.ID,
		CollectionID: q.Id,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.RemovedResponse{
		Removed: int32(cmd.Result.Deleted),
	}
	return nil
}

func RemoveVariants(ctx context.Context, q *wrapshop.RemoveVariantsEndpoint) error {
	cmd := &catalogmodelx.RemoveShopVariantsCommand{
		ShopID:          q.Context.Shop.ID,
		IDs:             q.Ids,
		ProductSourceID: q.Context.Shop.ProductSourceID,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.RemovedResponse{
		Removed: int32(cmd.Result.Removed),
	}
	return nil
}

func UpdateCollection(ctx context.Context, q *wrapshop.UpdateCollectionEndpoint) error {
	cmd := &catalogmodelx.UpdateShopCollectionCommand{
		Collection: pbshop.PbUpdateCollection(q.Context.Shop.ID, q.UpdateCollectionRequest),
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = pbshop.PbCollection(cmd.Result)
	return nil
}

func UpdateProductsCollection(ctx context.Context, q *wrapshop.UpdateProductsCollectionEndpoint) error {
	cmd := &catalogmodelx.AddProductsToShopCollectionCommand{
		ShopID:       q.Context.Shop.ID,
		ProductIDs:   q.ProductIds,
		CollectionID: q.CollectionId,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbshop.UpdateProductsCollectionResponse{
		Updated: int32(cmd.Result.Updated),
		Errors:  pbcm.PbErrors(cmd.Result.Errors),
	}
	return nil
}

func RemoveProductsCollection(ctx context.Context, q *wrapshop.RemoveProductsCollectionEndpoint) error {
	cmd := &catalogmodelx.RemoveProductsFromShopCollectionCommand{
		ShopID:       q.Context.Shop.ID,
		ProductIDs:   q.ProductIds,
		CollectionID: q.CollectionId,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = pbcm.Updated(cmd.Result.Updated)
	return nil
}

func AddProducts(ctx context.Context, q *wrapshop.AddProductsEndpoint) error {
	cmd := &catalogmodelx.AddShopProductsCommand{
		ShopID: q.Context.Shop.ID,
		IDs:    q.Ids,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbshop.AddProductsResponse{
		Products: PbShopProducts(cmd.Result.Products),
		Errors:   pbcm.PbErrors(cmd.Result.Errors),
	}
	return nil
}

func GetProduct(ctx context.Context, q *wrapshop.GetProductEndpoint) error {
	query := &catalog.GetShopProductWithVariantsByIDQuery{
		ProductID: q.Id,
		ShopID:    q.Context.Shop.ID,
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = PbShopProductWithVariants(query.Result)
	return nil
}

func GetProductsByIDs(ctx context.Context, q *wrapshop.GetProductsByIDsEndpoint) error {
	query := &catalog.ListShopProductsWithVariantsByIDsQuery{
		IDs:    q.Ids,
		ShopID: q.Context.Shop.ID,
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pbshop.ShopProductsResponse{
		Products: PbShopProductsWithVariants(query.Result.Products),
	}
	return nil
}

func GetProducts(ctx context.Context, q *wrapshop.GetProductsEndpoint) error {
	paging := q.Paging.CMPaging()
	query := &catalog.ListShopProductsWithVariantsQuery{
		ShopID:  q.Context.Shop.ID,
		Paging:  *paging,
		Filters: pbcm.ToFilters(q.Filters),
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = &pbshop.ShopProductsResponse{
		Paging:   pbcm.PbPageInfo(paging, query.Result.Count),
		Products: PbShopProductsWithVariants(query.Result.Products),
	}
	return nil
}

func RemoveProducts(ctx context.Context, q *wrapshop.RemoveProductsEndpoint) error {
	productSourceID := q.Context.Shop.ProductSourceID
	cmd := &catalogmodelx.RemoveShopProductsCommand{
		ShopID:          q.Context.Shop.ID,
		IDs:             q.Ids,
		ProductSourceID: productSourceID,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.RemovedResponse{
		Removed: int32(cmd.Result.Removed),
	}
	return nil
}

func UpdateProduct(ctx context.Context, q *wrapshop.UpdateProductEndpoint) error {
	shopID := q.Context.Shop.ID
	productSourceID := q.Context.Shop.ProductSourceID
	cmd := &catalogmodelx.UpdateShopProductCommand{
		ShopID:          shopID,
		Product:         pbshop.PbUpdateProductToModel(shopID, q.UpdateProductRequest),
		Code:            q.Code,
		ProductSourceID: productSourceID,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = PbShopProductWithVariants(cmd.Result)
	return nil
}

func UpdateProductsTags(ctx context.Context, q *wrapshop.UpdateProductsTagsEndpoint) error {
	shopID := q.Context.Shop.ID
	productSourceID := q.Context.Shop.ProductSourceID
	cmd := &catalogmodelx.UpdateShopProductsTagsCommand{
		ShopID:     shopID,
		ProductIDs: q.Ids,
		Update: &model.UpdateListRequest{
			Adds:       q.Adds,
			Deletes:    q.Deletes,
			ReplaceAll: q.ReplaceAll,
			DeleteAll:  q.DeleteAll,
		},
		ProductSourceID: productSourceID,
	}

	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.UpdatedResponse{
		Updated: int32(cmd.Result.Updated),
	}
	return nil
}

func CreateProductSource(ctx context.Context, q *wrapshop.CreateProductSourceEndpoint) error {
	shopID := q.Context.Shop.ID
	cmd := &catalogmodelx.CreateProductSourceCommand{
		ShopID: shopID,
		Name:   q.Name,
		Type:   q.Type.ToModel(),
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = PbProductSource(cmd.Result)
	return nil
}

func CreateVariant(ctx context.Context, q *wrapshop.CreateVariantEndpoint) error {
	cmd := &catalogmodelx.CreateVariantCommand{
		ShopID:            q.Context.Shop.ID,
		ProductSourceID:   q.ProductSourceId,
		ProductID:         q.ProductId,
		ProductName:       q.ProductName,
		Name:              q.Name,
		Description:       q.Description,
		ShortDesc:         q.ShortDesc,
		ImageURLs:         q.ImageUrls,
		Tags:              q.Tags,
		Status:            *q.Status.ToModel(),
		ListPrice:         q.ListPrice,
		SKU:               q.Sku,
		Code:              q.Code,
		QuantityAvailable: int(q.QuantityAvailable),
		QuantityOnHand:    int(q.QuantityOnHand),
		QuantityReserved:  int(q.QuantityReserved),
		CostPrice:         q.CostPrice,
		Attributes:        convertpb.AttributesTomodel(q.Attributes),
		DescHTML:          q.DescHtml,
	}

	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}

	q.Result = PbShopProductWithVariants(cmd.Result)
	return nil
}

func GetShopProductSources(ctx context.Context, q *wrapshop.GetShopProductSourcesEndpoint) error {
	query := &catalogmodelx.GetShopProductSourcesCommand{}
	if q.Context.User != nil {
		query.UserID = q.Context.User.ID
	} else {
		query.ShopID = q.Context.Shop.ID
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = &pbshop.ProductSourcesResponse{
		ProductSources: PbProductSources(query.Result),
	}
	return nil
}

func CreateProductSourceCategory(ctx context.Context, q *wrapshop.CreateProductSourceCategoryEndpoint) error {
	cmd := &catalogmodelx.CreateProductSourceCategoryCommand{
		ShopID:            q.Context.Shop.ID,
		Name:              q.Name,
		ParentID:          q.ParentId,
		ProductSourceID:   q.ProductSourceId,
		ProductSourceType: q.ProductSourceType.ToModel(),
	}

	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = convertpb.PbCategory(cmd.Result)
	return nil
}

func UpdateProductsPSCategory(ctx context.Context, q *wrapshop.UpdateProductsPSCategoryEndpoint) error {
	cmd := &catalogmodelx.UpdateProductsProductSourceCategoryCommand{
		CategoryID:      q.CategoryId,
		ProductIDs:      q.ProductIds,
		ShopID:          q.Context.Shop.ID,
		ProductSourceID: q.Context.Shop.ProductSourceID,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.UpdatedResponse{
		Updated: int32(cmd.Result.Updated),
	}
	return nil
}

func GetProductSourceCategory(ctx context.Context, q *wrapshop.GetProductSourceCategoryEndpoint) error {
	cmd := &catalogmodelx.GetProductSourceCategoryQuery{
		ShopID:     q.Context.Shop.ID,
		CategoryID: q.Id,
	}

	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}

	q.Result = convertpb.PbCategory(cmd.Result)
	return nil
}

func GetProductSourceCategories(ctx context.Context, q *wrapshop.GetProductSourceCategoriesEndpoint) error {
	cmd := &catalogmodelx.GetProductSourceCategoriesQuery{
		ProductSourceType: q.Type.ToModel(),
		ShopID:            q.Context.Shop.ID,
	}

	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}

	q.Result = &pbshop.CategoriesResponse{
		Categories: convertpb.PbCategories(cmd.Result.Categories),
	}
	return nil
}

func UpdateProductSourceCategory(ctx context.Context, q *wrapshop.UpdateProductSourceCategoryEndpoint) error {
	cmd := &catalogmodelx.UpdateShopProductSourceCategoryCommand{
		ID:       q.Id,
		ShopID:   q.Context.Shop.ID,
		ParentID: q.ParentId,
		Name:     q.Name,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = convertpb.PbCategory(cmd.Result)
	return nil
}

func RemoveProductSourceCategory(ctx context.Context, q *wrapshop.RemoveProductSourceCategoryEndpoint) error {
	cmd := &catalogmodelx.RemoveShopProductSourceCategoryCommand{
		ID:     q.Id,
		ShopID: q.Context.Shop.ID,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.RemovedResponse{
		Removed: int32(cmd.Result.Removed),
	}
	return nil
}

func UpdateProductImages(ctx context.Context, q *wrapshop.UpdateProductImagesEndpoint) error {
	shopID := q.Context.Shop.ID
	query := &catalog.GetShopProductByIDQuery{
		ProductID: q.Id,
		ShopID:    shopID,
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}

	r := &model.UpdateListRequest{
		Adds:       q.Adds,
		Deletes:    q.Deletes,
		ReplaceAll: q.ReplaceAll,
		DeleteAll:  q.DeleteAll,
	}

	imageURLs, err := pbcm.PatchImage(query.Result.ShopProduct.ImageURLs, r)
	if err != nil {
		return err
	}

	cmd := &catalogmodelx.UpdateShopProductCommand{
		ShopID: shopID,
		Product: &catalogmodel.ShopProduct{
			ProductID: q.Id,
			ShopID:    shopID,
			ImageURLs: imageURLs,
		},
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = PbShopProductWithVariants(cmd.Result)
	return nil
}

func UpdateVariantImages(ctx context.Context, q *wrapshop.UpdateVariantImagesEndpoint) error {
	shopID := q.Context.Shop.ID
	query := &catalogmodelx.GetShopVariantQuery{
		ShopID:    shopID,
		VariantID: q.Id,
	}
	var sourceImages []string
	if err := bus.Dispatch(ctx, query); err == nil {
		sourceImages = query.Result.ShopVariant.ImageURLs
	}

	r := &model.UpdateListRequest{
		Adds:       q.Adds,
		Deletes:    q.Deletes,
		ReplaceAll: q.ReplaceAll,
		DeleteAll:  q.DeleteAll,
	}
	imageURLs, err := pbcm.PatchImage(sourceImages, r)
	if err != nil {
		return err
	}

	cmd := &catalogmodelx.UpdateShopVariantCommand{
		ShopID: shopID,
		Variant: &catalogmodel.ShopVariant{
			VariantID: q.Id,
			ShopID:    shopID,
			ImageURLs: imageURLs,
		},
		ProductSourceID: q.Context.Shop.ProductSourceID,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = PbShopVariant(cmd.Result)
	return nil
}

func GetMoneyTransaction(ctx context.Context, q *wrapshop.GetMoneyTransactionEndpoint) error {
	query := &moneymodelx.GetMoneyTransaction{
		ShopID: q.Context.Shop.ID,
		ID:     q.Id,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = pborder.PbMoneyTransactionExtended(query.Result)
	return nil
}

func GetMoneyTransactions(ctx context.Context, q *wrapshop.GetMoneyTransactionsEndpoint) error {
	paging := q.Paging.CMPaging()
	query := &moneymodelx.GetMoneyTransactions{
		ShopID:  q.Context.Shop.ID,
		Paging:  paging,
		Filters: pbcm.ToFilters(q.Filters),
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pborder.MoneyTransactionsResponse{
		MoneyTransactions: pborder.PbMoneyTransactionExtendeds(query.Result.MoneyTransactions),
		Paging:            pbcm.PbPageInfo(paging, int32(query.Result.Total)),
	}
	return nil
}

func SummarizeFulfillments(ctx context.Context, q *wrapshop.SummarizeFulfillmentsEndpoint) error {
	query := &model.SummarizeFulfillmentsRequest{
		ShopID:   q.Context.Shop.ID,
		DateFrom: q.DateFrom,
		DateTo:   q.DateTo,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = &pbshop.SummarizeFulfillmentsResponse{
		Tables: pbshop.PbSummaryTables(query.Result.Tables),
	}
	return nil
}

func CalcBalance(ctx context.Context, q *wrapshop.CalcBalanceShopEndpoint) error {
	query := &model.GetBalanceShopCommand{
		ShopID: q.Context.Shop.ID,
	}

	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pbshop.CalcBalanceShopResponse{
		Balance: int32(query.Result.Amount),
	}
	return nil
}

func CreateDevice(ctx context.Context, q *wrapshop.CreateDeviceEndpoint) error {
	cmd := &notimodel.CreateDeviceArgs{
		UserID:           q.Context.UserID,
		AccountID:        q.Context.Shop.ID,
		DeviceID:         q.DeviceId,
		DeviceName:       q.DeviceName,
		ExternalDeviceID: q.ExternalDeviceId,
	}
	device, err := sqlstore.CreateDevice(ctx, cmd)
	if err != nil {
		return err
	}
	q.Result = pbetop.PbDevice(device)
	return nil
}

func DeleteDevice(ctx context.Context, q *wrapshop.DeleteDeviceEndpoint) error {
	device := &notimodel.Device{
		DeviceID:         q.DeviceId,
		ExternalDeviceID: q.ExternalDeviceId,
		AccountID:        q.Context.Shop.ID,
		UserID:           q.Context.UserID,
	}
	if err := sqlstore.DeleteDevice(ctx, device); err != nil {
		return err
	}
	q.Result = &pbcm.DeletedResponse{
		Deleted: 1,
	}
	return nil
}

func GetNotification(ctx context.Context, q *wrapshop.GetNotificationEndpoint) error {
	query := &notimodel.GetNotificationArgs{
		AccountID: q.Context.Shop.ID,
		ID:        q.Id,
	}
	noti, err := sqlstore.GetNotification(ctx, query)
	if err != nil {
		return err
	}
	q.Result = pbetop.PbNotification(noti)
	return nil
}

func GetNotifications(ctx context.Context, q *wrapshop.GetNotificationsEndpoint) error {
	paging := q.Paging.CMPaging()
	query := &notimodel.GetNotificationsArgs{
		Paging:    paging,
		AccountID: q.Context.Shop.ID,
	}
	notis, total, err := sqlstore.GetNotifications(ctx, query)
	if err != nil {
		return err
	}
	q.Result = &pbetop.NotificationsResponse{
		Notifications: pbetop.PbNotifications(notis),
		Paging:        pbcm.PbPageInfo(paging, int32(total)),
	}
	return nil
}

func UpdateNotifications(ctx context.Context, q *wrapshop.UpdateNotificationsEndpoint) error {
	cmd := &notimodel.UpdateNotificationsArgs{
		IDs:    q.Ids,
		IsRead: q.IsRead,
	}
	if err := sqlstore.UpdateNotifications(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.UpdatedResponse{
		Updated: int32(len(q.Ids)),
	}
	return nil
}

func GetShipnowFulfillment(ctx context.Context, q *wrapshop.GetShipnowFulfillmentEndpoint) error {
	query := &shipnow.GetShipnowFulfillmentQuery{
		Id:     q.Id,
		ShopId: q.Context.Shop.ID,
		Result: nil,
	}
	if err := shipnowQuery.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = pborder.Convert_core_ShipnowFulfillment_To_api_ShipnowFulfillment(query.Result.ShipnowFulfillment)
	return nil
}

func GetShipnowFulfillments(ctx context.Context, q *wrapshop.GetShipnowFulfillmentsEndpoint) error {
	shopIDs, err := api.MixAccount(q.Context.Claim, q.Mixed)
	if err != nil {
		return err
	}
	paging := q.Paging.CMPaging()

	query := &shipnow.GetShipnowFulfillmentsQuery{
		ShopIds: shopIDs,
		Paging:  paging,
		Filters: pbcm.ToFiltersPtr(q.Filters),
	}
	if err := shipnowQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pborder.ShipnowFulfillments{
		ShipnowFulfillments: pborder.Convert_core_ShipnowFulfillments_To_api_ShipnowFulfillments(query.Result.ShipnowFulfillments),
		Paging:              pbcm.PbPageInfo(paging, query.Result.Count),
	}
	return nil
}

func CreateShipnowFulfillment(ctx context.Context, q *wrapshop.CreateShipnowFulfillmentEndpoint) error {
	pickupAddress, err := q.PickupAddress.Fulfilled()
	if err != nil {
		return err
	}
	cmd := &shipnow.CreateShipnowFulfillmentCommand{
		OrderIds:            q.OrderIds,
		Carrier:             carrier.CarrierFromString(q.Carrier),
		ShopId:              q.Context.Shop.ID,
		ShippingServiceCode: q.ShippingServiceCode,
		ShippingServiceFee:  q.ShippingServiceFee,
		ShippingNote:        q.ShippingNote,
		RequestPickupAt:     nil,
		PickupAddress:       pborder.Convert_api_OrderAddress_To_core_OrderAddress(pickupAddress),
	}
	if err := shipnowAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = pborder.Convert_core_ShipnowFulfillment_To_api_ShipnowFulfillment(cmd.Result)
	return nil
}

func ConfirmShipnowFulfillment(ctx context.Context, q *wrapshop.ConfirmShipnowFulfillmentEndpoint) error {
	cmd := &shipnow.ConfirmShipnowFulfillmentCommand{
		Id:     q.Id,
		ShopId: q.Context.Shop.ID,
	}
	if err := shipnowAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = pborder.Convert_core_ShipnowFulfillment_To_api_ShipnowFulfillment(cmd.Result)
	return nil
}

func UpdateShipnowFulfillment(ctx context.Context, q *wrapshop.UpdateShipnowFulfillmentEndpoint) error {
	pickupAddress, err := q.PickupAddress.Fulfilled()
	if err != nil {
		return err
	}
	cmd := &shipnow.UpdateShipnowFulfillmentCommand{
		Id:                  q.Id,
		OrderIds:            q.OrderIds,
		Carrier:             carrier.CarrierFromString(q.Carrier),
		ShopId:              q.Context.Shop.ID,
		ShippingServiceCode: q.ShippingServiceCode,
		ShippingServiceFee:  q.ShippingServiceFee,
		ShippingNote:        q.ShippingNote,
		RequestPickupAt:     nil,
		PickupAddress:       pborder.Convert_api_OrderAddress_To_core_OrderAddress(pickupAddress),
	}
	if err := shipnowAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = pborder.Convert_core_ShipnowFulfillment_To_api_ShipnowFulfillment(cmd.Result)
	return nil
}

func CancelShipnowFulfillment(ctx context.Context, q *wrapshop.CancelShipnowFulfillmentEndpoint) error {
	cmd := &shipnow.CancelShipnowFulfillmentCommand{
		Id:           q.Id,
		ShopId:       q.Context.Shop.ID,
		CancelReason: q.CancelReason,
	}
	if err := shipnowAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}

	q.Result = &pbcm.UpdatedResponse{
		Updated: 1,
	}
	return nil
}

func GetShipnowServices(ctx context.Context, q *wrapshop.GetShipnowServicesEndpoint) error {
	pickupAddress, err := q.PickupAddress.Fulfilled()
	if err != nil {
		return err
	}
	var points []*shipnow.DeliveryPoint
	if len(q.DeliveryPoints) > 0 {
		for _, p := range q.DeliveryPoints {
			address, err := p.ShippingAddress.Fulfilled()
			if err != nil {
				return err
			}
			points = append(points, &shipnow.DeliveryPoint{
				ShippingAddress: pborder.Convert_api_OrderAddress_To_core_OrderAddress(address),
				ValueInfo: types.ValueInfo{
					CodAmount: p.CodAmount,
				},
			})
		}
	}

	cmd := &shipnow.GetShipnowServicesCommand{
		ShopId:         q.Context.Shop.ID,
		OrderIds:       q.OrderIds,
		PickupAddress:  pborder.Convert_api_OrderAddress_To_core_OrderAddress(pickupAddress),
		DeliveryPoints: points,
	}
	if err := shipnowAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pborder.GetShipnowServicesResponse{
		Services: pborder.Convert_core_ShipnowServices_To_api_ShipnowServices(cmd.Result.Services),
	}
	return nil
}

func CreateExternalAccountAhamove(ctx context.Context, q *wrapshop.CreateExternalAccountAhamoveEndpoint) error {
	query := &identity.GetUserByIDQuery{
		UserID: q.Context.Shop.OwnerID,
	}
	if err := identityQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	user := query.Result
	phone := user.Phone

	queryAddress := &address.GetAddressByIDQuery{
		ID: q.Context.Shop.AddressID,
	}
	if err := addressQuery.Dispatch(ctx, queryAddress); err != nil {
		return cm.Errorf(cm.FailedPrecondition, err, "Thiếu thông tin địa chỉ cửa hàng")
	}
	address := queryAddress.Result
	cmd := &identity.CreateExternalAccountAhamoveCommand{
		OwnerID: user.ID,
		Phone:   phone,
		Name:    user.FullName,
		Address: address.GetFullAddress(),
	}
	if err := identityAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = pbshop.Convert_core_XAhamoveAccount_To_api_XAhamoveAccount(cmd.Result)
	return nil
}

func GetExternalAccountAhamove(ctx context.Context, q *wrapshop.GetExternalAccountAhamoveEndpoint) error {
	queryUser := &identity.GetUserByIDQuery{
		UserID: q.Context.Shop.OwnerID,
	}
	if err := identityQuery.Dispatch(ctx, queryUser); err != nil {
		return err
	}
	user := queryUser.Result
	phone := user.Phone

	query := &identity.GetExternalAccountAhamoveQuery{
		Phone:   phone,
		OwnerID: user.ID,
	}
	if err := identityQuery.Dispatch(ctx, query); err != nil {
		return err
	}

	account := query.Result
	if !account.ExternalVerified && account.ExternalTicketID != "" {
		cmd := &identity.UpdateVerifiedExternalAccountAhamoveCommand{
			OwnerID: user.ID,
			Phone:   phone,
		}
		if err := identityAggr.Dispatch(ctx, cmd); err != nil {
			return err
		}
		account = cmd.Result
	}

	q.Result = pbshop.Convert_core_XAhamoveAccount_To_api_XAhamoveAccount(account)
	return nil
}

// func RequestVerifyExternalAccountAhamove(ctx context.Context, q *wrapshop.RequestVerifyExternalAccountAhamoveEndpoint) error {
// 	query := &model.GetUserByIDQuery{
// 		UserID: q.Context.Shop.OwnerID,
// 	}
// 	if err := bus.Dispatch(ctx, query); err != nil {
// 		return err
// 	}
// 	user := query.Result
// 	phone := user.Phone
//
// 	cmd := &identity.RequestVerifyExternalAccountAhamoveCommand{
// 		OwnerID: user.ID,
// 		Phone:   phone,
// 	}
// 	if err := identityAggr.Dispatch(ctx, cmd); err != nil {
// 		return err
// 	}
//
// 	q.Result = &pbcm.UpdatedResponse{
// 		Updated: 1,
// 	}
// 	return nil
// }
//
// func UpdateExternalAccountAhamoveVerificationImages(ctx context.Context, r *wrapshop.UpdateExternalAccountAhamoveVerificationImagesEndpoint) error {
// 	if err := validateImagesUrl(r.IdCardFrontImg, r.IdCardBackImg, r.PortraitImg); err != nil {
// 		return err
// 	}
//
// 	query := &model.GetUserByIDQuery{
// 		UserID: r.Context.Shop.OwnerID,
// 	}
// 	if err := bus.Dispatch(ctx, query); err != nil {
// 		return err
// 	}
// 	user := query.Result
// 	phone := user.Phone
//
// 	cmd := &identity.UpdateExternalAccountAhamoveVerificationImagesCommand{
// 		OwnerID:        user.ID,
// 		Phone:          phone,
// 		IDCardFrontImg: r.IdCardFrontImg,
// 		IDCardBackImg:  r.IdCardBackImg,
// 		PortraitImg:    r.PortraitImg,
// 	}
// 	if err := identityAggr.Dispatch(ctx, cmd); err != nil {
// 		return err
// 	}
//
// 	r.Result = &pbcm.UpdatedResponse{
// 		Updated: 1,
// 	}
// 	return nil
// }
//
// func validateImagesUrl(imgsUrl ...string) error {
// 	for _, url := range imgsUrl {
// 		if url == "" {
// 			continue
// 		}
// 		if !govalidator.IsURL(url) {
// 			return cm.Errorf(cm.InvalidArgument, nil, "Invalid url: %v", url)
// 		}
// 	}
// 	return nil
// }
