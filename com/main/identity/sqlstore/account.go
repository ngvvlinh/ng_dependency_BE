package sqlstore

import (
	"context"
	"time"

	"o.o/api/main/identity"
	identitytypes "o.o/api/main/identity/types"
	"o.o/api/meta"
	"o.o/api/top/types/etc/account_type"
	"o.o/api/top/types/etc/status3"
	"o.o/backend/com/main/identity/convert"
	identitymodel "o.o/backend/com/main/identity/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/backend/pkg/etop/model"
	"o.o/capi/dot"
)

type AccountStoreFactory func(context.Context) *AccountStore

func NewAccountStore(db *cmsql.Database) AccountStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *AccountStore {
		return &AccountStore{
			query: cmsql.NewQueryFactory(ctx, db),
			ctx:   ctx,
		}
	}
}

type AccountStore struct {
	query       cmsql.QueryFactory
	preds       []interface{}
	shopFt      ShopFilters
	affiliateFt AffiliateFilters
	accountFt   AccountFilters
	sqlstore.Paging
	filter         meta.Filters
	ctx            context.Context
	includeDeleted sqlstore.IncludeDeleted
}

func (s *AccountStore) extend() *AccountStore {
	s.shopFt.prefix = "s"
	return s
}

func (s *AccountStore) WithPaging(paging meta.Paging) *AccountStore {
	s.Paging.WithPaging(paging)
	return s
}

func (s *AccountStore) Filters(filters meta.Filters) *AccountStore {
	if s.filter == nil {
		s.filter = filters
	} else {
		s.filter = append(s.filter, filters...)
	}
	return s
}

func (s *AccountStore) ShopByID(id dot.ID) *AccountStore {
	s.preds = append(s.preds, s.shopFt.ByID(id))
	return s
}

func (s *AccountStore) ByAccountIds(ids ...dot.ID) *AccountStore {
	s.preds = append(s.preds, sq.In("id", ids))
	return s
}

func (s *AccountStore) ShopByIDs(ids ...dot.ID) *AccountStore {
	s.preds = append(s.preds, sq.In("id", ids))
	return s
}

func (s *AccountStore) AffiliateByID(id dot.ID) *AccountStore {
	s.preds = append(s.preds, s.affiliateFt.ByID(id))
	return s
}

func (s *AccountStore) ByType(ty account_type.AccountType) *AccountStore {
	s.preds = append(s.preds, s.accountFt.ByType(ty))
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

func (s *AccountStore) GetShopDB() (*identitymodel.Shop, error) {
	var shop identitymodel.Shop
	query := s.query().Where(s.preds)

	// FIX(Tuan): comment vụ check wlPartnerID
	//
	// Webhook NVC cần biết đơn thuộc wlPartnerID nào, hiện tại chưa lưu thông tin này trong order/ffm
	// Tạm thời gọi api GetShopByID từ ffm.ShopID ra để lấy wlPartnerID

	// query = s.FilterByWhiteLabelPartner(query, wl.GetWLPartnerID(s.ctx))
	err := query.ShouldGet(&shop)
	return &shop, err
}

func (s *AccountStore) GetShop() (*identity.Shop, error) {
	shop, err := s.GetShopDB()
	if err != nil {
		return nil, err
	}
	return convert.Shop(shop), nil
}

func (s *AccountStore) ListShopDBs() (res []*identitymodel.Shop, err error) {
	if len(s.Paging.Sort) == 0 {
		s.Paging.Sort = append(s.Paging.Sort, "-created_at")
	}
	query := s.query().Where(s.preds)
	query = s.FilterByWhiteLabelPartner(query, wl.GetWLPartnerID(s.ctx))
	query, err = sqlstore.LimitSort(query, &s.Paging, map[string]string{"created_at": "created_at"})
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filter, filterShopExtendedWhitelist)

	err = s.query().Where(s.preds).Find((*identitymodel.Shops)(&res))
	return
}

func (s *AccountStore) ListShops() ([]*identity.Shop, error) {
	shops, err := s.ListShopDBs()
	if err != nil {
		return nil, err
	}
	var res []*identity.Shop
	if err := scheme.Convert(shops, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (s *AccountStore) ListShopExtendedDBs() (res []*identitymodel.ShopExtended, err error) {
	if len(s.Paging.Sort) == 0 {
		s.Paging.Sort = append(s.Paging.Sort, "-created_at")
	}
	query := s.extend().query().Where(s.preds)
	query = s.FilterByWhiteLabelPartner(query, wl.GetWLPartnerID(s.ctx))
	query = s.includeDeleted.Check(query, s.shopFt.NotDeleted())
	query, err = sqlstore.LimitSort(query, &s.Paging, map[string]string{"created_at": "created_at"}, s.shopFt.prefix)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filter, filterShopExtendedWhitelist)

	err = query.Find((*identitymodel.ShopExtendeds)(&res))
	return
}

func (s *AccountStore) ListShopExtendeds() ([]*identity.ShopExtended, error) {
	shops, err := s.ListShopExtendedDBs()
	if err != nil {
		return nil, err
	}
	var res []*identity.ShopExtended
	if err := scheme.Convert(shops, &res); err != nil {
		return nil, err
	}
	return res, nil
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
	BankAccount *identitytypes.BankAccount
}

func (s *AccountStore) CreateAffiliate(args CreateAffiliateArgs) (*identity.Affiliate, error) {
	if args.OwnerID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing OwnerID")
	}

	id := model.NewAffiliateID()
	account := &identitymodel.Account{
		ID:       id,
		OwnerID:  args.OwnerID,
		Name:     args.Name,
		Type:     account_type.Affiliate,
		ImageURL: "",
		URLSlug:  "",
	}
	permission := &identitymodel.AccountUser{
		AccountID: id,
		UserID:    args.OwnerID,
		Status:    status3.P,
	}
	affiliate := &identitymodel.Affiliate{
		ID:          id,
		OwnerID:     args.OwnerID,
		Name:        args.Name,
		Phone:       args.Phone,
		Email:       args.Email,
		Status:      status3.P,
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
	BankAccount *identitytypes.BankAccount
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

func (s *AccountStore) FilterByWhiteLabelPartner(query cmsql.Query, wlPartnerID dot.ID) cmsql.Query {
	if wlPartnerID != 0 {
		return query.Where(s.shopFt.ByWLPartnerID(wlPartnerID))
	}
	return query.Where(s.shopFt.NotBelongWLPartner())
}

func (s *AccountStore) ListAccountDB() ([]*identitymodel.Account, error) {
	var accounts identitymodel.Accounts
	err := s.query().Where(s.preds).Find(&accounts)
	return accounts, err
}
