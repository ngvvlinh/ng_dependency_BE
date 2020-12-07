package etelecom

import (
	"context"
	"time"

	"o.o/api/etelecom/call_log_direction"
	"o.o/api/etelecom/call_state"
	"o.o/api/meta"
	cm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/status5"
	"o.o/capi/dot"
	"o.o/common/xerrors"
)

// +gen:api

type Aggregate interface {
	CreateExtension(context.Context, *CreateExtensionArgs) (*Extension, error)
	DeleteExtension(ctx context.Context, id dot.ID) error
	UpdateExternalExtensionInfo(context.Context, *UpdateExternalExtensionInfoArgs) error

	CreateCallLogFromCDR(context.Context, *CreateCallLogFromCDRArgs) (*CallLog, error)
}

type QueryService interface {
	GetHotline(context.Context, *GetHotlineArgs) (*Hotline, error)
	ListHotlines(context.Context, *ListHotlinesArgs) ([]*Hotline, error)
	ListBuiltinHotlines(context.Context, *cm.Empty) ([]*Hotline, error)

	GetExtension(context.Context, *GetExtensionArgs) (*Extension, error)
	ListExtensions(context.Context, *ListExtensionsArgs) ([]*Extension, error)

	// generate extension number => then use it for create external extension
	GetPrivateExtensionNumber(context.Context, *cm.Empty) (extensionNumber string, _ error)

	GetCallLogByExternalID(context.Context, *GetCallLogByExternalIDArgs) (*CallLog, error)

	ListCallLogs(context.Context, *ListCallLogsArgs) (*ListCallLogsResponse, error)
}

// +convert:create=Extension
type CreateExtensionArgs struct {
	ExtensionNumber   string
	UserID            dot.ID
	AccountID         dot.ID
	ExtensionPassword string
	HotlineID         dot.ID
}

func (args *CreateExtensionArgs) Validate() error {
	if args.UserID == 0 {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing user ID")
	}
	if args.AccountID == 0 {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing account ID")
	}
	if args.HotlineID == 0 {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing hotline ID")
	}
	return nil
}

// +convert:create=CallLog
type CreateCallLogFromCDRArgs struct {
	ExternalID         string
	StartedAt          time.Time
	EndedAt            time.Time
	Duration           int
	Caller             string
	Callee             string
	AudioURLs          []string
	ExternalDirection  string
	Direction          call_log_direction.CallLogDirection
	ExternalCallStatus string
	CallState          call_state.CallState
	CallStatus         status5.Status

	// use for find hotline_id & extension_id
	OwnerID      dot.ID
	ConnectionID dot.ID
}

func (args *CreateCallLogFromCDRArgs) Validate() error {
	if args.Callee == "" {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing callee")
	}
	if args.Caller == "" {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing caller")
	}
	if args.ConnectionID == 0 {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing connection ID")
	}
	return nil
}

type GetHotlineArgs struct {
	ID      dot.ID
	OwnerID dot.ID
}

type ListHotlinesArgs struct {
	OwnerID      dot.ID
	ConnectionID dot.ID
}

type GetExtensionArgs struct {
	ID        dot.ID
	UserID    dot.ID
	AccountID dot.ID
	HotlineID dot.ID
}

type ListExtensionsArgs struct {
	AccountIDs []dot.ID
	HotlineID  dot.ID
}

type UpdateExternalExtensionInfoArgs struct {
	ID                dot.ID
	HotlineID         dot.ID
	ExternalID        string
	ExtensionNumber   string
	ExtensionPassword string
}

type GetCallLogByExternalIDArgs struct {
	ExternalID string
}

type ListCallLogsArgs struct {
	HotlineIDs   []dot.ID
	ExtensionIDs []dot.ID
	Paging       meta.Paging
}

type ListCallLogsResponse struct {
	CallLogs []*CallLog
	Paging   meta.PageInfo
}
