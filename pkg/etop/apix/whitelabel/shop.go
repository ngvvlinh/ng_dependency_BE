package whitelabel

import (
	"etop.vn/api/main/catalog"
	catalogsqlstore "etop.vn/backend/com/main/catalog/sqlstore"
	customersqlstore "etop.vn/backend/com/shopping/customering/sqlstore"
	"etop.vn/backend/pkg/common/sql/cmsql"
)

var (
	db                                *cmsql.Database
	catalogAggregate                  *catalog.CommandBus
	shopProductStoreFactory           catalogsqlstore.ShopProductStoreFactory
	shopCollectionStoreFactory        catalogsqlstore.ShopCollectionStoreFactory
	shopProductCollectionStoreFactory catalogsqlstore.ShopProductCollectionStoreFactory
	shopVariantStoreFactory           catalogsqlstore.ShopVariantStoreFactory
	brandStoreFactory                 catalogsqlstore.ShopBrandStoreFactory
	categoryStoreFactory              catalogsqlstore.ShopCategoryStoreFactory
	customerStoreFactory              customersqlstore.CustomerStoreFactory
)

const (
	MaximumItems = 100
)

func Init(
	database *cmsql.Database,
	catalogA *catalog.CommandBus,
) {
	db = database
	catalogAggregate = catalogA
	shopProductStoreFactory = catalogsqlstore.NewShopProductStore(db)
	shopCollectionStoreFactory = catalogsqlstore.NewShopCollectionStore(db)
	shopProductCollectionStoreFactory = catalogsqlstore.NewShopProductCollectionStore(db)
	shopVariantStoreFactory = catalogsqlstore.NewShopVariantStore(db)
	brandStoreFactory = catalogsqlstore.NewShopBrandStore(db)
	categoryStoreFactory = catalogsqlstore.NewShopCategoryStore(db)
	customerStoreFactory = customersqlstore.NewCustomerStore(db)
}

type ImportService struct{}

var importService = &ImportService{}
