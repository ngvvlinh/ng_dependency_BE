package shop

import (
	"context"

	"o.o/api/main/catalog"
	"o.o/api/meta"
	api "o.o/api/top/int/shop"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/authorize/session"
)

type BrandService struct {
	session.Session

	CatalogQuery catalog.QueryBus
	CatalogAggr  catalog.CommandBus
}

func (s *BrandService) Clone() api.BrandService { res := *s; return &res }

func (s *BrandService) GetBrandByID(ctx context.Context, q *pbcm.IDRequest) (*api.Brand, error) {
	shopID := s.SS.Shop().ID
	query := &catalog.GetBrandByIDQuery{
		Id:     q.Id,
		ShopID: shopID,
	}
	if err := s.CatalogQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := PbBrand(query.Result)
	return result, nil
}

func (s *BrandService) GetBrandsByIDs(ctx context.Context, q *pbcm.IDsRequest) (*api.GetBrandsByIDsResponse, error) {
	shopID := s.SS.Shop().ID
	query := &catalog.GetBrandsByIDsQuery{
		Ids:    q.Ids,
		ShopID: shopID,
	}
	if err := s.CatalogQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &api.GetBrandsByIDsResponse{
		Brands: PbBrands(query.Result),
	}
	return result, nil
}

func (s *BrandService) GetBrands(ctx context.Context, q *api.GetBrandsRequest) (*api.GetBrandsResponse, error) {
	shopID := s.SS.Shop().ID
	query := &catalog.ListBrandsQuery{
		Paging: meta.Paging{
			Offset: q.Paging.Offset,
			Limit:  q.Paging.Limit,
		},
		ShopId: shopID,
	}
	if err := s.CatalogQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &api.GetBrandsResponse{
		Brands: PbBrands(query.Result.ShopBrands),
		Paging: cmapi.PbPaging(query.Paging),
	}
	return result, nil
}

func (s *BrandService) CreateBrand(ctx context.Context, q *api.CreateBrandRequest) (*api.Brand, error) {
	shopID := s.SS.Shop().ID
	cmd := &catalog.CreateBrandCommand{
		ShopID:      shopID,
		BrandName:   q.Name,
		Description: q.Description,
	}
	if err := s.CatalogAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := PbBrand(cmd.Result)
	return result, nil
}

func (s *BrandService) UpdateBrandInfo(ctx context.Context, q *api.UpdateBrandRequest) (*api.Brand, error) {
	shopID := s.SS.Shop().ID
	cmd := &catalog.UpdateBrandInfoCommand{
		ShopID:      shopID,
		ID:          q.Id,
		BrandName:   q.Name,
		Description: q.Description,
	}
	if err := s.CatalogAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := PbBrand(cmd.Result)
	return result, nil
}

func (s *BrandService) DeleteBrand(ctx context.Context, q *pbcm.IDsRequest) (*api.DeleteBrandResponse, error) {
	shopID := s.SS.Shop().ID
	cmd := &catalog.DeleteShopBrandCommand{
		ShopId: shopID,
		Ids:    q.Ids,
	}
	if err := s.CatalogAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &api.DeleteBrandResponse{
		Count: cmd.Result,
	}
	return result, nil
}
