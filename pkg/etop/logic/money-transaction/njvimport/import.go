package njvimport

import (
	"context"
	"mime/multipart"
	"strconv"

	"o.o/api/main/moneytx"
	"o.o/api/main/shipping"
	cm "o.o/backend/pkg/common"
	logicmoneytx "o.o/backend/pkg/etop/logic/money-transaction"
	moneytxtypes "o.o/backend/pkg/etop/logic/money-transaction/handlers/types"
	"o.o/common/l"
)

var _ moneytxtypes.ImportMoneyTx = &NJVImporter{}
var ll = l.New()

type NJVImporter struct {
	ShippingAggr  shipping.CommandBus
	ShippingQuery shipping.QueryBus
}

func (i *NJVImporter) ValidateAndReadFile(
	ctx context.Context, fileType string, file multipart.File,
) (lines []*moneytx.MoneyTransactionShippingExternalLine, err error) {
	rows, err := logicmoneytx.ReadXLSXFile(ctx, file)
	if err != nil {
		return nil, err
	}
	startRowNo, idx, err := validateSchema(ctx, rows)
	if err != nil {
		return nil, err
	}
	shippingLines, _ := parseRows(idx, startRowNo, rows)
	if len(shippingLines) == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "File không hợp lệ. Vui lòng kiểm tra lại.")
	}

	return ToMoneyTransactionShippingExternalLines(shippingLines), nil
}

func validateSchema(ctx context.Context, rows [][]string) (startRowNo int, idx indexes, err error) {
	for i, row := range rows {
		indexer, errs, _ := schema.ValidateSchema(ctx, &row)
		if len(errs) > 0 {
			continue
		}
		idx = idxes[0]
		idx.indexer = indexer
		startRowNo = i + 1
		return
	}
	return 0, idx, cm.Errorf(cm.FailedPrecondition, nil, "File import không hợp lệ")
}

func parseRows(idx indexes, startRowNo int, rows [][]string) (res []*RowMoneyTransaction, _ error) {
	for i := startRowNo; i < len(rows); i++ {
		row := rows[i]
		exCode := row[idx.orderNo]
		if exCode == "" {
			continue
		}

		codAmount, err := strconv.ParseFloat(row[idx.codAmount], 64)
		if err != nil {
			continue
		}

		res = append(res, &RowMoneyTransaction{
			ExCode:    exCode,
			CodAmount: codAmount,
		})
	}
	return
}
