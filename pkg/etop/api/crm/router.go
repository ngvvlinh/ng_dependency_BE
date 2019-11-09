package crm

import (
	"etop.vn/backend/pkg/common/httprpc"
	service "etop.vn/backend/zexp/api/root/int/crm"
)

// +gen:wrapper=etop.vn/backend/pb/services/crm
// +gen:wrapper:package=crm

func NewCrmServer(m httprpc.Muxer, secret string) {
	servers := []httprpc.Server{
		service.NewMiscServiceServer(NewMiscService(miscService, secret)),
		service.NewCrmServiceServer(NewCrmService(crmService, secret)),
		service.NewVtigerServiceServer(NewVtigerService(vtigerService)),
		service.NewVhtServiceServer(NewVhtService(vhtService)),
	}
	for _, s := range servers {
		m.Handle(s.PathPrefix(), s)
	}
}
