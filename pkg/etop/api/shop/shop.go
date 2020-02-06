package shop

import (
	"context"
	"time"

	"github.com/asaskevich/govalidator"

	haravanidentity "etop.vn/api/external/haravan/identity"
	paymentmanager "etop.vn/api/external/payment/manager"
	"etop.vn/api/main/address"
	"etop.vn/api/main/authorization"
	"etop.vn/api/main/catalog"
	"etop.vn/api/main/connectioning"
	"etop.vn/api/main/identity"
	"etop.vn/api/main/inventory"
	"etop.vn/api/main/ledgering"
	"etop.vn/api/main/location"
	"etop.vn/api/main/ordering"
	"etop.vn/api/main/purchaseorder"
	"etop.vn/api/main/purchaserefund"
	"etop.vn/api/main/receipting"
	"etop.vn/api/main/refund"
	"etop.vn/api/main/shipnow"
	carriertypes "etop.vn/api/main/shipnow/carrier/types"
	"etop.vn/api/main/shipping"
	shippingtypes "etop.vn/api/main/shipping/types"
	st "etop.vn/api/main/stocktaking"
	"etop.vn/api/meta"
	"etop.vn/api/shopping/addressing"
	"etop.vn/api/shopping/carrying"
	"etop.vn/api/shopping/customering"
	"etop.vn/api/shopping/suppliering"
	"etop.vn/api/shopping/tradering"
	"etop.vn/api/summary"
	"etop.vn/api/top/int/etop"
	"etop.vn/api/top/int/shop"
	apitypes "etop.vn/api/top/int/types"
	pbcm "etop.vn/api/top/types/common"
	"etop.vn/api/top/types/etc/payment_provider"
	"etop.vn/api/top/types/etc/payment_source"
	notimodel "etop.vn/backend/com/handler/notifier/model"
	catalogmodelx "etop.vn/backend/com/main/catalog/modelx"
	identitymodelx "etop.vn/backend/com/main/identity/modelx"
	moneymodelx "etop.vn/backend/com/main/moneytx/modelx"
	shippingcarrier "etop.vn/backend/com/main/shipping/carrier"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/cmapi"
	"etop.vn/backend/pkg/common/apifw/idemp"
	cmservice "etop.vn/backend/pkg/common/apifw/service"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/redis"
	"etop.vn/backend/pkg/etop/api"
	"etop.vn/backend/pkg/etop/api/convertpb"
	authorizeauth "etop.vn/backend/pkg/etop/authorize/auth"
	"etop.vn/backend/pkg/etop/logic/shipping_provider"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/etop/sqlstore"
	"etop.vn/backend/tools/pkg/acl"
	"etop.vn/capi"
	"etop.vn/capi/dot"
	"etop.vn/common/l"
)

var ll = l.New()

func init() {
	bus.AddHandler("api", miscService.VersionInfo)

	bus.AddHandler("api", inventoryService.CreateInventoryVoucher)
	bus.AddHandler("api", inventoryService.UpdateInventoryVoucher)
	bus.AddHandler("api", inventoryService.AdjustInventoryQuantity)
	bus.AddHandler("api", inventoryService.ConfirmInventoryVoucher)
	bus.AddHandler("api", inventoryService.CancelInventoryVoucher)
	bus.AddHandler("api", inventoryService.GetInventoryVariant)
	bus.AddHandler("api", inventoryService.GetInventoryVariants)
	bus.AddHandler("api", inventoryService.GetInventoryVariantsByVariantIDs)
	bus.AddHandler("api", inventoryService.GetInventoryVoucher)
	bus.AddHandler("api", inventoryService.GetInventoryVouchers)
	bus.AddHandler("api", inventoryService.GetInventoryVouchersByIDs)

	bus.AddHandler("api", productService.CreateProduct)
	bus.AddHandler("api", productService.CreateVariant)
	bus.AddHandler("api", productService.GetProduct)
	bus.AddHandler("api", productService.GetProducts)
	bus.AddHandler("api", productService.GetProductsByIDs)
	bus.AddHandler("api", productService.GetVariant)
	bus.AddHandler("api", productService.GetVariantsByIDs)
	bus.AddHandler("api", productService.RemoveProducts)
	bus.AddHandler("api", productService.RemoveVariants)
	bus.AddHandler("api", productService.UpdateProduct)
	bus.AddHandler("api", productService.UpdateProductImages)
	bus.AddHandler("api", productService.UpdateProductMetaFields)
	bus.AddHandler("api", productService.UpdateProductsStatus)
	bus.AddHandler("api", productService.UpdateProductsTags)
	bus.AddHandler("api", productService.UpdateVariant)
	bus.AddHandler("api", productService.UpdateVariantAttributes)
	bus.AddHandler("api", productService.UpdateVariantImages)
	bus.AddHandler("api", productService.UpdateVariantsStatus)

	bus.AddHandler("api", productSourceService.CreateVariant)
	bus.AddHandler("api", productSourceService.CreateProductSourceCategory)
	bus.AddHandler("api", productSourceService.UpdateProductsPSCategory)

	bus.AddHandler("api", productSourceService.GetProductSourceCategory)
	bus.AddHandler("api", productSourceService.GetProductSourceCategories)
	bus.AddHandler("api", productSourceService.UpdateProductSourceCategory)
	bus.AddHandler("api", productSourceService.RemoveProductSourceCategory)

	bus.AddHandler("api", moneyTransactionService.GetMoneyTransaction)
	bus.AddHandler("api", moneyTransactionService.GetMoneyTransactions)

	bus.AddHandler("api", summaryService.SummarizeFulfillments)
	bus.AddHandler("api", summaryService.SummarizePOS)
	bus.AddHandler("api", summaryService.CalcBalanceShop)

	bus.AddHandler("api", notificationService.CreateDevice)
	bus.AddHandler("api", notificationService.DeleteDevice)
	bus.AddHandler("api", notificationService.GetNotifications)
	bus.AddHandler("api", notificationService.GetNotification)
	bus.AddHandler("api", notificationService.UpdateNotifications)

	bus.AddHandler("api", shipnowService.GetShipnowFulfillment)
	bus.AddHandler("api", shipnowService.GetShipnowFulfillments)
	bus.AddHandler("api", shipnowService.CreateShipnowFulfillment)
	bus.AddHandler("api", shipnowService.ConfirmShipnowFulfillment)
	bus.AddHandler("api", shipnowService.UpdateShipnowFulfillment)
	bus.AddHandler("api", shipnowService.CancelShipnowFulfillment)
	bus.AddHandler("api", shipnowService.GetShipnowServices)
	bus.AddHandler("api", accountService.CreateExternalAccountAhamove)
	bus.AddHandler("api", accountService.GetExternalAccountAhamove)
	bus.AddHandler("api", accountService.RequestVerifyExternalAccountAhamove)
	bus.AddHandler("api", accountService.UpdateExternalAccountAhamoveVerification)

	bus.AddHandler("api", externalAccountService.GetExternalAccountHaravan)
	bus.AddHandler("api", externalAccountService.CreateExternalAccountHaravan)
	bus.AddHandler("api", externalAccountService.UpdateExternalAccountHaravanToken)
	bus.AddHandler("api", externalAccountService.ConnectCarrierServiceExternalAccountHaravan)
	bus.AddHandler("api", externalAccountService.DeleteConnectedCarrierServiceExternalAccountHaravan)

	bus.AddHandler("api", paymentService.PaymentTradingOrder)
	bus.AddHandler("api", paymentService.PaymentCheckReturnData)

	bus.AddHandler("api", categoryService.CreateCategory)
	bus.AddHandler("api", categoryService.GetCategory)
	bus.AddHandler("api", categoryService.GetCategories)
	bus.AddHandler("api", categoryService.UpdateCategory)
	bus.AddHandler("api", categoryService.DeleteCategory)

	bus.AddHandler("api", productService.UpdateProductCategory)
	bus.AddHandler("api", productService.RemoveProductCategory)

	bus.AddHandler("api", collectionService.GetCollection)
	bus.AddHandler("api", collectionService.GetCollections)
	bus.AddHandler("api", collectionService.CreateCollection)
	bus.AddHandler("api", collectionService.UpdateCollection)

	bus.AddHandler("api", productService.AddProductCollection)
	bus.AddHandler("api", productService.RemoveProductCollection)
	bus.AddHandler("api", collectionService.GetCollectionsByProductID)

	bus.AddHandler("api", stocktakeService.CancelStocktake)
	bus.AddHandler("api", stocktakeService.ConfirmStocktake)
	bus.AddHandler("api", stocktakeService.UpdateStocktake)
	bus.AddHandler("api", stocktakeService.CreateStocktake)
	bus.AddHandler("api", stocktakeService.GetStocktake)
	bus.AddHandler("api", stocktakeService.GetStocktakes)
	bus.AddHandler("api", stocktakeService.GetStocktakesByIDs)

	bus.AddHandler("api", brandService.GetBrandsByIDs)
	bus.AddHandler("api", brandService.DeleteBrand)
	bus.AddHandler("api", brandService.UpdateBrandInfo)
	bus.AddHandler("api", brandService.CreateBrand)
	bus.AddHandler("api", brandService.GetBrandByID)
	bus.AddHandler("api", brandService.GetBrands)
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
	supplierAggr         suppliering.CommandBus
	supplierQuery        suppliering.QueryBus
	carrierAggr          carrying.CommandBus
	carrierQuery         carrying.QueryBus
	traderQuery          tradering.QueryBus
	summaryQuery         summary.QueryBus
	eventBus             capi.EventBus
	receiptAggr          receipting.CommandBus
	receiptQuery         receipting.QueryBus
	inventoryAggregate   inventory.CommandBus
	inventoryQuery       inventory.QueryBus
	ledgerAggr           ledgering.CommandBus
	ledgerQuery          ledgering.QueryBus
	purchaseOrderAggr    purchaseorder.CommandBus
	purchaseOrderQuery   purchaseorder.QueryBus
	StocktakeQuery       st.QueryBus
	StocktakeAggregate   st.CommandBus
	shipmentManager      *shippingcarrier.ShipmentManager
	shippingAggregate    shipping.CommandBus
	RefundAggr           refund.CommandBus
	RefundQuery          refund.QueryBus
	PurchaseRefundAggr   purchaserefund.CommandBus
	PurchaseRefundQuery  purchaserefund.QueryBus
	connectionQuery      connectioning.QueryBus
	connectionAggr       connectioning.CommandBus
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
	supplierA suppliering.CommandBus,
	supplierQ suppliering.QueryBus,
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
	ledgerA ledgering.CommandBus,
	ledgerQ ledgering.QueryBus,
	purchaseOrderA purchaseorder.CommandBus,
	purchaseOrderQ purchaseorder.QueryBus,
	summary summary.QueryBus,
	StocktakeQ st.QueryBus,
	StocktakeA st.CommandBus,
	shipmentM *shippingcarrier.ShipmentManager,
	shippingA shipping.CommandBus,
	refundA refund.CommandBus,
	refundQ refund.QueryBus,
	purchaseRefundA purchaserefund.CommandBus,
	purchaseRefundQ purchaserefund.QueryBus,
	connectionQ connectioning.QueryBus,
	connectionA connectioning.CommandBus,
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
	supplierAggr = supplierA
	supplierQuery = supplierQ
	carrierAggr = carrierA
	carrierQuery = carrierQ
	traderQuery = traderQ
	eventBus = eventB
	summaryQuery = summary
	sd.Register(idempgroup.Shutdown)
	inventoryAggregate = inventoryA
	inventoryQuery = inventoryQ
	ledgerAggr = ledgerA
	ledgerQuery = ledgerQ
	purchaseOrderAggr = purchaseOrderA
	purchaseOrderQuery = purchaseOrderQ
	StocktakeQuery = StocktakeQ
	StocktakeAggregate = StocktakeA
	shipmentManager = shipmentM
	shippingAggregate = shippingA
	RefundAggr = refundA
	RefundQuery = refundQ
	PurchaseRefundAggr = purchaseRefundA
	PurchaseRefundQuery = purchaseRefundQ

	connectionQuery = connectionQ
	connectionAggr = connectionA
}

type MiscService struct{}
type InventoryService struct{}
type AccountService struct{}
type ExternalAccountService struct{}
type CollectionService struct{}
type CustomerService struct{}
type CustomerGroupService struct{}
type ProductService struct{}
type CategoryService struct{}
type ProductSourceService struct{}
type OrderService struct{}
type FulfillmentService struct{}
type ShipnowService struct{}
type HistoryService struct{}
type MoneyTransactionService struct{}
type SummaryService struct{}
type ExportService struct{}
type NotificationService struct{}
type AuthorizeService struct{}
type TradingService struct{}
type PaymentService struct{}
type ReceiptService struct{}
type SupplierService struct{}
type CarrierService struct{}
type BrandService struct{}
type LedgerService struct{}
type PurchaseOrderService struct{}
type StocktakeService struct{}
type ShipmentService struct{}
type ConnectionService struct{}
type RefundService struct{}
type PurchaseRefundService struct{}

var miscService = &MiscService{}
var inventoryService = &InventoryService{}
var accountService = &AccountService{}
var externalAccountService = &ExternalAccountService{}
var collectionService = &CollectionService{}
var customerService = &CustomerService{}
var customerGroupService = &CustomerGroupService{}
var productService = &ProductService{}
var categoryService = &CategoryService{}
var productSourceService = &ProductSourceService{}
var orderService = &OrderService{}
var fulfillmentService = &FulfillmentService{}
var shipnowService = &ShipnowService{}
var historyService = &HistoryService{}
var moneyTransactionService = &MoneyTransactionService{}
var summaryService = &SummaryService{}
var exportService = &ExportService{}
var notificationService = &NotificationService{}
var authorizeService = &AuthorizeService{}
var tradingService = &TradingService{}
var paymentService = &PaymentService{}
var receiptService = &ReceiptService{}
var supplierService = &SupplierService{}
var carrierService = &CarrierService{}
var brandService = &BrandService{}
var ledgerService = &LedgerService{}
var purchaseOrderService = &PurchaseOrderService{}
var stocktakeService = &StocktakeService{}
var shipmentService = &ShipmentService{}
var connectionService = &ConnectionService{}
var refundService = &RefundService{}
var purchaseRefundService = &PurchaseRefundService{}

func (s *MiscService) VersionInfo(ctx context.Context, q *VersionInfoEndpoint) error {
	q.Result = &pbcm.VersionInfoResponse{
		Service: "etop.Shop",
		Version: "0.1",
	}
	return nil
}

func (s *BrandService) GetBrandByID(ctx context.Context, q *GetBrandByIDEndpoint) error {
	shopID := q.Context.Shop.ID
	query := &catalog.GetBrandByIDQuery{
		Id:     q.Id,
		ShopID: shopID,
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = PbBrand(query.Result)
	return nil
}

func (s *BrandService) GetBrandsByIDs(ctx context.Context, q *GetBrandsByIDsEndpoint) error {
	shopID := q.Context.Shop.ID
	query := &catalog.GetBrandsByIDsQuery{
		Ids:    q.Ids,
		ShopID: shopID,
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &shop.GetBrandsByIDsResponse{
		Brands: PbBrands(query.Result),
	}
	return nil
}

func (s *BrandService) GetBrands(ctx context.Context, q *GetBrandsEndpoint) error {
	shopID := q.Context.Shop.ID
	query := &catalog.ListBrandsQuery{
		Paging: meta.Paging{
			Offset: q.Paging.Offset,
			Limit:  q.Paging.Limit,
		},
		ShopId: shopID,
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &shop.GetBrandsResponse{
		Brands: PbBrands(query.Result.ShopBrands),
		Paging: cmapi.PbPaging(query.Paging),
	}
	return nil
}

func (s *BrandService) CreateBrand(ctx context.Context, q *CreateBrandEndpoint) error {
	shopID := q.Context.Shop.ID
	cmd := &catalog.CreateBrandCommand{
		ShopID:      shopID,
		BrandName:   q.Name,
		Description: q.Description,
	}
	if err := catalogAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = PbBrand(cmd.Result)
	return nil
}

func (s *BrandService) UpdateBrandInfo(ctx context.Context, q *UpdateBrandInfoEndpoint) error {
	shopID := q.Context.Shop.ID
	cmd := &catalog.UpdateBrandInfoCommand{
		ShopID:      shopID,
		ID:          q.Id,
		BrandName:   q.Name,
		Description: q.Description,
	}
	if err := catalogAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = PbBrand(cmd.Result)
	return nil
}

func (s *BrandService) DeleteBrand(ctx context.Context, q *DeleteBrandEndpoint) error {
	shopID := q.Context.Shop.ID
	cmd := &catalog.DeleteShopBrandCommand{
		ShopId: shopID,
		Ids:    q.Ids,
	}
	if err := catalogAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &shop.DeleteBrandResponse{
		Count: cmd.Result,
	}
	return nil
}

func (s *ProductService) UpdateVariant(ctx context.Context, q *UpdateVariantEndpoint) error {
	shopID := q.Context.Shop.ID
	cmd := &catalog.UpdateShopVariantInfoCommand{
		ShopID:    shopID,
		VariantID: q.Id,
		Name:      q.Name,
		Code:      q.Code,
		Note:      q.Note,

		ShortDesc:    q.ShortDesc,
		Descripttion: q.Description,
		DescHTML:     q.DescHtml,

		CostPrice:   q.CostPrice,
		ListPrice:   q.ListPrice,
		RetailPrice: q.RetailPrice,
	}
	if err := catalogAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = PbShopVariant(cmd.Result)
	return nil
}

func (s *ProductService) UpdateVariantAttributes(ctx context.Context, q *UpdateVariantAttributesEndpoint) error {
	shopID := q.Context.Shop.ID
	cmd := &catalog.UpdateShopVariantAttributesCommand{
		ShopID:     shopID,
		VariantID:  q.VariantId,
		Attributes: q.Attributes,
	}
	if err := catalogAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = PbShopVariant(cmd.Result)
	return nil
}

func (s *ProductService) RemoveVariants(ctx context.Context, q *RemoveVariantsEndpoint) error {
	cmd := &catalog.DeleteShopVariantsCommand{
		ShopID: q.Context.Shop.ID,
		IDs:    q.Ids,
	}
	if err := catalogAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.RemovedResponse{
		Removed: cmd.Result,
	}
	return nil
}

func (s *ProductService) GetProduct(ctx context.Context, q *GetProductEndpoint) error {
	shopID := q.Context.Shop.ID
	query := &catalog.GetShopProductWithVariantsByIDQuery{
		ProductID: q.Id,
		ShopID:    q.Context.Shop.ID,
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	productPb, err := getProductQuantity(ctx, shopID, query.Result)
	if err != nil {
		return err
	}
	q.Result = productPb
	return nil
}

func (s *ProductService) GetProductsByIDs(ctx context.Context, q *GetProductsByIDsEndpoint) error {
	shopID := q.Context.Shop.ID
	query := &catalog.ListShopProductsWithVariantsByIDsQuery{
		IDs:    q.Ids,
		ShopID: shopID,
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	products, err := getProductsQuantity(ctx, shopID, query.Result.Products)
	if err != nil {
		return err
	}
	q.Result = &shop.ShopProductsResponse{
		Products: products,
	}
	return nil
}

func (s *ProductService) GetProducts(ctx context.Context, q *GetProductsEndpoint) error {
	paging := cmapi.CMPaging(q.Paging)
	shopID := q.Context.Shop.ID
	query := &catalog.ListShopProductsWithVariantsQuery{
		ShopID:  shopID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(q.Filters),
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	products, err := getProductsQuantity(ctx, shopID, query.Result.Products)
	if err != nil {
		return err
	}
	q.Result = &shop.ShopProductsResponse{
		Paging: cmapi.PbPaging(cm.Paging{
			Limit: query.Result.Paging.Limit,
			Sort:  query.Result.Paging.Sort,
		}),
		Products: products,
	}
	return nil
}

func (s *ProductService) CreateProduct(ctx context.Context, q *CreateProductEndpoint) error {
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
		BrandID:     q.BrandId,
		ProductType: q.ProductType.Apply(0),
		MetaFields:  metaFields,
	}
	if err := catalogAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = PbShopProductWithVariants(cmd.Result)
	return nil
}

func (s *ProductService) RemoveProducts(ctx context.Context, q *RemoveProductsEndpoint) error {
	cmd := &catalog.DeleteShopProductsCommand{
		ShopID: q.Context.Shop.ID,
		IDs:    q.Ids,
	}
	if err := catalogAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.RemovedResponse{
		Removed: cmd.Result,
	}
	return nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, q *UpdateProductEndpoint) error {
	shopID := q.Context.Shop.ID
	cmd := &catalog.UpdateShopProductInfoCommand{
		ShopID:    shopID,
		ProductID: q.Id,
		Code:      q.Code,
		Name:      q.Name,
		Unit:      q.Unit,
		Note:      q.Note,
		BrandID:   q.BrandId,

		ShortDesc:   q.ShortDesc,
		Description: q.Description,
		DescHTML:    q.DescHtml,
		CategoryID:  q.CategoryID,
		CostPrice:   q.CostPrice,
		ListPrice:   q.ListPrice,
		RetailPrice: q.RetailPrice,
		ProductType: q.ProductType,
	}
	if err := catalogAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = PbShopProductWithVariants(cmd.Result)
	return nil
}

func (s *ProductService) UpdateProductsStatus(ctx context.Context, q *UpdateProductsStatusEndpoint) error {
	shopID := q.Context.Shop.ID
	cmd := &catalog.UpdateShopProductStatusCommand{
		IDs:    q.Ids,
		ShopID: shopID,
		Status: int16(q.Status),
	}
	if err := catalogAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &shop.UpdateProductStatusResponse{Updated: cmd.Result}
	return nil
}

func (s *ProductService) UpdateVariantsStatus(ctx context.Context, q *UpdateVariantsStatusEndpoint) error {
	shopID := q.Context.Shop.ID
	cmd := &catalog.UpdateShopVariantStatusCommand{
		IDs:    q.Ids,
		ShopID: shopID,
		Status: int16(q.Status),
	}
	if err := catalogAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &shop.UpdateProductStatusResponse{Updated: cmd.Result}
	return nil
}

func (s *ProductService) UpdateProductsTags(ctx context.Context, q *UpdateProductsTagsEndpoint) error {
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
		Updated: cmd.Result.Updated,
	}
	return nil
}

func (s *ProductService) GetVariant(ctx context.Context, q *GetVariantEndpoint) error {
	query := &catalog.GetShopVariantQuery{
		Code:      q.Code,
		VariantID: q.ID,
		ShopID:    q.Context.Shop.ID,
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = PbShopVariant(query.Result)

	return nil
}

func (s *ProductService) GetVariantsByIDs(ctx context.Context, q *GetVariantsByIDsEndpoint) error {
	query := &catalog.ListShopVariantsWithProductByIDsQuery{
		IDs:    q.Ids,
		ShopID: q.Context.Shop.ID,
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = &shop.ShopVariantsResponse{Variants: PbShopVariantsWithProducts(query.Result.Variants)}

	return nil
}

func (s *ProductService) CreateVariant(ctx context.Context, q *CreateVariantEndpoint) error {
	cmd := &catalog.CreateShopVariantCommand{
		ShopID:     q.Context.Shop.ID,
		ProductID:  q.ProductId,
		Code:       q.Code,
		Name:       q.Name,
		ImageURLs:  q.ImageUrls,
		Note:       q.Note,
		Attributes: q.Attributes,
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

func (s *ProductSourceService) CreateVariant(ctx context.Context, q *DeprecatedCreateVariantEndpoint) error {
	cmd := &catalogmodelx.DeprecatedCreateVariantCommand{
		ShopID:      q.Context.Shop.ID,
		ProductID:   q.ProductId,
		ProductName: q.ProductName,
		Name:        q.Name,
		Description: q.Description,
		ShortDesc:   q.ShortDesc,
		ImageURLs:   q.ImageUrls,
		Tags:        q.Tags,
		Status:      q.Status,

		CostPrice:   q.CostPrice,
		ListPrice:   q.ListPrice,
		RetailPrice: q.RetailPrice,

		ProductCode:       q.Code,
		VariantCode:       q.Code,
		QuantityAvailable: q.QuantityAvailable,
		QuantityOnHand:    q.QuantityOnHand,
		QuantityReserved:  q.QuantityReserved,

		Attributes: q.Attributes,
		DescHTML:   q.DescHtml,
	}

	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}

	q.Result = PbShopProductWithVariants(cmd.Result)
	return nil
}

func (s *ProductSourceService) CreateProductSourceCategory(ctx context.Context, q *CreateProductSourceCategoryEndpoint) error {
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

func (s *ProductSourceService) UpdateProductsPSCategory(ctx context.Context, q *UpdateProductsPSCategoryEndpoint) error {
	cmd := &catalogmodelx.UpdateProductsShopCategoryCommand{
		CategoryID: q.CategoryId,
		ProductIDs: q.ProductIds,
		ShopID:     q.Context.Shop.ID,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.UpdatedResponse{
		Updated: cmd.Result.Updated,
	}
	return nil
}

func (s *ProductSourceService) GetProductSourceCategory(ctx context.Context, q *GetProductSourceCategoryEndpoint) error {
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

func (s *ProductSourceService) GetProductSourceCategories(ctx context.Context, q *GetProductSourceCategoriesEndpoint) error {
	cmd := &catalogmodelx.GetProductSourceCategoriesQuery{
		ShopID: q.Context.Shop.ID,
	}

	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}

	q.Result = &shop.CategoriesResponse{
		Categories: convertpb.PbCategories(cmd.Result.Categories),
	}
	return nil
}

func (s *ProductSourceService) UpdateProductSourceCategory(ctx context.Context, q *UpdateProductSourceCategoryEndpoint) error {
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

func (s *ProductSourceService) RemoveProductSourceCategory(ctx context.Context, q *RemoveProductSourceCategoryEndpoint) error {
	cmd := &catalogmodelx.RemoveShopCategoryCommand{
		ID:     q.Id,
		ShopID: q.Context.Shop.ID,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.RemovedResponse{
		Removed: cmd.Result.Removed,
	}
	return nil
}

func (s *ProductService) UpdateProductImages(ctx context.Context, q *UpdateProductImagesEndpoint) error {
	shopID := q.Context.Shop.ID

	var metaUpdate []*meta.UpdateSet
	if q.DeleteAll {
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

func (s *ProductService) UpdateProductMetaFields(ctx context.Context, q *UpdateProductMetaFieldsEndpoint) error {
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

func (s *ProductService) UpdateVariantImages(ctx context.Context, q *UpdateVariantImagesEndpoint) error {
	shopID := q.Context.Shop.ID

	var metaUpdate []*meta.UpdateSet
	if q.DeleteAll {
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

func (s *MoneyTransactionService) GetMoneyTransaction(ctx context.Context, q *GetMoneyTransactionEndpoint) error {
	query := &moneymodelx.GetMoneyTransaction{
		ShopID: q.Context.Shop.ID,
		ID:     q.Id,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = convertpb.PbMoneyTransactionExtended(query.Result)
	return nil
}

func (s *MoneyTransactionService) GetMoneyTransactions(ctx context.Context, q *GetMoneyTransactionsEndpoint) error {
	paging := cmapi.CMPaging(q.Paging)
	query := &moneymodelx.GetMoneyTransactions{
		ShopID:  q.Context.Shop.ID,
		Paging:  paging,
		Filters: cmapi.ToFilters(q.Filters),
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &apitypes.MoneyTransactionsResponse{
		MoneyTransactions: convertpb.PbMoneyTransactionExtendeds(query.Result.MoneyTransactions),
		Paging:            cmapi.PbPageInfo(paging),
	}
	return nil
}

func (s *SummaryService) SummarizeFulfillments(ctx context.Context, q *SummarizeFulfillmentsEndpoint) error {
	query := &model.SummarizeFulfillmentsRequest{
		ShopID:   q.Context.Shop.ID,
		DateFrom: q.DateFrom,
		DateTo:   q.DateTo,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = &shop.SummarizeFulfillmentsResponse{
		Tables: convertpb.PbSummaryTables(query.Result.Tables),
	}
	return nil
}

func (s *SummaryService) SummarizeTopShip(ctx context.Context, q *SummarizeTopShipEndpoint) error {
	dateFrom, dateTo, err := cm.ParseDateFromTo(q.DateFrom, q.DateTo)
	if err != nil {
		return err
	}
	query := &summary.SummaryTopShipQuery{
		ShopID:   q.Context.Shop.ID,
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	if err = summaryQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &shop.SummarizeTopShipResponse{
		Tables: convertpb.PbSummaryTablesNew(query.Result.ListTable),
	}
	return nil
}

func (s *SummaryService) SummarizePOS(ctx context.Context, q *SummarizePOSEndpoint) error {
	dateFrom, dateTo, err := cm.ParseDateFromTo(q.DateFrom, q.DateTo)
	if err != nil {
		return err
	}
	query := &summary.SummaryPOSQuery{
		ShopID:   q.Context.Shop.ID,
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	if err = summaryQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &shop.SummarizePOSResponse{
		Tables: convertpb.PbSummaryTablesNew(query.Result.ListTable),
	}
	return nil
}

func (s *SummaryService) CalcBalanceShop(ctx context.Context, q *CalcBalanceShopEndpoint) error {
	query := &model.GetBalanceShopCommand{
		ShopID: q.Context.Shop.ID,
	}

	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &shop.CalcBalanceShopResponse{
		Balance: query.Result.Amount,
	}
	return nil
}

func (s *NotificationService) CreateDevice(ctx context.Context, q *CreateDeviceEndpoint) error {
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
	q.Result = convertpb.PbDevice(device)
	return nil
}

func (s *NotificationService) DeleteDevice(ctx context.Context, q *DeleteDeviceEndpoint) error {
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

func (s *NotificationService) GetNotification(ctx context.Context, q *GetNotificationEndpoint) error {
	query := &notimodel.GetNotificationArgs{
		AccountID: q.Context.Shop.ID,
		ID:        q.Id,
	}
	noti, err := sqlstore.GetNotification(ctx, query)
	if err != nil {
		return err
	}
	q.Result = convertpb.PbNotification(noti)
	return nil
}

func (s *NotificationService) GetNotifications(ctx context.Context, q *GetNotificationsEndpoint) error {
	paging := cmapi.CMPaging(q.Paging)
	query := &notimodel.GetNotificationsArgs{
		Paging:    paging,
		AccountID: q.Context.Shop.ID,
	}
	notis, err := sqlstore.GetNotifications(ctx, query)
	if err != nil {
		return err
	}
	q.Result = &etop.NotificationsResponse{
		Notifications: convertpb.PbNotifications(notis),
		Paging:        cmapi.PbPageInfo(paging),
	}
	return nil
}

func (s *NotificationService) UpdateNotifications(ctx context.Context, q *UpdateNotificationsEndpoint) error {
	cmd := &notimodel.UpdateNotificationsArgs{
		IDs:    q.Ids,
		IsRead: q.IsRead,
	}
	if err := sqlstore.UpdateNotifications(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.UpdatedResponse{
		Updated: len(q.Ids),
	}
	return nil
}

func (s *ShipnowService) GetShipnowFulfillment(ctx context.Context, q *GetShipnowFulfillmentEndpoint) error {
	query := &shipnow.GetShipnowFulfillmentQuery{
		Id:     q.Id,
		ShopId: q.Context.Shop.ID,
		Result: nil,
	}
	if err := shipnowQuery.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = convertpb.Convert_core_ShipnowFulfillment_To_api_ShipnowFulfillment(query.Result.ShipnowFulfillment)
	return nil
}

func (s *ShipnowService) GetShipnowFulfillments(ctx context.Context, q *GetShipnowFulfillmentsEndpoint) error {
	shopIDs, err := api.MixAccount(q.Context.Claim, q.Mixed)
	if err != nil {
		return err
	}
	paging := cmapi.CMPaging(q.Paging)

	query := &shipnow.GetShipnowFulfillmentsQuery{
		ShopIds: shopIDs,
		Paging:  paging,
		Filters: cmapi.ToFiltersPtr(q.Filters),
	}
	if err := shipnowQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &apitypes.ShipnowFulfillments{
		ShipnowFulfillments: convertpb.Convert_core_ShipnowFulfillments_To_api_ShipnowFulfillments(query.Result.ShipnowFulfillments),
		Paging:              cmapi.PbPageInfo(paging),
	}
	return nil
}

func (s *ShipnowService) CreateShipnowFulfillment(ctx context.Context, q *CreateShipnowFulfillmentEndpoint) error {
	pickupAddress, err := convertpb.OrderAddressFulfilled(q.PickupAddress)
	if err != nil {
		return err
	}
	_carrier, _ := carriertypes.ParseCarrier(q.Carrier)
	cmd := &shipnow.CreateShipnowFulfillmentCommand{
		OrderIds:            q.OrderIds,
		Carrier:             _carrier,
		ShopId:              q.Context.Shop.ID,
		ShippingServiceCode: q.ShippingServiceCode,
		ShippingServiceFee:  q.ShippingServiceFee,
		ShippingNote:        q.ShippingNote,
		RequestPickupAt:     time.Time{},
		PickupAddress:       convertpb.Convert_api_OrderAddress_To_core_OrderAddress(pickupAddress),
	}
	if err := shipnowAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = convertpb.Convert_core_ShipnowFulfillment_To_api_ShipnowFulfillment(cmd.Result)
	return nil
}

func (s *ShipnowService) CreateShipnowFulfillmentV2(ctx context.Context, q *CreateShipnowFulfillmentV2Endpoint) error {
	pickupAddress, err := convertpb.OrderAddressFulfilled(q.PickupAddress)
	if err != nil {
		return err
	}

	var deliveryPoints []*shipnow.OrderShippingInfo
	for _, point := range q.DeliveryPoints {
		shippingAddress, err := convertpb.OrderAddressFulfilled(point.ShippingAddress)
		if err != nil {
			return err
		}
		p := &shipnow.OrderShippingInfo{
			OrderID:         point.OrderID,
			ShippingAddress: convertpb.Convert_api_OrderAddress_To_core_OrderAddress(shippingAddress),
			ShippingNote:    point.ShippingNote,
			WeightInfo: shippingtypes.WeightInfo{
				GrossWeight:      point.GrossWeight,
				ChargeableWeight: point.ChargeableWeight,
			},
			ValueInfo: shippingtypes.ValueInfo{
				CODAmount: point.CODAmount,
			},
		}
		deliveryPoints = append(deliveryPoints, p)
	}
	cmd := &shipnow.CreateShipnowFulfillmentV2Command{
		DeliveryPoints:      deliveryPoints,
		Carrier:             q.Carrier,
		ShopID:              q.Context.Shop.ID,
		ShippingServiceCode: q.ShippingServiceCode,
		ShippingServiceFee:  q.ShippingServiceFee,
		ShippingNote:        q.ShippingNote,
		PickupAddress:       convertpb.Convert_api_OrderAddress_To_core_OrderAddress(pickupAddress),
	}
	if err := shipnowAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = convertpb.Convert_core_ShipnowFulfillment_To_api_ShipnowFulfillment(cmd.Result)
	return nil
}

func (s *ShipnowService) ConfirmShipnowFulfillment(ctx context.Context, q *ConfirmShipnowFulfillmentEndpoint) error {
	cmd := &shipnow.ConfirmShipnowFulfillmentCommand{
		Id:     q.Id,
		ShopId: q.Context.Shop.ID,
	}
	if err := shipnowAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = convertpb.Convert_core_ShipnowFulfillment_To_api_ShipnowFulfillment(cmd.Result)
	return nil
}

func (s *ShipnowService) UpdateShipnowFulfillment(ctx context.Context, q *UpdateShipnowFulfillmentEndpoint) error {
	pickupAddress, err := convertpb.OrderAddressFulfilled(q.PickupAddress)
	if err != nil {
		return err
	}
	_carrier, _ := carriertypes.ParseCarrier(q.Carrier)
	cmd := &shipnow.UpdateShipnowFulfillmentCommand{
		Id:                  q.Id,
		OrderIds:            q.OrderIds,
		Carrier:             _carrier,
		ShopId:              q.Context.Shop.ID,
		ShippingServiceCode: q.ShippingServiceCode,
		ShippingServiceFee:  q.ShippingServiceFee,
		ShippingNote:        q.ShippingNote,
		RequestPickupAt:     time.Time{},
		PickupAddress:       convertpb.Convert_api_OrderAddress_To_core_OrderAddress(pickupAddress),
	}
	if err := shipnowAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = convertpb.Convert_core_ShipnowFulfillment_To_api_ShipnowFulfillment(cmd.Result)
	return nil
}

func (s *ShipnowService) CancelShipnowFulfillment(ctx context.Context, q *CancelShipnowFulfillmentEndpoint) error {
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

func (s *ShipnowService) GetShipnowServices(ctx context.Context, q *GetShipnowServicesEndpoint) error {
	pickupAddress, err := convertpb.OrderAddressFulfilled(q.PickupAddress)
	if err != nil {
		return err
	}
	var points []*shipnow.DeliveryPoint
	if len(q.DeliveryPoints) > 0 {
		for _, p := range q.DeliveryPoints {
			addr, err := convertpb.OrderAddressFulfilled(p.ShippingAddress)
			if err != nil {
				return err
			}
			points = append(points, &shipnow.DeliveryPoint{
				ShippingAddress: convertpb.Convert_api_OrderAddress_To_core_OrderAddress(addr),
				ValueInfo: shippingtypes.ValueInfo{
					CODAmount: p.CodAmount,
				},
			})
		}
	}

	cmd := &shipnow.GetShipnowServicesCommand{
		ShopId:         q.Context.Shop.ID,
		OrderIds:       q.OrderIds,
		PickupAddress:  convertpb.Convert_api_OrderAddress_To_core_OrderAddress(pickupAddress),
		DeliveryPoints: points,
	}
	if err := shipnowAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &apitypes.GetShipnowServicesResponse{
		Services: convertpb.Convert_core_ShipnowServices_To_api_ShipnowServices(cmd.Result.Services),
	}
	return nil
}

func (s *AccountService) CreateExternalAccountAhamove(ctx context.Context, q *CreateExternalAccountAhamoveEndpoint) error {
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
	q.Result = convertpb.Convert_core_XAccountAhamove_To_api_XAccountAhamove(cmd.Result, false)
	return nil
}

func (s *AccountService) GetExternalAccountAhamove(ctx context.Context, q *GetExternalAccountAhamoveEndpoint) error {
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

	var hideInfo bool
	if !authorization.IsContainsActionString(authorizeauth.ListActionsByRoles(q.Context.Roles), string(acl.ShopExternalAccountManage)) {
		hideInfo = true
	}
	q.Result = convertpb.Convert_core_XAccountAhamove_To_api_XAccountAhamove(account, hideInfo)
	return nil
}

func (s *AccountService) RequestVerifyExternalAccountAhamove(ctx context.Context, q *RequestVerifyExternalAccountAhamoveEndpoint) error {
	query := &identitymodelx.GetUserByIDQuery{
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

func (s *AccountService) UpdateExternalAccountAhamoveVerification(ctx context.Context, r *UpdateExternalAccountAhamoveVerificationEndpoint) error {
	if err := validateUrl(r.IdCardFrontImg, r.IdCardBackImg, r.PortraitImg, r.WebsiteUrl, r.FanpageUrl); err != nil {
		return err
	}
	if err := validateUrl(r.BusinessLicenseImgs...); err != nil {
		return err
	}
	if err := validateUrl(r.CompanyImgs...); err != nil {
		return err
	}

	query := &identitymodelx.GetUserByIDQuery{
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

func (s *ExternalAccountService) GetExternalAccountHaravan(ctx context.Context, r *GetExternalAccountHaravanEndpoint) error {
	query := &haravanidentity.GetExternalAccountHaravanByShopIDQuery{
		ShopID: r.Context.Shop.ID,
	}
	if err := haravanIdentityQuery.Dispatch(ctx, query); err != nil {
		return err
	}

	var hideInfo bool
	if !authorization.IsContainsActionString(authorizeauth.ListActionsByRoles(r.Context.Roles), string(acl.ShopExternalAccountManage)) {
		hideInfo = true
	}
	r.Result = convertpb.Convert_core_XAccountHaravan_To_api_XAccountHaravan(query.Result, hideInfo)
	return nil
}

func (s *ExternalAccountService) CreateExternalAccountHaravan(ctx context.Context, r *CreateExternalAccountHaravanEndpoint) error {
	cmd := &haravanidentity.CreateExternalAccountHaravanCommand{
		ShopID:      r.Context.Shop.ID,
		Subdomain:   r.Subdomain,
		Code:        r.Code,
		RedirectURI: r.RedirectUri,
	}
	if err := haravanIdentityAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = convertpb.Convert_core_XAccountHaravan_To_api_XAccountHaravan(cmd.Result, false)
	return nil
}

func (s *ExternalAccountService) UpdateExternalAccountHaravanToken(ctx context.Context, r *UpdateExternalAccountHaravanTokenEndpoint) error {
	cmd := &haravanidentity.UpdateExternalAccountHaravanTokenCommand{
		ShopID:      r.Context.Shop.ID,
		Subdomain:   r.Subdomain,
		RedirectURI: r.RedirectUri,
		Code:        r.Code,
	}
	if err := haravanIdentityAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = convertpb.Convert_core_XAccountHaravan_To_api_XAccountHaravan(cmd.Result, false)
	return nil
}

func (s *ExternalAccountService) ConnectCarrierServiceExternalAccountHaravan(ctx context.Context, r *ConnectCarrierServiceExternalAccountHaravanEndpoint) error {
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

func (s *ExternalAccountService) DeleteConnectedCarrierServiceExternalAccountHaravan(ctx context.Context, r *DeleteConnectedCarrierServiceExternalAccountHaravanEndpoint) error {
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

func (s *PaymentService) PaymentTradingOrder(ctx context.Context, q *PaymentTradingOrderEndpoint) error {
	if q.OrderId == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing OrderID")
	}
	if q.ReturnUrl == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing ReturnURL")
	}

	argGenCode := &paymentmanager.GenerateCodeCommand{
		PaymentSource: payment_source.PaymentSourceOrder,
		ID:            q.OrderId.String(),
	}
	if err := paymentCtrl.Dispatch(ctx, argGenCode); err != nil {
		return err
	}
	args := &paymentmanager.BuildUrlConnectPaymentGatewayCommand{
		OrderID:           argGenCode.Result,
		Desc:              q.Desc,
		ReturnURL:         q.ReturnUrl,
		TransactionAmount: q.Amount,
		Provider:          payment_provider.PaymentProvider(q.PaymentProvider),
	}

	if err := paymentCtrl.Dispatch(ctx, args); err != nil {
		return err
	}
	q.Result = &shop.PaymentTradingOrderResponse{
		Url: args.Result,
	}
	return nil
}

func (s *PaymentService) PaymentCheckReturnData(ctx context.Context, q *PaymentCheckReturnDataEndpoint) error {
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
		Amount:                q.Amount,
		ExternalTransactionID: q.ExternalTransactionId,
		Provider:              payment_provider.PaymentProvider(q.PaymentProvider),
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

func (s *CategoryService) CreateCategory(ctx context.Context, q *CreateCategoryEndpoint) error {
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

func (s *CategoryService) GetCategory(ctx context.Context, q *GetCategoryEndpoint) error {
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

func (s *CategoryService) GetCategories(ctx context.Context, q *GetCategoriesEndpoint) error {
	paging := cmapi.CMPaging(q.Paging)
	query := &catalog.ListShopCategoriesQuery{
		ShopID:  q.Context.Shop.ID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(q.Filters),
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = &shop.ShopCategoriesResponse{
		Paging:     cmapi.PbPageInfo(paging),
		Categories: PbShopCategories(query.Result.Categories),
	}
	return nil
}

func (s *CategoryService) UpdateCategory(ctx context.Context, q *UpdateCategoryEndpoint) error {
	shopID := q.Context.Shop.ID
	cmd := &catalog.UpdateShopCategoryCommand{
		ID:       q.Id,
		ShopID:   shopID,
		Name:     q.Name,
		ParentID: q.ParentId,
	}
	if err := catalogAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = PbShopCategory(cmd.Result)
	return nil
}

func (s *CategoryService) DeleteCategory(ctx context.Context, r *DeleteCategoryEndpoint) error {
	cmd := &catalog.DeleteShopCategoryCommand{
		ID:     r.Id,
		ShopID: r.Context.Shop.ID,
	}
	if err := catalogAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.DeletedResponse{Deleted: cmd.Result}
	return nil
}

func (s *ProductService) UpdateProductCategory(ctx context.Context, q *UpdateProductCategoryEndpoint) error {
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

func (s *CollectionService) GetCollection(ctx context.Context, q *GetCollectionEndpoint) error {
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

func (s *CollectionService) GetCollections(ctx context.Context, q *GetCollectionsEndpoint) error {
	paging := cmapi.CMPaging(q.Paging)
	query := &catalog.ListShopCollectionsQuery{
		ShopID:  q.Context.Shop.ID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(q.Filters),
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &shop.ShopCollectionsResponse{
		Paging:      cmapi.PbPageInfo(paging),
		Collections: PbShopCollections(query.Result.Collections),
	}
	return nil
}

func (s *CollectionService) UpdateCollection(ctx context.Context, q *UpdateCollectionEndpoint) error {
	shopID := q.Context.Shop.ID
	cmd := &catalog.UpdateShopCollectionCommand{
		ID:          q.Id,
		ShopID:      shopID,
		Name:        q.Name,
		Description: q.Description,
		DescHTML:    q.DescHtml,
		ShortDesc:   q.ShortDesc,
	}
	if err := catalogAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = PbShopCollection(cmd.Result)
	return nil
}

func (s *CollectionService) CreateCollection(ctx context.Context, q *CreateCollectionEndpoint) error {
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

func (s *ProductService) AddProductCollection(ctx context.Context, r *AddProductCollectionEndpoint) error {
	cmd := &catalog.AddShopProductCollectionCommand{
		ProductID:     r.ProductId,
		CollectionIDs: r.CollectionIds,
		ShopID:        r.Context.Shop.ID,
	}
	if err := catalogAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.UpdatedResponse{Updated: cmd.Result}
	return nil
}

func (s *ProductService) RemoveProductCollection(ctx context.Context, r *RemoveProductCollectionEndpoint) error {
	cmd := &catalog.RemoveShopProductCollectionCommand{
		ProductID:     r.ProductId,
		CollectionIDs: r.CollectionIds,
		ShopID:        r.Context.Shop.ID,
	}
	if err := catalogAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.RemovedResponse{Removed: cmd.Result}
	return nil
}

func (s *CollectionService) GetCollectionsByProductID(ctx context.Context, q *GetCollectionsByProductIDEndpoint) error {
	query := &catalog.ListShopCollectionsByProductIDQuery{
		ShopID:    q.Context.Shop.ID,
		ProductID: q.ProductId,
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &shop.CollectionsResponse{
		Collections: PbShopCollections(query.Result),
	}
	return nil
}

func (s *ProductService) RemoveProductCategory(ctx context.Context, r *RemoveProductCategoryEndpoint) error {
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

func getProductsQuantity(ctx context.Context, shopID dot.ID, products []*catalog.ShopProductWithVariants) ([]*shop.ShopProduct, error) {
	var variantIDs []dot.ID
	for _, valueProduct := range products {
		for _, valueVariant := range valueProduct.Variants {
			variantIDs = append(variantIDs, valueVariant.VariantID)
		}
	}
	inventoryVariants, err := getVariantsQuantity(ctx, shopID, variantIDs)
	if err != nil {
		return nil, err
	}
	return PbProductsQuantity(products, inventoryVariants), nil
}

func getProductQuantity(ctx context.Context, shopID dot.ID, shopProduct *catalog.ShopProductWithVariants) (*shop.ShopProduct, error) {
	var variantIDs []dot.ID
	for _, variant := range shopProduct.Variants {
		variantIDs = append(variantIDs, variant.VariantID)
	}
	inventoryVariants, err := getVariantsQuantity(ctx, shopID, variantIDs)
	if err != nil {
		return nil, err
	}
	shopProductPb := PbProductQuantity(shopProduct, inventoryVariants)
	return shopProductPb, nil
}

func getVariantsQuantity(ctx context.Context, shopID dot.ID, variantIDs []dot.ID) (map[dot.ID]*inventory.InventoryVariant, error) {

	var mapInventoryVariant = make(map[dot.ID]*inventory.InventoryVariant)
	if len(variantIDs) == 0 {
		return mapInventoryVariant, nil
	}
	q := &inventory.GetInventoryVariantsByVariantIDsQuery{
		ShopID:     shopID,
		VariantIDs: variantIDs,
	}
	if err := inventoryQuery.Dispatch(ctx, q); err != nil {
		return nil, err
	}

	for _, value := range q.Result.InventoryVariants {
		mapInventoryVariant[value.VariantID] = value
	}
	return mapInventoryVariant, nil
}

func (s *ProductService) GetVariantsBySupplierID(ctx context.Context, q *GetVariantsBySupplierIDEndpoint) error {
	query := &catalog.GetVariantsBySupplierIDQuery{
		SupplierID: q.SupplierId,
		ShopID:     q.Context.Shop.ID,
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &shop.ShopVariantsResponse{Variants: PbShopVariants(query.Result.Variants)}
	return nil
}
