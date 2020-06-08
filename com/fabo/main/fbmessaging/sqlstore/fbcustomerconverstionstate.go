package sqlstore

import (
	"context"

	"o.o/api/fabo/fbmessaging"
	"o.o/api/meta"
	"o.o/backend/com/fabo/main/fbmessaging/convert"
	"o.o/backend/com/fabo/main/fbmessaging/model"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type FbCustomerConversationStateStoreFactory func(ctx context.Context) *FbCustomerConversationStateStore

func NewFbCustomerConversationStateStore(db *cmsql.Database) FbCustomerConversationStateStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *FbCustomerConversationStateStore {
		return &FbCustomerConversationStateStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type FbCustomerConversationStateStore struct {
	ft FbCustomerConversationStateFilters

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
	sqlstore.Paging
}

func (s *FbCustomerConversationStateStore) ID(id dot.ID) *FbCustomerConversationStateStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *FbCustomerConversationStateStore) IDs(ids ...dot.ID) *FbCustomerConversationStateStore {
	s.preds = append(s.preds, sq.In("id", ids))
	return s
}

func (s *FbCustomerConversationStateStore) IsRead(isRead bool) *FbCustomerConversationStateStore {
	s.preds = append(s.preds, s.ft.ByIsReadPtr(&isRead))
	return s
}

func (s *FbCustomerConversationStateStore) GetFbCustomerConversationStateDB() (*model.FbCustomerConversationState, error) {
	query := s.query().Where(s.preds)

	var fbCustomerConversationState model.FbCustomerConversationState
	err := query.ShouldGet(&fbCustomerConversationState)
	return &fbCustomerConversationState, err
}

func (s *FbCustomerConversationStateStore) GetFbCustomerConversationState() (*fbmessaging.FbCustomerConversationState, error) {
	fbCustomerConversationState, err := s.GetFbCustomerConversationStateDB()
	if err != nil {
		return nil, err
	}
	result := &fbmessaging.FbCustomerConversationState{}
	err = scheme.Convert(fbCustomerConversationState, result)
	if err != nil {
		return nil, err
	}
	return result, err
}

func (s *FbCustomerConversationStateStore) UpdateIsRead(isRead bool) (int, error) {
	query := s.query().Where(s.preds)
	updateStatus, err := query.Table("fb_customer_conversation_state").UpdateMap(map[string]interface{}{
		"is_read": isRead,
	})
	return updateStatus, err
}

func (s *FbCustomerConversationStateStore) ListFbCustomerConversationStatesDB() ([]*model.FbCustomerConversationState, error) {
	query := s.query().Where(s.preds)

	var fbCustomerConversationStates model.FbCustomerConversationStates
	err := query.Find(&fbCustomerConversationStates)
	if err != nil {
		return nil, err
	}
	s.Paging.Apply(fbCustomerConversationStates)
	return fbCustomerConversationStates, nil
}

func (s *FbCustomerConversationStateStore) ListFbCustomerConversationStates() (result []*fbmessaging.FbCustomerConversationState, err error) {
	fbCustomerConversationStates, err := s.ListFbCustomerConversationStatesDB()
	if err != nil {
		return nil, err
	}
	if err = scheme.Convert(fbCustomerConversationStates, &result); err != nil {
		return nil, err
	}
	return
}

func (s *FbCustomerConversationStateStore) CreateFbCustomerConversationState(fbCustomerConversationState *fbmessaging.FbCustomerConversationState) error {
	sqlstore.MustNoPreds(s.preds)
	fbCustomerConversationStateDB := new(model.FbCustomerConversationState)
	if err := scheme.Convert(fbCustomerConversationState, fbCustomerConversationStateDB); err != nil {
		return err
	}

	_, err := s.query().Upsert(fbCustomerConversationStateDB)
	if err != nil {
		return err
	}

	var tempFbCustomerConversationState model.FbCustomerConversationState
	if err := s.query().Where(s.ft.ByID(fbCustomerConversationState.ID)).ShouldGet(&tempFbCustomerConversationState); err != nil {
		return err
	}
	fbCustomerConversationState.UpdatedAt = tempFbCustomerConversationState.UpdatedAt

	return nil
}

func (s *FbCustomerConversationStateStore) CreateFbCustomerConversationStates(fbCustomerConversationStates []*fbmessaging.FbCustomerConversationState) error {
	sqlstore.MustNoPreds(s.preds)
	fbCustomerConversationStatesDB := model.FbCustomerConversationStates(convert.Convert_fbmessaging_FbCustomerConversationStates_fbmessagingmodel_FbCustomerConversationStates(fbCustomerConversationStates))

	_, err := s.query().Upsert(&fbCustomerConversationStatesDB)
	if err != nil {
		return err
	}
	return nil
}
