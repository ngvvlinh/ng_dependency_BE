package update

import (
	"fmt"
	"net/http"
	"time"

	"etop.vn/capi/dot"

	"github.com/PuerkitoBio/goquery"

	shipmodel "etop.vn/backend/com/main/shipping/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/integration/shipping"
	ghnclient "etop.vn/backend/pkg/integration/shipping/ghn/client"
	"etop.vn/common/jsonx"
	"etop.vn/common/l"
)

var ll = l.New()

type OrderStateDetail struct {
	Icon    string `json:"icon"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Title   string `json:"title"`
}

type OrderState struct {
	Code               string           `json:"code"`
	GroupID            int              `json:"group_id"`
	CurrentStatus      string           `json:"current_status"`
	Detail             OrderStateDetail `json:"detail"`
	ActionAt           string           `json:"action_at"`
	ActionAtMilisecond int64            `json:"action_at_milisecond"`
	Last               bool             `json:"last"`
	SourceCollect      string           `json:"source_collect"`
	EstDone            string           `json:"est_done"`
}

type ResultOrder struct {
	Code   int    `json:"code"`
	Error  string `json:"Error"`
	Result struct {
		OrderTracking []*OrderState `json:"orderTracking"`
	} `json:"result"`
}

func CalcUpdateFulfillment(ffm *shipmodel.Fulfillment, msg *ghnclient.CallbackOrder, orderGHN *ghnclient.Order) *shipmodel.Fulfillment {
	if !shipping.CanUpdateFulfillmentFromWebhook(ffm) {
		return ffm
	}

	now := time.Now()
	state := ghnclient.State(msg.CurrentStatus)
	data, _ := jsonx.Marshal(msg)

	// GET LOGS
	ffm, _ = SyncTrackingOrder(ffm)

	update := &shipmodel.Fulfillment{
		ID:                        ffm.ID,
		ExternalShippingUpdatedAt: time.Now(),
		ExternalShippingState:     msg.CurrentStatus.String(),
		ExternalShippingStatus:    state.ToStatus5(ffm.ShippingState),
		ExternalShippingData:      data,
		ProviderShippingFeeLines:  ghnclient.CalcAndConvertShippingFeeLines(orderGHN.ShippingOrderCosts),
		ShippingState:             state.ToModel(ffm.ShippingState, msg),
		EtopDiscount:              ffm.EtopDiscount,
		ShippingStatus:            state.ToShippingStatus5(ffm.ShippingState),
		ExternalShippingLogs:      ffm.ExternalShippingLogs,
		ShippingCode:              ffm.ShippingCode,
	}

	if shipping.CanUpdateFulfillmentFeelines(ffm) {
		shippingFeeShopLines := model.GetShippingFeeShopLines(update.ProviderShippingFeeLines, ffm.EtopPriceRule, dot.Int(ffm.EtopAdjustedShippingFeeMain))
		shippingFeeShop := 0
		for _, line := range shippingFeeShopLines {
			shippingFeeShop += line.Cost
		}
		update.ShippingFeeShopLines = shippingFeeShopLines
		update.ShippingFeeShop = shipmodel.CalcShopShippingFee(shippingFeeShop, ffm)
	}

	// Only update status4 if the current status is not ending status
	newStatus := state.ToStatus5(ffm.ShippingState)

	// UpdateInfo ClosedAt
	if newStatus == model.S5Negative || newStatus == model.S5NegSuper || newStatus == model.S5Positive {
		if ffm.ExternalShippingClosedAt.IsZero() {
			update.ClosedAt = now
		}
		if ffm.ClosedAt.IsZero() {
			update.ClosedAt = now
		}
	}
	return update
}

// TODO: refactor, make new client for this method
func SyncTrackingOrders(ffms []*shipmodel.Fulfillment) ([]*shipmodel.Fulfillment, error) {
	rate := time.Second / 30
	burstLimit := 30
	tick := time.NewTicker(rate)
	defer tick.Stop()
	throttle := make(chan time.Time, burstLimit)
	go func() {
		for t := range tick.C {
			select {
			case throttle <- t:
			default:
			}
		} // does not exit after tick.Stop()
	}()

	ch := make(chan error, burstLimit)
	ll.Info("length GHN ffms :: ", l.Int("len", len(ffms)))
	var _ffms []*shipmodel.Fulfillment
	count := 0
	for _, ffm := range ffms {
		<-throttle // rate limit our Service.Method RPCs
		count++
		if count > burstLimit {
			time.Sleep(1 * time.Minute)
			count = 0
		}
		go func(ffm *shipmodel.Fulfillment) (_err error) {
			defer func() {
				ch <- _err
			}()
			ffm, _err = SyncTrackingOrder(ffm)
			if _err == nil {
				_ffms = append(_ffms, ffm)
			}
			return _err
		}(ffm)
	}
	var successCount, errorCount int

	for i, n := 0, len(ffms); i < n; i++ {
		err := <-ch
		if err == nil {
			successCount++
		} else {
			errorCount++
		}
	}
	ll.S.Infof("Sync fulfillment info success: %v/%v, errors %v/%v",
		successCount, len(ffms),
		errorCount, len(ffms))
	return _ffms, nil
}

func SyncTrackingOrder(ffm *shipmodel.Fulfillment) (*shipmodel.Fulfillment, error) {
	url := fmt.Sprintf("https://track.ghn.vn/order/tracking?code=%v", ffm.ShippingCode)
	res, err := http.Get(url)
	if err != nil {
		ll.Error("Error decoding url", l.Error(err), l.String("shipping_code", ffm.ShippingCode))
		return ffm, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		ll.Error("status code error: ", l.Int("StatusCode", res.StatusCode), l.String("Status", res.Status), l.String("shipping_code", ffm.ShippingCode))
		return ffm, cm.Error(cm.Unknown, "Status Code error", nil)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		ll.Error("status code error: %d %s", l.Error(err), l.String("shipping_code", ffm.ShippingCode))
		return ffm, err
	}
	var _err error
	doc.Find("body").Each(func(i int, s *goquery.Selection) {
		orderStr, ok := s.Attr("ng-init")
		if !ok {
			ll.Error("can not get order info", l.String("shipping_code", ffm.ShippingCode))
			_err = cm.Errorf(cm.Internal, nil, "Can not get order info: %v", ffm.ShippingCode)
		} else {
			// remove 'resultOrder='
			orderStr = orderStr[12:]
			var order ResultOrder
			if err = jsonx.Unmarshal([]byte(orderStr), &order); err != nil {
				ll.Error("Lỗi không xác định", l.Error(err), l.String("shipping_code", ffm.ShippingCode))
			}
			var logs []*model.ExternalShippingLog
			for _, orderState := range order.Result.OrderTracking {
				if orderState.GroupID != 5 {
					// trạng thái khác thanh toán
					logs = append(logs, &model.ExternalShippingLog{
						StateText: orderState.Detail.Title,
						Time:      orderState.ActionAt,
						Message:   orderState.Detail.Message,
					})
				}
			}
			ffm.ExternalShippingLogs = logs
		}
	})
	return ffm, _err
}
