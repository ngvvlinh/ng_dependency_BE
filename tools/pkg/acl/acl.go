package acl

import (
	"strings"

	"etop.vn/backend/pkg/etop/authorize/permission"
)

func init() {
	ACL2 := make(map[string]*permission.PermissionDecl)
	for key, p := range ACL {
		key2 := ConvertKey(key)
		delete(ACL, key)
		ACL2[key2] = p
	}
	ACL = ACL2
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
	Custom     = permission.Custom
	Secret     = permission.Secret

	User              = permission.User
	APIKey            = permission.APIKey
	APIPartnerShopKey = permission.APIPartnerShopKey

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

	ShopOrderCreate  permission.ActionType = "shop/order:create"
	ShopOrderConfirm permission.ActionType = "shop/order:confirm"
	ShopOrderUpdate  permission.ActionType = "shop/order:update"
	ShopOrderCancel  permission.ActionType = "shop/order:cancel"
	ShopOrderView    permission.ActionType = "shop/order:view"

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

	RelationshipInvitationCreate   permission.ActionType = "relationship/invitation:create"
	RelationshipInvitationView     permission.ActionType = "relationship/invitation:view"
	RelationshipInvitationDelete   permission.ActionType = "relationship/invitation:delete"
	RelationshipPermissionUpdate   permission.ActionType = "relationship/permission:update"
	RelationshipRelationshipUpdate permission.ActionType = "relationship/relationship:update"
	RelationshipRelationshipView   permission.ActionType = "relationship/relationship:view"
	RelationshipUserRemove         permission.ActionType = "relationship/user:remove"

	ShopStocktakeCreate  permission.ActionType = "shop/stocktake:create"
	ShopStocktakeUpdate  permission.ActionType = "shop/stocktake:update"
	ShopStocktakeConfirm permission.ActionType = "shop/stocktake:confirm"
	ShopStocktakeCancel  permission.ActionType = "shop/stocktake:cancel"
	ShopStocktakeView    permission.ActionType = "shop/stocktake:view"

	ShopRefundCreate  permission.ActionType = "shop/refund:create"
	ShopRefundUpdate  permission.ActionType = "shop/refund:update"
	ShopRefundConfirm permission.ActionType = "shop/refund:confirm"
	ShopRefundCancel  permission.ActionType = "shop/refund:cancel"
	ShopRefundView    permission.ActionType = "shop/refund:view"

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
)

// ACL declares access control list
var ACL = map[string]*permission.PermissionDecl{
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

	"etop.User/UpdatePermission": {Type: CurUsr},

	"etop.User/SendEmailVerification": {Type: CurUsr},
	"etop.User/SendPhoneVerification": {Type: Custom},
	"etop.User/VerifyEmailUsingToken": {Type: CurUsr},
	"etop.User/VerifyPhoneUsingToken": {Type: Custom},
	"etop.User/UpdateReferenceUser":   {Type: CurUsr},
	"etop.User/UpdateReferenceSale":   {Type: CurUsr},

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

	//-- external: partner --//

	"ext/partner.Misc/CurrentAccount":   {Type: Partner, Auth: APIKey},
	"ext/partner.History/GetChanges":    {Type: Partner, Auth: APIKey},
	"ext/partner.Misc/GetLocationList":  {Type: Partner, Auth: APIKey},
	"ext/partner.Shop/AuthorizeShop":    {Type: Partner, Auth: APIKey},
	"ext/partner.Webhook/CreateWebhook": {Type: Partner, Auth: APIKey},
	"ext/partner.Webhook/GetWebhooks":   {Type: Partner, Auth: APIKey},
	"ext/partner.Webhook/DeleteWebhook": {Type: Partner, Auth: APIKey},

	//-- external: partner using partnerShopKey --//
	"ext/partner.Shop/CurrentShop":               {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.Shipping/GetShippingServices":   {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.Shipping/CreateAndConfirmOrder": {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.Shipping/CancelOrder":           {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.Shipping/GetOrder":              {Type: Shop, Auth: APIPartnerShopKey},
	"ext/partner.Shipping/GetFulfillment":        {Type: Shop, Auth: APIPartnerShopKey},

	"ext/partner.Customer/GetCustomers": {Type: Shop, Auth: APIPartnerShopKey},

	"ext/partner.Product/GetProducts": {Type: Shop, Auth: APIPartnerShopKey},

	"ext/partner.Variant/GetVariants": {Type: Shop, Auth: APIPartnerShopKey},

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

	//-- integration --//

	"integration.Integration/Init":            {Type: Public},
	"integration.Integration/RequestLogin":    {Type: Protected, AuthPartner: Req, Captcha: "1"},
	"integration.Integration/LoginUsingToken": {Type: Protected, AuthPartner: Req},
	"integration.Integration/Register":        {Type: Protected, AuthPartner: Req},
	"integration.Integration/GrantAccess":     {Type: CurUsr, AuthPartner: Req},
	"integration.Integration/SessionInfo":     {Type: Protected, AuthPartner: Req},

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

	"admin.Fulfillment/GetFulfillment":                 {Type: EtopAdmin},
	"admin.Fulfillment/GetFulfillments":                {Type: EtopAdmin},
	"admin.Fulfillment/UpdateFulfillment":              {Type: EtopAdmin},
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
	"admin.MoneyTransaction/ConfirmMoneyTransactionShippingExternal":     {Type: EtopAdmin},
	"admin.MoneyTransaction/ConfirmMoneyTransactionShippingExternals":    {Type: EtopAdmin},
	"admin.MoneyTransaction/UpdateMoneyTransactionShippingExternal":      {Type: EtopAdmin},
	"admin.MoneyTransaction/GetMoneyTransactionShippingEtop":             {Type: EtopAdmin},
	"admin.MoneyTransaction/GetMoneyTransactionShippingEtops":            {Type: EtopAdmin},
	"admin.MoneyTransaction/CreateMoneyTransactionShippingEtop":          {Type: EtopAdmin},
	"admin.MoneyTransaction/UpdateMoneyTransactionShippingEtop":          {Type: EtopAdmin},
	"admin.MoneyTransaction/DeleteMoneyTransactionShippingEtop":          {Type: EtopAdmin},
	"admin.MoneyTransaction/ConfirmMoneyTransactionShippingEtop":         {Type: EtopAdmin},

	"admin.Shop/GetShop":  {Type: EtopAdmin},
	"admin.Shop/GetShops": {Type: EtopAdmin},

	"admin.Credit/GetCredit":     {Type: EtopAdmin},
	"admin.Credit/GetCredits":    {Type: EtopAdmin},
	"admin.Credit/CreateCredit":  {Type: EtopAdmin},
	"admin.Credit/UpdateCredit":  {Type: EtopAdmin},
	"admin.Credit/ConfirmCredit": {Type: EtopAdmin},
	"admin.Credit/DeleteCredit":  {Type: EtopAdmin},

	"admin.Notification/CreateNotifications": {Type: EtopAdmin},

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

	"shop.ExternalAccount/GetExternalAccountHaravan":                           {Type: Shop},
	"shop.ExternalAccount/CreateExternalAccountHaravan":                        {Type: Shop, Actions: actions(ShopExternalAccountManage)},
	"shop.ExternalAccount/UpdateExternalAccountHaravanToken":                   {Type: Shop, Actions: actions(ShopExternalAccountManage)},
	"shop.ExternalAccount/UpdateExternalAccountHaravan":                        {Type: Shop, Actions: actions(ShopExternalAccountManage)},
	"shop.ExternalAccount/ConnectCarrierServiceExternalAccountHaravan":         {Type: Shop, Actions: actions(ShopExternalAccountManage)},
	"shop.ExternalAccount/DeleteConnectedCarrierServiceExternalAccountHaravan": {Type: Shop, Actions: actions(ShopExternalAccountManage)},

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
	"shop.Order/CancelOrder":                        {Type: Shop, AuthPartner: Opt, Actions: actions(ShopOrderCancel)},
	"shop.Order/UpdateOrderPaymentStatus":           {Type: Shop, AuthPartner: Opt, Actions: actions(ShopOrderUpdate)},
	"shop.Order/UpdateOrderShippingInfo":            {Type: Shop, AuthPartner: Opt, Actions: actions(ShopOrderUpdate)},

	"shop.Fulfillment/GetPublicExternalShippingServices": {Type: Public},
	"shop.Fulfillment/GetPublicFulfillment":              {Type: Public},
	"shop.Fulfillment/GetExternalShippingServices":       {Type: Shop, AuthPartner: Opt, Actions: actions(ShopFulfillmentCreate)},
	"shop.Fulfillment/CancelFulfillment":                 {Type: Shop, AuthPartner: Opt, Actions: actions(ShopFulfillmentCancel)},
	"shop.Fulfillment/CreateFulfillmentsForOrder":        {Type: Shop, Actions: actions(ShopFulfillmentCreate)},
	"shop.Fulfillment/GetFulfillment":                    {Type: Shop, AuthPartner: Opt, Actions: actions(ShopFulfillmentView)},
	"shop.Fulfillment/GetFulfillments":                   {Type: Shop, AuthPartner: Opt, Actions: actions(ShopFulfillmentView)},
	"shop.Fulfillment/UpdateFulfillmentsShippingState":   {Type: Shop},

	"shop.Shipment/GetShippingServices": {Type: Shop, Actions: actions(ShopFulfillmentCreate)},
	"shop.Shipment/CreateFulfillments":  {Type: Shop, Actions: actions(ShopFulfillmentCreate)},

	"shop.Shipnow/GetShipnowFulfillment":     {Type: Shop, Actions: actions(ShopShipNowView)},
	"shop.Shipnow/GetShipnowFulfillments":    {Type: Shop, Actions: actions(ShopShipNowView)},
	"shop.Shipnow/CreateShipnowFulfillment":  {Type: Shop, Actions: actions(ShopShipNowCreate)},
	"shop.Shipnow/ConfirmShipnowFulfillment": {Type: Shop, Actions: actions(ShopShipNowConfirm)},
	"shop.Shipnow/UpdateShipnowFulfillment":  {Type: Shop, Actions: actions(ShopShipNowUpdate)},
	"shop.Shipnow/CancelShipnowFulfillment":  {Type: Shop, Actions: actions(ShopShipNowCancel)},
	"shop.Shipnow/GetShipnowServices":        {Type: Shop, Actions: actions(ShopShipNowView)},

	"shop.Brand/GetBrand":  {Type: Shop, AuthPartner: Opt, Actions: actions(ShopProductBasicInfoView)},
	"shop.Brand/GetBrands": {Type: Shop, AuthPartner: Opt, Actions: actions(ShopProductBasicInfoView)},

	"shop.History/GetFulfillmentHistory": {Type: Shop, AuthPartner: Opt},

	"shop.MoneyTransaction/GetMoneyTransaction":  {Type: Shop, Actions: actions(ShopMoneyTransactionView)},
	"shop.MoneyTransaction/GetMoneyTransactions": {Type: Shop, Actions: actions(ShopMoneyTransactionView)},

	"shop.Summary/SummarizeFulfillments": {Type: Shop, Actions: actions(ShopDashboardView)},
	"shop.Summary/SummarizePOS":          {Type: Shop, Actions: actions(ShopDashboardView)},
	"shop.Summary/SummarizeTopShip":      {Type: Shop},
	"shop.Summary/CalcBalanceShop":       {Type: Shop, AuthPartner: Opt, Actions: actions(ShopDashboardView)},

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
	"shop.Stocktake/UpdateStocktake":    {Type: Shop, Actions: actions(ShopStocktakeUpdate)},
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
	"etop.AccountRelationship/RemoveUser":         {Type: Shop, Actions: actions(RelationshipUserRemove)},

	"shop.Refund/CreateRefund":  {Type: Shop, Actions: actions(ShopRefundCreate)},
	"shop.Refund/UpdateRefund":  {Type: Shop, Actions: actions(ShopRefundUpdate)},
	"shop.Refund/CancelRefund":  {Type: Shop, Actions: actions(ShopRefundCancel)},
	"shop.Refund/ConfirmRefund": {Type: Shop, Actions: actions(ShopRefundConfirm)},

	"shop.Refund/GetRefund":       {Type: Shop, Actions: actions(ShopRefundView)},
	"shop.Refund/GetRefundsByIDs": {Type: Shop, Actions: actions(ShopRefundView)},
	"shop.Refund/GetRefunds":      {Type: Shop, Actions: actions(ShopRefundView)},
}

func actions(actions ...permission.ActionType) (actionsResult []permission.ActionType) {
	actionsResult = append(actionsResult, actions...)
	return
}
