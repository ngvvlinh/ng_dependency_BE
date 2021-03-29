package _all

import (
	"io/ioutil"
	"net/http"

	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/backend/pkg/integration/shipping/ghn"
	ghnwebhookv1 "o.o/backend/pkg/integration/shipping/ghn/webhook/v1"
	ghnwebhookv2 "o.o/backend/pkg/integration/shipping/ghn/webhook/v2"
	njvwebhook "o.o/backend/pkg/integration/shipping/ninjavan/webhook"
	"o.o/common/l"
)

var (
	ll = l.New()
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
	rt.POST("/callback-url/payment/payme", func(c *httpx.Context) error {
		body, err := ioutil.ReadAll(c.Req.Body)
		if err != nil {
			return cm.Error(cm.InvalidArgument, err.Error(), err)
		}

		ll.SendMessage("payme: ", string(body))

		defer func() {
			writer := c.SetResultRaw()
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(200)
		}()
		return nil
	})

	svr := &http.Server{
		Addr:    cfg.Address(),
		Handler: rt,
	}
	return svr
}
