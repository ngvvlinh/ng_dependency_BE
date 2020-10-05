package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCensoredWords(t *testing.T) {
	var tests = map[string]string{
		"su kien dien bien hoa binh tai vietnam": "su kien tai",
		"vuong dinh hue":                         "",
		"truong hoa binh":                        "",
		"giam han muc":                           "",
		"Nhà không có chó Khuyến Mãi":            "Nha khong co cho Khuyen M..",
		"quang cao hang mi pham":                 "hang mi pham",
		"hang mi pham duoc quang cao":            "hang mi pham duoc",
		"nha hang dang khuyen mai":               "nha hang dang khuyen m..",
		"cuc ki khung khiep khuyen mai abc":      "cuc ki khung khiep khuyen m.. abc",
		"khuyen mai abc":                         "khuyen m.. abc",
		"abc truong hoa binh khuyen mai abc":     "abc khuyen m.. abc",
	}
	count := 1
	for key, value := range tests {
		t.Run(fmt.Sprintf("util msg : %v", count), func(t *testing.T) {
			s := ModifyMsgPhone(key)
			assert.Equal(t, value, s)
		})
		count++
	}
}

func TestVietnameseToLatin(t *testing.T) {
	var tests = map[string]string{
		"Nam thích ăn vịt": "Nam thich an vit",
		"Nhà không có chó": "Nha khong co cho",
		"Công ty từng năm": "Cong ty tung nam",
		"Anh yeu em":       "Anh yeu em",
	}
	count := 1
	for key, value := range tests {
		t.Run(fmt.Sprintf("util msg : %v", count), func(t *testing.T) {
			s := removeAccent(key)
			assert.Equal(t, value, s)
		})
		count++
	}
}
