package xshop

import (
	"context"

	"o.o/backend/pkg/etop/apix/shopping"
)

func (s *VariantService) GetVariant(ctx context.Context, r *GetVariantEndpoint) error {
	resp, err := shopping.GetVariant(ctx, r.Context.Shop.ID, r.GetVariantRequest)
	r.Result = resp
	return err
}

func (s *VariantService) ListVariants(ctx context.Context, r *ListVariantsEndpoint) error {
	resp, err := shopping.ListVariants(ctx, r.Context.Shop.ID, r.ListVariantsRequest)
	r.Result = resp
	return err
}

func (s *VariantService) CreateVariant(ctx context.Context, r *CreateVariantEndpoint) error {
	resp, err := shopping.CreateVariant(ctx, r.Context.Shop.ID, r.Context.AuthPartnerID, r.CreateVariantRequest)
	r.Result = resp
	return err
}

func (s *VariantService) UpdateVariant(ctx context.Context, r *UpdateVariantEndpoint) error {
	resp, err := shopping.UpdateVariant(ctx, r.Context.Shop.ID, r.UpdateVariantRequest)
	r.Result = resp
	return err
}

func (s *VariantService) DeleteVariant(ctx context.Context, r *DeleteVariantEndpoint) error {
	resp, err := shopping.DeleteVariant(ctx, r.Context.Shop.ID, r.GetVariantRequest)
	r.Result = resp
	return err
}
