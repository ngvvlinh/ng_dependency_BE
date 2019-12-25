package shop

import (
	"context"

	"etop.vn/api/main/inventory"
	"etop.vn/api/meta"
	"etop.vn/api/shopping/tradering"
	"etop.vn/api/top/int/shop"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/cmapi"
	"etop.vn/capi/dot"
)

func (s *InventoryService) CreateInventoryVoucher(ctx context.Context, q *CreateInventoryVoucherEndpoint) error {
	cmd := &inventory.CreateInventoryVoucherByReferenceCommand{
		RefType:   inventory.InventoryRefType(q.RefType),
		RefID:     q.RefId,
		ShopID:    q.Context.Shop.ID,
		UserID:    q.Context.UserID,
		Type:      inventory.InventoryVoucherType(q.Type),
		OverStock: q.Context.Shop.InventoryOverstock.Apply(true),
	}
	err := inventoryAggregate.Dispatch(ctx, cmd)
	if err != nil {
		return err
	}
	q.Result = &shop.CreateInventoryVoucherResponse{
		InventoryVouchers: PbShopInventoryVouchers(cmd.Result),
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
	q.Result = &shop.ConfirmInventoryVoucherResponse{
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
	q.Result = &shop.CancelInventoryVoucherResponse{
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
		Title:       q.Title,
		ID:          q.Id,
		ShopID:      shopID,
		TotalAmount: q.TotalAmount,
		UpdatedBy:   userID,
		TraderID:    q.TraderId,
		Note:        q.Note,
		Lines:       items,
	}
	if err := inventoryAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &shop.UpdateInventoryVoucherResponse{
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
		Overstock: inventoryOverstock.Apply(true),
		ShopID:    shopID,
		Lines:     items,
		UserID:    userID,
		Note:      q.Note,
	}
	if err := inventoryAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &shop.AdjustInventoryQuantityResponse{
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
	q.Result = &shop.GetInventoryVariantsResponse{
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
	q.Result = &shop.GetInventoryVouchersByReferenceResponse{
		InventoryVouchers: PbShopInventoryVouchers(query.Result.InventoryVouchers),
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
	q.Result = &shop.GetInventoryVariantsResponse{
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
	paging := cmapi.CMPaging(q.Paging)
	query := &inventory.GetInventoryVouchersQuery{
		ShopID:  shopID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(q.Filters),
	}
	if err := inventoryQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	inventoryVouchers, err := s.checkValidateListTrader(ctx, shopID, query.Result.InventoryVoucher)
	if err != nil {
		return err
	}
	q.Result = &shop.GetInventoryVouchersResponse{
		InventoryVouchers: inventoryVouchers,
	}
	return nil
}

func (s *InventoryService) checkValidateListTrader(ctx context.Context, shopID dot.ID, inventoryVouchers []*inventory.InventoryVoucher) (result []*shop.InventoryVoucher, err error) {
	if inventoryVouchers == nil {
		return result, err
	}
	traderIDs := make([]dot.ID, 0, len(inventoryVouchers))
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
	var mapTraderValidate = make(map[dot.ID]bool)
	for _, trader := range traders {
		mapTraderValidate[trader.ID] = true
	}
	for _, inventoryVoucher := range inventoryVouchers {
		inventory := PbShopInventoryVoucher(inventoryVoucher)
		if inventoryVoucher.TraderID == 0 || inventory.Trader == nil {
			result = append(result, inventory)
			continue
		}
		if !mapTraderValidate[inventoryVoucher.TraderID] {
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
	q.Result = &shop.GetInventoryVouchersResponse{
		InventoryVouchers: PbShopInventoryVouchers(query.Result.InventoryVoucher),
	}
	return nil
}

func (s *InventoryService) UpdateInventoryVariantCostPrice(ctx context.Context, q *UpdateInventoryVariantCostPriceEndpoint) error {
	cmd := &inventory.UpdateInventoryVariantCostPriceCommand{
		ShopID:    q.Context.Shop.ID,
		VariantID: q.VariantId,
		CostPrice: q.CostPrice,
	}
	err := inventoryAggregate.Dispatch(ctx, cmd)
	if err != nil {
		return err
	}
	q.Result = &shop.UpdateInventoryVariantCostPriceResponse{
		InventoryVariant: PbInventory(cmd.Result),
	}
	return nil
}
