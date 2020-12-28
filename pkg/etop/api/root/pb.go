package root

import (
	"o.o/api/top/int/etop"
	"o.o/api/top/types/etc/account_tag"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/etop/authorize/claims"
	"o.o/capi/dot"
)

func MixAccount(claim claims.Claim, m *etop.MixedAccount) ([]dot.ID, error) {
	switch {
	case m == nil:
		// The same as default.

	case m.All:
		ids := make([]dot.ID, 0, len(claim.AccountIDs))
		for id := range claim.AccountIDs {
			ids = append(ids, id)
		}
		return ids, nil

	case m.AllShops:
		ids := make([]dot.ID, 0, len(claim.AccountIDs))
		for id, typ := range claim.AccountIDs {
			if typ == account_tag.TagShop {
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
	return []dot.ID{claim.AccountID}, nil
}
