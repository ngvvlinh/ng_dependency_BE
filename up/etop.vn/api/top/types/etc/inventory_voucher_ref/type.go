package inventory_voucher_ref

// +enum
// +enum:zero=null
type InventoryVoucherRef int

type NullInventoryVoucherRef struct {
	Enum  InventoryVoucherRef
	Valid bool
}

const (
	// +enum=unknown
	// +enum:RefName:Không nguồn
	Unknown InventoryVoucherRef = 0

	// +enum=refund
	// +enum:RefName:Trả hàng
	Refund InventoryVoucherRef = 1

	// +enum=purchase_refund
	// +enum:RefName:Trả hàng nhập
	PurchaseRefund InventoryVoucherRef = 2

	// +enum=stocktake
	// +enum:RefName:Kiểm kho
	StockTake InventoryVoucherRef = 3

	// +enum=purchase_order
	// +enum:RefName:Nhập hàng
	PurchaseOrder InventoryVoucherRef = 4

	// +enum=order
	// +enum:RefName:Bán hàng
	Order InventoryVoucherRef = 5

	// +enum=cancel_order
	// +enum:RefName: Hủy đơn hàng
	CancelOrder InventoryVoucherRef = 6
)
