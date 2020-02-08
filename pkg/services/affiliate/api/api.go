package api

import (
	"context"
	"strings"

	"etop.vn/api/main/catalog"
	"etop.vn/api/main/identity"
	"etop.vn/api/meta"
	"etop.vn/api/services/affiliate"
	apiaffiliate "etop.vn/api/top/services/affiliate"
	ordermodelx "etop.vn/backend/com/main/ordering/modelx"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/cmapi"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/api/convertpb"
	pbshop "etop.vn/backend/pkg/etop/api/shop"
	modeletop "etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
	"etop.vn/common/l"
)

func init() {
	bus.AddHandlers("",
		userService.UpdateReferral,

		tradingService.TradingGetProducts,
		tradingService.GetTradingProductPromotions,
		tradingService.CreateOrUpdateTradingCommissionSetting,
		tradingService.GetTradingProductPromotionByProductIDs,
		tradingService.CreateTradingProductPromotion,
		tradingService.UpdateTradingProductPromotion,

		shopService.GetProductPromotion,
		shopService.ShopGetProducts,
		shopService.CheckReferralCodeValid,

		affiliateService.GetCommissions,
		affiliateService.NotifyNewShopPurchase,
		affiliateService.GetTransactions,
		affiliateService.CreateOrUpdateAffiliateCommissionSetting,
		affiliateService.GetProductPromotionByProductID,
		affiliateService.AffiliateGetProducts,
		affiliateService.CreateReferralCode,
		affiliateService.GetReferralCodes,
		affiliateService.GetReferrals,
	)
}

var ll = l.New()

var (
	catalogQuery   catalog.QueryBus
	affiliateCmd   affiliate.CommandBus
	affiliateQuery affiliate.QueryBus
	identityQuery  identity.QueryBus
)

type UserService struct{}
type TradingService struct{}
type ShopService struct{}
type AffiliateService struct{}

var userService = &UserService{}
var tradingService = &TradingService{}
var shopService = &ShopService{}
var affiliateService = &AffiliateService{}

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

func (s *UserService) UpdateReferral(ctx context.Context, q *UpdateReferralEndpoint) error {
	cmd := &affiliate.CreateOrUpdateUserReferralCommand{
		UserID:           q.Context.UserID,
		ReferralCode:     q.ReferralCode,
		SaleReferralCode: q.SaleReferralCode,
	}
	if err := affiliateCmd.Dispatch(ctx, cmd); err != nil {
		return err
	}

	q.Result = &apiaffiliate.UserReferral{
		UserId:           cmd.Result.UserID,
		ReferralCode:     cmd.Result.ReferralCode,
		SaleReferralCode: cmd.Result.SaleReferralCode,
	}
	return nil
}

func (s *TradingService) TradingGetProducts(ctx context.Context, q *TradingGetProductsEndpoint) error {
	if q.Context.Shop.ID != modeletop.EtopTradingAccountID {
		return cm.Errorf(cm.PermissionDenied, nil, "PermissionDenied")
	}
	paging := cmapi.CMPaging(q.Paging)
	query := &catalog.ListShopProductsWithVariantsQuery{
		ShopID:  modeletop.EtopTradingAccountID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(q.Filters),
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}

	var productIds []dot.ID
	for _, product := range query.Result.Products {
		productIds = append(productIds, product.ProductID)
	}

	supplyCommissionSettingMap := GetSupplyCommissionSettingByProductIdsMap(ctx, modeletop.EtopTradingAccountID, productIds)
	productPromotionMap := GetShopProductPromotionMapByProductIDs(ctx, modeletop.EtopTradingAccountID, productIds)
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
		productResult := pbshop.PbShopProductWithVariants(product)
		productResult, err := pbshop.PopulateTradingProductWithInventoryCount(ctx, productResult)
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
	if q.Context.Shop.ID != modeletop.EtopTradingAccountID {
		return cm.Errorf(cm.PermissionDenied, nil, "PermissionDenied")
	}
	paging := cmapi.CMPaging(q.Paging)
	query := &affiliate.ListShopProductPromotionsQuery{
		ShopID:  modeletop.EtopTradingAccountID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(q.Filters),
	}

	if err := affiliateQuery.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = &apiaffiliate.GetProductPromotionsResponse{
		Paging:     cmapi.PbPageInfo(paging),
		Promotions: convertpb.PbProductPromotions(query.Result.Promotions),
	}
	return nil
}

func (s *TradingService) CreateOrUpdateTradingCommissionSetting(ctx context.Context, q *CreateOrUpdateTradingCommissionSettingEndpoint) error {
	if q.Context.Shop.ID != modeletop.EtopTradingAccountID {
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

	if err := affiliateCmd.Dispatch(ctx, cmd); err != nil {
		return err
	}

	q.Result = convertpb.PbSupplyCommissionSetting(cmd.Result)

	return nil
}

func (s *TradingService) GetTradingProductPromotionByProductIDs(ctx context.Context, q *GetTradingProductPromotionByProductIDsEndpoint) error {
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
	q.Result = &apiaffiliate.GetTradingProductPromotionByIDsResponse{
		Promotions: convertpb.PbProductPromotions(productPromotionsQ.Result),
	}
	return nil
}

func (s *TradingService) CreateTradingProductPromotion(ctx context.Context, q *CreateTradingProductPromotionEndpoint) error {
	if q.Context.Shop.ID != modeletop.EtopTradingAccountID {
		return cm.Errorf(cm.PermissionDenied, nil, "PermissionDenied")
	}

	if err := affiliateQuery.Dispatch(ctx, &affiliate.GetShopProductPromotionQuery{
		ShopID:    q.Context.Shop.ID,
		ProductID: q.ProductId,
	}); err == nil {
		return cm.Errorf(cm.AlreadyExists, nil, "Sản phẩm đã có chương trình khuyến mãi")
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
	q.Result = convertpb.PbProductPromotion(cmd.Result)
	return nil
}

func (s *TradingService) UpdateTradingProductPromotion(ctx context.Context, q *UpdateTradingProductPromotionEndpoint) error {
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
	q.Result = convertpb.PbProductPromotion(cmd.Result)
	return nil
}

func (s *ShopService) GetProductPromotion(ctx context.Context, q *GetProductPromotionEndpoint) error {
	promotionQuery := &affiliate.GetShopProductPromotionQuery{
		ShopID:    modeletop.EtopTradingAccountID,
		ProductID: q.ProductId,
	}
	if err := affiliateQuery.Dispatch(ctx, promotionQuery); err != nil {
		return err
	}
	var pbReferralDiscount *apiaffiliate.CommissionSetting
	if q.ReferralCode.Valid {
		commissionSetting, err := GetCommissionSettingByReferralCode(ctx, q.ReferralCode.String, q.ProductId)
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
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}

	var productIds []dot.ID
	for _, product := range query.Result.Products {
		productIds = append(productIds, product.ProductID)
	}
	productPromotionMap := GetShopProductPromotionMapByProductIDs(ctx, modeletop.EtopTradingAccountID, productIds)
	var products []*apiaffiliate.ShopProductResponse
	for _, product := range query.Result.Products {
		productPromotion := productPromotionMap[product.ProductID]
		var pbProductPromotion *apiaffiliate.ProductPromotion = nil
		if productPromotion != nil {
			pbProductPromotion = convertpb.PbProductPromotion(productPromotion)
		}
		productResult := pbshop.PbShopProductWithVariants(product)
		productResult, err := pbshop.PopulateTradingProductWithInventoryCount(ctx, productResult)
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
	if err := affiliateQuery.Dispatch(ctx, affiliateAccountReferralQ); err != nil {
		return cm.Errorf(cm.NotFound, nil, "Mã giới thiệu không hợp lệ")
	}

	if affiliateAccountReferralQ.Result.UserID == q.Context.Shop.OwnerID {
		return cm.Errorf(cm.ValidationFailed, nil, "Mã giới thiệu không hợp lệ")
	}

	promotionQuery := &affiliate.GetShopProductPromotionQuery{
		ShopID:    modeletop.EtopTradingAccountID,
		ProductID: q.ProductId,
	}
	_ = affiliateQuery.Dispatch(ctx, promotionQuery)

	commissionSetting, err := GetCommissionSettingByReferralCode(ctx, q.ReferralCode, q.ProductId)
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

func (s *AffiliateService) GetCommissions(ctx context.Context, q *GetCommissionsEndpoint) error {
	commissionQ := &affiliate.GetSellerCommissionsQuery{
		SellerID: q.Context.Affiliate.ID,
		Paging:   meta.Paging{},
		Filters:  cmapi.ToFilters(q.Filters),
	}
	if err := affiliateQuery.Dispatch(ctx, commissionQ); err != nil {
		return err
	}

	var pbCommissions []*apiaffiliate.SellerCommission

	for _, commission := range commissionQ.Result {
		pbCommission := convertpb.PbSellerCommission(commission)

		if commission.FromSellerID != 0 {
			affiliateQ := &identity.GetAffiliateByIDQuery{
				ID: commission.FromSellerID,
			}
			if err := identityQuery.Dispatch(ctx, affiliateQ); err == nil {
				pbCommission.FromSeller = convertpb.Convert_core_Affiliate_To_api_Affiliate(affiliateQ.Result)
			}
		}

		if commission.OrderID != 0 {
			orderQ := &ordermodelx.GetOrderQuery{
				OrderID:            commission.OrderID,
				ShopID:             commission.SupplyID,
				PartnerID:          0,
				IncludeFulfillment: false,
			}
			if err := bus.Dispatch(ctx, orderQ); err == nil {
				pbCommission.Order = convertpb.PbOrder(orderQ.Result.Order, nil, modeletop.TagEtop)
			}

			shopQ := &identity.GetShopByIDQuery{
				ID: commission.ShopID,
			}
			if err := identityQuery.Dispatch(ctx, shopQ); err == nil && pbCommission.Order != nil {
				pbCommission.Order.ShopName = shopQ.Result.Name
			}
			pbCommissions = append(pbCommissions, pbCommission)
		}
	}

	q.Result = &apiaffiliate.GetCommissionsResponse{
		Commissions: pbCommissions,
	}

	return nil
}

func (s *AffiliateService) NotifyNewShopPurchase(ctx context.Context, q *NotifyNewShopPurchaseEndpoint) error {
	panic("IMPLEMENT ME")
}

func (s *AffiliateService) GetTransactions(ctx context.Context, q *GetTransactionsEndpoint) error {
	panic("IMPLEMENT ME")
}

func (s *AffiliateService) CreateOrUpdateAffiliateCommissionSetting(ctx context.Context, q *CreateOrUpdateAffiliateCommissionSettingEndpoint) error {
	cmd := &affiliate.CreateOrUpdateCommissionSettingCommand{
		ProductID: q.ProductId,
		AccountID: q.Context.Affiliate.ID,
		Amount:    q.Amount,
		Unit:      q.Unit.Apply(""),
		Type:      "affiliate",
	}
	if err := affiliateCmd.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = convertpb.PbCommissionSetting(cmd.Result)
	return nil
}

func (s *AffiliateService) GetProductPromotionByProductID(ctx context.Context, q *GetProductPromotionByProductIDEndpoint) error {
	panic("IMPLEMENT ME")
}

func (s *AffiliateService) AffiliateGetProducts(ctx context.Context, q *AffiliateGetProductsEndpoint) error {
	paging := cmapi.CMPaging(q.Paging)
	query := &catalog.ListShopProductsWithVariantsQuery{
		ShopID:  modeletop.EtopTradingAccountID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(q.Filters),
	}
	if err := catalogQuery.Dispatch(ctx, query); err != nil {
		return err
	}

	var productIds []dot.ID
	for _, product := range query.Result.Products {
		productIds = append(productIds, product.ProductID)
	}

	tradingCommissionMap := GetSupplyCommissionSettingByProductIdsMap(ctx, modeletop.EtopTradingAccountID, productIds)
	affCommissionMap := GetShopCommissionSettingsByProducts(ctx, q.Context.Affiliate.ID, productIds)
	shopPromotionMap := GetShopProductPromotionMapByProductIDs(ctx, modeletop.EtopTradingAccountID, productIds)

	var products []*apiaffiliate.AffiliateProductResponse
	for _, product := range query.Result.Products {
		tradingCommissionSetting := tradingCommissionMap[product.ProductID]
		affCommissionSetting := affCommissionMap[product.ProductID]
		shopPromotion := shopPromotionMap[product.ProductID]

		var pbTradingCommissionSetting *apiaffiliate.CommissionSetting = nil
		if tradingCommissionSetting != nil {
			pbTradingCommissionSetting = &apiaffiliate.CommissionSetting{
				ProductId: tradingCommissionSetting.ProductID,
				Amount:    tradingCommissionSetting.Level1DirectCommission,
				Unit:      "percent",
			}
		}
		var pbAffCommissionSetting *apiaffiliate.CommissionSetting = nil
		if affCommissionSetting != nil {
			pbAffCommissionSetting = convertpb.PbCommissionSetting(affCommissionSetting)
		}
		var pbShopPromotion *apiaffiliate.ProductPromotion = nil
		if shopPromotion != nil {
			pbShopPromotion = convertpb.PbProductPromotion(shopPromotion)
		}

		products = append(products, &apiaffiliate.AffiliateProductResponse{
			Product:                    pbshop.PbShopProductWithVariants(product),
			ShopCommissionSetting:      pbTradingCommissionSetting,
			AffiliateCommissionSetting: pbAffCommissionSetting,
			Promotion:                  pbShopPromotion,
		})
	}

	q.Result = &apiaffiliate.AffiliateGetProductsResponse{
		Paging:   cmapi.PbPageInfo(paging),
		Products: products,
	}

	return nil
}

func GetShopCommissionSettingsByProducts(ctx context.Context, accountID dot.ID, productIds []dot.ID) map[dot.ID]*affiliate.CommissionSetting {
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

func GetSupplyCommissionSettingByProductIdsMap(ctx context.Context, shopID dot.ID, productIDs []dot.ID) map[dot.ID]*affiliate.SupplyCommissionSetting {
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

func GetShopProductPromotionMapByProductIDs(ctx context.Context, shopID dot.ID, productIDs []dot.ID) map[dot.ID]*affiliate.ProductPromotion {
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

func GetCommissionSettingByReferralCode(ctx context.Context, referralCode string, productID dot.ID) (*affiliate.CommissionSetting, error) {
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

func (s *AffiliateService) CreateReferralCode(ctx context.Context, q *CreateReferralCodeEndpoint) error {
	cmd := &affiliate.CreateAffiliateReferralCodeCommand{
		AffiliateAccountID: q.Context.Affiliate.ID,
		Code:               q.Code,
	}
	if err := affiliateCmd.Dispatch(ctx, cmd); err != nil {
		return err
	}

	q.Result = convertpb.PbReferralCode(cmd.Result)

	return nil
}

func (s *AffiliateService) GetReferralCodes(ctx context.Context, q *GetReferralCodesEndpoint) error {
	query := &affiliate.GetAffiliateAccountReferralCodesQuery{
		AffiliateAccountID: q.Context.Affiliate.ID,
	}
	if err := affiliateQuery.Dispatch(ctx, query); err != nil {
		return err
	}

	q.Result = &apiaffiliate.GetReferralCodesResponse{
		ReferralCodes: convertpb.PbReferralCodes(query.Result),
	}

	return nil
}

func (s *AffiliateService) GetReferrals(ctx context.Context, q *GetReferralsEndpoint) error {
	referralQ := &affiliate.GetReferralsByReferralIDQuery{
		ID: q.Context.Affiliate.ID,
	}
	if err := affiliateQuery.Dispatch(ctx, referralQ); err != nil {
		return err
	}

	var affiliateIDs []dot.ID
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

	var referrals []*apiaffiliate.Referral
	for _, aff := range affiliateQ.Result {
		pbAffiliate := convertpb.PbReferral(aff)
		referrals = append(referrals, pbAffiliate)
	}

	q.Result = &apiaffiliate.GetReferralsResponse{
		Referrals: referrals,
	}
	return nil
}
