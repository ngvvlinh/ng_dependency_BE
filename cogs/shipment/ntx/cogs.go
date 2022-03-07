package _ntx

import (
	"net/http"

	"o.o/api/main/identity"
	shippingcore "o.o/api/main/shipping"
	shippingcarrier "o.o/backend/com/main/shipping/carrier"
	"o.o/backend/pkg/common/apifw/httpx"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/integration/shipping/ntx/webhook"
)

type WebhookConfig cc.HTTP

type NTXWebhookServer *http.Server

func NewNTXWebhookServer(
	_cfg WebhookConfig,
	shipmentManager *shippingcarrier.ShipmentManager,
	identityQuery identity.QueryBus,
	shippingAggr shippingcore.CommandBus,
	webhook *webhook.Webhook,
) NTXWebhookServer {
	rt := httpx.New()
	rt.Use(httpx.RecoverAndLog(true))

	cfg := cc.HTTP(_cfg)
	webhook.Register(rt)
	svr := &http.Server{
		Addr:    cfg.Address(),
		Handler: rt,
	}
	return svr
}
