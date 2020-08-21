package api

import (
	"context"

	"o.o/api/main/catalog"
	"o.o/api/main/identity"
	"o.o/api/meta"
	"o.o/api/services/affiliate"
	apiaffiliate "o.o/api/top/services/affiliate"
	"o.o/api/top/types/etc/account_tag"
	ordermodelx "o.o/backend/com/main/ordering/modelx"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etc/idutil"
	"o.o/backend/pkg/etop/api/convertpb"
	product2 "o.o/backend/pkg/etop/api/shop/product"
	"o.o/capi/dot"
)

type AffiliateService struct {
	AffiliateAggr  affiliate.CommandBus
	CatalogQuery   catalog.QueryBus
	AffiliateQuery affiliate.QueryBus
	IdentityQuery  identity.QueryBus
}

func (s *AffiliateService) Clone() *AffiliateService { res := *s; return &res }

func (s *AffiliateService) GetCommissions(ctx context.Context, q *GetCommissionsEndpoint) error {
	commissionQ := &affiliate.GetSellerCommissionsQuery{
		SellerID: q.Context.Affiliate.ID,
		Paging:   meta.Paging{},
		Filters:  cmapi.ToFilters(q.Filters),
	}
	if err := s.AffiliateQuery.Dispatch(ctx, commissionQ); err != nil {
		return err
	}

	var pbCommissions []*apiaffiliate.SellerCommission

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
	if err := s.AffiliateAggr.Dispatch(ctx, cmd); err != nil {
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

	tradingCommissionMap := GetSupplyCommissionSettingByProductIdsMap(ctx, s.AffiliateQuery, idutil.EtopTradingAccountID, productIds)
	affCommissionMap := GetShopCommissionSettingsByProducts(ctx, s.AffiliateQuery, q.Context.Affiliate.ID, productIds)
	shopPromotionMap := GetShopProductPromotionMapByProductIDs(ctx, s.AffiliateQuery, idutil.EtopTradingAccountID, productIds)

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
			Product:                    product2.PbShopProductWithVariants(product),
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

func (s *AffiliateService) CreateReferralCode(ctx context.Context, q *CreateReferralCodeEndpoint) error {
	cmd := &affiliate.CreateAffiliateReferralCodeCommand{
		AffiliateAccountID: q.Context.Affiliate.ID,
		Code:               q.Code,
	}
	if err := s.AffiliateAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}

	q.Result = convertpb.PbReferralCode(cmd.Result)

	return nil
}

func (s *AffiliateService) GetReferralCodes(ctx context.Context, q *GetReferralCodesEndpoint) error {
	query := &affiliate.GetAffiliateAccountReferralCodesQuery{
		AffiliateAccountID: q.Context.Affiliate.ID,
	}
	if err := s.AffiliateQuery.Dispatch(ctx, query); err != nil {
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
	if err := s.AffiliateQuery.Dispatch(ctx, referralQ); err != nil {
		return err
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
