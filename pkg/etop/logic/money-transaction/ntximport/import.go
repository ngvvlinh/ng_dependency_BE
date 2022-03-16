package ntximport

import (
	"context"
	"mime/multipart"
	cm "o.o/backend/pkg/common"
	logicmoneytx "o.o/backend/pkg/etop/logic/money-transaction"
	"strconv"

	"o.o/api/main/moneytx"
	"o.o/api/main/shipping"
	moneytxtypes "o.o/backend/pkg/etop/logic/money-transaction/handlers/types"
	"o.o/common/l"
)

var _ moneytxtypes.ImportMoneyTx = &NTXImporter{}
var ll = l.New()

type NTXImporter struct {
	ShippingAggr  shipping.CommandBus
	ShippingQuery shipping.QueryBus
}

func (i *NTXImporter) ValidateAndReadFile(ctx context.Context, fileType string, file multipart.File) ([]*moneytx.MoneyTransactionShippingExternalLine, error) {
	rows, err := logicmoneytx.ReadXLSXFile(ctx, file)
	if err != nil {
		return nil, err
	}

	startRowNo, idx, err := validateSchema(ctx, rows)
	if err != nil {
		return nil, err
	}

	moneyTxs, _ := parseRows(idx, startRowNo, rows)
	if len(moneyTxs) == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "File không hợp lệ. Vui lòng kiểm tra lại.")
	}
	return moneyTxs, nil
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

func parseRows(idx indexes, startRowNo int, rows [][]string) (res []*moneytx.MoneyTransactionShippingExternalLine, _ error) {
	for i := startRowNo; i < len(rows); i++ {
		row := rows[i]
		stt := row[idx.stt]
		if _, err := strconv.Atoi(stt); err != nil {
			continue
		}

		shippingCode := row[idx.shippingCode]
		codAmount, err := strconv.ParseFloat(row[idx.codAmount], 64)
		if err != nil {
			ll.Error("Số tiền COD không hợp lệ", l.Int("Dòng", i+1))
			continue
		}

		res = append(res, &moneytx.MoneyTransactionShippingExternalLine{
			ExternalCode:     shippingCode,
			ExternalTotalCOD: int(codAmount),
		})
	}
	return
}
