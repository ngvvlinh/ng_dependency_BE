package api

import (
	"context"
	"strings"

	"o.o/api/main/catalog"
	"o.o/api/main/inventory"
	"o.o/api/services/affiliate"
	api "o.o/api/top/services/affiliate"
	pbcm "o.o/api/top/types/common"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etc/idutil"
	convertpball "o.o/backend/pkg/etop/api/convertpb/_all"
	"o.o/backend/pkg/etop/api/shop/product"
	"o.o/backend/pkg/etop/api/shop/trading"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/capi/dot"
)

type TradingService struct {
	session.Session

	AffiliateAggr  affiliate.CommandBus
	AffiliateQuery affiliate.QueryBus
	CatalogQuery   catalog.QueryBus
	InventoryQuery inventory.QueryBus
}

func (s *TradingService) Clone() api.TradingService { res := *s; return &res }

func (s *TradingService) TradingGetProducts(ctx context.Context, q *pbcm.CommonListRequest) (*api.SupplyGetProductsResponse, error) {
	if s.SS.Shop().ID != idutil.EtopTradingAccountID {
		return nil, cm.Errorf(cm.PermissionDenied, nil, "PermissionDenied")
	}
	paging := cmapi.CMPaging(q.Paging)
	query := &catalog.ListShopProductsWithVariantsQuery{
		ShopID:  idutil.EtopTradingAccountID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(q.Filters),
	}
	if err := s.CatalogQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}

	var productIds []dot.ID
	for _, product := range query.Result.Products {
		productIds = append(productIds, product.ProductID)
	}

	supplyCommissionSettingMap := GetSupplyCommissionSettingByProductIdsMap(ctx, s.AffiliateQuery, idutil.EtopTradingAccountID, productIds)
	productPromotionMap := GetShopProductPromotionMapByProductIDs(ctx, s.AffiliateQuery, idutil.EtopTradingAccountID, productIds)
	var products []*api.SupplyProductResponse
	for _, p := range query.Result.Products {
		supplyCommissionSetting := supplyCommissionSettingMap[p.ProductID]
		var pbSupplyCommissionSetting *api.SupplyCommissionSetting = nil
		if supplyCommissionSetting != nil {
			pbSupplyCommissionSetting = convertpball.PbSupplyCommissionSetting(supplyCommissionSetting)
		}
		productPromotion := productPromotionMap[p.ProductID]
		var pbProductPromotion *api.ProductPromotion = nil
		if productPromotion != nil {
			pbProductPromotion = convertpball.PbProductPromotion(productPromotion)
		}
		productResult := product.PbShopProductWithVariants(p)
		productResult, err := trading.PopulateTradingProductWithInventoryCount(ctx, s.InventoryQuery, productResult)
		if err != nil {
			return nil, err
		}
		products = append(products, &api.SupplyProductResponse{
			Product:                 productResult,
			SupplyCommissionSetting: pbSupplyCommissionSetting,
			Promotion:               pbProductPromotion,
		})
	}

	result := &api.SupplyGetProductsResponse{
		Paging:   cmapi.PbPageInfo(paging),
		Products: products,
	}
	return result, nil
}

func (s *TradingService) GetTradingProductPromotions(ctx context.Context, q *pbcm.CommonListRequest) (*api.GetProductPromotionsResponse, error) {
	if s.SS.Shop().ID != idutil.EtopTradingAccountID {
		return nil, cm.Errorf(cm.PermissionDenied, nil, "PermissionDenied")
	}
	paging := cmapi.CMPaging(q.Paging)
	query := &affiliate.ListShopProductPromotionsQuery{
		ShopID:  idutil.EtopTradingAccountID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(q.Filters),
	}

	if err := s.AffiliateQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}

	result := &api.GetProductPromotionsResponse{
		Paging:     cmapi.PbPageInfo(paging),
		Promotions: convertpball.PbProductPromotions(query.Result.Promotions),
	}
	return result, nil
}

func (s *TradingService) CreateOrUpdateTradingCommissionSetting(ctx context.Context, q *api.CreateOrUpdateTradingCommissionSettingRequest) (*api.SupplyCommissionSetting, error) {
	if s.SS.Shop().ID != idutil.EtopTradingAccountID {
		return nil, cm.Errorf(cm.PermissionDenied, nil, "PermissionDenied")
	}

	q.DependOn = strings.ToLower(q.DependOn)
	q.Group = strings.ToLower(q.Group)

	cmd := &affiliate.CreateOrUpdateSupplyCommissionSettingCommand{
		ShopID:                   s.SS.Shop().ID,
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
		return nil, err
	}

	result := convertpball.PbSupplyCommissionSetting(cmd.Result)

	return result, nil
}

func (s *TradingService) GetTradingProductPromotionByProductIDs(ctx context.Context, q *api.GetTradingProductPromotionByIDsRequest) (*api.GetTradingProductPromotionByIDsResponse, error) {
	if s.SS.Shop().ID != idutil.EtopTradingAccountID {
		return nil, cm.Errorf(cm.PermissionDenied, nil, "PermissionDenied")
	}
	productPromotionsQ := &affiliate.GetShopProductPromotionByProductIDsQuery{
		ShopID:     idutil.EtopTradingAccountID,
		ProductIDs: q.ProductIds,
	}
	if err := s.AffiliateQuery.Dispatch(ctx, productPromotionsQ); err != nil {
		return nil, err
	}
	result := &api.GetTradingProductPromotionByIDsResponse{
		Promotions: convertpball.PbProductPromotions(productPromotionsQ.Result),
	}
	return result, nil
}

func (s *TradingService) CreateTradingProductPromotion(ctx context.Context, q *api.CreateOrUpdateProductPromotionRequest) (*api.ProductPromotion, error) {
	if s.SS.Shop().ID != idutil.EtopTradingAccountID {
		return nil, cm.Errorf(cm.PermissionDenied, nil, "PermissionDenied")
	}

	if err := s.AffiliateQuery.Dispatch(ctx, &affiliate.GetShopProductPromotionQuery{
		ShopID:    s.SS.Shop().ID,
		ProductID: q.ProductId,
	}); err == nil {
		return nil, cm.Errorf(cm.AlreadyExists, nil, "Sản phẩm đã có chương trình khuyến mãi")
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
		return nil, err
	}
	result := convertpball.PbProductPromotion(cmd.Result)
	return result, nil
}

func (s *TradingService) UpdateTradingProductPromotion(ctx context.Context, q *api.CreateOrUpdateProductPromotionRequest) (*api.ProductPromotion, error) {
	if s.SS.Shop().ID != idutil.EtopTradingAccountID {
		return nil, cm.Errorf(cm.PermissionDenied, nil, "PermissionDenied")
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
		return nil, err
	}
	result := convertpball.PbProductPromotion(cmd.Result)
	return result, nil
}
