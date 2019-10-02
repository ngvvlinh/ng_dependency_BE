package affiliate

import (
	"etop.vn/api/main/identity"
	"etop.vn/api/services/affiliate"
	pbcm "etop.vn/backend/pb/common"
)

func PbCommissionSetting(m *affiliate.CommissionSetting) *CommissionSetting {
	if m == nil {
		return nil
	}
	out := &CommissionSetting{
		ProductId: m.ProductID,
		Amount:    m.Amount,
		Unit:      m.Unit,
		CreatedAt: pbcm.PbTime(m.CreatedAt),
		UpdatedAt: pbcm.PbTime(m.UpdatedAt),
	}
	return out
}

func PbProductPromotion(m *affiliate.ProductPromotion) *ProductPromotion {
	if m == nil {
		return nil
	}
	out := &ProductPromotion{
		Product:   nil,
		Id:        m.ID,
		ProductId: m.ProductID,
		Amount:    m.Amount,
		Unit:      m.Unit,
		Type:      m.Type,
	}
	return out
}

func PbProductPromotions(ms []*affiliate.ProductPromotion) []*ProductPromotion {
	var out []*ProductPromotion
	if len(ms) == 0 {
		return out
	}
	for _, pp := range ms {
		out = append(out, PbProductPromotion(pp))
	}
	return out
}

func PbReferralCode(m *affiliate.AffiliateReferralCode) *ReferralCode {
	if m == nil {
		return nil
	}
	return &ReferralCode{
		Code: m.Code,
	}
}

func PbReferralCodes(ms []*affiliate.AffiliateReferralCode) []*ReferralCode {
	var out []*ReferralCode
	if len(ms) == 0 {
		return out
	}
	for _, rc := range ms {
		out = append(out, PbReferralCode(rc))
	}
	return out
}

func PbReferral(m *identity.Affiliate) *Referral {
	if m == nil {
		return nil
	}
	return &Referral{
		Name:            m.Name,
		Phone:           m.Phone,
		Email:           m.Email,
		OrderCount:      0,
		TotalRevenue:    0,
		TotalCommission: 0,
		CreatedAt:       pbcm.PbTime(m.CreatedAt),
	}
}
