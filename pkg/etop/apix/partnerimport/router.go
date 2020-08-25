package partnerimport

import (
	"o.o/api/top/external/whitelabel"
	com "o.o/backend/com/main"
	catalogsqlstore "o.o/backend/com/main/catalog/sqlstore"
	customersqlstore "o.o/backend/com/shopping/customering/sqlstore"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/capi/httprpc"
)

type Servers []httprpc.Server

func NewServers(importService ImportService) Servers {
	servers := httprpc.MustNewServers(
		importService.Clone,
	)
	return servers
}

const (
	MaximumItems = 100
)

type ImportService struct {
	session.Session

	shopProductStoreFactory           catalogsqlstore.ShopProductStoreFactory
	shopCollectionStoreFactory        catalogsqlstore.ShopCollectionStoreFactory
	shopProductCollectionStoreFactory catalogsqlstore.ShopProductCollectionStoreFactory
	shopVariantStoreFactory           catalogsqlstore.ShopVariantStoreFactory
	brandStoreFactory                 catalogsqlstore.ShopBrandStoreFactory
	categoryStoreFactory              catalogsqlstore.ShopCategoryStoreFactory
	customerStoreFactory              customersqlstore.CustomerStoreFactory
}

func New(
	ss session.Session,
	db com.MainDB,
) ImportService {
	s := ImportService{
		Session:                           ss,
		shopProductStoreFactory:           catalogsqlstore.NewShopProductStore(db),
		shopCollectionStoreFactory:        catalogsqlstore.NewShopCollectionStore(db),
		shopProductCollectionStoreFactory: catalogsqlstore.NewShopProductCollectionStore(db),
		shopVariantStoreFactory:           catalogsqlstore.NewShopVariantStore(db),
		brandStoreFactory:                 catalogsqlstore.NewShopBrandStore(db),
		categoryStoreFactory:              catalogsqlstore.NewShopCategoryStore(db),
		customerStoreFactory:              customersqlstore.NewCustomerStore(db),
	}
	return s
}

func (s *ImportService) Clone() whitelabel.ImportService { res := *s; return &res }
