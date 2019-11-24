package vtpostimxlsx

import (
	"bytes"
	"context"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"

	txmodel "etop.vn/backend/com/main/moneytx/model"
	txmodelx "etop.vn/backend/com/main/moneytx/modelx"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/httpx"
	"etop.vn/backend/pkg/etop/api/convertpb"
	"etop.vn/backend/pkg/etop/model"
)

type VTPostMoneyTransactionShippingExternalLine struct {
	ExternalCode     string
	DeliveredAt      time.Time
	TotalCOD         int // tiền hàng
	TotalShippingFee int // tổng tiền cước
	Total            int // Số tiền phải trả
}

func (line *VTPostMoneyTransactionShippingExternalLine) ToModel() *txmodel.MoneyTransactionShippingExternalLine {
	return &txmodel.MoneyTransactionShippingExternalLine{
		ExternalCode:             line.ExternalCode,
		ExternalTotalCOD:         line.TotalCOD,
		ExternalClosedAt:         line.DeliveredAt,
		ExternalTotalShippingFee: line.TotalShippingFee,
	}
}

func ToMoneyTransactionShippingExternalLines(lines []*VTPostMoneyTransactionShippingExternalLine) []*txmodel.MoneyTransactionShippingExternalLine {
	if lines == nil {
		return nil
	}
	res := make([]*txmodel.MoneyTransactionShippingExternalLine, len(lines))
	for i, line := range lines {
		res[i] = line.ToModel()
	}
	return res
}

func HandleImportMoneyTransactions(c *httpx.Context) error {
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
	note := cm.GetFormValue(form.Value["note"])
	accountNumber := cm.GetFormValue(form.Value["account_number"])
	accountName := cm.GetFormValue(form.Value["account_name"])
	bankName := cm.GetFormValue(form.Value["bank_name"])
	invoiceNumber := cm.GetFormValue(form.Value["invoice_number"])

	if provider == nil || provider[0] == "" {
		return cm.Error(cm.InvalidArgument, "Missing Provider", nil)
	}
	shippingProvider := model.ShippingProvider(provider[0])
	if shippingProvider != model.TypeVTPost {
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
		return cm.Errorf(cm.InvalidArgument, err, "Không thể đọc được file. Vui lòng kiểm tra lại hoặc liên hệ hotro@etop.vn.").WithMeta("reason", "can not open file")
	}

	excelFile, err := excelize.OpenReader(bytes.NewReader(rawData))
	if err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "Không thể đọc được file. Vui lòng kiểm tra lại hoặc liên hệ hotro@etop.vn.").WithMeta("reason", "invalid file format")
	}
	sheetName := excelFile.GetSheetName(1)
	rows := excelFile.GetRows(sheetName)
	if len(rows) <= 1 {
		return cm.Errorf(cm.InvalidArgument, nil, "File không có nội dung. Vui lòng tải lại file import hoặc liên hệ hotro@etop.vn.").WithMeta("reason", "no rows")
	}

	var shippingLines []*VTPostMoneyTransactionShippingExternalLine
	for _, row := range rows {
		line, err := parseRow(row)
		if line != nil && err == nil {
			shippingLines = append(shippingLines, line)
		}
	}
	if len(shippingLines) == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "File không có nội dung. Vui lòng tải lại file import hoặc liên hệ hotro@etop.vn.").WithMeta("reason", "no rows")
	}
	ctx := context.Background()

	cmd := &txmodelx.CreateMoneyTransactionShippingExternal{
		Provider:       provider[0],
		ExternalPaidAt: externalPaidAt,
		Lines:          ToMoneyTransactionShippingExternalLines(shippingLines),
		Note:           note,
		InvoiceNumber:  invoiceNumber,
		BankAccount: &model.BankAccount{
			Name:          bankName,
			AccountNumber: accountNumber,
			AccountName:   accountName,
		},
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return cm.Error(cm.InvalidArgument, "unexpected error", err)
	}
	c.SetResult(convertpb.PbMoneyTransactionShippingExternalExtended(cmd.Result))
	return nil
}

func parseRow(row []string) (*VTPostMoneyTransactionShippingExternalLine, error) {
	externalCode := row[2]
	totalStr := row[5]
	if externalCode == "" || totalStr == "" {
		return nil, cm.Error(cm.InvalidArgument, "Row has wrong format", nil).WithMetap("row", row)
	}
	layout := "01-02-06"
	deliveredAtStr := row[6]
	deliveredAt, err := time.ParseInLocation(layout, strings.TrimSpace(deliveredAtStr), time.Local)
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
