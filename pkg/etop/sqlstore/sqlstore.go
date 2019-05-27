package sqlstore

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-xorm/builder"
	"github.com/lib/pq"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/httpreq"
	"etop.vn/backend/pkg/common/l"
	sq "etop.vn/backend/pkg/common/sql"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/backend/pkg/etop/model"
	notisqlstore "etop.vn/backend/pkg/notifier/sqlstore"
)

var (
	x                 cmsql.Database
	xNotifier         cmsql.Database
	ll                = l.New()
	deviceStore       *notisqlstore.DeviceStore
	notificationStore *notisqlstore.NotificationStore
)

type (
	M  map[string]interface{}
	Ms map[string]string

	Query = cmsql.Query
	Qx    = cmsql.QueryInterface
)

func Init(db cmsql.Database) {
	if x.DB() != nil {
		ll.Panic("Already initialized")
	}
	x = db
}

func InitDBNotifier(db cmsql.Database) {
	if xNotifier.DB() != nil {
		ll.Panic("Database Notifier already initialized")
	}
	xNotifier = db
	deviceStore = notisqlstore.NewDeviceStore(xNotifier)
	notificationStore = notisqlstore.NewNotificationStore(xNotifier)
}

type filterDeletable interface {
	Prefix() string
	ByDeletedAt(time.Time) *sq.ColumnFilter
}

type includeDeleted bool

func (d includeDeleted) filterDeleted(f filterDeletable) sq.WriterTo {
	if d {
		return nil
	}
	s := "deleted_at IS NULL"
	p := f.Prefix()
	if p != "" {
		s = p + "." + s
	}
	return sq.NewExpr(s)
}

type multiplelity bool

func (m multiplelity) ensureMultiplelity(countable interface{ Count() (uint64, error) }) error {
	n, err := countable.Count()
	if err != nil {
		return err
	}
	if !m && (n > 1) {
		return cm.Errorf(cm.Internal, nil, "unexpected number of changes")
	}
	return nil
}

func inTransaction(callback func(cmsql.QueryInterface) error) (err error) {
	return x.InTransaction(callback)
}

func LimitSort(s Query, p *cm.Paging, sortWhitelist map[string]string) (Query, error) {
	if p == nil {
		s = s.Limit(10000)
		return s, nil
	}
	if p.Limit <= 0 {
		p.Limit = 10000
	}
	s = s.Limit(uint64(p.Limit)).Offset(uint64(p.Offset))
	return Sort(s, p.Sort, sortWhitelist)
}

func IDs(items []int64) []interface{} {
	res := make([]interface{}, len(items))
	for i, item := range items {
		res[i] = item
	}
	return res
}

func FilterStatus(s Query, prefix string, query model.StatusQuery) Query {
	if query.Status != nil {
		s = s.Where(prefix+"status = ?", query.Status)
	}
	if query.EtopStatus != nil {
		s = s.Where(prefix+"etop_status = ?", query.EtopStatus)
	}
	return s
}

func Sort(s Query, sorts []string, whitelist map[string]string) (Query, error) {
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

func Filters(s Query, filters []cm.Filter, whitelist FilterWhitelist) (Query, bool, error) {
	ok := false
	for _, filter := range filters {
		filter.Name = strings.TrimSpace(filter.Name)
		names := reList.Split(filter.Name, -1)
		for i, name := range names {
			if name == "" {
				return Query{}, false, cm.Error(cm.InvalidArgument, "Invalid property name: "+filter.Name, nil)
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
			if !isString && !isNumber && !isStatus && !isBool && !isNullable {
				return Query{}, false, cm.Error(cm.InvalidArgument, "Exactly filter is not allowed for "+filter.Name, nil)
			}
			if countBool(isNumber, isStatus, isString, isBool, isNullable) != 1 {
				return Query{}, false, cm.Error(cm.InvalidArgument, "Exactly filter must contain the same type "+filter.Name, nil)
			}

			cfg := valueConfig{isNumber: isNumber, isStatus: isStatus, isBool: isBool, isNullable: isNullable}
			if op == "=" || op == "!=" {
				value, err := parseValue(filter.Value, cfg)
				if err != nil {
					return Query{}, false, err
				}
				if isNullable {
					s = buildQuery(s, names, value, func(name string) string {
						query := "IS NULL"
						if value.(bool) {
							query = "IS NOT NULL"
						}
						return whitelist.ToCol(name, "") + query
					})
				} else {
					s = buildQuery(s, names, value, func(name string) string {
						return whitelist.ToCol(name, "") + op + ` ?`
					})
				}

			} else {
				if isBool || isNullable {
					return Query{}, false, cm.Errorf(cm.InvalidArgument, nil, "Element filter is not allowed for "+filter.Name)
				}
				value, err := parseValueAsList(filter.Value, cfg)
				if err != nil {
					return Query{}, false, err
				}
				ors := make([]builder.Cond, len(names))
				for i, name := range names {
					col := whitelist.ToCol(name, "")
					ors[i] = builder.In(col, value)
				}
				sql, args, err := builder.ToSQL(builder.Or(ors...))
				if err != nil {
					return Query{}, false, err
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
				return Query{}, false, cm.Error(cm.InvalidArgument, "Compare filter is not allowed for "+filter.Name, nil)
			}
			if countBool(isNumber, isDate) != 1 {
				return Query{}, false, cm.Error(cm.InvalidArgument, "Compare filter must contain the same type "+filter.Name, nil)
			}

			value, err := parseValue(filter.Value, valueConfig{isNumber: isNumber, isDate: isDate})
			if err != nil {
				return Query{}, false, err
			}
			s = buildQuery(s, names, value, func(name string) string {
				return whitelist.ToCol(name, "") + ` ` + op + ` ?`
			})

		case "⊃", "c", "∩", "n":
			isContains := filter.Op == "⊃" || filter.Op == "c"
			isText := containsAll(whitelist.Contains, names)
			isArray := containsAll(whitelist.Arrays, names)
			if !isText && !isArray {
				return Query{}, false, cm.Error(cm.InvalidArgument, "Contains or intersect filter is not allowed for "+filter.Name, nil)
			}

			if isText {
				var value string
				if isContains {
					value = validate.NormalizeSearchQueryAnd(filter.Value)
				} else {
					value = validate.NormalizeSearchQueryOr(filter.Value)
				}

				s = buildQuery(s, names, value, func(name string) string {
					return whitelist.ToCol(name, "_norm") + ` @@ ?::tsquery`
				})

			} else {
				value, err := parseValueAsList(filter.Value, valueConfig{})
				if err != nil {
					return Query{}, false, err
				}

				var op string
				if isContains {
					op = "@>"
				} else {
					op = "&&"
				}
				s = buildQuery(s, names, pq.Array(value), func(name string) string {
					return whitelist.ToCol(name, "") + ` ` + op + ` ?`
				})
			}

		case "~=", "≃":
			isUnaccent := containsAll(whitelist.Unaccent, names)
			if !isUnaccent {
				return Query{}, false, cm.Errorf(cm.InvalidArgument, nil, "Almost equal is not allowed for %v", filter.Name)
			}
			value := validate.NormalizeUnaccent(filter.Value)
			s = buildQuery(s, names, value, func(name string) string {
				return whitelist.ToCol(name, "_norm_ua") + ` = ?`
			})

		default:
			return Query{}, false, cm.Error(cm.InvalidArgument, "Invalid filter operation", nil)
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

func buildQuery(s Query, names []string, value interface{}, fn func(name string) string) Query {
	if len(names) == 1 {
		return s.Where(fn(names[0]), value)
	}

	buf := make([]byte, 0, 64)
	args := make([]interface{}, len(names))
	for i, name := range names {
		if i > 0 {
			buf = append(buf, " OR "...)
		}
		buf = append(buf, fn(name)...)
		args[i] = value
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

func sqlPlaceholder(counter int) (*int, func(sql []byte) []byte) {
	fn := func(sql []byte) []byte {
		counter++
		sql = append(sql, ",$"...)
		sql = strconv.AppendInt(sql, int64(counter), 10)
		return sql
	}
	return &counter, fn
}

func appendInt(b []byte, i int) []byte {
	return strconv.AppendInt(b, int64(i), 10)
}

func appendInt64(b []byte, i int64) []byte {
	return strconv.AppendInt(b, i, 10)
}

func ignoreError(err error) {
}
