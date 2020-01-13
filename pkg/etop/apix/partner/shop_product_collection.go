package partner

import (
	"context"

	"etop.vn/api/main/catalog"
	"etop.vn/api/top/types/common"
	"etop.vn/capi/dot"
)

func (s *ProductCollectionService) GetCollection(ctx context.Context, r *GetCollectionEndpoint) error {
	panic("TODO")
}

func (s *ProductCollectionService) ListCollections(ctx context.Context, r *ListCollectionsEndpoint) error {
	panic("TODO")
}

func (s *ProductCollectionService) CreateCollection(ctx context.Context, r *CreateCollectionEndpoint) error {
	panic("TODO")
}

func (s *ProductCollectionService) UpdateCollection(ctx context.Context, r *UpdateCollectionEndpoint) error {
	panic("TODO")
}

func (s *ProductCollectionService) DeleteCollection(ctx context.Context, r *DeleteCollectionEndpoint) error {
	panic("TODO")
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
