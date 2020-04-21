package sqlstore

import (
	"context"

	"o.o/api/subscripting/subscriptionbill"
	"o.o/backend/com/subscripting/subscriptionbill/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type SubrBillLineStore struct {
	ft    SubscriptionBillLineFilters
	query func() cmsql.QueryInterface
	preds []interface{}

	includeDeleted sqlstore.IncludeDeleted
}

type SubrBillLineStoreFactory func(context.Context) *SubrBillLineStore

func NewSubrBillLineStore(db *cmsql.Database) SubrBillLineStoreFactory {
	return func(ctx context.Context) *SubrBillLineStore {
		return &SubrBillLineStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
		}
	}
}

func (s *SubrBillLineStore) SubscriptionBillID(id dot.ID) *SubrBillLineStore {
	s.preds = append(s.preds, s.ft.BySubscriptionBillID(id))
	return s
}

func (s *SubrBillLineStore) SubscriptionBillIDs(ids ...dot.ID) *SubrBillLineStore {
	s.preds = append(s.preds, sq.In("subscription_bill_id", ids))
	return s
}

func (s *SubrBillLineStore) CreateSubrBillLineDB(line *model.SubscriptionBillLine) error {
	if line.ID == 0 {
		line.ID = cm.NewID()
	}
	if line.SubscriptionBillID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing subscription bill ID")
	}
	return s.query().ShouldInsert(line)
}

func (s *SubrBillLineStore) CreateSubrBillLine(args *subscriptionbill.SubscriptionBillLine) error {
	var line model.SubscriptionBillLine
	if err := scheme.Convert(args, &line); err != nil {
		return err
	}
	return s.CreateSubrBillLineDB(&line)
}

func (s *SubrBillLineStore) ListSubrBillLinesDB() (res []*model.SubscriptionBillLine, err error) {
	err = s.query().Where(s.preds).Find((*model.SubscriptionBillLines)(&res))
	return
}
