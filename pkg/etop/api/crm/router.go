package crm

import (
	service "o.o/api/top/services/crm"
	"o.o/capi/httprpc"
)

// +gen:wrapper=o.o/api/top/services/crm
// +gen:wrapper:package=crm

func NewCrmServer(m httprpc.Muxer, secret string) {
	servers := []httprpc.Server{
		service.NewMiscServiceServer(WrapMiscService(miscService.Clone, secret)),
		service.NewCrmServiceServer(WrapCrmService(crmService.Clone, secret)),
		service.NewVtigerServiceServer(WrapVtigerService(vtigerService.Clone)),
		service.NewVhtServiceServer(WrapVhtService(vhtService.Clone)),
	}
	for _, s := range servers {
		m.Handle(s.PathPrefix(), s)
	}
}
