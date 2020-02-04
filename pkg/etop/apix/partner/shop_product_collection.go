package partner

import (
	"context"

	"etop.vn/backend/pkg/etop/apix/shopping"
)

func (s *ProductCollectionService) GetCollection(ctx context.Context, r *GetCollectionEndpoint) error {
	resp, err := shopping.GetCollection(ctx, r.Context.Shop.ID, r.GetCollectionRequest)
	r.Result = resp
	return err
}

func (s *ProductCollectionService) ListCollections(ctx context.Context, r *ListCollectionsEndpoint) error {
	resp, err := shopping.ListCollections(ctx, r.Context.Shop.ID, r.ListCollectionsRequest)
	r.Result = resp
	return err
}

func (s *ProductCollectionService) CreateCollection(ctx context.Context, r *CreateCollectionEndpoint) error {
	resp, err := shopping.CreateCollection(ctx, r.Context.Shop.ID, r.Context.AuthPartnerID, r.CreateCollectionRequest)
	r.Result = resp
	return err
}

func (s *ProductCollectionService) UpdateCollection(ctx context.Context, r *UpdateCollectionEndpoint) error {
	resp, err := shopping.UpdateCollection(ctx, r.Context.Shop.ID, r.UpdateCollectionRequest)
	r.Result = resp
	return err
}

func (s *ProductCollectionService) DeleteCollection(ctx context.Context, r *DeleteCollectionEndpoint) error {
	resp, err := shopping.DeleteCollection(ctx, r.Context.Shop.ID, r.GetCollectionRequest)
	r.Result = resp
	return err
}

func (s *ProductCollectionRelationshipService) ListRelationships(ctx context.Context, r *ProductCollectionListRelationshipsEndpoint) error {
	resp, err := shopping.ListRelationshipsProductCollection(ctx, r.Context.Shop.ID, r.ListProductCollectionRelationshipsRequest)
	r.Result = resp
	return err
}

func (s *ProductCollectionRelationshipService) CreateRelationship(ctx context.Context, r *ProductCollectionCreateRelationshipEndpoint) error {
	resp, err := shopping.CreateRelationshipProductCollection(ctx, r.Context.Shop.ID, r.CreateProductCollectionRelationshipRequest)
	r.Result = resp
	return err
}

func (s *ProductCollectionRelationshipService) DeleteRelationship(ctx context.Context, r *ProductCollectionDeleteRelationshipEndpoint) error {
	resp, err := shopping.DeleteRelationshipProductCollection(ctx, r.Context.Shop.ID, r.RemoveProductCollectionRequest)
	r.Result = resp
	return err
}
