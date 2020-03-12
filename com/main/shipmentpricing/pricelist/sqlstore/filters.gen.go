// +build !generator

// Code generated by generator sqlgen. DO NOT EDIT.

package sqlstore

import (
	time "time"

	sq "etop.vn/backend/pkg/common/sql/sq"
	dot "etop.vn/capi/dot"
)

type ShipmentPriceListFilters struct{ prefix string }

func NewShipmentPriceListFilters(prefix string) ShipmentPriceListFilters {
	return ShipmentPriceListFilters{prefix}
}

func (ft *ShipmentPriceListFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft ShipmentPriceListFilters) Prefix() string {
	return ft.prefix
}

func (ft *ShipmentPriceListFilters) ByID(ID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *ShipmentPriceListFilters) ByIDPtr(ID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *ShipmentPriceListFilters) ByName(Name string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "name",
		Value:  Name,
		IsNil:  Name == "",
	}
}

func (ft *ShipmentPriceListFilters) ByNamePtr(Name *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "name",
		Value:  Name,
		IsNil:  Name == nil,
		IsZero: Name != nil && (*Name) == "",
	}
}

func (ft *ShipmentPriceListFilters) ByDescription(Description string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "description",
		Value:  Description,
		IsNil:  Description == "",
	}
}

func (ft *ShipmentPriceListFilters) ByDescriptionPtr(Description *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "description",
		Value:  Description,
		IsNil:  Description == nil,
		IsZero: Description != nil && (*Description) == "",
	}
}

func (ft *ShipmentPriceListFilters) ByIsActive(IsActive bool) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "is_active",
		Value:  IsActive,
		IsNil:  bool(!IsActive),
	}
}

func (ft *ShipmentPriceListFilters) ByIsActivePtr(IsActive *bool) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "is_active",
		Value:  IsActive,
		IsNil:  IsActive == nil,
		IsZero: IsActive != nil && bool(!(*IsActive)),
	}
}

func (ft *ShipmentPriceListFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *ShipmentPriceListFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *ShipmentPriceListFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *ShipmentPriceListFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}

func (ft *ShipmentPriceListFilters) ByDeletedAt(DeletedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "deleted_at",
		Value:  DeletedAt,
		IsNil:  DeletedAt.IsZero(),
	}
}

func (ft *ShipmentPriceListFilters) ByDeletedAtPtr(DeletedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "deleted_at",
		Value:  DeletedAt,
		IsNil:  DeletedAt == nil,
		IsZero: DeletedAt != nil && (*DeletedAt).IsZero(),
	}
}

func (ft *ShipmentPriceListFilters) ByWLPartnerID(WLPartnerID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "wl_partner_id",
		Value:  WLPartnerID,
		IsNil:  WLPartnerID == 0,
	}
}

func (ft *ShipmentPriceListFilters) ByWLPartnerIDPtr(WLPartnerID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "wl_partner_id",
		Value:  WLPartnerID,
		IsNil:  WLPartnerID == nil,
		IsZero: WLPartnerID != nil && (*WLPartnerID) == 0,
	}
}
