package whitelabel

import (
	"o.o/api/main/catalog"
	catalogsqlstore "o.o/backend/com/main/catalog/sqlstore"
	customersqlstore "o.o/backend/com/shopping/customering/sqlstore"
	"o.o/backend/pkg/common/sql/cmsql"
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

func (s *ImportService) Clone() *ImportService { res := *s; return &res }

var importService = &ImportService{}
