// +build !generator

// Code generated by generator sqlgen. DO NOT EDIT.

package sqlstore

import (
	time "time"

	shipping_payment_type "o.o/api/top/types/etc/shipping_payment_type"
	try_on "o.o/api/top/types/etc/try_on"
	sq "o.o/backend/pkg/common/sql/sq"
	dot "o.o/capi/dot"
)

type ShopSettingFilters struct{ prefix string }

func NewShopSettingFilters(prefix string) ShopSettingFilters {
	return ShopSettingFilters{prefix}
}

func (ft *ShopSettingFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft ShopSettingFilters) Prefix() string {
	return ft.prefix
}

func (ft *ShopSettingFilters) ByShopID(ShopID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == 0,
	}
}

func (ft *ShopSettingFilters) ByShopIDPtr(ShopID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == nil,
		IsZero: ShopID != nil && (*ShopID) == 0,
	}
}

func (ft *ShopSettingFilters) ByPaymentTypeID(PaymentTypeID shipping_payment_type.ShippingPaymentType) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "payment_type_id",
		Value:  PaymentTypeID,
		IsNil:  PaymentTypeID == 0,
	}
}

func (ft *ShopSettingFilters) ByPaymentTypeIDPtr(PaymentTypeID *shipping_payment_type.ShippingPaymentType) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "payment_type_id",
		Value:  PaymentTypeID,
		IsNil:  PaymentTypeID == nil,
		IsZero: PaymentTypeID != nil && (*PaymentTypeID) == 0,
	}
}

func (ft *ShopSettingFilters) ByReturnAddressID(ReturnAddressID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "return_address_id",
		Value:  ReturnAddressID,
		IsNil:  ReturnAddressID == 0,
	}
}

func (ft *ShopSettingFilters) ByReturnAddressIDPtr(ReturnAddressID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "return_address_id",
		Value:  ReturnAddressID,
		IsNil:  ReturnAddressID == nil,
		IsZero: ReturnAddressID != nil && (*ReturnAddressID) == 0,
	}
}

func (ft *ShopSettingFilters) ByTryOn(TryOn try_on.TryOnCode) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "try_on",
		Value:  TryOn,
		IsNil:  TryOn == 0,
	}
}

func (ft *ShopSettingFilters) ByTryOnPtr(TryOn *try_on.TryOnCode) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "try_on",
		Value:  TryOn,
		IsNil:  TryOn == nil,
		IsZero: TryOn != nil && (*TryOn) == 0,
	}
}

func (ft *ShopSettingFilters) ByShippingNote(ShippingNote string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "shipping_note",
		Value:  ShippingNote,
		IsNil:  ShippingNote == "",
	}
}

func (ft *ShopSettingFilters) ByShippingNotePtr(ShippingNote *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "shipping_note",
		Value:  ShippingNote,
		IsNil:  ShippingNote == nil,
		IsZero: ShippingNote != nil && (*ShippingNote) == "",
	}
}

func (ft *ShopSettingFilters) ByWeight(Weight int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "weight",
		Value:  Weight,
		IsNil:  Weight == 0,
	}
}

func (ft *ShopSettingFilters) ByWeightPtr(Weight *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "weight",
		Value:  Weight,
		IsNil:  Weight == nil,
		IsZero: Weight != nil && (*Weight) == 0,
	}
}

func (ft *ShopSettingFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *ShopSettingFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *ShopSettingFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *ShopSettingFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}

func (ft *ShopSettingFilters) ByAllowConnectDirectShipment(AllowConnectDirectShipment bool) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "allow_connect_direct_shipment",
		Value:  AllowConnectDirectShipment,
		IsNil:  bool(!AllowConnectDirectShipment),
	}
}

func (ft *ShopSettingFilters) ByAllowConnectDirectShipmentPtr(AllowConnectDirectShipment *bool) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "allow_connect_direct_shipment",
		Value:  AllowConnectDirectShipment,
		IsNil:  AllowConnectDirectShipment == nil,
		IsZero: AllowConnectDirectShipment != nil && bool(!(*AllowConnectDirectShipment)),
	}
}
