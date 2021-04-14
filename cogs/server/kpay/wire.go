package server_kpay

import (
	"github.com/google/wire"
	kpaygatewayserver "o.o/backend/com/external/payment/kpay/gateway/server"
	"o.o/backend/pkg/common/apifw/httpx"
)

var WireSet = wire.NewSet(
	BuildKPayHandler,
)

type KPayHandler httpx.Server

func BuildKPayHandler(server *kpaygatewayserver.Server) KPayHandler {
	rt := httpx.New()
	rt.Use(httpx.RecoverAndLog(false))
	rt.POST("/api/payment/kpay/callback-url", server.Callback)
	return httpx.MakeServer("/api/payment/kpay/", rt)
}
