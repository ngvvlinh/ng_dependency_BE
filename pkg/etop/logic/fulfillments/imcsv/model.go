package imcsv

import "o.o/backend/pkg/common/imcsv"

var schema = imcsv.Schema{
	{
		Name:     "ed_code",
		Display:  "Mã nội bộ",
		Norm:     "ma noi bo",
		Optional: true,
	},
	{
		Name:    "customer_name",
		Display: "Tên khách hàng*",
		Norm:    "ten khach hang",
	},
	{
		Name:    "customer_phone",
		Display: "SĐT*",
		Norm:    "sdt",
	},
	{
		Name:    "shipping_address",
		Display: "Địa chỉ giao hàng*",
		Norm:    "dia chi giao hang",
	},
	{
		Name:     "province",
		Display:  "Tỉnh/ Thành phố",
		Norm:     "tinh thanh pho",
		Optional: true,
	},
	{
		Name:     "district",
		Display:  "Quận/ Huyện",
		Norm:     "quan huyen",
		Optional: true,
	},
	{
		Name:     "ward",
		Display:  "Phường/ Xã",
		Norm:     "phuong xa",
		Optional: true,
	},
	{
		Name:    "product_description",
		Display: "Mô tả sản phẩm*",
		Norm:    "mo ta san pham",
	},
	{
		Name:    "total_weight",
		Display: "Tổng khối lượng (g)*",
		Norm:    "tong khoi luong g",
	},
	{
		Name:    "basket_value",
		Display: "Tổng giá trị hàng hoá*",
		Norm:    "tong gia tri hang hoa",
	},
	{
		Name:     "include_insurance",
		Display:  "Bảo hiểm hàng hoá",
		Norm:     "bao hiem hang hoa",
		Optional: true,
	},
	{
		Name:    "cod_amount",
		Display: "COD*",
		Norm:    "cod",
	},
	{
		Name:     "shipping_note",
		Display:  "Ghi chú",
		Norm:     "ghi chu",
		Optional: true,
	},
}

var (
	indexer = schema.Indexer()

	idxFulfillmentEdCode  = indexer("ed_code")
	idxCustomerName       = indexer("customer_name")
	idxCustomerPhone      = indexer("customer_phone")
	idxShippingAddress    = indexer("shipping_address")
	idxProvince           = indexer("province")
	idxDistrict           = indexer("district")
	idxWard               = indexer("ward")
	idxProductDescription = indexer("product_description")
	idxTotalWeight        = indexer("total_weight")
	idxBasketValue        = indexer("basket_value")
	idxIncludeInsurance   = indexer("include_insurance")
	idxCODAmount          = indexer("cod_amount")
	idxShippingNote       = indexer("shipping_note")
)

type RowFulfillment struct {
	RowIndex int

	EdCode             string
	CustomerName       string
	CustomerPhone      string
	ShippingAddress    string
	Province           string
	District           string
	Ward               string
	ProductDescription string
	TotalWeight        int
	BasketValue        int
	IncludeInsurance   bool
	CODAmount          int
	ShippingNote       string
}
