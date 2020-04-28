package register

import (
	"fmt"

	"o.o/backend/cmd/etop-etl/register/table_name"
	addressmodel "o.o/backend/com/main/address/model"
	catalogmodel "o.o/backend/com/main/catalog/model"
	identitymodel "o.o/backend/com/main/identity/model"
	shopmodel "o.o/backend/com/main/identity/model"
	inventorymodel "o.o/backend/com/main/inventory/model"
	invitationmodel "o.o/backend/com/main/invitation/model"
	shopledgermodel "o.o/backend/com/main/ledgering/model"
	moneytxmodel "o.o/backend/com/main/moneytx/model"
	ordermodel "o.o/backend/com/main/ordering/model"
	purchaseorder "o.o/backend/com/main/purchaseorder/model"
	purchaserefundmodel "o.o/backend/com/main/purchaserefund/model"
	receiptmodel "o.o/backend/com/main/receipting/model"
	refundmodel "o.o/backend/com/main/refund/model"
	shipnowfulfillmentmodel "o.o/backend/com/main/shipnow/model"
	fulfillmentmodel "o.o/backend/com/main/shipping/model"
	shopstocktakemodel "o.o/backend/com/main/stocktaking/model"
	carriermodel "o.o/backend/com/shopping/carrying/model"
	customermodel "o.o/backend/com/shopping/customering/model"
	suppliermodel "o.o/backend/com/shopping/suppliering/model"
	tradermodel "o.o/backend/com/shopping/tradering/model"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/zexp/etl"
	accountconvert "o.o/backend/zexp/etl/main/account/convert"
	etlaccountmodel "o.o/backend/zexp/etl/main/account/model"
	accountuserconvert "o.o/backend/zexp/etl/main/accountuser/convert"
	etlaccountusermodel "o.o/backend/zexp/etl/main/accountuser/model"
	addressconvert "o.o/backend/zexp/etl/main/address/convert"
	etladdressmodel "o.o/backend/zexp/etl/main/address/model"
	customerconvert "o.o/backend/zexp/etl/main/customer/convert"
	etlcustomermodel "o.o/backend/zexp/etl/main/customer/model"
	fulfillmentconvert "o.o/backend/zexp/etl/main/fulfillment/convert"
	etlfulfillmentmodel "o.o/backend/zexp/etl/main/fulfillment/model"
	inventoryconvert "o.o/backend/zexp/etl/main/inventoryvariant/convert"
	etlinventorymodel "o.o/backend/zexp/etl/main/inventoryvariant/model"
	inventoryvoucherconvert "o.o/backend/zexp/etl/main/inventoryvoucher/convert"
	etlinventoryvouchermodel "o.o/backend/zexp/etl/main/inventoryvoucher/model"
	invitationconvert "o.o/backend/zexp/etl/main/invitation/convert"
	etlinvitationmodel "o.o/backend/zexp/etl/main/invitation/model"
	moneytxconvert "o.o/backend/zexp/etl/main/moneytransactionshipping/convert"
	etlmoneytxmodel "o.o/backend/zexp/etl/main/moneytransactionshipping/model"
	orderconvert "o.o/backend/zexp/etl/main/order/convert"
	etlordermodel "o.o/backend/zexp/etl/main/order/model"
	productshopcollectionconvert "o.o/backend/zexp/etl/main/productshopcollection/convert"
	etlproductshopcollectionmodel "o.o/backend/zexp/etl/main/productshopcollection/model"
	purchaseorderconvert "o.o/backend/zexp/etl/main/purchaseorder/convert"
	etlpurchaseordermodel "o.o/backend/zexp/etl/main/purchaseorder/model"
	purchaserefundconvert "o.o/backend/zexp/etl/main/purchaserefund/convert"
	etlpurchaserefundmodel "o.o/backend/zexp/etl/main/purchaserefund/model"
	receiptconvert "o.o/backend/zexp/etl/main/receipt/convert"
	etlreceiptmodel "o.o/backend/zexp/etl/main/receipt/model"
	refundconvert "o.o/backend/zexp/etl/main/refund/convert"
	etlrefundmodel "o.o/backend/zexp/etl/main/refund/model"
	shipnowfulfillmentconvert "o.o/backend/zexp/etl/main/shipnowfulfillment/convert"
	etlshipnowfulfillmentmodel "o.o/backend/zexp/etl/main/shipnowfulfillment/model"
	shopconvert "o.o/backend/zexp/etl/main/shop/convert"
	etlshopmodel "o.o/backend/zexp/etl/main/shop/model"
	shopbrandconvert "o.o/backend/zexp/etl/main/shopbrand/convert"
	etlshopbrandmodel "o.o/backend/zexp/etl/main/shopbrand/model"
	carrierconvert "o.o/backend/zexp/etl/main/shopcarrier/convert"
	etlcarriermodel "o.o/backend/zexp/etl/main/shopcarrier/model"
	shopcategoryconvert "o.o/backend/zexp/etl/main/shopcategory/convert"
	etlshopcategorymodel "o.o/backend/zexp/etl/main/shopcategory/model"
	shopcollectionconvert "o.o/backend/zexp/etl/main/shopcollection/convert"
	etlshopcollectionmodel "o.o/backend/zexp/etl/main/shopcollection/model"
	shopcustomergroupconvert "o.o/backend/zexp/etl/main/shopcustomergroup/convert"
	etlshopcustomergroupmodel "o.o/backend/zexp/etl/main/shopcustomergroup/model"
	shopcustomergroupcustomerconvert "o.o/backend/zexp/etl/main/shopcustomergroupcustomer/convert"
	etlshopcustomergroupcustomermodel "o.o/backend/zexp/etl/main/shopcustomergroupcustomer/model"
	shopledgerconvert "o.o/backend/zexp/etl/main/shopledger/convert"
	etlshopledgermodel "o.o/backend/zexp/etl/main/shopledger/model"
	shopproductconvert "o.o/backend/zexp/etl/main/shopproduct/convert"
	etlshopproductmodel "o.o/backend/zexp/etl/main/shopproduct/model"
	shopproductcollectionconvert "o.o/backend/zexp/etl/main/shopproductcollection/convert"
	etlshopproductcollectionmodel "o.o/backend/zexp/etl/main/shopproductcollection/model"
	shopstocktakeconvert "o.o/backend/zexp/etl/main/shopstocktake/convert"
	etlshopstocktakemodel "o.o/backend/zexp/etl/main/shopstocktake/model"
	shopsupplierconvert "o.o/backend/zexp/etl/main/shopsupplier/convert"
	etlshopsuppliermodel "o.o/backend/zexp/etl/main/shopsupplier/model"
	shoptraderconvert "o.o/backend/zexp/etl/main/shoptrader/convert"
	etlshoptradermodel "o.o/backend/zexp/etl/main/shoptrader/model"
	shoptraderaddressconvert "o.o/backend/zexp/etl/main/shoptraderaddress/convert"
	etlshoptraderaddressmodel "o.o/backend/zexp/etl/main/shoptraderaddress/model"
	shopvariantconvert "o.o/backend/zexp/etl/main/shopvariant/convert"
	etlshopvariantmodel "o.o/backend/zexp/etl/main/shopvariant/model"
	shopvariantsupplierconvert "o.o/backend/zexp/etl/main/shopvariantsupplier/convert"
	etlshopvariantsuppliermodel "o.o/backend/zexp/etl/main/shopvariantsupplier/model"
	userconvert "o.o/backend/zexp/etl/main/user/convert"
	etlusermodel "o.o/backend/zexp/etl/main/user/model"
	"o.o/capi/dot"
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
