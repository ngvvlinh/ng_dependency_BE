package server_vtpay

import (
	"github.com/google/wire"

	vtpaygatewayaggregate "o.o/backend/com/external/payment/vtpay/gateway/aggregate"
	vtpaygatewayserver "o.o/backend/com/external/payment/vtpay/gateway/server"
	"o.o/backend/pkg/common/apifw/httpx"
)

var WireSet = wire.NewSet(
	BuildVTPayHandler,
)

type VTPayHandler httpx.Server

func BuildVTPayHandler(server *vtpaygatewayserver.Server) VTPayHandler {
	buildRoute := vtpaygatewayaggregate.BuildGatewayRoute
	rt := httpx.New()
	rt.Use(httpx.RecoverAndLog(false))
	rt.POST("/api"+buildRoute(vtpaygatewayaggregate.PathValidateTransaction), server.ValidateTransaction)
	rt.POST("/api"+buildRoute(vtpaygatewayaggregate.PathGetResult), server.GetResult)
	return httpx.MakeServer("/api/payment/vtpay/", rt)
}
