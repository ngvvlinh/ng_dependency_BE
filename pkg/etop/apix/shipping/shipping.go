package shipping

import (
	"context"

	"o.o/api/main/connectioning"
	"o.o/api/main/location"
	exttypes "o.o/api/top/external/types"
	"o.o/api/top/int/types"
	locationlist "o.o/backend/com/main/location/list"
	shippingsharemodel "o.o/backend/com/main/shipping/sharemodel"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/etop/apix/convertpb"
	"o.o/capi/dot"
	"o.o/common/l"
)

// TODO: should not import location/list
func buildLocationList(locationBus location.QueryBus) *exttypes.LocationResponse {
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

func (s *Shipping) GetLocationList(ctx context.Context) (*exttypes.LocationResponse, error) {
	return locationList, nil
}

func (s *Shipping) GetShippingServices(ctx context.Context, accountID dot.ID, r *exttypes.GetShippingServicesRequest) (*exttypes.GetShippingServicesResponse, error) {
	if r.PickupAddress == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Cần cung cấp địa chỉ lấy hàng")
	}
	if r.ShippingAddress == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Cần cung cấp địa chỉ giao hàng")
	}
	req := &types.GetShippingServicesRequest{
		ConnectionIDs:    r.ConnectionIDs,
		FromProvince:     r.PickupAddress.Province,
		FromDistrict:     r.PickupAddress.District,
		FromWard:         r.PickupAddress.Ward,
		ToProvince:       r.ShippingAddress.Province,
		ToDistrict:       r.ShippingAddress.District,
		ToWard:           r.ShippingAddress.Ward,
		GrossWeight:      r.GrossWeight,
		ChargeableWeight: r.ChargeableWeight,
		Length:           r.Length,
		Width:            r.Width,
		Height:           r.Height,
		TotalCodAmount:   r.CodAmount,
		BasketValue:      r.BasketValue,
		IncludeInsurance: r.IncludeInsurance,
	}
	args, err := s.ShipmentManager.PrepareDataGetShippingServices(ctx, req)
	if err != nil {
		return nil, err
	}
	if args.FromWardCode == "" || args.ToWardCode == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Địa chỉ gửi/giao hàng không được thiếu phường/xã")
	}
	args.AccountID = accountID
	resp, err := s.ShipmentManager.GetShippingServices(ctx, args)
	if err != nil {
		return nil, err
	}
	if err := s.buildCodeForShippingServices(ctx, resp); err != nil {
		return nil, err
	}
	return &exttypes.GetShippingServicesResponse{
		Services: convertpb.PbShippingServices(resp),
	}, nil
}

func (s *Shipping) buildCodeForShippingServices(ctx context.Context, services []*shippingsharemodel.AvailableShippingService) error {
	// add connection code to service code to identify which connects
	// code format: XXXXYYYYYYYY (12 characters)
	for _, srv := range services {
		if srv.ConnectionInfo == nil {
			continue
		}
		connection, err := s.ShipmentManager.ConnectionManager.GetConnectionByID(ctx, srv.ConnectionInfo.ID)
		if err != nil {
			return err
		}
		srv.ProviderServiceID = connection.Code + srv.ProviderServiceID
	}
	return nil
}

func (s *Shipping) parseServiceCode(ctx context.Context, serviceCode string) (conn *connectioning.Connection, code string, _ error) {
	if len(serviceCode) <= 4 {
		return nil, "", cm.Errorf(cm.InvalidArgument, nil, "Shipping service code is invalid")
	}
	connCode, code := serviceCode[:4], serviceCode[4:]
	conn, err := s.ShipmentManager.ConnectionManager.GetConnectionByCode(ctx, connCode)
	if err != nil {
		return nil, "", err
	}
	return conn, code, nil
}
