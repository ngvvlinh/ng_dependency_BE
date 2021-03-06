package sqlstore

import (
	"context"
	"time"

	"o.o/api/main/identity"
	"o.o/api/meta"
	"o.o/api/top/types/etc/status3"
	"o.o/backend/com/main/identity/convert"
	identitymodel "o.o/backend/com/main/identity/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/backend/pkg/common/validate"
	"o.o/backend/pkg/etop/model"
	"o.o/capi/dot"
	"o.o/capi/filter"
)

type UserStoreFactory func(context.Context) *UserStore

func NewUserStore(db *cmsql.Database) UserStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *UserStore {
		return &UserStore{
			query: cmsql.NewQueryFactory(ctx, db),
			ctx:   ctx,
		}
	}
}

type UserStore struct {
	query cmsql.QueryFactory
	ft    UserFilters
	sqlstore.Paging
	preds []interface{}
	ctx   context.Context

	includeWLPartnerUser bool
}

func (s *UserStore) IncludeWLPartnerUser() *UserStore {
	s.includeWLPartnerUser = true
	return s
}

func (s *UserStore) ByID(id dot.ID) *UserStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *UserStore) ByIDs(ids []dot.ID) *UserStore {
	s.preds = append(s.preds, sq.In("id", ids))
	return s
}

func (s *UserStore) ByCreatedAt(from time.Time, to time.Time) *UserStore {
	s.preds = append(s.preds, sq.NewExpr("created_at > ? AND created_at <= ?", from, to))
	return s
}

func (s *UserStore) ByCreatedAtFrom(from time.Time) *UserStore {
	s.preds = append(s.preds, sq.NewExpr("created_at > ?", from))
	return s
}
func (s *UserStore) ByCreatedAtTo(to time.Time) *UserStore {
	s.preds = append(s.preds, sq.NewExpr("created_at <= ?", to))
	return s
}

func (s *UserStore) ByPhone(phone string) *UserStore {
	s.preds = append(s.preds, s.ft.ByPhone(phone))
	return s
}

func (s *UserStore) ByEmail(email string) *UserStore {
	s.preds = append(s.preds, s.ft.ByEmail(email))
	return s
}

func (s *UserStore) ByNameNorm(nameNorm filter.FullTextSearch) *UserStore {
	s.preds = append(s.preds, sq.NewExpr(`"name_norm" @@ ?::tsquery`, validate.NormalizeFullTextSearchQueryAnd(nameNorm)))
	return s
}

func (s *UserStore) ByFullNameNorm(name filter.FullTextSearch) *UserStore {
	s.preds = append(s.preds, s.ft.Filter(`full_name_norm @@ ?::tsquery`, validate.NormalizeFullTextSearchQueryAnd(name)))
	return s
}

func (s *UserStore) ByWLPartnerID(WLPartnerID dot.ID) *UserStore {
	s.preds = append(s.preds, s.ft.ByWLPartnerID(WLPartnerID))
	return s
}

func (s *UserStore) WithPaging(paging meta.Paging) *UserStore {
	s.Paging.WithPaging(paging)
	return s
}

func (s *UserStore) CreateUser(user *identity.User) (*identity.User, error) {
	sqlstore.MustNoPreds(s.preds)
	if user.ID == 0 {
		user.ID = cm.NewID()
	}
	userDB := &identitymodel.User{}
	if err := scheme.Convert(user, userDB); err != nil {
		return nil, err
	}
	userDB.WLPartnerID = wl.GetWLPartnerID(s.ctx)
	if err := s.query().ShouldInsert(userDB); err != nil {
		return nil, err
	}
	return s.ByID(user.ID).GetUser()
}

func (s *UserStore) GetUserDB() (*identitymodel.User, error) {
	var user identitymodel.User
	query := s.query().Where(s.preds)
	if !s.includeWLPartnerUser {
		query = s.FilterByWhiteLabelPartner(query, wl.GetWLPartnerID(s.ctx))
	}
	err := query.ShouldGet(&user)
	return &user, err
}

func (s *UserStore) GetUserFtRefSaffDB(ctx context.Context) (*identitymodel.UserFtRefSaff, error) {
	var user identitymodel.UserFtRefSaff
	query := s.query().Where(s.preds)
	if !s.includeWLPartnerUser {
		query = s.FilterByWhiteLabelPartner(query, wl.GetWLPartnerID(ctx))
	}
	err := query.ShouldGet(&user)
	return &user, err
}

func (s *UserStore) GetUserFtRefSaff(ctx context.Context) (*identity.UserFtRefSaff, error) {
	user, err := s.GetUserFtRefSaffDB(ctx)
	if err != nil {
		return nil, err
	}
	return convert.UserFtRefSaff(user), err
}

func (s *UserStore) UpdateUserEmail(email string) (int, error) {
	return s.query().Where(s.preds).Table("user").UpdateMap(
		map[string]interface{}{
			"email":             email,
			"email_verified_at": time.Now(),
		})
}

func (s *UserStore) UpdateUserPhone(phone string) (int, error) {
	return s.query().Where(s.preds).Table("user").UpdateMap(
		map[string]interface{}{
			"phone":             phone,
			"phone_verified_at": time.Now(),
		})
}

func (s *UserStore) UpdateUser(args *identity.User) error {
	var result = &identitymodel.User{}
	err := scheme.Convert(args, result)
	if err != nil {
		return err
	}
	err = s.UpdateUserDB(result)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStore) UpdateUserDB(args *identitymodel.User) error {
	query := s.query().Where(s.preds)
	return query.ShouldUpdate(args)
}

func (s *UserStore) UpdateUserStatus(userID dot.ID, blockedBy dot.ID, status status3.Status, blockReason string) (int, error) {
	var updateMap = map[string]interface{}{
		"status": status,
	}
	// update blocked_at and blocked_by if block user
	if status == status3.N {
		updateMap["blocked_at"] = time.Now()
		updateMap["blocked_by"] = blockedBy
		updateMap["blocked_reason"] = blockReason
	}
	return s.query().Where("user_id = ?", userID).Table("user").UpdateMap(updateMap)

}

func (s *UserStore) UnblockUser() (int, error) {
	return s.query().Where(s.preds).Table("user").UpdateMap(
		map[string]interface{}{
			"status":     1,
			"blocked_at": nil,
		})
}

func (s *UserStore) GetUser() (*identity.User, error) {
	result, err := s.GetUserDB()
	if err != nil {
		return nil, err
	}
	return convert.User(result), nil
}

func (s *UserStore) ListUsersDB() ([]*identitymodel.User, error) {
	query := s.query().Where(s.preds)
	// default sort by created_at
	if len(s.Paging.Sort) == 0 {
		s.Paging.Sort = append(s.Paging.Sort, "-created_at")
	}
	query, err := sqlstore.LimitSort(query, &s.Paging, SortUser)
	if err != nil {
		return nil, err
	}
	var users identitymodel.Users
	err = query.Find(&users)
	return users, err
}

func (s *UserStore) ListUsers() (users []*identity.User, err error) {
	usersDB, err := s.ListUsersDB()
	if err != nil {
		return nil, err
	}
	err = scheme.Convert(usersDB, &users)
	return users, err
}

type UpdateRefferenceIDArgs struct {
	UserID    dot.ID
	RefUserID dot.ID
	RefSaleID dot.ID
}

func (s *UserStore) UpdateUserRefferenceID(args *UpdateRefferenceIDArgs) error {
	if args.UserID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing userID")
	}
	if args.RefUserID == 0 && args.RefSaleID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing reference userID")
	}
	user := &identitymodel.User{
		RefUserID: args.RefUserID,
		RefSaleID: args.RefSaleID,
	}

	if err := s.query().Where(s.ft.ByID(args.UserID)).ShouldUpdate(user); err != nil {
		return err
	}
	return nil
}

func (s *UserStore) FilterByWhiteLabelPartner(query cmsql.Query, wlPartnerID dot.ID) cmsql.Query {
	if wlPartnerID != 0 {
		return query.Where(s.ft.ByWLPartnerID(wlPartnerID))
	}
	return query.Where(s.ft.NotBelongWLPartner())
}

func (s *UserStore) ListUserFtRefSaffDBs() ([]*identitymodel.UserFtRefSaff, error) {
	query := s.query().Where(s.preds)
	// default sort by created_at
	if len(s.Paging.Sort) == 0 {
		s.Paging.Sort = append(s.Paging.Sort, "-created_at")
	}
	query, err := sqlstore.LimitSort(query, &s.Paging, SortUser)
	if err != nil {
		return nil, err
	}
	var users identitymodel.UserFtRefSaffs
	err = query.Find(&users)
	return users, err
}

func (s *UserStore) ListUserFtRefSaffs() (users []*identity.UserFtRefSaff, err error) {
	usersDBs, err := s.ListUserFtRefSaffDBs()
	if err != nil {
		return nil, err
	}
	for _, userDB := range usersDBs {
		users = append(users, convert.UserFtRefSaff(userDB))
	}
	return users, err
}
