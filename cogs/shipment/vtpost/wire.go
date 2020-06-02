package _vtpost

import (
	"net/http"

	"github.com/google/wire"

	"o.o/api/main/identity"
	shippingcore "o.o/api/main/shipping"
	shippingcarrier "o.o/backend/com/main/shipping/carrier"
	"o.o/backend/pkg/common/apifw/httpx"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/integration/shipping/vtpost"
	"o.o/backend/pkg/integration/shipping/vtpost/webhook"
)

var WireSet = wire.NewSet(
	webhook.New,
	NewVTPostWebhookServer,
)

type WebhookConfig cc.HTTP

type VTPostWebhookServer *http.Server

func NewVTPostWebhookServer(
	_cfg WebhookConfig,
	shipmentManager *shippingcarrier.ShipmentManager,
	vtpostCarrier *vtpost.Carrier,
	identityQuery identity.QueryBus,
	shippingAggr shippingcore.CommandBus,
	webhook *webhook.Webhook,
) VTPostWebhookServer {
	rt := httpx.New()
	rt.Use(httpx.RecoverAndLog(true))

	cfg := cc.HTTP(WebhookConfig{})
	webhook.Register(rt)
	svr := &http.Server{
		Addr:    cfg.Address(),
		Handler: rt,
	}
	return svr
}
