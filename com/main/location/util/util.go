package util

import (
	corelocation "etop.vn/api/main/location"
	"etop.vn/api/top/types/etc/location_type"
	"etop.vn/api/top/types/etc/route_type"
	"etop.vn/backend/com/main/location"
	"etop.vn/capi/dot"
)

func CheckUrbanHCMHN(district *corelocation.District) bool {
	// HCM code: "79"
	// HN code: "01"
	if district.ProvinceCode != "79" && district.ProvinceCode != "01" {
		return false
	}
	if district.UrbanType != corelocation.Urban {
		return false
	}
	return true
}

func GetRegionRouteType(fromProvince, toProvince *corelocation.Province) route_type.RegionRouteType {
	if fromProvince.Region == toProvince.Region {
		return route_type.SameRegion
	}
	return route_type.DifferentRegion
}

func GetProvinceRouteType(fromProvince, toProvince *corelocation.Province) route_type.ProvinceRouteType {
	if fromProvince.Code == toProvince.Code {
		return route_type.SameProvince
	}
	return route_type.DifferentProvince
}

func GetShippingDistrictType(district *corelocation.District) route_type.UrbanType {
	switch district.UrbanType {
	case corelocation.Urban:
		return route_type.Urban
	case corelocation.Suburban1:
		return route_type.Suburban1
	case corelocation.Suburban2:
		return route_type.Suburban2
	default:
		return route_type.Suburban2
	}
}

func GetCustomRegionRouteType(fromCustomRegion, toCustomRegion dot.ID) route_type.CustomRegionRouteType {
	if fromCustomRegion == toCustomRegion {
		return route_type.SameCustomRegion
	}
	return route_type.DifferentCustomRegion
}

func GetRegion(provinceCode, districtCode string) location_type.RegionType {
	if provinceCode == "" {
		district := location.GetDistrictByCode(districtCode, corelocation.LocCodeTypeInternal)
		provinceCode = district.ProvinceCode
	}
	provice := location.GetProvinceByCode(provinceCode, corelocation.LocCodeTypeInternal)
	switch provice.Region {
	case corelocation.North:
		return location_type.NorthRegion
	case corelocation.South:
		return location_type.SouthRegion
	case corelocation.Middle:
		return location_type.MiddleRegion
	default:
		return 0
	}
}
