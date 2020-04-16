package convert

import (
	"time"

	"etop.vn/backend/com/fabo/api"
	"etop.vn/backend/com/fabo/main/fbpage/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/capi/dot"
)

func ConvertAccountsToFBPage(ins api.Accounts, shopID, userID dot.ID) model.FbPages {
	var outs []*model.FbPage

	for _, account := range ins.Accounts.Data {
		outs = append(outs, &model.FbPage{
			ID:           cm.NewID(),
			ExternalID:   account.Id,
			ShopID:       shopID,
			UserID:       userID,
			Name:         account.Name,
			Category:     account.Category,
			CategoryList: account.CategoryList,
			Tasks:        account.Tasks,
			Status:       0,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		})
	}

	return outs
}
