package etelecom

import (
	"context"
	"time"

	"o.o/api/etelecom/call_direction"
	"o.o/api/etelecom/call_state"
	"o.o/api/etelecom/mobile_network"
	"o.o/api/meta"
	cm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/connection_type"
	"o.o/api/top/types/etc/payment_method"
	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
	"o.o/common/xerrors"
)

// +gen:api

type Aggregate interface {
	CreateHotline(context.Context, *CreateHotlineArgs) (*Hotline, error)
	UpdateHotlineInfo(context.Context, *UpdateHotlineInfoArgs) error
	DeleteHotline(ctx context.Context, id dot.ID) error
	RemoveHotlineOutOfTenant(context.Context, *RemoveHotlineOutOfTenantArgs) error
	ActiveHotlineForTenant(context.Context, *ActiveHotlineForTenantArgs) error
	CreateExtension(context.Context, *CreateExtensionArgs) (*Extension, error)
	CreateExtensionBySubscription(context.Context, *CreateExtenstionBySubscriptionArgs) (*Extension, error)
	ExtendExtension(context.Context, *ExtendExtensionArgs) (*Extension, error)
	DeleteExtension(ctx context.Context, id dot.ID) error
	RemoveUserOfExtension(context.Context, *RemoveUserOfExtensionArgs) (int, error)
	UpdateExternalExtensionInfo(context.Context, *UpdateExternalExtensionInfoArgs) error
	AssignUserToExtension(context.Context, *AssignUserToExtensionArgs) error

	// use to import extension with expiresAt
	ImportExtensions(context.Context, *ImportExtensionsArgs) error

	UpdateCallLogPostage(context.Context, *UpdateCallLogPostageArgs) error
	CreateOrUpdateCallLogFromCDR(context.Context, *CreateOrUpdateCallLogFromCDRArgs) (*CallLog, error)
	CreateCallLog(context.Context, *CreateCallLogArgs) (*CallLog, error)

	CreateTenant(context.Context, *CreateTenantArgs) (*Tenant, error)
	DeleteTenant(ctx context.Context, id dot.ID) error
	// active tenant and assign hotline to tenant
	ActivateTenant(context.Context, *ActivateTenantArgs) (*Tenant, error)

	DestroyCallSession(context.Context, *DestroyCallSessionArgs) error
}

type QueryService interface {
	GetHotline(context.Context, *GetHotlineArgs) (*Hotline, error)
	ListHotlines(context.Context, *ListHotlinesArgs) (*ListHotlinesReponse, error)
	ListBuiltinHotlines(context.Context, *cm.Empty) ([]*Hotline, error)
	GetHotlineByHotlineNumber(context.Context, *GetHotlineByHotlineNumberArgs) (*Hotline, error)

	GetExtension(context.Context, *GetExtensionArgs) (*Extension, error)
	ListExtensions(context.Context, *ListExtensionsArgs) ([]*Extension, error)

	// generate extension number => then use it for create external extension
	GetPrivateExtensionNumber(context.Context, *cm.Empty) (extensionNumber string, _ error)

	GetCallLogByExternalID(context.Context, *GetCallLogByExternalIDArgs) (*CallLog, error)
	ListCallLogs(context.Context, *ListCallLogsArgs) (*ListCallLogsResponse, error)
	GetCallLog(ctx context.Context, ID dot.ID) (*CallLog, error)
	GetCallLogByCallee(context.Context, *GetCallLogByCalleeArgs) (*CallLog, error)

	ListTenants(context.Context, *ListTenantsArgs) (*ListTenantsResponse, error)
	GetTenantByID(ctx context.Context, ID dot.ID) (*Tenant, error)
	GetTenantByConnection(context.Context, *GetTenantByConnectionArgs) (*Tenant, error)
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
	UserID       dot.ID
	ConnectionID dot.ID
	CallTargets  []*CallTarget
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
	if args.ExternalSessionID == "" {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing external session ID")
	}
	return nil
}

type CallTarget struct {
	AddTime      time.Time
	AnsweredTime time.Time
	EndReason    string
	EndedTime    time.Time
	FailCode     int
	RingDuration int
	RingTime     time.Time
	Status       string
	TargetNumber string
	TrunkName    string
}

type CreateCallLogArgs struct {
	ExternalSessionID string
	Direction         call_direction.CallDirection
	Caller            string
	Callee            string
	ExtensionID       dot.ID
	AccountID         dot.ID
	OwnerID           dot.ID
	ContactID         dot.ID
	CallState         call_state.CallState
	StartedAt         time.Time
	EndedAt           time.Time
	Note              string
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
	if args.StartedAt.IsZero() != args.EndedAt.IsZero() {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Please provide both started_at and ended_at")
	}
	if args.StartedAt.After(args.EndedAt) {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Invalid started_at, ended_at")
	}
	return nil
}

type GetHotlineArgs struct {
	ID      dot.ID
	OwnerID dot.ID
}

type ListHotlinesArgs struct {
	OwnerID      dot.ID
	TenantID     dot.ID
	ConnectionID dot.ID
	Paging       meta.Paging
}

type ListHotlinesReponse struct {
	Hotlines []*Hotline
	Paging   meta.PageInfo
}

type GetExtensionArgs struct {
	ID              dot.ID
	UserID          dot.ID
	AccountID       dot.ID
	HotlineID       dot.ID
	SubscriptionID  dot.ID
	ExtensionNumber string
}

type ListExtensionsArgs struct {
	AccountIDs       []dot.ID
	HotlineIDs       []dot.ID
	TenantID         dot.ID
	ExtensionNumbers []string
}

type UpdateExternalExtensionInfoArgs struct {
	ID                dot.ID
	HotlineID         dot.ID
	ExternalID        string
	ExtensionNumber   string
	ExtensionPassword string
	TenantDomain      string
	TenantID          dot.ID
}

type GetCallLogByExternalIDArgs struct {
	ExternalID string
}

type ListCallLogsArgs struct {
	HotlineIDs     []dot.ID
	ExtensionIDs   []dot.ID
	UserID         dot.ID
	OwnerID        dot.ID
	AccountID      dot.ID
	CallerOrCallee string
	CallState      call_state.CallState
	DateFrom       time.Time
	DateTo         time.Time
	Direction      call_direction.CallDirection
	Paging         meta.Paging
}

type GetCallLogByCalleeArgs struct {
	HotlineIDs []dot.ID
	Callee     string
	Direction  call_direction.CallDirection
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
	return nil
}

type UpdateHotlineInfoArgs struct {
	ID               dot.ID
	IsFreeCharge     dot.NullBool
	Name             string
	Description      string
	Status           status3.NullStatus
	TenantID         dot.ID
	ConnectionID     dot.ID
	ConnectionMethod connection_type.ConnectionMethod
	OwnerID          dot.ID
	Network          mobile_network.MobileNetwork
}

type CreateExtenstionBySubscriptionArgs struct {
	SubscriptionID     dot.ID
	SubscriptionPlanID dot.ID
	PaymentMethod      payment_method.PaymentMethod
	AccountID          dot.ID
	ExtensionNumber    string
	UserID             dot.ID
	HotlineID          dot.ID
	OwnerID            dot.ID
}

func (args *CreateExtenstionBySubscriptionArgs) Validate() error {
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

type CreateTenantArgs struct {
	OwnerID      dot.ID
	AccountID    dot.ID
	ConnectionID dot.ID
}

type GetTenantArgs struct {
	ID           dot.ID
	OwnerID      dot.ID
	ConnectionID dot.ID
}

type ActivateTenantArgs struct {
	OwnerID            dot.ID
	TenantID           dot.ID
	HotlineID          dot.ID
	ConnectionID       dot.ID
	ConnectionProvider connection_type.ConnectionProvider
}

type DestroyCallSessionArgs struct {
	ExternalSessionID string
	OwnerID           dot.ID
}

func (r *DestroyCallSessionArgs) Validate() error {
	if r.ExternalSessionID == "" {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing external session ID")
	}
	if r.OwnerID == 0 {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing owner ID")
	}
	return nil
}

type ListTenantsArgs struct {
	OwnerID      dot.ID
	ConnectionID dot.ID
	Paging       meta.Paging
}

type ListTenantsResponse struct {
	Tenants []*Tenant
	Paging  meta.PageInfo
}

type RemoveUserOfExtensionArgs struct {
	AccountID   dot.ID
	ExtensionID dot.ID
	UserID      dot.ID
}

type AssignUserToExtensionArgs struct {
	AccountID   dot.ID
	UserID      dot.ID
	ExtensionID dot.ID
}

type GetTenantByConnectionArgs struct {
	OwnerID      dot.ID
	ConnectionID dot.ID
}

type GetHotlineByHotlineNumberArgs struct {
	Hotline string
	OwnerID dot.ID
}

type ImportExtensionsArgs struct {
	ImportExtensions []*ImportExtension
}

type ImportExtension struct {
	TenantID        dot.ID
	OwnerID         dot.ID
	AccountID       dot.ID
	ExtensionNumber string
	ExpiresAt       time.Time
	HotlineID       dot.ID
}

type RemoveHotlineOutOfTenantArgs struct {
	HotlineID dot.ID
	OwnerID   dot.ID
}

func (args *RemoveHotlineOutOfTenantArgs) Validate() error {
	if args.HotlineID == 0 {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing hotline ID")
	}
	if args.OwnerID == 0 {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing owner ID")
	}
	return nil
}

type ActiveHotlineForTenantArgs struct {
	HotlineID dot.ID
	OwnerID   dot.ID
	TenantID  dot.ID
}

type ListCallLogsExportArgs struct {
	DateFrom     time.Time
	DateTo       time.Time
	ExtensionIDs []string
	OwnerID      dot.ID
	UserID       dot.ID
	AccountID    dot.ID
}
