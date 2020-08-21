package types

import (
	"o.o/api/top/int/types/spreadsheet"
	"o.o/api/top/types/common"
	"o.o/capi/dot"
	"o.o/common/jsonx"
)

type FulfillmentResponse struct {
	EdCode             string `json:"ed_code"`
	CustomerName       string `json:"customer_name"`
	CustomerPhone      string `json:"customer_phone"`
	ShippingAddress    string `json:"shipping_address"`
	District           string `json:"district"`
	DistrictCode       string `json:"district_code"`
	Province           string `json:"province"`
	ProvinceCode       string `json:"province_code"`
	Ward               string `json:"ward"`
	WardCode           string `json:"ward_code"`
	ProductDescription string `json:"product_description"`
	TotalWeight        int    `json:"total_weight"` // gram
	BasketValue        int    `json:"basket_value"`
	IncludeInsurance   bool   `json:"include_insurance"`
	CODAmount          int    `json:"cod_amount"`
	ShippingNote       string `json:"shipping_note"`
}

func (m *FulfillmentResponse) String() string { return jsonx.MustMarshalToString(m) }

type ImportFulfillmentsResponse struct {
	Data            *spreadsheet.SpreadsheetData `json:"data"`
	Fulfillments    []*FulfillmentResponse       `json:"fulfillments"`
	SpecificColumns []*SpecificColumn            `json:"specific_columns"`
	CellErrors      []*common.Error              `json:"cell_errors"`
	ImportID        dot.ID                       `json:"import_id"`
}

func (m *ImportFulfillmentsResponse) String() string { return jsonx.MustMarshalToString(m) }

type SpecificColumn struct {
	Fields []string `json:"fields"`
	Label  string   `json:"label"`
	Name   string   `json:"name"`
}

func (m *SpecificColumn) String() string { return jsonx.MustMarshalToString(m) }
