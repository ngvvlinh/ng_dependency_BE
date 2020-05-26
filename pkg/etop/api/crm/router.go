package crm

import (
	service "o.o/api/top/services/crm"
	cc "o.o/backend/pkg/common/config"
	"o.o/capi/httprpc"
)

// +gen:wrapper=o.o/api/top/services/crm
// +gen:wrapper:package=crm

func NewCrmServer(m httprpc.Muxer, secret cc.SecretToken) {
	servers := []httprpc.Server{
		service.NewMiscServiceServer(WrapMiscService(miscService.Clone, string(secret))),
		service.NewCrmServiceServer(WrapCrmService(crmService.Clone, string(secret))),
		service.NewVtigerServiceServer(WrapVtigerService(vtigerService.Clone)),
		service.NewVhtServiceServer(WrapVhtService(vhtService.Clone)),
	}
	for _, s := range servers {
		m.Handle(s.PathPrefix(), s)
	}
}
