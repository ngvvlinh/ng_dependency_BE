package dhlimport

import (
	"time"

	"o.o/backend/pkg/common/imcsv"
)

var schema = imcsv.Schema{
	{
		Name:     "created_date",
		Display:  "Ngày Giao Hàng\nDel. Date",
		Norm:     "ngay giao hang del date",
		Optional: true,
	},
	{
		Name:     "customer_code",
		Display:  "Mã Khách Hàng\nMerchant",
		Norm:     "ma khach hang merchant",
		Optional: true,
	},
	{
		Name:     "customer_name",
		Display:  "Tên Khách Hàng\nMerchant Name",
		Norm:     "ten khach hang merchant name",
		Optional: true,
	},
	{
		Name:    "shipping_code",
		Display: "Mã Theo Dõi DHL\nDHL Tracking ID",
		Norm:    "ma theo doi dhl dhl tracking id",
	},
	{
		Name:    "cod_amount",
		Display: "Tiền Thu Hộ\nCOD Value",
		Norm:    "tien thu ho cod value",
	},
	{
		Name:     "currency",
		Display:  "Tiền Tệ Curr.",
		Norm:     "tien te curr",
		Optional: true,
	},
	{
		Name:     "date_of_payment",
		Display:  "Ngày Chuyển Khoản\nDate of Payment",
		Norm:     "ngay chuyen khoan date of payment",
		Optional: true,
	},
	{
		Name:    "shipment_id",
		Display: "Mã Đơn Hàng\nShipment ID",
		Norm:    "ma don hang shipment id",
	},
	{
		Name:     "consignee_name",
		Display:  "Tên Người Nhận Consigneee",
		Norm:     "ten nguoi nhan consigneee",
		Optional: true,
	},
}

var (
	indexer             = schema.Indexer()
	idxCreatedDate      = indexer("created_date")
	idxCustomerCode     = indexer("customer_code")
	idxCustomerCustomer = indexer("customer_name")
	idxShippingCode     = indexer("shipping_code")
	idxShipmentID       = indexer("shipment_id")
	idxCodAmount        = indexer("cod_amount")
	idxCurrency         = indexer("currency")
	idxDateOfPayment    = indexer("date_of_payment")
	idxConsigneeName    = indexer("consignee_name")
)

type RowMoneyTransaction struct {
	ExCode      string
	CreatedDate time.Time
	CODAmount   float64
}
