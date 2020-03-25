// +build !generator

// Code generated by generator sqlgen. DO NOT EDIT.

package sqlstore

import (
	"time"

	"etop.vn/api/top/types/etc/shipping_provider"
	"etop.vn/api/top/types/etc/status3"
	sq "etop.vn/backend/pkg/common/sql/sq"
	"etop.vn/capi/dot"
)

type MoneyTransactionShippingFilters struct{ prefix string }

func NewMoneyTransactionShippingFilters(prefix string) MoneyTransactionShippingFilters {
	return MoneyTransactionShippingFilters{prefix}
}

func (ft *MoneyTransactionShippingFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft MoneyTransactionShippingFilters) Prefix() string {
	return ft.prefix
}

func (ft *MoneyTransactionShippingFilters) ByID(ID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *MoneyTransactionShippingFilters) ByIDPtr(ID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *MoneyTransactionShippingFilters) ByShopID(ShopID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == 0,
	}
}

func (ft *MoneyTransactionShippingFilters) ByShopIDPtr(ShopID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "shop_id",
		Value:  ShopID,
		IsNil:  ShopID == nil,
		IsZero: ShopID != nil && (*ShopID) == 0,
	}
}

func (ft *MoneyTransactionShippingFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *MoneyTransactionShippingFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *MoneyTransactionShippingFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *MoneyTransactionShippingFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}

func (ft *MoneyTransactionShippingFilters) ByClosedAt(ClosedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "closed_at",
		Value:  ClosedAt,
		IsNil:  ClosedAt.IsZero(),
	}
}

func (ft *MoneyTransactionShippingFilters) ByClosedAtPtr(ClosedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "closed_at",
		Value:  ClosedAt,
		IsNil:  ClosedAt == nil,
		IsZero: ClosedAt != nil && (*ClosedAt).IsZero(),
	}
}

func (ft *MoneyTransactionShippingFilters) ByStatus(Status status3.Status) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == 0,
	}
}

func (ft *MoneyTransactionShippingFilters) ByStatusPtr(Status *status3.Status) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == nil,
		IsZero: Status != nil && (*Status) == 0,
	}
}

func (ft *MoneyTransactionShippingFilters) ByTotalCOD(TotalCOD int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "total_cod",
		Value:  TotalCOD,
		IsNil:  TotalCOD == 0,
	}
}

func (ft *MoneyTransactionShippingFilters) ByTotalCODPtr(TotalCOD *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "total_cod",
		Value:  TotalCOD,
		IsNil:  TotalCOD == nil,
		IsZero: TotalCOD != nil && (*TotalCOD) == 0,
	}
}

func (ft *MoneyTransactionShippingFilters) ByTotalAmount(TotalAmount int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "total_amount",
		Value:  TotalAmount,
		IsNil:  TotalAmount == 0,
	}
}

func (ft *MoneyTransactionShippingFilters) ByTotalAmountPtr(TotalAmount *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "total_amount",
		Value:  TotalAmount,
		IsNil:  TotalAmount == nil,
		IsZero: TotalAmount != nil && (*TotalAmount) == 0,
	}
}

func (ft *MoneyTransactionShippingFilters) ByTotalOrders(TotalOrders int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "total_orders",
		Value:  TotalOrders,
		IsNil:  TotalOrders == 0,
	}
}

func (ft *MoneyTransactionShippingFilters) ByTotalOrdersPtr(TotalOrders *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "total_orders",
		Value:  TotalOrders,
		IsNil:  TotalOrders == nil,
		IsZero: TotalOrders != nil && (*TotalOrders) == 0,
	}
}

func (ft *MoneyTransactionShippingFilters) ByCode(Code string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "code",
		Value:  Code,
		IsNil:  Code == "",
	}
}

func (ft *MoneyTransactionShippingFilters) ByCodePtr(Code *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "code",
		Value:  Code,
		IsNil:  Code == nil,
		IsZero: Code != nil && (*Code) == "",
	}
}

func (ft *MoneyTransactionShippingFilters) ByMoneyTransactionShippingExternalID(MoneyTransactionShippingExternalID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "money_transaction_shipping_external_id",
		Value:  MoneyTransactionShippingExternalID,
		IsNil:  MoneyTransactionShippingExternalID == 0,
	}
}

func (ft *MoneyTransactionShippingFilters) ByMoneyTransactionShippingExternalIDPtr(MoneyTransactionShippingExternalID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "money_transaction_shipping_external_id",
		Value:  MoneyTransactionShippingExternalID,
		IsNil:  MoneyTransactionShippingExternalID == nil,
		IsZero: MoneyTransactionShippingExternalID != nil && (*MoneyTransactionShippingExternalID) == 0,
	}
}

func (ft *MoneyTransactionShippingFilters) ByMoneyTransactionShippingEtopID(MoneyTransactionShippingEtopID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "money_transaction_shipping_etop_id",
		Value:  MoneyTransactionShippingEtopID,
		IsNil:  MoneyTransactionShippingEtopID == 0,
	}
}

func (ft *MoneyTransactionShippingFilters) ByMoneyTransactionShippingEtopIDPtr(MoneyTransactionShippingEtopID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "money_transaction_shipping_etop_id",
		Value:  MoneyTransactionShippingEtopID,
		IsNil:  MoneyTransactionShippingEtopID == nil,
		IsZero: MoneyTransactionShippingEtopID != nil && (*MoneyTransactionShippingEtopID) == 0,
	}
}

func (ft *MoneyTransactionShippingFilters) ByProvider(Provider string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "provider",
		Value:  Provider,
		IsNil:  Provider == "",
	}
}

func (ft *MoneyTransactionShippingFilters) ByProviderPtr(Provider *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "provider",
		Value:  Provider,
		IsNil:  Provider == nil,
		IsZero: Provider != nil && (*Provider) == "",
	}
}

func (ft *MoneyTransactionShippingFilters) ByConfirmedAt(ConfirmedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "confirmed_at",
		Value:  ConfirmedAt,
		IsNil:  ConfirmedAt.IsZero(),
	}
}

func (ft *MoneyTransactionShippingFilters) ByConfirmedAtPtr(ConfirmedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "confirmed_at",
		Value:  ConfirmedAt,
		IsNil:  ConfirmedAt == nil,
		IsZero: ConfirmedAt != nil && (*ConfirmedAt).IsZero(),
	}
}

func (ft *MoneyTransactionShippingFilters) ByEtopTransferedAt(EtopTransferedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "etop_transfered_at",
		Value:  EtopTransferedAt,
		IsNil:  EtopTransferedAt.IsZero(),
	}
}

func (ft *MoneyTransactionShippingFilters) ByEtopTransferedAtPtr(EtopTransferedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "etop_transfered_at",
		Value:  EtopTransferedAt,
		IsNil:  EtopTransferedAt == nil,
		IsZero: EtopTransferedAt != nil && (*EtopTransferedAt).IsZero(),
	}
}

func (ft *MoneyTransactionShippingFilters) ByNote(Note string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "note",
		Value:  Note,
		IsNil:  Note == "",
	}
}

func (ft *MoneyTransactionShippingFilters) ByNotePtr(Note *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "note",
		Value:  Note,
		IsNil:  Note == nil,
		IsZero: Note != nil && (*Note) == "",
	}
}

func (ft *MoneyTransactionShippingFilters) ByInvoiceNumber(InvoiceNumber string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "invoice_number",
		Value:  InvoiceNumber,
		IsNil:  InvoiceNumber == "",
	}
}

func (ft *MoneyTransactionShippingFilters) ByInvoiceNumberPtr(InvoiceNumber *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "invoice_number",
		Value:  InvoiceNumber,
		IsNil:  InvoiceNumber == nil,
		IsZero: InvoiceNumber != nil && (*InvoiceNumber) == "",
	}
}

func (ft *MoneyTransactionShippingFilters) ByType(Type string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "type",
		Value:  Type,
		IsNil:  Type == "",
	}
}

func (ft *MoneyTransactionShippingFilters) ByTypePtr(Type *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "type",
		Value:  Type,
		IsNil:  Type == nil,
		IsZero: Type != nil && (*Type) == "",
	}
}

func (ft *MoneyTransactionShippingFilters) ByRid(Rid dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "rid",
		Value:  Rid,
		IsNil:  Rid == 0,
	}
}

func (ft *MoneyTransactionShippingFilters) ByRidPtr(Rid *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "rid",
		Value:  Rid,
		IsNil:  Rid == nil,
		IsZero: Rid != nil && (*Rid) == 0,
	}
}

type MoneyTransactionShippingEtopFilters struct{ prefix string }

func NewMoneyTransactionShippingEtopFilters(prefix string) MoneyTransactionShippingEtopFilters {
	return MoneyTransactionShippingEtopFilters{prefix}
}

func (ft *MoneyTransactionShippingEtopFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft MoneyTransactionShippingEtopFilters) Prefix() string {
	return ft.prefix
}

func (ft *MoneyTransactionShippingEtopFilters) ByID(ID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *MoneyTransactionShippingEtopFilters) ByIDPtr(ID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *MoneyTransactionShippingEtopFilters) ByCode(Code string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "code",
		Value:  Code,
		IsNil:  Code == "",
	}
}

func (ft *MoneyTransactionShippingEtopFilters) ByCodePtr(Code *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "code",
		Value:  Code,
		IsNil:  Code == nil,
		IsZero: Code != nil && (*Code) == "",
	}
}

func (ft *MoneyTransactionShippingEtopFilters) ByTotalCOD(TotalCOD int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "total_cod",
		Value:  TotalCOD,
		IsNil:  TotalCOD == 0,
	}
}

func (ft *MoneyTransactionShippingEtopFilters) ByTotalCODPtr(TotalCOD *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "total_cod",
		Value:  TotalCOD,
		IsNil:  TotalCOD == nil,
		IsZero: TotalCOD != nil && (*TotalCOD) == 0,
	}
}

func (ft *MoneyTransactionShippingEtopFilters) ByTotalOrders(TotalOrders int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "total_orders",
		Value:  TotalOrders,
		IsNil:  TotalOrders == 0,
	}
}

func (ft *MoneyTransactionShippingEtopFilters) ByTotalOrdersPtr(TotalOrders *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "total_orders",
		Value:  TotalOrders,
		IsNil:  TotalOrders == nil,
		IsZero: TotalOrders != nil && (*TotalOrders) == 0,
	}
}

func (ft *MoneyTransactionShippingEtopFilters) ByTotalAmount(TotalAmount int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "total_amount",
		Value:  TotalAmount,
		IsNil:  TotalAmount == 0,
	}
}

func (ft *MoneyTransactionShippingEtopFilters) ByTotalAmountPtr(TotalAmount *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "total_amount",
		Value:  TotalAmount,
		IsNil:  TotalAmount == nil,
		IsZero: TotalAmount != nil && (*TotalAmount) == 0,
	}
}

func (ft *MoneyTransactionShippingEtopFilters) ByTotalFee(TotalFee int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "total_fee",
		Value:  TotalFee,
		IsNil:  TotalFee == 0,
	}
}

func (ft *MoneyTransactionShippingEtopFilters) ByTotalFeePtr(TotalFee *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "total_fee",
		Value:  TotalFee,
		IsNil:  TotalFee == nil,
		IsZero: TotalFee != nil && (*TotalFee) == 0,
	}
}

func (ft *MoneyTransactionShippingEtopFilters) ByTotalMoneyTransaction(TotalMoneyTransaction int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "total_money_transaction",
		Value:  TotalMoneyTransaction,
		IsNil:  TotalMoneyTransaction == 0,
	}
}

func (ft *MoneyTransactionShippingEtopFilters) ByTotalMoneyTransactionPtr(TotalMoneyTransaction *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "total_money_transaction",
		Value:  TotalMoneyTransaction,
		IsNil:  TotalMoneyTransaction == nil,
		IsZero: TotalMoneyTransaction != nil && (*TotalMoneyTransaction) == 0,
	}
}

func (ft *MoneyTransactionShippingEtopFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *MoneyTransactionShippingEtopFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *MoneyTransactionShippingEtopFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *MoneyTransactionShippingEtopFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}

func (ft *MoneyTransactionShippingEtopFilters) ByConfirmedAt(ConfirmedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "confirmed_at",
		Value:  ConfirmedAt,
		IsNil:  ConfirmedAt.IsZero(),
	}
}

func (ft *MoneyTransactionShippingEtopFilters) ByConfirmedAtPtr(ConfirmedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "confirmed_at",
		Value:  ConfirmedAt,
		IsNil:  ConfirmedAt == nil,
		IsZero: ConfirmedAt != nil && (*ConfirmedAt).IsZero(),
	}
}

func (ft *MoneyTransactionShippingEtopFilters) ByStatus(Status status3.Status) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == 0,
	}
}

func (ft *MoneyTransactionShippingEtopFilters) ByStatusPtr(Status *status3.Status) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == nil,
		IsZero: Status != nil && (*Status) == 0,
	}
}

func (ft *MoneyTransactionShippingEtopFilters) ByNote(Note string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "note",
		Value:  Note,
		IsNil:  Note == "",
	}
}

func (ft *MoneyTransactionShippingEtopFilters) ByNotePtr(Note *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "note",
		Value:  Note,
		IsNil:  Note == nil,
		IsZero: Note != nil && (*Note) == "",
	}
}

func (ft *MoneyTransactionShippingEtopFilters) ByInvoiceNumber(InvoiceNumber string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "invoice_number",
		Value:  InvoiceNumber,
		IsNil:  InvoiceNumber == "",
	}
}

func (ft *MoneyTransactionShippingEtopFilters) ByInvoiceNumberPtr(InvoiceNumber *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "invoice_number",
		Value:  InvoiceNumber,
		IsNil:  InvoiceNumber == nil,
		IsZero: InvoiceNumber != nil && (*InvoiceNumber) == "",
	}
}

type MoneyTransactionShippingExternalFilters struct{ prefix string }

func NewMoneyTransactionShippingExternalFilters(prefix string) MoneyTransactionShippingExternalFilters {
	return MoneyTransactionShippingExternalFilters{prefix}
}

func (ft *MoneyTransactionShippingExternalFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft MoneyTransactionShippingExternalFilters) Prefix() string {
	return ft.prefix
}

func (ft *MoneyTransactionShippingExternalFilters) ByID(ID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *MoneyTransactionShippingExternalFilters) ByIDPtr(ID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *MoneyTransactionShippingExternalFilters) ByCode(Code string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "code",
		Value:  Code,
		IsNil:  Code == "",
	}
}

func (ft *MoneyTransactionShippingExternalFilters) ByCodePtr(Code *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "code",
		Value:  Code,
		IsNil:  Code == nil,
		IsZero: Code != nil && (*Code) == "",
	}
}

func (ft *MoneyTransactionShippingExternalFilters) ByTotalCOD(TotalCOD int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "total_cod",
		Value:  TotalCOD,
		IsNil:  TotalCOD == 0,
	}
}

func (ft *MoneyTransactionShippingExternalFilters) ByTotalCODPtr(TotalCOD *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "total_cod",
		Value:  TotalCOD,
		IsNil:  TotalCOD == nil,
		IsZero: TotalCOD != nil && (*TotalCOD) == 0,
	}
}

func (ft *MoneyTransactionShippingExternalFilters) ByTotalOrders(TotalOrders int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "total_orders",
		Value:  TotalOrders,
		IsNil:  TotalOrders == 0,
	}
}

func (ft *MoneyTransactionShippingExternalFilters) ByTotalOrdersPtr(TotalOrders *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "total_orders",
		Value:  TotalOrders,
		IsNil:  TotalOrders == nil,
		IsZero: TotalOrders != nil && (*TotalOrders) == 0,
	}
}

func (ft *MoneyTransactionShippingExternalFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *MoneyTransactionShippingExternalFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *MoneyTransactionShippingExternalFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *MoneyTransactionShippingExternalFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}

func (ft *MoneyTransactionShippingExternalFilters) ByStatus(Status status3.Status) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == 0,
	}
}

func (ft *MoneyTransactionShippingExternalFilters) ByStatusPtr(Status *status3.Status) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "status",
		Value:  Status,
		IsNil:  Status == nil,
		IsZero: Status != nil && (*Status) == 0,
	}
}

func (ft *MoneyTransactionShippingExternalFilters) ByExternalPaidAt(ExternalPaidAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "external_paid_at",
		Value:  ExternalPaidAt,
		IsNil:  ExternalPaidAt.IsZero(),
	}
}

func (ft *MoneyTransactionShippingExternalFilters) ByExternalPaidAtPtr(ExternalPaidAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "external_paid_at",
		Value:  ExternalPaidAt,
		IsNil:  ExternalPaidAt == nil,
		IsZero: ExternalPaidAt != nil && (*ExternalPaidAt).IsZero(),
	}
}

func (ft *MoneyTransactionShippingExternalFilters) ByProvider(Provider shipping_provider.ShippingProvider) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "provider",
		Value:  Provider,
		IsNil:  Provider == 0,
	}
}

func (ft *MoneyTransactionShippingExternalFilters) ByProviderPtr(Provider *shipping_provider.ShippingProvider) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "provider",
		Value:  Provider,
		IsNil:  Provider == nil,
		IsZero: Provider != nil && (*Provider) == 0,
	}
}

func (ft *MoneyTransactionShippingExternalFilters) ByNote(Note string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "note",
		Value:  Note,
		IsNil:  Note == "",
	}
}

func (ft *MoneyTransactionShippingExternalFilters) ByNotePtr(Note *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "note",
		Value:  Note,
		IsNil:  Note == nil,
		IsZero: Note != nil && (*Note) == "",
	}
}

func (ft *MoneyTransactionShippingExternalFilters) ByInvoiceNumber(InvoiceNumber string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "invoice_number",
		Value:  InvoiceNumber,
		IsNil:  InvoiceNumber == "",
	}
}

func (ft *MoneyTransactionShippingExternalFilters) ByInvoiceNumberPtr(InvoiceNumber *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "invoice_number",
		Value:  InvoiceNumber,
		IsNil:  InvoiceNumber == nil,
		IsZero: InvoiceNumber != nil && (*InvoiceNumber) == "",
	}
}

type MoneyTransactionShippingExternalLineFilters struct{ prefix string }

func NewMoneyTransactionShippingExternalLineFilters(prefix string) MoneyTransactionShippingExternalLineFilters {
	return MoneyTransactionShippingExternalLineFilters{prefix}
}

func (ft *MoneyTransactionShippingExternalLineFilters) Filter(pred string, args ...interface{}) sq.WriterTo {
	return sq.Filter(&ft.prefix, pred, args...)
}

func (ft MoneyTransactionShippingExternalLineFilters) Prefix() string {
	return ft.prefix
}

func (ft *MoneyTransactionShippingExternalLineFilters) ByID(ID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == 0,
	}
}

func (ft *MoneyTransactionShippingExternalLineFilters) ByIDPtr(ID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "id",
		Value:  ID,
		IsNil:  ID == nil,
		IsZero: ID != nil && (*ID) == 0,
	}
}

func (ft *MoneyTransactionShippingExternalLineFilters) ByExternalCode(ExternalCode string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "external_code",
		Value:  ExternalCode,
		IsNil:  ExternalCode == "",
	}
}

func (ft *MoneyTransactionShippingExternalLineFilters) ByExternalCodePtr(ExternalCode *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "external_code",
		Value:  ExternalCode,
		IsNil:  ExternalCode == nil,
		IsZero: ExternalCode != nil && (*ExternalCode) == "",
	}
}

func (ft *MoneyTransactionShippingExternalLineFilters) ByExternalCustomer(ExternalCustomer string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "external_customer",
		Value:  ExternalCustomer,
		IsNil:  ExternalCustomer == "",
	}
}

func (ft *MoneyTransactionShippingExternalLineFilters) ByExternalCustomerPtr(ExternalCustomer *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "external_customer",
		Value:  ExternalCustomer,
		IsNil:  ExternalCustomer == nil,
		IsZero: ExternalCustomer != nil && (*ExternalCustomer) == "",
	}
}

func (ft *MoneyTransactionShippingExternalLineFilters) ByExternalAddress(ExternalAddress string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "external_address",
		Value:  ExternalAddress,
		IsNil:  ExternalAddress == "",
	}
}

func (ft *MoneyTransactionShippingExternalLineFilters) ByExternalAddressPtr(ExternalAddress *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "external_address",
		Value:  ExternalAddress,
		IsNil:  ExternalAddress == nil,
		IsZero: ExternalAddress != nil && (*ExternalAddress) == "",
	}
}

func (ft *MoneyTransactionShippingExternalLineFilters) ByExternalTotalCOD(ExternalTotalCOD int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "external_total_cod",
		Value:  ExternalTotalCOD,
		IsNil:  ExternalTotalCOD == 0,
	}
}

func (ft *MoneyTransactionShippingExternalLineFilters) ByExternalTotalCODPtr(ExternalTotalCOD *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "external_total_cod",
		Value:  ExternalTotalCOD,
		IsNil:  ExternalTotalCOD == nil,
		IsZero: ExternalTotalCOD != nil && (*ExternalTotalCOD) == 0,
	}
}

func (ft *MoneyTransactionShippingExternalLineFilters) ByExternalCreatedAt(ExternalCreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "external_created_at",
		Value:  ExternalCreatedAt,
		IsNil:  ExternalCreatedAt.IsZero(),
	}
}

func (ft *MoneyTransactionShippingExternalLineFilters) ByExternalCreatedAtPtr(ExternalCreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "external_created_at",
		Value:  ExternalCreatedAt,
		IsNil:  ExternalCreatedAt == nil,
		IsZero: ExternalCreatedAt != nil && (*ExternalCreatedAt).IsZero(),
	}
}

func (ft *MoneyTransactionShippingExternalLineFilters) ByExternalClosedAt(ExternalClosedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "external_closed_at",
		Value:  ExternalClosedAt,
		IsNil:  ExternalClosedAt.IsZero(),
	}
}

func (ft *MoneyTransactionShippingExternalLineFilters) ByExternalClosedAtPtr(ExternalClosedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "external_closed_at",
		Value:  ExternalClosedAt,
		IsNil:  ExternalClosedAt == nil,
		IsZero: ExternalClosedAt != nil && (*ExternalClosedAt).IsZero(),
	}
}

func (ft *MoneyTransactionShippingExternalLineFilters) ByEtopFulfillmentIDRaw(EtopFulfillmentIDRaw string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "etop_fulfillment_id_raw",
		Value:  EtopFulfillmentIDRaw,
		IsNil:  EtopFulfillmentIDRaw == "",
	}
}

func (ft *MoneyTransactionShippingExternalLineFilters) ByEtopFulfillmentIDRawPtr(EtopFulfillmentIDRaw *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "etop_fulfillment_id_raw",
		Value:  EtopFulfillmentIDRaw,
		IsNil:  EtopFulfillmentIDRaw == nil,
		IsZero: EtopFulfillmentIDRaw != nil && (*EtopFulfillmentIDRaw) == "",
	}
}

func (ft *MoneyTransactionShippingExternalLineFilters) ByEtopFulfillmentID(EtopFulfillmentID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "etop_fulfillment_id",
		Value:  EtopFulfillmentID,
		IsNil:  EtopFulfillmentID == 0,
	}
}

func (ft *MoneyTransactionShippingExternalLineFilters) ByEtopFulfillmentIDPtr(EtopFulfillmentID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "etop_fulfillment_id",
		Value:  EtopFulfillmentID,
		IsNil:  EtopFulfillmentID == nil,
		IsZero: EtopFulfillmentID != nil && (*EtopFulfillmentID) == 0,
	}
}

func (ft *MoneyTransactionShippingExternalLineFilters) ByNote(Note string) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "note",
		Value:  Note,
		IsNil:  Note == "",
	}
}

func (ft *MoneyTransactionShippingExternalLineFilters) ByNotePtr(Note *string) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "note",
		Value:  Note,
		IsNil:  Note == nil,
		IsZero: Note != nil && (*Note) == "",
	}
}

func (ft *MoneyTransactionShippingExternalLineFilters) ByMoneyTransactionShippingExternalID(MoneyTransactionShippingExternalID dot.ID) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "money_transaction_shipping_external_id",
		Value:  MoneyTransactionShippingExternalID,
		IsNil:  MoneyTransactionShippingExternalID == 0,
	}
}

func (ft *MoneyTransactionShippingExternalLineFilters) ByMoneyTransactionShippingExternalIDPtr(MoneyTransactionShippingExternalID *dot.ID) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "money_transaction_shipping_external_id",
		Value:  MoneyTransactionShippingExternalID,
		IsNil:  MoneyTransactionShippingExternalID == nil,
		IsZero: MoneyTransactionShippingExternalID != nil && (*MoneyTransactionShippingExternalID) == 0,
	}
}

func (ft *MoneyTransactionShippingExternalLineFilters) ByCreatedAt(CreatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt.IsZero(),
	}
}

func (ft *MoneyTransactionShippingExternalLineFilters) ByCreatedAtPtr(CreatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "created_at",
		Value:  CreatedAt,
		IsNil:  CreatedAt == nil,
		IsZero: CreatedAt != nil && (*CreatedAt).IsZero(),
	}
}

func (ft *MoneyTransactionShippingExternalLineFilters) ByUpdatedAt(UpdatedAt time.Time) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt.IsZero(),
	}
}

func (ft *MoneyTransactionShippingExternalLineFilters) ByUpdatedAtPtr(UpdatedAt *time.Time) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "updated_at",
		Value:  UpdatedAt,
		IsNil:  UpdatedAt == nil,
		IsZero: UpdatedAt != nil && (*UpdatedAt).IsZero(),
	}
}

func (ft *MoneyTransactionShippingExternalLineFilters) ByExternalTotalShippingFee(ExternalTotalShippingFee int) *sq.ColumnFilter {
	return &sq.ColumnFilter{
		Prefix: &ft.prefix,
		Column: "external_total_shipping_fee",
		Value:  ExternalTotalShippingFee,
		IsNil:  ExternalTotalShippingFee == 0,
	}
}

func (ft *MoneyTransactionShippingExternalLineFilters) ByExternalTotalShippingFeePtr(ExternalTotalShippingFee *int) *sq.ColumnFilterPtr {
	return &sq.ColumnFilterPtr{
		Prefix: &ft.prefix,
		Column: "external_total_shipping_fee",
		Value:  ExternalTotalShippingFee,
		IsNil:  ExternalTotalShippingFee == nil,
		IsZero: ExternalTotalShippingFee != nil && (*ExternalTotalShippingFee) == 0,
	}
}
