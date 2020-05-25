package shop

import (
	"o.o/api/main/catalog"
	"o.o/api/main/inventory"
	"o.o/api/webserver"
)

type WebServerService struct {
	CatalogQuery   catalog.QueryBus
	WebserverAggr  webserver.CommandBus
	WebserverQuery webserver.QueryBus
	InventoryQuery inventory.QueryBus
}

func (s *WebServerService) Clone() *WebServerService { res := *s; return &res }
