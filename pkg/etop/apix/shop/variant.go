package xshop

import (
	"context"

	"o.o/backend/pkg/etop/apix/shopping"
)

type VariantService struct {
	Shopping *shopping.Shopping
}

func (s *VariantService) Clone() *VariantService { res := *s; return &res }

func (s *VariantService) GetVariant(ctx context.Context, r *GetVariantEndpoint) error {
	resp, err := s.Shopping.GetVariant(ctx, r.Context.Shop.ID, r.GetVariantRequest)
	r.Result = resp
	return err
}

func (s *VariantService) ListVariants(ctx context.Context, r *ListVariantsEndpoint) error {
	resp, err := s.Shopping.ListVariants(ctx, r.Context.Shop.ID, r.ListVariantsRequest)
	r.Result = resp
	return err
}

func (s *VariantService) CreateVariant(ctx context.Context, r *CreateVariantEndpoint) error {
	resp, err := s.Shopping.CreateVariant(ctx, r.Context.Shop.ID, r.Context.AuthPartnerID, r.CreateVariantRequest)
	r.Result = resp
	return err
}

func (s *VariantService) UpdateVariant(ctx context.Context, r *UpdateVariantEndpoint) error {
	resp, err := s.Shopping.UpdateVariant(ctx, r.Context.Shop.ID, r.UpdateVariantRequest)
	r.Result = resp
	return err
}

func (s *VariantService) DeleteVariant(ctx context.Context, r *DeleteVariantEndpoint) error {
	resp, err := s.Shopping.DeleteVariant(ctx, r.Context.Shop.ID, r.GetVariantRequest)
	r.Result = resp
	return err
}
