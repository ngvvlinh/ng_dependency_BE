package types

import (
	"strconv"

	"o.o/api/main/location"
)

const CountryVietnam = "Việt Nam"

var UrbanTypes = []location.UrbanType{location.Urban, location.Suburban1, location.Suburban2}

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
	Region      location.VietnamRegion
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
	UrbanType    location.UrbanType
	VTPostID     int
	HaravanCode  string
}

type Ward struct {
	Name        string
	NameNorm    string
	Alias       []string
	VTPostID    int
	GhnCode     string
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
		return strconv.Itoa(d.GhnID)
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
