package partner

import (
	"context"

	"etop.vn/api/main/inventory"
	"etop.vn/api/top/external/types"
	"etop.vn/backend/pkg/etop/apix/convertpb"
)

func (s *InventoryService) ListInventoryLevels(ctx context.Context, r *ListInventoryLevelsEndpoint) error {
	query := &inventory.ListInventoryVariantsByVariantIDsQuery{
		ShopID:     r.Context.Shop.ID,
		VariantIDs: r.VariantIds,
	}
	if err := inventoryQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = &types.InventoryLevelsResponse{InventoryLevels: convertpb.PbInventoryLevels(query.Result.InventoryVariants)}
	return nil
}
