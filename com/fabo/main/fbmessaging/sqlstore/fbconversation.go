package sqlstore

import (
	"context"
	"time"

	"o.o/api/fabo/fbmessaging"
	"o.o/api/meta"
	"o.o/backend/com/fabo/main/fbmessaging/convert"
	"o.o/backend/com/fabo/main/fbmessaging/model"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type FbExternalConversationStoreFactory func(ctx context.Context) *FbExternalConversationStore

func NewFbExternalConversationStore(db *cmsql.Database) FbExternalConversationStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *FbExternalConversationStore {
		return &FbExternalConversationStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type FbExternalConversationStore struct {
	ft FbExternalConversationFilters

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
	sqlstore.Paging

	includeDeleted sqlstore.IncludeDeleted
}

func (s *FbExternalConversationStore) WithPaging(paging meta.Paging) *FbExternalConversationStore {
	s.Paging.WithPaging(paging)
	return s
}

func (s *FbExternalConversationStore) ExternalID(externalID string) *FbExternalConversationStore {
	s.preds = append(s.preds, s.ft.ByExternalID(externalID))
	return s
}

func (s *FbExternalConversationStore) ExternalIDs(externalIDs []string) *FbExternalConversationStore {
	s.preds = append(s.preds, sq.In("external_id", externalIDs))
	return s
}

func (s *FbExternalConversationStore) ExternalPageID(externalPageID string) *FbExternalConversationStore {
	s.preds = append(s.preds, s.ft.ByExternalPageID(externalPageID))
	return s
}

func (s *FbExternalConversationStore) ExternalUserID(externalUserID string) *FbExternalConversationStore {
	s.preds = append(s.preds, s.ft.ByExternalUserID(externalUserID))
	return s
}

func (s *FbExternalConversationStore) ID(ID dot.ID) *FbExternalConversationStore {
	s.preds = append(s.preds, s.ft.ByID(ID))
	return s
}

func (s *FbExternalConversationStore) CreateFbExternalConversation(fbExternalConversation *fbmessaging.FbExternalConversation) error {
	sqlstore.MustNoPreds(s.preds)
	fbExternalConversationDB := new(model.FbExternalConversation)
	if err := scheme.Convert(fbExternalConversation, fbExternalConversationDB); err != nil {
		return err
	}

	_, err := s.query().Upsert(fbExternalConversationDB)
	if err != nil {
		return err
	}

	var tempFbExternalConversation model.FbExternalConversation
	if err := s.query().Where(s.ft.ByID(fbExternalConversation.ID)).ShouldGet(&tempFbExternalConversation); err != nil {
		return err
	}
	fbExternalConversation.CreatedAt = tempFbExternalConversation.CreatedAt
	fbExternalConversation.UpdatedAt = tempFbExternalConversation.UpdatedAt

	return nil
}

func (s *FbExternalConversationStore) CreateFbExternalConversations(fbExternalConversations []*fbmessaging.FbExternalConversation) error {
	sqlstore.MustNoPreds(s.preds)
	fbExternalConversationsDB := model.FbExternalConversations(convert.Convert_fbmessaging_FbExternalConversations_fbmessagingmodel_FbExternalConversations(fbExternalConversations))

	_, err := s.query().Upsert(&fbExternalConversationsDB)
	if err != nil {
		return err
	}
	return nil
}

func (s *FbExternalConversationStore) GetFbExternalConversationDB() (*model.FbExternalConversation, error) {
	query := s.query().Where(s.preds)

	var fbExternalConversation model.FbExternalConversation
	err := query.ShouldGet(&fbExternalConversation)
	return &fbExternalConversation, err
}

func (s *FbExternalConversationStore) GetFbExternalConversation() (*fbmessaging.FbExternalConversation, error) {
	fbExternalConversation, err := s.GetFbExternalConversationDB()
	if err != nil {
		return nil, err
	}
	result := &fbmessaging.FbExternalConversation{}
	err = scheme.Convert(fbExternalConversation, result)
	if err != nil {
		return nil, err
	}
	return result, err
}

func (s *FbExternalConversationStore) ListFbExternalConversationsDB() ([]*model.FbExternalConversation, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	if !s.Paging.IsCursorPaging() && len(s.Paging.Sort) == 0 {
		s.Paging.Sort = []string{"-created_at"}
	}
	query, err := sqlstore.LimitSort(query, &s.Paging, SortFbExternalConversation, s.ft.prefix)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterFbExternalConversation)
	if err != nil {
		return nil, err
	}

	var fbExternalConversations model.FbExternalConversations
	err = query.Find(&fbExternalConversations)
	if err != nil {
		return nil, err
	}
	s.Paging.Apply(fbExternalConversations)
	return fbExternalConversations, nil
}

func (s *FbExternalConversationStore) UpdateFbExternalConversation(fbExternalConversation *fbmessaging.FbExternalConversation) error {
	fbExternalConversationDB := new(model.FbExternalConversation)
	if err := scheme.Convert(fbExternalConversation, fbExternalConversationDB); err != nil {
		return err
	}

	return s.query().Where(s.preds).ShouldUpdate(fbExternalConversationDB)
}

func (s *FbExternalConversationStore) ListFbExternalConversations() (result []*fbmessaging.FbExternalConversation, err error) {
	fbExternalConversations, err := s.ListFbExternalConversationsDB()
	if err != nil {
		return nil, err
	}
	if err = scheme.Convert(fbExternalConversations, &result); err != nil {
		return nil, err
	}
	return
}

func (s *FbExternalConversationStore) SoftDelete() (int, error) {
	query := s.query().Where(s.preds)
	return query.Table("fb_external_conversation").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
	})
}
