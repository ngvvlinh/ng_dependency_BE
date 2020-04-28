// Code generated by `go run ./scripts/ghn_location/main.go` DO NOT EDIT.

// Total VTPost Province not found: 0/63
// Total VTPost: 0/63

package list

import (
	"o.o/api/main/location"
	"o.o/backend/com/main/location/types"
)

var Countries = []*types.Country{{Name: types.CountryVietnam}}

var Provinces = []*types.Province{
	{
		Code:        "01",
		Region:      location.North,
		Name:        "Thành phố Hà Nội",
		VTPostID:    1,
		HaravanCode: "HI",
		Special:     true,
		Alias:       []string{"HN"},
	},
	{
		Code:        "02",
		Region:      location.North,
		Name:        "Tỉnh Hà Giang",
		VTPostID:    18,
		HaravanCode: "HG",
	},
	{
		Code:        "04",
		Region:      location.North,
		Name:        "Tỉnh Cao Bằng",
		VTPostID:    19,
		HaravanCode: "CB",
	},
	{
		Code:        "06",
		Region:      location.North,
		Name:        "Tỉnh Bắc Kạn",
		VTPostID:    21,
		HaravanCode: "BK",
		Alias:       []string{"Bắc Cạn"},
	},
	{
		Code:        "08",
		Region:      location.North,
		Name:        "Tỉnh Tuyên Quang",
		VTPostID:    23,
		HaravanCode: "TQ",
	},
	{
		Code:        "10",
		Region:      location.North,
		Name:        "Tỉnh Lào Cai",
		VTPostID:    20,
		HaravanCode: "LO",
	},
	{
		Code:        "11",
		Region:      location.North,
		Name:        "Tỉnh Điện Biên",
		VTPostID:    64,
		HaravanCode: "DB",
	},
	{
		Code:        "12",
		Region:      location.North,
		Name:        "Tỉnh Lai Châu",
		VTPostID:    29,
		HaravanCode: "LI",
	},
	{
		Code:        "14",
		Region:      location.North,
		Name:        "Tỉnh Sơn La",
		VTPostID:    30,
		HaravanCode: "SL",
	},
	{
		Code:        "15",
		Region:      location.North,
		Name:        "Tỉnh Yên Bái",
		VTPostID:    24,
		HaravanCode: "YB",
	},
	{
		Code:        "17",
		Region:      location.North,
		Name:        "Tỉnh Hoà Bình",
		VTPostID:    31,
		HaravanCode: "HO",
	},
	{
		Code:        "19",
		Region:      location.North,
		Name:        "Tỉnh Thái Nguyên",
		VTPostID:    25,
		HaravanCode: "TY",
	},
	{
		Code:        "20",
		Region:      location.North,
		Name:        "Tỉnh Lạng Sơn",
		VTPostID:    22,
		HaravanCode: "LS",
	},
	{
		Code:        "22",
		Region:      location.North,
		Name:        "Tỉnh Quảng Ninh",
		VTPostID:    28,
		HaravanCode: "QN",
	},
	{
		Code:        "24",
		Region:      location.North,
		Name:        "Tỉnh Bắc Giang",
		VTPostID:    27,
		HaravanCode: "BG",
	},
	{
		Code:        "25",
		Region:      location.North,
		Name:        "Tỉnh Phú Thọ",
		VTPostID:    26,
		HaravanCode: "PT",
	},
	{
		Code:        "26",
		Region:      location.North,
		Name:        "Tỉnh Vĩnh Phúc",
		VTPostID:    10,
		HaravanCode: "VT",
	},
	{
		Code:        "27",
		Region:      location.North,
		Name:        "Tỉnh Bắc Ninh",
		VTPostID:    11,
		HaravanCode: "BN",
	},
	{
		Code:        "30",
		Region:      location.North,
		Name:        "Tỉnh Hải Dương",
		VTPostID:    12,
		HaravanCode: "HD",
	},
	{
		Code:        "31",
		Region:      location.North,
		Name:        "Thành phố Hải Phòng",
		VTPostID:    3,
		HaravanCode: "HP",
	},
	{
		Code:        "33",
		Region:      location.North,
		Name:        "Tỉnh Hưng Yên",
		VTPostID:    13,
		HaravanCode: "HY",
	},
	{
		Code:        "34",
		Region:      location.North,
		Name:        "Tỉnh Thái Bình",
		VTPostID:    16,
		HaravanCode: "TB",
	},
	{
		Code:        "35",
		Region:      location.North,
		Name:        "Tỉnh Hà Nam",
		VTPostID:    14,
		HaravanCode: "HM",
	},
	{
		Code:        "36",
		Region:      location.North,
		Name:        "Tỉnh Nam Định",
		VTPostID:    15,
		HaravanCode: "ND",
	},
	{
		Code:        "37",
		Region:      location.North,
		Name:        "Tỉnh Ninh Bình",
		VTPostID:    17,
		HaravanCode: "NB",
	},
	{
		Code:        "38",
		Region:      location.North,
		Name:        "Tỉnh Thanh Hóa",
		VTPostID:    32,
		HaravanCode: "TH",
	},
	{
		Code:        "40",
		Region:      location.North,
		Name:        "Tỉnh Nghệ An",
		VTPostID:    33,
		HaravanCode: "NA",
	},
	{
		Code:        "42",
		Region:      location.North,
		Name:        "Tỉnh Hà Tĩnh",
		VTPostID:    34,
		HaravanCode: "HT",
	},
	{
		Code:        "44",
		Region:      location.Middle,
		Name:        "Tỉnh Quảng Bình",
		VTPostID:    35,
		HaravanCode: "QB",
	},
	{
		Code:        "45",
		Region:      location.Middle,
		Name:        "Tỉnh Quảng Trị",
		VTPostID:    36,
		HaravanCode: "QT",
	},
	{
		Code:        "46",
		Region:      location.Middle,
		Name:        "Tỉnh Thừa Thiên Huế",
		VTPostID:    37,
		HaravanCode: "TT",
		Alias:       []string{"Huế"},
	},
	{
		Code:        "48",
		Region:      location.Middle,
		Name:        "Thành phố Đà Nẵng",
		VTPostID:    4,
		HaravanCode: "DA",
		Special:     true,
	},
	{
		Code:        "49",
		Region:      location.Middle,
		Name:        "Tỉnh Quảng Nam",
		VTPostID:    38,
		HaravanCode: "QM",
	},
	{
		Code:        "51",
		Region:      location.Middle,
		Name:        "Tỉnh Quảng Ngãi",
		VTPostID:    39,
		HaravanCode: "QG",
	},
	{
		Code:        "52",
		Region:      location.Middle,
		Name:        "Tỉnh Bình Định",
		VTPostID:    40,
		HaravanCode: "BD",
	},
	{
		Code:        "54",
		Region:      location.Middle,
		Name:        "Tỉnh Phú Yên",
		VTPostID:    41,
		HaravanCode: "PY",
	},
	{
		Code:        "56",
		Region:      location.Middle,
		Name:        "Tỉnh Khánh Hòa",
		VTPostID:    42,
		HaravanCode: "KH",
	},
	{
		Code:        "62",
		Region:      location.Middle,
		Name:        "Tỉnh Kon Tum",
		VTPostID:    43,
		HaravanCode: "KT",
	},
	{
		Code:        "64",
		Region:      location.Middle,
		Name:        "Tỉnh Gia Lai",
		VTPostID:    44,
		HaravanCode: "GL",
	},
	{
		Code:        "66",
		Region:      location.Middle,
		Name:        "Tỉnh Đắk Lắk",
		VTPostID:    45,
		HaravanCode: "DC",
		Alias:       []string{"Daklak", "Tỉnh Đắc Lắk", "Tỉnh Đắk Lắc", "Tỉnh Đắc Lắc"},
	},
	{
		Code:        "67",
		Region:      location.Middle,
		Name:        "Tỉnh Đắk Nông",
		VTPostID:    9,
		HaravanCode: "DO",
		Alias:       []string{"Tỉnh Đắc Nông"},
	},
	{
		Code:        "58",
		Region:      location.South,
		Name:        "Tỉnh Ninh Thuận",
		VTPostID:    47,
		HaravanCode: "NT",
	},
	{
		Code:        "60",
		Region:      location.South,
		Name:        "Tỉnh Bình Thuận",
		VTPostID:    52,
		HaravanCode: "BU",
	},
	{
		Code:        "68",
		Region:      location.South,
		Name:        "Tỉnh Lâm Đồng",
		VTPostID:    46,
		HaravanCode: "LD",
	},
	{
		Code:        "70",
		Region:      location.South,
		Name:        "Tỉnh Bình Phước",
		VTPostID:    48,
		HaravanCode: "BP",
	},
	{
		Code:        "72",
		Region:      location.South,
		Name:        "Tỉnh Tây Ninh",
		VTPostID:    49,
		HaravanCode: "TN",
	},
	{
		Code:        "74",
		Region:      location.South,
		Name:        "Tỉnh Bình Dương",
		VTPostID:    50,
		HaravanCode: "BI",
	},
	{
		Code:        "75",
		Region:      location.South,
		Name:        "Tỉnh Đồng Nai",
		VTPostID:    51,
		HaravanCode: "DN",
	},
	{
		Code:        "77",
		Region:      location.South,
		Name:        "Tỉnh Bà Rịa - Vũng Tàu",
		VTPostID:    53,
		HaravanCode: "BV",
		Alias:       []string{"Vũng Tàu", "Bà Rịa", "BRVT", "BR-VT"},
	},
	{
		Code:        "79",
		Region:      location.South,
		Name:        "Thành phố Hồ Chí Minh",
		VTPostID:    2,
		HaravanCode: "HC",
		Special:     true,
		Alias:       []string{"HCM", "TPHCM"},
	},
	{
		Code:        "80",
		Region:      location.South,
		Name:        "Tỉnh Long An",
		VTPostID:    54,
		HaravanCode: "LA",
	},
	{
		Code:        "82",
		Region:      location.South,
		Name:        "Tỉnh Tiền Giang",
		VTPostID:    6,
		HaravanCode: "TG",
	},
	{
		Code:        "83",
		Region:      location.South,
		Name:        "Tỉnh Bến Tre",
		VTPostID:    58,
		HaravanCode: "BT",
	},
	{
		Code:        "84",
		Region:      location.South,
		Name:        "Tỉnh Trà Vinh",
		VTPostID:    60,
		HaravanCode: "TV",
	},
	{
		Code:        "86",
		Region:      location.South,
		Name:        "Tỉnh Vĩnh Long",
		VTPostID:    57,
		HaravanCode: "VL",
	},
	{
		Code:        "87",
		Region:      location.South,
		Name:        "Tỉnh Đồng Tháp",
		VTPostID:    55,
		HaravanCode: "DT",
	},
	{
		Code:        "89",
		Region:      location.South,
		Name:        "Tỉnh An Giang",
		VTPostID:    56,
		HaravanCode: "AG",
	},
	{
		Code:        "91",
		Region:      location.South,
		Name:        "Tỉnh Kiên Giang",
		VTPostID:    59,
		HaravanCode: "KG",
	},
	{
		Code:        "92",
		Region:      location.South,
		Name:        "Thành phố Cần Thơ",
		VTPostID:    5,
		HaravanCode: "CN",
	},
	{
		Code:        "93",
		Region:      location.South,
		Name:        "Tỉnh Hậu Giang",
		VTPostID:    8,
		HaravanCode: "HU",
	},
	{
		Code:        "94",
		Region:      location.South,
		Name:        "Tỉnh Sóc Trăng",
		VTPostID:    61,
		HaravanCode: "ST",
	},
	{
		Code:        "95",
		Region:      location.South,
		Name:        "Tỉnh Bạc Liêu",
		VTPostID:    62,
		HaravanCode: "BL",
	},
	{
		Code:        "96",
		Region:      location.South,
		Name:        "Tỉnh Cà Mau",
		VTPostID:    63,
		HaravanCode: "CM",
	},
}
