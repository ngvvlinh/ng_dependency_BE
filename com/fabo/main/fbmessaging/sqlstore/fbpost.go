package sqlstore

import (
	"context"
	"time"

	"github.com/lib/pq"

	"o.o/api/fabo/fbmessaging"
	"o.o/api/fabo/fbmessaging/fb_live_video_status"
	"o.o/api/fabo/fbmessaging/fb_status_type"
	"o.o/api/meta"
	"o.o/backend/com/fabo/main/fbmessaging/convert"
	"o.o/backend/com/fabo/main/fbmessaging/model"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type FbExternalPostStoreFactory func(ctx context.Context) *FbExternalPostStore

var scheme = conversion.Build(convert.RegisterConversions)

func NewFbExternalPostStore(db *cmsql.Database) FbExternalPostStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *FbExternalPostStore {
		return &FbExternalPostStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type FbExternalPostStore struct {
	ft FbExternalPostFilters

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
	sqlstore.Paging

	includeDeleted sqlstore.IncludeDeleted
}

func (s *FbExternalPostStore) WithPaging(paging meta.Paging) *FbExternalPostStore {
	s.Paging.WithPaging(paging)
	return s
}

func (s *FbExternalPostStore) ID(ID dot.ID) *FbExternalPostStore {
	s.preds = append(s.preds, s.ft.ByID(ID))
	return s
}

func (s *FbExternalPostStore) IDs(IDs []dot.ID) *FbExternalPostStore {
	s.preds = append(s.preds, sq.In("id", IDs))
	return s
}

func (s *FbExternalPostStore) ExternalStatusType(statusType fb_status_type.FbStatusType) *FbExternalPostStore {
	s.preds = append(s.preds, s.ft.ByStatusType(statusType))
	return s
}

func (s *FbExternalPostStore) LiveVideoStatus(liveVideoStatus fb_live_video_status.FbLiveVideoStatus) *FbExternalPostStore {
	s.preds = append(s.preds, s.ft.ByLiveVideoStatus(liveVideoStatus))
	return s
}

func (s *FbExternalPostStore) ExternalID(externalID string) *FbExternalPostStore {
	s.preds = append(s.preds, s.ft.ByExternalID(externalID))
	return s
}

func (s *FbExternalPostStore) ExternalUserID(externalUserID string) *FbExternalPostStore {
	s.preds = append(s.preds, s.ft.ByExternalUserID(externalUserID))
	return s
}

func (s *FbExternalPostStore) ExternalParentID(externalParentID string) *FbExternalPostStore {
	s.preds = append(s.preds, s.ft.ByExternalParentID(externalParentID))
	return s
}

func (s *FbExternalPostStore) ExternalCreatedTime(created time.Time) *FbExternalPostStore {
	s.preds = append(s.preds, s.ft.ByExternalCreatedTime(created))
	return s
}

func (s *FbExternalPostStore) IsLiveVideo(isLiveVideo bool) *FbExternalPostStore {
	s.preds = append(s.preds, s.ft.ByIsLiveVideo(isLiveVideo))
	return s
}

func (s *FbExternalPostStore) ExternalIDs(externalIDs []string) *FbExternalPostStore {
	s.preds = append(s.preds, sq.In("external_id", externalIDs))
	return s
}

func (s *FbExternalPostStore) ExternalPageIDs(externalPageIDs []string) *FbExternalPostStore {
	s.preds = append(s.preds, sq.In("external_page_id", externalPageIDs))
	return s
}

func (s *FbExternalPostStore) ExternalPageIDsOrExternalUserID(externalPageIDs []string, externalUserID string) *FbExternalPostStore {
	s.preds = append(s.preds, sq.NewExpr("external_page_id IN (?) OR external_user_id = ?", pq.Array(externalPageIDs), externalUserID))
	return s
}

func (s *FbExternalPostStore) CreateFbExternalPost(fbExternalPost *fbmessaging.FbExternalPost) error {
	fbExternalPostDB := new(model.FbExternalPost)
	if err := scheme.Convert(fbExternalPost, fbExternalPostDB); err != nil {
		return err
	}

	_, err := s.query().Insert(fbExternalPostDB)
	if err != nil {
		return err
	}

	var tempFbExternalPost model.FbExternalPost
	if err := s.query().Where(s.ft.ByID(fbExternalPost.ID)).ShouldGet(&tempFbExternalPost); err != nil {
		return err
	}
	fbExternalPost.CreatedAt = tempFbExternalPost.CreatedAt
	fbExternalPost.UpdatedAt = tempFbExternalPost.UpdatedAt

	return nil
}

func (s *FbExternalPostStore) UpsertFbExternalPost(fbExternalPost *fbmessaging.FbExternalPost) error {
	sqlstore.MustNoPreds(s.preds)
	fbExternalPostDB := new(model.FbExternalPost)
	if err := scheme.Convert(fbExternalPost, fbExternalPostDB); err != nil {
		return err
	}

	_, err := s.query().Upsert(fbExternalPostDB)
	if err != nil {
		return err
	}

	var tempFbExternalPost model.FbExternalPost
	if err := s.query().Where(s.ft.ByID(fbExternalPost.ID)).ShouldGet(&tempFbExternalPost); err != nil {
		return err
	}
	fbExternalPost.CreatedAt = tempFbExternalPost.CreatedAt
	fbExternalPost.UpdatedAt = tempFbExternalPost.UpdatedAt

	return nil
}

func (s *FbExternalPostStore) UpsertFbExternalPosts(fbExternalPosts []*fbmessaging.FbExternalPost) error {
	sqlstore.MustNoPreds(s.preds)
	fbExternalPostsDB := model.FbExternalPosts(convert.Convert_fbmessaging_FbExternalPosts_fbmessagingmodel_FbExternalPosts(fbExternalPosts))

	_, err := s.query().Upsert(&fbExternalPostsDB)
	if err != nil {
		return err
	}
	return nil
}

func (s *FbExternalPostStore) UpdateFbExternalPost(fbExternalPost *fbmessaging.FbExternalPost) error {
	fbExternalPostDB := new(model.FbExternalPost)
	if err := scheme.Convert(fbExternalPost, fbExternalPostDB); err != nil {
		return err
	}

	return s.query().Where(s.preds).ShouldUpdate(fbExternalPostDB)
}

// func (s *FbExternalPostStore) UpdateTotalComments(externalID string) error {
// 	sql := `UPDATE fb_external_post
// SET total_comments = (
// 	SELECT count(external_id)
// 	FROM fb_external_comment
// 	WHERE external_post_id = ?
// )
// WHERE external_id = ?`
// 	_, err := s.query().SQL(sql, externalID, externalID).Exec()
// 	return err
// }

func (s *FbExternalPostStore) GetFbExternalPostDB() (*model.FbExternalPost, error) {
	query := s.query().Where(s.preds)

	var fbExternalPost model.FbExternalPost
	err := query.ShouldGet(&fbExternalPost)
	return &fbExternalPost, err
}

func (s *FbExternalPostStore) GetFbExternalPost() (*fbmessaging.FbExternalPost, error) {
	fbExternalPost, err := s.GetFbExternalPostDB()
	if err != nil {
		return nil, err
	}
	result := &fbmessaging.FbExternalPost{}
	err = scheme.Convert(fbExternalPost, result)
	if err != nil {
		return nil, err
	}
	return result, err
}

func (s *FbExternalPostStore) ListFbExternalPostsDB() ([]*model.FbExternalPost, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	if !s.Paging.IsCursorPaging() && len(s.Paging.Sort) == 0 {
		s.Paging.Sort = []string{"-created_at"}
	}
	query, err := sqlstore.LimitSort(query, &s.Paging, SortFbExternalPost, s.ft.prefix)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterFbExternalPost)
	if err != nil {
		return nil, err
	}

	var fbExternalPosts model.FbExternalPosts
	err = query.Find(&fbExternalPosts)
	if err != nil {
		return nil, err
	}
	s.Paging.Apply(fbExternalPosts)

	// calc total comment
	externalPostIDs := []string{}
	for _, post := range fbExternalPosts {
		externalPostIDs = append(externalPostIDs, post.ExternalID)
	}
	postExtendeds, err := s.CalcTotalComments(externalPostIDs)
	if err != nil {
		return nil, err
	}
	if postExtendeds != nil {
		for _, post := range fbExternalPosts {
			if postTotalComment, ok := postExtendeds[post.ExternalID]; ok {
				post.TotalComments = postTotalComment
			}
		}
	}

	return fbExternalPosts, nil
}

func (s *FbExternalPostStore) CalcTotalComments(externalPostIDs []string) (map[string]int, error) {
	if len(externalPostIDs) == 0 {
		return nil, nil
	}
	var posts model.FbExternalPostFtTotalComments
	var res = make(map[string]int)
	query := s.query().From("fb_external_comment").In("external_post_id", externalPostIDs).GroupBy("external_post_id")

	if err := query.Find(&posts); err != nil {
		return nil, err
	}
	for _, post := range posts {
		res[post.ExternalPostID] = post.Count
	}
	return res, nil
}

func (s *FbExternalPostStore) ListFbExternalPosts() (result []*fbmessaging.FbExternalPost, err error) {
	fbExternalPosts, err := s.ListFbExternalPostsDB()
	if err != nil {
		return nil, err
	}
	if err = scheme.Convert(fbExternalPosts, &result); err != nil {
		return nil, err
	}
	return
}

func (s *FbExternalPostStore) UpdatePostMessageAndPicture(message string, picture string) error {
	query := s.query().Where(s.preds)
	return query.Table("fb_external_post").ShouldUpdateMap(map[string]interface{}{
		"external_message": message,
		"external_picture": picture,
	})
}

func (s *FbExternalPostStore) SoftDelete() (int, error) {
	query := s.query().Where(s.preds)
	return query.Table("fb_external_post").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
	})
}
