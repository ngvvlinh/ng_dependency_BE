package shop

import (
	"context"

	"etop.vn/api/main/catalog"
	"etop.vn/api/main/inventory"
	"etop.vn/api/main/stocktaking"
	"etop.vn/api/meta"
	pbcm "etop.vn/backend/pb/common"
	pbshop "etop.vn/backend/pb/etop/shop"
	cm "etop.vn/backend/pkg/common"
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

func (s *StocktakeService) AttachShopVariantsInformation(ctx context.Context, shopID int64, stocktakeLines []*stocktaking.StocktakeLine) error {
	var variantIDs []int64
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
	var mapVariants = make(map[int64]*catalog.ShopVariant)
	var mapProductIDs = make(map[int64]int64)
	for _, value := range queryVariants.Result.Variants {
		mapVariants[value.VariantID] = value
		mapProductIDs[value.ProductID] = value.ProductID
	}
	var productIDs []int64
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

	var mapProducts = make(map[int64]*catalog.ShopProduct)
	for _, value := range queryProducts.Result.Products {
		mapProducts[value.ProductID] = value
	}
	for key, value := range stocktakeLines {
		if mapVariants[value.VariantID] == nil {
			return cm.Errorf(cm.InvalidArgument, nil, "Phiên bản không tồn tại")
		}
		product := mapProducts[mapVariants[value.VariantID].ProductID]
		stocktakeLines[key] = ConvertInfoVariants(stocktakeLines[key], mapVariants[value.VariantID], product)
	}
	return nil
}

func ConvertInfoVariants(stocktakeLine *stocktaking.StocktakeLine, shopVariant *catalog.ShopVariant, shopProduct *catalog.ShopProduct) *stocktaking.StocktakeLine {
	stocktakeLine.VariantName = shopVariant.Name
	stocktakeLine.Code = shopVariant.Code
	stocktakeLine.ProductID = shopProduct.ProductID
	stocktakeLine.ProductName = shopProduct.Name
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
		OverStock:            cm.BoolDefault(overStock, true),
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
		ShopID: shopID,
		Id:     q.Id,
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
	paging := q.Paging.CMPaging()
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
		Paging:     pbcm.PbPaging(query.Page, query.Result.Total),
	}
	return nil
}
