package snappyimport

import "o.o/backend/pkg/common/imcsv"

var schema = imcsv.Schema{
	{
		Name:    "shipping_code",
		Display: "Mã VĐ",
		Norm:    "ma vd",
	},
	{
		Name:     "shipment_id",
		Display:  "Mã VĐ KH",
		Norm:     "ma vd kh",
		Optional: true,
	},
	{
		Name:     "district",
		Display:  "Quận/ Huyện",
		Norm:     "quan/ huyen",
		Optional: true,
	},
	{
		Name:     "customer",
		Display:  "Người nhận",
		Norm:     "nguoi nhan",
		Optional: true,
	},
	{
		Name:     "address",
		Display:  "Địa chỉ",
		Norm:     "dia chi",
		Optional: true,
	},
	{
		Name:     "phone",
		Display:  "SĐT",
		Norm:     "sdt",
		Optional: true,
	},
	{
		Name:     "shipping_state",
		Display:  "Trạng thái",
		Norm:     "trang thai",
		Optional: true,
	},
	{
		Name:     "weight",
		Display:  "Khối lượng hàng (g)",
		Norm:     "khoi luong hang (g)",
		Optional: true,
	},
	{
		Name:     "shipping_fee",
		Display:  "Tổng phí phí",
		Norm:     "tong chi phi",
		Optional: false,
	},
	{
		Name:     "cod_amount",
		Display:  "Tiền thu hộ",
		Norm:     "tien thu ho",
		Optional: false,
	},
	{
		Name:     "total_amount",
		Display:  "Tổng ",
		Norm:     "tong",
		Optional: false,
	},
}

var (
	indexer         = schema.Indexer()
	idxShipmentID   = indexer("shipment_id")
	idxShippingCode = indexer("shipping_code")
	idxShippingFee  = indexer("shipping_fee")
	idxCODAmount    = indexer("cod_amount")
)
