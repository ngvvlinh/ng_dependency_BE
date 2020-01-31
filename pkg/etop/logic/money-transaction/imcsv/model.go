package imcsv

import (
	"etop.vn/backend/pkg/common/imcsv"
)

var schemaV0 = imcsv.Schema{
	{
		Name:    "ex_code",
		Display: "Mã đơn hàng",
		Norm:    "ma don hang",
	}, {
		Name:    "etop_code",
		Display: "Mã đơn hàng KH",
		Norm:    "ma don hang kh",
	}, {
		Name:    "created_at",
		Display: "Ngày tạo",
		Norm:    "ngay tao",
	}, {
		Name:    "closed_at",
		Display: "Ngày giao hàng thành công",
		Norm:    "ngay giao hang thanh cong",
	}, {
		Name:    "customer",
		Display: "Người nhận",
		Norm:    "nguoi nhan",
	}, {
		Name:    "address",
		Display: "Địa chỉ",
		Norm:    "dia chi",
	}, {
		Name:    "total_cod",
		Display: "COD",
		Norm:    "cod",
	},
}

type indexes struct {
	indexer imcsv.Indexer

	exCode    int
	etopCode  int
	createdAt int
	closedAt  int
	customer  int
	address   int
	totalCOD  int
}

var (
	schemas = []imcsv.Schema{schemaV0}
	idxes   = []indexes{
		initIndexes(schemaV0),
	}
	nColumn = len(schemaV0)
)

func initIndexes(schema imcsv.Schema) indexes {
	indexer := schema.Indexer()
	return indexes{
		exCode:    indexer("ex_code"),
		etopCode:  indexer("etop_code"),
		createdAt: indexer("created_at"),
		closedAt:  indexer("closed_at"),
		customer:  indexer("customer"),
		address:   indexer("address"),
		totalCOD:  indexer("total_cod"),
	}
}

type RowMoneyTransaction struct {
	ExCode    string
	EtopCode  string
	CreatedAt string
	ClosedAt  string
	Customer  string
	Address   string
	TotalCOD  string
}

func validateSchema(headerRow *[]string) (schema imcsv.Schema, idx indexes, errs []error, err error) {
	i, indexer, errs, err := imcsv.ValidateAgainstSchemas(headerRow, schemas)
	if err != nil || errs != nil {
		return
	}
	idx = idxes[i]
	idx.indexer = indexer
	return schemas[i], idx, nil, nil
}
