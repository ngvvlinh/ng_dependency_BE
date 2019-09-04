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
