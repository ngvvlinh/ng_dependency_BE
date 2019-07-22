package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"
	pbts "github.com/golang/protobuf/ptypes/timestamp"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/common/l"
	"etop.vn/common/xerrors"
)

var ll = l.New()

func (m *Error) Error() string {
	return m.Msg
}

var marshaler = jsonpb.Marshaler{EmitDefaults: true, OrigName: true}

func Marshal(v proto.Message) ([]byte, error) {
	var buf bytes.Buffer
	err := marshaler.Marshal(&buf, v)
	return buf.Bytes(), err
}

func MustMarshalToString(v proto.Message) string {
	var b strings.Builder
	err := marshaler.Marshal(&b, v)
	if err != nil {
		panic(err)
	}
	return b.String()
}

func PbTime(t time.Time) *pbts.Timestamp {
	if t.Year() < 1990 {
		return nil
	}
	ts, err := ptypes.TimestampProto(t.Truncate(time.Second))
	if err != nil {
		ll.Error("Invalid timestamp", l.Time("t", t), l.Error(err))
	}
	return ts
}

func PbTimeToModel(t *pbts.Timestamp) time.Time {
	if t == nil {
		return time.Time{}
	}
	ts, err := ptypes.Timestamp(t)
	if err != nil {
		ll.Error("Invalid timestamp")
	}
	return ts
}

func BareInt32(v *int32) int32 {
	if v == nil {
		return 0
	}
	return *v
}

func BareString(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func PbPtrInt32(v int) *int32 {
	i := int32(v)
	return &i
}

func PatchInt32(v int, p *int32) int {
	if p == nil {
		return v
	}
	return int(*p)
}

func (p *Paging) CMPaging(sorts ...string) *cm.Paging {
	if p == nil {
		return &cm.Paging{
			Limit: 100,
			Sort:  sorts,
		}
	}
	res := &cm.Paging{
		Offset: p.Offset,
		Limit:  p.Limit,
	}
	if res.Limit <= 0 {
		res.Limit = 100
	} else if res.Limit >= 1000 {
		res.Limit = 1000
	}
	if p.Sort != "" {
		res.Sort = strings.Split(p.Sort, ",")
	}
	for _, s := range sorts {
		if !cm.StringsContain(res.Sort, s) {
			res.Sort = append(res.Sort, s)
		}
	}
	return res
}

func PbPageInfo(p *cm.Paging, total int32) *PageInfo {
	return &PageInfo{
		Total: total,
		Limit: p.Limit,
		Sort:  p.Sort,
	}
}

func ToFilters(filters []*Filter) []cm.Filter {
	res := make([]cm.Filter, len(filters))
	for i, filter := range filters {
		res[i] = cm.Filter{
			Name:  filter.Name,
			Op:    filter.Op,
			Value: filter.Value,
		}
	}
	return res
}

func ToFiltersPtr(filters []*Filter) []*cm.Filter {
	res := make([]*cm.Filter, len(filters))
	for i, filter := range filters {
		res[i] = &cm.Filter{
			Name:  filter.Name,
			Op:    filter.Op,
			Value: filter.Value,
		}
	}
	return res
}

func PbError(err error) *Error {
	if err == nil {
		return &Error{
			Code: "ok",
		}
	}
	if xerr, ok := err.(*model.Error); ok {
		if xerr == nil {
			return &Error{
				Code: "ok",
			}
		}
		return &Error{
			Code: xerr.Code,
			Msg:  xerr.Msg,
			Meta: xerr.Meta,
		}
	}
	if xerr, ok := err.(*Error); ok {
		return xerr
	}
	xerr := xerrors.TwirpError(err)
	return &Error{
		Msg:  xerr.Msg(),
		Code: string(xerr.Code()),
		Meta: xerr.MetaMap(),
	}
}

func PbCustomError(err error) *Error {
	if err == nil {
		return nil
	}
	if xerr, ok := err.(*model.Error); ok {
		if xerr == nil {
			return nil
		}
		return &Error{
			Code: xerr.Code,
			Msg:  xerr.Msg,
			Meta: xerr.Meta,
		}
	}
	xerr := xerrors.TwirpError(err)
	return &Error{
		Msg:  xerr.Msg(),
		Code: string(xerr.Code()),
		Meta: xerr.MetaMap(),
	}
}

func PbErrors(errs []error) []*Error {
	res := make([]*Error, len(errs))
	for i, err := range errs {
		res[i] = PbError(err)
	}
	return res
}

func PbErrorsFromModel(errs []*model.Error) []*Error {
	res := make([]*Error, len(errs))
	for i, err := range errs {
		res[i] = PbError(err)
	}
	return res
}

func CountErrors(errs []*Error) (c int) {
	for _, err := range errs {
		if err.Code != "ok" {
			c++
		}
	}
	return c
}

func Updated(updated int) *UpdatedResponse {
	return &UpdatedResponse{
		Updated: int32(updated),
	}
}

// PbRawJSON ...
func PbRawJSON(v interface{}) *RawJSONObject {
	data, _ := json.Marshal(v)
	return &RawJSONObject{Data: data}
}

// RawJSONObjectMsg ...
func RawJSONObjectMsg(data []byte) *RawJSONObject {
	if len(data) == 0 {
		return nil
	}
	return &RawJSONObject{Data: data}
}

var _jsonEmptyObject = []byte(`{}`)

// MarshalJSONPB implements JSONPBMarshaler
func (m *RawJSONObject) MarshalJSONPB(_ *jsonpb.Marshaler) ([]byte, error) {
	if len(m.Data) == 0 {
		return _jsonEmptyObject, nil
	}
	return m.Data, nil
}

// UnmarshalJSONPB implements JSONPBUnmarshaler
func (m *RawJSONObject) UnmarshalJSONPB(_ *jsonpb.Unmarshaler, data []byte) error {
	if len(data) < 2 || data[0] != '{' || data[len(data)-1] != '}' {
		return errors.New("expect JSON object")
	}
	m.Data = data
	return nil
}

func PatchImage(sourceImages []string, cmd *model.UpdateListRequest) ([]string, error) {
	for _, imgURL := range cmd.Adds {
		if !govalidator.IsURL(imgURL) {
			return nil, cm.Error(cm.InvalidArgument, "invalid url: "+imgURL, nil)
		}
	}
	for _, imgURL := range cmd.ReplaceAll {
		if !govalidator.IsURL(imgURL) {
			return nil, cm.Error(cm.InvalidArgument, "invalid url: "+imgURL, nil)
		}
	}

	return cmd.Patch(sourceImages), nil
}

func Message(code string, msg string) *MessageResponse {
	return &MessageResponse{
		Code: code,
		Msg:  msg,
	}
}

func ErrorToModel(err *Error) *model.Error {
	if err == nil {
		return &model.Error{
			Code: "ok",
		}
	}
	return &model.Error{
		Msg:  err.Msg,
		Code: err.Code,
		Meta: err.Meta,
	}
}

func ErrorsToModel(errs []*Error) []*model.Error {
	res := make([]*model.Error, len(errs))
	for i, err := range errs {
		res[i] = ErrorToModel(err)
	}
	return res
}
