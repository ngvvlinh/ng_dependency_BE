package aggregate

import (
	"context"
	"fmt"
	"time"

	"etop.vn/api/main/catalog/types"
	"etop.vn/api/main/inventory"
	"etop.vn/api/main/purchaseorder"
	"etop.vn/api/main/purchaserefund"
	"etop.vn/api/main/refund"
	"etop.vn/api/main/stocktaking"
	"etop.vn/api/shopping/tradering"
	"etop.vn/api/top/types/etc/inventory_auto"
	"etop.vn/api/top/types/etc/inventory_type"
	"etop.vn/api/top/types/etc/inventory_voucher_ref"
	"etop.vn/api/top/types/etc/status3"
	catalogconvert "etop.vn/backend/com/main/catalog/convert"
	"etop.vn/backend/com/main/inventory/convert"
	"etop.vn/backend/com/main/inventory/model"
	"etop.vn/backend/com/main/inventory/sqlstore"
	ordermodelx "etop.vn/backend/com/main/ordering/modelx"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/conversion"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/capi/dot"
)

var _ inventory.Aggregate = &InventoryAggregate{}
var scheme = conversion.Build(convert.RegisterConversions, catalogconvert.RegisterConversions)

type InventoryAggregate struct {
	InventoryStore        sqlstore.InventoryFactory
	InventoryVoucherStore sqlstore.InventoryVoucherFactory
	traderQuery           tradering.QueryBus
	EventBus              bus.Bus
	db                    *cmsql.Database
	PurchaseOrderQuery    purchaseorder.QueryBus
	StocktakeQuery        stocktaking.QueryBus
	RefundQuery           refund.QueryBus
	PurchaseRefundQuery   purchaserefund.QueryBus
}

func NewAggregateInventory(eventBus bus.Bus,
	db *cmsql.Database,
	traderQuery tradering.QueryBus,
	purchaseOrderQuery purchaseorder.QueryBus,
	StocktakeQuery stocktaking.QueryBus,
	refundQuery refund.QueryBus,
	purchaserRefundQuery purchaserefund.QueryBus,
) *InventoryAggregate {
	return &InventoryAggregate{
		InventoryStore:        sqlstore.NewInventoryStore(db),
		InventoryVoucherStore: sqlstore.NewInventoryVoucherStore(db),
		EventBus:              eventBus,
		traderQuery:           traderQuery,
		db:                    db,
		PurchaseOrderQuery:    purchaseOrderQuery,
		StocktakeQuery:        StocktakeQuery,
		RefundQuery:           refundQuery,
		PurchaseRefundQuery:   purchaserRefundQuery,
	}
}

func (q *InventoryAggregate) MessageBus() inventory.CommandBus {
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
			if (value.Status == status3.P || value.Status == status3.Z) && value.Rollback == inventoryVoucher.Rollback {
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
	err = checkInventoryVoucherRefType(&voucher)
	if err != nil {
		return nil, err
	}
	if voucher.TraderID != 0 {
		err := q.validateTrader(ctx, voucher.ShopID, &voucher)
		if err != nil {
			return nil, err
		}
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

func (q *InventoryAggregate) validateTrader(ctx context.Context, shopID dot.ID, voucher *inventory.InventoryVoucher) error {
	query := &tradering.GetTraderInfoByIDQuery{
		ID:     voucher.TraderID,
		ShopID: shopID,
	}
	if err := q.traderQuery.Dispatch(ctx, query); err != nil {
		return cm.MapError(err).
			Map(cm.NotFound, cm.FailedPrecondition, "Đối tác không hợp lệ").
			Throw()
	}
	voucher.Trader = &inventory.Trader{
		ID:       query.Result.ID,
		Type:     query.Result.Type,
		FullName: query.Result.FullName,
		Phone:    query.Result.Phone,
	}
	return nil
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

func (q *InventoryAggregate) UpdateInventoryVoucher(ctx context.Context, args *inventory.UpdateInventoryVoucherArgs) (*inventory.InventoryVoucher, error) {
	if args.ShopID == 0 || args.ID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing value requirement")
	}
	inventoryVoucher, err := q.InventoryVoucherStore(ctx).ShopID(args.ShopID).ID(args.ID).Get()
	if err != nil {
		return nil, err
	}
	if inventoryVoucher.Status != status3.Z {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "This inventory is already confirmed or cancelled")
	}
	if inventoryVoucher.Type == inventory_type.Out {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Can not update inventory delivery voucher")
	}
	event := &inventory.InventoryVoucherUpdatingEvent{
		ShopID: args.ShopID,
		Line:   args.Lines,
	}
	err = q.EventBus.Publish(ctx, event)
	if err != nil {
		return nil, err
	}
	var totalAmount = 0
	for _, value := range args.Lines {
		if errValidate := validateInventoryVoucherItem(value); errValidate != nil {
			return nil, errValidate
		}
		totalAmount = totalAmount + value.Quantity*value.Price
	}
	if args.TotalAmount != totalAmount {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Tổng giá trị phiếu không hợp lệ")
	}

	if err := scheme.Convert(args, inventoryVoucher); err != nil {
		return nil, err
	}
	if args.TraderID.Apply(inventoryVoucher.TraderID) != inventoryVoucher.TraderID {
		err := q.validateTrader(ctx, inventoryVoucher.ShopID, inventoryVoucher)
		if err != nil {
			return nil, err
		}
	}
	err = q.InventoryVoucherStore(ctx).ShopID(args.ShopID).ID(args.ID).UpdateInventoryVoucherAll(inventoryVoucher)
	if err != nil {
		return nil, err
	}

	return q.InventoryVoucherStore(ctx).ShopID(args.ShopID).ID(args.ID).Get()
}

func (q *InventoryAggregate) AdjustInventoryQuantity(ctx context.Context, overStock bool, args *inventory.AdjustInventoryQuantityArgs) (*inventory.AdjustInventoryQuantityRespone, error) {
	if args.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing value requirement")
	}
	var linesCheckin []*inventory.InventoryVoucherItem
	var linesCheckout []*inventory.InventoryVoucherItem
	var listVariantID []dot.ID
	var err error
	linesCheckin, linesCheckout, listVariantID, err = q.DevideInOutInventoryVoucher(ctx, args)
	if err != nil {
		return nil, err
	}
	if args.Title == "" {
		args.Title = "Phiếu cân bằng kho"
	}
	var inventoryVoucherInID dot.ID
	if len(linesCheckin) > 0 {
		inventoryVoucherInID, err = q.CreateVoucherForAdjustInventoryQuantity(ctx, overStock, args, linesCheckin, inventory_type.In)
		if err != nil {
			return nil, err
		}
	}
	var inventoryVoucherOutID dot.ID
	if len(linesCheckout) > 0 {
		inventoryVoucherOutID, err = q.CreateVoucherForAdjustInventoryQuantity(ctx, overStock, args, linesCheckout, inventory_type.Out)
		if err != nil {
			return nil, err
		}
	}

	inventoryVouchers, err := q.InventoryVoucherStore(ctx).ShopID(args.ShopID).IDs(inventoryVoucherInID, inventoryVoucherOutID).ListInventoryVoucher()
	if err != nil {
		return nil, err
	}
	resultUpdate, err := q.InventoryStore(ctx).ShopID(args.ShopID).VariantIDs(listVariantID...).ListInventory()
	if err != nil {
		return nil, err
	}
	return &inventory.AdjustInventoryQuantityRespone{
		InventoryVariants: resultUpdate,
		InventoryVouchers: inventoryVouchers,
	}, nil
}

func (q *InventoryAggregate) DevideInOutInventoryVoucher(ctx context.Context,
	args *inventory.AdjustInventoryQuantityArgs) ([]*inventory.InventoryVoucherItem,
	[]*inventory.InventoryVoucherItem,
	[]dot.ID, error) {
	var listVariantID []dot.ID
	var linesCheckin []*inventory.InventoryVoucherItem
	var linesCheckout []*inventory.InventoryVoucherItem

	for _, value := range args.Lines {
		if value.QuantitySummary < 0 {
			return nil, nil, nil, cm.Errorf(cm.InvalidArgument, nil, "Số lượng sản phẩm cân bằng kho phải lớn hơn 0")
		}
		listVariantID = append(listVariantID, value.VariantID)
		result, err := q.InventoryStore(ctx).ShopID(args.ShopID).VariantID(value.VariantID).Get()
		if err != nil && cm.ErrorCode(err) == cm.NotFound {
			linesCheckin = append(linesCheckin, &inventory.InventoryVoucherItem{
				VariantID: value.VariantID,
				Price:     value.CostPrice,
				Quantity:  value.QuantitySummary,
			})
			continue
		}
		if err != nil {
			return nil, nil, nil, err
		}
		if value.QuantitySummary > (result.QuantityOnHand + result.QuantityPicked) {
			linesCheckin = append(linesCheckin, &inventory.InventoryVoucherItem{
				VariantID: value.VariantID,
				Price:     result.CostPrice,
				Quantity:  value.QuantitySummary - (result.QuantityOnHand + result.QuantityPicked),
			})
		} else if value.QuantitySummary < (result.QuantityOnHand + result.QuantityPicked) {
			linesCheckout = append(linesCheckout, &inventory.InventoryVoucherItem{
				VariantID: value.VariantID,
				Price:     result.CostPrice,
				Quantity:  (result.QuantityOnHand + result.QuantityPicked) - value.QuantitySummary,
			})
		}
	}
	return linesCheckin, linesCheckout, listVariantID, nil
}

func (q *InventoryAggregate) CreateVoucherForAdjustInventoryQuantity(ctx context.Context, overStock bool, info *inventory.AdjustInventoryQuantityArgs,
	lines []*inventory.InventoryVoucherItem,
	typeVoucher inventory_type.InventoryVoucherType) (dot.ID, error) {
	var totalValue = 0
	for _, value := range lines {
		totalValue = totalValue + value.Price*value.Quantity
	}
	result, err := q.CreateInventoryVoucher(ctx, overStock, &inventory.CreateInventoryVoucherArgs{
		Title:       info.Title,
		ShopID:      info.ShopID,
		CreatedBy:   info.UserID,
		TotalAmount: totalValue,
		Type:        typeVoucher,
		Note:        info.Note,
		Lines:       lines,
	})
	if err != nil {
		return 0, err
	}
	return result.ID, nil
}

func (q *InventoryAggregate) ConfirmInventoryVoucher(ctx context.Context, args *inventory.ConfirmInventoryVoucherArgs) (*inventory.InventoryVoucher, error) {
	if args.ShopID == 0 || args.ID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing value requirement")
	}
	inventoryVoucher, err := q.InventoryVoucherStore(ctx).ShopID(args.ShopID).ID(args.ID).GetDB()
	if err != nil {
		return nil, err
	}
	if inventoryVoucher.Status != status3.Z {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Inventory voucher already confirmed or cancelled")
	}
	for _, value := range inventoryVoucher.Lines {
		var data *inventory.InventoryVariant
		data, err = q.InventoryStore(ctx).ShopID(args.ShopID).VariantID(value.VariantID).Get()
		if err != nil {
			return nil, err
		}
		if inventoryVoucher.Type == inventory_type.Out {
			// if purchase refund -> change cost_price
			if inventory_voucher_ref.PurchaseRefund == inventoryVoucher.RefType {
				currentQuantity := data.QuantityPicked + data.QuantityOnHand
				currentValue := currentQuantity * data.CostPrice
				outValue := value.Quantity * value.Price
				if currentQuantity-value.Quantity != 0 {
					data.CostPrice = (currentValue - outValue) / (currentQuantity - value.Quantity)
				}
			}
			data.QuantityPicked = data.QuantityPicked - value.Quantity
		} else if inventoryVoucher.Type == inventory_type.In {
			// if TraderID = 0 -> stocktake, TraderID != 0 -> purchase order
			if inventoryVoucher.TraderID != 0 {
				if data.QuantityOnHand < 0 {
					data.CostPrice = value.Price
				} else {
					// Update costPirce from Purchase Order
					data.CostPrice = AvgValue(data.CostPrice, value.Price, data.QuantityOnHand, value.Quantity)
				}
			}
			data.QuantityOnHand = data.QuantityOnHand + value.Quantity
		}
		err = q.InventoryStore(ctx).VariantID(value.VariantID).ShopID(args.ShopID).UpdateInventoryVariantAll(data)
		if err != nil {
			return nil, err
		}
	}
	inventoryVoucher.Status = status3.P
	inventoryVoucher.ConfirmedAt = time.Now()

	err = q.InventoryVoucherStore(ctx).ShopID(args.ShopID).ID(args.ID).UpdateInventoryVoucherAllDB(inventoryVoucher)
	if err != nil {
		return nil, err
	}
	return q.InventoryVoucherStore(ctx).ShopID(args.ShopID).ID(args.ID).Get()
}

func AvgValue(value1 int, value2 int, quantity1 int, quantity2 int) int {
	if quantity1+quantity2 == 0 {
		return 0
	}
	return int((int64(value1)*int64(quantity1) + int64(value2)*int64(quantity2)) / (int64(quantity1) + int64(quantity2)))
}

func (q *InventoryAggregate) CancelInventoryVoucher(ctx context.Context, args *inventory.CancelInventoryVoucherArgs) (*inventory.InventoryVoucher, error) {
	if args.ShopID == 0 || args.ID == 0 || args.CancelReason == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing value requirement")
	}
	inventoryVoucher, err := q.InventoryVoucherStore(ctx).ShopID(args.ShopID).ID(args.ID).GetDB()
	if err != nil {
		return nil, err
	}
	if inventoryVoucher.Status != status3.Z {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Inventory voucher already confirmed or cancelled")
	}
	if inventoryVoucher.Type == inventory_type.Out {
		for _, value := range inventoryVoucher.Lines {
			var data *inventory.InventoryVariant
			data, err = q.InventoryStore(ctx).ShopID(args.ShopID).VariantID(value.VariantID).Get()
			if err != nil {
				return nil, err
			}
			data.CostPrice = AvgValue(data.CostPrice, value.Price, data.QuantityOnHand, value.Quantity)
			data.QuantityPicked = data.QuantityPicked - value.Quantity
			data.QuantityOnHand = data.QuantityOnHand + value.Quantity

			err = q.InventoryStore(ctx).VariantID(value.VariantID).ShopID(args.ShopID).UpdateInventoryVariantAll(data)
			if err != nil {
				return nil, err
			}
		}
	}
	inventoryVoucher.Status = status3.N
	inventoryVoucher.CancelledAt = time.Now()
	inventoryVoucher.CancelReason = args.CancelReason
	err = q.InventoryVoucherStore(ctx).ShopID(args.ShopID).ID(args.ID).UpdateInventoryVoucherAllDB(inventoryVoucher)
	if err != nil {
		return nil, err
	}
	inventoryVoucherConfirmed, err := q.InventoryVoucherStore(ctx).ShopID(args.ShopID).ID(args.ID).Get()
	return inventoryVoucherConfirmed, err
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

func checkInventoryVoucherRefType(inventoryVoucher *inventory.InventoryVoucher) error {
	switch inventoryVoucher.RefType {
	case inventory_voucher_ref.Order:
		if inventoryVoucher.Type != inventory_type.In && inventoryVoucher.Type != inventory_type.Out {
			return cm.Error(cm.InvalidArgument, "'type' không đúng. Bán hàng chỉ có thể là 'in' hoặc 'out'", nil)
		}
	case inventory_voucher_ref.StockTake:
		if inventoryVoucher.Type != inventory_type.Out && inventoryVoucher.Type != inventory_type.In {
			return cm.Error(cm.InvalidArgument, "'type' không đúng.Kiểm kho chỉ có thể là 'in' hoặc 'out'", nil)
		}
	case inventory_voucher_ref.PurchaseOrder:
		if inventoryVoucher.Type == inventory_type.Out && !inventoryVoucher.Rollback {
			return cm.Error(cm.InvalidArgument, "'type' không đúng.Nhập hàng chỉ có thể là 'in'", nil)
		}
		if inventoryVoucher.Type == inventory_type.In && inventoryVoucher.Rollback {
			return cm.Error(cm.InvalidArgument, "'type' không đúng. Hủy nhập hàng chỉ có thể là type 'out'", nil)
		}
	case inventory_voucher_ref.Refund:
		if inventoryVoucher.Type == inventory_type.Out && !inventoryVoucher.Rollback {
			return cm.Error(cm.InvalidArgument, "'type' không đúng.Trả hàng chỉ có thể là 'in'", nil)
		}
		if inventoryVoucher.Type == inventory_type.In && inventoryVoucher.Rollback {
			return cm.Error(cm.InvalidArgument, "'type' không đúng. Hủy trả hàng chỉ có thể là 'out'", nil)
		}
	case inventory_voucher_ref.PurchaseRefund:
		if inventoryVoucher.Type != inventory_type.In && inventoryVoucher.Rollback {
			return cm.Error(cm.InvalidArgument, "'type' không đúng. Trả hàng nhập chỉ có thể là 'out'", nil)
		}
		if inventoryVoucher.Type != inventory_type.Out && !inventoryVoucher.Rollback {
			return cm.Error(cm.InvalidArgument, "'type' không đúng. Hủy trả hàng nhập chỉ có thể là 'in'", nil)
		}
	case inventory_voucher_ref.Unknown:
		if inventoryVoucher.RefName != "" || inventoryVoucher.RefID != 0 {
			return cm.Error(cm.InvalidArgument, "'ref_type','ref_id' hoặc 'ref_name' không đúng. Vui lòng kiểm tra lại", nil)
		}
		return nil
	default:
		return cm.Error(cm.InvalidArgument, "'ref_type' không đúng. Vui lòng nhập đúng ref_type", nil)

	}
	return nil
}

func (q *InventoryAggregate) CreateInventoryVoucherByQuantityChange(ctx context.Context, args *inventory.CreateInventoryVoucherByQuantityChangeRequest) (*inventory.CreateInventoryVoucherByQuantityChangeResponse, error) {
	var inventoryVariantIDs []dot.ID
	var inventoryVoucherIn []*inventory.InventoryVoucherItem
	var inventoryVoucherOut []*inventory.InventoryVoucherItem
	for _, value := range args.Lines {
		inventoryVariantIDs = append(inventoryVariantIDs, value.ItemInfo.VariantID)
	}
	listVariant, err := q.InventoryStore(ctx).ShopID(args.ShopID).VariantIDs(inventoryVariantIDs...).ListInventory()
	if err != nil {
		return nil, err
	}
	var mapInventoryVariantInfo = make(map[dot.ID]*inventory.InventoryVariant)
	for _, value := range listVariant {
		mapInventoryVariantInfo[value.VariantID] = value
	}

	for _, value := range args.Lines {
		inventoryVoucherItem := &inventory.InventoryVoucherItem{
			ProductID:   value.ItemInfo.ProductID,
			ProductName: value.ItemInfo.ProductName,
			VariantID:   value.ItemInfo.VariantID,
			VariantName: value.ItemInfo.VariantName,
			Price:       value.ItemInfo.Price,
			Code:        value.ItemInfo.Code,
			ImageURL:    value.ItemInfo.ImageURL,
			Attributes:  value.ItemInfo.Attributes,
		}
		if value.QuantityChange > 0 {
			inventoryVoucherItem.Quantity = value.QuantityChange
			if mapInventoryVariantInfo[value.ItemInfo.VariantID] != nil {
				inventoryVoucherItem.Price = mapInventoryVariantInfo[value.ItemInfo.VariantID].CostPrice
			}
			inventoryVoucherIn = append(inventoryVoucherIn, inventoryVoucherItem)
		} else if value.QuantityChange < 0 {
			inventoryVoucherItem.Quantity = value.QuantityChange * -1
			if mapInventoryVariantInfo[value.ItemInfo.VariantID] != nil {
				inventoryVoucherItem.Price = mapInventoryVariantInfo[value.ItemInfo.VariantID].CostPrice
			}
			inventoryVoucherOut = append(inventoryVoucherOut, inventoryVoucherItem)
		}

	}
	var typeIn = &inventory.InventoryVoucher{}
	if len(inventoryVoucherIn) != 0 {
		typeIn, err = q.CreateInventoryVoucher(ctx, args.Overstock, &inventory.CreateInventoryVoucherArgs{
			ShopID:    args.ShopID,
			CreatedBy: args.CreatedBy,
			Title:     args.Title,
			RefID:     args.RefID,
			RefType:   args.RefType,
			RefName:   args.RefName,
			RefCode:   args.RefCode,
			TraderID:  0,
			Type:      inventory_type.In,
			Note:      args.NoteIn,
			Lines:     inventoryVoucherIn,
		})
		if err != nil {
			return nil, err
		}
	}
	var typeOut = &inventory.InventoryVoucher{}
	if len(inventoryVoucherOut) != 0 {
		typeOut, err = q.CreateInventoryVoucher(ctx, args.Overstock, &inventory.CreateInventoryVoucherArgs{
			ShopID:    args.ShopID,
			CreatedBy: args.CreatedBy,
			Title:     args.Title,
			RefID:     args.RefID,
			RefType:   args.RefType,
			RefName:   args.RefName,
			RefCode:   args.RefCode,
			TraderID:  0,
			Type:      inventory_type.Out,
			Note:      args.NoteOut,
			Lines:     inventoryVoucherOut,
		})
		if err != nil {
			return nil, err
		}
	}

	return &inventory.CreateInventoryVoucherByQuantityChangeResponse{
		TypeIn:  typeIn,
		TypeOut: typeOut,
	}, nil
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
	for _, value := range inventoryVouchers {
		if value.Status == status3.P {
			POExists = true
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

func (q *InventoryAggregate) CreateInventoryVoucherByReference(ctx context.Context, args *inventory.CreateInventoryVoucherByReferenceArgs) ([]*inventory.InventoryVoucher, error) {
	switch args.RefType {
	case inventory_voucher_ref.PurchaseOrder:
		return q.CreateInventoryVoucherByPurchaseOrder(ctx, args)
	case inventory_voucher_ref.StockTake:
		return q.CreateInventoryVoucherByStockTake(ctx, args)
	case inventory_voucher_ref.Order:
		return q.CreateInventoryVoucherByOrder(ctx, args)
	case inventory_voucher_ref.Refund:
		return q.CreateInventoryVoucherByRefund(ctx, args)
	case inventory_voucher_ref.PurchaseRefund:
		return q.CreateInventoryVoucherByPurchaseRefund(ctx, args)
	default:
		return nil, cm.Error(cm.InvalidArgument, "wrong ref_type", nil)
	}
}

func (q *InventoryAggregate) CreateInventoryVoucherByPurchaseRefund(ctx context.Context, args *inventory.CreateInventoryVoucherByReferenceArgs) ([]*inventory.InventoryVoucher, error) {
	var items []*inventory.InventoryVoucherItem
	queryPurchaseRefund := &purchaserefund.GetPurchaseRefundByIDQuery{
		ID:     args.RefID,
		ShopID: args.ShopID,
	}
	if err := q.PurchaseRefundQuery.Dispatch(ctx, queryPurchaseRefund); err != nil {
		return nil, err
	}
	if queryPurchaseRefund.Result.Status != status3.P {
		return nil, cm.Error(cm.InvalidArgument, "không thể tạo phiếu kiểm kho cho Refund chưa được xác nhận.", nil)
	}
	for _, value := range queryPurchaseRefund.Result.Lines {
		price := value.PaymentPrice + value.Adjustment
		if queryPurchaseRefund.Result.TotalAdjustment != 0 {
			valueLine := (value.Quantity * (value.PaymentPrice + value.Adjustment))
			price += (queryPurchaseRefund.Result.TotalAdjustment * valueLine / queryPurchaseRefund.Result.BasketValue) / value.Quantity
		}
		items = append(items, &inventory.InventoryVoucherItem{
			Price:       price,
			ProductID:   value.ProductID,
			ProductName: value.ProductName,
			VariantID:   value.VariantID,
			VariantName: value.ProductName,
			Quantity:    value.Quantity,
			Code:        value.Code,
			ImageURL:    value.ImageURL,
			Attributes:  value.Attributes,
		})
	}
	inventoryVoucherCreateRequest := &inventory.CreateInventoryVoucherArgs{
		ShopID:    args.ShopID,
		CreatedBy: args.UserID,
		Title:     "Xuất kho khi trả hàng nhập",
		RefID:     args.RefID,
		RefType:   args.RefType,
		RefName:   args.RefType.GetLabelRefName(),
		RefCode:   queryPurchaseRefund.Result.Code,
		TraderID:  queryPurchaseRefund.Result.SupplierID,
		Type:      inventory_type.Out,
		Note:      fmt.Sprintf("Tạo phiếu xuất kho theo phiếu trả hàng nhập %v", queryPurchaseRefund.Result.Code),
		Lines:     items,
	}
	createResult, err := q.CreateInventoryVoucher(ctx, args.OverStock, inventoryVoucherCreateRequest)
	if err != nil {
		return nil, err
	}
	var listInventoryVoucher []*inventory.InventoryVoucher
	listInventoryVoucher = append(listInventoryVoucher, createResult)
	return listInventoryVoucher, nil
}

func (q *InventoryAggregate) CreateInventoryVoucherByRefund(ctx context.Context, args *inventory.CreateInventoryVoucherByReferenceArgs) ([]*inventory.InventoryVoucher, error) {
	var items []*inventory.InventoryVoucherItem
	queryRefund := &refund.GetRefundByIDQuery{
		ID:     args.RefID,
		ShopID: args.ShopID,
	}
	if err := q.RefundQuery.Dispatch(ctx, queryRefund); err != nil {
		return nil, err
	}
	if queryRefund.Result.Status != status3.P {
		return nil, cm.Error(cm.InvalidArgument, "không thể tạo phiếu kiểm kho cho Refund chưa được xác nhận.", nil)
	}
	for _, value := range queryRefund.Result.Lines {
		items = append(items, &inventory.InventoryVoucherItem{
			ProductID:   value.ProductID,
			ProductName: value.ProductName,
			VariantID:   value.VariantID,
			VariantName: value.ProductName,
			Quantity:    value.Quantity,
			Code:        value.Code,
			ImageURL:    value.ImageURL,
			Attributes:  value.Attributes,
		})
	}
	inventoryVoucherCreateRequest := &inventory.CreateInventoryVoucherArgs{
		ShopID:    args.ShopID,
		CreatedBy: args.UserID,
		Title:     "Nhập kho khi nhập hàng",
		RefID:     args.RefID,
		RefType:   args.RefType,
		RefName:   args.RefType.GetLabelRefName(),
		RefCode:   queryRefund.Result.Code,
		TraderID:  queryRefund.Result.CustomerID,
		Type:      inventory_type.In,
		Note:      fmt.Sprintf("Tạo phiếu nhập kho theo đơn trả hàng %v", queryRefund.Result.Code),
		Lines:     items,
	}
	createResult, err := q.CreateInventoryVoucher(ctx, args.OverStock, inventoryVoucherCreateRequest)
	if err != nil {
		return nil, err
	}
	var listInventoryVoucher []*inventory.InventoryVoucher
	listInventoryVoucher = append(listInventoryVoucher, createResult)
	return listInventoryVoucher, nil
}

func (q *InventoryAggregate) CreateInventoryVoucherByPurchaseOrder(ctx context.Context, args *inventory.CreateInventoryVoucherByReferenceArgs) ([]*inventory.InventoryVoucher, error) {
	var items []*inventory.InventoryVoucherItem
	// check order_id exit
	queryPurchaseOrder := &purchaseorder.GetPurchaseOrderByIDQuery{
		ID:     args.RefID,
		ShopID: args.ShopID,
	}
	if err := q.PurchaseOrderQuery.Dispatch(ctx, queryPurchaseOrder); err != nil {
		return nil, err
	}
	if queryPurchaseOrder.Result.Status != status3.P {
		return nil, cm.Error(cm.InvalidArgument, "không thể tạo phiếu kiểm kho cho Purchase Order chưa được xác nhận.", nil)
	}
	// GET info and put it to cmd
	for _, value := range queryPurchaseOrder.Result.Lines {
		adjustment := queryPurchaseOrder.Result.TotalFee - queryPurchaseOrder.Result.TotalDiscount
		price := value.PaymentPrice - value.Discount
		if adjustment != 0 {
			price += (adjustment * (value.Quantity * (value.PaymentPrice - value.Discount)) / queryPurchaseOrder.Result.BasketValue) / value.Quantity
		}
		items = append(items, &inventory.InventoryVoucherItem{
			ProductID:   value.ProductID,
			ProductName: value.ProductName,
			VariantID:   value.VariantID,
			Quantity:    value.Quantity,
			Price:       price,
			Code:        value.Code,
			ImageURL:    value.ImageUrl,
			Attributes:  value.Attributes,
		})
	}
	inventoryVoucherCreateRequest := &inventory.CreateInventoryVoucherArgs{
		ShopID:    args.ShopID,
		CreatedBy: args.UserID,
		Title:     "Nhập kho khi nhập hàng",
		RefID:     args.RefID,
		RefType:   args.RefType,
		RefName:   args.RefType.GetLabelRefName(),
		RefCode:   queryPurchaseOrder.Result.Code,
		TraderID:  queryPurchaseOrder.Result.SupplierID,
		Type:      inventory_type.In,
		Note:      fmt.Sprintf("Tạo phiếu nhập kho theo đơn nhập hàng %v", queryPurchaseOrder.Result.Code),
		Lines:     items,
	}
	createResult, err := q.CreateInventoryVoucher(ctx, args.OverStock, inventoryVoucherCreateRequest)
	if err != nil {
		return nil, err
	}
	var listInventoryVoucher []*inventory.InventoryVoucher
	listInventoryVoucher = append(listInventoryVoucher, createResult)
	return listInventoryVoucher, err
}

func (q *InventoryAggregate) CreateInventoryVoucherByOrder(ctx context.Context, args *inventory.CreateInventoryVoucherByReferenceArgs) ([]*inventory.InventoryVoucher, error) {
	var items []*inventory.InventoryVoucherItem
	// check order_id exit
	queryOrder := &ordermodelx.GetOrderQuery{
		OrderID: args.RefID,
		ShopID:  args.ShopID,
	}
	if err := bus.Dispatch(ctx, queryOrder); err != nil {
		return nil, err
	}
	// GET info and put it to cmd
	for _, value := range queryOrder.Result.Order.Lines {
		if value.VariantID != 0 {
			var attributes []*types.Attribute
			err := scheme.Convert(value.Attributes, &attributes)
			if err != nil {
				return nil, err
			}
			items = append(items, &inventory.InventoryVoucherItem{
				ProductID:   value.ProductID,
				ProductName: value.ProductName,
				VariantID:   value.VariantID,
				Quantity:    value.Quantity,
				Code:        value.Code,
				ImageURL:    value.ImageURL,
				Attributes:  attributes,
			})
		}
	}
	if len(items) == 0 {
		return []*inventory.InventoryVoucher{}, nil
	}
	inventoryVoucherCreateRequest := &inventory.CreateInventoryVoucherArgs{
		ShopID:    args.ShopID,
		CreatedBy: args.UserID,
		RefID:     args.RefID,
		RefType:   args.RefType,
		RefCode:   queryOrder.Result.Order.Code,
		TraderID:  queryOrder.Result.Order.CustomerID,
		Type:      args.Type,
		Lines:     items,
	}
	switch inventoryVoucherCreateRequest.Type {
	case inventory_type.Out:
		inventoryVoucherCreateRequest.RefName = inventory_voucher_ref.Order.GetLabelRefName()
		inventoryVoucherCreateRequest.Title = "Xuất kho khi bán hàng"
		inventoryVoucherCreateRequest.Note = fmt.Sprintf("Tạo phiếu xuất kho theo đơn hàng %v", queryOrder.Result.Order.Code)
	case inventory_type.In:
		inventoryVoucherCreateRequest.Title = "Nhập kho khi Hủy bán hàng"
		inventoryVoucherCreateRequest.RefName = "Hủy " + inventory_voucher_ref.Order.GetLabelRefName()
		inventoryVoucherCreateRequest.Note = fmt.Sprintf("Tạo phiếu nhập kho theo đơn hàng %v", queryOrder.Result.Order.Code)
	}
	createResult, err := q.CreateInventoryVoucher(ctx, args.OverStock, inventoryVoucherCreateRequest)
	if err != nil {
		return nil, err
	}
	var listInventoryVoucher []*inventory.InventoryVoucher
	listInventoryVoucher = append(listInventoryVoucher, createResult)
	return listInventoryVoucher, err
}

func (q *InventoryAggregate) CreateInventoryVoucherByStockTake(ctx context.Context, args *inventory.CreateInventoryVoucherByReferenceArgs) ([]*inventory.InventoryVoucher, error) {
	// check order_id exit
	queryStocktake := &stocktaking.GetStocktakeByIDQuery{
		Id:     args.RefID,
		ShopID: args.ShopID,
	}
	if err := q.StocktakeQuery.Dispatch(ctx, queryStocktake); err != nil {
		return nil, err
	}
	if queryStocktake.Result.Status != status3.P {
		return nil, cm.Error(cm.InvalidArgument, "không thể tạo phiếu kiểm kho cho stocktake chưa được xác nhận.", nil)
	}
	// GET info and put it to cmd
	var inventoryVariantChange []*inventory.InventoryVariantQuantityChange
	for _, value := range queryStocktake.Result.Lines {
		inventoryVariantChange = append(inventoryVariantChange, &inventory.InventoryVariantQuantityChange{
			ItemInfo: &inventory.InventoryVoucherItem{
				ProductID:   value.ProductID,
				ProductName: value.ProductName,
				VariantID:   value.VariantID,
				VariantName: value.VariantName,
				Price:       value.CostPrice,
				Code:        value.Code,
				ImageURL:    value.ImageURL,
				Attributes:  value.Attributes,
			},
			QuantityChange: value.NewQuantity - value.OldQuantity,
		})
	}
	inventoryVoucherCreateRequest := &inventory.CreateInventoryVoucherByQuantityChangeRequest{
		ShopID:    args.ShopID,
		RefID:     args.RefID,
		RefType:   inventory_voucher_ref.StockTake,
		RefName:   inventory_voucher_ref.StockTake.GetLabelRefName(),
		Title:     "Phiếu kiểm kho",
		RefCode:   queryStocktake.Result.Code,
		Overstock: args.OverStock,
		CreatedBy: args.UserID,
		Lines:     inventoryVariantChange,
		NoteIn:    fmt.Sprintf("Tạo phiếu nhập kho theo phiếu kiểm kho %v", queryStocktake.Result.Code),
		NoteOut:   fmt.Sprintf("Tạo phiếu xuất kho theo phiếu kiểm kho  %v", queryStocktake.Result.Code),
	}
	createResult, err := q.CreateInventoryVoucherByQuantityChange(ctx, inventoryVoucherCreateRequest)
	if err != nil {
		return nil, err
	}
	var listInventoryVoucher []*inventory.InventoryVoucher
	if createResult.TypeIn.ID != 0 {
		listInventoryVoucher = append(listInventoryVoucher, createResult.TypeIn)
	}
	if createResult.TypeOut.ID != 0 {
		listInventoryVoucher = append(listInventoryVoucher, createResult.TypeOut)
	}
	return listInventoryVoucher, err
}

func (q *InventoryAggregate) CancelInventoryByRefID(ctx context.Context, args *inventory.CancelInventoryByRefIDRequest) (*inventory.CancelInventoryByRefIDResponse, error) {
	if args.AutoInventoryVoucher == inventory_auto.Unknown {
		return nil, nil
	}

	var inventoryVouchers []*inventory.InventoryVoucher
	inventoryVouchersData, err := q.InventoryVoucherStore(ctx).ShopID(args.ShopID).RefID(args.RefID).ListInventoryVoucher()
	if err != nil {
		return nil, err
	}

	// return if have any inventory voucher rollback
	for _, value := range inventoryVouchersData {
		if value.Rollback == true {
			return &inventory.CancelInventoryByRefIDResponse{
				InventoryVouchers: inventoryVouchers,
			}, nil
		}
	}

	for _, value := range inventoryVouchersData {
		switch value.Status {
		case status3.P:
			var typeInventoryVoucher inventory_type.InventoryVoucherType
			if value.Type == inventory_type.Out {
				typeInventoryVoucher = inventory_type.In
			} else {
				typeInventoryVoucher = inventory_type.Out
			}
			newVoucher := &inventory.CreateInventoryVoucherArgs{
				ShopID:      value.ShopID,
				CreatedBy:   value.CreatedBy,
				Title:       "Phiếu xuất nhập kho",
				RefID:       value.RefID,
				RefType:     value.RefType,
				RefName:     "Hủy " + value.RefName,
				RefCode:     value.RefCode,
				TraderID:    value.TraderID,
				TotalAmount: value.TotalAmount,
				Type:        typeInventoryVoucher,
				Lines:       value.Lines,
				Rollback:    true,
				Note:        getNoteCancelInventoryVoucher(value.RefType, value.RefCode),
			}
			result, err := q.CreateInventoryVoucher(ctx, args.InventoryOverStock, newVoucher)
			if err != nil {
				return nil, err
			}
			inventoryVouchers = append(inventoryVouchers, result)
		case status3.Z:
			cancelResult, err := q.CancelInventoryVoucher(ctx, &inventory.CancelInventoryVoucherArgs{
				ShopID:       value.ShopID,
				ID:           value.ID,
				UpdatedBy:    args.UpdateBy,
				CancelReason: getCancelReason(value.RefType, value.RefCode),
			})
			if err != nil {
				return nil, err
			}
			inventoryVouchers = append(inventoryVouchers, cancelResult)
		case status3.N:
			continue
		}
	}
	if args.AutoInventoryVoucher == inventory_auto.Confirm {
		for _, value := range inventoryVouchers {
			if value.Status == status3.Z {
				_, err = q.ConfirmInventoryVoucher(ctx, &inventory.ConfirmInventoryVoucherArgs{
					ShopID:    args.ShopID,
					ID:        value.ID,
					UpdatedBy: args.UpdateBy,
				})
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return &inventory.CancelInventoryByRefIDResponse{
		InventoryVouchers: inventoryVouchers,
	}, nil
}

func getNoteCancelInventoryVoucher(ref inventory_voucher_ref.InventoryVoucherRef, refCode string) string {
	switch ref {
	case inventory_voucher_ref.Order:
		return fmt.Sprintf("Tạo phiếu nhập kho theo đơn hàng %v", refCode)
	case inventory_voucher_ref.Refund:
		return fmt.Sprintf("Tạo phiếu xuất kho theo đơn trả hàng %v", refCode)
	case inventory_voucher_ref.PurchaseOrder:
		return fmt.Sprintf("Tạo phiếu xuất kho theo đơn nhập hàng %v", refCode)
	case inventory_voucher_ref.PurchaseRefund:
		return fmt.Sprintf("Tạo phiếu nhập kho theo đơn trả hàng nhập %v", refCode)
	case inventory_voucher_ref.StockTake:
		return fmt.Sprintf("Tạo phiếu xuất/nhập kho theo phiếu kiểm kho %v", refCode)
	default:
		return ""
	}
}

func getCancelReason(ref inventory_voucher_ref.InventoryVoucherRef, refCode string) string {
	switch ref {
	case inventory_voucher_ref.Order:
		return fmt.Sprintf("Hủy đơn hàng %v", refCode)
	case inventory_voucher_ref.Refund:
		return fmt.Sprintf("Hủy đơn trả hàng %v", refCode)
	case inventory_voucher_ref.PurchaseOrder:
		return fmt.Sprintf("Hủy đơn nhập hàng %v", refCode)
	case inventory_voucher_ref.PurchaseRefund:
		return fmt.Sprintf("Hủy đơn trả hàng nhập %v", refCode)
	case inventory_voucher_ref.StockTake:
		return fmt.Sprintf("Hủy phiếu kiểm kho %v", refCode)
	default:
		return ""
	}
}
