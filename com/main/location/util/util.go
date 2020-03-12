package util

import (
	"etop.vn/api/main/location"
	"etop.vn/api/top/types/etc/route_type"
	"etop.vn/capi/dot"
)

func CheckUrbanHCMHN(district *location.District) bool {
	// HCM code: "79"
	// HN code: "01"
	if district.ProvinceCode != "79" && district.ProvinceCode != "01" {
		return false
	}
	if district.UrbanType != location.Urban {
		return false
	}
	return true
}

func GetRegionRouteType(fromProvince, toProvince *location.Province) route_type.RegionRouteType {
	if fromProvince.Region == toProvince.Region {
		return route_type.SameRegion
	}
	return route_type.DifferentRegion
}

func GetProvinceRouteType(fromProvince, toProvince *location.Province) route_type.ProvinceRouteType {
	if fromProvince.Code == toProvince.Code {
		return route_type.SameProvince
	}
	return route_type.DifferentProvince
}

func GetShippingDistrictType(district *location.District) route_type.UrbanType {
	switch district.UrbanType {
	case location.Urban:
		return route_type.Urban
	case location.Suburban1:
		return route_type.Suburban1
	case location.Suburban2:
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
