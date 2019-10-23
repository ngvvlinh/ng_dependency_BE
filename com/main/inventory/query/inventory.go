package query

import (
	"context"

	"etop.vn/api/main/inventory"
	"etop.vn/api/meta"
	"etop.vn/backend/com/main/inventory/convert"
	"etop.vn/backend/com/main/inventory/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/capi"
)

var _ inventory.QueryService = &InventoryQueryService{}

type InventoryQueryService struct {
	InventoryStore        sqlstore.InventoryFactory
	InventoryVoucherStore sqlstore.InventoryVoucherFactory
	EventBus              capi.EventBus
}

func NewQueryInventory(eventBus capi.EventBus, db *cmsql.Database) *InventoryQueryService {
	return &InventoryQueryService{
		InventoryStore:        sqlstore.NewInventoryStore(db),
		InventoryVoucherStore: sqlstore.NewInventoryVoucherStore(db),
		EventBus:              eventBus,
	}
}

func (q *InventoryQueryService) MessageBus() inventory.QueryBus {
	b := bus.New()
	return inventory.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *InventoryQueryService) GetInventories(ctx context.Context, args *inventory.GetInventoryRequest) (*inventory.GetInventoriesResponse, error) {
	if args.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing value requirement")
	}
	var page cm.Paging
	page.Limit = args.Paging.Limit
	page.Offset = args.Paging.Offset
	result, err := q.InventoryStore(ctx).ShopID(args.ShopID).Paging(&page).ListInventoryDB()
	if err != nil {
		return nil, err
	}
	return &inventory.GetInventoriesResponse{Inventories: convert.InventoryVariantsFromModel(result)}, nil
}

func (q *InventoryQueryService) GetInventoryVouchers(ctx context.Context, ShopID int64, Paging *meta.Paging) (*inventory.GetInventoryVouchersResponse, error) {
	if ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing value requirement")
	}
	var page cm.Paging
	page.Limit = Paging.Limit
	page.Offset = Paging.Offset
	result, err := q.InventoryVoucherStore(ctx).ShopID(ShopID).Paging(&page).ListInventoryVoucherDB()
	if err != nil {
		return nil, err
	}
	var inventoryVoucherItems []*inventory.InventoryVoucher
	for _, value := range result {
		inventoryVoucherItems = append(inventoryVoucherItems, convert.Convert_inventorymodel_InventoryVoucher_inventory_InventoryVoucher(value, nil))
	}
	return &inventory.GetInventoryVouchersResponse{InventoryVoucher: inventoryVoucherItems}, nil
}

func (q *InventoryQueryService) GetInventoriesByVariantIDs(ctx context.Context, args *inventory.GetInventoriesByVariantIDsArgs) (*inventory.GetInventoriesResponse, error) {
	if args.ShopID == 0 || args.VariantIDs == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing value requirement")
	}
	result, err := q.InventoryStore(ctx).ShopID(args.ShopID).VariantIDs(args.VariantIDs...).ListInventoryDB()
	if err != nil {
		return nil, err
	}
	return &inventory.GetInventoriesResponse{Inventories: convert.InventoryVariantsFromModel(result)}, nil
}

func (q *InventoryQueryService) GetInventoryVouchersByIDs(ctx context.Context, args *inventory.GetInventoryVouchersByIDArgs) (*inventory.GetInventoryVouchersResponse, error) {
	if args.ShopID == 0 || args.IDs == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing value requirement")
	}
	result, err := q.InventoryVoucherStore(ctx).ShopID(args.ShopID).IDs(args.IDs...).ListInventoryVoucherDB()
	if err != nil {
		return nil, err
	}
	var inventoryVoucherItems []*inventory.InventoryVoucher
	for _, value := range result {
		inventoryVoucherItems = append(inventoryVoucherItems, convert.Convert_inventorymodel_InventoryVoucher_inventory_InventoryVoucher(value, nil))
	}
	return &inventory.GetInventoryVouchersResponse{InventoryVoucher: inventoryVoucherItems}, nil
}

func (q *InventoryQueryService) GetInventoryVoucher(ctx context.Context, ShopID int64, ID int64) (*inventory.InventoryVoucher, error) {
	result, err := q.InventoryVoucherStore(ctx).ID(ID).ShopID(ShopID).Get()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (q *InventoryQueryService) GetInventory(ctx context.Context, ShopID int64, VariantID int64) (*inventory.InventoryVariant, error) {
	result, err := q.InventoryStore(ctx).VariantIDs(VariantID).ShopID(ShopID).Get()
	if err != nil {
		return nil, err
	}
	return result, nil
}
