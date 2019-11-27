package shop

import (
	"context"

	"etop.vn/api/main/catalog"
	"etop.vn/api/main/inventory"
	"etop.vn/api/main/stocktaking"
	"etop.vn/api/meta"
	pbshop "etop.vn/api/pb/etop/shop"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmapi"
	"etop.vn/capi/dot"
)

func (s *StocktakeService) CreateStocktake(ctx context.Context, q *CreateStocktakeEndpoint) error {
	shopID := q.Context.Shop.ID
	UserID := q.Context.UserID

	var lines []*stocktaking.StocktakeLine
	for _, value := range q.Lines {
		lines = append(lines, &stocktaking.StocktakeLine{
			VariantID:   value.VariantId,
			OldQuantity: value.OldQuantity,
			NewQuantity: value.NewQuantity,
		})
	}
	err := s.AttachShopVariantsInformation(ctx, shopID, lines)
	if err != nil {
		return err
	}
	cmd := &stocktaking.CreateStocktakeCommand{
		ShopID:        shopID,
		TotalQuantity: q.TotalQuantity,
		CreatedBy:     UserID,
		Lines:         lines,
		Note:          q.Note,
	}
	err = StocktakeAggregate.Dispatch(ctx, cmd)
	if err != nil {
		return err
	}
	q.Result = PbStocktake(cmd.Result)
	return nil
}

func (s *StocktakeService) AttachShopVariantsInformation(ctx context.Context, shopID dot.ID, stocktakeLines []*stocktaking.StocktakeLine) error {
	var variantIDs []dot.ID
	for _, value := range stocktakeLines {
		variantIDs = append(variantIDs, value.VariantID)
	}
	queryVariants := &catalog.ListShopVariantsByIDsQuery{
		IDs:    variantIDs,
		ShopID: shopID,
		Result: nil,
	}
	err := catalogQuery.Dispatch(ctx, queryVariants)
	if err != nil {
		return err
	}
	var mapVariants = make(map[dot.ID]*catalog.ShopVariant)
	var mapProductIDs = make(map[dot.ID]dot.ID)
	for _, value := range queryVariants.Result.Variants {
		mapVariants[value.VariantID] = value
		mapProductIDs[value.ProductID] = value.ProductID
	}
	var productIDs []dot.ID
	for key, _ := range mapProductIDs {
		productIDs = append(productIDs, key)
	}
	queryProducts := &catalog.ListShopProductsByIDsQuery{
		IDs:    productIDs,
		ShopID: shopID,
	}
	err = catalogQuery.Dispatch(ctx, queryProducts)
	if err != nil {
		return err
	}

	var mapProducts = make(map[dot.ID]*catalog.ShopProduct)
	for _, value := range queryProducts.Result.Products {
		mapProducts[value.ProductID] = value
	}

	queryInventoryVariant := &inventory.GetInventoryVariantsByVariantIDsQuery{
		ShopID:     shopID,
		VariantIDs: variantIDs,
	}
	err = inventoryQuery.Dispatch(ctx, queryInventoryVariant)
	if err != nil {
		return err
	}

	var mapInventoryVariant = make(map[dot.ID]*inventory.InventoryVariant)
	for _, value := range queryInventoryVariant.Result.InventoryVariants {
		mapInventoryVariant[value.VariantID] = value
	}

	for key, value := range stocktakeLines {
		if mapVariants[value.VariantID] == nil {
			return cm.Errorf(cm.InvalidArgument, nil, "Phiên bản không tồn tại")
		}
		product := mapProducts[mapVariants[value.VariantID].ProductID]
		inventoryVariant := mapInventoryVariant[value.VariantID]
		stocktakeLines[key] = ConvertInfoVariants(stocktakeLines[key], mapVariants[value.VariantID], product, inventoryVariant)
	}
	return nil
}

func ConvertInfoVariants(stocktakeLine *stocktaking.StocktakeLine, shopVariant *catalog.ShopVariant, shopProduct *catalog.ShopProduct, inventoryVariant *inventory.InventoryVariant) *stocktaking.StocktakeLine {
	stocktakeLine.VariantName = shopVariant.Name
	stocktakeLine.Code = shopVariant.Code
	stocktakeLine.ProductID = shopProduct.ProductID
	stocktakeLine.ProductName = shopProduct.Name
	if inventoryVariant != nil {
		stocktakeLine.CostPrice = inventoryVariant.CostPrice
	}
	var attributes []*stocktaking.Attribute
	for _, value := range shopVariant.Attributes {
		attributes = append(attributes, &stocktaking.Attribute{
			Name:  value.Name,
			Value: value.Value,
		})
	}
	if len(shopVariant.ImageURLs) > 0 {
		stocktakeLine.ImageURL = shopVariant.ImageURLs[0]
	}
	if stocktakeLine.ImageURL == "" {
		if len(shopProduct.ImageURLs) > 0 {
			stocktakeLine.ImageURL = shopProduct.ImageURLs[0]
		}
	}
	stocktakeLine.Attributes = attributes
	return stocktakeLine
}

func (s *StocktakeService) UpdateStocktake(
	ctx context.Context, q *UpdateStocktakeEndpoint) error {
	shopID := q.Context.Shop.ID
	UserID := q.Context.UserID
	var lines []*stocktaking.StocktakeLine
	for _, value := range q.Lines {
		lines = append(lines, &stocktaking.StocktakeLine{
			VariantID:   value.VariantId,
			OldQuantity: value.OldQuantity,
			NewQuantity: value.NewQuantity,
		})
	}
	err := s.AttachShopVariantsInformation(ctx, shopID, lines)
	if err != nil {
		return err
	}
	cmd := &stocktaking.UpdateStocktakeCommand{
		ShopID:        shopID,
		TotalQuantity: q.TotalQuantity,
		ID:            q.Id,
		UpdatedBy:     UserID,
		Lines:         lines,
		Note:          q.Note,
	}
	err = StocktakeAggregate.Dispatch(ctx, cmd)
	if err != nil {
		return err
	}
	q.Result = PbStocktake(cmd.Result)
	return nil
}

func (s *StocktakeService) ConfirmStocktake(ctx context.Context, q *ConfirmStocktakeEndpoint) error {
	shopID := q.Context.Shop.ID
	userID := q.Context.UserID
	overStock := q.Context.Shop.InventoryOverstock
	cmd := &stocktaking.ConfirmStocktakeCommand{
		ID:                   q.Id,
		ShopID:               shopID,
		ConfirmedBy:          userID,
		OverStock:            overStock.Apply(true),
		AutoInventoryVoucher: inventory.AutoInventoryVoucher(q.AutoInventoryVoucher),
		Result:               nil,
	}
	err := StocktakeAggregate.Dispatch(ctx, cmd)
	if err != nil {
		return err
	}
	q.Result = PbStocktake(cmd.Result)
	return nil
}

func (s *StocktakeService) CancelStocktake(ctx context.Context, q *CancelStocktakeEndpoint) error {
	shopID := q.Context.Shop.ID
	cmd := &stocktaking.CancelStocktakeCommand{
		ShopID:       shopID,
		ID:           q.Id,
		CancelReason: q.CancelReason,
	}
	err := StocktakeAggregate.Dispatch(ctx, cmd)
	if err != nil {
		return err
	}
	q.Result = PbStocktake(cmd.Result)
	return nil
}

func (s *StocktakeService) GetStocktake(ctx context.Context, q *GetStocktakeEndpoint) error {
	shopID := q.Context.Shop.ID
	query := &stocktaking.GetStocktakeByIDQuery{
		ShopID: shopID,
		Id:     q.Id,
	}
	err := StocktakeQuery.Dispatch(ctx, query)
	if err != nil {
		return err
	}
	q.Result = PbStocktake(query.Result)
	return nil
}

func (s *StocktakeService) GetStocktakesByIDs(ctx context.Context, q *GetStocktakesByIDsEndpoint) error {
	shopID := q.Context.Shop.ID
	query := &stocktaking.GetStocktakesByIDsQuery{
		ShopID: shopID,
		Ids:    q.Ids,
	}
	err := StocktakeQuery.Dispatch(ctx, query)
	if err != nil {
		return err
	}
	q.Result = &pbshop.GetStocktakesByIDsResponse{
		Stocktakes: PbStocktakes(query.Result),
	}
	return nil
}

func (s *StocktakeService) GetStocktakes(ctx context.Context, q *GetStocktakesEndpoint) error {
	shopID := q.Context.Shop.ID
	var filters []meta.Filter
	for _, value := range q.Filters {
		filters = append(filters, meta.Filter{
			Name:  value.Name,
			Op:    value.Op,
			Value: value.Value,
		})
	}
	paging := cmapi.CMPaging(q.Paging)
	query := &stocktaking.ListStocktakeQuery{
		Page:   *paging,
		ShopID: shopID,
		Filter: filters,
	}
	err := StocktakeQuery.Dispatch(ctx, query)
	if err != nil {
		return err
	}
	q.Result = &pbshop.GetStocktakesResponse{
		Stocktakes: PbStocktakes(query.Result.Stocktakes),
		Paging:     cmapi.PbPaging(query.Page, query.Result.Total),
	}
	return nil
}
