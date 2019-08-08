package location

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"etop.vn/api/main/location"
	"etop.vn/backend/com/main/location/list"
	"etop.vn/backend/com/main/location/types"
)

func TestNormalizeDistrict(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"Q.Bình Tân", "binh tan"},
		{"Quận3", "3"},
		{"quan3", "3"},
		{"q3", "3"},
		{"q03", "3"},
		{"q003", "3"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := NormalizeDistrict(tt.name)
			assert.Equal(t, tt.expected, output)
		})
	}
}

func TestNormalizeWard(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"TT Tam Sơn", "tam son"},
		{"TX Trà Vinh", "tra vinh"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := NormalizeWard(tt.name)
			assert.Equal(t, tt.expected, output)
		})
	}
}

func TestAlias(t *testing.T) {
	tests := []struct {
		name  string
		alias []string
	}{
		{"Đắk Lắk", []string{"Đắc Lắk", "Đắk Lắc", "Đắc Lắc"}},
		{"Ea H'leo", []string{"Ea Hleo"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var alias []string
			alias = appendAlias(alias, tt.name)
			assert.EqualValues(t, tt.alias, alias)
		})
	}
}

func TestFindLocation(t *testing.T) {
	tests := []struct {
		name     string
		location string
		expected string
	}{
		{
			"Quận 3",
			"Hồ Chí Minh, Quận 3, Phường 03",
			"Thành phố Hồ Chí Minh,Quận 3,Phường 03 - 79 770 27154",
		},
		{
			"Duplicated ward name (1)",
			"Cao Bằng, Huyện Bảo Lâm, Xã Quảng Lâm",
			"Tỉnh Cao Bằng,Huyện Bảo Lâm,Xã Quảng Lâm - 04 042 01303",
		},
		{
			"Duplicated ward name (2)",
			"Điện Biên, Huyện Mường Nhé, Xã Quảng Lâm",
			"Tỉnh Điện Biên,Huyện Mường Nhé,Xã Quảng Lâm - 11 096 03164",
		},
		{
			"Fuzzy matching (1)",
			"dien bien, huyen muong nhe, xã quảng lâm",
			"Tỉnh Điện Biên,Huyện Mường Nhé,Xã Quảng Lâm - 11 096 03164",
		},
		{
			"Fuzzy matching (2)",
			"dien bien, muong nhe, quảng lâm",
			"Tỉnh Điện Biên,Huyện Mường Nhé,Xã Quảng Lâm - 11 096 03164",
		},
		{
			"Fuzzy matching (3) - Alias",
			"Vũng Tàu, Huyện Đất Đỏ, Thị trấn Phước Hải",
			"Tỉnh Bà Rịa - Vũng Tàu,Huyện Đất Đỏ,Thị trấn Phước Hải - 77 753 26692",
		},
		{
			"Fuzzy matching (4) - Alias",
			"Bà Rịa, Huyện Đất Đỏ, Thị trấn Phước Hải",
			"Tỉnh Bà Rịa - Vũng Tàu,Huyện Đất Đỏ,Thị trấn Phước Hải - 77 753 26692",
		},
		{
			"Fuzzy matching (5) - Alias",
			"Tỉnh Lai Châu, Nậm Nhùm,",
			"Tỉnh Lai Châu,Huyện Nậm Nhùn - 12 112",
		},
		{
			"Fuzzy matching: only province",
			"dien bien, muong nha, quảng lâm",
			"Tỉnh Điện Biên - 11",
		},
		{
			"Fuzzy matching: special name",
			"dak lak, ea h leo, xa ea h'leo",
			"Tỉnh Đắk Lắk,Huyện Ea H'leo,Xã Ea H'leo - 66 645 24184",
		},
		{
			"Fuzzy matching: Quận 03 (leading zero)",
			"Hồ Chí Minh, Quận 03, Phường 004",
			"Thành phố Hồ Chí Minh,Quận 3,Phường 04 - 79 770 27148",
		},
		{
			"Fuzzy matching: Quận3 (no space)",
			"Hồ Chí Minh, Quận3,Phường4",
			"Thành phố Hồ Chí Minh,Quận 3,Phường 04 - 79 770 27148",
		},
		{
			"Fuzzy matching: quan3 (no accent)",
			"Hồ Chí Minh, quan3,phuong4",
			"Thành phố Hồ Chí Minh,Quận 3,Phường 04 - 79 770 27148",
		},
		{
			"Fuzzy matching: q3 (only q)",
			"Hồ Chí Minh, q3,p4",
			"Thành phố Hồ Chí Minh,Quận 3,Phường 04 - 79 770 27148",
		},
		{
			"Fuzzy matching: q.3 (only q.)",
			"Hồ Chí Minh, q.3,p. 4",
			"Thành phố Hồ Chí Minh,Quận 3,Phường 04 - 79 770 27148",
		},
		{
			"Fuzzy matching: tx.phu tho (only tx.)",
			"Phu tho, tx.phu tho,",
			"Tỉnh Phú Thọ,Thị xã Phú Thọ - 25 228",
		},
		{
			"Fuzzy matching: q03 (only q with leading zero)",
			"Hồ Chí Minh, q03,p04",
			"Thành phố Hồ Chí Minh,Quận 3,Phường 04 - 79 770 27148",
		},
		{
			"Fuzzy matching: Phường 3 and Phường 03",
			"Hồ Chí Minh,Quận Tân Bình,Phường 3",
			"Thành phố Hồ Chí Minh,Quận Tân Bình,Phường 03 - 79 766 26980",
		},
		{
			"Fuzzy matching: 3 (only number)",
			"Hồ Chí Minh,3,4",
			"Thành phố Hồ Chí Minh,Quận 3,Phường 04 - 79 770 27148",
		},
		{
			"No match",
			"somewhere,,",
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			locs := strings.Split(tt.location, ",")
			loc := FindLocation(locs[0], locs[1], locs[2])
			output := locToString(loc)
			assert.Equal(t, tt.expected, output)
		})
	}
}

func locToString(loc Location) string {
	if loc.Province == nil {
		return ""
	}
	s := loc.Province.Name
	c := " - " + loc.Province.Code
	if loc.District != nil {
		s += "," + loc.District.Name
		c += " " + loc.District.Code
	}
	if loc.Ward != nil {
		s += "," + loc.Ward.Name
		c += " " + loc.Ward.Code
	}
	return s + c
}

func TestUrban(t *testing.T) {
	districts := districtsIndexProvinceCode["79"] // HCM
	for _, d := range districts {
		assert.Equal(t, "79", d.ProvinceCode)
		assert.NotEqual(t, types.UrbanType(0), d.UrbanType)
	}

	assert.Equal(t, 15, len(strings.Split(HCMUrbanCodes, ",")))
	assert.Equal(t, 4, len(strings.Split(HCMSuburban1Codes, ",")))
	assert.Equal(t, 5, len(strings.Split(HCMSuburban2Codes, ",")))
}

func TestGetAllLocations(t *testing.T) {
	locationBus := New().MessageBus()
	ctx := context.Background()

	t.Run("no input (error)", func(t *testing.T) {
		query := &location.GetAllLocationsQuery{}
		err := locationBus.Dispatch(ctx, query)
		assert.EqualError(t, err, "must provide exactly 1 argument")
	})

	t.Run("all", func(t *testing.T) {
		query := &location.GetAllLocationsQuery{All: true}
		err := locationBus.Dispatch(ctx, query)
		require.NoError(t, err)
		assert.Equal(t, len(query.Result.Provinces), len(list.Provinces))
		assert.Equal(t, len(query.Result.Districts), len(list.Districts))
		assert.Equal(t, len(query.Result.Wards), len(list.Wards))
	})

	t.Run("ProvinceCode", func(t *testing.T) {
		provinceCode := HCMProvinceCode
		query := &location.GetAllLocationsQuery{ProvinceCode: provinceCode}
		err := locationBus.Dispatch(ctx, query)
		require.NoError(t, err)
		assert.Nil(t, query.Result.Provinces)
		assert.Equal(t, len(query.Result.Districts), len(districtsIndexProvinceCode[provinceCode]))
		assert.Nil(t, query.Result.Wards)
		for _, district := range query.Result.Districts {
			assert.Equal(t, district.ProvinceCode, provinceCode)
		}
	})

	t.Run("DistrictCode", func(t *testing.T) {
		t.Run("Normal", func(t *testing.T) {
			districtCode := "769" // Quận 2, HCM
			query := &location.GetAllLocationsQuery{DistrictCode: districtCode}
			err := locationBus.Dispatch(ctx, query)
			require.NoError(t, err)
			assert.Nil(t, query.Result.Provinces)
			assert.Nil(t, query.Result.Districts)
			assert.Equal(t, len(query.Result.Wards), len(wardsIndexDistrictCode[districtCode]))
			for _, ward := range query.Result.Wards {
				assert.Equal(t, ward.DistrictCode, districtCode)
			}
		})

		// returned wards is not nil, and has length equal to zero
		t.Run("District with no ward", func(t *testing.T) {
			districtCode := "318" // Huyện đảo Bạch Long Vĩ
			query := &location.GetAllLocationsQuery{DistrictCode: districtCode}
			err := locationBus.Dispatch(ctx, query)
			require.NoError(t, err)
			assert.NotNil(t, query.Result.Wards)
			assert.Zero(t, len(query.Result.Wards))
		})
	})
}

func TestGetLocation(t *testing.T) {
	locationBus := New().MessageBus()
	ctx := context.Background()

	t.Run("no input (error)", func(t *testing.T) {
		query := &location.GetLocationQuery{}
		err := locationBus.Dispatch(ctx, query)
		assert.EqualError(t, err, "empty request")
	})

	t.Run("not found (error)", func(t *testing.T) {
		query := &location.GetLocationQuery{DistrictCode: "989"}
		err := locationBus.Dispatch(ctx, query)
		assert.EqualError(t, err, "không tìm thấy quận/huyện")
	})

	t.Run("ProvinceCode", func(t *testing.T) {
		provinceCode := HCMProvinceCode
		query := &location.GetLocationQuery{ProvinceCode: provinceCode}
		err := locationBus.Dispatch(ctx, query)
		require.NoError(t, err)
		assert.Nil(t, query.Result.Ward)
		assert.Nil(t, query.Result.District)

		province := query.Result.Province
		require.NotNil(t, province)
		assert.Equal(t, province.Code, provinceCode)
		assert.Equal(t, province.Name, "Thành phố Hồ Chí Minh")
	})

	t.Run("DistrictCode", func(t *testing.T) {
		districtCode := "867" // Thành phố Sa Đéc
		provinceCode := "87"  // Tỉnh Đồng Tháp
		query0 := &location.GetLocationQuery{DistrictCode: districtCode}
		query1 := &location.GetLocationQuery{
			DistrictCode: districtCode,
			ProvinceCode: provinceCode,
		}

		tests := []*location.GetLocationQuery{query0, query1}
		for i, query := range tests {
			t.Run(fmt.Sprint(i), func(t *testing.T) {
				err := locationBus.Dispatch(ctx, query)
				require.NoError(t, err)
				assert.Nil(t, query.Result.Ward)

				district, province := query.Result.District, query.Result.Province
				require.NotNil(t, district)
				require.NotNil(t, province)
				assert.Equal(t, district.Code, districtCode)
				assert.Equal(t, district.Name, "Thành phố Sa Đéc")
				assert.Equal(t, province.Code, provinceCode)
				assert.Equal(t, province.Name, "Tỉnh Đồng Tháp")
			})
		}
	})

	t.Run("DistrictCode (conflicted)", func(t *testing.T) {
		districtCode := "867" // Thành phố Sa Đéc
		query := &location.GetLocationQuery{
			DistrictCode: districtCode,
			ProvinceCode: "86", // conflict
		}
		err := locationBus.Dispatch(ctx, query)
		assert.EqualError(t, err, "mã tỉnh/thành phố không thống nhất")
	})

	t.Run("WardCode", func(t *testing.T) {
		wardCode := "10180"   // Xã Hồng Dương
		districtCode := "278" // Huyện Thanh Oai
		provinceCode := "01"  // Thành Phố Hà Nội
		query0 := &location.GetLocationQuery{WardCode: wardCode}
		query1 := &location.GetLocationQuery{
			WardCode:     wardCode,
			DistrictCode: districtCode,
		}
		query2 := &location.GetLocationQuery{
			WardCode:     wardCode,
			ProvinceCode: provinceCode,
		}
		query3 := &location.GetLocationQuery{
			WardCode:     wardCode,
			DistrictCode: districtCode,
			ProvinceCode: provinceCode,
		}
		tests := []*location.GetLocationQuery{query0, query1, query2, query3}
		for i, query := range tests {
			t.Run(fmt.Sprint(i), func(t *testing.T) {
				err := locationBus.Dispatch(ctx, query)
				require.NoError(t, err)

				ward := query.Result.Ward
				district := query.Result.District
				province := query.Result.Province
				require.NotNil(t, ward)
				require.NotNil(t, district)
				require.NotNil(t, province)
				assert.Equal(t, ward.Code, wardCode)
				assert.Equal(t, ward.Name, "Xã Hồng Dương")
				assert.Equal(t, district.Code, districtCode)
				assert.Equal(t, district.Name, "Huyện Thanh Oai")
				assert.Equal(t, province.Code, provinceCode)
				assert.Equal(t, province.Name, "Thành phố Hà Nội")
			})
		}
	})

	t.Run("WardCode (conflicted)", func(t *testing.T) {
		wardCode := "10180" // Xã Hồng Dương
		query := &location.GetLocationQuery{
			WardCode:     wardCode,
			DistrictCode: "279",
		}
		err := locationBus.Dispatch(ctx, query)
		assert.EqualError(t, err, "mã quận/huyện không thống nhất")
	})
}
