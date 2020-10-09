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
		Name:     "customer_account",
		Display:  "Tên Tài Khoản Lấy Hàng\n",
		Norm:     "ten tai khoan lay hang",
		Optional: true,
	},
	{
		Name:     "sales_channel",
		Display:  "Kênh Bán Hàng\nSales Channel",
		Norm:     "kenh ban hang sales channel",
		Optional: true,
	},
	{
		Name:    "shipment_id",
		Display: "Mã Đơn Hàng\nShipment ID",
		Norm:    "ma don hang shipment id",
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
	idxCustomerAccount  = indexer("customer_account")
	idxSalesChannel     = indexer("sales_channel")
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
