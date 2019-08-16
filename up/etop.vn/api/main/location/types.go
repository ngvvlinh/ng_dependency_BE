package location

const CountryVietnam = "Việt Nam"

type VietnamRegion int32

const (
	North  VietnamRegion = 1
	Middle VietnamRegion = 2
	South  VietnamRegion = 3
)

type UrbanType int32

const (
	Unknown   UrbanType = 0
	Urban     UrbanType = -1
	Suburban1 UrbanType = 1
	Suburban2 UrbanType = 2
)

func (g VietnamRegion) Name() string {
	switch g {
	case North:
		return "Miền Bắc"
	case Middle:
		return "Miền Trung"
	case South:
		return "Miền Nam"
	default:
		return "?"
	}
}

func (a UrbanType) Name() string {
	switch a {
	case Urban:
		return "Nội thành"
	case Suburban1:
		return "Ngoại thành 1"
	case Suburban2:
		return "Ngoại thành 2"
	default:
		return "?"
	}
}

type LocationCodeType int32

const (
	LocCodeTypeInternal LocationCodeType = 0
	LocCodeTypeGHN      LocationCodeType = 1
	LocCodeTypeVTPost   LocationCodeType = 2
	LocCodeTypeHaravan  LocationCodeType = 3
)

type Province struct {
	Name   string
	Code   string
	Region VietnamRegion
	Extra
}

type District struct {
	Name         string
	Code         string
	ProvinceCode string
	UrbanType    UrbanType
	Extra
}

type Ward struct {
	Name         string
	Code         string
	DistrictCode string
	Extra
}

type Extra struct {
	Special     bool
	GhnId       int32
	VtpostId    int32
	HaravanCode string
}
