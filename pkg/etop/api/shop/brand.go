package shop

import (
	"context"

	"o.o/api/main/catalog"
	"o.o/api/meta"
	"o.o/api/top/int/shop"
	"o.o/backend/pkg/common/apifw/cmapi"
)

type BrandService struct {
	CatalogQuery catalog.QueryBus
	CatalogAggr  catalog.CommandBus
}

func (s *BrandService) Clone() *BrandService { res := *s; return &res }

func (s *BrandService) GetBrandByID(ctx context.Context, q *GetBrandByIDEndpoint) error {
	shopID := q.Context.Shop.ID
	query := &catalog.GetBrandByIDQuery{
		Id:     q.Id,
		ShopID: shopID,
	}
	if err := s.CatalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = PbBrand(query.Result)
	return nil
}

func (s *BrandService) GetBrandsByIDs(ctx context.Context, q *GetBrandsByIDsEndpoint) error {
	shopID := q.Context.Shop.ID
	query := &catalog.GetBrandsByIDsQuery{
		Ids:    q.Ids,
		ShopID: shopID,
	}
	if err := s.CatalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &shop.GetBrandsByIDsResponse{
		Brands: PbBrands(query.Result),
	}
	return nil
}

func (s *BrandService) GetBrands(ctx context.Context, q *GetBrandsEndpoint) error {
	shopID := q.Context.Shop.ID
	query := &catalog.ListBrandsQuery{
		Paging: meta.Paging{
			Offset: q.Paging.Offset,
			Limit:  q.Paging.Limit,
		},
		ShopId: shopID,
	}
	if err := s.CatalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &shop.GetBrandsResponse{
		Brands: PbBrands(query.Result.ShopBrands),
		Paging: cmapi.PbPaging(query.Paging),
	}
	return nil
}

func (s *BrandService) CreateBrand(ctx context.Context, q *CreateBrandEndpoint) error {
	shopID := q.Context.Shop.ID
	cmd := &catalog.CreateBrandCommand{
		ShopID:      shopID,
		BrandName:   q.Name,
		Description: q.Description,
	}
	if err := s.CatalogAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = PbBrand(cmd.Result)
	return nil
}

func (s *BrandService) UpdateBrandInfo(ctx context.Context, q *UpdateBrandInfoEndpoint) error {
	shopID := q.Context.Shop.ID
	cmd := &catalog.UpdateBrandInfoCommand{
		ShopID:      shopID,
		ID:          q.Id,
		BrandName:   q.Name,
		Description: q.Description,
	}
	if err := s.CatalogAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = PbBrand(cmd.Result)
	return nil
}

func (s *BrandService) DeleteBrand(ctx context.Context, q *DeleteBrandEndpoint) error {
	shopID := q.Context.Shop.ID
	cmd := &catalog.DeleteShopBrandCommand{
		ShopId: shopID,
		Ids:    q.Ids,
	}
	if err := s.CatalogAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &shop.DeleteBrandResponse{
		Count: cmd.Result,
	}
	return nil
}
