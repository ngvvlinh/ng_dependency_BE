package partner

import (
	"context"

	api "o.o/api/top/external/partner"
	externaltypes "o.o/api/top/external/types"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/etop/apix/shopping"
	"o.o/backend/pkg/etop/authorize/session"
)

type VariantService struct {
	session.Session

	Shopping *shopping.Shopping
}

func (s *VariantService) Clone() api.VariantService { res := *s; return &res }

func (s *VariantService) GetVariant(ctx context.Context, r *externaltypes.GetVariantRequest) (*externaltypes.ShopVariant, error) {
	resp, err := s.Shopping.GetVariant(ctx, s.SS.Shop().ID, r)
	return resp, err
}

func (s *VariantService) ListVariants(ctx context.Context, r *externaltypes.ListVariantsRequest) (*externaltypes.ShopVariantsResponse, error) {
	resp, err := s.Shopping.ListVariants(ctx, s.SS.Shop().ID, r)
	return resp, err
}

func (s *VariantService) CreateVariant(ctx context.Context, r *externaltypes.CreateVariantRequest) (*externaltypes.ShopVariant, error) {
	resp, err := s.Shopping.CreateVariant(ctx, s.SS.Shop().ID, s.SS.Claim().AuthPartnerID, r)
	return resp, err
}

func (s *VariantService) UpdateVariant(ctx context.Context, r *externaltypes.UpdateVariantRequest) (*externaltypes.ShopVariant, error) {
	resp, err := s.Shopping.UpdateVariant(ctx, s.SS.Shop().ID, r)
	return resp, err
}

func (s *VariantService) DeleteVariant(ctx context.Context, r *externaltypes.GetVariantRequest) (*pbcm.Empty, error) {
	resp, err := s.Shopping.DeleteVariant(ctx, s.SS.Shop().ID, r)
	return resp, err
}
