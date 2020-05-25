package api

import (
	"context"

	"o.o/api/services/affiliate"
	"o.o/capi/dot"
)

func GetShopCommissionSettingsByProducts(ctx context.Context, affiliateQuery affiliate.QueryBus, accountID dot.ID, productIds []dot.ID) map[dot.ID]*affiliate.CommissionSetting {
	getShopCommissionByProductIDsQuery := &affiliate.GetCommissionByProductIDsQuery{
		AccountID:  accountID,
		ProductIDs: productIds,
	}
	if err := affiliateQuery.Dispatch(ctx, getShopCommissionByProductIDsQuery); err != nil {
		return map[dot.ID]*affiliate.CommissionSetting{}
	}

	shopCommissionMap := map[dot.ID]*affiliate.CommissionSetting{}
	for _, e := range getShopCommissionByProductIDsQuery.Result {
		shopCommissionMap[e.ProductID] = e
	}

	return shopCommissionMap
}

func GetSupplyCommissionSettingByProductIdsMap(ctx context.Context, affiliateQuery affiliate.QueryBus, shopID dot.ID, productIDs []dot.ID) map[dot.ID]*affiliate.SupplyCommissionSetting {
	supplyCommissionSettingsQ := &affiliate.GetSupplyCommissionSettingsByProductIDsQuery{
		ShopID:     shopID,
		ProductIDs: productIDs,
	}
	if err := affiliateQuery.Dispatch(ctx, supplyCommissionSettingsQ); err != nil {
		return map[dot.ID]*affiliate.SupplyCommissionSetting{}
	}
	supplyCommissionSettingMap := map[dot.ID]*affiliate.SupplyCommissionSetting{}
	for _, e := range supplyCommissionSettingsQ.Result {
		supplyCommissionSettingMap[e.ProductID] = e
	}
	return supplyCommissionSettingMap
}

func GetShopProductPromotionMapByProductIDs(ctx context.Context, affiliateQuery affiliate.QueryBus, shopID dot.ID, productIDs []dot.ID) map[dot.ID]*affiliate.ProductPromotion {
	promotionsQ := &affiliate.GetShopProductPromotionByProductIDsQuery{
		ShopID:     shopID,
		ProductIDs: productIDs,
	}
	if err := affiliateQuery.Dispatch(ctx, promotionsQ); err != nil {
		return map[dot.ID]*affiliate.ProductPromotion{}
	}

	promotionMap := map[dot.ID]*affiliate.ProductPromotion{}
	for _, e := range promotionsQ.Result {
		promotionMap[e.ProductID] = e
	}
	return promotionMap
}

func GetCommissionSettingByReferralCode(ctx context.Context, affiliateQuery affiliate.QueryBus, referralCode string, productID dot.ID) (*affiliate.CommissionSetting, error) {
	referralQ := &affiliate.GetAffiliateAccountReferralByCodeQuery{Code: referralCode}
	if err := affiliateQuery.Dispatch(ctx, referralQ); err != nil {
		return nil, err
	}
	commissionSettingQ := &affiliate.GetCommissionByProductIDQuery{
		AccountID: referralQ.Result.AffiliateID,
		ProductID: productID,
	}
	_ = affiliateQuery.Dispatch(ctx, commissionSettingQ)
	return commissionSettingQ.Result, nil
}
