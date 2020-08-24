package notify

import (
	"context"

	"o.o/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateUserNotifySetting(context.Context, *CreateUserNotifySettingArgs) (*UserNotiSetting, error)
	GetOrCreateUserNotifySetting(context.Context, *GetOrCreateUserNotifySettingArgs) (*UserNotiSetting, error)
	DisableTopic(context.Context, *DisableTopicArgs) (*UserNotiSetting, error)
	EnableTopic(context.Context, *EnableTopicArgs) (*UserNotiSetting, error)
}

type QueryService interface {
	GetUserNotiSetting(context.Context, *GetUserNotiSettingArgs) (*UserNotiSetting, error)
}

type GetOrCreateUserNotifySettingArgs struct {
	UserID        dot.ID
	DisableTopics []string
}

type CreateUserNotifySettingArgs struct {
	UserID        dot.ID
	DisableTopics []string
}

type DisableTopicArgs struct {
	UserID dot.ID
	Topic  string
}

type EnableTopicArgs struct {
	UserID dot.ID
	Topic  string
}

type DisableShopArgs struct {
	UserID dot.ID
	ShopID dot.ID
}

type EnableShopArgs struct {
	UserID dot.ID
	ShopID dot.ID
}

type GetUserNotiSettingArgs struct {
	UserID dot.ID
}
