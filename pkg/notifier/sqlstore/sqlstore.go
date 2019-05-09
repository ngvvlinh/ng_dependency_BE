package sqlstore

import (
	"strings"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/l"
)

type (
	Ms    map[string]string
	Query = cmsql.Query
	Qx    = cmsql.QueryInterface
)

var (
	ll = l.New()
)

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
