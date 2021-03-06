package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"o.o/api/main/identity"
	shippingcore "o.o/api/main/shipping"
	"o.o/api/meta"
	"o.o/api/top/types/etc/shipping_provider"
	"o.o/backend/com/etc/logging/shippingwebhook"
	logmodel "o.o/backend/com/etc/logging/shippingwebhook/model"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/shipping/carrier"
	shippingconvert "o.o/backend/com/main/shipping/convert"
	shipmodel "o.o/backend/com/main/shipping/model"
	"o.o/backend/com/main/shipping/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/etop/model"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/backend/pkg/integration/shipping"
	"o.o/backend/pkg/integration/shipping/ghtk"
	ghtkclient "o.o/backend/pkg/integration/shipping/ghtk/client"
	ghtkupdate "o.o/backend/pkg/integration/shipping/ghtk/update"
	"o.o/capi/dot"
	"o.o/common/l"
)

var ll = l.New().WithChannel(meta.ChannelShipmentCarrier)

type Webhook struct {
	db                     *cmsql.Database
	shipmentManager        *carrier.ShipmentManager
	identityQS             identity.QueryBus
	shippingAggr           shippingcore.CommandBus
	shipmentWebhookLogAggr *shippingwebhook.Aggregate

	OrderStore sqlstore.OrderStoreInterface
}

func New(db com.MainDB,
	shipmentManager *carrier.ShipmentManager,
	identityQ identity.QueryBus, shippingA shippingcore.CommandBus,
	shipmentWebhookLogAggr *shippingwebhook.Aggregate,
	OrderStore sqlstore.OrderStoreInterface,
) *Webhook {
	wh := &Webhook{
		db:                     db,
		shipmentManager:        shipmentManager,
		identityQS:             identityQ,
		shippingAggr:           shippingA,
		shipmentWebhookLogAggr: shipmentWebhookLogAggr,
		OrderStore:             OrderStore,
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

	var ffm *shipmodel.Fulfillment
	var err error
	ctx := c.Req.Context()
	defer func() {
		// save to database etop_log
		wh.saveLogsWebhook(ctx, msg, _err, ffm)
	}()

	ffm, err = wh.validateDataAndGetFfm(ctx, msg)
	if err != nil {
		return err
	}

	ctx, err = shipping.WebhookWlWrapContext(ctx, ffm.ShopID, wh.identityQS)
	if err != nil {
		return err
	}

	err = wh.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		// request get order in GHTK
		// get provider shipping fee lines
		// if error occurred, ignore it, update data from webhook
		updateFfm, _ := wh.shipmentManager.RefreshFulfillment(ctx, ffm)
		// tr???ng th??i ph??? c???a ????n ghtk n???m trong data webhook
		// state_id ??? webhook kh??c state_id khi get order ghtk
		updateFfm, err = ghtkupdate.CalcUpdateFulfillmentFromWebhook(ffm, &msg, updateFfm)
		if err != nil {
			return cm.Errorf(cm.FailedPrecondition, err, err.Error()).WithMeta("result", "ignore")
		}

		updateFfm.LastSyncAt = t0
		// UpdateInfo other time
		updateFfm = shipping.CalcOtherTimeBaseOnState(updateFfm, ffm, t0)

		// update shipping fee lines
		// GHTK tr??? v??? kh???i l?????ng ????n v??? kg
		weight := int(msg.Weight * 1000)
		updateFeeLinesArgs := &shipping.UpdateShippingFeeLinesArgs{
			FfmID:            ffm.ID,
			Weight:           weight,
			State:            updateFfm.ShippingState,
			ProviderFeeLines: updateFfm.ProviderShippingFeeLines,
		}
		if err := shipping.UpdateShippingFeeLines(ctx, wh.shippingAggr, updateFeeLinesArgs); err != nil {
			msg := "?????????\n???? GHTK: ????n %v c?? thay ?????i c?????c ph??. Kh??ng th??? c???p nh???t. Vui l??ng ki???m tra l???i. ???? \nL???i: %v\n?????????"
			ll.SendMessage(fmt.Sprintf(msg, ffm.ShippingCode, err.Error()))
		}

		update := &shippingcore.UpdateFulfillmentExternalShippingInfoCommand{
			FulfillmentID:             ffm.ID,
			ShippingState:             updateFfm.ShippingState,
			ShippingSubstate:          updateFfm.ShippingSubstate,
			ShippingStatus:            updateFfm.ShippingStatus,
			ExternalShippingData:      updateFfm.ExternalShippingData,
			ExternalShippingState:     updateFfm.ExternalShippingState,
			ExternalShippingStatus:    updateFfm.ExternalShippingStatus,
			ExternalShippingUpdatedAt: updateFfm.ExternalShippingUpdatedAt,
			ExternalShippingLogs:      shippingconvert.Convert_shippingmodel_ExternalShippingLogs_shipping_ExternalShippingLogs(updateFfm.ExternalShippingLogs),
			ExternalShippingStateCode: updateFfm.ExternalShippingStateCode,
			Weight:                    weight,
			ClosedAt:                  updateFfm.ClosedAt,
			LastSyncAt:                updateFfm.LastSyncAt,
			ShippingCreatedAt:         updateFfm.ShippingCreatedAt,
			ShippingPickingAt:         updateFfm.ShippingPickingAt,
			ShippingHoldingAt:         updateFfm.ShippingHoldingAt,
			ShippingDeliveringAt:      updateFfm.ShippingDeliveringAt,
			ShippingDeliveredAt:       updateFfm.ShippingDeliveredAt,
			ShippingReturningAt:       updateFfm.ShippingReturningAt,
			ShippingReturnedAt:        updateFfm.ShippingReturnedAt,
			ShippingCancelledAt:       updateFfm.ShippingCancelledAt,
			ExternalShippingNote:      updateFfm.ExternalShippingNote,
			ExternalShippingSubState:  updateFfm.ExternalShippingSubState,
		}
		if err := wh.shippingAggr.Dispatch(ctx, update); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	c.SetResult(map[string]string{"code": "ok"})
	return nil
}

func (wh *Webhook) saveLogsWebhook(ctx context.Context, msg ghtkclient.CallbackOrder, err error, ffm *shipmodel.Fulfillment) {
	statusID := int(msg.StatusID)
	stateID := ghtkclient.StateID(statusID)
	shippingState := stateID.ToModel().String()

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
		Error:                    model.ToError(err),
	}
	if ffm != nil {
		webhookData.ConnectionID = ffm.ConnectionID
	}
	if err := enc.Encode(msg); err == nil {
		webhookData.Data = buf.Bytes()
	}
	if err := wh.shipmentWebhookLogAggr.CreateShippingWebhookLog(ctx, webhookData); err != nil {
		ll.Error("Insert db etop_log error", l.Error(err))
	}
}

func (wh *Webhook) validateDataAndGetFfm(ctx context.Context, msg ghtkclient.CallbackOrder) (ffm *shipmodel.Fulfillment, err error) {
	if msg.PartnerID == "" {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "PartnerID is empty")
	}
	ffmID, err := dot.ParseID(msg.PartnerID.String())
	if err != nil {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "PartnerID is invalid: %v", msg.PartnerID).WithMeta("result", "ignore")
	}
	if ffmID == 0 {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "PartnerID is zero").WithMeta("result", "ignore")
	}

	query := &modelx.GetFulfillmentQuery{
		ShippingProvider: shipping_provider.GHTK,
		FulfillmentID:    ffmID,
	}
	if err := wh.OrderStore.GetFulfillment(ctx, query); err != nil {
		return nil, cm.MapError(err).
			Wrapf(cm.NotFound, "Fulfillment not found: %v", ffmID).
			DefaultInternal().WithMeta("result", "ignore")
	}

	return query.Result, nil
}
