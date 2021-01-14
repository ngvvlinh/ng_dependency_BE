package sqlstore

import (
	"context"

	"github.com/lib/pq"

	"o.o/api/fabo/fbusering"
	"o.o/api/meta"
	"o.o/api/top/types/etc/status3"
	fbpagemodel "o.o/backend/com/fabo/main/fbpage/model"
	"o.o/backend/com/fabo/main/fbuser/convert"
	"o.o/backend/com/fabo/main/fbuser/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type FbExternalUserStoreFactory func(ctx context.Context) *FbExternalUserStore

var scheme = conversion.Build(convert.RegisterConversions)

func NewFbExternalUserStore(db *cmsql.Database) FbExternalUserStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *FbExternalUserStore {
		return &FbExternalUserStore{
			db:    db,
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type FbExternalUserStore struct {
	db *cmsql.Database
	ft FbExternalUserFilters

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
	sqlstore.Paging

	includeDeleted sqlstore.IncludeDeleted
}

func (s *FbExternalUserStore) ExternalIDs(externalIDs []string) *FbExternalUserStore {
	s.preds = append(s.preds, sq.In("external_id", externalIDs))
	return s
}

func (s *FbExternalUserStore) WithPaging(paging meta.Paging) *FbExternalUserStore {
	s.Paging.WithPaging(paging)
	return s
}

func (s *FbExternalUserStore) Filters(filters meta.Filters) *FbExternalUserStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *FbExternalUserStore) ExternalID(externalID string) *FbExternalUserStore {
	s.preds = append(s.preds, s.ft.ByExternalID(externalID))
	return s
}

func (s *FbExternalUserStore) ExternalPageID(externalPageID string) *FbExternalUserStore {
	s.preds = append(s.preds, s.ft.ByExternalPageID(externalPageID))
	return s
}

func (s *FbExternalUserStore) Status(status status3.Status) *FbExternalUserStore {
	s.preds = append(s.preds, s.ft.ByStatus(status))
	return s
}

func (s *FbExternalUserStore) UpdateStatus(status int) (int, error) {
	query := s.query().Where(s.preds)
	updateStatus, err := query.Table("fb_external_user").UpdateMap(map[string]interface{}{
		"status": status,
	})
	return updateStatus, err
}

func (s *FbExternalUserStore) CreateFbExternalUser(fbExternalUser *fbusering.FbExternalUser) error {
	sqlstore.MustNoPreds(s.preds)
	fbExternalUserDB := new(model.FbExternalUser)
	if err := scheme.Convert(fbExternalUser, fbExternalUserDB); err != nil {
		return err
	}
	_, err := s.query().Upsert(fbExternalUserDB)
	if err != nil {
		return err
	}

	var tempFbExternalUser model.FbExternalUser
	if err := s.query().Where(s.ft.ByExternalID(fbExternalUser.ExternalID)).ShouldGet(&tempFbExternalUser); err != nil {
		return err
	}
	fbExternalUser.CreatedAt = tempFbExternalUser.CreatedAt
	fbExternalUser.UpdatedAt = tempFbExternalUser.UpdatedAt

	return nil
}

func (s *FbExternalUserStore) CreateFbExternalUsers(fbExternalUsers []*fbusering.FbExternalUser) error {
	sqlstore.MustNoPreds(s.preds)
	fbExternalUsersDB := model.FbExternalUsers(convert.Convert_fbusering_FbExternalUsers_fbusermodel_FbExternalUsers(fbExternalUsers))

	_, err := s.query().Upsert(&fbExternalUsersDB)
	if err != nil {
		return err
	}
	return nil
}

func (s *FbExternalUserStore) GetFbExternalUserDB() (*model.FbExternalUser, error) {
	query := s.query().Where(s.preds)

	var fbExternalUser model.FbExternalUser
	err := query.ShouldGet(&fbExternalUser)
	return &fbExternalUser, err
}

func (s *FbExternalUserStore) GetFbExternalUser() (*fbusering.FbExternalUser, error) {
	fbExternalUser, err := s.GetFbExternalUserDB()
	if err != nil {
		return nil, err
	}
	result := &fbusering.FbExternalUser{}
	err = scheme.Convert(fbExternalUser, result)
	if err != nil {
		return nil, err
	}
	return result, err
}

func (s *FbExternalUserStore) ListFbExternalUsersDB() ([]*model.FbExternalUser, error) {
	query := s.query().Where(s.preds)
	if !s.Paging.IsCursorPaging() && len(s.Paging.Sort) == 0 {
		s.Paging.Sort = []string{"-created_at"}
	}
	query, err := sqlstore.LimitSort(query, &s.Paging, SortFbExternalUser, s.ft.prefix)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterFbExternalUser)
	if err != nil {
		return nil, err
	}

	var fbExternalUsers model.FbExternalUsers
	err = query.Find(&fbExternalUsers)
	if err != nil {
		return nil, err
	}
	s.Paging.Apply(fbExternalUsers)
	return fbExternalUsers, nil
}

func (s *FbExternalUserStore) ListFbExternalUsers() (result []*fbusering.FbExternalUser, err error) {
	fbExternalUsers, err := s.ListFbExternalUsersDB()
	if err != nil {
		return nil, err
	}
	if err = scheme.Convert(fbExternalUsers, &result); err != nil {
		return nil, err
	}
	return
}

func (s *FbExternalUserStore) GetFbExternalUserByShopID(
	extUserID dot.ID,
	shopID dot.ID,
) (*fbusering.FbExternalUser, error) {
	extUser, err := s.ExternalID(extUserID.String()).GetFbExternalUser()
	if err != nil {
		return nil, err
	}

	// Dựa vào ExternalPageID để biết fb_external_user có thuộc shop_id này không
	extPage := &fbpagemodel.FbExternalPage{}
	err = s.query().
		Table("fb_external_page").
		Where("shop_id = ? AND external_id = ?", shopID, extUser.ExternalPageID).
		ShouldGet(extPage)
	if err != nil {
		if cm.ErrorCode(err) == cm.NotFound {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Khách hàng không thuộc shop")
		}
		return nil, err
	}

	return extUser, nil
}

func (s *FbExternalUserStore) UpdateUserTags(tagIDs []dot.ID) error {
	updatedIDs := map[string]interface{}{
		"tag_ids": pq.Array(tagIDs),
	}
	return s.query().Table("fb_external_user").Where(s.preds).ShouldUpdateMap(updatedIDs)
}
