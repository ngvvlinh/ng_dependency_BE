package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"o.o/api/main/catalog/types"
)

func TestNormalizeAttributes(t *testing.T) {
	tests := []struct {
		attrs  []*types.Attribute
		expect string
	}{
		{
			attrs:  []*types.Attribute{},
			expect: "_",
		},
		{
			attrs: []*types.Attribute{
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
