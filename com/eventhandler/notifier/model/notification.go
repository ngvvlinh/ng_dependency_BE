package model

import (
	"encoding/json"

	"o.o/api/top/types/etc/notifier_entity"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/extservice/onesignal"
	"o.o/capi/dot"
)

type CreateNotificationArgs struct {
	AccountID        dot.ID
	UserID           dot.ID
	Title            string
	Message          string
	EntityID         dot.ID
	Entity           notifier_entity.NotifierEntity
	SendNotification bool
	MetaData         json.RawMessage
	TopicType        string
}

type CreateNotificationsArgs struct {
	SendAll          bool
	AccountIDs       []dot.ID
	Title            string
	Message          string
	EntityID         dot.ID
	Entity           notifier_entity.NotifierEntity
	SendNotification bool
	MetaData         json.RawMessage
	TopicType        string
}

type GetNotificationArgs struct {
	AccountID dot.ID
	ID        dot.ID
}

type GetNotificationsArgs struct {
	AccountID dot.ID
	Paging    *cm.Paging
}

type UpdateNotificationsArgs struct {
	IDs    []dot.ID
	IsRead bool
}

type SendNotificationCommand struct {
	Request *CreateNotificationRequest

	Result *onesignal.CreateNotificationResponse
}

type CreateNotificationRequest struct {
	Title             string
	Content           string
	Data              json.RawMessage
	ExternalDeviceIDs []string
	WebURL            string
}

func (c *CreateNotificationRequest) ToOnesignalModel() *onesignal.CreateNotificationRequest {
	return &onesignal.CreateNotificationRequest{
		IncludePlayerIDs: c.ExternalDeviceIDs,
		Headings: onesignal.MultipleContentLanguages{
			VI: c.Title,
			EN: c.Title,
		},
		Contents: onesignal.MultipleContentLanguages{
			VI: c.Content,
			EN: c.Content,
		},
		Data:   c.Data,
		WebURL: c.WebURL,
	}
}

type CreateUserNotiSettingArgs struct {
	UserID dot.ID
}
