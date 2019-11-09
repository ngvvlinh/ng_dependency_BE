package admin

import (
	"etop.vn/backend/pkg/common/httprpc"
	service "etop.vn/backend/zexp/api/root/int/sadmin"
)

// +gen:wrapper=etop.vn/backend/pb/etop/sadmin
// +gen:wrapper:package=sadmin

func NewSadminServer(m httprpc.Muxer) {
	servers := []httprpc.Server{
		service.NewMiscServiceServer(NewMiscService(miscService)),
		service.NewUserServiceServer(NewUserService(userService)),
	}
	for _, s := range servers {
		m.Handle(s.PathPrefix(), s)
	}
}
