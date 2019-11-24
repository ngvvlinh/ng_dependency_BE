package shipping

import (
	"context"

	"etop.vn/api/main/location"
	pborder "etop.vn/api/pb/etop/order"
	pbexternal "etop.vn/api/pb/external"
	servicelocation "etop.vn/backend/com/main/location"
	locationlist "etop.vn/backend/com/main/location/list"
	ordersqlstore "etop.vn/backend/com/main/ordering/sqlstore"
	shipsqlstore "etop.vn/backend/com/main/shipping/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/etop/apix/convertpb"
	"etop.vn/backend/pkg/etop/logic/shipping_provider"
	"etop.vn/capi/dot"
	"etop.vn/common/l"
)

var shippingCtrl *shipping_provider.ProviderManager
var locationBus = servicelocation.New().MessageBus()
var locationList = buildLocationList()
var orderStore ordersqlstore.OrderStoreFactory
var fulfillmentStore shipsqlstore.FulfillmentStoreFactory

func Init(_shippingCtrl *shipping_provider.ProviderManager, _orderStore ordersqlstore.OrderStoreFactory, ffmStore shipsqlstore.FulfillmentStoreFactory) {
	shippingCtrl = _shippingCtrl
	orderStore = _orderStore
	fulfillmentStore = ffmStore
}

// TODO: should not import location/list
func buildLocationList() *pbexternal.LocationResponse {
	provinces := make([]pbexternal.Province, len(locationlist.Provinces))
	for i, p := range locationlist.Provinces {
		districtsQuery := &location.GetAllLocationsQuery{ProvinceCode: p.Code}
		if err := locationBus.Dispatch(context.Background(), districtsQuery); err != nil {
			ll.Panic("unexpected", l.Error(err))
		}
		ds := districtsQuery.Result.Districts
		districts := make([]pbexternal.District, len(ds))

		for i, d := range ds {
			wardsQuery := &location.GetAllLocationsQuery{DistrictCode: d.Code}
			if err := locationBus.Dispatch(context.Background(), wardsQuery); err != nil {
				ll.Panic("unexpected", l.Error(err))
			}
			ws := wardsQuery.Result.Wards
			wards := make([]pbexternal.Ward, len(ws))
			for i, w := range ws {
				wards[i] = pbexternal.Ward{Name: w.Name}
			}
			districts[i] = pbexternal.District{
				Name:  d.Name,
				Wards: wards,
			}
		}

		provinces[i] = pbexternal.Province{
			Name:      p.Name,
			Districts: districts,
		}
	}
	return &pbexternal.LocationResponse{
		Provinces: provinces,
	}
}

func GetLocationList(ctx context.Context) (*pbexternal.LocationResponse, error) {
	return locationList, nil
}

func GetShippingServices(ctx context.Context, accountID dot.ID, r *pbexternal.GetShippingServicesRequest) (*pbexternal.GetShippingServicesResponse, error) {
	if r.PickupAddress == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Cần cung cấp địa chỉ lấy hàng")
	}
	if r.ShippingAddress == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Cần cung cấp địa chỉ giao hàng")
	}

	req := &pborder.GetExternalShippingServicesRequest{
		Provider:         0,
		Carrier:          0,
		FromDistrictCode: "",
		FromProvinceCode: "",
		ToDistrictCode:   "",
		ToProvinceCode:   "",
		FromProvince:     r.PickupAddress.Province,
		FromDistrict:     r.PickupAddress.District,
		ToProvince:       r.ShippingAddress.Province,
		ToDistrict:       r.ShippingAddress.District,
		Weight:           0,
		GrossWeight:      r.GrossWeight,
		ChargeableWeight: r.ChargeableWeight,
		Length:           r.Length,
		Width:            r.Width,
		Height:           r.Height,
		Value:            r.BasketValue,
		TotalCodAmount:   r.CodAmount,
		CodAmount:        r.CodAmount,
		BasketValue:      r.BasketValue,
		IncludeInsurance: r.IncludeInsurance,
	}
	services, err := shippingCtrl.GetExternalShippingServices(ctx, accountID, req)
	if err != nil {
		return nil, err
	}
	return &pbexternal.GetShippingServicesResponse{
		Services: convertpb.PbShippingServices(services),
	}, nil
}
