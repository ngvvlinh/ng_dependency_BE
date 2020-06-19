package sqlstore

import "o.o/backend/pkg/common/sql/sqlstore"

var FilterEquation = sqlstore.FilterWhitelist{
	Equals:   []string{"result"},
	Contains: []string{"equation"},
}
