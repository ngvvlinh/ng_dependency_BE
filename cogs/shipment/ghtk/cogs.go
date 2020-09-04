package _ghtk

import (
	"net/http"

	"o.o/api/main/identity"
	shippingcore "o.o/api/main/shipping"
	shippingcarrier "o.o/backend/com/main/shipping/carrier"
	"o.o/backend/pkg/common/apifw/httpx"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/integration/shipping/ghtk/webhook"
)

type WebhookConfig cc.HTTP

type GHTKWebhookServer *http.Server

func NewGHTKWebhookServer(
	_cfg WebhookConfig,
	shipmentManager *shippingcarrier.ShipmentManager,
	identityQuery identity.QueryBus,
	shippingAggr shippingcore.CommandBus,
	webhook *webhook.Webhook,
) GHTKWebhookServer {
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
