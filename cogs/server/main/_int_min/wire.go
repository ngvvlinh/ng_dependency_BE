package server_int_min

import (
	"net/http"

	"github.com/google/wire"

	_main "o.o/backend/cogs/server/main"
	"o.o/backend/pkg/etop/api"
	"o.o/backend/pkg/etop/api/admin"
	"o.o/backend/pkg/etop/api/sadmin"
	"o.o/backend/pkg/etop/api/shop"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/etop/middlewares"
	"o.o/backend/tools/pkg/acl"
	"o.o/capi/httprpc"
)

var WireSet = wire.NewSet(
	BuildHandlers,
)

func BuildHandlers(
	rootServers api.Servers,
	shopServers shop.Servers,
	adminServers admin.Servers,
	sadminServers sadmin.Servers,
) _main.IntHandlers {
	logging := middlewares.NewLogging()
	ssHooks := session.NewHook(acl.GetACL())

	var servers []httprpc.Server
	servers = append(servers, rootServers...)
	servers = append(servers, shopServers...)
	servers = append(servers, adminServers...)
	servers = append(servers, httprpc.WithHooks(sadminServers, ssHooks, logging)...)
	servers = httprpc.WithHooks(servers)

	mux := http.NewServeMux()
	for _, s := range servers {
		mux.Handle(s.PathPrefix(), s)
	}
	return servers
}
