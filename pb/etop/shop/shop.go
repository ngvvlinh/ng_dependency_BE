package shop

import (
	"context"
	"net/url"
	"strconv"

	"github.com/golang/protobuf/jsonpb"

	"etop.vn/api/main/identity"
	"etop.vn/api/main/ledgering"
	"etop.vn/api/main/location"
	"etop.vn/api/main/receipting"
	"etop.vn/api/shopping/addressing"
	"etop.vn/api/shopping/carrying"
	"etop.vn/api/shopping/customering"
	"etop.vn/api/shopping/suppliering"
	summary "etop.vn/api/summary"
	catalogmodel "etop.vn/backend/com/main/catalog/model"
	"etop.vn/backend/pb/common"
	pbcm "etop.vn/backend/pb/common"
	pbetop "etop.vn/backend/pb/etop"
	pbs3 "etop.vn/backend/pb/etop/etc/status3"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/common/jsonx"
)

func PbSummaryTables(tables []*model.SummaryTable) []*SummaryTable {
	res := make([]*SummaryTable, len(tables))
	for i, table := range tables {
		res[i] = &SummaryTable{
			Label:   table.Label,
			Tags:    table.Tags,
			Columns: PbSummaryColRow(table.Cols),
			Rows:    PbSummaryColRow(table.Rows),
			Data:    PbSummaryData(table.Data),
		}
	}
	return res
}

func PbSummaryColRow(items []model.SummaryColRow) []SummaryColRow {
	res := make([]SummaryColRow, len(items))
	for i, item := range items {
		res[i] = SummaryColRow{
			Label:  item.Label,
			Spec:   item.Spec,
			Unit:   item.Unit,
			Indent: int32(item.Indent),
		}
	}
	return res
}

func PbSummaryData(data []model.SummaryItem) []SummaryItem {
	res := make([]SummaryItem, len(data))
	for i, item := range data {
		res[i] = SummaryItem{
			Spec:  item.Spec,
			Value: int32(item.Value),
			Unit:  item.Unit,
		}
	}
	return res
}

// From up
func PbSummaryTablesNew(tables []*summary.SummaryTable) []*SummaryTable {
	res := make([]*SummaryTable, len(tables))
	for i, table := range tables {
		res[i] = &SummaryTable{
			Label:   table.Label,
			Tags:    table.Tags,
			Columns: PbSummaryColRowNew(table.Cols),
			Rows:    PbSummaryColRowNew(table.Rows),
			Data:    PbSummaryDataNew(table.Data),
		}
	}
	return res
}

func PbSummaryColRowNew(items []summary.SummaryColRow) []SummaryColRow {
	res := make([]SummaryColRow, len(items))
	for i, item := range items {
		res[i] = SummaryColRow{
			Label:  item.Label,
			Spec:   item.Spec,
			Unit:   item.Unit,
			Indent: int32(item.Indent),
		}
	}
	return res
}

func PbSummaryDataNew(data []summary.SummaryItem) []SummaryItem {
	res := make([]SummaryItem, len(data))
	for i, item := range data {
		res[i] = SummaryItem{
			Spec:      item.Spec,
			Value:     int32(item.Value),
			Unit:      item.Unit,
			ImageUrls: item.ImageUrls,
			Label:     item.Label,
		}
	}
	return res
}

// MarshalJSONPB implements JSONPBMarshaler
func (m *SummaryTable) MarshalJSONPB(_ *jsonpb.Marshaler) ([]byte, error) {
	ncol := len(m.Columns)
	data := make([][]SummaryItem, len(m.Rows))
	for r := range m.Rows {
		data[r] = m.Data[r*ncol : (r+1)*ncol]
	}
	res := SummaryTableJSON{
		Label:   m.Label,
		Tags:    m.Tags,
		Columns: m.Columns,
		Rows:    m.Rows,
		Data:    data,
	}
	return jsonx.Marshal(res)
}

// UnmarshalJSONPB implements JSONPBUnmarshaler
func (m *SummaryTable) UnmarshalJSONPB(_ *jsonpb.Unmarshaler, data []byte) error {
	var tmp SummaryTableJSON
	if err := jsonx.Unmarshal(data, &tmp); err != nil {
		return err
	}
	ncol := len(tmp.Columns)
	mdata := make([]SummaryItem, len(tmp.Rows)*ncol)
	for r := range tmp.Rows {
		copy(mdata[r*ncol:], tmp.Data[r])
	}
	*m = SummaryTable{
		Label:   tmp.Label,
		Tags:    tmp.Tags,
		Columns: tmp.Columns,
		Rows:    tmp.Rows,
		Data:    mdata,
	}
	return nil
}

type SummaryTableJSON struct {
	Label   string          `json:"label"`
	Tags    []string        `json:"tags"`
	Columns []SummaryColRow `json:"columns"`
	Rows    []SummaryColRow `json:"rows"`
	Data    [][]SummaryItem `json:"data"`
}

func (m *ImportProductsResponse) HasErrors() []*common.Error {
	if len(m.CellErrors) > 0 {
		return m.CellErrors
	}
	return m.ImportErrors
}

func PbExportAttempts(ms []*model.ExportAttempt) []*ExportItem {
	res := make([]*ExportItem, len(ms))
	for i, m := range ms {
		res[i] = PbExportAttempt(m)
	}
	return res
}

func PbExportAttempt(m *model.ExportAttempt) *ExportItem {
	return &ExportItem{
		Id:           m.ID,
		Filename:     m.FileName,
		ExportType:   m.ExportType,
		DownloadUrl:  m.DownloadURL,
		AccountId:    m.AccountID,
		UserId:       m.UserID,
		CreatedAt:    common.PbTime(m.CreatedAt),
		DeletedAt:    common.PbTime(m.DeletedAt),
		RequestQuery: m.RequestQuery,
		MimeType:     m.MimeType,
		Status:       0,
		ExportErrors: common.PbErrorsFromModel(m.Errors),
		Error:        common.PbError(m.GetAbortedError()),
	}
}

func PbAuthorizedPartner(item *model.Partner, shop *model.Shop) *AuthorizedPartnerResponse {
	redirectUrl := ""
	if item.AvailableFromEtopConfig != nil {
		redirectUrl = item.AvailableFromEtopConfig.RedirectUrl
	}
	rUrl := GenerateRedirectAuthorizedPartnerURL(redirectUrl, shop)
	return &AuthorizedPartnerResponse{
		Partner:     pbetop.PbPublicAccountInfo(item),
		RedirectUrl: rUrl,
	}
}

func PbAuthorizedPartners(items []*model.Partner, shop *model.Shop) []*AuthorizedPartnerResponse {
	res := make([]*AuthorizedPartnerResponse, len(items))
	for i, item := range items {
		res[i] = PbAuthorizedPartner(item, shop)
	}
	return res
}

func GenerateRedirectAuthorizedPartnerURL(redirectUrl string, shop *model.Shop) string {
	u, err := url.Parse(redirectUrl)
	if err != nil {
		return ""
	}
	query := u.Query()
	query.Set("shop_id", strconv.FormatInt(shop.ID, 10))
	query.Set("email", shop.Email)
	u.RawQuery = query.Encode()
	res, _ := url.QueryUnescape(u.String())
	return res
}

func (a *Attribute) ToModel() *catalogmodel.ProductAttribute {
	if a == nil {
		return nil
	}
	return &catalogmodel.ProductAttribute{
		Name:  a.Name,
		Value: a.Value,
	}
}

func PbCustomer(m *customering.ShopCustomer) *Customer {
	return &Customer{
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
		CreatedAt: pbcm.PbTime(m.CreatedAt),
		UpdatedAt: pbcm.PbTime(m.UpdatedAt),
		Status:    pbs3.Pb(model.Status3(m.Status)),
	}
}

func PbCustopmerGroup(m *customering.ShopCustomerGroup) *CustomerGroup {
	return &CustomerGroup{
		Id:   m.ID,
		Name: m.Name,
	}
}

func PbCustomerGroups(ms []*customering.ShopCustomerGroup) []*CustomerGroup {
	res := make([]*CustomerGroup, len(ms))
	for i, m := range ms {
		res[i] = PbCustopmerGroup(m)
	}
	return res
}

func PbCustomers(ms []*customering.ShopCustomer) []*Customer {
	res := make([]*Customer, len(ms))
	for i, m := range ms {
		res[i] = PbCustomer(m)
	}
	return res
}

func PbSupplier(m *suppliering.ShopSupplier) *Supplier {
	return &Supplier{
		Id:        m.ID,
		ShopId:    m.ShopID,
		FullName:  m.FullName,
		Note:      m.Note,
		Status:    pbs3.Pb(model.Status3(m.Status)),
		CreatedAt: pbcm.PbTime(m.CreatedAt),
		UpdatedAt: pbcm.PbTime(m.UpdatedAt),
	}
}

func PbSuppliers(ms []*suppliering.ShopSupplier) []*Supplier {
	res := make([]*Supplier, len(ms))
	for i, m := range ms {
		res[i] = PbSupplier(m)
	}
	return res
}

func PbCarrier(m *carrying.ShopCarrier) *Carrier {
	return &Carrier{
		Id:        m.ID,
		ShopId:    m.ShopID,
		FullName:  m.FullName,
		Note:      m.Note,
		Status:    pbs3.Pb(model.Status3(m.Status)),
		CreatedAt: pbcm.PbTime(m.CreatedAt),
		UpdatedAt: pbcm.PbTime(m.UpdatedAt),
	}
}

func PbCarriers(ms []*carrying.ShopCarrier) []*Carrier {
	res := make([]*Carrier, len(ms))
	for i, m := range ms {
		res[i] = PbCarrier(m)
	}
	return res
}

func PbShopAddress(ctx context.Context, in *addressing.ShopTraderAddress, locationBus location.QueryBus) (*CustomerAddress, error) {
	query := &location.GetLocationQuery{
		DistrictCode: in.DistrictCode,
		WardCode:     in.WardCode,
	}
	if err := locationBus.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	province, district, ward := query.Result.Province, query.Result.District, query.Result.Ward
	out := &CustomerAddress{
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
		Coordinates:  pbetop.PbCoordinates(in.Coordinates),
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

func PbShopAddresses(ctx context.Context, ms []*addressing.ShopTraderAddress, locationBus location.QueryBus) ([]*CustomerAddress, error) {
	var err error
	res := make([]*CustomerAddress, len(ms))
	for i, m := range ms {
		res[i], err = PbShopAddress(ctx, m, locationBus)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func PbReceipt(m *receipting.Receipt) *Receipt {
	return &Receipt{
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
		Status:      pbs3.Pb(model.Status3(m.Status)),
		CreatedBy:   m.CreatedBy,
		CreatedType: string(m.CreatedType),
		PaidAt:      pbcm.PbTime(m.PaidAt),
		ConfirmedAt: pbcm.PbTime(m.ConfirmedAt),
		CancelledAt: pbcm.PbTime(m.CancelledAt),
		CreatedAt:   pbcm.PbTime(m.CreatedAt),
		UpdatedAt:   pbcm.PbTime(m.UpdatedAt),
	}
}

func PbReceipts(ms []*receipting.Receipt) []*Receipt {
	res := make([]*Receipt, len(ms))
	for i, m := range ms {
		res[i] = PbReceipt(m)
	}
	return res
}

func PbTrader(m *receipting.Trader) *Trader {
	if m == nil {
		return nil
	}
	return &Trader{
		Id:       m.ID,
		Type:     m.Type,
		FullName: m.FullName,
		Phone:    m.Phone,
		Deleted:  false,
	}
}

func PbTraders(ms []*receipting.Trader) []*Trader {
	res := make([]*Trader, len(ms))
	for i, m := range ms {
		res[i] = PbTrader(m)
	}
	return res
}

func PbReceiptLine(m *receipting.ReceiptLine) *ReceiptLine {
	if m == nil {
		return nil
	}
	return &ReceiptLine{
		RefId:  m.RefID,
		Title:  m.Title,
		Amount: m.Amount,
	}
}

func PbReceiptLines(ms []*receipting.ReceiptLine) []*ReceiptLine {
	res := make([]*ReceiptLine, len(ms))
	for i, m := range ms {
		res[i] = PbReceiptLine(m)
	}
	return res
}

func PbLedger(m *ledgering.ShopLedger) *Ledger {
	if m == nil {
		return nil
	}
	return &Ledger{
		Id:          m.ID,
		Name:        m.Name,
		BankAccount: PbBankAccount(m.BankAccount),
		Note:        m.Note,
		Type:        m.Type,
		CreatedBy:   m.CreatedBy,
		CreatedAt:   pbcm.PbTime(m.CreatedAt),
		UpdatedAt:   pbcm.PbTime(m.UpdatedAt),
	}
}

func PbLedgers(ms []*ledgering.ShopLedger) []*Ledger {
	res := make([]*Ledger, len(ms))
	for i, m := range ms {
		res[i] = PbLedger(m)
	}
	return res
}

func PbBankAccount(m *identity.BankAccount) *pbetop.BankAccount {
	if m == nil {
		return nil
	}
	return &pbetop.BankAccount{
		Name:          m.Name,
		Province:      m.Province,
		Branch:        m.Branch,
		AccountNumber: m.AccountNumber,
		AccountName:   m.AccountName,
	}
}
