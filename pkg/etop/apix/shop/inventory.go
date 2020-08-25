package xshop

import (
	"context"

	api "o.o/api/top/external/shop"
	externaltypes "o.o/api/top/external/types"
	"o.o/backend/pkg/etop/apix/shopping"
	"o.o/backend/pkg/etop/authorize/session"
)

type InventoryService struct {
	session.Session

	Shopping *shopping.Shopping
}

func (s *InventoryService) Clone() api.InventoryService { res := *s; return &res }

func (s *InventoryService) ListInventoryLevels(ctx context.Context, r *externaltypes.ListInventoryLevelsRequest) (*externaltypes.InventoryLevelsResponse, error) {
	resp, err := s.Shopping.ListInventoryLevels(ctx, s.SS.Shop().ID, r)
	return resp, err
}
