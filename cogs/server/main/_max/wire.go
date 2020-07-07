package server_max

import (
	"github.com/google/wire"

	_main "o.o/backend/cogs/server/main"
	"o.o/backend/pkg/etop/api"
	"o.o/backend/pkg/etop/api/admin"
	affapi "o.o/backend/pkg/etop/api/affiliate"
	"o.o/backend/pkg/etop/api/integration"
	"o.o/backend/pkg/etop/api/sadmin"
	"o.o/backend/pkg/etop/api/shop"
	"o.o/backend/pkg/etop/apix/partner"
	"o.o/backend/pkg/etop/apix/partnercarrier"
	"o.o/backend/pkg/etop/apix/partnerimport"
	xshop "o.o/backend/pkg/etop/apix/shop"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/etop/middlewares"
	serviceaffapi "o.o/backend/pkg/services/affiliate/api"
	"o.o/backend/tools/pkg/acl"
	"o.o/capi/httprpc"
)

var WireSet = wire.NewSet(
	BuildIntHandlers,
	BuildExtHandlers,
)

func BuildIntHandlers(
	rootServers api.Servers,
	shopServers shop.Servers,
	adminServers admin.Servers,
	sadminServers sadmin.Servers,
	integrationServers integration.Servers,
	affServer affapi.Servers,
	saffServer serviceaffapi.Servers,
) (hs _main.IntHandlers) {
	logging := middlewares.NewLogging()
	ssHooks := session.NewHook(acl.GetACL())

	hs = append(hs, rootServers...)
	hs = append(hs, shopServers...)
	hs = append(hs, adminServers...)
	hs = append(hs, sadminServers...)
	hs = httprpc.WithHooks(hs, ssHooks, logging)

	hs = append(hs, integrationServers...)
	hs = append(hs, affServer...)
	hs = append(hs, saffServer...)
	return hs
}

func BuildExtHandlers(
	partnerServers partner.Servers,
	xshopServers xshop.Servers,
	carrierServers partnercarrier.Servers,
	partnerImportServers partnerimport.Servers,
) (hs _main.ExtHandlers) {
	logging := middlewares.NewLogging()
	ssExtHooks := session.NewHook(acl.GetExtACL())

	hs = append(hs, partnerServers...)
	hs = append(hs, xshopServers...)
	hs = append(hs, httprpc.WithHooks(carrierServers, ssExtHooks, logging)...)
	hs = append(hs, partnerImportServers...)
	hs = httprpc.WithHooks(hs)
	return
}
