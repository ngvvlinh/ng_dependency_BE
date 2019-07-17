package location

import (
	"regexp"
	"strings"

	"etop.vn/api/main/location"

	"golang.org/x/text/unicode/norm"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/backend/pkg/services/location/list"
	"etop.vn/backend/pkg/services/location/types"
	"etop.vn/common/l"
)

const HCMUrbanCodes = "760,764,765,766,767,768,769,770,771,772,773,774,775,776,778"
const HCMSuburban1Codes = "761,762,763,777"
const HCMSuburban2Codes = "783,784,785,786,787"
const (
	HCMProvinceCode       = "79"
	HNProvinceCode        = "01"
	BinhDuongProvinceCode = "74"
	DongNaiProvinceCode   = "75"
	VungTauProvinceCode   = "77"
)

var (
	provinceIndexName        = make(map[string]*types.Province)
	provinceIndexNamN        = make(map[string]*types.Province)
	provinceIndexNamX        = make(map[string]*types.Province)
	provinceIndexCode        = make(map[string]*types.Province)
	provinceIndexHaravanCode = make(map[string]*types.Province)

	districtIndexName          = make(map[string][]*types.District)
	districtIndexNamN          = make(map[string][]*types.District)
	districtIndexNamX          = make(map[string][]*types.District)
	districtIndexCode          = make(map[string]*types.District)
	districtIndexGhnID         = make(map[int32]*types.District)
	districtsIndexProvinceCode = make(map[string][]*types.District)
	districtIndexHaravanCode   = make(map[string]*types.District)

	wardIndexName          = make(map[string][]*types.Ward)
	wardIndexNamN          = make(map[string][]*types.Ward)
	wardIndexNamX          = make(map[string][]*types.Ward)
	wardIndexCode          = make(map[string]*types.Ward)
	wardsIndexDistrictCode = make(map[string][]*types.Ward)
	wardIndexHaravanCode   = make(map[string]*types.Ward)

	reNumber           = regexp.MustCompile(`[01-9]+`)
	urbanTypeIndexNamN = make(map[string]types.UrbanType)

	GroupProvinceCodes = []string{BinhDuongProvinceCode, DongNaiProvinceCode, VungTauProvinceCode}
)

func init() {
	for _, urbanType := range types.UrbanTypes {
		namN := validate.NormalizeSearchSimple(urbanType.Name())
		urbanTypeIndexNamN[namN] = urbanType
	}

	for _, province := range list.Provinces {
		province.Name = norm.NFC.String(province.Name)
		province.NameNorm = validate.NormalizeSearchSimple(province.Name)

		name := strings.ToLower(province.Name)
		namN := province.NameNorm
		namX := validate.NormalizeSearchSimple(normalizeProvince(name))
		provinceIndexName[name] = assignProvince(provinceIndexName[name], province)
		provinceIndexNamN[namN] = assignProvince(provinceIndexNamN[namN], province)
		provinceIndexNamX[namX] = assignProvince(provinceIndexNamX[namX], province)
		provinceIndexHaravanCode[province.HaravanCode] = province

		province.Alias = appendAlias(province.Alias, province.Name)
		for _, alias := range province.Alias {
			name := strings.ToLower(norm.NFC.String(alias))
			namN := validate.NormalizeSearchSimple(name)
			namX := normalizeProvince(namN)
			provinceIndexName[name] = assignProvince(provinceIndexName[name], province)
			provinceIndexNamN[namN] = assignProvince(provinceIndexNamN[namN], province)
			provinceIndexNamX[namX] = assignProvince(provinceIndexNamX[namX], province)
		}

		code := province.Code
		if provinceIndexCode[code] != nil {
			ll.Fatal("Duplicated province code", l.String("code", code))
		}
		provinceIndexCode[code] = province

		// init district list as non-nil slice
		districtsIndexProvinceCode[code] = []*types.District{}
	}

	for _, district := range list.Districts {
		district.Name = norm.NFC.String(district.Name)
		district.NameNorm = validate.NormalizeSearchSimple(district.Name)

		name := strings.ToLower(district.Name)
		namN := district.NameNorm
		namX := normalizeDistrict(namN)
		districtIndexName[name] = append(districtIndexName[name], district)
		districtIndexNamN[namN] = append(districtIndexNamN[namN], district)
		districtIndexNamX[namX] = append(districtIndexNamX[namX], district)
		if district.GhnID != 0 {
			districtIndexGhnID[district.GhnID] = district
		}

		if strings.HasPrefix(name, "quận") {
			s := reNumber.FindString(name)
			if s != "" {
				s = strings.TrimLeft(s, "0")
				district.Alias = append(district.Alias, s)
			}
		}

		district.Alias = appendAlias(district.Alias, district.Name)
		for _, alias := range district.Alias {
			name := strings.ToLower(norm.NFC.String(alias))
			namN := validate.NormalizeSearchSimple(name)
			namX := normalizeDistrict(namN)
			districtIndexName[name] = append(districtIndexName[name], district)
			districtIndexNamN[namN] = append(districtIndexNamN[namN], district)
			districtIndexNamX[namX] = append(districtIndexNamX[namX], district)
		}

		// merge districts by province code
		provinceCode := district.ProvinceCode
		districtsIndexProvinceCode[provinceCode] = append(districtsIndexProvinceCode[provinceCode], district)

		code := district.Code
		if districtIndexCode[code] != nil {
			ll.Fatal("Duplicated district code", l.String("code", code))
		}
		districtIndexCode[code] = district
		districtIndexHaravanCode[district.HaravanCode] = district

		// init ward list as non-nil slice
		wardsIndexDistrictCode[code] = []*types.Ward{}
	}

	for _, ward := range list.Wards {
		ward.Name = norm.NFC.String(ward.Name)
		ward.NameNorm = validate.NormalizeSearchSimple(ward.Name)

		name := strings.ToLower(ward.Name)
		namN := ward.NameNorm
		namX := normalizeWard(namN)
		wardIndexName[name] = append(wardIndexName[name], ward)
		wardIndexNamN[namN] = append(wardIndexNamN[namN], ward)
		wardIndexNamX[namX] = append(wardIndexNamX[namX], ward)

		if strings.HasPrefix(name, "phường") {
			s := reNumber.FindString(name)
			if s != "" {
				s = strings.TrimLeft(s, "0")
				ward.Alias = append(ward.Alias, s)
			}
		}

		ward.Alias = appendAlias(ward.Alias, ward.Name)
		for _, alias := range ward.Alias {
			name := strings.ToLower(norm.NFC.String(alias))
			namN := validate.NormalizeSearchSimple(name)
			namX := normalizeWard(namN)
			wardIndexName[name] = append(wardIndexName[name], ward)
			wardIndexNamN[namN] = append(wardIndexNamN[namN], ward)
			wardIndexNamX[namX] = append(wardIndexNamX[namX], ward)
		}

		// merge wards by district code
		districtCode := ward.DistrictCode
		wardsIndexDistrictCode[districtCode] = append(wardsIndexDistrictCode[districtCode], ward)

		code := ward.Code
		if wardIndexCode[code] != nil {
			ll.Fatal("Duplicated ward code", l.String("code", code))
		}
		wardIndexCode[code] = ward
		wardIndexHaravanCode[ward.HaravanCode] = ward
	}

	for _, code := range strings.Split(HCMUrbanCodes, ",") {
		district := districtIndexCode[code]
		if district == nil {
			ll.S.Fatal("Invalid urban district code: ", code)
		}
		district.UrbanType = types.Urban
	}
	for _, code := range strings.Split(HCMSuburban1Codes, ",") {
		district := districtIndexCode[code]
		if district == nil {
			ll.S.Fatal("Invalid suburban 1 district code: ", code)
		}
		district.UrbanType = types.Suburban1
	}
	for _, code := range strings.Split(HCMSuburban2Codes, ",") {
		district := districtIndexCode[code]
		if district == nil {
			ll.S.Fatal("Invalid suburban 2 district code: ", code)
		}
		district.UrbanType = types.Suburban2
	}
}

func appendAlias(alias []string, s string) []string {
	var c int
	var nparts []string

	parts := strings.Split(s, " ")
	for i, p := range parts {
		if p == "" {
			panic("Invalid name: " + s)
		}

		// Đắk -> Đắc (Đắk Lắk)
		if p[len(p)-1] == 'k' {
			// Replace single part (Đắk Lắk -> Đắc Lắk)
			np := p[:len(p)-1] + "c"
			ns := strings.Replace(s, p, np, -1)
			if !cm.StringsContain(alias, ns) {
				alias = append(alias, ns)
			}

			// Replace the whole name (Đắk Lắk -> Đắc Lắc)
			if nparts == nil {
				nparts = make([]string, len(parts))
				copy(nparts, parts)
			}
			nparts[i] = np
			c++
		}
	}
	if c >= 2 {
		// Replace the whole name (Đắk Lắk -> Đắc Lắc)
		ns := strings.Join(nparts, " ")
		if !cm.StringsContain(alias, ns) {
			alias = append(alias, ns)
		}
	}

	// Ea H'leo -> Ea Hleo
	if strings.Contains(s, "'") {
		ns := strings.Replace(s, "'", "", -1)
		if !cm.StringsContain(alias, ns) {
			alias = append(alias, ns)
		}
	}

	return alias
}

func assignProvince(lhs, rhs *types.Province) *types.Province {
	if lhs != nil {
		ll.Debug("duplicated province", l.Object("province", rhs))
		ll.Panic("Duplicated province", l.Object("province", rhs))
	}
	return rhs
}

type duplicatedParents struct {
	parents    []int
	duplicated []int
}

func (d *duplicatedParents) Check(parent int) {
	for _, p := range d.parents {
		if p == parent {
			d.duplicated = append(d.duplicated, parent)
			return
		}
	}

	d.parents = append(d.parents, parent)
}

type duplicatedCounter struct {
	typ      string
	count    int
	names    []string
	dparents []*duplicatedParents
}

func newDuplidatedCounter(typ string) *duplicatedCounter {
	return &duplicatedCounter{typ: typ}
}

func (c *duplicatedCounter) Add(name string) *duplicatedParents {
	c.count++
	c.names = append(c.names, name)

	dp := &duplicatedParents{}
	c.dparents = append(c.dparents, dp)
	return dp
}

func (c *duplicatedCounter) PrintResult() {
	ll.S.Infof("Duplicated %v name: %v", c.typ, c.count)
	if c.count == 0 || c.count >= 20 {
		return
	}

	s := strings.Join(c.names, ",")
	for _, d := range c.dparents {
		if len(d.duplicated) == 0 {
			s += " (No duplicated parent)"
			ll.Info(s)
			return
		}

		ll.S.Warnf("Duplicated parents: %v", d.duplicated)
	}
}

func normalizeLower(name string) string {
	// workaround case: "Nam Ðịnh" -> "nam ðịnh"
	name = strings.Replace(name, "Ð", "Đ", -1)
	s := strings.TrimSpace(strings.ToLower(norm.NFC.String(name)))
	return s
}

// Name must be lowercase. Only remove prefix if the next character is a number or
// a space or a special character.
func normalizePrefix(name string, prefixes ...string) string {
	for _, prefix := range prefixes {
		normName := strings.TrimPrefix(name, prefix)
		if normName != name && normName != "" {
			next := normName[0]
			switch {
			case next >= 'a' && next <= 'z':
				continue // do not trim prefix
			case next >= '0' && next <= '9':
				return normName
			default:
				return strings.TrimSpace(normName[1:])
			}
		}
	}
	return name
}

func NormalizeProvince(name string) string {
	name = validate.NormalizeSearchSimple(name)
	name = normalizeProvince(name)
	return name
}

func NormalizeDistrict(name string) string {
	name = validate.NormalizeSearchSimple(name)
	name = normalizeDistrict(name)
	return name
}

func NormalizeWard(name string) string {
	name = validate.NormalizeSearchSimple(name)
	name = normalizeWard(name)
	return name
}

func normalizeProvince(name string) string {
	return normalizePrefix(name,
		"tỉnh", "thành phố",
		"tinh", "thanh pho",
		"tp", "t p",
	)
}

func normalizeDistrict(name string) string {
	district := normalizePrefix(name,
		"quan", "huyen dao", "huyen",
		"q", "hd", "h d", "h",
		"thanh pho", "tp", "t p", "tp",
		"thi xa", "tx", "t x",
	)
	return strings.TrimLeft(district, "0")
}

func normalizeWard(name string) string {
	ward := normalizePrefix(name,
		"phuong", "xa", "p", "x",
		"thi tran", "tt", "t t", "tx", "t x",
	)
	s := strings.TrimLeft(ward, "0")
	return s
}

type Location struct {
	Province *types.Province
	District *types.District
	Ward     *types.Ward

	// In case multiple location were found and ambiguous,
	// this field will contain all location found.
	// For example, "Huyện Long Mỹ" and "Thị Xã Long Mỹ".
	//
	// If we can find exact ward, we will set AllDistrics to nil (no longer ambiguous).
	OtherDistricts []*types.District
}

// FindLocation tries to return location from input address.
// First, it tries finding by exact name.
// Then, it tries finding by normalized name.
//
// It always returns non-null Location, but
func FindLocation(province, district, ward string) Location {
	return findLocation(province, district, ward)
}

type debugInfoStruct struct {
	items [][]debugInfoStep
}

type debugInfoStep struct {
	String string
	Result bool
}

// This field is used for debug only and will be ignored for production build.
// Not threadsafe.
var debugInfo debugInfoStruct

func writeDebug(typ int, s string, result bool) {
	if debug {
		debugInfo.items[typ] = append(
			debugInfo.items[typ],
			debugInfoStep{s, result},
		)
	}
}

func findLocation(province, district, ward string) (L Location) {
	if debug {
		debugInfo.items = make([][]debugInfoStep, 3)
	}

	province = normalizeLower(province)
	district = normalizeLower(district)
	ward = normalizeLower(ward)
	if province == "" {
		return L
	}

	// Find province by name, then by normalized name
	{
		prov := provinceIndexName[province]
		writeDebug(0, province, prov != nil)

		var normName string
		if prov == nil {
			normName = validate.NormalizeSearchSimple(province)
			prov = provinceIndexNamN[normName]
			writeDebug(0, normName, prov != nil)
		}
		if prov == nil {
			normName = normalizeProvince(normName)
			prov = provinceIndexNamX[normName]
			writeDebug(0, normName, prov != nil)
		}
		if prov == nil {
			return L
		}
		L.Province = prov
	}

	// Find district by name, then by normalized name
	{
		if district == "" {
			return L
		}
		districts := districtIndexName[district]
		writeDebug(1, district, len(districts) > 0)

		var normName string
		if len(districts) == 0 {
			normName = validate.NormalizeSearchSimple(district)
			districts = districtIndexNamN[normName]
			writeDebug(1, normName, len(districts) > 0)
		}
		if len(districts) == 0 {
			normName = normalizeDistrict(normName)
			districts = districtIndexNamX[normName]
			writeDebug(1, normName, len(districts) > 0)
		}

		switch len(districts) {
		case 0:
			// Found province, but no district
			return L

		default:
			// Distinct by province code
			count := 0
			for _, d := range districts {
				if d.ProvinceCode == L.Province.Code {
					// "Huyện Long Mỹ" and "Thị Xã Long Mỹ".
					if count > 0 {
						L.OtherDistricts = append(L.OtherDistricts, d)
					}

					L.District = d
					count++
				}
			}
			if L.District == nil {
				return L
			}
		}
	}

	// Find ward by name, then by normalized name
	{
		if ward == "" {
			return L
		}

		// Because of "Phường 3" and "Phường 03", we only do normalization search on ward.
		normName := normalizeWard(validate.NormalizeSearchSimple(ward))
		wards := wardIndexNamX[normName]
		writeDebug(2, normName, len(wards) > 0)

		switch len(wards) {
		case 0:
			// Found province and district, but no ward
			return L

		default:
			// Distinct by district code
			for _, w := range wards {
				if w.DistrictCode == L.District.Code {
					L.Ward = w
					break
				}
			}
			if L.Ward != nil {
				L.OtherDistricts = nil // no longer ambiguous
				return L
			}

			// Continue find in the province
			foundWards := make(map[string]*types.Ward)
			provinceCode := L.Province.Code
			for _, w := range wards {
				district := districtIndexCode[w.DistrictCode]
				if district.ProvinceCode == provinceCode {
					foundWards[w.Code] = w
					writeDebug(2, "<found in district "+district.Name+">", true)
				}
			}
			switch len(foundWards) {
			case 1:
				for _, w := range foundWards {
					L.Ward = w
					break
				}
			case 0:
				writeDebug(2, "<not found in province>", false)
			default:
				writeDebug(2, "<ambiguous in province>", false)
			}
		}
	}
	return L
}

func FindWardByDistrictCode(name string, districtCode string) *types.Ward {
	name = normalizeLower(name)
	wards := wardIndexName[name]

	var normName string
	if len(wards) == 0 {
		normName = validate.NormalizeSearchSimple(name)
		wards = wardIndexNamN[normName]
	}
	if len(wards) == 0 {
		normName = normalizeWard(normName)
		wards = wardIndexNamX[normName]
	}
	for _, w := range wards {
		if districtCode == w.DistrictCode {
			return w
		}
	}
	return nil
}

func GetDistrictByCode(code string, codeType location.LocationCodeType) *types.District {
	switch codeType {
	case location.LocCodeTypeInternal:
		return districtIndexCode[code]
	case location.LocCodeTypeHaravan:
		return districtIndexHaravanCode[code]
	default:
		// TODO: handle other cases
		return nil
	}
}

func GetProvinceByCode(code string, codeType location.LocationCodeType) *types.Province {
	switch codeType {
	case location.LocCodeTypeInternal:
		return provinceIndexCode[code]
	case location.LocCodeTypeHaravan:
		return provinceIndexHaravanCode[code]
	default:
		// TODO: handle other cases
		return nil
	}
}

func GetWardByCode(code string, codeType location.LocationCodeType) *types.Ward {
	switch codeType {
	case location.LocCodeTypeInternal:
		return wardIndexCode[code]
	case location.LocCodeTypeHaravan:
		return wardIndexHaravanCode[code]
	default:
		// TODO: handle other cases
		return nil
	}
}

func GetDistrictByGhnID(ghnID int32) *types.District {
	return districtIndexGhnID[ghnID]
}

func GetDistrictsByProvinceCode(code string) ([]*types.District, bool) {
	districts, ok := districtsIndexProvinceCode[code]
	return districts, ok
}

func GetWardsByDistrictCode(code string) ([]*types.Ward, bool) {
	wards, ok := wardsIndexDistrictCode[code]
	return wards, ok
}

func CheckValidLocation(code string, codeType string) error {
	codeType = strings.ToLower(codeType)
	switch codeType {
	case "province":
		if _, ok := provinceIndexCode[code]; !ok {
			return cm.Error(cm.NotFound, "Province code does not exist", nil)
		}
	case "district":
		if _, ok := districtIndexCode[code]; !ok {
			return cm.Error(cm.NotFound, "District code does not exist", nil)
		}
	case "ward":
		if _, ok := wardIndexCode[code]; !ok {
			return cm.Error(cm.NotFound, "Ward code does not exist", nil)
		}
	default:
		return cm.Error(cm.NotFound, "Wrong type", nil)
	}
	return nil
}

func GetUrbanType(s string) types.UrbanType {
	namN := validate.NormalizeSearchSimple(s)
	return urbanTypeIndexNamN[namN]
}
