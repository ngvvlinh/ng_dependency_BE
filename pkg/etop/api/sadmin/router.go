package admin

import (
	"etop.vn/backend/pkg/common/httprpc"
	service "etop.vn/backend/zexp/api/root/int/sadmin"
)

// +gen:wrapper=etop.vn/backend/zexp/api/root/int/sadmin
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
