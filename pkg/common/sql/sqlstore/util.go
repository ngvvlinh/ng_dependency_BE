package sqlstore

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/go-xorm/builder"
	"github.com/lib/pq"

	"etop.vn/api/meta"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/httpreq"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/pkg/common/sql/sq"
	"etop.vn/backend/pkg/common/sql/sq/core"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/common/strs"
	"etop.vn/common/xerrors"
)

type IncludeDeleted bool

func (includeDeleted IncludeDeleted) Check(query cmsql.Query, ft sq.WriterTo) cmsql.Query {
	if includeDeleted {
		return query
	}
	return query.Where(ft)
}

type FilterWhitelist struct {
	Arrays   []string
	Bools    []string
	Contains []string
	Dates    []string
	Equals   []string
	Nullable []string
	Numbers  []string
	Status   []string
	Unaccent []string

	PrefixOrRename map[string]string
}

func (f *FilterWhitelist) ToCol(col, suffix string) string {
	if p, ok := f.PrefixOrRename[col]; ok {
		// put a dot to rename (for example: "shop_name": "s.name")
		if strings.Contains(p, ".") {
			return p
		}
		return p + `."` + col + suffix + `" `
	}
	return `"` + col + suffix + `" `
}

type Paging struct {
	Offset int
	Limit  int
	Sort   []string

	First  int
	Last   int
	Before string
	After  string

	Next string
	Prev string
}

func (p *Paging) GetPaging() meta.PageInfo {
	if p == nil {
		return meta.PageInfo{}
	}
	return meta.PageInfo{
		Offset: p.Offset,
		Limit:  p.Limit,
		Sort:   p.Sort,
		First:  p.First,
		Last:   p.Last,
		Before: p.Before,
		After:  p.After,
		Next:   p.Next,
		Prev:   p.Prev,
	}
}

func (p *Paging) WithPaging(pg meta.Paging) {
	*p = Paging{
		Offset: pg.Offset,
		Limit:  pg.Limit,
		Sort:   pg.Sort,
		First:  pg.First,
		Last:   pg.Last,
		Before: pg.Before,
		After:  pg.After,
	}
}

func ConvertPaging(pg *meta.Paging) *Paging {
	if pg == nil {
		return nil
	}
	return &Paging{
		Offset: pg.Offset,
		Limit:  pg.Limit,
		Sort:   pg.Sort,
		First:  pg.First,
		Last:   pg.Last,
		Before: pg.Before,
		After:  pg.After,
	}
}

var cacheStructPaging = sync.Map{}

func getPagingFieldMapping(st reflect.Type) (mapping PagingFieldMapping) {
	if idxes, ok := cacheStructPaging.Load(st); ok {
		return idxes.(PagingFieldMapping)
	}
	for i, n := 0, st.NumField(); i < n; i++ {
		field := st.Field(i)
		pagingField := field.Tag.Get("paging")
		pagingColumn := parseTagForColumn(field.Tag)
		if pagingColumn == "" {
			pagingColumn = strs.ToSnake(field.Name)
		}
		if v := enumPagingFieldValue[pagingField]; v > 0 {
			mapping[v] = PagingFieldMappingItem{
				Index: i,
			}
		} else {
			mapping[v] = PagingFieldMappingItem{Index: -1}
		}
	}
	cacheStructPaging.Store(st, mapping)
	return
}

func parseTagForColumn(stag reflect.StructTag) string {
	tag := stag.Get("sq")
	if tag == "" {
		return ""
	}
	if tag[0] == '\'' {
		return tag[1 : len(tag)-1]
	}
	return ""
}

func getPrevAndNext(ss interface{}, fields []PagingCursorItem, descOrderBy bool) (prev, next string) {
	s := reflect.ValueOf(ss)
	if s.Kind() != reflect.Slice {
		panic(fmt.Sprintf("unsupported type %T (must be slice type)", s))
	}
	if s.Len() == 0 {
		return
	}

	prevCode := make([]PagingCursorItem, len(fields))
	nextCode := make([]PagingCursorItem, len(fields))
	firstPtr := reflect.Indirect(s.Index(0))
	lastPtr := reflect.Indirect(s.Index(s.Len() - 1))
	mapping := getPagingFieldMapping(firstPtr.Type())
	for i := range fields {
		field := fields[i].Field
		prevCode[i] = PagingCursorItem{field, mapping.GetFieldValue(firstPtr, field)}
		nextCode[i] = PagingCursorItem{field, mapping.GetFieldValue(lastPtr, field)}
	}
	return encodeCursor(prevCode, descOrderBy), encodeCursor(nextCode, descOrderBy)
}

func (p *Paging) Apply(s interface{}) error {
	if !p.IsCursorPaging() {
		return nil
	}
	cursor, err := p.decodeCursor()
	if err != nil {
		return err
	}
	if cursor.Reverse {
		reverseAny(s)
	}
	p.Prev, p.Next = getPrevAndNext(s, cursor.Items, cursor.DescOrderBy)
	return nil
}

func (p *Paging) IsCursorPaging() bool {
	return p.First != 0 || p.Last != 0 || p.Before != "" || p.After != ""
}

func reverseBytes(b []byte) {
	for i, n, l := 0, len(b)/2, len(b)-1; i < n; i++ {
		b[i], b[l-i] = b[l-i], b[i]
	}
}

func encodeCursor(cursorItems []PagingCursorItem, descOrderBy bool) string {
	var b bytes.Buffer
	for _, item := range cursorItems {
		desc, ok := pagingFieldDescs[item.Field]
		if !ok {
			panic(fmt.Sprintf("unexpected cursor field %v", item.Field))
		}
		b.WriteByte(byte(item.Field))
		if err := desc.Encode(&b, item.Value); err != nil {
			panic(fmt.Sprintf("unexpected cursor encoding: %v (%v)", err, item.Value))
		}
	}
	data := b.Bytes()
	if descOrderBy { // flip field sign for encoding desc order by
		data[0] = -data[0]
	}
	reverseBytes(data)
	return base64.URLEncoding.EncodeToString(data)
}

func decodeCursor(code string) (cursorItems []PagingCursorItem, descOrderBy bool, _ error) {
	data, err := base64.URLEncoding.DecodeString(code)
	if err != nil {
		return nil, false, err
	}
	if len(data) == 0 {
		return nil, false, cm.Errorf(cm.InvalidArgument, nil, "empty cursor")
	}
	reverseBytes(data)
	if int8(data[0]) < 0 { // unflip sign for decoding desc order by
		descOrderBy = true
		data[0] = -data[0]
	}
	r := bytes.NewBuffer(data)
	var hasID bool
	for r.Len() > 0 {
		b, err := r.ReadByte()
		if err != nil {
			return nil, false, cm.Errorf(cm.InvalidArgument, nil, "invalid cursor")
		}
		pagingField := PagingField(b)
		if pagingField == PagingID {
			hasID = true
		}
		desc, ok := pagingFieldDescs[pagingField]
		if !ok {
			return nil, false, cm.Errorf(cm.InvalidArgument, nil, "invalid cursor")
		}
		value, err := desc.Decode(r)
		if err != nil {
			return nil, false, cm.Errorf(cm.InvalidArgument, err, "invalid cursor")
		}
		cursorItems = append(cursorItems, PagingCursorItem{pagingField, value})
	}
	if !hasID {
		return nil, false, cm.Errorf(cm.InvalidArgument, nil, "paging is invalid")
	}

	return cursorItems, descOrderBy, nil
}

func (p *Paging) decodeCursor() (PagingCursor, error) {
	var result []PagingCursorItem
	var hasPagingFieldID bool
	if p.First != 0 || p.Last != 0 {
		var hasSort, isFirst bool
		hasSort = len(p.Sort) == 0 || !strings.HasPrefix(p.Sort[0], "-")
		isFirst = p.First != 0

		if len(p.Sort) > 0 {
			s := strings.TrimPrefix(p.Sort[0], "-")
			pagingField := ParsePagingFieldWithDefault(s, 0)
			if pagingField == PagingID {
				hasPagingFieldID = true
			}
			if pagingField == PagingUnknown {
				return PagingCursor{}, cm.Errorf(cm.InvalidArgument, nil, "sorting by %v is not supported", s)
			}
			result = append(result, PagingCursorItem{pagingField, 0})
		}
		if !hasPagingFieldID {
			result = append(result, PagingCursorItem{PagingID, 0})
		}
		return PagingCursor{
			Items:       result,
			Reverse:     !isFirst,
			DescOrderBy: hasSort != isFirst,
		}, nil
	}
	cursorItems, descOrderBy, err := decodeCursor(cm.Coalesce(p.Before, p.After))
	return PagingCursor{
		Items:       cursorItems,
		Reverse:     p.Before != "",
		DescOrderBy: descOrderBy,
	}, err
}

func (p *Paging) buildSort() error {
	if !p.IsCursorPaging() {
		return nil
	}

	result, err := p.decodeCursor()
	if err != nil {
		return err
	}
	sort := make([]string, 0, len(result.Items))
	for _, item := range result.Items {
		temp := item.Field.Name()
		if result.DescOrderBy {
			temp = "-" + temp
		}
		sort = append(sort, temp)
	}
	p.Sort = sort
	return nil
}

func (p *Paging) Validate() error {
	if p == nil {
		return nil
	}
	if p.Limit < 0 || p.Limit > 1000 {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "paging is invalid (0 < Limit < 1000)")
	}
	if len(p.Sort) > 1 {
		return cm.Errorf(cm.InvalidArgument, nil, "paging is invalid (Sort support only one field at the same time)")
	}
	counter := 0
	if p.First != 0 {
		if p.First < 0 || p.First > 1000 {
			return xerrors.Errorf(xerrors.InvalidArgument, nil, "paging is invalid (First must have a value between 0 and 1000)")
		}
		counter++
	}
	if p.Last != 0 {
		if p.Last < 0 || p.Last > 1000 {
			return xerrors.Errorf(xerrors.InvalidArgument, nil, "paging is invalid (Last must have a value between 0 and 1000)")
		}
		counter++
	}
	if p.Before != "" {
		if p.Limit == 0 {
			return xerrors.Errorf(xerrors.InvalidArgument, nil, "paging is invalid (Before must be used together with Limit)")
		}
		if len(p.Sort) > 0 {
			return xerrors.Errorf(xerrors.InvalidArgument, nil, "paging is invalid (Before must not be used together with Sort)")
		}
		counter++
	}
	if p.After != "" {
		if p.Limit == 0 {
			return xerrors.Errorf(xerrors.InvalidArgument, nil, "paging is invalid (After must be used together with Limit)")
		}
		if len(p.Sort) > 0 {
			return xerrors.Errorf(xerrors.InvalidArgument, nil, "paging is invalid (After must not be used together with Sort)")
		}
		counter++
	}
	if counter == 1 {
		return nil
	} else {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "paging is invalid")
	}
}

func validateSort(p *Paging, sortWhitelist map[string]string) error {
	if p == nil {
		return nil
	}
	if p.IsCursorPaging() {
		if len(p.Sort) > 1 {
			return cm.Errorf(cm.InvalidArgument, nil, "paging is invalid (Sort support only one field at the same time)")
		}
		for _, key := range p.Sort {
			field := strings.TrimPrefix(key, "-")
			if _, ok := sortWhitelist[field]; !ok {
				return cm.Errorf(cm.InvalidArgument, nil, "Sort by %v is not allowed", field)
			}
			if _, ok := enumPagingFieldValue[field]; !ok {
				return cm.Errorf(cm.InvalidArgument, nil, "Sort by %v is not allowed", field)
			}
		}
	}
	return nil
}

func reverseAny(s interface{}) {
	n := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}

func applyCursorPaging(
	s cmsql.Query, p *Paging,
	sortWhitelist map[string]string, prefix string,
) (cmsql.Query, error) {
	if err := validateSort(p, sortWhitelist); err != nil {
		return cmsql.Query{}, err
	}
	if err := p.Validate(); err != nil {
		return cmsql.Query{}, err
	}
	if p.First != 0 {
		p.Limit = p.First
	}
	if p.Last != 0 {
		p.Limit = p.Last
	}
	// check before, after
	if p.Before != "" || p.After != "" {
		writerTo, err := build(p, sortWhitelist, prefix)
		if err != nil {
			return cmsql.Query{}, err
		}
		s = s.Where(writerTo)
	}
	return cmsql.Query{}, nil
}

func LimitSort(s cmsql.Query, p *Paging, sortWhitelist map[string]string) (cmsql.Query, error) {
	if p == nil {
		s = s.Limit(1000)
		return s, nil
	}
	if p.IsCursorPaging() {
		query, err := applyCursorPaging(s, p, sortWhitelist, "")
		if err != nil {
			return query, err
		}
	} else {
		if p.Offset < 0 {
			return s, cm.Errorf(cm.InvalidArgument, nil, "invalid offset")
		}
		s = s.Offset(uint64(p.Offset))
	}

	if p.Limit < 0 || p.Limit > 10000 {
		return s, cm.Errorf(cm.InvalidArgument, nil, "invalid limit")
	}
	if p.Limit == 0 {
		p.Limit = 1000
	}
	s = s.Limit(uint64(p.Limit))
	if err := p.buildSort(); err != nil {
		return cmsql.Query{}, err
	}
	return Sort(s, p.Sort, sortWhitelist)
}

type CursorPagingCondition struct {
	prefix    string
	operation string
	cols      []string
	args      []interface{}
}

func (p CursorPagingCondition) WriteSQLTo(w core.SQLWriter) error {
	w.WriteByte('(')
	for _, col := range p.cols {
		if p.prefix != "" {
			w.WriteQueryString(p.prefix)
			w.WriteByte('.')
		}
		w.WriteQueryString(col)
		w.WriteByte(',')
	}
	w.TrimLast(1)
	w.WriteByte(')')

	w.WriteByte(' ')
	w.WriteQueryString(p.operation)
	w.WriteByte(' ')

	w.WriteByte('(')
	w.WriteMarkers(len(p.args))
	w.WriteByte(')')

	w.WriteArgs(p.args)
	return nil
}

func build(p *Paging, sortWhitelist map[string]string, prefix string) (sq.WriterTo, error) {
	if p == nil {
		return nil, nil
	}

	cursor, err := p.decodeCursor()
	if err != nil {
		return nil, err
	}
	condition := CursorPagingCondition{
		prefix:    prefix,
		operation: ">",
	}
	if p.Before != "" {
		condition.operation = "<"
	}

	pagingIDIndex := -1
	for i, item := range cursor.Items {
		key := item.Field.Name()
		if key == PagingID.Name() {
			pagingIDIndex = i
			continue
		}

		columnName, ok := sortWhitelist[key]
		if !ok {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "sort by %v is not allowed", key)
		}
		condition.cols = append(condition.cols, columnName)
		condition.args = append(condition.args, item.Value)
	}
	if pagingIDIndex != -1 {
		condition.cols = append(condition.cols, sortWhitelist[PagingID.Name()])
		condition.args = append(condition.args, cursor.Items[pagingIDIndex].Value)
	}

	return condition, nil
}

func PrefixedLimitSort(s cmsql.Query, p *Paging, sortWhitelist map[string]string, prefix string) (cmsql.Query, error) {
	if p == nil {
		s = s.Limit(1000)
		return s, nil
	}
	if p.IsCursorPaging() {
		query, err := applyCursorPaging(s, p, sortWhitelist, prefix)
		if err != nil {
			return query, err
		}
	} else {
		if p.Offset < 0 {
			return s, cm.Errorf(cm.InvalidArgument, nil, "invalid offset")
		}
		s = s.Offset(uint64(p.Offset))
	}

	if p.Limit < 0 || p.Limit > 10000 {
		return s, cm.Errorf(cm.InvalidArgument, nil, "invalid limit")
	}
	if p.Limit == 0 {
		p.Limit = 1000
	}
	s = s.Limit(uint64(p.Limit))
	if err := p.buildSort(); err != nil {
		return cmsql.Query{}, err
	}
	return PrefixSort(s, p.Sort, sortWhitelist, prefix)
}

func Sort(s cmsql.Query, sorts []string, whitelist map[string]string) (cmsql.Query, error) {
	for _, sort := range sorts {
		sort = strings.TrimSpace(sort)
		if sort == "" {
			continue
		}

		field := sort
		desc := ""
		if sort[0] == '-' {
			field = sort[1:]
			desc = " DESC"
		}

		if sortField, ok := whitelist[field]; ok {
			if sortField == "" {
				sortField = field
			}
			s = s.OrderBy(sortField + desc)
		} else {
			return s, cm.Errorf(cm.InvalidArgument, nil, "Sort by %v is not allowed", field)
		}
	}
	return s, nil
}

func PrefixSort(s cmsql.Query, sorts []string, whitelist map[string]string, prefix string) (cmsql.Query, error) {
	if prefix != "" && !strings.Contains(prefix, ".") {
		prefix = prefix + "."
	}
	for _, sort := range sorts {
		sort = strings.TrimSpace(sort)
		if sort == "" {
			continue
		}

		field := sort
		desc := ""
		if sort[0] == '-' {
			field = sort[1:]
			desc = " DESC"
		}

		if sortField, ok := whitelist[field]; ok {
			if sortField == "" {
				sortField = field
			}
			s = s.OrderBy(prefix + sortField + desc)
		} else {
			return s, cm.Errorf(cm.InvalidArgument, nil, "Sort by %v is not allowed", field)
		}
	}
	return s, nil
}

func Filters(s cmsql.Query, filters []cm.Filter, whitelist FilterWhitelist) (cmsql.Query, bool, error) {
	ok := false
	for _, filter := range filters {
		filter.Name = strings.TrimSpace(filter.Name)
		names := reList.Split(filter.Name, -1)
		for i, name := range names {
			if name == "" {
				return cmsql.Query{}, false, cm.Error(cm.InvalidArgument, "Invalid property name: "+filter.Name, nil)
			}

			if strings.HasPrefix(name, "s_") {
				name = "ed_" + name[2:]
			} else if strings.HasPrefix(name, "x_") {
				name = "external_" + name[2:]
			} else if strings.HasPrefix(name, "e_") {
				name = "etop_" + name[2:]
			}
			names[i] = name
		}
		ok = true

		switch filter.Op {
		case "=", "≠", "!=", "∈", "in":
			op := filter.Op
			if op == "≠" {
				op = "!="
			}
			isNumber := containsAll(whitelist.Numbers, names)
			isStatus := containsAll(whitelist.Status, names)
			isString := containsAll(whitelist.Equals, names)
			isNullable := containsAll(whitelist.Nullable, names)
			isBool := containsAll(whitelist.Bools, names)
			isArray := containsAll(whitelist.Arrays, names)
			if !isString && !isNumber && !isStatus && !isBool && !isNullable && !isArray {
				return cmsql.Query{}, false, cm.Error(cm.InvalidArgument, "Exactly filter is not allowed for "+filter.Name, nil)
			}
			if countBool(isNumber, isStatus, isString, isBool, isNullable, isArray) != 1 {
				return cmsql.Query{}, false, cm.Error(cm.InvalidArgument, "Exactly filter must contain the same type "+filter.Name, nil)
			}

			if isArray && op != "=" && op != "!=" && op != "≠" {
				return cmsql.Query{}, false, cm.Error(cm.InvalidArgument, "Array not support for "+op+" operation", nil)
			}
			if isArray && filter.Value != "{}" && filter.Value != "Ø" {
				return cmsql.Query{}, false, cm.Error(cm.InvalidArgument, "Array support for "+op+" operation only when value is {} or Ø", nil)
			}
			if isArray {
				value := "{}"
				if filter.Op == "=" {
					s = buildQuery(s, names, value, simpleQuery(func(name string) string {
						return whitelist.ToCol(name, "") + ` = ? OR ` + whitelist.ToCol(name, "") + " IS NULL"
					}))
				} else {
					s = buildQuery(s, names, value, simpleQuery(func(name string) string {
						return whitelist.ToCol(name, "") + ` != ? AND ` + whitelist.ToCol(name, "") + " IS NOT NULL"
					}))
				}
				break
			}

			cfg := valueConfig{isNumber: isNumber, isStatus: isStatus, isBool: isBool, isNullable: isNullable}
			if op == "=" || op == "!=" {
				value, err := parseValue(filter.Value, cfg)
				if err != nil {
					return cmsql.Query{}, false, err
				}
				if isNullable {
					s = buildQuery(s, names, nil, simpleQuery(func(name string) string {
						query := "IS NULL"
						if value.(bool) {
							query = "IS NOT NULL"
						}
						return whitelist.ToCol(name, "") + ` ` + query
					}))
				} else {
					s = buildQuery(s, names, value, simpleQuery(func(name string) string {
						return whitelist.ToCol(name, "") + op + ` ?`
					}))
				}

			} else {
				if isBool || isNullable {
					return cmsql.Query{}, false, cm.Errorf(cm.InvalidArgument, nil, "Element filter is not allowed for "+filter.Name)
				}
				value, err := parseValueAsList(filter.Value, cfg)
				if err != nil {
					return cmsql.Query{}, false, err
				}
				ors := make([]builder.Cond, len(names))
				for i, name := range names {
					col := whitelist.ToCol(name, "")
					ors[i] = builder.In(col, value)
				}
				sql, args, err := builder.ToSQL(builder.Or(ors...))
				if err != nil {
					return cmsql.Query{}, false, err
				}
				s = s.Where(combineArgs(sql, args)...)
			}

		case "<", ">", "<=", ">=", "≤", "≥":
			op := filter.Op
			if op == "≤" {
				op = "<="
			} else if op == "≥" {
				op = ">="
			}
			isNumber := containsAll(whitelist.Numbers, names)
			isDate := containsAll(whitelist.Dates, names)
			if !isNumber && !isDate {
				return cmsql.Query{}, false, cm.Error(cm.InvalidArgument, "Compare filter is not allowed for "+filter.Name, nil)
			}
			if countBool(isNumber, isDate) != 1 {
				return cmsql.Query{}, false, cm.Error(cm.InvalidArgument, "Compare filter must contain the same type "+filter.Name, nil)
			}

			value, err := parseValue(filter.Value, valueConfig{isNumber: isNumber, isDate: isDate})
			if err != nil {
				return cmsql.Query{}, false, err
			}
			s = buildQuery(s, names, value, simpleQuery(func(name string) string {
				return whitelist.ToCol(name, "") + ` ` + op + ` ?`
			}))

		case "⊃", "c", "∩", "n":
			isContains := filter.Op == "⊃" || filter.Op == "c"
			isText := containsAll(whitelist.Contains, names)
			isArray := containsAll(whitelist.Arrays, names)
			if !isText && !isArray {
				return cmsql.Query{}, false, cm.Error(cm.InvalidArgument, "Contains or intersect filter is not allowed for "+filter.Name, nil)
			}

			if isText {
				var value string
				if isContains {
					value = validate.NormalizeSearchQueryAnd(filter.Value)
				} else {
					value = validate.NormalizeSearchQueryOr(filter.Value)
				}

				s = buildQuery(s, names, value, simpleQuery(func(name string) string {
					return whitelist.ToCol(name, "_norm") + ` @@ ?::tsquery`
				}))

			} else {
				value, err := parseValueAsList(filter.Value, valueConfig{})
				if err != nil {
					return cmsql.Query{}, false, err
				}

				var op string
				if isContains {
					op = "@>"
				} else {
					op = "&&"
				}
				s = buildQuery(s, names, pq.Array(value), simpleQuery(func(name string) string {
					return whitelist.ToCol(name, "") + ` ` + op + ` ?`
				}))
			}

		case "~=", "≃":
			isUnaccent := containsAll(whitelist.Unaccent, names)
			if !isUnaccent {
				return cmsql.Query{}, false, cm.Errorf(cm.InvalidArgument, nil, "Almost equal is not allowed for %v", filter.Name)
			}
			value := validate.NormalizeUnaccent(filter.Value)
			s = buildQuery(s, names, value, simpleQuery(func(name string) string {
				return whitelist.ToCol(name, "_norm_ua") + ` = ?`
			}))

		default:
			return cmsql.Query{}, false, cm.Error(cm.InvalidArgument, "Invalid filter operation", nil)
		}
	}
	return s, ok, nil
}

func combineArgs(sql string, args []interface{}) []interface{} {
	res := make([]interface{}, len(args)+1)
	res[0] = sql
	copy(res[1:], args)
	return res
}

func simpleQuery(fn func(name string) string) func(value interface{}, name string) (string, interface{}) {
	return func(value interface{}, name string) (string, interface{}) {
		return fn(name), value
	}
}

func buildQuery(s cmsql.Query, names []string, value interface{}, fn func(value interface{}, name string) (string, interface{})) cmsql.Query {
	if len(names) == 1 {
		_expr, _args := fn(value, names[0])
		if _args == nil {
			return s.Where(_expr)
		}
		return s.Where(_expr, _args)
	}

	buf := make([]byte, 0, 64)
	args := make([]interface{}, 0, len(names))
	for i, name := range names {
		if i > 0 {
			buf = append(buf, " OR "...)
		}
		_expr, _args := fn(value, name)
		buf = append(buf, _expr...)
		if _args != nil {
			args = append(args, _args)
		}
	}
	s = s.Where(combineArgs(string(buf), args)...)
	return s
}

func countBool(A ...bool) int {
	c := 0
	for _, a := range A {
		if a {
			c++
		}
	}
	return c
}

type valueConfig struct {
	isNumber   bool
	isDate     bool
	isStatus   bool
	isBool     bool
	isNullable bool
}

func parseValue(v string, cfg valueConfig) (interface{}, error) {
	if cfg.isNullable {
		n, err := strconv.ParseBool(v)
		if err != nil {
			return 0, cm.Error(cm.InvalidArgument, "Invalid bool: "+v, nil)
		}
		// nullable will be handled specially at caller
		return n, nil
	}
	if cfg.isBool {
		n, err := strconv.ParseBool(v)
		if err != nil {
			return 0, cm.Error(cm.InvalidArgument, "Invalid bool: "+v, nil)
		}
		return n, nil
	}
	if cfg.isNumber {
		n, err := strconv.Atoi(v)
		if err != nil {
			return 0, cm.Error(cm.InvalidArgument, "Invalid number: "+v, nil)
		}
		return n, nil
	}
	if cfg.isDate {
		t, ok := httpreq.ParseAsISO8601([]byte(v))
		if !ok {
			return 0, cm.Error(cm.InvalidArgument, "Invalid date: "+v, nil)
		}
		return t, nil
	}
	if cfg.isStatus {
		switch v {
		case "P", "1":
			return 1, nil
		case "Z", "0":
			return 0, nil
		case "N", "-1":
			return -1, nil
		case "S", "2":
			return 2, nil
		case "NS", "-2":
			return -2, nil
		}
	}
	return v, nil
}

var reList = regexp.MustCompile(`[\s,]+`)

func parseValueAsList(v string, cfg valueConfig) ([]interface{}, error) {
	ss := reList.Split(v, -1)
	res := make([]interface{}, 0, len(ss))
	for _, s := range ss {
		if s == "" {
			continue
		}
		resi, err := parseValue(s, cfg)
		if err != nil {
			return nil, err
		}
		res = append(res, resi)
	}
	if len(res) == 0 {
		return nil, cm.Error(cm.InvalidArgument, "Empty value", nil)
	}
	return res, nil
}

func contains(ss []string, item string) bool {
	for _, s := range ss {
		if s == item {
			return true
		}
	}
	return false
}

func containsAll(ss []string, items []string) bool {
	for _, item := range items {
		if !contains(ss, item) {
			return false
		}
	}
	return true
}

func containsID(ss []int64, item int64) bool {
	for _, s := range ss {
		if s == item {
			return true
		}
	}
	return false
}

func MustNoPreds(preds []interface{}) {
	if len(preds) != 0 {
		panic("must provide no preds")
	}
}
