package vtpostWebhook

import (
	"time"

	"etop.vn/api/top/types/etc/shipping_provider"
	logmodel "etop.vn/backend/com/etc/logging/webhook/model"
	"etop.vn/backend/com/main/shipping/modelx"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/httpreq"
	"etop.vn/backend/pkg/common/httpx"
	"etop.vn/backend/pkg/integration/shipping"
	"etop.vn/backend/pkg/integration/shipping/vtpost"
	vtpostclient "etop.vn/backend/pkg/integration/shipping/vtpost/client"
	"etop.vn/capi/dot"
	"etop.vn/common/jsonx"
	"etop.vn/common/l"
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

func (wh *Webhook) Callback(c *httpx.Context) error {
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
	logID := cm.NewID()
	{
		// save to database etop_log
		data, _ := jsonx.Marshal(orderData)
		webhookData := &logmodel.ShippingProviderWebhook{
			ID:                       logID,
			ShippingProvider:         shipping_provider.VTPost.String(),
			Data:                     data,
			ShippingCode:             orderData.OrderNumber,
			ExternalShippingState:    orderData.StatusName,
			ExternalShippingSubState: vtpostclient.SubStateMap[statusCode],
		}
		if _, err := wh.dbLogs.Insert(webhookData); err != nil {
			ll.Error("Insert db etop_log error", l.Error(err))
		}
	}

	query := &modelx.GetFulfillmentQuery{
		ShippingProvider:     shipping_provider.VTPost,
		ExternalShippingCode: orderData.OrderNumber,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return cm.MapError(err).
			Wrapf(cm.NotFound, "VTPost: Fulfillment not found: %v", orderData.OrderNumber).
			DefaultInternal().WithMeta("result", "ignore")
	}
	ffm := query.Result
	// gặp các hành trình này 501 giao thành công. 503, tiêu hủy.
	// 504 hoàn thành công. 201 hủy phiếu gửi(Viettelpost thực hiện).
	// 107, hủy đơn(Khách hang thực hiện)
	// => Không update trạng thái đơn nữa.
	if cm.StringsContain(EndStatesCode, ffm.ExternalShippingStateCode) {
		return cm.Errorf(cm.FailedPrecondition, nil, "This ffm was done. Cannot update it.").WithMeta("result", "ignore")
	}
	{
		// update database etop_log
		webhookData := &logmodel.ShippingProviderWebhook{
			ID:            logID,
			ShippingState: vtpostStatus.ToModel(ffm.ShippingState).String(),
		}
		_, _ = wh.dbLogs.Where("id = ?", logID).Update(webhookData)
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
