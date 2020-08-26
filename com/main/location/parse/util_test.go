package parse

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"o.o/backend/com/main/location"
)

func TestBuildProvince(t *testing.T) {
	_, m := buildProvinces()
	ks := make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	for _, k := range ks {
		var names []string
		for _, v := range m[k] {
			names = append(names, location.ProvinceIndexCode[v].Name)
		}
		// fmt.Println(k, names)
	}

	require.Equal(t, len(m["brvt"]), 1)
	assert.Equal(t, location.ProvinceIndexCode[m["brvt"][0]].Name, "Tỉnh Bà Rịa - Vũng Tàu")

	require.Equal(t, len(m["vt"]), 1)
	assert.Equal(t, location.ProvinceIndexCode[m["vt"][0]].Name, "Tỉnh Bà Rịa - Vũng Tàu")

	require.Equal(t, len(m["hn"]), 2)
}

func TestBuildDistricts(t *testing.T) {
	_, mName, _ := buildDistricts()
	ks := make([]string, 0, len(mName))
	for k := range mName {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	for _, k := range ks {
		var names []string
		for _, v := range mName[k] {
			names = append(names, location.DistrictIndexCode[v].Name)
		}
		// fmt.Println(k, names)
	}
}

func TestExtractAlias(t *testing.T) {
	assert.Equal(t, "n", extractAbbr("nam"))
	assert.Equal(t, "lc", extractAbbr("lao cai"))
	assert.Equal(t, "tth", extractAbbr("thua thien hue"))
}
