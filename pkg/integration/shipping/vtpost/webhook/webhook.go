package webhook

import (
	"context"
	"time"

	"o.o/api/main/identity"
	shippingcore "o.o/api/main/shipping"
	"o.o/api/top/types/etc/shipping_provider"
	"o.o/backend/com/etc/logging/shippingwebhook"
	logmodel "o.o/backend/com/etc/logging/shippingwebhook/model"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/shipping/carrier"
	shippingconvert "o.o/backend/com/main/shipping/convert"
	shipmodel "o.o/backend/com/main/shipping/model"
	"o.o/backend/com/main/shipping/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/httpreq"
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/etop/model"
	"o.o/backend/pkg/integration/shipping"
	"o.o/backend/pkg/integration/shipping/vtpost"
	vtpostclient "o.o/backend/pkg/integration/shipping/vtpost/client"
	"o.o/capi/dot"
	"o.o/common/jsonx"
	"o.o/common/l"
)

type (
	String = httpreq.String
	Int    = httpreq.Int
)

var ll = l.New()
var EndStatesCode = []string{"501", "503", "504", "201", "107"}

type Webhook struct {
	db                     *cmsql.Database
	carrier                *vtpost.Carrier
	shipmentManager        *carrier.ShipmentManager
	identityQS             identity.QueryBus
	shippingAggr           shippingcore.CommandBus
	shipmentWebhookLogAggr *shippingwebhook.Aggregate
}

func New(db com.MainDB,
	carrier *vtpost.Carrier,
	shipmentM *carrier.ShipmentManager, identityQ identity.QueryBus,
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
	rt.POST("/webhook/vtpost/callback/:id", wh.Callback)
}

func (wh *Webhook) Callback(c *httpx.Context) (_err error) {
	t0 := time.Now()
	var msg vtpostclient.CallbackOrder
	if err := c.DecodeJson(&msg); err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "VTPost: Can not decode JSON callback")
	}
	ll.Debug("VPPOST callback", l.Object("msg", msg))
	orderData := msg.Data
	statusCode := orderData.OrderStatus
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
		updateFfm, err := vtpost.CalcUpdateFulfillment(ffm, orderData)
		if err != nil {
			return cm.Errorf(cm.FailedPrecondition, err, err.Error()).WithMeta("result", "ignore")
		}
		updateFfm.LastSyncAt = t0
		// UpdateInfo other time
		updateFfm = shipping.CalcOtherTimeBaseOnState(updateFfm, ffm, t0)

		// update shipping fee lines
		weight := orderData.ProductWeight
		updateFeeLinesArgs := &shipping.UpdateShippingFeeLinesArgs{
			FfmID:            ffm.ID,
			Weight:           weight,
			State:            updateFfm.ShippingState,
			ProviderFeeLines: updateFfm.ProviderShippingFeeLines,
		}
		if err := shipping.UpdateShippingFeeLines(ctx, wh.shippingAggr, updateFeeLinesArgs); err != nil {
			ll.S.Errorf("Lỗi cập nhật cước phí VTPost: %v", err.Error())
		}

		note := orderData.Note
		subState := vtpostclient.SubStateMap[statusCode]
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
			Weight:                    weight,
			ClosedAt:                  updateFfm.ClosedAt,
			LastSyncAt:                updateFfm.LastSyncAt,
			ShippingCreatedAt:         updateFfm.ShippingCreatedAt,
			ShippingPickingAt:         updateFfm.ShippingPickingAt,
			ShippingDeliveringAt:      updateFfm.ShippingDeliveringAt,
			ShippingDeliveredAt:       updateFfm.ShippingDeliveredAt,
			ShippingReturningAt:       updateFfm.ShippingReturningAt,
			ShippingReturnedAt:        updateFfm.ShippingReturnedAt,
			ShippingCancelledAt:       updateFfm.ShippingCancelledAt,
			ExternalShippingNote:      dot.String(note),
			ExternalShippingSubState:  dot.String(subState),
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

func (wh *Webhook) validateDataAndGetFfm(ctx context.Context, msg vtpostclient.CallbackOrder) (ffm *shipmodel.Fulfillment, err error) {
	orderData := msg.Data
	query := &modelx.GetFulfillmentQuery{
		ShippingProvider:     shipping_provider.VTPost,
		ExternalShippingCode: orderData.OrderNumber,
	}

	if err := bus.Dispatch(ctx, query); err != nil {
		return nil, cm.MapError(err).
			Wrapf(cm.NotFound, "VTPost: Fulfillment not found: %v", orderData.OrderNumber).
			DefaultInternal().WithMeta("result", "ignore")
	}
	ffm = query.Result
	// gặp các hành trình này 501 giao thành công. 503, tiêu hủy.
	// 504 hoàn thành công. 201 hủy phiếu gửi(Viettelpost thực hiện).
	// 107, hủy đơn(Khách hang thực hiện)
	// => Không update trạng thái đơn nữa.
	if cm.StringsContain(EndStatesCode, ffm.ExternalShippingStateCode) {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "This ffm was done. Cannot update it.").WithMeta("result", "ignore")
	}

	providerServiceID := ffm.ProviderServiceID
	_, _, err = vtpost.ParseServiceID(providerServiceID)
	if err != nil {
		return nil, cm.Errorf(cm.FailedPrecondition, err, "VTPost: Can not parse ProviderServiceID in fulfillment.").WithMeta("result", "ignore")
	}
	return
}

func (wh *Webhook) saveLogsWebhook(ctx context.Context, msg vtpostclient.CallbackOrder, err error, ffm *shipmodel.Fulfillment) {
	orderData := msg.Data
	data, _ := jsonx.Marshal(orderData)
	statusCode := orderData.OrderStatus
	vtpostStatus := vtpostclient.ToVTPostShippingState(statusCode)
	webhookData := &logmodel.ShippingProviderWebhook{
		ID:                       cm.NewID(),
		ShippingProvider:         shipping_provider.VTPost.String(),
		Data:                     data,
		ShippingCode:             orderData.OrderNumber,
		ExternalShippingState:    orderData.StatusName,
		ExternalShippingSubState: vtpostclient.SubStateMap[statusCode],
		Error:                    model.ToError(err),
	}
	if ffm != nil {
		webhookData.ShippingState = vtpostStatus.ToModel(ffm.ShippingState).String()
		webhookData.ConnectionID = ffm.ConnectionID
	}
	if err := wh.shipmentWebhookLogAggr.CreateShippingWebhookLog(ctx, webhookData); err != nil {
		ll.Error("Insert db etop_log error", l.Error(err))
	}
}
