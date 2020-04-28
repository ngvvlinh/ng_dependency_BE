// +build !generator

// Code generated by generator sqlgen. DO NOT EDIT.

package sqlstore

import (
	time "time"

	status4 "o.o/api/top/types/etc/status4"
	sq "o.o/backend/pkg/common/sql/sq"
	etopmodel "o.o/backend/pkg/etop/model"
	dot "o.o/capi/dot"
)

type CallbackFilters struct{ prefix string }

func NewCallbackFilters(prefix string) CallbackFilters {
	return CallbackFilters{prefix}
}

func (ft *CallbackFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft CallbackFilters) Prefix() string {
	return ft.prefix
}

func (ft *CallbackFilters) ByID(ID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *CallbackFilters) ByIDPtr(ID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *CallbackFilters) ByWebhookID(WebhookID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "webhook_id",
		Value:  WebhookID,
		IsNil:  WebhookID == 0,
	}
}

func (ft *CallbackFilters) ByWebhookIDPtr(WebhookID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "webhook_id",
		Value:  WebhookID,
		IsNil:  WebhookID == nil,
		IsZero: WebhookID != nil && (*WebhookID) == 0,
	}
}

func (ft *CallbackFilters) ByAccountID(AccountID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "account_id",
		Value:  AccountID,
		IsNil:  AccountID == 0,
	}
}

func (ft *CallbackFilters) ByAccountIDPtr(AccountID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "account_id",
		Value:  AccountID,
		IsNil:  AccountID == nil,
		IsZero: AccountID != nil && (*AccountID) == 0,
	}
}

func (ft *CallbackFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *CallbackFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

type CodeFilters struct{ prefix string }

func NewCodeFilters(prefix string) CodeFilters {
	return CodeFilters{prefix}
}

func (ft *CodeFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft CodeFilters) Prefix() string {
	return ft.prefix
}

func (ft *CodeFilters) ByCode(Code string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "code",
		Value:  Code,
		IsNil:  Code == "",
	}
}

func (ft *CodeFilters) ByCodePtr(Code *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "code",
		Value:  Code,
		IsNil:  Code == nil,
		IsZero: Code != nil && (*Code) == "",
	}
}

func (ft *CodeFilters) ByType(Type etopmodel.CodeType) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "type",
		Value:  Type,
		IsNil:  Type == "",
	}
}

func (ft *CodeFilters) ByTypePtr(Type *etopmodel.CodeType) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "type",
		Value:  Type,
		IsNil:  Type == nil,
		IsZero: Type != nil && (*Type) == "",
	}
}

func (ft *CodeFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *CodeFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

type ExportAttemptFilters struct{ prefix string }

func NewExportAttemptFilters(prefix string) ExportAttemptFilters {
	return ExportAttemptFilters{prefix}
}

func (ft *ExportAttemptFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft ExportAttemptFilters) Prefix() string {
	return ft.prefix
}

func (ft *ExportAttemptFilters) ByID(ID string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == "",
	}
}

func (ft *ExportAttemptFilters) ByIDPtr(ID *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == "",
	}
}

func (ft *ExportAttemptFilters) ByUserID(UserID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "user_id",
		Value:  UserID,
		IsNil:  UserID == 0,
	}
}

func (ft *ExportAttemptFilters) ByUserIDPtr(UserID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "user_id",
		Value:  UserID,
		IsNil:  UserID == nil,
		IsZero: UserID != nil && (*UserID) == 0,
	}
}

func (ft *ExportAttemptFilters) ByAccountID(AccountID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "account_id",
		Value:  AccountID,
		IsNil:  AccountID == 0,
	}
}

func (ft *ExportAttemptFilters) ByAccountIDPtr(AccountID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "account_id",
		Value:  AccountID,
		IsNil:  AccountID == nil,
		IsZero: AccountID != nil && (*AccountID) == 0,
	}
}

func (ft *ExportAttemptFilters) ByExportType(ExportType string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "export_type",
		Value:  ExportType,
		IsNil:  ExportType == "",
	}
}

func (ft *ExportAttemptFilters) ByExportTypePtr(ExportType *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "export_type",
		Value:  ExportType,
		IsNil:  ExportType == nil,
		IsZero: ExportType != nil && (*ExportType) == "",
	}
}

func (ft *ExportAttemptFilters) ByFileName(FileName string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "filename",
		Value:  FileName,
		IsNil:  FileName == "",
	}
}

func (ft *ExportAttemptFilters) ByFileNamePtr(FileName *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "filename",
		Value:  FileName,
		IsNil:  FileName == nil,
		IsZero: FileName != nil && (*FileName) == "",
	}
}

func (ft *ExportAttemptFilters) ByStoredFile(StoredFile string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "stored_file",
		Value:  StoredFile,
		IsNil:  StoredFile == "",
	}
}

func (ft *ExportAttemptFilters) ByStoredFilePtr(StoredFile *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "stored_file",
		Value:  StoredFile,
		IsNil:  StoredFile == nil,
		IsZero: StoredFile != nil && (*StoredFile) == "",
	}
}

func (ft *ExportAttemptFilters) ByDownloadURL(DownloadURL string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "download_url",
		Value:  DownloadURL,
		IsNil:  DownloadURL == "",
	}
}

func (ft *ExportAttemptFilters) ByDownloadURLPtr(DownloadURL *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "download_url",
		Value:  DownloadURL,
		IsNil:  DownloadURL == nil,
		IsZero: DownloadURL != nil && (*DownloadURL) == "",
	}
}

func (ft *ExportAttemptFilters) ByRequestQuery(RequestQuery string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "request_query",
		Value:  RequestQuery,
		IsNil:  RequestQuery == "",
	}
}

func (ft *ExportAttemptFilters) ByRequestQueryPtr(RequestQuery *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "request_query",
		Value:  RequestQuery,
		IsNil:  RequestQuery == nil,
		IsZero: RequestQuery != nil && (*RequestQuery) == "",
	}
}

func (ft *ExportAttemptFilters) ByMimeType(MimeType string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "mime_type",
		Value:  MimeType,
		IsNil:  MimeType == "",
	}
}

func (ft *ExportAttemptFilters) ByMimeTypePtr(MimeType *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "mime_type",
		Value:  MimeType,
		IsNil:  MimeType == nil,
		IsZero: MimeType != nil && (*MimeType) == "",
	}
}

func (ft *ExportAttemptFilters) ByStatus(Status status4.Status) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == 0,
	}
}

func (ft *ExportAttemptFilters) ByStatusPtr(Status *status4.Status) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == nil,
		IsZero: Status != nil && (*Status) == 0,
	}
}

func (ft *ExportAttemptFilters) ByNTotal(NTotal int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "n_total",
		Value:  NTotal,
		IsNil:  NTotal == 0,
	}
}

func (ft *ExportAttemptFilters) ByNTotalPtr(NTotal *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "n_total",
		Value:  NTotal,
		IsNil:  NTotal == nil,
		IsZero: NTotal != nil && (*NTotal) == 0,
	}
}

func (ft *ExportAttemptFilters) ByNExported(NExported int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "n_exported",
		Value:  NExported,
		IsNil:  NExported == 0,
	}
}

func (ft *ExportAttemptFilters) ByNExportedPtr(NExported *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "n_exported",
		Value:  NExported,
		IsNil:  NExported == nil,
		IsZero: NExported != nil && (*NExported) == 0,
	}
}

func (ft *ExportAttemptFilters) ByNError(NError int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "n_error",
		Value:  NError,
		IsNil:  NError == 0,
	}
}

func (ft *ExportAttemptFilters) ByNErrorPtr(NError *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "n_error",
		Value:  NError,
		IsNil:  NError == nil,
		IsZero: NError != nil && (*NError) == 0,
	}
}

func (ft *ExportAttemptFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *ExportAttemptFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *ExportAttemptFilters) ByDeletedAt(DeletedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "deleted_at",
		Value:  DeletedAt,
		IsNil:  DeletedAt.IsZero(),
	}
}

func (ft *ExportAttemptFilters) ByDeletedAtPtr(DeletedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "deleted_at",
		Value:  DeletedAt,
		IsNil:  DeletedAt == nil,
		IsZero: DeletedAt != nil && (*DeletedAt).IsZero(),
	}
}

func (ft *ExportAttemptFilters) ByStartedAt(StartedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "started_at",
		Value:  StartedAt,
		IsNil:  StartedAt.IsZero(),
	}
}

func (ft *ExportAttemptFilters) ByStartedAtPtr(StartedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "started_at",
		Value:  StartedAt,
		IsNil:  StartedAt == nil,
		IsZero: StartedAt != nil && (*StartedAt).IsZero(),
	}
}

func (ft *ExportAttemptFilters) ByDoneAt(DoneAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "done_at",
		Value:  DoneAt,
		IsNil:  DoneAt.IsZero(),
	}
}

func (ft *ExportAttemptFilters) ByDoneAtPtr(DoneAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "done_at",
		Value:  DoneAt,
		IsNil:  DoneAt == nil,
		IsZero: DoneAt != nil && (*DoneAt).IsZero(),
	}
}

func (ft *ExportAttemptFilters) ByExpiresAt(ExpiresAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "expires_at",
		Value:  ExpiresAt,
		IsNil:  ExpiresAt.IsZero(),
	}
}

func (ft *ExportAttemptFilters) ByExpiresAtPtr(ExpiresAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "expires_at",
		Value:  ExpiresAt,
		IsNil:  ExpiresAt == nil,
		IsZero: ExpiresAt != nil && (*ExpiresAt).IsZero(),
	}
}

type ImportAttemptFilters struct{ prefix string }

func NewImportAttemptFilters(prefix string) ImportAttemptFilters {
	return ImportAttemptFilters{prefix}
}

func (ft *ImportAttemptFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft ImportAttemptFilters) Prefix() string {
	return ft.prefix
}

func (ft *ImportAttemptFilters) ByID(ID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *ImportAttemptFilters) ByIDPtr(ID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *ImportAttemptFilters) ByUserID(UserID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "user_id",
		Value:  UserID,
		IsNil:  UserID == 0,
	}
}

func (ft *ImportAttemptFilters) ByUserIDPtr(UserID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "user_id",
		Value:  UserID,
		IsNil:  UserID == nil,
		IsZero: UserID != nil && (*UserID) == 0,
	}
}

func (ft *ImportAttemptFilters) ByAccountID(AccountID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "account_id",
		Value:  AccountID,
		IsNil:  AccountID == 0,
	}
}

func (ft *ImportAttemptFilters) ByAccountIDPtr(AccountID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "account_id",
		Value:  AccountID,
		IsNil:  AccountID == nil,
		IsZero: AccountID != nil && (*AccountID) == 0,
	}
}

func (ft *ImportAttemptFilters) ByOriginalFile(OriginalFile string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "original_file",
		Value:  OriginalFile,
		IsNil:  OriginalFile == "",
	}
}

func (ft *ImportAttemptFilters) ByOriginalFilePtr(OriginalFile *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "original_file",
		Value:  OriginalFile,
		IsNil:  OriginalFile == nil,
		IsZero: OriginalFile != nil && (*OriginalFile) == "",
	}
}

func (ft *ImportAttemptFilters) ByStoredFile(StoredFile string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "stored_file",
		Value:  StoredFile,
		IsNil:  StoredFile == "",
	}
}

func (ft *ImportAttemptFilters) ByStoredFilePtr(StoredFile *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "stored_file",
		Value:  StoredFile,
		IsNil:  StoredFile == nil,
		IsZero: StoredFile != nil && (*StoredFile) == "",
	}
}

func (ft *ImportAttemptFilters) ByType(Type etopmodel.ImportType) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "type",
		Value:  Type,
		IsNil:  Type == "",
	}
}

func (ft *ImportAttemptFilters) ByTypePtr(Type *etopmodel.ImportType) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "type",
		Value:  Type,
		IsNil:  Type == nil,
		IsZero: Type != nil && (*Type) == "",
	}
}

func (ft *ImportAttemptFilters) ByNCreated(NCreated int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "n_created",
		Value:  NCreated,
		IsNil:  NCreated == 0,
	}
}

func (ft *ImportAttemptFilters) ByNCreatedPtr(NCreated *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "n_created",
		Value:  NCreated,
		IsNil:  NCreated == nil,
		IsZero: NCreated != nil && (*NCreated) == 0,
	}
}

func (ft *ImportAttemptFilters) ByNUpdated(NUpdated int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "n_updated",
		Value:  NUpdated,
		IsNil:  NUpdated == 0,
	}
}

func (ft *ImportAttemptFilters) ByNUpdatedPtr(NUpdated *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "n_updated",
		Value:  NUpdated,
		IsNil:  NUpdated == nil,
		IsZero: NUpdated != nil && (*NUpdated) == 0,
	}
}

func (ft *ImportAttemptFilters) ByNError(NError int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "n_error",
		Value:  NError,
		IsNil:  NError == 0,
	}
}

func (ft *ImportAttemptFilters) ByNErrorPtr(NError *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "n_error",
		Value:  NError,
		IsNil:  NError == nil,
		IsZero: NError != nil && (*NError) == 0,
	}
}

func (ft *ImportAttemptFilters) ByStatus(Status status4.Status) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == 0,
	}
}

func (ft *ImportAttemptFilters) ByStatusPtr(Status *status4.Status) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == nil,
		IsZero: Status != nil && (*Status) == 0,
	}
}

func (ft *ImportAttemptFilters) ByErrorType(ErrorType string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "error_type",
		Value:  ErrorType,
		IsNil:  ErrorType == "",
	}
}

func (ft *ImportAttemptFilters) ByErrorTypePtr(ErrorType *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "error_type",
		Value:  ErrorType,
		IsNil:  ErrorType == nil,
		IsZero: ErrorType != nil && (*ErrorType) == "",
	}
}

func (ft *ImportAttemptFilters) ByDurationMs(DurationMs int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "duration_ms",
		Value:  DurationMs,
		IsNil:  DurationMs == 0,
	}
}

func (ft *ImportAttemptFilters) ByDurationMsPtr(DurationMs *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "duration_ms",
		Value:  DurationMs,
		IsNil:  DurationMs == nil,
		IsZero: DurationMs != nil && (*DurationMs) == 0,
	}
}

func (ft *ImportAttemptFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *ImportAttemptFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

type ShippingSourceFilters struct{ prefix string }

func NewShippingSourceFilters(prefix string) ShippingSourceFilters {
	return ShippingSourceFilters{prefix}
}

func (ft *ShippingSourceFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft ShippingSourceFilters) Prefix() string {
	return ft.prefix
}

func (ft *ShippingSourceFilters) ByID(ID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *ShippingSourceFilters) ByIDPtr(ID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *ShippingSourceFilters) ByName(Name string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "name",
		Value:  Name,
		IsNil:  Name == "",
	}
}

func (ft *ShippingSourceFilters) ByNamePtr(Name *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "name",
		Value:  Name,
		IsNil:  Name == nil,
		IsZero: Name != nil && (*Name) == "",
	}
}

func (ft *ShippingSourceFilters) ByUsername(Username string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "username",
		Value:  Username,
		IsNil:  Username == "",
	}
}

func (ft *ShippingSourceFilters) ByUsernamePtr(Username *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "username",
		Value:  Username,
		IsNil:  Username == nil,
		IsZero: Username != nil && (*Username) == "",
	}
}

func (ft *ShippingSourceFilters) ByType(Type string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "type",
		Value:  Type,
		IsNil:  Type == "",
	}
}

func (ft *ShippingSourceFilters) ByTypePtr(Type *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "type",
		Value:  Type,
		IsNil:  Type == nil,
		IsZero: Type != nil && (*Type) == "",
	}
}

func (ft *ShippingSourceFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *ShippingSourceFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *ShippingSourceFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *ShippingSourceFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}

type ShippingSourceInternalFilters struct{ prefix string }

func NewShippingSourceInternalFilters(prefix string) ShippingSourceInternalFilters {
	return ShippingSourceInternalFilters{prefix}
}

func (ft *ShippingSourceInternalFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft ShippingSourceInternalFilters) Prefix() string {
	return ft.prefix
}

func (ft *ShippingSourceInternalFilters) ByID(ID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *ShippingSourceInternalFilters) ByIDPtr(ID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *ShippingSourceInternalFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *ShippingSourceInternalFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *ShippingSourceInternalFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *ShippingSourceInternalFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}

func (ft *ShippingSourceInternalFilters) ByLastSyncAt(LastSyncAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "last_sync_at",
		Value:  LastSyncAt,
		IsNil:  LastSyncAt.IsZero(),
	}
}

func (ft *ShippingSourceInternalFilters) ByLastSyncAtPtr(LastSyncAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "last_sync_at",
		Value:  LastSyncAt,
		IsNil:  LastSyncAt == nil,
		IsZero: LastSyncAt != nil && (*LastSyncAt).IsZero(),
	}
}

func (ft *ShippingSourceInternalFilters) ByAccessToken(AccessToken string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "access_token",
		Value:  AccessToken,
		IsNil:  AccessToken == "",
	}
}

func (ft *ShippingSourceInternalFilters) ByAccessTokenPtr(AccessToken *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "access_token",
		Value:  AccessToken,
		IsNil:  AccessToken == nil,
		IsZero: AccessToken != nil && (*AccessToken) == "",
	}
}

func (ft *ShippingSourceInternalFilters) ByExpiresAt(ExpiresAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "expires_at",
		Value:  ExpiresAt,
		IsNil:  ExpiresAt.IsZero(),
	}
}

func (ft *ShippingSourceInternalFilters) ByExpiresAtPtr(ExpiresAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "expires_at",
		Value:  ExpiresAt,
		IsNil:  ExpiresAt == nil,
		IsZero: ExpiresAt != nil && (*ExpiresAt).IsZero(),
	}
}

type WebhookFilters struct{ prefix string }

func NewWebhookFilters(prefix string) WebhookFilters {
	return WebhookFilters{prefix}
}

func (ft *WebhookFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft WebhookFilters) Prefix() string {
	return ft.prefix
}

func (ft *WebhookFilters) ByID(ID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *WebhookFilters) ByIDPtr(ID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *WebhookFilters) ByAccountID(AccountID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "account_id",
		Value:  AccountID,
		IsNil:  AccountID == 0,
	}
}

func (ft *WebhookFilters) ByAccountIDPtr(AccountID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "account_id",
		Value:  AccountID,
		IsNil:  AccountID == nil,
		IsZero: AccountID != nil && (*AccountID) == 0,
	}
}

func (ft *WebhookFilters) ByURL(URL string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "url",
		Value:  URL,
		IsNil:  URL == "",
	}
}

func (ft *WebhookFilters) ByURLPtr(URL *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "url",
		Value:  URL,
		IsNil:  URL == nil,
		IsZero: URL != nil && (*URL) == "",
	}
}

func (ft *WebhookFilters) ByMetadata(Metadata string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "metadata",
		Value:  Metadata,
		IsNil:  Metadata == "",
	}
}

func (ft *WebhookFilters) ByMetadataPtr(Metadata *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "metadata",
		Value:  Metadata,
		IsNil:  Metadata == nil,
		IsZero: Metadata != nil && (*Metadata) == "",
	}
}

func (ft *WebhookFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *WebhookFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *WebhookFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *WebhookFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}

func (ft *WebhookFilters) ByDeletedAt(DeletedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "deleted_at",
		Value:  DeletedAt,
		IsNil:  DeletedAt.IsZero(),
	}
}

func (ft *WebhookFilters) ByDeletedAtPtr(DeletedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "deleted_at",
		Value:  DeletedAt,
		IsNil:  DeletedAt == nil,
		IsZero: DeletedAt != nil && (*DeletedAt).IsZero(),
	}
}
