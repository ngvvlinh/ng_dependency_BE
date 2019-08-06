package shop

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/golang/protobuf/jsonpb"

	"etop.vn/api/shopping/customering"
	"etop.vn/backend/pb/common"
	pbcm "etop.vn/backend/pb/common"
	pbetop "etop.vn/backend/pb/etop"
	pbs3 "etop.vn/backend/pb/etop/etc/status3"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/etop/model"
	catalogmodel "etop.vn/backend/pkg/services/catalog/model"
)

func PbUpdateVariantToModel(shopID int64, p *UpdateVariantRequest) *catalogmodel.ShopVariant {
	return &catalogmodel.ShopVariant{
		ShopID:      shopID,
		VariantID:   p.Id,
		Name:        p.Name,
		Code:        cm.Coalesce(p.Code, p.Sku),
		Description: p.Description,
		ShortDesc:   p.ShortDesc,
		DescHTML:    p.DescHtml,
		Note:        p.Note,
		RetailPrice: p.RetailPrice,
		ListPrice:   p.ListPrice,
		CostPrice:   p.CostPrice,
	}
}

func PbUpdateProductToModel(shopID int64, p *UpdateProductRequest) *catalogmodel.ShopProduct {
	return &catalogmodel.ShopProduct{
		ShopID:      shopID,
		ProductID:   p.Id,
		Name:        p.Name,
		Description: p.Description,
		ShortDesc:   p.ShortDesc,
		DescHTML:    p.DescHtml,
	}
}

func PbCollections(items []*catalogmodel.ShopCollection) []*Collection {
	res := make([]*Collection, len(items))
	for i, item := range items {
		res[i] = PbCollection(item)
	}
	return res
}

func PbCollection(c *catalogmodel.ShopCollection) *Collection {
	return &Collection{
		Id:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		DescHtml:    c.DescHTML,
		ShortDesc:   c.ShortDesc,
		CreatedAt:   common.PbTime(c.CreatedAt),
		UpdatedAt:   common.PbTime(c.UpdatedAt),
	}
}

func PbCreateCollection(shopID int64, p *CreateCollectionRequest) *catalogmodel.ShopCollection {
	return &catalogmodel.ShopCollection{
		ShopID:      shopID,
		Name:        p.Name,
		DescHTML:    p.DescHtml,
		Description: p.Description,
		ShortDesc:   p.ShortDesc,
	}
}

func PbUpdateCollection(shopID int64, p *UpdateCollectionRequest) *catalogmodel.ShopCollection {
	return &catalogmodel.ShopCollection{
		ID:          p.Id,
		ShopID:      shopID,
		Name:        p.Name,
		DescHTML:    p.DescHtml,
		Description: p.Description,
		ShortDesc:   p.ShortDesc,
	}
}

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
	return json.Marshal(res)
}

// UnmarshalJSONPB implements JSONPBUnmarshaler
func (m *SummaryTable) UnmarshalJSONPB(_ *jsonpb.Unmarshaler, data []byte) error {
	var tmp SummaryTableJSON
	if err := json.Unmarshal(data, &tmp); err != nil {
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

func (a *Attribute) ToModel() catalogmodel.ProductAttribute {
	if a == nil {
		return catalogmodel.ProductAttribute{}
	}
	return catalogmodel.ProductAttribute{
		Name:  a.Name,
		Value: a.Value,
	}
}

func PbCustomer(m *customering.ShopCustomer) *Customer {
	return &Customer{
		Id:        m.ID,
		ShopId:    m.ShopID,
		Code:      m.Code,
		FullName:  m.FullName,
		Note:      m.Note,
		Phone:     m.Phone,
		Email:     m.Email,
		CreatedAt: pbcm.PbTime(m.CreatedAt),
		UpdatedAt: pbcm.PbTime(m.UpdatedAt),
		Status:    pbs3.Pb(model.Status3(m.Status)),
	}
}

func PbCustomers(ms []*customering.ShopCustomer) []*Customer {
	res := make([]*Customer, len(ms))
	for i, m := range ms {
		res[i] = PbCustomer(m)
	}
	return res
}
