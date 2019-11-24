package integration

import (
	service "etop.vn/api/top/int/integration"
	"etop.vn/capi/httprpc"
)

// +gen:wrapper=etop.vn/api/top/int/integration
// +gen:wrapper:package=integration

func NewIntegrationServer(m httprpc.Muxer) {
	servers := []httprpc.Server{
		service.NewMiscServiceServer(WrapMiscService(miscService)),
		service.NewIntegrationServiceServer(WrapIntegrationService(integrationService)),
	}
	for _, s := range servers {
		m.Handle(s.PathPrefix(), s)
	}
}
