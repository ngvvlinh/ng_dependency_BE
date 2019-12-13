package ahamovewebhook

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"time"

	"etop.vn/api/main/ordering"
	"etop.vn/api/main/shipnow"
	shipnowtypes "etop.vn/api/main/shipnow/types"
	"etop.vn/api/top/types/etc/status4"
	"etop.vn/api/top/types/etc/status5"
	"etop.vn/backend/com/etc/logging/webhook/model"
	shipnowmodel "etop.vn/backend/com/main/shipnow/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/httpx"
	"etop.vn/backend/pkg/integration/shipnow/ahamove"
	"etop.vn/backend/pkg/integration/shipnow/ahamove/client"
	"etop.vn/backend/pkg/integration/shipping"
	"etop.vn/capi/dot"
	"etop.vn/common/l"
)

var ll = l.New()
var PaymentStates = []shipnowtypes.State{shipnowtypes.StateDelivering, shipnowtypes.StateDelivered, shipnowtypes.StateReturning, shipnowtypes.StateReturned}

type Webhook struct {
	db           cmsql.Transactioner
	dbLogs       *cmsql.Database
	carrier      *ahamove.Carrier
	shipnowQuery shipnow.QueryBus
	shipnow      shipnow.CommandBus
	order        ordering.CommandBus
	orderQuery   ordering.QueryBus
}

func New(db *cmsql.Database, dbLogs *cmsql.Database, carrier *ahamove.Carrier, shipnowQS shipnow.QueryBus, shipnowAggr shipnow.CommandBus, orderAggr ordering.CommandBus, orderQS ordering.QueryBus) *Webhook {
	wh := &Webhook{
		db:           db,
		dbLogs:       dbLogs,
		carrier:      carrier,
		shipnowQuery: shipnowQS,
		shipnow:      shipnowAggr,
		order:        orderAggr,
		orderQuery:   orderQS,
	}
	return wh
}

func (wh *Webhook) Register(rt *httpx.Router) {
	rt.POST("/webhook/ahamove/callback/:id", wh.Callback)
}

func (wh *Webhook) Callback(c *httpx.Context) error {
	// t0 := time.Now()
	var msg client.Order
	if err := c.DecodeJson(&msg); err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "Can not decode JSON callback")
	}
	ll.Logger.Info("ahamove order webhook", l.Object("msg", msg))
	status := client.OrderState(msg.Status)
	shippingState := status.ToCoreState().String()
	{
		// save to database etop_log
		buf := new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		webhookData := &model.ShippingProviderWebhook{
			ID:                    cm.NewID(),
			ShippingProvider:      shipnowmodel.Ahamove.ToString(),
			ShippingCode:          msg.ID,
			ExternalShippingState: msg.Status,
			ShippingState:         shippingState,
		}
		if err := enc.Encode(msg); err == nil {
			webhookData.Data = buf.Bytes()
		}
		if _, err := wh.dbLogs.Insert(webhookData); err != nil {
			ll.Error("Insert db etop_log error", l.Error(err))
		}
	}

	ctx := c.Req.Context()
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

func IsPaymentState(s shipnowtypes.State) bool {
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

	update := &shipnow.UpdateShipnowFulfillmentCarrierInfoCommand{
		Id:                   ffm.Id,
		ShippingState:        shippingState,
		ShippingStatus:       shipnowtypes.StateToStatus5(shippingState),
		TotalFee:             int(orderMsg.TotalFee),
		ShippingPickingAt:    shipnowTimestamp.ShippingPickingAt,
		ShippingDeliveringAt: shipnowTimestamp.ShippingDeliveringAt,
		ShippingDeliveredAt:  shipnowTimestamp.ShippingDeliveredAt,
		ShippingCancelledAt:  shipnowTimestamp.ShippingCancelledAt,
		FeeLines:             nil, // update if needed
		CarrierFeeLines:      nil, // update if needed
		CancelReason:         orderMsg.CancelComment,
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

func (wh *Webhook) ProcessOrder(ctx context.Context, point *client.DeliveryPoint, shippingState shipnowtypes.State, paymentStatus status4.Status) error {
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
	if point.Status == "" && shippingState == shipnowtypes.StateCancelled {
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
	updateOrder := &ordering.UpdateOrderShippingStatusCommand{
		ID:                         orderID,
		FulfillmentShippingStates:  []string{orderState.String()},
		FulfillmentShippingStatus:  dStatus.ToStatus5(),
		FulfillmentPaymentStatuses: []int{int(paymentStatus)},
		EtopPaymentStatus:          paymentStatus,
	}

	if err := wh.order.Dispatch(ctx, updateOrder); err != nil {
		return err
	}
	return nil
}
