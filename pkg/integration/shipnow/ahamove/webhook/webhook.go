package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"time"

	"o.o/api/main/ordering"
	"o.o/api/main/shipnow"
	"o.o/api/top/types/etc/shipnow_state"
	"o.o/api/top/types/etc/status4"
	"o.o/api/top/types/etc/status5"
	"o.o/backend/com/etc/logging/shippingwebhook"
	"o.o/backend/com/etc/logging/shippingwebhook/model"
	com "o.o/backend/com/main"
	shipnowmodel "o.o/backend/com/main/shipnow/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/backend/pkg/common/sql/cmsql"
	etopmodel "o.o/backend/pkg/etop/model"
	"o.o/backend/pkg/integration/shipnow/ahamove"
	"o.o/backend/pkg/integration/shipnow/ahamove/client"
	"o.o/backend/pkg/integration/shipping"
	"o.o/capi/dot"
	"o.o/common/l"
)

var ll = l.New()
var PaymentStates = []shipnow_state.State{shipnow_state.StateDelivering, shipnow_state.StateDelivered, shipnow_state.StateReturning, shipnow_state.StateReturned}

type Webhook struct {
	db                     *cmsql.Database
	carrier                *ahamove.Carrier
	shipnowQuery           shipnow.QueryBus
	shipnow                shipnow.CommandBus
	order                  ordering.CommandBus
	orderQuery             ordering.QueryBus
	shipmentWebhookLogAggr *shippingwebhook.Aggregate
}

func New(db com.MainDB,
	carrier *ahamove.Carrier,
	shipnowQS shipnow.QueryBus, shipnowAggr shipnow.CommandBus,
	orderAggr ordering.CommandBus, orderQS ordering.QueryBus,
	shipmentWebhookLogAggr *shippingwebhook.Aggregate,
) *Webhook {
	wh := &Webhook{
		db:                     db,
		carrier:                carrier,
		shipnowQuery:           shipnowQS,
		shipnow:                shipnowAggr,
		order:                  orderAggr,
		orderQuery:             orderQS,
		shipmentWebhookLogAggr: shipmentWebhookLogAggr,
	}
	return wh
}

func (wh *Webhook) Register(rt *httpx.Router) {
	rt.POST("/webhook/ahamove/callback/:id", wh.Callback)
}

func (wh *Webhook) Callback(c *httpx.Context) (_err error) {
	var msg client.Order
	if err := c.DecodeJson(&msg); err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "Can not decode JSON callback")
	}
	ll.Logger.Info("ahamove order webhook", l.Object("msg", msg))
	ctx := c.Req.Context()

	defer func() {
		// save to database etop_log
		wh.saveLogsWebhook(ctx, msg, _err)
	}()

	// 1JFU54-1
	code := strings.Split(msg.ID, "-")[0]
	query := &shipnow.GetShipnowFulfillmentByShippingCodeQuery{
		ShippingCode: code,
	}
	if err := wh.shipnowQuery.Dispatch(ctx, query); err != nil {
		return cm.MapError(err).
			Wrapf(cm.NotFound, "Ahamove: Fulfillment not found: %v", code).DefaultInternal().WithMeta("result", "ignore")
	}
	if err := wh.ProcessAhamoveWebhook(ctx, query.Result.ShipnowFulfillment, msg); err != nil {
		return err
	}
	c.SetResult(map[string]string{"code": "ok"})
	return nil
}

func (wh *Webhook) ProcessAhamoveWebhook(ctx context.Context, ffm *shipnow.ShipnowFulfillment, orderMsg client.Order) error {
	if ffm.Status != status5.Z && ffm.Status != status5.S {
		return cm.Errorf(cm.FailedPrecondition, nil, "Can not update this shipnow").WithMeta("result", "ignore")
	}
	err := wh.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		updateFfm, err := wh.ProcessShipnowFulfillment(ctx, ffm, orderMsg)
		if err != nil {
			return err
		}
		for i, point := range orderMsg.Path {
			if i == 0 {
				// ignore first path: pickup address
				continue
			}
			if err := wh.ProcessOrder(ctx, point, updateFfm.ShippingState, updateFfm.EtopPaymentStatus); err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func IsPaymentState(s shipnow_state.State) bool {
	for _, state := range PaymentStates {
		if state == s {
			return true
		}
	}
	return false
}

func (wh *Webhook) ProcessShipnowFulfillment(ctx context.Context, ffm *shipnow.ShipnowFulfillment, orderMsg client.Order) (*shipnow.ShipnowFulfillment, error) {
	t0 := time.Now()
	status := client.OrderState(orderMsg.Status)
	shippingState := status.ToCoreState()
	shipnowTimestamp := shipping.CalcShipnowTimeBaseOnState(ffm, shippingState, t0)
	paymentStatus := ffm.EtopPaymentStatus

	deliveryPoints := ffm.DeliveryPoints
	for i, point := range orderMsg.Path {
		if i == 0 {
			// ignore first path: pickup address
			continue
		}
		dStatus := client.DeliveryStatus(point.Status)
		orderState := dStatus.ToCoreState(shippingState)
		for _, _point := range deliveryPoints {
			if _point.OrderCode == point.TrackingNumber {
				_point.ShippingState = orderState
				break
			}
		}
	}

	update := &shipnow.UpdateShipnowFulfillmentCarrierInfoCommand{
		ID:                   ffm.ID,
		ShippingState:        shippingState,
		ShippingStatus:       shippingState.ToStatus5(),
		TotalFee:             int(orderMsg.TotalFee),
		ShippingPickingAt:    shipnowTimestamp.ShippingPickingAt,
		ShippingDeliveringAt: shipnowTimestamp.ShippingDeliveringAt,
		ShippingDeliveredAt:  shipnowTimestamp.ShippingDeliveredAt,
		ShippingCancelledAt:  shipnowTimestamp.ShippingCancelledAt,
		FeeLines:             nil, // update if needed
		CarrierFeeLines:      nil, // update if needed
		CancelReason:         orderMsg.CancelComment,
		DriverPhone:          orderMsg.SupplierId,
		DriverName:           orderMsg.SupplierName,
		DeliveryPoints:       deliveryPoints,
	}
	if IsPaymentState(shippingState) && ffm.EtopPaymentStatus != status4.P {
		// EtopPaymentStatus: Ahamove khong doi soat, thanh toán ngay khi lấy hàng
		update.CodEtopTransferedAt = time.Now()
		update.EtopPaymentStatus = status4.P
		paymentStatus = update.EtopPaymentStatus
	}
	update.Status = shipnow.ShipnowStatus(update.ShippingState, paymentStatus)

	if err := wh.shipnow.Dispatch(ctx, update); err != nil {
		return nil, err
	}
	return update.Result, nil
}

func (wh *Webhook) ProcessOrder(ctx context.Context, point *client.DeliveryPoint, shippingState shipnow_state.State, paymentStatus status4.Status) error {
	trackingNumber := point.TrackingNumber
	if trackingNumber == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing tracking number (order_code)").WithMeta("result", "ignore")
	}
	query := &ordering.GetOrderByCodeQuery{
		Code: trackingNumber,
	}
	if err := wh.orderQuery.Dispatch(ctx, query); err != nil {
		if cm.ErrorCode(err) == cm.NotFound {
			return cm.Errorf(cm.NotFound, nil, "NotFound: order does not exist (order_code: %v)", trackingNumber).WithMeta("result", "ignore")
		}
		return err
	}
	orderID := query.Result.ID

	// case shipnow does not assign to any driver
	// => release order (ETOP)
	if point.Status == "" && shippingState == shipnow_state.StateCancelled {
		cmd := &ordering.ReleaseOrdersForFfmCommand{
			OrderIDs: []dot.ID{orderID},
		}
		if err := wh.order.Dispatch(ctx, cmd); err != nil {
			return err
		}
		return nil
	}

	dStatus := client.DeliveryStatus(point.Status)
	orderState := dStatus.ToCoreState(shippingState)
	return wh.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		updateOrder := &ordering.UpdateOrderShippingStatusCommand{
			ID:                         orderID,
			FulfillmentShippingStates:  []string{orderState.String()},
			FulfillmentShippingStatus:  dStatus.ToStatus5(),
			FulfillmentPaymentStatuses: []int{int(paymentStatus)},
			FulfillmentStatuses:        []int{dStatus.ToStatus5().Enum()},
			EtopPaymentStatus:          paymentStatus,
		}

		if err := wh.order.Dispatch(ctx, updateOrder); err != nil {
			return err
		}

		// update payment_status trong đơn hàng
		// xem như đơn hàng đã thanh toán (ko quản lý phiếu thu/chi)
		updateOrderPayment := &ordering.UpdateOrderPaymentStatusCommand{
			OrderID:       orderID,
			PaymentStatus: paymentStatus.Wrap(),
		}
		if err := wh.order.Dispatch(ctx, updateOrderPayment); err != nil {
			return err
		}
		return nil
	})
}

func (wh *Webhook) saveLogsWebhook(ctx context.Context, msg client.Order, err error) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)

	status := client.OrderState(msg.Status)
	shippingState := status.ToCoreState().String()
	webhookData := &model.ShippingProviderWebhook{
		ID:                    cm.NewID(),
		ShippingProvider:      shipnowmodel.Ahamove.String(),
		ShippingCode:          msg.ID,
		ExternalShippingState: msg.Status,
		ShippingState:         shippingState,
		Error:                 etopmodel.ToError(err),
	}
	if err := enc.Encode(msg); err == nil {
		webhookData.Data = buf.Bytes()
	}
	if err := wh.shipmentWebhookLogAggr.CreateShippingWebhookLog(ctx, webhookData); err != nil {
		ll.Error("Insert db etop_log error", l.Error(err))
	}
}
