package vtpostWebhook

import (
	"time"

	"o.o/api/top/types/etc/shipping_provider"
	logmodel "o.o/backend/com/etc/logging/webhook/model"
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
	dbLogs  *cmsql.Database
	carrier *vtpost.Carrier
}

func New(dbLogs *cmsql.Database, carrier *vtpost.Carrier) *Webhook {
	wh := &Webhook{
		dbLogs:  dbLogs,
		carrier: carrier,
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
	ctx := c.Req.Context()
	orderData := msg.Data
	statusCode := orderData.OrderStatus
	vtpostStatus := vtpostclient.ToVTPostShippingState(statusCode)
	var ffm *shipmodel.Fulfillment

	defer func() {
		// save to database etop_log
		data, _ := jsonx.Marshal(orderData)
		webhookData := &logmodel.ShippingProviderWebhook{
			ID:                       cm.NewID(),
			ShippingProvider:         shipping_provider.VTPost.String(),
			Data:                     data,
			ShippingCode:             orderData.OrderNumber,
			ExternalShippingState:    orderData.StatusName,
			ExternalShippingSubState: vtpostclient.SubStateMap[statusCode],
			Error:                    model.ToError(_err),
		}
		if ffm != nil {
			webhookData.ShippingState = vtpostStatus.ToModel(ffm.ShippingState).String()
		}
		if _, err := wh.dbLogs.Insert(webhookData); err != nil {
			ll.Error("Insert db etop_log error", l.Error(err))
		}
	}()

	query := &modelx.GetFulfillmentQuery{
		ShippingProvider:     shipping_provider.VTPost,
		ExternalShippingCode: orderData.OrderNumber,
	}

	if err := bus.Dispatch(ctx, query); err != nil {
		return cm.MapError(err).
			Wrapf(cm.NotFound, "VTPost: Fulfillment not found: %v", orderData.OrderNumber).
			DefaultInternal().WithMeta("result", "ignore")
	}

	ffm = query.Result
	// gặp các hành trình này 501 giao thành công. 503, tiêu hủy.
	// 504 hoàn thành công. 201 hủy phiếu gửi(Viettelpost thực hiện).
	// 107, hủy đơn(Khách hang thực hiện)
	// => Không update trạng thái đơn nữa.
	if cm.StringsContain(EndStatesCode, ffm.ExternalShippingStateCode) {
		return cm.Errorf(cm.FailedPrecondition, nil, "This ffm was done. Cannot update it.").WithMeta("result", "ignore")
	}

	providerServiceID := ffm.ProviderServiceID
	_, _, err := vtpost.ParseServiceID(providerServiceID)
	if err != nil {
		return cm.Errorf(cm.FailedPrecondition, err, "VTPost: Can not parse ProviderServiceID in fulfillment.").WithMeta("result", "ignore")
	}

	updateFfm := vtpost.CalcUpdateFulfillment(ffm, orderData)
	updateFfm.LastSyncAt = t0
	// UpdateInfo other time
	updateFfm = shipping.CalcOtherTimeBaseOnState(updateFfm, ffm, t0)
	note := orderData.Note
	subState := vtpostclient.SubStateMap[statusCode]
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
