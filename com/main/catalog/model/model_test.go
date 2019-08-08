package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeAttributes(t *testing.T) {
	tests := []struct {
		attrs  ProductAttributes
		expect string
	}{
		{
			attrs:  ProductAttributes{},
			expect: "_",
		},
		{
			attrs: ProductAttributes{
				{"Màu sắc", "XANH"},
				{"Màu sắc", "Đỏ"},
				{"Kích thước", ""},
				{"", "Lớn"},
				{"Kích thước", "Nhỏ"},
			},
			expect: "mau_sac=xanh mau_sac=do kich_thuoc=nho",
		},
	}
	for _, tt := range tests {
		t.Run(tt.expect, func(t *testing.T) {
			_, norm := NormalizeAttributes(tt.attrs)
			assert.Equal(t, tt.expect, norm)
		})
	}
}
