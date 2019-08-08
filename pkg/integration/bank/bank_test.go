package bank

import (
	"context"
	"testing"

	"etop.vn/api/main/location"
	servicelocation "etop.vn/backend/com/main/location"
	locationlist "etop.vn/backend/com/main/location/list"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var locationBus = servicelocation.New().MessageBus()

func TestProvinceCode(t *testing.T) {
	mapProvinces := make(map[string]bool)

	t.Run("Province code must be valid", func(t *testing.T) {
		for _, p := range Provinces {
			// Special cases
			province := p.TenTinhThanh
			if province == "Vung Tau" {
				province = "Ba Ria Vung Tau"
			} else if province == "TCTD VN o nuoc ngoai" {
				continue
			}

			query := &location.FindLocationQuery{Province: province}
			err := locationBus.Dispatch(context.Background(), query)
			require.NoError(t, err)
			require.Equal(t, query.Result.Province.Code, p.MaTinh)
			mapProvinces[p.MaTinh] = true
		}
	})
	t.Run("All provinces must have banks", func(t *testing.T) {
		for _, p := range locationlist.Provinces {
			assert.True(t, mapProvinces[p.Code], "No bank in province %v", p.Name)
		}
	})
}
