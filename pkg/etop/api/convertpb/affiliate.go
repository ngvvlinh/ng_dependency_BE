package convertpb

import (
	"etop.vn/api/main/identity"
	"etop.vn/api/services/affiliate"
	etop "etop.vn/api/top/int/etop"
	affiliate3 "etop.vn/api/top/services/affiliate"
	"etop.vn/backend/pkg/common/cmapi"
	"etop.vn/backend/pkg/etop/model"
)

func PbCommissionSetting(m *affiliate.CommissionSetting) *affiliate3.CommissionSetting {
	if m == nil {
		return nil
	}
	out := &affiliate3.CommissionSetting{
		ProductId: m.ProductID,
		Amount:    m.Amount,
		Unit:      m.Unit,
		CreatedAt: cmapi.PbTime(m.CreatedAt),
		UpdatedAt: cmapi.PbTime(m.UpdatedAt),
	}
	return out
}

func PbProductPromotion(m *affiliate.ProductPromotion) *affiliate3.ProductPromotion {
	if m == nil {
		return nil
	}
	out := &affiliate3.ProductPromotion{
		Product:   nil,
		Id:        m.ID,
		ProductId: m.ProductID,
		Amount:    m.Amount,
		Unit:      m.Unit,
		Type:      m.Type,
	}
	return out
}

func PbProductPromotions(ms []*affiliate.ProductPromotion) []*affiliate3.ProductPromotion {
	var out []*affiliate3.ProductPromotion
	if len(ms) == 0 {
		return out
	}
	for _, pp := range ms {
		out = append(out, PbProductPromotion(pp))
	}
	return out
}

func PbReferralCode(m *affiliate.AffiliateReferralCode) *affiliate3.ReferralCode {
	if m == nil {
		return nil
	}
	return &affiliate3.ReferralCode{
		Code: m.Code,
	}
}

func PbReferralCodes(ms []*affiliate.AffiliateReferralCode) []*affiliate3.ReferralCode {
	var out []*affiliate3.ReferralCode
	if len(ms) == 0 {
		return out
	}
	for _, rc := range ms {
		out = append(out, PbReferralCode(rc))
	}
	return out
}

func PbReferral(m *identity.Affiliate) *affiliate3.Referral {
	if m == nil {
		return nil
	}
	return &affiliate3.Referral{
		Name:            m.Name,
		Phone:           m.Phone,
		Email:           m.Email,
		OrderCount:      0,
		TotalRevenue:    0,
		TotalCommission: 0,
		CreatedAt:       cmapi.PbTime(m.CreatedAt),
	}
}

func PbSupplyCommissionSetting(m *affiliate.SupplyCommissionSetting) *affiliate3.SupplyCommissionSetting {
	if m == nil {
		return nil
	}
	return &affiliate3.SupplyCommissionSetting{
		ProductId:                m.ProductID,
		Level1DirectCommission:   m.Level1DirectCommission,
		Level1IndirectCommission: m.Level1IndirectCommission,
		Level2DirectCommission:   m.Level2DirectCommission,
		Level2IndirectCommission: m.Level2IndirectCommission,
		DependOn:                 m.DependOn,
		Level1LimitCount:         m.Level1LimitCount,
		MLifetimeDuration:        PbDuration(m.MLifetimeDuration),
		MLevel1LimitDuration:     PbDuration(m.MLevel1LimitDuration),
		Group:                    m.Group,
		CreatedAt:                cmapi.PbTime(m.CreatedAt),
		UpdatedAt:                cmapi.PbTime(m.UpdatedAt),
	}
}

func PbDuration(m *affiliate.Duration) *affiliate3.SupplyCommissionSettingDurationObject {
	if m == nil {
		return nil
	}
	return &affiliate3.SupplyCommissionSettingDurationObject{
		Duration: m.Duration,
		Type:     m.Type,
	}
}

func PbSellerCommission(m *affiliate.SellerCommission) *affiliate3.SellerCommission {
	if m == nil {
		return nil
	}
	return &affiliate3.SellerCommission{
		Id:          m.ID,
		Value:       m.Amount,
		Description: m.Description,
		Note:        m.Note,
		Status:      Pb4(model.Status4(m.Status)),
		Type:        m.Type,
		OValue:      m.OValue,
		OBaseValue:  m.OBaseValue,
		Product:     nil,
		Order:       nil,
		FromSeller:  nil,
		ValidAt:     cmapi.PbTime(m.ValidAt),
		CreatedAt:   cmapi.PbTime(m.CreatedAt),
		UpdatedAt:   cmapi.PbTime(m.UpdatedAt),
	}
}

func Convert_core_Affiliate_To_api_Affiliate(in *identity.Affiliate) *etop.Affiliate {
	if in == nil {
		return nil
	}
	return &etop.Affiliate{
		Id:          in.ID,
		Name:        in.Name,
		Status:      Pb3(model.Status3(in.Status)),
		IsTest:      in.IsTest != 0,
		Phone:       in.Phone,
		Email:       in.Email,
		BankAccount: Convert_core_BankAccount_To_api_BankAccount(in.BankAccount),
	}
}
