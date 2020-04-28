package ghnWebhook

import (
	"bytes"
	"encoding/json"
	"time"

	"o.o/api/top/types/etc/shipping_provider"
	logmodel "o.o/backend/com/etc/logging/webhook/model"
	"o.o/backend/com/main/shipping/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/etop/model"
	"o.o/backend/pkg/integration/shipping"
	"o.o/backend/pkg/integration/shipping/ghn"
	ghnclient "o.o/backend/pkg/integration/shipping/ghn/client"
	"o.o/backend/pkg/integration/shipping/ghn/update"
	"o.o/capi/dot"
	"o.o/common/l"
)

var ll = l.New()

type Webhook struct {
	dbLogs  *cmsql.Database
	carrier *ghn.Carrier
}

func New(dbLogs *cmsql.Database, carrier *ghn.Carrier) *Webhook {
	wh := &Webhook{
		dbLogs:  dbLogs,
		carrier: carrier,
	}
	return wh
}

func (wh *Webhook) Register(rt *httpx.Router) {
	rt.POST("/webhook/ghn/callback/:id", wh.Callback)
}

func (wh *Webhook) Callback(c *httpx.Context) (_err error) {
	t0 := time.Now()
	var msg ghnclient.CallbackOrder
	if err := c.DecodeJson(&msg); err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "GHN: can not decode JSON callback")
	}

	defer func() {
		// save to database etop_log
		buf := new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		webhookData := &logmodel.ShippingProviderWebhook{
			ID:                    cm.NewID(),
			ShippingProvider:      shipping_provider.GHN.String(),
			ShippingCode:          msg.OrderCode.String(),
			ExternalShippingState: msg.CurrentStatus.String(),
			Error:                 model.ToError(_err),
		}
		if err := enc.Encode(msg); err == nil {
			webhookData.Data = buf.Bytes()
		}
		if _, err := wh.dbLogs.Insert(webhookData); err != nil {
			ll.Error("Insert db etop_log error", l.Error(err))
		}
	}()

	if msg.ExternalCode == "" {
		return cm.Errorf(cm.FailedPrecondition, nil, "ExternalCode is empty")
	}
	ffmID, err := dot.ParseID(msg.ExternalCode.String())
	if err != nil {
		return cm.Errorf(cm.FailedPrecondition, nil, "ExternalCode is invalid: %v", msg.ExternalCode)
	}
	if ffmID == 0 {
		return cm.Errorf(cm.FailedPrecondition, nil, "ExternalCode is zero")
	}

	ctx := c.Req.Context()
	query := &modelx.GetFulfillmentQuery{
		ShippingProvider: shipping_provider.GHN,
		FulfillmentID:    ffmID,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return cm.MapError(err).
			Wrapf(cm.NotFound, "ExternalCode not found: %v", ffmID).
			DefaultInternal()
	}
	ffm := query.Result

	updateFfm := update.CalcUpdateFulfillment(ffm, &msg)
	updateFfm.LastSyncAt = t0
	// UpdateInfo other time
	updateFfm = shipping.CalcOtherTimeBaseOnState(updateFfm, ffm, t0)
	updateCmd := &modelx.UpdateFulfillmentCommand{
		Fulfillment: updateFfm,
	}
	if err := bus.Dispatch(ctx, updateCmd); err != nil {
		return err
	}

	c.SetResult(map[string]string{
		"code": "ok",
	})
	return nil
}
