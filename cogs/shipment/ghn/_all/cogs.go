package _all

import (
	"net/http"

	"o.o/api/main/identity"
	shippingcore "o.o/api/main/shipping"
	shippingcarrier "o.o/backend/com/main/shipping/carrier"
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/backend/pkg/integration/shipping/ghn"
	ghnwebhookv1 "o.o/backend/pkg/integration/shipping/ghn/webhook/v1"
	ghnwebhookv2 "o.o/backend/pkg/integration/shipping/ghn/webhook/v2"
)

type GHNWebhookServer *http.Server

func NewGHNWebhookServer(
	cfg ghn.WebhookConfig,
	shipmentManager *shippingcarrier.ShipmentManager,
	identityQuery identity.QueryBus,
	shippingAggr shippingcore.CommandBus,
	webhookv1 *ghnwebhookv1.Webhook,
	webhookv2 *ghnwebhookv2.Webhook,
) GHNWebhookServer {
	rt := httpx.New()
	rt.Use(httpx.RecoverAndLog(true))

	webhookv1.Register(rt)
	webhookv2.Register(rt)
	svr := &http.Server{
		Addr:    cfg.Address(),
		Handler: rt,
	}
	return svr
}
