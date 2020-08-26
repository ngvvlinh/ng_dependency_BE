package v2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"o.o/api/main/identity"
	shippingcore "o.o/api/main/shipping"
	"o.o/api/meta"
	"o.o/api/top/types/etc/connection_type"
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
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/etop/model"
	"o.o/backend/pkg/integration/shipping"
	"o.o/backend/pkg/integration/shipping/ghn"
	ghnclient "o.o/backend/pkg/integration/shipping/ghn/clientv2"
	update "o.o/backend/pkg/integration/shipping/ghn/update/v2"
	"o.o/capi/dot"
	"o.o/common/l"
)

var ll = l.New().WithChannel(meta.ChannelShipmentCarrier)

type MainDB *cmsql.Database // TODO(vu): call the right service

type Webhook struct {
	db                     *cmsql.Database
	carrier                *ghn.Carrier
	shipmentManager        *carrier.ShipmentManager
	identityQS             identity.QueryBus
	shippingAggr           shippingcore.CommandBus
	shipmentWebhookLogAggr *shippingwebhook.Aggregate
}

func New(
	db com.MainDB,
	carrier *ghn.Carrier,
	shipmentM *carrier.ShipmentManager,
	identityQ identity.QueryBus,
	shippingA shippingcore.CommandBus,
	shipmentWebhookLogAggr *shippingwebhook.Aggregate,
) *Webhook {
	wh := &Webhook{
		db:                     db,
		carrier:                carrier,
		shipmentManager:        shipmentM,
		identityQS:             identityQ,
		shippingAggr:           shippingA,
		shipmentWebhookLogAggr: shipmentWebhookLogAggr,
	}
	return wh
}

func (wh *Webhook) Register(rt *httpx.Router) {
	rt.POST("/webhook/ghn/v2/callback/:id", wh.Callback)
}

func (wh *Webhook) Callback(c *httpx.Context) (_err error) {
	t0 := time.Now()
	var msg ghnclient.CallbackOrder
	if err := c.DecodeJson(&msg); err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "GHN: can not decode JSON callback")
	}
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
		updateFfm, err := update.CalcUpdateFulfillment(ffm, &msg)
		if err != nil {
			return cm.Errorf(cm.FailedPrecondition, err, err.Error()).WithMeta("result", "ignore")
		}
		updateFfm.LastSyncAt = t0
		// UpdateInfo other time
		updateFfm = shipping.CalcOtherTimeBaseOnState(updateFfm, ffm, t0)

		// update shipping fee lines
		updateFeeLinesArgs := &shipping.UpdateShippingFeeLinesArgs{
			FfmID:            ffm.ID,
			Weight:           msg.Weight.Int(),
			State:            updateFfm.ShippingState,
			ProviderFeeLines: updateFfm.ProviderShippingFeeLines,
		}
		if err := shipping.UpdateShippingFeeLines(ctx, wh.shippingAggr, updateFeeLinesArgs); err != nil {
			msg := "–––\n👹 GHN: đơn %v có thay đổi cước phí. Không thể cập nhật. Vui lòng kiểm tra lại. 👹\n- Weight: %v\n- State: %v\n- Lỗi: %v\n–––"
			ll.SendMessage(fmt.Sprintf(msg, ffm.ShippingCode, updateFeeLinesArgs.Weight, updateFeeLinesArgs.State, err.Error()))
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
			ShippingHoldingAt:         updateFfm.ShippingHoldingAt,
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

		// updateCOD
		updateCODAmountArgs := &shipping.UpdateFfmCODAmountArgs{
			NewCODAmount:  msg.CODAmount.Int(),
			Ffm:           ffm,
			CarrierName:   shipping_provider.GHN.String(),
			ShippingState: updateFfm.ShippingState,
		}
		shipping.ValidateAndUpdateFulfillmentCOD(ctx, wh.shippingAggr, updateCODAmountArgs)

		// đối soát GHN direct
		if err := wh.updateFulfillmentsCODTransferedAt(ctx, msg, ffm); err != nil {
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

func (wh *Webhook) updateFulfillmentsCODTransferedAt(ctx context.Context, msg ghnclient.CallbackOrder, ffm *shipmodel.Fulfillment) error {
	if !msg.CODTransferDate.IsZero() && ffm.ConnectionMethod == connection_type.ConnectionMethodDirect {
		updateCODTransferDateCmd := &shippingcore.UpdateFulfillmentsCODTransferedAtCommand{
			FulfillmentIDs:  []dot.ID{ffm.ID},
			CODTransferedAt: msg.CODTransferDate.ToTime(),
		}
		if err := wh.shippingAggr.Dispatch(ctx, updateCODTransferDateCmd); err != nil {
			return err
		}
	}
	return nil
}

func (wh *Webhook) validateDataAndGetFfm(ctx context.Context, msg ghnclient.CallbackOrder) (ffm *shipmodel.Fulfillment, err error) {
	orderCode := msg.OrderCode
	if orderCode == "" {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "OrderCode is empty")
	}

	query := &modelx.GetFulfillmentQuery{
		ShippingProvider: shipping_provider.GHN,
		ShippingCode:     orderCode.String(),
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return nil, cm.MapError(err).
			Wrapf(cm.NotFound, "OrderCode not found: %v", orderCode).
			DefaultInternal()
	}
	return query.Result, nil
}

func (wh *Webhook) saveLogsWebhook(ctx context.Context, msg ghnclient.CallbackOrder, err error, ffm *shipmodel.Fulfillment) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	shippingState := ghnclient.State(msg.Status).ToModel()
	webhookData := &logmodel.ShippingProviderWebhook{
		ID:                    cm.NewID(),
		ShippingProvider:      shipping_provider.GHN.String(),
		ShippingCode:          msg.OrderCode.String(),
		ExternalShippingState: msg.Status.String(),
		ShippingState:         shippingState.String(),
		Error:                 model.ToError(err),
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
