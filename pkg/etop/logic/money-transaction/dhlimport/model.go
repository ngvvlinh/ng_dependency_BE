package dhlimport

import (
	"time"

	"o.o/backend/pkg/common/imcsv"
)

var schema = imcsv.Schema{
	{
		Name:    "created_date",
		Display: "Ngày Giao Hàng\nDel. Date",
		Norm:    "ngay giao hang del date",
	},
	{
		Name:    "customer_code",
		Display: "Mã Khách Hàng\nMerchant",
		Norm:    "ma khach hang merchant",
	},
	{
		Name:    "customer_name",
		Display: "Tên Khách Hàng\nMerchant Name",
		Norm:    "ten khach hang merchant name",
	},
	{
		Name:    "customer_account",
		Display: "Tên Tài Khoản Lấy Hàng\n",
		Norm:    "ten tai khoan lay hang",
	},
	{
		Name:    "sales_channel",
		Display: "Kênh Bán Hàng\nSales Channel",
		Norm:    "kenh ban hang sales channel",
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
		Name:    "currency",
		Display: "Tiền Tệ Curr.",
		Norm:    "tien te curr",
	},
	{
		Name:    "date_of_payment",
		Display: "Ngày Chuyển Khoản\nDate of Payment",
		Norm:    "ngay chuyen khoan date of payment",
	},
	{
		Name:    "consignee_name",
		Display: "Tên Người Nhận Consigneee",
		Norm:    "ten nguoi nhan consigneee",
	},
}

type indexes struct {
	indexer imcsv.Indexer

	createdDate      int
	customerCode     int
	customerCustomer int
	customerAccount  int
	salesChannel     int
	shipmentID       int
	codAmount        int
	currency         int
	dateOfPayment    int
	consigneeName    int
}

var (
	schemas = []imcsv.Schema{schema}
	idxes   = []indexes{
		initIndexes(schema),
	}
	nColumn = len(schema)
)

func initIndexes(schema imcsv.Schema) indexes {
	indexer := schema.Indexer()
	return indexes{
		createdDate:      indexer("created_date"),
		customerCode:     indexer("customer_code"),
		customerCustomer: indexer("customer_name"),
		customerAccount:  indexer("customer_account"),
		salesChannel:     indexer("sales_channel"),
		shipmentID:       indexer("shipment_id"),
		codAmount:        indexer("cod_amount"),
		currency:         indexer("currency"),
		dateOfPayment:    indexer("date_of_payment"),
		consigneeName:    indexer("consignee_name"),
	}
}

type RowMoneyTransaction struct {
	ExCode      string
	CreatedDate time.Time
	CODAmount   float64
}
