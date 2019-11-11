package crm

import (
	"etop.vn/backend/pkg/common/httprpc"
	service "etop.vn/backend/zexp/api/root/int/crm"
)

// +gen:wrapper=etop.vn/backend/pb/services/crm
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
