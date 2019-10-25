package shop

import (
	"context"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"

	haravanidentity "etop.vn/api/external/haravan/identity"
	"etop.vn/api/external/payment"
	paymentmanager "etop.vn/api/external/payment/manager"
	"etop.vn/api/main/address"
	"etop.vn/api/main/catalog"
	"etop.vn/api/main/identity"
	"etop.vn/api/main/inventory"
	"etop.vn/api/main/location"
	"etop.vn/api/main/ordering"
	"etop.vn/api/main/receipting"
	"etop.vn/api/main/shipnow"
	carriertypes "etop.vn/api/main/shipnow/carrier/types"
	"etop.vn/api/main/shipping/types"
	"etop.vn/api/meta"
	"etop.vn/api/shopping/addressing"
	"etop.vn/api/shopping/carrying"
	"etop.vn/api/shopping/customering"
	"etop.vn/api/shopping/tradering"
	"etop.vn/api/shopping/vendoring"
	"etop.vn/api/summary"
	notimodel "etop.vn/backend/com/handler/notifier/model"
	catalogmodelx "etop.vn/backend/com/main/catalog/modelx"
	moneymodelx "etop.vn/backend/com/main/moneytx/modelx"
	pbcm "etop.vn/backend/pb/common"
	pbetop "etop.vn/backend/pb/etop"
	pborder "etop.vn/backend/pb/etop/order"
	pbshop "etop.vn/backend/pb/etop/shop"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/idemp"
	"etop.vn/backend/pkg/common/redis"
	cmservice "etop.vn/backend/pkg/common/service"
	"etop.vn/backend/pkg/etop/api"
	"etop.vn/backend/pkg/etop/api/convertpb"
	"etop.vn/backend/pkg/etop/logic/shipping_provider"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/etop/sqlstore"
	wrapshop "etop.vn/backend/wrapper/etop/shop"
	"etop.vn/capi"
	. "etop.vn/capi/dot"
	"etop.vn/common/l"
)

var (
	ll = l.New()
	s  = &Service{}
)

func init() {
	bus.AddHandler("api", s.VersionInfo)

	bus.AddHandler("api", s.CreateInventoryVoucher)
	bus.AddHandler("api", s.UpdateInventoryVoucher)
	bus.AddHandler("api", s.AdjustInventoryQuantity)
	bus.AddHandler("api", s.ConfirmInventoryVoucher)
	bus.AddHandler("api", s.CancelInventoryVoucher)
	bus.AddHandler("api", s.GetInventory)
	bus.AddHandler("api", s.GetInventories)
	bus.AddHandler("api", s.GetInventoryByVariantIDs)
	bus.AddHandler("api", s.GetInventoryVoucher)
	bus.AddHandler("api", s.GetInventoryVouchers)
	bus.AddHandler("api", s.GetInventoryVoucherByIDs)

	bus.AddHandler("api", s.RemoveVariants)
	bus.AddHandler("api", s.UpdateVariant)
	bus.AddHandler("api", s.UpdateVariantAttributes)
	bus.AddHandler("api", s.UpdateVariantsStatus)

	bus.AddHandler("api", s.AddProducts)
	bus.AddHandler("api", s.GetProduct)
	bus.AddHandler("api", s.GetProducts)
	bus.AddHandler("api", s.GetProductsByIDs)
	bus.AddHandler("api", s.CreateProduct)
	bus.AddHandler("api", s.UpdateProduct)
	bus.AddHandler("api", s.UpdateProductsStatus)
	bus.AddHandler("api", s.UpdateProductsTags)
	bus.AddHandler("api", s.RemoveProducts)

	bus.AddHandler("api", s.CreateVariant)
	bus.AddHandler("api", s.GetVariant)
	bus.AddHandler("api", s.GetVariantsByIDs)
	bus.AddHandler("api", s.DeprecatedCreateVariant)
	bus.AddHandler("api", s.CreateProductSourceCategory)
	bus.AddHandler("api", s.UpdateProductsPSCategory)
	bus.AddHandler("api", s.UpdateProductImages)
	bus.AddHandler("api", s.UpdateProductMetaFields)
	bus.AddHandler("api", s.GetProductSourceCategory)
	bus.AddHandler("api", s.GetProductSourceCategories)
	bus.AddHandler("api", s.UpdateProductSourceCategory)
	bus.AddHandler("api", s.RemoveProductSourceCategory)
	bus.AddHandler("api", s.UpdateVariantImages)

	bus.AddHandler("api", s.GetMoneyTransaction)
	bus.AddHandler("api", s.GetMoneyTransactions)

	bus.AddHandler("api", s.SummarizeFulfillments)
	bus.AddHandler("api", s.CalcBalance)
	bus.AddHandler("api", s.CreateDevice)
	bus.AddHandler("api", s.DeleteDevice)
	bus.AddHandler("api", s.GetNotifications)
	bus.AddHandler("api", s.GetNotification)
	bus.AddHandler("api", s.UpdateNotifications)

	bus.AddHandler("api", s.GetShipnowFulfillment)
	bus.AddHandler("api", s.GetShipnowFulfillments)
	bus.AddHandler("api", s.CreateShipnowFulfillment)
	bus.AddHandler("api", s.ConfirmShipnowFulfillment)
	bus.AddHandler("api", s.UpdateShipnowFulfillment)
	bus.AddHandler("api", s.CancelShipnowFulfillment)
	bus.AddHandler("api", s.GetShipnowServices)
	bus.AddHandler("api", s.CreateExternalAccountAhamove)
	bus.AddHandler("api", s.GetExternalAccountAhamove)
	bus.AddHandler("api", s.RequestVerifyExternalAccountAhamove)
	bus.AddHandler("api", s.UpdateExternalAccountAhamoveVerification)

	bus.AddHandler("api", s.GetExternalAccountHaravan)
	bus.AddHandler("api", s.CreateExternalAccountHaravan)
	bus.AddHandler("api", s.UpdateExternalAccountHaravanToken)
	bus.AddHandler("api", s.ConnectCarrierServiceExternalAccountHaravan)
	bus.AddHandler("api", s.DeleteConnectedCarrierServiceExternalAccountHaravan)
	bus.AddHandler("api", s.PaymentTradingOrder)
	bus.AddHandler("api", s.PaymentCheckReturnData)

	bus.AddHandler("api", s.CreateCategory)
	bus.AddHandler("api", s.GetCategory)
	bus.AddHandler("api", s.GetCategories)
	bus.AddHandler("api", s.UpdateCategory)
	bus.AddHandler("api", s.DeleteCategory)

	bus.AddHandler("api", s.UpdateProductCategory)
	bus.AddHandler("api", s.RemoveProductCategory)

	bus.AddHandler("api", s.GetCollection)
	bus.AddHandler("api", s.GetCollections)
	bus.AddHandler("api", s.CreateCollection)
	bus.AddHandler("api", s.UpdateCollection)

	bus.AddHandler("api", s.AddProductCollections)
	bus.AddHandler("api", s.RemoveProductCollection)
	bus.AddHandler("api", s.GetCollectionsByProductID)
	bus.AddHandler("api", s.GetSummarizePOS)
}

const (
	PrefixIdemp = "IdempShop"
)

var (
	locationQuery        location.QueryBus
	idempgroup           *idemp.RedisGroup
	idempgroupReceipt    *idemp.RedisGroup
	shipnowAggr          shipnow.CommandBus
	shipnowQuery         shipnow.QueryBus
	identityAggr         identity.CommandBus
	identityQuery        identity.QueryBus
	addressQuery         address.QueryBus
	shippingCtrl         *shipping_provider.ProviderManager
	catalogQuery         catalog.QueryBus
	catalogAggr          catalog.CommandBus
	haravanIdentityAggr  haravanidentity.CommandBus
	haravanIdentityQuery haravanidentity.QueryBus
	customerQuery        customering.QueryBus
	customerAggr         customering.CommandBus
	orderAggr            ordering.CommandBus
	orderQuery           ordering.QueryBus
	traderAddressAggr    addressing.CommandBus
	traderAddressQuery   addressing.QueryBus
	paymentCtrl          paymentmanager.CommandBus
	vendorAggr           vendoring.CommandBus
	vendorQuery          vendoring.QueryBus
	carrierAggr          carrying.CommandBus
	carrierQuery         carrying.QueryBus
	traderQuery          tradering.QueryBus
	summaryQuery         summary.QueryBus
	eventBus             capi.EventBus
	receiptAggr          receipting.CommandBus
	receiptQuery         receipting.QueryBus
	inventoryAggregate   inventory.CommandBus
	inventoryQuery       inventory.QueryBus
)

func Init(
	locationQ location.QueryBus,
	catalogQueryBus catalog.QueryBus,
	catalogCommandBus catalog.CommandBus,
	shipnow shipnow.CommandBus,
	shipnowQS shipnow.QueryBus,
	identity identity.CommandBus,
	identityQS identity.QueryBus,
	addressQS address.QueryBus,
	providerManager *shipping_provider.ProviderManager,
	haravanIdentity haravanidentity.CommandBus,
	haravanIdentityQS haravanidentity.QueryBus,
	customerA customering.CommandBus,
	customerQS customering.QueryBus,
	traderAddressA addressing.CommandBus,
	traderAddressQ addressing.QueryBus,
	orderA ordering.CommandBus,
	orderQ ordering.QueryBus,
	paymentManager paymentmanager.CommandBus,
	vendorA vendoring.CommandBus,
	vendorQ vendoring.QueryBus,
	carrierA carrying.CommandBus,
	carrierQ carrying.QueryBus,
	traderQ tradering.QueryBus,
	eventB capi.EventBus,
	receiptA receipting.CommandBus,
	receiptQS receipting.QueryBus,
	sd cmservice.Shutdowner,
	rd redis.Store,
	inventoryA inventory.CommandBus,
	inventoryQ inventory.QueryBus,

	summary summary.QueryBus,
) {
	idempgroup = idemp.NewRedisGroup(rd, PrefixIdemp, 5*60)
	locationQuery = locationQ
	catalogQuery = catalogQueryBus
	catalogAggr = catalogCommandBus
	shippingCtrl = providerManager
	shipnowAggr = shipnow
	shipnowQuery = shipnowQS
	identityAggr = identity
	identityQuery = identityQS
	addressQuery = addressQS
	haravanIdentityAggr = haravanIdentity
	haravanIdentityQuery = haravanIdentityQS
	customerQuery = customerQS
	customerAggr = customerA
	traderAddressAggr = traderAddressA
	traderAddressQuery = traderAddressQ
	orderAggr = orderA
	orderQuery = orderQ
	receiptAggr = receiptA
	receiptQuery = receiptQS
	paymentCtrl = paymentManager
	vendorAggr = vendorA
	vendorQuery = vendorQ
	carrierAggr = carrierA
	carrierQuery = carrierQ
	traderQuery = traderQ
	eventBus = eventB
	summaryQuery = summary
	sd.Register(idempgroup.Shutdown)
	inventoryAggregate = inventoryA
	inventoryQuery = inventoryQ
}

type Service struct{}

func (s *Service) VersionInfo(ctx context.Context, q *wrapshop.VersionInfoEndpoint) error {
	q.Result = &pbcm.VersionInfoResponse{
		Service: "etop.Shop",
		Version: "0.1",
	}
	return nil
}

func (s *Service) CreateInventoryVoucher(ctx context.Context, q *wrapshop.CreateInventoryVoucherEndpoint) error {
	shopID := q.Context.Shop.ID
	userID := q.Context.UserID
	inventoryOverstock := q.Context.Shop.InventoryOverstock
	var items []*inventory.InventoryVoucherItem
	for _, value := range q.Lines {
		items = append(items, &inventory.InventoryVoucherItem{
			VariantID: value.VariantId,
			Price:     value.Price,
			Quantity:  value.Quantity,
		})
	}
	cmd := &inventory.CreateInventoryVoucherCommand{
		Title:       q.Title,
		ShopID:      shopID,
		Overstock:   cm.BoolDefault(inventoryOverstock, true),
		TotalAmount: q.TotalAmount,
		CreatedBy:   userID,
		TraderID:    q.TraderId,
		Type:        q.Type,
		Note:        q.Note,
		Lines:       items,
	}
	if err := inventoryAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbshop.CreateInventoryVoucherResponse{
		Inventory: PbShopInventoryVoucher(cmd.Result),
	}
	return nil
}

func (s *Service) ConfirmInventoryVoucher(ctx context.Context, q *wrapshop.ConfirmInventoryVoucherEndpoint) error {
	shopID := q.Context.Shop.ID
	userID := q.Context.UserID

	cmd := &inventory.ConfirmInventoryVoucherCommand{
		ShopID:    shopID,
		ID:        q.Id,
		UpdatedBy: userID,
	}
	if err := inventoryAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbshop.ConfirmInventoryVoucherResponse{
		Inventory: PbShopInventoryVoucher(cmd.Result),
	}
	return nil
}

func (s *Service) CancelInventoryVoucher(ctx context.Context, q *wrapshop.CancelInventoryVoucherEndpoint) error {
	shopID := q.Context.Shop.ID
	userID := q.Context.UserID

	cmd := &inventory.CancelInventoryVoucherCommand{
		ShopID:    shopID,
		ID:        q.Id,
		UpdatedBy: userID,
		Reason:    q.Reason,
	}
	if err := inventoryAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbshop.CancelInventoryVoucherResponse{
		Inventory: PbShopInventoryVoucher(cmd.Result),
	}
	return nil
}

func (s *Service) UpdateInventoryVoucher(ctx context.Context, q *wrapshop.UpdateInventoryVoucherEndpoint) error {
	shopID := q.Context.Shop.ID
	userID := q.Context.UserID
	var items []*inventory.InventoryVoucherItem
	for _, value := range q.Lines {
		items = append(items, &inventory.InventoryVoucherItem{
			VariantID: value.VariantId,
			Price:     value.Price,
			Quantity:  value.Quantity,
		})
	}
	cmd := &inventory.UpdateInventoryVoucherCommand{
		Title:       q.Title,
		ID:          q.Id,
		ShopID:      shopID,
		TotalAmount: q.TotalAmount,
		UpdatedBy:   userID,
		TraderID:    q.TraderId,
		Note:        q.Note,
		Lines:       items,
	}
	if err := inventoryAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbshop.UpdateInventoryVoucherResponse{
		Inventory: PbShopInventoryVoucher(cmd.Result),
	}
	return nil
}

func (s *Service) AdjustInventoryQuantity(ctx context.Context, q *wrapshop.AdjustInventoryQuantityEndpoint) error {
	shopID := q.Context.Shop.ID
	userID := q.Context.UserID
	inventoryOverstock := q.Context.Shop.InventoryOverstock
	var items []*inventory.InventoryVariant
	for _, value := range q.Inventories {
		items = append(items, &inventory.InventoryVariant{
			ShopID:          shopID,
			VariantID:       value.VariantId,
			QuantityOnHand:  value.QuantityOnHand,
			QuantitySummary: value.Quantity,
			QuantityPicked:  value.QuantityPicked,
		})
	}
	cmd := &inventory.AdjustInventoryQuantityCommand{
		Overstock: cm.BoolDefault(inventoryOverstock, true),
		ShopID:    shopID,
		Lines:     items,
		UserID:    userID,
		Note:      q.Note,
	}
	if err := inventoryAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbshop.AdjustInventoryQuantityResponse{
		Inventories:       PbInventories(cmd.Result.Inventory),
		InventoryVouchers: PbShopInventoryVouchers(cmd.Result.InventoryVouchers),
	}
	return nil
}

func (s *Service) GetInventories(ctx context.Context, q *wrapshop.GetInventoriesEndpoint) error {
	shopID := q.Context.Shop.ID
	query := &inventory.GetInventoriesQuery{
		ShopID: shopID,
		Paging: &meta.Paging{
			Offset: q.Paging.Offset,
			Limit:  q.Paging.Limit,
		},
	}
	if err := inventoryQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pbshop.GetInventoriesResponse{
		Inventories: PbInventories(query.Result.Inventories),
	}
	return nil
}

func (s *Service) GetInventory(ctx context.Context, q *wrapshop.GetInventoryEndpoint) error {
	shopID := q.Context.Shop.ID
	query := &inventory.GetInventoryQuery{
		ShopID:    shopID,
		VariantID: q.VariantId,
	}
	if err := inventoryQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = PbInventory(query.Result)
	return nil
}

func (s *Service) GetInventoryByVariantIDs(ctx context.Context, q *wrapshop.GetInventoriesByVariantIDsEndpoint) error {
	shopID := q.Context.Shop.ID
	query := &inventory.GetInventoriesByVariantIDsQuery{
		ShopID:     shopID,
		VariantIDs: q.VariantIds,
	}
	if err := inventoryQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pbshop.GetInventoriesResponse{
		Inventories: PbInventories(query.Result.Inventories),
	}
	return nil
}

func (s *Service) GetInventoryVoucher(ctx context.Context, q *wrapshop.GetInventoryVoucherEndpoint) error {
	shopID := q.Context.Shop.ID
	query := &inventory.GetInventoryVoucherQuery{
		ShopID: shopID,
		ID:     q.Id,
	}
	if err := inventoryQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = PbShopInventoryVoucher(query.Result)
	return nil
}

func (s *Service) GetInventoryVouchers(ctx context.Context, q *wrapshop.GetInventoryVouchersEndpoint) error {
	shopID := q.Context.Shop.ID
	query := &inventory.GetInventoryVouchersQuery{
		ShopID: shopID,
		Paging: &meta.Paging{
			Offset: q.Paging.Offset,
			Limit:  q.Paging.Limit,
		},
	}
	if err := inventoryQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pbshop.GetInventoryVouchersResponse{
		InventoryVouchers: PbShopInventoryVouchers(query.Result.InventoryVoucher),
	}
	return nil
}

func (s *Service) GetInventoryVoucherByIDs(ctx context.Context, q *wrapshop.GetInventoryVouchersByIDsEndpoint) error {
	shopID := q.Context.Shop.ID
	query := &inventory.GetInventoryVouchersByIDsQuery{
		ShopID: shopID,
		IDs:    q.Ids,
	}
	if err := inventoryQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pbshop.GetInventoryVouchersResponse{
		InventoryVouchers: PbShopInventoryVouchers(query.Result.InventoryVoucher),
	}
	return nil
}

func (s *Service) UpdateVariant(ctx context.Context, q *wrapshop.UpdateVariantEndpoint) error {
	shopID := q.Context.Shop.ID
	cmd := &catalog.UpdateShopVariantInfoCommand{
		ShopID:    shopID,
		VariantID: q.Id,
		Name:      PString(q.Name),
		Code:      PString(q.Code),
		Note:      PString(q.Note),

		ShortDesc:    PString(q.ShortDesc),
		Descripttion: PString(q.Description),
		DescHTML:     PString(q.DescHtml),

		CostPrice:   PInt32(q.CostPrice),
		ListPrice:   PInt32(q.ListPrice),
		RetailPrice: PInt32(q.RetailPrice),
	}
	if err := catalogAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = PbShopVariant(cmd.Result)
	return nil
}

func (s *Service) UpdateVariantAttributes(ctx context.Context, q *wrapshop.UpdateVariantAttributesEndpoint) error {
	shopID := q.Context.Shop.ID

	var attributes catalog.Attributes
	for _, value := range q.Attributes {
		attributes = append(attributes, &catalog.Attribute{
			Name:  value.Name,
			Value: value.Value,
		})
	}
	cmd := &catalog.UpdateShopVariantAttributesCommand{
		ShopID:     shopID,
		VariantID:  q.VariantId,
		Attributes: attributes,
	}
	if err := catalogAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = PbShopVariant(cmd.Result)
	return nil
}

func (s *Service) RemoveVariants(ctx context.Context, q *wrapshop.RemoveVariantsEndpoint) error {
	cmd := &catalog.DeleteShopVariantsCommand{
		ShopID: q.Context.Shop.ID,
		IDs:    q.Ids,
	}
	if err := catalogAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.RemovedResponse{
		Removed: int32(cmd.Result),
	}
	return nil
}

func (s *Service) GetProduct(ctx context.Context, q *wrapshop.GetProductEndpoint) error {
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

func (s *Service) GetProductsByIDs(ctx context.Context, q *wrapshop.GetProductsByIDsEndpoint) error {
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

func (s *Service) GetProducts(ctx context.Context, q *wrapshop.GetProductsEndpoint) error {
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
		Paging:   pbcm.PbPaging(cm.Paging(query.Result.Paging), query.Result.Count),
		Products: PbShopProductsWithVariants(query.Result.Products),
	}
	return nil
}

func (s *Service) CreateProduct(ctx context.Context, q *wrapshop.CreateProductEndpoint) error {
	metaFields := []*catalog.MetaField{}

	for _, metaField := range q.MetaFields {
		metaFields = append(metaFields, &catalog.MetaField{
			Key:   metaField.Key,
			Value: metaField.Value,
		})
	}
	cmd := &catalog.CreateShopProductCommand{
		ShopID:    q.Context.Shop.ID,
		Code:      q.Code,
		Name:      q.Name,
		Unit:      q.Unit,
		ImageURLs: q.ImageUrls,
		Note:      q.Note,
		DescriptionInfo: catalog.DescriptionInfo{
			ShortDesc:   q.ShortDesc,
			Description: q.Description,
			DescHTML:    q.DescHtml,
		},
		PriceInfo: catalog.PriceInfo{
			CostPrice:   q.CostPrice,
			ListPrice:   q.ListPrice,
			RetailPrice: q.RetailPrice,
		},
		VendorID:    q.VendorId,
		ProductType: q.ProductType.ToProductType(),
		MetaFields:  metaFields,
	}
	if err := catalogAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = PbShopProductWithVariants(cmd.Result)
	return nil
}

func (s *Service) RemoveProducts(ctx context.Context, q *wrapshop.RemoveProductsEndpoint) error {
	cmd := &catalog.DeleteShopProductsCommand{
		ShopID: q.Context.Shop.ID,
		IDs:    q.Ids,
	}
	if err := catalogAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.RemovedResponse{
		Removed: int32(cmd.Result),
	}
	return nil
}

func (s *Service) UpdateProduct(ctx context.Context, q *wrapshop.UpdateProductEndpoint) error {
	shopID := q.Context.Shop.ID
	cmd := &catalog.UpdateShopProductInfoCommand{
		ShopID:    shopID,
		ProductID: q.Id,
		VendorID:  q.VendorId,
		Code:      PString(q.Code),
		Name:      PString(q.Name),
		Unit:      PString(q.Unit),
		Note:      PString(q.Note),

		ShortDesc:   PString(q.ShortDesc),
		Description: PString(q.Description),
		DescHTML:    PString(q.DescHtml),

		CostPrice:   PInt32(q.CostPrice),
		ListPrice:   PInt32(q.ListPrice),
		RetailPrice: PInt32(q.RetailPrice),
	}
	if err := catalogAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = PbShopProductWithVariants(cmd.Result)
	return nil
}

func (s *Service) UpdateProductsStatus(ctx context.Context, q *wrapshop.UpdateProductsStatusEndpoint) error {
	shopID := q.Context.Shop.ID
	cmd := &catalog.UpdateShopProductStatusCommand{
		IDs:    q.Ids,
		ShopID: shopID,
		Status: int16(q.Status),
	}
	if err := catalogAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbshop.UpdateProductStatusResponse{Updated: int32(cmd.Result)}
	return nil
}

func (s *Service) UpdateVariantsStatus(ctx context.Context, q *wrapshop.UpdateVariantsStatusEndpoint) error {
	shopID := q.Context.Shop.ID
	cmd := &catalog.UpdateShopVariantStatusCommand{
		IDs:    q.Ids,
		ShopID: shopID,
		Status: int16(q.Status),
	}
	if err := catalogAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbshop.UpdateProductStatusResponse{Updated: int32(cmd.Result)}
	return nil
}

func (s *Service) UpdateProductsTags(ctx context.Context, q *wrapshop.UpdateProductsTagsEndpoint) error {
	shopID := q.Context.Shop.ID
	cmd := &catalogmodelx.UpdateShopProductsTagsCommand{
		ShopID:     shopID,
		ProductIDs: q.Ids,
		Update: &model.UpdateListRequest{
			Adds:       q.Adds,
			Deletes:    q.Deletes,
			ReplaceAll: q.ReplaceAll,
			DeleteAll:  q.DeleteAll,
		},
	}

	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.UpdatedResponse{
		Updated: int32(cmd.Result.Updated),
	}
	return nil
}

func (s *Service) GetVariant(ctx context.Context, q *wrapshop.GetVariantEndpoint) error {
	return cm.ErrTODO
}

func (s *Service) GetVariantsByIDs(ctx context.Context, q *wrapshop.GetVariantsByIDsEndpoint) error {
	return cm.ErrTODO
}

func (s *Service) CreateVariant(ctx context.Context, q *wrapshop.CreateVariantEndpoint) error {
	cmd := &catalog.CreateShopVariantCommand{
		ShopID:     q.Context.Shop.ID,
		ProductID:  q.ProductId,
		Code:       q.Code,
		Name:       q.Name,
		ImageURLs:  q.ImageUrls,
		Note:       q.Note,
		Attributes: convertpb.PbAttributesToDomain(q.Attributes),
		DescriptionInfo: catalog.DescriptionInfo{
			ShortDesc:   q.ShortDesc,
			Description: q.Description,
			DescHTML:    q.DescHtml,
		},
		PriceInfo: catalog.PriceInfo{
			CostPrice:   q.CostPrice,
			ListPrice:   q.ListPrice,
			RetailPrice: q.RetailPrice,
		},
	}
	if err := catalogAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = PbShopVariant(cmd.Result)
	return nil
}

func (s *Service) DeprecatedCreateVariant(ctx context.Context, q *wrapshop.DeprecatedCreateVariantEndpoint) error {
	cmd := &catalogmodelx.DeprecatedCreateVariantCommand{
		ShopID:      q.Context.Shop.ID,
		ProductID:   q.ProductId,
		ProductName: q.ProductName,
		Name:        q.Name,
		Description: q.Description,
		ShortDesc:   q.ShortDesc,
		ImageURLs:   q.ImageUrls,
		Tags:        q.Tags,
		Status:      *q.Status.ToModel(),

		CostPrice:   q.CostPrice,
		ListPrice:   q.ListPrice,
		RetailPrice: q.RetailPrice,

		ProductCode:       q.Code,
		VariantCode:       q.Code,
		QuantityAvailable: int(q.QuantityAvailable),
		QuantityOnHand:    int(q.QuantityOnHand),
		QuantityReserved:  int(q.QuantityReserved),

		Attributes: convertpb.AttributesTomodel(q.Attributes),
		DescHTML:   q.DescHtml,
	}

	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}

	q.Result = PbShopProductWithVariants(cmd.Result)
	return nil
}

func (s *Service) CreateProductSourceCategory(ctx context.Context, q *wrapshop.CreateProductSourceCategoryEndpoint) error {
	cmd := &catalogmodelx.CreateShopCategoryCommand{
		ShopID:   q.Context.Shop.ID,
		Name:     q.Name,
		ParentID: q.ParentId,
	}

	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = convertpb.PbCategory(cmd.Result)
	return nil
}

func (s *Service) UpdateProductsPSCategory(ctx context.Context, q *wrapshop.UpdateProductsPSCategoryEndpoint) error {
	cmd := &catalogmodelx.UpdateProductsShopCategoryCommand{
		CategoryID: q.CategoryId,
		ProductIDs: q.ProductIds,
		ShopID:     q.Context.Shop.ID,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.UpdatedResponse{
		Updated: int32(cmd.Result.Updated),
	}
	return nil
}

func (s *Service) GetProductSourceCategory(ctx context.Context, q *wrapshop.GetProductSourceCategoryEndpoint) error {
	cmd := &catalogmodelx.GetShopCategoryQuery{
		ShopID:     q.Context.Shop.ID,
		CategoryID: q.Id,
	}

	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}

	q.Result = convertpb.PbCategory(cmd.Result)
	return nil
}

func (s *Service) GetProductSourceCategories(ctx context.Context, q *wrapshop.GetProductSourceCategoriesEndpoint) error {
	cmd := &catalogmodelx.GetProductSourceCategoriesQuery{
		ShopID: q.Context.Shop.ID,
	}

	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}

	q.Result = &pbshop.CategoriesResponse{
		Categories: convertpb.PbCategories(cmd.Result.Categories),
	}
	return nil
}

func (s *Service) UpdateProductSourceCategory(ctx context.Context, q *wrapshop.UpdateProductSourceCategoryEndpoint) error {
	cmd := &catalogmodelx.UpdateShopCategoryCommand{
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

func (s *Service) RemoveProductSourceCategory(ctx context.Context, q *wrapshop.RemoveProductSourceCategoryEndpoint) error {
	cmd := &catalogmodelx.RemoveShopCategoryCommand{
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

func (s *Service) UpdateProductImages(ctx context.Context, q *wrapshop.UpdateProductImagesEndpoint) error {
	shopID := q.Context.Shop.ID

	var metaUpdate []*meta.UpdateSet
	if q.DeleteAll == true {
		metaUpdate = append(metaUpdate, &meta.UpdateSet{
			Op: meta.OpDeleteAll,
		})
	}
	if q.ReplaceAll != nil {
		metaUpdate = append(metaUpdate, &meta.UpdateSet{
			Op:      meta.OpReplaceAll,
			Changes: q.ReplaceAll,
		})
	}
	metaUpdate = append(metaUpdate, &meta.UpdateSet{
		Op:      meta.OpAdd,
		Changes: q.Adds,
	})
	metaUpdate = append(metaUpdate, &meta.UpdateSet{
		Op:      meta.OpRemove,
		Changes: q.Deletes,
	})

	cmd := catalog.UpdateShopProductImagesCommand{
		ShopID:  shopID,
		ID:      q.Id,
		Updates: metaUpdate,
	}

	if err := catalogAggr.Dispatch(ctx, &cmd); err != nil {
		return err
	}
	q.Result = PbShopProductWithVariants(cmd.Result)
	return nil
}

func (s *Service) UpdateProductMetaFields(ctx context.Context, q *wrapshop.UpdateProductMetaFieldsEndpoint) error {
	metaFields := []*catalog.MetaField{}
	for _, metaField := range q.MetaFields {
		metaFields = append(metaFields, &catalog.MetaField{
			Key:   metaField.Key,
			Value: metaField.Value,
		})
	}
	cmd := catalog.UpdateShopProductMetaFieldsCommand{
		ID:         q.Id,
		ShopID:     q.Context.Shop.ID,
		MetaFields: metaFields,
	}
	if err := catalogAggr.Dispatch(ctx, &cmd); err != nil {
		return err
	}
	q.Result = PbShopProductWithVariants(cmd.Result)
	return nil
}

func (s *Service) UpdateVariantImages(ctx context.Context, q *wrapshop.UpdateVariantImagesEndpoint) error {
	shopID := q.Context.Shop.ID

	var metaUpdate []*meta.UpdateSet
	if q.DeleteAll == true {
		metaUpdate = append(metaUpdate, &meta.UpdateSet{
			Op: meta.OpDeleteAll,
		})
	}
	if q.ReplaceAll != nil {
		metaUpdate = append(metaUpdate, &meta.UpdateSet{
			Op:      meta.OpReplaceAll,
			Changes: q.ReplaceAll,
		})
	}
	metaUpdate = append(metaUpdate, &meta.UpdateSet{
		Op:      meta.OpAdd,
		Changes: q.Adds,
	})
	metaUpdate = append(metaUpdate, &meta.UpdateSet{
		Op:      meta.OpRemove,
		Changes: q.Deletes,
	})

	cmd := catalog.UpdateShopVariantImagesCommand{
		ShopID:  shopID,
		ID:      q.Id,
		Updates: metaUpdate,
	}
	if err := catalogAggr.Dispatch(ctx, &cmd); err != nil {
		return err
	}
	q.Result = PbShopVariant(cmd.Result)
	return nil
}

func (s *Service) GetMoneyTransaction(ctx context.Context, q *wrapshop.GetMoneyTransactionEndpoint) error {
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

func (s *Service) GetMoneyTransactions(ctx context.Context, q *wrapshop.GetMoneyTransactionsEndpoint) error {
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

func (s *Service) SummarizeFulfillments(ctx context.Context, q *wrapshop.SummarizeFulfillmentsEndpoint) error {
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

func (s *Service) GetSummarizePOS(ctx context.Context, q *wrapshop.SummarizePOSEndpoint) error {
	dateFrom, dateTo, err := cm.ParseDateFromTo(q.DateFrom, q.DateTo)
	if err != nil {
		return err
	}
	query := &summary.SummaryPOSQuery{
		ShopID:   q.Context.Shop.ID,
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	if err := summaryQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pbshop.SummarizePOSResponse{
		Tables: pbshop.PbSummaryTablesNew(query.Result.ListTable),
	}
	return nil
}

func (s *Service) CalcBalance(ctx context.Context, q *wrapshop.CalcBalanceShopEndpoint) error {
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

func (s *Service) CreateDevice(ctx context.Context, q *wrapshop.CreateDeviceEndpoint) error {
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

func (s *Service) DeleteDevice(ctx context.Context, q *wrapshop.DeleteDeviceEndpoint) error {
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

func (s *Service) GetNotification(ctx context.Context, q *wrapshop.GetNotificationEndpoint) error {
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

func (s *Service) GetNotifications(ctx context.Context, q *wrapshop.GetNotificationsEndpoint) error {
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

func (s *Service) UpdateNotifications(ctx context.Context, q *wrapshop.UpdateNotificationsEndpoint) error {
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

func (s *Service) GetShipnowFulfillment(ctx context.Context, q *wrapshop.GetShipnowFulfillmentEndpoint) error {
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

func (s *Service) GetShipnowFulfillments(ctx context.Context, q *wrapshop.GetShipnowFulfillmentsEndpoint) error {
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

func (s *Service) CreateShipnowFulfillment(ctx context.Context, q *wrapshop.CreateShipnowFulfillmentEndpoint) error {
	pickupAddress, err := q.PickupAddress.Fulfilled()
	if err != nil {
		return err
	}
	cmd := &shipnow.CreateShipnowFulfillmentCommand{
		OrderIds:            q.OrderIds,
		Carrier:             carriertypes.CarrierFromString(q.Carrier),
		ShopId:              q.Context.Shop.ID,
		ShippingServiceCode: q.ShippingServiceCode,
		ShippingServiceFee:  q.ShippingServiceFee,
		ShippingNote:        q.ShippingNote,
		RequestPickupAt:     time.Time{},
		PickupAddress:       pborder.Convert_api_OrderAddress_To_core_OrderAddress(pickupAddress),
	}
	if err := shipnowAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = pborder.Convert_core_ShipnowFulfillment_To_api_ShipnowFulfillment(cmd.Result)
	return nil
}

func (s *Service) ConfirmShipnowFulfillment(ctx context.Context, q *wrapshop.ConfirmShipnowFulfillmentEndpoint) error {
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

func (s *Service) UpdateShipnowFulfillment(ctx context.Context, q *wrapshop.UpdateShipnowFulfillmentEndpoint) error {
	pickupAddress, err := q.PickupAddress.Fulfilled()
	if err != nil {
		return err
	}
	cmd := &shipnow.UpdateShipnowFulfillmentCommand{
		Id:                  q.Id,
		OrderIds:            q.OrderIds,
		Carrier:             carriertypes.CarrierFromString(q.Carrier),
		ShopId:              q.Context.Shop.ID,
		ShippingServiceCode: q.ShippingServiceCode,
		ShippingServiceFee:  q.ShippingServiceFee,
		ShippingNote:        q.ShippingNote,
		RequestPickupAt:     time.Time{},
		PickupAddress:       pborder.Convert_api_OrderAddress_To_core_OrderAddress(pickupAddress),
	}
	if err := shipnowAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = pborder.Convert_core_ShipnowFulfillment_To_api_ShipnowFulfillment(cmd.Result)
	return nil
}

func (s *Service) CancelShipnowFulfillment(ctx context.Context, q *wrapshop.CancelShipnowFulfillmentEndpoint) error {
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

func (s *Service) GetShipnowServices(ctx context.Context, q *wrapshop.GetShipnowServicesEndpoint) error {
	pickupAddress, err := q.PickupAddress.Fulfilled()
	if err != nil {
		return err
	}
	var points []*shipnow.DeliveryPoint
	if len(q.DeliveryPoints) > 0 {
		for _, p := range q.DeliveryPoints {
			addr, err := p.ShippingAddress.Fulfilled()
			if err != nil {
				return err
			}
			points = append(points, &shipnow.DeliveryPoint{
				ShippingAddress: pborder.Convert_api_OrderAddress_To_core_OrderAddress(addr),
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

func (s *Service) CreateExternalAccountAhamove(ctx context.Context, q *wrapshop.CreateExternalAccountAhamoveEndpoint) error {
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
	addr := queryAddress.Result
	cmd := &identity.CreateExternalAccountAhamoveCommand{
		OwnerID: user.ID,
		Phone:   phone,
		Name:    user.FullName,
		Address: addr.GetFullAddress(),
	}
	if err := identityAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = pbshop.Convert_core_XAccountAhamove_To_api_XAccountAhamove(cmd.Result)
	return nil
}

func (s *Service) GetExternalAccountAhamove(ctx context.Context, q *wrapshop.GetExternalAccountAhamoveEndpoint) error {
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

	q.Result = pbshop.Convert_core_XAccountAhamove_To_api_XAccountAhamove(account)
	return nil
}

func (s *Service) RequestVerifyExternalAccountAhamove(ctx context.Context, q *wrapshop.RequestVerifyExternalAccountAhamoveEndpoint) error {
	query := &model.GetUserByIDQuery{
		UserID: q.Context.Shop.OwnerID,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	user := query.Result
	phone := user.Phone

	cmd := &identity.RequestVerifyExternalAccountAhamoveCommand{
		OwnerID: user.ID,
		Phone:   phone,
	}
	if err := identityAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}

	q.Result = &pbcm.UpdatedResponse{
		Updated: 1,
	}
	return nil
}

func (s *Service) UpdateExternalAccountAhamoveVerification(ctx context.Context, r *wrapshop.UpdateExternalAccountAhamoveVerificationEndpoint) error {
	if err := validateUrl(r.IdCardFrontImg, r.IdCardBackImg, r.PortraitImg, r.WebsiteUrl, r.FanpageUrl); err != nil {
		return err
	}
	if err := validateUrl(r.BusinessLicenseImgs...); err != nil {
		return err
	}
	if err := validateUrl(r.CompanyImgs...); err != nil {
		return err
	}

	query := &model.GetUserByIDQuery{
		UserID: r.Context.Shop.OwnerID,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	user := query.Result
	phone := user.Phone

	cmd := &identity.UpdateExternalAccountAhamoveVerificationCommand{
		OwnerID:             user.ID,
		Phone:               phone,
		IDCardFrontImg:      r.IdCardFrontImg,
		IDCardBackImg:       r.IdCardBackImg,
		PortraitImg:         r.PortraitImg,
		WebsiteURL:          r.WebsiteUrl,
		FanpageURL:          r.FanpageUrl,
		CompanyImgs:         r.CompanyImgs,
		BusinessLicenseImgs: r.BusinessLicenseImgs,
	}
	if err := identityAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}

	r.Result = &pbcm.UpdatedResponse{
		Updated: 1,
	}
	return nil
}

func validateUrl(imgsUrl ...string) error {
	for _, url := range imgsUrl {
		if url == "" {
			continue
		}
		if !govalidator.IsURL(url) {
			return cm.Errorf(cm.InvalidArgument, nil, "Invalid url: %v", url)
		}
	}
	return nil
}

func (s *Service) GetExternalAccountHaravan(ctx context.Context, r *wrapshop.GetExternalAccountHaravanEndpoint) error {
	query := &haravanidentity.GetExternalAccountHaravanByShopIDQuery{
		ShopID: r.Context.Shop.ID,
	}
	if err := haravanIdentityQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = pbshop.Convert_core_XAccountHaravan_To_api_XAccountHaravan(query.Result)
	return nil
}

func (s *Service) CreateExternalAccountHaravan(ctx context.Context, r *wrapshop.CreateExternalAccountHaravanEndpoint) error {
	cmd := &haravanidentity.CreateExternalAccountHaravanCommand{
		ShopID:      r.Context.Shop.ID,
		Subdomain:   r.Subdomain,
		Code:        r.Code,
		RedirectURI: r.RedirectUri,
	}
	if err := haravanIdentityAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = pbshop.Convert_core_XAccountHaravan_To_api_XAccountHaravan(cmd.Result)
	return nil
}

func (s *Service) UpdateExternalAccountHaravanToken(ctx context.Context, r *wrapshop.UpdateExternalAccountHaravanTokenEndpoint) error {
	cmd := &haravanidentity.UpdateExternalAccountHaravanTokenCommand{
		ShopID:      r.Context.Shop.ID,
		Subdomain:   r.Subdomain,
		RedirectURI: r.RedirectUri,
		Code:        r.Code,
	}
	if err := haravanIdentityAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = pbshop.Convert_core_XAccountHaravan_To_api_XAccountHaravan(cmd.Result)
	return nil
}

func (s *Service) ConnectCarrierServiceExternalAccountHaravan(ctx context.Context, r *wrapshop.ConnectCarrierServiceExternalAccountHaravanEndpoint) error {
	cmd := &haravanidentity.ConnectCarrierServiceExternalAccountHaravanCommand{
		ShopID: r.Context.Shop.ID,
	}
	if err := haravanIdentityAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.UpdatedResponse{
		Updated: 1,
	}
	return nil
}

func (s *Service) DeleteConnectedCarrierServiceExternalAccountHaravan(ctx context.Context, r *wrapshop.DeleteConnectedCarrierServiceExternalAccountHaravanEndpoint) error {
	cmd := &haravanidentity.DeleteConnectedCarrierServiceExternalAccountHaravanCommand{
		ShopID: r.Context.Shop.ID,
	}
	if err := haravanIdentityAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.DeletedResponse{
		Deleted: 1,
	}
	return nil
}

func (s *Service) PaymentTradingOrder(ctx context.Context, q *wrapshop.PaymentTradingOrderEndpoint) error {
	if q.OrderId == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing OrderID")
	}
	if q.ReturnUrl == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing ReturnURL")
	}

	argGenCode := &paymentmanager.GenerateCodeCommand{
		PaymentSource: payment.PaymentSourceOrder,
		ID:            strconv.FormatInt(q.OrderId, 10),
	}
	if err := paymentCtrl.Dispatch(ctx, argGenCode); err != nil {
		return err
	}
	args := &paymentmanager.BuildUrlConnectPaymentGatewayCommand{
		OrderID:           argGenCode.Result,
		Desc:              q.Desc,
		ReturnURL:         q.ReturnUrl,
		TransactionAmount: int(q.Amount),
		Provider:          q.PaymentProvider.ToPaymentProvider(),
	}

	if err := paymentCtrl.Dispatch(ctx, args); err != nil {
		return err
	}
	q.Result = &pbshop.PaymentTradingOrderResponse{
		Url: args.Result,
	}
	return nil
}

func (s *Service) PaymentCheckReturnData(ctx context.Context, q *wrapshop.PaymentCheckReturnDataEndpoint) error {
	if q.Id == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Mã giao dịch không được để trống")
	}
	if q.Code == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Mã 'Code' không được để trống")
	}
	args := &paymentmanager.CheckReturnDataCommand{
		ID:                    q.Id,
		Code:                  q.Code,
		PaymentStatus:         q.PaymentStatus,
		Amount:                int(q.Amount),
		ExternalTransactionID: q.ExternalTransactionId,
		Provider:              q.PaymentProvider.ToPaymentProvider(),
	}
	if err := paymentCtrl.Dispatch(ctx, args); err != nil {
		return err
	}
	q.Result = &pbcm.MessageResponse{
		Code: "ok",
		Msg:  args.Result.Msg,
	}
	return nil
}

func (s *Service) CreateCategory(ctx context.Context, q *wrapshop.CreateCategoryEndpoint) error {
	cmd := &catalog.CreateShopCategoryCommand{
		ShopID:   q.Context.Shop.ID,
		Name:     q.Name,
		ParentID: q.ParentId,
	}
	if err := catalogAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = PbShopCategory(cmd.Result)
	return nil
}

func (s *Service) GetCategory(ctx context.Context, q *wrapshop.GetCategoryEndpoint) error {
	query := &catalog.GetShopCategoryQuery{
		ID:     q.Id,
		ShopID: q.Context.Shop.ID,
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = PbShopCategory(query.Result)
	return nil
}

func (s *Service) GetCategories(ctx context.Context, q *wrapshop.GetCategoriesEndpoint) error {
	paging := q.Paging.CMPaging()
	query := &catalog.ListShopCategoriesQuery{
		ShopID:  q.Context.Shop.ID,
		Paging:  *paging,
		Filters: pbcm.ToFilters(q.Filters),
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = &pbshop.ShopCategoriesResponse{
		Paging:     pbcm.PbPageInfo(paging, query.Result.Count),
		Categories: PbShopCategories(query.Result.Categories),
	}
	return nil
}

func (s *Service) UpdateCategory(ctx context.Context, q *wrapshop.UpdateCategoryEndpoint) error {
	shopID := q.Context.Shop.ID
	cmd := &catalog.UpdateShopCategoryCommand{
		ID:       q.Id,
		ShopID:   shopID,
		Name:     PString(q.Name),
		ParentID: q.ParentId,
	}
	if err := catalogAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = PbShopCategory(cmd.Result)
	return nil
}

func (s *Service) DeleteCategory(ctx context.Context, r *wrapshop.DeleteCategoryEndpoint) error {
	cmd := &catalog.DeleteShopCategoryCommand{
		ID:     r.Id,
		ShopID: r.Context.Shop.ID,
	}
	if err := catalogAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.DeletedResponse{Deleted: int32(cmd.Result)}
	return nil
}

func (s *Service) UpdateProductCategory(ctx context.Context, q *wrapshop.UpdateProductCategoryEndpoint) error {
	shopID := q.Context.Shop.ID
	cmd := &catalog.UpdateShopProductCategoryCommand{
		ProductID:  q.ProductId,
		CategoryID: q.CategoryId,
		ShopID:     shopID,
	}
	if err := catalogAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = PbShopProductWithVariants(cmd.Result)
	return nil
}

func (s *Service) GetCollection(ctx context.Context, q *wrapshop.GetCollectionEndpoint) error {
	query := &catalog.GetShopCollectionQuery{
		ID:     q.Id,
		ShopID: q.Context.Shop.ID,
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = PbShopCollection(query.Result)
	return nil
}

func (s *Service) GetCollections(ctx context.Context, q *wrapshop.GetCollectionsEndpoint) error {
	paging := q.Paging.CMPaging()
	query := &catalog.ListShopCollectionsQuery{
		ShopID:  q.Context.Shop.ID,
		Paging:  *paging,
		Filters: pbcm.ToFilters(q.Filters),
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pbshop.ShopCollectionsResponse{
		Paging:      pbcm.PbPageInfo(paging, query.Result.Count),
		Collections: PbShopCollections(query.Result.Collections),
	}
	return nil
}

func (s *Service) UpdateCollection(ctx context.Context, q *wrapshop.UpdateCollectionEndpoint) error {
	shopID := q.Context.Shop.ID
	cmd := &catalog.UpdateShopCollectionCommand{
		ID:          q.Id,
		ShopID:      shopID,
		Name:        PString(q.Name),
		Description: PString(q.Description),
		DescHTML:    PString(q.DescHtml),
		ShortDesc:   PString(q.ShortDesc),
	}
	if err := catalogAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = PbShopCollection(cmd.Result)
	return nil
}

func (s *Service) CreateCollection(ctx context.Context, q *wrapshop.CreateCollectionEndpoint) error {
	cmd := &catalog.CreateShopCollectionCommand{
		ShopID:      q.Context.Shop.ID,
		Name:        q.Name,
		DescHTML:    q.DescHtml,
		Description: q.Description,
		ShortDesc:   q.ShortDesc,
	}
	if err := catalogAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = PbShopCollection(cmd.Result)
	return nil
}

func (s *Service) AddProductCollections(ctx context.Context, r *wrapshop.AddProductCollectionEndpoint) error {
	cmd := &catalog.AddShopProductCollectionCommand{
		ProductID:     r.ProductId,
		CollectionIDs: r.CollectionIds,
		ShopID:        r.Context.Shop.ID,
	}
	if err := catalogAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.UpdatedResponse{Updated: int32(cmd.Result)}
	return nil
}

func (s *Service) RemoveProductCollection(ctx context.Context, r *wrapshop.RemoveProductCollectionEndpoint) error {
	cmd := &catalog.RemoveShopProductCollectionCommand{
		ProductID:     r.ProductId,
		CollectionIDs: r.CollectionIds,
		ShopID:        r.Context.Shop.ID,
	}
	if err := catalogAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.RemovedResponse{Removed: int32(cmd.Result)}
	return nil
}

func (s *Service) GetCollectionsByProductID(ctx context.Context, q *wrapshop.GetCollectionsByProductIDEndpoint) error {
	query := &catalog.ListShopCollectionsByProductIDQuery{
		ShopID:    q.Context.Shop.ID,
		ProductID: q.ProductId,
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pbshop.CollectionsResponse{
		Collections: PbShopCollections(query.Result),
	}
	return nil
}

func (s *Service) RemoveProductCategory(ctx context.Context, r *wrapshop.RemoveProductCategoryEndpoint) error {
	cmd := &catalog.RemoveShopProductCategoryCommand{
		ShopID:    r.Context.Shop.ID,
		ProductID: r.Id,
	}
	if err := catalogAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = PbShopProductWithVariants(cmd.Result)
	return nil
}
