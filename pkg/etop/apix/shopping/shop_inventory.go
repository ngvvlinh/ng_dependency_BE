package shopping

import (
	"context"

	"o.o/api/main/inventory"
	"o.o/api/top/external/types"
	externaltypes "o.o/api/top/external/types"
	"o.o/backend/pkg/etop/apix/convertpb"
	"o.o/capi/dot"
)

func ListInventoryLevels(ctx context.Context, shopID dot.ID, request *externaltypes.ListInventoryLevelsRequest) (*externaltypes.InventoryLevelsResponse, error) {
	var IDs []dot.ID
	if len(request.Filter.VariantID) != 0 {
		IDs = request.Filter.VariantID
	}
	query := &inventory.ListInventoryVariantsByVariantIDsQuery{
		ShopID:     shopID,
		VariantIDs: IDs,
	}
	if err := inventoryQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	return &types.InventoryLevelsResponse{InventoryLevels: convertpb.PbInventoryLevels(query.Result.InventoryVariants)}, nil
}
