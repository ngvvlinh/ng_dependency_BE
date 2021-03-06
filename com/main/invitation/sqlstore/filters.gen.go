// +build !generator

// Code generated by generator sqlgen. DO NOT EDIT.

package sqlstore

import (
	time "time"

	status3 "o.o/api/top/types/etc/status3"
	sq "o.o/backend/pkg/common/sql/sq"
	dot "o.o/capi/dot"
)

type InvitationFilters struct{ prefix string }

func NewInvitationFilters(prefix string) InvitationFilters {
	return InvitationFilters{prefix}
}

func (ft *InvitationFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft InvitationFilters) Prefix() string {
	return ft.prefix
}

func (ft *InvitationFilters) ByID(ID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *InvitationFilters) ByIDPtr(ID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *InvitationFilters) ByAccountID(AccountID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "account_id",
		Value:  AccountID,
		IsNil:  AccountID == 0,
	}
}

func (ft *InvitationFilters) ByAccountIDPtr(AccountID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "account_id",
		Value:  AccountID,
		IsNil:  AccountID == nil,
		IsZero: AccountID != nil && (*AccountID) == 0,
	}
}

func (ft *InvitationFilters) ByEmail(Email string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "email",
		Value:  Email,
		IsNil:  Email == "",
	}
}

func (ft *InvitationFilters) ByEmailPtr(Email *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "email",
		Value:  Email,
		IsNil:  Email == nil,
		IsZero: Email != nil && (*Email) == "",
	}
}

func (ft *InvitationFilters) ByPhone(Phone string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "phone",
		Value:  Phone,
		IsNil:  Phone == "",
	}
}

func (ft *InvitationFilters) ByPhonePtr(Phone *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "phone",
		Value:  Phone,
		IsNil:  Phone == nil,
		IsZero: Phone != nil && (*Phone) == "",
	}
}

func (ft *InvitationFilters) ByFullName(FullName string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "full_name",
		Value:  FullName,
		IsNil:  FullName == "",
	}
}

func (ft *InvitationFilters) ByFullNamePtr(FullName *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "full_name",
		Value:  FullName,
		IsNil:  FullName == nil,
		IsZero: FullName != nil && (*FullName) == "",
	}
}

func (ft *InvitationFilters) ByShortName(ShortName string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "short_name",
		Value:  ShortName,
		IsNil:  ShortName == "",
	}
}

func (ft *InvitationFilters) ByShortNamePtr(ShortName *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "short_name",
		Value:  ShortName,
		IsNil:  ShortName == nil,
		IsZero: ShortName != nil && (*ShortName) == "",
	}
}

func (ft *InvitationFilters) ByToken(Token string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "token",
		Value:  Token,
		IsNil:  Token == "",
	}
}

func (ft *InvitationFilters) ByTokenPtr(Token *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "token",
		Value:  Token,
		IsNil:  Token == nil,
		IsZero: Token != nil && (*Token) == "",
	}
}

func (ft *InvitationFilters) ByStatus(Status status3.Status) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == 0,
	}
}

func (ft *InvitationFilters) ByStatusPtr(Status *status3.Status) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == nil,
		IsZero: Status != nil && (*Status) == 0,
	}
}

func (ft *InvitationFilters) ByInvitedBy(InvitedBy dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "invited_by",
		Value:  InvitedBy,
		IsNil:  InvitedBy == 0,
	}
}

func (ft *InvitationFilters) ByInvitedByPtr(InvitedBy *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "invited_by",
		Value:  InvitedBy,
		IsNil:  InvitedBy == nil,
		IsZero: InvitedBy != nil && (*InvitedBy) == 0,
	}
}

func (ft *InvitationFilters) ByAcceptedAt(AcceptedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "accepted_at",
		Value:  AcceptedAt,
		IsNil:  AcceptedAt.IsZero(),
	}
}

func (ft *InvitationFilters) ByAcceptedAtPtr(AcceptedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "accepted_at",
		Value:  AcceptedAt,
		IsNil:  AcceptedAt == nil,
		IsZero: AcceptedAt != nil && (*AcceptedAt).IsZero(),
	}
}

func (ft *InvitationFilters) ByRejectedAt(RejectedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "rejected_at",
		Value:  RejectedAt,
		IsNil:  RejectedAt.IsZero(),
	}
}

func (ft *InvitationFilters) ByRejectedAtPtr(RejectedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "rejected_at",
		Value:  RejectedAt,
		IsNil:  RejectedAt == nil,
		IsZero: RejectedAt != nil && (*RejectedAt).IsZero(),
	}
}

func (ft *InvitationFilters) ByExpiresAt(ExpiresAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "expires_at",
		Value:  ExpiresAt,
		IsNil:  ExpiresAt.IsZero(),
	}
}

func (ft *InvitationFilters) ByExpiresAtPtr(ExpiresAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "expires_at",
		Value:  ExpiresAt,
		IsNil:  ExpiresAt == nil,
		IsZero: ExpiresAt != nil && (*ExpiresAt).IsZero(),
	}
}

func (ft *InvitationFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *InvitationFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *InvitationFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *InvitationFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}

func (ft *InvitationFilters) ByDeletedAt(DeletedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "deleted_at",
		Value:  DeletedAt,
		IsNil:  DeletedAt.IsZero(),
	}
}

func (ft *InvitationFilters) ByDeletedAtPtr(DeletedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "deleted_at",
		Value:  DeletedAt,
		IsNil:  DeletedAt == nil,
		IsZero: DeletedAt != nil && (*DeletedAt).IsZero(),
	}
}

func (ft *InvitationFilters) ByRid(Rid dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "rid",
		Value:  Rid,
		IsNil:  Rid == 0,
	}
}

func (ft *InvitationFilters) ByRidPtr(Rid *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "rid",
		Value:  Rid,
		IsNil:  Rid == nil,
		IsZero: Rid != nil && (*Rid) == 0,
	}
}

func (ft *InvitationFilters) ByInvitationURL(InvitationURL string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "invitation_url",
		Value:  InvitationURL,
		IsNil:  InvitationURL == "",
	}
}

func (ft *InvitationFilters) ByInvitationURLPtr(InvitationURL *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "invitation_url",
		Value:  InvitationURL,
		IsNil:  InvitationURL == nil,
		IsZero: InvitationURL != nil && (*InvitationURL) == "",
	}
}
