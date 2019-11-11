package integration

import (
	"etop.vn/backend/pkg/common/httprpc"
	service "etop.vn/backend/zexp/api/root/int/integration"
)

// +gen:wrapper=etop.vn/backend/pb/etop/integration
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
