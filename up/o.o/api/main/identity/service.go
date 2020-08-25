package identity

import (
	"context"

	identitytypes "o.o/api/main/identity/types"
	"o.o/api/meta"
	"o.o/api/top/types/etc/account_type"
	"o.o/capi/dot"
	"o.o/capi/filter"
)

// +gen:api

type Aggregate interface {
	// -- User -- //

	UpdateUserReferenceUserID(context.Context, *UpdateUserReferenceUserIDArgs) error

	UpdateUserReferenceSaleID(context.Context, *UpdateUserReferenceSaleIDArgs) error

	UpdateUserEmail(ctx context.Context, userID dot.ID, email string) error

	UpdateUserPhone(ctx context.Context, userID dot.ID, phone string) error

	// -- External Account Ahamove -- //
	// TODO: move External Account Ahamove to its own module

	CreateExternalAccountAhamove(context.Context, *CreateExternalAccountAhamoveArgs) (*ExternalAccountAhamove, error)

	RequestVerifyExternalAccountAhamove(context.Context, *RequestVerifyExternalAccountAhamoveArgs) (*RequestVerifyExternalAccountAhamoveResult, error)

	UpdateVerifiedExternalAccountAhamove(context.Context, *UpdateVerifiedExternalAccountAhamoveArgs) (*ExternalAccountAhamove, error)

	UpdateExternalAccountAhamoveVerification(context.Context, *UpdateExternalAccountAhamoveVerificationArgs) (*ExternalAccountAhamove, error)

	// -- Affiliate -- //

	CreateAffiliate(context.Context, *CreateAffiliateArgs) (*Affiliate, error)

	UpdateAffiliateInfo(context.Context, *UpdateAffiliateInfoArgs) (*Affiliate, error)

	UpdateAffiliateBankAccount(context.Context, *UpdateAffiliateBankAccountArgs) (*Affiliate, error)

	DeleteAffiliate(context.Context, *DeleteAffiliateArgs) error

	// -- Block, Unblock User -- //

	BlockUser(context.Context, *BlockUserArgs) (*User, error)

	UnblockUser(ctx context.Context, userID dot.ID) (*User, error)

	UpdateUserRef(context.Context, *UpdateUserRefArgs) (*UserRefSaff, error)

	UpdateShipFromAddressID(context.Context, *UpdateShipFromAddressArgs) error
}

type QueryService interface {

	// -- Account -- //
	GetAccountByID(ctx context.Context, ID dot.ID) (*Account, error)

	// -- Shop -- //

	GetShopByID(ctx context.Context, ID dot.ID) (*Shop, error)

	ListShopsByIDs(ctx context.Context, IDs []dot.ID) ([]*Shop, error)

	ListShopExtendeds(context.Context, *ListShopQuery) (*ListShopExtendedsResponse, error)

	// -- User -- //

	GetUserByID(context.Context, *GetUserByIDQueryArgs) (*User, error)

	GetUserFtRefSaffByID(context.Context, *GetUserByIDQueryArgs) (*UserFtRefSaff, error)

	GetUsersByAccount(ctx context.Context, accountID dot.ID) ([]*AccountUser, error)

	GetUserByPhone(ctx context.Context, phone string) (*User, error)

	GetUserByEmail(ctx context.Context, email string) (*User, error)

	GetUsersByIDs(ctx context.Context, IDs []dot.ID) ([]*User, error)

	GetUserByPhoneOrEmail(context.Context, *GetUserByPhoneOrEmailArgs) (*User, error)

	// -- ExternalAccountAhamove -- //

	GetExternalAccountAhamove(context.Context, *GetExternalAccountAhamoveArgs) (*ExternalAccountAhamove, error)

	GetExternalAccountAhamoveByExternalID(context.Context, *GetExternalAccountAhamoveByExternalIDQueryArgs) (*ExternalAccountAhamove, error)

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
}

type GetUserByPhoneOrEmailArgs struct {
	Phone string
	Email string
}

type ListUsersArgs struct {
	Name      string
	Phone     string
	Email     string
	CreatedAt filter.Date
	Paging    meta.Paging
}

type ListUserFtRefSaffsArgs struct {
	Name      string
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

type GetExternalAccountAhamoveArgs struct {
	OwnerID dot.ID
	Phone   string
}

type GetExternalAccountAhamoveByExternalIDQueryArgs struct {
	ExternalID string
}

//-- commands --//
type CreateExternalAccountAhamoveArgs struct {
	OwnerID dot.ID // user id
	Phone   string
	Name    string
	Address string
}

type RequestVerifyExternalAccountAhamoveArgs struct {
	OwnerID dot.ID
	Phone   string
}

type RequestVerifyExternalAccountAhamoveResult struct {
}

type UpdateVerifiedExternalAccountAhamoveArgs struct {
	OwnerID dot.ID
	Phone   string
}

type UpdateExternalAccountAhamoveVerificationArgs struct {
	OwnerID             dot.ID
	Phone               string
	IDCardFrontImg      string
	IDCardBackImg       string
	PortraitImg         string
	WebsiteURL          string
	FanpageURL          string
	CompanyImgs         []string
	BusinessLicenseImgs []string
}

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
	Paging  meta.Paging
	Filters meta.Filters
	Name    filter.FullTextSearch
	ShopIDs []dot.ID
}

type ListShopExtendedsResponse struct {
	Shops  []*ShopExtended
	Paging meta.PageInfo
}

type UpdateUserRefArgs struct {
	UserID  dot.ID
	RefAff  string
	RefSale string
}

type UpdateShipFromAddressArgs struct {
	ID                dot.ID
	ShipFromAddressID dot.ID
}

type UpdateDefaultAddressArgs struct {
	ID dot.ID
}
