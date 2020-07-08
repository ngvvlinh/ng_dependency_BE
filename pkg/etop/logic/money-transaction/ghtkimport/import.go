package ghtkimport

import (
	"context"
	"mime/multipart"

	"o.o/api/main/moneytx"
	"o.o/api/top/types/etc/shipping_provider"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	logicmoneytx "o.o/backend/pkg/etop/logic/money-transaction"
	moneytxtypes "o.o/backend/pkg/etop/logic/money-transaction/handlers/types"
)

var _ moneytxtypes.ImportMoneyTx = &GHTKImporter{}

type GHTKImporter struct{}

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
	fulfillments, err := UpdateShippingFeeFulfillmentsFromImportFile(ctx, shippingLines, shipping_provider.GHTK)
	if err != nil {
		return nil, err
	}
	if err := updateFulfillments(ctx, fulfillments); err != nil {
		return nil, err
	}
	return ToMoneyTransactionShippingExternalLines(shippingLines), nil
}
