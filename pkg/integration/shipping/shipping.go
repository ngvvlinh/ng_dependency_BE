package shipping

import (
	"fmt"
	"time"

	mdlocation "etop.vn/api/main/location"
	"etop.vn/api/main/shipnow"
	shipnowtypes "etop.vn/api/main/shipnow/types"
	locationutil "etop.vn/backend/com/main/location/util"
	shipmodel "etop.vn/backend/com/main/shipping/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/etop/model"
)

func CalcPickTime(shippingProvider model.ShippingProvider, t time.Time) time.Time {
	// VTPOST: thời gian lấy hàng dự kiến tạo trước 16h
	// => lấy trong ngày (trước 18h), tạo sau 16h => lấy trước 18h ngày hôm sau
	h, m := t.Hour(), t.Minute()
	hm := cm.Clock(h, m)
	var res time.Time
	switch {
	case hm < cm.Clock(10, 0):
		if shippingProvider == model.TypeVTPost {
			res = t.Add(-hm).Add(cm.Clock(18, 0))
		} else {
			res = t.Add(-hm).Add(cm.Clock(12, 0))
		}
	case hm < cm.Clock(16, 0):
		res = t.Add(-hm).Add(cm.Clock(18, 0))
	default:
		if shippingProvider == model.TypeVTPost {
			res = t.Add(-hm).Add(cm.Clock(18, 0) + 24*time.Hour)
		} else {
			res = t.Add(-hm).Add(cm.Clock(12, 0) + 24*time.Hour)
		}
	}

	res = res.Truncate(time.Hour)
	return res
}

func CalcServicesTime(shippingProvider model.ShippingProvider, fromDistrict *mdlocation.District, toDistrict *mdlocation.District, services []*model.AvailableShippingService) []*model.AvailableShippingService {
	for _, service := range services {
		service = CalcServiceTime(shippingProvider, fromDistrict, toDistrict, service)
	}
	return services
}

func CalcServiceTime(shippingProvider model.ShippingProvider, fromDistrict *mdlocation.District, toDistrict *mdlocation.District, service *model.AvailableShippingService) *model.AvailableShippingService {
	// GHN, GHTK
	// Thời gian lấy hàng dự kiến:
	// Chỉ lấy vào chủ chủ nhật trường hợp nội thành HCM,HN.
	// Các TH khác, chuyển qua sáng thứ 2 (dời thời gian giao lui lại 1 ngày)
	//
	// VTPOST: Không lấy hàng CN
	//
	// ALL: Thời gian giao hàng dự kiến của NVC nếu sau 19h
	// => chuyển qua 12h ngày hôm sau
	// Chủ nhật chỉ giao nội thành HCM, HN
	// VTPOST: Chủ nhật ko giao
	pickAt := service.ExpectedPickAt
	deliveryAt := service.ExpectedDeliveryAt

	weekDayPickAt := int(pickAt.Weekday())
	if weekDayPickAt == 7 {
		if shippingProvider == model.TypeVTPost {
			pickAt = pickAt.Add(time.Hour * 24)
			service.ExpectedPickAt = pickAt
			deliveryAt = deliveryAt.Add(time.Hour * 24)
			service.ExpectedDeliveryAt = deliveryAt
		} else if !locationutil.CheckUrbanHCMHN(fromDistrict) {
			pickAt = pickAt.Add(time.Hour * 24)
			service.ExpectedPickAt = time.Date(pickAt.Year(), pickAt.Month(), pickAt.Day(), 12, 0, 0, 0, pickAt.Location())
			deliveryAt = deliveryAt.Add(time.Hour * 24)
			service.ExpectedDeliveryAt = deliveryAt
		}
	}

	if deliveryAt.Hour() >= 19 {
		deliveryAt = deliveryAt.Add(time.Hour * 24)
		deliveryAt = time.Date(deliveryAt.Year(), deliveryAt.Month(), deliveryAt.Day(), 12, 0, 0, 0, deliveryAt.Location())
		service.ExpectedDeliveryAt = deliveryAt
	}
	weekDayDeliveryAt := int(deliveryAt.Weekday())
	if weekDayDeliveryAt == 7 {
		if shippingProvider == model.TypeVTPost || !locationutil.CheckUrbanHCMHN(toDistrict) {
			deliveryAt = deliveryAt.Add(time.Hour * 24)
			service.ExpectedDeliveryAt = deliveryAt
		}
	}
	return service
}

func CalcDeliveryTime(shippingProvider model.ShippingProvider, toDistrict *mdlocation.District, deliveryAt time.Time) time.Time {
	// Thời gian giao hàng dự kiến của NVC nếu sau 19h => chuyển qua 12h ngày hôm sau
	// Chủ nhật chỉ giao nội thành HCM, HN
	// VTPOST: Chủ nhật ko giao
	if deliveryAt.Hour() >= 19 {
		deliveryAt = deliveryAt.Add(time.Hour * 24)
		deliveryAt = time.Date(deliveryAt.Year(), deliveryAt.Month(), deliveryAt.Day(), 12, 0, 0, 0, deliveryAt.Location())
	}
	weekDayDeliveryAt := int(deliveryAt.Weekday())
	if weekDayDeliveryAt == 7 {
		if shippingProvider == model.TypeVTPost || !locationutil.CheckUrbanHCMHN(toDistrict) {
			deliveryAt = deliveryAt.Add(time.Hour * 24)
		}
	}
	return deliveryAt
}

func CalcOtherTimeBaseOnState(update *shipmodel.Fulfillment, oldFfm *shipmodel.Fulfillment, t time.Time) *shipmodel.Fulfillment {
	state := update.ShippingState
	switch state {
	case model.StateCreated:
		if oldFfm.ShippingCreatedAt.IsZero() {
			update.ShippingCreatedAt = t
		}
	case model.StatePicking:
		if oldFfm.ShippingPickingAt.IsZero() {
			update.ShippingPickingAt = t
		}
	case model.StateHolding:
		if oldFfm.ShippingHoldingAt.IsZero() {
			update.ShippingHoldingAt = t
		}
	case model.StateDelivering:
		if oldFfm.ShippingDeliveringAt.IsZero() {
			update.ShippingDeliveringAt = t
		}
	case model.StateDelivered:
		if oldFfm.ExternalShippingDeliveredAt.IsZero() {
			update.ExternalShippingDeliveredAt = t
		}
		if oldFfm.ShippingDeliveredAt.IsZero() {
			update.ShippingDeliveredAt = t
		}
	case model.StateReturning:
		if oldFfm.ShippingReturningAt.IsZero() {
			update.ShippingReturningAt = t
		}
	case model.StateReturned:
		if oldFfm.ExternalShippingReturnedAt.IsZero() {
			update.ExternalShippingReturnedAt = t
		}
		if oldFfm.ShippingReturnedAt.IsZero() {
			update.ShippingReturnedAt = t
		}
	case model.StateCancelled:
		if oldFfm.ExternalShippingCancelledAt.IsZero() {
			update.ExternalShippingCancelledAt = t
		}
		if oldFfm.ShippingCancelledAt.IsZero() {
			update.ShippingCancelledAt = t
		}
	default:
	}
	return update
}

func CanUpdateFulfillmentFromWebhook(ffm *shipmodel.Fulfillment) bool {
	return ffm.Status == model.S5Zero ||
		ffm.Status == model.S5SuperPos ||

		// returning has status -2 (NS) and we allow updating it via webhook
		ffm.ShippingState == model.StateReturning
}

type ShipnowTimestamp struct {
	ShippingCreatedAt    time.Time
	ShippingPickingAt    time.Time
	ShippingDeliveringAt time.Time
	ShippingDeliveredAt  time.Time
	ShippingCancelledAt  time.Time
}

func CalcShipnowTimeBaseOnState(ffm *shipnow.ShipnowFulfillment, state shipnowtypes.State, t time.Time) (res ShipnowTimestamp) {
	switch state {
	case shipnowtypes.StateCreated:
		if ffm.ShippingCreatedAt == nil || ffm.ShippingCreatedAt.ToTime().IsZero() {
			res.ShippingCreatedAt = t
		}
	case shipnowtypes.StatePicking:
		if ffm.ShippingPickingAt == nil || ffm.ShippingPickingAt.ToTime().IsZero() {
			res.ShippingPickingAt = t
		}
	case shipnowtypes.StateDelivering:
		if ffm.ShippingDeliveringAt == nil || ffm.ShippingDeliveringAt.ToTime().IsZero() {
			res.ShippingDeliveringAt = t
		}
	case shipnowtypes.StateDelivered:
		if ffm.ShippingDeliveredAt == nil || ffm.ShippingDeliveredAt.ToTime().IsZero() {
			res.ShippingDeliveredAt = t
		}
	case shipnowtypes.StateCancelled:
		if ffm.ShippingCancelledAt == nil || ffm.ShippingCancelledAt.ToTime().IsZero() {
			res.ShippingCancelledAt = t
		}
	default:
	}
	return
}

func ChangeWeightNote(oldWeight, newWeight int) string {
	return fmt.Sprintf("Khối lượng thay đổi từ %v thành %v", oldWeight, newWeight)
}
