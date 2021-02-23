package handlers

import (
	"context"
	"strings"
	"time"

	"o.o/api/main/connectioning"
	identitytypes "o.o/api/main/identity/types"
	"o.o/api/main/moneytx"
	"o.o/api/top/types/etc/connection_type"
	"o.o/api/top/types/etc/shipping_provider"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/imcsv"
	"o.o/backend/pkg/common/validate"
	convertpball "o.o/backend/pkg/etop/api/convertpb/_all"
	"o.o/backend/pkg/etop/logic/money-transaction/dhlimport"
	"o.o/backend/pkg/etop/logic/money-transaction/ghnimport"
	"o.o/backend/pkg/etop/logic/money-transaction/ghtkimport"
	moneytxtypes "o.o/backend/pkg/etop/logic/money-transaction/handlers/types"
	"o.o/backend/pkg/etop/logic/money-transaction/jtexpressimport"
	"o.o/backend/pkg/etop/logic/money-transaction/njvimport"
	"o.o/backend/pkg/etop/logic/money-transaction/snappyimport"
	"o.o/backend/pkg/etop/logic/money-transaction/vtpostimport"
	"o.o/capi/dot"
)

type ImportService struct {
	MoneyTxAggr     moneytx.CommandBus
	ConnectionQuery connectioning.QueryBus

	VTPostImporter    *vtpostimport.VTPostImporter
	GHTKImporter      *ghtkimport.GHTKImporter
	GHNImporter       *ghnimport.GHNImporter
	JTExpressImporter *jtexpressimport.JTImporter
	DHLImporter       *dhlimport.DHLImporter
	NJVImporter       *njvimport.NJVImporter
	SnappyImporter    *snappyimport.SnappyImporter
}

var (
	JTExpressNameNormContains     = []string{"jt", "j t", "jtexpress"}
	SnappyExpressNameNormContains = []string{"snappy"}
)

func (s *ImportService) HandleImportMoneyTxs(c *httpx.Context) error {
	form, err := c.MultipartForm()
	ctx := bus.Ctx()
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

	connectionIDStr := imcsv.GetFormValue(form.Value["connection_id"])
	connectionID, err := dot.ParseID(connectionIDStr)
	if err != nil {
		return cm.Errorf(cm.InvalidArgument, nil, "Connection ID does not valid")
	}

	externalPaidAtStr := imcsv.GetFormValue(form.Value["external_paid_at"])
	var externalPaidAt time.Time
	if externalPaidAtStr != "" {
		externalPaidAt, err = time.Parse(time.RFC3339, externalPaidAtStr)
		if err != nil {
			return cm.Error(cm.InvalidArgument, "externalPaidAt is invalid! Use format: `2018-07-17T09:25:13.193Z`", err)
		}
	}

	note := GetFormValue(form.Value["note"])
	accountNumber := GetFormValue(form.Value["account_number"])
	accountName := GetFormValue(form.Value["account_name"])
	bankName := GetFormValue(form.Value["bank_name"])
	invoiceNumber := GetFormValue(form.Value["invoice_number"])

	driver, carrier, err := s.getCarrierImporter(ctx, connectionID)
	if err != nil {
		return err
	}
	lines, err := driver.ValidateAndReadFile(ctx, fileTypes[0], file)
	if err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "File không đúng định dạng")
	}
	cmd := &moneytx.CreateMoneyTxShippingExternalCommand{
		Provider:       carrier,
		ConnectionID:   connectionID,
		ExternalPaidAt: externalPaidAt,
		Lines:          lines,
		Note:           note,
		InvoiceNumber:  invoiceNumber,
		BankAccount: &identitytypes.BankAccount{
			Name:          bankName,
			AccountNumber: accountNumber,
			AccountName:   accountName,
		},
	}

	if err = s.MoneyTxAggr.Dispatch(ctx, cmd); err != nil {
		return cm.Error(cm.InvalidArgument, "unexpected error", err)
	}
	c.SetResult(convertpball.PbMoneyTxShippingExternalFtLine(cmd.Result))
	return nil
}

func (s *ImportService) getCarrierImporter(ctx context.Context, connectionID dot.ID) (importer moneytxtypes.ImportMoneyTx, carrier shipping_provider.ShippingProvider, _ error) {
	query := &connectioning.GetConnectionByIDQuery{
		ID: connectionID,
	}
	if err := s.ConnectionQuery.Dispatch(ctx, query); err != nil {
		return nil, 0, err
	}
	conn := query.Result
	if conn.ConnectionMethod != connection_type.ConnectionMethodBuiltin {
		return nil, 0, cm.Errorf(cm.InvalidArgument, nil, "Connection ID does not valid")
	}

	switch conn.ConnectionProvider {
	case connection_type.ConnectionProviderGHN:
		return s.GHNImporter, shipping_provider.GHN, nil
	case connection_type.ConnectionProviderGHTK:
		return s.GHTKImporter, shipping_provider.GHTK, nil
	case connection_type.ConnectionProviderVTP:
		return s.VTPostImporter, shipping_provider.VTPost, nil
	case connection_type.ConnectionProviderNinjaVan:
		return s.NJVImporter, shipping_provider.NinjaVan, nil
	case connection_type.ConnectionProviderDHL:
		return s.DHLImporter, shipping_provider.DHL, nil
	case connection_type.ConnectionProviderPartner:
		nameNorm := validate.NormalizeSearchSimple(query.Result.Name)
		// checkJTExpress
		if checkContains(nameNorm, JTExpressNameNormContains) {
			return s.JTExpressImporter, shipping_provider.Partner, nil
		}

		// check snappy express
		if checkContains(nameNorm, SnappyExpressNameNormContains) {
			return s.SnappyImporter, shipping_provider.Partner, nil
		}

		// not found
		return nil, 0, cm.Errorf(cm.InvalidArgument, nil, "Connection ID does not valid")
	default:
		return nil, 0, cm.Errorf(cm.InvalidArgument, nil, "Connection ID does not valid")
	}
}

func GetFormValue(ss []string) string {
	if ss == nil {
		return ""
	}
	return ss[0]
}

func checkContains(connNameNorm string, contains []string) bool {
	for _, name := range contains {
		if strings.Contains(connNameNorm, name) {
			return true
		}
	}
	return false
}
