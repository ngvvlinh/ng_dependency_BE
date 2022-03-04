package bankstatement

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseBankStatementDescription(t *testing.T) {
	testcases := []struct {
		name       string
		desc       string
		want       *ShopInfo
		errMessage string
	}{
		{
			name:       "Empty description",
			desc:       "",
			want:       nil,
			errMessage: "Missing description",
		},
		{
			name:       "Wrong format",
			desc:       "Hello this is a wrong format description",
			want:       nil,
			errMessage: "Description does not match format",
		},
		{
			name:       "Wrong shop code format",
			desc:       "XYT2nn 0909090909",
			want:       nil,
			errMessage: "Description does not match format",
		},
		{
			name:       "Wrong phone number format",
			desc:       "XYT2 09090909",
			want:       nil,
			errMessage: "Description does not match format",
		},
		{
			name: "Correct format",
			desc: "XYT2 0909090909",
			want: &ShopInfo{ShopCode: "XYT2", UserPhone: "0909090909"},
		},
		{
			name: "Correct format with mord ending white space",
			desc: "XYT2 0909090909     ",
			want: &ShopInfo{ShopCode: "XYT2", UserPhone: "0909090909"},
		},
		{
			name: "Correct format with begining white space",
			desc: " XYT2 0909090909",
			want: &ShopInfo{ShopCode: "XYT2", UserPhone: "0909090909"},
		},
		{
			name: "Correct format with 84 phone number",
			desc: " XYT2 84909090909",
			want: &ShopInfo{ShopCode: "XYT2", UserPhone: "0909090909"},
		},
		{
			name: "Correct format with +84 phone number",
			desc: " XYT2 +84909090909",
			want: &ShopInfo{ShopCode: "XYT2", UserPhone: "0909090909"},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			res, err := ParseBankStatementDescription(tt.desc)
			if err != nil {
				assert.Equal(t, tt.errMessage, err.Error())
			} else {
				assert.Equal(t, tt.want, res)
			}
		})
	}
}
