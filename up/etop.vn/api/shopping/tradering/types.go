package tradering

var (
	CustomerType = "customer"
	SupplierType = "supplier"
	CarrierType  = "carrier"
)

type ShopTrader struct {
	ID       int64
	ShopID   int64
	Type     string
	FullName string
	Phone    string
}
