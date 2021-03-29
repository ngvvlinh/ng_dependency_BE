package usersetting

import (
	"context"

	"o.o/api/meta"
	"o.o/api/top/types/etc/charge_type"
	"o.o/capi/dot"
)

// +gen:api

type Aggregate interface {
	UpdateUserSetting(context.Context, *UpdateUserSettingArgs) (*UserSetting, error)
}

type QueryService interface {
	GetUserSetting(ctx context.Context, userID dot.ID) (*UserSetting, error)
	ListUserSettings(context.Context, *ListUserSettingsArgs) (*ListUserSettingsResponse, error)
}

type UpdateUserSettingArgs struct {
	UserID              dot.ID
	ExtensionChargeType charge_type.ChargeType
}

type ListUserSettingsArgs struct {
	UserIDs []dot.ID
	Paging  meta.Paging
}

type ListUserSettingsResponse struct {
	UserSettings []*UserSetting
	Paging       meta.PageInfo
}
