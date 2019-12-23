package ghtkWebhook

import (
	"bytes"
	"encoding/json"
	"strconv"
	"time"

	"etop.vn/api/main/connectioning"
	"etop.vn/api/top/types/etc/shipping_provider"
	logmodel "etop.vn/backend/com/etc/logging/webhook/model"
	"etop.vn/backend/com/main/shipping/carrier"
	"etop.vn/backend/com/main/shipping/modelx"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/httpx"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/integration/shipping"
	"etop.vn/backend/pkg/integration/shipping/ghtk"
	ghtkclient "etop.vn/backend/pkg/integration/shipping/ghtk/client"
	ghtkdriver "etop.vn/backend/pkg/integration/shipping/ghtk/driver"
	ghtkupdate "etop.vn/backend/pkg/integration/shipping/ghtk/update"
	"etop.vn/capi/dot"
	"etop.vn/common/l"
)

var ll = l.New()

type Webhook struct {
	dbLogs          *cmsql.Database
	carrier         *ghtk.Carrier
	shipmentManager *carrier.ShipmentManager
}

func New(dbLogs *cmsql.Database, carrier *ghtk.Carrier, shipmentManager *carrier.ShipmentManager) *Webhook {
	wh := &Webhook{
		dbLogs:          dbLogs,
		carrier:         carrier,
		shipmentManager: shipmentManager,
	}
	return wh
}

func (wh *Webhook) Register(rt *httpx.Router) {
	rt.POST("/webhook/ghtk/callback/:id", wh.Callback)
}

func (wh *Webhook) Callback(c *httpx.Context) (_err error) {
	t0 := time.Now()
	var msg ghtkclient.CallbackOrder
	if err := c.DecodeJson(&msg); err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "GHTK: can not decode JSON callback")
	}
	ll.Logger.Info("ghtk webhook", l.Object("msg", msg))
	statusID := int(msg.StatusID)
	stateID := ghtkclient.StateID(statusID)
	shippingState := stateID.ToModel().String()

	defer func() {
		// save to database etop_log
		buf := new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		webhookData := &logmodel.ShippingProviderWebhook{
			ID:                       cm.NewID(),
			ShippingProvider:         shipping_provider.GHTK.String(),
			ShippingCode:             ghtk.NormalizeGHTKCode(msg.LabelID.String()),
			ExternalShippingState:    ghtkclient.StateMapping[stateID],
			ExternalShippingSubState: ghtkclient.SubStateMapping[stateID],
			ShippingState:            shippingState,
			Error:                    model.ToError(_err),
		}
		if err := enc.Encode(msg); err == nil {
			webhookData.Data = buf.Bytes()
		}
		if _, err := wh.dbLogs.Insert(webhookData); err != nil {
			ll.Error("Insert db etop_log error", l.Error(err))
		}
	}()

	if msg.PartnerID == "" {
		return cm.Errorf(cm.FailedPrecondition, nil, "PartnerID is empty").WithMeta("result", "ignore")
	}
	ffmID, err := dot.ParseID(msg.PartnerID.String())
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

	// backward compatible
	// set default driver when ffm.ConnectionID = 0
	if ffm.ConnectionID == 0 {
		queryShopConn := &connectioning.GetShopConnectionByIDQuery{
			ConnectionID: connectioning.DefaultTopshipGHTKConnectionID,
		}
		if err := wh.shipmentManager.ConnectionQS.Dispatch(ctx, queryShopConn); err != nil {
			return cm.Errorf(cm.InvalidArgument, err, "Không thể lấy default driver cho ghtk, err = %v", err)
		}
		cfg := ghtkclient.GhtkAccount{
			Token: queryShopConn.Result.Token,
		}
		driver := ghtkdriver.New(wh.shipmentManager.Env, cfg, wh.shipmentManager.LocationQS)
		wh.shipmentManager.SetDriver(driver)

		defer func() {
			wh.shipmentManager.ResetDriver()
		}()
	}

	updateFfm, err := wh.shipmentManager.UpdateFulfillment(ctx, ffm)
	if err != nil {
		return err
	}
	// trạng thái phụ của đơn ghtk nằm trong data webhook
	ghtkupdate.CalcUpdateFulfillmentFromWebhook(ffm, &msg, updateFfm)

	updateFfm.LastSyncAt = t0
	// UpdateInfo other time
	updateFfm = shipping.CalcOtherTimeBaseOnState(updateFfm, ffm, t0)
	// Thêm trạng thái đơn vào note
	note, _ := strconv.Unquote("\"" + msg.Reason.String() + "\"")
	subState := ghtkclient.SubStateMapping[stateID]
	updateCmd := &modelx.UpdateFulfillmentCommand{
		Fulfillment:              updateFfm,
		ExternalShippingNote:     dot.String(note),
		ExternalShippingSubState: dot.String(subState),
	}
	if err := bus.Dispatch(ctx, updateCmd); err != nil {
		return err
	}

	c.SetResult(map[string]string{"code": "ok"})
	return nil
}
