package api

import (
	"o.o/capi/httprpc"
	"o.o/common/l"
)

var ll = l.New()

type Servers []httprpc.Server

func NewPgeventServer(
	miscService *MiscService,
	eventService *EventService,
) Servers {
	servers := httprpc.MustNewServers(
		miscService.Clone,
		eventService.Clone,
	)
	return servers
}
