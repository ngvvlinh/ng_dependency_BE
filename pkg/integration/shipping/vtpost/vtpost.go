package vtpost

import (
	"strconv"
	"time"

	"o.o/api/main/location"
	"o.o/api/top/types/etc/shipping_fee_type"
	"o.o/api/top/types/etc/status5"
	shipmodel "o.o/backend/com/main/shipping/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/etop/model"
	"o.o/backend/pkg/integration/shipping"
	vtpostclient "o.o/backend/pkg/integration/shipping/vtpost/client"
	"o.o/capi/dot"
	"o.o/common/jsonx"
	"o.o/common/l"
)

var ll = l.New()

const (
	VTPostCodePublic = 'D'
)

func CalcUpdateFulfillment(ffm *shipmodel.Fulfillment, orderMsg vtpostclient.CallbackOrderData) (*shipmodel.Fulfillment, error) {
	if !shipping.CanUpdateFulfillment(ffm) {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Can not update fulfillment (id = %v, shipping_code = %v)", ffm.ID, ffm.ShippingCode)
	}

	now := time.Now()
	data, _ := jsonx.Marshal(orderMsg)
	statusCode := orderMsg.OrderStatus
	vtpostStatus := vtpostclient.ToVTPostShippingState(statusCode)
	update := &shipmodel.Fulfillment{
		ID:                        ffm.ID,
		ExternalShippingUpdatedAt: now,
		ExternalShippingData:      data,
		ExternalShippingState:     orderMsg.StatusName,
		ExternalShippingStateCode: strconv.Itoa(statusCode),
		ExternalShippingStatus:    vtpostStatus.ToStatus5(),
		ShippingState:             vtpostStatus.ToModel(ffm.ShippingState),
		ShippingSubstate:          vtpostclient.ToShippingSubState(statusCode).Wrap(),
		EtopDiscount:              ffm.EtopDiscount,
		ShippingStatus:            vtpostStatus.ToShippingStatus5(ffm.ShippingState),
		ExternalShippingNote:      dot.String(orderMsg.Note),
		ExternalShippingSubState:  dot.String(vtpostclient.SubStateMap[statusCode]),
	}

	// Update price + weight
	if ffm.TotalWeight != orderMsg.ProductWeight {
		changeWeightNote := shipping.ChangeWeightNote(ffm.TotalWeight, orderMsg.ProductWeight)
		update.TotalWeight = orderMsg.ProductWeight
		update.AdminNote = ffm.AdminNote + "\n" + changeWeightNote
	}
	if ffm.ShippingFeeShop != orderMsg.MoneyTotal {
		// keep all shipping fee lines except shippingFeeMain
		mainFee := orderMsg.MoneyTotal
		for _, line := range ffm.ProviderShippingFeeLines {
			if line.ShippingFeeType == shipping_fee_type.Main {
				continue
			}
			mainFee = mainFee - line.Cost
		}
		if mainFee >= 0 {
			for _, line := range ffm.ProviderShippingFeeLines {
				if line.ShippingFeeType == shipping_fee_type.Main {
					line.Cost = mainFee
				}
			}
		}
		update.ProviderShippingFeeLines = ffm.ProviderShippingFeeLines
	}

	// Only update status5 if the current status is not ending status
	newStatus := vtpostStatus.ToStatus5()
	// UpdateInfo ClosedAt
	if newStatus == status5.N || newStatus == status5.NS || newStatus == status5.P {
		if ffm.ExternalShippingClosedAt.IsZero() {
			update.ClosedAt = now
		}
		if ffm.ClosedAt.IsZero() {
			update.ClosedAt = now
		}
	}
	return update, nil
}

func CalcDeliveryDuration(orderService vtpostclient.VTPostOrderServiceCode, fromProvince, toProvince *location.Province, fromDistrict, toDistrict *location.District) (duration time.Duration) {
	serviceName := orderService.Name()
	switch serviceName {
	case model.ShippingServiceNameFaster:
		duration = CalcDeliveryDurationFastService(fromProvince, toProvince, fromDistrict, toDistrict)
	case model.ShippingServiceNameStandard:
		duration = CalcDeliveryDurationStandardService(fromProvince, toProvince, fromDistrict, toDistrict)
	}
	return duration
}

func CalcDeliveryDurationStandardService(fromProvince, toProvince *location.Province, fromDistrict, toDistrict *location.District) (duration time.Duration) {
	switch {
	// Nội tỉnh
	case fromProvince.Code == toProvince.Code:
		switch toDistrict.UrbanType {
		case location.Urban:
			duration = 72 * time.Hour
		case location.Suburban1:
			duration = 96 * time.Hour
		default:
			duration = 120 * time.Hour
		}
		return duration

	// Nội miền
	case fromProvince.Region == toProvince.Region:
		switch toDistrict.UrbanType {
		case location.Urban:
			duration = 120 * time.Hour
		case location.Suburban1:
			duration = 144 * time.Hour
		default:
			duration = 168 * time.Hour
		}
		return duration

	// Khác miền
	case fromProvince.Region != toProvince.Region:
		switch toDistrict.UrbanType {
		case location.Urban:
			duration = 168 * time.Hour
		case location.Suburban1:
			duration = 192 * time.Hour
		default:
			duration = 216 * time.Hour
		}
		return duration
	// Tỉnh thành khác, nội miền, khác miền, Nhanh (không hỗ trợ gói Chuẩn)
	default:
		duration = 216 * time.Hour
		return duration
	}
}

// TODO: move back to location?
const (
	HCMProvinceCode       = "79"
	BinhDuongProvinceCode = "74"
	DongNaiProvinceCode   = "75"
	VungTauProvinceCode   = "77"
)

var groupProvinceCodes = []string{BinhDuongProvinceCode, DongNaiProvinceCode, VungTauProvinceCode}

func CalcDeliveryDurationFastService(fromProvince, toProvince *location.Province, fromDistrict, toDistrict *location.District) (duration time.Duration) {
	switch {
	// Nội tỉnh
	case fromProvince.Code == toProvince.Code:
		switch toDistrict.UrbanType {
		case location.Urban:
			duration = 24 * time.Hour
		case location.Suburban1:
			duration = 48 * time.Hour
		default:
			duration = 72 * time.Hour
		}
		return duration

	// HCM <=> Binh Duong, Dong Nai, Ba Ria Vung Tau
	case fromProvince.Code == HCMProvinceCode && cm.StringsContain(groupProvinceCodes, toProvince.Code) ||
		cm.StringsContain(groupProvinceCodes, fromProvince.Code) && toProvince.Code == HCMProvinceCode:
		switch toDistrict.UrbanType {
		case location.Urban:
			duration = 48 * time.Hour
		case location.Suburban1:
			duration = 48 * time.Hour
		default:
			duration = 72 * time.Hour
		}
		return duration

	// Nội miền
	case fromProvince.Region == toProvince.Region:
		switch toDistrict.UrbanType {
		case location.Urban:
			duration = 48 * time.Hour
		case location.Suburban1:
			duration = 48 * time.Hour
		default:
			duration = 120 * time.Hour
		}
		return duration

	// HCM <=> HN; DN <=> HCM, HN
	case fromProvince.Region != toProvince.Region && fromProvince.Special && toProvince.Special:
		switch toDistrict.UrbanType {
		case location.Urban:
			duration = 48 * time.Hour
		case location.Suburban1:
			duration = 48 * time.Hour
		default:
			duration = 72 * time.Hour
		}
		return duration

	// Khác miền
	case fromProvince.Region != toProvince.Region:
		switch toDistrict.UrbanType {
		case location.Urban:
			duration = 48 * time.Hour
		case location.Suburban1:
			duration = 72 * time.Hour
		default:
			duration = 120 * time.Hour
		}
		return duration
	// Tỉnh thành khác, nội miền, khác miền, Nhanh (không hỗ trợ gói Chuẩn)
	default:
		duration = 120 * time.Hour
		return duration
	}
}
