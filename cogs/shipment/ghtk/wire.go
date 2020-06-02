package _ghtk

import (
	"net/http"

	"github.com/google/wire"

	"o.o/api/main/identity"
	shippingcore "o.o/api/main/shipping"
	shippingcarrier "o.o/backend/com/main/shipping/carrier"
	"o.o/backend/pkg/common/apifw/httpx"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/integration/shipping/ghtk"
	"o.o/backend/pkg/integration/shipping/ghtk/webhook"
)

var WireSet = wire.NewSet(
	webhook.New,
	NewGHTKWebhookServer,
)

type WebhookConfig cc.HTTP

type GHTKWebhookServer *http.Server

func NewGHTKWebhookServer(
	_cfg WebhookConfig,
	shipmentManager *shippingcarrier.ShipmentManager,
	ghtkCarrier *ghtk.Carrier,
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
