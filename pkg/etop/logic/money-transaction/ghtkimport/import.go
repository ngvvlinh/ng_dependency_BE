package ghtkimport

import (
	"context"
	"mime/multipart"

	"o.o/api/main/moneytx"
	"o.o/api/main/shipping"
	"o.o/api/top/types/etc/shipping_fee_type"
	"o.o/api/top/types/etc/status5"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	logicmoneytx "o.o/backend/pkg/etop/logic/money-transaction"
	moneytxtypes "o.o/backend/pkg/etop/logic/money-transaction/handlers/types"
	"o.o/capi/dot"
)

var _ moneytxtypes.ImportMoneyTx = &GHTKImporter{}

type GHTKImporter struct {
	ShippingAggr  shipping.CommandBus
	ShippingQuery shipping.QueryBus
}

func (i *GHTKImporter) ValidateAndReadFile(ctx context.Context, fileType string, file multipart.File) ([]*moneytx.MoneyTransactionShippingExternalLine, error) {
	rows, err := logicmoneytx.ReadXLSXFile(ctx, file)
	if err != nil {
		return nil, err
	}

	headerIndexMap := getHeaderIndex(rows)
	if err := checkHeaderIndex(headerIndexMap); err != nil {
		return nil, err
	}

	var shippingLines []*GHTKMoneyTransactionShippingExternalLine
	for _, row := range rows {
		line, err := parseRow(row, headerIndexMap)
		if line != nil && err == nil {
			shippingLines = append(shippingLines, line)
		}
	}
	if len(shippingLines) == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "File không có nội dung. Vui lòng tải lại file import hoặc liên hệ %v.", wl.X(ctx).CSEmail).WithMeta("reason", "no rows")
	}

	// update Fulfillments shipping fee (insurance, return, discount, change address)
	fulfillments, err := i.updateShippingFeeFulfillmentsFromImportFile(ctx, shippingLines)
	if err != nil {
		return nil, err
	}
	if err := i.updateShippingFeeFulfillments(ctx, fulfillments); err != nil {
		return nil, err
	}
	return ToMoneyTransactionShippingExternalLines(shippingLines), nil
}

type FulfillmentUpdate struct {
	ID                       dot.ID
	ProviderShippingFeeLines []*shipping.ShippingFeeLine
}

func (i *GHTKImporter) updateShippingFeeFulfillments(ctx context.Context, ffms []*FulfillmentUpdate) error {
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

func (i *GHTKImporter) updateShippingFeeFulfillmentsFromImportFile(ctx context.Context, lines []*GHTKMoneyTransactionShippingExternalLine) ([]*FulfillmentUpdate, error) {
	ffmShippingCodes := make([]string, len(lines))
	for i, line := range lines {
		ffmShippingCodes[i] = line.ExternalCode
	}
	query := &shipping.ListFulfillmentsByShippingCodesQuery{
		Codes: ffmShippingCodes,
	}
	if err := i.ShippingQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}

	updatesMap := make(map[string]*FulfillmentUpdate, len(query.Result))
	ffmsByShippingCode := make(map[string]*shipping.Fulfillment, len(query.Result))
	for _, ffm := range query.Result {
		// ignore ffms that finished
		if (ffm.Status != status5.Z && ffm.Status != status5.S) ||
			!ffm.CODEtopTransferedAt.IsZero() {
			continue
		}

		feeLines := ffm.ProviderShippingFeeLines
		var newFeeLines []*shipping.ShippingFeeLine
		for _, feeLine := range feeLines {
			if feeLine.ShippingFeeType == shipping_fee_type.Main {
				// keep the shipping fee type main (phí dịch vụ)
				newFeeLines = append(newFeeLines, feeLine)
				break
			}
		}
		updatesMap[ffm.ShippingCode] = &FulfillmentUpdate{
			ID:                       ffm.ID,
			ProviderShippingFeeLines: newFeeLines,
		}
		ffmsByShippingCode[ffm.ShippingCode] = ffm
	}
	for _, line := range lines {
		update := updatesMap[line.ExternalCode]
		ffm := ffmsByShippingCode[line.ExternalCode]
		if update == nil || ffm == nil {
			continue
		}
		if line.InsuranceFee != 0 {
			update.ProviderShippingFeeLines = append(update.ProviderShippingFeeLines, &shipping.ShippingFeeLine{
				ShippingFeeType: shipping_fee_type.Insurance,
				Cost:            line.InsuranceFee,
			})
		}
		if line.ReturnFee != 0 {
			update.ProviderShippingFeeLines = append(update.ProviderShippingFeeLines, &shipping.ShippingFeeLine{
				ShippingFeeType: shipping_fee_type.Return,
				Cost:            line.ReturnFee,
			})
		}
		if line.Discount != 0 {
			cost := line.Discount
			if cost > 0 {
				cost = -cost
			}
			update.ProviderShippingFeeLines = append(update.ProviderShippingFeeLines, &shipping.ShippingFeeLine{
				Cost:            cost,
				ShippingFeeType: shipping_fee_type.Discount,
			})
		}
		if line.ChangeAddressFee != 0 {
			update.ProviderShippingFeeLines = append(update.ProviderShippingFeeLines, &shipping.ShippingFeeLine{
				ShippingFeeType: shipping_fee_type.AddressChange,
				Cost:            line.ChangeAddressFee,
			})
		}

		updatesMap[line.ExternalCode] = update
	}
	fulfillments := make([]*FulfillmentUpdate, 0, len(updatesMap))
	for _, ffm := range updatesMap {
		fulfillments = append(fulfillments, ffm)
	}
	return fulfillments, nil
}
