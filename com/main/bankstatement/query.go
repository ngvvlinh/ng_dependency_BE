package bankstatement

import (
	"context"

	"o.o/api/main/bankstatement"
	"o.o/api/main/identity"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/bankstatement/sqlstore"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
)

var _ bankstatement.QueryService = &BankStatementQueryService{}

type BankStatementQueryService struct {
	dbMain *cmsql.Database
	store  sqlstore.BankStatementFactory
}

func NewQueryBankStatement(
	dbMain com.MainDB,
	identityQ identity.QueryBus,
) *BankStatementQueryService {
	return &BankStatementQueryService{
		dbMain: dbMain,
		store:  sqlstore.NewBankStatementStore(dbMain),
	}
}

func BankStatementQueryServiceMessageBus(q *BankStatementQueryService) bankstatement.QueryBus {
	b := bus.New()
	return bankstatement.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *BankStatementQueryService) GetBankStatement(ctx context.Context, args *bankstatement.GetBankStatementArgs) (*bankstatement.BankStatement, error) {
	return q.store(ctx).ID(args.ID).Get()
}
