package ws

import (
	"o.o/api/main/catalog"
	"o.o/api/main/inventory"
	api "o.o/api/top/int/shop"
	"o.o/api/webserver"
	"o.o/backend/pkg/etop/authorize/session"
)

type WebServerService struct {
	session.Session

	CatalogQuery   catalog.QueryBus
	WebserverAggr  webserver.CommandBus
	WebserverQuery webserver.QueryBus
	InventoryQuery inventory.QueryBus
}

func (s *WebServerService) Clone() api.WebServerService { res := *s; return &res }
