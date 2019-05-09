package shipnow

import (
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/services/shipnow/sqlstore"
)

type QueryService struct {
	shipnowStore *sqlstore.ShipnowStore
}

func NewQueryService(db cmsql.Database) *QueryService {
	return &QueryService{
		shipnowStore: sqlstore.NewShipnowStore(db),
	}
}
