package shipping

import (
	"context"
	"fmt"
	"time"

	"o.o/api/main/identity"
	"o.o/api/main/location"
	"o.o/api/main/shipnow"
	shippingcore "o.o/api/main/shipping"
	"o.o/api/meta"
	"o.o/api/top/types/etc/connection_type"
	"o.o/api/top/types/etc/shipnow_state"
	"o.o/api/top/types/etc/shipping"
	"o.o/api/top/types/etc/shipping_provider"
	"o.o/api/top/types/etc/status5"
	locationutil "o.o/backend/com/main/location/util"
	shippingconvert "o.o/backend/com/main/shipping/convert"
	shipmodel "o.o/backend/com/main/shipping/model"
	shippingsharemodel "o.o/backend/com/main/shipping/sharemodel"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/capi/dot"
	"o.o/common/l"
)

var ll = l.New().WithChannel(meta.ChannelShipmentCarrier)

func CalcPickTime(shippingProvider shipping_provider.ShippingProvider, t time.Time) time.Time {
	// VTPOST: thời gian lấy hàng dự kiến tạo trước 16h
	// => lấy trong ngày (trước 18h), tạo sau 16h => lấy trước 18h ngày hôm sau
	h, m := t.Hour(), t.Minute()
	hm := cm.Clock(h, m)
	var res time.Time
	switch {
	case hm < cm.Clock(10, 0):
		if shippingProvider == shipping_provider.VTPost {
			res = t.Add(-hm).Add(cm.Clock(18, 0))
		} else {
			res = t.Add(-hm).Add(cm.Clock(12, 0))
		}
	case hm < cm.Clock(16, 0):
		res = t.Add(-hm).Add(cm.Clock(18, 0))
	default:
		if shippingProvider == shipping_provider.VTPost {
			res = t.Add(-hm).Add(cm.Clock(18, 0) + 24*time.Hour)
		} else {
			res = t.Add(-hm).Add(cm.Clock(12, 0) + 24*time.Hour)
		}
	}

	res = res.Truncate(time.Hour)
	return res
}

func CalcServicesTime(shippingProvider shipping_provider.ShippingProvider, fromDistrict *location.District, toDistrict *location.District, services []*shippingsharemodel.AvailableShippingService) []*shippingsharemodel.AvailableShippingService {
	for _, service := range services {
		_ = CalcServiceTime(shippingProvider, fromDistrict, toDistrict, service)
	}
	return services
}

func CalcServiceTime(shippingProvider shipping_provider.ShippingProvider, fromDistrict *location.District, toDistrict *location.District, service *shippingsharemodel.AvailableShippingService) *shippingsharemodel.AvailableShippingService {
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
		if shippingProvider == shipping_provider.VTPost {
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
		if shippingProvider == shipping_provider.VTPost || !locationutil.CheckUrbanHCMHN(toDistrict) {
			deliveryAt = deliveryAt.Add(time.Hour * 24)
			service.ExpectedDeliveryAt = deliveryAt
		}
	}
	return service
}

func CalcDeliveryTime(shippingProvider shipping_provider.ShippingProvider, toDistrict *location.District, deliveryAt time.Time) time.Time {
	// Thời gian giao hàng dự kiến của NVC nếu sau 19h => chuyển qua 12h ngày hôm sau
	// Chủ nhật chỉ giao nội thành HCM, HN
	// VTPOST: Chủ nhật ko giao
	if deliveryAt.Hour() >= 19 {
		deliveryAt = deliveryAt.Add(time.Hour * 24)
		deliveryAt = time.Date(deliveryAt.Year(), deliveryAt.Month(), deliveryAt.Day(), 12, 0, 0, 0, deliveryAt.Location())
	}
	weekDayDeliveryAt := int(deliveryAt.Weekday())
	if weekDayDeliveryAt == 7 {
		if shippingProvider == shipping_provider.VTPost || !locationutil.CheckUrbanHCMHN(toDistrict) {
			deliveryAt = deliveryAt.Add(time.Hour * 24)
		}
	}
	return deliveryAt
}

func CalcOtherTimeBaseOnState(update *shipmodel.Fulfillment, oldFfm *shipmodel.Fulfillment, t time.Time) *shipmodel.Fulfillment {
	state := update.ShippingState
	switch state {
	case shipping.Created:
		if oldFfm.ShippingCreatedAt.IsZero() {
			update.ShippingCreatedAt = t
		}
	case shipping.Picking:
		if oldFfm.ShippingPickingAt.IsZero() {
			update.ShippingPickingAt = t
		}
	case shipping.Holding:
		if oldFfm.ShippingHoldingAt.IsZero() {
			update.ShippingHoldingAt = t
		}
	case shipping.Delivering:
		if oldFfm.ShippingDeliveringAt.IsZero() {
			update.ShippingDeliveringAt = t
		}
	case shipping.Delivered:
		if oldFfm.ExternalShippingDeliveredAt.IsZero() {
			update.ExternalShippingDeliveredAt = t
		}
		if oldFfm.ShippingDeliveredAt.IsZero() {
			update.ShippingDeliveredAt = t
		}
	case shipping.Returning:
		if oldFfm.ShippingReturningAt.IsZero() {
			update.ShippingReturningAt = t
		}
	case shipping.Returned:
		if oldFfm.ExternalShippingReturnedAt.IsZero() {
			update.ExternalShippingReturnedAt = t
		}
		if oldFfm.ShippingReturnedAt.IsZero() {
			update.ShippingReturnedAt = t
		}
	case shipping.Cancelled:
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

func CanUpdateFulfillment(ffm *shipmodel.Fulfillment) bool {
	return ffm.Status == status5.Z ||
		ffm.Status == status5.S ||

		// returning has status -2 (NS) and we allow updating it via webhook
		ffm.ShippingState == shipping.Returning
}

func CanUpdateFulfillmentFeelines(ffm *shipmodel.Fulfillment) bool {
	if ffm.MoneyTransactionShippingExternalID != 0 ||
		ffm.MoneyTransactionID != 0 {
		return false
	}
	return true
}

type ShipnowTimestamp struct {
	ShippingCreatedAt    time.Time
	ShippingPickingAt    time.Time
	ShippingDeliveringAt time.Time
	ShippingDeliveredAt  time.Time
	ShippingCancelledAt  time.Time
}

func CalcShipnowTimeBaseOnState(ffm *shipnow.ShipnowFulfillment, state shipnow_state.State, t time.Time) (res ShipnowTimestamp) {
	switch state {
	case shipnow_state.StateCreated:
		if ffm.ShippingCreatedAt.IsZero() {
			res.ShippingCreatedAt = t
		}
	case shipnow_state.StatePicking:
		if ffm.ShippingPickingAt.IsZero() {
			res.ShippingPickingAt = t
		}
	case shipnow_state.StateDelivering:
		if ffm.ShippingDeliveringAt.IsZero() {
			res.ShippingDeliveringAt = t
		}
	case shipnow_state.StateDelivered:
		if ffm.ShippingDeliveredAt.IsZero() {
			res.ShippingDeliveredAt = t
		}
	case shipnow_state.StateCancelled:
		if ffm.ShippingCancelledAt.IsZero() {
			res.ShippingCancelledAt = t
		}
	default:
	}
	return
}

func ChangeWeightNote(oldWeight, newWeight int) string {
	return fmt.Sprintf("Khối lượng thay đổi từ %v thành %v", oldWeight, newWeight)
}

func WebhookWlWrapContext(ctx context.Context, shopID dot.ID, identityQS identity.QueryBus) (context.Context, error) {
	// Get WLPartnerID to wrap context from shop (GetShopByID)
	// TODO: ffm contains WLPartnerID
	queryShop := &identity.GetShopByIDQuery{
		ID: shopID,
	}
	if err := identityQS.Dispatch(ctx, queryShop); err != nil {
		return nil, err
	}
	ctx = wl.WrapContextByPartnerID(ctx, queryShop.Result.WLPartnerID)
	return ctx, nil
}

type UpdateShippingFeeLinesArgs struct {
	FfmID            dot.ID
	Weight           int
	State            shipping.State
	ProviderFeeLines []*shippingsharemodel.ShippingFeeLine
}

func UpdateShippingFeeLines(ctx context.Context, shippingAggr shippingcore.CommandBus, args *UpdateShippingFeeLinesArgs) error {
	providerFeeLinesCore := shippingconvert.Convert_sharemodel_ShippingFeeLines_shipping_ShippingFeeLines(args.ProviderFeeLines)
	cmd := &shippingcore.UpdateFulfillmentShippingFeesFromWebhookCommand{
		FulfillmentID:    args.FfmID,
		NewWeight:        args.Weight,
		NewState:         args.State,
		ProviderFeeLines: providerFeeLinesCore,
	}
	return shippingAggr.Dispatch(ctx, cmd)
}

type UpdateFfmCODAmountArgs struct {
	NewCODAmount int
	Ffm          *shipmodel.Fulfillment
	CarrierName  string
}

// ValidateAndUpdateFulfillmentCOD
//
// Cập nhật COD Amount (chỉ sử dụng khi nhận webhook)
// Nếu phát sinh lỗi, bắn ra telegram để thông báo, không trả về lỗi
func ValidateAndUpdateFulfillmentCOD(ctx context.Context, shippingAggr shippingcore.CommandBus, args *UpdateFfmCODAmountArgs) {
	newCODAmount := args.NewCODAmount
	ffm := args.Ffm
	if newCODAmount != ffm.TotalCODAmount {
		switch ffm.ConnectionMethod {
		case connection_type.ConnectionMethodDirect:
			updateFulfillmentShippingFeesCmd := &shippingcore.UpdateFulfillmentShippingFeesCommand{
				FulfillmentID:  ffm.ID,
				TotalCODAmount: dot.Int(newCODAmount),
			}
			if err := shippingAggr.Dispatch(ctx, updateFulfillmentShippingFeesCmd); err != nil {
				ll.SendMessage("–––\n👹 %v: đơn %v cập nhật thay đổi COD thất bại. 👹 \n Lỗi: %v \n––", args.CarrierName, ffm.ShippingCode, err.Error())
				return
			}
		default:
			str := "–––\n👹 %v: đơn %v có thay đổi COD. Không thể cập nhật, vui lòng kiểm tra lại. 👹 \n- COD hiện tại: %v \n- COD mới: %v\n–––"
			ll.SendMessage(fmt.Sprintf(str, args.CarrierName, ffm.ShippingCode, ffm.TotalCODAmount, newCODAmount))
		}
	}
	return
}
