package ghtk

import (
	"strconv"
	"time"

	"o.o/api/main/location"
	"o.o/api/top/types/etc/shipping_fee_type"
	"o.o/api/top/types/etc/status5"
	shipmodel "o.o/backend/com/main/shipping/model"
	shippingsharemodel "o.o/backend/com/main/shipping/sharemodel"
	"o.o/backend/pkg/integration/shipping"
	ghtkclient "o.o/backend/pkg/integration/shipping/ghtk/client"
	"o.o/capi/dot"
	"o.o/common/jsonx"
	"o.o/common/l"
)

var ll = l.New()

func CalcUpdateFulfillment(ffm *shipmodel.Fulfillment, msg *ghtkclient.CallbackOrder, ghtkOrder *ghtkclient.OrderInfo) *shipmodel.Fulfillment {
	if !shipping.CanUpdateFulfillment(ffm) {
		return ffm
	}

	now := time.Now()
	data, _ := jsonx.Marshal(ghtkOrder)
	var statusID int
	if msg == nil {
		statusID, _ = strconv.Atoi(ghtkOrder.Status.String())
	} else {
		statusID = int(msg.StatusID)
	}
	stateID := ghtkclient.StateID(statusID)
	update := &shipmodel.Fulfillment{
		ID:                        ffm.ID,
		ExternalShippingUpdatedAt: now,
		ExternalShippingData:      data,
		ExternalShippingState:     ghtkclient.StateMapping[stateID],
		ExternalShippingStatus:    stateID.ToStatus5(),
		ExternalShippingCode:      ghtkOrder.LabelID.String(),
		ShippingState:             stateID.ToModel(),
		EtopDiscount:              ffm.EtopDiscount,
		ShippingStatus:            stateID.ToStatus5(),
	}

	// make sure can not update ffm's shipping fee when it belong to a money transaction
	if shipping.CanUpdateFulfillmentFeelines(ffm) {
		update.ProviderShippingFeeLines = CalcAndConvertShippingFeeLines(ghtkOrder)
		shippingFeeShopLines := shippingsharemodel.GetShippingFeeShopLines(update.ProviderShippingFeeLines, ffm.EtopPriceRule, dot.Int(ffm.EtopAdjustedShippingFeeMain))
		shippingFeeShop := 0
		for _, line := range shippingFeeShopLines {
			shippingFeeShop += line.Cost
		}
		update.ShippingFeeShopLines = shippingFeeShopLines
		update.ShippingFeeShop = shipmodel.CalcShopShippingFee(shippingFeeShop, ffm)
	}

	// Only update status4 if the current status is not ending status
	newStatus := stateID.ToStatus5()
	// UpdateInfo ClosedAt
	if newStatus == status5.N || newStatus == status5.NS || newStatus == status5.P {
		if ffm.ExternalShippingClosedAt.IsZero() {
			update.ClosedAt = now
		}
		if ffm.ClosedAt.IsZero() {
			update.ClosedAt = now
		}
	}
	return update
}

func CalcAndConvertShippingFeeLines(order *ghtkclient.OrderInfo) []*shippingsharemodel.ShippingFeeLine {
	var res []*shippingsharemodel.ShippingFeeLine
	insuranceFee := int(order.Insurance)
	fee := int(order.ShipMoney)
	shippingFeeMain := fee - insuranceFee

	// shipping fee
	res = append(res, &shippingsharemodel.ShippingFeeLine{
		ShippingFeeType:      shipping_fee_type.Main,
		Cost:                 shippingFeeMain,
		ExternalShippingCode: order.LabelID.String(),
	})
	// insurance fee
	if insuranceFee > 0 {
		res = append(res, &shippingsharemodel.ShippingFeeLine{
			ShippingFeeType:      shipping_fee_type.Insurance,
			Cost:                 insuranceFee,
			ExternalShippingCode: order.LabelID.String(),
		})
	}
	return res
}

func CalcDeliveryDuration(transport ghtkclient.TransportType, from, to *location.Province) time.Duration {
	switch {
	// Nội tỉnh
	case from.Code == to.Code:
		return 6 * time.Hour

	// HN, HCM, ĐN, Nội miền
	case from.Region == to.Region && from.Extra.Special:
		return 24 * time.Hour

	// HN, HCM, ĐN, Khác miền, Đặc biệt
	case from.Region != to.Region && from.Extra.Special && to.Extra.Special:
		if transport == ghtkclient.TransportFly {
			return 24 * time.Hour
		} else {
			return 96 * time.Hour
		}

	// HN, HCM, ĐH, Khác miền
	case from.Region != to.Region && from.Extra.Special:
		if transport == ghtkclient.TransportFly {
			return 48 * time.Hour
		} else {
			return 120 * time.Hour
		}

	// Tỉnh thành khác, nội miền, khác miền, Nhanh (không hỗ trợ gói Chuẩn)
	default:
		return 48 * time.Hour
	}
}

/*
func SyncOrders(ffms []*shipmodel.Fulfillment) ([]*shipmodel.Fulfillment, error) {
	rate := time.Second / 30
	burstLimit := 30
	ctx := bus.Ctx()
	tick := time.NewTicker(rate)
	defer tick.Stop()
	throttle := make(chan time.Time, burstLimit)
	go func() {
		for t := range tick.C {
			select {
			case throttle <- t:
			default:
			}
		}
	}()
	ch := make(chan error, burstLimit)
	ll.Info("Length GHTK SyncOrders", l.Int("len", len(ffms)))
	var _ffms []*shipmodel.Fulfillment
	count := 0
	for _, ffm := range ffms {
		<-throttle
		count++
		if count > burstLimit {
			time.Sleep(20 * time.Second)
			count = 0
		}
		go func(ffm *shipmodel.Fulfillment) (_err error) {
			defer func() {
				ch <- _err
			}()
			// get order info to update service fee
			ghtkCmd := &GetOrderCommand{
				ServiceID: ffm.ProviderServiceID,
				LabelID:   ffm.ShippingCode,
			}
			if ghtkErr := bus.Dispatch(ctx, ghtkCmd); ghtkErr != nil {
				ll.Error("GHTK get order error :: ", l.String("shipping_code", ffm.ShippingCode), l.Error(ghtkErr))
				return ghtkErr
			}
			updateFfm := CalcUpdateFulfillment(ffm, nil, &ghtkCmd.Result.Order)
			_ffms = append(_ffms, updateFfm)
			return nil
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
	ll.S.Infof("Sync fulfillments GHTK info success: %v/%v, errors %v/%v", successCount, len(ffms), errorCount, len(ffms))
	return _ffms, nil
}
*/
