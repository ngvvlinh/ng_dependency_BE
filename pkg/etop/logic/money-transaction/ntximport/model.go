package ntximport

import "o.o/backend/pkg/common/imcsv"

var schema = imcsv.Schema{
	{
		Name:    "stt",
		Display: "STT",
		Norm:    "stt",
	},
	{
		Name:    "sender_phone",
		Display: "ĐT người gửi",
		Norm:    "dt nguoi gui",
	},
	{
		Name:    "waybill_number",
		Display: "Số Vận Đơn",
		Norm:    "so van don",
	},
	{
		Name:    "shipping_code",
		Display: "Số Tham Chiếu",
		Norm:    "so tham chieu",
	},
	{
		Name:    "ref_code",
		Display: "DO Khách Hàng",
		Norm:    "do khach hang",
	},
	{
		Name:    "payment_method",
		Display: "HTTT",
		Norm:    "httt",
	},
	{
		Name:    "cod_amount",
		Display: "Tiền Thu Hộ",
		Norm:    "tien thu ho",
	},
	{
		Name:    "shipping_charges",
		Display: "Cước vận chuyển cấn trừ",
		Norm:    "cuoc van chuyen can tru",
	},
	{
		Name:    "total_amount",
		Display: "Tổng TT",
		Norm:    "tong tt",
	},
}

type indexes struct {
	indexer imcsv.Indexer

	stt          int
	shippingCode int
	codAmount    int
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
		stt:          indexer("stt"),
		shippingCode: indexer("shipping_code"),
		codAmount:    indexer("cod_amount"),
	}
}
