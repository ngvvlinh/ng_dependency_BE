package convertpb

import (
	"context"
	"net/url"

	"etop.vn/api/main/ledgering"
	"etop.vn/api/main/location"
	"etop.vn/api/main/purchaseorder"
	"etop.vn/api/main/receipting"
	"etop.vn/api/shopping/addressing"
	"etop.vn/api/shopping/carrying"
	"etop.vn/api/shopping/customering"
	"etop.vn/api/shopping/suppliering"
	"etop.vn/api/summary"
	apishop "etop.vn/api/top/int/shop"
	identitymodel "etop.vn/backend/com/main/identity/model"
	identitysharemodel "etop.vn/backend/com/main/identity/sharemodel"
	"etop.vn/backend/pkg/common/apifw/cmapi"
	"etop.vn/backend/pkg/etop/model"
)

func PbSummaryTables(tables []*model.SummaryTable) []*apishop.SummaryTable {
	res := make([]*apishop.SummaryTable, len(tables))
	for i, table := range tables {
		res[i] = &apishop.SummaryTable{
			Label:   table.Label,
			Tags:    table.Tags,
			Columns: PbSummaryColRow(table.Cols),
			Rows:    PbSummaryColRow(table.Rows),
			Data:    PbSummaryData(table.Data),
		}
	}
	return res
}

func PbSummaryColRow(items []model.SummaryColRow) []apishop.SummaryColRow {
	res := make([]apishop.SummaryColRow, len(items))
	for i, item := range items {
		res[i] = apishop.SummaryColRow{
			Label:  item.Label,
			Spec:   item.Spec,
			Unit:   item.Unit,
			Indent: item.Indent,
		}
	}
	return res
}

func PbSummaryData(data []model.SummaryItem) []apishop.SummaryItem {
	res := make([]apishop.SummaryItem, len(data))
	for i, item := range data {
		res[i] = apishop.SummaryItem{
			Spec:  item.Spec,
			Value: item.Value,
			Unit:  item.Unit,
		}
	}
	return res
}

// From up
func PbSummaryTablesNew(tables []*summary.SummaryTable) []*apishop.SummaryTable {
	res := make([]*apishop.SummaryTable, len(tables))
	for i, table := range tables {
		res[i] = &apishop.SummaryTable{
			Label:   table.Label,
			Tags:    table.Tags,
			Columns: PbSummaryColRowNew(table.Cols),
			Rows:    PbSummaryColRowNew(table.Rows),
			Data:    PbSummaryDataNew(table.Data),
		}
	}
	return res
}

func PbSummaryColRowNew(items []summary.SummaryColRow) []apishop.SummaryColRow {
	res := make([]apishop.SummaryColRow, len(items))
	for i, item := range items {
		res[i] = apishop.SummaryColRow{
			Label:  item.Label,
			Spec:   item.Spec,
			Unit:   item.Unit,
			Indent: item.Indent,
		}
	}
	return res
}

func PbSummaryDataNew(data []summary.SummaryItem) []apishop.SummaryItem {
	res := make([]apishop.SummaryItem, len(data))
	for i, item := range data {
		res[i] = apishop.SummaryItem{
			Spec:      item.Spec,
			Value:     item.Value,
			Unit:      item.Unit,
			ImageUrls: item.ImageUrls,
			Label:     item.Label,
		}
	}
	return res
}

func PbExportAttempts(ms []*model.ExportAttempt) []*apishop.ExportItem {
	res := make([]*apishop.ExportItem, len(ms))
	for i, m := range ms {
		res[i] = PbExportAttempt(m)
	}
	return res
}

func PbExportAttempt(m *model.ExportAttempt) *apishop.ExportItem {
	return &apishop.ExportItem{
		Id:           m.ID,
		Filename:     m.FileName,
		ExportType:   m.ExportType,
		DownloadUrl:  m.DownloadURL,
		AccountId:    m.AccountID,
		UserId:       m.UserID,
		CreatedAt:    cmapi.PbTime(m.CreatedAt),
		DeletedAt:    cmapi.PbTime(m.DeletedAt),
		RequestQuery: m.RequestQuery,
		MimeType:     m.MimeType,
		Status:       0,
		ExportErrors: cmapi.PbErrorsFromModel(m.Errors),
		Error:        cmapi.PbError(m.GetAbortedError()),
	}
}

func PbAuthorizedPartner(item *identitymodel.Partner, s *identitymodel.Shop) *apishop.AuthorizedPartnerResponse {
	redirectUrl := ""
	if item.AvailableFromEtopConfig != nil {
		redirectUrl = item.AvailableFromEtopConfig.RedirectUrl
	}
	rUrl := GenerateRedirectAuthorizedPartnerURL(redirectUrl, s)
	return &apishop.AuthorizedPartnerResponse{
		Partner:     PbPublicAccountInfo(item),
		RedirectUrl: rUrl,
	}
}

func PbAuthorizedPartners(items []*identitymodel.Partner, s *identitymodel.Shop) []*apishop.AuthorizedPartnerResponse {
	res := make([]*apishop.AuthorizedPartnerResponse, len(items))
	for i, item := range items {
		res[i] = PbAuthorizedPartner(item, s)
	}
	return res
}

func GenerateRedirectAuthorizedPartnerURL(redirectUrl string, shop *identitymodel.Shop) string {
	u, err := url.Parse(redirectUrl)
	if err != nil {
		return ""
	}
	query := u.Query()
	query.Set("shop_id", shop.ID.String())
	query.Set("email", shop.Email)
	u.RawQuery = query.Encode()
	res, _ := url.QueryUnescape(u.String())
	return res
}

func PbCustomer(m *customering.ShopCustomer) *apishop.Customer {
	return &apishop.Customer{
		Id:        m.ID,
		ShopId:    m.ShopID,
		GroupIds:  m.GroupIDs,
		FullName:  m.FullName,
		Code:      m.Code,
		Note:      m.Note,
		Phone:     m.Phone,
		Email:     m.Email,
		Gender:    m.Gender,
		Type:      m.Type,
		Birthday:  m.Birthday,
		CreatedAt: cmapi.PbTime(m.CreatedAt),
		UpdatedAt: cmapi.PbTime(m.UpdatedAt),
		Status:    m.Status,
	}
}

func PbCustopmerGroup(m *customering.ShopCustomerGroup) *apishop.CustomerGroup {
	return &apishop.CustomerGroup{
		Id:   m.ID,
		Name: m.Name,
	}
}

func PbCustomerGroups(ms []*customering.ShopCustomerGroup) []*apishop.CustomerGroup {
	res := make([]*apishop.CustomerGroup, len(ms))
	for i, m := range ms {
		res[i] = PbCustopmerGroup(m)
	}
	return res
}

func PbCustomers(ms []*customering.ShopCustomer) []*apishop.Customer {
	res := make([]*apishop.Customer, len(ms))
	for i, m := range ms {
		res[i] = PbCustomer(m)
	}
	return res
}

func PbSupplier(m *suppliering.ShopSupplier) *apishop.Supplier {
	if m == nil {
		return nil
	}
	return &apishop.Supplier{
		Id:                m.ID,
		ShopId:            m.ShopID,
		FullName:          m.FullName,
		Note:              m.Note,
		Code:              m.Code,
		Phone:             m.Phone,
		Email:             m.Email,
		CompanyName:       m.CompanyName,
		TaxNumber:         m.TaxNumber,
		HeadquaterAddress: m.HeadquaterAddress,

		Status:    m.Status,
		CreatedAt: cmapi.PbTime(m.CreatedAt),
		UpdatedAt: cmapi.PbTime(m.UpdatedAt),
	}
}

func PbSuppliers(ms []*suppliering.ShopSupplier) []*apishop.Supplier {
	res := make([]*apishop.Supplier, len(ms))
	for i, m := range ms {
		res[i] = PbSupplier(m)
	}
	return res
}

func PbCarrier(m *carrying.ShopCarrier) *apishop.Carrier {
	return &apishop.Carrier{
		Id:        m.ID,
		ShopId:    m.ShopID,
		FullName:  m.FullName,
		Note:      m.Note,
		Status:    m.Status,
		CreatedAt: cmapi.PbTime(m.CreatedAt),
		UpdatedAt: cmapi.PbTime(m.UpdatedAt),
	}
}

func PbCarriers(ms []*carrying.ShopCarrier) []*apishop.Carrier {
	res := make([]*apishop.Carrier, len(ms))
	for i, m := range ms {
		res[i] = PbCarrier(m)
	}
	return res
}

func PbShopAddress(ctx context.Context, in *addressing.ShopTraderAddress, locationBus location.QueryBus) (*apishop.CustomerAddress, error) {
	query := &location.GetLocationQuery{
		DistrictCode: in.DistrictCode,
		WardCode:     in.WardCode,
	}
	if err := locationBus.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	province, district, ward := query.Result.Province, query.Result.District, query.Result.Ward
	out := &apishop.CustomerAddress{
		Id:           in.ID,
		District:     "",
		DistrictCode: in.DistrictCode,
		Ward:         "",
		Company:      in.Company,
		WardCode:     in.WardCode,
		Address1:     in.Address1,
		Address2:     in.Address2,
		FullName:     in.FullName,
		Phone:        in.Phone,
		Email:        in.Email,
		Position:     in.Position,
		Coordinates:  PbCoordinates(in.Coordinates),
	}
	if ward != nil {
		out.Ward = ward.Name
	}
	if district != nil {
		out.District = district.Name
	}
	if province != nil {
		out.Province = province.Name
		out.ProvinceCode = province.Code
	}
	return out, nil
}

func PbShopAddresses(ctx context.Context, ms []*addressing.ShopTraderAddress, locationBus location.QueryBus) ([]*apishop.CustomerAddress, error) {
	var err error
	res := make([]*apishop.CustomerAddress, len(ms))
	for i, m := range ms {
		res[i], err = PbShopAddress(ctx, m, locationBus)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func PbReceipt(m *receipting.Receipt) *apishop.Receipt {
	return &apishop.Receipt{
		Id:           m.ID,
		ShopId:       m.ShopID,
		TraderId:     m.TraderID,
		Code:         m.Code,
		Title:        m.Title,
		Type:         m.Type,
		Description:  m.Description,
		Amount:       m.Amount,
		LedgerId:     m.LedgerID,
		RefType:      m.RefType,
		CancelReason: m.CancelReason,
		Lines:        PbReceiptLines(m.Lines),
		Trader:       PbTrader(m.Trader),
		Status:       m.Status,
		CreatedBy:    m.CreatedBy,
		CreatedType:  m.Mode,
		Mode:         m.Mode,
		PaidAt:       cmapi.PbTime(m.PaidAt),
		ConfirmedAt:  cmapi.PbTime(m.ConfirmedAt),
		CancelledAt:  cmapi.PbTime(m.CancelledAt),
		CreatedAt:    cmapi.PbTime(m.CreatedAt),
		UpdatedAt:    cmapi.PbTime(m.UpdatedAt),
	}
}

func PbReceipts(ms []*receipting.Receipt) []*apishop.Receipt {
	res := make([]*apishop.Receipt, len(ms))
	for i, m := range ms {
		res[i] = PbReceipt(m)
	}
	return res
}

func PbTrader(m *receipting.Trader) *apishop.Trader {
	if m == nil {
		return nil
	}
	return &apishop.Trader{
		Id:       m.ID,
		Type:     m.Type,
		FullName: m.FullName,
		Phone:    m.Phone,
		Deleted:  false,
	}
}

func PbTraders(ms []*receipting.Trader) []*apishop.Trader {
	res := make([]*apishop.Trader, len(ms))
	for i, m := range ms {
		res[i] = PbTrader(m)
	}
	return res
}

func PbReceiptLine(m *receipting.ReceiptLine) *apishop.ReceiptLine {
	if m == nil {
		return nil
	}
	return &apishop.ReceiptLine{
		RefId:  m.RefID,
		Title:  m.Title,
		Amount: m.Amount,
	}
}

func PbReceiptLines(ms []*receipting.ReceiptLine) []*apishop.ReceiptLine {
	res := make([]*apishop.ReceiptLine, len(ms))
	for i, m := range ms {
		res[i] = PbReceiptLine(m)
	}
	return res
}

func PbLedger(m *ledgering.ShopLedger) *apishop.Ledger {
	if m == nil {
		return nil
	}
	return &apishop.Ledger{
		Id:          m.ID,
		Name:        m.Name,
		BankAccount: PbBankAccount((*identitysharemodel.BankAccount)(m.BankAccount)),
		Note:        m.Note,
		Type:        m.Type,
		CreatedBy:   m.CreatedBy,
		CreatedAt:   cmapi.PbTime(m.CreatedAt),
		UpdatedAt:   cmapi.PbTime(m.UpdatedAt),
	}
}

func PbLedgers(ms []*ledgering.ShopLedger) []*apishop.Ledger {
	res := make([]*apishop.Ledger, len(ms))
	for i, m := range ms {
		res[i] = PbLedger(m)
	}
	return res
}

func PbPurchaseOrderLine(m *purchaseorder.PurchaseOrderLine) *apishop.PurchaseOrderLine {
	if m == nil {
		return nil
	}
	return &apishop.PurchaseOrderLine{
		ProductName:  m.ProductName,
		ImageUrl:     m.ImageUrl,
		ProductId:    m.ProductID,
		Code:         m.Code,
		Attributes:   m.Attributes,
		VariantId:    m.VariantID,
		Quantity:     m.Quantity,
		PaymentPrice: m.PaymentPrice,
	}
}

func PbPurchaseOrderLines(ms []*purchaseorder.PurchaseOrderLine) []*apishop.PurchaseOrderLine {
	res := make([]*apishop.PurchaseOrderLine, len(ms))
	for i, m := range ms {
		res[i] = PbPurchaseOrderLine(m)
	}
	return res
}

func PbPurchaseOrder(m *purchaseorder.PurchaseOrder) *apishop.PurchaseOrder {
	if m == nil {
		return nil
	}
	return &apishop.PurchaseOrder{
		Id:              m.ID,
		ShopId:          m.ShopID,
		SupplierId:      m.SupplierID,
		Supplier:        PbPurchaseOrderSupplier(m.Supplier),
		BasketValue:     m.BasketValue,
		DiscountLines:   m.DiscountLines,
		TotalDiscount:   m.TotalDiscount,
		FeeLines:        m.FeeLines,
		TotalFee:        m.TotalFee,
		TotalAmount:     m.TotalAmount,
		Code:            m.Code,
		Note:            m.Note,
		Status:          m.Status,
		Lines:           PbPurchaseOrderLines(m.Lines),
		PaidAmount:      m.PaidAmount,
		CreatedBy:       m.CreatedBy,
		CancelledReason: m.CancelReason,
		ConfirmedAt:     cmapi.PbTime(m.ConfirmedAt),
		CancelledAt:     cmapi.PbTime(m.CancelledAt),
		CreatedAt:       cmapi.PbTime(m.CreatedAt),
		UpdatedAt:       cmapi.PbTime(m.UpdatedAt),
	}
}

func PbPurchaseOrders(ms []*purchaseorder.PurchaseOrder) []*apishop.PurchaseOrder {
	res := make([]*apishop.PurchaseOrder, len(ms))
	for i, m := range ms {
		res[i] = PbPurchaseOrder(m)
	}
	return res
}

func PbPurchaseOrderSupplier(m *purchaseorder.PurchaseOrderSupplier) *apishop.PurchaseOrderSupplier {
	if m == nil {
		return nil
	}
	return &apishop.PurchaseOrderSupplier{
		FullName:           m.FullName,
		Phone:              m.Phone,
		Email:              m.Email,
		CompanyName:        m.CompanyName,
		TaxNumber:          m.TaxNumber,
		HeadquarterAddress: m.HeadquarterAddress,
		Deleted:            m.Deleted,
	}
}
