package sqlstore

import (
	"context"
	"fmt"
	"time"

	"o.o/api/fabo/fbmessaging"
	"o.o/api/fabo/fbmessaging/fb_customer_conversation_type"
	"o.o/api/meta"
	fbsearchmodel "o.o/backend/com/fabo/main/fbcustomerconversationsearch/model"
	"o.o/backend/com/fabo/main/fbmessaging/convert"
	"o.o/backend/com/fabo/main/fbmessaging/model"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/backend/pkg/common/validate"
	"o.o/capi/dot"
	"o.o/common/l"
)

var (
	ll = l.New()
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

func (s *FbCustomerConversationStore) WithPaging(
	paging meta.Paging) *FbCustomerConversationStore {
	s.Paging.WithPaging(paging)
	return s
}

func (s *FbCustomerConversationStore) ID(id dot.ID) *FbCustomerConversationStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *FbCustomerConversationStore) IDs(ids []dot.ID) *FbCustomerConversationStore {
	s.preds = append(s.preds, sq.In("id", ids))
	return s
}

func (s *FbCustomerConversationStore) ExternalID(externalID string) *FbCustomerConversationStore {
	s.preds = append(s.preds, s.ft.ByExternalID(externalID))
	return s
}

func (s *FbCustomerConversationStore) ExternalIDs(externalIDs []string) *FbCustomerConversationStore {
	s.preds = append(s.preds, sq.In("external_id", externalIDs))
	return s
}

func (s *FbCustomerConversationStore) ExternalPageIDs(externalPageIDs []string) *FbCustomerConversationStore {
	s.preds = append(s.preds, sq.In("external_page_id", externalPageIDs))
	return s
}

func (s *FbCustomerConversationStore) ExternalUserIDs(extUserIDs []string) *FbCustomerConversationStore {
	s.preds = append(s.preds, sq.In("external_user_id", extUserIDs))
	return s
}

func (s *FbCustomerConversationStore) ExternalIDAndExternalUserID(externalID, externalUserID string) *FbCustomerConversationStore {
	s.preds = append(s.preds, sq.NewExpr("external_id = ? AND external_user_id = ?", externalID, externalUserID))
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

func (s *FbCustomerConversationStore) FbExternalUserID(fbExternalUserID string) *FbCustomerConversationStore {
	s.preds = append(s.preds, s.ft.ByExternalUserID(fbExternalUserID))
	return s
}

func (s *FbCustomerConversationStore) Type(typ fb_customer_conversation_type.FbCustomerConversationType) *FbCustomerConversationStore {
	s.preds = append(s.preds, s.ft.ByType(typ.Enum()))
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

	// prepare data for search
	customerConvSearch := &fbsearchmodel.FbCustomerConversationSearch{
		ID:                   fbCustomerConversationDB.ID,
		ExternalUserNameNorm: normalizeText(fbCustomerConversationDB.ExternalUserName),
		CreatedAt:            fbCustomerConversationDB.CreatedAt,
		ExternalPageID:       fbCustomerConversationDB.ExternalPageID,
	}
	_, err = s.query().Upsert(customerConvSearch)
	if err != nil {
		ll.Error(fmt.Sprintf("create fb_customer_conversation_search got error: %v", err))
	}

	return nil
}

func (s *FbCustomerConversationStore) CreateFbCustomerConversations(fbCustomerConversations []*fbmessaging.FbCustomerConversation) error {
	sqlstore.MustNoPreds(s.preds)
	fbCustomerConversationsDB := model.FbCustomerConversations(convert.Convert_fbmessaging_FbCustomerConversations_fbmessagingmodel_FbCustomerConversations(fbCustomerConversations))

	_, err := s.query().Upsert(&fbCustomerConversationsDB)
	if err != nil {
		return err
	}

	// prepare data for search
	var customerConvSearchs fbsearchmodel.FbCustomerConversationSearchs
	for _, convDB := range fbCustomerConversationsDB {
		customerConvSearchs = append(customerConvSearchs, &fbsearchmodel.FbCustomerConversationSearch{
			ID:                   convDB.ID,
			ExternalUserNameNorm: normalizeText(convDB.ExternalUserName),
			CreatedAt:            convDB.CreatedAt,
			ExternalPageID:       convDB.ExternalPageID,
		})
	}
	_, err = s.query().Upsert(&customerConvSearchs)
	if err != nil {
		ll.Error(fmt.Sprintf("create fb_customer_conversation_search got error: %v", err))
	}

	return nil
}

func (s *FbCustomerConversationStore) GetFbCustomerConversationDB() (*model.FbCustomerConversation, error) {
	query := s.query().Where(s.preds)

	var fbCustomerConversation model.FbCustomerConversation
	err := query.ShouldGet(&fbCustomerConversation)
	return &fbCustomerConversation, err
}

func (s *FbCustomerConversationStore) GetFbCustomerConversation() (*fbmessaging.FbCustomerConversation, error) {
	fbCustomerConversation, err := s.GetFbCustomerConversationDB()
	if err != nil {
		return nil, err
	}
	result := &fbmessaging.FbCustomerConversation{}
	err = scheme.Convert(fbCustomerConversation, result)
	if err != nil {
		return nil, err
	}
	return result, err
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

func (s *FbCustomerConversationStore) SoftDelete() (int, error) {
	query := s.query().Where(s.preds)
	return query.Table("fb_customer_conversation").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
	})
}

func normalizeText(s string) string {
	return validate.NormalizedSearchToTsVector(validate.NormalizeSearch(s))
}
