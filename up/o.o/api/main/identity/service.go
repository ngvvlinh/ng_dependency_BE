package identity

import (
	"context"
	"time"

	"o.o/api/main/address"
	identitytypes "o.o/api/main/identity/types"
	"o.o/api/meta"
	"o.o/api/top/types/etc/account_type"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/user_source"
	"o.o/capi/dot"
	"o.o/capi/filter"
	"o.o/common/xerrors"
)

// +gen:api

type Aggregate interface {
	// -- User -- //

	UpdateUserReferenceUserID(context.Context, *UpdateUserReferenceUserIDArgs) error

	UpdateUserReferenceSaleID(context.Context, *UpdateUserReferenceSaleIDArgs) error

	UpdateUserEmail(ctx context.Context, userID dot.ID, email string) error

	UpdateUserPhone(ctx context.Context, userID dot.ID, phone string) error

	// if phone is not existed
	// create new user & create a default shop for this user
	RegisterSimplify(context.Context, *RegisterSimplifyArgs) error

	// -- Affiliate -- //

	CreateAffiliate(context.Context, *CreateAffiliateArgs) (*Affiliate, error)

	UpdateAffiliateInfo(context.Context, *UpdateAffiliateInfoArgs) (*Affiliate, error)

	UpdateAffiliateBankAccount(context.Context, *UpdateAffiliateBankAccountArgs) (*Affiliate, error)

	DeleteAffiliate(context.Context, *DeleteAffiliateArgs) error

	// -- Block, Unblock User -- //

	BlockUser(context.Context, *BlockUserArgs) (*User, error)

	UnblockUser(ctx context.Context, userID dot.ID) (*User, error)

	UpdateShipFromAddressID(context.Context, *UpdateShipFromAddressArgs) error

	UpdateUserRef(context.Context, *UpdateUserRefArgs) (*UserRefSaff, error)

	// -- Shop -- //
	CreateShop(context.Context, *CreateShopArgs) (*Shop, error)

	UpdateShopInfo(context.Context, *UpdateShopInfoArgs) error

	// -- Account User -- //
	CreateAccountUser(context.Context, *CreateAccountUserArgs) (*AccountUser, error)

	UpdateAccountUserPermission(context.Context, *UpdateAccountUserPermissionArgs) error
}

type QueryService interface {

	// -- Account -- //
	GetAccountByID(ctx context.Context, ID dot.ID) (*Account, error)

	// -- Shop -- //

	GetShopByID(ctx context.Context, ID dot.ID) (*Shop, error)

	ListShopsByIDs(context.Context, *ListShopsByIDsArgs) ([]*Shop, error)

	ListShopExtendeds(context.Context, *ListShopQuery) (*ListShopExtendedsResponse, error)

	// -- User -- //

	GetUserByID(context.Context, *GetUserByIDQueryArgs) (*User, error)

	GetUserFtRefSaffByID(context.Context, *GetUserByIDQueryArgs) (*UserFtRefSaff, error)

	GetUsersByAccount(ctx context.Context, accountID dot.ID) ([]*AccountUser, error)

	GetUserByPhone(ctx context.Context, phone string) (*User, error)

	GetUserByEmail(ctx context.Context, email string) (*User, error)

	GetUsersByIDs(ctx context.Context, IDs []dot.ID) ([]*User, error)

	GetUserByPhoneOrEmail(context.Context, *GetUserByPhoneOrEmailArgs) (*User, error)

	// -- Affiliate -- //

	GetAffiliateByID(ctx context.Context, ID dot.ID) (*Affiliate, error)

	GetAffiliateWithPermission(ctx context.Context, AffiliateID dot.ID, UserID dot.ID) (*GetAffiliateWithPermissionResult, error)

	GetAffiliatesByIDs(context.Context, *GetAffiliatesByIDsArgs) ([]*Affiliate, error)

	GetAffiliatesByOwnerID(context.Context, *GetAffiliatesByOwnerIDArgs) ([]*Affiliate, error)

	ListPartnersForWhiteLabel(context.Context, *meta.Empty) ([]*Partner, error)

	GetPartnerByID(context.Context, *GetPartnerByIDArgs) (*Partner, error)

	GetUsers(context.Context, *ListUsersArgs) (*UsersResponse, error)

	GetUserFtRefSaffs(context.Context, *ListUserFtRefSaffsArgs) (*UserFtRefSaffsResponse, error)

	GetAllAccountsByUsers(context.Context, *GetAllAccountUsersArg) ([]*AccountUser, error)

	ListUsersByWLPartnerID(context.Context, *ListUsersByWLPartnerID) ([]*User, error)

	ListUsersByIDsAndNameNorm(context.Context, *ListUsersByIDsAndNameNormArgs) ([]*User, error)

	GetAccountUser(ctx context.Context, UserID, AccountID dot.ID) (*AccountUser, error)

	ListPartnerRelationsBySubjectIDs(context.Context, *ListPartnerRelationsBySubjectIDsArgs) ([]*PartnerRelation, error)
}

type Account struct {
	ID       dot.ID
	OwnerID  dot.ID
	Name     string
	Type     account_type.AccountType
	ImageURL string
	URLSlug  string

	Rid dot.ID
}

type BlockUserArgs struct {
	UserID      dot.ID
	BlockBy     dot.ID
	BlockReason string
}

//-- queries --//
type GetUserByIDQueryArgs struct {
	UserID dot.ID
}

type GetAllAccountUsersArg struct {
	UserIDs []dot.ID
	Type    account_type.NullAccountType
	Roles   []string
}

type GetUserByPhoneOrEmailArgs struct {
	Phone string
	Email string
}

type ListUsersArgs struct {
	Name      filter.FullTextSearch
	Phone     string
	Email     string
	CreatedAt filter.Date
	Paging    meta.Paging
}

type ListUserFtRefSaffsArgs struct {
	Name      filter.FullTextSearch
	Phone     string
	Email     string
	RefAff    string
	RefSale   string
	CreatedAt filter.Date
	Paging    meta.Paging
}

type UsersResponse struct {
	ListUsers []*User
	Paging    meta.PageInfo
}

type UserFtRefSaffsResponse struct {
	ListUsers []*UserFtRefSaff
	Paging    meta.PageInfo
}

//-- commands --//

type UpdateUserReferenceUserIDArgs struct {
	UserID       dot.ID
	RefUserPhone string
}

type UpdateUserReferenceSaleIDArgs struct {
	UserID       dot.ID
	RefSalePhone string
}

type GetAffiliateWithPermissionResult struct {
	Affiliate  *Affiliate
	Permission Permission
}

type GetAffiliatesByIDsArgs struct {
	AffiliateIDs []dot.ID
}

type GetAffiliatesByOwnerIDArgs struct {
	ID dot.ID
}

type CreateAffiliateArgs struct {
	Name        string
	OwnerID     dot.ID
	Phone       string
	Email       string
	IsTest      bool
	BankAccount *identitytypes.BankAccount
}

type UpdateAffiliateInfoArgs struct {
	ID      dot.ID
	OwnerID dot.ID
	Phone   string
	Email   string
	Name    string
}

type DeleteAffiliateArgs struct {
	ID      dot.ID
	OwnerID dot.ID
}

type UpdateAffiliateBankAccountArgs struct {
	ID          dot.ID
	OwnerID     dot.ID
	BankAccount *identitytypes.BankAccount
}

type GetPartnerByIDArgs struct {
	ID dot.ID
}

type ListUsersByWLPartnerID struct {
	ID dot.ID
}

type ListShopQuery struct {
	Paging               meta.Paging
	Filters              meta.Filters
	Name                 filter.FullTextSearch
	ShopIDs              []dot.ID
	OwnerID              dot.ID
	IncludeWLPartnerShop bool
}

type ListShopExtendedsResponse struct {
	Shops  []*ShopExtended
	Paging meta.PageInfo
}

type UpdateUserRefArgs struct {
	UserID  dot.ID
	RefAff  dot.NullString
	RefSale dot.NullString
}

type UpdateShipFromAddressArgs struct {
	ID                dot.ID
	ShipFromAddressID dot.ID
}

type UpdateDefaultAddressArgs struct {
	ID dot.ID
}

type ListShopsByIDsArgs struct {
	IDs                     []dot.ID
	IsPriorMoneyTransaction bool
	IncludeWLPartnerShop    bool
}

type CreateShopArgs struct {
	ID                          dot.ID
	Name                        string
	OwnerID                     dot.ID
	AddressID                   dot.ID
	Address                     *address.Address
	Phone                       string
	BankAccount                 *identitytypes.BankAccount
	WebsiteURL                  dot.NullString
	ImageURL                    string
	Email                       string
	AutoCreateFFM               bool
	IsTest                      bool
	URLSlug                     string
	CompanyInfo                 *identitytypes.CompanyInfo
	MoneyTransactionRRule       string
	SurveyInfo                  []*SurveyInfo
	ShippingServicePickStrategy []*ShippingServiceSelectStrategyItem
}

type CreateUserArgs struct {
	UserID                  dot.ID
	FullName                string
	ShortName               string
	Email                   string
	Phone                   string
	Password                string
	Status                  status3.Status
	AgreeTOS                bool
	AgreeEmailInfo          bool
	IsTest                  bool
	Source                  user_source.UserSource
	RefSale                 string
	RefAff                  string
	PhoneVerifiedAt         time.Time
	PhoneVerificationSentAt time.Time
}

type ListUsersByIDsAndNameNormArgs struct {
	IDs      []dot.ID
	NameNorm filter.FullTextSearch
}

type UpdateShopInfoArgs struct {
	ShopID                  dot.ID
	MoneyTransactionRrule   string
	IsPriorMoneyTransaction dot.NullBool
}

type ListPartnerRelationsBySubjectIDsArgs struct {
	SubjectIDs  []dot.ID
	SubjectType SubjectType
}

type RegisterSimplifyArgs struct {
	Phone               string
	Password            string
	FullName            string
	IsCreateDefaultShop bool
}

// +convert:create=AccountUser
type CreateAccountUserArgs struct {
	AccountID dot.ID
	UserID    dot.ID

	Status status3.Status // 1: activated, -1: rejected/disabled, 0: pending
	Permission

	FullName  string
	ShortName string
	Position  string

	InvitationSentAt     time.Time
	InvitationSentBy     dot.ID
	InvitationAcceptedAt time.Time
}

func (args *CreateAccountUserArgs) Validate() error {
	if args.UserID == 0 {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing user ID")
	}
	if args.AccountID == 0 {
		return xerrors.Errorf(xerrors.InvalidArgument, nil, "Missing account ID")
	}
	if len(args.Roles) == 0 {
		return xerrors.Errorf(xerrors.FailedPrecondition, nil, "Require at least 1 role")
	}
	return nil
}

type UpdateAccountUserPermissionArgs struct {
	AccountID dot.ID
	UserID    dot.ID
	Permission
}
