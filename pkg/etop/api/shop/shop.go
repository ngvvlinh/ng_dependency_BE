package shop

import (
	"context"

	"etop.vn/backend/pkg/services/moneytx/modelx"

	cmP "etop.vn/backend/pb/common"
	pbcm "etop.vn/backend/pb/common"
	etopP "etop.vn/backend/pb/etop"
	orderP "etop.vn/backend/pb/etop/order"
	shopP "etop.vn/backend/pb/etop/shop"
	supplierP "etop.vn/backend/pb/etop/supplier"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/idemp"
	"etop.vn/backend/pkg/common/l"
	"etop.vn/backend/pkg/common/redis"
	cmService "etop.vn/backend/pkg/common/service"
	"etop.vn/backend/pkg/etop/logic/shipping_provider"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/etop/sqlstore"
	notimodel "etop.vn/backend/pkg/notifier/model"
	shopW "etop.vn/backend/wrapper/etop/shop"
)

var ll = l.New()

func init() {
	bus.AddHandler("api", AddVariants)
	bus.AddHandler("api", CreateCollection)
	bus.AddHandler("api", DeleteCollection)
	bus.AddHandler("api", GetCollection)
	bus.AddHandler("api", GetCollections)
	bus.AddHandler("api", GetCollectionsByIDs)
	bus.AddHandler("api", GetPriceRules)
	bus.AddHandler("api", GetVariant)
	bus.AddHandler("api", GetVariants)
	bus.AddHandler("api", GetVariantsByIDs)
	bus.AddHandler("api", RemoveVariants)
	bus.AddHandler("api", RemoveProductsCollection)
	bus.AddHandler("api", UpdateCollection)
	bus.AddHandler("api", UpdatePriceRules)
	bus.AddHandler("api", UpdateVariant)
	bus.AddHandler("api", UpdateProducts)
	bus.AddHandler("api", UpdateProductsCollection)
	bus.AddHandler("api", UpdateVariantsStatus)
	bus.AddHandler("api", UpdateVariantsTags)
	bus.AddHandler("api", VersionInfo)
	bus.AddHandler("api", GetBrand)
	bus.AddHandler("api", GetBrands)

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
	bus.AddHandler("api", ConnectProductSource)
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
var shippingCtrl *shipping_provider.ProviderManager

func Init(shippingProviderCtrl *shipping_provider.ProviderManager, sd cmService.Shutdowner, rd redis.Store) {
	shippingCtrl = shippingProviderCtrl
	idempgroup = idemp.NewRedisGroup(rd, PrefixIdemp, 5*60)
	sd.Register(idempgroup.Shutdown)
}

func VersionInfo(ctx context.Context, q *shopW.VersionInfoEndpoint) error {
	q.Result = &cmP.VersionInfoResponse{
		Service: "etop.Shop",
		Version: "0.1",
	}
	return nil
}

func GetCollection(ctx context.Context, q *shopW.GetCollectionEndpoint) error {
	query := &model.GetShopCollectionQuery{
		ShopID:       q.Context.Shop.ID,
		CollectionID: q.Id,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = shopP.PbCollection(query.Result)
	return nil
}

func GetCollectionsByIDs(ctx context.Context, q *shopW.GetCollectionsByIDsEndpoint) error {
	query := &model.GetShopCollectionsQuery{
		ShopID:        q.Context.Shop.ID,
		CollectionIDs: q.Ids,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &shopP.CollectionsResponse{
		Collections: shopP.PbCollections(query.Result.Collections),
	}
	return nil
}

func GetCollections(ctx context.Context, q *shopW.GetCollectionsEndpoint) error {
	query := &model.GetShopCollectionsQuery{
		ShopID: q.Context.Shop.ID,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &shopP.CollectionsResponse{
		Collections: shopP.PbCollections(query.Result.Collections),
	}
	return nil
}

func GetVariant(ctx context.Context, q *shopW.GetVariantEndpoint) error {
	query := &model.GetShopVariantQuery{
		ShopID:    q.Context.Shop.ID,
		VariantID: q.Id,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = PbShopVariant(query.Result)
	return nil
}

func GetVariantsByIDs(ctx context.Context, q *shopW.GetVariantsByIDsEndpoint) error {
	query := &model.GetShopVariantsQuery{
		ShopID:     q.Context.Shop.ID,
		VariantIDs: q.Ids,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &shopP.ShopVariantsResponse{
		Variants: PbShopVariants(query.Result.Variants),
	}
	return nil
}

func GetVariants(ctx context.Context, q *shopW.GetVariantsEndpoint) error {
	paging := q.Paging.CMPaging()
	query := &model.GetShopVariantsQuery{
		ShopID:  q.Context.Shop.ID,
		Paging:  paging,
		Filters: cmP.ToFilters(q.Filters),
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &shopP.ShopVariantsResponse{
		Paging:   cmP.PbPageInfo(paging, query.Result.Total),
		Variants: PbShopVariants(query.Result.Variants),
	}
	return nil
}

func UpdateVariant(ctx context.Context, q *shopW.UpdateVariantEndpoint) error {
	shopID := q.Context.Shop.ID
	productSourceID := q.Context.Shop.ProductSourceID
	cmd := &model.UpdateShopVariantCommand{
		ShopID:          shopID,
		Variant:         shopP.PbUpdateVariantToModel(shopID, q.UpdateVariantRequest),
		CostPrice:       int(q.CostPrice),
		Inventory:       int(q.Inventory),
		EdCode:          q.Sku,
		Attributes:      supplierP.AttributesTomodel(q.Attributes),
		ProductSourceID: productSourceID,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = PbShopVariant(cmd.Result)
	return nil
}

func UpdateProducts(ctx context.Context, q *shopW.UpdateVariantsEndpoint) error {
	return cm.ErrTODO
}

func UpdateVariantsStatus(ctx context.Context, q *shopW.UpdateVariantsStatusEndpoint) error {
	if q.Status == nil {
		return cm.Error(cm.InvalidArgument, "Missing status", nil)
	}

	shopID := q.Context.Shop.ID
	productSourceID := q.Context.Shop.ProductSourceID
	cmd := &model.UpdateShopVariantsStatusCommand{
		ShopID:          shopID,
		VariantIDs:      q.Ids,
		ProductSourceID: productSourceID,
	}
	cmd.Update.Status = q.Status.ToModel()
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &cmP.UpdatedResponse{
		Updated: int32(cmd.Result.Updated),
	}
	return nil
}

func UpdateVariantsTags(ctx context.Context, q *shopW.UpdateVariantsTagsEndpoint) error {
	shopID := q.Context.Shop.ID
	productSourceID := q.Context.Shop.ProductSourceID
	cmd := &model.UpdateShopVariantsTagsCommand{
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
	q.Result = &cmP.UpdatedResponse{
		Updated: int32(cmd.Result.Updated),
	}
	return nil
}

func AddVariants(ctx context.Context, q *shopW.AddVariantsEndpoint) error {
	cmd := &model.AddShopVariantsCommand{
		ShopID: q.Context.Shop.ID,
		IDs:    q.Ids,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &shopP.AddVariantsResponse{
		Variants: PbShopVariants(cmd.Result.Variants),
		Errors:   cmP.PbErrors(cmd.Result.Errors),
	}
	return nil
}

func GetPriceRules(ctx context.Context, q *shopW.GetPriceRulesEndpoint) error {
	return cm.ErrTODO
}

func UpdatePriceRules(ctx context.Context, q *shopW.UpdatePriceRulesEndpoint) error {
	return cm.ErrTODO
}

func CreateCollection(ctx context.Context, q *shopW.CreateCollectionEndpoint) error {
	cmd := &model.CreateShopCollectionCommand{
		Collection: shopP.PbCreateCollection(q.Context.Shop.ID, q.CreateCollectionRequest),
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = shopP.PbCollection(cmd.Result)
	return nil
}

func DeleteCollection(ctx context.Context, q *shopW.DeleteCollectionEndpoint) error {
	cmd := &model.RemoveShopCollectionCommand{
		ShopID:       q.Context.Shop.ID,
		CollectionID: q.Id,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &cmP.RemovedResponse{
		Removed: int32(cmd.Result.Deleted),
	}
	return nil
}

func RemoveVariants(ctx context.Context, q *shopW.RemoveVariantsEndpoint) error {
	cmd := &model.RemoveShopVariantsCommand{
		ShopID:          q.Context.Shop.ID,
		IDs:             q.Ids,
		ProductSourceID: q.Context.Shop.ProductSourceID,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &cmP.RemovedResponse{
		Removed: int32(cmd.Result.Removed),
	}
	return nil
}

func UpdateCollection(ctx context.Context, q *shopW.UpdateCollectionEndpoint) error {
	cmd := &model.UpdateShopCollectionCommand{
		Collection: shopP.PbUpdateCollection(q.Context.Shop.ID, q.UpdateCollectionRequest),
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = shopP.PbCollection(cmd.Result)
	return nil
}

func UpdateProductsCollection(ctx context.Context, q *shopW.UpdateProductsCollectionEndpoint) error {
	cmd := &model.AddProductsToShopCollectionCommand{
		ShopID:       q.Context.Shop.ID,
		ProductIDs:   q.ProductIds,
		CollectionID: q.CollectionId,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &shopP.UpdateProductsCollectionResponse{
		Updated: int32(cmd.Result.Updated),
		Errors:  cmP.PbErrors(cmd.Result.Errors),
	}
	return nil
}

func RemoveProductsCollection(ctx context.Context, q *shopW.RemoveProductsCollectionEndpoint) error {
	cmd := &model.RemoveProductsFromShopCollectionCommand{
		ShopID:       q.Context.Shop.ID,
		ProductIDs:   q.ProductIds,
		CollectionID: q.CollectionId,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = cmP.Updated(cmd.Result.Updated)
	return nil
}

func GetBrand(ctx context.Context, q *shopW.GetBrandEndpoint) error {
	query := &model.GetProductBrandQuery{
		ID: q.Id,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = supplierP.PbBrandExt(query.Result)
	return nil
}

func GetBrands(ctx context.Context, q *shopW.GetBrandsEndpoint) error {
	query := &model.GetProductBrandsQuery{}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = &supplierP.BrandsResponse{
		Brands: supplierP.PbBrandsExt(query.Result.Brands),
	}
	return nil
}

func AddProducts(ctx context.Context, q *shopW.AddProductsEndpoint) error {
	cmd := &model.AddShopProductsCommand{
		ShopID: q.Context.Shop.ID,
		IDs:    q.Ids,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &shopP.AddProductsResponse{
		Products: PbShopProducts(cmd.Result.Products),
		Errors:   cmP.PbErrors(cmd.Result.Errors),
	}
	return nil
}

func GetProduct(ctx context.Context, q *shopW.GetProductEndpoint) error {
	productSourceID := q.Context.Shop.ProductSourceID
	query := &model.GetShopProductQuery{
		ShopID:          q.Context.Shop.ID,
		ProductID:       q.Id,
		ProductSourceID: productSourceID,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = PbShopProductFtVariant(query.Result)
	return nil
}

func GetProductsByIDs(ctx context.Context, q *shopW.GetProductsByIDsEndpoint) error {
	productSourceID := q.Context.Shop.ProductSourceID
	query := &model.GetShopProductsQuery{
		ShopID:          q.Context.Shop.ID,
		ProductIDs:      q.Ids,
		ProductSourceID: productSourceID,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &shopP.ShopProductsResponse{
		Products: PbShopProductsFtVariant(query.Result.Products),
	}
	return nil
}

func GetProducts(ctx context.Context, q *shopW.GetProductsEndpoint) error {
	paging := q.Paging.CMPaging()
	productSourceID := q.Context.Shop.ProductSourceID
	query := &model.GetShopProductsQuery{
		ShopID:          q.Context.Shop.ID,
		Paging:          paging,
		Filters:         cmP.ToFilters(q.Filters),
		ProductSourceID: productSourceID,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &shopP.ShopProductsResponse{
		Paging:   cmP.PbPageInfo(paging, query.Result.Total),
		Products: PbShopProductsFtVariant(query.Result.Products),
	}
	return nil
}

func RemoveProducts(ctx context.Context, q *shopW.RemoveProductsEndpoint) error {
	productSourceID := q.Context.Shop.ProductSourceID
	cmd := &model.RemoveShopProductsCommand{
		ShopID:          q.Context.Shop.ID,
		IDs:             q.Ids,
		ProductSourceID: productSourceID,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &cmP.RemovedResponse{
		Removed: int32(cmd.Result.Removed),
	}
	return nil
}

func UpdateProduct(ctx context.Context, q *shopW.UpdateProductEndpoint) error {
	shopID := q.Context.Shop.ID
	productSourceID := q.Context.Shop.ProductSourceID
	cmd := &model.UpdateShopProductCommand{
		ShopID:          shopID,
		Product:         shopP.PbUpdateProductToModel(shopID, q.UpdateProductRequest),
		Code:            q.Code,
		ProductSourceID: productSourceID,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = PbShopProductFtVariant(cmd.Result)
	return nil
}

func UpdateProductsStatus(ctx context.Context, q *shopW.UpdateProductsStatusEndpoint) error {
	if q.Status == nil {
		return cm.Error(cm.InvalidArgument, "Missing status", nil)
	}

	shopID := q.Context.Shop.ID
	productSourceID := q.Context.Shop.ProductSourceID
	cmd := &model.UpdateShopProductsStatusCommand{
		ShopID:          shopID,
		ProductIDs:      q.Ids,
		ProductSourceID: productSourceID,
	}
	cmd.Update.Status = q.Status.ToModel()
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &cmP.UpdatedResponse{
		Updated: int32(cmd.Result.Updated),
	}
	return nil
}

func UpdateProductsTags(ctx context.Context, q *shopW.UpdateProductsTagsEndpoint) error {
	shopID := q.Context.Shop.ID
	productSourceID := q.Context.Shop.ProductSourceID
	cmd := &model.UpdateShopProductsTagsCommand{
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
	q.Result = &cmP.UpdatedResponse{
		Updated: int32(cmd.Result.Updated),
	}
	return nil
}

func CreateProductSource(ctx context.Context, q *shopW.CreateProductSourceEndpoint) error {
	shopID := q.Context.Shop.ID
	cmd := &model.CreateProductSourceCommand{
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

func CreateVariant(ctx context.Context, q *shopW.CreateVariantEndpoint) error {
	cmd := &model.CreateVariantCommand{
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
		ListPrice:         int(q.ListPrice),
		SKU:               q.Sku,
		Code:              q.Code,
		QuantityAvailable: int(q.QuantityAvailable),
		QuantityOnHand:    int(q.QuantityOnHand),
		QuantityReserved:  int(q.QuantityReserved),
		CostPrice:         int(q.CostPrice),
		Attributes:        supplierP.AttributesTomodel(q.Attributes),
		DescHTML:          q.DescHtml,
	}

	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}

	q.Result = PbShopProductFtVariant(cmd.Result)
	return nil
}

func GetShopProductSources(ctx context.Context, q *shopW.GetShopProductSourcesEndpoint) error {
	query := &model.GetShopProductSourcesCommand{}
	if q.Context.User != nil {
		query.UserID = q.Context.User.ID
	} else {
		query.ShopID = q.Context.Shop.ID
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = &shopP.ProductSourcesResponse{
		ProductSources: PbProductSources(query.Result),
	}
	return nil
}

func ConnectProductSource(ctx context.Context, q *shopW.ConnectProductSourceEndpoint) error {
	cmd := &model.ConnectProductSourceCommand{
		ShopID:          q.Context.Shop.ID,
		ProductSourceID: q.ProductSourceId,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}

	q.Result = &cmP.UpdatedResponse{
		Updated: int32(cmd.Result.Updated),
	}
	return nil
}

// func RemoveProductSource(ctx context.Context, q *shopW.RemoveProductSourceEndpoint) error {
// 	cmd := &model.RemoveProductSourceCommand{
// 		AccountID: q.Context.Shop.ID,
// 	}

// 	if err := bus.Dispatch(ctx, cmd); err != nil {
// 		return err
// 	}

// 	q.Result = &cmP.UpdatedResponse{
// 		Updated: int32(cmd.Result.Updated),
// 	}

// 	return nil
// }

func CreateProductSourceCategory(ctx context.Context, q *shopW.CreateProductSourceCategoryEndpoint) error {
	cmd := &model.CreateProductSourceCategoryCommand{
		ShopID:            q.Context.Shop.ID,
		Name:              q.Name,
		ParentID:          q.ParentId,
		ProductSourceID:   q.ProductSourceId,
		ProductSourceType: q.ProductSourceType.ToModel(),
	}

	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = PbProductSourceCategory(cmd.Result)
	return nil
}

func UpdateProductsPSCategory(ctx context.Context, q *shopW.UpdateProductsPSCategoryEndpoint) error {
	cmd := &model.UpdateProductsProductSourceCategoryCommand{
		CategoryID:      q.CategoryId,
		ProductIDs:      q.ProductIds,
		ShopID:          q.Context.Shop.ID,
		ProductSourceID: q.Context.Shop.ProductSourceID,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &cmP.UpdatedResponse{
		Updated: int32(cmd.Result.Updated),
	}
	return nil
}

func GetProductSourceCategory(ctx context.Context, q *shopW.GetProductSourceCategoryEndpoint) error {
	cmd := &model.GetProductSourceCategoryQuery{
		ShopID:     q.Context.Shop.ID,
		CategoryID: q.Id,
	}

	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}

	q.Result = supplierP.PbCategory(cmd.Result)
	return nil
}

func GetProductSourceCategories(ctx context.Context, q *shopW.GetProductSourceCategoriesEndpoint) error {
	cmd := &model.GetProductSourceCategoriesExtendedQuery{
		ProductSourceType: q.Type.ToModel(),
		ShopID:            q.Context.Shop.ID,
	}

	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}

	q.Result = &supplierP.CategoriesResponse{
		Categories: supplierP.PbCategories(cmd.Result.Categories),
	}
	return nil
}

func UpdateProductSourceCategory(ctx context.Context, q *shopW.UpdateProductSourceCategoryEndpoint) error {
	cmd := &model.UpdateShopProductSourceCategoryCommand{
		ID:       q.Id,
		ShopID:   q.Context.Shop.ID,
		ParentID: q.ParentId,
		Name:     q.Name,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = supplierP.PbCategory(cmd.Result)
	return nil
}

func RemoveProductSourceCategory(ctx context.Context, q *shopW.RemoveProductSourceCategoryEndpoint) error {
	cmd := &model.RemoveShopProductSourceCategoryCommand{
		ID:     q.Id,
		ShopID: q.Context.Shop.ID,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &cmP.RemovedResponse{
		Removed: int32(cmd.Result.Removed),
	}
	return nil
}

func UpdateProductImages(ctx context.Context, q *shopW.UpdateProductImagesEndpoint) error {
	shopID := q.Context.Shop.ID
	query := &model.GetShopProductQuery{
		ShopID:    shopID,
		ProductID: q.Id,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}

	r := &model.UpdateListRequest{
		Adds:       q.Adds,
		Deletes:    q.Deletes,
		ReplaceAll: q.ReplaceAll,
		DeleteAll:  q.DeleteAll,
	}

	imageURLs, err := cmP.PatchImage(query.Result.ShopProduct.ImageURLs, r)
	if err != nil {
		return err
	}

	cmd := &model.UpdateShopProductCommand{
		ShopID: shopID,
		Product: &model.ShopProduct{
			ProductID: q.Id,
			ShopID:    shopID,
			ImageURLs: imageURLs,
		},
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = PbShopProductFtVariant(cmd.Result)
	return nil
}

func UpdateVariantImages(ctx context.Context, q *shopW.UpdateVariantImagesEndpoint) error {
	shopID := q.Context.Shop.ID
	query := &model.GetShopVariantQuery{
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
	imageURLs, err := cmP.PatchImage(sourceImages, r)
	if err != nil {
		return err
	}

	cmd := &model.UpdateShopVariantCommand{
		ShopID: shopID,
		Variant: &model.ShopVariant{
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

func GetMoneyTransaction(ctx context.Context, q *shopW.GetMoneyTransactionEndpoint) error {
	query := &modelx.GetMoneyTransaction{
		ShopID: q.Context.Shop.ID,
		ID:     q.Id,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = orderP.PbMoneyTransactionExtended(query.Result)
	return nil
}

func GetMoneyTransactions(ctx context.Context, q *shopW.GetMoneyTransactionsEndpoint) error {
	paging := q.Paging.CMPaging()
	query := &modelx.GetMoneyTransactions{
		ShopID:  q.Context.Shop.ID,
		Paging:  paging,
		Filters: pbcm.ToFilters(q.Filters),
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &orderP.MoneyTransactionsResponse{
		MoneyTransactions: orderP.PbMoneyTransactionExtendeds(query.Result.MoneyTransactions),
		Paging:            cmP.PbPageInfo(paging, query.Result.Total),
	}
	return nil
}

func SummarizeFulfillments(ctx context.Context, q *shopW.SummarizeFulfillmentsEndpoint) error {
	query := &model.SummarizeFulfillmentsRequest{
		ShopID:   q.Context.Shop.ID,
		DateFrom: q.DateFrom,
		DateTo:   q.DateTo,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = &shopP.SummarizeFulfillmentsResponse{
		Tables: shopP.PbSummaryTables(query.Result.Tables),
	}
	return nil
}

func CalcBalance(ctx context.Context, q *shopW.CalcBalanceShopEndpoint) error {
	query := &model.GetBalanceShopCommand{
		ShopID: q.Context.Shop.ID,
	}

	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &shopP.CalcBalanceShopResponse{
		Balance: int32(query.Result.Amount),
	}
	return nil
}

func CreateDevice(ctx context.Context, q *shopW.CreateDeviceEndpoint) error {
	cmd := &notimodel.CreateDeviceArgs{
		AccountID:        q.Context.Shop.ID,
		DeviceID:         q.DeviceId,
		DeviceName:       q.DeviceName,
		ExternalDeviceID: q.ExternalDeviceId,
	}
	device, err := sqlstore.CreateDevice(ctx, cmd)
	if err != nil {
		return err
	}
	q.Result = etopP.PbDevice(device)
	return nil
}

func DeleteDevice(ctx context.Context, q *shopW.DeleteDeviceEndpoint) error {
	device := &notimodel.Device{
		DeviceID:  q.DeviceId,
		AccountID: q.Context.Shop.ID,
	}
	if err := sqlstore.DeleteDevice(ctx, device); err != nil {
		return err
	}
	q.Result = &pbcm.DeletedResponse{
		Deleted: 1,
	}
	return nil
}

func GetNotification(ctx context.Context, q *shopW.GetNotificationEndpoint) error {
	query := &notimodel.GetNotificationArgs{
		AccountID: q.Context.Shop.ID,
		ID:        q.Id,
	}
	noti, err := sqlstore.GetNotification(ctx, query)
	if err != nil {
		return err
	}
	q.Result = etopP.PbNotification(noti)
	return nil
}

func GetNotifications(ctx context.Context, q *shopW.GetNotificationsEndpoint) error {
	paging := q.Paging.CMPaging()
	query := &notimodel.GetNotificationsArgs{
		Paging:    paging,
		AccountID: q.Context.Shop.ID,
	}
	notis, total, err := sqlstore.GetNotifications(ctx, query)
	if err != nil {
		return err
	}
	q.Result = &etopP.NotificationsResponse{
		Notifications: etopP.PbNotifications(notis),
		Paging:        cmP.PbPageInfo(paging, total),
	}
	return nil
}

func UpdateNotifications(ctx context.Context, q *shopW.UpdateNotificationsEndpoint) error {
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
