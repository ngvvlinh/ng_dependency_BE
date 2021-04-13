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

	ShopConnectionCreate permission.ActionType = "shop/connection:create"
	ShopConnectionUpdate permission.ActionType = "shop/connection:update"
	ShopConnectionDelete permission.ActionType = "shop/connection:delete"
	ShopConnectionView   permission.ActionType = "shop/connection:view"

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
	ShopFulfillmentUpdate permission.ActionType = "shop/fulfillment:update"
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

	// shop ticket
	ShopTicketCreate permission.ActionType = "shop/shop_ticket:create"
	ShopTicketUpdate permission.ActionType = "shop/shop_ticket:update"
	ShopTicketAssign permission.ActionType = "shop/shop_ticket:assign"
	ShopTicketView   permission.ActionType = "shop/shop_ticket:view"

	// shop ticket comment
	ShopTicketCommentCreate permission.ActionType = "shop/shop_ticket_comment:create"
	ShopTicketCommentUpdate permission.ActionType = "shop/shop_ticket_comment:update"
	ShopTicketCommentView   permission.ActionType = "shop/shop_ticket_comment:view"
	ShopTicketCommentDelete permission.ActionType = "shop/shop_ticket_comment:delete"

	// shop ticket label
	ShopTicketLabelCreate permission.ActionType = "shop/shop_ticket_label:create"
	ShopTicketLabelUpdate permission.ActionType = "shop/shop_ticket_label:update"
	ShopTicketLabelView   permission.ActionType = "shop/shop_ticket_label:view"

	ShopPaymentCreate permission.ActionType = "shop/payment:create"
	ShopPaymentView   permission.ActionType = "shop/payment:view"

	// shop extension etelecom
	ShopExtensionCreate permission.ActionType = "shop/extension:create"
	ShopExtensionDelete permission.ActionType = "shop/extension:delete"
	ShopExtensionView   permission.ActionType = "shop/extension:view"
	ShopHotlineView     permission.ActionType = "shop/hotline:view"
	ShopCallLogView     permission.ActionType = "shop/calllog:view"
	ShopCallLogCreate   permission.ActionType = "shop/calllog:create"

	// shop credit
	ShopCreditCreate permission.ActionType = "shop/credit:create"

	// etelecom user setting
	ShopEtelecomUserSettingView permission.ActionType = "shop/etelecom_user_setting:view"

	// shop subscription
	ShopSubscriptionProductView permission.ActionType = "shop/subscription_product:view"
	ShopSubscriptionPlanView    permission.ActionType = "shop/subscription_plan:view"
	ShopSubscriptionView        permission.ActionType = "shop/subscription:view"
	ShopSubscriptionCreate      permission.ActionType = "shop/subscription:create"
	ShopSubscriptionUpdate      permission.ActionType = "shop/subscription:update"

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

	// shop transaction
	ShopTransactionView permission.ActionType = "shop/transaction:view"

	// Fabo
	FbCommentView   permission.ActionType = "facebook/comment:view"
	FbCommentCreate permission.ActionType = "facebook/comment:create"
	FbMessageCreate permission.ActionType = "facebook/message:create"
	FbMessageView   permission.ActionType = "facebook/message:view"
	FbPostCreate    permission.ActionType = "facebook/post:create"
	FbUserCreate    permission.ActionType = "facebook/fbuser:create" // Liên kết fb_user với customer
	FbUserView      permission.ActionType = "facebook/fbuser:view"
	FbUserUpdate    permission.ActionType = "facebook/fbuser:update"
	FbFanpageCreate permission.ActionType = "facebook/fanpage:create"
	FbFanpageDelete permission.ActionType = "facebook/fanpage:delete"
	FbFanpageView   permission.ActionType = "facebook/fanpage:view"

	FbShopTagCreate permission.ActionType = "facebook/shoptag:create"
	FbShopTagUpdate permission.ActionType = "facebook/shoptag:update"
	FbShopTagView   permission.ActionType = "facebook/shoptag:view"
	FbShopTagDelete permission.ActionType = "facebook/shoptag:delete"

	FbMessageTemplateCreate permission.ActionType = "facebook/message_template:create"
	FbMessageTemplateUpdate permission.ActionType = "facebook/message_template:update"
	FbMessageTemplateView   permission.ActionType = "facebook/message_template:view"
	FbMessageTemplateDelete permission.ActionType = "facebook/message_template:delete"

	// Admin Credit
	AdminCreditCreate  permission.ActionType = "admin/credit:create"
	AdminCreditView    permission.ActionType = "admin/credit:view"
	AdminCreditUpdate  permission.ActionType = "admin/credit:update"
	AdminCreditConfirm permission.ActionType = "admin/credit:confirm"
	AdminCreditDelete  permission.ActionType = "admin/credit:delete"

	// Admin Fulfillment
	AdminFulfillmentView               permission.ActionType = "admin/fulfillment:view"
	AdminFulfillmentCreate             permission.ActionType = "admin/fulfillment:create"
	AdminFulfillmentStateUpdate        permission.ActionType = "admin/fulfillment_state:update"
	AdminFulfillmentUpdate             permission.ActionType = "admin/fulfillment:update"
	AdminFulfillmentShippingFeesUpdate permission.ActionType = "admin/fulfillment_shipping_fees:update"

	// Admin Account
	AdminAdminUserCreate permission.ActionType = "admin/admin_user:create"
	AdminAdminUserUpdate permission.ActionType = "admin/admin_user:update"
	AdminAdminUserView   permission.ActionType = "admin/admin_user:view"
	AdminAdminUserDelete permission.ActionType = "admin/admin_user:delete"
	AdminPartnerCreate   permission.ActionType = "admin/partner:create"

	// Admin ticket
	AdminLeadTicketAssign permission.ActionType = "admin/admin_ticket:assign"
	AdminTicketCreate     permission.ActionType = "admin/admin_ticket:create"
	AdminTicketUpdate     permission.ActionType = "admin/admin_ticket:update"
	LeadAdminTicketReopen permission.ActionType = "admin/admin_lead_ticket:update"
	AdminTicketView       permission.ActionType = "admin/admin_ticket:view"

	AdminTicketCommentCreate permission.ActionType = "admin/admin_ticket:create"
	AdminTicketCommentUpdate permission.ActionType = "admin/admin_ticket:update"
	AdminTicketCommentDelete permission.ActionType = "admin/admin_ticket:delete"
	AdminTicketCommentView   permission.ActionType = "admin/admin_ticket:view"

	AdminTicketLabelCreate permission.ActionType = "admin/admin_ticket_label:create"
	AdminTicketLabelUpdate permission.ActionType = "admin/admin_ticket_label:update"
	AdminTicketLabelView   permission.ActionType = "admin/admin_ticket_label:view"

	// Admin MoneyTransaction
	AdminMoneyTransactionCreate  permission.ActionType = "admin/money_transaction:create"
	AdminMoneyTransactionView    permission.ActionType = "admin/money_transaction:view"
	AdminMoneyTransactionUpdate  permission.ActionType = "admin/money_transaction:update"
	AdminMoneyTransactionConfirm permission.ActionType = "admin/money_transaction:confirm"

	AdminMoneyTransactionShippingEtopView    permission.ActionType = "admin/money_transaction_shipping_etop:view"
	AdminMoneyTransactionShippingEtopCreate  permission.ActionType = "admin/money_transaction_shipping_etop:create"
	AdminMoneyTransactionShippingEtopUpdate  permission.ActionType = "admin/money_transaction_shipping_etop:update"
	AdminMoneyTransactionShippingEtopConfirm permission.ActionType = "admin/money_transaction_shipping_etop:confirm"
	AdminMoneyTransactionShippingEtopDelete  permission.ActionType = "admin/money_transaction_shipping_etop:delete"

	AdminMoneyTransactionShippingExternalView        permission.ActionType = "admin/money_transaction_shipping_external:view"
	AdminMoneyTransactionShippingExternalUpdate      permission.ActionType = "admin/money_transaction_shipping_external:update"
	AdminMoneyTransactionShippingExternalConfirm     permission.ActionType = "admin/money_transaction_shipping_external:confirm"
	AdminMoneyTransactionShippingExternalDelete      permission.ActionType = "admin/money_transaction_shipping_external:delete"
	AdminMoneyTransactionShippingExternalLinesRemove permission.ActionType = "admin/money_transaction_shipping_external_lines:remove"

	// Admin Order
	AdminOrderView permission.ActionType = "admin/order:view"

	// Admin Shop ShipmentPrice
	AdminShopShipmentPriceListCreate permission.ActionType = "admin/shop_shipment_price_list:create"
	AdminShopShipmentPriceListUpdate permission.ActionType = "admin/shop_shipment_price_list:update"
	AdminShopShipmentPriceListDelete permission.ActionType = "admin/shop_shipment_price_list:delete"
	AdminShopShipmentPriceListView   permission.ActionType = "admin/shop_shipment_price_list:view"

	// Admin ShipmentPrice Services
	AdminShipmentServiceView   permission.ActionType = "admin/shipment_service:view"
	AdminShipmentServiceCreate permission.ActionType = "admin/shipment_service:create"
	AdminShipmentServiceUpdate permission.ActionType = "admin/shipment_service:update"
	AdminShipmentServiceDelete permission.ActionType = "admin/shipment_service:delete"

	// Admin ShipmentPrice List
	AdminShipmentPriceListCreate permission.ActionType = "admin/shipment_price_list:create"
	AdminShipmentPriceListUpdate permission.ActionType = "admin/shipment_price_list:update"
	AdminShipmentPriceListDelete permission.ActionType = "admin/shipment_price_list:delete"
	AdminShipmentPriceListView   permission.ActionType = "admin/shipment_price_list:view"

	// Admin ShipmentPrice List Promotion
	AdminShipmentPriceListPromotionCreate permission.ActionType = "admin/shipment_price_list_promotion:create"
	AdminShipmentPriceListPromotionUpdate permission.ActionType = "admin/shipment_price_list_promotion:update"
	AdminShipmentPriceListPromotionDelete permission.ActionType = "admin/shipment_price_list_promotion:delete"
	AdminShipmentPriceListPromotionView   permission.ActionType = "admin/shipment_price_list_promotion:view"

	// Admin ShipmentPrice
	AdminShipmentPriceCreate permission.ActionType = "admin/shipment_price:create"
	AdminShipmentPriceUpdate permission.ActionType = "admin/shipment_price:update"
	AdminShipmentPriceDelete permission.ActionType = "admin/shipment_price:delete"
	AdminShipmentPriceView   permission.ActionType = "admin/shipment_price:view"

	// Admin Shop
	AdminShopView   permission.ActionType = "admin/shop:view"
	AdminShopUpdate permission.ActionType = "admin/shop:update"

	// Admin User
	AdminUserView      permission.ActionType = "admin/user:view"
	AdminUserBlock     permission.ActionType = "admin/user:block"
	AdminUpdateUserRef permission.ActionType = "admin/user_ref:update"

	// AdminSubscription
	AdminSubscriptionProductCreate permission.ActionType = "admin/subscription_product:create"
	AdminSubscriptionProductView   permission.ActionType = "admin/subscription_product:view"
	AdminSubscriptionProductDelete permission.ActionType = "admin/subscription_product:delete"
	AdminSubscriptionPlanCreate    permission.ActionType = "admin/subscription_plan:create"
	AdminSubscriptionPlanUpdate    permission.ActionType = "admin/subscription_plan:update"
	AdminSubscriptionPlanView      permission.ActionType = "admin/subscription_plan:view"
	AdminSubscriptionPlanDelete    permission.ActionType = "admin/subscription_plan:delete"

	AdminSubscriptionView           permission.ActionType = "admin/subscription:view"
	AdminSubscriptionCreate         permission.ActionType = "admin/subscription:create"
	AdminSubscriptionUpdate         permission.ActionType = "admin/subscription:update"
	AdminSubscriptionCancel         permission.ActionType = "admin/subscription:cancel"
	AdminSubscriptionActive         permission.ActionType = "admin/subscription:active"
	AdminSubscriptionDelete         permission.ActionType = "admin/subscription:delete"
	AdminInvoiceView                permission.ActionType = "admin/invoice:view"
	AdminInvoiceCreate              permission.ActionType = "admin/invoice:create"
	AdminManualPaymentInvoiceCreate permission.ActionType = "admin/invoice_manual_payment:create"
	AdminInvoiceDelete              permission.ActionType = "admin/invoice:delete"

	// AdminMisc
	AdminMiscLoginAccount permission.ActionType = "admin/misc_account:login"

	// AdminConnection
	AdminConnectionView              permission.ActionType = "admin/connection:view"
	AdminConnectionConfirm           permission.ActionType = "admin/connection:confirm"
	AdminConnectionDisable           permission.ActionType = "admin/connection:disable"
	AdminConnectionBuiltinCreate     permission.ActionType = "admin/connection_builtin:create"
	AdminConnectionShopBuiltinCreate permission.ActionType = "admin/shop_connection_builtin:create"
	AdminConnectionShopUpdate        permission.ActionType = "admin/shop_connection:update"
	AdminConnectionServiceView       permission.ActionType = "admin/connection_service:view"

	// Admin custom region
	AdminCustomRegionCreate permission.ActionType = "admin/custom_region:create"
	AdminCustomRegionUpdate permission.ActionType = "admin/custom_region:update"
	AdminCustomRegionDelete permission.ActionType = "admin/custom_region:delete"
	AdminCustomRegionView   permission.ActionType = "admin/custom_region:view"

	// admin extension etelecom
	AdminHotlineCreate permission.ActionType = "admin/hotline:create"
	AdminHotlineUpdate permission.ActionType = "admin/hotline:update"

	// admin etelecom user setting
	AdminEtelecomUserSettingUpdate permission.ActionType = "admin/etelecom_user_setting:update"
	AdminEtelecomUserSettingView   permission.ActionType = "admin/etelecom_user_setting:view"
)

// ACL declares access control list
var _acl = map[string]*permission.Decl{
	//-- sadmin --//

	"sadmin.User/CreateUser":           {Type: SuperAdmin},
	"sadmin.User/ResetPassword":        {Type: SuperAdmin},
	"sadmin.User/LoginAsAccount":       {Type: SuperAdmin},
	"sadmin.Misc/VersionInfo":          {Type: SuperAdmin},
	"sadmin.Webhook/RegisterWebhook":   {Type: SuperAdmin},
	"sadmin.Webhook/UnregisterWebhook": {Type: SuperAdmin},

	//-- common --//

	"       etop.Misc/VersionInfo": {Type: Public},
	"      admin.Misc/VersionInfo": {Type: Public},
	"       shop.Misc/VersionInfo": {Type: Public},
	"ext/partner.Misc/VersionInfo": {Type: Public},
	"   ext/shop.Misc/VersionInfo": {Type: Partner, Auth: APIKey},
	"integration.Misc/VersionInfo": {Type: Public},

	"admin.Misc/AdminLoginAsAccount": {Type: EtopAdmin, Actions: actions(AdminMiscLoginAccount)},

	"etop.User/Register":                 {Type: Public},
	"etop.User/RegisterUsingToken":       {Type: Public},
	"etop.User/RequestRegisterSimplify":  {Type: Public},
	"etop.User/RegisterSimplify":         {Type: Public},
	"etop.User/Login":                    {Type: Public},
	"etop.User/ResetPassword":            {Type: Public, Captcha: "custom"},
	"etop.User/ChangePasswordUsingToken": {Type: Public},
	"etop.User/ChangePassword":           {Type: CurUsr},
	"etop.User/GetNotifySetting":         {Type: CurUsr},
	"etop.User/EnableNotifyTopic":        {Type: CurUsr},
	"etop.User/DisableNotifyTopic":       {Type: CurUsr},
	"etop.User/ChangeRefAff":             {Type: CurUsr},
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

	"etop.User/WebphoneRequestLogin": {Type: Public},
	"etop.User/WebphoneLogin":        {Type: Public},

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
	"etop.Bank/GetBankProvinces":          {Type: CurUsr},
	"etop.Bank/GetBankBranches":           {Type: CurUsr},
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

	"ext/shop.Shipnow/GetShipnowServices":       {Type: Shop, Auth: APIKey},
	"ext/shop.Shipnow/CreateShipnowFulfillment": {Type: Shop, Auth: APIKey},
	"ext/shop.Shipnow/CancelShipnowFulfillment": {Type: Shop, Auth: APIKey},
	"ext/shop.Shipnow/GetShipnowFulfillment":    {Type: Shop, Auth: APIKey},

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

	//-- vnpost --//

	"ext/vnposts/ping":              {Type: Shop, Auth: APIKey},
	"ext/vnposts/getservicesvnpost": {Type: Shop, Auth: APIKey},
	"ext/vnposts/createordervnpost": {Type: Shop, Auth: APIKey},
	"ext/vnposts/cancelordervnpost": {Type: Shop, Auth: APIKey},
	"ext/vnposts/getordervnpost":    {Type: Shop, Auth: APIKey},

	"ext/vnposts/webhook/createwebhook": {Type: Shop, Auth: APIKey},
	"ext/vnposts/webhook/getwebhooks":   {Type: Shop, Auth: APIKey},
	"ext/vnposts/webhook/deletewebhook": {Type: Shop, Auth: APIKey},

	//-- vht --/

	// Only support VNPost wl partner
	"ext/vht.User/RegisterUser": {Type: Partner, Auth: APIKey},

	//-- integration --//

	"integration.Integration/Init":              {Type: Public},
	"integration.Integration/RequestLogin":      {Type: Protected, AuthPartner: Req, Captcha: "1"},
	"integration.Integration/LoginUsingToken":   {Type: Protected, AuthPartner: Req},
	"integration.Integration/LoginUsingTokenWL": {Type: Protected, AuthPartner: Req},
	"integration.Integration/Register":          {Type: Protected, AuthPartner: Req},
	"integration.Integration/GrantAccess":       {Type: CurUsr, AuthPartner: Req},
	"integration.Integration/SessionInfo":       {Type: Protected, AuthPartner: Req},

	//-- admin --//

	"admin.Account/CreatePartner":   {Type: EtopAdmin, Actions: actions(AdminPartnerCreate)},
	"admin.Account/GenerateAPIKey":  {Type: EtopAdmin},
	"admin.Account/CreateAdminUser": {Type: EtopAdmin, Actions: actions(AdminAdminUserCreate)},
	"admin.Account/UpdateAdminUser": {Type: EtopAdmin, Actions: actions(AdminAdminUserUpdate)},
	"admin.Account/GetAdminUsers":   {Type: EtopAdmin, Actions: actions(AdminAdminUserView)},
	"admin.Account/DeleteAdminUser": {Type: EtopAdmin, Actions: actions(AdminAdminUserDelete)},

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

	"admin.Order/GetOrder":       {Type: EtopAdmin, Actions: actions(AdminOrderView)},
	"admin.Order/GetOrders":      {Type: EtopAdmin, Actions: actions(AdminOrderView)},
	"admin.Order/GetOrdersByIDs": {Type: EtopAdmin, Actions: actions(AdminOrderView)},

	"admin.User/GetUsers":      {Type: EtopAdmin, Actions: actions(AdminUserView)},
	"admin.User/GetUser":       {Type: EtopAdmin, Actions: actions(AdminUserView)},
	"admin.User/GetUsersByIDs": {Type: EtopAdmin, Actions: actions(AdminUserView)},
	"admin.User/BlockUser":     {Type: EtopAdmin, Actions: actions(AdminUserBlock)},
	"admin.User/UnblockUser":   {Type: EtopAdmin, Actions: actions(AdminUserBlock)},
	"admin.User/UpdateUserRef": {Type: EtopAdmin, Actions: actions(AdminUpdateUserRef)},

	"admin.Fulfillment/GetFulfillment":                 {Type: EtopAdmin, Actions: actions(AdminFulfillmentView)},
	"admin.Fulfillment/GetFulfillments":                {Type: EtopAdmin, Actions: actions(AdminFulfillmentView)},
	"admin.Fulfillment/UpdateFulfillment":              {Type: EtopAdmin, Actions: actions(AdminFulfillmentUpdate)}, // deprecated
	"admin.Fulfillment/UpdateFulfillmentInfo":          {Type: EtopAdmin, Actions: actions(AdminFulfillmentUpdate)},
	"admin.Fulfillment/UpdateFulfillmentCODAmount":     {Type: EtopAdmin, Actions: actions(AdminFulfillmentUpdate)},
	"admin.Fulfillment/AddShippingFee":                 {Type: EtopAdmin, Actions: actions(AdminFulfillmentUpdate)},
	"admin.Fulfillment/UpdateFulfillmentShippingState": {Type: EtopAdmin, Actions: actions(AdminFulfillmentStateUpdate)},
	"admin.Fulfillment/UpdateFulfillmentShippingFees":  {Type: EtopAdmin, Actions: actions(AdminFulfillmentShippingFeesUpdate)},

	"admin.MoneyTransaction/CreateMoneyTransaction":  {Type: EtopAdmin, Actions: actions(AdminMoneyTransactionCreate)},
	"admin.MoneyTransaction/GetMoneyTransaction":     {Type: EtopAdmin, Actions: actions(AdminMoneyTransactionView)},
	"admin.MoneyTransaction/GetMoneyTransactions":    {Type: EtopAdmin, Actions: actions(AdminMoneyTransactionView)},
	"admin.MoneyTransaction/ConfirmMoneyTransaction": {Type: EtopAdmin, Actions: actions(AdminMoneyTransactionConfirm)},
	"admin.MoneyTransaction/UpdateMoneyTransaction":  {Type: EtopAdmin, Actions: actions(AdminMoneyTransactionUpdate)},

	"admin.MoneyTransaction/GetMoneyTransactionShippingExternal":         {Type: EtopAdmin, Actions: actions(AdminMoneyTransactionShippingExternalView)},
	"admin.MoneyTransaction/GetMoneyTransactionShippingExternals":        {Type: EtopAdmin, Actions: actions(AdminMoneyTransactionShippingExternalView)},
	"admin.MoneyTransaction/RemoveMoneyTransactionShippingExternalLines": {Type: EtopAdmin, Actions: actions(AdminMoneyTransactionShippingExternalLinesRemove)},
	"admin.MoneyTransaction/DeleteMoneyTransactionShippingExternal":      {Type: EtopAdmin, Actions: actions(AdminMoneyTransactionShippingExternalDelete)},
	"admin.MoneyTransaction/ConfirmMoneyTransactionShippingExternals":    {Type: EtopAdmin, Actions: actions(AdminMoneyTransactionShippingExternalConfirm)},
	"admin.MoneyTransaction/UpdateMoneyTransactionShippingExternal":      {Type: EtopAdmin, Actions: actions(AdminMoneyTransactionShippingExternalUpdate)},
	"admin.MoneyTransaction/SplitMoneyTransactionShippingExternal":       {Type: EtopAdmin, Actions: actions(AdminMoneyTransactionShippingExternalUpdate)},

	"admin.MoneyTransaction/GetMoneyTransactionShippingEtop":     {Type: EtopAdmin, Actions: actions(AdminMoneyTransactionShippingEtopView)},
	"admin.MoneyTransaction/GetMoneyTransactionShippingEtops":    {Type: EtopAdmin, Actions: actions(AdminMoneyTransactionShippingEtopView)},
	"admin.MoneyTransaction/CreateMoneyTransactionShippingEtop":  {Type: EtopAdmin, Actions: actions(AdminMoneyTransactionShippingEtopCreate)},
	"admin.MoneyTransaction/UpdateMoneyTransactionShippingEtop":  {Type: EtopAdmin, Actions: actions(AdminMoneyTransactionShippingEtopUpdate)},
	"admin.MoneyTransaction/DeleteMoneyTransactionShippingEtop":  {Type: EtopAdmin, Actions: actions(AdminMoneyTransactionShippingEtopDelete)},
	"admin.MoneyTransaction/ConfirmMoneyTransactionShippingEtop": {Type: EtopAdmin, Actions: actions(AdminMoneyTransactionShippingEtopConfirm)},

	"admin.Shop/GetShop":        {Type: EtopAdmin, Actions: actions(AdminShopView)},
	"admin.Shop/GetShops":       {Type: EtopAdmin, Actions: actions(AdminUserView)},
	"admin.Shop/GetShopsByIDs":  {Type: EtopAdmin, Actions: actions(AdminUserView)},
	"admin.Shop/UpdateShopInfo": {Type: EtopAdmin, Actions: actions(AdminShopUpdate)},

	"admin.Credit/GetCredit":     {Type: EtopAdmin, Actions: actions(AdminCreditView)},
	"admin.Credit/GetCredits":    {Type: EtopAdmin, Actions: actions(AdminCreditView)},
	"admin.Credit/CreateCredit":  {Type: EtopAdmin, Actions: actions(AdminCreditCreate)},
	"admin.Credit/ConfirmCredit": {Type: EtopAdmin, Actions: actions(AdminCreditConfirm)},
	"admin.Credit/DeleteCredit":  {Type: EtopAdmin, Actions: actions(AdminCreditDelete)},

	"admin.Notification/CreateNotifications": {Type: EtopAdmin},

	"admin.Connection/GetConnections":            {Type: EtopAdmin, Actions: actions(AdminConnectionView)},
	"admin.Connection/ConfirmConnection":         {Type: EtopAdmin, Actions: actions(AdminConnectionConfirm)},
	"admin.Connection/DisableConnection":         {Type: EtopAdmin, Actions: actions(AdminConnectionDisable)},
	"admin.Connection/CreateBuiltinConnection":   {Type: EtopAdmin, Actions: actions(AdminConnectionBuiltinCreate)},
	"admin.Connection/GetBuiltinShopConnections": {Type: EtopAdmin, Actions: actions(AdminConnectionShopBuiltinCreate)},
	"admin.Connection/UpdateShopConnection":      {Type: EtopAdmin, Actions: actions(AdminConnectionShopUpdate)},
	"admin.Connection/GetConnectionServices":     {Type: EtopAdmin, Actions: actions(AdminConnectionServiceView)},

	"admin.ShipmentPrice/GetShipmentServices":                      {Type: EtopAdmin, Actions: actions(AdminShipmentServiceView)},
	"admin.ShipmentPrice/GetShipmentService":                       {Type: EtopAdmin, Actions: actions(AdminShipmentServiceView)},
	"admin.ShipmentPrice/CreateShipmentService":                    {Type: EtopAdmin, Actions: actions(AdminShipmentServiceCreate)},
	"admin.ShipmentPrice/UpdateShipmentService":                    {Type: EtopAdmin, Actions: actions(AdminShipmentServiceUpdate)},
	"admin.ShipmentPrice/DeleteShipmentService":                    {Type: EtopAdmin, Actions: actions(AdminShipmentServiceDelete)},
	"admin.ShipmentPrice/UpdateShipmentServicesAvailableLocations": {Type: EtopAdmin, Actions: actions(AdminShipmentServiceUpdate)},
	"admin.ShipmentPrice/UpdateShipmentServicesBlacklistLocations": {Type: EtopAdmin, Actions: actions(AdminShipmentServiceUpdate)},

	"admin.ShipmentPrice/GetShipmentPriceLists":       {Type: EtopAdmin, Actions: actions(AdminShipmentPriceListView)},
	"admin.ShipmentPrice/GetShipmentPriceList":        {Type: EtopAdmin, Actions: actions(AdminShipmentPriceListView)},
	"admin.ShipmentPrice/CreateShipmentPriceList":     {Type: EtopAdmin, Actions: actions(AdminShipmentPriceListCreate)},
	"admin.ShipmentPrice/UpdateShipmentPriceList":     {Type: EtopAdmin, Actions: actions(AdminShipmentPriceListUpdate)},
	"admin.ShipmentPrice/SetDefaultShipmentPriceList": {Type: EtopAdmin, Actions: actions(AdminShipmentPriceListUpdate)},
	"admin.ShipmentPrice/DeleteShipmentPriceList":     {Type: EtopAdmin, Actions: actions(AdminShipmentPriceListDelete)},

	"admin.ShipmentPrice/GetShipmentPriceListPromotions":   {Type: EtopAdmin, Actions: actions(AdminShipmentPriceListPromotionView)},
	"admin.ShipmentPrice/GetShipmentPriceListPromotion":    {Type: EtopAdmin, Actions: actions(AdminShipmentPriceListPromotionView)},
	"admin.ShipmentPrice/CreateShipmentPriceListPromotion": {Type: EtopAdmin, Actions: actions(AdminShipmentPriceListPromotionCreate)},
	"admin.ShipmentPrice/UpdateShipmentPriceListPromotion": {Type: EtopAdmin, Actions: actions(AdminShipmentPriceListPromotionUpdate)},
	"admin.ShipmentPrice/DeleteShipmentPriceListPromotion": {Type: EtopAdmin, Actions: actions(AdminShipmentPriceListPromotionDelete)},

	"admin.ShipmentPrice/GetShipmentPrice":                  {Type: EtopAdmin, Actions: actions(AdminShipmentPriceView)},
	"admin.ShipmentPrice/GetShipmentPrices":                 {Type: EtopAdmin, Actions: actions(AdminShipmentPriceView)},
	"admin.ShipmentPrice/CreateShipmentPrice":               {Type: EtopAdmin, Actions: actions(AdminShipmentPriceCreate)},
	"admin.ShipmentPrice/UpdateShipmentPrice":               {Type: EtopAdmin, Actions: actions(AdminShipmentPriceUpdate)},
	"admin.ShipmentPrice/DeleteShipmentPrice":               {Type: EtopAdmin, Actions: actions(AdminShipmentPriceDelete)},
	"admin.ShipmentPrice/UpdateShipmentPricesPriorityPoint": {Type: EtopAdmin, Actions: actions(AdminShipmentPriceUpdate)},

	"admin.ShipmentPrice/GetShopShipmentPriceLists":   {Type: EtopAdmin, Actions: actions(AdminShopShipmentPriceListView)},
	"admin.ShipmentPrice/GetShopShipmentPriceList":    {Type: EtopAdmin, Actions: actions(AdminShopShipmentPriceListView)},
	"admin.ShipmentPrice/CreateShopShipmentPriceList": {Type: EtopAdmin, Actions: actions(AdminShopShipmentPriceListCreate)},
	"admin.ShipmentPrice/UpdateShopShipmentPriceList": {Type: EtopAdmin, Actions: actions(AdminShopShipmentPriceListUpdate)},
	"admin.ShipmentPrice/DeleteShopShipmentPriceList": {Type: EtopAdmin, Actions: actions(AdminShopShipmentPriceListDelete)},

	"admin.ShipmentPrice/GetShippingServices": {Type: EtopAdmin},

	"admin.Location/GetCustomRegion":    {Type: EtopAdmin, Actions: actions(AdminCustomRegionView)},
	"admin.Location/GetCustomRegions":   {Type: EtopAdmin, Actions: actions(AdminCustomRegionView)},
	"admin.Location/CreateCustomRegion": {Type: EtopAdmin, Actions: actions(AdminCustomRegionCreate)},
	"admin.Location/UpdateCustomRegion": {Type: EtopAdmin, Actions: actions(AdminCustomRegionUpdate)},
	"admin.Location/DeleteCustomRegion": {Type: EtopAdmin, Actions: actions(AdminCustomRegionDelete)},

	"admin.Subscription/CreateSubscriptionProduct": {Type: EtopAdmin, Actions: actions(AdminSubscriptionProductCreate)},
	"admin.Subscription/GetSubscriptionProducts":   {Type: EtopAdmin, Actions: actions(AdminSubscriptionProductView)},
	"admin.Subscription/DeleteSubscriptionProduct": {Type: EtopAdmin, Actions: actions(AdminSubscriptionProductDelete)},
	"admin.Subscription/CreateSubscriptionPlan":    {Type: EtopAdmin, Actions: actions(AdminSubscriptionPlanCreate)},
	"admin.Subscription/UpdateSubscriptionPlan":    {Type: EtopAdmin, Actions: actions(AdminSubscriptionPlanUpdate)},
	"admin.Subscription/GetSubscriptionPlans":      {Type: EtopAdmin, Actions: actions(AdminSubscriptionPlanView)},
	"admin.Subscription/DeleteSubscriptionPlan":    {Type: EtopAdmin, Actions: actions(AdminSubscriptionPlanDelete)},

	"admin.Subscription/GetSubscription":        {Type: EtopAdmin, Actions: actions(AdminSubscriptionView)},
	"admin.Subscription/GetSubscriptions":       {Type: EtopAdmin, Actions: actions(AdminSubscriptionView)},
	"admin.Subscription/CreateSubscription":     {Type: EtopAdmin, Actions: actions(AdminSubscriptionCreate)},
	"admin.Subscription/UpdateSubscriptionInfo": {Type: EtopAdmin, Actions: actions(AdminSubscriptionUpdate)},
	"admin.Subscription/CancelSubscription":     {Type: EtopAdmin, Actions: actions(AdminSubscriptionCancel)},
	"admin.Subscription/ActivateSubscription":   {Type: EtopAdmin, Actions: actions(AdminSubscriptionActive)},
	"admin.Subscription/DeleteSubscription":     {Type: EtopAdmin, Actions: actions(AdminSubscriptionDelete)},

	"admin.Invoice/GetInvoices":          {Type: EtopAdmin, Actions: actions(AdminInvoiceView)},
	"admin.Invoice/CreateInvoice":        {Type: EtopAdmin, Actions: actions(AdminInvoiceCreate)},
	"admin.Invoice/ManualPaymentInvoice": {Type: EtopAdmin, Actions: actions(AdminManualPaymentInvoiceCreate)},
	"admin.Invoice/DeleteInvoice":        {Type: EtopAdmin, Actions: actions(AdminInvoiceDelete)},

	"admin.Etelecom/CreateHotline": {Type: EtopAdmin, Actions: actions(AdminHotlineCreate)},
	"admin.Etelecom/UpdateHotline": {Type: EtopAdmin, Actions: actions(AdminHotlineUpdate)},

	//-- shop --//

	"shop.Account/RegisterShop": {Type: CurUsr, AuthPartner: Opt},

	// permission: owner

	"shop.Account/UpdateShop": {Type: Shop, Actions: actions(ShopSettingsShopInfoUpdate)},
	"shop.Account/DeleteShop": {Type: Shop, Actions: actions(ShopAccountDelete)},

	// permission: admin

	"shop.Account/SetDefaultAddress": {Type: Shop, Actions: actions(ShopSettingsShopInfoUpdate)},

	"shop.Account/CreateExternalAccountAhamove":                   {Type: Shop, Actions: actions(ShopExternalAccountManage)},
	"shop.Account/GetExternalAccountAhamove":                      {Type: Shop},
	"shop.Account/RequestVerifyExternalAccountAhamove":            {Type: Shop, Actions: actions(ShopExternalAccountManage)},
	"shop.Account/UpdateExternalAccountAhamoveVerification":       {Type: Shop, Actions: actions(ShopExternalAccountManage)},
	"shop.Account/UpdateExternalAccountAhamoveVerificationImages": {Type: Shop, Actions: actions(ShopExternalAccountManage)},

	"shop.AccountShipnow/GetAccountShipnow": {Type: Shop},

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

	"shop.Credit/CreateCredit": {Type: Shop, Actions: actions(ShopCreditCreate)},

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
	"shop.Order/CreateOrderSimplify":                {Type: Shop, AuthPartner: Opt, Actions: actions(ShopOrderCreate)},
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

	"shop.Shipment/GetShippingServices":          {Type: Shop, Actions: actions(ShopFulfillmentCreate)},
	"shop.Shipment/CreateFulfillments":           {Type: Shop, Actions: actions(ShopFulfillmentCreate)},
	"shop.Shipment/CancelFulfillment":            {Type: Shop, Actions: actions(ShopFulfillmentCancel)},
	"shop.Shipment/UpdateFulfillmentInfo":        {Type: Shop, Actions: actions(ShopFulfillmentUpdate)},
	"shop.Shipment/UpdateFulfillmentCOD":         {Type: Shop, Actions: actions(ShopFulfillmentUpdate)},
	"shop.Shipment/CreateFulfillmentsFromImport": {Type: Shop},

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

	"shop.Summary/SummarizeShop":    {Type: Shop, Actions: actions(ShopDashboardView)},
	"shop.Summary/SummarizePOS":     {Type: Shop, Actions: actions(ShopDashboardView)},
	"shop.Summary/SummarizeTopShip": {Type: Shop},
	"shop.Summary/CalcBalanceUser":  {Type: Shop, AuthPartner: Opt, Actions: actions(UserBalanceView)},

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
	"shop.Payment/GetExternalPaymentUrl":  {Type: Shop},

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

	"shop.Connection/GetConnections":             {Type: Shop, Actions: actions(ShopConnectionView)},
	"shop.Connection/GetAvailableConnections":    {Type: Shop, Actions: actions(ShopConnectionView)},
	"shop.Connection/GetShopConnections":         {Type: Shop, Actions: actions(ShopConnectionView)},
	"shop.Connection/LoginShopConnection":        {Type: Shop, Actions: actions(ShopConnectionCreate)},
	"shop.Connection/LoginShopConnectionWithOTP": {Type: Shop, Actions: actions(ShopConnectionUpdate)},
	"shop.Connection/RegisterShopConnection":     {Type: Shop, Actions: actions(ShopConnectionCreate)},
	"shop.Connection/DeleteShopConnection":       {Type: Shop, Actions: actions(ShopConnectionDelete)},
	"shop.Connection/UpdateShopConnection":       {Type: Shop, Actions: actions(ShopConnectionUpdate)},

	"shop.Subscription/GetSubscriptionProducts": {Type: Shop, Actions: actions(ShopSubscriptionProductView)},
	"shop.Subscription/GetSubscriptionPlans":    {Type: Shop, Actions: actions(ShopSubscriptionPlanView)},
	"shop.Subscription/GetSubscription":         {Type: Shop, Actions: actions(ShopSubscriptionView)},
	"shop.Subscription/GetSubscriptions":        {Type: Shop, Actions: actions(ShopSubscriptionView)},
	"shop.Subscription/CreateSubscription":      {Type: Shop, Actions: actions(ShopSubscriptionCreate)},
	"shop.Subscription/UpdateSubscriptionInfo":  {Type: Shop, Actions: actions(ShopSubscriptionUpdate)},

	//-- pgevent --//
	"pgevent.Misc/VersionInfo":     {Type: Secret},
	"pgevent.Event/GenerateEvents": {Type: Secret},

	//-- pghandler --//
	"handler.Misc/VersionInfo":   {Type: Secret},
	"handler.Webhook/ResetState": {Type: Secret},

	//-- exporter --//
	"exporter.Misc/VersionInfo": {Type: Secret},

	// -- ticket -- //

	"shop.Ticket/CreateTicket":            {Type: Shop, Actions: actions(ShopTicketCreate)},
	"shop.Ticket/GetTickets":              {Type: Shop, Actions: actions(ShopTicketView)},
	"shop.Ticket/GetTicketsByRefTicketID": {Type: Shop, Actions: actions(ShopTicketView)},
	"shop.Ticket/GetTicket":               {Type: Shop, Actions: actions(ShopTicketView)},
	"shop.Ticket/AssignTicket":            {Type: Shop, Actions: actions(ShopTicketAssign)},
	"shop.Ticket/ConfirmTicket":           {Type: Shop, Actions: actions(ShopTicketUpdate)},
	"shop.Ticket/CloseTicket":             {Type: Shop, Actions: actions(ShopTicketUpdate)},
	"shop.Ticket/ReopenTicket":            {Type: Shop, Actions: actions(ShopTicketUpdate)},
	"shop.Ticket/UpdateTicketRefTicketID": {Type: Shop, Actions: actions(ShopTicketUpdate)},
	"shop.Ticket/CreateTicketComment":     {Type: Shop, Actions: actions(ShopTicketCommentCreate)},
	"shop.Ticket/UpdateTicketComment":     {Type: Shop, Actions: actions(ShopTicketCommentUpdate)},
	"shop.Ticket/DeleteTicketComment":     {Type: Shop, Actions: actions(ShopTicketCommentDelete)},
	"shop.Ticket/GetTicketComments":       {Type: Shop, Actions: actions(ShopTicketCommentView)},
	"shop.Ticket/CreateTicketLabel":       {Type: Shop, Actions: actions(ShopTicketLabelCreate)},
	"shop.Ticket/UpdateTicketLabel":       {Type: Shop, Actions: actions(ShopTicketLabelUpdate)},
	"shop.Ticket/DeleteTicketLabel":       {Type: Shop, Actions: actions(ShopTicketLabelUpdate)},
	"shop.Ticket/GetTicketLabels":         {Type: Shop, Actions: actions(ShopTicketLabelView)},

	"admin.Ticket/CreateTicket":            {Type: EtopAdmin, Actions: actions(AdminTicketCreate)},
	"admin.Ticket/GetTickets":              {Type: EtopAdmin, Actions: actions(AdminTicketView)},
	"admin.Ticket/GetTicketsByRefTicketID": {Type: EtopAdmin, Actions: actions(AdminTicketView)},
	"admin.Ticket/GetTicket":               {Type: EtopAdmin, Actions: actions(AdminTicketView)},
	"admin.Ticket/AssignTicket":            {Type: EtopAdmin, Actions: actions(AdminLeadTicketAssign)},
	"admin.Ticket/ConfirmTicket":           {Type: EtopAdmin, Actions: actions(AdminTicketUpdate)},
	"admin.Ticket/CloseTicket":             {Type: EtopAdmin, Actions: actions(AdminTicketUpdate)},
	"admin.Ticket/ReopenTicket":            {Type: EtopAdmin, Actions: actions(LeadAdminTicketReopen)},
	"admin.Ticket/UpdateTicketRefTicketID": {Type: EtopAdmin, Actions: actions(AdminTicketUpdate)},

	"admin.Ticket/CreateTicketComment": {Type: EtopAdmin, Actions: actions(AdminTicketCommentCreate)},
	"admin.Ticket/UpdateTicketComment": {Type: EtopAdmin, Actions: actions(AdminTicketCommentUpdate)},
	"admin.Ticket/DeleteTicketComment": {Type: EtopAdmin, Actions: actions(AdminTicketCommentDelete)},
	"admin.Ticket/GetTicketComments":   {Type: EtopAdmin, Actions: actions(AdminTicketCommentView)},

	"admin.Ticket/CreateTicketLabel": {Type: EtopAdmin, Actions: actions(AdminTicketLabelCreate)},
	"admin.Ticket/UpdateTicketLabel": {Type: EtopAdmin, Actions: actions(AdminTicketLabelUpdate)},
	"admin.Ticket/DeleteTicketLabel": {Type: EtopAdmin, Actions: actions(AdminTicketLabelUpdate)},
	"admin.Ticket/GetTicketLabels":   {Type: EtopAdmin, Actions: actions(AdminTicketLabelView)},

	// Etelecom User Setting
	"admin.Etelecom/GetUserSettings":   {Type: EtopAdmin, Actions: actions(AdminEtelecomUserSettingView)},
	"admin.Etelecom/UpdateUserSetting": {Type: EtopAdmin, Actions: actions(AdminEtelecomUserSettingUpdate)},

	"etop.Ticket/GetTicketLabels": {Type: Public},

	// -- contact -- //
	"shop.Contact/GetContact":    {Type: Shop},
	"shop.Contact/GetContacts":   {Type: Shop},
	"shop.Contact/CreateContact": {Type: Shop},
	"shop.Contact/UpdateContact": {Type: Shop},
	"shop.Contact/DeleteContact": {Type: Shop},

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
	"etop.AccountRelationship/ResendInvitation":   {Type: Shop, Actions: actions(RelationshipInvitationCreate)},
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

	// Transaction
	"shop.Transaction/GetTransactions": {Type: Shop, Actions: actions(ShopTransactionView)},
	"shop.Transaction/GetTransaction":  {Type: Shop, Actions: actions(ShopTransactionView)},

	// Etelecom
	"shop.Etelecom/CreateExtension":               {Type: Shop, Actions: actions(ShopExtensionCreate)},
	"shop.Etelecom/CreateExtensionBySubscription": {Type: Shop, Actions: actions(ShopExtensionCreate)},
	"shop.Etelecom/ExtendExtension":               {Type: Shop, Actions: actions(ShopExtensionCreate)},
	"shop.Etelecom/GetExtensions":                 {Type: Shop, Actions: actions(ShopExtensionView)},
	"shop.Etelecom/CreateUserAndAssignExtension":  {Type: Shop, Actions: actions(ShopExtensionCreate)},

	"shop.Etelecom/GetHotlines": {Type: Shop, Actions: actions(ShopHotlineView)},

	"shop.Etelecom/GetCallLogs":   {Type: Shop, Actions: actions(ShopCallLogView)},
	"shop.Etelecom/CreateCallLog": {Type: Shop, Actions: actions(ShopCallLogCreate)},

	"shop.Etelecom/SummaryEtelecom": {Type: Shop, Actions: actions(ShopDashboardView)},

	// Etelecom User Setting
	"etelecom.User/GetUserSetting": {Type: Shop, Actions: actions(ShopEtelecomUserSettingView)},

	// Setting
	"shop.Setting/CreateSetting": {Type: Shop},
	"shop.Setting/UpdateSetting": {Type: Shop},
	"shop.Setting/GetSetting":    {Type: Shop},

	// -- Fabo --
	"fabo.Page/ConnectPages":                                 {Type: Shop, Actions: actions(FbFanpageCreate)},
	"fabo.Page/RemovePages":                                  {Type: Shop, Auth: User, IncludeFaboInfo: true, Actions: actions(FbFanpageDelete)},
	"fabo.Page/ListPages":                                    {Type: Shop, Auth: User, IncludeFaboInfo: true, Actions: actions(FbFanpageView)},
	"fabo.Page/CheckPermissions":                             {Type: Shop, Auth: User, IncludeFaboInfo: true, Actions: actions(FbFanpageView)},
	"fabo.Page/ListPosts":                                    {Type: Shop, Auth: User, IncludeFaboInfo: true, Actions: actions(FbPostCreate)},
	"fabo.CustomerConversation/ListCustomerConversations":    {Type: Shop, Auth: User, IncludeFaboInfo: true, Actions: actions(FbCommentView, FbMessageView)},
	"fabo.CustomerConversation/SearchCustomerConversations":  {Type: Shop, Auth: User, IncludeFaboInfo: true, Actions: actions(FbCommentView, FbMessageView)},
	"fabo.CustomerConversation/GetCustomerConversationByID":  {Type: Shop, Auth: User, IncludeFaboInfo: true, Actions: actions(FbCommentView, FbMessageView)},
	"fabo.CustomerConversation/ListMessages":                 {Type: Shop, Auth: User, IncludeFaboInfo: true, Actions: actions(FbMessageView)},
	"fabo.CustomerConversation/ListCommentsByExternalPostID": {Type: Shop, Auth: User, IncludeFaboInfo: true, Actions: actions(FbCommentView)},
	"fabo.CustomerConversation/UpdateReadStatus":             {Type: Shop, Auth: User, IncludeFaboInfo: true, Actions: actions(FbCommentView, FbMessageView)},
	"fabo.CustomerConversation/SendMessage":                  {Type: Shop, Auth: User, IncludeFaboInfo: true, Actions: actions(FbMessageCreate)},
	"fabo.CustomerConversation/SendComment":                  {Type: Shop, Auth: User, IncludeFaboInfo: true, Actions: actions(FbCommentCreate)},
	"fabo.CustomerConversation/CreatePost":                   {Type: Shop, Auth: User, IncludeFaboInfo: true, Actions: actions(FbPostCreate)},
	"fabo.CustomerConversation/MessageTemplateVariables":     {Type: Shop, Auth: User, IncludeFaboInfo: true, Actions: actions(FbMessageTemplateView)},
	"fabo.CustomerConversation/MessageTemplates":             {Type: Shop, Auth: User, IncludeFaboInfo: true, Actions: actions(FbMessageTemplateView)},
	"fabo.CustomerConversation/CreateMessageTemplate":        {Type: Shop, Auth: User, IncludeFaboInfo: true, Actions: actions(FbMessageTemplateCreate)},
	"fabo.CustomerConversation/UpdateMessageTemplate":        {Type: Shop, Auth: User, IncludeFaboInfo: true, Actions: actions(FbMessageTemplateUpdate)},
	"fabo.CustomerConversation/DeleteMessageTemplate":        {Type: Shop, Auth: User, IncludeFaboInfo: true, Actions: actions(FbMessageTemplateDelete)},
	"fabo.CustomerConversation/LikeOrUnLikeComment":          {Type: Shop, Auth: User, IncludeFaboInfo: true},
	"fabo.CustomerConversation/HideOrUnHideComment":          {Type: Shop, Auth: User, IncludeFaboInfo: true},
	"fabo.CustomerConversation/SendPrivateReply":             {Type: Shop, Auth: User, IncludeFaboInfo: true},
	"fabo.CustomerConversation/ListLiveVideos":               {Type: Shop, Auth: User, IncludeFaboInfo: true},

	"fabo.Shop/CreateTag": {Type: Shop, Auth: User, IncludeFaboInfo: true, Actions: actions(FbShopTagCreate)},
	"fabo.Shop/UpdateTag": {Type: Shop, Auth: User, IncludeFaboInfo: true, Actions: actions(FbShopTagUpdate)},
	"fabo.Shop/DeleteTag": {Type: Shop, Auth: User, IncludeFaboInfo: true, Actions: actions(FbShopTagDelete)},
	"fabo.Shop/GetTags":   {Type: Shop, Auth: User, IncludeFaboInfo: true, Actions: actions(FbShopTagView)},

	"fabo.Customer/UpdateTags": {Type: Shop, Auth: User, IncludeFaboInfo: true, Actions: actions(FbUserUpdate)},

	// -- Fabo Customer --
	"fabo.Customer/CreateFbUserCustomer":     {Type: Shop, Actions: actions(FbUserCreate)},
	"fabo.Customer/ListFbUsers":              {Type: Shop},
	"fabo.Customer/GetFbUser":                {Type: Shop, Actions: actions(FbUserView)},
	"fabo.Customer/ListCustomersWithFbUsers": {Type: Shop},

	// -- Fabo external
	"fabo.ExtraShipment/CustomerReturnRate": {Type: Shop},

	// -- Fabo summary
	"fabo.Summary/SummaryShop": {Type: Shop, Auth: User, Actions: actions(ShopDashboardView)},

	// -- Fabo demo
	"fabo.Demo/ListLiveVideos": {Type: Shop},
	"fabo.Demo/ListFeeds":      {Type: Shop},
}

func actions(actions ...permission.ActionType) (actionsResult []permission.ActionType) {
	actionsResult = append(actionsResult, actions...)
	return
}
