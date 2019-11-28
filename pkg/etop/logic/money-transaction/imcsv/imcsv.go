package imcsv

import (
	"context"
	"strings"
	"time"

	"etop.vn/backend/pkg/common/imcsv"

	"github.com/valyala/tsvreader"

	txmodel "etop.vn/backend/com/main/moneytx/model"
	txmodelx "etop.vn/backend/com/main/moneytx/modelx"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/httpx"
	"etop.vn/backend/pkg/etop/api/convertpb"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/common/xerrors"
)

type Error struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`

	Meta map[string]interface{} `json:"meta,omitempty"`
}

func (e *Error) Error() string {
	return e.Msg
}

func NewError(code xerrors.Code, msg string) error {
	err := &Error{
		Code: code.String(),
		Msg:  msg,
	}
	return err
}

type MoneyTransactionShippingExternalLine struct {
	ExCode      string
	CreatedAt   time.Time
	DeliveredAt time.Time
	Customer    string
	Address     string
	TotalCOD    int
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
	defer func() { _ = file.Close() }()
	r := tsvreader.New(file)

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

	transactionLines := make(map[string]*txmodel.MoneyTransactionShippingExternalLine)
	row := 0
	for r.Next() {
		row++
		if row == 1 {
			r.SkipCol()
			r.SkipCol()
			r.SkipCol()
			r.SkipCol()
			r.SkipCol()
			r.SkipCol()
			r.SkipCol()
			r.SkipCol()
			continue
		}
		code := r.String()
		etopOrderCode := r.String()
		createdAtStr := r.String()
		closedAtStr := r.String()
		customer := r.String()
		address := r.String()
		totalCOD := r.Int()
		r.SkipCol()

		layout := "01/02/06 15:04"
		createdAt, err := time.ParseInLocation(layout, strings.TrimSpace(createdAtStr), time.Local)
		if err != nil {
			return cm.Errorf(cm.InvalidArgument, err, "CreatedAt is invalid!")
		}
		closedAt, err := time.ParseInLocation(layout, strings.TrimSpace(closedAtStr), time.Local)
		if err != nil {
			return cm.Errorf(cm.InvalidArgument, err, "UpdatedAt is invalid!")
		}

		transactionLines[code] = &txmodel.MoneyTransactionShippingExternalLine{
			ExternalCode:         code,
			EtopFulfillmentIdRaw: etopOrderCode,
			ExternalCreatedAt:    createdAt,
			ExternalClosedAt:     closedAt,
			ExternalCustomer:     customer,
			ExternalAddress:      address,
			ExternalTotalCOD:     totalCOD,
		}
	}
	if err := r.Error(); err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "unexpected error")
	}
	lines := make([]*txmodel.MoneyTransactionShippingExternalLine, 0, len(transactionLines))
	for _, line := range transactionLines {
		lines = append(lines, line)
	}

	cmd := &txmodelx.CreateMoneyTransactionShippingExternal{
		Provider:       provider[0],
		ExternalPaidAt: externalPaidAt,
		Lines:          lines,
		Note:           note,
		InvoiceNumber:  invoiceNumber,
		BankAccount: &model.BankAccount{
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
