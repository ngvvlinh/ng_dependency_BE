package etelecom

import (
	"context"
	"time"

	"o.o/api/etelecom/call_direction"
	"o.o/api/etelecom/call_state"
	"o.o/api/etelecom/mobile_network"
	"o.o/api/meta"
	cm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/payment_method"
	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
	"o.o/common/xerrors"
)

// +gen:api

type Aggregate interface {
	CreateHotline(context.Context, *CreateHotlineArgs) (*Hotline, error)
	UpdateHotlineInfo(context.Context, *UpdateHotlineInfoArgs) error

	CreateExtension(context.Context, *CreateExtensionArgs) (*Extension, error)
	CreateExtensionBySubscription(context.Context, *CreateExtenstionBySubscriptionArgs) (*Extension, error)
	ExtendExtension(context.Context, *ExtendExtensionArgs) (*Extension, error)
	DeleteExtension(ctx context.Context, id dot.ID) error
	UpdateExternalExtensionInfo(context.Context, *UpdateExternalExtensionInfoArgs) error

	UpdateCallLogPostage(context.Context, *UpdateCallLogPostageArgs) error
	CreateOrUpdateCallLogFromCDR(context.Context, *CreateOrUpdateCallLogFromCDRArgs) (*CallLog, error)
	CreateCallLog(context.Context, *CreateCallLogArgs) (*CallLog, error)
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
	SubscriptionID    dot.ID
	ExpiresAt         time.Time
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
type CreateOrUpdateCallLogFromCDRArgs struct {
	ExternalID         string
	StartedAt          time.Time
	EndedAt            time.Time
	Duration           int
	Caller             string
	Callee             string
	AudioURLs          []string
	ExternalDirection  string
	Direction          call_direction.CallDirection
	ExternalCallStatus string
	CallState          call_state.CallState
	ExternalSessionID  string
	ExtensionID        dot.ID
	HotlineID          dot.ID

	// use for find hotline_id & extension_id
	OwnerID      dot.ID
	ConnectionID dot.ID
}

func (args *CreateOrUpdateCallLogFromCDRArgs) Validate() error {
	if args.Callee == "" {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing callee")
	}
	if args.Caller == "" {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing caller")
	}
	if args.ConnectionID == 0 {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing connection ID")
	}
	if args.HotlineID == 0 {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing hotline ID")
	}
	if args.ExternalSessionID == "" {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing external session ID")
	}
	return nil
}

type CreateCallLogArgs struct {
	ExternalSessionID string
	Direction         call_direction.CallDirection
	Caller            string
	Callee            string
	ExtensionID       dot.ID
	AccountID         dot.ID
	ContactID         dot.ID
	CallState         call_state.CallState
}

func (args *CreateCallLogArgs) Validate() error {
	if args.ExtensionID == 0 {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing extension ID")
	}
	if args.ExternalSessionID == "" {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing external session ID")
	}
	if args.Direction == 0 {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing call direction")
	}
	if args.Callee == "" || args.Caller == "" {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing caller or callee")
	}
	if args.AccountID == 0 {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing account ID")
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
	ID             dot.ID
	UserID         dot.ID
	AccountID      dot.ID
	HotlineID      dot.ID
	SubscriptionID dot.ID
}

type ListExtensionsArgs struct {
	AccountIDs       []dot.ID
	HotlineIDs       []dot.ID
	ExtensionNumbers []string
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

type CreateExtenstionBySubscriptionArgs struct {
	SubscriptionID     dot.ID
	SubscriptionPlanID dot.ID
	PaymentMethod      payment_method.PaymentMethod
	AccountID          dot.ID

	UserID    dot.ID
	HotlineID dot.ID
	OwnerID   dot.ID
}

func (args *CreateExtenstionBySubscriptionArgs) Validate() error {
	if args.UserID == 0 {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing user ID")
	}
	if args.AccountID == 0 {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing account ID")
	}
	if args.HotlineID == 0 {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing hotline ID")
	}
	if args.SubscriptionID == 0 && args.SubscriptionPlanID == 0 {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Please provide subscription_id or subscription_plan_id ")
	}
	return nil
}

type ExtendExtensionArgs struct {
	ExtensionID dot.ID
	UserID      dot.ID
	AccountID   dot.ID

	SubscriptionID     dot.ID
	SubscriptionPlanID dot.ID
	PaymentMethod      payment_method.PaymentMethod
}

func (args *ExtendExtensionArgs) Validate() error {
	if args.UserID == 0 {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing user ID")
	}
	if args.AccountID == 0 {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing account ID")
	}
	if args.ExtensionID == 0 {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing extension ID")
	}
	if args.PaymentMethod != payment_method.Balance {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Phương thức thanh toán không hợp lệ. Chỉ hỗ trợ thanh toán bằng số dư.")
	}
	return nil
}
