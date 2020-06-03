package acl

import (
	"strings"

	"o.o/backend/pkg/etop/authorize/permission"
)

var ACL map[string]*permission.Decl

func init() {
	ACL = make(map[string]*permission.Decl)
	for key, p := range _acl {
		key2 := ConvertKey(key)
		ACL[key2] = p
	}
}

// GetACL is a slightly different function which append / to key to make it
// compatible with authroize/session
//
// TODO(vu): handle ext/
func GetACL() map[string]*permission.Decl {
	res := make(map[string]*permission.Decl)
	for key, p := range _acl {
		key2 := "/" + strings.TrimSpace(key)
		res[key2] = p
	}
	return res
}

func GetExtACL() map[string]*permission.Decl {
	res := make(map[string]*permission.Decl)
	for key, p := range _acl {
		extPrefix := "ext/"
		if strings.HasPrefix(key, extPrefix) {
			key2 := strings.TrimPrefix(key, extPrefix)
			key3 := "/" + strings.TrimSpace(key2)
			res[key3] = p
		}
	}
	return res
}

func ConvertKey(key string) string {
	key2 := strings.TrimSpace(key)
	idx := strings.LastIndex(key2, "/")
	key2 = key2[:idx] + "Service" + key2[idx:]
	return key2
}

const (
	Public     = permission.Public
	Protected  = permission.Protected
	CurUsr     = permission.CurUsr
	Partner    = permission.Partner
	Shop       = permission.Shop
	Affiliate  = permission.Affiliate
	EtopAdmin  = permission.EtopAdmin
	SuperAdmin = permission.SuperAdmin
	Secret     = permission.Secret

	User                 = permission.User
	APIKey               = permission.APIKey
	APIPartnerCarrierKey = permission.APIPartnerCarrierKey
	APIPartnerShopKey    = permission.APIPartnerShopKey

	Req = permission.Required
	Opt = permission.Optional

	// Actions

	ShopCarrierCreate permission.ActionType = "shop/carrier:create"
	ShopCarrierUpdate permission.ActionType = "shop/carrier:update"
	ShopCarrierDelete permission.ActionType = "shop/carrier:delete"
	ShopCarrierView   permission.ActionType = "shop/carrier:view"

	ShopCategoryCreate permission.ActionType = "shop/category:create"
	ShopCategoryUpdate permission.ActionType = "shop/category:update"
	ShopCategoryDelete permission.ActionType = "shop/category:delete"
	ShopCategoryView   permission.ActionType = "shop/category:view"

	ShopCollectionCreate permission.ActionType = "shop/collection:create"
	ShopCollectionView   permission.ActionType = "shop/collection:view"
	ShopCollectionUpdate permission.ActionType = "shop/collection:update"

	ShopCustomerCreate      permission.ActionType = "shop/customer:create"
	ShopCustomerUpdate      permission.ActionType = "shop/customer:update"
	ShopCustomerDelete      permission.ActionType = "shop/customer:delete"
	ShopCustomerView        permission.ActionType = "shop/customer:view"
	ShopCustomerManage      permission.ActionType = "shop/customer:manage"
	ShopCustomerGroupManage permission.ActionType = "shop/customer_group:manage"

	ShopDashboardView permission.ActionType = "shop/dashboard:view"

	ShopExternalAccountManage permission.ActionType = "shop/external_account:manage"

	ShopFulfillmentCreate permission.ActionType = "shop/fulfillment:create"
	ShopFulfillmentCancel permission.ActionType = "shop/fulfillment:cancel"
	ShopFulfillmentView   permission.ActionType = "shop/fulfillment:view"
	ShopFulfillmentExport permission.ActionType = "shop/fulfillment:export"

	ShopInventoryView    permission.ActionType = "shop/inventory:view"
	ShopInventoryCreate  permission.ActionType = "shop/inventory:create"
	ShopInventoryUpdate  permission.ActionType = "shop/inventory:update"
	ShopInventoryConfirm permission.ActionType = "shop/inventory:confirm"
	ShopInventoryCancel  permission.ActionType = "shop/inventory:cancel"

	ShopLedgerCreate permission.ActionType = "shop/ledger:create"
	ShopLedgerView   permission.ActionType = "shop/ledger:view"
	ShopLedgerUpdate permission.ActionType = "shop/ledger:update"
	ShopLedgerDelete permission.ActionType = "shop/ledger:delete"

	ShopMoneyTransactionView   permission.ActionType = "shop/money_transaction:view"
	ShopMoneyTransactionExport permission.ActionType = "shop/money_transaction:export"

	ShopOrderCreate   permission.ActionType = "shop/order:create"
	ShopOrderConfirm  permission.ActionType = "shop/order:confirm"
	ShopOrderComplete permission.ActionType = "shop/order:complete"
	ShopOrderUpdate   permission.ActionType = "shop/order:update"
	ShopOrderCancel   permission.ActionType = "shop/order:cancel"
	ShopOrderView     permission.ActionType = "shop/order:view"
	ShopOrderImport   permission.ActionType = "shop/order:import"
	ShopOrderExport   permission.ActionType = "shop/order:export"

	ShopPurchaseOrderCreate  permission.ActionType = "shop/purchase_order:create"
	ShopPurchaseOrderConfirm permission.ActionType = "shop/purchase_order:confirm"
	ShopPurchaseOrderUpdate  permission.ActionType = "shop/purchase_order:update"
	ShopPurchaseOrderCancel  permission.ActionType = "shop/purchase_order:cancel"
	ShopPurchaseOrderView    permission.ActionType = "shop/purchase_order:view"

	ShopProductCreate          permission.ActionType = "shop/product:create"
	ShopProductDelete          permission.ActionType = "shop/product:delete"
	ShopProductImport          permission.ActionType = "shop/product:import"
	ShopProductCostPriceView   permission.ActionType = "shop/product/cost_price:view"
	ShopProductCostPriceUpdate permission.ActionType = "shop/product/cost_price:update"
	ShopProductBasicInfoUpdate permission.ActionType = "shop/product/basic_info:update"
	ShopProductBasicInfoView   permission.ActionType = "shop/product/basic_info:view"

	ShopReceiptView    permission.ActionType = "shop/receipt:view"
	ShopReceiptCreate  permission.ActionType = "shop/receipt:create"
	ShopReceiptUpdate  permission.ActionType = "shop/receipt:update"
	ShopReceiptConfirm permission.ActionType = "shop/receipt:confirm"
	ShopReceiptCancel  permission.ActionType = "shop/receipt:cancel"

	ShopSettingsShopInfoUpdate        permission.ActionType = "shop/settings/shop_info:update"
	ShopSettingsShopInfoView          permission.ActionType = "shop/settings/shop_info:view"
	ShopSettingsCompanyInfoView       permission.ActionType = "shop/settings/company_info:view"
	ShopSettingsCompanyInfoUpdate     permission.ActionType = "shop/settings/company_info:update"
	ShopSettingsBankInfoView          permission.ActionType = "shop/settings/bank_info:view"
	ShopSettingsBankInfoUpdate        permission.ActionType = "shop/settings/bank_info:update"
	ShopSettingsPaymentScheduleView   permission.ActionType = "shop/settings/payment_schedule:view"
	ShopSettingsPaymentScheduleUpdate permission.ActionType = "shop/settings/payment_schedule:update"
	ShopSettingsTopshipBalanceView    permission.ActionType = "shop/settings/topship_balance:view"
	ShopSettingsWalletBalanceView     permission.ActionType = "shop/settings/wallet_balance:view"
	ShopSettingsShippingSettingView   permission.ActionType = "shop/settings/shipping_setting:view"
	ShopSettingsShippingSettingUpdate permission.ActionType = "shop/settings/shipping_setting:update"

	ShopShipNowView    permission.ActionType = "shop/shipnow:view"
	ShopShipNowCreate  permission.ActionType = "shop/shipnow:create"
	ShopShipNowUpdate  permission.ActionType = "shop/shipnow:update"
	ShopShipNowCancel  permission.ActionType = "shop/shipnow:cancel"
	ShopShipNowConfirm permission.ActionType = "shop/shipnow:confirm"

	ShopBalanceView   permission.ActionType = "shop/balance:view"
	ShopAccountDelete permission.ActionType = "shop/account:delete"

	UserBalanceView permission.ActionType = "shop/user/balance:view"

	RelationshipInvitationCreate   permission.ActionType = "relationship/invitation:create"
	RelationshipInvitationView     permission.ActionType = "relationship/invitation:view"
	RelationshipInvitationDelete   permission.ActionType = "relationship/invitation:delete"
	RelationshipPermissionUpdate   permission.ActionType = "relationship/permission:update"
	RelationshipRelationshipUpdate permission.ActionType = "relationship/relationship:update"
	RelationshipRelationshipView   permission.ActionType = "relationship/relationship:view"
	RelationshipRelationshipRemove permission.ActionType = "relationship/relationship:remove"

	ShopStocktakeCreate     permission.ActionType = "shop/stocktake:create"
	ShopStocktakeUpdate     permission.ActionType = "shop/stocktake:update"
	ShopStocktakeConfirm    permission.ActionType = "shop/stocktake:confirm"
	ShopStocktakeCancel     permission.ActionType = "shop/stocktake:cancel"
	ShopStocktakeView       permission.ActionType = "shop/stocktake:view"
	ShopStocktakeSelfUpdate permission.ActionType = "shop/stocktake:self_update"

	ShopRefundCreate  permission.ActionType = "shop/refund:create"
	ShopRefundUpdate  permission.ActionType = "shop/refund:update"
	ShopRefundConfirm permission.ActionType = "shop/refund:confirm"
	ShopRefundCancel  permission.ActionType = "shop/refund:cancel"
	ShopRefundView    permission.ActionType = "shop/refund:view"

	ShopPurchaseRefundCreate  permission.ActionType = "shop/purchaserefund:create"
	ShopPurchaseRefundUpdate  permission.ActionType = "shop/purchaserefund:update"
	ShopPurchaseRefundConfirm permission.ActionType = "shop/purchaserefund:confirm"
	ShopPurchaseRefundCancel  permission.ActionType = "shop/purchaserefund:cancel"
	ShopPurchaseRefundView    permission.ActionType = "shop/purchaserefund:view"

	ShopSupplierCreate permission.ActionType = "shop/supplier:create"
	ShopSupplierUpdate permission.ActionType = "shop/supplier:update"
	ShopSupplierDelete permission.ActionType = "shop/supplier:delete"
	ShopSupplierView   permission.ActionType = "shop/supplier:view"

	TradingOrderView       permission.ActionType = "trading/order:view"
	TradingOrderCreate     permission.ActionType = "trading/order:create"
	ShopTradingProductView permission.ActionType = "shop/trading/product:view"
	ShopTradingOrderCreate permission.ActionType = "shop/trading/order:create"
	ShopTradingOrderView   permission.ActionType = "shop/trading/order:view"

	ShopPaymentCreate permission.ActionType = "shop/payment:create"
	ShopPaymentView   permission.ActionType = "shop/payment:view"

	WsWebsiteCreate permission.ActionType = "shop/webserver/wswebsite:create"
	WsWebsiteUpdate permission.ActionType = "shop/webserver/wswebsite:update"
	WsWebsiteView   permission.ActionType = "shop/webserver/wswebsite:view"

	WsProductCreate permission.ActionType = "shop/webserver/wsproduct:create"
	WsProductUpdate permission.ActionType = "shop/webserver/wsproduct:update"
	WsProductView   permission.ActionType = "shop/webserver/wsproduct:view"

	WsCategoryCreate permission.ActionType = "shop/webserver/wscategory:create"
	WsCategoryUpdate permission.ActionType = "shop/webserver/wscategory:update"
	WsCategoryView   permission.ActionType = "shop/webserver/wscategory:view"

	WsPageCreate permission.ActionType = "shop/webserver/wspage:create"
	WsPageUpdate permission.ActionType = "shop/webserver/wspage:update"
	WsPageDelete permission.ActionType = "shop/webserver/wspage:delete"
	WsPageView   permission.ActionType = "shop/webserver/wspage:view"
)

// ACL declares access control list
var _acl = map[string]*permission.Decl{
	//-- sadmin --//

	"sadmin.User/CreateUser":     {Type: SuperAdmin},
	"sadmin.User/ResetPassword":  {Type: SuperAdmin},
	"sadmin.User/LoginAsAccount": {Type: SuperAdmin},
	"sadmin.Misc/VersionInfo":    {Type: SuperAdmin},

	//-- common --//

	"       etop.Misc/VersionInfo": {Type: Public},
	"      admin.Misc/VersionInfo": {Type: Public},
	"       shop.Misc/VersionInfo": {Type: Public},
	"ext/partner.Misc/VersionInfo": {Type: Public},
	"   ext/shop.Misc/VersionInfo": {Type: Partner, Auth: APIKey},
	"integration.Misc/VersionInfo": {Type: Public},

	"admin.Misc/AdminLoginAsAccount": {Type: EtopAdmin},

	"etop.User/Register":                 {Type: Public},
	"etop.User/RegisterUsingToken":       {Type: Public},
	"etop.User/Login":                    {Type: Public},
	"etop.User/ResetPassword":            {Type: Public, Captcha: "custom"},
	"etop.User/ChangePasswordUsingToken": {Type: Public},
	"etop.User/ChangePassword":           {Type: CurUsr},
	"etop.User/InitSession":              {Type: Public},
	"etop.User/CheckUserRegistration":    {Type: Public, Captcha: "1"},
	"etop.User/SessionInfo":              {Type: CurUsr},
	"etop.User/SwitchAccount":            {Type: CurUsr},
	"etop.User/SendSTokenEmail":          {Type: CurUsr},
	"etop.User/UpgradeAccessToken":       {Type: CurUsr},
	"etop.User/UpdateUserEmail":          {Type: CurUsr},
	"etop.User/UpdateUserPhone":          {Type: CurUsr},

	"etop.User/UpdatePermission": {Type: CurUsr},

	"etop.User/SendEmailVerification":              {Type: CurUsr},
	"etop.User/SendEmailVerificationUsingOTP":      {Type: CurUsr},
	"etop.User/SendPhoneVerification":              {Type: Public},
	"etop.User/VerifyEmailUsingToken":              {Type: CurUsr},
	"etop.User/VerifyEmailUsingOTP":                {Type: CurUsr},
	"etop.User/VerifyPhoneUsingToken":              {Type: Public},
	"etop.User/VerifyPhoneResetPasswordUsingToken": {Type: Public},
	"etop.User/UpdateReferenceUser":                {Type: CurUsr},
	"etop.User/UpdateReferenceSale":                {Type: CurUsr},

	"etop.Relationship/InviteUserToAccount":          {Type: CurUsr},
	"etop.Relationship/AnswerInvitation":             {Type: CurUsr},
	"etop.Relationship/GetUsersInCurrentAccounts":    {Type: CurUsr},
	"etop.Relationship/LeaveAccount":                 {Type: CurUsr},
	"etop.Relationship/RemoveUserFromCurrentAccount": {Type: CurUsr},

	"etop.Account/UpdateURLSlug":        {Type: CurUsr},
	"etop.Account/GetPublicPartnerInfo": {Type: Protected},
	"etop.Account/GetPublicPartners":    {Type: Protected},

	"etop.Location/GetProvinces":           {Type: Public},
	"etop.Location/GetDistricts":           {Type: Public},
	"etop.Location/GetDistrictsByProvince": {Type: Public},
	"etop.Location/GetWards":               {Type: Public},
	"etop.Location/GetWardsByDistrict":     {Type: Public},
	"etop.Location/ParseLocation":          {Type: Protected},

	"etop.Bank/GetBanks":                  {Type: CurUsr},
	"etop.Bank/GetProvincesByBank":        {Type: CurUsr},
	"etop.Bank/GetBranchesByBankProvince": {Type: CurUsr},

	"etop.Address/CreateAddress": {Type: CurUsr},
	"etop.Address/GetAddresses":  {Type: CurUsr, AuthPartner: Opt},
	"etop.Address/UpdateAddress": {Type: CurUsr},
	"etop.Address/RemoveAddress": {Type: CurUsr},

	"etop.Ecom/SessionInfo": {Type: Public},

	//-- external: partner --//

	"ext/partner.Misc/CurrentAccount":   {Type: Partner, Auth: APIKey},
	"ext/partner.History/GetChanges":    {Type: Partner, Auth: APIKey},
	"ext/partner.Misc/GetLocationList":  {Type: Partner, Auth: APIKey},
	"ext/partner.Shop/AuthorizeShop":    {Type: Partner, Auth: APIKey},
	"ext/partner.Webhook/CreateWebhook": {Type: Partner, Auth: APIKey},
	"ext/partner.Webhook/GetWebhooks":   {Type: Partner, Auth: APIKey},
	"ext/partner.Webhook/DeleteWebhook": {Type: Partner, Auth: APIKey},

	//-- external: partner using partnerShopKey --//
	"ext/partner.Import/Products":           {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.Import/Brands":             {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.Import/Categories":         {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.Import/Customers":          {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.Import/Variants":           {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.Import/Collections":        {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.Import/ProductCollections": {Type: Shop, Auth: APIPartnerShopKey},

	"ext/partner.Shop/CurrentShop":               {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.Shipping/GetShippingServices":   {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.Shipping/CreateAndConfirmOrder": {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.Shipping/CancelOrder":           {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.Shipping/GetOrder":              {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.Shipping/GetFulfillment":        {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.Shipping/ListFulfillments":      {Type: Shop, Auth: APIPartnerShopKey},

	"ext/partner.Order/CreateAndConfirmOrder": {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.Order/CreateOrder":           {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.Order/ConfirmOrder":          {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.Order/CancelOrder":           {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.Order/GetOrder":              {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.Order/ListOrders":            {Type: Shop, Auth: APIPartnerShopKey},

	"ext/partner.Customer/GetCustomer":    {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.Customer/ListCustomers":  {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.Customer/CreateCustomer": {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.Customer/UpdateCustomer": {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.Customer/DeleteCustomer": {Type: Shop, Auth: APIPartnerShopKey},

	"ext/partner.CustomerAddress/GetAddress":    {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.CustomerAddress/ListAddresses": {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.CustomerAddress/CreateAddress": {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.CustomerAddress/UpdateAddress": {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.CustomerAddress/DeleteAddress": {Type: Shop, Auth: APIPartnerShopKey},

	"ext/partner.CustomerGroup/CreateGroup": {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.CustomerGroup/GetGroup":    {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.CustomerGroup/ListGroups":  {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.CustomerGroup/UpdateGroup": {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.CustomerGroup/DeleteGroup": {Type: Shop, Auth: APIPartnerShopKey},

	"ext/partner.CustomerGroupRelationship/ListRelationships":  {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.CustomerGroupRelationship/CreateRelationship": {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.CustomerGroupRelationship/DeleteRelationship": {Type: Shop, Auth: APIPartnerShopKey},

	"ext/partner.Fulfillment/GetFulfillment":    {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.Fulfillment/ListFulfillments":  {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.Fulfillment/CreateFulfillment": {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.Fulfillment/CancelFulfillment": {Type: Shop, Auth: APIPartnerShopKey},

	"ext/partner.Inventory/ListInventoryLevels": {Type: Shop, Auth: APIPartnerShopKey},

	"ext/partner.Product/GetProduct":    {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.Product/ListProducts":  {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.Product/CreateProduct": {Type: Shop, Auth: APIPartnerShopKey},
	// (status, tag, images, metaFields, category)
	"ext/partner.Product/UpdateProduct": {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.Product/DeleteProduct": {Type: Shop, Auth: APIPartnerShopKey},
	// (tags, status)
	"ext/partner.ProductCollection/CreateCollection": {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.ProductCollection/UpdateCollection": {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.ProductCollection/DeleteCollection": {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.ProductCollection/GetCollection":    {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.ProductCollection/ListCollections":  {Type: Shop, Auth: APIPartnerShopKey},

	"ext/partner.ProductCollectionRelationship/ListRelationships":  {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.ProductCollectionRelationship/CreateRelationship": {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.ProductCollectionRelationship/DeleteRelationship": {Type: Shop, Auth: APIPartnerShopKey},

	"ext/partner.Variant/GetVariant":    {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.Variant/DeleteVariant": {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.Variant/ListVariants":  {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.Variant/CreateVariant": {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.Variant/UpdateVariant": {Type: Shop, Auth: APIPartnerShopKey},

	//-- external: carrier -- //
	"ext/carrier.Misc/GetLocationList":                {Type: Partner, Auth: APIPartnerCarrierKey},
	"ext/carrier.Misc/CurrentAccount":                 {Type: Partner, Auth: APIPartnerCarrierKey},
	"ext/carrier.ShipmentConnection/GetConnections":   {Type: Partner, Auth: APIPartnerCarrierKey},
	"ext/carrier.ShipmentConnection/CreateConnection": {Type: Partner, Auth: APIPartnerCarrierKey},
	"ext/carrier.ShipmentConnection/UpdateConnection": {Type: Partner, Auth: APIPartnerCarrierKey},
	"ext/carrier.ShipmentConnection/DeleteConnection": {Type: Partner, Auth: APIPartnerCarrierKey},
	"ext/carrier.Shipment/UpdateFulfillment":          {Type: Partner, Auth: APIPartnerCarrierKey},

	//-- external: shop --//

	"ext/shop.Misc/CurrentAccount":            {Type: Shop, Auth: APIKey},
	"ext/shop.Misc/GetLocationList":           {Type: Shop, Auth: APIKey},
	"ext/shop.History/GetChanges":             {Type: Shop, Auth: APIKey},
	"ext/shop.Webhook/CreateWebhook":          {Type: Shop, Auth: APIKey},
	"ext/shop.Webhook/GetWebhooks":            {Type: Shop, Auth: APIKey},
	"ext/shop.Webhook/DeleteWebhook":          {Type: Shop, Auth: APIKey},
	"ext/shop.Shipping/GetShippingServices":   {Type: Shop, Auth: APIKey},
	"ext/shop.Shipping/CreateAndConfirmOrder": {Type: Shop, Auth: APIKey},
	"ext/shop.Shipping/CancelOrder":           {Type: Shop, Auth: APIKey},
	"ext/shop.Shipping/GetOrder":              {Type: Shop, Auth: APIKey},
	"ext/shop.Shipping/GetFulfillment":        {Type: Shop, Auth: APIKey},

	"ext/shop.Order/CreateOrder":  {Type: Shop, Auth: APIKey},
	"ext/shop.Order/ConfirmOrder": {Type: Shop, Auth: APIKey},
	"ext/shop.Order/CancelOrder":  {Type: Shop, Auth: APIKey},
	"ext/shop.Order/GetOrder":     {Type: Shop, Auth: APIKey},
	"ext/shop.Order/ListOrders":   {Type: Shop, Auth: APIKey},

	"ext/shop.Customer/GetCustomer":    {Type: Shop, Auth: APIKey},
	"ext/shop.Customer/ListCustomers":  {Type: Shop, Auth: APIKey},
	"ext/shop.Customer/CreateCustomer": {Type: Shop, Auth: APIKey},
	"ext/shop.Customer/UpdateCustomer": {Type: Shop, Auth: APIKey},
	"ext/shop.Customer/DeleteCustomer": {Type: Shop, Auth: APIKey},

	"ext/shop.CustomerAddress/GetAddress":    {Type: Shop, Auth: APIKey},
	"ext/shop.CustomerAddress/ListAddresses": {Type: Shop, Auth: APIKey},
	"ext/shop.CustomerAddress/CreateAddress": {Type: Shop, Auth: APIKey},
	"ext/shop.CustomerAddress/UpdateAddress": {Type: Shop, Auth: APIKey},
	"ext/shop.CustomerAddress/DeleteAddress": {Type: Shop, Auth: APIKey},

	"ext/shop.CustomerGroup/CreateGroup": {Type: Shop, Auth: APIKey},
	"ext/shop.CustomerGroup/GetGroup":    {Type: Shop, Auth: APIKey},
	"ext/shop.CustomerGroup/ListGroups":  {Type: Shop, Auth: APIKey},
	"ext/shop.CustomerGroup/UpdateGroup": {Type: Shop, Auth: APIKey},
	"ext/shop.CustomerGroup/DeleteGroup": {Type: Shop, Auth: APIKey},

	"ext/shop.CustomerGroupRelationship/ListRelationships":  {Type: Shop, Auth: APIKey},
	"ext/shop.CustomerGroupRelationship/CreateRelationship": {Type: Shop, Auth: APIKey},
	"ext/shop.CustomerGroupRelationship/DeleteRelationship": {Type: Shop, Auth: APIKey},

	"ext/shop.Fulfillment/GetFulfillment":   {Type: Shop, Auth: APIKey},
	"ext/shop.Fulfillment/ListFulfillments": {Type: Shop, Auth: APIKey},

	"ext/shop.Inventory/ListInventoryLevels": {Type: Shop, Auth: APIKey},

	"ext/shop.Product/GetProduct":    {Type: Shop, Auth: APIKey},
	"ext/shop.Product/ListProducts":  {Type: Shop, Auth: APIKey},
	"ext/shop.Product/CreateProduct": {Type: Shop, Auth: APIKey},
	// (status, tag, images, metaFields, category)
	"ext/shop.Product/UpdateProduct": {Type: Shop, Auth: APIKey},
	"ext/shop.Product/DeleteProduct": {Type: Shop, Auth: APIKey},
	// (tags, status)
	"ext/shop.ProductCollection/CreateCollection": {Type: Shop, Auth: APIKey},
	"ext/shop.ProductCollection/UpdateCollection": {Type: Shop, Auth: APIKey},
	"ext/shop.ProductCollection/DeleteCollection": {Type: Shop, Auth: APIKey},
	"ext/shop.ProductCollection/GetCollection":    {Type: Shop, Auth: APIKey},
	"ext/shop.ProductCollection/ListCollections":  {Type: Shop, Auth: APIKey},

	"ext/shop.ProductCollectionRelationship/ListRelationships":  {Type: Shop, Auth: APIKey},
	"ext/shop.ProductCollectionRelationship/CreateRelationship": {Type: Shop, Auth: APIKey},
	"ext/shop.ProductCollectionRelationship/DeleteRelationship": {Type: Shop, Auth: APIKey},

	"ext/shop.Variant/GetVariant":    {Type: Shop, Auth: APIKey},
	"ext/shop.Variant/DeleteVariant": {Type: Shop, Auth: APIKey},
	"ext/shop.Variant/ListVariants":  {Type: Shop, Auth: APIKey},
	"ext/shop.Variant/CreateVariant": {Type: Shop, Auth: APIKey},
	"ext/shop.Variant/UpdateVariant": {Type: Shop, Auth: APIKey},

	//-- integration --//

	"integration.Integration/Init":              {Type: Public},
	"integration.Integration/RequestLogin":      {Type: Protected, AuthPartner: Req, Captcha: "1"},
	"integration.Integration/LoginUsingToken":   {Type: Protected, AuthPartner: Req},
	"integration.Integration/LoginUsingTokenWL": {Type: Protected, AuthPartner: Req},
	"integration.Integration/Register":          {Type: Protected, AuthPartner: Req},
	"integration.Integration/GrantAccess":       {Type: CurUsr, AuthPartner: Req},
	"integration.Integration/SessionInfo":       {Type: Protected, AuthPartner: Req},

	//-- admin --//

	"admin.Account/CreatePartner":  {Type: EtopAdmin},
	"admin.Account/GenerateAPIKey": {Type: EtopAdmin},

	"admin.Category/CreateCategory":         {Type: EtopAdmin},
	"admin.Category/GetCategories":          {Type: EtopAdmin},
	"admin.Category/UpdateProductsCategory": {Type: EtopAdmin},
	"admin.Category/RemoveProductsCategory": {Type: EtopAdmin},

	"admin.Product/GetVariant":           {Type: EtopAdmin},
	"admin.Product/GetVariantsByIDs":     {Type: EtopAdmin},
	"admin.Product/UpdateVariantsStatus": {Type: EtopAdmin},
	"admin.Product/GetProduct":           {Type: EtopAdmin},
	"admin.Product/GetProducts":          {Type: EtopAdmin},
	"admin.Product/GetProductsByIDs":     {Type: EtopAdmin},
	"admin.Product/UpdateProductsStatus": {Type: EtopAdmin},
	"admin.Product/UpdateProduct":        {Type: EtopAdmin},
	"admin.Product/UpdateVariant":        {Type: EtopAdmin},
	"admin.Product/UpdateVariantImages":  {Type: EtopAdmin},
	"admin.Product/UpdateProductImages":  {Type: EtopAdmin},

	"admin.Order/GetOrder":       {Type: EtopAdmin},
	"admin.Order/GetOrders":      {Type: EtopAdmin},
	"admin.Order/GetOrdersByIDs": {Type: EtopAdmin},

	"admin.User/GetUsers":      {Type: EtopAdmin},
	"admin.User/GetUser":       {Type: EtopAdmin},
	"admin.User/GetUsersByIDs": {Type: EtopAdmin},

	"admin.Fulfillment/GetFulfillment":                 {Type: EtopAdmin},
	"admin.Fulfillment/GetFulfillments":                {Type: EtopAdmin},
	"admin.Fulfillment/UpdateFulfillment":              {Type: EtopAdmin},
	"admin.Fulfillment/UpdateFulfillmentInfo":          {Type: EtopAdmin},
	"admin.Fulfillment/UpdateFulfillmentShippingState": {Type: EtopAdmin},
	"admin.Fulfillment/UpdateFulfillmentShippingFee":   {Type: EtopAdmin},

	"admin.MoneyTransaction/GetMoneyTransaction":                         {Type: EtopAdmin},
	"admin.MoneyTransaction/GetMoneyTransactions":                        {Type: EtopAdmin},
	"admin.MoneyTransaction/ConfirmMoneyTransaction":                     {Type: EtopAdmin},
	"admin.MoneyTransaction/UpdateMoneyTransaction":                      {Type: EtopAdmin},
	"admin.MoneyTransaction/GetMoneyTransactionShippingExternal":         {Type: EtopAdmin},
	"admin.MoneyTransaction/GetMoneyTransactionShippingExternals":        {Type: EtopAdmin},
	"admin.MoneyTransaction/RemoveMoneyTransactionShippingExternalLines": {Type: EtopAdmin},
	"admin.MoneyTransaction/DeleteMoneyTransactionShippingExternal":      {Type: EtopAdmin},
	"admin.MoneyTransaction/ConfirmMoneyTransactionShippingExternals":    {Type: EtopAdmin},
	"admin.MoneyTransaction/UpdateMoneyTransactionShippingExternal":      {Type: EtopAdmin},
	"admin.MoneyTransaction/GetMoneyTransactionShippingEtop":             {Type: EtopAdmin},
	"admin.MoneyTransaction/GetMoneyTransactionShippingEtops":            {Type: EtopAdmin},
	"admin.MoneyTransaction/CreateMoneyTransactionShippingEtop":          {Type: EtopAdmin},
	"admin.MoneyTransaction/UpdateMoneyTransactionShippingEtop":          {Type: EtopAdmin},
	"admin.MoneyTransaction/DeleteMoneyTransactionShippingEtop":          {Type: EtopAdmin},
	"admin.MoneyTransaction/ConfirmMoneyTransactionShippingEtop":         {Type: EtopAdmin},

	"admin.Shop/GetShop":       {Type: EtopAdmin},
	"admin.Shop/GetShops":      {Type: EtopAdmin},
	"admin.Shop/GetShopsByIDs": {Type: EtopAdmin},

	"admin.Credit/GetCredit":     {Type: EtopAdmin},
	"admin.Credit/GetCredits":    {Type: EtopAdmin},
	"admin.Credit/CreateCredit":  {Type: EtopAdmin},
	"admin.Credit/UpdateCredit":  {Type: EtopAdmin},
	"admin.Credit/ConfirmCredit": {Type: EtopAdmin},
	"admin.Credit/DeleteCredit":  {Type: EtopAdmin},

	"admin.Notification/CreateNotifications": {Type: EtopAdmin},

	"admin.Connection/GetConnections":              {Type: EtopAdmin},
	"admin.Connection/ConfirmConnection":           {Type: EtopAdmin},
	"admin.Connection/DisableConnection":           {Type: EtopAdmin},
	"admin.Connection/CreateBuiltinConnection":     {Type: EtopAdmin},
	"admin.Connection/GetBuiltinShopConnections":   {Type: EtopAdmin},
	"admin.Connection/UpdateBuiltinShopConnection": {Type: EtopAdmin},
	"admin.Connection/GetConnectionServices":       {Type: EtopAdmin},

	"admin.ShipmentPrice/GetShipmentServices":                      {Type: EtopAdmin},
	"admin.ShipmentPrice/GetShipmentService":                       {Type: EtopAdmin},
	"admin.ShipmentPrice/CreateShipmentService":                    {Type: EtopAdmin},
	"admin.ShipmentPrice/UpdateShipmentService":                    {Type: EtopAdmin},
	"admin.ShipmentPrice/DeleteShipmentService":                    {Type: EtopAdmin},
	"admin.ShipmentPrice/UpdateShipmentServicesAvailableLocations": {Type: EtopAdmin},
	"admin.ShipmentPrice/UpdateShipmentServicesBlacklistLocations": {Type: EtopAdmin},

	"admin.ShipmentPrice/GetShipmentPriceLists":       {Type: EtopAdmin},
	"admin.ShipmentPrice/GetShipmentPriceList":        {Type: EtopAdmin},
	"admin.ShipmentPrice/CreateShipmentPriceList":     {Type: EtopAdmin},
	"admin.ShipmentPrice/UpdateShipmentPriceList":     {Type: EtopAdmin},
	"admin.ShipmentPrice/SetDefaultShipmentPriceList": {Type: EtopAdmin},
	"admin.ShipmentPrice/DeleteShipmentPriceList":     {Type: EtopAdmin},

	"admin.ShipmentPrice/GetShipmentPrice":                  {Type: EtopAdmin},
	"admin.ShipmentPrice/GetShipmentPrices":                 {Type: EtopAdmin},
	"admin.ShipmentPrice/CreateShipmentPrice":               {Type: EtopAdmin},
	"admin.ShipmentPrice/UpdateShipmentPrice":               {Type: EtopAdmin},
	"admin.ShipmentPrice/DeleteShipmentPrice":               {Type: EtopAdmin},
	"admin.ShipmentPrice/UpdateShipmentPricesPriorityPoint": {Type: EtopAdmin},

	"admin.ShipmentPrice/GetShopShipmentPriceLists":   {Type: EtopAdmin},
	"admin.ShipmentPrice/GetShopShipmentPriceList":    {Type: EtopAdmin},
	"admin.ShipmentPrice/CreateShopShipmentPriceList": {Type: EtopAdmin},
	"admin.ShipmentPrice/UpdateShopShipmentPriceList": {Type: EtopAdmin},
	"admin.ShipmentPrice/DeleteShopShipmentPriceList": {Type: EtopAdmin},

	"admin.ShipmentPrice/GetShippingServices": {Type: EtopAdmin},

	"admin.Location/GetCustomRegion":    {Type: EtopAdmin},
	"admin.Location/GetCustomRegions":   {Type: EtopAdmin},
	"admin.Location/CreateCustomRegion": {Type: EtopAdmin},
	"admin.Location/UpdateCustomRegion": {Type: EtopAdmin},
	"admin.Location/DeleteCustomRegion": {Type: EtopAdmin},

	"admin.Subscription/CreateSubscriptionProduct": {Type: EtopAdmin},
	"admin.Subscription/GetSubscriptionProducts":   {Type: EtopAdmin},
	"admin.Subscription/DeleteSubscriptionProduct": {Type: EtopAdmin},
	"admin.Subscription/CreateSubscriptionPlan":    {Type: EtopAdmin},
	"admin.Subscription/UpdateSubscriptionPlan":    {Type: EtopAdmin},
	"admin.Subscription/GetSubscriptionPlans":      {Type: EtopAdmin},
	"admin.Subscription/DeleteSubscriptionPlan":    {Type: EtopAdmin},

	"admin.Subscription/GetSubscription":               {Type: EtopAdmin},
	"admin.Subscription/GetSubscriptions":              {Type: EtopAdmin},
	"admin.Subscription/CreateSubscription":            {Type: EtopAdmin},
	"admin.Subscription/UpdateSubscriptionInfo":        {Type: EtopAdmin},
	"admin.Subscription/CancelSubscription":            {Type: EtopAdmin},
	"admin.Subscription/ActivateSubscription":          {Type: EtopAdmin},
	"admin.Subscription/DeleteSubscription":            {Type: EtopAdmin},
	"admin.Subscription/GetSubscriptionBills":          {Type: EtopAdmin},
	"admin.Subscription/CreateSubscriptionBill":        {Type: EtopAdmin},
	"admin.Subscription/ManualPaymentSubscriptionBill": {Type: EtopAdmin},
	"admin.Subscription/DeleteSubscriptionBill":        {Type: EtopAdmin},
	//-- shop --//

	"shop.Account/RegisterShop": {Type: CurUsr, AuthPartner: Opt},

	// permission: owner

	"shop.Account/UpdateShop": {Type: Shop, Actions: actions(ShopSettingsShopInfoUpdate)},
	"shop.Account/DeleteShop": {Type: Shop, Actions: actions(ShopAccountDelete)},

	// permission: admin

	"shop.Account/SetDefaultAddress": {Type: Shop, Actions: actions(ShopSettingsShopInfoUpdate)},
	"shop.Account/GetBalanceShop":    {Type: Shop, Actions: actions(ShopBalanceView)},

	"shop.Account/CreateExternalAccountAhamove":                   {Type: Shop, Actions: actions(ShopExternalAccountManage)},
	"shop.Account/GetExternalAccountAhamove":                      {Type: Shop},
	"shop.Account/RequestVerifyExternalAccountAhamove":            {Type: Shop, Actions: actions(ShopExternalAccountManage)},
	"shop.Account/UpdateExternalAccountAhamoveVerification":       {Type: Shop, Actions: actions(ShopExternalAccountManage)},
	"shop.Account/UpdateExternalAccountAhamoveVerificationImages": {Type: Shop, Actions: actions(ShopExternalAccountManage)},

	"shop.Browse/BrowseCategories":    {Type: Shop},
	"shop.Browse/BrowseProduct":       {Type: Shop},
	"shop.Browse/BrowseVariant":       {Type: Shop},
	"shop.Browse/BrowseProducts":      {Type: Shop},
	"shop.Browse/BrowseVariants":      {Type: Shop},
	"shop.Browse/BrowseProductsByIDs": {Type: Shop},
	"shop.Browse/BrowseVariantsByIDs": {Type: Shop},

	"shop.Collection/GetCollection":             {Type: Shop, Actions: actions(ShopCollectionView)},
	"shop.Collection/GetCollections":            {Type: Shop, Actions: actions(ShopCollectionView)},
	"shop.Collection/CreateCollection":          {Type: Shop, Actions: actions(ShopCollectionCreate)},
	"shop.Collection/UpdateCollection":          {Type: Shop, Actions: actions(ShopCollectionUpdate)},
	"shop.Collection/GetCollectionsByProductID": {Type: Shop, Actions: actions(ShopCollectionView)},

	"shop.Customer/CreateCustomer":          {Type: Shop, Actions: actions(ShopCustomerCreate)},
	"shop.Customer/UpdateCustomer":          {Type: Shop, Actions: actions(ShopCustomerUpdate)},
	"shop.Customer/DeleteCustomer":          {Type: Shop, Actions: actions(ShopCustomerDelete)},
	"shop.Customer/GetCustomer":             {Type: Shop, Actions: actions(ShopCustomerView)},
	"shop.Customer/GetCustomerDetails":      {Type: Shop, Actions: actions(ShopCustomerView)},
	"shop.Customer/GetCustomers":            {Type: Shop, Actions: actions(ShopCustomerView)},
	"shop.Customer/GetCustomersByIDs":       {Type: Shop, Actions: actions(ShopCustomerView)},
	"shop.Customer/BatchSetCustomersStatus": {Type: Shop, Actions: actions(ShopCustomerUpdate)},

	"shop.Customer/GetCustomerAddresses":      {Type: Shop, Actions: actions(ShopCustomerView)},
	"shop.Customer/CreateCustomerAddress":     {Type: Shop, Actions: actions(ShopCustomerCreate, ShopCustomerUpdate)},
	"shop.Customer/UpdateCustomerAddress":     {Type: Shop, Actions: actions(ShopCustomerUpdate)},
	"shop.Customer/DeleteCustomerAddress":     {Type: Shop, Actions: actions(ShopCustomerUpdate, ShopCustomerDelete)},
	"shop.Customer/SetDefaultCustomerAddress": {Type: Shop, Actions: actions(ShopCustomerUpdate)},

	"shop.Customer/AddCustomersToGroup":      {Type: Shop, Actions: actions(ShopCustomerManage)},
	"shop.Customer/RemoveCustomersFromGroup": {Type: Shop, Actions: actions(ShopCustomerManage)},

	"shop.CustomerGroup/CreateCustomerGroup": {Type: Shop, Actions: actions(ShopCustomerGroupManage)},
	"shop.CustomerGroup/GetCustomerGroup":    {Type: Shop, Actions: actions(ShopCustomerView)},
	"shop.CustomerGroup/GetCustomerGroups":   {Type: Shop, Actions: actions(ShopCustomerView)},
	"shop.CustomerGroup/UpdateCustomerGroup": {Type: Shop, Actions: actions(ShopCustomerGroupManage)},

	"shop.Category/GetCategory":    {Type: Shop, Actions: actions(ShopCategoryView)},
	"shop.Category/GetCategories":  {Type: Shop, Actions: actions(ShopCategoryView)},
	"shop.Category/CreateCategory": {Type: Shop, Actions: actions(ShopCategoryCreate)},
	"shop.Category/UpdateCategory": {Type: Shop, Actions: actions(ShopCategoryUpdate)},
	"shop.Category/DeleteCategory": {Type: Shop, Actions: actions(ShopCategoryDelete)},

	"shop.Product/GetProduct":              {Type: Shop, AuthPartner: Opt, Actions: actions(ShopProductBasicInfoView)},
	"shop.Product/GetProducts":             {Type: Shop, AuthPartner: Opt, Actions: actions(ShopProductBasicInfoView)},
	"shop.Product/GetProductsByIDs":        {Type: Shop, AuthPartner: Opt, Actions: actions(ShopProductBasicInfoView)},
	"shop.Product/RemoveProducts":          {Type: Shop, Actions: actions(ShopProductDelete)},
	"shop.Product/UpdateProduct":           {Type: Shop, Actions: actions(ShopProductBasicInfoUpdate)},
	"shop.Product/UpdateProductsStatus":    {Type: Shop, Actions: actions(ShopProductBasicInfoUpdate)},
	"shop.Product/UpdateProductsTags":      {Type: Shop, Actions: actions(ShopProductBasicInfoUpdate)},
	"shop.Product/UpdateProductImages":     {Type: Shop, Actions: actions(ShopProductBasicInfoUpdate)},
	"shop.Product/UpdateProductMetaFields": {Type: Shop, Actions: actions(ShopProductBasicInfoUpdate)},
	"shop.Product/CreateProduct":           {Type: Shop, Actions: actions(ShopProductCreate)},
	"shop.Product/UpdateProductStatus":     {Type: Shop, Actions: actions(ShopProductBasicInfoUpdate)},
	"shop.Product/UpdateProductCategory":   {Type: Shop, Actions: actions(ShopProductBasicInfoUpdate)},
	"shop.Product/RemoveProductCategory":   {Type: Shop, Actions: actions(ShopProductBasicInfoUpdate, ShopProductDelete)},
	"shop.Product/AddProductCollection":    {Type: Shop, Actions: actions(ShopProductCreate, ShopProductBasicInfoUpdate)},
	"shop.Product/RemoveProductCollection": {Type: Shop, Actions: actions(ShopProductBasicInfoUpdate, ShopProductDelete)},

	"shop.Product/GetVariant":                  {Type: Shop, AuthPartner: Opt, Actions: actions(ShopProductBasicInfoView)},
	"shop.Product/GetVariants":                 {Type: Shop, AuthPartner: Opt, Actions: actions(ShopProductBasicInfoView)},
	"shop.Product/GetVariantsByIDs":            {Type: Shop, AuthPartner: Opt, Actions: actions(ShopProductBasicInfoView)},
	"shop.Product/GetVariantsBySupplierID":     {Type: Shop, Actions: actions(ShopProductBasicInfoView)},
	"shop.Product/CreateVariant":               {Type: Shop, Actions: actions(ShopProductCreate)},
	"shop.Product/AddVariants":                 {Type: Shop, AuthPartner: Opt, Actions: actions(ShopProductCreate, ShopProductBasicInfoUpdate)},
	"shop.Product/RemoveVariants":              {Type: Shop, Actions: actions(ShopProductDelete, ShopProductBasicInfoUpdate)},
	"shop.Product/UpdateVariant":               {Type: Shop, Actions: actions(ShopProductBasicInfoUpdate)},
	"shop.Product/UpdateVariants":              {Type: Shop, Actions: actions(ShopProductBasicInfoUpdate)},
	"shop.Product/UpdateVariantsStatus":        {Type: Shop, Actions: actions(ShopProductBasicInfoUpdate)},
	"shop.Product/UpdateVariantsTags":          {Type: Shop, Actions: actions(ShopProductBasicInfoUpdate)},
	"shop.Product/UpdateVariantImages":         {Type: Shop, Actions: actions(ShopProductBasicInfoUpdate)},
	"shop.Product/UpdateVariantAttributes":     {Type: Shop, Actions: actions(ShopProductBasicInfoUpdate)},
	"shop.ProductSource/GetShopProductSources": {Type: Shop, AuthPartner: Opt, Actions: actions(ShopProductBasicInfoView)},
	"shop.ProductSource/CreateProductSource":   {Type: Shop, Actions: actions(ShopProductCreate, ShopProductBasicInfoUpdate)},
	"shop.ProductSource/ConnectProductSource":  {Type: Shop, Actions: actions(ShopProductCreate, ShopProductBasicInfoUpdate)},

	"shop.ProductSource/GetProductSourceCategories":  {Type: Shop, AuthPartner: Opt, Actions: actions(ShopProductBasicInfoView)},
	"shop.ProductSource/GetProductSourceCategory":    {Type: Shop, AuthPartner: Opt, Actions: actions(ShopProductBasicInfoView)},
	"shop.ProductSource/CreateVariant":               {Type: Shop, Actions: actions(ShopProductCreate, ShopProductBasicInfoUpdate), AuthPartner: Opt, Rename: "DeprecatedCreateVariant"},
	"shop.ProductSource/CreateProductSourceCategory": {Type: Shop, Actions: actions(ShopProductCreate, ShopProductBasicInfoUpdate)},
	"shop.ProductSource/UpdateProductsPSCategory":    {Type: Shop, Actions: actions(ShopProductCreate, ShopProductBasicInfoUpdate)},
	"shop.ProductSource/UpdateProductSourceCategory": {Type: Shop, Actions: actions(ShopProductCreate, ShopProductBasicInfoUpdate)},
	"shop.ProductSource/RemoveProductSourceCategory": {Type: Shop, Actions: actions(ShopProductDelete, ShopProductBasicInfoUpdate)},

	"shop.Price/GetPriceRules":    {Type: Shop},
	"shop.Price/UpdatePriceRules": {Type: Shop},

	"shop.Order/CreateOrder":                        {Type: Shop, AuthPartner: Opt, Actions: actions(ShopOrderCreate)},
	"shop.Order/GetOrder":                           {Type: Shop, AuthPartner: Opt, Actions: actions(ShopOrderView)},
	"shop.Order/GetOrders":                          {Type: Shop, AuthPartner: Opt, Actions: actions(ShopOrderView)},
	"shop.Order/GetOrdersByIDs":                     {Type: Shop, AuthPartner: Opt, Actions: actions(ShopOrderView)},
	"shop.Order/GetOrdersByReceiptID":               {Type: Shop, AuthPartner: Opt, Actions: actions(ShopOrderView)},
	"shop.Order/UpdateOrder":                        {Type: Shop, AuthPartner: Opt, Actions: actions(ShopOrderUpdate)},
	"shop.Order/UpdateOrdersStatus":                 {Type: Shop, AuthPartner: Opt, Actions: actions(ShopOrderUpdate)},
	"shop.Order/ConfirmOrder":                       {Type: Shop, AuthPartner: Opt, Actions: actions(ShopOrderConfirm)},
	"shop.Order/ConfirmOrderAndCreateFulfillments":  {Type: Shop, AuthPartner: Opt, Actions: actions(ShopOrderConfirm)},
	"shop.Order/ConfirmOrdersAndCreateFulfillments": {Type: Shop, AuthPartner: Opt, Actions: actions(ShopOrderConfirm)},
	"shop.Order/CompleteOrder":                      {Type: Shop, AuthPartner: Opt, Actions: actions(ShopOrderComplete)},
	"shop.Order/CancelOrder":                        {Type: Shop, AuthPartner: Opt, Actions: actions(ShopOrderCancel)},
	"shop.Order/UpdateOrderPaymentStatus":           {Type: Shop, AuthPartner: Opt, Actions: actions(ShopOrderUpdate)},
	"shop.Order/UpdateOrderShippingInfo":            {Type: Shop, AuthPartner: Opt, Actions: actions(ShopOrderUpdate)},

	"shop.Fulfillment/GetPublicExternalShippingServices": {Type: Public},
	"shop.Fulfillment/GetPublicFulfillment":              {Type: Public},
	"shop.Fulfillment/GetFulfillmentsByIDs":              {Type: Shop, AuthPartner: Opt, Actions: actions(ShopFulfillmentView)},
	"shop.Fulfillment/GetExternalShippingServices":       {Type: Shop, AuthPartner: Opt, Actions: actions(ShopFulfillmentCreate)},
	"shop.Fulfillment/CancelFulfillment":                 {Type: Shop, AuthPartner: Opt, Actions: actions(ShopFulfillmentCancel)},
	"shop.Fulfillment/CreateFulfillmentsForOrder":        {Type: Shop, Actions: actions(ShopFulfillmentCreate)},
	"shop.Fulfillment/GetFulfillment":                    {Type: Shop, AuthPartner: Opt, Actions: actions(ShopFulfillmentView)},
	"shop.Fulfillment/GetFulfillments":                   {Type: Shop, AuthPartner: Opt, Actions: actions(ShopFulfillmentView)},
	"shop.Fulfillment/UpdateFulfillmentsShippingState":   {Type: Shop},

	"shop.Shipment/GetShippingServices": {Type: Shop, Actions: actions(ShopFulfillmentCreate)},
	"shop.Shipment/CreateFulfillments":  {Type: Shop, Actions: actions(ShopFulfillmentCreate)},
	"shop.Shipment/CancelFulfillment":   {Type: Shop, Actions: actions(ShopFulfillmentCancel)},

	"shop.Shipnow/GetShipnowFulfillment":      {Type: Shop, Actions: actions(ShopShipNowView)},
	"shop.Shipnow/GetShipnowFulfillments":     {Type: Shop, Actions: actions(ShopShipNowView)},
	"shop.Shipnow/CreateShipnowFulfillment":   {Type: Shop, Actions: actions(ShopShipNowCreate)},
	"shop.Shipnow/CreateShipnowFulfillmentV2": {Type: Shop, Actions: actions(ShopShipNowCreate)},
	"shop.Shipnow/ConfirmShipnowFulfillment":  {Type: Shop, Actions: actions(ShopShipNowConfirm)},
	"shop.Shipnow/UpdateShipnowFulfillment":   {Type: Shop, Actions: actions(ShopShipNowUpdate)},
	"shop.Shipnow/CancelShipnowFulfillment":   {Type: Shop, Actions: actions(ShopShipNowCancel)},
	"shop.Shipnow/GetShipnowServices":         {Type: Shop, Actions: actions(ShopShipNowView)},

	"shop.Brand/GetBrand":  {Type: Shop, AuthPartner: Opt, Actions: actions(ShopProductBasicInfoView)},
	"shop.Brand/GetBrands": {Type: Shop, AuthPartner: Opt, Actions: actions(ShopProductBasicInfoView)},

	"shop.History/GetFulfillmentHistory": {Type: Shop, AuthPartner: Opt},

	"shop.MoneyTransaction/GetMoneyTransaction":  {Type: Shop, Actions: actions(ShopMoneyTransactionView)},
	"shop.MoneyTransaction/GetMoneyTransactions": {Type: Shop, Actions: actions(ShopMoneyTransactionView)},

	"shop.Summary/SummarizeFulfillments": {Type: Shop, Actions: actions(ShopDashboardView)},
	"shop.Summary/SummarizePOS":          {Type: Shop, Actions: actions(ShopDashboardView)},
	"shop.Summary/SummarizeTopShip":      {Type: Shop},
	"shop.Summary/CalcBalanceShop":       {Type: Shop, AuthPartner: Opt, Actions: actions(ShopDashboardView)},
	"shop.Summary/CalcBalanceUser":       {Type: Shop, AuthPartner: Opt, Actions: actions(UserBalanceView)},

	"shop.Export/GetExports":    {Type: Shop},
	"shop.Export/RequestExport": {Type: Shop},

	"shop.Notification/CreateDevice":        {Type: Shop},
	"shop.Notification/DeleteDevice":        {Type: Shop},
	"shop.Notification/GetNotification":     {Type: Shop},
	"shop.Notification/GetNotifications":    {Type: Shop},
	"shop.Notification/UpdateNotifications": {Type: Shop},

	"shop.Authorize/GetAuthorizedPartners": {Type: Shop},
	"shop.Authorize/GetAvailablePartners":  {Type: Shop},
	"shop.Authorize/AuthorizePartner":      {Type: Shop, Actions: actions(ShopExternalAccountManage)},

	//-- Receipt --//
	"shop.Receipt/CreateReceipt":           {Type: Shop, Actions: actions(ShopReceiptCreate)},
	"shop.Receipt/UpdateReceipt":           {Type: Shop, Actions: actions(ShopReceiptUpdate)},
	"shop.Receipt/GetReceipt":              {Type: Shop, Actions: actions(ShopReceiptView)},
	"shop.Receipt/GetReceipts":             {Type: Shop, Actions: actions(ShopReceiptView)},
	"shop.Receipt/GetReceiptsByLedgerType": {Type: Shop, Actions: actions(ShopReceiptView)},
	"shop.Receipt/ConfirmReceipt":          {Type: Shop, Actions: actions(ShopReceiptConfirm)},
	"shop.Receipt/CancelReceipt":           {Type: Shop, Actions: actions(ShopReceiptCancel)},

	"shop.Trading/TradingGetProduct":  {Type: Shop, Actions: actions(ShopTradingProductView)},
	"shop.Trading/TradingGetProducts": {Type: Shop, Actions: actions(ShopTradingProductView)},
	"shop.Trading/TradingCreateOrder": {Type: Shop, Actions: actions(ShopTradingOrderCreate)},
	"shop.Trading/TradingGetOrder":    {Type: Shop, Actions: actions(ShopTradingOrderView)},
	"shop.Trading/TradingGetOrders":   {Type: Shop, Actions: actions(ShopTradingOrderView)},

	"shop.Payment/PaymentTradingOrder":    {Type: Shop, Actions: actions(ShopTradingOrderCreate)},
	"shop.Payment/PaymentCheckReturnData": {Type: Shop, Actions: actions(ShopTradingOrderCreate)},

	"shop.Inventory/CreateInventoryVoucher":          {Type: Shop, Actions: actions(ShopInventoryCreate)},
	"shop.Inventory/ConfirmInventoryVoucher":         {Type: Shop, Actions: actions(ShopInventoryConfirm)},
	"shop.Inventory/CancelInventoryVoucher":          {Type: Shop, Actions: actions(ShopInventoryCancel)},
	"shop.Inventory/UpdateInventoryVoucher":          {Type: Shop, Actions: actions(ShopInventoryUpdate)},
	"shop.Inventory/AdjustInventoryQuantity":         {Type: Shop, Actions: actions(ShopInventoryConfirm)},
	"shop.Inventory/UpdateInventoryVariantCostPrice": {Type: Shop, Actions: actions(ShopProductCostPriceUpdate)},

	"shop.Inventory/GetInventoryVariants":             {Type: Shop, Actions: actions(ShopInventoryView)},
	"shop.Inventory/GetInventoryVariantsByVariantIDs": {Type: Shop, Actions: actions(ShopInventoryView)},
	"shop.Inventory/GetInventoryVouchersByIDs":        {Type: Shop, Actions: actions(ShopInventoryView)},
	"shop.Inventory/GetInventoryVouchers":             {Type: Shop, Actions: actions(ShopInventoryView)},
	"shop.Inventory/GetInventoryVoucher":              {Type: Shop, Actions: actions(ShopInventoryView)},
	"shop.Inventory/GetInventoryVouchersByReference":  {Type: Shop, Actions: actions(ShopInventoryView)},
	"shop.Inventory/GetInventoryVariant":              {Type: Shop, Actions: actions(ShopInventoryView)},

	"shop.Brand/CreateBrand":     {Type: Shop, Actions: actions(ShopProductCreate)},
	"shop.Brand/UpdateBrandInfo": {Type: Shop, Actions: actions(ShopProductBasicInfoUpdate)},
	"shop.Brand/DeleteBrand":     {Type: Shop, Actions: actions(ShopProductDelete)},
	"shop.Brand/GetBrandByID":    {Type: Shop, Actions: actions(ShopProductBasicInfoView)},
	"shop.Brand/GetBrandsByIDs":  {Type: Shop, Actions: actions(ShopProductBasicInfoView)},
	"shop.Brand/ListBrands":      {Type: Shop, Actions: actions(ShopProductBasicInfoView)},

	"shop.Connection/GetConnections":          {Type: Shop},
	"shop.Connection/GetAvailableConnections": {Type: Shop},
	"shop.Connection/GetShopConnections":      {Type: Shop},
	"shop.Connection/LoginShopConnection":     {Type: Shop},
	"shop.Connection/RegisterShopConnection":  {Type: Shop},
	"shop.Connection/DeleteShopConnection":    {Type: Shop},
	"shop.Connection/UpdateShopConnection":    {Type: Shop},

	"shop.Subscription/GetSubscription":  {Type: Shop},
	"shop.Subscription/GetSubscriptions": {Type: Shop},
	//-- pgevent --//
	"pgevent.Misc/VersionInfo":     {Type: Secret},
	"pgevent.Event/GenerateEvents": {Type: Secret},

	//-- pghandler --//
	"handler.Misc/VersionInfo":   {Type: Secret},
	"handler.Webhook/ResetState": {Type: Secret},

	//-- exporter --//
	"exporter.Misc/VersionInfo": {Type: Secret},

	//-- crm-service --//
	"crm.User/GetUserInfo":                         {Type: Shop},
	"crm.Vtiger/GetContacts":                       {Type: EtopAdmin},
	"crm.Vtiger/CreateOrUpdateContact":             {Type: Shop},
	"crm.Vtiger/GetTickets":                        {Type: Shop},
	"crm.Vtiger/CreateTicket":                      {Type: Shop},
	"crm.Vtiger/UpdateTicket":                      {Type: Shop},
	"crm.Vtiger/GetCategories":                     {Type: Shop},
	"crm.Vtiger/GetStatus":                         {Type: EtopAdmin},
	"crm.Vtiger/CreateOrUpdateLead":                {Type: Shop},
	"crm.Vtiger/CountTicketByStatus":               {Type: EtopAdmin},
	"crm.Vtiger/GetTicketStatusCount":              {Type: EtopAdmin},
	"crm.Vht/GetCallHistories":                     {Type: EtopAdmin},
	"crm.Vht/CreateOrUpdateCallHistoryBySDKCallID": {Type: EtopAdmin},
	"crm.Vht/CreateOrUpdateCallHistoryByCallID":    {Type: EtopAdmin},

	//-- crm --//
	"crm.Misc/VersionInfo":                  {Type: Secret},
	"crm.Crm/RefreshFulfillmentFromCarrier": {Type: Secret},
	"crm.Crm/SendNotification":              {Type: Secret},

	//-- affiliate --//
	"affiliate.Misc/VersionInfo":                   {Type: Public},
	"affiliate.Account/RegisterAffiliate":          {Type: CurUsr},
	"affiliate.Account/UpdateAffiliate":            {Type: Affiliate},
	"affiliate.Account/UpdateAffiliateBankAccount": {Type: Affiliate},
	"affiliate.Account/DeleteAffiliate":            {Type: Affiliate},

	// affiliate for etopTrading
	//"affiliate.Affiliate/VersionInfo":                          {Type: Shop},
	"affiliate.User/UpdateReferral": {Type: CurUsr},

	"affiliate.Trading/TradingGetProducts":                     {Type: Shop},
	"affiliate.Trading/CreateOrUpdateTradingCommissionSetting": {Type: Shop},
	"affiliate.Trading/GetTradingProductPromotions":            {Type: Shop},
	"affiliate.Trading/GetTradingProductPromotionByProductIDs": {Type: Shop},
	"affiliate.Trading/CreateTradingProductPromotion":          {Type: Shop},
	"affiliate.Trading/UpdateTradingProductPromotion":          {Type: Shop},

	// affiliate shop
	"affiliate.Shop/GetProductPromotion":    {Type: Shop},
	"affiliate.Shop/ShopGetProducts":        {Type: Shop},
	"affiliate.Shop/CheckReferralCodeValid": {Type: Shop},

	// affiliate
	"affiliate.Affiliate/GetCommissions":                           {Type: Affiliate},
	"affiliate.Affiliate/NotifyNewShopPurchase":                    {Type: Secret},
	"affiliate.Affiliate/GetTransactions":                          {Type: Affiliate},
	"affiliate.Affiliate/CreateOrUpdateAffiliateCommissionSetting": {Type: Affiliate},
	"affiliate.Affiliate/GetProductPromotionByProductID":           {Type: Affiliate},
	"affiliate.Affiliate/CreateProductPromotion":                   {Type: Affiliate},
	"affiliate.Affiliate/UpdateProductPromotion":                   {Type: Affiliate},
	"affiliate.Affiliate/TradingGetProducts":                       {Type: Shop},
	"affiliate.Affiliate/AffiliateGetProducts":                     {Type: Affiliate},
	"affiliate.Affiliate/CreateReferralCode":                       {Type: Affiliate},
	"affiliate.Affiliate/GetReferralCodes":                         {Type: Affiliate},
	"affiliate.Affiliate/GetReferrals":                             {Type: Affiliate},

	// supplier:
	"shop.Supplier/GetSupplier":             {Type: Shop, Actions: actions(ShopSupplierView)},
	"shop.Supplier/GetSuppliers":            {Type: Shop, Actions: actions(ShopSupplierView)},
	"shop.Supplier/GetSuppliersByIDs":       {Type: Shop, Actions: actions(ShopSupplierView)},
	"shop.Supplier/CreateSupplier":          {Type: Shop, Actions: actions(ShopSupplierCreate)},
	"shop.Supplier/UpdateSupplier":          {Type: Shop, Actions: actions(ShopSupplierUpdate)},
	"shop.Supplier/DeleteSupplier":          {Type: Shop, Actions: actions(ShopSupplierDelete)},
	"shop.Supplier/GetSuppliersByVariantID": {Type: Shop, Actions: actions(ShopSupplierView)},

	// carrier:
	"shop.Carrier/GetCarrier":       {Type: Shop, Actions: actions(ShopCarrierView)},
	"shop.Carrier/GetCarriers":      {Type: Shop, Actions: actions(ShopCarrierView)},
	"shop.Carrier/GetCarriersByIDs": {Type: Shop, Actions: actions(ShopCarrierView)},
	"shop.Carrier/CreateCarrier":    {Type: Shop, Actions: actions(ShopCarrierCreate)},
	"shop.Carrier/UpdateCarrier":    {Type: Shop, Actions: actions(ShopCarrierUpdate)},
	"shop.Carrier/DeleteCarrier":    {Type: Shop, Actions: actions(ShopCarrierDelete)},

	// Ledger:
	"shop.Ledger/GetLedger":    {Type: Shop, Actions: actions(ShopLedgerView)},
	"shop.Ledger/GetLedgers":   {Type: Shop, Actions: actions(ShopLedgerView)},
	"shop.Ledger/CreateLedger": {Type: Shop, Actions: actions(ShopLedgerCreate)},
	"shop.Ledger/UpdateLedger": {Type: Shop, Actions: actions(ShopLedgerUpdate)},
	"shop.Ledger/DeleteLedger": {Type: Shop, Actions: actions(ShopLedgerDelete)},

	// PurchaseOrder:
	"shop.PurchaseOrder/GetPurchaseOrder":             {Type: Shop, Actions: actions(ShopPurchaseOrderView)},
	"shop.PurchaseOrder/GetPurchaseOrders":            {Type: Shop, Actions: actions(ShopPurchaseOrderView)},
	"shop.PurchaseOrder/GetPurchaseOrdersByIDs":       {Type: Shop, Actions: actions(ShopPurchaseOrderView)},
	"shop.PurchaseOrder/GetPurchaseOrdersByReceiptID": {Type: Shop, Actions: actions(ShopPurchaseOrderView)},
	"shop.PurchaseOrder/CreatePurchaseOrder":          {Type: Shop, Actions: actions(ShopPurchaseOrderCreate)},
	"shop.PurchaseOrder/UpdatePurchaseOrder":          {Type: Shop, Actions: actions(ShopPurchaseOrderUpdate)},
	"shop.PurchaseOrder/DeletePurchaseOrder":          {Type: Shop},
	"shop.PurchaseOrder/ConfirmPurchaseOrder":         {Type: Shop, Actions: actions(ShopPurchaseOrderConfirm)},
	"shop.PurchaseOrder/CancelPurchaseOrder":          {Type: Shop, Actions: actions(ShopPurchaseOrderCancel)},

	// Stocktake:
	"shop.Stocktake/CreateStocktake":    {Type: Shop, Auth: User, Actions: actions(ShopStocktakeCreate)},
	"shop.Stocktake/UpdateStocktake":    {Type: Shop, Actions: actions(ShopStocktakeUpdate, ShopStocktakeSelfUpdate)},
	"shop.Stocktake/ConfirmStocktake":   {Type: Shop, Actions: actions(ShopStocktakeConfirm)},
	"shop.Stocktake/CancelStocktake":    {Type: Shop, Actions: actions(ShopStocktakeCancel)},
	"shop.Stocktake/GetStocktake":       {Type: Shop, Actions: actions(ShopStocktakeView)},
	"shop.Stocktake/GetStocktakesByIDs": {Type: Shop, Actions: actions(ShopStocktakeView)},
	"shop.Stocktake/GetStocktakes":      {Type: Shop, Actions: actions(ShopStocktakeView)},

	// Relationship:
	"etop.UserRelationship/AcceptInvitation":     {Type: CurUsr},
	"etop.UserRelationship/RejectInvitation":     {Type: CurUsr},
	"etop.UserRelationship/GetInvitationByToken": {Type: Public},
	"etop.UserRelationship/GetInvitations":       {Type: CurUsr},
	"etop.UserRelationship/LeaveAccount":         {Type: CurUsr},

	// Authorization:
	"etop.AccountRelationship/CreateInvitation":   {Type: Shop, Actions: actions(RelationshipInvitationCreate)},
	"etop.AccountRelationship/GetInvitations":     {Type: Shop, Actions: actions(RelationshipInvitationView)},
	"etop.AccountRelationship/DeleteInvitation":   {Type: Shop, Actions: actions(RelationshipInvitationDelete)},
	"etop.AccountRelationship/UpdatePermission":   {Type: Shop, Actions: actions(RelationshipPermissionUpdate)},
	"etop.AccountRelationship/UpdateRelationship": {Type: Shop, Actions: actions(RelationshipRelationshipUpdate)},
	"etop.AccountRelationship/GetRelationships":   {Type: Shop, Actions: actions(RelationshipRelationshipView)},
	"etop.AccountRelationship/RemoveUser":         {Type: Shop, Actions: actions(RelationshipRelationshipRemove)},

	"shop.Refund/CreateRefund":    {Type: Shop, Actions: actions(ShopRefundCreate)},
	"shop.Refund/UpdateRefund":    {Type: Shop, Actions: actions(ShopRefundUpdate)},
	"shop.Refund/CancelRefund":    {Type: Shop, Actions: actions(ShopRefundCancel)},
	"shop.Refund/ConfirmRefund":   {Type: Shop, Actions: actions(ShopRefundConfirm)},
	"shop.Refund/GetRefund":       {Type: Shop, Actions: actions(ShopRefundView)},
	"shop.Refund/GetRefundsByIDs": {Type: Shop, Actions: actions(ShopRefundView)},
	"shop.Refund/GetRefunds":      {Type: Shop, Actions: actions(ShopRefundView)},

	"shop.PurchaseRefund/CreatePurchaseRefund":    {Type: Shop, Actions: actions(ShopPurchaseRefundCreate)},
	"shop.PurchaseRefund/UpdatePurchaseRefund":    {Type: Shop, Actions: actions(ShopPurchaseRefundUpdate)},
	"shop.PurchaseRefund/CancelPurchaseRefund":    {Type: Shop, Actions: actions(ShopPurchaseRefundCancel)},
	"shop.PurchaseRefund/ConfirmPurchaseRefund":   {Type: Shop, Actions: actions(ShopPurchaseRefundConfirm)},
	"shop.PurchaseRefund/GetPurchaseRefund":       {Type: Shop, Actions: actions(ShopPurchaseRefundView)},
	"shop.PurchaseRefund/GetPurchaseRefundsByIDs": {Type: Shop, Actions: actions(ShopPurchaseRefundView)},
	"shop.PurchaseRefund/GetPurchaseRefunds":      {Type: Shop, Actions: actions(ShopPurchaseRefundView)},

	"shop.WebServer/CreateWsWebsite":    {Type: Shop, Actions: actions(WsWebsiteCreate)},
	"shop.WebServer/UpdateWsWebsite":    {Type: Shop, Actions: actions(WsWebsiteUpdate)},
	"shop.WebServer/GetWsWebsite":       {Type: Shop, Actions: actions(WsWebsiteView)},
	"shop.WebServer/GetWsWebsites":      {Type: Shop, Actions: actions(WsWebsiteView)},
	"shop.WebServer/GetWsWebsitesByIDs": {Type: Shop, Actions: actions(WsWebsiteView)},

	"shop.WebServer/CreateOrUpdateWsProduct": {Type: Shop, Actions: actions(WsProductCreate)},
	"shop.WebServer/GetWsProduct":            {Type: Shop, Actions: actions(WsProductView)},
	"shop.WebServer/GetWsProducts":           {Type: Shop, Actions: actions(WsProductView)},
	"shop.WebServer/GetWsProductsByIDs":      {Type: Shop, Actions: actions(WsProductView)},

	"shop.WebServer/CreateOrUpdateWsCategory": {Type: Shop, Actions: actions(WsCategoryCreate)},
	"shop.WebServer/GetWsCategory":            {Type: Shop, Actions: actions(WsCategoryView)},
	"shop.WebServer/GetWsCategories":          {Type: Shop, Actions: actions(WsCategoryView)},
	"shop.WebServer/GetWsCategoriesByIDs":     {Type: Shop, Actions: actions(WsCategoryView)},

	"shop.WebServer/CreateWsPage":    {Type: Shop, Actions: actions(WsPageCreate)},
	"shop.WebServer/UpdateWsPage":    {Type: Shop, Actions: actions(WsPageUpdate)},
	"shop.WebServer/DeleteWsPage":    {Type: Shop, Actions: actions(WsPageDelete)},
	"shop.WebServer/GetWsPage":       {Type: Shop, Actions: actions(WsPageView)},
	"shop.WebServer/GetWsPages":      {Type: Shop, Actions: actions(WsPageView)},
	"shop.WebServer/GetWsPagesByIDs": {Type: Shop, Actions: actions(WsPageView)},

	// -- Fabo --
	"fabo.Page/ConnectPages":                                 {Type: Shop},
	"fabo.Page/RemovePages":                                  {Type: Shop, Auth: User, IncludeFaboInfo: true},
	"fabo.Page/ListPages":                                    {Type: Shop, Auth: User, IncludeFaboInfo: true},
	"fabo.CustomerConversation/ListCustomerConversations":    {Type: Shop, Auth: User, IncludeFaboInfo: true},
	"fabo.CustomerConversation/ListMessages":                 {Type: Shop, Auth: User, IncludeFaboInfo: true},
	"fabo.CustomerConversation/ListCommentsByExternalPostID": {Type: Shop, Auth: User, IncludeFaboInfo: true},
	"fabo.CustomerConversation/UpdateReadStatus":             {Type: Shop, Auth: User, IncludeFaboInfo: true},
	"fabo.CustomerConversation/SendMessage":                  {Type: Shop, Auth: User, IncludeFaboInfo: true},
	"fabo.CustomerConversation/SendComment":                  {Type: Shop, Auth: User, IncludeFaboInfo: true},

	// -- Fabo Customer --
	"fabo.Customer/CreateFbUserCustomer":     {Type: Shop},
	"fabo.Customer/ListFbUsers":              {Type: Shop},
	"fabo.Customer/GetFbUser":                {Type: Shop},
	"fabo.Customer/ListCustomersWithFbUsers": {Type: Shop},
}

func actions(actions ...permission.ActionType) (actionsResult []permission.ActionType) {
	actionsResult = append(actionsResult, actions...)
	return
}
