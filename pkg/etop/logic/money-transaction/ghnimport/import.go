package ghnimport

import (
	"context"
	"mime/multipart"

	"o.o/api/main/moneytx"
	cm "o.o/backend/pkg/common"
	logicmoneytx "o.o/backend/pkg/etop/logic/money-transaction"
	moneytxtypes "o.o/backend/pkg/etop/logic/money-transaction/handlers/types"
)

var _ moneytxtypes.ImportMoneyTx = &GHNImporter{}

type GHNImporter struct{}

func (i *GHNImporter) ValidateAndReadFile(ctx context.Context, fileType string, file multipart.File) (lines []*moneytx.MoneyTransactionShippingExternalLine, err error) {
	var rows [][]string
	switch fileType {
	case xlsxFileType:
		rows, err = logicmoneytx.ReadXLSXFile(ctx, file)
	case csvFileType:
		rows, err = logicmoneytx.ReadCSVFile(ctx, file, nColumn)
	default:
		return nil, cm.Errorf(cm.InvalidArgument, nil, "File không đúng định dạng. Chỉ hỗ trợ file csv, xlsx")
	}
	if err != nil {
		return nil, err
	}

	schema, idx, _errs, err := validateSchema(ctx, &rows[0])
	if err != nil {
		return nil, err
	}
	if len(_errs) > 0 {
		return nil, _errs[0]
	}
	rowMoneyTxes, err := parseRows(schema, idx, rows)
	if err != nil {
		return nil, err
	}
	for _, r := range rowMoneyTxes {
		line, err := r.ToModel()
		if err != nil {
			return nil, err
		}
		lines = append(lines, line)
	}
	return
}
