package shopping

import (
	"context"

	"etop.vn/api/main/inventory"
	"etop.vn/api/top/external/types"
	externaltypes "etop.vn/api/top/external/types"
	"etop.vn/backend/pkg/etop/apix/convertpb"
	"etop.vn/capi/dot"
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
