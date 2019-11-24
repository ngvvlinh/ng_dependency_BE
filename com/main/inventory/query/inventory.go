package query

import (
	"context"

	"etop.vn/api/main/etop"
	"etop.vn/api/main/inventory"
	"etop.vn/backend/com/main/inventory/convert"
	"etop.vn/backend/com/main/inventory/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/capi"
	"etop.vn/capi/dot"
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

func (q *InventoryQueryService) GetInventoryVariants(ctx context.Context, args *inventory.GetInventoryRequest) (*inventory.GetInventoryVariantsResponse, error) {
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
	return &inventory.GetInventoryVariantsResponse{InventoryVariants: convert.InventoryVariantsFromModel(result)}, nil
}

func (q *InventoryQueryService) GetInventoryVouchers(ctx context.Context, args *inventory.ListInventoryVouchersArgs) (*inventory.GetInventoryVouchersResponse, error) {
	if args.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing value requirement")
	}
	result, err := q.InventoryVoucherStore(ctx).ShopID(args.ShopID).Filters(args.Filters).Paging(args.Paging).ListInventoryVoucherDB()
	if err != nil {
		return nil, err
	}
	var inventoryVoucherItems []*inventory.InventoryVoucher
	for _, value := range result {
		inventoryVoucherItems = append(inventoryVoucherItems, convert.Convert_inventorymodel_InventoryVoucher_inventory_InventoryVoucher(value, nil))
	}
	return &inventory.GetInventoryVouchersResponse{InventoryVoucher: inventoryVoucherItems}, nil
}

func (q *InventoryQueryService) GetInventoryVariantsByVariantIDs(ctx context.Context, args *inventory.GetInventoryVariantsByVariantIDsArgs) (*inventory.GetInventoryVariantsResponse, error) {
	if args.ShopID == 0 || args.VariantIDs == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing value requirement")
	}
	result, err := q.InventoryStore(ctx).ShopID(args.ShopID).VariantIDs(args.VariantIDs...).ListInventoryDB()
	if err != nil {
		return nil, err
	}
	return &inventory.GetInventoryVariantsResponse{InventoryVariants: convert.InventoryVariantsFromModel(result)}, nil
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

func (q *InventoryQueryService) GetInventoryVoucher(ctx context.Context, ShopID dot.ID, ID dot.ID) (*inventory.InventoryVoucher, error) {
	result, err := q.InventoryVoucherStore(ctx).ID(ID).ShopID(ShopID).Get()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (q *InventoryQueryService) GetInventoryVariant(ctx context.Context, ShopID dot.ID, VariantID dot.ID) (*inventory.InventoryVariant, error) {
	result, err := q.InventoryStore(ctx).VariantIDs(VariantID).ShopID(ShopID).Get()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (q *InventoryQueryService) GetInventoryVouchersByRefIDs(
	ctx context.Context, RefIDs []dot.ID, ShopID dot.ID,
) (*inventory.GetInventoryVouchersResponse, error) {
	inventoryVouchers, err := q.InventoryVoucherStore(ctx).RefIDs(RefIDs...).ShopID(ShopID).ListInventoryVoucher()
	if err != nil {
		return nil, err
	}
	return &inventory.GetInventoryVouchersResponse{InventoryVoucher: inventoryVouchers}, nil
}

func (q *InventoryQueryService) GetInventoryVoucherByReference(ctx context.Context, ShopID dot.ID, refID dot.ID, refType inventory.InventoryRefType) (*inventory.GetInventoryVoucherByReferenceResponse, error) {
	result, err := q.InventoryVoucherStore(ctx).ShopID(ShopID).RefID(refID).RefType(string(refType)).ListInventoryVoucher()
	if err != nil {
		return nil, err
	}
	var status etop.Status4 = etop.S4Negative
	if len(result) == 0 {
		status = etop.S4Zero
	} else {
		for _, value := range result {
			if value.Status == etop.S3Zero {
				status = etop.S4SuperPos
				break
			}
			if value.Status == etop.S3Positive {
				status = etop.S4Positive
				break
			}
		}
	}
	return &inventory.GetInventoryVoucherByReferenceResponse{
		InventoryVouchers: result,
		Status:            status,
	}, nil
}
