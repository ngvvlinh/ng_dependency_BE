package main

import (
	"o.o/backend/pkg/etop/api"
	"o.o/backend/pkg/etop/api/shop"
	"o.o/backend/pkg/etop/apix/partner"
	xshop "o.o/backend/pkg/etop/apix/shop"
	"o.o/backend/pkg/etop/logic/shipping_provider"
	affapi "o.o/backend/pkg/services/affiliate/api"
	"o.o/capi/httprpc"
)

func NewServers(
	apiServers api.Servers,
	shopServers shop.Servers,
	affServer affapi.Servers,
	partnerServers partner.Servers,
	xshopServers xshop.Servers,
) (servers []httprpc.Server) {
	servers = append(servers, apiServers...)
	servers = append(servers, shopServers...)
	servers = append(servers, affServer...)

	// TODO: remove this later
	serversExt = append(serversExt, partnerServers...)
	serversExt = append(serversExt, xshopServers...)
	return servers
}

func SupportedCarrierDrivers() []shipping_provider.CarrierDriver {
	return []shipping_provider.CarrierDriver{ghnCarrier, ghtkCarrier, vtpostCarrier}
}
