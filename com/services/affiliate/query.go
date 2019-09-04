package affiliate

import (
	"context"

	"etop.vn/api/services/affiliate"
	"etop.vn/backend/com/services/affiliate/sqlstore"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/common/bus"
)

var _ affiliate.QueryService = &QueryService{}

type QueryService struct {
	etopCommissionSetting sqlstore.CommissionSettingStoreFactory
}

func NewQuery(db cmsql.Database) *QueryService {
	return &QueryService{
		etopCommissionSetting: sqlstore.NewShopCommissionSettingStore(db),
	}
}

func (a *QueryService) MessageBus() affiliate.QueryBus {
	b := bus.New()
	return affiliate.NewQueryServiceHandler(a).RegisterHandlers(b)
}

func (q *QueryService) GetCommissionByProductIDs(ctx context.Context, args *affiliate.GetCommissionByProductIDsArgs) ([]*affiliate.CommissionSetting, error) {
	return q.etopCommissionSetting(ctx).AccountID(args.AccountID).ProductIDs(args.ProductIDs).GetShopCommissionSettings()
}
