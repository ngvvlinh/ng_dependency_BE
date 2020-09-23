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
	ninjavanclient "o.o/backend/pkg/integration/shipping/ninjavan/client"
	"o.o/backend/pkg/integration/shipping/ninjavan/update"
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

func New(
	db com.MainDB,
	shipmentM *carrier.ShipmentManager,
	identityQ identity.QueryBus,
	shippingA shippingcore.CommandBus,
	shipmentWebhookLogAggr *shippingwebhook.Aggregate,
	OrderStore sqlstore.OrderStoreInterface,
) *Webhook {
	wh := &Webhook{
		db:                     db,
		shipmentManager:        shipmentM,
		identityQS:             identityQ,
		shippingAggr:           shippingA,
		shipmentWebhookLogAggr: shipmentWebhookLogAggr,
		OrderStore:             OrderStore,
	}
	return wh
}

func (wh *Webhook) Register(rt *httpx.Router) {
	rt.POST("/webhook/ninjavan/callback/:id", wh.Callback)
}

func (wh *Webhook) Callback(c *httpx.Context) (_err error) {
	t0 := time.Now()
	var msg ninjavanclient.CallbackOrder
	if err := c.DecodeJson(&msg); err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "Ninja Van: can not decode JSON callback")
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

		weight := msg.GetWeight(ffm.ChargeableWeight)
		updateFeeLinesArgs := &shipping.UpdateShippingFeeLinesArgs{
			FfmID:  ffm.ID,
			Weight: weight,
			State:  updateFfm.ShippingState,
		}
		if err := shipping.UpdateShippingFeeLines(ctx, wh.shippingAggr, updateFeeLinesArgs); err != nil {
			msg := "–––\n👹 NinjaVan: đơn %v có thay đổi cước phí. Không thể cập nhật. Vui lòng kiểm tra lại. 👹 \nLỗi: %v\n–––"
			ll.SendMessage(fmt.Sprintf(msg, ffm.ShippingCode, err.Error()))
		}

		updateFfmArgs := &shippingcore.UpdateFulfillmentExternalShippingInfoCommand{
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
		if err := wh.shippingAggr.Dispatch(ctx, updateFfmArgs); err != nil {
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

func (wh *Webhook) validateDataAndGetFfm(ctx context.Context, msg ninjavanclient.CallbackOrder) (ffm *shipmodel.Fulfillment, err error) {
	shippingCode := msg.TrackingID
	if shippingCode == "" {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "OrderCode is empty")
	}

	query := &modelx.GetFulfillmentQuery{
		ShippingProvider: shipping_provider.NinjaVan,
		ShippingCode:     shippingCode.String(),
	}
	if err := wh.OrderStore.GetFulfillment(ctx, query); err != nil {
		return nil, cm.MapError(err).
			Wrapf(cm.NotFound, "OrderCode not found: %v", shippingCode).
			DefaultInternal()
	}
	return query.Result, nil
}

func (wh *Webhook) saveLogsWebhook(ctx context.Context, msg ninjavanclient.CallbackOrder, err error, ffm *shipmodel.Fulfillment) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	shippingState := ninjavanclient.State(msg.Status).ToModel()
	webhookData := &logmodel.ShippingProviderWebhook{
		ID:                    cm.NewID(),
		ShippingProvider:      shipping_provider.NinjaVan.String(),
		ShippingCode:          msg.TrackingID.String(),
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
