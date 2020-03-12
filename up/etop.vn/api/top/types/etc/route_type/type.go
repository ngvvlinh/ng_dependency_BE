package route_type

// +enum
// +enum:zero=null
type RegionRouteType int

type NullRegionRouteType struct {
	Enum  RegionRouteType
	Valid bool
}

const (
	// +enum=noi_mien
	SameRegion RegionRouteType = 1

	// +enum=lien_mien
	DifferentRegion RegionRouteType = 2
)

// +enum
// +enum:zero=null
type CustomRegionRouteType int

type NullCustomRegionRouteType struct {
	Enum  CustomRegionRouteType
	Valid bool
}

const (
	// +enum=noi_vung
	SameCustomRegion CustomRegionRouteType = 1

	// +enum=lien_vung
	DifferentCustomRegion CustomRegionRouteType = 2
)

// +enum
// +enum:zero=null
type ProvinceRouteType int

type NullProvinceRouteType struct {
	Enum  ProvinceRouteType
	Valid bool
}

const (
	// +enum=noi_tinh
	SameProvince ProvinceRouteType = 1

	// +enum=lien_tinh
	DifferentProvince ProvinceRouteType = 2
)

// +enum
// +enum:zero=null
type UrbanType int

type NullUrbanType struct {
	Enum  UrbanType
	Valid bool
}

const (
	// +enum=unknown
	Unknown UrbanType = 0

	// +enum=noi_thanh
	Urban UrbanType = 1

	// +enum=ngoai_thanh_1
	Suburban1 UrbanType = 2

	// +enum=ngoai_thanh_2
	Suburban2 UrbanType = 3
)
