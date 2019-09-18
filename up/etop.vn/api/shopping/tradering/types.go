package tradering

var (
	CustomerType = "customer"
	VendorType   = "vendor"
	CarrierType  = "carrier"
)

type ShopTrader struct {
	ID     int64
	ShopID int64
	Type   string
}
