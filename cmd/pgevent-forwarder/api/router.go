package api

import (
	"etop.vn/backend/pkg/common/httprpc"
	service "etop.vn/backend/zexp/api/root/int/pgevent"
)

// +gen:wrapper=etop.vn/backend/pb/services/pgevent
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
