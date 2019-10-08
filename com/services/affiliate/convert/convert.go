package convert

import (
	"etop.vn/api/services/affiliate"
	"etop.vn/backend/com/services/affiliate/model"
)

func CommissionSetting(in *model.CommissionSetting) *affiliate.CommissionSetting {
	if in == nil {
		return nil
	}

	return &affiliate.CommissionSetting{
		ProductID: in.ProductID,
		Amount:    in.Amount,
		Unit:      in.Unit,
		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
	}
}

func CommissionSettings(ins []*model.CommissionSetting) []*affiliate.CommissionSetting {
	if len(ins) == 0 {
		return []*affiliate.CommissionSetting{}
	}
	var results []*affiliate.CommissionSetting
	for _, in := range ins {
		results = append(results, CommissionSetting(in))
	}
	return results
}

func ProductPromotion(in *model.ProductPromotion) *affiliate.ProductPromotion {
	if in == nil {
		return nil
	}

	return &affiliate.ProductPromotion{
		ID:          in.ID,
		ProductID:   in.ProductID,
		Amount:      in.Amount,
		Unit:        in.Unit,
		Code:        in.Code,
		Description: in.Description,
		Note:        in.Note,
		Type:        in.Type,
		CreatedAt:   in.CreatedAt,
		UpdatedAt:   in.UpdatedAt,
	}
}

func ProductPromotions(ins []*model.ProductPromotion) []*affiliate.ProductPromotion {
	if len(ins) == 0 {
		return []*affiliate.ProductPromotion{}
	}
	var results []*affiliate.ProductPromotion
	for _, in := range ins {
		results = append(results, ProductPromotion(in))
	}
	return results
}

func AffiliateCommission(in *model.AffiliateCommission) *affiliate.AffiliateCommission {
	if in == nil {
		return nil
	}
	result := &affiliate.AffiliateCommission{
		ID:              in.ID,
		AffiliateID:     in.AffiliateID,
		FromAffiliateID: in.FromAffiliateID,
		ProductID:       in.ProductID,
		OrderId:         in.OrderId,
		Value:           in.Value,
		Description:     in.Description,
		Note:            in.Note,
		Type:            in.Type,
		Status:          in.Status,
		ValidAt:         in.ValidAt,
		CreatedAt:       in.CreatedAt,
		UpdatedAt:       in.UpdatedAt,
	}
	return result
}

func AffiliateCommissions(ins []*model.AffiliateCommission) []*affiliate.AffiliateCommission {
	if len(ins) == 0 {
		return []*affiliate.AffiliateCommission{}
	}
	var results []*affiliate.AffiliateCommission
	for _, in := range ins {
		results = append(results, AffiliateCommission(in))
	}
	return results
}

func AffiliateReferralCode(in *model.AffiliateReferralCode) *affiliate.AffiliateReferralCode {
	if in == nil {
		return nil
	}
	result := &affiliate.AffiliateReferralCode{
		ID:          in.ID,
		Code:        in.Code,
		AffiliateID: in.AffiliateID,
		UserID:      in.UserID,
		CreatedAt:   in.CreatedAt,
		UpdatedAt:   in.UpdatedAt,
	}
	return result
}

func AffiliateReferralCodes(ins []*model.AffiliateReferralCode) []*affiliate.AffiliateReferralCode {
	var results []*affiliate.AffiliateReferralCode
	if len(ins) == 0 {
		return results
	}
	for _, arc := range ins {
		results = append(results, AffiliateReferralCode(arc))
	}
	return results
}

func UserReferral(in *model.UserReferral) *affiliate.UserReferral {
	if in == nil {
		return nil
	}
	result := &affiliate.UserReferral{
		UserID:         in.UserID,
		ReferralID:     in.ReferralID,
		SaleReferralID: in.SaleReferralID,
		ReferralAt:     in.ReferralAt,
		ReferralSaleAt: in.SaleReferralAt,
		CreatedAt:      in.CreatedAt,
		UpdatedAt:      in.UpdatedAt,
	}
	return result
}

func UserReferrals(ins []*model.UserReferral) []*affiliate.UserReferral {
	var results []*affiliate.UserReferral
	if len(ins) == 0 {
		return results
	}
	for _, ur := range ins {
		results = append(results, UserReferral(ur))
	}
	return results
}

func Duration(in *model.DurationJSON) *affiliate.Duration {
	if in == nil {
		return nil
	}
	return &affiliate.Duration{
		Type:     in.Type,
		Duration: in.Duration,
	}
}

func SupplyCommissionSetting(in *model.SupplyCommissionSetting) *affiliate.SupplyCommissionSetting {
	if in == nil {
		return nil
	}
	return &affiliate.SupplyCommissionSetting{
		ShopID:                   in.ShopID,
		ProductID:                in.ProductID,
		Level1DirectCommission:   in.Level1DirectCommission,
		Level1IndirectCommission: in.Level1IndirectCommission,
		Level2DirectCommission:   in.Level2DirectCommission,
		Level2IndirectCommission: in.Level2IndirectCommission,
		DependOn:                 in.DependOn,
		Level1LimitCount:         in.Level1LimitCount,
		Level1LimitDuration:      in.Level1LimitDuration,
		MLevel1LimitDuration:     Duration(in.MLevel1LimitDuration),
		LifetimeDuration:         in.LifetimeDuration,
		MLifetimeDuration:        Duration(in.MLifetimeDuration),
		CreatedAt:                in.CreatedAt,
		UpdatedAt:                in.UpdatedAt,
	}
}
