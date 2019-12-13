package sqlstore

import (
	"context"
	"fmt"
	"time"

	"etop.vn/api/main/ledgering"
	"etop.vn/api/meta"
	"etop.vn/api/top/types/etc/ledger_type"
	identityconvert "etop.vn/backend/com/main/identity/convert"
	"etop.vn/backend/com/main/ledgering/convert"
	"etop.vn/backend/com/main/ledgering/model"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/conversion"
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/common/sqlstore"
	"etop.vn/capi/dot"
)

type LedgerStoreFactory func(ctx context.Context) *LedgerStore

var scheme = conversion.Build(convert.RegisterConversions)

func NewLedgerStore(db *cmsql.Database) LedgerStoreFactory {
	return func(ctx context.Context) *LedgerStore {
		return &LedgerStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
		}
	}
}

type LedgerStore struct {
	ft ShopLedgerFilters

	query   func() cmsql.QueryInterface
	preds   []interface{}
	filters meta.Filters
	paging  meta.Paging

	includeDeleted sqlstore.IncludeDeleted
}

func (s *LedgerStore) Paging(paging meta.Paging) *LedgerStore {
	s.paging = paging
	return s
}

func (s *LedgerStore) GetPaing() meta.PageInfo {
	return meta.FromPaging(s.paging)
}

func (s *LedgerStore) Filters(filters meta.Filters) *LedgerStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *LedgerStore) ID(id dot.ID) *LedgerStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *LedgerStore) IDs(ids ...dot.ID) *LedgerStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "id", ids))
	return s
}

func (s *LedgerStore) ShopID(id dot.ID) *LedgerStore {
	s.preds = append(s.preds, s.ft.ByShopID(id))
	return s
}

func (s *LedgerStore) Type(ledgerType ledger_type.LedgerType) *LedgerStore {
	s.preds = append(s.preds, s.ft.ByType(ledgerType))
	return s
}

func (s *LedgerStore) AccountNumber(accountNumber string) *LedgerStore {
	condition := fmt.Sprintf("bank_account @> '{\"account_number\": \"%v\"}'", accountNumber)
	s.preds = append(s.preds, sq.NewExpr(condition))
	return s
}

func (s *LedgerStore) Count() (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	return query.Count((*model.ShopLedger)(nil))
}

func (s *LedgerStore) SoftDelete() (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	_deleted, err := query.Table("shop_ledger").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
	})
	return _deleted, err
}

func (s *LedgerStore) CreateLedger(ledger *ledgering.ShopLedger) error {
	sqlstore.MustNoPreds(s.preds)
	ledgerDB := new(model.ShopLedger)
	if err := scheme.Convert(ledger, ledgerDB); err != nil {
		return err
	}
	ledgerDB.BankAccount = identityconvert.BankAccountDB(ledger.BankAccount)
	if _, err := s.query().Insert(ledgerDB); err != nil {
		return err
	}

	var tempLedger model.ShopLedger
	if err := s.query().Where(s.ft.ByID(ledger.ID), s.ft.ByShopID(ledger.ShopID)).ShouldGet(&tempLedger); err != nil {
		return err
	}

	ledger.CreatedAt = tempLedger.CreatedAt
	ledger.UpdatedAt = tempLedger.UpdatedAt
	return nil
}

func (s *LedgerStore) GetLedgerDB() (*model.ShopLedger, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())

	var shopLedger model.ShopLedger
	err := query.ShouldGet(&shopLedger)
	return &shopLedger, err
}

func (s *LedgerStore) GetLedger() (ledgerResult *ledgering.ShopLedger, _ error) {
	ledger, err := s.GetLedgerDB()
	if err != nil {
		return nil, err
	}
	ledgerResult = convert.Convert_ledgeringmodel_ShopLedger_ledgering_ShopLedger(ledger, ledgerResult)
	ledgerResult.BankAccount = identityconvert.BankAccount(ledger.BankAccount)
	return ledgerResult, nil
}

func (s *LedgerStore) UpdateLedgerDB(ledger *model.ShopLedger) error {
	sqlstore.MustNoPreds(s.preds)
	err := s.query().Where(
		s.ft.ByID(ledger.ID),
		s.ft.ByShopID(ledger.ShopID),
	).UpdateAll().ShouldUpdate(ledger)
	return err
}

func (s *LedgerStore) ListLedgersDB() ([]*model.ShopLedger, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	query, err := sqlstore.LimitSort(query, &s.paging, SortLedger)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterLedger)
	if err != nil {
		return nil, err
	}

	var ledgers model.ShopLedgers
	err = query.Find(&ledgers)
	return ledgers, err
}

func (s *LedgerStore) ListLedgers() (ledgersResult []*ledgering.ShopLedger, _ error) {
	ledgers, err := s.ListLedgersDB()
	if err != nil {
		return nil, err
	}

	for _, ledger := range ledgers {
		var ledgerResult *ledgering.ShopLedger
		ledgerResult = convert.Convert_ledgeringmodel_ShopLedger_ledgering_ShopLedger(ledger, ledgerResult)
		ledgerResult.BankAccount = identityconvert.BankAccount(ledger.BankAccount)
		ledgersResult = append(ledgersResult, ledgerResult)
	}
	return ledgersResult, nil
}
