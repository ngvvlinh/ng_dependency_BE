package sqlstore

import (
	"context"
	"time"

	"o.o/api/main/identity"
	"o.o/backend/com/main/identity/convert"
	identitymodel "o.o/backend/com/main/identity/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sq/core"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/backend/pkg/common/validate"
	"o.o/backend/pkg/etop/model"
	"o.o/capi/dot"
)

type AccountUserStoreFactory func(context.Context) *AccountUserStore

func NewAccountUserStore(db *cmsql.Database) AccountUserStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *AccountUserStore {
		return &AccountUserStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type AccountUserStore struct {
	query cmsql.QueryFactory
	preds []interface{}
	ft    AccountUserFilters

	includeDeleted sqlstore.IncludeDeleted
}

func (s *AccountUserStore) ByAccountID(id dot.ID) *AccountUserStore {
	s.preds = append(s.preds, s.ft.ByAccountID(id))
	return s
}

func (s *AccountUserStore) ByUserIDs(ids []dot.ID) *AccountUserStore {
	s.preds = append(s.preds, sq.In("user_id", ids))
	return s
}

func (s *AccountUserStore) ByUserID(id dot.ID) *AccountUserStore {
	s.preds = append(s.preds, s.ft.ByUserID(id))
	return s
}

func (s *AccountUserStore) ByRoles(roles ...string) *AccountUserStore {
	s.preds = append(s.preds, sq.NewExpr("roles @> ?", core.Array{V: roles}))
	return s
}

func (s *AccountUserStore) GetAccountUserDB() (*identitymodel.AccountUser, error) {
	var acc identitymodel.AccountUser
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	err := query.ShouldGet(&acc)
	return &acc, err
}

func (s *AccountUserStore) GetAccountUser() (*identity.AccountUser, error) {
	accountUsersDB, err := s.GetAccountUserDB()
	if err != nil {
		return nil, err
	}
	var res *identity.AccountUser
	res = convert.Convert_identitymodel_AccountUser_identity_AccountUser(accountUsersDB, res)
	return res, nil
}

type DeleteAccountUserArgs struct {
	AccountID dot.ID
	UserID    dot.ID
}

func (s *AccountUserStore) DeleteAccountUser(args DeleteAccountUserArgs) error {
	if args.AccountID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing AccountID")
	}
	if args.UserID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing UserID")
	}

	update := &identitymodel.AccountUser{
		DeletedAt: time.Now(),
	}

	if err := s.query().Where(s.ft.ByUserID(args.UserID)).
		Where(s.ft.ByAccountID(args.AccountID)).
		ShouldDelete(update); err != nil {
		return err
	}
	return nil
}

func (s *AccountUserStore) SoftDeleteAccountUsers() (int, error) {
	query := s.query().Where(s.preds)
	return query.Update(&identitymodel.AccountUser{DeletedAt: time.Now()})
}

func (s *AccountUserStore) ListAccountUserDBs() ([]*identitymodel.AccountUser, error) {
	query := s.query().Where(s.preds)

	var accountUser identitymodel.AccountUsers
	err := query.Find(&accountUser)
	return accountUser, err
}

func (s *AccountUserStore) ListAccountUsers() ([]*identity.AccountUser, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())

	var accountUser identitymodel.AccountUsers
	err := query.Find(&accountUser)
	return convert.Convert_identitymodel_AccountUsers_identity_AccountUsers(accountUser), err
}

func (s *AccountUserStore) CreateAccountUser(au *identity.AccountUser) error {
	if au.UserID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing UserID")
	}
	var user = new(identitymodel.User)
	if err := s.query().Table("user").Where("id = ?", au.UserID).ShouldGet(user); err != nil {
		return err
	}
	sqlstore.MustNoPreds(s.preds)
	var auDB identitymodel.AccountUser
	if err := scheme.Convert(au, &auDB); err != nil {
		return err
	}
	s.normalizeSearchFields(&auDB, user)
	return s.query().ShouldInsert(&auDB)
}

func (s *AccountUserStore) UpdateAccountUser(au *identity.AccountUser) error {
	query := s.query().Where(s.preds)
	var auDB identitymodel.AccountUser
	if err := scheme.Convert(au, &auDB); err != nil {
		return err
	}
	return query.ShouldUpdate(&auDB)
}

func (s *AccountUserStore) UpdateExtensionNumberNorm(extNumberNorm string) (int, error) {
	return s.query().Where(s.preds).Table("account_user").UpdateMap(
		map[string]interface{}{
			"extension_number_norm": extNumberNorm,
		})
}

func (s *AccountUserStore) normalizeSearchFields(auDB *identitymodel.AccountUser, user *identitymodel.User) {
	auDB.Phone = user.Phone
	auDB.FullNameNorm = validate.NormalizeSearchCharacter(user.FullName)
	auDB.PhoneNorm = validate.NormalizeSearchCharacter(user.Phone)
}
