package integration

import (
	service "o.o/api/top/int/integration"
	"o.o/capi/httprpc"
)

// +gen:wrapper=o.o/api/top/int/integration
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
