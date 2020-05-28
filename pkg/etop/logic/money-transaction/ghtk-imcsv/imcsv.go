package ghtkimcsv

import (
	"bytes"
	"context"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"

	identitytypes "o.o/api/main/identity/types"
	"o.o/api/main/moneytx"
	"o.o/api/top/types/etc/shipping_fee_type"
	"o.o/api/top/types/etc/shipping_provider"
	"o.o/api/top/types/etc/status5"
	shipmodel "o.o/backend/com/main/shipping/model"
	shipmodelx "o.o/backend/com/main/shipping/modelx"
	shippingsharemodel "o.o/backend/com/main/shipping/sharemodel"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/imcsv"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/integration/shipping/ghtk"
	"o.o/capi/dot"
)

/*
	// Cấu trúc file excel
	[]string{
  0  "STT",
  1  "Mã đơn hàng",
  2  "Mã đơn hàng shop",
  3  "Thông tin khách hàng",
  4  "Tổng tiền thu hộ",
  5  "Shop trả trước", 											=> bỏ qua
  6  "Phí bảo hiểm",
  7  "Phí dịch vụ",
  8  "Phí đồng kiểm",
  9  "Phí chuyển hoàn",
  10  "Khuyến mãi",
  11  "Phí thay đổi địa chỉ giao",
  12  "Thanh toán",
  13  "Ngày tạo",
  14  "Ngày hoàn thành",
  15  "",
  16  "",
  }
*/

type Import struct {
	MoneyTxAggr moneytx.CommandBus
}

type GHTKMoneyTransactionShippingExternalLine struct {
	ExternalCode     string
	ShopCode         string
	Customer         string
	TotalCOD         int
	InsuranceFee     int
	ShippingFee      int
	ReturnFee        int
	Discount         int
	ChangeAddressFee int
	Total            int // after sub fees
	CreatedAt        time.Time
	DeliveredAt      time.Time
}

func (line *GHTKMoneyTransactionShippingExternalLine) ToModel() *moneytx.MoneyTransactionShippingExternalLine {
	return &moneytx.MoneyTransactionShippingExternalLine{
		ExternalCode:         line.ExternalCode,
		ExternalCustomer:     line.Customer,
		ExternalTotalCOD:     line.TotalCOD,
		ExternalCreatedAt:    line.CreatedAt,
		ExternalClosedAt:     line.DeliveredAt,
		EtopFulfillmentIDRaw: line.ShopCode,
	}
}

func ToMoneyTransactionShippingExternalLines(lines []*GHTKMoneyTransactionShippingExternalLine) []*moneytx.MoneyTransactionShippingExternalLine {
	if lines == nil {
		return nil
	}
	res := make([]*moneytx.MoneyTransactionShippingExternalLine, len(lines))
	for i, line := range lines {
		res[i] = line.ToModel()
	}
	return res
}

const (
	ExternalCode     = "Mã đơn hàng"
	ShopCode         = "Mã đơn hàng shop"
	CustomerInfo     = "Thông tin khách hàng"
	TotalCOD         = "Tổng tiền thu hộ"
	InsuranceFee     = "Phí bảo hiểm"
	ShippingFee      = "Phí dịch vụ"
	ReturnFee        = "Phí chuyển hoàn"
	Discount         = "Khuyến mãi"
	ChangeAddressFee = "Phí thay đổi địa chỉ giao"
	Total            = "Thanh toán"
	CreatedAt        = "Ngày tạo"
	DeliveredAt      = "Ngày hoàn thành"
)

var (
	HeaderStrings = []string{ExternalCode, ShopCode, CustomerInfo, TotalCOD,
		InsuranceFee, ShippingFee, ReturnFee, Discount, ChangeAddressFee, Total, CreatedAt, DeliveredAt}
)

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

	var externalPaidAt time.Time
	if externalPaidAtStr != nil {
		externalPaidAt, err = time.Parse(time.RFC3339, externalPaidAtStr[0])
		if err != nil {
			return cm.Error(cm.InvalidArgument, "externalPaidAt is invalid! Use format: `2018-07-17T09:25:13.193Z`", err)
		}
	}

	rawData, err := ioutil.ReadAll(file)
	if err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "Không thể đọc được file. Vui lòng kiểm tra lại hoặc liên hệ %v", wl.X(c.Context()).CSEmail).WithMeta("reason", "can not open file")
	}

	excelFile, err := excelize.OpenReader(bytes.NewReader(rawData))
	if err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "Không thể đọc được file. Vui lòng kiểm tra lại hoặc liên hệ %v", wl.X(c.Context()).CSEmail).WithMeta("reason", "invalid file format")
	}
	sheetName := excelFile.GetSheetName(1)
	rows := excelFile.GetRows(sheetName)
	if len(rows) <= 1 {
		return cm.Errorf(cm.InvalidArgument, nil, "File không có nội dung. Vui lòng tải lại file import hoặc liên hệ %v", wl.X(c.Context()).CSEmail).WithMeta("reason", "no rows")
	}

	headerIndexMap := getHeaderIndex(rows)
	if err := checkHeaderIndex(headerIndexMap); err != nil {
		return err
	}

	var shippingLines []*GHTKMoneyTransactionShippingExternalLine
	for _, row := range rows {
		line, err := parseRow(row, headerIndexMap)
		if line != nil && err == nil {
			shippingLines = append(shippingLines, line)
		}
	}
	if len(shippingLines) == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "File không có nội dung. Vui lòng tải lại file import hoặc liên hệ %v", wl.X(c.Context()).CSEmail).WithMeta("reason", "no rows")
	}
	ctx := bus.Ctx()
	fulfillments, err := UpdateShippingFeeFulfillmentsFromImportFile(ctx, shippingLines, shippingProvider)
	if err != nil {
		return err
	}
	// update Fulfillments shipping fee (insurance, return, discount, change address)
	if err := updateFulfillments(ctx, fulfillments); err != nil {
		return err
	}

	cmd := &moneytx.CreateMoneyTxShippingExternalCommand{
		Provider:       shippingProvider,
		ExternalPaidAt: externalPaidAt,
		Lines:          ToMoneyTransactionShippingExternalLines(shippingLines),
		BankAccount: &identitytypes.BankAccount{
			Name:          bankName,
			AccountNumber: accountNumber,
			AccountName:   accountName,
		},
		Note:          note,
		InvoiceNumber: invoiceNumber,
	}
	if err := im.MoneyTxAggr.Dispatch(ctx, cmd); err != nil {
		return cm.Error(cm.InvalidArgument, "unexpected error", err)
	}
	c.SetResult(convertpb.PbMoneyTxShippingExternalFtLine(cmd.Result))
	return nil
}

func parseRow(row []string, headerIndexMap map[string]int) (*GHTKMoneyTransactionShippingExternalLine, error) {
	externalCode := row[headerIndexMap[ExternalCode]]
	customer := row[headerIndexMap[CustomerInfo]]
	totalStr := row[headerIndexMap[Total]]
	if externalCode == "" || customer == "" || totalStr == "" {
		return nil, cm.Error(cm.InvalidArgument, "Row has wrong format", nil).WithMetap("row", row)
	}
	layout := "15:04, 02-01-2006"
	createdAtStr := row[headerIndexMap[CreatedAt]]
	deliveredAtStr := row[headerIndexMap[DeliveredAt]]
	createdAt, err := time.ParseInLocation(layout, strings.TrimSpace(createdAtStr), time.Local)
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "CreatedAt is invalid!").WithMetap("row", row)
	}
	deliveredAt, err := time.ParseInLocation(layout, strings.TrimSpace(deliveredAtStr), time.Local)
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "DeliveredAt is invalid!").WithMetap("row", row)
	}

	totalCOD, err := strconv.ParseFloat(row[headerIndexMap[TotalCOD]], 64)
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "TotalCOD is invalid!").WithMetap("total COD", row[headerIndexMap[TotalCOD]]).WithMetap("row", row)
	}
	insuranceFee, err := strconv.ParseFloat(row[headerIndexMap[InsuranceFee]], 64)
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "InsuranceFee is invalid!").WithMetap("row", row)
	}
	shippingFee, err := strconv.ParseFloat(row[headerIndexMap[ShippingFee]], 64)
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "ShippingFee is invalid!").WithMetap("row", row)
	}
	returnFee, err := strconv.ParseFloat(row[headerIndexMap[ReturnFee]], 64)
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "ReturnFee is invalid!").WithMetap("row", row)
	}
	discount, err := strconv.ParseFloat(row[headerIndexMap[Discount]], 64)
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Discount is invalid!").WithMetap("row", row)
	}
	changeAddressFee, err := strconv.ParseFloat(row[headerIndexMap[ChangeAddressFee]], 64)
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "ChangeAddressFee is invalid!").WithMetap("row", row)
	}
	total, err := strconv.ParseFloat(row[headerIndexMap[Total]], 64)
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Total is invalid!").WithMetap("row", row)
	}

	return &GHTKMoneyTransactionShippingExternalLine{
		ExternalCode:     ghtk.NormalizeGHTKCode(externalCode),
		ShopCode:         row[headerIndexMap[ShopCode]],
		Customer:         customer,
		TotalCOD:         int(totalCOD),
		InsuranceFee:     int(insuranceFee),
		ShippingFee:      int(shippingFee),
		ReturnFee:        int(returnFee),
		Discount:         int(discount),
		ChangeAddressFee: int(changeAddressFee),
		Total:            int(total),
		CreatedAt:        createdAt,
		DeliveredAt:      deliveredAt,
	}, nil
}

func getHeaderIndex(rows [][]string) map[string]int {
	HeaderIndexMap := make(map[string]int, len(HeaderStrings))
	found := false
	for _, row := range rows {
		if found {
			break
		}
		for i, value := range row {
			if len(HeaderIndexMap) == len(HeaderStrings) {
				found = true
				break
			}
			for _, title := range HeaderStrings {
				if value == title {
					HeaderIndexMap[value] = i
				}
			}
		}
	}
	return HeaderIndexMap
}

func checkHeaderIndex(headerIndexMap map[string]int) error {
	for _, title := range HeaderStrings {
		if headerIndexMap[title] == 0 {
			return cm.Errorf(cm.InvalidArgument, nil, "Cột không đúng cấu trúc yêu cầu (mong đợi: %v).", title)
		}
	}
	return nil
}

func UpdateShippingFeeFulfillmentsFromImportFile(ctx context.Context, lines []*GHTKMoneyTransactionShippingExternalLine, shippingProvider shipping_provider.ShippingProvider) ([]*shipmodel.Fulfillment, error) {
	if shippingProvider != shipping_provider.GHTK {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Đơn vị vận chuyển phải là GHTK.").WithMeta("shipping_provider", shippingProvider.String())
	}
	ffmShippingCodes := make([]string, len(lines))
	for i, line := range lines {
		ffmShippingCodes[i] = line.ExternalCode
	}
	cmd := &shipmodelx.GetFulfillmentsQuery{
		ShippingCodes: ffmShippingCodes,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	updatesMap := make(map[string]*shipmodel.Fulfillment, len(cmd.Result.Fulfillments))
	ffmsByShippingCode := make(map[string]*shipmodel.Fulfillment, len(cmd.Result.Fulfillments))
	for _, ffm := range cmd.Result.Fulfillments {
		// ignore ffms that finished
		if (ffm.Status != status5.Z && ffm.Status != status5.S) ||
			!ffm.CODEtopTransferedAt.IsZero() {
			continue
		}

		feeLines := ffm.ProviderShippingFeeLines
		var newFeeLines []*shippingsharemodel.ShippingFeeLine
		for _, feeLine := range feeLines {
			if feeLine.ShippingFeeType == shipping_fee_type.Main {
				// keep the shipping fee type main (phí dịch vụ)
				newFeeLines = append(newFeeLines, feeLine)
				break
			}
		}
		updatesMap[ffm.ShippingCode] = &shipmodel.Fulfillment{
			ID:                       ffm.ID,
			ShippingFeeShopLines:     newFeeLines,
			ProviderShippingFeeLines: newFeeLines,
			ShippingFeeShop:          calcTotalFee(newFeeLines),
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
			update.ProviderShippingFeeLines = append(update.ProviderShippingFeeLines, &shippingsharemodel.ShippingFeeLine{
				ShippingFeeType:      shipping_fee_type.Insurance,
				Cost:                 line.InsuranceFee,
				ExternalShippingCode: line.ExternalCode,
			})
		}
		if line.ReturnFee != 0 {
			update.ProviderShippingFeeLines = append(update.ProviderShippingFeeLines, &shippingsharemodel.ShippingFeeLine{
				ShippingFeeType:      shipping_fee_type.Return,
				Cost:                 line.ReturnFee,
				ExternalShippingCode: line.ExternalCode,
			})
		}
		if line.Discount != 0 {
			cost := line.Discount
			if cost > 0 {
				cost = -cost
			}
			update.ProviderShippingFeeLines = append(update.ProviderShippingFeeLines, &shippingsharemodel.ShippingFeeLine{
				Cost:                 cost,
				ShippingFeeType:      shipping_fee_type.Discount,
				ExternalShippingCode: line.ExternalCode,
			})
		}
		if line.ChangeAddressFee != 0 {
			update.ProviderShippingFeeLines = append(update.ProviderShippingFeeLines, &shippingsharemodel.ShippingFeeLine{
				ShippingFeeType:      shipping_fee_type.AddressChange,
				Cost:                 line.ChangeAddressFee,
				ExternalShippingCode: line.ExternalCode,
			})
		}
		update.ShippingFeeShopLines = shippingsharemodel.GetShippingFeeShopLines(update.ProviderShippingFeeLines, ffm.EtopPriceRule, dot.Int(ffm.EtopAdjustedShippingFeeMain))
		totalFee := calcTotalFee(update.ShippingFeeShopLines)
		update.ShippingFeeShop = shipmodel.CalcShopShippingFee(totalFee, update)

		updatesMap[line.ExternalCode] = update
	}
	fulfillments := make([]*shipmodel.Fulfillment, 0, len(updatesMap))
	for _, ffm := range updatesMap {
		fulfillments = append(fulfillments, ffm)
	}
	return fulfillments, nil
}

func calcTotalFee(lines []*shippingsharemodel.ShippingFeeLine) int {
	res := 0
	for _, line := range lines {
		res += line.Cost
	}
	return res
}

func updateFulfillments(ctx context.Context, fulfillments []*shipmodel.Fulfillment) error {
	cmd := &shipmodelx.UpdateFulfillmentsCommand{
		Fulfillments: fulfillments,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return cm.Errorf(cm.Internal, err, "Không thể cập nhật ffm")
	}
	return nil
}
