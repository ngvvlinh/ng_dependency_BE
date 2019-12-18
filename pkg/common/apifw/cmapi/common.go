package cmapi

import (
	"strings"
	"time"

	"github.com/asaskevich/govalidator"

	meta "etop.vn/api/meta"
	"etop.vn/api/top/types/common"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
	"etop.vn/common/jsonx"
	"etop.vn/common/l"
	"etop.vn/common/xerrors"
)

var ll = l.New()

func PbTime(t time.Time) dot.Time {
	return dot.Time(t)
}

func PbTimeToModel(t dot.Time) time.Time {
	return time.Time(t)
}

func BareInt(v *int) int {
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

func PbPtrInt(v int) dot.NullInt {
	return dot.Int(v)
}

func CMPaging(p *common.Paging, sorts ...string) *cm.Paging {
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

func PbPageInfo(p *cm.Paging) *common.PageInfo {
	return &common.PageInfo{
		Limit: p.Limit,
		Sort:  p.Sort,
	}
}

func PbPaging(p cm.Paging) *common.PageInfo {
	return &common.PageInfo{
		Limit: p.Limit,
		Sort:  p.Sort,
	}
}

func ToFilters(filters []*common.Filter) []cm.Filter {
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

func ToFiltersPtr(filters []*common.Filter) []*cm.Filter {
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

func PbError(err error) *common.Error {
	if err == nil {
		return &common.Error{
			Code: "ok",
		}
	}
	if xerr, ok := err.(*model.Error); ok {
		if xerr == nil {
			return &common.Error{
				Code: "ok",
			}
		}
		return &common.Error{
			Code: xerr.Code,
			Msg:  xerr.Msg,
			Meta: xerr.Meta,
		}
	}
	if xerr, ok := err.(*common.Error); ok {
		return xerr
	}
	xerr := xerrors.TwirpError(err)
	return &common.Error{
		Msg:  xerr.Msg(),
		Code: string(xerr.Code()),
		Meta: xerr.MetaMap(),
	}
}

func PbCustomError(err error) *common.Error {
	if err == nil {
		return nil
	}
	if xerr, ok := err.(*model.Error); ok {
		if xerr == nil {
			return nil
		}
		return &common.Error{
			Code: xerr.Code,
			Msg:  xerr.Msg,
			Meta: xerr.Meta,
		}
	}
	xerr := xerrors.TwirpError(err)
	return &common.Error{
		Msg:  xerr.Msg(),
		Code: string(xerr.Code()),
		Meta: xerr.MetaMap(),
	}
}

func PbErrors(errs []error) []*common.Error {
	res := make([]*common.Error, len(errs))
	for i, err := range errs {
		res[i] = PbError(err)
	}
	return res
}

func PbErrorsFromModel(errs []*model.Error) []*common.Error {
	res := make([]*common.Error, len(errs))
	for i, err := range errs {
		res[i] = PbError(err)
	}
	return res
}

func CountErrors(errs []*common.Error) (c int) {
	for _, err := range errs {
		if err.Code != "ok" {
			c++
		}
	}
	return c
}

func Updated(updated int) *common.UpdatedResponse {
	return &common.UpdatedResponse{
		Updated: updated,
	}
}

// PbRawJSON ...
func PbRawJSON(v interface{}) *common.RawJSONObject {
	data, _ := jsonx.Marshal(v)
	return &common.RawJSONObject{Data: data}
}

// RawJSONObjectMsg ...
func RawJSONObjectMsg(data []byte) *common.RawJSONObject {
	if len(data) == 0 {
		return nil
	}
	return &common.RawJSONObject{Data: data}
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

func Message(code string, msg string) *common.MessageResponse {
	return &common.MessageResponse{
		Code: code,
		Msg:  msg,
	}
}

func ErrorToModel(err *common.Error) *model.Error {
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

func ErrorsToModel(errs []*common.Error) []*model.Error {
	res := make([]*model.Error, len(errs))
	for i, err := range errs {
		res[i] = ErrorToModel(err)
	}
	return res
}

func PagingToModel(paging *common.Paging, defaultOffset int, defaultLimit int, maxLimit int) *meta.Paging {
	if defaultLimit == 0 {
		defaultLimit = 10
	}
	if defaultOffset == 0 {
		defaultOffset = 1
	}
	if maxLimit == 0 {
		maxLimit = 1000
	}
	if paging == nil {
		return &meta.Paging{
			Offset: defaultOffset,
			Limit:  defaultLimit,
		}
	}
	if paging.Limit > maxLimit {
		paging.Limit = maxLimit
	}
	arrSort := []string{paging.Sort}
	return &meta.Paging{
		Offset: paging.Offset,
		Limit:  paging.Limit,
		Sort:   arrSort,
	}
}

func CMCursorPaging(p *common.CursorPaging) (*cm.Paging, error) {
	if p == nil {
		return &cm.Paging{
			First: 100,
		}, nil
	}

	var sort []string
	if p.Sort != "" {
		sort = strings.Split(p.Sort, ",")
		if len(sort) > 1 {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "paging is invalid (Sort support only one field at the same time")
		}
		if sort[0] == "id" || sort[0] == "-id" {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "paging is invalid (Sort doesn't support \"id\"")
		}
	}

	counter := 0
	if p.First != 0 {
		if p.First < 0 || p.First > 1000 {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "paging is invalid (First must have a value between 1 and 1000")
		}
		counter++
	}
	if p.Last != 0 {
		if p.Last < 0 || p.Last > 1000 {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "paging is invalid (Last must have a value between 1 and 1000")
		}
		counter++
	}
	if p.Before != "" {
		if p.Limit != 0 {
			counter++
		} else {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "paging is invalid (Before doesn't allow Limit null)")
		}
	}
	if p.After != "" {
		if p.Limit != 0 {
			counter++
		} else {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "paging is invalid (After doesn't allow Limit null)")
		}
	}
	if counter != 1 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "paging is invalid")
	}

	res := &cm.Paging{
		First:  p.First,
		Last:   p.Last,
		Before: p.Before,
		After:  p.After,

		Limit: p.Limit,
		Sort:  sort,
	}

	return res, nil
}
