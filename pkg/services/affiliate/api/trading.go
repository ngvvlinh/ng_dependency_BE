package api

import (
	"context"
	"strings"

	"o.o/api/main/catalog"
	"o.o/api/main/inventory"
	"o.o/api/services/affiliate"
	apiaffiliate "o.o/api/top/services/affiliate"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etc/idutil"
	"o.o/backend/pkg/etop/api/convertpb"
	product2 "o.o/backend/pkg/etop/api/shop/product"
	"o.o/backend/pkg/etop/api/shop/trading"
	"o.o/capi/dot"
)

type TradingService struct {
	AffiliateAggr  affiliate.CommandBus
	AffiliateQuery affiliate.QueryBus
	CatalogQuery   catalog.QueryBus
	InventoryQuery inventory.QueryBus
}

func (s *TradingService) Clone() *TradingService { res := *s; return &res }

func (s *TradingService) TradingGetProducts(ctx context.Context, q *TradingGetProductsEndpoint) error {
	if q.Context.Shop.ID != idutil.EtopTradingAccountID {
		return cm.Errorf(cm.PermissionDenied, nil, "PermissionDenied")
	}
	paging := cmapi.CMPaging(q.Paging)
	query := &catalog.ListShopProductsWithVariantsQuery{
		ShopID:  idutil.EtopTradingAccountID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(q.Filters),
	}
	if err := s.CatalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}

	var productIds []dot.ID
	for _, product := range query.Result.Products {
		productIds = append(productIds, product.ProductID)
	}

	supplyCommissionSettingMap := GetSupplyCommissionSettingByProductIdsMap(ctx, s.AffiliateQuery, idutil.EtopTradingAccountID, productIds)
	productPromotionMap := GetShopProductPromotionMapByProductIDs(ctx, s.AffiliateQuery, idutil.EtopTradingAccountID, productIds)
	var products []*apiaffiliate.SupplyProductResponse
	for _, product := range query.Result.Products {
		supplyCommissionSetting := supplyCommissionSettingMap[product.ProductID]
		var pbSupplyCommissionSetting *apiaffiliate.SupplyCommissionSetting = nil
		if supplyCommissionSetting != nil {
			pbSupplyCommissionSetting = convertpb.PbSupplyCommissionSetting(supplyCommissionSetting)
		}
		productPromotion := productPromotionMap[product.ProductID]
		var pbProductPromotion *apiaffiliate.ProductPromotion = nil
		if productPromotion != nil {
			pbProductPromotion = convertpb.PbProductPromotion(productPromotion)
		}
		productResult := product2.PbShopProductWithVariants(product)
		productResult, err := trading.PopulateTradingProductWithInventoryCount(ctx, s.InventoryQuery, productResult)
		if err != nil {
			return err
		}
		products = append(products, &apiaffiliate.SupplyProductResponse{
			Product:                 productResult,
			SupplyCommissionSetting: pbSupplyCommissionSetting,
			Promotion:               pbProductPromotion,
		})
	}

	q.Result = &apiaffiliate.SupplyGetProductsResponse{
		Paging:   cmapi.PbPageInfo(paging),
		Products: products,
	}
	return nil
}

func (s *TradingService) GetTradingProductPromotions(ctx context.Context, q *GetTradingProductPromotionsEndpoint) error {
	if q.Context.Shop.ID != idutil.EtopTradingAccountID {
		return cm.Errorf(cm.PermissionDenied, nil, "PermissionDenied")
	}
	paging := cmapi.CMPaging(q.Paging)
	query := &affiliate.ListShopProductPromotionsQuery{
		ShopID:  idutil.EtopTradingAccountID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(q.Filters),
	}

	if err := s.AffiliateQuery.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = &apiaffiliate.GetProductPromotionsResponse{
		Paging:     cmapi.PbPageInfo(paging),
		Promotions: convertpb.PbProductPromotions(query.Result.Promotions),
	}
	return nil
}

func (s *TradingService) CreateOrUpdateTradingCommissionSetting(ctx context.Context, q *CreateOrUpdateTradingCommissionSettingEndpoint) error {
	if q.Context.Shop.ID != idutil.EtopTradingAccountID {
		return cm.Errorf(cm.PermissionDenied, nil, "PermissionDenied")
	}

	q.DependOn = strings.ToLower(q.DependOn)
	q.Group = strings.ToLower(q.Group)

	cmd := &affiliate.CreateOrUpdateSupplyCommissionSettingCommand{
		ShopID:                   q.Context.Shop.ID,
		ProductID:                q.ProductId,
		Level1DirectCommission:   q.Level1DirectCommission,
		Level1IndirectCommission: q.Level1IndirectCommission,
		Level2DirectCommission:   q.Level2DirectCommission,
		Level2IndirectCommission: q.Level2IndirectCommission,
		DependOn:                 q.DependOn,
		Level1LimitCount:         q.Level1LimitCount,
		Level1LimitDuration:      q.Level1LimitDuration,
		Level1LimitDurationType:  q.Level1LimitDurationType,
		LifetimeDuration:         q.LifetimeDuration,
		LifetimeDurationType:     q.LifetimeDurationType,
		Group:                    q.Group,
	}

	if err := s.AffiliateAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}

	q.Result = convertpb.PbSupplyCommissionSetting(cmd.Result)

	return nil
}

func (s *TradingService) GetTradingProductPromotionByProductIDs(ctx context.Context, q *GetTradingProductPromotionByProductIDsEndpoint) error {
	if q.Context.Shop.ID != idutil.EtopTradingAccountID {
		return cm.Errorf(cm.PermissionDenied, nil, "PermissionDenied")
	}
	productPromotionsQ := &affiliate.GetShopProductPromotionByProductIDsQuery{
		ShopID:     idutil.EtopTradingAccountID,
		ProductIDs: q.ProductIds,
	}
	if err := s.AffiliateQuery.Dispatch(ctx, productPromotionsQ); err != nil {
		return err
	}
	q.Result = &apiaffiliate.GetTradingProductPromotionByIDsResponse{
		Promotions: convertpb.PbProductPromotions(productPromotionsQ.Result),
	}
	return nil
}

func (s *TradingService) CreateTradingProductPromotion(ctx context.Context, q *CreateTradingProductPromotionEndpoint) error {
	if q.Context.Shop.ID != idutil.EtopTradingAccountID {
		return cm.Errorf(cm.PermissionDenied, nil, "PermissionDenied")
	}

	if err := s.AffiliateQuery.Dispatch(ctx, &affiliate.GetShopProductPromotionQuery{
		ShopID:    q.Context.Shop.ID,
		ProductID: q.ProductId,
	}); err == nil {
		return cm.Errorf(cm.AlreadyExists, nil, "Sản phẩm đã có chương trình khuyến mãi")
	}

	cmd := &affiliate.CreateProductPromotionCommand{
		ShopID:      idutil.EtopTradingAccountID,
		ProductID:   q.ProductId,
		Amount:      q.Amount,
		Code:        q.Code,
		Description: q.Description,
		Unit:        q.Unit,
		Note:        q.Note,
		Type:        q.Type,
	}
	if err := s.AffiliateAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = convertpb.PbProductPromotion(cmd.Result)
	return nil
}

func (s *TradingService) UpdateTradingProductPromotion(ctx context.Context, q *UpdateTradingProductPromotionEndpoint) error {
	if q.Context.Shop.ID != idutil.EtopTradingAccountID {
		return cm.Errorf(cm.PermissionDenied, nil, "PermissionDenied")
	}
	cmd := &affiliate.UpdateProductPromotionCommand{
		ID:          q.Id,
		Amount:      q.Amount,
		Unit:        q.Unit,
		Code:        q.Code,
		Description: q.Description,
		Note:        q.Note,
		Type:        q.Type,
	}
	if err := s.AffiliateAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = convertpb.PbProductPromotion(cmd.Result)
	return nil
}
