package partnerimport

import (
	"o.o/api/main/catalog"
	"o.o/api/top/external/whitelabel"
	catalogsqlstore "o.o/backend/com/main/catalog/sqlstore"
	customersqlstore "o.o/backend/com/shopping/customering/sqlstore"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/httprpc"
)

// +gen:wrapper=o.o/api/top/external/whitelabel
// +gen:wrapper:package=partner
// +gen:wrapper:prefix=ext

type Servers []httprpc.Server

func NewServers(importService ImportService) Servers {
	servers := []httprpc.Server{
		whitelabel.NewImportServiceServer(WrapImportService(importService.Clone)),
	}
	return servers
}

const (
	MaximumItems = 100
)

type ImportService struct {
	shopProductStoreFactory           catalogsqlstore.ShopProductStoreFactory
	shopCollectionStoreFactory        catalogsqlstore.ShopCollectionStoreFactory
	shopProductCollectionStoreFactory catalogsqlstore.ShopProductCollectionStoreFactory
	shopVariantStoreFactory           catalogsqlstore.ShopVariantStoreFactory
	brandStoreFactory                 catalogsqlstore.ShopBrandStoreFactory
	categoryStoreFactory              catalogsqlstore.ShopCategoryStoreFactory
	customerStoreFactory              customersqlstore.CustomerStoreFactory
}

func New(
	db *cmsql.Database,
	catalogAggr catalog.CommandBus,
) ImportService {
	s := ImportService{}
	s.shopProductStoreFactory = catalogsqlstore.NewShopProductStore(db)
	s.shopCollectionStoreFactory = catalogsqlstore.NewShopCollectionStore(db)
	s.shopProductCollectionStoreFactory = catalogsqlstore.NewShopProductCollectionStore(db)
	s.shopVariantStoreFactory = catalogsqlstore.NewShopVariantStore(db)
	s.brandStoreFactory = catalogsqlstore.NewShopBrandStore(db)
	s.categoryStoreFactory = catalogsqlstore.NewShopCategoryStore(db)
	s.customerStoreFactory = customersqlstore.NewCustomerStore(db)
	return s
}

func (s *ImportService) Clone() *ImportService { res := *s; return &res }
