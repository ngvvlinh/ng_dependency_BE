package api

import (
	"context"

	"o.o/api/main/catalog"
	"o.o/api/main/inventory"
	"o.o/api/services/affiliate"
	api "o.o/api/top/services/affiliate"
	pbcm "o.o/api/top/types/common"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etc/idutil"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/api/shop/product"
	"o.o/backend/pkg/etop/api/shop/trading"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/capi/dot"
)

type ShopService struct {
	session.Session

	CatalogQuery   catalog.QueryBus
	InventoryQuery inventory.QueryBus
	AffiliateQuery affiliate.QueryBus
}

func (s *ShopService) Clone() api.ShopService { res := *s; return &res }

func (s *ShopService) GetProductPromotion(ctx context.Context, q *api.GetProductPromotionRequest) (*api.GetProductPromotionResponse, error) {
	promotionQuery := &affiliate.GetShopProductPromotionQuery{
		ShopID:    idutil.EtopTradingAccountID,
		ProductID: q.ProductId,
	}
	if err := s.AffiliateQuery.Dispatch(ctx, promotionQuery); err != nil {
		return nil, err
	}
	var pbReferralDiscount *api.CommissionSetting
	if q.ReferralCode.Valid {
		commissionSetting, err := GetCommissionSettingByReferralCode(ctx, s.AffiliateQuery, q.ReferralCode.String, q.ProductId)
		if err == nil {
			pbReferralDiscount = convertpb.PbCommissionSetting(commissionSetting)
		}
	}
	result := &api.GetProductPromotionResponse{
		Promotion:        convertpb.PbProductPromotion(promotionQuery.Result),
		ReferralDiscount: pbReferralDiscount,
	}
	return result, nil
}

func (s *ShopService) ShopGetProducts(ctx context.Context, q *pbcm.CommonListRequest) (*api.ShopGetProductsResponse, error) {
	paging := cmapi.CMPaging(q.Paging)
	query := &catalog.ListShopProductsWithVariantsQuery{
		ShopID:  idutil.EtopTradingAccountID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(q.Filters),
		Name:    q.Filter.Name,
	}
	if err := s.CatalogQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}

	var productIds []dot.ID
	for _, product := range query.Result.Products {
		productIds = append(productIds, product.ProductID)
	}
	productPromotionMap := GetShopProductPromotionMapByProductIDs(ctx, s.AffiliateQuery, idutil.EtopTradingAccountID, productIds)
	var products []*api.ShopProductResponse
	for _, p := range query.Result.Products {
		productPromotion := productPromotionMap[p.ProductID]
		var pbProductPromotion *api.ProductPromotion = nil
		if productPromotion != nil {
			pbProductPromotion = convertpb.PbProductPromotion(productPromotion)
		}
		productResult := product.PbShopProductWithVariants(p)
		productResult, err := trading.PopulateTradingProductWithInventoryCount(ctx, s.InventoryQuery, productResult)
		if err != nil {
			return nil, err
		}
		products = append(products, &api.ShopProductResponse{
			Product:   productResult,
			Promotion: pbProductPromotion,
		})
	}
	result := &api.ShopGetProductsResponse{
		Paging:   cmapi.PbPageInfo(paging),
		Products: products,
	}
	return result, nil
}

func (s *ShopService) CheckReferralCodeValid(ctx context.Context, q *api.CheckReferralCodeValidRequest) (*api.GetProductPromotionResponse, error) {
	affiliateAccountReferralQ := &affiliate.GetAffiliateAccountReferralByCodeQuery{
		Code: q.ReferralCode,
	}
	if err := s.AffiliateQuery.Dispatch(ctx, affiliateAccountReferralQ); err != nil {
		return nil, cm.Errorf(cm.NotFound, nil, "Mã giới thiệu không hợp lệ")
	}

	if affiliateAccountReferralQ.Result.UserID == s.SS.Shop().OwnerID {
		return nil, cm.Errorf(cm.ValidationFailed, nil, "Mã giới thiệu không hợp lệ")
	}

	promotionQuery := &affiliate.GetShopProductPromotionQuery{
		ShopID:    idutil.EtopTradingAccountID,
		ProductID: q.ProductId,
	}
	_ = s.AffiliateQuery.Dispatch(ctx, promotionQuery)

	commissionSetting, err := GetCommissionSettingByReferralCode(ctx, s.AffiliateQuery, q.ReferralCode, q.ProductId)
	if err != nil {
		return nil, cm.Errorf(cm.ValidationFailed, nil, "Không thể sử dụng mã giới thiệu của chính bạn")
	}
	pbReferralDiscount := convertpb.PbCommissionSetting(commissionSetting)
	result := &api.GetProductPromotionResponse{
		Promotion:        convertpb.PbProductPromotion(promotionQuery.Result),
		ReferralDiscount: pbReferralDiscount,
	}
	return result, nil
}
