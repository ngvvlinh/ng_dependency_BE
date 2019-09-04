package affiliate

import (
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
