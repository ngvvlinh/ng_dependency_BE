package aggregatex

import (
	"context"
	"fmt"
	"time"

	"o.o/api/main/catalog"
	"o.o/api/main/catalog/types"
	"o.o/api/main/inventory"
	"o.o/api/main/stocktaking"
	"o.o/api/top/types/etc/inventory_auto"
	"o.o/api/top/types/etc/inventory_type"
	"o.o/api/top/types/etc/inventory_voucher_ref"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/stocktake_type"
	"o.o/backend/com/main/inventory/util"
	ordermodelx "o.o/backend/com/main/ordering/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/capi/dot"
)

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

	err = q.InventoryVoucherStore(ctx).ShopID(args.ShopID).ID(args.ID).UpdateInventoryVoucherAll(inventoryVoucher)
	if err != nil {
		return nil, err
	}
	result, err := q.InventoryVoucherStore(ctx).ShopID(args.ShopID).ID(args.ID).Get()
	if err != nil {
		return nil, err
	}
	return q.populateInventoryVoucher(ctx, result)
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
	inventoryVouchers, err = q.populateInventoryVouchers(ctx, inventoryVouchers)
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
	result, err := q.InventoryVoucherStore(ctx).ShopID(args.ShopID).ID(args.ID).Get()
	if err != nil {
		return nil, err
	}
	return q.populateInventoryVoucher(ctx, result)
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
	if err != nil {
		return nil, err
	}
	return q.populateInventoryVoucher(ctx, inventoryVoucherConfirmed)
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
			RefCode:   args.RefCode,
			Type:      inventory_type.In,
			Lines:     inventoryVoucherIn,
		})
		if err != nil {
			return nil, err
		}
		typeIn, err = q.populateInventoryVoucher(ctx, typeIn)
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
			RefCode:   args.RefCode,
			Type:      inventory_type.Out,
			Lines:     inventoryVoucherOut,
		})
		if err != nil {
			return nil, err
		}
		typeOut, err = q.populateInventoryVoucher(ctx, typeOut)
		if err != nil {
			return nil, err
		}
	}
	return &inventory.CreateInventoryVoucherByQuantityChangeResponse{
		TypeIn:  typeIn,
		TypeOut: typeOut,
	}, nil
}

func (q *InventoryAggregate) CreateInventoryVoucherByReference(ctx context.Context, args *inventory.CreateInventoryVoucherByReferenceArgs) ([]*inventory.InventoryVoucher, error) {
	switch args.RefType {
	case inventory_voucher_ref.StockTake:
		return q.CreateInventoryVoucherByStockTake(ctx, args)
	case inventory_voucher_ref.Order:
		return q.CreateInventoryVoucherByOrder(ctx, args)
	default:
		return nil, cm.Error(cm.InvalidArgument, "wrong ref_type", nil)
	}
}

func (q *InventoryAggregate) CreateInventoryVoucherByOrder(ctx context.Context, args *inventory.CreateInventoryVoucherByReferenceArgs) ([]*inventory.InventoryVoucher, error) {
	var items []*inventory.InventoryVoucherItem
	// check order_id exit
	queryOrder := &ordermodelx.GetOrderQuery{
		OrderID: args.RefID,
		ShopID:  args.ShopID,
	}
	if err := q.OrderStore.GetOrder(ctx, queryOrder); err != nil {
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
		inventoryVoucherCreateRequest.Title = "Xuất kho khi bán hàng"
	case inventory_type.In:
		inventoryVoucherCreateRequest.Title = "Nhập kho khi Hủy bán hàng"
	}
	createResult, err := q.CreateInventoryVoucher(ctx, args.OverStock, inventoryVoucherCreateRequest)
	if err != nil {
		return nil, err
	}
	var listInventoryVoucher []*inventory.InventoryVoucher
	listInventoryVoucher = append(listInventoryVoucher, createResult)
	return q.populateInventoryVouchers(ctx, listInventoryVoucher)
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
		Title:     "Phiếu kiểm kho",
		RefCode:   queryStocktake.Result.Code,
		Overstock: args.OverStock,
		CreatedBy: args.UserID,
		Lines:     inventoryVariantChange,
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
	return q.populateInventoryVouchers(ctx, listInventoryVoucher)
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
		if isRollBack(value) == true {
			return &inventory.CancelInventoryByRefIDResponse{
				InventoryVouchers: inventoryVouchers,
			}, nil
		}
	}

	for _, value := range inventoryVouchersData {
		switch value.Status {
		case status3.P:
			var typeInventoryVoucher inventory_type.InventoryVoucherType
			var result *inventory.InventoryVoucher
			if value.Type == inventory_type.Out {
				typeInventoryVoucher = inventory_type.In
			} else {
				typeInventoryVoucher = inventory_type.Out
			}
			newVoucher := &inventory.CreateInventoryVoucherArgs{
				ShopID:      value.ShopID,
				CreatedBy:   args.UpdateBy,
				Title:       "Phiếu xuất nhập kho",
				RefID:       value.RefID,
				RefType:     value.RefType,
				RefCode:     value.RefCode,
				TraderID:    value.TraderID,
				TotalAmount: value.TotalAmount,
				Type:        typeInventoryVoucher,
				Lines:       value.Lines,
			}

			if args.RefType == inventory_voucher_ref.PurchaseOrder {
				var variantIDs []dot.ID
				var linesTemp []*inventory.InventoryVoucherItem
				mapCheckVariant := make(map[dot.ID]bool)
				for _, line := range newVoucher.Lines {
					variantIDs = append(variantIDs, line.VariantID)
				}
				var queryVariants = &catalog.ListShopVariantsByIDsQuery{
					ShopID: newVoucher.ShopID,
					IDs:    variantIDs,
				}
				if err := q.CatalogQuery.Dispatch(ctx, queryVariants); err != nil {
					return nil, err
				}
				for _, variant := range queryVariants.Result.Variants {
					mapCheckVariant[variant.VariantID] = true
				}
				for _, line := range newVoucher.Lines {
					if mapCheckVariant[line.VariantID] {
						linesTemp = append(linesTemp, line)
					}
				}
				if len(linesTemp) != 0 {
					newVoucher.Lines = linesTemp
					result, err = q.CreateInventoryVoucher(ctx, args.InventoryOverStock, newVoucher)
					if err != nil {
						return nil, err
					}
				}

			} else {
				result, err = q.CreateInventoryVoucher(ctx, args.InventoryOverStock, newVoucher)
				if err != nil {
					return nil, err
				}
			}
			if result != nil {
				inventoryVouchers = append(inventoryVouchers, result)
			}
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
	if len(inventoryVouchers) != 0 {
		inventoryVouchers, err = q.populateInventoryVouchers(ctx, inventoryVouchers)
		if err != nil {
			return nil, err
		}
	}
	return &inventory.CancelInventoryByRefIDResponse{
		InventoryVouchers: inventoryVouchers,
	}, nil
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

func (q *InventoryAggregate) populateInventoryVouchers(ctx context.Context, args []*inventory.InventoryVoucher) ([]*inventory.InventoryVoucher, error) {
	var stocktakeIDs []dot.ID
	for _, v := range args {
		if v.RefType == inventory_voucher_ref.StockTake {
			stocktakeIDs = append(stocktakeIDs, v.RefID)
		}
	}
	var mapStocktake = make(map[dot.ID]stocktake_type.StocktakeType)
	if len(stocktakeIDs) > 0 {
		cmdStocktake := &stocktaking.GetStocktakesByIDsQuery{
			Ids:    stocktakeIDs,
			ShopID: args[0].ShopID,
		}
		err := q.StocktakeQuery.Dispatch(ctx, cmdStocktake)
		if err != nil {
			return nil, err
		}
		for _, value := range cmdStocktake.Result {
			mapStocktake[value.ID] = value.Type
		}
	}
	return util.PopulateInventoryVouchers(args, mapStocktake)
}

func (q *InventoryAggregate) populateInventoryVoucher(ctx context.Context, arg *inventory.InventoryVoucher) (*inventory.InventoryVoucher, error) {
	var stocktake *stocktaking.ShopStocktake
	if arg.RefType == inventory_voucher_ref.StockTake {
		cmdStocktake := &stocktaking.GetStocktakeByIDQuery{
			Id:     arg.RefID,
			ShopID: arg.ShopID,
		}
		err := q.StocktakeQuery.Dispatch(ctx, cmdStocktake)
		if err != nil {
			return nil, err
		}
		stocktake = cmdStocktake.Result
	}
	return util.PopulateInventoryVoucher(arg, stocktake)
}

func isRollBack(arg *inventory.InventoryVoucher) bool {
	switch arg.RefType {
	case inventory_voucher_ref.Order:
		if arg.Type == inventory_type.In {
			return true
		}
	case inventory_voucher_ref.Refund:
		if arg.Type == inventory_type.Out {
			return true
		}
	case inventory_voucher_ref.PurchaseOrder:
		if arg.Type == inventory_type.Out {
			return true
		}
	case inventory_voucher_ref.PurchaseRefund:
		if arg.Type == inventory_type.In {
			return true
		}
	}
	return false
}
