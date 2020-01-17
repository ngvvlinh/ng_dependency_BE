package partner

import (
	"context"

	"etop.vn/api/main/catalog"
	externaltypes "etop.vn/api/top/external/types"
	"etop.vn/api/top/types/common"
	"etop.vn/backend/pkg/common/apifw/cmapi"
	"etop.vn/backend/pkg/etop/apix/convertpb"
	"etop.vn/capi/dot"
)

func (s *ProductCollectionService) GetCollection(ctx context.Context, r *GetCollectionEndpoint) error {
	query := &catalog.GetShopCollectionQuery{
		ID:     r.ID,
		ShopID: r.Context.Shop.ID,
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = convertpb.PbShopProductCollection(query.Result)
	return nil
}

func (s *ProductCollectionService) ListCollections(ctx context.Context, r *ListCollectionsEndpoint) error {
	paging, err := cmapi.CMCursorPaging(r.Paging)
	if err != nil {
		return err
	}
	query := &catalog.ListShopCollectionsByIDsQuery{
		IDs:    r.Filter.ID,
		ShopID: r.Context.Shop.ID,
		Paging: *paging,
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = &externaltypes.ProductCollectionsResponse{
		Collections: convertpb.PbShopProductCollections(query.Result.Collections),
		Paging:      convertpb.PbPageInfo(paging, &query.Result.Paging),
	}
	return nil
}

func (s *ProductCollectionService) CreateCollection(ctx context.Context, r *CreateCollectionEndpoint) error {
	cmd := &catalog.CreateShopCollectionCommand{
		ShopID:      r.Context.Shop.ID,
		Name:        r.Name,
		Description: r.Description,
		ShortDesc:   r.ShortDesc,
	}
	if err := catalogAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = convertpb.PbShopProductCollection(cmd.Result)
	return nil
}

func (s *ProductCollectionService) UpdateCollection(ctx context.Context, r *UpdateCollectionEndpoint) error {
	panic("TODO")
}

func (s *ProductCollectionService) DeleteCollection(ctx context.Context, r *DeleteCollectionEndpoint) error {
	cmd := &catalog.DeleteShopCollectionCommand{
		Id:     r.ID,
		ShopId: r.Context.Shop.ID,
	}
	if err := catalogAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &common.Empty{}
	return nil
}

func (s *ProductCollectionRelationshipService) ListRelationships(ctx context.Context, r *ProductCollectionListRelationshipsEndpoint) error {
	panic("TODO")
}

func (s *ProductCollectionRelationshipService) CreateRelationship(ctx context.Context, r *ProductCollectionCreateRelationshipEndpoint) error {
	cmd := &catalog.AddShopProductCollectionCommand{
		ProductID:     r.ProductId,
		ShopID:        r.Context.Shop.ID,
		CollectionIDs: []dot.ID{r.CollectionId},
	}
	if err := catalogAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &common.Empty{}
	return nil
}

func (s *ProductCollectionRelationshipService) DeleteRelationship(ctx context.Context, r *ProductCollectionDeleteRelationshipEndpoint) error {
	cmd := &catalog.RemoveShopProductCollectionCommand{
		ProductID:     r.ProductId,
		ShopID:        r.Context.Shop.ID,
		CollectionIDs: []dot.ID{r.CollectionId},
	}
	if err := catalogAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &common.Empty{}
	return nil
}
