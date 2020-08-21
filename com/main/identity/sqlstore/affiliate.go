package sqlstore

import (
	"context"
	"time"

	"o.o/api/main/identity"
	identitytypes "o.o/api/main/identity/types"
	"o.o/api/meta"
	"o.o/api/top/types/etc/account_type"
	status3 "o.o/api/top/types/etc/status3"
	"o.o/backend/com/main/identity/convert"
	identitymodel "o.o/backend/com/main/identity/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/backend/pkg/etc/idutil"
	"o.o/backend/pkg/etop/model"
	"o.o/capi/dot"
)

type AffiliateStoreFactory func(context.Context) *AffiliateStore

func NewAffiliateStore(db *cmsql.Database) AffiliateStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *AffiliateStore {
		return &AffiliateStore{
			query: cmsql.NewQueryFactory(ctx, db),
			ctx:   ctx,
		}
	}
}

type AffiliateStore struct {
	query       cmsql.QueryFactory
	affiliateFt AffiliateFilters
	preds       []interface{}
	sqlstore.Paging
	filter         meta.Filters
	ctx            context.Context
	includeDeleted sqlstore.IncludeDeleted
}

func (s *AffiliateStore) ByID(id dot.ID) *AffiliateStore {
	s.preds = append(s.preds, s.affiliateFt.ByID(id))
	return s
}

func (s *AffiliateStore) ByOwnerID(id dot.ID) *AffiliateStore {
	s.preds = append(s.preds, s.affiliateFt.ByOwnerID(id))
	return s
}

func (s *AffiliateStore) ByIDs(ids ...dot.ID) *AffiliateStore {
	s.preds = append(s.preds, sq.In("id", ids))
	return s
}

func (s *AffiliateStore) GetAffiliateDB() (*identitymodel.Affiliate, error) {
	var affiliate identitymodel.Affiliate
	err := s.query().Where(s.preds).ShouldGet(&affiliate)
	return &affiliate, err
}

func (s *AffiliateStore) GetAffiliate() (*identity.Affiliate, error) {
	affiliate, err := s.GetAffiliateDB()
	if err != nil {
		return nil, err
	}
	return convert.Affiliate(affiliate), nil
}

func (s *AffiliateStore) ListAffiliates() ([]*identity.Affiliate, error) {
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

func (s *AffiliateStore) CreateAffiliate(args CreateAffiliateArgs) (*identity.Affiliate, error) {
	if args.OwnerID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing OwnerID")
	}

	id := idutil.NewAffiliateID()
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
	return s.ByID(id).GetAffiliate()
}

type UpdateAffiliateArgs struct {
	ID          dot.ID
	OwnerID     dot.ID
	Phone       string
	Email       string
	Name        string
	BankAccount *identitytypes.BankAccount
}

func (s *AffiliateStore) UpdateAffiliate(args UpdateAffiliateArgs) (*identity.Affiliate, error) {
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
	return s.ByID(args.ID).GetAffiliate()
}

type DeleteAffiliateArgs struct {
	ID      dot.ID
	OwnerID dot.ID
}

func (s *AffiliateStore) DeleteAffiliate(args DeleteAffiliateArgs) error {
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
