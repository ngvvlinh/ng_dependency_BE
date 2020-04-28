package bank

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"o.o/api/main/location"
	servicelocation "o.o/backend/com/main/location"
	locationlist "o.o/backend/com/main/location/list"
)

var locationBus = servicelocation.New(nil).MessageBus()

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
