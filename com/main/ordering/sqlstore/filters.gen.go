// Generated by common/sq. DO NOT EDIT.

package sqlstore

import (
	"time"

	"etop.vn/api/top/types/etc/ghn_note_code"
	"etop.vn/api/top/types/etc/order_source"
	"etop.vn/api/top/types/etc/payment_method"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/api/top/types/etc/status4"
	"etop.vn/api/top/types/etc/status5"
	"etop.vn/api/top/types/etc/try_on"
	m "etop.vn/backend/com/main/ordering/model"
	"etop.vn/backend/pkg/common/sql/sq"
	"etop.vn/capi/dot"
)

type OrderFilters struct{ prefix string }

func NewOrderFilters(prefix string) OrderFilters {
	return OrderFilters{prefix}
}

func (ft *OrderFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft OrderFilters) Prefix() string {
	return ft.prefix
}

func (ft *OrderFilters) ByID(ID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *OrderFilters) ByIDPtr(ID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *OrderFilters) ByShopID(ShopID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == 0,
	}
}

func (ft *OrderFilters) ByShopIDPtr(ShopID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == nil,
		IsZero: ShopID != nil && (*ShopID) == 0,
	}
}

func (ft *OrderFilters) ByCode(Code string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "code",
		Value:  Code,
		IsNil:  Code == "",
	}
}

func (ft *OrderFilters) ByCodePtr(Code *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "code",
		Value:  Code,
		IsNil:  Code == nil,
		IsZero: Code != nil && (*Code) == "",
	}
}

func (ft *OrderFilters) ByEdCode(EdCode string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "ed_code",
		Value:  EdCode,
		IsNil:  EdCode == "",
	}
}

func (ft *OrderFilters) ByEdCodePtr(EdCode *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "ed_code",
		Value:  EdCode,
		IsNil:  EdCode == nil,
		IsZero: EdCode != nil && (*EdCode) == "",
	}
}

func (ft *OrderFilters) ByPartnerID(PartnerID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "partner_id",
		Value:  PartnerID,
		IsNil:  PartnerID == 0,
	}
}

func (ft *OrderFilters) ByPartnerIDPtr(PartnerID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "partner_id",
		Value:  PartnerID,
		IsNil:  PartnerID == nil,
		IsZero: PartnerID != nil && (*PartnerID) == 0,
	}
}

func (ft *OrderFilters) ByCurrency(Currency string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "currency",
		Value:  Currency,
		IsNil:  Currency == "",
	}
}

func (ft *OrderFilters) ByCurrencyPtr(Currency *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "currency",
		Value:  Currency,
		IsNil:  Currency == nil,
		IsZero: Currency != nil && (*Currency) == "",
	}
}

func (ft *OrderFilters) ByPaymentMethod(PaymentMethod payment_method.PaymentMethod) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "payment_method",
		Value:  PaymentMethod,
		IsNil:  PaymentMethod == 0,
	}
}

func (ft *OrderFilters) ByPaymentMethodPtr(PaymentMethod *payment_method.PaymentMethod) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "payment_method",
		Value:  PaymentMethod,
		IsNil:  PaymentMethod == nil,
		IsZero: PaymentMethod != nil && (*PaymentMethod) == 0,
	}
}

func (ft *OrderFilters) ByCustomerName(CustomerName string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "customer_name",
		Value:  CustomerName,
		IsNil:  CustomerName == "",
	}
}

func (ft *OrderFilters) ByCustomerNamePtr(CustomerName *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "customer_name",
		Value:  CustomerName,
		IsNil:  CustomerName == nil,
		IsZero: CustomerName != nil && (*CustomerName) == "",
	}
}

func (ft *OrderFilters) ByCustomerPhone(CustomerPhone string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "customer_phone",
		Value:  CustomerPhone,
		IsNil:  CustomerPhone == "",
	}
}

func (ft *OrderFilters) ByCustomerPhonePtr(CustomerPhone *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "customer_phone",
		Value:  CustomerPhone,
		IsNil:  CustomerPhone == nil,
		IsZero: CustomerPhone != nil && (*CustomerPhone) == "",
	}
}

func (ft *OrderFilters) ByCustomerEmail(CustomerEmail string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "customer_email",
		Value:  CustomerEmail,
		IsNil:  CustomerEmail == "",
	}
}

func (ft *OrderFilters) ByCustomerEmailPtr(CustomerEmail *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "customer_email",
		Value:  CustomerEmail,
		IsNil:  CustomerEmail == nil,
		IsZero: CustomerEmail != nil && (*CustomerEmail) == "",
	}
}

func (ft *OrderFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *OrderFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *OrderFilters) ByProcessedAt(ProcessedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "processed_at",
		Value:  ProcessedAt,
		IsNil:  ProcessedAt.IsZero(),
	}
}

func (ft *OrderFilters) ByProcessedAtPtr(ProcessedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "processed_at",
		Value:  ProcessedAt,
		IsNil:  ProcessedAt == nil,
		IsZero: ProcessedAt != nil && (*ProcessedAt).IsZero(),
	}
}

func (ft *OrderFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *OrderFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}

func (ft *OrderFilters) ByClosedAt(ClosedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "closed_at",
		Value:  ClosedAt,
		IsNil:  ClosedAt.IsZero(),
	}
}

func (ft *OrderFilters) ByClosedAtPtr(ClosedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "closed_at",
		Value:  ClosedAt,
		IsNil:  ClosedAt == nil,
		IsZero: ClosedAt != nil && (*ClosedAt).IsZero(),
	}
}

func (ft *OrderFilters) ByConfirmedAt(ConfirmedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "confirmed_at",
		Value:  ConfirmedAt,
		IsNil:  ConfirmedAt.IsZero(),
	}
}

func (ft *OrderFilters) ByConfirmedAtPtr(ConfirmedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "confirmed_at",
		Value:  ConfirmedAt,
		IsNil:  ConfirmedAt == nil,
		IsZero: ConfirmedAt != nil && (*ConfirmedAt).IsZero(),
	}
}

func (ft *OrderFilters) ByCancelledAt(CancelledAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "cancelled_at",
		Value:  CancelledAt,
		IsNil:  CancelledAt.IsZero(),
	}
}

func (ft *OrderFilters) ByCancelledAtPtr(CancelledAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "cancelled_at",
		Value:  CancelledAt,
		IsNil:  CancelledAt == nil,
		IsZero: CancelledAt != nil && (*CancelledAt).IsZero(),
	}
}

func (ft *OrderFilters) ByCancelReason(CancelReason string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "cancel_reason",
		Value:  CancelReason,
		IsNil:  CancelReason == "",
	}
}

func (ft *OrderFilters) ByCancelReasonPtr(CancelReason *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "cancel_reason",
		Value:  CancelReason,
		IsNil:  CancelReason == nil,
		IsZero: CancelReason != nil && (*CancelReason) == "",
	}
}

func (ft *OrderFilters) ByCustomerConfirm(CustomerConfirm status3.Status) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "customer_confirm",
		Value:  CustomerConfirm,
		IsNil:  CustomerConfirm == 0,
	}
}

func (ft *OrderFilters) ByCustomerConfirmPtr(CustomerConfirm *status3.Status) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "customer_confirm",
		Value:  CustomerConfirm,
		IsNil:  CustomerConfirm == nil,
		IsZero: CustomerConfirm != nil && (*CustomerConfirm) == 0,
	}
}

func (ft *OrderFilters) ByShopConfirm(ShopConfirm status3.Status) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "shop_confirm",
		Value:  ShopConfirm,
		IsNil:  ShopConfirm == 0,
	}
}

func (ft *OrderFilters) ByShopConfirmPtr(ShopConfirm *status3.Status) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "shop_confirm",
		Value:  ShopConfirm,
		IsNil:  ShopConfirm == nil,
		IsZero: ShopConfirm != nil && (*ShopConfirm) == 0,
	}
}

func (ft *OrderFilters) ByConfirmStatus(ConfirmStatus status3.Status) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "confirm_status",
		Value:  ConfirmStatus,
		IsNil:  ConfirmStatus == 0,
	}
}

func (ft *OrderFilters) ByConfirmStatusPtr(ConfirmStatus *status3.Status) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "confirm_status",
		Value:  ConfirmStatus,
		IsNil:  ConfirmStatus == nil,
		IsZero: ConfirmStatus != nil && (*ConfirmStatus) == 0,
	}
}

func (ft *OrderFilters) ByFulfillmentShippingStatus(FulfillmentShippingStatus status5.Status) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "fulfillment_shipping_status",
		Value:  FulfillmentShippingStatus,
		IsNil:  FulfillmentShippingStatus == 0,
	}
}

func (ft *OrderFilters) ByFulfillmentShippingStatusPtr(FulfillmentShippingStatus *status5.Status) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "fulfillment_shipping_status",
		Value:  FulfillmentShippingStatus,
		IsNil:  FulfillmentShippingStatus == nil,
		IsZero: FulfillmentShippingStatus != nil && (*FulfillmentShippingStatus) == 0,
	}
}

func (ft *OrderFilters) ByEtopPaymentStatus(EtopPaymentStatus status4.Status) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "etop_payment_status",
		Value:  EtopPaymentStatus,
		IsNil:  EtopPaymentStatus == 0,
	}
}

func (ft *OrderFilters) ByEtopPaymentStatusPtr(EtopPaymentStatus *status4.Status) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "etop_payment_status",
		Value:  EtopPaymentStatus,
		IsNil:  EtopPaymentStatus == nil,
		IsZero: EtopPaymentStatus != nil && (*EtopPaymentStatus) == 0,
	}
}

func (ft *OrderFilters) ByStatus(Status status5.Status) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == 0,
	}
}

func (ft *OrderFilters) ByStatusPtr(Status *status5.Status) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == nil,
		IsZero: Status != nil && (*Status) == 0,
	}
}

func (ft *OrderFilters) ByTotalItems(TotalItems int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "total_items",
		Value:  TotalItems,
		IsNil:  TotalItems == 0,
	}
}

func (ft *OrderFilters) ByTotalItemsPtr(TotalItems *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "total_items",
		Value:  TotalItems,
		IsNil:  TotalItems == nil,
		IsZero: TotalItems != nil && (*TotalItems) == 0,
	}
}

func (ft *OrderFilters) ByBasketValue(BasketValue int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "basket_value",
		Value:  BasketValue,
		IsNil:  BasketValue == 0,
	}
}

func (ft *OrderFilters) ByBasketValuePtr(BasketValue *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "basket_value",
		Value:  BasketValue,
		IsNil:  BasketValue == nil,
		IsZero: BasketValue != nil && (*BasketValue) == 0,
	}
}

func (ft *OrderFilters) ByTotalWeight(TotalWeight int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "total_weight",
		Value:  TotalWeight,
		IsNil:  TotalWeight == 0,
	}
}

func (ft *OrderFilters) ByTotalWeightPtr(TotalWeight *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "total_weight",
		Value:  TotalWeight,
		IsNil:  TotalWeight == nil,
		IsZero: TotalWeight != nil && (*TotalWeight) == 0,
	}
}

func (ft *OrderFilters) ByTotalTax(TotalTax int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "total_tax",
		Value:  TotalTax,
		IsNil:  TotalTax == 0,
	}
}

func (ft *OrderFilters) ByTotalTaxPtr(TotalTax *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "total_tax",
		Value:  TotalTax,
		IsNil:  TotalTax == nil,
		IsZero: TotalTax != nil && (*TotalTax) == 0,
	}
}

func (ft *OrderFilters) ByOrderDiscount(OrderDiscount int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "order_discount",
		Value:  OrderDiscount,
		IsNil:  OrderDiscount == 0,
	}
}

func (ft *OrderFilters) ByOrderDiscountPtr(OrderDiscount *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "order_discount",
		Value:  OrderDiscount,
		IsNil:  OrderDiscount == nil,
		IsZero: OrderDiscount != nil && (*OrderDiscount) == 0,
	}
}

func (ft *OrderFilters) ByTotalDiscount(TotalDiscount int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "total_discount",
		Value:  TotalDiscount,
		IsNil:  TotalDiscount == 0,
	}
}

func (ft *OrderFilters) ByTotalDiscountPtr(TotalDiscount *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "total_discount",
		Value:  TotalDiscount,
		IsNil:  TotalDiscount == nil,
		IsZero: TotalDiscount != nil && (*TotalDiscount) == 0,
	}
}

func (ft *OrderFilters) ByShopShippingFee(ShopShippingFee int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "shop_shipping_fee",
		Value:  ShopShippingFee,
		IsNil:  ShopShippingFee == 0,
	}
}

func (ft *OrderFilters) ByShopShippingFeePtr(ShopShippingFee *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "shop_shipping_fee",
		Value:  ShopShippingFee,
		IsNil:  ShopShippingFee == nil,
		IsZero: ShopShippingFee != nil && (*ShopShippingFee) == 0,
	}
}

func (ft *OrderFilters) ByTotalFee(TotalFee int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "total_fee",
		Value:  TotalFee,
		IsNil:  TotalFee == 0,
	}
}

func (ft *OrderFilters) ByTotalFeePtr(TotalFee *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "total_fee",
		Value:  TotalFee,
		IsNil:  TotalFee == nil,
		IsZero: TotalFee != nil && (*TotalFee) == 0,
	}
}

func (ft *OrderFilters) ByShopCOD(ShopCOD int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "shop_cod",
		Value:  ShopCOD,
		IsNil:  ShopCOD == 0,
	}
}

func (ft *OrderFilters) ByShopCODPtr(ShopCOD *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "shop_cod",
		Value:  ShopCOD,
		IsNil:  ShopCOD == nil,
		IsZero: ShopCOD != nil && (*ShopCOD) == 0,
	}
}

func (ft *OrderFilters) ByTotalAmount(TotalAmount int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "total_amount",
		Value:  TotalAmount,
		IsNil:  TotalAmount == 0,
	}
}

func (ft *OrderFilters) ByTotalAmountPtr(TotalAmount *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "total_amount",
		Value:  TotalAmount,
		IsNil:  TotalAmount == nil,
		IsZero: TotalAmount != nil && (*TotalAmount) == 0,
	}
}

func (ft *OrderFilters) ByOrderNote(OrderNote string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "order_note",
		Value:  OrderNote,
		IsNil:  OrderNote == "",
	}
}

func (ft *OrderFilters) ByOrderNotePtr(OrderNote *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "order_note",
		Value:  OrderNote,
		IsNil:  OrderNote == nil,
		IsZero: OrderNote != nil && (*OrderNote) == "",
	}
}

func (ft *OrderFilters) ByShopNote(ShopNote string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "shop_note",
		Value:  ShopNote,
		IsNil:  ShopNote == "",
	}
}

func (ft *OrderFilters) ByShopNotePtr(ShopNote *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "shop_note",
		Value:  ShopNote,
		IsNil:  ShopNote == nil,
		IsZero: ShopNote != nil && (*ShopNote) == "",
	}
}

func (ft *OrderFilters) ByShippingNote(ShippingNote string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "shipping_note",
		Value:  ShippingNote,
		IsNil:  ShippingNote == "",
	}
}

func (ft *OrderFilters) ByShippingNotePtr(ShippingNote *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "shipping_note",
		Value:  ShippingNote,
		IsNil:  ShippingNote == nil,
		IsZero: ShippingNote != nil && (*ShippingNote) == "",
	}
}

func (ft *OrderFilters) ByOrderSourceType(OrderSourceType order_source.Source) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "order_source_type",
		Value:  OrderSourceType,
		IsNil:  OrderSourceType == 0,
	}
}

func (ft *OrderFilters) ByOrderSourceTypePtr(OrderSourceType *order_source.Source) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "order_source_type",
		Value:  OrderSourceType,
		IsNil:  OrderSourceType == nil,
		IsZero: OrderSourceType != nil && (*OrderSourceType) == 0,
	}
}

func (ft *OrderFilters) ByOrderSourceID(OrderSourceID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "order_source_id",
		Value:  OrderSourceID,
		IsNil:  OrderSourceID == 0,
	}
}

func (ft *OrderFilters) ByOrderSourceIDPtr(OrderSourceID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "order_source_id",
		Value:  OrderSourceID,
		IsNil:  OrderSourceID == nil,
		IsZero: OrderSourceID != nil && (*OrderSourceID) == 0,
	}
}

func (ft *OrderFilters) ByExternalOrderID(ExternalOrderID string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "external_order_id",
		Value:  ExternalOrderID,
		IsNil:  ExternalOrderID == "",
	}
}

func (ft *OrderFilters) ByExternalOrderIDPtr(ExternalOrderID *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "external_order_id",
		Value:  ExternalOrderID,
		IsNil:  ExternalOrderID == nil,
		IsZero: ExternalOrderID != nil && (*ExternalOrderID) == "",
	}
}

func (ft *OrderFilters) ByReferenceURL(ReferenceURL string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "reference_url",
		Value:  ReferenceURL,
		IsNil:  ReferenceURL == "",
	}
}

func (ft *OrderFilters) ByReferenceURLPtr(ReferenceURL *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "reference_url",
		Value:  ReferenceURL,
		IsNil:  ReferenceURL == nil,
		IsZero: ReferenceURL != nil && (*ReferenceURL) == "",
	}
}

func (ft *OrderFilters) ByExternalURL(ExternalURL string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "external_url",
		Value:  ExternalURL,
		IsNil:  ExternalURL == "",
	}
}

func (ft *OrderFilters) ByExternalURLPtr(ExternalURL *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "external_url",
		Value:  ExternalURL,
		IsNil:  ExternalURL == nil,
		IsZero: ExternalURL != nil && (*ExternalURL) == "",
	}
}

func (ft *OrderFilters) ByIsOutsideEtop(IsOutsideEtop bool) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "is_outside_etop",
		Value:  IsOutsideEtop,
		IsNil:  bool(!IsOutsideEtop),
	}
}

func (ft *OrderFilters) ByIsOutsideEtopPtr(IsOutsideEtop *bool) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "is_outside_etop",
		Value:  IsOutsideEtop,
		IsNil:  IsOutsideEtop == nil,
		IsZero: IsOutsideEtop != nil && bool(!(*IsOutsideEtop)),
	}
}

func (ft *OrderFilters) ByGhnNoteCode(GhnNoteCode ghn_note_code.GHNNoteCode) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "ghn_note_code",
		Value:  GhnNoteCode,
		IsNil:  GhnNoteCode == 0,
	}
}

func (ft *OrderFilters) ByGhnNoteCodePtr(GhnNoteCode *ghn_note_code.GHNNoteCode) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "ghn_note_code",
		Value:  GhnNoteCode,
		IsNil:  GhnNoteCode == nil,
		IsZero: GhnNoteCode != nil && (*GhnNoteCode) == 0,
	}
}

func (ft *OrderFilters) ByTryOn(TryOn try_on.TryOnCode) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "try_on",
		Value:  TryOn,
		IsNil:  TryOn == 0,
	}
}

func (ft *OrderFilters) ByTryOnPtr(TryOn *try_on.TryOnCode) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "try_on",
		Value:  TryOn,
		IsNil:  TryOn == nil,
		IsZero: TryOn != nil && (*TryOn) == 0,
	}
}

func (ft *OrderFilters) ByCustomerNameNorm(CustomerNameNorm string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "customer_name_norm",
		Value:  CustomerNameNorm,
		IsNil:  CustomerNameNorm == "",
	}
}

func (ft *OrderFilters) ByCustomerNameNormPtr(CustomerNameNorm *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "customer_name_norm",
		Value:  CustomerNameNorm,
		IsNil:  CustomerNameNorm == nil,
		IsZero: CustomerNameNorm != nil && (*CustomerNameNorm) == "",
	}
}

func (ft *OrderFilters) ByProductNameNorm(ProductNameNorm string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "product_name_norm",
		Value:  ProductNameNorm,
		IsNil:  ProductNameNorm == "",
	}
}

func (ft *OrderFilters) ByProductNameNormPtr(ProductNameNorm *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "product_name_norm",
		Value:  ProductNameNorm,
		IsNil:  ProductNameNorm == nil,
		IsZero: ProductNameNorm != nil && (*ProductNameNorm) == "",
	}
}

func (ft *OrderFilters) ByFulfillmentType(FulfillmentType m.FulfillType) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "fulfillment_type",
		Value:  FulfillmentType,
		IsNil:  FulfillmentType == 0,
	}
}

func (ft *OrderFilters) ByFulfillmentTypePtr(FulfillmentType *m.FulfillType) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "fulfillment_type",
		Value:  FulfillmentType,
		IsNil:  FulfillmentType == nil,
		IsZero: FulfillmentType != nil && (*FulfillmentType) == 0,
	}
}

func (ft *OrderFilters) ByTradingShopID(TradingShopID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "trading_shop_id",
		Value:  TradingShopID,
		IsNil:  TradingShopID == 0,
	}
}

func (ft *OrderFilters) ByTradingShopIDPtr(TradingShopID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "trading_shop_id",
		Value:  TradingShopID,
		IsNil:  TradingShopID == nil,
		IsZero: TradingShopID != nil && (*TradingShopID) == 0,
	}
}

func (ft *OrderFilters) ByPaymentStatus(PaymentStatus status4.Status) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "payment_status",
		Value:  PaymentStatus,
		IsNil:  PaymentStatus == 0,
	}
}

func (ft *OrderFilters) ByPaymentStatusPtr(PaymentStatus *status4.Status) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "payment_status",
		Value:  PaymentStatus,
		IsNil:  PaymentStatus == nil,
		IsZero: PaymentStatus != nil && (*PaymentStatus) == 0,
	}
}

func (ft *OrderFilters) ByPaymentID(PaymentID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "payment_id",
		Value:  PaymentID,
		IsNil:  PaymentID == 0,
	}
}

func (ft *OrderFilters) ByPaymentIDPtr(PaymentID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "payment_id",
		Value:  PaymentID,
		IsNil:  PaymentID == nil,
		IsZero: PaymentID != nil && (*PaymentID) == 0,
	}
}

func (ft *OrderFilters) ByCustomerID(CustomerID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "customer_id",
		Value:  CustomerID,
		IsNil:  CustomerID == 0,
	}
}

func (ft *OrderFilters) ByCustomerIDPtr(CustomerID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "customer_id",
		Value:  CustomerID,
		IsNil:  CustomerID == nil,
		IsZero: CustomerID != nil && (*CustomerID) == 0,
	}
}

func (ft *OrderFilters) ByCreatedBy(CreatedBy dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_by",
		Value:  CreatedBy,
		IsNil:  CreatedBy == 0,
	}
}

func (ft *OrderFilters) ByCreatedByPtr(CreatedBy *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_by",
		Value:  CreatedBy,
		IsNil:  CreatedBy == nil,
		IsZero: CreatedBy != nil && (*CreatedBy) == 0,
	}
}

type OrderLineFilters struct{ prefix string }

func NewOrderLineFilters(prefix string) OrderLineFilters {
	return OrderLineFilters{prefix}
}

func (ft *OrderLineFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft OrderLineFilters) Prefix() string {
	return ft.prefix
}

func (ft *OrderLineFilters) ByOrderID(OrderID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "order_id",
		Value:  OrderID,
		IsNil:  OrderID == 0,
	}
}

func (ft *OrderLineFilters) ByOrderIDPtr(OrderID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "order_id",
		Value:  OrderID,
		IsNil:  OrderID == nil,
		IsZero: OrderID != nil && (*OrderID) == 0,
	}
}

func (ft *OrderLineFilters) ByVariantID(VariantID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "variant_id",
		Value:  VariantID,
		IsNil:  VariantID == 0,
	}
}

func (ft *OrderLineFilters) ByVariantIDPtr(VariantID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "variant_id",
		Value:  VariantID,
		IsNil:  VariantID == nil,
		IsZero: VariantID != nil && (*VariantID) == 0,
	}
}

func (ft *OrderLineFilters) ByProductName(ProductName string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "product_name",
		Value:  ProductName,
		IsNil:  ProductName == "",
	}
}

func (ft *OrderLineFilters) ByProductNamePtr(ProductName *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "product_name",
		Value:  ProductName,
		IsNil:  ProductName == nil,
		IsZero: ProductName != nil && (*ProductName) == "",
	}
}

func (ft *OrderLineFilters) ByProductID(ProductID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "product_id",
		Value:  ProductID,
		IsNil:  ProductID == 0,
	}
}

func (ft *OrderLineFilters) ByProductIDPtr(ProductID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "product_id",
		Value:  ProductID,
		IsNil:  ProductID == nil,
		IsZero: ProductID != nil && (*ProductID) == 0,
	}
}

func (ft *OrderLineFilters) ByShopID(ShopID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == 0,
	}
}

func (ft *OrderLineFilters) ByShopIDPtr(ShopID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == nil,
		IsZero: ShopID != nil && (*ShopID) == 0,
	}
}

func (ft *OrderLineFilters) ByWeight(Weight int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "weight",
		Value:  Weight,
		IsNil:  Weight == 0,
	}
}

func (ft *OrderLineFilters) ByWeightPtr(Weight *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "weight",
		Value:  Weight,
		IsNil:  Weight == nil,
		IsZero: Weight != nil && (*Weight) == 0,
	}
}

func (ft *OrderLineFilters) ByQuantity(Quantity int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "quantity",
		Value:  Quantity,
		IsNil:  Quantity == 0,
	}
}

func (ft *OrderLineFilters) ByQuantityPtr(Quantity *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "quantity",
		Value:  Quantity,
		IsNil:  Quantity == nil,
		IsZero: Quantity != nil && (*Quantity) == 0,
	}
}

func (ft *OrderLineFilters) ByListPrice(ListPrice int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "list_price",
		Value:  ListPrice,
		IsNil:  ListPrice == 0,
	}
}

func (ft *OrderLineFilters) ByListPricePtr(ListPrice *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "list_price",
		Value:  ListPrice,
		IsNil:  ListPrice == nil,
		IsZero: ListPrice != nil && (*ListPrice) == 0,
	}
}

func (ft *OrderLineFilters) ByRetailPrice(RetailPrice int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "retail_price",
		Value:  RetailPrice,
		IsNil:  RetailPrice == 0,
	}
}

func (ft *OrderLineFilters) ByRetailPricePtr(RetailPrice *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "retail_price",
		Value:  RetailPrice,
		IsNil:  RetailPrice == nil,
		IsZero: RetailPrice != nil && (*RetailPrice) == 0,
	}
}

func (ft *OrderLineFilters) ByPaymentPrice(PaymentPrice int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "payment_price",
		Value:  PaymentPrice,
		IsNil:  PaymentPrice == 0,
	}
}

func (ft *OrderLineFilters) ByPaymentPricePtr(PaymentPrice *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "payment_price",
		Value:  PaymentPrice,
		IsNil:  PaymentPrice == nil,
		IsZero: PaymentPrice != nil && (*PaymentPrice) == 0,
	}
}

func (ft *OrderLineFilters) ByLineAmount(LineAmount int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "line_amount",
		Value:  LineAmount,
		IsNil:  LineAmount == 0,
	}
}

func (ft *OrderLineFilters) ByLineAmountPtr(LineAmount *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "line_amount",
		Value:  LineAmount,
		IsNil:  LineAmount == nil,
		IsZero: LineAmount != nil && (*LineAmount) == 0,
	}
}

func (ft *OrderLineFilters) ByTotalDiscount(TotalDiscount int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "total_discount",
		Value:  TotalDiscount,
		IsNil:  TotalDiscount == 0,
	}
}

func (ft *OrderLineFilters) ByTotalDiscountPtr(TotalDiscount *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "total_discount",
		Value:  TotalDiscount,
		IsNil:  TotalDiscount == nil,
		IsZero: TotalDiscount != nil && (*TotalDiscount) == 0,
	}
}

func (ft *OrderLineFilters) ByTotalLineAmount(TotalLineAmount int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "total_line_amount",
		Value:  TotalLineAmount,
		IsNil:  TotalLineAmount == 0,
	}
}

func (ft *OrderLineFilters) ByTotalLineAmountPtr(TotalLineAmount *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "total_line_amount",
		Value:  TotalLineAmount,
		IsNil:  TotalLineAmount == nil,
		IsZero: TotalLineAmount != nil && (*TotalLineAmount) == 0,
	}
}

func (ft *OrderLineFilters) ByImageURL(ImageURL string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "image_url",
		Value:  ImageURL,
		IsNil:  ImageURL == "",
	}
}

func (ft *OrderLineFilters) ByImageURLPtr(ImageURL *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "image_url",
		Value:  ImageURL,
		IsNil:  ImageURL == nil,
		IsZero: ImageURL != nil && (*ImageURL) == "",
	}
}

func (ft *OrderLineFilters) ByIsOutsideEtop(IsOutsideEtop bool) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "is_outside_etop",
		Value:  IsOutsideEtop,
		IsNil:  bool(!IsOutsideEtop),
	}
}

func (ft *OrderLineFilters) ByIsOutsideEtopPtr(IsOutsideEtop *bool) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "is_outside_etop",
		Value:  IsOutsideEtop,
		IsNil:  IsOutsideEtop == nil,
		IsZero: IsOutsideEtop != nil && bool(!(*IsOutsideEtop)),
	}
}

func (ft *OrderLineFilters) ByCode(Code string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "code",
		Value:  Code,
		IsNil:  Code == "",
	}
}

func (ft *OrderLineFilters) ByCodePtr(Code *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "code",
		Value:  Code,
		IsNil:  Code == nil,
		IsZero: Code != nil && (*Code) == "",
	}
}
