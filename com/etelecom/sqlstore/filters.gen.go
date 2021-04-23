// +build !generator

// Code generated by generator sqlgen. DO NOT EDIT.

package sqlstore

import (
	time "time"

	call_direction "o.o/api/etelecom/call_direction"
	call_state "o.o/api/etelecom/call_state"
	connection_type "o.o/api/top/types/etc/connection_type"
	status3 "o.o/api/top/types/etc/status3"
	status5 "o.o/api/top/types/etc/status5"
	sq "o.o/backend/pkg/common/sql/sq"
	dot "o.o/capi/dot"
)

type CallLogFilters struct{ prefix string }

func NewCallLogFilters(prefix string) CallLogFilters {
	return CallLogFilters{prefix}
}

func (ft *CallLogFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft CallLogFilters) Prefix() string {
	return ft.prefix
}

func (ft *CallLogFilters) ByID(ID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *CallLogFilters) ByIDPtr(ID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *CallLogFilters) ByExternalID(ExternalID string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "external_id",
		Value:  ExternalID,
		IsNil:  ExternalID == "",
	}
}

func (ft *CallLogFilters) ByExternalIDPtr(ExternalID *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "external_id",
		Value:  ExternalID,
		IsNil:  ExternalID == nil,
		IsZero: ExternalID != nil && (*ExternalID) == "",
	}
}

func (ft *CallLogFilters) ByAccountID(AccountID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "account_id",
		Value:  AccountID,
		IsNil:  AccountID == 0,
	}
}

func (ft *CallLogFilters) ByAccountIDPtr(AccountID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "account_id",
		Value:  AccountID,
		IsNil:  AccountID == nil,
		IsZero: AccountID != nil && (*AccountID) == 0,
	}
}

func (ft *CallLogFilters) ByOwnerID(OwnerID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "owner_id",
		Value:  OwnerID,
		IsNil:  OwnerID == 0,
	}
}

func (ft *CallLogFilters) ByOwnerIDPtr(OwnerID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "owner_id",
		Value:  OwnerID,
		IsNil:  OwnerID == nil,
		IsZero: OwnerID != nil && (*OwnerID) == 0,
	}
}

func (ft *CallLogFilters) ByUserID(UserID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "user_id",
		Value:  UserID,
		IsNil:  UserID == 0,
	}
}

func (ft *CallLogFilters) ByUserIDPtr(UserID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "user_id",
		Value:  UserID,
		IsNil:  UserID == nil,
		IsZero: UserID != nil && (*UserID) == 0,
	}
}

func (ft *CallLogFilters) ByStartedAt(StartedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "started_at",
		Value:  StartedAt,
		IsNil:  StartedAt.IsZero(),
	}
}

func (ft *CallLogFilters) ByStartedAtPtr(StartedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "started_at",
		Value:  StartedAt,
		IsNil:  StartedAt == nil,
		IsZero: StartedAt != nil && (*StartedAt).IsZero(),
	}
}

func (ft *CallLogFilters) ByEndedAt(EndedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "ended_at",
		Value:  EndedAt,
		IsNil:  EndedAt.IsZero(),
	}
}

func (ft *CallLogFilters) ByEndedAtPtr(EndedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "ended_at",
		Value:  EndedAt,
		IsNil:  EndedAt == nil,
		IsZero: EndedAt != nil && (*EndedAt).IsZero(),
	}
}

func (ft *CallLogFilters) ByDuration(Duration int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "duration",
		Value:  Duration,
		IsNil:  Duration == 0,
	}
}

func (ft *CallLogFilters) ByDurationPtr(Duration *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "duration",
		Value:  Duration,
		IsNil:  Duration == nil,
		IsZero: Duration != nil && (*Duration) == 0,
	}
}

func (ft *CallLogFilters) ByCaller(Caller string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "caller",
		Value:  Caller,
		IsNil:  Caller == "",
	}
}

func (ft *CallLogFilters) ByCallerPtr(Caller *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "caller",
		Value:  Caller,
		IsNil:  Caller == nil,
		IsZero: Caller != nil && (*Caller) == "",
	}
}

func (ft *CallLogFilters) ByCallee(Callee string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "callee",
		Value:  Callee,
		IsNil:  Callee == "",
	}
}

func (ft *CallLogFilters) ByCalleePtr(Callee *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "callee",
		Value:  Callee,
		IsNil:  Callee == nil,
		IsZero: Callee != nil && (*Callee) == "",
	}
}

func (ft *CallLogFilters) ByExternalDirection(ExternalDirection string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "external_direction",
		Value:  ExternalDirection,
		IsNil:  ExternalDirection == "",
	}
}

func (ft *CallLogFilters) ByExternalDirectionPtr(ExternalDirection *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "external_direction",
		Value:  ExternalDirection,
		IsNil:  ExternalDirection == nil,
		IsZero: ExternalDirection != nil && (*ExternalDirection) == "",
	}
}

func (ft *CallLogFilters) ByDirection(Direction call_direction.CallDirection) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "direction",
		Value:  Direction,
		IsNil:  Direction == 0,
	}
}

func (ft *CallLogFilters) ByDirectionPtr(Direction *call_direction.CallDirection) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "direction",
		Value:  Direction,
		IsNil:  Direction == nil,
		IsZero: Direction != nil && (*Direction) == 0,
	}
}

func (ft *CallLogFilters) ByExtensionID(ExtensionID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "extension_id",
		Value:  ExtensionID,
		IsNil:  ExtensionID == 0,
	}
}

func (ft *CallLogFilters) ByExtensionIDPtr(ExtensionID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "extension_id",
		Value:  ExtensionID,
		IsNil:  ExtensionID == nil,
		IsZero: ExtensionID != nil && (*ExtensionID) == 0,
	}
}

func (ft *CallLogFilters) ByHotlineID(HotlineID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "hotline_id",
		Value:  HotlineID,
		IsNil:  HotlineID == 0,
	}
}

func (ft *CallLogFilters) ByHotlineIDPtr(HotlineID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "hotline_id",
		Value:  HotlineID,
		IsNil:  HotlineID == nil,
		IsZero: HotlineID != nil && (*HotlineID) == 0,
	}
}

func (ft *CallLogFilters) ByExternalCallStatus(ExternalCallStatus string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "external_call_status",
		Value:  ExternalCallStatus,
		IsNil:  ExternalCallStatus == "",
	}
}

func (ft *CallLogFilters) ByExternalCallStatusPtr(ExternalCallStatus *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "external_call_status",
		Value:  ExternalCallStatus,
		IsNil:  ExternalCallStatus == nil,
		IsZero: ExternalCallStatus != nil && (*ExternalCallStatus) == "",
	}
}

func (ft *CallLogFilters) ByContactID(ContactID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "contact_id",
		Value:  ContactID,
		IsNil:  ContactID == 0,
	}
}

func (ft *CallLogFilters) ByContactIDPtr(ContactID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "contact_id",
		Value:  ContactID,
		IsNil:  ContactID == nil,
		IsZero: ContactID != nil && (*ContactID) == 0,
	}
}

func (ft *CallLogFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *CallLogFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *CallLogFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *CallLogFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}

func (ft *CallLogFilters) ByCallState(CallState call_state.CallState) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "call_state",
		Value:  CallState,
		IsNil:  CallState == 0,
	}
}

func (ft *CallLogFilters) ByCallStatePtr(CallState *call_state.CallState) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "call_state",
		Value:  CallState,
		IsNil:  CallState == nil,
		IsZero: CallState != nil && (*CallState) == 0,
	}
}

func (ft *CallLogFilters) ByCallStatus(CallStatus status5.Status) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "call_status",
		Value:  CallStatus,
		IsNil:  CallStatus == 0,
	}
}

func (ft *CallLogFilters) ByCallStatusPtr(CallStatus *status5.Status) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "call_status",
		Value:  CallStatus,
		IsNil:  CallStatus == nil,
		IsZero: CallStatus != nil && (*CallStatus) == 0,
	}
}

func (ft *CallLogFilters) ByDurationPostage(DurationPostage int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "duration_postage",
		Value:  DurationPostage,
		IsNil:  DurationPostage == 0,
	}
}

func (ft *CallLogFilters) ByDurationPostagePtr(DurationPostage *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "duration_postage",
		Value:  DurationPostage,
		IsNil:  DurationPostage == nil,
		IsZero: DurationPostage != nil && (*DurationPostage) == 0,
	}
}

func (ft *CallLogFilters) ByPostage(Postage int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "postage",
		Value:  Postage,
		IsNil:  Postage == 0,
	}
}

func (ft *CallLogFilters) ByPostagePtr(Postage *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "postage",
		Value:  Postage,
		IsNil:  Postage == nil,
		IsZero: Postage != nil && (*Postage) == 0,
	}
}

func (ft *CallLogFilters) ByExternalSessionID(ExternalSessionID string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "external_session_id",
		Value:  ExternalSessionID,
		IsNil:  ExternalSessionID == "",
	}
}

func (ft *CallLogFilters) ByExternalSessionIDPtr(ExternalSessionID *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "external_session_id",
		Value:  ExternalSessionID,
		IsNil:  ExternalSessionID == nil,
		IsZero: ExternalSessionID != nil && (*ExternalSessionID) == "",
	}
}

type ExtensionFilters struct{ prefix string }

func NewExtensionFilters(prefix string) ExtensionFilters {
	return ExtensionFilters{prefix}
}

func (ft *ExtensionFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft ExtensionFilters) Prefix() string {
	return ft.prefix
}

func (ft *ExtensionFilters) ByID(ID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *ExtensionFilters) ByIDPtr(ID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *ExtensionFilters) ByUserID(UserID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "user_id",
		Value:  UserID,
		IsNil:  UserID == 0,
	}
}

func (ft *ExtensionFilters) ByUserIDPtr(UserID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "user_id",
		Value:  UserID,
		IsNil:  UserID == nil,
		IsZero: UserID != nil && (*UserID) == 0,
	}
}

func (ft *ExtensionFilters) ByAccountID(AccountID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "account_id",
		Value:  AccountID,
		IsNil:  AccountID == 0,
	}
}

func (ft *ExtensionFilters) ByAccountIDPtr(AccountID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "account_id",
		Value:  AccountID,
		IsNil:  AccountID == nil,
		IsZero: AccountID != nil && (*AccountID) == 0,
	}
}

func (ft *ExtensionFilters) ByHotlineID(HotlineID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "hotline_id",
		Value:  HotlineID,
		IsNil:  HotlineID == 0,
	}
}

func (ft *ExtensionFilters) ByHotlineIDPtr(HotlineID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "hotline_id",
		Value:  HotlineID,
		IsNil:  HotlineID == nil,
		IsZero: HotlineID != nil && (*HotlineID) == 0,
	}
}

func (ft *ExtensionFilters) ByExtensionNumber(ExtensionNumber string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "extension_number",
		Value:  ExtensionNumber,
		IsNil:  ExtensionNumber == "",
	}
}

func (ft *ExtensionFilters) ByExtensionNumberPtr(ExtensionNumber *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "extension_number",
		Value:  ExtensionNumber,
		IsNil:  ExtensionNumber == nil,
		IsZero: ExtensionNumber != nil && (*ExtensionNumber) == "",
	}
}

func (ft *ExtensionFilters) ByExtensionPassword(ExtensionPassword string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "extension_password",
		Value:  ExtensionPassword,
		IsNil:  ExtensionPassword == "",
	}
}

func (ft *ExtensionFilters) ByExtensionPasswordPtr(ExtensionPassword *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "extension_password",
		Value:  ExtensionPassword,
		IsNil:  ExtensionPassword == nil,
		IsZero: ExtensionPassword != nil && (*ExtensionPassword) == "",
	}
}

func (ft *ExtensionFilters) ByTenantDomain(TenantDomain string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "tenant_domain",
		Value:  TenantDomain,
		IsNil:  TenantDomain == "",
	}
}

func (ft *ExtensionFilters) ByTenantDomainPtr(TenantDomain *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "tenant_domain",
		Value:  TenantDomain,
		IsNil:  TenantDomain == nil,
		IsZero: TenantDomain != nil && (*TenantDomain) == "",
	}
}

func (ft *ExtensionFilters) ByTenantID(TenantID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "tenant_id",
		Value:  TenantID,
		IsNil:  TenantID == 0,
	}
}

func (ft *ExtensionFilters) ByTenantIDPtr(TenantID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "tenant_id",
		Value:  TenantID,
		IsNil:  TenantID == nil,
		IsZero: TenantID != nil && (*TenantID) == 0,
	}
}

func (ft *ExtensionFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *ExtensionFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *ExtensionFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *ExtensionFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}

func (ft *ExtensionFilters) ByDeletedAt(DeletedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "deleted_at",
		Value:  DeletedAt,
		IsNil:  DeletedAt.IsZero(),
	}
}

func (ft *ExtensionFilters) ByDeletedAtPtr(DeletedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "deleted_at",
		Value:  DeletedAt,
		IsNil:  DeletedAt == nil,
		IsZero: DeletedAt != nil && (*DeletedAt).IsZero(),
	}
}

func (ft *ExtensionFilters) BySubscriptionID(SubscriptionID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "subscription_id",
		Value:  SubscriptionID,
		IsNil:  SubscriptionID == 0,
	}
}

func (ft *ExtensionFilters) BySubscriptionIDPtr(SubscriptionID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "subscription_id",
		Value:  SubscriptionID,
		IsNil:  SubscriptionID == nil,
		IsZero: SubscriptionID != nil && (*SubscriptionID) == 0,
	}
}

func (ft *ExtensionFilters) ByExpiresAt(ExpiresAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "expires_at",
		Value:  ExpiresAt,
		IsNil:  ExpiresAt.IsZero(),
	}
}

func (ft *ExtensionFilters) ByExpiresAtPtr(ExpiresAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "expires_at",
		Value:  ExpiresAt,
		IsNil:  ExpiresAt == nil,
		IsZero: ExpiresAt != nil && (*ExpiresAt).IsZero(),
	}
}

type HotlineFilters struct{ prefix string }

func NewHotlineFilters(prefix string) HotlineFilters {
	return HotlineFilters{prefix}
}

func (ft *HotlineFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft HotlineFilters) Prefix() string {
	return ft.prefix
}

func (ft *HotlineFilters) ByID(ID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *HotlineFilters) ByIDPtr(ID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *HotlineFilters) ByOwnerID(OwnerID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "owner_id",
		Value:  OwnerID,
		IsNil:  OwnerID == 0,
	}
}

func (ft *HotlineFilters) ByOwnerIDPtr(OwnerID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "owner_id",
		Value:  OwnerID,
		IsNil:  OwnerID == nil,
		IsZero: OwnerID != nil && (*OwnerID) == 0,
	}
}

func (ft *HotlineFilters) ByName(Name string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "name",
		Value:  Name,
		IsNil:  Name == "",
	}
}

func (ft *HotlineFilters) ByNamePtr(Name *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "name",
		Value:  Name,
		IsNil:  Name == nil,
		IsZero: Name != nil && (*Name) == "",
	}
}

func (ft *HotlineFilters) ByHotline(Hotline string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "hotline",
		Value:  Hotline,
		IsNil:  Hotline == "",
	}
}

func (ft *HotlineFilters) ByHotlinePtr(Hotline *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "hotline",
		Value:  Hotline,
		IsNil:  Hotline == nil,
		IsZero: Hotline != nil && (*Hotline) == "",
	}
}

func (ft *HotlineFilters) ByNetwork(Network string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "network",
		Value:  Network,
		IsNil:  Network == "",
	}
}

func (ft *HotlineFilters) ByNetworkPtr(Network *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "network",
		Value:  Network,
		IsNil:  Network == nil,
		IsZero: Network != nil && (*Network) == "",
	}
}

func (ft *HotlineFilters) ByConnectionID(ConnectionID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "connection_id",
		Value:  ConnectionID,
		IsNil:  ConnectionID == 0,
	}
}

func (ft *HotlineFilters) ByConnectionIDPtr(ConnectionID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "connection_id",
		Value:  ConnectionID,
		IsNil:  ConnectionID == nil,
		IsZero: ConnectionID != nil && (*ConnectionID) == 0,
	}
}

func (ft *HotlineFilters) ByConnectionMethod(ConnectionMethod connection_type.ConnectionMethod) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "connection_method",
		Value:  ConnectionMethod,
		IsNil:  ConnectionMethod == 0,
	}
}

func (ft *HotlineFilters) ByConnectionMethodPtr(ConnectionMethod *connection_type.ConnectionMethod) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "connection_method",
		Value:  ConnectionMethod,
		IsNil:  ConnectionMethod == nil,
		IsZero: ConnectionMethod != nil && (*ConnectionMethod) == 0,
	}
}

func (ft *HotlineFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *HotlineFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *HotlineFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *HotlineFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}

func (ft *HotlineFilters) ByDeletedAt(DeletedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "deleted_at",
		Value:  DeletedAt,
		IsNil:  DeletedAt.IsZero(),
	}
}

func (ft *HotlineFilters) ByDeletedAtPtr(DeletedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "deleted_at",
		Value:  DeletedAt,
		IsNil:  DeletedAt == nil,
		IsZero: DeletedAt != nil && (*DeletedAt).IsZero(),
	}
}

func (ft *HotlineFilters) ByStatus(Status status3.Status) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == 0,
	}
}

func (ft *HotlineFilters) ByStatusPtr(Status *status3.Status) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == nil,
		IsZero: Status != nil && (*Status) == 0,
	}
}

func (ft *HotlineFilters) ByDescription(Description string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "description",
		Value:  Description,
		IsNil:  Description == "",
	}
}

func (ft *HotlineFilters) ByDescriptionPtr(Description *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "description",
		Value:  Description,
		IsNil:  Description == nil,
		IsZero: Description != nil && (*Description) == "",
	}
}

func (ft *HotlineFilters) ByTenantID(TenantID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "tenant_id",
		Value:  TenantID,
		IsNil:  TenantID == 0,
	}
}

func (ft *HotlineFilters) ByTenantIDPtr(TenantID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "tenant_id",
		Value:  TenantID,
		IsNil:  TenantID == nil,
		IsZero: TenantID != nil && (*TenantID) == 0,
	}
}

type TenantFilters struct{ prefix string }

func NewTenantFilters(prefix string) TenantFilters {
	return TenantFilters{prefix}
}

func (ft *TenantFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft TenantFilters) Prefix() string {
	return ft.prefix
}

func (ft *TenantFilters) ByID(ID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *TenantFilters) ByIDPtr(ID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *TenantFilters) ByOwnerID(OwnerID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "owner_id",
		Value:  OwnerID,
		IsNil:  OwnerID == 0,
	}
}

func (ft *TenantFilters) ByOwnerIDPtr(OwnerID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "owner_id",
		Value:  OwnerID,
		IsNil:  OwnerID == nil,
		IsZero: OwnerID != nil && (*OwnerID) == 0,
	}
}

func (ft *TenantFilters) ByName(Name string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "name",
		Value:  Name,
		IsNil:  Name == "",
	}
}

func (ft *TenantFilters) ByNamePtr(Name *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "name",
		Value:  Name,
		IsNil:  Name == nil,
		IsZero: Name != nil && (*Name) == "",
	}
}

func (ft *TenantFilters) ByDomain(Domain string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "domain",
		Value:  Domain,
		IsNil:  Domain == "",
	}
}

func (ft *TenantFilters) ByDomainPtr(Domain *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "domain",
		Value:  Domain,
		IsNil:  Domain == nil,
		IsZero: Domain != nil && (*Domain) == "",
	}
}

func (ft *TenantFilters) ByPassword(Password string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "password",
		Value:  Password,
		IsNil:  Password == "",
	}
}

func (ft *TenantFilters) ByPasswordPtr(Password *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "password",
		Value:  Password,
		IsNil:  Password == nil,
		IsZero: Password != nil && (*Password) == "",
	}
}

func (ft *TenantFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *TenantFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *TenantFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *TenantFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}

func (ft *TenantFilters) ByDeletedAt(DeletedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "deleted_at",
		Value:  DeletedAt,
		IsNil:  DeletedAt.IsZero(),
	}
}

func (ft *TenantFilters) ByDeletedAtPtr(DeletedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "deleted_at",
		Value:  DeletedAt,
		IsNil:  DeletedAt == nil,
		IsZero: DeletedAt != nil && (*DeletedAt).IsZero(),
	}
}

func (ft *TenantFilters) ByConnectionID(ConnectionID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "connection_id",
		Value:  ConnectionID,
		IsNil:  ConnectionID == 0,
	}
}

func (ft *TenantFilters) ByConnectionIDPtr(ConnectionID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "connection_id",
		Value:  ConnectionID,
		IsNil:  ConnectionID == nil,
		IsZero: ConnectionID != nil && (*ConnectionID) == 0,
	}
}

func (ft *TenantFilters) ByConnectionMethod(ConnectionMethod connection_type.ConnectionMethod) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "connection_method",
		Value:  ConnectionMethod,
		IsNil:  ConnectionMethod == 0,
	}
}

func (ft *TenantFilters) ByConnectionMethodPtr(ConnectionMethod *connection_type.ConnectionMethod) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "connection_method",
		Value:  ConnectionMethod,
		IsNil:  ConnectionMethod == nil,
		IsZero: ConnectionMethod != nil && (*ConnectionMethod) == 0,
	}
}
