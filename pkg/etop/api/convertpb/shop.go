package convertpb

import (
	"context"
	"net/url"

	"etop.vn/api/main/ledgering"
	"etop.vn/api/main/location"
	"etop.vn/api/main/purchaseorder"
	"etop.vn/api/main/receipting"
	"etop.vn/api/pb/etop/shop"
	"etop.vn/api/shopping/addressing"
	"etop.vn/api/shopping/carrying"
	"etop.vn/api/shopping/customering"
	"etop.vn/api/shopping/suppliering"
	"etop.vn/api/summary"
	"etop.vn/backend/pkg/common/cmapi"
	"etop.vn/backend/pkg/etop/model"
)

func PbSummaryTables(tables []*model.SummaryTable) []*shop.SummaryTable {
	res := make([]*shop.SummaryTable, len(tables))
	for i, table := range tables {
		res[i] = &shop.SummaryTable{
			Label:   table.Label,
			Tags:    table.Tags,
			Columns: PbSummaryColRow(table.Cols),
			Rows:    PbSummaryColRow(table.Rows),
			Data:    PbSummaryData(table.Data),
		}
	}
	return res
}

func PbSummaryColRow(items []model.SummaryColRow) []shop.SummaryColRow {
	res := make([]shop.SummaryColRow, len(items))
	for i, item := range items {
		res[i] = shop.SummaryColRow{
			Label:  item.Label,
			Spec:   item.Spec,
			Unit:   item.Unit,
			Indent: int32(item.Indent),
		}
	}
	return res
}

func PbSummaryData(data []model.SummaryItem) []shop.SummaryItem {
	res := make([]shop.SummaryItem, len(data))
	for i, item := range data {
		res[i] = shop.SummaryItem{
			Spec:  item.Spec,
			Value: int32(item.Value),
			Unit:  item.Unit,
		}
	}
	return res
}

// From up
func PbSummaryTablesNew(tables []*summary.SummaryTable) []*shop.SummaryTable {
	res := make([]*shop.SummaryTable, len(tables))
	for i, table := range tables {
		res[i] = &shop.SummaryTable{
			Label:   table.Label,
			Tags:    table.Tags,
			Columns: PbSummaryColRowNew(table.Cols),
			Rows:    PbSummaryColRowNew(table.Rows),
			Data:    PbSummaryDataNew(table.Data),
		}
	}
	return res
}

func PbSummaryColRowNew(items []summary.SummaryColRow) []shop.SummaryColRow {
	res := make([]shop.SummaryColRow, len(items))
	for i, item := range items {
		res[i] = shop.SummaryColRow{
			Label:  item.Label,
			Spec:   item.Spec,
			Unit:   item.Unit,
			Indent: int32(item.Indent),
		}
	}
	return res
}

func PbSummaryDataNew(data []summary.SummaryItem) []shop.SummaryItem {
	res := make([]shop.SummaryItem, len(data))
	for i, item := range data {
		res[i] = shop.SummaryItem{
			Spec:      item.Spec,
			Value:     int32(item.Value),
			Unit:      item.Unit,
			ImageUrls: item.ImageUrls,
			Label:     item.Label,
		}
	}
	return res
}

func PbExportAttempts(ms []*model.ExportAttempt) []*shop.ExportItem {
	res := make([]*shop.ExportItem, len(ms))
	for i, m := range ms {
		res[i] = PbExportAttempt(m)
	}
	return res
}

func PbExportAttempt(m *model.ExportAttempt) *shop.ExportItem {
	return &shop.ExportItem{
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

func PbAuthorizedPartner(item *model.Partner, s *model.Shop) *shop.AuthorizedPartnerResponse {
	redirectUrl := ""
	if item.AvailableFromEtopConfig != nil {
		redirectUrl = item.AvailableFromEtopConfig.RedirectUrl
	}
	rUrl := GenerateRedirectAuthorizedPartnerURL(redirectUrl, s)
	return &shop.AuthorizedPartnerResponse{
		Partner:     PbPublicAccountInfo(item),
		RedirectUrl: rUrl,
	}
}

func PbAuthorizedPartners(items []*model.Partner, s *model.Shop) []*shop.AuthorizedPartnerResponse {
	res := make([]*shop.AuthorizedPartnerResponse, len(items))
	for i, item := range items {
		res[i] = PbAuthorizedPartner(item, s)
	}
	return res
}

func GenerateRedirectAuthorizedPartnerURL(redirectUrl string, shop *model.Shop) string {
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

func PbCustomer(m *customering.ShopCustomer) *shop.Customer {
	return &shop.Customer{
		Id:        m.ID,
		ShopId:    m.ShopID,
		GroupIds:  m.GroupIDs,
		FullName:  m.FullName,
		Code:      m.Code,
		Note:      m.Note,
		Phone:     m.Phone,
		Email:     m.Email,
		Gender:    m.Gender,
		Type:      string(m.Type),
		Birthday:  m.Birthday,
		CreatedAt: cmapi.PbTime(m.CreatedAt),
		UpdatedAt: cmapi.PbTime(m.UpdatedAt),
		Status:    Pb3(model.Status3(m.Status)),
	}
}

func PbCustopmerGroup(m *customering.ShopCustomerGroup) *shop.CustomerGroup {
	return &shop.CustomerGroup{
		Id:   m.ID,
		Name: m.Name,
	}
}

func PbCustomerGroups(ms []*customering.ShopCustomerGroup) []*shop.CustomerGroup {
	res := make([]*shop.CustomerGroup, len(ms))
	for i, m := range ms {
		res[i] = PbCustopmerGroup(m)
	}
	return res
}

func PbCustomers(ms []*customering.ShopCustomer) []*shop.Customer {
	res := make([]*shop.Customer, len(ms))
	for i, m := range ms {
		res[i] = PbCustomer(m)
	}
	return res
}

func PbSupplier(m *suppliering.ShopSupplier) *shop.Supplier {
	if m == nil {
		return nil
	}
	return &shop.Supplier{
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

		Status:    Pb3(model.Status3(m.Status)),
		CreatedAt: cmapi.PbTime(m.CreatedAt),
		UpdatedAt: cmapi.PbTime(m.UpdatedAt),
	}
}

func PbSuppliers(ms []*suppliering.ShopSupplier) []*shop.Supplier {
	res := make([]*shop.Supplier, len(ms))
	for i, m := range ms {
		res[i] = PbSupplier(m)
	}
	return res
}

func PbCarrier(m *carrying.ShopCarrier) *shop.Carrier {
	return &shop.Carrier{
		Id:        m.ID,
		ShopId:    m.ShopID,
		FullName:  m.FullName,
		Note:      m.Note,
		Status:    Pb3(model.Status3(m.Status)),
		CreatedAt: cmapi.PbTime(m.CreatedAt),
		UpdatedAt: cmapi.PbTime(m.UpdatedAt),
	}
}

func PbCarriers(ms []*carrying.ShopCarrier) []*shop.Carrier {
	res := make([]*shop.Carrier, len(ms))
	for i, m := range ms {
		res[i] = PbCarrier(m)
	}
	return res
}

func PbShopAddress(ctx context.Context, in *addressing.ShopTraderAddress, locationBus location.QueryBus) (*shop.CustomerAddress, error) {
	query := &location.GetLocationQuery{
		DistrictCode: in.DistrictCode,
		WardCode:     in.WardCode,
	}
	if err := locationBus.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	province, district, ward := query.Result.Province, query.Result.District, query.Result.Ward
	out := &shop.CustomerAddress{
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

func PbShopAddresses(ctx context.Context, ms []*addressing.ShopTraderAddress, locationBus location.QueryBus) ([]*shop.CustomerAddress, error) {
	var err error
	res := make([]*shop.CustomerAddress, len(ms))
	for i, m := range ms {
		res[i], err = PbShopAddress(ctx, m, locationBus)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func PbReceipt(m *receipting.Receipt) *shop.Receipt {
	return &shop.Receipt{
		Id:          m.ID,
		ShopId:      m.ShopID,
		TraderId:    m.TraderID,
		Code:        m.Code,
		Title:       m.Title,
		Type:        string(m.Type),
		Description: m.Description,
		Amount:      m.Amount,
		LedgerId:    m.LedgerID,
		RefType:     string(m.RefType),
		Lines:       PbReceiptLines(m.Lines),
		Trader:      PbTrader(m.Trader),
		Status:      Pb3(model.Status3(m.Status)),
		CreatedBy:   m.CreatedBy,
		CreatedType: string(m.CreatedType),
		PaidAt:      cmapi.PbTime(m.PaidAt),
		ConfirmedAt: cmapi.PbTime(m.ConfirmedAt),
		CancelledAt: cmapi.PbTime(m.CancelledAt),
		CreatedAt:   cmapi.PbTime(m.CreatedAt),
		UpdatedAt:   cmapi.PbTime(m.UpdatedAt),
	}
}

func PbReceipts(ms []*receipting.Receipt) []*shop.Receipt {
	res := make([]*shop.Receipt, len(ms))
	for i, m := range ms {
		res[i] = PbReceipt(m)
	}
	return res
}

func PbTrader(m *receipting.Trader) *shop.Trader {
	if m == nil {
		return nil
	}
	return &shop.Trader{
		Id:       m.ID,
		Type:     m.Type,
		FullName: m.FullName,
		Phone:    m.Phone,
		Deleted:  false,
	}
}

func PbTraders(ms []*receipting.Trader) []*shop.Trader {
	res := make([]*shop.Trader, len(ms))
	for i, m := range ms {
		res[i] = PbTrader(m)
	}
	return res
}

func PbReceiptLine(m *receipting.ReceiptLine) *shop.ReceiptLine {
	if m == nil {
		return nil
	}
	return &shop.ReceiptLine{
		RefId:  m.RefID,
		Title:  m.Title,
		Amount: int(m.Amount),
	}
}

func PbReceiptLines(ms []*receipting.ReceiptLine) []*shop.ReceiptLine {
	res := make([]*shop.ReceiptLine, len(ms))
	for i, m := range ms {
		res[i] = PbReceiptLine(m)
	}
	return res
}

func PbLedger(m *ledgering.ShopLedger) *shop.Ledger {
	if m == nil {
		return nil
	}
	return &shop.Ledger{
		Id:          m.ID,
		Name:        m.Name,
		BankAccount: PbBankAccount((*model.BankAccount)(m.BankAccount)),
		Note:        m.Note,
		Type:        m.Type,
		CreatedBy:   m.CreatedBy,
		CreatedAt:   cmapi.PbTime(m.CreatedAt),
		UpdatedAt:   cmapi.PbTime(m.UpdatedAt),
	}
}

func PbLedgers(ms []*ledgering.ShopLedger) []*shop.Ledger {
	res := make([]*shop.Ledger, len(ms))
	for i, m := range ms {
		res[i] = PbLedger(m)
	}
	return res
}

func PbPurchaseOrderLine(m *purchaseorder.PurchaseOrderLine) *shop.PurchaseOrderLine {
	if m == nil {
		return nil
	}
	return &shop.PurchaseOrderLine{
		ProductName:  m.ProductName,
		ImageUrl:     m.ImageUrl,
		ProductId:    m.ProductID,
		Code:         m.Code,
		Attributes:   PbAttributes(m.Attributes),
		VariantId:    m.VariantID,
		Quantity:     m.Quantity,
		PaymentPrice: m.PaymentPrice,
	}
}

func PbPurchaseOrderLines(ms []*purchaseorder.PurchaseOrderLine) []*shop.PurchaseOrderLine {
	res := make([]*shop.PurchaseOrderLine, len(ms))
	for i, m := range ms {
		res[i] = PbPurchaseOrderLine(m)
	}
	return res
}

func PbPurchaseOrder(m *purchaseorder.PurchaseOrder) *shop.PurchaseOrder {
	if m == nil {
		return nil
	}
	return &shop.PurchaseOrder{
		Id:              m.ID,
		ShopId:          m.ShopID,
		SupplierId:      m.SupplierID,
		Supplier:        PbPurchaseOrderSupplier(m.Supplier),
		BasketValue:     m.BasketValue,
		TotalDiscount:   m.TotalDiscount,
		TotalAmount:     m.TotalAmount,
		Code:            m.Code,
		Note:            m.Note,
		Status:          Pb3(model.Status3(m.Status)),
		Lines:           PbPurchaseOrderLines(m.Lines),
		PaidAmount:      m.PaidAmount,
		CreatedBy:       m.CreatedBy,
		CancelledReason: m.CancelledReason,
		ConfirmedAt:     cmapi.PbTime(m.ConfirmedAt),
		CancelledAt:     cmapi.PbTime(m.CancelledAt),
		CreatedAt:       cmapi.PbTime(m.CreatedAt),
		UpdatedAt:       cmapi.PbTime(m.UpdatedAt),
	}
}

func PbPurchaseOrders(ms []*purchaseorder.PurchaseOrder) []*shop.PurchaseOrder {
	res := make([]*shop.PurchaseOrder, len(ms))
	for i, m := range ms {
		res[i] = PbPurchaseOrder(m)
	}
	return res
}

func PbPurchaseOrderSupplier(m *purchaseorder.PurchaseOrderSupplier) *shop.PurchaseOrderSupplier {
	if m == nil {
		return nil
	}
	return &shop.PurchaseOrderSupplier{
		FullName:           m.FullName,
		Phone:              m.Phone,
		Email:              m.Email,
		CompanyName:        m.CompanyName,
		TaxNumber:          m.TaxNumber,
		HeadquarterAddress: m.HeadquarterAddress,
		Deleted:            m.Deleted,
	}
}
