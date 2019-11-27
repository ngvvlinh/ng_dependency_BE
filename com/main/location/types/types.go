package types

import (
	"strconv"

	"etop.vn/api/main/location"
)

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

	Code        string
	Region      Region
	VTPostID    int
	HaravanCode string
}

type District struct {
	Name     string
	NameNorm string
	Alias    []string

	Code         string
	ProvinceCode string
	GhnID        int
	UrbanType    UrbanType
	VTPostID     int
	HaravanCode  string
}

type Ward struct {
	Name        string
	NameNorm    string
	Alias       []string
	VTPostID    int
	HaravanCode string

	Code         string
	DistrictCode string
}

type VTPostWard struct {
	WardsID          int
	WardsName        string
	DistrictID       int
	EtopDistrictCode string
}

func (p *Province) GetProvinceIndex(codeType location.LocationCodeType) string {
	switch codeType {
	case location.LocCodeTypeHaravan:
		return p.HaravanCode
	case location.LocCodeTypeInternal:
		return p.Code
	default:
		return ""
	}
}

func (d *District) GetDistrictIndex(codeType location.LocationCodeType) string {
	switch codeType {
	case location.LocCodeTypeHaravan:
		return d.HaravanCode
	case location.LocCodeTypeGHN:
		if d.GhnID == 0 {
			return ""
		}
		return strconv.Itoa(int(d.GhnID))
	case location.LocCodeTypeInternal:
		return d.Code
	default:
		return ""
	}
}

func (w *Ward) GetWardIndex(codeType location.LocationCodeType) string {
	switch codeType {
	case location.LocCodeTypeHaravan:
		return w.HaravanCode
	case location.LocCodeTypeInternal:
		return w.Code
	default:
		return ""
	}
}
