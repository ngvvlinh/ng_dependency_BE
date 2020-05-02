package admin

import (
	service "o.o/api/top/int/sadmin"
	"o.o/capi/httprpc"
)

// +gen:wrapper=o.o/api/top/int/sadmin
// +gen:wrapper:package=sadmin

func NewSadminServer(m httprpc.Muxer) {
	servers := []httprpc.Server{
		service.NewMiscServiceServer(WrapMiscService(miscService.Clone)),
		service.NewUserServiceServer(WrapUserService(userService.Clone)),
	}
	for _, s := range servers {
		m.Handle(s.PathPrefix(), s)
	}
}
