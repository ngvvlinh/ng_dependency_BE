package sqlstore

import (
	"context"
	"fmt"
	"time"

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

type FbExternalCommentStoreFactory func(ctx context.Context) *FbExternalCommentStore

func NewFbExternalCommentStore(db *cmsql.Database) FbExternalCommentStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *FbExternalCommentStore {
		return &FbExternalCommentStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type FbExternalCommentStore struct {
	ft FbExternalCommentFilters

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
	sqlstore.Paging

	includeDeleted sqlstore.IncludeDeleted
}

func (s *FbExternalCommentStore) WithPaging(paging meta.Paging) *FbExternalCommentStore {
	s.Paging.WithPaging(paging)
	return s
}

func (s *FbExternalCommentStore) ExternalPostID(externalPostID string) *FbExternalCommentStore {
	s.preds = append(s.preds, s.ft.ByExternalPostID(externalPostID))
	return s
}

func (s *FbExternalCommentStore) ExternalUserID(externalUserID string) *FbExternalCommentStore {
	s.preds = append(s.preds, s.ft.ByExternalUserID(externalUserID))
	return s
}

func (s *FbExternalCommentStore) ExternalID(externalID string) *FbExternalCommentStore {
	s.preds = append(s.preds, s.ft.ByExternalID(externalID))
	return s
}

func (s *FbExternalCommentStore) ExternalIDs(externalIDs []string) *FbExternalCommentStore {
	s.preds = append(s.preds, sq.In("external_id", externalIDs))
	return s
}

func (s *FbExternalCommentStore) ExternalIDOrExternalParentID(externalID, externalParentID string) *FbExternalCommentStore {
	s.preds = append(s.preds, sq.NewExpr("external_id = ? OR external_parent_id = ?", externalID, externalParentID))
	return s
}

func (s *FbExternalCommentStore) ExternalParentIDIsNull() *FbExternalCommentStore {
	s.preds = append(s.preds, sq.NewExpr("external_parent_id is NULL"))
	return s
}

func (s *FbExternalCommentStore) ExternalParentID(externalID string) *FbExternalCommentStore {
	s.preds = append(s.preds, s.ft.ByExternalParentID(externalID))
	return s
}

func (s *FbExternalCommentStore) ExternalPageID(externalPageID string) *FbExternalCommentStore {
	s.preds = append(s.preds, s.ft.ByExternalPageID(externalPageID))
	return s
}

func (s *FbExternalCommentStore) ID(ID dot.ID) *FbExternalCommentStore {
	s.preds = append(s.preds, s.ft.ByID(ID))
	return s
}

func (s *FbExternalCommentStore) CreateOrUpdateFbExternalComment(fbExternalComment *fbmessaging.FbExternalComment) error {
	sqlstore.MustNoPreds(s.preds)
	fbExternalCommentDB := new(model.FbExternalComment)
	if err := scheme.Convert(fbExternalComment, fbExternalCommentDB); err != nil {
		return err
	}

	_, err := s.query().Upsert(fbExternalCommentDB)
	if err != nil {
		return err
	}

	// prepare data for search
	externalUserID := fbExternalCommentDB.ExternalUserID
	if externalUserID == fbExternalCommentDB.ExternalPageID {
		externalUserID = fbExternalCommentDB.ExternalParentUserID
	}
	commentSearch := &fbsearchmodel.FbExternalCommentSearch{
		ID:                  fbExternalCommentDB.ID,
		ExternalMessageNorm: normalizeText(fbExternalCommentDB.ExternalMessage),
		ExternalPageID:      fbExternalComment.ExternalPageID,
		CreatedAt:           fbExternalCommentDB.CreatedAt,
		ExternalPostID:      fbExternalCommentDB.ExternalPostID,
		ExternalUserID:      externalUserID,
	}
	if _, err = s.query().Upsert(commentSearch); err != nil {
		ll.Error(fmt.Sprintf("create fb_external_comment_search got error: %v", err))
	}

	return nil
}

func (s *FbExternalCommentStore) CreateFbExternalComments(fbExternalComments []*fbmessaging.FbExternalComment) error {
	sqlstore.MustNoPreds(s.preds)
	fbExternalCommentsDB := model.FbExternalComments(convert.Convert_fbmessaging_FbExternalComments_fbmessagingmodel_FbExternalComments(fbExternalComments))

	_, err := s.query().Upsert(&fbExternalCommentsDB)
	if err != nil {
		return err
	}

	// prepare data for search
	var commentSearchs fbsearchmodel.FbExternalCommentSearchs
	for _, cmt := range fbExternalCommentsDB {
		externalUserID := cmt.ExternalUserID
		if externalUserID == cmt.ExternalPageID {
			externalUserID = cmt.ExternalParentUserID
		}
		commentSearchs = append(commentSearchs, &fbsearchmodel.FbExternalCommentSearch{
			ID:                  cmt.ID,
			ExternalUserID:      externalUserID,
			ExternalMessageNorm: normalizeText(cmt.ExternalMessage),
			ExternalPageID:      cmt.ExternalPageID,
			CreatedAt:           cmt.CreatedAt,
			ExternalPostID:      cmt.ExternalPostID,
		})
	}
	if _, err = s.query().Upsert(&commentSearchs); err != nil {
		ll.Error(fmt.Sprintf("create fb_external_comment_search got error: %v", err))
	}

	return nil
}

func (s *FbExternalCommentStore) ListFbExternalCommentsDB() ([]*model.FbExternalComment, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	if !s.Paging.IsCursorPaging() && len(s.Paging.Sort) == 0 {
		s.Paging.Sort = []string{"-created_at"}
	}
	query, err := sqlstore.LimitSort(query, &s.Paging, SortFbExternalComment, s.ft.prefix)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterFbExternalComment)
	if err != nil {
		return nil, err
	}

	var fbExternalComments model.FbExternalComments
	err = query.Find(&fbExternalComments)
	if err != nil {
		return nil, err
	}
	s.Paging.Apply(fbExternalComments)
	return fbExternalComments, nil
}

func (s *FbExternalCommentStore) ListFbExternalComments() (result []*fbmessaging.FbExternalComment, err error) {
	fbExternalComments, err := s.ListFbExternalCommentsDB()
	if err != nil {
		return nil, err
	}
	if err = scheme.Convert(fbExternalComments, &result); err != nil {
		return nil, err
	}
	return
}

func (s *FbExternalCommentStore) ListFbExternalCommentsOfPage(externalPageID string) (result []*fbmessaging.FbExternalComment, err error) {
	s.preds = append(s.preds, sq.NewExpr(fmt.Sprintf(`
	((external_user_id = '%s' AND  external_parent_user_id IS NULL) OR 
	(external_user_id = '%s' AND external_parent_user_id = '%s'))`, externalPageID, externalPageID, externalPageID)))

	fbExternalComments, err := s.ListFbExternalCommentsDB()
	if err != nil {
		return nil, err
	}
	if err = scheme.Convert(fbExternalComments, &result); err != nil {
		return nil, err
	}
	return
}

func (s *FbExternalCommentStore) ExternalPageIDAndExternalUserID(externalPageID, externalUserID string) *FbExternalCommentStore {
	s.preds = append(s.preds, sq.NewExpr(fmt.Sprintf(`
		external_user_id = '%s' OR 
		(external_parent_user_id = '%s' AND (external_user_id = '%s' OR external_user_id = '%s'))
	`, externalUserID, externalUserID, externalUserID, externalPageID)))
	return s
}

func (s *FbExternalCommentStore) GetLatestExternalComment(
	externalPageID, externalPostID, externalUserID string,
) (*fbmessaging.FbExternalComment, error) {
	var fbExternalComment model.FbExternalComment

	query := s.query().
		Where(fmt.Sprintf(`
			external_post_id = '%s' AND 
			(
				(external_user_id = '%s') OR
				(external_user_id = '%s' AND external_parent_user_id = '%s')
			)
		`, externalPostID, externalUserID, externalPageID, externalUserID)).
		OrderBy("external_created_time desc", "id asc").
		Limit(1)
	if err := query.ShouldGet(&fbExternalComment); err != nil {
		return nil, err
	}

	result := fbmessaging.FbExternalComment{}
	if err := scheme.Convert(&fbExternalComment, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *FbExternalCommentStore) GetFbExternalCommentDB() (*model.FbExternalComment, error) {
	query := s.query().Where(s.preds)

	var fbExternalComment model.FbExternalComment
	err := query.ShouldGet(&fbExternalComment)
	return &fbExternalComment, err
}

func (s *FbExternalCommentStore) GetFbExternalComment() (*fbmessaging.FbExternalComment, error) {
	fbExternalComment, err := s.GetFbExternalCommentDB()
	if err != nil {
		return nil, err
	}
	result := &fbmessaging.FbExternalComment{}
	err = scheme.Convert(fbExternalComment, result)
	if err != nil {
		return nil, err
	}
	return result, err
}

func (s *FbExternalCommentStore) GetLatestCustomerExternalComment(
	externalPostID, externalUserID, externalPageID string,
) (*fbmessaging.FbExternalComment, error) {
	var fbExternalComment model.FbExternalComment

	if err := s.query().
		Where(s.ft.ByExternalPostID(externalPostID)).
		Where(s.ft.ByExternalUserID(externalUserID)).
		OrderBy("external_created_time DESC", "id DESC").
		Limit(1).
		ShouldGet(&fbExternalComment); err != nil {
		return nil, err
	}

	result := fbmessaging.FbExternalComment{}
	if err := scheme.Convert(&fbExternalComment, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *FbExternalCommentStore) GetLatestUpdatedActiveComment() (*fbmessaging.FbExternalComment, error) {
	var fbExternalComment model.FbExternalComment

	if err := s.query().
		Where(s.preds).
		Where(sq.NewExpr("deleted_at is NULL")).
		OrderBy("updated_at DESC").
		Limit(1).
		ShouldGet(&fbExternalComment); err != nil {
		return nil, err
	}

	result := fbmessaging.FbExternalComment{}
	if err := scheme.Convert(&fbExternalComment, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *FbExternalCommentStore) UpdateMessage(
	message string,
) (int, error) {
	query := s.query().Where(s.preds)
	return query.Table("fb_external_comment").UpdateMap(map[string]interface{}{
		"external_message": message,
	})
}

func (s *FbExternalCommentStore) SoftDelete() (int, error) {
	query := s.query().Where(s.preds)
	return query.Table("fb_external_comment").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
	})
}
