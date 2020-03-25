package register

import (
	"fmt"

	"etop.vn/backend/cmd/etop-etl/register/table_name"
	addressmodel "etop.vn/backend/com/main/address/model"
	catalogmodel "etop.vn/backend/com/main/catalog/model"
	identitymodel "etop.vn/backend/com/main/identity/model"
	shopmodel "etop.vn/backend/com/main/identity/model"
	inventorymodel "etop.vn/backend/com/main/inventory/model"
	invitationmodel "etop.vn/backend/com/main/invitation/model"
	shopledgermodel "etop.vn/backend/com/main/ledgering/model"
	moneytxmodel "etop.vn/backend/com/main/moneytx/model"
	ordermodel "etop.vn/backend/com/main/ordering/model"
	purchaseorder "etop.vn/backend/com/main/purchaseorder/model"
	purchaserefundmodel "etop.vn/backend/com/main/purchaserefund/model"
	receiptmodel "etop.vn/backend/com/main/receipting/model"
	refundmodel "etop.vn/backend/com/main/refund/model"
	shipnowfulfillmentmodel "etop.vn/backend/com/main/shipnow/model"
	fulfillmentmodel "etop.vn/backend/com/main/shipping/model"
	shopstocktakemodel "etop.vn/backend/com/main/stocktaking/model"
	carriermodel "etop.vn/backend/com/shopping/carrying/model"
	customermodel "etop.vn/backend/com/shopping/customering/model"
	suppliermodel "etop.vn/backend/com/shopping/suppliering/model"
	tradermodel "etop.vn/backend/com/shopping/tradering/model"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/pkg/common/sql/sq"
	"etop.vn/backend/zexp/etl"
	accountconvert "etop.vn/backend/zexp/etl/main/account/convert"
	etlaccountmodel "etop.vn/backend/zexp/etl/main/account/model"
	accountuserconvert "etop.vn/backend/zexp/etl/main/accountuser/convert"
	etlaccountusermodel "etop.vn/backend/zexp/etl/main/accountuser/model"
	addressconvert "etop.vn/backend/zexp/etl/main/address/convert"
	etladdressmodel "etop.vn/backend/zexp/etl/main/address/model"
	customerconvert "etop.vn/backend/zexp/etl/main/customer/convert"
	etlcustomermodel "etop.vn/backend/zexp/etl/main/customer/model"
	fulfillmentconvert "etop.vn/backend/zexp/etl/main/fulfillment/convert"
	etlfulfillmentmodel "etop.vn/backend/zexp/etl/main/fulfillment/model"
	inventoryconvert "etop.vn/backend/zexp/etl/main/inventoryvariant/convert"
	etlinventorymodel "etop.vn/backend/zexp/etl/main/inventoryvariant/model"
	inventoryvoucherconvert "etop.vn/backend/zexp/etl/main/inventoryvoucher/convert"
	etlinventoryvouchermodel "etop.vn/backend/zexp/etl/main/inventoryvoucher/model"
	invitationconvert "etop.vn/backend/zexp/etl/main/invitation/convert"
	etlinvitationmodel "etop.vn/backend/zexp/etl/main/invitation/model"
	moneytxconvert "etop.vn/backend/zexp/etl/main/moneytransactionshipping/convert"
	etlmoneytxmodel "etop.vn/backend/zexp/etl/main/moneytransactionshipping/model"
	orderconvert "etop.vn/backend/zexp/etl/main/order/convert"
	etlordermodel "etop.vn/backend/zexp/etl/main/order/model"
	productshopcollectionconvert "etop.vn/backend/zexp/etl/main/productshopcollection/convert"
	etlproductshopcollectionmodel "etop.vn/backend/zexp/etl/main/productshopcollection/model"
	purchaseorderconvert "etop.vn/backend/zexp/etl/main/purchaseorder/convert"
	etlpurchaseordermodel "etop.vn/backend/zexp/etl/main/purchaseorder/model"
	purchaserefundconvert "etop.vn/backend/zexp/etl/main/purchaserefund/convert"
	etlpurchaserefundmodel "etop.vn/backend/zexp/etl/main/purchaserefund/model"
	receiptconvert "etop.vn/backend/zexp/etl/main/receipt/convert"
	etlreceiptmodel "etop.vn/backend/zexp/etl/main/receipt/model"
	refundconvert "etop.vn/backend/zexp/etl/main/refund/convert"
	etlrefundmodel "etop.vn/backend/zexp/etl/main/refund/model"
	shipnowfulfillmentconvert "etop.vn/backend/zexp/etl/main/shipnowfulfillment/convert"
	etlshipnowfulfillmentmodel "etop.vn/backend/zexp/etl/main/shipnowfulfillment/model"
	shopconvert "etop.vn/backend/zexp/etl/main/shop/convert"
	etlshopmodel "etop.vn/backend/zexp/etl/main/shop/model"
	shopbrandconvert "etop.vn/backend/zexp/etl/main/shopbrand/convert"
	etlshopbrandmodel "etop.vn/backend/zexp/etl/main/shopbrand/model"
	carrierconvert "etop.vn/backend/zexp/etl/main/shopcarrier/convert"
	etlcarriermodel "etop.vn/backend/zexp/etl/main/shopcarrier/model"
	shopcategoryconvert "etop.vn/backend/zexp/etl/main/shopcategory/convert"
	etlshopcategorymodel "etop.vn/backend/zexp/etl/main/shopcategory/model"
	shopcollectionconvert "etop.vn/backend/zexp/etl/main/shopcollection/convert"
	etlshopcollectionmodel "etop.vn/backend/zexp/etl/main/shopcollection/model"
	shopcustomergroupconvert "etop.vn/backend/zexp/etl/main/shopcustomergroup/convert"
	etlshopcustomergroupmodel "etop.vn/backend/zexp/etl/main/shopcustomergroup/model"
	shopcustomergroupcustomerconvert "etop.vn/backend/zexp/etl/main/shopcustomergroupcustomer/convert"
	etlshopcustomergroupcustomermodel "etop.vn/backend/zexp/etl/main/shopcustomergroupcustomer/model"
	shopledgerconvert "etop.vn/backend/zexp/etl/main/shopledger/convert"
	etlshopledgermodel "etop.vn/backend/zexp/etl/main/shopledger/model"
	shopproductconvert "etop.vn/backend/zexp/etl/main/shopproduct/convert"
	etlshopproductmodel "etop.vn/backend/zexp/etl/main/shopproduct/model"
	shopproductcollectionconvert "etop.vn/backend/zexp/etl/main/shopproductcollection/convert"
	etlshopproductcollectionmodel "etop.vn/backend/zexp/etl/main/shopproductcollection/model"
	shopstocktakeconvert "etop.vn/backend/zexp/etl/main/shopstocktake/convert"
	etlshopstocktakemodel "etop.vn/backend/zexp/etl/main/shopstocktake/model"
	shopsupplierconvert "etop.vn/backend/zexp/etl/main/shopsupplier/convert"
	etlshopsuppliermodel "etop.vn/backend/zexp/etl/main/shopsupplier/model"
	shoptraderconvert "etop.vn/backend/zexp/etl/main/shoptrader/convert"
	etlshoptradermodel "etop.vn/backend/zexp/etl/main/shoptrader/model"
	shoptraderaddressconvert "etop.vn/backend/zexp/etl/main/shoptraderaddress/convert"
	etlshoptraderaddressmodel "etop.vn/backend/zexp/etl/main/shoptraderaddress/model"
	shopvariantconvert "etop.vn/backend/zexp/etl/main/shopvariant/convert"
	etlshopvariantmodel "etop.vn/backend/zexp/etl/main/shopvariant/model"
	shopvariantsupplierconvert "etop.vn/backend/zexp/etl/main/shopvariantsupplier/convert"
	etlshopvariantsuppliermodel "etop.vn/backend/zexp/etl/main/shopvariantsupplier/model"
	userconvert "etop.vn/backend/zexp/etl/main/user/convert"
	etlusermodel "etop.vn/backend/zexp/etl/main/user/model"
	"etop.vn/capi/dot"
)

type ETLRegisterFunc func(ng *etl.ETLEngine, DB, dstDB *cmsql.Database, args ...interface{})

var mapFieldRegisters = map[table_name.TableName]ETLRegisterFunc{
	table_name.User:                      registerUser,
	table_name.Account:                   registerAccount,
	table_name.ShopCustomer:              registerShopCustomer,
	table_name.Order:                     registerOrder,
	table_name.Shop:                      registerShop,
	table_name.Fulfillment:               registerFulfillment,
	table_name.ShopBrand:                 registerShopBrand,
	table_name.ShopProduct:               registerShopProduct,
	table_name.AccountUser:               registerAccountUser,
	table_name.Address:                   registerAddress,
	table_name.InventoryVariant:          registerInventoryVariant,
	table_name.InventoryVoucher:          registerInventoryVoucher,
	table_name.Invitation:                registerInvitation,
	table_name.MoneyTransactionShipping:  registerMoneyTransactionShipping,
	table_name.ProductShopCollection:     registerProductShopCollection,
	table_name.PurchaseOrder:             registerPurchaseOrder,
	table_name.PurchaseRefund:            registerPurchaseRefund,
	table_name.Receipt:                   registerReceipt,
	table_name.Refund:                    registerRefund,
	table_name.ShipNowFufillment:         registerShipNowFulfillment,
	table_name.ShopCarrier:               registerShopCarrier,
	table_name.ShopCategory:              registerShopCategory,
	table_name.ShopCollection:            registerShopCollection,
	table_name.ShopCustomerGroup:         registerShopCustomerGroup,
	table_name.ShopCustomerGroupCustomer: registerShopCustomerGroupCustomer,
	table_name.ShopLedger:                registerShopLedger,
	table_name.ShopProductCollection:     registerShopProductCollection,
	table_name.ShopStocktake:             registerShopStocktake,
	table_name.ShopSupplier:              registerShopSupplier,
	table_name.ShopTrader:                registerShopTrader,
	table_name.ShopTraderAddress:         registerShopTraderAddress,
	table_name.ShopVariant:               registerShopVariant,
	table_name.ShopVariantSupplier:       registerShopVariantSupplier,
}

func GetRegisterFuncFromTableName(name table_name.TableName) ETLRegisterFunc {
	return mapFieldRegisters[name]
}

func registerUser(ng *etl.ETLEngine, DB, dstDB *cmsql.Database, args ...interface{}) {
	userIDs := args[0]
	ng.Register(DB, (*identitymodel.Users)(nil), dstDB, (*etlusermodel.Users)(nil))
	ng.RegisterConversion(userconvert.RegisterConversions)
	ng.RegisterQuery(etl.ETLQuery{
		OrderBy: etl.OrderByRidASC,
		Where:   []interface{}{sq.In("id", userIDs.([]dot.ID))},
		Limit:   100,
	})
	ng.Bootstrap()
}

func registerAccount(ng *etl.ETLEngine, DB, dstDB *cmsql.Database, args ...interface{}) {
	accountIDs := args[0]
	ng.Register(DB, (*identitymodel.Accounts)(nil), dstDB, (*etlaccountmodel.Accounts)(nil))
	ng.RegisterConversion(accountconvert.RegisterConversions)
	ng.RegisterQuery(etl.ETLQuery{
		OrderBy: etl.OrderByRidASC,
		Where:   []interface{}{sq.In("id", accountIDs.([]dot.ID))},
		Limit:   100,
	})
	ng.Bootstrap()
}

func registerShopCustomer(ng *etl.ETLEngine, DB, dstDB *cmsql.Database, args ...interface{}) {
	accountIDs := args[0]
	ng.Register(DB, (*customermodel.ShopCustomers)(nil), dstDB, (*etlcustomermodel.ShopCustomers)(nil))
	ng.RegisterConversion(customerconvert.RegisterConversions)
	ng.RegisterQuery(etl.ETLQuery{
		OrderBy: etl.OrderByRidASC,
		Where:   []interface{}{sq.In("shop_id", accountIDs.([]dot.ID)), sq.NewExpr("deleted_at is null")},
		Limit:   100,
	})
	ng.Bootstrap()
}

func registerOrder(ng *etl.ETLEngine, DB, dstDB *cmsql.Database, args ...interface{}) {
	accountIDs := args[0]
	ng.Register(DB, (*ordermodel.Orders)(nil), dstDB, (*etlordermodel.Orders)(nil))
	ng.RegisterConversion(orderconvert.RegisterConversions)
	ng.RegisterQuery(etl.ETLQuery{
		OrderBy: etl.OrderByRidASC,
		Where:   []interface{}{sq.In("shop_id", accountIDs.([]dot.ID))},
		Limit:   100,
	})
	ng.Bootstrap()
}

func registerShop(ng *etl.ETLEngine, DB, dstDB *cmsql.Database, args ...interface{}) {
	accountIDs := args[0]
	ng.Register(DB, (*shopmodel.Shops)(nil), dstDB, (*etlshopmodel.Shops)(nil))
	ng.RegisterConversion(shopconvert.RegisterConversions)
	ng.RegisterQuery(etl.ETLQuery{
		OrderBy: etl.OrderByRidASC,
		Where:   []interface{}{sq.In("id", accountIDs.([]dot.ID)), sq.NewExpr("deleted_at is null")},
		Limit:   100,
	})
}

func registerFulfillment(ng *etl.ETLEngine, DB, dstDB *cmsql.Database, args ...interface{}) {
	accountIDs := args[0]
	ng.Register(DB, (*fulfillmentmodel.Fulfillments)(nil), dstDB, (*etlfulfillmentmodel.Fulfillments)(nil))
	ng.RegisterConversion(fulfillmentconvert.RegisterConversions)
	ng.RegisterQuery(etl.ETLQuery{
		OrderBy: etl.OrderByRidASC,
		Where:   []interface{}{sq.In("shop_id", accountIDs.([]dot.ID))},
		Limit:   100,
	})
}

func registerShopBrand(ng *etl.ETLEngine, DB, dstDB *cmsql.Database, args ...interface{}) {
	accountIDs := args[0]
	ng.Register(DB, (*catalogmodel.ShopBrands)(nil), dstDB, (*etlshopbrandmodel.ShopBrands)(nil))
	ng.RegisterConversion(shopbrandconvert.RegisterConversions)
	ng.RegisterQuery(etl.ETLQuery{
		OrderBy: etl.OrderByRidASC,
		Where:   []interface{}{sq.In("shop_id", accountIDs.([]dot.ID)), sq.NewExpr("deleted_at is null")},
		Limit:   100,
	})
}

func registerShopProduct(ng *etl.ETLEngine, DB, dstdB *cmsql.Database, args ...interface{}) {
	accountIDs := args[0]
	ng.Register(DB, (*catalogmodel.ShopProducts)(nil), dstdB, (*etlshopproductmodel.ShopProducts)(nil))
	ng.RegisterConversion(shopproductconvert.RegisterConversions)
	ng.RegisterQuery(etl.ETLQuery{
		OrderBy: etl.OrderByRidASC,
		Where:   []interface{}{sq.In("shop_id", accountIDs.([]dot.ID)), sq.NewExpr("deleted_at is null")},
		Limit:   100,
	})
}

func registerAccountUser(ng *etl.ETLEngine, DB, dstDB *cmsql.Database, args ...interface{}) {
	userIDs := args[0]
	ng.Register(DB, (*identitymodel.AccountUsers)(nil), dstDB, (*etlaccountusermodel.AccountUsers)(nil))
	ng.RegisterConversion(accountuserconvert.RegisterConversions)
	ng.RegisterQuery(etl.ETLQuery{
		OrderBy: etl.OrderByRidASC,
		Where:   []interface{}{sq.In("user_id", userIDs.([]dot.ID)), sq.NewExpr("deleted_at is null")},
		Limit:   100,
	})
}

func registerAddress(ng *etl.ETLEngine, DB, dstDB *cmsql.Database, args ...interface{}) {
	accountIDs := args[0]
	ng.Register(DB, (*addressmodel.Addresses)(nil), dstDB, (*etladdressmodel.Addresses)(nil))
	ng.RegisterConversion(addressconvert.RegisterConversions)
	ng.RegisterQuery(etl.ETLQuery{
		OrderBy: etl.OrderByRidASC,
		Where:   []interface{}{sq.In("account_id", accountIDs.([]dot.ID))},
		Limit:   100,
	})
}

func registerInventoryVariant(ng *etl.ETLEngine, DB, dstDB *cmsql.Database, args ...interface{}) {
	accountIDs := args[0]
	ng.Register(DB, (*inventorymodel.InventoryVariants)(nil), dstDB, (*etlinventorymodel.InventoryVariants)(nil))
	ng.RegisterConversion(inventoryconvert.RegisterConversions)
	ng.RegisterQuery(etl.ETLQuery{
		OrderBy: etl.OrderByRidASC,
		Where:   []interface{}{sq.In("shop_id", accountIDs.([]dot.ID))},
		Limit:   100,
	})
}

func registerInventoryVoucher(ng *etl.ETLEngine, DB, dstDB *cmsql.Database, args ...interface{}) {
	accountIDs := args[0]
	ng.Register(DB, (*inventorymodel.InventoryVouchers)(nil), dstDB, (*etlinventoryvouchermodel.InventoryVouchers)(nil))
	ng.RegisterConversion(inventoryvoucherconvert.RegisterConversions)
	ng.RegisterQuery(etl.ETLQuery{
		OrderBy: etl.OrderByRidASC,
		Where:   []interface{}{sq.In("shop_id", accountIDs.([]dot.ID))},
		Limit:   100,
	})
}

func registerInvitation(ng *etl.ETLEngine, DB, dstDB *cmsql.Database, args ...interface{}) {
	accountIDs := args[0]
	ng.Register(DB, (*invitationmodel.Invitations)(nil), dstDB, (*etlinvitationmodel.Invitations)(nil))
	ng.RegisterConversion(invitationconvert.RegisterConversions)
	ng.RegisterQuery(etl.ETLQuery{
		OrderBy: etl.OrderByRidASC,
		Where:   []interface{}{sq.In("account_id", accountIDs.([]dot.ID)), sq.NewExpr("deleted_at is null")},
		Limit:   100,
	})
}

func registerMoneyTransactionShipping(ng *etl.ETLEngine, DB, dstDB *cmsql.Database, args ...interface{}) {
	accountIDs := args[0]
	ng.Register(DB, (*moneytxmodel.MoneyTransactionShippings)(nil), dstDB, (*etlmoneytxmodel.MoneyTransactionShippings)(nil))
	ng.RegisterConversion(moneytxconvert.RegisterConversions)
	ng.RegisterQuery(etl.ETLQuery{
		OrderBy: etl.OrderByRidASC,
		Where:   []interface{}{sq.In("shop_id", accountIDs.([]dot.ID))},
		Limit:   100,
	})
}

func registerProductShopCollection(ng *etl.ETLEngine, DB, dstDB *cmsql.Database, args ...interface{}) {
	accountIDs := args[0]
	ng.Register(DB, (*catalogmodel.ProductShopCollections)(nil), dstDB, (*etlproductshopcollectionmodel.ProductShopCollections)(nil))
	ng.RegisterConversion(productshopcollectionconvert.RegisterConversions)
	ng.RegisterQuery(etl.ETLQuery{
		OrderBy: etl.OrderByRidASC,
		Where:   []interface{}{sq.In("shop_id", accountIDs.([]dot.ID))},
		Limit:   100,
	})
}

func registerPurchaseOrder(ng *etl.ETLEngine, DB, dstDB *cmsql.Database, args ...interface{}) {
	accountIDs := args[0]
	ng.Register(DB, (*purchaseorder.PurchaseOrders)(nil), dstDB, (*etlpurchaseordermodel.PurchaseOrders)(nil))
	ng.RegisterConversion(purchaseorderconvert.RegisterConversions)
	ng.RegisterQuery(etl.ETLQuery{
		OrderBy: etl.OrderByRidASC,
		Where:   []interface{}{sq.In("shop_id", accountIDs.([]dot.ID)), sq.NewExpr("deleted_at is null")},
		Limit:   100,
	})
}

func registerPurchaseRefund(ng *etl.ETLEngine, DB, dstDB *cmsql.Database, args ...interface{}) {
	accountIDs := args[0]
	ng.Register(DB, (*purchaserefundmodel.PurchaseRefunds)(nil), dstDB, (*etlpurchaserefundmodel.PurchaseRefunds)(nil))
	ng.RegisterConversion(purchaserefundconvert.RegisterConversions)
	ng.RegisterQuery(etl.ETLQuery{
		OrderBy: etl.OrderByRidASC,
		Where:   []interface{}{sq.In("shop_id", accountIDs.([]dot.ID))},
		Limit:   100,
	})
}

func registerReceipt(ng *etl.ETLEngine, DB, dstDB *cmsql.Database, args ...interface{}) {
	accountIDs := args[0]
	ng.Register(DB, (*receiptmodel.Receipts)(nil), dstDB, (*etlreceiptmodel.Receipts)(nil))
	ng.RegisterConversion(receiptconvert.RegisterConversions)
	ng.RegisterQuery(etl.ETLQuery{
		OrderBy: etl.OrderByRidASC,
		Where:   []interface{}{sq.In("shop_id", accountIDs.([]dot.ID)), sq.NewExpr("deleted_at is null")},
		Limit:   100,
	})
}

func registerRefund(ng *etl.ETLEngine, DB, dstDB *cmsql.Database, args ...interface{}) {
	accountIDs := args[0]
	ng.Register(DB, (*refundmodel.Refunds)(nil), dstDB, (*etlrefundmodel.Refunds)(nil))
	ng.RegisterConversion(refundconvert.RegisterConversions)
	ng.RegisterQuery(etl.ETLQuery{
		OrderBy: etl.OrderByRidASC,
		Where:   []interface{}{sq.In("shop_id", accountIDs.([]dot.ID))},
		Limit:   100,
	})
}

func registerShipNowFulfillment(ng *etl.ETLEngine, DB, dstDB *cmsql.Database, args ...interface{}) {
	accountIDs := args[0]
	ng.Register(DB, (*shipnowfulfillmentmodel.ShipnowFulfillments)(nil), dstDB, (*etlshipnowfulfillmentmodel.ShipnowFulfillments)(nil))
	ng.RegisterConversion(shipnowfulfillmentconvert.RegisterConversions)
	ng.RegisterQuery(etl.ETLQuery{
		OrderBy: etl.OrderByRidASC,
		Where:   []interface{}{sq.In("shop_id", accountIDs.([]dot.ID))},
		Limit:   100,
	})
}

func registerShopCarrier(ng *etl.ETLEngine, DB, dstDB *cmsql.Database, args ...interface{}) {
	accountIDs := args[0]
	ng.Register(DB, (*carriermodel.ShopCarriers)(nil), dstDB, (*etlcarriermodel.ShopCarriers)(nil))
	ng.RegisterConversion(carrierconvert.RegisterConversions)
	ng.RegisterQuery(etl.ETLQuery{
		OrderBy: etl.OrderByRidASC,
		Where:   []interface{}{sq.In("shop_id", accountIDs.([]dot.ID)), sq.NewExpr("deleted_at is null")},
		Limit:   100,
	})
}

func registerShopCategory(ng *etl.ETLEngine, DB, dstDB *cmsql.Database, args ...interface{}) {
	accountIDs := args[0]
	ng.Register(DB, (*catalogmodel.ShopCategories)(nil), dstDB, (*etlshopcategorymodel.ShopCategories)(nil))
	ng.RegisterConversion(shopcategoryconvert.RegisterConversions)
	ng.RegisterQuery(etl.ETLQuery{
		OrderBy: etl.OrderByRidASC,
		Where:   []interface{}{sq.In("shop_id", accountIDs.([]dot.ID)), sq.NewExpr("deleted_at is null")},
		Limit:   100,
	})
}

func registerShopCollection(ng *etl.ETLEngine, DB, dstDB *cmsql.Database, args ...interface{}) {
	accountIDs := args[0]
	ng.Register(DB, (*catalogmodel.ShopCollections)(nil), dstDB, (*etlshopcollectionmodel.ShopCollections)(nil))
	ng.RegisterConversion(shopcollectionconvert.RegisterConversions)
	ng.RegisterQuery(etl.ETLQuery{
		OrderBy: etl.OrderByRidASC,
		Where:   []interface{}{sq.In("shop_id", accountIDs.([]dot.ID)), sq.NewExpr("deleted_at is null")},
		Limit:   100,
	})
}

func registerShopCustomerGroup(ng *etl.ETLEngine, DB, dstDB *cmsql.Database, args ...interface{}) {
	accountIDs := args[0]
	ng.Register(DB, (*customermodel.ShopCustomerGroups)(nil), dstDB, (*etlshopcustomergroupmodel.ShopCustomerGroups)(nil))
	ng.RegisterConversion(shopcustomergroupconvert.RegisterConversions)
	ng.RegisterQuery(etl.ETLQuery{
		OrderBy: etl.OrderByRidASC,
		Where:   []interface{}{sq.In("shop_id", accountIDs.([]dot.ID)), sq.NewExpr("deleted_at is null")},
		Limit:   100,
	})
}

func registerShopCustomerGroupCustomer(ng *etl.ETLEngine, DB, dstDB *cmsql.Database, args ...interface{}) {
	accountIDs := args[0]
	ng.Register(DB, (*customermodel.ShopCustomerGroupCustomers)(nil), dstDB, (*etlshopcustomergroupcustomermodel.ShopCustomerGroupCustomers)(nil))
	ng.RegisterConversion(shopcustomergroupcustomerconvert.RegisterConversions)
	ng.RegisterQuery(etl.ETLQuery{
		OrderBy: etl.OrderByRidASC,
		Where: []interface{}{
			sq.NewExpr(fmt.Sprintf("customer_id in (select customer_id from shop_customer where shop_id in (%s) and deleted_at is null)", dot.JoinIDs(accountIDs.([]dot.ID))))},
		Limit: 100,
	})
}

func registerShopLedger(ng *etl.ETLEngine, DB, dstDB *cmsql.Database, args ...interface{}) {
	accountIDs := args[0]
	ng.Register(DB, (*shopledgermodel.ShopLedgers)(nil), dstDB, (*etlshopledgermodel.ShopLedgers)(nil))
	ng.RegisterConversion(shopledgerconvert.RegisterConversions)
	ng.RegisterQuery(etl.ETLQuery{
		OrderBy: etl.OrderByRidASC,
		Where:   []interface{}{sq.In("shop_id", accountIDs.([]dot.ID)), sq.NewExpr("deleted_at is null")},
		Limit:   100,
	})
}

func registerShopProductCollection(ng *etl.ETLEngine, DB, dstDB *cmsql.Database, args ...interface{}) {
	accountIDs := args[0]
	ng.Register(DB, (*catalogmodel.ShopProductCollections)(nil), dstDB, (*etlshopproductcollectionmodel.ShopProductCollections)(nil))
	ng.RegisterConversion(shopproductcollectionconvert.RegisterConversions)
	ng.RegisterQuery(etl.ETLQuery{
		OrderBy: etl.OrderByRidASC,
		Where:   []interface{}{sq.NewExpr(fmt.Sprintf("product_id in (select product_id from shop_product where shop_id in (%s) and deleted_at is null)", dot.JoinIDs(accountIDs.([]dot.ID))))},
		Limit:   100,
	})
}

func registerShopStocktake(ng *etl.ETLEngine, DB, dstDB *cmsql.Database, args ...interface{}) {
	accountIDs := args[0]
	ng.Register(DB, (*shopstocktakemodel.ShopStocktakes)(nil), dstDB, (*etlshopstocktakemodel.ShopStocktakes)(nil))
	ng.RegisterConversion(shopstocktakeconvert.RegisterConversions)
	ng.RegisterQuery(etl.ETLQuery{
		OrderBy: etl.OrderByRidASC,
		Where:   []interface{}{sq.In("shop_id", accountIDs.([]dot.ID))},
		Limit:   100,
	})
}

func registerShopSupplier(ng *etl.ETLEngine, DB, dstDB *cmsql.Database, args ...interface{}) {
	accountIDs := args[0]
	ng.Register(DB, (*suppliermodel.ShopSuppliers)(nil), dstDB, (*etlshopsuppliermodel.ShopSuppliers)(nil))
	ng.RegisterConversion(shopsupplierconvert.RegisterConversions)
	ng.RegisterQuery(etl.ETLQuery{
		OrderBy: etl.OrderByRidASC,
		Where:   []interface{}{sq.In("shop_id", accountIDs.([]dot.ID)), sq.NewExpr("deleted_at is null")},
		Limit:   100,
	})
}

func registerShopTrader(ng *etl.ETLEngine, DB, dstDB *cmsql.Database, args ...interface{}) {
	accountIDs := args[0]
	ng.Register(DB, (*tradermodel.ShopTraders)(nil), dstDB, (*etlshoptradermodel.ShopTraders)(nil))
	ng.RegisterConversion(shoptraderconvert.RegisterConversions)
	ng.RegisterQuery(etl.ETLQuery{
		OrderBy: etl.OrderByRidASC,
		Where:   []interface{}{sq.In("shop_id", accountIDs.([]dot.ID)), sq.NewExpr("deleted_at is null")},
		Limit:   100,
	})
}

func registerShopTraderAddress(ng *etl.ETLEngine, DB, dstDB *cmsql.Database, args ...interface{}) {
	accountIDs := args[0]
	ng.Register(DB, (*customermodel.ShopTraderAddresses)(nil), dstDB, (*etlshoptraderaddressmodel.ShopTraderAddresses)(nil))
	ng.RegisterConversion(shoptraderaddressconvert.RegisterConversions)
	ng.RegisterQuery(etl.ETLQuery{
		OrderBy: etl.OrderByRidASC,
		Where:   []interface{}{sq.In("shop_id", accountIDs.([]dot.ID))},
		Limit:   100,
	})
}

func registerShopVariant(ng *etl.ETLEngine, DB, dstDB *cmsql.Database, args ...interface{}) {
	accountIDs := args[0]
	ng.Register(DB, (*catalogmodel.ShopVariants)(nil), dstDB, (*etlshopvariantmodel.ShopVariants)(nil))
	ng.RegisterConversion(shopvariantconvert.RegisterConversions)
	ng.RegisterQuery(etl.ETLQuery{
		OrderBy: etl.OrderByRidASC,
		Where:   []interface{}{sq.In("shop_id", accountIDs.([]dot.ID)), sq.NewExpr("deleted_at is null")},
		Limit:   100,
	})
}

func registerShopVariantSupplier(ng *etl.ETLEngine, DB, dstDB *cmsql.Database, args ...interface{}) {
	accountIDs := args[0]
	ng.Register(DB, (*catalogmodel.ShopVariantSuppliers)(nil), dstDB, (*etlshopvariantsuppliermodel.ShopVariantSuppliers)(nil))
	ng.RegisterConversion(shopvariantsupplierconvert.RegisterConversions)
	ng.RegisterQuery(etl.ETLQuery{
		OrderBy: etl.OrderByRidASC,
		Where:   []interface{}{sq.In("shop_id", accountIDs.([]dot.ID))},
		Limit:   100,
	})
}
