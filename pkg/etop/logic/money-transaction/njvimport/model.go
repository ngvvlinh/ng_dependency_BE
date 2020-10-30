package njvimport

import (
	"o.o/api/main/moneytx"
	"o.o/backend/pkg/common/imcsv"
)

var schema = imcsv.Schema{
	{
		Name:     "no",
		Display:  "No.",
		Norm:     "no",
		Optional: true,
	},
	{
		Name:    "order_no",
		Display: "Order No.",
		Norm:    "order no",
	},
	{
		Name:    "cod_amount",
		Display: "COD Amount",
		Norm:    "cod amount",
	},
	{
		Name:    "date_of_deliver",
		Display: "Date of deliver",
		Norm:    "date of deliver",
	},
	{
		Name:    "status",
		Display: "Status",
		Norm:    "status",
	},
	{
		Name:    "report_sent_date",
		Display: "Report sent date",
		Norm:    "report sent date",
	},
}

type indexes struct {
	indexer imcsv.Indexer

	no             int
	orderNo        int
	codAmount      int
	dateOfDeliver  int
	status         int
	reportSentDate int
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
		no:             indexer("no"),
		orderNo:        indexer("order_no"),
		codAmount:      indexer("cod_amount"),
		dateOfDeliver:  indexer("date_of_deliver"),
		status:         indexer("status"),
		reportSentDate: indexer("report_sent_date"),
	}
}

type RowMoneyTransaction struct {
	ExCode    string
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
