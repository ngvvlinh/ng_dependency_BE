package api

import (
	service "etop.vn/api/root/services/pgevent"
	"etop.vn/capi/httprpc"
)

// +gen:wrapper=etop.vn/api/root/services/pgevent
// +gen:wrapper:package=pgevent

func NewPgeventServer(m httprpc.Muxer, secret string) {
	servers := []httprpc.Server{
		service.NewMiscServiceServer(WrapMiscService(miscService, secret)),
		service.NewEventServiceServer(WrapEventService(eventService, secret)),
	}
	for _, s := range servers {
		m.Handle(s.PathPrefix(), s)
	}
}
