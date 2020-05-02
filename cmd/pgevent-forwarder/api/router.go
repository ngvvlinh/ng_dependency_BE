package api

import (
	service "o.o/api/top/services/pgevent"
	"o.o/capi/httprpc"
)

// +gen:wrapper=o.o/api/top/services/pgevent
// +gen:wrapper:package=pgevent

func NewPgeventServer(m httprpc.Muxer, secret string) {
	servers := []httprpc.Server{
		service.NewMiscServiceServer(WrapMiscService(miscService.Clone, secret)),
		service.NewEventServiceServer(WrapEventService(eventService.Clone, secret)),
	}
	for _, s := range servers {
		m.Handle(s.PathPrefix(), s)
	}
}
