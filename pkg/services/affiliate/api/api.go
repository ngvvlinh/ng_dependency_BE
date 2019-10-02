package api

import (
	"context"

	"etop.vn/api/main/identity"

	cm "etop.vn/backend/pkg/common"

	"etop.vn/api/main/catalog"
	"etop.vn/api/services/affiliate"
	pbcm "etop.vn/backend/pb/common"
	pbaff "etop.vn/backend/pb/services/affiliate"
	pbshop "etop.vn/backend/pkg/etop/api/shop"
	modeletop "etop.vn/backend/pkg/etop/model"
	wrapaff "etop.vn/backend/wrapper/services/affiliate"
	"etop.vn/common/bus"
	"etop.vn/common/l"
)

func init() {
	bus.AddHandlers("",
		UpdateReferral,

		TradingGetProducts,
		CreateOrUpdateTradingCommissionSetting,
		GetTradingProductPromotions,
		GetTradingProductPromotionByProductIDs,
		CreateTradingProductPromotion,
		UpdateTradingProductPromotion,

		GetProductPromotion,

		GetCommissions,
		NotifyNewShopPurchase,
		GetTransactions,
		CreateOrUpdateAffiliateCommissionSetting,
		GetProductPromotionByProductID,
		AffiliateGetProducts,
		CreateReferralCode,
		GetReferralCodes,
		GetReferrals,
	)
}

var ll = l.New()

var (
	catalogQuery   catalog.QueryBus
	affiliateCmd   affiliate.CommandBus
	affiliateQuery affiliate.QueryBus
	identityQuery  identity.QueryBus
)

func Init(
	affCmd affiliate.CommandBus,
	affQuery affiliate.QueryBus,
	catQuery catalog.QueryBus,
	idenQuery identity.QueryBus,
) {
	affiliateCmd = affCmd
	catalogQuery = catQuery
	affiliateQuery = affQuery
	identityQuery = idenQuery
}

func UpdateReferral(ctx context.Context, q *wrapaff.UpdateReferralEndpoint) error {
	cmd := &affiliate.CreateOrUpdateUserReferralCommand{
		UserID:           q.Context.UserID,
		ReferralCode:     q.ReferralCode,
		SaleReferralCode: q.SaleReferralCode,
	}
	if err := affiliateCmd.Dispatch(ctx, cmd); err != nil {
		return err
	}

	q.Result = &pbaff.UserReferral{
		UserId:           cmd.Result.UserID,
		ReferralCode:     cmd.Result.ReferralCode,
		SaleReferralCode: cmd.Result.SaleReferralCode,
	}
	return nil
}

func TradingGetProducts(ctx context.Context, q *wrapaff.TradingGetProductsEndpoint) error {
	if q.Context.Shop.ID != modeletop.EtopTradingAccountID {
		return cm.Errorf(cm.PermissionDenied, nil, "PermissionDenied")
	}
	paging := q.Paging.CMPaging()
	query := &catalog.ListShopProductsWithVariantsQuery{
		ShopID:  modeletop.EtopTradingAccountID,
		Paging:  *paging,
		Filters: pbcm.ToFilters(q.Filters),
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}

	var productIds []int64
	for _, product := range query.Result.Products {
		productIds = append(productIds, product.ProductID)
	}

	shopCommissionSettingMap := GetShopCommissionSettingsByProducts(ctx, modeletop.EtopTradingAccountID, productIds)
	var products []*pbaff.ShopProductResponse
	for _, product := range query.Result.Products {
		shopCommissionSetting := shopCommissionSettingMap[product.ProductID]
		var pbShopCommissionSetting *pbaff.CommissionSetting = nil
		if shopCommissionSetting != nil {
			pbShopCommissionSetting = pbaff.PbCommissionSetting(shopCommissionSetting)
		}
		products = append(products, &pbaff.ShopProductResponse{
			Product:               pbshop.PbShopProductWithVariants(product),
			ShopCommissionSetting: pbShopCommissionSetting,
		})
	}

	q.Result = &pbaff.ShopGetProductsResponse{
		Paging:   pbcm.PbPageInfo(paging, query.Result.Count),
		Products: products,
	}
	return nil
}

func CreateOrUpdateTradingCommissionSetting(ctx context.Context, q *wrapaff.CreateOrUpdateTradingCommissionSettingEndpoint) error {
	if q.Context.Shop.ID != modeletop.EtopTradingAccountID {
		return cm.Errorf(cm.PermissionDenied, nil, "PermissionDenied")
	}

	cmd := &affiliate.CreateOrUpdateCommissionSettingCommand{
		ProductID: q.ProductId,
		AccountID: modeletop.EtopTradingAccountID, // TODO test public api
		Amount:    q.Amount,
		Unit:      *q.Unit,
		Type:      "shop",
	}
	if err := affiliateCmd.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = pbaff.PbCommissionSetting(cmd.Result)
	return nil
}

func GetTradingProductPromotions(ctx context.Context, q *wrapaff.GetTradingProductPromotionsEndpoint) error {
	if q.Context.Shop.ID != modeletop.EtopTradingAccountID {
		return cm.Errorf(cm.PermissionDenied, nil, "PermissionDenied")
	}
	paging := q.Paging.CMPaging()
	query := &affiliate.ListShopProductPromotionsQuery{
		ShopID:  modeletop.EtopTradingAccountID,
		Paging:  *paging,
		Filters: pbcm.ToFilters(q.Filters),
	}

	if err := affiliateQuery.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = &pbaff.GetProductPromotionsResponse{
		Paging:     pbcm.PbPageInfo(paging, query.Result.Count),
		Promotions: pbaff.PbProductPromotions(query.Result.Promotions),
	}
	return nil
}

func GetTradingProductPromotionByProductIDs(ctx context.Context, q *wrapaff.GetTradingProductPromotionByProductIDsEndpoint) error {
	if q.Context.Shop.ID != modeletop.EtopTradingAccountID {
		return cm.Errorf(cm.PermissionDenied, nil, "PermissionDenied")
	}
	productPromotionsQ := &affiliate.GetShopProductPromotionByProductIDsQuery{
		ShopID:     modeletop.EtopTradingAccountID,
		ProductIDs: q.ProductIds,
	}
	if err := affiliateQuery.Dispatch(ctx, productPromotionsQ); err != nil {
		return err
	}
	q.Result = &pbaff.GetTradingProductPromotionByIDsResponse{
		Promotions: pbaff.PbProductPromotions(productPromotionsQ.Result),
	}
	return nil
}

func CreateTradingProductPromotion(ctx context.Context, q *wrapaff.CreateTradingProductPromotionEndpoint) error {
	if q.Context.Shop.ID != modeletop.EtopTradingAccountID {
		return cm.Errorf(cm.PermissionDenied, nil, "PermissionDenied")
	}
	cmd := &affiliate.CreateProductPromotionCommand{
		ShopID:      modeletop.EtopTradingAccountID,
		ProductID:   q.ProductId,
		Amount:      q.Amount,
		Code:        q.Code,
		Description: q.Description,
		Unit:        q.Unit,
		Note:        q.Note,
		Type:        q.Type,
	}
	if err := affiliateCmd.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = pbaff.PbProductPromotion(cmd.Result)
	return nil
}

func UpdateTradingProductPromotion(ctx context.Context, q *wrapaff.UpdateTradingProductPromotionEndpoint) error {
	if q.Context.Shop.ID != modeletop.EtopTradingAccountID {
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
	if err := affiliateCmd.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = pbaff.PbProductPromotion(cmd.Result)
	return nil
}

func GetProductPromotion(ctx context.Context, q *wrapaff.GetProductPromotionEndpoint) error {
	promotionQuery := &affiliate.GetShopProductPromotionQuery{
		ShopID:    modeletop.EtopTradingAccountID,
		ProductID: q.ProductId,
	}
	if err := affiliateQuery.Dispatch(ctx, promotionQuery); err != nil {
		return err
	}
	var pbReferralDiscount *pbaff.CommissionSetting
	if q.ReferralCode != nil {
		commissionSetting, err := GetCommissionSettingByReferralCode(ctx, *q.ReferralCode, q.ProductId)
		if err == nil {
			pbReferralDiscount = pbaff.PbCommissionSetting(commissionSetting)
		}
	}
	q.Result = &pbaff.GetProductPromotionResponse{
		Promotion:        pbaff.PbProductPromotion(promotionQuery.Result),
		ReferralDiscount: pbReferralDiscount,
	}
	return nil
}

func GetCommissions(ctx context.Context, q *wrapaff.GetCommissionsEndpoint) error {
	q.Result = &pbaff.GetCommissionsResponse{Message: "hello"}
	return nil
}

func NotifyNewShopPurchase(ctx context.Context, q *wrapaff.NotifyNewShopPurchaseEndpoint) error {
	panic("IMPLEMENT ME")
}

func GetTransactions(ctx context.Context, q *wrapaff.GetTransactionsEndpoint) error {
	panic("IMPLEMENT ME")
}

func CreateOrUpdateAffiliateCommissionSetting(ctx context.Context, q *wrapaff.CreateOrUpdateAffiliateCommissionSettingEndpoint) error {
	cmd := &affiliate.CreateOrUpdateCommissionSettingCommand{
		ProductID: q.ProductId,
		AccountID: q.Context.Affiliate.ID,
		Amount:    q.Amount,
		Unit:      *q.Unit,
		Type:      "affiliate",
	}
	if err := affiliateCmd.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = pbaff.PbCommissionSetting(cmd.Result)
	return nil
}

func GetProductPromotionByProductID(ctx context.Context, q *wrapaff.GetProductPromotionByProductIDEndpoint) error {
	panic("IMPLEMENT ME")
}

func AffiliateGetProducts(ctx context.Context, q *wrapaff.AffiliateGetProductsEndpoint) error {
	paging := q.Paging.CMPaging()
	query := &catalog.ListShopProductsWithVariantsQuery{
		ShopID:  modeletop.EtopTradingAccountID,
		Paging:  *paging,
		Filters: pbcm.ToFilters(q.Filters),
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}

	var productIds []int64
	for _, product := range query.Result.Products {
		productIds = append(productIds, product.ProductID)
	}

	tradingCommissionMap := GetShopCommissionSettingsByProducts(ctx, modeletop.EtopTradingAccountID, productIds)
	affCommissionMap := GetShopCommissionSettingsByProducts(ctx, q.Context.Affiliate.ID, productIds)
	shopPromotionMap := GetShopProductPromotionMapByProductIDs(ctx, modeletop.EtopTradingAccountID, productIds)

	var products []*pbaff.AffiliateProductResponse
	for _, product := range query.Result.Products {
		tradingCommissionSetting := tradingCommissionMap[product.ProductID]
		affCommissionSetting := affCommissionMap[product.ProductID]
		shopPromotion := shopPromotionMap[product.ProductID]

		var pbTradingCommissionSetting *pbaff.CommissionSetting = nil
		if tradingCommissionSetting != nil {
			pbTradingCommissionSetting = pbaff.PbCommissionSetting(tradingCommissionSetting)
		}
		var pbAffCommissionSetting *pbaff.CommissionSetting = nil
		if affCommissionSetting != nil {
			pbAffCommissionSetting = pbaff.PbCommissionSetting(affCommissionSetting)
		}
		var pbShopPromotion *pbaff.ProductPromotion = nil
		if shopPromotion != nil {
			pbShopPromotion = pbaff.PbProductPromotion(shopPromotion)
		}

		products = append(products, &pbaff.AffiliateProductResponse{
			Product:                    pbshop.PbShopProductWithVariants(product),
			ShopCommissionSetting:      pbTradingCommissionSetting,
			AffiliateCommissionSetting: pbAffCommissionSetting,
			Promotion:                  pbShopPromotion,
		})
	}

	q.Result = &pbaff.AffiliateGetProductsResponse{
		Paging:   pbcm.PbPageInfo(paging, query.Result.Count),
		Products: products,
	}

	return nil
}

func GetShopCommissionSettingsByProducts(ctx context.Context, accountID int64, productIds []int64) map[int64]*affiliate.CommissionSetting {
	getShopCommissionByProductIDsQuery := &affiliate.GetCommissionByProductIDsQuery{
		AccountID:  accountID,
		ProductIDs: productIds,
	}
	if err := affiliateQuery.Dispatch(ctx, getShopCommissionByProductIDsQuery); err != nil {
		return map[int64]*affiliate.CommissionSetting{}
	}

	shopCommissionMap := map[int64]*affiliate.CommissionSetting{}
	for _, e := range getShopCommissionByProductIDsQuery.Result {
		shopCommissionMap[e.ProductID] = e
	}

	return shopCommissionMap
}

func GetShopProductPromotionMapByProductIDs(ctx context.Context, shopID int64, productIDs []int64) map[int64]*affiliate.ProductPromotion {
	promotionsQ := &affiliate.GetShopProductPromotionByProductIDsQuery{
		ShopID:     shopID,
		ProductIDs: productIDs,
	}
	if err := affiliateQuery.Dispatch(ctx, promotionsQ); err != nil {
		return map[int64]*affiliate.ProductPromotion{}
	}

	promotionMap := map[int64]*affiliate.ProductPromotion{}
	for _, e := range promotionsQ.Result {
		promotionMap[e.ProductID] = e
	}
	return promotionMap
}

func GetCommissionSettingByReferralCode(ctx context.Context, referralCode string, productID int64) (*affiliate.CommissionSetting, error) {
	referralQ := &affiliate.GetAffiliateAccountReferralByCodeQuery{Code: referralCode}
	if err := affiliateQuery.Dispatch(ctx, referralQ); err != nil {
		return nil, err
	}
	commissionSettingQ := &affiliate.GetCommissionByProductIDQuery{
		AccountID: referralQ.Result.AffiliateID,
		ProductID: productID,
	}
	if err := affiliateQuery.Dispatch(ctx, commissionSettingQ); err != nil {
		return nil, err
	}
	return commissionSettingQ.Result, nil
}

func CreateReferralCode(ctx context.Context, q *wrapaff.CreateReferralCodeEndpoint) error {
	cmd := &affiliate.CreateAffiliateReferralCodeCommand{
		AffiliateAccountID: q.Context.Affiliate.ID,
		Code:               q.Code,
	}
	if err := affiliateCmd.Dispatch(ctx, cmd); err != nil {
		return err
	}

	q.Result = pbaff.PbReferralCode(cmd.Result)

	return nil
}

func GetReferralCodes(ctx context.Context, q *wrapaff.GetReferralCodesEndpoint) error {
	query := &affiliate.GetAffiliateAccountReferralCodesQuery{
		AffiliateAccountID: q.Context.Affiliate.ID,
	}
	if err := affiliateQuery.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = &pbaff.GetReferralCodesResponse{
		ReferralCodes: pbaff.PbReferralCodes(query.Result),
	}

	return nil
}

func GetReferrals(ctx context.Context, q *wrapaff.GetReferralsEndpoint) error {
	referralQ := &affiliate.GetReferralsByReferralIDQuery{
		ID: q.Context.Affiliate.ID,
	}
	if err := affiliateQuery.Dispatch(ctx, referralQ); err != nil {
		return err
	}

	var affiliateIDs []int64
	for _, userReferral := range referralQ.Result {
		userQ := &identity.GetAffiliatesByOwnerIDQuery{
			ID: userReferral.UserID,
		}
		if err := identityQuery.Dispatch(ctx, userQ); err == nil {
			affiliateIDs = append(affiliateIDs, userQ.Result[0].ID)
		}
	}

	affiliateQ := &identity.GetAffiliatesByIDsQuery{AffiliateIDs: affiliateIDs}
	if err := identityQuery.Dispatch(ctx, affiliateQ); err != nil {
		return err
	}

	var referrals []*pbaff.Referral
	for _, aff := range affiliateQ.Result {
		pbAffiliate := pbaff.PbReferral(aff)
		referrals = append(referrals, pbAffiliate)
	}

	q.Result = &pbaff.GetReferralsResponse{
		Referrals: referrals,
	}
	return nil
}
