// Generated by common/sq. DO NOT EDIT.

package sqlstore

import (
	"time"

	"etop.vn/backend/pkg/common/sq"
)

type ExternalAccountAhamoveFilters struct{ prefix string }

func NewExternalAccountAhamoveFilters(prefix string) ExternalAccountAhamoveFilters {
	return ExternalAccountAhamoveFilters{prefix}
}

func (ft ExternalAccountAhamoveFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft ExternalAccountAhamoveFilters) Prefix() string {
	return ft.prefix
}

func (ft ExternalAccountAhamoveFilters) ByID(ID int64) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft ExternalAccountAhamoveFilters) ByIDPtr(ID *int64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft ExternalAccountAhamoveFilters) ByOwnerID(OwnerID int64) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "owner_id",
		Value:  OwnerID,
		IsNil:  OwnerID == 0,
	}
}

func (ft ExternalAccountAhamoveFilters) ByOwnerIDPtr(OwnerID *int64) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "owner_id",
		Value:  OwnerID,
		IsNil:  OwnerID == nil,
		IsZero: OwnerID != nil && (*OwnerID) == 0,
	}
}

func (ft ExternalAccountAhamoveFilters) ByPhone(Phone string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "phone",
		Value:  Phone,
		IsNil:  Phone == "",
	}
}

func (ft ExternalAccountAhamoveFilters) ByPhonePtr(Phone *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "phone",
		Value:  Phone,
		IsNil:  Phone == nil,
		IsZero: Phone != nil && (*Phone) == "",
	}
}

func (ft ExternalAccountAhamoveFilters) ByName(Name string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "name",
		Value:  Name,
		IsNil:  Name == "",
	}
}

func (ft ExternalAccountAhamoveFilters) ByNamePtr(Name *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "name",
		Value:  Name,
		IsNil:  Name == nil,
		IsZero: Name != nil && (*Name) == "",
	}
}

func (ft ExternalAccountAhamoveFilters) ByExternalID(ExternalID string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "external_id",
		Value:  ExternalID,
		IsNil:  ExternalID == "",
	}
}

func (ft ExternalAccountAhamoveFilters) ByExternalIDPtr(ExternalID *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "external_id",
		Value:  ExternalID,
		IsNil:  ExternalID == nil,
		IsZero: ExternalID != nil && (*ExternalID) == "",
	}
}

func (ft ExternalAccountAhamoveFilters) ByExternalVerified(ExternalVerified bool) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "external_verified",
		Value:  ExternalVerified,
		IsNil:  bool(!ExternalVerified),
	}
}

func (ft ExternalAccountAhamoveFilters) ByExternalVerifiedPtr(ExternalVerified *bool) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "external_verified",
		Value:  ExternalVerified,
		IsNil:  ExternalVerified == nil,
		IsZero: ExternalVerified != nil && bool(!(*ExternalVerified)),
	}
}

func (ft ExternalAccountAhamoveFilters) ByExternalCreatedAt(ExternalCreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "external_created_at",
		Value:  ExternalCreatedAt,
		IsNil:  ExternalCreatedAt.IsZero(),
	}
}

func (ft ExternalAccountAhamoveFilters) ByExternalCreatedAtPtr(ExternalCreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "external_created_at",
		Value:  ExternalCreatedAt,
		IsNil:  ExternalCreatedAt == nil,
		IsZero: ExternalCreatedAt != nil && (*ExternalCreatedAt).IsZero(),
	}
}

func (ft ExternalAccountAhamoveFilters) ByExternalToken(ExternalToken string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "external_token",
		Value:  ExternalToken,
		IsNil:  ExternalToken == "",
	}
}

func (ft ExternalAccountAhamoveFilters) ByExternalTokenPtr(ExternalToken *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "external_token",
		Value:  ExternalToken,
		IsNil:  ExternalToken == nil,
		IsZero: ExternalToken != nil && (*ExternalToken) == "",
	}
}

func (ft ExternalAccountAhamoveFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft ExternalAccountAhamoveFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft ExternalAccountAhamoveFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft ExternalAccountAhamoveFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}

func (ft ExternalAccountAhamoveFilters) ByLastSendVerifiedAt(LastSendVerifiedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "last_send_verified_at",
		Value:  LastSendVerifiedAt,
		IsNil:  LastSendVerifiedAt.IsZero(),
	}
}

func (ft ExternalAccountAhamoveFilters) ByLastSendVerifiedAtPtr(LastSendVerifiedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "last_send_verified_at",
		Value:  LastSendVerifiedAt,
		IsNil:  LastSendVerifiedAt == nil,
		IsZero: LastSendVerifiedAt != nil && (*LastSendVerifiedAt).IsZero(),
	}
}

func (ft ExternalAccountAhamoveFilters) ByExternalTicketID(ExternalTicketID string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "external_ticket_id",
		Value:  ExternalTicketID,
		IsNil:  ExternalTicketID == "",
	}
}

func (ft ExternalAccountAhamoveFilters) ByExternalTicketIDPtr(ExternalTicketID *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "external_ticket_id",
		Value:  ExternalTicketID,
		IsNil:  ExternalTicketID == nil,
		IsZero: ExternalTicketID != nil && (*ExternalTicketID) == "",
	}
}

func (ft ExternalAccountAhamoveFilters) ByIDCardFrontImg(IDCardFrontImg string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id_card_front_img",
		Value:  IDCardFrontImg,
		IsNil:  IDCardFrontImg == "",
	}
}

func (ft ExternalAccountAhamoveFilters) ByIDCardFrontImgPtr(IDCardFrontImg *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id_card_front_img",
		Value:  IDCardFrontImg,
		IsNil:  IDCardFrontImg == nil,
		IsZero: IDCardFrontImg != nil && (*IDCardFrontImg) == "",
	}
}

func (ft ExternalAccountAhamoveFilters) ByIDCardBackImg(IDCardBackImg string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id_card_back_img",
		Value:  IDCardBackImg,
		IsNil:  IDCardBackImg == "",
	}
}

func (ft ExternalAccountAhamoveFilters) ByIDCardBackImgPtr(IDCardBackImg *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id_card_back_img",
		Value:  IDCardBackImg,
		IsNil:  IDCardBackImg == nil,
		IsZero: IDCardBackImg != nil && (*IDCardBackImg) == "",
	}
}

func (ft ExternalAccountAhamoveFilters) ByPortraitImg(PortraitImg string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "portrait_img",
		Value:  PortraitImg,
		IsNil:  PortraitImg == "",
	}
}

func (ft ExternalAccountAhamoveFilters) ByPortraitImgPtr(PortraitImg *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "portrait_img",
		Value:  PortraitImg,
		IsNil:  PortraitImg == nil,
		IsZero: PortraitImg != nil && (*PortraitImg) == "",
	}
}

func (ft ExternalAccountAhamoveFilters) ByUploadedAt(UploadedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "uploaded_at",
		Value:  UploadedAt,
		IsNil:  UploadedAt.IsZero(),
	}
}

func (ft ExternalAccountAhamoveFilters) ByUploadedAtPtr(UploadedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "uploaded_at",
		Value:  UploadedAt,
		IsNil:  UploadedAt == nil,
		IsZero: UploadedAt != nil && (*UploadedAt).IsZero(),
	}
}
