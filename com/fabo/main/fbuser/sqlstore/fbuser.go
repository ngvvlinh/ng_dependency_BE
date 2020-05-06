package sqlstore

import (
	"context"

	"o.o/api/fabo/fbusering"
	"o.o/api/meta"
	"o.o/api/top/types/etc/status3"
	"o.o/backend/com/fabo/main/fbuser/convert"
	"o.o/backend/com/fabo/main/fbuser/model"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
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

func (s *FbExternalUserStore) ExternalID(externalID string) *FbExternalUserStore {
	s.preds = append(s.preds, s.ft.ByExternalID(externalID))
	return s
}

func (s *FbExternalUserStore) UserID(userID dot.ID) *FbExternalUserStore {
	s.preds = append(s.preds, s.ft.ByUserID(userID))
	return s
}

func (s *FbExternalUserStore) Status(status status3.Status) *FbExternalUserStore {
	s.preds = append(s.preds, s.ft.ByStatus(status))
	return s
}

func (s *FbExternalUserStore) UpdateStatus(status int) (int, error) {
	query := s.query().Where(s.preds)
	updateStatus, err := query.Table("fb_user").UpdateMap(map[string]interface{}{
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
	if err := s.query().Where(s.ft.ByID(fbExternalUser.ID)).ShouldGet(&tempFbExternalUser); err != nil {
		return err
	}
	fbExternalUser.CreatedAt = tempFbExternalUser.CreatedAt
	fbExternalUser.UpdatedAt = tempFbExternalUser.UpdatedAt

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
