// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package api

import (
	"github.com/google/wire"

	"o.o/api/main/catalog"
	"o.o/api/main/identity"
	"o.o/api/main/inventory"
	"o.o/api/services/affiliate"
	cc "o.o/backend/pkg/common/config"
)

// Injectors from wire.go:

func BuildServers(secret cc.SecretToken, affAggr affiliate.CommandBus, affQuery affiliate.QueryBus, catQuery catalog.QueryBus, idenQuery identity.QueryBus, inventoryQuery inventory.QueryBus) Servers {
	userService := &UserService{
		AffiliateAggr: affAggr,
	}
	tradingService := &TradingService{
		AffiliateAggr:  affAggr,
		AffiliateQuery: affQuery,
		CatalogQuery:   catQuery,
		InventoryQuery: inventoryQuery,
	}
	shopService := &ShopService{
		CatalogQuery:   catQuery,
		InventoryQuery: inventoryQuery,
		AffiliateQuery: affQuery,
	}
	affiliateService := &AffiliateService{
		AffiliateAggr:  affAggr,
		CatalogQuery:   catQuery,
		AffiliateQuery: affQuery,
		IdentityQuery:  idenQuery,
	}
	servers := NewServers(secret, userService, tradingService, shopService, affiliateService)
	return servers
}

// wire.go:

var WireSet = wire.NewSet(wire.Struct(new(UserService), "*"), wire.Struct(new(TradingService), "*"), wire.Struct(new(ShopService), "*"), wire.Struct(new(AffiliateService), "*"), NewServers)
