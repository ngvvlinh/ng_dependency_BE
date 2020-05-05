package sqlstore

import (
	"context"
	"time"

	"o.o/api/subscripting/subscriptionplan"
	"o.o/backend/com/subscripting/subscriptionplan/convert"
	"o.o/backend/com/subscripting/subscriptionplan/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type SubrPlanStore struct {
	ft    SubscriptionPlanFilters
	query func() cmsql.QueryInterface
	preds []interface{}
	ctx   context.Context

	includeDeleted sqlstore.IncludeDeleted
}

type SubrPlanStoreFactory func(context.Context) *SubrPlanStore

func NewSubrPlanStore(db *cmsql.Database) SubrPlanStoreFactory {
	return func(ctx context.Context) *SubrPlanStore {
		return &SubrPlanStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
			ctx: ctx,
		}
	}
}

var scheme = conversion.Build(convert.RegisterConversions)

func (ft *SubscriptionPlanFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

func (ft *SubscriptionPlanFilters) NotBelongWLPartner() sq.WriterTo {
	return ft.Filter("$.wl_partner_id IS NULL")
}

func (s *SubrPlanStore) ID(id dot.ID) *SubrPlanStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *SubrPlanStore) ProductIDs(ids ...dot.ID) *SubrPlanStore {
	s.preds = append(s.preds, sq.In("product_id", ids))
	return s
}

func (s *SubrPlanStore) FreePlan() *SubrPlanStore {
	freePrice := 0
	s.preds = append(s.preds, s.ft.ByPricePtr(&freePrice))
	return s
}

func (s *SubrPlanStore) GetSubrPlanDB() (*model.SubscriptionPlan, error) {
	query := s.query().Where(s.preds).OrderBy("created_at DESC")
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	var res model.SubscriptionPlan
	if err := query.ShouldGet(&res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *SubrPlanStore) GetSubrPlan() (*subscriptionplan.SubscriptionPlan, error) {
	subr, err := s.GetSubrPlanDB()
	if err != nil {
		return nil, err
	}
	var res subscriptionplan.SubscriptionPlan
	if err := scheme.Convert(subr, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *SubrPlanStore) ListSubscriptionsDB() ([]*model.SubscriptionPlan, error) {
	query := s.query().Where(s.preds).OrderBy("created_at DESC")
	query = s.ByWhiteLabelPartner(s.ctx, query)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())

	var res model.SubscriptionPlans
	if err := query.Find(&res); err != nil {
		return nil, err
	}
	return res, nil
}

func (s *SubrPlanStore) ListSubscriptions() (res []*subscriptionplan.SubscriptionPlan, err error) {
	subrs, err := s.ListSubscriptionsDB()
	if err != nil {
		return nil, err
	}
	if err := scheme.Convert(subrs, &res); err != nil {
		return nil, err
	}
	return
}

func (s *SubrPlanStore) CreateSubscriptionDB(SubrPlan *model.SubscriptionPlan) error {
	sqlstore.MustNoPreds(s.preds)
	if SubrPlan.ID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing ID")
	}
	SubrPlan.WLPartnerID = wl.GetWLPartnerID(s.ctx)
	return s.query().ShouldInsert(SubrPlan)
}

func (s *SubrPlanStore) CreateSubscription(SubrPlan *subscriptionplan.SubscriptionPlan) error {
	var subrDB model.SubscriptionPlan
	if err := scheme.Convert(SubrPlan, &subrDB); err != nil {
		return err
	}
	return s.CreateSubscriptionDB(&subrDB)
}

func (s *SubrPlanStore) UpdateSubrPlanDB(SubrPlan *model.SubscriptionPlan) error {
	if len(s.preds) == 0 {
		return cm.Errorf(cm.FailedPrecondition, nil, "Must provide preds")
	}
	query := s.query().Where(s.preds)
	query = s.ByWhiteLabelPartner(s.ctx, query)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	return query.ShouldUpdate(SubrPlan)
}

func (s *SubrPlanStore) UpdateSubrPlan(subrPlan *subscriptionplan.SubscriptionPlan) error {
	var subrDB model.SubscriptionPlan
	if err := scheme.Convert(subrPlan, &subrDB); err != nil {
		return err
	}
	return s.UpdateSubrPlanDB(&subrDB)
}

func (s *SubrPlanStore) SoftDelete() (int, error) {
	query := s.query().Where(s.preds)
	query = s.ByWhiteLabelPartner(s.ctx, query)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	return query.Table("subscription_plan").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
	})
}

func (s *SubrPlanStore) ByWhiteLabelPartner(ctx context.Context, query cmsql.Query) cmsql.Query {
	// partner := wl.X(ctx)
	// if partner.IsWhiteLabel() {
	// 	return query.Where(s.ft.ByWLPartnerID(partner.ID))
	// }
	// return query.Where(s.ft.NotBelongWLPartner())
	return query
}
