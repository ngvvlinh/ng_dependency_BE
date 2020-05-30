package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"o.o/api/main/identity"
	shippingcore "o.o/api/main/shipping"
	"o.o/api/top/types/etc/shipping_provider"
	logmodel "o.o/backend/com/etc/logging/webhook/model"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/shipping/carrier"
	shippingconvert "o.o/backend/com/main/shipping/convert"
	shipmodel "o.o/backend/com/main/shipping/model"
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

type MainDB *cmsql.Database // TODO(vu): call the right service
type LogDB *cmsql.Database  // TODO(vu): move to new service

type Webhook struct {
	db              *cmsql.Database
	dbLogs          *cmsql.Database
	carrier         *ghn.Carrier
	shipmentManager *carrier.ShipmentManager
	identityQS      identity.QueryBus
	shippingAggr    shippingcore.CommandBus
}

func New(db com.MainDB, dbLogs com.LogDB, carrier *ghn.Carrier, shipmentM *carrier.ShipmentManager, identityQ identity.QueryBus, shippingA shippingcore.CommandBus) *Webhook {
	wh := &Webhook{
		db:              db,
		dbLogs:          dbLogs,
		carrier:         carrier,
		shipmentManager: shipmentM,
		identityQS:      identityQ,
		shippingAggr:    shippingA,
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
	var ffm *shipmodel.Fulfillment
	defer func() {
		// save to database etop_log
		wh.saveLogsWebhook(msg, ffm, _err)
	}()

	ctx := c.Req.Context()
	ffm, err := wh.validateDataAndGetFfm(ctx, msg)
	if err != nil {
		return err
	}

	ctx, err = shipping.WebhookWlWrapContext(ctx, ffm.ShopID, wh.identityQS)
	if err != nil {
		return err
	}

	err = wh.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		updateFfm, err := update.CalcUpdateFulfillment(ffm, &msg)
		if err != nil {
			return cm.Errorf(cm.FailedPrecondition, err, err.Error()).WithMeta("result", "ignore")
		}
		updateFfm.LastSyncAt = t0
		// UpdateInfo other time
		updateFfm = shipping.CalcOtherTimeBaseOnState(updateFfm, ffm, t0)

		// update shipping fee lines
		if err := shipping.UpdateShippingFeeLines(ctx, wh.shippingAggr, ffm.ID, msg.Weight.Int(), updateFfm.ProviderShippingFeeLines); err != nil {
			ll.S.Errorf("Lỗi cập nhật cước phí GHN: %v", err.Error())
		}

		// update info
		update := &shippingcore.UpdateFulfillmentExternalShippingInfoCommand{
			FulfillmentID:             ffm.ID,
			ShippingState:             updateFfm.ShippingState,
			ShippingStatus:            updateFfm.ShippingStatus,
			ExternalShippingData:      updateFfm.ExternalShippingData,
			ExternalShippingState:     updateFfm.ExternalShippingState,
			ExternalShippingStatus:    updateFfm.ExternalShippingStatus,
			ExternalShippingUpdatedAt: updateFfm.ExternalShippingUpdatedAt,
			ExternalShippingLogs:      shippingconvert.Convert_shippingmodel_ExternalShippingLogs_shipping_ExternalShippingLogs(updateFfm.ExternalShippingLogs),
			ExternalShippingStateCode: updateFfm.ExternalShippingStateCode,
			Weight:                    msg.Weight.Int(),
			ClosedAt:                  updateFfm.ClosedAt,
			LastSyncAt:                updateFfm.LastSyncAt,
			ShippingCreatedAt:         updateFfm.ShippingCreatedAt,
			ShippingPickingAt:         updateFfm.ShippingPickingAt,
			ShippingDeliveringAt:      updateFfm.ShippingDeliveringAt,
			ShippingDeliveredAt:       updateFfm.ShippingDeliveredAt,
			ShippingReturningAt:       updateFfm.ShippingReturningAt,
			ShippingReturnedAt:        updateFfm.ShippingReturnedAt,
			ShippingCancelledAt:       updateFfm.ShippingCancelledAt,
			ExternalShippingNote:      dot.String(updateFfm.ExternalShippingNote),
		}
		if err := wh.shippingAggr.Dispatch(ctx, update); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	c.SetResult(map[string]string{
		"code": "ok",
	})
	return nil
}

func (wh *Webhook) validateDataAndGetFfm(ctx context.Context, msg ghnclient.CallbackOrder) (ffm *shipmodel.Fulfillment, err error) {
	if msg.ExternalCode == "" {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "ExternalCode is empty")
	}
	ffmID, err := dot.ParseID(msg.ExternalCode.String())
	if err != nil {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "ExternalCode is invalid: %v", msg.ExternalCode)
	}
	if ffmID == 0 {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "ExternalCode is zero")
	}

	query := &modelx.GetFulfillmentQuery{
		ShippingProvider: shipping_provider.GHN,
		FulfillmentID:    ffmID,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return nil, cm.MapError(err).
			Wrapf(cm.NotFound, "ExternalCode not found: %v", ffmID).
			DefaultInternal()
	}
	return query.Result, nil
}

func (wh *Webhook) saveLogsWebhook(msg ghnclient.CallbackOrder, ffm *shipmodel.Fulfillment, err error) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	state := ghnclient.State(msg.CurrentStatus)
	var isReturnOrder bool
	if msg.ReturnInfo != "" {
		isReturnOrder = true
	}
	shippingState := state.ToModel(ffm.ShippingState, isReturnOrder)
	webhookData := &logmodel.ShippingProviderWebhook{
		ID:                    cm.NewID(),
		ShippingProvider:      shipping_provider.GHN.String(),
		ShippingCode:          msg.OrderCode.String(),
		ExternalShippingState: msg.CurrentStatus.String(),
		ShippingState:         shippingState.String(),
		Error:                 model.ToError(err),
	}
	if err := enc.Encode(msg); err == nil {
		webhookData.Data = buf.Bytes()
	}
	if _, err := wh.dbLogs.Insert(webhookData); err != nil {
		ll.Error("Insert db etop_log error", l.Error(err))
	}
}
