package ghtkWebhook

import (
	"bytes"
	"encoding/json"
	"strconv"
	"time"

	logmodel "etop.vn/backend/com/etc/log/webhook/model"
	"etop.vn/backend/com/main/shipping/modelx"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/httpx"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/integration/shipping"
	"etop.vn/backend/pkg/integration/shipping/ghtk"
	ghtkclient "etop.vn/backend/pkg/integration/shipping/ghtk/client"
	"etop.vn/common/l"
)

var ll = l.New()

type Webhook struct {
	dbLogs  *cmsql.Database
	carrier *ghtk.Carrier
}

func New(dbLogs *cmsql.Database, carrier *ghtk.Carrier) *Webhook {
	wh := &Webhook{
		dbLogs:  dbLogs,
		carrier: carrier,
	}
	return wh
}

func (wh *Webhook) Register(rt *httpx.Router) {
	rt.POST("/webhook/ghtk/callback/:id", wh.Callback)
}

func (wh *Webhook) Callback(c *httpx.Context) error {
	t0 := time.Now()
	var msg ghtkclient.CallbackOrder
	if err := c.DecodeJson(&msg); err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "Can not decode JSON callback")
	}
	statusID := int(msg.StatusID)
	stateID := ghtkclient.StateID(statusID)
	shippingState := string(stateID.ToModel())
	{
		// save to database etop_log
		buf := new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		webhookData := &logmodel.ShippingProviderWebhook{
			ID:                       cm.NewID(),
			ShippingProvider:         model.TypeGHTK.ToString(),
			ShippingCode:             ghtk.NormalizeGHTKCode(msg.LabelID.String()),
			ExternalShippingState:    ghtkclient.StateMapping[stateID],
			ExternalShippingSubState: ghtkclient.SubStateMapping[stateID],
			ShippingState:            shippingState,
		}
		if err := enc.Encode(msg); err == nil {
			webhookData.Data = buf.Bytes()
		}
		if _, err := wh.dbLogs.Insert(webhookData); err != nil {
			ll.Error("Insert db etop_log error", l.Error(err))
		}
	}

	if msg.PartnerID == "" {
		return cm.Errorf(cm.FailedPrecondition, nil, "PartnerID is empty").WithMeta("result", "ignore")
	}
	ffmID, err := strconv.ParseInt(msg.PartnerID.String(), 10, 64)
	if err != nil {
		return cm.Errorf(cm.FailedPrecondition, nil, "PartnerID is invalid: %v", msg.PartnerID).WithMeta("result", "ignore")
	}
	if ffmID == 0 {
		return cm.Errorf(cm.FailedPrecondition, nil, "PartnerID is zero").WithMeta("result", "ignore")
	}

	ctx := c.Req.Context()
	query := &modelx.GetFulfillmentQuery{
		FulfillmentID: ffmID,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return cm.MapError(err).
			Wrapf(cm.NotFound, "Fulfillment not found: %v", ffmID).
			DefaultInternal().WithMeta("result", "ignore")
	}

	ffm := query.Result
	providerServiceID := ffm.ProviderServiceID
	_, _, err = ghtk.ParseServiceID(providerServiceID)
	if err != nil {
		return cm.Errorf(cm.FailedPrecondition, err, "Can not parse ProviderServiceID in fulfillment.").WithMeta("result", "ignore")
	}
	// get order info to update service fee
	ghtkCmd := &ghtk.GetOrderCommand{
		ServiceID: ffm.ProviderServiceID,
		LabelID:   msg.LabelID.String(),
	}
	if err := wh.carrier.GetOrder(ctx, ghtkCmd); err != nil {
		return err
	}

	updateFfm := ghtk.CalcUpdateFulfillment(ffm, &msg, &ghtkCmd.Result.Order)
	updateFfm.LastSyncAt = t0
	// UpdateInfo other time
	updateFfm = shipping.CalcOtherTimeBaseOnState(updateFfm, ffm, t0)
	// Thêm trạng thái đơn vào note
	note, _ := strconv.Unquote("\"" + msg.Reason.String() + "\"")
	subState := ghtkclient.SubStateMapping[stateID]
	updateCmd := &modelx.UpdateFulfillmentCommand{
		Fulfillment:              updateFfm,
		ExternalShippingNote:     cm.PString(note),
		ExternalShippingSubState: cm.PString(subState),
	}
	if err := bus.Dispatch(ctx, updateCmd); err != nil {
		return err
	}

	c.SetResult(map[string]string{"code": "ok"})
	return nil
}
