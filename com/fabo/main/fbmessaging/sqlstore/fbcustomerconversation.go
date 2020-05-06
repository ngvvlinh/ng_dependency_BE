package sqlstore

import (
	"context"

	"o.o/api/fabo/fbmessaging"
	"o.o/api/fabo/fbmessaging/fb_customer_conversation_type"
	"o.o/api/meta"
	"o.o/backend/com/fabo/main/fbmessaging/convert"
	"o.o/backend/com/fabo/main/fbmessaging/model"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type FbCustomerConversationStoreFactory func(ctx context.Context) *FbCustomerConversationStore

func NewFbCustomerConversationStore(db *cmsql.Database) FbCustomerConversationStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *FbCustomerConversationStore {
		return &FbCustomerConversationStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type FbCustomerConversationStore struct {
	ft FbCustomerConversationFilters

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
	sqlstore.Paging

	includeDeleted sqlstore.IncludeDeleted
}

func (s *FbCustomerConversationStore) WithPaging(paging meta.Paging) *FbCustomerConversationStore {
	s.Paging.WithPaging(paging)
	return s
}

func (s *FbCustomerConversationStore) ExternalIDs(externalIDs []string) *FbCustomerConversationStore {
	s.preds = append(s.preds, sq.In("external_id", externalIDs))
	return s
}

func (s *FbCustomerConversationStore) FbPageID(fbPageID dot.ID) *FbCustomerConversationStore {
	s.preds = append(s.preds, s.ft.ByFbPageID(fbPageID))
	return s
}

func (s *FbCustomerConversationStore) FbPageIDs(fbPageIDs []dot.ID) *FbCustomerConversationStore {
	s.preds = append(s.preds, sq.In("fb_page_id", fbPageIDs))
	return s
}

func (s *FbCustomerConversationStore) FbExternalID(fbExternalID string) *FbCustomerConversationStore {
	s.preds = append(s.preds, s.ft.ByExternalID(fbExternalID))
	return s
}

func (s *FbCustomerConversationStore) Type(typ fb_customer_conversation_type.FbCustomerConversationType) *FbCustomerConversationStore {
	s.preds = append(s.preds, s.ft.ByType(typ.Enum()))
	return s
}

func (s *FbCustomerConversationStore) IsRead(isRead bool) *FbCustomerConversationStore {
	s.preds = append(s.preds, s.ft.ByIsRead(isRead))
	return s
}

func (s *FbCustomerConversationStore) CreateFbCustomerConversation(fbCustomerConversation *fbmessaging.FbCustomerConversation) error {
	sqlstore.MustNoPreds(s.preds)
	fbCustomerConversationDB := new(model.FbCustomerConversation)
	if err := scheme.Convert(fbCustomerConversation, fbCustomerConversationDB); err != nil {
		return err
	}

	_, err := s.query().Upsert(fbCustomerConversationDB)
	if err != nil {
		return err
	}

	var tempFbCustomerConversation model.FbCustomerConversation
	if err := s.query().Where(s.ft.ByID(fbCustomerConversation.ID)).ShouldGet(&tempFbCustomerConversation); err != nil {
		return err
	}
	fbCustomerConversation.CreatedAt = tempFbCustomerConversation.CreatedAt
	fbCustomerConversation.UpdatedAt = tempFbCustomerConversation.UpdatedAt

	return nil
}

func (s *FbCustomerConversationStore) CreateFbCustomerConversations(fbCustomerConversations []*fbmessaging.FbCustomerConversation) error {
	sqlstore.MustNoPreds(s.preds)
	fbCustomerConversationsDB := model.FbCustomerConversations(convert.Convert_fbmessaging_FbCustomerConversations_fbmessagingmodel_FbCustomerConversations(fbCustomerConversations))

	_, err := s.query().Upsert(&fbCustomerConversationsDB)
	if err != nil {
		return err
	}
	return nil
}

func (s *FbCustomerConversationStore) ListFbCustomerConversationsDB() ([]*model.FbCustomerConversation, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	if !s.Paging.IsCursorPaging() && len(s.Paging.Sort) == 0 {
		s.Paging.Sort = []string{"-created_at"}
	}
	query, err := sqlstore.LimitSort(query, &s.Paging, SortFbCustomerConversation, s.ft.prefix)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterFbCustomerConversation)
	if err != nil {
		return nil, err
	}

	var fbCustomerConversations model.FbCustomerConversations
	err = query.Find(&fbCustomerConversations)
	if err != nil {
		return nil, err
	}
	s.Paging.Apply(fbCustomerConversations)
	return fbCustomerConversations, nil
}

func (s *FbCustomerConversationStore) ListFbCustomerConversations() (result []*fbmessaging.FbCustomerConversation, err error) {
	fbCustomerConversations, err := s.ListFbCustomerConversationsDB()
	if err != nil {
		return nil, err
	}
	if err = scheme.Convert(fbCustomerConversations, &result); err != nil {
		return nil, err
	}
	return
}

func (s *FbCustomerConversationStore) IncludeDeleted() *FbCustomerConversationStore {
	s.includeDeleted = true
	return s
}
