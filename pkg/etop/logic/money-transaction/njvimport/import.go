package njvimport

import (
	"context"
	"mime/multipart"
	"strconv"

	"o.o/api/main/moneytx"
	"o.o/api/main/shipping"
	shippingtypes "o.o/api/main/shipping/types"
	"o.o/api/top/types/etc/shipping_fee_type"
	"o.o/api/top/types/etc/status5"
	cm "o.o/backend/pkg/common"
	logicmoneytx "o.o/backend/pkg/etop/logic/money-transaction"
	moneytxtypes "o.o/backend/pkg/etop/logic/money-transaction/handlers/types"
	"o.o/capi/dot"
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

	fulfillments, err := i.updateShippingFeeFulfillmentsFromImportFile(ctx, shippingLines)
	if err != nil {
		return nil, err
	}
	if err := i.updateShippingFeeFulfillments(ctx, fulfillments); err != nil {
		return nil, err
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
		exCode := row[idx.trackingID]
		if exCode == "" {
			continue
		}

		mainFee, err := strconv.ParseFloat(row[idx.mainFee], 64)
		if err != nil {
			continue
		}

		insFee, err := strconv.ParseFloat(row[idx.insFee], 64)
		if err != nil {
			continue
		}

		rtsFee, err := strconv.ParseFloat(row[idx.rtsFee], 64)
		if err != nil {
			continue
		}

		vatFee, err := strconv.ParseFloat(row[idx.vatFee], 64)
		if err != nil {
			continue
		}

		codFee, err := strconv.ParseFloat(row[idx.codFee], 64)
		if err != nil {
			continue
		}

		codAmount, err := strconv.ParseFloat(row[idx.codAmount], 64)
		if err != nil {
			continue
		}

		res = append(res, &RowMoneyTransaction{
			ExCode:    exCode,
			MainFee:   mainFee,
			RtsFee:    rtsFee,
			InsFee:    insFee,
			VatFee:    vatFee,
			CodFee:    codFee,
			CodAmount: codAmount,
		})
	}
	return
}

type FulfillmentUpdate struct {
	ID                       dot.ID
	ProviderShippingFeeLines []*shippingtypes.ShippingFeeLine
}

func (i *NJVImporter) updateShippingFeeFulfillments(ctx context.Context, ffms []*FulfillmentUpdate) error {
	for _, ffm := range ffms {
		cmd := &shipping.UpdateFulfillmentShippingFeesCommand{
			FulfillmentID:            ffm.ID,
			ProviderShippingFeeLines: ffm.ProviderShippingFeeLines,
		}
		if err := i.ShippingAggr.Dispatch(ctx, cmd); err != nil {
			return err
		}
	}
	return nil
}

func (i *NJVImporter) updateShippingFeeFulfillmentsFromImportFile(ctx context.Context, lines []*RowMoneyTransaction) ([]*FulfillmentUpdate, error) {
	ffmShippingCodes := make([]string, len(lines))
	for i, line := range lines {
		ffmShippingCodes[i] = line.ExCode
	}
	query := &shipping.ListFulfillmentsByShippingCodesQuery{
		Codes: ffmShippingCodes,
	}
	if err := i.ShippingQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}

	updatesMap := make(map[string]*FulfillmentUpdate, len(query.Result))
	for _, ffm := range query.Result {
		// ignore ffms that finished
		if (ffm.Status != status5.Z && ffm.Status != status5.S) ||
			!ffm.CODEtopTransferedAt.IsZero() {
			continue
		}

		updatesMap[ffm.ShippingCode] = &FulfillmentUpdate{
			ID:                       ffm.ID,
			ProviderShippingFeeLines: []*shippingtypes.ShippingFeeLine{},
		}
	}

	// multiple with 1.1
	// 1.1 = (100% + 10%(VAT)) / 100%
	for _, line := range lines {
		update := updatesMap[line.ExCode]
		if update == nil {
			continue
		}

		if line.MainFee != 0 {
			update.ProviderShippingFeeLines = append(update.ProviderShippingFeeLines, &shippingtypes.ShippingFeeLine{
				ShippingFeeType: shipping_fee_type.Main,
				Cost:            int(line.MainFee * 1.1),
			})
		}

		if line.CodFee != 0 {
			update.ProviderShippingFeeLines = append(update.ProviderShippingFeeLines, &shippingtypes.ShippingFeeLine{
				ShippingFeeType: shipping_fee_type.Cods,
				Cost:            int(line.CodFee * 1.1),
			})
		}

		if line.InsFee != 0 {
			update.ProviderShippingFeeLines = append(update.ProviderShippingFeeLines, &shippingtypes.ShippingFeeLine{
				ShippingFeeType: shipping_fee_type.Insurance,
				Cost:            int(line.InsFee * 1.1),
			})
		}

		if line.RtsFee != 0 {
			update.ProviderShippingFeeLines = append(update.ProviderShippingFeeLines, &shippingtypes.ShippingFeeLine{
				ShippingFeeType: shipping_fee_type.Return,
				Cost:            int(line.RtsFee * 1.1),
			})
		}

		updatesMap[line.ExCode] = update
	}
	fulfillments := make([]*FulfillmentUpdate, 0, len(updatesMap))
	for _, ffm := range updatesMap {
		fulfillments = append(fulfillments, ffm)
	}
	return fulfillments, nil
}
