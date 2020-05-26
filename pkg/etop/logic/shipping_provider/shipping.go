package shipping_provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"o.o/api/main/location"
	"o.o/api/top/int/types"
	pbsp "o.o/api/top/types/etc/shipping_provider"
	shippingprovider "o.o/api/top/types/etc/shipping_provider"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/etop/model"
	"o.o/capi/dot"
	"o.o/common/l"
)

var ll = l.New()

var blockCarriers = map[shippingprovider.ShippingProvider]*struct {
	DateFrom time.Time
	DateTo   time.Time
}{
	shippingprovider.VTPost: {
		DateFrom: time.Date(2020, 1, 17, 0, 0, 0, 0, time.Local),
		DateTo:   time.Date(2020, 2, 9, 0, 0, 0, 0, time.Local),
	},
}

func (ctrl *CarrierManager) GetExternalShippingServices(ctx context.Context, accountID dot.ID, q *types.GetExternalShippingServicesRequest) ([]*model.AvailableShippingService, error) {
	fromQuery := &location.FindOrGetLocationQuery{
		ProvinceCode: q.FromProvinceCode,
		DistrictCode: q.FromDistrictCode,
		Province:     q.FromProvince,
		District:     q.FromDistrict,
	}
	toQuery := &location.FindOrGetLocationQuery{
		ProvinceCode: q.ToProvinceCode,
		DistrictCode: q.ToDistrictCode,
		Province:     q.ToProvince,
		District:     q.ToDistrict,
	}
	if err := ctrl.location.Dispatch(ctx, fromQuery); err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "địa chỉ gửi không hợp lệ: %v", err)
	}
	if err := ctrl.location.Dispatch(ctx, toQuery); err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "địa chỉ nhận không hợp lệ: %v", err)
	}

	fromDistrict, fromProvince := fromQuery.Result.District, fromQuery.Result.Province
	toDistrict, toProvince := toQuery.Result.District, toQuery.Result.Province
	if fromDistrict == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "địa chỉ gửi không hợp lệ")
	}
	if toDistrict == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "địa chỉ nhận không hợp lệ")
	}

	var res []*model.AvailableShippingService
	weight := q.Weight
	if q.GrossWeight != 0 {
		weight = q.GrossWeight
	}
	length := q.Length
	width := q.Width
	height := q.Height
	chargeableWeight := q.ChargeableWeight
	calculatedChargeableWeight := model.CalcChargeableWeight(weight, length, width, height)
	if chargeableWeight == 0 {
		chargeableWeight = calculatedChargeableWeight

	} else if chargeableWeight < calculatedChargeableWeight {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Khối lượng tính phí không hợp lệ.").
			WithMetap("chargeable_weight", chargeableWeight).
			WithMetap("gross_weight", q.GrossWeight).
			WithMetap("volumetric_weight (= length*width*height / 5)", length*width*height/5).
			WithMetap("expected chargeable_weight (= MAX(gross_weight, volumetric_weight))", calculatedChargeableWeight)
	}

	value := q.Value
	if q.BasketValue != 0 {
		value = q.BasketValue
	}
	includeInsurance := q.IncludeInsurance.Apply(false)

	totalCODAmount := q.TotalCodAmount
	if q.CodAmount != 0 {
		totalCODAmount = q.CodAmount
	}

	args := GetShippingServicesArgs{
		AccountID:        accountID,
		FromDistrictCode: fromDistrict.Code,
		ToDistrictCode:   toDistrict.Code,
		ChargeableWeight: chargeableWeight,
		Length:           length,
		Width:            width,
		Height:           height,
		IncludeInsurance: includeInsurance,
		BasketValue:      value,
		CODAmount:        totalCODAmount,
	}

	switch q.Provider {
	case pbsp.All, pbsp.Unknown:
		ch := make(chan []*model.AvailableShippingService, 2)
		go func() {
			defer catchAndRecover()

			var services []*model.AvailableShippingService
			var err error
			defer func() { sendServices(ch, services, err) }()
			services, err = ctrl.GetShippingProviderDriver(shippingprovider.GHN).GetAllShippingServices(ctx, args)
		}()
		go func() {
			defer catchAndRecover()

			var services []*model.AvailableShippingService
			var err error
			defer func() { sendServices(ch, services, err) }()
			services, err = ctrl.GetShippingProviderDriver(shippingprovider.GHTK).GetAllShippingServices(ctx, args)
		}()
		// go func() {
		// 	var services []*model.AvailableShippingService
		//
		// 	if err := checkBlockCarrier(shippingprovider.VTPost); err != nil {
		// 		sendServices(ch, nil, nil)
		// 		return
		// 	}
		//
		// 	var err error
		// 	defer func() { sendServices(ch, services, err) }()
		// 	services, err = ctrl.VTPost.GetAllShippingServices(ctx, args)
		// }()
		for i := 0; i < 2; i++ {
			res = append(res, <-ch...)
		}

	default:
		driver := ctrl.GetShippingProviderDriver(q.Provider)
		if driver == nil {
			return nil, cm.Error(cm.InvalidArgument, "Invalid carrier", nil)
		}

		services, err := driver.GetAllShippingServices(ctx, args)
		if err != nil {
			return nil, err
		}
		res = append(res, services...)
	}

	if len(res) == 0 {
		return nil, cm.Errorf(cm.ExternalServiceError, nil,
			"Tuyến giao hàng từ địa chỉ %v, %v đến địa chỉ %v, %v không được hỗ trợ bởi đơn vị vận chuyển nào",
			fromDistrict.Name, fromProvince.Name, toDistrict.Name, toProvince.Name).
			Log("District", l.String("from_code", fromDistrict.Code), l.String("to_code", toDistrict.Code))
	}

	res = CompactServices(res)
	return res, nil
}

func sendServices(ch chan<- []*model.AvailableShippingService, services []*model.AvailableShippingService, err error) {
	if err == nil {
		ch <- services
	} else {
		ch <- nil
	}
}

// CompactServices Loại bỏ các service không sử dụng
// Trường hợp:
// - Có gói TopShip: chỉ sử dụng gói TopShip
// - Mỗi NVC phải có 2 dịch vụ: Nhanh và Chuẩn, ưu tiên gói TopShip
// - Không có gói TopShip: Sử dụng gói của NVC như bình thường
func CompactServices(services []*model.AvailableShippingService) []*model.AvailableShippingService {
	var res []*model.AvailableShippingService
	carrierServicesIndex := make(map[string][]*model.AvailableShippingService)
	for _, s := range services {
		connectionID := dot.ID(0)
		if s.ConnectionInfo != nil {
			connectionID = s.ConnectionInfo.ID
		}
		key := fmt.Sprintf("%v_%v_%v", s.Provider.String(), s.Name, connectionID)
		carrierServicesIndex[key] = append(carrierServicesIndex[key], s)
	}
	for _, carrierServices := range carrierServicesIndex {
		var ss []*model.AvailableShippingService
		for _, s := range carrierServices {
			if s.Source == model.TypeShippingSourceEtop {
				ss = append(ss, s)
			}
		}
		if len(ss) > 0 {
			res = append(res, ss...)
		} else {
			res = append(res, carrierServices...)
		}
	}
	return res
}

func catchAndRecover() {
	e := recover()
	if e != nil {
		ll.Error("panic (recovered)", l.Object("error", e), l.Stack())
	}
}

func checkBlockCarrier(carrier shippingprovider.ShippingProvider) error {
	now := time.Now()
	blockInfo := blockCarriers[carrier]
	if blockInfo == nil {
		return nil
	}
	carrierName := strings.ToUpper(carrier.String())
	if now.After(blockInfo.DateFrom) && now.Before(blockInfo.DateTo) {
		return cm.Errorf(cm.FailedPrecondition, nil, "%v ngừng lấy hàng từ ngày %v đến ngày %v. Bạn không thể tạo đơn %v trong thời gian này!", carrierName, blockInfo.DateFrom.Format("02-01-2006"), blockInfo.DateTo.Add(-24*time.Hour).Format("02-01-2006"), carrierName)
	}
	return nil
}
