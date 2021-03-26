package _all

import (
	"net/http"
	
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/backend/pkg/integration/shipping/ghn"
	ghnwebhookv1 "o.o/backend/pkg/integration/shipping/ghn/webhook/v1"
	ghnwebhookv2 "o.o/backend/pkg/integration/shipping/ghn/webhook/v2"
	njvwebhook "o.o/backend/pkg/integration/shipping/ninjavan/webhook"
)

type ShipmentWebhookServer *http.Server

// TODO(Tuan): all shipment webhook (ghtk, vtpost) use the same port
func NewShipmentWebhookServer(
	cfg ghn.WebhookConfig,
	ghnwebhookv1 *ghnwebhookv1.Webhook,
	ghnwebhookv2 *ghnwebhookv2.Webhook,
	njvwh *njvwebhook.Webhook,
) ShipmentWebhookServer {
	rt := httpx.New()
	rt.Use(httpx.RecoverAndLog(true))

	// ghn
	ghnwebhookv1.Register(rt)
	ghnwebhookv2.Register(rt)
	// ninjavan
	njvwh.Register(rt)
	
	// test callback url for payment
	rt.GET("/callback-url/payment/payme", func(c *httpx.Context) error {
		r := c.Req
		c.SetResult(r.URL.Query())
		return nil
	})

	svr := &http.Server{
		Addr:    cfg.Address(),
		Handler: rt,
	}
	return svr
}
