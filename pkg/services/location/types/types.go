package types

import "etop.vn/api/main/location"

const CountryVietnam = "Việt Nam"

type Region = location.VietnamRegion

const (
	North  = location.North
	Middle = location.Middle
	South  = location.South
)

type UrbanType = location.UrbanType

const (
	Urban     = location.Urban
	Suburban1 = location.Suburban1
	Suburban2 = location.Suburban2
)

var UrbanTypes = []UrbanType{Urban, Suburban1, Suburban2}

type Country struct {
	Name     string
	NameNorm string
}

type Province struct {
	Name     string
	NameNorm string
	Alias    []string
	Special  bool

	Code     string
	Region   Region
	VTPostID int32
}

type District struct {
	Name     string
	NameNorm string
	Alias    []string

	Code         string
	ProvinceCode string
	GhnID        int32
	UrbanType    UrbanType
	VTPostID     int32
}

type Ward struct {
	Name     string
	NameNorm string
	Alias    []string
	VTPostID int32

	Code         string
	DistrictCode string
}

type VTPostWard struct {
	WardsID          int32
	WardsName        string
	DistrictID       int32
	EtopDistrictCode string
}
