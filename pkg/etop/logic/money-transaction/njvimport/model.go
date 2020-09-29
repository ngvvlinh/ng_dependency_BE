package njvimport

import (
	"o.o/api/main/moneytx"
	"o.o/backend/pkg/common/imcsv"
)

var schema = imcsv.Schema{
	{
		Name:    "tracking_id",
		Display: "tracking_id",
		Norm:    "tracking id",
	},
	{
		Name:    "created",
		Display: "Created",
		Norm:    "created",
	},
	{
		Name:    "status",
		Display: "Status",
		Norm:    "status",
	},
	{
		Name:    "updated",
		Display: "Updated",
		Norm:    "updated",
	},
	{
		Name:    "from_address_1",
		Display: "From address 1",
		Norm:    "from address 1",
	},
	{
		Name:    "to_address_1",
		Display: "To_address1",
		Norm:    "to address1",
	},
	{
		Name:    "from_address_2",
		Display: "From address 2",
		Norm:    "from address 2",
	},
	{
		Name:    "to_address_2",
		Display: "To_address2",
		Norm:    "to address2",
	},
	{
		Name:    "to_address_3",
		Display: "To_address3",
		Norm:    "to address3",
	},
	{
		Name:    "delivery_type_id",
		Display: "delivery_type_id",
		Norm:    "delivery type id",
	},
	{
		Name:    "weight",
		Display: "Weight (kg)",
		Norm:    "weight kg",
	},
	{
		Name:    "cod",
		Display: "COD",
		Norm:    "cod",
	},
	{
		Name:    "classification",
		Display: "Classification",
		Norm:    "classification",
	},
	{
		Name:    "zones",
		Display: "Zones",
		Norm:    "zones",
	},
	{
		Name:    "main_fee",
		Display: "Fees",
		Norm:    "fees",
	},
	{
		Name:    "rts_fee",
		Display: "RTS",
		Norm:    "rts",
	},
	{
		Name:    "ins_fee",
		Display: "ins",
		Norm:    "ins",
	},
	{
		Name:    "vat_fee",
		Display: "VAT",
		Norm:    "vat",
	},
	{
		Name:    "cod_fee",
		Display: "Phí thu hộ",
		Norm:    "phi thu ho",
	},
	{
		Name:    "total_fee",
		Display: "Total",
		Norm:    "total",
	},
}

type indexes struct {
	indexer imcsv.Indexer

	trackingID     int
	created        int
	status         int
	updated        int
	fromAddress1   int
	toAddress1     int
	fromAddress2   int
	toAddress2     int
	toAddress3     int
	deliveryTypeID int
	weight         int
	codAmount      int
	classification int
	zones          int
	mainFee        int
	rtsFee         int
	insFee         int
	vatFee         int
	codFee         int
	totalFee       int
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
		trackingID:     indexer("tracking_id"),
		created:        indexer("created"),
		status:         indexer("status"),
		updated:        indexer("updated"),
		fromAddress1:   indexer("from_address_1"),
		toAddress1:     indexer("to_address_1"),
		fromAddress2:   indexer("from_address_2"),
		toAddress2:     indexer("to_address_2"),
		toAddress3:     indexer("to_address_3"),
		deliveryTypeID: indexer("delivery_type_id"),
		weight:         indexer("weight"),
		codAmount:      indexer("cod"),
		classification: indexer("classification"),
		zones:          indexer("zones"),
		mainFee:        indexer("main_fee"),
		rtsFee:         indexer("rts_fee"),
		insFee:         indexer("ins_fee"),
		vatFee:         indexer("vat_fee"),
		codFee:         indexer("cod_fee"),
		totalFee:       indexer("total_fee"),
	}
}

type RowMoneyTransaction struct {
	ExCode    string
	MainFee   float64
	RtsFee    float64 // return fee
	InsFee    float64 // insurance fee
	VatFee    float64 // VAT fee
	CodFee    float64
	CodAmount float64
}

func ToMoneyTransactionShippingExternalLines(lines []*RowMoneyTransaction) []*moneytx.MoneyTransactionShippingExternalLine {
	if lines == nil {
		return nil
	}
	res := make([]*moneytx.MoneyTransactionShippingExternalLine, len(lines))
	for i, line := range lines {
		res[i] = line.ToModel()
	}
	return res
}

func (line *RowMoneyTransaction) ToModel() *moneytx.MoneyTransactionShippingExternalLine {
	return &moneytx.MoneyTransactionShippingExternalLine{
		ExternalCode:     line.ExCode,
		ExternalTotalCOD: int(line.CodAmount),
	}
}
