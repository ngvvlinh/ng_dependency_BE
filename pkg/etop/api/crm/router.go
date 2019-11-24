package crm

import (
	service "etop.vn/api/top/services/crm"
	"etop.vn/capi/httprpc"
)

// +gen:wrapper=etop.vn/api/top/services/crm
// +gen:wrapper:package=crm

func NewCrmServer(m httprpc.Muxer, secret string) {
	servers := []httprpc.Server{
		service.NewMiscServiceServer(WrapMiscService(miscService, secret)),
		service.NewCrmServiceServer(WrapCrmService(crmService, secret)),
		service.NewVtigerServiceServer(WrapVtigerService(vtigerService)),
		service.NewVhtServiceServer(WrapVhtService(vhtService)),
	}
	for _, s := range servers {
		m.Handle(s.PathPrefix(), s)
	}
}
