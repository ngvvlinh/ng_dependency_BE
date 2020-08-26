package parse

import (
	"strings"

	"o.o/backend/com/main/location"
	"o.o/backend/com/main/location/list"
	"o.o/backend/com/main/location/types"
)

var maps = buildMaps()

type mapsType struct {
	provByName mapCodes
	distByName mapCodes
	distByProv mapCodes
	wardByName mapCodes
	wardByProv mapCodes
	wardByDist mapCodes

	maxWordProvince int
	maxWordDistrict int
	maxWordWard     int

	districtsByProvinceCode mapDistrictsByProvinceCode
	wardsByDistrictCode     mapWardsByDistrictCode
}

type mapCodes map[string][]string
type mapProvincesByName map[string]*[]*types.Province
type mapDistrictsByName map[string]*[]*types.District
type mapDistrictsByProvinceCode map[string]*[]*types.District
type mapWardsByName map[string]*[]*types.Ward
type mapWardsByDistrictCode map[string]*[]*types.Ward

func (m mapCodes) append(key, code string) {
	for _, si := range m[key] {
		if si == code {
			return
		}
	}
	m[key] = append(m[key], code)
}

func (m mapProvincesByName) append(key string, p *types.Province) {
	s := m[key]
	if s == nil {
		var ss []*types.Province
		s = &ss
		m[key] = s
	}
	for _, pi := range *s {
		if pi.Code == p.Code {
			return
		}
	}
	*s = append(*s, p)
}

func (m mapDistrictsByName) append(key string, p *types.District) {
	s := m[key]
	if s == nil {
		var ss []*types.District
		s = &ss
		m[key] = s
	}
	for _, pi := range *s {
		if pi.Code == p.Code {
			return
		}
	}
	*s = append(*s, p)
}

func (m mapDistrictsByProvinceCode) append(key string, p *types.District) {
	s := m[key]
	if s == nil {
		var ss []*types.District
		s = &ss
		m[key] = s
	}
	for _, pi := range *s {
		if pi.Code == p.Code {
			return
		}
	}
	*s = append(*s, p)
}

func (m mapWardsByName) append(key string, p *types.Ward) {
	s := m[key]
	if s == nil {
		var ss []*types.Ward
		s = &ss
		m[key] = s
	}
	for _, pi := range *s {
		if pi.Code == p.Code {
			return
		}
	}
	*s = append(*s, p)
}

func (m mapWardsByDistrictCode) append(key string, p *types.Ward) {
	s := m[key]
	if s == nil {
		var ss []*types.Ward
		s = &ss
		m[key] = s
	}
	for _, pi := range *s {
		if pi.Code == p.Code {
			return
		}
	}
	*s = append(*s, p)
}

func extractAbbr(input string) string {
	s := make([]byte, 0, 4)
	for i := range input {
		c := input[i]
		if i == 0 || input[i-1] == ' ' {
			s = append(s, c)
		}
	}
	return string(s)
}

func alternativeNames(input string) (result []string) {
	result = append(result, extractAbbr(input)) // abbr

	parts := strings.Split(input, " ")
	result = append(result, strings.Join(parts, "")) // missing space
	for i, c := range input {
		if c == ' ' {
			result = append(result, input[:i-1]+input[i:])
		}
	}
	return
}

func countWords(input string) int {
	c := 0
	for i := range input {
		if input[i] == ' ' {
			c++
		}
	}
	return c
}

func buildMaps() mapsType {
	m := mapsType{}
	m.maxWordProvince, m.provByName = buildProvinces()
	m.maxWordDistrict, m.distByName, m.distByProv = buildDistricts()
	m.maxWordWard, m.wardByName, m.wardByDist, m.wardByProv = buildWards()
	return m
}

func buildProvinces() (maxWords int, byName mapCodes) {
	byName = mapCodes{}
	for _, item := range list.Provinces {
		if n := countWords(item.Name); n > maxWords {
			maxWords = n
		}
		_name := location.NormalizeProvince(item.Name)
		byName.append(_name, item.Code)
		for _, alt := range alternativeNames(_name) {
			byName.append(alt, item.Code)
		}

		for _, alias := range item.Alias {
			_alias := location.NormalizeProvince(alias)
			byName.append(_alias, item.Code)
			for _, alt := range alternativeNames(_alias) {
				byName.append(alt, item.Code)
			}
		}
	}
	return
}

func buildDistricts() (maxWords int, byName mapCodes, byCode mapCodes) {
	byName = mapCodes{}
	byCode = mapCodes{}
	for _, item := range list.Districts {
		if n := countWords(item.Name); n > maxWords {
			maxWords = n
		}
		_name := location.NormalizeDistrict(item.Name)
		byName.append(_name, item.Code)
		for _, alt := range alternativeNames(_name) {
			byName.append(alt, item.Code)
		}
		byCode.append(item.ProvinceCode, item.Code)

		for _, alias := range item.Alias {
			_alias := location.NormalizeDistrict(alias)
			byName.append(_alias, item.Code)
			for _, alt := range alternativeNames(_alias) {
				byName.append(alt, item.Code)
			}
		}
	}
	return
}

func buildWards() (maxWords int, byName, byDist, byProv mapCodes) {
	byName = mapCodes{}
	byDist = mapCodes{}
	byProv = mapCodes{}
	for _, item := range list.Wards {
		if n := countWords(item.Name); n > maxWords {
			maxWords = n
		}
		_name := location.NormalizeWard(item.Name)
		byName.append(_name, item.Code)
		for _, alt := range alternativeNames(_name) {
			byName.append(alt, item.Code)
		}
		byDist.append(item.DistrictCode, item.Code)
		prov := location.DistrictIndexCode[item.DistrictCode].ProvinceCode
		byProv.append(prov, item.Code)

		for _, alias := range item.Alias {
			_alias := location.NormalizeWard(alias)
			byName.append(_alias, item.Code)
			for _, alt := range alternativeNames(_alias) {
				byName.append(alt, item.Code)
			}
		}
	}
	return
}
