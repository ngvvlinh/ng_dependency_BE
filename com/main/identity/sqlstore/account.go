package sqlstore

import (
	"context"
	"time"

	"etop.vn/api/top/types/etc/account_type"

	"etop.vn/api/main/identity"
	"etop.vn/backend/com/main/identity/convert"
	identitymodel "etop.vn/backend/com/main/identity/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/etop/sqlstore"
	"etop.vn/capi/dot"
)

type AccountStoreFactory func(context.Context) *AccountStore

func NewAccountStore(db *cmsql.Database) AccountStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *AccountStore {
		return &AccountStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type AccountStore struct {
	query       cmsql.QueryFactory
	preds       []interface{}
	shopFt      sqlstore.ShopFilters
	affiliateFt AffiliateFilters
}

func (s *AccountStore) ShopByID(id dot.ID) *AccountStore {
	s.preds = append(s.preds, s.shopFt.ByID(id))
	return s
}

func (s *AccountStore) AffiliateByID(id dot.ID) *AccountStore {
	s.preds = append(s.preds, s.affiliateFt.ByID(id))
	return s
}

func (s *AccountStore) AffiliatesByOwnerID(id dot.ID) *AccountStore {
	s.preds = append(s.preds, s.affiliateFt.ByOwnerID(id))
	return s
}

func (s *AccountStore) AffiliatesByIDs(ids ...dot.ID) *AccountStore {
	s.preds = append(s.preds, sq.In("id", ids))
	return s
}

func (s *AccountStore) GetShopDB() (*model.Shop, error) {
	var shop model.Shop
	err := s.query().Where(s.preds).ShouldGet(&shop)
	return &shop, err
}

func (s *AccountStore) GetShop() (*identity.Shop, error) {
	shop, err := s.GetShopDB()
	if err != nil {
		return nil, err
	}
	return convert.Shop(shop), nil
}

func (s *AccountStore) GetAffiliateDB() (*identitymodel.Affiliate, error) {
	var affiliate identitymodel.Affiliate
	err := s.query().Where(s.preds).ShouldGet(&affiliate)
	return &affiliate, err
}

func (s *AccountStore) GetAffiliate() (*identity.Affiliate, error) {
	affiliate, err := s.GetAffiliateDB()
	if err != nil {
		return nil, err
	}
	return convert.Affiliate(affiliate), nil
}

func (s *AccountStore) GetAffiliates() ([]*identity.Affiliate, error) {
	var affiliates identitymodel.Affiliates
	err := s.query().Where(s.preds).Find(&affiliates)
	return convert.Affiliates(affiliates), err
}

type CreateAffiliateArgs struct {
	Name        string
	OwnerID     dot.ID
	Phone       string
	Email       string
	IsTest      bool
	BankAccount *identity.BankAccount
}

func (s *AccountStore) CreateAffiliate(args CreateAffiliateArgs) (*identity.Affiliate, error) {
	if args.OwnerID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing OwnerID")
	}

	id := model.NewAffiliateID()
	account := &model.Account{
		ID:       id,
		OwnerID:  args.OwnerID,
		Name:     args.Name,
		Type:     account_type.Affiliate,
		ImageURL: "",
		URLSlug:  "",
	}
	permission := &model.AccountUser{
		AccountID: id,
		UserID:    args.OwnerID,
		Status:    model.StatusActive,
	}
	affiliate := &identitymodel.Affiliate{
		ID:          id,
		OwnerID:     args.OwnerID,
		Name:        args.Name,
		Phone:       args.Phone,
		Email:       args.Email,
		Status:      model.StatusActive,
		BankAccount: convert.BankAccountDB(args.BankAccount),
	}
	if args.IsTest {
		affiliate.IsTest = 1
	}
	if err := s.query().ShouldInsert(account, affiliate, permission); err != nil {
		return nil, err
	}
	return s.AffiliateByID(id).GetAffiliate()
}

type UpdateAffiliateArgs struct {
	ID          dot.ID
	OwnerID     dot.ID
	Phone       string
	Email       string
	Name        string
	BankAccount *identity.BankAccount
}

func (s *AccountStore) UpdateAffiliate(args UpdateAffiliateArgs) (*identity.Affiliate, error) {
	if args.ID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing Affiliate ID")
	}
	update := &identitymodel.Affiliate{
		Name:        args.Name,
		Phone:       args.Phone,
		Email:       args.Email,
		BankAccount: convert.BankAccountDB(args.BankAccount),
	}
	if err := s.query().Where(s.affiliateFt.ByID(args.ID)).
		Where(s.affiliateFt.ByOwnerID(args.OwnerID)).
		ShouldUpdate(update); err != nil {
		return nil, err
	}
	return s.AffiliateByID(args.ID).GetAffiliate()
}

type DeleteAffiliateArgs struct {
	ID      dot.ID
	OwnerID dot.ID
}

func (s *AccountStore) DeleteAffiliate(args DeleteAffiliateArgs) error {
	if args.ID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing ID")
	}
	now := time.Now()
	updateAff := &identitymodel.Affiliate{
		DeletedAt: now,
	}
	if err := s.query().Where(s.affiliateFt.ByID(args.ID)).Where(s.affiliateFt.ByOwnerID(args.OwnerID)).ShouldUpdate(updateAff); err != nil {
		return err
	}
	return nil
}
