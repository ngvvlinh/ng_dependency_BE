package inventory

import (
	"context"

	"o.o/api/main/inventory"
	"o.o/api/meta"
	"o.o/api/shopping/tradering"
	"o.o/api/top/int/shop"
	api "o.o/api/top/int/shop"
	pbcm "o.o/api/top/types/common"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/capi/dot"
)

type InventoryService struct {
	session.Session

	// TraderQuery may be nil for fabo
	TraderQuery tradering.QueryBus

	InventoryAggr  inventory.CommandBus
	InventoryQuery inventory.QueryBus
}

func (s *InventoryService) Clone() api.InventoryService { res := *s; return &res }

func (s *InventoryService) CreateInventoryVoucher(ctx context.Context, q *api.CreateInventoryVoucherRequest) (*api.CreateInventoryVoucherResponse, error) {
	cmd := &inventory.CreateInventoryVoucherByReferenceCommand{
		RefType:   q.RefType,
		RefID:     q.RefId,
		ShopID:    s.SS.Shop().ID,
		UserID:    s.SS.Claim().UserID,
		Type:      q.Type,
		OverStock: s.SS.Shop().InventoryOverstock.Apply(true),
	}
	err := s.InventoryAggr.Dispatch(ctx, cmd)
	if err != nil {
		return nil, err
	}
	result := &api.CreateInventoryVoucherResponse{
		InventoryVouchers: PbShopInventoryVouchers(cmd.Result),
	}
	return result, nil
}

func (s *InventoryService) ConfirmInventoryVoucher(ctx context.Context, q *api.ConfirmInventoryVoucherRequest) (*api.ConfirmInventoryVoucherResponse, error) {
	shopID := s.SS.Shop().ID
	userID := s.SS.Claim().UserID

	cmd := &inventory.ConfirmInventoryVoucherCommand{
		ShopID:    shopID,
		ID:        q.Id,
		UpdatedBy: userID,
	}
	if err := s.InventoryAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &api.ConfirmInventoryVoucherResponse{
		InventoryVoucher: PbShopInventoryVoucher(cmd.Result),
	}
	return result, nil
}

func (s *InventoryService) CancelInventoryVoucher(ctx context.Context, q *api.CancelInventoryVoucherRequest) (*api.CancelInventoryVoucherResponse, error) {
	shopID := s.SS.Shop().ID
	userID := s.SS.Claim().UserID

	cmd := &inventory.CancelInventoryVoucherCommand{
		ShopID:       shopID,
		ID:           q.Id,
		UpdatedBy:    userID,
		CancelReason: q.CancelReason,
	}
	if err := s.InventoryAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &api.CancelInventoryVoucherResponse{
		Inventory: PbShopInventoryVoucher(cmd.Result),
	}
	return result, nil
}

func (s *InventoryService) UpdateInventoryVoucher(ctx context.Context, q *api.UpdateInventoryVoucherRequest) (*api.UpdateInventoryVoucherResponse, error) {
	shopID := s.SS.Shop().ID
	userID := s.SS.Claim().UserID
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
	if err := s.InventoryAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &api.UpdateInventoryVoucherResponse{
		InventoryVoucher: PbShopInventoryVoucher(cmd.Result),
	}
	return result, nil
}

func (s *InventoryService) AdjustInventoryQuantity(ctx context.Context, q *api.AdjustInventoryQuantityRequest) (*api.AdjustInventoryQuantityResponse, error) {
	shopID := s.SS.Shop().ID
	userID := s.SS.Claim().UserID
	inventoryOverstock := s.SS.Shop().InventoryOverstock
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
	if err := s.InventoryAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &api.AdjustInventoryQuantityResponse{
		InventoryVariants: PbInventoryVariants(cmd.Result.InventoryVariants),
		InventoryVouchers: PbShopInventoryVouchers(cmd.Result.InventoryVouchers),
	}
	return result, nil
}

func (s *InventoryService) GetInventoryVariants(ctx context.Context, q *api.GetInventoryVariantsRequest) (*api.GetInventoryVariantsResponse, error) {
	shopID := s.SS.Shop().ID
	query := &inventory.GetInventoryVariantsQuery{
		ShopID: shopID,
		Paging: &meta.Paging{
			Offset: q.Paging.Offset,
			Limit:  q.Paging.Limit,
		},
	}
	if err := s.InventoryQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &api.GetInventoryVariantsResponse{
		InventoryVariants: PbInventoryVariants(query.Result.InventoryVariants),
	}
	return result, nil
}

func (s *InventoryService) GetInventoryVouchersByReference(ctx context.Context, q *api.GetInventoryVouchersByReferenceRequest) (*api.GetInventoryVouchersByReferenceResponse, error) {
	shopID := s.SS.Shop().ID
	query := &inventory.GetInventoryVoucherByReferenceQuery{
		ShopID:  shopID,
		RefID:   q.RefId,
		RefType: q.RefType,
	}
	if err := s.InventoryQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &api.GetInventoryVouchersByReferenceResponse{
		InventoryVouchers: PbShopInventoryVouchers(query.Result.InventoryVouchers),
	}
	return result, nil
}

func (s *InventoryService) GetInventoryVariant(ctx context.Context, q *api.GetInventoryVariantRequest) (*api.InventoryVariant, error) {
	shopID := s.SS.Shop().ID
	query := &inventory.GetInventoryVariantQuery{
		ShopID:    shopID,
		VariantID: q.VariantId,
	}
	if err := s.InventoryQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := PbInventory(query.Result)
	return result, nil
}

func (s *InventoryService) GetInventoryVariantsByVariantIDs(ctx context.Context, q *api.GetInventoryVariantsByVariantIDsRequest) (*api.GetInventoryVariantsResponse, error) {
	shopID := s.SS.Shop().ID
	query := &inventory.GetInventoryVariantsByVariantIDsQuery{
		ShopID:     shopID,
		VariantIDs: q.VariantIds,
	}
	if err := s.InventoryQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &api.GetInventoryVariantsResponse{
		InventoryVariants: PbInventoryVariants(query.Result.InventoryVariants),
	}
	return result, nil
}

func (s *InventoryService) GetInventoryVoucher(ctx context.Context, q *pbcm.IDRequest) (*api.InventoryVoucher, error) {
	shopID := s.SS.Shop().ID
	query := &inventory.GetInventoryVoucherQuery{
		ShopID: shopID,
		ID:     q.Id,
	}
	if err := s.InventoryQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}

	result := PbShopInventoryVoucher(query.Result)
	if result.TraderId != 0 {
		getTrader := &tradering.GetTraderByIDQuery{
			ID:     result.TraderId,
			ShopID: shopID,
		}
		if err := s.TraderQuery.Dispatch(ctx, getTrader); err != nil {
			if cm.ErrorCode(err) != cm.NotFound {
				return nil, err
			}
			result.Trader.Deleted = true
		}
	}
	return result, nil
}

func (s *InventoryService) GetInventoryVouchers(ctx context.Context, q *api.GetInventoryVouchersRequest) (*api.GetInventoryVouchersResponse, error) {
	shopID := s.SS.Shop().ID
	paging := cmapi.CMPaging(q.Paging)
	query := &inventory.GetInventoryVouchersQuery{
		ShopID:  shopID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(q.Filters),
	}
	if err := s.InventoryQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	inventoryVouchers, err := s.checkValidateListTrader(ctx, shopID, query.Result.InventoryVoucher)
	if err != nil {
		return nil, err
	}
	result := &api.GetInventoryVouchersResponse{
		InventoryVouchers: inventoryVouchers,
	}
	return result, nil
}

func (s *InventoryService) checkValidateListTrader(ctx context.Context, shopID dot.ID, inventoryVouchers []*inventory.InventoryVoucher) (result []*shop.InventoryVoucher, err error) {
	if inventoryVouchers == nil {
		return result, err
	}
	traderIDs := make([]dot.ID, 0, len(inventoryVouchers))
	count := 0
	for _, trader := range inventoryVouchers {
		if trader.TraderID != 0 {
			count++
		}
		traderIDs = append(traderIDs, trader.TraderID)
	}

	mapTraderValidate := map[dot.ID]bool{}
	if count > 0 {
		queryTraders := &tradering.ListTradersByIDsQuery{
			IDs:    traderIDs,
			ShopID: shopID,
		}
		if err := s.TraderQuery.Dispatch(ctx, queryTraders); err != nil {
			return result, err
		}
		traders := queryTraders.Result.Traders
		for _, trader := range traders {
			mapTraderValidate[trader.ID] = true
		}
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

func (s *InventoryService) GetInventoryVouchersByIDs(ctx context.Context, q *api.GetInventoryVouchersByIDsRequest) (*api.GetInventoryVouchersResponse, error) {
	shopID := s.SS.Shop().ID
	query := &inventory.GetInventoryVouchersByIDsQuery{
		ShopID: shopID,
		IDs:    q.Ids,
	}
	if err := s.InventoryQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &api.GetInventoryVouchersResponse{
		InventoryVouchers: PbShopInventoryVouchers(query.Result.InventoryVoucher),
	}
	return result, nil
}

func (s *InventoryService) UpdateInventoryVariantCostPrice(ctx context.Context, q *api.UpdateInventoryVariantCostPriceRequest) (*api.UpdateInventoryVariantCostPriceResponse, error) {
	cmd := &inventory.UpdateInventoryVariantCostPriceCommand{
		ShopID:    s.SS.Shop().ID,
		VariantID: q.VariantId,
		CostPrice: q.CostPrice,
	}
	err := s.InventoryAggr.Dispatch(ctx, cmd)
	if err != nil {
		return nil, err
	}
	result := &api.UpdateInventoryVariantCostPriceResponse{
		InventoryVariant: PbInventory(cmd.Result),
	}
	return result, nil
}
