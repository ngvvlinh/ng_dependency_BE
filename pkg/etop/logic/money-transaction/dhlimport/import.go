package dhlimport

import (
	"context"
	"mime/multipart"
	"strconv"

	"o.o/api/main/moneytx"
	"o.o/api/main/shipping"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/imcsv"
	logicmoneytx "o.o/backend/pkg/etop/logic/money-transaction"
	moneytxtypes "o.o/backend/pkg/etop/logic/money-transaction/handlers/types"
	"o.o/capi/dot"
	"o.o/common/l"
)

var _ moneytxtypes.ImportMoneyTx = &DHLImporter{}
var ll = l.New()

type DHLImporter struct {
	ShippingQuery shipping.QueryBus
}

func (D DHLImporter) ValidateAndReadFile(
	ctx context.Context, fileType string, file multipart.File,
) ([]*moneytx.MoneyTransactionShippingExternalLine, error) {
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

func parseRows(idx imcsv.Indexer, startRowNo int, rows [][]string) (res []*moneytx.MoneyTransactionShippingExternalLine, _ error) {
	for i := startRowNo; i < len(rows); i++ {
		row := rows[i]
		shipmentID := idx.GetCell(row, idxShipmentID)
		if shipmentID == "" {
			continue
		}
		ffmID, err := strconv.ParseInt(shipmentID, 10, 64)
		if err != nil {
			continue
		}

		shippingCode := idx.GetCell(row, idxShippingCode)
		if shippingCode == "" {
			continue
		}

		codAmount, err := strconv.ParseFloat(idx.GetCell(row, idxCodAmount), 64)
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

func validateSchema(ctx context.Context, rows [][]string) (startRowNo int, idx imcsv.Indexer, err error) {
	for i, row := range rows {
		indexer, errs, _ := schema.ValidateSchema(ctx, &row)
		if len(errs) > 0 {
			continue
		}
		idx = indexer
		startRowNo = i + 1
		return
	}
	return 0, idx, cm.Errorf(cm.FailedPrecondition, nil, "File import không hợp lệ")
}
