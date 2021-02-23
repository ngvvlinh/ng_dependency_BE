package snappyimport

import (
	"context"
	"mime/multipart"
	"strconv"

	"o.o/api/main/moneytx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/imcsv"
	logicmoneytx "o.o/backend/pkg/etop/logic/money-transaction"
	moneytxtypes "o.o/backend/pkg/etop/logic/money-transaction/handlers/types"
	"o.o/capi/dot"
)

var _ moneytxtypes.ImportMoneyTx = &SnappyImporter{}

type SnappyImporter struct{}

func (s *SnappyImporter) ValidateAndReadFile(ctx context.Context, fileType string, file multipart.File) ([]*moneytx.MoneyTransactionShippingExternalLine, error) {
	rows, err := logicmoneytx.ReadXLSXFile(ctx, file)
	if err != nil {
		return nil, err
	}
	startRowNo, idx, err := validateSchema(ctx, rows)
	if err != nil {
		return nil, err
	}

	moneyTxs, err := parseRows(idx, startRowNo, rows)
	if err != nil || len(moneyTxs) == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "File không hợp lệ. Vui lòng kiểm tra lại.")
	}
	return moneyTxs, nil
}

func validateSchema(ctx context.Context, rows [][]string) (startRowNo int, idx imcsv.Indexer, err error) {
	for i, row := range rows {
		_indexer, errs, _ := schema.ValidateSchema(ctx, &row)
		if len(errs) > 0 {
			continue
		}
		idx = _indexer
		startRowNo = i + 1
		return
	}
	return 0, idx, cm.Errorf(cm.FailedPrecondition, nil, "File import không hợp lệ")
}

func parseRows(idx imcsv.Indexer, startRowNo int, rows [][]string) (res []*moneytx.MoneyTransactionShippingExternalLine, _ error) {
	for i := startRowNo; i < len(rows); i++ {
		row := rows[i]

		shippingCode := idx.GetCell(row, idxShippingCode)
		if shippingCode == "" {
			continue
		}

		shipmentID := idx.GetCell(row, idxShipmentID)
		var ffmID int64
		var err error
		if shipmentID != "" {
			ffmID, err = strconv.ParseInt(shipmentID, 10, 64)
			if err != nil {
				continue
			}
		}

		codAmount, err := strconv.ParseFloat(idx.GetCell(row, idxCODAmount), 64)
		if err != nil {
			continue
		}

		res = append(res, &moneytx.MoneyTransactionShippingExternalLine{
			EtopFulfillmentID: dot.ID(ffmID),
			ExternalCode:      shippingCode,
			ExternalTotalCOD:  int(codAmount),
		})
	}
	return
}
