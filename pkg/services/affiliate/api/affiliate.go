package api

import (
	"context"

	"o.o/api/main/catalog"
	"o.o/api/main/identity"
	"o.o/api/meta"
	"o.o/api/services/affiliate"
	api "o.o/api/top/services/affiliate"
	pbcm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/account_tag"
	ordermodelx "o.o/backend/com/main/ordering/modelx"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etc/idutil"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/api/shop/product"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/capi/dot"
)

type AffiliateService struct {
	session.Session

	AffiliateAggr  affiliate.CommandBus
	CatalogQuery   catalog.QueryBus
	AffiliateQuery affiliate.QueryBus
	IdentityQuery  identity.QueryBus
}

func (s *AffiliateService) Clone() api.AffiliateService { res := *s; return &res }

func (s *AffiliateService) GetCommissions(ctx context.Context, q *pbcm.CommonListRequest) (*api.GetCommissionsResponse, error) {
	commissionQ := &affiliate.GetSellerCommissionsQuery{
		SellerID: s.SS.Affiliate().ID,
		Paging:   meta.Paging{},
		Filters:  cmapi.ToFilters(q.Filters),
	}
	if err := s.AffiliateQuery.Dispatch(ctx, commissionQ); err != nil {
		return nil, err
	}

	var pbCommissions []*api.SellerCommission

	for _, commission := range commissionQ.Result {
		pbCommission := convertpb.PbSellerCommission(commission)

		if commission.FromSellerID != 0 {
			affiliateQ := &identity.GetAffiliateByIDQuery{
				ID: commission.FromSellerID,
			}
			if err := s.IdentityQuery.Dispatch(ctx, affiliateQ); err == nil {
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
				pbCommission.Order = convertpb.PbOrder(orderQ.Result.Order, nil, account_tag.TagEtop)
			}

			shopQ := &identity.GetShopByIDQuery{
				ID: commission.ShopID,
			}
			if err := s.IdentityQuery.Dispatch(ctx, shopQ); err == nil && pbCommission.Order != nil {
				pbCommission.Order.ShopName = shopQ.Result.Name
			}
			pbCommissions = append(pbCommissions, pbCommission)
		}
	}

	result := &api.GetCommissionsResponse{
		Commissions: pbCommissions,
	}

	return result, nil
}

func (s *AffiliateService) NotifyNewShopPurchase(ctx context.Context, q *api.NotifyNewShopPurchaseRequest) (*api.NotifyNewShopPurchaseResponse, error) {
	panic("IMPLEMENT ME")
}

func (s *AffiliateService) GetTransactions(ctx context.Context, q *pbcm.CommonListRequest) (*api.GetTransactionsResponse, error) {
	panic("IMPLEMENT ME")
}

func (s *AffiliateService) CreateOrUpdateAffiliateCommissionSetting(ctx context.Context, q *api.CreateOrUpdateCommissionSettingRequest) (*api.CommissionSetting, error) {
	cmd := &affiliate.CreateOrUpdateCommissionSettingCommand{
		ProductID: q.ProductId,
		AccountID: s.SS.Affiliate().ID,
		Amount:    q.Amount,
		Unit:      q.Unit.Apply(""),
		Type:      "affiliate",
	}
	if err := s.AffiliateAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := convertpb.PbCommissionSetting(cmd.Result)
	return result, nil
}

func (s *AffiliateService) GetProductPromotionByProductID(ctx context.Context, q *api.GetProductPromotionByProductIDRequest) (*api.GetProductPromotionByProductIDResponse, error) {
	panic("IMPLEMENT ME")
}

func (s *AffiliateService) AffiliateGetProducts(ctx context.Context, q *pbcm.CommonListRequest) (*api.AffiliateGetProductsResponse, error) {
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

	tradingCommissionMap := GetSupplyCommissionSettingByProductIdsMap(ctx, s.AffiliateQuery, idutil.EtopTradingAccountID, productIds)
	affCommissionMap := GetShopCommissionSettingsByProducts(ctx, s.AffiliateQuery, s.SS.Affiliate().ID, productIds)
	shopPromotionMap := GetShopProductPromotionMapByProductIDs(ctx, s.AffiliateQuery, idutil.EtopTradingAccountID, productIds)

	var products []*api.AffiliateProductResponse
	for _, p := range query.Result.Products {
		tradingCommissionSetting := tradingCommissionMap[p.ProductID]
		affCommissionSetting := affCommissionMap[p.ProductID]
		shopPromotion := shopPromotionMap[p.ProductID]

		var pbTradingCommissionSetting *api.CommissionSetting = nil
		if tradingCommissionSetting != nil {
			pbTradingCommissionSetting = &api.CommissionSetting{
				ProductId: tradingCommissionSetting.ProductID,
				Amount:    tradingCommissionSetting.Level1DirectCommission,
				Unit:      "percent",
			}
		}
		var pbAffCommissionSetting *api.CommissionSetting = nil
		if affCommissionSetting != nil {
			pbAffCommissionSetting = convertpb.PbCommissionSetting(affCommissionSetting)
		}
		var pbShopPromotion *api.ProductPromotion = nil
		if shopPromotion != nil {
			pbShopPromotion = convertpb.PbProductPromotion(shopPromotion)
		}

		products = append(products, &api.AffiliateProductResponse{
			Product:                    product.PbShopProductWithVariants(p),
			ShopCommissionSetting:      pbTradingCommissionSetting,
			AffiliateCommissionSetting: pbAffCommissionSetting,
			Promotion:                  pbShopPromotion,
		})
	}

	result := &api.AffiliateGetProductsResponse{
		Paging:   cmapi.PbPageInfo(paging),
		Products: products,
	}

	return result, nil
}

func (s *AffiliateService) CreateReferralCode(ctx context.Context, q *api.CreateReferralCodeRequest) (*api.ReferralCode, error) {
	cmd := &affiliate.CreateAffiliateReferralCodeCommand{
		AffiliateAccountID: s.SS.Affiliate().ID,
		Code:               q.Code,
	}
	if err := s.AffiliateAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	result := convertpb.PbReferralCode(cmd.Result)

	return result, nil
}

func (s *AffiliateService) GetReferralCodes(ctx context.Context, q *pbcm.CommonListRequest) (*api.GetReferralCodesResponse, error) {
	query := &affiliate.GetAffiliateAccountReferralCodesQuery{
		AffiliateAccountID: s.SS.Affiliate().ID,
	}
	if err := s.AffiliateQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}

	result := &api.GetReferralCodesResponse{
		ReferralCodes: convertpb.PbReferralCodes(query.Result),
	}

	return result, nil
}

func (s *AffiliateService) GetReferrals(ctx context.Context, q *pbcm.CommonListRequest) (*api.GetReferralsResponse, error) {
	referralQ := &affiliate.GetReferralsByReferralIDQuery{
		ID: s.SS.Affiliate().ID,
	}
	if err := s.AffiliateQuery.Dispatch(ctx, referralQ); err != nil {
		return nil, err
	}

	var affiliateIDs []dot.ID
	for _, userReferral := range referralQ.Result {
		userQ := &identity.GetAffiliatesByOwnerIDQuery{
			ID: userReferral.UserID,
		}
		if err := s.IdentityQuery.Dispatch(ctx, userQ); err == nil {
			affiliateIDs = append(affiliateIDs, userQ.Result[0].ID)
		}
	}

	affiliateQ := &identity.GetAffiliatesByIDsQuery{AffiliateIDs: affiliateIDs}
	if err := s.IdentityQuery.Dispatch(ctx, affiliateQ); err != nil {
		return nil, err
	}

	var referrals []*api.Referral
	for _, aff := range affiliateQ.Result {
		pbAffiliate := convertpb.PbReferral(aff)
		referrals = append(referrals, pbAffiliate)
	}

	result := &api.GetReferralsResponse{
		Referrals: referrals,
	}
	return result, nil
}
