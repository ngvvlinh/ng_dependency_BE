package shop

import (
	"context"

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
	"etop.vn/backend/pkg/etop/logic/shipping_provider"
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
	bus.AddHandler("api", AddVariants)
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
	bus.AddHandler("api", UpdateVariantsStatus)
	bus.AddHandler("api", UpdateVariantsTags)
	bus.AddHandler("api", VersionInfo)

	bus.AddHandler("api", AddProducts)
	bus.AddHandler("api", GetProduct)
	bus.AddHandler("api", GetProducts)
	bus.AddHandler("api", GetProductsByIDs)
	bus.AddHandler("api", UpdateProduct)
	bus.AddHandler("api", UpdateProductsStatus)
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

	bus.AddHandler("api", CreateShipNowFulfillment)
	bus.AddHandler("api", GetShipNowFulfillment)
	bus.AddHandler("api", GetShipNowFulfillments)
	bus.AddHandler("api", ConfirmShipNowFulfillment)
	bus.AddHandler("api", UpdateShipNowFulfillment)
	bus.AddHandler("api", CancelShipNowFulfillment)
}

const PrefixIdemp = "IdempOrder"

var idempgroup *idemp.RedisGroup
var shippingCtrl *shipping_provider.ProviderManager
var shipnowAggr shipnow.Aggregate
var shipnowQuery shipnow.QueryService
var catalogQuery catalog.QueryBus

func Init(catalogQueryBus catalog.QueryBus, shipnowAggregate shipnow.Aggregate, shipnowQueryService shipnow.QueryService, shippingProviderCtrl *shipping_provider.ProviderManager, sd cmservice.Shutdowner, rd redis.Store) {
	shippingCtrl = shippingProviderCtrl
	idempgroup = idemp.NewRedisGroup(rd, PrefixIdemp, 5*60)
	catalogQuery = catalogQueryBus
	shipnowAggr = shipnowAggregate
	shipnowQuery = shipnowQueryService
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

func UpdateVariantsStatus(ctx context.Context, q *wrapshop.UpdateVariantsStatusEndpoint) error {
	if q.Status == nil {
		return cm.Error(cm.InvalidArgument, "Missing status", nil)
	}

	shopID := q.Context.Shop.ID
	productSourceID := q.Context.Shop.ProductSourceID
	cmd := &catalogmodelx.UpdateShopVariantsStatusCommand{
		ShopID:          shopID,
		VariantIDs:      q.Ids,
		ProductSourceID: productSourceID,
	}
	cmd.Update.Status = q.Status.ToModel()
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.UpdatedResponse{
		Updated: int32(cmd.Result.Updated),
	}
	return nil
}

func UpdateVariantsTags(ctx context.Context, q *wrapshop.UpdateVariantsTagsEndpoint) error {
	shopID := q.Context.Shop.ID
	productSourceID := q.Context.Shop.ProductSourceID
	cmd := &catalogmodelx.UpdateShopVariantsTagsCommand{
		ShopID:     shopID,
		VariantIDs: q.Ids,
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

func AddVariants(ctx context.Context, q *wrapshop.AddVariantsEndpoint) error {
	cmd := &catalogmodelx.AddShopVariantsCommand{
		ShopID: q.Context.Shop.ID,
		IDs:    q.Ids,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbshop.AddVariantsResponse{
		Variants: PbShopVariants(cmd.Result.Variants),
		Errors:   pbcm.PbErrors(cmd.Result.Errors),
	}
	return nil
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

func UpdateProductsStatus(ctx context.Context, q *wrapshop.UpdateProductsStatusEndpoint) error {
	if q.Status == nil {
		return cm.Error(cm.InvalidArgument, "Missing status", nil)
	}

	shopID := q.Context.Shop.ID
	productSourceID := q.Context.Shop.ProductSourceID
	cmd := &catalogmodelx.UpdateShopProductsStatusCommand{
		ShopID:          shopID,
		ProductIDs:      q.Ids,
		ProductSourceID: productSourceID,
	}
	cmd.Update.Status = q.Status.ToModel()
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.UpdatedResponse{
		Updated: int32(cmd.Result.Updated),
	}
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

func GetShipNowFulfillment(ctx context.Context, q *wrapshop.GetShipnowFulfillmentEndpoint) error {
	args := &shipnow.GetShipnowFulfillmentQueryArgs{
		Id: q.Id,
	}
	ffm, err := shipnowQuery.GetShipnowFulfillment(ctx, args)
	if err != nil {
		return err
	}
	q.Result = Convert_core_ShipnowFulfillment_To_api_ShipnowFulfillment(ffm.ShipnowFulfillment)
	return nil
}

func GetShipNowFulfillments(ctx context.Context, q *wrapshop.GetShipnowFulfillmentsEndpoint) error {
	args := &shipnow.GetShipnowFulfillmentsQueryArgs{
		ShopId: q.Context.Shop.ID,
	}
	ffms, err := shipnowQuery.GetShipnowFulfillments(ctx, args)
	if err != nil {
		return err
	}
	q.Result = &pborder.ShipnowFulfillments{
		ShipnowFulfillments: Convert_core_ShipnowFulfillments_To_api_ShipnowFulfillments(ffms.ShipnowFulfillments),
	}
	return nil
}

func CreateShipNowFulfillment(ctx context.Context, q *wrapshop.CreateShipnowFulfillmentEndpoint) error {
	pickupAddress, err := q.PickupAddress.Fulfilled()
	if err != nil {
		return err
	}
	args := &shipnow.CreateShipnowFulfillmentArgs{
		OrderIds:            q.OrderIds,
		Carrier:             q.Carrier,
		ShopId:              q.Context.Shop.ID,
		ShippingServiceCode: q.ShippingServiceCode,
		ShippingServiceFee:  q.ShippingServiceFee,
		ShippingNote:        q.ShippingNote,
		RequestPickupAt:     nil,
		PickupAddress:       Convert_api_OrderAddress_To_core_OrderAddress(pickupAddress),
	}
	res, err := shipnowAggr.CreateShipnowFulfillment(ctx, args)
	if err != nil {
		return err
	}
	q.Result = Convert_core_ShipnowFulfillment_To_api_ShipnowFulfillment(res)
	return nil
}

func ConfirmShipNowFulfillment(ctx context.Context, q *wrapshop.ConfirmShipnowFulfillmentEndpoint) error {
	args := &shipnow.ConfirmShipnowFulfillmentArgs{
		Id:     q.Id,
		ShopId: q.Context.Shop.ID,
	}
	res, err := shipnowAggr.ConfirmShipnowFulfillment(ctx, args)
	if err != nil {
		return err
	}
	q.Result = Convert_core_ShipnowFulfillment_To_api_ShipnowFulfillment(res)
	return nil
}

func UpdateShipNowFulfillment(ctx context.Context, q *wrapshop.UpdateShipnowFulfillmentEndpoint) error {
	pickupAddress, err := q.PickupAddress.Fulfilled()
	if err != nil {
		return err
	}
	args := &shipnow.UpdateShipnowFulfillmentArgs{
		Id:                  q.Id,
		OrderIds:            q.OrderIds,
		Carrier:             q.Carrier,
		ShopId:              q.Context.Shop.ID,
		ShippingServiceCode: q.ShippingServiceCode,
		ShippingServiceFee:  q.ShippingServiceFee,
		ShippingNote:        q.ShippingNote,
		RequestPickupAt:     nil,
		PickupAddress:       Convert_api_OrderAddress_To_core_OrderAddress(pickupAddress),
	}
	res, err := shipnowAggr.UpdateShipnowFulfillment(ctx, args)
	if err != nil {
		return err
	}
	q.Result = Convert_core_ShipnowFulfillment_To_api_ShipnowFulfillment(res)
	return nil
}

func CancelShipNowFulfillment(ctx context.Context, q *wrapshop.CancelShipnowFulfillmentEndpoint) error {
	args := &shipnow.CancelShipnowFulfillmentArgs{
		Id:           q.Id,
		ShopId:       q.Context.Shop.ID,
		CancelReason: q.CancelReason,
	}
	if _, err := shipnowAggr.CancelShipnowFulfillment(ctx, args); err != nil {
		return err
	}
	q.Result = &pbcm.UpdatedResponse{
		Updated: 1,
	}
	return nil
}
