// +build !generator

// Code generated by generator sqlgen. DO NOT EDIT.

package sqlstore

import (
	time "time"

	sq "o.o/backend/pkg/common/sql/sq"
	dot "o.o/capi/dot"
)

type VhtCallHistoryFilters struct{ prefix string }

func NewVhtCallHistoryFilters(prefix string) VhtCallHistoryFilters {
	return VhtCallHistoryFilters{prefix}
}

func (ft *VhtCallHistoryFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft VhtCallHistoryFilters) Prefix() string {
	return ft.prefix
}

func (ft *VhtCallHistoryFilters) ByCdrID(CdrID string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "cdr_id",
		Value:  CdrID,
		IsNil:  CdrID == "",
	}
}

func (ft *VhtCallHistoryFilters) ByCdrIDPtr(CdrID *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "cdr_id",
		Value:  CdrID,
		IsNil:  CdrID == nil,
		IsZero: CdrID != nil && (*CdrID) == "",
	}
}

func (ft *VhtCallHistoryFilters) ByCallID(CallID string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "call_id",
		Value:  CallID,
		IsNil:  CallID == "",
	}
}

func (ft *VhtCallHistoryFilters) ByCallIDPtr(CallID *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "call_id",
		Value:  CallID,
		IsNil:  CallID == nil,
		IsZero: CallID != nil && (*CallID) == "",
	}
}

func (ft *VhtCallHistoryFilters) BySipCallID(SipCallID string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "sip_call_id",
		Value:  SipCallID,
		IsNil:  SipCallID == "",
	}
}

func (ft *VhtCallHistoryFilters) BySipCallIDPtr(SipCallID *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "sip_call_id",
		Value:  SipCallID,
		IsNil:  SipCallID == nil,
		IsZero: SipCallID != nil && (*SipCallID) == "",
	}
}

func (ft *VhtCallHistoryFilters) BySdkCallID(SdkCallID string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "sdk_call_id",
		Value:  SdkCallID,
		IsNil:  SdkCallID == "",
	}
}

func (ft *VhtCallHistoryFilters) BySdkCallIDPtr(SdkCallID *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "sdk_call_id",
		Value:  SdkCallID,
		IsNil:  SdkCallID == nil,
		IsZero: SdkCallID != nil && (*SdkCallID) == "",
	}
}

func (ft *VhtCallHistoryFilters) ByCause(Cause string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "cause",
		Value:  Cause,
		IsNil:  Cause == "",
	}
}

func (ft *VhtCallHistoryFilters) ByCausePtr(Cause *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "cause",
		Value:  Cause,
		IsNil:  Cause == nil,
		IsZero: Cause != nil && (*Cause) == "",
	}
}

func (ft *VhtCallHistoryFilters) ByQ850Cause(Q850Cause string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "q850_cause",
		Value:  Q850Cause,
		IsNil:  Q850Cause == "",
	}
}

func (ft *VhtCallHistoryFilters) ByQ850CausePtr(Q850Cause *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "q850_cause",
		Value:  Q850Cause,
		IsNil:  Q850Cause == nil,
		IsZero: Q850Cause != nil && (*Q850Cause) == "",
	}
}

func (ft *VhtCallHistoryFilters) ByFromExtension(FromExtension string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "from_extension",
		Value:  FromExtension,
		IsNil:  FromExtension == "",
	}
}

func (ft *VhtCallHistoryFilters) ByFromExtensionPtr(FromExtension *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "from_extension",
		Value:  FromExtension,
		IsNil:  FromExtension == nil,
		IsZero: FromExtension != nil && (*FromExtension) == "",
	}
}

func (ft *VhtCallHistoryFilters) ByToExtension(ToExtension string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "to_extension",
		Value:  ToExtension,
		IsNil:  ToExtension == "",
	}
}

func (ft *VhtCallHistoryFilters) ByToExtensionPtr(ToExtension *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "to_extension",
		Value:  ToExtension,
		IsNil:  ToExtension == nil,
		IsZero: ToExtension != nil && (*ToExtension) == "",
	}
}

func (ft *VhtCallHistoryFilters) ByFromNumber(FromNumber string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "from_number",
		Value:  FromNumber,
		IsNil:  FromNumber == "",
	}
}

func (ft *VhtCallHistoryFilters) ByFromNumberPtr(FromNumber *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "from_number",
		Value:  FromNumber,
		IsNil:  FromNumber == nil,
		IsZero: FromNumber != nil && (*FromNumber) == "",
	}
}

func (ft *VhtCallHistoryFilters) ByToNumber(ToNumber string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "to_number",
		Value:  ToNumber,
		IsNil:  ToNumber == "",
	}
}

func (ft *VhtCallHistoryFilters) ByToNumberPtr(ToNumber *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "to_number",
		Value:  ToNumber,
		IsNil:  ToNumber == nil,
		IsZero: ToNumber != nil && (*ToNumber) == "",
	}
}

func (ft *VhtCallHistoryFilters) ByDuration(Duration int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "duration",
		Value:  Duration,
		IsNil:  Duration == 0,
	}
}

func (ft *VhtCallHistoryFilters) ByDurationPtr(Duration *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "duration",
		Value:  Duration,
		IsNil:  Duration == nil,
		IsZero: Duration != nil && (*Duration) == 0,
	}
}

func (ft *VhtCallHistoryFilters) ByDirection(Direction int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "direction",
		Value:  Direction,
		IsNil:  Direction == 0,
	}
}

func (ft *VhtCallHistoryFilters) ByDirectionPtr(Direction *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "direction",
		Value:  Direction,
		IsNil:  Direction == nil,
		IsZero: Direction != nil && (*Direction) == 0,
	}
}

func (ft *VhtCallHistoryFilters) ByTimeStarted(TimeStarted time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "time_started",
		Value:  TimeStarted,
		IsNil:  TimeStarted.IsZero(),
	}
}

func (ft *VhtCallHistoryFilters) ByTimeStartedPtr(TimeStarted *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "time_started",
		Value:  TimeStarted,
		IsNil:  TimeStarted == nil,
		IsZero: TimeStarted != nil && (*TimeStarted).IsZero(),
	}
}

func (ft *VhtCallHistoryFilters) ByTimeConnected(TimeConnected time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "time_connected",
		Value:  TimeConnected,
		IsNil:  TimeConnected.IsZero(),
	}
}

func (ft *VhtCallHistoryFilters) ByTimeConnectedPtr(TimeConnected *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "time_connected",
		Value:  TimeConnected,
		IsNil:  TimeConnected == nil,
		IsZero: TimeConnected != nil && (*TimeConnected).IsZero(),
	}
}

func (ft *VhtCallHistoryFilters) ByTimeEnded(TimeEnded time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "time_ended",
		Value:  TimeEnded,
		IsNil:  TimeEnded.IsZero(),
	}
}

func (ft *VhtCallHistoryFilters) ByTimeEndedPtr(TimeEnded *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "time_ended",
		Value:  TimeEnded,
		IsNil:  TimeEnded == nil,
		IsZero: TimeEnded != nil && (*TimeEnded).IsZero(),
	}
}

func (ft *VhtCallHistoryFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *VhtCallHistoryFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *VhtCallHistoryFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *VhtCallHistoryFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}

func (ft *VhtCallHistoryFilters) ByRecordingPath(RecordingPath string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "recording_path",
		Value:  RecordingPath,
		IsNil:  RecordingPath == "",
	}
}

func (ft *VhtCallHistoryFilters) ByRecordingPathPtr(RecordingPath *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "recording_path",
		Value:  RecordingPath,
		IsNil:  RecordingPath == nil,
		IsZero: RecordingPath != nil && (*RecordingPath) == "",
	}
}

func (ft *VhtCallHistoryFilters) ByRecordingURL(RecordingURL string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "recording_url",
		Value:  RecordingURL,
		IsNil:  RecordingURL == "",
	}
}

func (ft *VhtCallHistoryFilters) ByRecordingURLPtr(RecordingURL *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "recording_url",
		Value:  RecordingURL,
		IsNil:  RecordingURL == nil,
		IsZero: RecordingURL != nil && (*RecordingURL) == "",
	}
}

func (ft *VhtCallHistoryFilters) ByRecordFileSize(RecordFileSize int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "record_file_size",
		Value:  RecordFileSize,
		IsNil:  RecordFileSize == 0,
	}
}

func (ft *VhtCallHistoryFilters) ByRecordFileSizePtr(RecordFileSize *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "record_file_size",
		Value:  RecordFileSize,
		IsNil:  RecordFileSize == nil,
		IsZero: RecordFileSize != nil && (*RecordFileSize) == 0,
	}
}

func (ft *VhtCallHistoryFilters) ByEtopAccountID(EtopAccountID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "etop_account_id",
		Value:  EtopAccountID,
		IsNil:  EtopAccountID == 0,
	}
}

func (ft *VhtCallHistoryFilters) ByEtopAccountIDPtr(EtopAccountID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "etop_account_id",
		Value:  EtopAccountID,
		IsNil:  EtopAccountID == nil,
		IsZero: EtopAccountID != nil && (*EtopAccountID) == 0,
	}
}

func (ft *VhtCallHistoryFilters) ByVtigerAccountID(VtigerAccountID string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "vtiger_account_id",
		Value:  VtigerAccountID,
		IsNil:  VtigerAccountID == "",
	}
}

func (ft *VhtCallHistoryFilters) ByVtigerAccountIDPtr(VtigerAccountID *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "vtiger_account_id",
		Value:  VtigerAccountID,
		IsNil:  VtigerAccountID == nil,
		IsZero: VtigerAccountID != nil && (*VtigerAccountID) == "",
	}
}

func (ft *VhtCallHistoryFilters) BySyncStatus(SyncStatus string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "sync_status",
		Value:  SyncStatus,
		IsNil:  SyncStatus == "",
	}
}

func (ft *VhtCallHistoryFilters) BySyncStatusPtr(SyncStatus *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "sync_status",
		Value:  SyncStatus,
		IsNil:  SyncStatus == nil,
		IsZero: SyncStatus != nil && (*SyncStatus) == "",
	}
}

func (ft *VhtCallHistoryFilters) ByOData(OData string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "o_data",
		Value:  OData,
		IsNil:  OData == "",
	}
}

func (ft *VhtCallHistoryFilters) ByODataPtr(OData *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "o_data",
		Value:  OData,
		IsNil:  OData == nil,
		IsZero: OData != nil && (*OData) == "",
	}
}

func (ft *VhtCallHistoryFilters) BySearchNorm(SearchNorm string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "search_norm",
		Value:  SearchNorm,
		IsNil:  SearchNorm == "",
	}
}

func (ft *VhtCallHistoryFilters) BySearchNormPtr(SearchNorm *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "search_norm",
		Value:  SearchNorm,
		IsNil:  SearchNorm == nil,
		IsZero: SearchNorm != nil && (*SearchNorm) == "",
	}
}
