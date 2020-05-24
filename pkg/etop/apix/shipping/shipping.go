package shipping

import (
	"context"

	"o.o/api/main/connectioning"
	"o.o/api/main/location"
	"o.o/api/main/shipping"
	exttypes "o.o/api/top/external/types"
	"o.o/api/top/int/types"
	servicelocation "o.o/backend/com/main/location"
	locationlist "o.o/backend/com/main/location/list"
	ordersqlstore "o.o/backend/com/main/ordering/sqlstore"
	shippingcarrier "o.o/backend/com/main/shipping/carrier"
	shipsqlstore "o.o/backend/com/main/shipping/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/etop/apix/convertpb"
	"o.o/backend/pkg/etop/logic/shipping_provider"
	"o.o/backend/pkg/etop/model"
	"o.o/capi/dot"
	"o.o/common/l"
)

var (
	shippingCtrl     *shipping_provider.ProviderManager
	locationBus      = servicelocation.QueryMessageBus(servicelocation.New(nil))
	locationList     = buildLocationList()
	orderStore       ordersqlstore.OrderStoreFactory
	fulfillmentStore shipsqlstore.FulfillmentStoreFactory
	shipmentManager  *shippingcarrier.ShipmentManager
	shippingAggr     shipping.CommandBus
	shippingQuery    shipping.QueryBus
	connectionQS     connectioning.QueryBus
)

func Init(_shippingCtrl *shipping_provider.ProviderManager, _orderStore ordersqlstore.OrderStoreFactory, ffmStore shipsqlstore.FulfillmentStoreFactory, shipmentM *shippingcarrier.ShipmentManager, shippingA shipping.CommandBus, shippingQ shipping.QueryBus, connectionQueryService connectioning.QueryBus) {
	shippingCtrl = _shippingCtrl
	orderStore = _orderStore
	fulfillmentStore = ffmStore
	shipmentManager = shipmentM
	shippingAggr = shippingA
	shippingQuery = shippingQ
	connectionQS = connectionQueryService
}

// TODO: should not import location/list
func buildLocationList() *exttypes.LocationResponse {
	provinces := make([]exttypes.Province, len(locationlist.Provinces))
	for i, p := range locationlist.Provinces {
		districtsQuery := &location.GetAllLocationsQuery{ProvinceCode: p.Code}
		if err := locationBus.Dispatch(context.Background(), districtsQuery); err != nil {
			ll.Panic("unexpected", l.Error(err))
		}
		ds := districtsQuery.Result.Districts
		districts := make([]exttypes.District, len(ds))

		for i, d := range ds {
			wardsQuery := &location.GetAllLocationsQuery{DistrictCode: d.Code}
			if err := locationBus.Dispatch(context.Background(), wardsQuery); err != nil {
				ll.Panic("unexpected", l.Error(err))
			}
			ws := wardsQuery.Result.Wards
			wards := make([]exttypes.Ward, len(ws))
			for i, w := range ws {
				wards[i] = exttypes.Ward{Name: w.Name}
			}
			districts[i] = exttypes.District{
				Name:  d.Name,
				Wards: wards,
			}
		}

		provinces[i] = exttypes.Province{
			Name:      p.Name,
			Districts: districts,
		}
	}
	return &exttypes.LocationResponse{
		Provinces: provinces,
	}
}

func GetLocationList(ctx context.Context) (*exttypes.LocationResponse, error) {
	return locationList, nil
}

func GetShippingServices(ctx context.Context, accountID dot.ID, r *exttypes.GetShippingServicesRequest) (*exttypes.GetShippingServicesResponse, error) {
	if r.PickupAddress == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Cần cung cấp địa chỉ lấy hàng")
	}
	if r.ShippingAddress == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Cần cung cấp địa chỉ giao hàng")
	}
	req := &types.GetShippingServicesRequest{
		ConnectionIDs:    r.ConnectionIDs,
		FromDistrictCode: "",
		FromProvinceCode: "",
		ToDistrictCode:   "",
		ToProvinceCode:   "",
		FromProvince:     r.PickupAddress.Province,
		FromDistrict:     r.PickupAddress.District,
		ToProvince:       r.ShippingAddress.Province,
		ToDistrict:       r.ShippingAddress.District,
		GrossWeight:      r.GrossWeight,
		ChargeableWeight: r.ChargeableWeight,
		Length:           r.Length,
		Width:            r.Width,
		Height:           r.Height,
		TotalCodAmount:   r.CodAmount,
		BasketValue:      r.BasketValue,
		IncludeInsurance: r.IncludeInsurance,
	}
	args, err := shipmentManager.PrepareDataGetShippingServices(ctx, req)
	if err != nil {
		return nil, err
	}
	resp, err := shipmentManager.GetShippingServices(ctx, accountID, args)
	if err != nil {
		return nil, err
	}
	if err := buildCodeForShippingServices(ctx, resp); err != nil {
		return nil, err
	}
	return &exttypes.GetShippingServicesResponse{
		Services: convertpb.PbShippingServices(resp),
	}, nil
}

func buildCodeForShippingServices(ctx context.Context, services []*model.AvailableShippingService) error {
	// add connection code to service code to identify which connects
	// code format: XXXXYYYYYYYY (12 characters)
	for _, s := range services {
		if s.ConnectionInfo == nil {
			continue
		}
		connection, err := shipmentManager.GetConnectionByID(ctx, s.ConnectionInfo.ID)
		if err != nil {
			return err
		}
		s.ProviderServiceID = connection.Code + s.ProviderServiceID
	}
	return nil
}

func parseServiceCode(ctx context.Context, serviceCode string) (conn *connectioning.Connection, code string, _ error) {
	if len(serviceCode) <= 8 {
		return nil, "", cm.Errorf(cm.InvalidArgument, nil, "Shipping service code is invalid")
	}
	connCode, code := serviceCode[:4], serviceCode[4:]
	conn, err := shipmentManager.GetConnectionByCode(ctx, connCode)
	if err != nil {
		return nil, "", err
	}
	return conn, code, nil
}
