package shop

import (
	"context"

	cmP "etop.vn/backend/pb/common"
	etopP "etop.vn/backend/pb/etop"
	shopP "etop.vn/backend/pb/etop/shop"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/model"
	shopW "etop.vn/backend/wrapper/etop/shop"
)

func init() {
	bus.AddHandler("api", BrowseCategories)
	bus.AddHandler("api", BrowseVariant)
	bus.AddHandler("api", BrowseVariantsByIDs)
	bus.AddHandler("api", BrowseVariants)
	bus.AddHandler("api", BrowseProduct)
	bus.AddHandler("api", BrowseProductsByIDs)
	bus.AddHandler("api", BrowseProducts)
}

func BrowseCategories(ctx context.Context, q *shopW.BrowseCategoriesEndpoint) error {
	query := &model.GetEtopCategoriesQuery{Status: model.S3Positive.P()}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &etopP.CategoriesResponse{
		Categories: etopP.PbCategories(query.Result.Categories),
	}
	return nil
}

func BrowseProduct(ctx context.Context, q *shopW.BrowseProductEndpoint) error {
	query := &model.GetProductQuery{
		ProductID: q.Id,
	}
	query.Status = model.S3Positive.P()
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = PbEtopProduct(query.Result)
	return nil
}

func BrowseProductsByIDs(ctx context.Context, q *shopW.BrowseProductsByIDsEndpoint) error {
	query := &model.GetProductsExtendedQuery{
		IDs:               q.Ids,
		ProductSourceType: model.ProductSourceKiotViet,
	}
	query.Status = model.S3Positive.P()
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = &shopP.EtopProductsResponse{
		Products: PbEtopProducts(query.Result.Products),
	}
	return nil
}

func BrowseProducts(ctx context.Context, q *shopW.BrowseProductsEndpoint) error {
	query := &model.GetProductsExtendedQuery{
		Paging:            q.Paging.CMPaging(),
		Filters:           cmP.ToFilters(q.Filters),
		ProductSourceType: model.ProductSourceKiotViet,
	}
	query.Status = model.S3Positive.P()
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = &shopP.EtopProductsResponse{
		Products: PbEtopProducts(query.Result.Products),
		Paging:   cmP.PbPageInfo(query.Paging, query.Result.Total),
	}
	return nil
}

func BrowseVariant(ctx context.Context, q *shopW.BrowseVariantEndpoint) error {
	query := &model.GetVariantQuery{
		VariantID: q.Id,
	}
	query.Status = model.S3Positive.P()
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = PbEtopVariant(query.Result)
	return nil
}

func BrowseVariantsByIDs(ctx context.Context, q *shopW.BrowseVariantsByIDsEndpoint) error {
	query := &model.GetVariantsExtendedQuery{
		IDs: q.Ids,
	}
	query.Status = model.S3Positive.P()
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = &shopP.EtopVariantsResponse{
		Variants: PbEtopVariants(query.Result.Variants),
	}
	return nil
}

func BrowseVariants(ctx context.Context, q *shopW.BrowseVariantsEndpoint) error {
	query := &model.GetVariantsExtendedQuery{
		Paging:  q.Paging.CMPaging(),
		Filters: cmP.ToFilters(q.Filters),
	}
	query.Status = model.S3Positive.P()
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = &shopP.EtopVariantsResponse{
		Variants: PbEtopVariants(query.Result.Variants),
		Paging:   cmP.PbPageInfo(query.Paging, query.Result.Total),
	}
	return nil
}
