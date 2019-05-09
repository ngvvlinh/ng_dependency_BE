package model

import (
	"encoding/json"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/onesignal"
)

type CreateNotificationArgs struct {
	AccountID        int64
	Title            string
	Message          string
	EntityID         int64
	Entity           NotiEntity
	SendNotification bool
}

type CreateNotificationsArgs struct {
	SendAll          bool
	AccountIDs       []int64
	Title            string
	Message          string
	EntityID         int64
	Entity           NotiEntity
	SendNotification bool
}

type GetNotificationArgs struct {
	AccountID int64
	ID        int64
}

type GetNotificationsArgs struct {
	AccountID int64
	Paging    *cm.Paging
}

type UpdateNotificationsArgs struct {
	IDs    []int64
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
