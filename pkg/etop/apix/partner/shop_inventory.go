package partner

import (
	"context"

	"o.o/backend/pkg/etop/apix/shopping"
)

func (s *InventoryService) ListInventoryLevels(ctx context.Context, r *ListInventoryLevelsEndpoint) error {
	resp, err := shopping.ListInventoryLevels(ctx, r.Context.Shop.ID, r.ListInventoryLevelsRequest)
	r.Result = resp
	return err
}
