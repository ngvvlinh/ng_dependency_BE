package convert

import (
	"time"

	"etop.vn/backend/com/fabo/api"
	"etop.vn/backend/com/fabo/main/fbuser/model"
	cm "etop.vn/backend/pkg/common"
)

func ConvertMeToUser(me *api.Me) *model.FbUser {
	if me == nil {
		return nil
	}
	return &model.FbUser{
		ID:         cm.NewID(),
		ExternalID: me.ID,
		UserID:     0,
		Info: model.FBUserInfo{
			Name:      me.Name,
			FirstName: me.FirstName,
			LastName:  me.LastName,
			ShortName: me.ShortName,
			ImgURL:    me.Picture.Data.Url,
		},
		Status:    0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
