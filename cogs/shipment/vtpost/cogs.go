package _vtpost

import (
	"net/http"

	"o.o/api/main/identity"
	shippingcore "o.o/api/main/shipping"
	shippingcarrier "o.o/backend/com/main/shipping/carrier"
	"o.o/backend/pkg/common/apifw/httpx"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/integration/shipping/vtpost/webhook"
)

type WebhookConfig cc.HTTP

type VTPostWebhookServer *http.Server

func NewVTPostWebhookServer(
	_cfg WebhookConfig,
	shipmentManager *shippingcarrier.ShipmentManager,
	identityQuery identity.QueryBus,
	shippingAggr shippingcore.CommandBus,
	webhook *webhook.Webhook,
) VTPostWebhookServer {
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
