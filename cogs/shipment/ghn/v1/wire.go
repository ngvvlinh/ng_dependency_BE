package v1

import (
	"net/http"

	"github.com/google/wire"

	"o.o/api/main/identity"
	shippingcore "o.o/api/main/shipping"
	shippingcarrier "o.o/backend/com/main/shipping/carrier"
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/backend/pkg/integration/shipping/ghn"
	ghnwebhook "o.o/backend/pkg/integration/shipping/ghn/webhook/v1"
)

var WireSet = wire.NewSet(
	ghnwebhook.WireSet,
	NewGHNWebhookServer,
)

type GHNWebhookServer *http.Server

func NewGHNWebhookServer(
	cfg ghn.WebhookConfig,
	shipmentManager *shippingcarrier.ShipmentManager,
	ghnCarrier *ghn.Carrier,
	identityQuery identity.QueryBus,
	shippingAggr shippingcore.CommandBus,
	webhook *ghnwebhook.Webhook,
) GHNWebhookServer {
	rt := httpx.New()
	rt.Use(httpx.RecoverAndLog(true))

	webhook.Register(rt)
	svr := &http.Server{
		Addr:    cfg.Address(),
		Handler: rt,
	}
	return svr
}
