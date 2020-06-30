package location

import (
	"context"

	"o.o/api/main/location"
	"o.o/api/meta"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/location/convert"
	"o.o/backend/com/main/location/list"
	"o.o/backend/com/main/location/sqlstore"
	"o.o/backend/com/main/location/types"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/capi/dot"
	"o.o/common/l"
	"o.o/common/xerrors"
)

var ll = l.New()
var _ location.QueryService = &Query{}

type Query struct {
	customRegionStore sqlstore.CustomRegionFactory
}

func New(db com.MainDB) *Query {
	return &Query{
		customRegionStore: sqlstore.NewCustomRegionStore(db),
	}
}

func QueryMessageBus(im *Query) location.QueryBus {
	b := bus.New()
	return location.NewQueryServiceHandler(im).RegisterHandlers(b)
}

func (im *Query) GetAllLocations(ctx context.Context, query *location.GetAllLocationsQueryArgs) (result *location.GetAllLocationsQueryResult, err error) {
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
	if count != 1 || execFunc == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "must provide exactly 1 argument")
	}
	return result, execFunc()
}

func (im *Query) GetLocation(ctx context.Context, query *location.GetLocationQueryArgs) (result *location.LocationQueryResult, _err error) {
	switch query.LocationCodeType {
	case location.LocCodeTypeVTPost:
		return nil, cm.Error(cm.Unimplemented, "Address type does not valid", nil)
	default:

	}
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
		ward = GetWardByCode(query.WardCode, query.LocationCodeType)
		if ward == nil {
			return nil, cm.Errorf(cm.NotFound, nil, "không tìm thấy phường/xã")
		}
		district = GetDistrictByCode(ward.DistrictCode, location.LocCodeTypeInternal)
		if district != nil {
			province = GetProvinceByCode(district.ProvinceCode, location.LocCodeTypeInternal)
		}
	}
	if query.DistrictCode != "" {
		if district != nil {
			districtCode := district.GetDistrictIndex(query.LocationCodeType)
			if districtCode != query.DistrictCode {
				return nil, cm.Errorf(cm.InvalidArgument, nil, "mã quận/huyện không thống nhất")
			}
		}
		district = GetDistrictByCode(query.DistrictCode, query.LocationCodeType)
		if district == nil {
			return nil, cm.Errorf(cm.NotFound, nil, "không tìm thấy quận/huyện")
		}
		province = GetProvinceByCode(district.ProvinceCode, location.LocCodeTypeInternal)
	}
	if query.ProvinceCode != "" {
		if province != nil {
			provinceCode := province.GetProvinceIndex(query.LocationCodeType)
			if provinceCode != query.ProvinceCode {
				return nil, cm.Errorf(cm.InvalidArgument, nil, "mã tỉnh/thành phố không thống nhất")
			}
		}
		province = GetProvinceByCode(query.ProvinceCode, query.LocationCodeType)
		if province == nil {
			return nil, cm.Errorf(cm.NotFound, nil, "không tìm thấy tỉnh/thành phố")
		}
	}
	return result, nil
}

func (im *Query) FindLocation(ctx context.Context, query *location.FindLocationQueryArgs) (*location.LocationQueryResult, error) {
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

func (im *Query) FindOrGetLocation(ctx context.Context, query *location.FindOrGetLocationQueryArgs) (*location.LocationQueryResult, error) {
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

func (im *Query) GetCustomRegion(ctx context.Context, id dot.ID) (*location.CustomRegion, error) {
	return im.customRegionStore(ctx).ID(id).GetCustomRegion()
}

func (im *Query) ListCustomRegionsByCode(ctx context.Context, provinceCode string, districtCode string) ([]*location.CustomRegion, error) {
	if provinceCode == "" && districtCode == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing required params")
	}
	if provinceCode == "" {
		district := GetDistrictByCode(districtCode, location.LocCodeTypeInternal)
		provinceCode = district.ProvinceCode
	}
	return im.customRegionStore(ctx).ProvinceCode(provinceCode).ListCustomRegions()
}

func (im *Query) ListCustomRegions(ctx context.Context, _ *meta.Empty) ([]*location.CustomRegion, error) {
	return im.customRegionStore(ctx).ListCustomRegions()
}
