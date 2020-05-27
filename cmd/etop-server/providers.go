package main

import (
	"o.o/backend/pkg/etop/api"
	"o.o/backend/pkg/etop/api/shop"
	"o.o/backend/pkg/etop/apix/partner"
	"o.o/backend/pkg/etop/apix/partnercarrier"
	xshop "o.o/backend/pkg/etop/apix/shop"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/etop/logic/shipping_provider"
	affapi "o.o/backend/pkg/services/affiliate/api"
	"o.o/backend/tools/pkg/acl"
	"o.o/capi/httprpc"
)

func NewServers(
	apiServers api.Servers,
	shopServers shop.Servers,
	affServer affapi.Servers,
	partnerServers partner.Servers,
	xshopServers xshop.Servers,
	carrierServers partnercarrier.Servers,
) (servers []httprpc.Server) {
	servers = append(servers, apiServers...)
	servers = append(servers, shopServers...)
	servers = append(servers, affServer...)

	// TODO: remove this later
	extHooks := session.NewHook(acl.GetExtACL())
	serversExt = append(serversExt, partnerServers...)
	serversExt = append(serversExt, xshopServers...)
	serversExt = httprpc.WithHooks(serversExt, extHooks)
	return servers
}

func SupportedCarrierDrivers() []shipping_provider.CarrierDriver {
	return []shipping_provider.CarrierDriver{ghnCarrier, ghtkCarrier, vtpostCarrier}
}
