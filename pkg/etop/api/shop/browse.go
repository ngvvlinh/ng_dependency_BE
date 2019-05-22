package shop

import (
	"context"

	pbcm "etop.vn/backend/pb/common"
	pbetop "etop.vn/backend/pb/etop"
	pbshop "etop.vn/backend/pb/etop/shop"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/model"
	catalogmodelx "etop.vn/backend/pkg/services/catalog/modelx"
	wrapshop "etop.vn/backend/wrapper/etop/shop"
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

func BrowseCategories(ctx context.Context, q *wrapshop.BrowseCategoriesEndpoint) error {
	query := &model.GetEtopCategoriesQuery{Status: model.S3Positive.P()}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pbetop.CategoriesResponse{
		Categories: pbetop.PbCategories(query.Result.Categories),
	}
	return nil
}

func BrowseProduct(ctx context.Context, q *wrapshop.BrowseProductEndpoint) error {
	query := &catalogmodelx.GetProductQuery{
		ProductID: q.Id,
	}
	query.Status = model.S3Positive.P()
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = PbEtopProduct(query.Result)
	return nil
}

func BrowseProductsByIDs(ctx context.Context, q *wrapshop.BrowseProductsByIDsEndpoint) error {
	query := &catalogmodelx.GetProductsExtendedQuery{
		IDs:               q.Ids,
		ProductSourceType: model.ProductSourceKiotViet,
	}
	query.Status = model.S3Positive.P()
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = &pbshop.EtopProductsResponse{
		Products: PbEtopProducts(query.Result.Products),
	}
	return nil
}

func BrowseProducts(ctx context.Context, q *wrapshop.BrowseProductsEndpoint) error {
	query := &catalogmodelx.GetProductsExtendedQuery{
		Paging:            q.Paging.CMPaging(),
		Filters:           pbcm.ToFilters(q.Filters),
		ProductSourceType: model.ProductSourceKiotViet,
	}
	query.Status = model.S3Positive.P()
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = &pbshop.EtopProductsResponse{
		Products: PbEtopProducts(query.Result.Products),
		Paging:   pbcm.PbPageInfo(query.Paging, query.Result.Total),
	}
	return nil
}

func BrowseVariant(ctx context.Context, q *wrapshop.BrowseVariantEndpoint) error {
	query := &catalogmodelx.GetVariantQuery{
		VariantID: q.Id,
	}
	query.Status = model.S3Positive.P()
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = PbEtopVariant(query.Result)
	return nil
}

func BrowseVariantsByIDs(ctx context.Context, q *wrapshop.BrowseVariantsByIDsEndpoint) error {
	query := &catalogmodelx.GetVariantsExtendedQuery{
		IDs: q.Ids,
	}
	query.Status = model.S3Positive.P()
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = &pbshop.EtopVariantsResponse{
		Variants: PbEtopVariants(query.Result.Variants),
	}
	return nil
}

func BrowseVariants(ctx context.Context, q *wrapshop.BrowseVariantsEndpoint) error {
	query := &catalogmodelx.GetVariantsExtendedQuery{
		Paging:  q.Paging.CMPaging(),
		Filters: pbcm.ToFilters(q.Filters),
	}
	query.Status = model.S3Positive.P()
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = &pbshop.EtopVariantsResponse{
		Variants: PbEtopVariants(query.Result.Variants),
		Paging:   pbcm.PbPageInfo(query.Paging, query.Result.Total),
	}
	return nil
}
