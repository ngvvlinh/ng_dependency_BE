package sqlstore

import (
	"context"
	"time"

	"o.o/api/subscripting/subscriptionproduct"
	"o.o/backend/com/subscripting/subscriptionproduct/convert"
	"o.o/backend/com/subscripting/subscriptionproduct/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type SubrProductStore struct {
	ft    SubscriptionProductFilters
	query func() cmsql.QueryInterface
	preds []interface{}
	ctx   context.Context

	includeDeleted sqlstore.IncludeDeleted
}

type SubrProductStoreFactory func(context.Context) *SubrProductStore

func NewSubscriptionProductStore(db *cmsql.Database) SubrProductStoreFactory {
	return func(ctx context.Context) *SubrProductStore {
		return &SubrProductStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
			ctx: ctx,
		}
	}
}

var scheme = conversion.Build(convert.RegisterConversions)

func (ft *SubscriptionProductFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

func (ft *SubscriptionProductFilters) NotBelongWLPartner() sq.WriterTo {
	return ft.Filter("$.wl_partner_id IS NULL")
}

func (s *SubrProductStore) ID(id dot.ID) *SubrProductStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *SubrProductStore) GetSubrProductDB() (*model.SubscriptionProduct, error) {
	query := s.query().Where(s.preds)
	query = s.ByWhiteLabelPartner(s.ctx, query)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	var res model.SubscriptionProduct
	if err := query.ShouldGet(&res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *SubrProductStore) GetSubrProduct() (*subscriptionproduct.SubscriptionProduct, error) {
	subr, err := s.GetSubrProductDB()
	if err != nil {
		return nil, err
	}
	var res subscriptionproduct.SubscriptionProduct
	if err := scheme.Convert(subr, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *SubrProductStore) ListSubscriptionsDB() ([]*model.SubscriptionProduct, error) {
	query := s.query().Where(s.preds).OrderBy("created_at DESC")
	query = s.ByWhiteLabelPartner(s.ctx, query)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())

	var res model.SubscriptionProducts
	if err := query.Find(&res); err != nil {
		return nil, err
	}
	return res, nil
}

func (s *SubrProductStore) ListSubscriptions() (res []*subscriptionproduct.SubscriptionProduct, err error) {
	subrs, err := s.ListSubscriptionsDB()
	if err != nil {
		return nil, err
	}
	if err := scheme.Convert(subrs, &res); err != nil {
		return nil, err
	}
	return
}

func (s *SubrProductStore) CreateSubscriptionDB(subrProduct *model.SubscriptionProduct) error {
	sqlstore.MustNoPreds(s.preds)
	if subrProduct.ID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing ID")
	}
	subrProduct.WLPartnerID = wl.GetWLPartnerID(s.ctx)
	return s.query().ShouldInsert(subrProduct)
}

func (s *SubrProductStore) CreateSubscription(subrProduct *subscriptionproduct.SubscriptionProduct) error {
	var subrDB model.SubscriptionProduct
	if err := scheme.Convert(subrProduct, &subrDB); err != nil {
		return err
	}
	return s.CreateSubscriptionDB(&subrDB)
}

func (s *SubrProductStore) UpdateSubrProductDB(subrProduct *model.SubscriptionProduct) error {
	if len(s.preds) == 0 {
		return cm.Errorf(cm.FailedPrecondition, nil, "Must provide preds")
	}
	query := s.query().Where(s.preds)
	query = s.ByWhiteLabelPartner(s.ctx, query)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	return query.ShouldUpdate(subrProduct)
}

func (s *SubrProductStore) UpdateSubrProduct(subrProduct *subscriptionproduct.SubscriptionProduct) error {
	var subrDB model.SubscriptionProduct
	if err := scheme.Convert(subrProduct, &subrDB); err != nil {
		return err
	}
	return s.UpdateSubrProductDB(&subrDB)
}

func (s *SubrProductStore) SoftDelete() (int, error) {
	query := s.query().Where(s.preds)
	query = s.ByWhiteLabelPartner(s.ctx, query)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	return query.Table("subscription_product").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
	})
}

func (s *SubrProductStore) ByWhiteLabelPartner(ctx context.Context, query cmsql.Query) cmsql.Query {
	partner := wl.X(ctx)
	if partner.IsWhiteLabel() {
		return query.Where(s.ft.ByWLPartnerID(partner.ID))
	}
	return query.Where(s.ft.NotBelongWLPartner())
}
