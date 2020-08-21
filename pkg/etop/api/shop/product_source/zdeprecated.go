package product_source

import (
	"context"

	api "o.o/api/top/int/shop"
	pbcm "o.o/api/top/types/common"
)

// deprecated
func (s *ProductSourceService) CreateProductSource(ctx context.Context, q *api.CreateProductSourceRequest) (*api.ProductSource, error) {
	result := &api.ProductSource{
		Id:     s.SS.Shop().ID,
		Status: 1,
	}
	return result, nil
}

// deprecated: 2018.07.24+14
func (s *ProductSourceService) GetShopProductSources(ctx context.Context, q *pbcm.Empty) (*api.ProductSourcesResponse, error) {
	result := &api.ProductSourcesResponse{
		ProductSources: []*api.ProductSource{
			{
				Id:     s.SS.Shop().ID,
				Status: 1,
			},
		},
	}
	return result, nil
}
