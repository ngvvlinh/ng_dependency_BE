package query

import (
	"context"

	"o.o/api/main/inventory"
	"o.o/api/main/stocktaking"
	"o.o/api/top/types/etc/inventory_voucher_ref"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/status4"
	"o.o/api/top/types/etc/stocktake_type"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/inventory/convert"
	"o.o/backend/com/main/inventory/sqlstore"
	"o.o/backend/com/main/inventory/util"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/capi"
	"o.o/capi/dot"
)

var _ inventory.QueryService = &InventoryQueryService{}

type InventoryQueryService struct {
	InventoryStore        sqlstore.InventoryFactory
	InventoryVoucherStore sqlstore.InventoryVoucherFactory
	EventBus              capi.EventBus
	StocktakeQuery        stocktaking.QueryBus
}

func NewQueryInventory(stocktakeQuery stocktaking.QueryBus, eventBus capi.EventBus, db com.MainDB) *InventoryQueryService {
	return &InventoryQueryService{
		InventoryStore:        sqlstore.NewInventoryStore(db),
		InventoryVoucherStore: sqlstore.NewInventoryVoucherStore(db),
		EventBus:              eventBus,
		StocktakeQuery:        stocktakeQuery,
	}
}

func InventoryQueryServiceMessageBus(q *InventoryQueryService) inventory.QueryBus {
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
	result, err := q.InventoryStore(ctx).ShopID(args.ShopID).WithPaging(&page).ListInventory()
	if err != nil {
		return nil, err
	}
	return &inventory.GetInventoryVariantsResponse{InventoryVariants: result}, nil
}

func (q *InventoryQueryService) GetInventoryVouchers(ctx context.Context, args *inventory.ListInventoryVouchersArgs) (*inventory.GetInventoryVouchersResponse, error) {
	if args.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing value requirement")
	}
	result, err := q.InventoryVoucherStore(ctx).ShopID(args.ShopID).Filters(args.Filters).WithPaging(args.Paging).ListInventoryVoucherDB()
	if err != nil {
		return nil, err
	}
	var inventoryVoucherItems []*inventory.InventoryVoucher
	for _, value := range result {
		inventoryVoucherItems = append(inventoryVoucherItems, convert.Convert_inventorymodel_InventoryVoucher_inventory_InventoryVoucher(value, nil))
	}
	inventoryVoucherItems, err = q.populateInventoryVouchers(ctx, inventoryVoucherItems)
	if err != nil {
		return nil, err
	}
	return &inventory.GetInventoryVouchersResponse{InventoryVoucher: inventoryVoucherItems}, nil
}

func (q *InventoryQueryService) GetInventoryVariantsByVariantIDs(ctx context.Context, args *inventory.GetInventoryVariantsByVariantIDsArgs) (*inventory.GetInventoryVariantsResponse, error) {
	if args.ShopID == 0 || args.VariantIDs == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing value requirement")
	}
	result, err := q.InventoryStore(ctx).ShopID(args.ShopID).VariantIDs(args.VariantIDs...).ListInventory()
	if err != nil {
		return nil, err
	}
	return &inventory.GetInventoryVariantsResponse{InventoryVariants: result}, nil
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
	inventoryVoucherItems, err = q.populateInventoryVouchers(ctx, inventoryVoucherItems)
	if err != nil {
		return nil, err
	}
	return &inventory.GetInventoryVouchersResponse{InventoryVoucher: inventoryVoucherItems}, nil
}

func (q *InventoryQueryService) GetInventoryVoucher(ctx context.Context, ShopID dot.ID, ID dot.ID) (*inventory.InventoryVoucher, error) {
	result, err := q.InventoryVoucherStore(ctx).ID(ID).ShopID(ShopID).Get()
	if err != nil {
		return nil, err
	}
	return q.populateInventoryVoucher(ctx, result)
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
	inventoryVouchers, err = q.populateInventoryVouchers(ctx, inventoryVouchers)
	if err != nil {
		return nil, err
	}
	return &inventory.GetInventoryVouchersResponse{InventoryVoucher: inventoryVouchers}, nil
}

func (q *InventoryQueryService) GetInventoryVoucherByReference(ctx context.Context, ShopID dot.ID, refID dot.ID, refType inventory_voucher_ref.InventoryVoucherRef) (*inventory.GetInventoryVoucherByReferenceResponse, error) {
	result, err := q.InventoryVoucherStore(ctx).ShopID(ShopID).RefID(refID).RefType(refType).ListInventoryVoucher()
	if err != nil {
		return nil, err
	}
	var status = status4.N
	if len(result) == 0 {
		status = status4.Z
	} else {
		for _, value := range result {
			if value.Status == status3.Z {
				status = status4.S
				break
			}
			if value.Status == status3.P {
				status = status4.P
				break
			}
		}
	}
	result, err = q.populateInventoryVouchers(ctx, result)
	if err != nil {
		return nil, err
	}
	return &inventory.GetInventoryVoucherByReferenceResponse{
		InventoryVouchers: result,
		Status:            status,
	}, nil
}

func (q *InventoryQueryService) ListInventoryVariantsByVariantIDs(
	ctx context.Context, args *inventory.ListInventoryVariantsByVariantIDsArgs,
) (*inventory.GetInventoryVariantsResponse, error) {
	inventoryVariants, err := q.InventoryStore(ctx).ShopID(args.ShopID).VariantIDs(args.VariantIDs...).ListInventory()
	if err != nil {
		return nil, err
	}
	return &inventory.GetInventoryVariantsResponse{InventoryVariants: inventoryVariants}, nil
}

func (q *InventoryQueryService) populateInventoryVouchers(ctx context.Context, args []*inventory.InventoryVoucher) ([]*inventory.InventoryVoucher, error) {
	var stocktakeIDs []dot.ID
	for _, v := range args {
		if v.RefType == inventory_voucher_ref.StockTake {
			stocktakeIDs = append(stocktakeIDs, v.RefID)
		}
	}
	var mapStocktake = make(map[dot.ID]stocktake_type.StocktakeType)
	if len(stocktakeIDs) > 0 {
		queryStocktake := &stocktaking.GetStocktakesByIDsQuery{
			Ids:    stocktakeIDs,
			ShopID: args[0].ShopID,
		}
		err := q.StocktakeQuery.Dispatch(ctx, queryStocktake)
		if err != nil {
			return nil, err
		}
		for _, value := range queryStocktake.Result {
			mapStocktake[value.ID] = value.Type
		}
	}
	return util.PopulateInventoryVouchers(args, mapStocktake)
}

func (q *InventoryQueryService) populateInventoryVoucher(ctx context.Context, arg *inventory.InventoryVoucher) (*inventory.InventoryVoucher, error) {
	var stocktake *stocktaking.ShopStocktake
	if arg.RefType == inventory_voucher_ref.StockTake {
		queryStocktake := &stocktaking.GetStocktakeByIDQuery{
			Id:     arg.RefID,
			ShopID: arg.ShopID,
		}
		err := q.StocktakeQuery.Dispatch(ctx, queryStocktake)
		if err != nil {
			return nil, err
		}
		stocktake = queryStocktake.Result
	}
	return util.PopulateInventoryVoucher(arg, stocktake)
}
