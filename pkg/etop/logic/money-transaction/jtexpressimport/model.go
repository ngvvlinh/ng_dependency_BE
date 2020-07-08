package jtexpressimport

import (
	"o.o/backend/pkg/common/imcsv"
)

var schema = imcsv.Schema{
	{
		Name:    "no",
		Display: "STT",
		Norm:    "stt",
	},
	{
		Name:     "created_at",
		Display:  "Ngày gửi",
		Norm:     "ngay gui",
		Optional: true,
	},
	{
		Name:    "shipping_code",
		Display: "Số vận đơn",
		Norm:    "so van don",
	},
	{
		Name:     "etop_code",
		Display:  "Mã vận đơn liên kết",
		Norm:     "ma van don lien ket",
		Optional: true,
	},
	{
		Name:     "from_province",
		Display:  "Tỉnh gửi",
		Norm:     "tinh gui",
		Optional: true,
	},
	{
		Name:     "to_province",
		Display:  "Tỉnh đến",
		Norm:     "tinh den",
		Optional: true,
	},
	{
		Name:    "weight",
		Display: "Trọng lượng",
		Norm:    "trong luong",
	},
	{
		Name:    "cod_amount",
		Display: "Số tiền COD",
		Norm:    "so tien cod",
	},
}

type indexes struct {
	indexer imcsv.Indexer

	no           int
	createdAt    int
	shippingCode int
	etopCode     int
	fromProvince int
	toProvince   int
	weight       int
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
		no:           indexer("no"),
		createdAt:    indexer("created_at"),
		shippingCode: indexer("shipping_code"),
		etopCode:     indexer("etop_code"),
		fromProvince: indexer("from_province"),
		toProvince:   indexer("to_province"),
		weight:       indexer("weight"),
		codAmount:    indexer("cod_amount"),
	}
}
