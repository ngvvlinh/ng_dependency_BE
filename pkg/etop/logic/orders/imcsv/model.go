package imcsv

import (
	"etop.vn/api/main/catalog"
	"etop.vn/api/top/types/etc/ghn_note_code"
	ordermodel "etop.vn/backend/com/main/ordering/model"
	"etop.vn/backend/pkg/common/imcsv"
	"etop.vn/backend/pkg/common/validate"
)

var schema = imcsv.Schema{
	{
		Name:    "order_ed_code",
		Display: "Mã đơn hàng",
		Norm:    "ma don hang",
	},
	{
		Name:    "customer_name",
		Display: "Tên khách hàng",
		Norm:    "ten khach hang",
	},
	{
		Name:    "customer_phone",
		Display: "Số điện thoại",
		Norm:    "so dien thoai",
	},
	{
		Name:    "variant_ed_code",
		Display: "Mã phiên bản sản phẩm",
		Norm:    "ma phien ban san pham",
		Line:    true,
	},
	{
		Name:    "variant_name",
		Display: "Tên phiên bản sản phẩm",
		Norm:    "ten phien ban san pham",
		Line:    true,
	},
	{
		Name:    "line_quantity",
		Display: "Số lượng",
		Norm:    "so luong",
		Line:    true,
	},
	{
		Name:    "variant_retail_price",
		Display: "Đơn giá",
		Norm:    "don gia",
		Line:    true,
	},
	{
		Name:    "line_amount",
		Display: "Thành tiền (trước giảm giá)",
		Norm:    "thanh tien truoc giam gia",
		Line:    true,
	},
	{
		Name:     "line_discount_percent",
		Display:  "Giảm giá (%)",
		Norm:     "giam gia",
		Contains: "%",
		Line:     true,
	},
	{
		Name:     "line_discount_value",
		Display:  "Giảm giá (₫)",
		Norm:     "giam gia",
		Contains: "₫",
		Line:     true,
	},
	{
		Name:    "line_total_amount",
		Display: "Thành tiền (sau giảm giá)",
		Norm:    "thanh tien sau giam gia",
		Line:    true,
	},
	{
		Name:    "total_items",
		Display: "Tổng số lượng sản phẩm",
		Norm:    "tong so luong san pham",
	},
	{
		Name:    "total_weight",
		Display: "Tổng khối lượng (kg)",
		Norm:    "tong khoi luong kg",
	},
	{
		Name:    "shipping_address",
		Display: "Địa chỉ giao hàng",
		Norm:    "dia chi giao hang",
	},
	{
		Name:    "shipping_province",
		Display: "Địa chỉ giao hàng (tỉnh)",
		Norm:    "dia chi giao hang tinh",
	},
	{
		Name:    "shipping_district",
		Display: "Địa chỉ giao hàng (quận/huyện)",
		Norm:    "dia chi giao hang quan huyen",
	},
	{
		Name:    "shipping_ward",
		Display: "Địa chỉ giao hàng (phường/xã)",
		Norm:    "dia chi giao hang phuong xa",
	},
	{
		Name:    "shipping_note",
		Display: "Ghi chú giao hàng",
		Norm:    "ghi chu giao hang",
	},
	{
		Name:    "shipping_note_ghn",
		Display: "Ghi chú xem hàng",
		Norm:    "ghi chu xem hang",
	},
	{
		Name:    "basket_value",
		Display: "Tổng tiền hàng (trước giảm giá)",
		Norm:    "tong tien hang truoc giam gia",
	},
	{
		Name:    "basket_value_discounted",
		Display: "Tổng tiền hàng (sau giảm giá)",
		Norm:    "tong tien hang sau giam gia",
	},
	{
		Name:    "order_discount",
		Display: "Giảm giá đơn hàng (₫)",
		Norm:    "giam gia don hang",
	},
	{
		Name:     "fee_line_shipping",
		Display:  "Phí giao hàng tính cho khách",
		Norm:     "phi giao hang tinh cho khach",
		Optional: true,
	},
	{
		Name:     "fee_line_tax",
		Display:  "Thuế",
		Norm:     "thue",
		Optional: true,
	},
	{
		Name:     "fee_line_other",
		Display:  "Phí khác",
		Norm:     "phi khac",
		Optional: true,
	},
	{
		Name:     "total_fee",
		Display:  "Tổng phí",
		Norm:     "tong phi",
		Optional: true,
	},
	{
		Name:    "total_amount",
		Display: "Tổng tiền thanh toán",
		Norm:    "tong tien thanh toan",
	},
	{
		Name:    "shop_cod",
		Display: "Tiền thu hộ (COD)", // Tiền thu hộ (COD, nếu không thu hộ để trống)
		Norm:    "tien thu ho cod",
	},
	{
		Name:    "is_cod",
		Display: "Thu hộ",
		Norm:    "thu ho",
	},
	{
		Name:    "_",
		Display: "Vui lòng không sửa cột này",
		Norm:    "vui long khong sua",
		Hidden:  true,
	},
	{
		Name:    "_lines",
		Display: "Số dòng sản phẩm",
		Norm:    "so dong san pham",
		Hidden:  true,
	},
}

var (
	indexer = schema.Indexer()

	idxOrderEdCode           = indexer("order_ed_code")
	idxCustomerName          = indexer("customer_name")
	idxCustomerPhone         = indexer("customer_phone")
	idxVariantEdCode         = indexer("variant_ed_code")
	idxVariantName           = indexer("variant_name")
	idxLineQuantity          = indexer("line_quantity")
	idxVariantRetailPrice    = indexer("variant_retail_price")
	idxLineAmount            = indexer("line_amount")
	idxLineDiscountPercent   = indexer("line_discount_percent")
	idxLineDiscountValue     = indexer("line_discount_value")
	idxLineTotalAmount       = indexer("line_total_amount")
	idxTotalItems            = indexer("total_items")
	idxTotalWeight           = indexer("total_weight")
	idxShippingAddress       = indexer("shipping_address")
	idxShippingProvince      = indexer("shipping_province")
	idxShippingDistrict      = indexer("shipping_district")
	idxShippingWard          = indexer("shipping_ward")
	idxShippingNote          = indexer("shipping_note")
	idxShippingNoteGhn       = indexer("shipping_note_ghn")
	idxBasketValue           = indexer("basket_value")
	idxBasketValueDiscounted = indexer("basket_value_discounted")
	idxOrderDiscount         = indexer("order_discount")
	idxFeeLineShipping       = indexer("fee_line_shipping")
	idxFeeLineTax            = indexer("fee_line_tax")
	idxFeeLineOther          = indexer("fee_line_other")
	idxTotalFee              = indexer("total_fee")
	idxTotalAmount           = indexer("total_amount")
	idxShopCod               = indexer("shop_cod")
	idxIsCod                 = indexer("is_cod")
	idxUnderscore            = indexer("_")
	idxLines                 = indexer("_lines")
)

type RowOrder struct {
	RowIndex int

	OrderEdCode      string
	CustomerName     string
	CustomerPhone    string
	ShippingAddress  string
	ShippingProvince string
	ShippingDistrict string
	ShippingWard     string
	ShippingNote     string
	GHNNoteCode      ghn_note_code.GHNNoteCode

	BasketValue           int
	BasketValueDiscounted int
	OrderDiscount         int
	TotalAmount           int

	FeeLineShipping int
	FeeLineTax      int
	FeeLineOther    int
	TotalFee        int

	ShopCOD int
	IsCOD   bool

	TotalWeight int // converted to g
	TotalItems  int

	Lines []*RowOrderLine
}

func (m *RowOrder) Validate(idx imcsv.Indexer, mode Mode) (errs []error) {
	var col int
	r := m.RowIndex

	col = idxCustomerName
	if m.CustomerName == "" {
		errs = append(errs, imcsv.CellError(idx, r, col, "%v không được để trống.", schema[col].Display))
	}

	col = idxCustomerPhone
	if m.CustomerPhone == "" {
		errs = append(errs, imcsv.CellError(idx, r, col, "%v không được để trống.", schema[col].Display))
	} else {
		_, ok := validate.NormalizePhone(m.CustomerPhone)
		if !ok {
			errs = append(errs, imcsv.CellError(idx, r, col, "%v không hợp lệ.", schema[col].Display))
		}
	}

	col = idxTotalWeight
	if m.TotalWeight < 0 {
		errs = append(errs, imcsv.CellError(idx, r, col, "%v không hợp lệ.", schema[col].Display))

	} else if mode == ModeFulfillment && m.TotalWeight == 0 {
		errs = append(errs, imcsv.CellError(idx, r, col, "%v không được để trống.", schema[col].Display))
	}

	col = idxShippingAddress
	if mode == ModeFulfillment && m.ShippingAddress == "" {
		errs = append(errs, imcsv.CellError(idx, r, col, "%v không được để trống.", schema[col].Display))
	}

	col = idxShippingProvince
	if mode == ModeFulfillment && m.ShippingProvince == "" {
		errs = append(errs, imcsv.CellError(idx, r, col, "%v không được để trống.", schema[col].Display))
	}

	col = idxShippingDistrict
	if mode == ModeFulfillment && m.ShippingDistrict == "" {
		errs = append(errs, imcsv.CellError(idx, r, col, "%v không được để trống.", schema[col].Display))
	}

	// ward is no longer required
	//
	// col = idxShippingWard
	// if mode == ModeFulfillment && m.ShippingWard == "" {
	// 	errs = append(errs, imcsv.CellError(idx, r, col, "%v không được để trống.", schema[col].Display))
	// }

	col = idxShippingNoteGhn
	if mode == ModeFulfillment && m.GHNNoteCode == 0 {
		errs = append(errs, imcsv.CellError(idx, r, col, "%v không được để trống.", schema[col].Display))
	}

	return errs
}

func (m *RowOrder) GetFeeLines() []ordermodel.OrderFeeLine {
	var lines []ordermodel.OrderFeeLine
	if m.FeeLineShipping != 0 {
		lines = append(lines, ordermodel.OrderFeeLine{
			Amount: m.FeeLineShipping,
			Desc:   "Phí giao hàng tính cho khách",
			Name:   "Phí giao hàng tính cho khách",
			Type:   ordermodel.OrderFeeShipping,
		})
	}
	if m.FeeLineTax != 0 {
		lines = append(lines, ordermodel.OrderFeeLine{
			Amount: m.FeeLineTax,
			Desc:   "Thuế",
			Name:   "Thuế",
			Type:   ordermodel.OrderFeeTax,
		})
	}
	if m.FeeLineOther != 0 {
		lines = append(lines, ordermodel.OrderFeeLine{
			Amount: m.FeeLineOther,
			Desc:   "Phí khác",
			Name:   "Phí khác",
			Type:   ordermodel.OrderFeeOther,
		})
	}
	return lines
}

type RowOrderLine struct {
	RowIndex int

	VariantEdCode string
	VariantName   string

	RetailPrice         int
	PaymentPrice        int
	Quantity            int
	LineAmount          int
	LineDiscountPercent float64
	LineDiscountValue   int
	LineTotalDiscount   int
	LineTotalAmount     int

	XVariant *catalog.ShopVariantWithProduct
}

func (m *RowOrderLine) Validate(idx imcsv.Indexer, mode Mode) (errs []error) {
	var col int
	r := m.RowIndex

	if m.VariantEdCode == "" && m.VariantName == "" {
		errs = append(errs, imcsv.CellError(idx, r, col, "%v không được để trống.", schema[idxVariantName].Display))
	}

	col = idxVariantRetailPrice
	if m.RetailPrice <= 0 {
		errs = append(errs, imcsv.CellError(idx, r, col, "%v không hợp lệ.", schema[col].Display))
	}

	col = idxLineQuantity
	if m.Quantity <= 0 {
		errs = append(errs, imcsv.CellError(idx, r, col, "%v không hợp lệ.", schema[col].Display))
	}

	col = idxLineAmount
	if m.LineAmount != m.RetailPrice*m.Quantity {
		errs = append(errs, imcsv.CellError(idx, r, col, "%v không hợp lệ", schema[col].Display))
	}

	col = idxLineDiscountPercent
	if m.LineDiscountPercent < 0 || m.LineDiscountPercent > 100 {
		errs = append(errs, imcsv.CellError(idx, r, col, "%v không hợp lệ", schema[col].Display))
	}

	col = idxLineDiscountValue
	if m.LineDiscountValue < 0 {
		errs = append(errs, imcsv.CellError(idx, r, col, "%v không hợp lệ", schema[col].Display))
	}
	if m.LineDiscountValue > m.RetailPrice {
		errs = append(errs, imcsv.CellError(idx, r, col, "%v không được lớn hơn %v", schema[col].Display, schema[idxVariantRetailPrice].Display))
	}

	col = idxLineTotalAmount
	if m.LineTotalAmount < 0 {
		errs = append(errs, imcsv.CellError(idx, r, col, "%v không hợp lệ", schema[col].Display))
	}

	totalAmount := m.LineAmount - int(float64(m.LineAmount)*m.LineDiscountPercent) - m.LineDiscountValue*m.Quantity
	if m.LineTotalAmount != totalAmount {
		errs = append(errs, imcsv.CellError(idx, r, col, "%v phải bằng %v trừ đi cả hai loại giảm giá", schema[col].Display, schema[idxLineAmount].Display))
	}

	// Now update the payment price,
	// and update back (because of the division)
	// TODO(qv): this may cause incorrect result if the division has remaining
	m.PaymentPrice = totalAmount / m.Quantity
	m.LineTotalAmount = m.PaymentPrice * m.Quantity
	m.LineTotalDiscount = m.LineAmount - m.LineTotalAmount
	return errs
}
