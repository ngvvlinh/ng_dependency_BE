package sqlstore

import (
	"context"
	"time"

	"o.o/api/main/credit"
	"o.o/api/top/types/etc/status3"
	"o.o/backend/com/main/credit/convert"
	"o.o/backend/com/main/credit/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sq/core"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

var Sort = map[string]string{
	"updated_at": "updated_at",
	"created_at": "created_at",
}

func (ft *CreditFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

var scheme = conversion.Build(convert.RegisterConversions)

type CreditFactory func(context.Context) *CreditStore

func NewCreditStore(db *cmsql.Database) CreditFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *CreditStore {
		return &CreditStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type CreditStore struct {
	query cmsql.QueryFactory
	ft    CreditFilters
	sqlstore.Paging
	preds []interface{}

	includeDeleted sqlstore.IncludeDeleted
}

func (s *CreditStore) ID(id dot.ID) *CreditStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *CreditStore) WithPaging(paging *cm.Paging) *CreditStore {
	s.Paging.WithPaging(*paging)
	return s
}

func (s *CreditStore) IDs(ids ...dot.ID) *CreditStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "variant_id", ids))
	return s
}

func (s *CreditStore) ShopID(id dot.ID) *CreditStore {
	s.preds = append(s.preds, s.ft.ByShopID(id))
	return s
}

func (s *CreditStore) UpdateCreditAllDB(args *model.Credit) error {
	query := s.query().Where(s.preds)
	return query.UpdateAll().ShouldUpdate(args)
}

func (s *CreditStore) UpdateCreditAll(args *credit.Credit) error {
	creditCore := &model.Credit{}
	if err := scheme.Convert(args, creditCore); err != nil {
		return err
	}
	return s.UpdateCreditAllDB(creditCore)
}

func (s *CreditStore) Create(credit *credit.Credit) error {
	sqlstore.MustNoPreds(s.preds)
	creditDB := new(model.Credit)
	if err := scheme.Convert(credit, creditDB); err != nil {
		return err
	}
	if _, err := s.query().Insert(creditDB); err != nil {
		return err
	}
	return nil
}

func (s *CreditStore) Get() (*credit.Credit, error) {
	result, err := s.GetDB()
	if err != nil {
		return nil, err
	}
	creditCore := &credit.Credit{}
	if err := scheme.Convert(result, creditCore); err != nil {
		return nil, err
	}
	return creditCore, nil
}

func (s *CreditStore) GetDB() (*model.Credit, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	var creditModel model.Credit
	err := query.ShouldGet(&creditModel)
	return &creditModel, err
}

func (s *CreditStore) ListCreditDBs() ([]*model.Credit, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	if len(s.Paging.Sort) == 0 {
		s.Paging.Sort = []string{"-updated_at"}
	}
	query, err := sqlstore.LimitSort(query, &s.Paging, Sort)
	if err != nil {
		return nil, err
	}

	var addrs model.Credits
	err = query.Find(&addrs)
	return addrs, err
}

func (s *CreditStore) ListCredit() ([]*credit.Credit, error) {
	result, err := s.ListCreditDBs()
	if err != nil {
		return nil, err
	}
	var credits []*credit.Credit
	if err := scheme.Convert(result, &credits); err != nil {
		return nil, err
	}
	return credits, nil
}

func (s *CreditStore) SoftDelete() (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	_deleted, err := query.Table("credit").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
		"status":     status3.N,
	})
	return _deleted, err
}

func (s *CreditStore) SumCredit() (int, error) {
	var totalAmount core.Int
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	if err := query.SQL("SELECT SUM(amount) from credit").Where("status = 1 AND paid_at is not NULL").Scan(&totalAmount); err != nil {
		return 0, err
	}
	return int(totalAmount), nil
}
