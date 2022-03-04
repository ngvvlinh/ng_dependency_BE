package sqlstore

import (
	"context"

	"o.o/api/main/bankstatement"
	"o.o/backend/com/main/bankstatement/convert"
	"o.o/backend/com/main/bankstatement/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

var scheme = conversion.Build(convert.RegisterConversions)

type BankStatementFactory func(context.Context) *BankStatementStore

func NewBankStatementStore(db *cmsql.Database) BankStatementFactory {
	return func(ctx context.Context) *BankStatementStore {
		return &BankStatementStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type BankStatementStore struct {
	query cmsql.QueryFactory
	ft    BankStatementFilters
	sqlstore.Paging
	preds []interface{}

	includeDeleted sqlstore.IncludeDeleted
}

func (s *BankStatementStore) ID(id dot.ID) *BankStatementStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *BankStatementStore) ExternalTrxnID(trxnID string) *BankStatementStore {
	s.preds = append(s.preds, s.ft.ByExternalTransactionID(trxnID))
	return s
}

func (s *BankStatementStore) WithPaging(paging *cm.Paging) *BankStatementStore {
	s.Paging.WithPaging(*paging)
	return s
}

func (s *BankStatementStore) Create(args *bankstatement.BankStatement) error {
	sqlstore.MustNoPreds(s.preds)
	bankStatement := convert.Convert_bankstatement_BankStatement_bankstatementmodel_BankStatement(args, nil)

	if _, err := s.query().Insert(bankStatement); err != nil {
		return err
	}
	return nil
}

func (s *BankStatementStore) Get() (*bankstatement.BankStatement, error) {
	result, err := s.GetDB()
	if err != nil {
		return nil, err
	}
	bankStatement := &bankstatement.BankStatement{}
	if err := scheme.Convert(result, bankStatement); err != nil {
		return nil, err
	}
	return bankStatement, nil
}

func (s *BankStatementStore) GetDB() (*model.BankStatement, error) {
	query := s.query().Where(s.preds)
	var res model.BankStatement
	err := query.ShouldGet(&res)
	return &res, err
}
