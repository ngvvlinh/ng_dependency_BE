package common

import (
	"etop.vn/capi/dot"
	"etop.vn/common/jsonx"
)

type Empty struct {
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return jsonx.MustMarshalToString(m) }

type VersionInfoResponse struct {
	Service   string   `json:"service"`
	Version   string   `json:"version"`
	UpdatedAt dot.Time `json:"updated_at"`
}

func (m *VersionInfoResponse) Reset()         { *m = VersionInfoResponse{} }
func (m *VersionInfoResponse) String() string { return jsonx.MustMarshalToString(m) }

type IDRequest struct {
	// @required
	Id dot.ID `json:"id"`
}

func (m *IDRequest) Reset()         { *m = IDRequest{} }
func (m *IDRequest) String() string { return jsonx.MustMarshalToString(m) }

type CodeRequest struct {
	// @required
	Code string `json:"code"`
}

func (m *CodeRequest) Reset()         { *m = CodeRequest{} }
func (m *CodeRequest) String() string { return jsonx.MustMarshalToString(m) }

type NameRequest struct {
	// @required
	Name string `json:"name"`
}

func (m *NameRequest) Reset()         { *m = NameRequest{} }
func (m *NameRequest) String() string { return jsonx.MustMarshalToString(m) }

type IDsRequest struct {
	// @required
	Ids []dot.ID `json:"ids"`
}

func (m *IDsRequest) Reset()         { *m = IDsRequest{} }
func (m *IDsRequest) String() string { return jsonx.MustMarshalToString(m) }

type StatusResponse struct {
	Status string `json:"status"`
}

func (m *StatusResponse) Reset()         { *m = StatusResponse{} }
func (m *StatusResponse) String() string { return jsonx.MustMarshalToString(m) }

type IDMRequest struct {
	Id     dot.ID  `json:"id"`
	Paging *Paging `json:"paging"`
}

func (m *IDMRequest) Reset()         { *m = IDMRequest{} }
func (m *IDMRequest) String() string { return jsonx.MustMarshalToString(m) }

type Paging struct {
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	Sort   string `json:"sort"`
}

func (m *Paging) Reset()         { *m = Paging{} }
func (m *Paging) String() string { return jsonx.MustMarshalToString(m) }

type ForwardPaging struct {
	Since string `json:"since"`
	Limit int    `json:"limit"`
}

func (m *ForwardPaging) Reset()         { *m = ForwardPaging{} }
func (m *ForwardPaging) String() string { return jsonx.MustMarshalToString(m) }

type PageInfo struct {
	Total int      `json:"total"`
	Limit int      `json:"limit"`
	Sort  []string `json:"sort"`
}

func (m *PageInfo) Reset()         { *m = PageInfo{} }
func (m *PageInfo) String() string { return jsonx.MustMarshalToString(m) }

type ForwardPageInfo struct {
	Since   string `json:"since"`
	Limit   int    `json:"limit"`
	HasNext bool   `json:"has_next"`
}

func (m *ForwardPageInfo) Reset()         { *m = ForwardPageInfo{} }
func (m *ForwardPageInfo) String() string { return jsonx.MustMarshalToString(m) }

type RawJSONObject struct {
	Data []byte `json:"data"`
}

func (m *RawJSONObject) Reset()         { *m = RawJSONObject{} }
func (m *RawJSONObject) String() string { return jsonx.MustMarshalToString(m) }

type Error struct {
	Code string            `json:"code"`
	Msg  string            `json:"msg"`
	Meta map[string]string `json:"meta" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *Error) Reset()         { *m = Error{} }
func (m *Error) String() string { return jsonx.MustMarshalToString(m) }

type ErrorsResponse struct {
	Errors []*Error `json:"errors"`
}

func (m *ErrorsResponse) Reset()         { *m = ErrorsResponse{} }
func (m *ErrorsResponse) String() string { return jsonx.MustMarshalToString(m) }

type UpdatedResponse struct {
	Updated int `json:"updated"`
}

func (m *UpdatedResponse) Reset()         { *m = UpdatedResponse{} }
func (m *UpdatedResponse) String() string { return jsonx.MustMarshalToString(m) }

type RemovedResponse struct {
	Removed int `json:"removed"`
}

func (m *RemovedResponse) Reset()         { *m = RemovedResponse{} }
func (m *RemovedResponse) String() string { return jsonx.MustMarshalToString(m) }

type DeletedResponse struct {
	Deleted int `json:"deleted"`
}

func (m *DeletedResponse) Reset()         { *m = DeletedResponse{} }
func (m *DeletedResponse) String() string { return jsonx.MustMarshalToString(m) }

type MessageResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

func (m *MessageResponse) Reset()         { *m = MessageResponse{} }
func (m *MessageResponse) String() string { return jsonx.MustMarshalToString(m) }

type Filter struct {
	// Comma separated properties: "name,s_name"
	Name string `json:"name"`
	// Can be = ≠ (!=) < ≤ (<=) > ≥ (>=) ⊃ (c) ∈ (in) ∩ (n)
	//
	// - Text or set: ⊃ ∩
	// - Exactly: = ≠ ∈
	// - Numeric: = ≠ ∈ < ≤ > ≥
	// - Array: = ≠ (only with value is {})
	Op string `json:"op"`
	// Must always be string
	Value string `json:"value"`
}

func (m *Filter) Reset()         { *m = Filter{} }
func (m *Filter) String() string { return jsonx.MustMarshalToString(m) }

type CommonListRequest struct {
	Paging  *Paging   `json:"paging"`
	Filters []*Filter `json:"filters"`
}

func (m *CommonListRequest) Reset()         { *m = CommonListRequest{} }
func (m *CommonListRequest) String() string { return jsonx.MustMarshalToString(m) }

type MetaField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (m *MetaField) Reset()         { *m = MetaField{} }
func (m *MetaField) String() string { return jsonx.MustMarshalToString(m) }