package admin

import (
	service "etop.vn/api/root/int/sadmin"
	"etop.vn/capi/httprpc"
)

// +gen:wrapper=etop.vn/api/root/int/sadmin
// +gen:wrapper:package=sadmin

func NewSadminServer(m httprpc.Muxer) {
	servers := []httprpc.Server{
		service.NewMiscServiceServer(WrapMiscService(miscService)),
		service.NewUserServiceServer(WrapUserService(userService)),
	}
	for _, s := range servers {
		m.Handle(s.PathPrefix(), s)
	}
}
