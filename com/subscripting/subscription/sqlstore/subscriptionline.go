package sqlstore

import (
	"context"

	"o.o/api/subscripting/subscription"
	"o.o/backend/com/subscripting/subscription/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type SubscriptionLineStore struct {
	ft    SubscriptionLineFilters
	query func() cmsql.QueryInterface
	preds []interface{}

	includeDeleted sqlstore.IncludeDeleted
}

type SubscriptionLineStoreFactory func(context.Context) *SubscriptionLineStore

func NewSubscriptionLineStore(db *cmsql.Database) SubscriptionLineStoreFactory {
	return func(ctx context.Context) *SubscriptionLineStore {
		return &SubscriptionLineStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
		}
	}
}

func (s *SubscriptionLineStore) ID(id dot.ID) *SubscriptionLineStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *SubscriptionLineStore) IDs(ids []dot.ID) *SubscriptionLineStore {
	s.preds = append(s.preds, sq.In("id", ids))
	return s
}

func (s *SubscriptionLineStore) SubscriptionID(id dot.ID) *SubscriptionLineStore {
	s.preds = append(s.preds, s.ft.BySubscriptionID(id))
	return s
}

func (s *SubscriptionLineStore) SubscriptionIDs(ids ...dot.ID) *SubscriptionLineStore {
	s.preds = append(s.preds, sq.In("subscription_id", ids))
	return s
}

func (s *SubscriptionLineStore) CreateSubscriptionLineDB(line *model.SubscriptionLine) error {
	if line.ID == 0 {
		line.ID = cm.NewID()
	}
	if line.SubscriptionID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing subscription ID").WithMetap("func", "CreateSubscriptionLineDB")
	}
	return s.query().ShouldInsert(line)
}

func (s *SubscriptionLineStore) CreateSubscriptionLine(args *subscription.SubscriptionLine) error {
	var line model.SubscriptionLine
	if err := scheme.Convert(args, &line); err != nil {
		return err
	}
	return s.CreateSubscriptionLineDB(&line)
}

func (s *SubscriptionLineStore) UpdateSubscriptionLineDB(line *model.SubscriptionLine) error {
	if len(s.preds) == 0 {
		return cm.Errorf(cm.FailedPrecondition, nil, "must provide preds")
	}
	return s.query().Where(s.preds).ShouldUpdate(line)
}

func (s *SubscriptionLineStore) UpdateSubscriptionLine(line *subscription.SubscriptionLine) error {
	var lineDB model.SubscriptionLine
	if err := scheme.Convert(line, &lineDB); err != nil {
		return err
	}
	return s.UpdateSubscriptionLineDB(&lineDB)
}

func (s *SubscriptionLineStore) ListSubscriptionLinesDB() (res []*model.SubscriptionLine, err error) {
	err = s.query().Where(s.preds).Find((*model.SubscriptionLines)(&res))
	return
}

func (s *SubscriptionLineStore) DeleteSubscriptionLine() error {
	if len(s.preds) == 0 {
		return cm.Errorf(cm.FailedPrecondition, nil, "must provide preds")
	}
	return s.query().Where(s.preds).ShouldDelete(&model.SubscriptionLine{})
}
