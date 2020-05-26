package partner

import (
	"net/url"

	service "o.o/api/top/external/partner"
	"o.o/backend/pkg/common/apifw/idemp"
	cmService "o.o/backend/pkg/common/apifw/service"
	"o.o/backend/pkg/common/authorization/auth"
	"o.o/backend/pkg/common/redis"
	"o.o/capi/httprpc"
	"o.o/common/l"
)

// +gen:wrapper=o.o/api/top/external/partner
// +gen:wrapper:package=partner
// +gen:wrapper:prefix=ext

var (
	idempgroup *idemp.RedisGroup
	authStore  auth.Generator
	authURL    string

	ll = l.New()
)

const PrefixIdempPartnerAPI = "IdempPartnerAPI"

const ttlShopRequest = 15 * 60 // 15 minutes
const msgShopRequest = `Sử dụng mã này để hỏi quyền tạo đơn hàng với tư cách shop (có hiệu lực trong 15 phút)`
const msgShopKey = `Sử dụng mã này để tạo đơn hàng với tư cách shop (có hiệu lực khi shop vẫn tiếp tục sử dụng dịch vụ của đối tác)`
const msgUserKey = `Sử dụng mã này để truy cập hệ thống với tư cách user (có hiệu lực khi user vẫn tiếp tục sử dụng dịch vụ của đối tác)`

type Servers []httprpc.Server

type AuthURL string

func NewServers(
	sd cmService.Shutdowner,
	rd redis.Store,
	s auth.Generator,
	_authURL AuthURL,
	miscService *MiscService,
	shopService *ShopService,
	webhookService *WebhookService,
	historyService *HistoryService,
	shippingService *ShippingService,
	orderService *OrderService,
	fulfillmentService *FulfillmentService,
	customerService *CustomerService,
	customerAddressService *CustomerAddressService,
	customerGroupService *CustomerGroupService,
	customerGroupRelationshipService *CustomerGroupRelationshipService,
	inventoryService *InventoryService,
	variantService *VariantService,
	productService *ProductService,
	productCollectionService *ProductCollectionService,
	productCollectionRelationshipService *ProductCollectionRelationshipService,
) Servers {

	authURL = string(_authURL)
	if authURL == "" {
		ll.Panic("no auth_url")
	}
	if _, err := url.Parse(authURL); err != nil {
		ll.Panic("invalid auth_url", l.String("url", authURL))
	}

	idempgroup = idemp.NewRedisGroup(rd, PrefixIdempPartnerAPI, 0)
	sd.Register(idempgroup.Shutdown)

	servers := []httprpc.Server{
		service.NewMiscServiceServer(WrapMiscService(miscService.Clone)),
		service.NewShopServiceServer(WrapShopService(shopService.Clone)),
		service.NewWebhookServiceServer(WrapWebhookService(webhookService.Clone)),
		service.NewHistoryServiceServer(WrapHistoryService(historyService.Clone)),
		service.NewShippingServiceServer(WrapShippingService(shippingService.Clone)),
		service.NewOrderServiceServer(WrapOrderService(orderService.Clone)),
		service.NewFulfillmentServiceServer(WrapFulfillmentService(fulfillmentService.Clone)),
		service.NewCustomerServiceServer(WrapCustomerService(customerService.Clone)),
		service.NewCustomerAddressServiceServer(WrapCustomerAddressService(customerAddressService.Clone)),
		service.NewCustomerGroupServiceServer(WrapCustomerGroupService(customerGroupService.Clone)),
		service.NewCustomerGroupRelationshipServiceServer(WrapCustomerGroupRelationshipService(customerGroupRelationshipService.Clone)),
		service.NewInventoryServiceServer(WrapInventoryService(inventoryService.Clone)),
		service.NewVariantServiceServer(WrapVariantService(variantService.Clone)),
		service.NewProductServiceServer(WrapProductService(productService.Clone)),
		service.NewProductCollectionServiceServer(WrapProductCollectionService(productCollectionService.Clone)),
		service.NewProductCollectionRelationshipServiceServer(WrapProductCollectionRelationshipService(productCollectionRelationshipService.Clone)),
	}
	return servers
}
