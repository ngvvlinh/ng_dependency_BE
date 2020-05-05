package sqlstore

import (
	"context"
	"time"

	"o.o/api/meta"
	"o.o/api/subscripting/subscription"
	"o.o/api/top/types/etc/status3"
	"o.o/backend/com/subscripting/subscription/convert"
	"o.o/backend/com/subscripting/subscription/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sq/core"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type SubscriptionStore struct {
	ft      SubscriptionFilters
	query   func() cmsql.QueryInterface
	preds   []interface{}
	filters meta.Filters
	sqlstore.Paging
	ctx                   context.Context
	subscriptionLineStore *SubscriptionLineStore

	includeDeleted sqlstore.IncludeDeleted
}

type SubscriptionStoreFactory func(context.Context) *SubscriptionStore

func NewSubscriptionStore(db *cmsql.Database) SubscriptionStoreFactory {
	return func(ctx context.Context) *SubscriptionStore {
		return &SubscriptionStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
			subscriptionLineStore: NewSubscriptionLineStore(db)(ctx),
			ctx:                   ctx,
		}
	}
}

var scheme = conversion.Build(convert.RegisterConversions)

func (ft *SubscriptionFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

func (ft *SubscriptionFilters) NotBelongWLPartner() sq.WriterTo {
	return ft.Filter("$.wl_partner_id IS NULL")
}

func (s *SubscriptionStore) WithPaging(paging meta.Paging) *SubscriptionStore {
	s.Paging.WithPaging(paging)
	return s
}

func (s *SubscriptionStore) Filters(filters meta.Filters) *SubscriptionStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *SubscriptionStore) ID(id dot.ID) *SubscriptionStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *SubscriptionStore) Status(status status3.Status) *SubscriptionStore {
	s.preds = append(s.preds, s.ft.ByStatus(status))
	return s
}

func (s *SubscriptionStore) AccountID(id dot.ID) *SubscriptionStore {
	s.preds = append(s.preds, s.ft.ByAccountID(id))
	return s
}

func (s *SubscriptionStore) OptionalAccountID(id dot.ID) *SubscriptionStore {
	s.preds = append(s.preds, s.ft.ByAccountID(id).Optional())
	return s
}

func (s *SubscriptionStore) PlanIDs(ids ...dot.ID) *SubscriptionStore {
	s.preds = append(s.preds, sq.NewExpr("plan_ids <@ ?", core.Array{
		V: ids,
	}))
	return s
}

func (s *SubscriptionStore) GetSubscriptionFtLineDB() (*model.SubscriptionFtLine, error) {
	query := s.query().Where(s.preds).OrderBy("created_at DESC")
	query = s.ByWhiteLabelPartner(s.ctx, query)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())

	var subr model.Subscription
	if err := query.ShouldGet(&subr); err != nil {
		return nil, err
	}
	lines, err := s.subscriptionLineStore.SubscriptionID(subr.ID).ListSubscriptionLinesDB()
	if err != nil {
		return nil, err
	}
	return &model.SubscriptionFtLine{
		Subscription: &subr,
		Lines:        lines,
	}, nil
}

func (s *SubscriptionStore) GetSubscriptionFtLine() (*subscription.SubscriptionFtLine, error) {
	subr, err := s.GetSubscriptionFtLineDB()
	if err != nil {
		return nil, err
	}
	var res subscription.SubscriptionFtLine
	if err := scheme.Convert(subr, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *SubscriptionStore) ListSubscriptionFtLinesDB() ([]*model.SubscriptionFtLine, error) {
	query := s.query().Where(s.preds)
	query = s.ByWhiteLabelPartner(s.ctx, query)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	if len(s.Paging.Sort) == 0 {
		s.Paging.Sort = []string{"-created_at"}
	}
	query, err := sqlstore.LimitSort(query, &s.Paging, SortSubscription, s.ft.prefix)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterSubscription)
	if err != nil {
		return nil, err
	}

	var subrs model.Subscriptions
	if err := query.Find(&subrs); err != nil {
		return nil, err
	}
	subrIDs := make([]dot.ID, len(subrs))
	for i, subr := range subrs {
		subrIDs[i] = subr.ID
	}
	subrLines, err := s.subscriptionLineStore.SubscriptionIDs(subrIDs...).ListSubscriptionLinesDB()
	if err != nil {
		return nil, err
	}
	subrLinesMap := make(map[dot.ID][]*model.SubscriptionLine)
	for _, line := range subrLines {
		subrLinesMap[line.SubscriptionID] = append(subrLinesMap[line.SubscriptionID], line)
	}

	var res []*model.SubscriptionFtLine
	for _, subr := range subrs {
		res = append(res, &model.SubscriptionFtLine{
			Subscription: subr,
			Lines:        subrLinesMap[subr.ID],
		})
	}
	return res, nil
}

func (s *SubscriptionStore) ListSubscriptionFtLines() (res []*subscription.SubscriptionFtLine, err error) {
	subrs, err := s.ListSubscriptionFtLinesDB()
	if err != nil {
		return nil, err
	}
	if err := scheme.Convert(subrs, &res); err != nil {
		return nil, err
	}
	return
}

func (s *SubscriptionStore) CreateSubscriptionDB(subr *model.Subscription) error {
	sqlstore.MustNoPreds(s.preds)
	if subr.ID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing ID")
	}
	return s.query().ShouldInsert(subr)
}

func (s *SubscriptionStore) CreateSubscription(subr *subscription.Subscription) error {
	var subrDB model.Subscription
	if err := scheme.Convert(subr, &subrDB); err != nil {
		return err
	}
	return s.CreateSubscriptionDB(&subrDB)
}

func (s *SubscriptionStore) UpdateSubscriptionDB(subr *model.Subscription) error {
	if len(s.preds) == 0 {
		return cm.Errorf(cm.FailedPrecondition, nil, "Must provide preds")
	}
	query := s.query().Where(s.preds)
	query = s.ByWhiteLabelPartner(s.ctx, query)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	return query.ShouldUpdate(subr)
}

func (s *SubscriptionStore) UpdateSubscription(subr *subscription.Subscription) error {
	var subrDB model.Subscription
	if err := scheme.Convert(subr, &subrDB); err != nil {
		return err
	}
	return s.UpdateSubscriptionDB(&subrDB)
}

func (s *SubscriptionStore) UpdateSubscriptionStatus(id, accountID dot.ID, status status3.Status) error {
	return s.ID(id).OptionalAccountID(accountID).
		query().Table("subscription").
		Where(s.preds).
		ShouldUpdateMap(map[string]interface{}{
			"status": status.Enum(),
		})
}

func (s *SubscriptionStore) SoftDelete() (int, error) {
	query := s.query().Where(s.preds)
	query = s.ByWhiteLabelPartner(s.ctx, query)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	return query.Table("subscription").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
	})
}

func (s *SubscriptionStore) ByWhiteLabelPartner(ctx context.Context, query cmsql.Query) cmsql.Query {
	// partner := wl.X(ctx)
	// if partner.IsWhiteLabel() {
	// 	return query.Where(s.ft.ByWLPartnerID(partner.ID))
	// }
	// return query.Where(s.ft.NotBelongWLPartner())
	return query
}
