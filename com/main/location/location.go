package location

import (
	"regexp"
	"strings"

	"golang.org/x/text/unicode/norm"

	"o.o/api/main/location"
	"o.o/backend/com/main/location/list"
	"o.o/backend/com/main/location/types"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/validate"
	"o.o/common/l"
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
	provinceIndexName = make(map[string]*types.Province)
	provinceIndexNamN = make(map[string]*types.Province)
	provinceIndexNamX = make(map[string]*types.Province)
	ProvinceIndexCode = make(map[string]*types.Province)

	districtIndexName          = make(map[string][]*types.District)
	districtIndexNamN          = make(map[string][]*types.District)
	districtIndexNamX          = make(map[string][]*types.District)
	DistrictIndexCode          = make(map[string]*types.District)
	districtsIndexProvinceCode = make(map[string][]*types.District)

	wardIndexName          = make(map[string][]*types.Ward)
	wardIndexNamN          = make(map[string][]*types.Ward)
	wardIndexNamX          = make(map[string][]*types.Ward)
	WardIndexCode          = make(map[string]*types.Ward)
	wardsIndexDistrictCode = make(map[string][]*types.Ward)

	reNumber           = regexp.MustCompile(`[01-9]+`)
	urbanTypeIndexNamN = make(map[string]location.UrbanType)

	GroupProvinceCodes = []string{BinhDuongProvinceCode, DongNaiProvinceCode, VungTauProvinceCode}
)

type LocationIndex struct {
	ProvinceIndex map[string]*types.Province
	DistrictIndex map[string]*types.District
	WardIndex     map[string]*types.Ward
}

var LocationIndexType = make(map[location.LocationCodeType]*LocationIndex)

// define external code types
var LocationCodeTypes = []location.LocationCodeType{
	location.LocCodeTypeGHN,
	location.LocCodeTypeVTPost,
	location.LocCodeTypeHaravan,
	location.LocCodeTypeInternal,
}

func init() {
	for _, urbanType := range types.UrbanTypes {
		namN := validate.NormalizeSearchSimple(urbanType.Name())
		urbanTypeIndexNamN[namN] = urbanType
	}

	for _, codeType := range LocationCodeTypes {
		LocationIndexType[codeType] = &LocationIndex{
			ProvinceIndex: make(map[string]*types.Province),
			DistrictIndex: make(map[string]*types.District),
			WardIndex:     make(map[string]*types.Ward),
		}
	}

	for _, province := range list.Provinces {
		province.Name = norm.NFC.String(province.Name)
		province.NameNorm = validate.NormalizeSearchSimple(province.Name)

		name := strings.ToLower(province.Name)
		namN := province.NameNorm
		namX := NormalizeProvince(name)
		provinceIndexName[name] = assignProvince(provinceIndexName[name], province)
		provinceIndexNamN[namN] = assignProvince(provinceIndexNamN[namN], province)
		provinceIndexNamX[namX] = assignProvince(provinceIndexNamX[namX], province)
		for _, codeType := range LocationCodeTypes {
			_index := province.GetProvinceIndex(codeType)
			if _index == "" {
				continue
			}
			LocationIndexType[codeType].ProvinceIndex[_index] = province
		}

		province.Alias = appendAlias(province.Alias, province.Name)
		for _, alias := range province.Alias {
			name := strings.ToLower(norm.NFC.String(alias))
			namN := validate.NormalizeSearchSimple(name)
			namX := NormalizeProvince(namN)
			provinceIndexName[name] = assignProvince(provinceIndexName[name], province)
			provinceIndexNamN[namN] = assignProvince(provinceIndexNamN[namN], province)
			provinceIndexNamX[namX] = assignProvince(provinceIndexNamX[namX], province)
		}

		code := province.Code
		if ProvinceIndexCode[code] != nil {
			ll.Fatal("Duplicated province code", l.String("code", code))
		}
		ProvinceIndexCode[code] = province

		// init district list as non-nil slice
		districtsIndexProvinceCode[code] = []*types.District{}
	}

	for _, district := range list.Districts {
		district.Name = norm.NFC.String(district.Name)
		district.NameNorm = validate.NormalizeSearchSimple(district.Name)

		name := strings.ToLower(district.Name)
		namN := district.NameNorm
		namX := NormalizeDistrict(namN)
		districtIndexName[name] = append(districtIndexName[name], district)
		districtIndexNamN[namN] = append(districtIndexNamN[namN], district)
		districtIndexNamX[namX] = append(districtIndexNamX[namX], district)

		if strings.HasPrefix(name, "qu???n") {
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
		if DistrictIndexCode[code] != nil {
			ll.Fatal("Duplicated district code", l.String("code", code))
		}
		DistrictIndexCode[code] = district
		for _, codeType := range LocationCodeTypes {
			_index := district.GetDistrictIndex(codeType)
			if _index == "" {
				continue
			}
			LocationIndexType[codeType].DistrictIndex[_index] = district
		}

		// init ward list as non-nil slice
		wardsIndexDistrictCode[code] = []*types.Ward{}
	}

	for _, ward := range list.Wards {
		ward.ProvinceCode = DistrictIndexCode[ward.DistrictCode].ProvinceCode
		ward.Name = norm.NFC.String(ward.Name)
		ward.NameNorm = validate.NormalizeSearchSimple(ward.Name)

		name := strings.ToLower(ward.Name)
		namN := ward.NameNorm
		namX := NormalizeWard(namN)
		wardIndexName[name] = append(wardIndexName[name], ward)
		wardIndexNamN[namN] = append(wardIndexNamN[namN], ward)
		wardIndexNamX[namX] = append(wardIndexNamX[namX], ward)

		if strings.HasPrefix(name, "ph?????ng") {
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
			namX := NormalizeWard(namN)
			wardIndexName[name] = append(wardIndexName[name], ward)
			wardIndexNamN[namN] = append(wardIndexNamN[namN], ward)
			wardIndexNamX[namX] = append(wardIndexNamX[namX], ward)
		}

		// merge wards by district code
		districtCode := ward.DistrictCode
		wardsIndexDistrictCode[districtCode] = append(wardsIndexDistrictCode[districtCode], ward)

		code := ward.Code
		if WardIndexCode[code] != nil {
			ll.Fatal("Duplicated ward code", l.String("code", code))
		}
		WardIndexCode[code] = ward
		for _, codeType := range LocationCodeTypes {
			_index := ward.GetWardIndex(codeType)
			if _index == "" {
				continue
			}
			LocationIndexType[codeType].WardIndex[_index] = ward
		}
	}

	for _, code := range strings.Split(HCMUrbanCodes, ",") {
		district := DistrictIndexCode[code]
		if district == nil {
			ll.S.Fatal("Invalid urban district code: ", code)
			panic("unexpected")
		}
		district.UrbanType = location.Urban
	}
	for _, code := range strings.Split(HCMSuburban1Codes, ",") {
		district := DistrictIndexCode[code]
		if district == nil {
			ll.S.Fatal("Invalid suburban 1 district code: ", code)
			panic("unexpected")
		}
		district.UrbanType = location.Suburban1
	}
	for _, code := range strings.Split(HCMSuburban2Codes, ",") {
		district := DistrictIndexCode[code]
		if district == nil {
			ll.S.Fatal("Invalid suburban 2 district code: ", code)
			panic("unexpected")
		}
		district.UrbanType = location.Suburban2
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

		// ?????k -> ?????c (?????k L???k)
		if p[len(p)-1] == 'k' {
			// Replace single part (?????k L???k -> ?????c L???k)
			np := p[:len(p)-1] + "c"
			ns := strings.Replace(s, p, np, -1)
			if !cm.StringsContain(alias, ns) {
				alias = append(alias, ns)
			}

			// Replace the whole name (?????k L???k -> ?????c L???c)
			if nparts == nil {
				nparts = make([]string, len(parts))
				copy(nparts, parts)
			}
			nparts[i] = np
			c++
		}
	}
	if c >= 2 {
		// Replace the whole name (?????k L???k -> ?????c L???c)
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
	// workaround case: "Nam ?????nh" -> "nam ?????nh"
	name = strings.Replace(name, "??", "??", -1)
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

var PrefixProvince = []string{"tinh", "t", "thanh pho", "tp", "t p"}
var PrefixDistrict = []string{
	"quan", "q", "huyen dao", "hd", "h d", "huyen", "h",
	"thanh pho", "tp", "t p", "tp", "thi xa", "tx", "t x",
}
var PrefixWard = []string{"phuong", "xa", "p", "x", "thi tran", "tt", "t t"}

func normalizeProvince(name string) string {
	return normalizePrefix(name, PrefixProvince...)
}

func normalizeDistrict(name string) string {
	district := normalizePrefix(name, PrefixDistrict...)
	return strings.TrimLeft(district, "0")
}

func normalizeWard(name string) string {
	ward := normalizePrefix(name, PrefixWard...)
	s := strings.TrimLeft(ward, "0")
	return s
}

type Location struct {
	Province *types.Province
	District *types.District
	Ward     *types.Ward

	// In case multiple locations were found and ambiguous, this field will
	// contain all locations found. For example, "Huy???n Long M???" and "Th??? X??
	// Long M???".
	//
	// If we can find exact ward, we will set AllDistricts to nil (no longer
	// ambiguous).
	OtherDistricts []*types.District
}

// FindLocation tries to return location from input address. First, it tries
// finding by exact name. Then, it tries finding by normalized name.
//
// It always returns non-null Location, but each value may be empty
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
					// "Huy???n Long M???" and "Th??? X?? Long M???".
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

		// Because of "Ph?????ng 3" and "Ph?????ng 03", we only do normalization search on ward.
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
				district := DistrictIndexCode[w.DistrictCode]
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

func ParseLocation(rawLocation string) (L Location) {
	if rawLocation == "" {
		return
	}

	{
		_location := []rune(rawLocation)
		for i := 0; i < len(_location); i++ {
			if _location[i] == '.' || _location[i] == ',' {
				_location[i] = ' '
			}
		}
		rawLocation = string(_location)
	}

	L = parseLocation(rawLocation)
	if L.Province == nil && L.District != nil {
		L.Province = ProvinceIndexCode[L.District.ProvinceCode]
	}

	return
}

func parseLocation(rawLocation string) (L Location) {
	// Start parse location from province
	normLocation := validate.NormalizeSearchSimple(rawLocation)

	var L1 Location
	{
		currLocation := strings.Split(normLocation, " ")
		for i := len(currLocation) - 1; i >= 0; i-- {
			for j := i; j < len(currLocation); j++ {
				var tempL Location
				province := strings.Join(currLocation[i:j+1], " ")

				prov := provinceIndexNamN[province]

				if prov == nil {
					province = normalizeProvince(province)
					prov = provinceIndexNamX[province]
				}

				if prov == nil {
					continue
				}

				tempL.Province = prov
				if L1.Province == nil {
					L1 = tempL
				}

				currLocation = removeElements(i, j, currLocation)
				for _i := len(currLocation) - 1; _i >= 0; _i-- {
					for _j := _i; _j < len(currLocation); _j++ {
						district := strings.Join(currLocation[_i:_j+1], " ")

						districts := districtIndexNamN[district]

						if len(districts) == 0 {
							district = normalizeDistrict(district)
							districts = districtIndexNamX[district]
						}

						if len(districts) == 0 {
							continue
						}

						for _, district := range districts {
							if district.ProvinceCode == prov.Code {
								tempL.District = district
								if L1.District == nil {
									L1 = tempL
								}
								break
							}
						}

						if tempL.District == nil {
							continue
						}

						currLocation = removeElements(_i, _j, currLocation)
						for _ii := len(currLocation) - 1; _ii >= 0; _ii-- {
							for _jj := _ii; _jj < len(currLocation); _jj++ {
								ward := strings.Join(currLocation[_ii:_jj+1], " ")

								wards := wardIndexNamN[ward]

								if len(wards) == 0 {
									ward = NormalizeWard(ward)
									wards = wardIndexNamX[ward]
								}

								if len(wards) == 0 {
									continue
								}

								for _, ward := range wards {
									if ward.DistrictCode == tempL.District.Code {
										tempL.Ward = ward
										return tempL
									}
								}
							}
						}
					}
				}
			}
		}
	}

	if L1.Province != nil {
		return L1
	}

	var L2 Location
	{
		currLocation := strings.Split(normLocation, " ")
		for _i := len(currLocation) - 1; _i >= 0; _i-- {
			for _j := _i; _j < len(currLocation); _j++ {
				var tempL Location
				district := strings.Join(currLocation[_i:_j+1], " ")
				district = normalizeLower(district)

				normName := validate.NormalizeSearchSimple(district)
				districts := districtIndexNamN[normName]

				if len(districts) == 0 {
					normName = normalizeDistrict(normName)
					districts = districtIndexNamX[normName]
				}

				if len(districts) == 0 {
					continue
				}

				for _, district := range districts {
					tempL.District = district
					if L2.District == nil {
						L2 = tempL
					}
					break
				}

				if tempL.District == nil {
					continue
				}

				currLocation = removeElements(_i, _j, currLocation)
				for _ii := len(currLocation) - 1; _ii >= 0; _ii-- {
					for _jj := _ii; _jj < len(currLocation); _jj++ {
						ward := strings.Join(currLocation[_ii:_jj+1], " ")
						ward = normalizeLower(ward)

						normName = validate.NormalizeSearchSimple(ward)
						wards := wardIndexNamN[normName]

						if len(wards) == 0 {
							normName = NormalizeWard(normName)
							wards = wardIndexNamX[normName]
						}

						if len(wards) == 0 {
							continue
						}

						for _, ward := range wards {
							if ward.DistrictCode == tempL.District.Code {
								tempL.Ward = ward
								return tempL
							}
						}
					}
				}
			}
		}
	}

	return L2
}

func removeElements(start, end int, arrs []string) []string {
	var result []string

	for i := 0; i < len(arrs); i++ {
		if start <= i && i <= end {
			continue
		}
		result = append(result, arrs[i])
	}

	return result
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
	return LocationIndexType[codeType].DistrictIndex[code]
}

func GetProvinceByCode(code string, codeType location.LocationCodeType) *types.Province {
	return LocationIndexType[codeType].ProvinceIndex[code]
}

func GetWardByCode(code string, codeType location.LocationCodeType) *types.Ward {
	return LocationIndexType[codeType].WardIndex[code]
}

func GetDistrictsByProvinceCode(code string) ([]*types.District, bool) {
	districts, ok := districtsIndexProvinceCode[code]
	return districts, ok
}

func GetWardsByDistrictCode(code string) ([]*types.Ward, bool) {
	wards, ok := wardsIndexDistrictCode[code]
	return wards, ok
}

func GetUrbanType(s string) location.UrbanType {
	namN := validate.NormalizeSearchSimple(s)
	return urbanTypeIndexNamN[namN]
}
