package api

import (
	"context"

	"o.o/api/main/catalog"
	"o.o/api/main/inventory"
	"o.o/api/services/affiliate"
	apiaffiliate "o.o/api/top/services/affiliate"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/api/convertpb"
	pbshop "o.o/backend/pkg/etop/api/shop"
	modeletop "o.o/backend/pkg/etop/model"
	"o.o/capi/dot"
)

type ShopService struct {
	CatalogQuery   catalog.QueryBus
	InventoryQuery inventory.QueryBus
	AffiliateQuery affiliate.QueryBus
}

func (s *ShopService) Clone() *ShopService { res := *s; return &res }

func (s *ShopService) GetProductPromotion(ctx context.Context, q *GetProductPromotionEndpoint) error {
	promotionQuery := &affiliate.GetShopProductPromotionQuery{
		ShopID:    modeletop.EtopTradingAccountID,
		ProductID: q.ProductId,
	}
	if err := s.AffiliateQuery.Dispatch(ctx, promotionQuery); err != nil {
		return err
	}
	var pbReferralDiscount *apiaffiliate.CommissionSetting
	if q.ReferralCode.Valid {
		commissionSetting, err := GetCommissionSettingByReferralCode(ctx, s.AffiliateQuery, q.ReferralCode.String, q.ProductId)
		if err == nil {
			pbReferralDiscount = convertpb.PbCommissionSetting(commissionSetting)
		}
	}
	q.Result = &apiaffiliate.GetProductPromotionResponse{
		Promotion:        convertpb.PbProductPromotion(promotionQuery.Result),
		ReferralDiscount: pbReferralDiscount,
	}
	return nil
}

func (s *ShopService) ShopGetProducts(ctx context.Context, q *ShopGetProductsEndpoint) error {
	paging := cmapi.CMPaging(q.Paging)
	query := &catalog.ListShopProductsWithVariantsQuery{
		ShopID:  modeletop.EtopTradingAccountID,
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
	productPromotionMap := GetShopProductPromotionMapByProductIDs(ctx, s.AffiliateQuery, modeletop.EtopTradingAccountID, productIds)
	var products []*apiaffiliate.ShopProductResponse
	for _, product := range query.Result.Products {
		productPromotion := productPromotionMap[product.ProductID]
		var pbProductPromotion *apiaffiliate.ProductPromotion = nil
		if productPromotion != nil {
			pbProductPromotion = convertpb.PbProductPromotion(productPromotion)
		}
		productResult := pbshop.PbShopProductWithVariants(product)
		productResult, err := pbshop.PopulateTradingProductWithInventoryCount(ctx, s.InventoryQuery, productResult)
		if err != nil {
			return err
		}
		products = append(products, &apiaffiliate.ShopProductResponse{
			Product:   productResult,
			Promotion: pbProductPromotion,
		})
	}
	q.Result = &apiaffiliate.ShopGetProductsResponse{
		Paging:   cmapi.PbPageInfo(paging),
		Products: products,
	}
	return nil
}

func (s *ShopService) CheckReferralCodeValid(ctx context.Context, q *CheckReferralCodeValidEndpoint) error {
	affiliateAccountReferralQ := &affiliate.GetAffiliateAccountReferralByCodeQuery{
		Code: q.ReferralCode,
	}
	if err := s.AffiliateQuery.Dispatch(ctx, affiliateAccountReferralQ); err != nil {
		return cm.Errorf(cm.NotFound, nil, "Mã giới thiệu không hợp lệ")
	}

	if affiliateAccountReferralQ.Result.UserID == q.Context.Shop.OwnerID {
		return cm.Errorf(cm.ValidationFailed, nil, "Mã giới thiệu không hợp lệ")
	}

	promotionQuery := &affiliate.GetShopProductPromotionQuery{
		ShopID:    modeletop.EtopTradingAccountID,
		ProductID: q.ProductId,
	}
	_ = s.AffiliateQuery.Dispatch(ctx, promotionQuery)

	commissionSetting, err := GetCommissionSettingByReferralCode(ctx, s.AffiliateQuery, q.ReferralCode, q.ProductId)
	if err != nil {
		return cm.Errorf(cm.ValidationFailed, nil, "Không thể sử dụng mã giới thiệu của chính bạn")
	}
	pbReferralDiscount := convertpb.PbCommissionSetting(commissionSetting)
	q.Result = &apiaffiliate.GetProductPromotionResponse{
		Promotion:        convertpb.PbProductPromotion(promotionQuery.Result),
		ReferralDiscount: pbReferralDiscount,
	}
	return nil
}
