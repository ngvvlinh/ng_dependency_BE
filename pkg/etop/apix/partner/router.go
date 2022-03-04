package partner

import (
	"net/url"

	"o.o/backend/pkg/common/apifw/idemp"
	"o.o/backend/pkg/common/authorization/auth"
	"o.o/backend/pkg/common/redis"
	"o.o/capi/httprpc"
	"o.o/common/l"
)

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
	internalService *InternalService,
) (Servers, func()) {
	authStore = auth.NewGenerator(rd)
	authURL = string(_authURL)
	if authURL == "" {
		ll.Panic("no auth_url")
	}
	if _, err := url.Parse(authURL); err != nil {
		ll.Panic("invalid auth_url", l.String("url", authURL))
	}

	idempgroup = idemp.NewRedisGroup(rd, PrefixIdempPartnerAPI, 0)
	servers := httprpc.MustNewServers(
		miscService.Clone,
		shopService.Clone,
		webhookService.Clone,
		historyService.Clone,
		shippingService.Clone,
		orderService.Clone,
		fulfillmentService.Clone,
		customerService.Clone,
		customerAddressService.Clone,
		customerGroupService.Clone,
		customerGroupRelationshipService.Clone,
		inventoryService.Clone,
		variantService.Clone,
		productService.Clone,
		productCollectionService.Clone,
		productCollectionRelationshipService.Clone,
		internalService.Clone,
	)
	return servers, idempgroup.Shutdown
}
