package etelecom

import (
	"context"

	"o.o/api/top/int/etelecom/types"
	cm "o.o/api/top/types/common"
)

// +gen:apix
// +gen:swagger:doc-path=etop/etelecom

// +apix:path=/shop.Etelecom
type EtelecomService interface {
	GetExtensions(context.Context, *types.GetExtensionsRequest) (*types.GetExtensionsResponse, error)
	CreateExtension(context.Context, *types.CreateExtensionRequest) (*types.Extension, error)
	CreateExtensionBySubscription(context.Context, *types.CreateExtensionBySubscriptionRequest) (*types.Extension, error)
	ExtendExtension(context.Context, *types.ExtendExtensionRequest) (*types.Extension, error)

	GetHotlines(context.Context, *cm.Empty) (*types.GetHotLinesResponse, error)
	GetCallLogs(context.Context, *types.GetCallLogsRequest) (*types.GetCallLogsResponse, error)
	CreateCallLog(context.Context, *CreateCallLogRequest) (*types.CallLog, error)

	SummaryEtelecom(context.Context, *SummaryEtelecomRequest) (*SummaryEtelecomResponse, error)

	CreateUserAndAssignExtension(context.Context, *CreateUserAndAssignExtensionRequest) (*cm.MessageResponse, error)
}

// +apix:path=/etelecom.User
type EtelecomUserService interface {
	UpdateUserSetting(context.Context, *UpdateUserSettingRequest) (*cm.UpdatedResponse, error)
	GetUserSetting(context.Context, *cm.Empty) (*EtelecomUserSetting, error)
}
