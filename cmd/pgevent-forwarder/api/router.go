package api

import (
	"etop.vn/backend/pkg/common/httprpc"
	service "etop.vn/backend/zexp/api/root/int/pgevent"
)

// +gen:wrapper=etop.vn/backend/pb/services/pgevent
// +gen:wrapper:package=pgevent

func NewPgeventServer(m httprpc.Muxer, secret string) {
	servers := []httprpc.Server{
		service.NewMiscServiceServer(NewMiscService(miscService, secret)),
		service.NewEventServiceServer(NewEventService(eventService, secret)),
	}
	for _, s := range servers {
		m.Handle(s.PathPrefix(), s)
	}
}
