package shipping_provider

import (
	"fmt"
	"strings"
	"time"

	shippingprovider "o.o/api/top/types/etc/shipping_provider"
	shippingsharemodel "o.o/backend/com/main/shipping/sharemodel"
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

func sendServices(ch chan<- []*shippingsharemodel.AvailableShippingService, services []*shippingsharemodel.AvailableShippingService, err error) {
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
func CompactServices(services []*shippingsharemodel.AvailableShippingService) []*shippingsharemodel.AvailableShippingService {
	var res []*shippingsharemodel.AvailableShippingService
	carrierServicesIndex := make(map[string][]*shippingsharemodel.AvailableShippingService)
	for _, s := range services {
		connectionID := dot.ID(0)
		if s.ConnectionInfo != nil {
			connectionID = s.ConnectionInfo.ID
		}
		key := fmt.Sprintf("%v_%v_%v", s.Provider.String(), s.Name, connectionID)
		carrierServicesIndex[key] = append(carrierServicesIndex[key], s)
	}
	for _, carrierServices := range carrierServicesIndex {
		var ss []*shippingsharemodel.AvailableShippingService
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
