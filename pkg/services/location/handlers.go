package location

import (
	"context"

	"etop.vn/api/main/location"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/services/location/convert"
	"etop.vn/backend/pkg/services/location/list"
	"etop.vn/backend/pkg/services/location/types"
	"etop.vn/common/bus"
	"etop.vn/common/l"
	"etop.vn/common/xerrors"
)

var ll = l.New()
var _ location.LocationQueryService = &Impl{}

type Impl struct {
}

func New() *Impl {
	return &Impl{}
}

func (im *Impl) MessageBus() location.QueryBus {
	b := bus.New()
	return location.NewLocationQueryServiceHandler(im).RegisterHandlers(b)
}

func (im *Impl) GetAllLocations(ctx context.Context, query *location.GetAllLocationsQueryArgs) (result *location.GetAllLocationsQueryResult, err error) {
	result = &location.GetAllLocationsQueryResult{}
	var execFunc func() error

	count := 0
	if query.All {
		count++
		execFunc = func() error {
			return xerrors.FirstError(
				convert.Provinces(list.Provinces, &result.Provinces),
				convert.Districts(list.Districts, &result.Districts),
				convert.Wards(list.Wards, &result.Wards),
			)
		}
	}
	if query.ProvinceCode != "" {
		count++
		execFunc = func() error {
			districts, ok := GetDistrictsByProvinceCode(query.ProvinceCode)
			if !ok {
				return cm.Errorf(cm.InvalidArgument, nil, "invalid province code").
					WithMetap("code", query.ProvinceCode)
			}
			return convert.Districts(districts, &result.Districts)
		}
	}
	if query.DistrictCode != "" {
		count++
		execFunc = func() error {
			wards, ok := GetWardsByDistrictCode(query.DistrictCode)
			if !ok {
				return cm.Errorf(cm.InvalidArgument, nil, "invalid district code").
					WithMetap("code", query.DistrictCode)
			}
			return convert.Wards(wards, &result.Wards)
		}
	}
	if count != 1 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "must provide exactly 1 argument")
	}
	return result, execFunc()
}

func (im *Impl) GetLocation(ctx context.Context, query *location.GetLocationQueryArgs) (result *location.LocationQueryResult, _err error) {
	var ward *types.Ward
	var district *types.District
	var province *types.Province
	defer func() {
		if _err != nil {
			return
		}
		result = &location.LocationQueryResult{}
		_err = xerrors.FirstErrorWithMsg(
			"can not convert location",
			convert.PtrWard(ward, &result.Ward),
			convert.PtrDistrict(district, &result.District),
			convert.PtrProvince(province, &result.Province),
		)
	}()

	if query.WardCode == "" && query.DistrictCode == "" && query.ProvinceCode == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "empty request")
	}
	if query.WardCode != "" {
		ward = GetWardByCode(query.WardCode)
		if ward == nil {
			return nil, cm.Errorf(cm.NotFound, nil, "không tìm thấy phường/xã")
		}
		district = GetDistrictByCode(ward.DistrictCode)
		province = GetProvinceByCode(district.ProvinceCode)
	}
	if query.DistrictCode != "" {
		if district != nil && district.Code != query.DistrictCode {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "mã quận/huyện không thống nhất")
		}
		district = GetDistrictByCode(query.DistrictCode)
		if district == nil {
			return nil, cm.Errorf(cm.NotFound, nil, "không tìm thấy quận/huyện")
		}
		province = GetProvinceByCode(district.ProvinceCode)
	}
	if query.ProvinceCode != "" {
		if province != nil && province.Code != query.ProvinceCode {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "mã tỉnh/thành phố không thống nhất")
		}
		province = GetProvinceByCode(query.ProvinceCode)
		if province == nil {
			return nil, cm.Errorf(cm.NotFound, nil, "không tìm thấy tỉnh/thành phố")
		}
	}
	return result, nil
}

func (im *Impl) FindLocation(ctx context.Context, query *location.FindLocationQueryArgs) (*location.LocationQueryResult, error) {
	if query.Province == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Thiếu thông tin tỉnh")
	}

	loc := FindLocation(query.Province, query.District, query.Ward)
	result := &location.LocationQueryResult{}
	err := xerrors.FirstErrorWithMsg(
		"can not convert location",
		convert.PtrProvince(loc.Province, &result.Province),
		convert.PtrDistrict(loc.District, &result.District),
		convert.PtrWard(loc.Ward, &result.Ward),
	)
	return result, err
}

func (im *Impl) FindOrGetLocation(ctx context.Context, query *location.FindOrGetLocationQueryArgs) (*location.LocationQueryResult, error) {
	switch {
	case query.ProvinceCode != "" || query.DistrictCode != "" || query.WardCode != "":
		locationQuery := &location.GetLocationQueryArgs{
			ProvinceCode: query.ProvinceCode,
			DistrictCode: query.DistrictCode,
			WardCode:     query.WardCode,
		}
		return im.GetLocation(ctx, locationQuery)

	case query.ProvinceCode == "" && query.DistrictCode == "" && query.WardCode == "":
		locationQuery := &location.FindLocationQueryArgs{
			Province: query.Province,
			District: query.District,
			Ward:     query.Ward,
		}
		return im.FindLocation(ctx, locationQuery)

	default:
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Must provide all code or leave it all empty")
	}
}
