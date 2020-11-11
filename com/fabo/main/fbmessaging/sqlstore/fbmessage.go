package sqlstore

import (
	"context"
	"fmt"
	"strings"

	"o.o/api/fabo/fbmessaging"
	"o.o/api/meta"
	fbsearchmodel "o.o/backend/com/fabo/main/fbcustomerconversationsearch/model"
	"o.o/backend/com/fabo/main/fbmessaging/convert"
	"o.o/backend/com/fabo/main/fbmessaging/model"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type FbExternalMessageStoreFactory func(ctx context.Context) *FbExternalMessageStore

func NewFbExternalMessageStore(db *cmsql.Database) FbExternalMessageStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *FbExternalMessageStore {
		return &FbExternalMessageStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type FbExternalMessageStore struct {
	ft FbExternalMessageFilters

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
	sqlstore.Paging

	includeDeleted sqlstore.IncludeDeleted
}

func (s *FbExternalMessageStore) WithPaging(paging meta.Paging) *FbExternalMessageStore {
	s.Paging.WithPaging(paging)
	return s
}

func (s *FbExternalMessageStore) IDs(IDs []dot.ID) *FbExternalMessageStore {
	s.preds = append(s.preds, sq.In("id", IDs))
	return s
}

func (s *FbExternalMessageStore) ExternalConversationIDs(externalConversationIDs []string) *FbExternalMessageStore {
	s.preds = append(s.preds, sq.In("external_conversation_id", externalConversationIDs))
	return s
}

func (s *FbExternalMessageStore) ExternalIDs(externalIDs []string) *FbExternalMessageStore {
	s.preds = append(s.preds, sq.In("external_id", externalIDs))
	return s
}

func (s *FbExternalMessageStore) ExternalPageIDs(externalPageIDs []string) *FbExternalMessageStore {
	s.preds = append(s.preds, sq.In("external_page_id", externalPageIDs))
	return s
}

func (s *FbExternalMessageStore) ExternalID(externalID string) *FbExternalMessageStore {
	s.preds = append(s.preds, s.ft.ByExternalID(externalID))
	return s
}

func (s *FbExternalMessageStore) ID(ID dot.ID) *FbExternalMessageStore {
	s.preds = append(s.preds, s.ft.ByID(ID))
	return s
}

func (s *FbExternalMessageStore) CreateFbExternalMessage(fbExternalMessage *fbmessaging.FbExternalMessage) error {
	sqlstore.MustNoPreds(s.preds)
	fbExternalMessageDB := new(model.FbExternalMessage)
	if err := scheme.Convert(fbExternalMessage, fbExternalMessageDB); err != nil {
		return err
	}

	_, err := s.query().Insert(fbExternalMessageDB)
	if err != nil {
		return err
	}

	var tempFbExternalMessage model.FbExternalMessage
	if err := s.query().Where(s.ft.ByID(fbExternalMessage.ID)).ShouldGet(&tempFbExternalMessage); err != nil {
		return err
	}
	fbExternalMessage.CreatedAt = tempFbExternalMessage.CreatedAt
	fbExternalMessage.UpdatedAt = tempFbExternalMessage.UpdatedAt

	// prepare for search
	messageSearch := &fbsearchmodel.FbExternalMessageSearch{
		ID:                     fbExternalMessageDB.ID,
		ExternalMessageNorm:    normalizeText(fbExternalMessageDB.ExternalMessage),
		ExternalPageID:         fbExternalMessageDB.ExternalPageID,
		CreatedAt:              fbExternalMessageDB.CreatedAt,
		ExternalConversationID: fbExternalMessageDB.ExternalConversationID,
	}
	if _, err = s.query().Upsert(messageSearch); err != nil {
		ll.Error(fmt.Sprintf("create fb_external_message_search got error: %v", err))
	}

	return nil
}

func (s *FbExternalMessageStore) CreateFbExternalMessages(fbExternalMessages []*fbmessaging.FbExternalMessage) error {
	sqlstore.MustNoPreds(s.preds)
	fbExternalMessagesDB := model.FbExternalMessages(convert.Convert_fbmessaging_FbExternalMessages_fbmessagingmodel_FbExternalMessages(fbExternalMessages))

	_, err := s.query().Upsert(&fbExternalMessagesDB)
	if err != nil {
		return err
	}

	// prepare for search
	var messageSearchs fbsearchmodel.FbExternalMessageSearchs
	for _, msg := range fbExternalMessagesDB {
		messageSearchs = append(messageSearchs, &fbsearchmodel.FbExternalMessageSearch{
			ID:                     msg.ID,
			ExternalMessageNorm:    normalizeText(msg.ExternalMessage),
			CreatedAt:              msg.CreatedAt,
			ExternalConversationID: msg.ExternalConversationID,
			ExternalPageID:         msg.ExternalPageID,
		})
	}
	if _, err = s.query().Upsert(&messageSearchs); err != nil {
		ll.Error(fmt.Sprintf("create fb_external_message_search got error: %v", err))
	}

	return nil
}

func (s *FbExternalMessageStore) ListFbExternalMessagesDB() ([]*model.FbExternalMessage, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	if !s.Paging.IsCursorPaging() && len(s.Paging.Sort) == 0 {
		s.Paging.Sort = []string{"-created_at"}
	}
	query, err := sqlstore.LimitSort(query, &s.Paging, SortFbExternalMessage, s.ft.prefix)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterFbExternalMessage)
	if err != nil {
		return nil, err
	}

	var fbExternalMessages model.FbExternalMessages
	err = query.OrderBy("external_timestamp desc").Find(&fbExternalMessages)
	if err != nil {
		return nil, err
	}
	s.Paging.Apply(fbExternalMessages)
	return fbExternalMessages, nil
}

func (s *FbExternalMessageStore) ListFbExternalMessages() (result []*fbmessaging.FbExternalMessage, err error) {
	fbExternalConversations, err := s.ListFbExternalMessagesDB()
	if err != nil {
		return nil, err
	}
	if err = scheme.Convert(fbExternalConversations, &result); err != nil {
		return nil, err
	}
	return
}

func (s *FbExternalMessageStore) ListLatestExternalMessages(externalConversationIDs []string) (result []*fbmessaging.FbExternalMessage, err error) {
	if len(externalConversationIDs) == 0 {
		return nil, nil
	}

	rows, err := s.query().
		SQL(`
			SELECT DISTINCT ON (external_conversation_id) id
			FROM fb_external_message
			WHERE external_conversation_id IN (?)
			ORDER BY external_conversation_id, external_created_time DESC, external_timestamp DESC
		`, strings.Join(externalConversationIDs, "','")).
		Query()
	if err != nil {
		return nil, err
	}

	var fbExternalMessageIDs []dot.ID
	var id dot.ID
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		fbExternalMessageIDs = append(fbExternalMessageIDs, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	var fbExternalMessages model.FbExternalMessages
	if err := s.query().
		Where(sq.In("id", fbExternalMessageIDs)).
		Find(&fbExternalMessages); err != nil {
		return nil, err
	}

	if err := scheme.Convert([]*model.FbExternalMessage(fbExternalMessages), &result); err != nil {
		return nil, err
	}
	return
}

func (s *FbExternalMessageStore) ListLatestCustomerExternalMessages(externalConversationIDs []string) (result []*fbmessaging.FbExternalMessage, err error) {
	if len(externalConversationIDs) == 0 {
		return nil, nil
	}

	rows, err := s.query().
		SQL(`
			SELECT DISTINCT ON (external_conversation_id) id
			FROM fb_external_message
			WHERE
				external_conversation_id IN (?) AND
				external_from_id <> external_page_id
			ORDER BY external_conversation_id, external_created_time DESC, external_timestamp DESC
		`, strings.Join(externalConversationIDs, "','")).
		Query()
	if err != nil {
		return nil, err
	}

	var fbExternalMessageIDs []dot.ID
	var id dot.ID
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		fbExternalMessageIDs = append(fbExternalMessageIDs, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	var fbExternalMessages model.FbExternalMessages
	if err := s.query().
		Where(sq.In("id", fbExternalMessageIDs)).
		Find(&fbExternalMessages); err != nil {
		return nil, err
	}

	if err := scheme.Convert([]*model.FbExternalMessage(fbExternalMessages), &result); err != nil {
		return nil, err
	}
	return
}

func (s *FbExternalMessageStore) GetFbExternalMessageDB() (*model.FbExternalMessage, error) {
	query := s.query().Where(s.preds)

	var fbExternalMessage model.FbExternalMessage
	err := query.ShouldGet(&fbExternalMessage)
	return &fbExternalMessage, err
}

func (s *FbExternalMessageStore) GetFbExternalMessage() (*fbmessaging.FbExternalMessage, error) {
	fbExternalMessage, err := s.GetFbExternalMessageDB()
	if err != nil {
		return nil, err
	}
	result := &fbmessaging.FbExternalMessage{}
	err = scheme.Convert(fbExternalMessage, result)
	if err != nil {
		return nil, err
	}
	return result, err
}
