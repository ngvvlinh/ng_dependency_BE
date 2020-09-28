package sqlstore

import (
	"context"

	"o.o/api/main/identity"
	"o.o/backend/com/main/identity/convert"
	identitymodel "o.o/backend/com/main/identity/model"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/backend/pkg/etop/model"
	"o.o/capi/dot"
)

type UserRefSaffStoreFactory func(context.Context) *UserRefSaffStore

type UserRefSaffStore struct {
	query cmsql.QueryFactory
	ft    UserRefSaffFilters
	sqlstore.Paging
	preds []interface{}
}

func NewUserRefSaffStore(db *cmsql.Database) UserRefSaffStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *UserRefSaffStore {
		return &UserRefSaffStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

func (s *UserRefSaffStore) CreateUserRefSaffDB(ref *identitymodel.UserRefSaff) error {
	return s.query().ShouldInsert(ref)
}

func (s *UserRefSaffStore) CreateUserRefSaff(ref *identity.UserRefSaff) error {
	var refModel *identitymodel.UserRefSaff
	refModel = convert.Convert_identity_UserRefSaff_identitymodel_UserRefSaff(ref, refModel)
	return s.query().ShouldInsert(refModel)
}

func (s *UserRefSaffStore) ByUserID(id dot.ID) *UserRefSaffStore {
	s.preds = append(s.preds, s.ft.ByUserID(id))
	return s
}

func (s *UserRefSaffStore) ByRefSale(phone string) *UserRefSaffStore {
	s.preds = append(s.preds, s.ft.ByRefSale(phone))
	return s
}

func (s *UserRefSaffStore) ByRefAff(phone string) *UserRefSaffStore {
	s.preds = append(s.preds, s.ft.ByRefAff(phone))
	return s
}

func (s *UserRefSaffStore) GetUserRefSaffDB() (*identitymodel.UserRefSaff, error) {
	userRef := &identitymodel.UserRefSaff{}
	err := s.query().Where(s.preds...).ShouldGet(userRef)
	if err != nil {
		return nil, err
	}
	return userRef, nil
}

func (s *UserRefSaffStore) GetUserRefSaff() (*identity.UserRefSaff, error) {
	userRefModel, err := s.GetUserRefSaffDB()
	if err != nil {
		return nil, err
	}

	var userRef *identity.UserRefSaff
	userRef = convert.Convert_identitymodel_UserRefSaff_identity_UserRefSaff(userRefModel, userRef)
	return userRef, nil
}

func (s *UserRefSaffStore) Update(ref *identity.UserRefSaff) error {
	var update map[string]interface{}
	if ref.RefAff.Valid {
		update["ref_aff"] = ref.RefAff
	}
	if ref.RefSale.Valid {
		update["ref_sale"] = ref.RefAff
	}
	if len(update) == 0 {
		return nil
	}
	return s.query().Where(s.preds...).
		Table("user_ref_saff").
		ShouldUpdateMap(update)
}

func (s *UserRefSaffStore) ListUserRefSaffDB() (identitymodel.UserRefSaffs, error) {
	query := s.query().Where(s.preds)
	query, err := sqlstore.LimitSort(query, &s.Paging, SortUser)
	if err != nil {
		return nil, err
	}
	var userSaffs identitymodel.UserRefSaffs
	err = query.Find(&userSaffs)
	return userSaffs, err
}

func (s *UserRefSaffStore) ListUserRefSaff() ([]*identity.UserRefSaff, error) {
	userSaffDB, err := s.ListUserRefSaffDB()
	if err != nil {
		return nil, err
	}
	res := convert.Convert_identitymodel_UserRefSaffs_identity_UserRefSaffs(userSaffDB)
	return res, nil
}
