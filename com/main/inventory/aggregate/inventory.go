package aggregate

import (
	"context"
	"time"

	"etop.vn/api/main/etop"
	"etop.vn/api/main/inventory"
	"etop.vn/api/shopping/tradering"
	"etop.vn/backend/com/main/inventory/convert"
	"etop.vn/backend/com/main/inventory/model"
	"etop.vn/backend/com/main/inventory/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/conversion"
)

var _ inventory.Aggregate = &InventoryAggregate{}
var scheme = conversion.Build(convert.RegisterConversions)

type InventoryAggregate struct {
	InventoryStore        sqlstore.InventoryFactory
	InventoryVoucherStore sqlstore.InventoryVoucherFactory
	traderQuery           tradering.QueryBus
	EventBus              bus.Bus
	db                    *cmsql.Database
}

func NewAggregateInventory(eventBus bus.Bus, db *cmsql.Database, traderQuery tradering.QueryBus) *InventoryAggregate {
	return &InventoryAggregate{
		InventoryStore:        sqlstore.NewInventoryStore(db),
		InventoryVoucherStore: sqlstore.NewInventoryVoucherStore(db),
		EventBus:              eventBus,
		traderQuery:           traderQuery,
		db:                    db,
	}
}

func (q *InventoryAggregate) MessageBus() inventory.CommandBus {
	b := bus.New()
	return inventory.NewAggregateHandler(q).RegisterHandlers(b)
}

func (q *InventoryAggregate) CreateInventoryVoucher(ctx context.Context, Overstock bool, inventoryVoucher *inventory.CreateInventoryVoucherArgs) (*inventory.InventoryVoucher, error) {
	if inventoryVoucher.ShopID == 0 || inventoryVoucher.Type == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing value requirement")
	}
	if inventoryVoucher.RefID != 0 {
		inventoryVoucherRefIDs, err := q.InventoryVoucherStore(ctx).ShopID(inventoryVoucher.ShopID).RefID(inventoryVoucher.RefID).ListInventoryVoucher()
		if err != nil {
			return nil, err
		}
		for _, value := range inventoryVoucherRefIDs {
			if value.Status == etop.S3Positive || value.Status == etop.S3Zero {
				return nil, cm.Errorf(cm.InvalidArgument, nil, "Phiếu xuất nhập kho cho ref_id đã tồn tại, Vui lòng kiểm tra lại ", value.ID)
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
	var totalAmount int32 = 0
	var listInventoryModel []*inventory.InventoryVariant
	totalAmount, listInventoryModel, err = q.PreInventoryVariantForVoucher(ctx, Overstock, inventoryVoucher)
	if err != nil {
		return nil, err
	}
	if inventoryVoucher.TotalAmount != totalAmount {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Tổng giá trị phiếu không hợp lệ")
	}
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
		var maxCodeNorm int32
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
		voucher.Code = convert.GenerateCode(int(codeNorm))
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

func (q *InventoryAggregate) validateTrader(ctx context.Context, shopID int64, voucher *inventory.InventoryVoucher) error {
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

func (q *InventoryAggregate) PreInventoryVariantForVoucher(ctx context.Context, overStock bool, args *inventory.CreateInventoryVoucherArgs) (totalAmount int32, listInventoryVariants []*inventory.InventoryVariant, err error) {

	totalAmount = 0
	var inventoryvariant *inventory.InventoryVariant

	// Check have been existed variant_id in database table inventory_variant
	for _, value := range args.Lines {
		if errValidate := validateInventoryVoucherItem(value); errValidate != nil {
			return 0, nil, errValidate
		}
		totalAmount = totalAmount + value.Price*value.Quantity
		inventoryvariant, err = q.InventoryStore(ctx).ShopID(args.ShopID).VariantID(value.VariantID).Get()
		if err != nil && cm.ErrorCode(err) == cm.NotFound {

			// Create InventoryVariant follow variant_id if it have not been exit
			err = q.CreateInventoryVariant(ctx, &inventory.CreateInventoryVariantArgs{
				ShopID:    args.ShopID,
				VariantID: value.VariantID,
			})
		}
		if err != nil && cm.ErrorCode(err) != cm.NotFound {
			return 0, nil, err
		}

		// if InventoryVoucher is type 'out' move InventoryVariant quantity from onhand -> picked
		if args.Type == inventory.InventoryVoucherTypeOut {
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
			err = q.CreateInventoryVariant(ctx, &inventory.CreateInventoryVariantArgs{
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
		if args.Type == inventory.InventoryVoucherTypeOut {
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
	dbResult, err := q.InventoryVoucherStore(ctx).ShopID(args.ShopID).ID(args.ID).Get()
	if err != nil {
		return nil, err
	}
	if dbResult.Status != etop.S3Zero {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "This inventory is already confirmed or cancelled")
	}
	if dbResult.Type == inventory.InventoryVoucherTypeOut {
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
	var totalAmount int32 = 0
	for _, value := range args.Lines {
		if errValidate := validateInventoryVoucherItem(value); errValidate != nil {
			return nil, errValidate
		}
		totalAmount = totalAmount + value.Quantity*value.Price
	}
	if args.TotalAmount != totalAmount {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Tổng giá trị phiếu không hợp lệ")
	}

	updateInventoryCore := convert.ApplyUpdateInventoryVoucher(args, dbResult)
	if args.TraderID.Valid {
		if args.TraderID.Int64 != updateInventoryCore.TraderID {
			err := q.validateTrader(ctx, updateInventoryCore.ShopID, updateInventoryCore)
			if err != nil {
				return nil, err
			}
		}
	}
	err = q.InventoryVoucherStore(ctx).ShopID(args.ShopID).ID(args.ID).UpdateInventoryVoucherAll(updateInventoryCore)
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
	var listVariantID []int64
	var err error
	linesCheckin, linesCheckout, listVariantID, err = q.DevideInOutInventoryVoucher(ctx, args)
	if err != nil {
		return nil, err
	}
	if args.Title == "" {
		args.Title = "Phiếu cân bằng kho"
	}
	var inventoryVoucherInID int64
	if len(linesCheckin) > 0 {
		inventoryVoucherInID, err = q.CreateVoucherForAdjustInventoryQuantity(ctx, overStock, args, linesCheckin, inventory.InventoryVoucherTypeIn)
		if err != nil {
			return nil, err
		}
	}
	var inventoryVoucherOutID int64
	if len(linesCheckout) > 0 {
		inventoryVoucherOutID, err = q.CreateVoucherForAdjustInventoryQuantity(ctx, overStock, args, linesCheckout, inventory.InventoryVoucherTypeOut)
		if err != nil {
			return nil, err
		}
	}

	inventoryVouchers, err := q.InventoryVoucherStore(ctx).ShopID(args.ShopID).IDs(inventoryVoucherInID, inventoryVoucherOutID).ListInventoryVoucherDB()
	if err != nil {
		return nil, err
	}
	resultUpdate, err := q.InventoryStore(ctx).ShopID(args.ShopID).VariantIDs(listVariantID...).ListInventoryDB()
	if err != nil {
		return nil, err
	}
	return &inventory.AdjustInventoryQuantityRespone{
		InventoryVariants: convert.InventoryVariantsFromModel(resultUpdate),
		InventoryVouchers: convert.InventoryVouchersFromModel(inventoryVouchers),
	}, nil
}

func (q *InventoryAggregate) DevideInOutInventoryVoucher(ctx context.Context,
	args *inventory.AdjustInventoryQuantityArgs) ([]*inventory.InventoryVoucherItem,
	[]*inventory.InventoryVoucherItem,
	[]int64, error) {
	var listVariantID []int64
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
				Price:     value.PurchasePrice,
				Quantity:  value.QuantitySummary,
			})
			continue
		}
		if err != nil {
			return nil, nil, []int64{}, err
		}
		if value.QuantitySummary > (result.QuantityOnHand + result.QuantityPicked) {
			linesCheckin = append(linesCheckin, &inventory.InventoryVoucherItem{
				VariantID: value.VariantID,
				Price:     result.PurchasePrice,
				Quantity:  value.QuantitySummary - (result.QuantityOnHand + result.QuantityPicked),
			})
		} else if value.QuantitySummary < (result.QuantityOnHand + result.QuantityPicked) {
			linesCheckout = append(linesCheckout, &inventory.InventoryVoucherItem{
				VariantID: value.VariantID,
				Price:     result.PurchasePrice,
				Quantity:  (result.QuantityOnHand + result.QuantityPicked) - value.QuantitySummary,
			})
		}
	}
	return linesCheckin, linesCheckout, listVariantID, nil
}

func (q *InventoryAggregate) CreateVoucherForAdjustInventoryQuantity(ctx context.Context, overStock bool, info *inventory.AdjustInventoryQuantityArgs,
	lines []*inventory.InventoryVoucherItem,
	typeVoucher inventory.InventoryVoucherType) (int64, error) {
	var totalValue int32 = 0
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
	if inventoryVoucher.Status != etop.S3Zero {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Inventory voucher already confirmed or cancelled")
	}
	for _, value := range inventoryVoucher.Lines {
		var data *inventory.InventoryVariant
		data, err = q.InventoryStore(ctx).ShopID(args.ShopID).VariantID(value.VariantID).Get()
		if err != nil {
			return nil, err
		}
		if inventoryVoucher.Type == string(inventory.InventoryVoucherTypeOut) {
			data.QuantityPicked = data.QuantityPicked - value.Quantity
		} else if inventoryVoucher.Type == string(inventory.InventoryVoucherTypeIn) {
			if inventoryVoucher.TraderID != 0 {
				if data.QuantityOnHand < 0 {
					data.PurchasePrice = value.Price
				} else {
					data.PurchasePrice = AvgValue(data.PurchasePrice, value.Price, data.QuantityOnHand, value.Quantity)
				}
			}
			data.QuantityOnHand = data.QuantityOnHand + value.Quantity
		}
		err = q.InventoryStore(ctx).VariantID(value.VariantID).ShopID(args.ShopID).UpdateInventoryVariantAll(data)
		if err != nil {
			return nil, err
		}
	}
	inventoryVoucher.Status = etop.S3Positive
	inventoryVoucher.ConfirmedAt = time.Now()

	err = q.InventoryVoucherStore(ctx).ShopID(args.ShopID).ID(args.ID).UpdateInventoryVoucherAllDB(inventoryVoucher)
	if err != nil {
		return nil, err
	}
	return q.InventoryVoucherStore(ctx).ShopID(args.ShopID).ID(args.ID).Get()
}

func AvgValue(value1 int32, value2 int32, quantity1 int32, quantity2 int32) int32 {
	if quantity1+quantity2 == 0 {
		return 0
	}
	return (value1 + value2) / (quantity1 + quantity2)
}

func (q *InventoryAggregate) CancelInventoryVoucher(ctx context.Context, args *inventory.CancelInventoryVoucherArgs) (*inventory.InventoryVoucher, error) {
	if args.ShopID == 0 || args.ID == 0 || args.Reason == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing value requirement")
	}
	inventoryVoucher, err := q.InventoryVoucherStore(ctx).ShopID(args.ShopID).ID(args.ID).GetDB()
	if err != nil {
		return nil, err
	}
	if inventoryVoucher.Status != etop.S3Zero {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Inventory voucher already confirmed or cancelled")
	}
	if inventoryVoucher.Type == string(inventory.InventoryVoucherTypeOut) {
		for _, value := range inventoryVoucher.Lines {
			var data *inventory.InventoryVariant
			data, err = q.InventoryStore(ctx).ShopID(args.ShopID).VariantID(value.VariantID).Get()
			if err != nil {
				return nil, err
			}
			data.PurchasePrice = AvgValue(data.PurchasePrice, value.Price, data.QuantityOnHand, value.Quantity)
			data.QuantityPicked = data.QuantityPicked - value.Quantity
			data.QuantityOnHand = data.QuantityOnHand + value.Quantity

			err = q.InventoryStore(ctx).VariantID(value.VariantID).ShopID(args.ShopID).UpdateInventoryVariantAll(data)
			if err != nil {
				return nil, err
			}
		}
	}
	inventoryVoucher.Status = etop.S3Negative
	inventoryVoucher.CancelledAt = time.Now()
	inventoryVoucher.CancelReason = args.Reason
	err = q.InventoryVoucherStore(ctx).ShopID(args.ShopID).ID(args.ID).UpdateInventoryVoucherAllDB(inventoryVoucher)
	if err != nil {
		return nil, err
	}
	inventoryVoucherConfirmed, err := q.InventoryVoucherStore(ctx).ShopID(args.ShopID).ID(args.ID).Get()
	return inventoryVoucherConfirmed, err
}

func (q *InventoryAggregate) CreateInventoryVariant(ctx context.Context, args *inventory.CreateInventoryVariantArgs) error {
	if args.ShopID == 0 && args.VariantID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing value requirement")
	}
	_, err := q.InventoryStore(ctx).ShopID(args.ShopID).VariantID(args.VariantID).Get()
	if err != nil && cm.ErrorCode(err) == cm.NotFound {
		err = q.InventoryStore(ctx).Create(&model.InventoryVariant{
			ShopID:         args.ShopID,
			VariantID:      args.VariantID,
			QuantityOnHand: 0,
			QuantityPicked: 0,
			PurchasePrice:  0,
		})
		return err
	}
	if err != nil {
		return err
	}
	return nil
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
	case inventory.RefTypeOrder:
		if inventoryVoucher.Type != inventory.InventoryVoucherTypeOut {
			return cm.Error(cm.InvalidArgument, "'type' không đúng. Bán hàng chỉ có thể là 'out'", nil)
		}
		inventoryVoucher.RefName = inventory.RefNameOrder
	case inventory.RefTypeStockTake:
		if inventoryVoucher.Type != inventory.InventoryVoucherTypeOut && inventoryVoucher.Type != inventory.InventoryVoucherTypeIn {
			return cm.Error(cm.InvalidArgument, "'type' không đúng.Kiểm kho chỉ có thể là 'in' hoặc 'out'", nil)
		}
		inventoryVoucher.RefName = inventory.RefNameStockTake
	case inventory.RefTypePurchaseOrder:
		if inventoryVoucher.Type != inventory.InventoryVoucherTypeIn {
			return cm.Error(cm.InvalidArgument, "'type' không đúng.Nhập hàng chỉ có thể là 'in'", nil)
		}
		inventoryVoucher.RefName = inventory.RefNamePurchaseOrder
	case inventory.RefTypeReturns:
		if inventoryVoucher.Type != inventory.InventoryVoucherTypeIn {
			return cm.Error(cm.InvalidArgument, "'type' không đúng.Trả hàng chỉ có thể là 'in'", nil)
		}
		inventoryVoucher.RefName = inventory.RefNameReturns
	case "":
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
	var inventoryVariantIDs []int64
	var inventoryVoucherIn []*inventory.InventoryVoucherItem
	var inventoryVoucherOut []*inventory.InventoryVoucherItem
	for _, value := range args.Variants {
		inventoryVariantIDs = append(inventoryVariantIDs, value.VariantID)
	}
	listVariant, err := q.InventoryStore(ctx).ShopID(args.ShopID).VariantIDs(inventoryVariantIDs...).ListInventory()
	if err != nil {
		return nil, err
	}
	var mapInventoryVariantInfo = make(map[int64]*inventory.InventoryVariant)
	for _, value := range listVariant {
		mapInventoryVariantInfo[value.VariantID] = value
	}

	var totalAmountIn int32 = 0
	var totalAmountOut int32 = 0
	for _, value := range args.Variants {
		if value.QuantityChange > 0 {
			inventoryVoucherItem := &inventory.InventoryVoucherItem{
				VariantID: value.VariantID,
				Quantity:  value.QuantityChange,
			}
			if mapInventoryVariantInfo[value.VariantID] != nil {
				inventoryVoucherItem.Price = mapInventoryVariantInfo[value.VariantID].PurchasePrice
			}
			inventoryVoucherIn = append(inventoryVoucherIn, inventoryVoucherItem)
			totalAmountIn += value.QuantityChange * inventoryVoucherItem.Price
		} else {
			inventoryVoucherItem := &inventory.InventoryVoucherItem{
				VariantID: value.VariantID,
				Quantity:  value.QuantityChange * -1,
			}
			if mapInventoryVariantInfo[value.VariantID] != nil {
				inventoryVoucherItem.Price = mapInventoryVariantInfo[value.VariantID].PurchasePrice
			}
			inventoryVoucherOut = append(inventoryVoucherOut, inventoryVoucherItem)
			totalAmountOut += value.QuantityChange * inventoryVoucherItem.Price
		}

	}
	totalAmountOut = totalAmountOut * -1
	var typeIn = &inventory.InventoryVoucher{}
	if len(inventoryVoucherIn) != 0 {
		typeIn, err = q.CreateInventoryVoucher(ctx, args.Overstock, &inventory.CreateInventoryVoucherArgs{
			ShopID:      args.ShopID,
			CreatedBy:   args.CreatedBy,
			Title:       args.Title,
			RefID:       args.RefID,
			RefType:     args.RefType,
			RefName:     args.RefName,
			TraderID:    0,
			TotalAmount: totalAmountIn,
			Type:        "in",
			Note:        args.Note,
			Lines:       inventoryVoucherIn,
		})
		if err != nil {
			return nil, err
		}
	}
	var typeOut = &inventory.InventoryVoucher{}
	if len(inventoryVoucherOut) != 0 {
		typeOut, err = q.CreateInventoryVoucher(ctx, args.Overstock, &inventory.CreateInventoryVoucherArgs{
			ShopID:      args.ShopID,
			CreatedBy:   args.CreatedBy,
			Title:       args.Title,
			RefID:       args.RefID,
			RefType:     args.RefType,
			RefName:     args.RefName,
			TraderID:    0,
			TotalAmount: totalAmountIn,
			Type:        "out",
			Note:        args.Note,
			Lines:       inventoryVoucherOut,
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
