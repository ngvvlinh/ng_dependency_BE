package aggregatex

import (
	"context"
	"time"

	"o.o/api/main/catalog"
	"o.o/api/main/inventory"
	"o.o/api/main/stocktaking"
	"o.o/api/top/types/etc/inventory_type"
	"o.o/api/top/types/etc/inventory_voucher_ref"
	"o.o/api/top/types/etc/status3"
	com "o.o/backend/com/main"
	catalogconvert "o.o/backend/com/main/catalog/convert"
	"o.o/backend/com/main/inventory/convert"
	"o.o/backend/com/main/inventory/model"
	invstore "o.o/backend/com/main/inventory/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/capi"
	"o.o/capi/dot"
)

var _ inventory.Aggregate = &InventoryAggregate{}
var scheme = conversion.Build(convert.RegisterConversions, catalogconvert.RegisterConversions)

type InventoryAggregate struct {
	InventoryStore        invstore.InventoryFactory
	InventoryVoucherStore invstore.InventoryVoucherFactory
	EventBus              capi.EventBus
	db                    *cmsql.Database
	StocktakeQuery        stocktaking.QueryBus
	CatalogQuery          catalog.QueryBus

	OrderStore sqlstore.OrderStoreInterface
}

func NewAggregateInventory(
	eventBus capi.EventBus,
	db com.MainDB,
	StocktakeQuery stocktaking.QueryBus,
	CatalogQ catalog.QueryBus,
	OrderStore sqlstore.OrderStoreInterface,
) *InventoryAggregate {
	return &InventoryAggregate{
		InventoryStore:        invstore.NewInventoryStore(db),
		InventoryVoucherStore: invstore.NewInventoryVoucherStore(db),
		EventBus:              eventBus,
		db:                    db,
		StocktakeQuery:        StocktakeQuery,
		CatalogQuery:          CatalogQ,
		OrderStore:            OrderStore,
	}
}

func InventoryAggregateMessageBus(q *InventoryAggregate) inventory.CommandBus {
	b := bus.New()
	return inventory.NewAggregateHandler(q).RegisterHandlers(b)
}

func (q *InventoryAggregate) CreateInventoryVoucher(ctx context.Context, Overstock bool, inventoryVoucher *inventory.CreateInventoryVoucherArgs) (*inventory.InventoryVoucher, error) {
	if inventoryVoucher.ShopID == 0 || inventoryVoucher.Type == inventory_type.Unknown {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing value requirement")
	}
	if inventoryVoucher.RefID != 0 {
		inventoryVoucherRefIDs, err := q.InventoryVoucherStore(ctx).ShopID(inventoryVoucher.ShopID).Type(inventoryVoucher.Type).RefID(inventoryVoucher.RefID).ListInventoryVoucher()
		if err != nil {
			return nil, err
		}
		for _, value := range inventoryVoucherRefIDs {
			if value.Status == status3.P || value.Status == status3.Z {
				return nil, cm.Errorf(cm.InvalidArgument, nil, "Phiếu xuất nhập kho cho đơn %v đã tồn tại, Vui lòng kiểm tra lại.", value.Code)
			}
		}
	}
	event := &inventory.InventoryVoucherCreatingEvent{
		ShopID: inventoryVoucher.ShopID,
		Line:   inventoryVoucher.Lines,
	}
	err := q.EventBus.Publish(ctx, event)
	if err != nil {
		return nil, err
	}
	var totalAmount = 0
	var listInventoryModel []*inventory.InventoryVariant
	totalAmount, listInventoryModel, err = q.PreInventoryVariantForVoucher(ctx, Overstock, inventoryVoucher)
	if err != nil {
		return nil, err
	}
	inventoryVoucher.TotalAmount = totalAmount
	var voucher inventory.InventoryVoucher
	if err = scheme.Convert(inventoryVoucher, &voucher); err != nil {
		return nil, err
	}

	err = q.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		for _, value := range listInventoryModel {
			err = q.InventoryStore(ctx).ShopID(inventoryVoucher.ShopID).VariantID(value.VariantID).UpdateInventoryVariantAll(value)
			if err != nil {
				return err
			}
		}
		var maxCodeNorm int
		inventoryVoucherTemp, err := q.InventoryVoucherStore(ctx).ShopID(inventoryVoucher.ShopID).GetInventoryVoucherByMaximumCodeNorm()
		switch cm.ErrorCode(err) {
		case cm.NoError:
			maxCodeNorm = inventoryVoucherTemp.CodeNorm
		case cm.NotFound:
			// no-op
		default:
			return err
		}
		if maxCodeNorm >= convert.MaxCodeNorm {
			return cm.Errorf(cm.InvalidArgument, nil, "Vui lòng nhập mã")
		}
		codeNorm := maxCodeNorm + 1
		voucher.Code = convert.GenerateCode(codeNorm)
		voucher.CodeNorm = codeNorm
		err = q.InventoryVoucherStore(ctx).Create(&voucher)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return q.InventoryVoucherStore(ctx).ShopID(inventoryVoucher.ShopID).ID(voucher.ID).Get()
}

func (q *InventoryAggregate) PreInventoryVariantForVoucher(ctx context.Context, overStock bool, args *inventory.CreateInventoryVoucherArgs) (totalAmount int, listInventoryVariants []*inventory.InventoryVariant, err error) {

	totalAmount = 0
	var inventoryvariant *inventory.InventoryVariant

	// Check have been existed variant_id in database table inventory_variant
	for key, value := range args.Lines {
		if errValidate := validateInventoryVoucherItem(value); errValidate != nil {
			return 0, nil, errValidate
		}
		inventoryvariant, err = q.InventoryStore(ctx).ShopID(args.ShopID).VariantID(value.VariantID).Get()
		if err != nil && cm.ErrorCode(err) == cm.NotFound {
			// Create InventoryVariant follow variant_id if it have not been exit
			inventoryvariant, err = q.CreateInventoryVariant(ctx, &inventory.CreateInventoryVariantArgs{
				ShopID:    args.ShopID,
				VariantID: value.VariantID,
			})
			if err != nil {
				return 0, nil, err
			}
		}
		if err != nil && cm.ErrorCode(err) != cm.NotFound {
			return 0, nil, err
		}

		if args.RefType == inventory_voucher_ref.Order || args.RefType == inventory_voucher_ref.StockTake || args.RefType == inventory_voucher_ref.Refund {
			args.Lines[key].Price = inventoryvariant.CostPrice
		}
		totalAmount = totalAmount + args.Lines[key].Price*value.Quantity

		// if InventoryVoucher is type 'out' move InventoryVariant quantity from onhand -> picked
		if args.Type == inventory_type.Out {
			if !overStock && inventoryvariant.QuantityOnHand < value.Quantity {
				return 0, nil, cm.Error(cm.InvalidArgument, "not enough quantity in inventory", nil)
			}
			inventoryvariant.QuantityOnHand = inventoryvariant.QuantityOnHand - value.Quantity
			inventoryvariant.QuantityPicked = inventoryvariant.QuantityPicked + value.Quantity
			listInventoryVariants = append(listInventoryVariants, inventoryvariant)
		}
	}
	return totalAmount, listInventoryVariants, err
}

func (q *InventoryAggregate) CheckInventoryVariantsQuantity(ctx context.Context, args *inventory.CheckInventoryVariantQuantityRequest) error {
	for _, value := range args.Lines {
		inventoryCore, err := q.InventoryStore(ctx).ShopID(args.ShopID).VariantID(value.VariantID).Get()
		if err != nil && cm.ErrorCode(err) == cm.NotFound {
			_, err = q.CreateInventoryVariant(ctx, &inventory.CreateInventoryVariantArgs{
				ShopID:    args.ShopID,
				VariantID: value.VariantID,
			})
			if !args.InventoryOverStock && value.Quantity > 0 {
				return cm.Error(cm.InvalidArgument, "not enough quantity in inventory", nil)
			}
		}
		if err != nil && cm.ErrorCode(err) != cm.NotFound {
			return err
		}
		if args.Type == inventory_type.Out {
			if !args.InventoryOverStock && inventoryCore.QuantityOnHand < value.Quantity {
				return cm.Error(cm.InvalidArgument, "not enough quantity in inventory", nil)
			}
		}
	}
	return nil
}

func (q *InventoryAggregate) CreateInventoryVariant(ctx context.Context, args *inventory.CreateInventoryVariantArgs) (*inventory.InventoryVariant, error) {
	if args.ShopID == 0 && args.VariantID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing value requirement")
	}
	_, err := q.InventoryStore(ctx).ShopID(args.ShopID).VariantID(args.VariantID).Get()
	if err != nil && cm.ErrorCode(err) == cm.NotFound {
		err = q.InventoryStore(ctx).Create(&model.InventoryVariant{
			ShopID:         args.ShopID,
			VariantID:      args.VariantID,
			QuantityOnHand: 0,
			QuantityPicked: 0,
			CostPrice:      0,
		})
	}
	if err != nil {
		return nil, err
	}
	result, err := q.InventoryStore(ctx).ShopID(args.ShopID).VariantID(args.VariantID).Get()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func validateInventoryVoucherItem(args *inventory.InventoryVoucherItem) error {
	if args.Price < 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Giá sản phẩm không được âm")
	}
	if args.Quantity < 1 {
		return cm.Errorf(cm.InvalidArgument, nil, "Số lượng nhập xuất phải lớn hơn 0")
	}
	return nil
}

func (q *InventoryAggregate) UpdateInventoryVariantCostPrice(ctx context.Context, args *inventory.UpdateInventoryVariantCostPriceRequest) (*inventory.InventoryVariant, error) {
	if args.ShopID == 0 || args.VariantID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing shop_id, variant_id")
	}
	inventoryVouchers, err := q.InventoryVoucherStore(ctx).ShopID(args.ShopID).RefType(inventory_voucher_ref.PurchaseOrder).VariantID(args.VariantID).ListInventoryVoucher()
	if err != nil {
		return nil, err
	}
	POExists := false
	var purchaseOrderID dot.ID
	for _, value := range inventoryVouchers {
		if value.Status == status3.P {
			POExists = true
			purchaseOrderID = value.RefID
			break
		}
	}
	for _, value := range inventoryVouchers {
		if value.Status == status3.P && value.RefID == purchaseOrderID && value.Type == inventory_type.Out {
			POExists = false
			break
		}
	}
	if POExists {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Không thể cập nhập giá vốn cho phiên bản đã có phiếu nhập hàng đã xác nhận")
	}
	inventoryVariant, err := q.InventoryStore(ctx).ShopID(args.ShopID).VariantID(args.VariantID).Get()
	switch cm.ErrorCode(err) {
	case cm.NoError:
	case cm.NotFound:
		inventoryVariant, err = q.CreateInventoryVariant(ctx, &inventory.CreateInventoryVariantArgs{
			ShopID:    args.ShopID,
			VariantID: args.VariantID,
		})
		if err != nil {
			return nil, err
		}
	default:
		return nil, err
	}
	inventoryVariant.CostPrice = args.CostPrice
	inventoryVariant.UpdatedAt = time.Now()
	err = q.InventoryStore(ctx).VariantID(args.VariantID).ShopID(args.ShopID).UpdateInventoryVariantAll(inventoryVariant)
	if err != nil {
		return nil, err
	}
	return q.InventoryStore(ctx).ShopID(args.ShopID).VariantID(args.VariantID).Get()
}
