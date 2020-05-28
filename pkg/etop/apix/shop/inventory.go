package xshop

import (
	"context"

	"o.o/backend/pkg/etop/apix/shopping"
)

type InventoryService struct {
	Shopping *shopping.Shopping
}

func (s *InventoryService) Clone() *InventoryService { res := *s; return &res }

func (s *InventoryService) ListInventoryLevels(ctx context.Context, r *ListInventoryLevelsEndpoint) error {
	resp, err := s.Shopping.ListInventoryLevels(ctx, r.Context.Shop.ID, r.ListInventoryLevelsRequest)
	r.Result = resp
	return err
}
