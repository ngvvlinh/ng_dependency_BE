package sqlstore

import (
	"context"
	"time"

	"o.o/api/main/transaction"
	"o.o/api/meta"
	"o.o/api/top/types/etc/service_classify"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/subject_referral"
	"o.o/api/top/types/etc/transaction_type"
	"o.o/backend/com/main/transaction/convert"
	transactionmodel "o.o/backend/com/main/transaction/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sq/core"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type TransactionStoreFactory func(context.Context) *TransactionStore

func NewTransactionStore(db *cmsql.Database) TransactionStoreFactory {
	transactionmodel.SQLVerifySchema(db)
	return func(ctx context.Context) *TransactionStore {
		return &TransactionStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type TransactionStore struct {
	query cmsql.QueryFactory
	preds []interface{}
	ft    TransactionFilters
	sqlstore.Paging
}

var SortTransaction = map[string]string{
	"id":         "id",
	"created_at": "created_at",
	"updated_at": "updated_at",
}

func (s *TransactionStore) ID(id dot.ID) *TransactionStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *TransactionStore) AccountID(id dot.ID) *TransactionStore {
	s.preds = append(s.preds, s.ft.ByAccountID(id))
	return s
}

func (s *TransactionStore) AccountIDs(ids ...dot.ID) *TransactionStore {
	s.preds = append(s.preds, sq.In("account_id", ids))
	return s
}

func (s *TransactionStore) Classify(classify service_classify.ServiceClassify) *TransactionStore {
	s.preds = append(s.preds, s.ft.ByClassifyPtr(&classify))
	return s
}

func (s *TransactionStore) ByConfirmedTransaction() *TransactionStore {
	s.preds = append(s.preds, s.ft.ByStatus(status3.P))
	return s
}

func (s *TransactionStore) OptionalTransactionType(transactionType transaction_type.TransactionType) *TransactionStore {
	s.preds = append(s.preds, s.ft.ByType(transactionType).Optional())
	return s
}

func (s *TransactionStore) ReferralType(_type subject_referral.SubjectReferral) *TransactionStore {
	s.preds = append(s.preds, s.ft.ByReferralType(_type))
	return s
}

func (s *TransactionStore) ReferralID(id dot.ID) *TransactionStore {
	s.preds = append(s.preds, sq.NewExpr("referral_ids @> ?", core.Array{V: []dot.ID{id}}))
	return s
}

func (s *TransactionStore) BetweenDateFromAndDateTo(dateFrom time.Time, dateTo time.Time) *TransactionStore {
	s.preds = append(s.preds, sq.NewExpr("created_at BETWEEN ? AND ?", dateFrom, dateTo))
	return s
}

func (s *TransactionStore) WithPaging(paging meta.Paging) *TransactionStore {
	s.Paging.WithPaging(paging)
	return s
}

func (s *TransactionStore) GetTransactionDB() (*transactionmodel.Transaction, error) {
	var transaction transactionmodel.Transaction
	err := s.query().Where(s.preds).ShouldGet(&transaction)
	return &transaction, err
}

func (s *TransactionStore) GetTransaction() (*transaction.Transaction, error) {
	trx, err := s.GetTransactionDB()
	if err != nil {
		return nil, err
	}
	res := convert.Convert_transactionmodel_Transaction_transaction_Transaction(trx, nil)
	return res, nil
}

func (s *TransactionStore) ListTransactionsDB() (res []*transactionmodel.Transaction, err error) {
	query := s.query().Where(s.preds)
	if len(s.Paging.Sort) == 0 {
		s.Paging.Sort = []string{"-created_at"}
	}
	query, err = sqlstore.LimitSort(query, &s.Paging, SortTransaction)
	if err != nil {
		return nil, err
	}
	if err = query.Find((*transactionmodel.Transactions)(&res)); err != nil {
		return nil, err
	}
	s.Paging.Apply(res)
	return
}

func (s *TransactionStore) ListTransactions() ([]*transaction.Transaction, error) {
	transactions, err := s.ListTransactionsDB()
	if err != nil {
		return nil, err
	}
	return convert.Convert_transactionmodel_Transactions_transaction_Transactions(transactions), err
}

func (s *TransactionStore) Count() (int, error) {
	query := s.query().Where(s.preds)
	count, err := query.Count((*transactionmodel.Transaction)(nil))
	return count, err
}

func (s *TransactionStore) CreateTransactionDB(trxn *transactionmodel.Transaction) error {
	return s.query().ShouldInsert(trxn)
}

func (s *TransactionStore) CreateTransaction(trxn *transaction.Transaction) (*transaction.Transaction, error) {
	sqlstore.MustNoPreds(s.preds)
	trxnDB := convert.Convert_transaction_Transaction_transactionmodel_Transaction(trxn, nil)
	if trxnDB.ID == 0 {
		trxnDB.ID = cm.NewID()
	}
	if err := s.CreateTransactionDB(trxnDB); err != nil {
		return nil, err
	}
	return s.ID(trxnDB.ID).GetTransaction()
}

type UpdateTransactionStatusArgs struct {
	ID        dot.ID
	AccountID dot.ID
	Status    status3.Status
}

func (s *TransactionStore) UpdateTransactionStatus(args *UpdateTransactionStatusArgs) (*transaction.Transaction, error) {
	sqlstore.MustNoPreds(s.preds)
	update := &transactionmodel.Transaction{
		Status: args.Status,
	}

	if err := s.query().Table("transaction").Where(s.ft.ByID(args.ID)).Where(s.ft.ByAccountID(args.AccountID)).ShouldUpdate(update); err != nil {
		return nil, err
	}
	return s.ID(args.ID).GetTransaction()
}

func (s *TransactionStore) GetBalance() (int, error) {
	var totalAmount core.Int
	if err := s.query().SQL("SELECT SUM(amount) from transaction").Where(s.preds).Where("status = ?", status3.P).Scan(&totalAmount); err != nil {
		return 0, err
	}
	return int(totalAmount), nil
}

func (s *TransactionStore) DeleteTransaction() error {
	return s.query().Where(s.preds).ShouldDelete(&transactionmodel.Transaction{})
}
