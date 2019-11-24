package ghnWebhook

import (
	"bytes"
	"encoding/json"
	"time"

	logmodel "etop.vn/backend/com/etc/log/webhook/model"
	"etop.vn/backend/com/main/shipping/modelx"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/httpx"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/integration/shipping"
	"etop.vn/backend/pkg/integration/shipping/ghn"
	ghnclient "etop.vn/backend/pkg/integration/shipping/ghn/client"
	"etop.vn/backend/pkg/integration/shipping/ghn/update"
	"etop.vn/capi/util"
	"etop.vn/common/l"
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

func (wh *Webhook) Callback(c *httpx.Context) error {
	t0 := time.Now()
	var msg ghnclient.CallbackOrder
	if err := c.DecodeJson(&msg); err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "...")
	}

	{
		// save to database etop_log
		buf := new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		webhookData := &logmodel.ShippingProviderWebhook{
			ID:                    cm.NewID(),
			ShippingProvider:      model.TypeGHN.ToString(),
			ShippingCode:          msg.OrderCode.String(),
			ExternalShippingState: msg.CurrentStatus.String(),
		}
		if err := enc.Encode(msg); err == nil {
			webhookData.Data = buf.Bytes()
		}
		if _, err := wh.dbLogs.Insert(webhookData); err != nil {
			ll.Error("Insert db etop_log error", l.Error(err))
		}
	}

	if msg.ExternalCode == "" {
		return cm.Errorf(cm.FailedPrecondition, nil, "ExternalCode is empty")
	}
	ffmID, err := util.ParseID(msg.ExternalCode.String())
	if err != nil {
		return cm.Errorf(cm.FailedPrecondition, nil, "ExternalCode is invalid: %v", msg.ExternalCode)
	}
	if ffmID == 0 {
		return cm.Errorf(cm.FailedPrecondition, nil, "ExternalCode is zero")
	}

	ctx := c.Req.Context()
	query := &modelx.GetFulfillmentQuery{
		FulfillmentID: ffmID,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return cm.MapError(err).
			Wrapf(cm.NotFound, "ExternalCode not found: %v", ffmID).
			DefaultInternal()
	}
	ffm := query.Result

	// get order GHN info to update service fee
	providerServiceID := ffm.ProviderServiceID
	ghnCmd := &ghn.RequestGetOrderCommand{
		ServiceID: providerServiceID,
		Request: &ghnclient.OrderCodeRequest{
			OrderCode: query.Result.ExternalShippingCode,
		},
	}
	if err := wh.carrier.GetOrder(ctx, ghnCmd); err != nil {
		return err
	}

	updateFfm := update.CalcUpdateFulfillment(ffm, &msg, ghnCmd.Result)
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
