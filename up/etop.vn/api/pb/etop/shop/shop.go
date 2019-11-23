package shop

import (
	"github.com/golang/protobuf/jsonpb"

	"etop.vn/api/pb/common"
	"etop.vn/common/jsonx"
)

type SummaryTableJSON struct {
	Label   string          `json:"label"`
	Tags    []string        `json:"tags"`
	Columns []SummaryColRow `json:"columns"`
	Rows    []SummaryColRow `json:"rows"`
	Data    [][]SummaryItem `json:"data"`
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

func (m *ImportProductsResponse) HasErrors() []*common.Error {
	if len(m.CellErrors) > 0 {
		return m.CellErrors
	}
	return m.ImportErrors
}
