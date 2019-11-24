package sqlstore

import (
	"context"
	"database/sql"

	etoptypes "etop.vn/api/main/etop"
	"etop.vn/api/main/transaction"
	"etop.vn/api/meta"
	"etop.vn/backend/com/main/transaction/convert"
	transactionmodel "etop.vn/backend/com/main/transaction/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/sqlstore"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
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
	query  cmsql.QueryFactory
	preds  []interface{}
	ft     TransactionFilters
	paging meta.Paging
}

var SortTransaction = map[string]string{
	"created_at": "",
	"updated_at": "",
}

func (s *TransactionStore) ID(id dot.ID) *TransactionStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *TransactionStore) AccountID(id dot.ID) *TransactionStore {
	s.preds = append(s.preds, s.ft.ByAccountID(id))
	return s
}

func (s *TransactionStore) ByConfirmedTransaction() *TransactionStore {
	s.preds = append(s.preds, s.ft.ByStatus(int(model.S3Positive)))
	return s
}

func (s *TransactionStore) OptionalTransactionType(transactionType string) *TransactionStore {
	s.preds = append(s.preds, s.ft.ByType(transactionType).Optional())
	return s
}

func (s *TransactionStore) Paging(paging meta.Paging) *TransactionStore {
	s.paging = paging
	return s
}

func (s *TransactionStore) GetPaging() meta.PageInfo {
	return meta.FromPaging(s.paging)
}

func (s *TransactionStore) GetTransactionDB() (*transactionmodel.Transaction, error) {
	var transaction transactionmodel.Transaction
	err := s.query().Where(s.preds).ShouldGet(&transaction)
	return &transaction, err
}

func (s *TransactionStore) GetTransaction() (*transaction.Transaction, error) {
	transaction, err := s.GetTransactionDB()
	if err != nil {
		return nil, err
	}
	return convert.Transaction(transaction), nil
}

func (s *TransactionStore) ListTransactionsDB() ([]*transactionmodel.Transaction, error) {
	query := s.query().Where(s.preds)
	query, err := sqlstore.LimitSort(query, &s.paging, SortTransaction)
	if err != nil {
		return nil, err
	}
	var transactions transactionmodel.Transactions
	err = query.Find(&transactions)
	return transactions, err
}

func (s *TransactionStore) ListTransactions() ([]*transaction.Transaction, error) {
	transactions, err := s.ListTransactionsDB()
	if err != nil {
		return nil, err
	}
	return convert.Transactions(transactions), err
}

func (s *TransactionStore) Count() (int, error) {
	query := s.query().Where(s.preds)
	count, err := query.Count((*transactionmodel.Transaction)(nil))
	return int(count), err
}

func (s *TransactionStore) CreateTransaction(trxn *transaction.Transaction) (*transaction.Transaction, error) {
	sqlstore.MustNoPreds(s.preds)
	trxnDB := convert.TransactionDB(trxn)
	if trxnDB.ID == 0 {
		trxnDB.ID = cm.NewID()
	}
	_, err := s.query().Insert(trxnDB)
	if err != nil {
		return nil, err
	}
	return s.ID(trxnDB.ID).GetTransaction()
}

type UpdateTransactionStatusArgs struct {
	ID        dot.ID
	AccountID dot.ID
	Status    etoptypes.Status3
}

func (s *TransactionStore) UpdateTransactionStatus(args *UpdateTransactionStatusArgs) (*transaction.Transaction, error) {
	sqlstore.MustNoPreds(s.preds)
	update := &transactionmodel.Transaction{
		Status: int(args.Status),
	}

	if err := s.query().Table("transaction").Where(s.ft.ByID(args.ID)).Where(s.ft.ByAccountID(args.AccountID)).ShouldUpdate(update); err != nil {
		return nil, err
	}
	return s.ID(args.ID).GetTransaction()
}

func (s *TransactionStore) GetBalance() (int, error) {
	var totalAmount sql.NullInt64
	if err := s.query().SQL("SELECT SUM(amount) from transaction").Where(s.preds).Scan(&totalAmount); err != nil {
		return 0, err
	}
	return int(totalAmount.Int64), nil
}
