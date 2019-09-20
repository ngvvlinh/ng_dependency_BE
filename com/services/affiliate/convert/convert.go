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
