package vtpostimport

import (
	"bytes"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"

	identitytypes "o.o/api/main/identity/types"
	"o.o/api/main/moneytx"
	"o.o/api/top/types/etc/shipping_provider"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/imcsv"
	convertpball "o.o/backend/pkg/etop/api/convertpb/_all"
)

type Import struct {
	MoneyTxAggr moneytx.CommandBus
}

func (im *Import) HandleImportMoneyTransactions(c *httpx.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return cm.Errorf(cm.InvalidArgument, nil, "Invalid request")
	}

	files := form.File["files"]
	switch len(files) {
	case 0:
		return cm.Errorf(cm.InvalidArgument, nil, "No file")
	case 1:
		// continue
	default:
		return cm.Errorf(cm.InvalidArgument, nil, "Too many files")
	}

	file, err := files[0].Open()
	if err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "Can not read file")
	}
	defer file.Close()

	provider := form.Value["provider"]
	externalPaidAtStr := form.Value["external_paid_at"]
	note := imcsv.GetFormValue(form.Value["note"])
	accountNumber := imcsv.GetFormValue(form.Value["account_number"])
	accountName := imcsv.GetFormValue(form.Value["account_name"])
	bankName := imcsv.GetFormValue(form.Value["bank_name"])
	invoiceNumber := imcsv.GetFormValue(form.Value["invoice_number"])

	if provider == nil || provider[0] == "" {
		return cm.Error(cm.InvalidArgument, "Missing Provider", nil)
	}
	shippingProvider, ok := shipping_provider.ParseShippingProvider(provider[0])
	if !ok {
		return cm.Errorf(cm.InvalidArgument, nil, "invalid carrier %v", provider[0])
	}

	if shippingProvider != shipping_provider.VTPost {
		return cm.Error(cm.InvalidArgument, "Đơn vị vận chuyển phải là vtpost.", nil)
	}

	var externalPaidAt time.Time
	if externalPaidAtStr != nil {
		externalPaidAt, err = time.Parse(time.RFC3339, externalPaidAtStr[0])
		if err != nil {
			return cm.Error(cm.InvalidArgument, "externalPaidAt is invalid! Use format: `2018-07-17T09:25:13.193Z`", err)
		}
	}

	rawData, err := ioutil.ReadAll(file)
	if err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "Không thể đọc được file. Vui lòng kiểm tra lại hoặc liên hệ %v.", wl.X(c.Context()).CSEmail).WithMeta("reason", "can not open file")
	}

	excelFile, err := excelize.OpenReader(bytes.NewReader(rawData))
	if err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "Không thể đọc được file. Vui lòng kiểm tra lại hoặc liên hệ %v.", wl.X(c.Context()).CSEmail).WithMeta("reason", "invalid file format")
	}
	sheetName := excelFile.GetSheetName(1)
	rows := excelFile.GetRows(sheetName)
	if len(rows) <= 1 {
		return cm.Errorf(cm.InvalidArgument, nil, "File không có nội dung. Vui lòng tải lại file import hoặc liên hệ %v.", wl.X(c.Context()).CSEmail).WithMeta("reason", "no rows")
	}

	var shippingLines []*VTPostMoneyTransactionShippingExternalLine
	for _, row := range rows {
		line, err := parseRow(row)
		if line != nil && err == nil {
			shippingLines = append(shippingLines, line)
		}
	}
	if len(shippingLines) == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "File không có nội dung. Vui lòng tải lại file import hoặc liên hệ %v.", wl.X(c.Context()).CSEmail).WithMeta("reason", "no rows")
	}
	ctx := bus.Ctx()

	cmd := &moneytx.CreateMoneyTxShippingExternalCommand{
		Provider:       shippingProvider,
		ExternalPaidAt: externalPaidAt,
		Lines:          ToMoneyTransactionShippingExternalLines(shippingLines),
		Note:           note,
		InvoiceNumber:  invoiceNumber,
		BankAccount: &identitytypes.BankAccount{
			Name:          bankName,
			AccountNumber: accountNumber,
			AccountName:   accountName,
		},
	}
	if err := im.MoneyTxAggr.Dispatch(ctx, cmd); err != nil {
		return cm.Error(cm.InvalidArgument, "unexpected error", err)
	}
	c.SetResult(convertpball.PbMoneyTxShippingExternalFtLine(cmd.Result))
	return nil
}

func parseRow(row []string) (*VTPostMoneyTransactionShippingExternalLine, error) {
	externalCode := row[2]
	totalStr := row[5]
	if externalCode == "" || totalStr == "" {
		return nil, cm.Error(cm.InvalidArgument, "Row has wrong format", nil).WithMetap("row", row)
	}

	deliveredAtStr := row[6]
	deliveredAt, err := parseDateTime(deliveredAtStr)
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "DeliveredAt is invalid!").WithMetap("row", row)
	}

	totalCOD, err := strconv.ParseFloat(row[3], 64)
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "TotalCOD is invalid!").WithMetap("total COD", row[3]).WithMetap("row", row)
	}
	totalshippingFee, err := strconv.ParseFloat(row[4], 64)
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Total ShippingFee is invalid!").WithMetap("Total ShippingFee", row[4]).WithMetap("row", row)
	}
	total, err := strconv.ParseFloat(totalStr, 64)
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Total is invalid!").WithMetap("row", row)
	}

	return &VTPostMoneyTransactionShippingExternalLine{
		ExternalCode:     externalCode,
		DeliveredAt:      deliveredAt,
		TotalCOD:         int(totalCOD),
		TotalShippingFee: int(totalshippingFee),
		Total:            int(total),
	}, nil
}

func parseDateTime(dateTimeStr string) (res time.Time, err error) {
	dateTimeStr = strings.TrimSpace(dateTimeStr)
	for _, layout := range dateTimeLayouts {
		res, err = time.ParseInLocation(layout, dateTimeStr, time.Local)
		if err == nil {
			return
		}
	}
	return
}
