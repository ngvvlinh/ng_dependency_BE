package crm

import (
	service "o.o/api/top/services/crm"
	"o.o/capi/httprpc"
)

// +gen:wrapper=o.o/api/top/services/crm
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
