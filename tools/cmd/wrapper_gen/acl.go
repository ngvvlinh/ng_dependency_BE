package main

import (
	"strings"

	"etop.vn/backend/pkg/etop/authorize/permission"
)

func init() {
	ACL2 := make(map[string]*permission.PermissionDecl)
	for key, p := range ACL {
		key2 := strings.TrimSpace(key)
		idx := strings.LastIndex(key2, "/")
		key2 = key2[:idx] + "Service" + key2[idx:]

		delete(ACL, key)
		ACL2[key2] = p
	}
	ACL = ACL2
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

	Owner = permission.RoleOwner // Can add/remove user
	Admin = permission.RoleAdmin // Can update settings
	Staff = permission.RoleStaff // Can create orders
	_____ = ""                   // Viewer, readonly

	User              = permission.User
	APIKey            = permission.APIKey
	APIPartnerShopKey = permission.APIPartnerShopKey

	Req = permission.Required
	Opt = permission.Optional
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
	"etop.User/Login":                    {Type: Public},
	"etop.User/ResetPassword":            {Type: Public},
	"etop.User/ChangePasswordUsingToken": {Type: Public},
	"etop.User/ChangePassword":           {Type: CurUsr},

	"etop.User/SessionInfo":        {Type: CurUsr},
	"etop.User/SwitchAccount":      {Type: CurUsr},
	"etop.User/SendSTokenEmail":    {Type: CurUsr},
	"etop.User/UpgradeAccessToken": {Type: CurUsr},

	"etop.User/UpdatePermission": {Type: CurUsr},

	"etop.User/SendEmailVerification": {Type: CurUsr},
	"etop.User/SendPhoneVerification": {Type: CurUsr},
	"etop.User/VerifyEmailUsingToken": {Type: CurUsr},
	"etop.User/VerifyPhoneUsingToken": {Type: CurUsr},
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

	"admin.Fulfillment/GetFulfillment":    {Type: EtopAdmin},
	"admin.Fulfillment/GetFulfillments":   {Type: EtopAdmin},
	"admin.Fulfillment/UpdateFulfillment": {Type: EtopAdmin},

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

	"shop.Account/UpdateShop": {Type: Shop, Role: Admin},
	"shop.Account/DeleteShop": {Type: Shop, Role: Owner},

	// permission: admin

	"shop.Account/SetDefaultAddress": {Type: Shop, Role: Staff},
	"shop.Account/GetBalanceShop":    {Type: Shop, Role: Staff},

	"shop.Account/CreateExternalAccountAhamove":                   {Type: Shop},
	"shop.Account/GetExternalAccountAhamove":                      {Type: Shop},
	"shop.Account/RequestVerifyExternalAccountAhamove":            {Type: Shop},
	"shop.Account/UpdateExternalAccountAhamoveVerification":       {Type: Shop},
	"shop.Account/UpdateExternalAccountAhamoveVerificationImages": {Type: Shop},

	"shop.ExternalAccount/GetExternalAccountHaravan":                           {Type: Shop},
	"shop.ExternalAccount/CreateExternalAccountHaravan":                        {Type: Shop},
	"shop.ExternalAccount/UpdateExternalAccountHaravanToken":                   {Type: Shop},
	"shop.ExternalAccount/UpdateExternalAccountHaravan":                        {Type: Shop},
	"shop.ExternalAccount/ConnectCarrierServiceExternalAccountHaravan":         {Type: Shop},
	"shop.ExternalAccount/DeleteConnectedCarrierServiceExternalAccountHaravan": {Type: Shop},

	"shop.Browse/BrowseCategories":    {Type: Shop},
	"shop.Browse/BrowseProduct":       {Type: Shop},
	"shop.Browse/BrowseVariant":       {Type: Shop},
	"shop.Browse/BrowseProducts":      {Type: Shop},
	"shop.Browse/BrowseVariants":      {Type: Shop},
	"shop.Browse/BrowseProductsByIDs": {Type: Shop},
	"shop.Browse/BrowseVariantsByIDs": {Type: Shop},

	"shop.Collection/GetCollection":             {Type: Shop, Role: _____},
	"shop.Collection/GetCollections":            {Type: Shop, Role: _____},
	"shop.Collection/CreateCollection":          {Type: Shop, Role: Staff},
	"shop.Collection/UpdateCollection":          {Type: Shop, Role: Staff},
	"shop.Collection/GetCollectionsByProductID": {Type: Shop},

	"shop.Customer/CreateCustomer":          {Type: Shop},
	"shop.Customer/UpdateCustomer":          {Type: Shop},
	"shop.Customer/DeleteCustomer":          {Type: Shop},
	"shop.Customer/GetCustomer":             {Type: Shop},
	"shop.Customer/GetCustomerDetails":      {Type: Shop},
	"shop.Customer/GetCustomers":            {Type: Shop},
	"shop.Customer/GetCustomersByIDs":       {Type: Shop},
	"shop.Customer/BatchSetCustomersStatus": {Type: Shop},

	"shop.Customer/GetCustomerAddresses":      {Type: Shop},
	"shop.Customer/CreateCustomerAddress":     {Type: Shop},
	"shop.Customer/UpdateCustomerAddress":     {Type: Shop},
	"shop.Customer/DeleteCustomerAddress":     {Type: Shop},
	"shop.Customer/SetDefaultCustomerAddress": {Type: Shop},

	"shop.Customer/AddCustomersToGroup":      {Type: Shop},
	"shop.Customer/RemoveCustomersFromGroup": {Type: Shop},

	"shop.CustomerGroup/CreateCustomerGroup": {Type: Shop},
	"shop.CustomerGroup/GetCustomerGroup":    {Type: Shop},
	"shop.CustomerGroup/GetCustomerGroups":   {Type: Shop},
	"shop.CustomerGroup/UpdateCustomerGroup": {Type: Shop},

	"shop.Category/GetCategory":    {Type: Shop},
	"shop.Category/GetCategories":  {Type: Shop},
	"shop.Category/CreateCategory": {Type: Shop},
	"shop.Category/UpdateCategory": {Type: Shop},
	"shop.Category/DeleteCategory": {Type: Shop},

	"shop.Product/GetProduct":              {Type: Shop, Role: _____, AuthPartner: Opt},
	"shop.Product/GetProducts":             {Type: Shop, Role: _____, AuthPartner: Opt},
	"shop.Product/GetProductsByIDs":        {Type: Shop, Role: _____, AuthPartner: Opt},
	"shop.Product/RemoveProducts":          {Type: Shop, Role: Staff},
	"shop.Product/UpdateProduct":           {Type: Shop, Role: Staff},
	"shop.Product/UpdateProductsStatus":    {Type: Shop, Role: Staff},
	"shop.Product/UpdateProductsTags":      {Type: Shop, Role: Staff},
	"shop.Product/UpdateProductImages":     {Type: Shop, Role: Staff},
	"shop.Product/UpdateProductMetaFields": {Type: Shop, Role: Staff},
	"shop.Product/CreateProduct":           {Type: Shop, Role: Staff},
	"shop.Product/UpdateProductStatus":     {Type: Shop, Role: Staff},
	"shop.Product/AddProducts":             {Type: Shop, Role: Staff, AuthPartner: Opt},
	"shop.Product/UpdateProductCategory":   {Type: Shop},
	"shop.Product/RemoveProductCategory":   {Type: Shop},
	"shop.Product/AddProductCollection":    {Type: Shop},
	"shop.Product/RemoveProductCollection": {Type: Shop},

	"shop.Product/GetVariant":                  {Type: Shop, Role: _____, AuthPartner: Opt},
	"shop.Product/GetVariants":                 {Type: Shop, Role: _____, AuthPartner: Opt},
	"shop.Product/GetVariantsByIDs":            {Type: Shop, Role: _____, AuthPartner: Opt},
	"shop.Product/CreateVariant":               {Type: Shop, Role: Staff},
	"shop.Product/AddVariants":                 {Type: Shop, Role: Staff, AuthPartner: Opt},
	"shop.Product/RemoveVariants":              {Type: Shop, Role: Staff},
	"shop.Product/UpdateVariant":               {Type: Shop, Role: Staff},
	"shop.Product/UpdateVariants":              {Type: Shop, Role: Staff},
	"shop.Product/UpdateVariantsStatus":        {Type: Shop, Role: Staff},
	"shop.Product/UpdateVariantsTags":          {Type: Shop, Role: Staff},
	"shop.Product/UpdateVariantImages":         {Type: Shop, Role: Staff},
	"shop.Product/UpdateVariantAttributes":     {Type: Shop, Role: Staff},
	"shop.ProductSource/GetShopProductSources": {Type: Shop, Role: _____, AuthPartner: Opt},
	"shop.ProductSource/CreateProductSource":   {Type: Shop, Role: Staff},
	"shop.ProductSource/ConnectProductSource":  {Type: Shop, Role: Staff},

	"shop.ProductSource/GetProductSourceCategories":  {Type: Shop, Role: _____, AuthPartner: Opt},
	"shop.ProductSource/GetProductSourceCategory":    {Type: Shop, Role: _____, AuthPartner: Opt},
	"shop.ProductSource/CreateVariant":               {Type: Shop, Role: Staff, AuthPartner: Opt, Rename: "DeprecatedCreateVariant"},
	"shop.ProductSource/CreateProductSourceCategory": {Type: Shop, Role: Staff},
	"shop.ProductSource/UpdateProductsPSCategory":    {Type: Shop, Role: Staff},
	"shop.ProductSource/UpdateProductSourceCategory": {Type: Shop, Role: Staff},
	"shop.ProductSource/RemoveProductSourceCategory": {Type: Shop, Role: Admin},

	"shop.Price/GetPriceRules":    {Type: Shop, Role: _____},
	"shop.Price/UpdatePriceRules": {Type: Shop, Role: _____},

	"shop.Order/CreateOrder":                        {Type: Shop, Role: Staff, AuthPartner: Opt},
	"shop.Order/GetOrder":                           {Type: Shop, Role: _____, AuthPartner: Opt},
	"shop.Order/GetOrders":                          {Type: Shop, Role: _____, AuthPartner: Opt},
	"shop.Order/GetOrdersByIDs":                     {Type: Shop, Role: _____, AuthPartner: Opt},
	"shop.Order/GetOrdersByReceiptID":               {Type: Shop, Role: _____, AuthPartner: Opt},
	"shop.Order/UpdateOrder":                        {Type: Shop, Role: Staff, AuthPartner: Opt},
	"shop.Order/UpdateOrdersStatus":                 {Type: Shop, Role: Staff, AuthPartner: Opt},
	"shop.Order/ConfirmOrderAndCreateFulfillments":  {Type: Shop, Role: Staff, AuthPartner: Opt},
	"shop.Order/ConfirmOrdersAndCreateFulfillments": {Type: Shop, Role: Staff, AuthPartner: Opt},
	"shop.Order/CancelOrder":                        {Type: Shop, Role: Staff, AuthPartner: Opt},
	"shop.Order/UpdateOrderPaymentStatus":           {Type: Shop, Role: Staff, AuthPartner: Opt},

	"shop.Fulfillment/GetPublicExternalShippingServices": {Type: Public, Role: _____},
	"shop.Fulfillment/GetPublicFulfillment":              {Type: Public, Role: _____},
	"shop.Fulfillment/GetExternalShippingServices":       {Type: Shop, Role: _____, AuthPartner: Opt},
	"shop.Fulfillment/CancelFulfillment":                 {Type: Shop, Role: Staff, AuthPartner: Opt},
	"shop.Fulfillment/CreateFulfillmentsForOrder":        {Type: Shop, Role: Staff},
	"shop.Fulfillment/GetFulfillment":                    {Type: Shop, Role: _____, AuthPartner: Opt},
	"shop.Fulfillment/GetFulfillments":                   {Type: Shop, Role: _____, AuthPartner: Opt},
	"shop.Fulfillment/UpdateFulfillmentsShippingState":   {Type: Shop, Role: Staff},

	"shop.Shipnow/GetShipnowFulfillment":     {Type: Shop},
	"shop.Shipnow/GetShipnowFulfillments":    {Type: Shop},
	"shop.Shipnow/CreateShipnowFulfillment":  {Type: Shop},
	"shop.Shipnow/ConfirmShipnowFulfillment": {Type: Shop},
	"shop.Shipnow/UpdateShipnowFulfillment":  {Type: Shop},
	"shop.Shipnow/CancelShipnowFulfillment":  {Type: Shop},
	"shop.Shipnow/GetShipnowServices":        {Type: Shop},

	"shop.Brand/GetBrand":  {Type: Shop, Role: _____, AuthPartner: Opt},
	"shop.Brand/GetBrands": {Type: Shop, Role: _____, AuthPartner: Opt},

	"shop.History/GetFulfillmentHistory": {Type: Shop, Role: _____, AuthPartner: Opt},

	"shop.MoneyTransaction/GetMoneyTransaction":  {Type: Shop, Role: _____},
	"shop.MoneyTransaction/GetMoneyTransactions": {Type: Shop, Role: _____},

	"shop.Summary/SummarizeFulfillments": {Type: Shop, Role: Admin},
	"shop.Summary/SummarizePOS":          {Type: Shop},
	"shop.Summary/CalcBalanceShop":       {Type: Shop, AuthPartner: Opt},

	"shop.Export/GetExports":    {Type: Shop, Auth: User},
	"shop.Export/RequestExport": {Type: Shop},

	"shop.Notification/CreateDevice":        {Type: Shop},
	"shop.Notification/DeleteDevice":        {Type: Shop},
	"shop.Notification/GetNotification":     {Type: Shop},
	"shop.Notification/GetNotifications":    {Type: Shop},
	"shop.Notification/UpdateNotifications": {Type: Shop},

	"shop.Authorize/GetAuthorizedPartners": {Type: Shop},
	"shop.Authorize/GetAvailablePartners":  {Type: Shop},
	"shop.Authorize/AuthorizePartner":      {Type: Shop},

	//-- Receipt --//
	"shop.Receipt/CreateReceipt":           {Type: Shop},
	"shop.Receipt/UpdateReceipt":           {Type: Shop},
	"shop.Receipt/DeleteReceipt":           {Type: Shop},
	"shop.Receipt/GetReceipt":              {Type: Shop},
	"shop.Receipt/GetReceipts":             {Type: Shop},
	"shop.Receipt/GetReceiptsByLedgerType": {Type: Shop},
	"shop.Receipt/ConfirmReceipt":          {Type: Shop},
	"shop.Receipt/CancelReceipt":           {Type: Shop},

	"shop.Trading/TradingPaymentOrder": {Type: Shop},
	"shop.Trading/TradingGetProduct":   {Type: Shop},
	"shop.Trading/TradingGetProducts":  {Type: Shop},
	"shop.Trading/TradingCreateOrder":  {Type: Shop},
	"shop.Trading/TradingGetOrder":     {Type: Shop},
	"shop.Trading/TradingGetOrders":    {Type: Shop},

	"shop.Payment/PaymentTradingOrder":    {Type: Shop},
	"shop.Payment/PaymentCheckReturnData": {Type: Shop},

	"shop.Inventory/CreateInventoryVoucher":  {Type: Shop},
	"shop.Inventory/ConfirmInventoryVoucher": {Type: Shop},
	"shop.Inventory/CancelInventoryVoucher":  {Type: Shop},
	"shop.Inventory/UpdateInventoryVoucher":  {Type: Shop},
	"shop.Inventory/AdjustInventoryQuantity": {Type: Shop},

	"shop.Inventory/GetInventories":             {Type: Shop},
	"shop.Inventory/GetInventoriesByVariantIDs": {Type: Shop},
	"shop.Inventory/GetInventoryVouchersByIDs":  {Type: Shop},
	"shop.Inventory/GetInventoryVouchers":       {Type: Shop},
	"shop.Inventory/GetInventoryVoucher":        {Type: Shop},
	"shop.Inventory/GetInventory":               {Type: Shop},

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

	// vendor:
	"shop.Vendor/GetVendor":       {Type: Shop},
	"shop.Vendor/GetVendors":      {Type: Shop},
	"shop.Vendor/GetVendorsByIDs": {Type: Shop},
	"shop.Vendor/CreateVendor":    {Type: Shop},
	"shop.Vendor/UpdateVendor":    {Type: Shop},
	"shop.Vendor/DeleteVendor":    {Type: Shop},

	// carrier:
	"shop.Carrier/GetCarrier":       {Type: Shop},
	"shop.Carrier/GetCarriers":      {Type: Shop},
	"shop.Carrier/GetCarriersByIDs": {Type: Shop},
	"shop.Carrier/CreateCarrier":    {Type: Shop},
	"shop.Carrier/UpdateCarrier":    {Type: Shop},
	"shop.Carrier/DeleteCarrier":    {Type: Shop},

	// Ledger:
	"shop.Ledger/GetLedger":    {Type: Shop},
	"shop.Ledger/GetLedgers":   {Type: Shop},
	"shop.Ledger/CreateLedger": {Type: Shop},
	"shop.Ledger/UpdateLedger": {Type: Shop},
	"shop.Ledger/DeleteLedger": {Type: Shop},
}
