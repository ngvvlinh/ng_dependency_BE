package sqlstore

import (
	"context"
	"time"

	"o.o/api/meta"
	"o.o/api/subscripting/subscriptionbill"
	"o.o/backend/com/subscripting/subscriptionbill/convert"
	"o.o/backend/com/subscripting/subscriptionbill/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type SubrBillStore struct {
	ft      SubscriptionBillFilters
	query   func() cmsql.QueryInterface
	preds   []interface{}
	filters meta.Filters
	sqlstore.Paging
	ctx               context.Context
	subrBillLineStore *SubrBillLineStore

	includeDeleted sqlstore.IncludeDeleted
}

type SubrBillStoreFactory func(context.Context) *SubrBillStore

func NewSubrBillStore(db *cmsql.Database) SubrBillStoreFactory {
	return func(ctx context.Context) *SubrBillStore {
		return &SubrBillStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
			subrBillLineStore: NewSubrBillLineStore(db)(ctx),
			ctx:               ctx,
		}
	}
}

var scheme = conversion.Build(convert.RegisterConversions)

func (ft *SubscriptionBillFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

func (ft *SubscriptionBillFilters) NotBelongWLPartner() sq.WriterTo {
	return ft.Filter("$.wl_partner_id IS NULL")
}

func (s *SubrBillStore) WithPaging(paging meta.Paging) *SubrBillStore {
	s.Paging.WithPaging(paging)
	return s
}

func (s *SubrBillStore) Filters(filters meta.Filters) *SubrBillStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *SubrBillStore) ID(id dot.ID) *SubrBillStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *SubrBillStore) AccountID(id dot.ID) *SubrBillStore {
	s.preds = append(s.preds, s.ft.ByAccountID(id))
	return s
}

func (s *SubrBillStore) OptionalAccountID(id dot.ID) *SubrBillStore {
	s.preds = append(s.preds, s.ft.ByAccountID(id).Optional())
	return s
}

func (s *SubrBillStore) OptionalSubscriptionID(id dot.ID) *SubrBillStore {
	s.preds = append(s.preds, s.ft.BySubscriptionID(id).Optional())
	return s
}

func (s *SubrBillStore) GetSubrBillFtLineDB() (*model.SubscriptionBillFtLine, error) {
	query := s.query().Where(s.preds)
	query = s.ByWhiteLabelPartner(s.ctx, query)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	var subrBill model.SubscriptionBill
	if err := query.ShouldGet(&subrBill); err != nil {
		return nil, err
	}

	lines, err := s.subrBillLineStore.SubscriptionBillID(subrBill.ID).ListSubrBillLinesDB()
	if err != nil {
		return nil, err
	}

	return &model.SubscriptionBillFtLine{
		SubscriptionBill: &subrBill,
		Lines:            lines,
	}, nil
}

func (s *SubrBillStore) GetSubrBillFtLine() (*subscriptionbill.SubscriptionBillFtLine, error) {
	subr, err := s.GetSubrBillFtLineDB()
	if err != nil {
		return nil, err
	}
	var res subscriptionbill.SubscriptionBillFtLine
	if err := scheme.Convert(subr, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *SubrBillStore) ListSubrBillFtLinesDB() ([]*model.SubscriptionBillFtLine, error) {
	query := s.query().Where(s.preds)
	query = s.ByWhiteLabelPartner(s.ctx, query)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	if len(s.Paging.Sort) == 0 {
		s.Paging.Sort = []string{"-created_at"}
	}
	query, err := sqlstore.LimitSort(query, &s.Paging, SortSubscriptionBill, s.ft.prefix)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterSubscriptionBill)
	if err != nil {
		return nil, err
	}

	var subrBills model.SubscriptionBills
	if err := query.Find(&subrBills); err != nil {
		return nil, err
	}
	subrBillIDs := make([]dot.ID, len(subrBills))
	for i, subrBill := range subrBills {
		subrBillIDs[i] = subrBill.ID
	}
	subrBillLines, err := s.subrBillLineStore.SubscriptionBillIDs(subrBillIDs...).ListSubrBillLinesDB()
	if err != nil {
		return nil, err
	}
	subrBillLinesMap := make(map[dot.ID][]*model.SubscriptionBillLine)
	for _, line := range subrBillLines {
		subrBillLinesMap[line.SubscriptionBillID] = append(subrBillLinesMap[line.SubscriptionBillID], line)
	}

	var res []*model.SubscriptionBillFtLine
	for _, subrBill := range subrBills {
		res = append(res, &model.SubscriptionBillFtLine{
			SubscriptionBill: subrBill,
			Lines:            subrBillLinesMap[subrBill.ID],
		})
	}
	return res, nil
}

func (s *SubrBillStore) ListSubrBillFtLines() (res []*subscriptionbill.SubscriptionBillFtLine, err error) {
	subrs, err := s.ListSubrBillFtLinesDB()
	if err != nil {
		return nil, err
	}
	if err := scheme.Convert(subrs, &res); err != nil {
		return nil, err
	}
	return
}

func (s *SubrBillStore) CreateSubrBillDB(subrBill *model.SubscriptionBill) error {
	sqlstore.MustNoPreds(s.preds)
	if subrBill.ID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing ID")
	}
	subrBill.WLPartnerID = wl.GetWLPartnerID(s.ctx)
	return s.query().ShouldInsert(subrBill)
}

func (s *SubrBillStore) CreateSubrBill(SubrPlan *subscriptionbill.SubscriptionBill) error {
	var subrDB model.SubscriptionBill
	if err := scheme.Convert(SubrPlan, &subrDB); err != nil {
		return err
	}
	return s.CreateSubrBillDB(&subrDB)
}

func (s *SubrBillStore) UpdateSubrBillDB(subrBill *model.SubscriptionBill) error {
	if len(s.preds) == 0 {
		return cm.Errorf(cm.FailedPrecondition, nil, "Must provide preds")
	}
	query := s.query().Where(s.preds)
	query = s.ByWhiteLabelPartner(s.ctx, query)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	return query.ShouldUpdate(subrBill)
}

func (s *SubrBillStore) SoftDelete() (int, error) {
	query := s.query().Where(s.preds)
	query = s.ByWhiteLabelPartner(s.ctx, query)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	return query.Table("subscription_bill").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
	})
}

func (s *SubrBillStore) ByWhiteLabelPartner(ctx context.Context, query cmsql.Query) cmsql.Query {
	partner := wl.X(ctx)
	if partner.IsWhiteLabel() {
		return query.Where(s.ft.ByWLPartnerID(partner.ID))
	}
	return query.Where(s.ft.NotBelongWLPartner())
}
