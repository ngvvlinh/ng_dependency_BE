package imcsv

import (
	"bytes"
	"context"
	"io/ioutil"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/valyala/tsvreader"

	identitysharemodel "etop.vn/backend/com/main/identity/sharemodel"
	txmodel "etop.vn/backend/com/main/moneytx/model"
	txmodelx "etop.vn/backend/com/main/moneytx/modelx"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/httpx"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/imcsv"
	"etop.vn/backend/pkg/etop/api/convertpb"
)

const (
	xlsxFileType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	csvFileType  = "text/csv"
)

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

	fileTypes := files[0].Header["Content-Type"]
	file, err := files[0].Open()
	if err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "Can not read file")
	}
	defer func() { _ = file.Close() }()

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
	var externalPaidAt time.Time
	if externalPaidAtStr != nil {
		externalPaidAt, err = time.Parse(time.RFC3339, externalPaidAtStr[0])
		if err != nil {
			return cm.Error(cm.InvalidArgument, "externalPaidAt is invalid! Use format: `2018-07-17T09:25:13.193Z`", err)
		}
	}

	var lines []*txmodel.MoneyTransactionShippingExternalLine
	var rows [][]string
	fileType := fileTypes[0]
	switch fileType {
	case xlsxFileType:
		rows, err = ReadXLSXFile(file)
	case csvFileType:
		rows, err = ReadCSVFile(file)
	default:
		return cm.Errorf(cm.InvalidArgument, nil, "File không đúng định dạng")
	}

	if err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "File không đúng định dạng")
	}

	schema, idx, _errs, err := validateSchema(c.Context(), &rows[0])
	if err != nil {
		return err
	}
	if len(_errs) > 0 {
		return _errs[0]
	}
	rowMoneyTxes, err := parseRows(schema, idx, rows)
	if err != nil {
		return err
	}
	for _, r := range rowMoneyTxes {
		line, err := r.ToModel()
		if err != nil {
			return err
		}
		lines = append(lines, line)
	}

	cmd := &txmodelx.CreateMoneyTransactionShippingExternal{
		Provider:       provider[0],
		ExternalPaidAt: externalPaidAt,
		Lines:          lines,
		Note:           note,
		InvoiceNumber:  invoiceNumber,
		BankAccount: &identitysharemodel.BankAccount{
			Name:          bankName,
			AccountNumber: accountNumber,
			AccountName:   accountName,
		},
	}

	ctx := context.Background()
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return cm.Error(cm.InvalidArgument, "unexpected error", err)
	}
	c.SetResult(convertpb.PbMoneyTransactionShippingExternalExtended(cmd.Result))
	return nil
}

func ReadCSVFile(file multipart.File) (rows [][]string, _ error) {
	r := tsvreader.New(file)
	for r.Next() {
		var row []string
		for i := 0; i < nColumn; i++ {
			row = append(row, r.String())
		}
		r.SkipCol()
		rows = append(rows, row)
	}
	if err := r.Error(); err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "unexpected error")
	}
	return
}

func ReadXLSXFile(file multipart.File) (rows [][]string, _ error) {
	rawData, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Không thể đọc được file. Vui lòng kiểm tra lại hoặc liên hệ hotro@etop.vn.").WithMeta("reason", "can not open file")
	}

	excelFile, err := excelize.OpenReader(bytes.NewReader(rawData))
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Không thể đọc được file. Vui lòng kiểm tra lại hoặc liên hệ hotro@etop.vn.").WithMeta("reason", "invalid file format")
	}
	sheetName := excelFile.GetSheetName(1)
	rows = excelFile.GetRows(sheetName)
	if len(rows) <= 1 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "File không có nội dung. Vui lòng tải lại file import hoặc liên hệ hotro@etop.vn.").WithMeta("reason", "no rows")
	}
	return
}

func parseRows(schema imcsv.Schema, idx indexes, rows [][]string) (res []*RowMoneyTransaction, _ error) {
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) < len(schemas) {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Số cột không đúng cấu trúc yêu cầu.").WithMetap("row", row).WithMetap("index", i)
		}
		rowMoneyTx := &RowMoneyTransaction{
			ExCode:    row[idx.exCode],
			EtopCode:  row[idx.etopCode],
			CreatedAt: row[idx.createdAt],
			ClosedAt:  row[idx.closedAt],
			Customer:  row[idx.customer],
			Address:   row[idx.address],
			TotalCOD:  row[idx.totalCOD],
		}
		res = append(res, rowMoneyTx)
	}
	return
}

func (m *RowMoneyTransaction) ToModel() (*txmodel.MoneyTransactionShippingExternalLine, error) {
	layout := "01/02/06 15:04"
	createdAt, err := time.ParseInLocation(layout, strings.TrimSpace(m.CreatedAt), time.Local)
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "CreatedAt is invalid!").WithMetap("row", m)
	}
	closedAt, err := time.ParseInLocation(layout, strings.TrimSpace(m.ClosedAt), time.Local)
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "UpdatedAt is invalid!").WithMetap("row", m)
	}

	totalCOD, err := strconv.ParseFloat(m.TotalCOD, 64)
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "TotalCOD is invalid!").WithMetap("total COD", m.TotalCOD).WithMetap("row", m)
	}
	return &txmodel.MoneyTransactionShippingExternalLine{
		ExternalCode:         m.ExCode,
		EtopFulfillmentIDRaw: m.EtopCode,
		ExternalCreatedAt:    createdAt,
		ExternalClosedAt:     closedAt,
		ExternalCustomer:     m.Customer,
		ExternalAddress:      m.Address,
		ExternalTotalCOD:     int(totalCOD),
	}, nil
}
