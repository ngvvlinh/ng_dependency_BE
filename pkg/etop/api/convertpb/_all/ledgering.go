package _all

import (
	"o.o/api/main/ledgering"
	apishop "o.o/api/top/int/shop"
	identitysharemodel "o.o/backend/com/main/identity/sharemodel"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/api/convertpb"
)

func PbLedger(m *ledgering.ShopLedger) *apishop.Ledger {
	if m == nil {
		return nil
	}
	return &apishop.Ledger{
		Id:          m.ID,
		Name:        m.Name,
		BankAccount: convertpb.PbBankAccount((*identitysharemodel.BankAccount)(m.BankAccount)),
		Note:        m.Note,
		Type:        m.Type,
		CreatedBy:   m.CreatedBy,
		CreatedAt:   cmapi.PbTime(m.CreatedAt),
		UpdatedAt:   cmapi.PbTime(m.UpdatedAt),
	}
}

func PbLedgers(ms []*ledgering.ShopLedger) []*apishop.Ledger {
	res := make([]*apishop.Ledger, len(ms))
	for i, m := range ms {
		res[i] = PbLedger(m)
	}
	return res
}
