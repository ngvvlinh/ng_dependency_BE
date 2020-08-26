package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := [][2]string{
		{
			"123 phường 7 quận 5 TP HCM",
			"Thành phố Hồ Chí Minh|Quận 5|Phường 07|123",
		},
		{
			"123 p 7 q 5 HCM",
			"Thành phố Hồ Chí Minh|Quận 5|Phường 07|123",
		},
		{
			"123 p7 q5 HCM",
			"Thành phố Hồ Chí Minh|Quận 5|Phường 07|123",
		},
		{
			"123 p7 q5 tphcm q5 tphcm",
			"Thành phố Hồ Chí Minh|Quận 5|Phường 07|123",
		},
		{
			"123 p7 tanbinh HCM",
			"Thành phố Hồ Chí Minh|Quận Tân Bình|Phường 07|123",
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			best, _ := ParseAddress(tt[0])
			assert.Equal(t, tt[1], best.String())
		})
	}
}
