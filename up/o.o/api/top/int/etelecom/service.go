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

	GetHotlines(context.Context, *types.GetHotLinesRequest) (*types.GetHotLinesResponse, error)
	GetCallLogs(context.Context, *types.GetCallLogsRequest) (*types.GetCallLogsResponse, error)
	CreateCallLog(context.Context, *CreateCallLogRequest) (*types.CallLog, error)

	SummaryEtelecom(context.Context, *SummaryEtelecomRequest) (*SummaryEtelecomResponse, error)

	CreateUserAndAssignExtension(context.Context, *CreateUserAndAssignExtensionRequest) (*cm.MessageResponse, error)

	CreateTenant(context.Context, *CreateTenantRequest) (*types.Tenant, error)
	GetTenant(context.Context, *GetTenantRequest) (*types.Tenant, error)

	RemoveUserOfExtension(context.Context, *RemoveUserOfExtensionRequest) (*cm.UpdatedResponse, error)
	AssignUserToExtension(context.Context, *AssignUserToExtensionRequest) (*cm.UpdatedResponse, error)
	DestroyCallSession(context.Context, *DestroyCallSessionRequest) (*cm.UpdatedResponse, error)

	ActionCall(context.Context, *cm.Empty) (*types.ActionCallResponse, error)
}

// +apix:path=/etelecom.User
type EtelecomUserService interface {
	GetUserSetting(context.Context, *cm.Empty) (*types.EtelecomUserSetting, error)
}
