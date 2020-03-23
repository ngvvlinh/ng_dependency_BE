package register

import (
	"etop.vn/backend/cmd/etop-etl/register/table_name"
	catalogmodel "etop.vn/backend/com/main/catalog/model"
	identitymodel "etop.vn/backend/com/main/identity/model"
	shopmodel "etop.vn/backend/com/main/identity/model"
	ordermodel "etop.vn/backend/com/main/ordering/model"
	fulfillmentmodel "etop.vn/backend/com/main/shipping/model"
	customermodel "etop.vn/backend/com/shopping/customering/model"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/pkg/common/sql/sq"
	"etop.vn/backend/zexp/etl"
	accountconvert "etop.vn/backend/zexp/etl/main/account/convert"
	etlaccountmodel "etop.vn/backend/zexp/etl/main/account/model"
	customerconvert "etop.vn/backend/zexp/etl/main/customer/convert"
	etlcustomermodel "etop.vn/backend/zexp/etl/main/customer/model"
	fulfillmentconvert "etop.vn/backend/zexp/etl/main/fulfillment/convert"
	etlfulfillmentmodel "etop.vn/backend/zexp/etl/main/fulfillment/model"
	orderconvert "etop.vn/backend/zexp/etl/main/order/convert"
	etlordermodel "etop.vn/backend/zexp/etl/main/order/model"
	shopconvert "etop.vn/backend/zexp/etl/main/shop/convert"
	etlshopmodel "etop.vn/backend/zexp/etl/main/shop/model"
	shopbrandconvert "etop.vn/backend/zexp/etl/main/shopbrand/convert"
	etlshopbrandmodel "etop.vn/backend/zexp/etl/main/shopbrand/model"
	shopproductconvert "etop.vn/backend/zexp/etl/main/shopproduct/convert"
	etlshopproductmodel "etop.vn/backend/zexp/etl/main/shopproduct/model"
	userconvert "etop.vn/backend/zexp/etl/main/user/convert"
	etlusermodel "etop.vn/backend/zexp/etl/main/user/model"
	"etop.vn/capi/dot"
)

type ETLRegisterFunc func(ng *etl.ETLEngine, DB, dstDB *cmsql.Database, args ...interface{})

var mapFieldRegisters = map[table_name.TableName]ETLRegisterFunc{
	table_name.User:         registerUser,
	table_name.Account:      registerAccount,
	table_name.ShopCustomer: registerShopCustomer,
	table_name.Order:        registerOrder,
	table_name.Shop:         registerShop,
	table_name.Fulfillment:  registerFulfillment,
	table_name.ShopBrand:    registerShopBrand,
	table_name.ShopProduct:  registerShopProduct,
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
		Where:   []interface{}{sq.In("shop_id", accountIDs.([]dot.ID))},
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
		Where:   []interface{}{sq.In("id", accountIDs.([]dot.ID))},
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
		Where:   []interface{}{sq.In("shop_id", accountIDs.([]dot.ID))},
		Limit:   100,
	})
}

func registerShopProduct(ng *etl.ETLEngine, DB, dstdB *cmsql.Database, args ...interface{}) {
	accountIDs := args[0]
	ng.Register(DB, (*catalogmodel.ShopProducts)(nil), dstdB, (*etlshopproductmodel.ShopProducts)(nil))
	ng.RegisterConversion(shopproductconvert.RegisterConversions)
	ng.RegisterQuery(etl.ETLQuery{
		OrderBy: etl.OrderByRidASC,
		Where:   []interface{}{sq.In("shop_id", accountIDs.([]dot.ID))},
		Limit:   100,
	})
}
