// +build !generator

// Code generated by generator sqlgen. DO NOT EDIT.

package sqlstore

import (
	time "time"

	status5 "o.o/api/top/types/etc/status5"
	ticket_ref_type "o.o/api/top/types/etc/ticket/ticket_ref_type"
	ticket_source "o.o/api/top/types/etc/ticket/ticket_source"
	ticket_state "o.o/api/top/types/etc/ticket/ticket_state"
	sq "o.o/backend/pkg/common/sql/sq"
	dot "o.o/capi/dot"
)

type TicketFilters struct{ prefix string }

func NewTicketFilters(prefix string) TicketFilters {
	return TicketFilters{prefix}
}

func (ft *TicketFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft TicketFilters) Prefix() string {
	return ft.prefix
}

func (ft *TicketFilters) ByID(ID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *TicketFilters) ByIDPtr(ID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *TicketFilters) ByCode(Code string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "code",
		Value:  Code,
		IsNil:  Code == "",
	}
}

func (ft *TicketFilters) ByCodePtr(Code *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "code",
		Value:  Code,
		IsNil:  Code == nil,
		IsZero: Code != nil && (*Code) == "",
	}
}

func (ft *TicketFilters) ByAccountID(AccountID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "account_id",
		Value:  AccountID,
		IsNil:  AccountID == 0,
	}
}

func (ft *TicketFilters) ByAccountIDPtr(AccountID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "account_id",
		Value:  AccountID,
		IsNil:  AccountID == nil,
		IsZero: AccountID != nil && (*AccountID) == 0,
	}
}

func (ft *TicketFilters) ByTitle(Title string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "title",
		Value:  Title,
		IsNil:  Title == "",
	}
}

func (ft *TicketFilters) ByTitlePtr(Title *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "title",
		Value:  Title,
		IsNil:  Title == nil,
		IsZero: Title != nil && (*Title) == "",
	}
}

func (ft *TicketFilters) ByDescription(Description string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "description",
		Value:  Description,
		IsNil:  Description == "",
	}
}

func (ft *TicketFilters) ByDescriptionPtr(Description *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "description",
		Value:  Description,
		IsNil:  Description == nil,
		IsZero: Description != nil && (*Description) == "",
	}
}

func (ft *TicketFilters) ByNote(Note string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "note",
		Value:  Note,
		IsNil:  Note == "",
	}
}

func (ft *TicketFilters) ByNotePtr(Note *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "note",
		Value:  Note,
		IsNil:  Note == nil,
		IsZero: Note != nil && (*Note) == "",
	}
}

func (ft *TicketFilters) ByRefID(RefID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "ref_id",
		Value:  RefID,
		IsNil:  RefID == 0,
	}
}

func (ft *TicketFilters) ByRefIDPtr(RefID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "ref_id",
		Value:  RefID,
		IsNil:  RefID == nil,
		IsZero: RefID != nil && (*RefID) == 0,
	}
}

func (ft *TicketFilters) ByRefType(RefType ticket_ref_type.TicketRefType) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "ref_type",
		Value:  RefType,
		IsNil:  RefType == 0,
	}
}

func (ft *TicketFilters) ByRefTypePtr(RefType *ticket_ref_type.TicketRefType) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "ref_type",
		Value:  RefType,
		IsNil:  RefType == nil,
		IsZero: RefType != nil && (*RefType) == 0,
	}
}

func (ft *TicketFilters) ByRefCode(RefCode string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "ref_code",
		Value:  RefCode,
		IsNil:  RefCode == "",
	}
}

func (ft *TicketFilters) ByRefCodePtr(RefCode *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "ref_code",
		Value:  RefCode,
		IsNil:  RefCode == nil,
		IsZero: RefCode != nil && (*RefCode) == "",
	}
}

func (ft *TicketFilters) BySource(Source ticket_source.TicketSource) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "source",
		Value:  Source,
		IsNil:  Source == 0,
	}
}

func (ft *TicketFilters) BySourcePtr(Source *ticket_source.TicketSource) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "source",
		Value:  Source,
		IsNil:  Source == nil,
		IsZero: Source != nil && (*Source) == 0,
	}
}

func (ft *TicketFilters) ByState(State ticket_state.TicketState) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "state",
		Value:  State,
		IsNil:  State == 0,
	}
}

func (ft *TicketFilters) ByStatePtr(State *ticket_state.TicketState) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "state",
		Value:  State,
		IsNil:  State == nil,
		IsZero: State != nil && (*State) == 0,
	}
}

func (ft *TicketFilters) ByStatus(Status status5.Status) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == 0,
	}
}

func (ft *TicketFilters) ByStatusPtr(Status *status5.Status) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == nil,
		IsZero: Status != nil && (*Status) == 0,
	}
}

func (ft *TicketFilters) ByCreatedBy(CreatedBy dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_by",
		Value:  CreatedBy,
		IsNil:  CreatedBy == 0,
	}
}

func (ft *TicketFilters) ByCreatedByPtr(CreatedBy *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_by",
		Value:  CreatedBy,
		IsNil:  CreatedBy == nil,
		IsZero: CreatedBy != nil && (*CreatedBy) == 0,
	}
}

func (ft *TicketFilters) ByUpdatedBy(UpdatedBy dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_by",
		Value:  UpdatedBy,
		IsNil:  UpdatedBy == 0,
	}
}

func (ft *TicketFilters) ByUpdatedByPtr(UpdatedBy *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_by",
		Value:  UpdatedBy,
		IsNil:  UpdatedBy == nil,
		IsZero: UpdatedBy != nil && (*UpdatedBy) == 0,
	}
}

func (ft *TicketFilters) ByConfirmedBy(ConfirmedBy dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "confirmed_by",
		Value:  ConfirmedBy,
		IsNil:  ConfirmedBy == 0,
	}
}

func (ft *TicketFilters) ByConfirmedByPtr(ConfirmedBy *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "confirmed_by",
		Value:  ConfirmedBy,
		IsNil:  ConfirmedBy == nil,
		IsZero: ConfirmedBy != nil && (*ConfirmedBy) == 0,
	}
}

func (ft *TicketFilters) ByClosedBy(ClosedBy dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "closed_by",
		Value:  ClosedBy,
		IsNil:  ClosedBy == 0,
	}
}

func (ft *TicketFilters) ByClosedByPtr(ClosedBy *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "closed_by",
		Value:  ClosedBy,
		IsNil:  ClosedBy == nil,
		IsZero: ClosedBy != nil && (*ClosedBy) == 0,
	}
}

func (ft *TicketFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *TicketFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *TicketFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *TicketFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}

func (ft *TicketFilters) ByConfirmedAt(ConfirmedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "confirmed_at",
		Value:  ConfirmedAt,
		IsNil:  ConfirmedAt.IsZero(),
	}
}

func (ft *TicketFilters) ByConfirmedAtPtr(ConfirmedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "confirmed_at",
		Value:  ConfirmedAt,
		IsNil:  ConfirmedAt == nil,
		IsZero: ConfirmedAt != nil && (*ConfirmedAt).IsZero(),
	}
}

func (ft *TicketFilters) ByClosedAt(ClosedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "closed_at",
		Value:  ClosedAt,
		IsNil:  ClosedAt.IsZero(),
	}
}

func (ft *TicketFilters) ByClosedAtPtr(ClosedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "closed_at",
		Value:  ClosedAt,
		IsNil:  ClosedAt == nil,
		IsZero: ClosedAt != nil && (*ClosedAt).IsZero(),
	}
}

type TicketCommentFilters struct{ prefix string }

func NewTicketCommentFilters(prefix string) TicketCommentFilters {
	return TicketCommentFilters{prefix}
}

func (ft *TicketCommentFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft TicketCommentFilters) Prefix() string {
	return ft.prefix
}

func (ft *TicketCommentFilters) ByID(ID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *TicketCommentFilters) ByIDPtr(ID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *TicketCommentFilters) ByTicketID(TicketID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "ticket_id",
		Value:  TicketID,
		IsNil:  TicketID == 0,
	}
}

func (ft *TicketCommentFilters) ByTicketIDPtr(TicketID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "ticket_id",
		Value:  TicketID,
		IsNil:  TicketID == nil,
		IsZero: TicketID != nil && (*TicketID) == 0,
	}
}

func (ft *TicketCommentFilters) ByCreatedBy(CreatedBy dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_by",
		Value:  CreatedBy,
		IsNil:  CreatedBy == 0,
	}
}

func (ft *TicketCommentFilters) ByCreatedByPtr(CreatedBy *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_by",
		Value:  CreatedBy,
		IsNil:  CreatedBy == nil,
		IsZero: CreatedBy != nil && (*CreatedBy) == 0,
	}
}

func (ft *TicketCommentFilters) ByAccountID(AccountID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "account_id",
		Value:  AccountID,
		IsNil:  AccountID == 0,
	}
}

func (ft *TicketCommentFilters) ByAccountIDPtr(AccountID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "account_id",
		Value:  AccountID,
		IsNil:  AccountID == nil,
		IsZero: AccountID != nil && (*AccountID) == 0,
	}
}

func (ft *TicketCommentFilters) ByParentID(ParentID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "parent_id",
		Value:  ParentID,
		IsNil:  ParentID == 0,
	}
}

func (ft *TicketCommentFilters) ByParentIDPtr(ParentID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "parent_id",
		Value:  ParentID,
		IsNil:  ParentID == nil,
		IsZero: ParentID != nil && (*ParentID) == 0,
	}
}

func (ft *TicketCommentFilters) ByMessage(Message string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "message",
		Value:  Message,
		IsNil:  Message == "",
	}
}

func (ft *TicketCommentFilters) ByMessagePtr(Message *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "message",
		Value:  Message,
		IsNil:  Message == nil,
		IsZero: Message != nil && (*Message) == "",
	}
}

func (ft *TicketCommentFilters) ByImageUrl(ImageUrl string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "image_url",
		Value:  ImageUrl,
		IsNil:  ImageUrl == "",
	}
}

func (ft *TicketCommentFilters) ByImageUrlPtr(ImageUrl *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "image_url",
		Value:  ImageUrl,
		IsNil:  ImageUrl == nil,
		IsZero: ImageUrl != nil && (*ImageUrl) == "",
	}
}

func (ft *TicketCommentFilters) ByDeletedAt(DeletedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "deleted_at",
		Value:  DeletedAt,
		IsNil:  DeletedAt.IsZero(),
	}
}

func (ft *TicketCommentFilters) ByDeletedAtPtr(DeletedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "deleted_at",
		Value:  DeletedAt,
		IsNil:  DeletedAt == nil,
		IsZero: DeletedAt != nil && (*DeletedAt).IsZero(),
	}
}

func (ft *TicketCommentFilters) ByDeletedBy(DeletedBy dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "deleted_by",
		Value:  DeletedBy,
		IsNil:  DeletedBy == 0,
	}
}

func (ft *TicketCommentFilters) ByDeletedByPtr(DeletedBy *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "deleted_by",
		Value:  DeletedBy,
		IsNil:  DeletedBy == nil,
		IsZero: DeletedBy != nil && (*DeletedBy) == 0,
	}
}

func (ft *TicketCommentFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *TicketCommentFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *TicketCommentFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *TicketCommentFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}

type TicketLabelFilters struct{ prefix string }

func NewTicketLabelFilters(prefix string) TicketLabelFilters {
	return TicketLabelFilters{prefix}
}

func (ft *TicketLabelFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft TicketLabelFilters) Prefix() string {
	return ft.prefix
}

func (ft *TicketLabelFilters) ByID(ID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *TicketLabelFilters) ByIDPtr(ID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *TicketLabelFilters) ByName(Name string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "name",
		Value:  Name,
		IsNil:  Name == "",
	}
}

func (ft *TicketLabelFilters) ByNamePtr(Name *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "name",
		Value:  Name,
		IsNil:  Name == nil,
		IsZero: Name != nil && (*Name) == "",
	}
}

func (ft *TicketLabelFilters) ByCode(Code string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "code",
		Value:  Code,
		IsNil:  Code == "",
	}
}

func (ft *TicketLabelFilters) ByCodePtr(Code *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "code",
		Value:  Code,
		IsNil:  Code == nil,
		IsZero: Code != nil && (*Code) == "",
	}
}

func (ft *TicketLabelFilters) ByColor(Color string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "color",
		Value:  Color,
		IsNil:  Color == "",
	}
}

func (ft *TicketLabelFilters) ByColorPtr(Color *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "color",
		Value:  Color,
		IsNil:  Color == nil,
		IsZero: Color != nil && (*Color) == "",
	}
}

func (ft *TicketLabelFilters) ByParentID(ParentID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "parent_id",
		Value:  ParentID,
		IsNil:  ParentID == 0,
	}
}

func (ft *TicketLabelFilters) ByParentIDPtr(ParentID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "parent_id",
		Value:  ParentID,
		IsNil:  ParentID == nil,
		IsZero: ParentID != nil && (*ParentID) == 0,
	}
}

type TicketSearchFilters struct{ prefix string }

func NewTicketSearchFilters(prefix string) TicketSearchFilters {
	return TicketSearchFilters{prefix}
}

func (ft *TicketSearchFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft TicketSearchFilters) Prefix() string {
	return ft.prefix
}

func (ft *TicketSearchFilters) ByID(ID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *TicketSearchFilters) ByIDPtr(ID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *TicketSearchFilters) ByTitleNorm(TitleNorm string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "title_norm",
		Value:  TitleNorm,
		IsNil:  TitleNorm == "",
	}
}

func (ft *TicketSearchFilters) ByTitleNormPtr(TitleNorm *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "title_norm",
		Value:  TitleNorm,
		IsNil:  TitleNorm == nil,
		IsZero: TitleNorm != nil && (*TitleNorm) == "",
	}
}