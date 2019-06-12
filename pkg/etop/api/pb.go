package api

import (
	pbetop "etop.vn/backend/pb/etop"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/etop/authorize/claims"
	"etop.vn/backend/pkg/etop/model"
)

func MixAccount(claim *claims.Claim, m *pbetop.MixedAccount) ([]int64, error) {
	switch {
	case m == nil:
		// The same as default.

	case m.All:
		ids := make([]int64, 0, len(claim.AccountIDs))
		for id := range claim.AccountIDs {
			ids = append(ids, id)
		}
		return ids, nil

	case m.AllShops:
		ids := make([]int64, 0, len(claim.AccountIDs))
		for id, typ := range claim.AccountIDs {
			if typ == model.TagShop {
				ids = append(ids, id)
			}
		}
		return ids, nil

	case m.AllSuppliers:
		ids := make([]int64, 0, len(claim.AccountIDs))
		for id, typ := range claim.AccountIDs {
			if typ == model.TagSupplier {
				ids = append(ids, id)
			}
		}
		return ids, nil

	case len(m.Ids) > 0:
		accountIDs := claim.AccountIDs
		for _, id := range m.Ids {
			if accountIDs[id] == 0 {
				return nil, cm.Errorf(cm.InvalidArgument, nil, "Invalid MixedAccount.ids (%v)", id)
			}
		}
		return m.Ids, nil
	}

	if claim.AccountID == 0 {
		return nil, cm.Error(cm.InvalidArgument, "Expected account token", nil)
	}
	return []int64{claim.AccountID}, nil
}
