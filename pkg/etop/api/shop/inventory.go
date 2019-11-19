package shop

import (
	"context"
	"fmt"

	"etop.vn/api/main/etop"
	"etop.vn/api/main/inventory"
	stocktake "etop.vn/api/main/stocktaking"
	"etop.vn/api/meta"
	"etop.vn/api/shopping/tradering"
	ordermodelx "etop.vn/backend/com/main/ordering/modelx"
	pbcm "etop.vn/backend/pb/common"
	pbs4 "etop.vn/backend/pb/etop/etc/status4"
	pbshop "etop.vn/backend/pb/etop/shop"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/model"
	. "etop.vn/capi/dot"
)

func (s *InventoryService) CreateInventoryVoucher(ctx context.Context, q *CreateInventoryVoucherEndpoint) error {
	if q.RefType == string(inventory.RefTypeStockTake) {
		cmd, err := PreCreateInventoryVoucherRefStocktake(ctx, q)
		if err != nil {
			return err
		}
		if err := inventoryAggregate.Dispatch(ctx, cmd); err != nil {
			return err
		}
		var inventoryVouchers []*pbshop.InventoryVoucher
		if cmd.Result.TypeIn.ID != 0 {
			inventoryVouchers = append(inventoryVouchers, PbShopInventoryVoucher(cmd.Result.TypeIn))
		}
		if cmd.Result.TypeOut.ID != 0 {
			inventoryVouchers = append(inventoryVouchers, PbShopInventoryVoucher(cmd.Result.TypeOut))
		}
		q.Result = &pbshop.CreateInventoryVoucherResponse{
			InventoryVouchers: inventoryVouchers,
		}
		return nil
	}
	{
		cmd, err := PreCreateInventoryVoucher(ctx, q)
		if err != nil {
			return err
		}
		if err := inventoryAggregate.Dispatch(ctx, cmd); err != nil {
			return err
		}
		q.Result = &pbshop.CreateInventoryVoucherResponse{
			InventoryVouchers: []*pbshop.InventoryVoucher{PbShopInventoryVoucher(cmd.Result)},
		}
	}
	return nil
}

func (s *InventoryService) ConfirmInventoryVoucher(ctx context.Context, q *ConfirmInventoryVoucherEndpoint) error {
	shopID := q.Context.Shop.ID
	userID := q.Context.UserID

	cmd := &inventory.ConfirmInventoryVoucherCommand{
		ShopID:    shopID,
		ID:        q.Id,
		UpdatedBy: userID,
	}
	if err := inventoryAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbshop.ConfirmInventoryVoucherResponse{
		InventoryVoucher: PbShopInventoryVoucher(cmd.Result),
	}
	return nil
}

func (s *InventoryService) CancelInventoryVoucher(ctx context.Context, q *CancelInventoryVoucherEndpoint) error {
	shopID := q.Context.Shop.ID
	userID := q.Context.UserID

	cmd := &inventory.CancelInventoryVoucherCommand{
		ShopID:    shopID,
		ID:        q.Id,
		UpdatedBy: userID,
		Reason:    q.Reason,
	}
	if err := inventoryAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbshop.CancelInventoryVoucherResponse{
		Inventory: PbShopInventoryVoucher(cmd.Result),
	}
	return nil
}

func (s *InventoryService) UpdateInventoryVoucher(ctx context.Context, q *UpdateInventoryVoucherEndpoint) error {
	shopID := q.Context.Shop.ID
	userID := q.Context.UserID
	var items []*inventory.InventoryVoucherItem
	for _, value := range q.Lines {
		items = append(items, &inventory.InventoryVoucherItem{
			VariantID: value.VariantId,
			Price:     value.Price,
			Quantity:  value.Quantity,
		})
	}
	cmd := &inventory.UpdateInventoryVoucherCommand{
		Title:       PString(q.Title),
		ID:          q.Id,
		ShopID:      shopID,
		TotalAmount: q.TotalAmount,
		UpdatedBy:   userID,
		TraderID:    PInt64(q.TraderId),
		Note:        PString(q.Note),
		Lines:       items,
	}
	if err := inventoryAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbshop.UpdateInventoryVoucherResponse{
		InventoryVoucher: PbShopInventoryVoucher(cmd.Result),
	}
	return nil
}

func (s *InventoryService) AdjustInventoryQuantity(ctx context.Context, q *AdjustInventoryQuantityEndpoint) error {
	shopID := q.Context.Shop.ID
	userID := q.Context.UserID
	inventoryOverstock := q.Context.Shop.InventoryOverstock
	var items []*inventory.InventoryVariant
	for _, value := range q.InventoryVariants {
		items = append(items, &inventory.InventoryVariant{
			ShopID:          shopID,
			VariantID:       value.VariantId,
			QuantityOnHand:  value.QuantityOnHand,
			QuantitySummary: value.Quantity,
			QuantityPicked:  value.QuantityPicked,
		})
	}
	cmd := &inventory.AdjustInventoryQuantityCommand{
		Overstock: cm.BoolDefault(inventoryOverstock, true),
		ShopID:    shopID,
		Lines:     items,
		UserID:    userID,
		Note:      q.Note,
	}
	if err := inventoryAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbshop.AdjustInventoryQuantityResponse{
		InventoryVariants: PbInventoryVariants(cmd.Result.InventoryVariants),
		InventoryVouchers: PbShopInventoryVouchers(cmd.Result.InventoryVouchers),
	}
	return nil
}

func (s *InventoryService) GetInventoryVariants(ctx context.Context, q *GetInventoryVariantsEndpoint) error {
	shopID := q.Context.Shop.ID
	query := &inventory.GetInventoryVariantsQuery{
		ShopID: shopID,
		Paging: &meta.Paging{
			Offset: q.Paging.Offset,
			Limit:  q.Paging.Limit,
		},
	}
	if err := inventoryQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pbshop.GetInventoryVariantsResponse{
		InventoryVariants: PbInventoryVariants(query.Result.InventoryVariants),
	}
	return nil
}

func (s *InventoryService) GetInventoryVouchersByReference(ctx context.Context, q *GetInventoryVouchersByReferenceEndpoint) error {
	shopID := q.Context.Shop.ID
	query := &inventory.GetInventoryVoucherByReferenceQuery{
		ShopID:  shopID,
		RefID:   q.RefId,
		RefType: inventory.InventoryRefType(q.RefType),
	}
	if err := inventoryQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pbshop.GetInventoryVouchersByReferenceResponse{
		InventoryVoucher: PbShopInventoryVouchers(query.Result.InventoryVouchers),
		Status:           pbs4.Pb(model.Status4(query.Result.Status)),
	}
	return nil
}

func (s *InventoryService) GetInventoryVariant(ctx context.Context, q *GetInventoryVariantEndpoint) error {
	shopID := q.Context.Shop.ID
	query := &inventory.GetInventoryVariantQuery{
		ShopID:    shopID,
		VariantID: q.VariantId,
	}
	if err := inventoryQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = PbInventory(query.Result)
	return nil
}

func (s *InventoryService) GetInventoryVariantsByVariantIDs(ctx context.Context, q *GetInventoryVariantsByVariantIDsEndpoint) error {
	shopID := q.Context.Shop.ID
	query := &inventory.GetInventoryVariantsByVariantIDsQuery{
		ShopID:     shopID,
		VariantIDs: q.VariantIds,
	}
	if err := inventoryQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pbshop.GetInventoryVariantsResponse{
		InventoryVariants: PbInventoryVariants(query.Result.InventoryVariants),
	}
	return nil
}

func (s *InventoryService) GetInventoryVoucher(ctx context.Context, q *GetInventoryVoucherEndpoint) error {
	shopID := q.Context.Shop.ID
	query := &inventory.GetInventoryVoucherQuery{
		ShopID: shopID,
		ID:     q.Id,
	}
	if err := inventoryQuery.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = PbShopInventoryVoucher(query.Result)
	if q.Result.TraderId != 0 {
		getTrader := &tradering.GetTraderByIDQuery{
			ID:     q.Result.TraderId,
			ShopID: shopID,
		}
		if err := traderQuery.Dispatch(ctx, getTrader); err != nil {
			if cm.ErrorCode(err) != cm.NotFound {
				return err
			}
			q.Result.Trader.Deleted = true
		}
	}
	return nil
}

func (s *InventoryService) GetInventoryVouchers(ctx context.Context, q *GetInventoryVouchersEndpoint) error {
	shopID := q.Context.Shop.ID
	paging := q.Paging.CMPaging()
	query := &inventory.GetInventoryVouchersQuery{
		ShopID:  shopID,
		Paging:  *paging,
		Filters: pbcm.ToFilters(q.Filters),
	}
	if err := inventoryQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	inventoryVouchers, err := s.checkValidateListTrader(ctx, shopID, query.Result.InventoryVoucher)
	if err != nil {
		return err
	}
	q.Result = &pbshop.GetInventoryVouchersResponse{
		InventoryVouchers: inventoryVouchers,
	}
	return nil
}

func (s *InventoryService) checkValidateListTrader(ctx context.Context, shopID int64, inventoryVouchers []*inventory.InventoryVoucher) (result []*pbshop.InventoryVoucher, err error) {
	if inventoryVouchers == nil {
		return result, err
	}
	traderIDs := make([]int64, 0, len(inventoryVouchers))
	for _, trader := range inventoryVouchers {
		traderIDs = append(traderIDs, trader.TraderID)
	}
	queryTraders := &tradering.ListTradersByIDsQuery{
		IDs:    traderIDs,
		ShopID: shopID,
	}
	if err := traderQuery.Dispatch(ctx, queryTraders); err != nil {
		return result, err
	}
	traders := queryTraders.Result.Traders
	var mapTraderValidate = make(map[int64]bool)
	for _, trader := range traders {
		mapTraderValidate[trader.ID] = true
	}
	for _, inventoryVoucher := range inventoryVouchers {
		inventory := PbShopInventoryVoucher(inventoryVoucher)
		if !mapTraderValidate[inventoryVoucher.TraderID] && inventoryVoucher.TraderID != 0 {
			inventory.Trader.Deleted = true
		}
		result = append(result, inventory)
	}
	return result, err
}

func (s *InventoryService) GetInventoryVouchersByIDs(ctx context.Context, q *GetInventoryVouchersByIDsEndpoint) error {
	shopID := q.Context.Shop.ID
	query := &inventory.GetInventoryVouchersByIDsQuery{
		ShopID: shopID,
		IDs:    q.Ids,
	}
	if err := inventoryQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &pbshop.GetInventoryVouchersResponse{
		InventoryVouchers: PbShopInventoryVouchers(query.Result.InventoryVoucher),
	}
	return nil
}

func PreCreateInventoryVoucherRefStocktake(ctx context.Context, q *CreateInventoryVoucherEndpoint) (*inventory.CreateInventoryVoucherByQuantityChangeCommand, error) {
	// check order_id exit
	queryStocktake := &stocktake.GetStocktakeByIDQuery{
		Id:     q.RefId,
		ShopID: q.Context.Shop.ID,
	}
	if err := StocktakeQuery.Dispatch(ctx, queryStocktake); err != nil {
		return nil, err
	}
	if queryStocktake.Result.Status != etop.S3Positive {
		return nil, cm.Error(cm.InvalidArgument, "không thể tạo phiếu kiểm kho cho stocktake chưa được xác nhận.", nil)
	}
	// GET info and put it to cmd
	var inventoryVariantChange []*inventory.InventoryVariantQuantityChange
	for _, value := range queryStocktake.Result.Lines {
		inventoryVariantChange = append(inventoryVariantChange, &inventory.InventoryVariantQuantityChange{
			VariantID:      value.VariantID,
			QuantityChange: value.NewQuantity - value.OldQuantity,
		})
	}
	cmd := &inventory.CreateInventoryVoucherByQuantityChangeCommand{
		ShopID:    q.Context.Shop.ID,
		RefID:     q.RefId,
		RefType:   inventory.RefTypeStockTake,
		RefName:   inventory.RefNameStockTake,
		Title:     "Phiếu kiểm kho",
		Overstock: cm.BoolDefault(q.Context.Shop.InventoryOverstock, true),
		CreatedBy: q.Context.UserID,
		Variants:  inventoryVariantChange,
	}
	cmd.Note = fmt.Sprintf("Tạo phiếu xuất nhập kho theo phiếu kiểm kho mã %v", queryStocktake.Result.ID)
	return cmd, nil
}

func PreCreateInventoryVoucherRefOrder(ctx context.Context, cmd *inventory.CreateInventoryVoucherCommand) error {
	var items []*inventory.InventoryVoucherItem

	// check order_id exit
	queryOrder := &ordermodelx.GetOrderQuery{
		OrderID: cmd.RefID,
		ShopID:  cmd.ShopID,
	}
	if err := bus.Dispatch(ctx, queryOrder); err != nil {
		return err
	}
	// GET info and put it to cmd
	for _, value := range queryOrder.Result.Order.Lines {
		items = append(items, &inventory.InventoryVoucherItem{
			VariantID: value.VariantID,
			Quantity:  int32(value.Quantity),
		})
	}
	cmd.Title = "Xuất kho khi bán hàng"
	cmd.Lines = items
	cmd.Type = "out"
	cmd.TraderID = queryOrder.Result.Order.CustomerID
	cmd.Note = fmt.Sprintf("Tạo phiếu nhập kho theo đơn đặt hàng mã %v", queryOrder.Result.Order.ID)
	return nil
}

func PreCreateInventoryVoucher(ctx context.Context, q *CreateInventoryVoucherEndpoint) (*inventory.CreateInventoryVoucherCommand, error) {
	shopID := q.Context.Shop.ID
	inventoryOverstock := q.Context.Shop.InventoryOverstock

	cmd := &inventory.CreateInventoryVoucherCommand{
		ShopID:    shopID,
		Overstock: cm.BoolDefault(inventoryOverstock, true),
		RefType:   inventory.InventoryRefType(q.RefType),
		RefID:     q.RefId,
		CreatedBy: q.Context.UserID,
	}
	// Check ref ID here "order", "purchaseorder", "stocktake", "return"
	// modify value flow reftype
	switch cmd.RefType {
	case inventory.RefTypeOrder:
		if err := PreCreateInventoryVoucherRefOrder(ctx, cmd); err != nil {
			return nil, err
		}
	case inventory.RefTypeReturns:
		return nil, cm.Error(cm.InvalidArgument, "not support ref_type = 'return' now", nil)
	case inventory.RefTypePurchaseOrder:
		return nil, cm.Error(cm.InvalidArgument, "not support ref_type = 'perchaseorder' now", nil)
	default:
		return nil, cm.Error(cm.InvalidArgument, "wrong ref_type", nil)
	}
	return cmd, nil
}
