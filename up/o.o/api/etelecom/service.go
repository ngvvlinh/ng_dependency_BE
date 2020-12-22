package etelecom

import (
	"context"
	"time"

	"o.o/api/etelecom/call_log_direction"
	"o.o/api/etelecom/call_state"
	"o.o/api/etelecom/mobile_network"
	"o.o/api/meta"
	cm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/status5"
	"o.o/capi/dot"
	"o.o/common/xerrors"
)

// +gen:api

type Aggregate interface {
	CreateHotline(context.Context, *CreateHotlineArgs) (*Hotline, error)
	UpdateHotlineInfo(context.Context, *UpdateHotlineInfoArgs) error

	CreateExtension(context.Context, *CreateExtensionArgs) (*Extension, error)
	DeleteExtension(ctx context.Context, id dot.ID) error
	UpdateExternalExtensionInfo(context.Context, *UpdateExternalExtensionInfoArgs) error

	UpdateCallLogPostage(context.Context, *UpdateCallLogPostageArgs) error
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
	GetCallLog(ctx context.Context, ID dot.ID) (*CallLog, error)
}

// +convert:create=Extension
type CreateExtensionArgs struct {
	ExtensionNumber   string
	UserID            dot.ID
	AccountID         dot.ID
	ExtensionPassword string
	HotlineID         dot.ID
	OwnerID           dot.ID
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
	TenantDomain      string
}

type GetCallLogByExternalIDArgs struct {
	ExternalID string
}

type ListCallLogsArgs struct {
	HotlineIDs   []dot.ID
	ExtensionIDs []dot.ID
	AccountID    dot.ID
	Paging       meta.Paging
}

type ListCallLogsResponse struct {
	CallLogs []*CallLog
	Paging   meta.PageInfo
}

type UpdateCallLogPostageArgs struct {
	ID                 dot.ID
	DurationForPostage int
	Postage            int
}

// +convert:create=Hotline
type CreateHotlineArgs struct {
	OwnerID      dot.ID
	Name         string
	Hotline      string
	Network      mobile_network.MobileNetwork
	ConnectionID dot.ID
	Status       status3.Status
	Description  string
	IsFreeCharge dot.NullBool
}

func (args *CreateHotlineArgs) Validate() error {
	if args.Hotline == "" {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing hotline number")
	}
	if args.ConnectionID == 0 {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing connection ID")
	}
	if args.Network == 0 {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing network")
	}
	return nil
}

type UpdateHotlineInfoArgs struct {
	ID           dot.ID
	IsFreeCharge dot.NullBool
	Name         string
	Description  string
	Status       status3.NullStatus
}
