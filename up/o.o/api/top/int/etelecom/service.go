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

	GetHotlines(context.Context, *cm.Empty) (*types.GetHotLinesResponse, error)
	GetCallLogs(context.Context, *types.GetCallLogsRequest) (*types.GetCallLogsResponse, error)

	SummaryEtelecom(context.Context, *SummaryEtelecomRequest) (*SummaryEtelecomResponse, error)

	CreateUserAndAssignExtension(context.Context, *CreateUserAndAssignExtensionRequest) (*cm.MessageResponse, error)
}
