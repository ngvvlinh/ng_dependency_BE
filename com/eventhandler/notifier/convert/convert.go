package convert

import (
	"o.o/api/main/notify"
	"o.o/backend/com/eventhandler/notifier/model"
)

// +gen:convert: o.o/backend/com/eventhandler/notifier/model -> o.o/api/main/notify/types
// +gen:convert: o.o/backend/com/eventhandler/notifier/model

func Convert_model_UserNotiSetting_To_api_UserNotiSetting(setting *model.UserNotiSetting) *notify.UserNotiSetting {
	return &notify.UserNotiSetting{
		UserID:        setting.UserID,
		DisableTopics: setting.DisableTopics,
	}
}

func Convert_api_UserNotiSetting_To_model_UserNotiSetting(setting *notify.UserNotiSetting) *model.UserNotiSetting {
	return &model.UserNotiSetting{
		UserID:        setting.UserID,
		DisableTopics: setting.DisableTopics,
	}
}
