package vtpostimport

import (
	"context"
	"mime/multipart"

	"o.o/api/main/moneytx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	logicmoneytx "o.o/backend/pkg/etop/logic/money-transaction"
	moneytxtypes "o.o/backend/pkg/etop/logic/money-transaction/handlers/types"
)

var _ moneytxtypes.ImportMoneyTx = &VTPostImporter{}

type VTPostImporter struct{}

func (i *VTPostImporter) ValidateAndReadFile(ctx context.Context, fileType string, file multipart.File) ([]*moneytx.MoneyTransactionShippingExternalLine, error) {
	rows, err := logicmoneytx.ReadXLSXFile(ctx, file)
	if err != nil {
		return nil, err
	}

	var shippingLines []*VTPostMoneyTransactionShippingExternalLine
	for _, row := range rows {
		line, err := parseRow(row)
		if line != nil && err == nil {
			shippingLines = append(shippingLines, line)
		}
	}
	if len(shippingLines) == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "File không có nội dung. Vui lòng tải lại file import hoặc liên hệ %v.", wl.X(ctx).CSEmail).WithMeta("reason", "no rows")
	}
	return ToMoneyTransactionShippingExternalLines(shippingLines), nil
}
