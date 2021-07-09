// +build !generator

// Code generated by generator api. DO NOT EDIT.

package identity

import (
	context "context"
	time "time"

	address "o.o/api/main/address"
	identitytypes "o.o/api/main/identity/types"
	meta "o.o/api/meta"
	account_type "o.o/api/top/types/etc/account_type"
	status3 "o.o/api/top/types/etc/status3"
	capi "o.o/capi"
	dot "o.o/capi/dot"
	filter "o.o/capi/filter"
)

type CommandBus struct{ bus capi.Bus }
type QueryBus struct{ bus capi.Bus }

func NewCommandBus(bus capi.Bus) CommandBus { return CommandBus{bus} }
func NewQueryBus(bus capi.Bus) QueryBus     { return QueryBus{bus} }

func (b CommandBus) Dispatch(ctx context.Context, msg interface{ command() }) error {
	return b.bus.Dispatch(ctx, msg)
}
func (b QueryBus) Dispatch(ctx context.Context, msg interface{ query() }) error {
	return b.bus.Dispatch(ctx, msg)
}

type BlockUserCommand struct {
	UserID      dot.ID
	BlockBy     dot.ID
	BlockReason string

	Result *User `json:"-"`
}

func (h AggregateHandler) HandleBlockUser(ctx context.Context, msg *BlockUserCommand) (err error) {
	msg.Result, err = h.inner.BlockUser(msg.GetArgs(ctx))
	return err
}

type CreateAccountUserCommand struct {
	AccountID            dot.ID
	UserID               dot.ID
	Status               status3.Status
	Permission           Permission
	FullName             string
	ShortName            string
	Position             string
	InvitationSentAt     time.Time
	InvitationSentBy     dot.ID
	InvitationAcceptedAt time.Time

	Result *AccountUser `json:"-"`
}

func (h AggregateHandler) HandleCreateAccountUser(ctx context.Context, msg *CreateAccountUserCommand) (err error) {
	msg.Result, err = h.inner.CreateAccountUser(msg.GetArgs(ctx))
	return err
}

type CreateAffiliateCommand struct {
	Name        string
	OwnerID     dot.ID
	Phone       string
	Email       string
	IsTest      bool
	BankAccount *identitytypes.BankAccount

	Result *Affiliate `json:"-"`
}

func (h AggregateHandler) HandleCreateAffiliate(ctx context.Context, msg *CreateAffiliateCommand) (err error) {
	msg.Result, err = h.inner.CreateAffiliate(msg.GetArgs(ctx))
	return err
}

type CreateShopCommand struct {
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

	Result *Shop `json:"-"`
}

func (h AggregateHandler) HandleCreateShop(ctx context.Context, msg *CreateShopCommand) (err error) {
	msg.Result, err = h.inner.CreateShop(msg.GetArgs(ctx))
	return err
}

type DeleteAccountCommand struct {
	AccountID dot.ID
	OwnerID   dot.ID

	Result struct {
	} `json:"-"`
}

func (h AggregateHandler) HandleDeleteAccount(ctx context.Context, msg *DeleteAccountCommand) (err error) {
	return h.inner.DeleteAccount(msg.GetArgs(ctx))
}

type DeleteAccountUsersCommand struct {
	AccountID dot.ID
	UserID    dot.ID

	Result int `json:"-"`
}

func (h AggregateHandler) HandleDeleteAccountUsers(ctx context.Context, msg *DeleteAccountUsersCommand) (err error) {
	msg.Result, err = h.inner.DeleteAccountUsers(msg.GetArgs(ctx))
	return err
}

type DeleteAffiliateCommand struct {
	ID      dot.ID
	OwnerID dot.ID

	Result struct {
	} `json:"-"`
}

func (h AggregateHandler) HandleDeleteAffiliate(ctx context.Context, msg *DeleteAffiliateCommand) (err error) {
	return h.inner.DeleteAffiliate(msg.GetArgs(ctx))
}

type RegisterSimplifyCommand struct {
	Phone               string
	Password            string
	FullName            string
	Email               string
	CompanyName         string
	IsCreateDefaultShop bool

	Result struct {
	} `json:"-"`
}

func (h AggregateHandler) HandleRegisterSimplify(ctx context.Context, msg *RegisterSimplifyCommand) (err error) {
	return h.inner.RegisterSimplify(msg.GetArgs(ctx))
}

type UnblockUserCommand struct {
	UserID dot.ID

	Result *User `json:"-"`
}

func (h AggregateHandler) HandleUnblockUser(ctx context.Context, msg *UnblockUserCommand) (err error) {
	msg.Result, err = h.inner.UnblockUser(msg.GetArgs(ctx))
	return err
}

type UpdateAccountUserPermissionCommand struct {
	AccountID  dot.ID
	UserID     dot.ID
	Permission Permission

	Result struct {
	} `json:"-"`
}

func (h AggregateHandler) HandleUpdateAccountUserPermission(ctx context.Context, msg *UpdateAccountUserPermissionCommand) (err error) {
	return h.inner.UpdateAccountUserPermission(msg.GetArgs(ctx))
}

type UpdateAffiliateBankAccountCommand struct {
	ID          dot.ID
	OwnerID     dot.ID
	BankAccount *identitytypes.BankAccount

	Result *Affiliate `json:"-"`
}

func (h AggregateHandler) HandleUpdateAffiliateBankAccount(ctx context.Context, msg *UpdateAffiliateBankAccountCommand) (err error) {
	msg.Result, err = h.inner.UpdateAffiliateBankAccount(msg.GetArgs(ctx))
	return err
}

type UpdateAffiliateInfoCommand struct {
	ID      dot.ID
	OwnerID dot.ID
	Phone   string
	Email   string
	Name    string

	Result *Affiliate `json:"-"`
}

func (h AggregateHandler) HandleUpdateAffiliateInfo(ctx context.Context, msg *UpdateAffiliateInfoCommand) (err error) {
	msg.Result, err = h.inner.UpdateAffiliateInfo(msg.GetArgs(ctx))
	return err
}

type UpdateExtensionNumberNormCommand struct {
	AccountID       dot.ID
	UserID          dot.ID
	ExtensionNumber string

	Result struct {
	} `json:"-"`
}

func (h AggregateHandler) HandleUpdateExtensionNumberNorm(ctx context.Context, msg *UpdateExtensionNumberNormCommand) (err error) {
	return h.inner.UpdateExtensionNumberNorm(msg.GetArgs(ctx))
}

type UpdateShipFromAddressIDCommand struct {
	ID                dot.ID
	ShipFromAddressID dot.ID

	Result struct {
	} `json:"-"`
}

func (h AggregateHandler) HandleUpdateShipFromAddressID(ctx context.Context, msg *UpdateShipFromAddressIDCommand) (err error) {
	return h.inner.UpdateShipFromAddressID(msg.GetArgs(ctx))
}

type UpdateShopInfoCommand struct {
	ShopID                  dot.ID
	MoneyTransactionRrule   string
	IsPriorMoneyTransaction dot.NullBool

	Result struct {
	} `json:"-"`
}

func (h AggregateHandler) HandleUpdateShopInfo(ctx context.Context, msg *UpdateShopInfoCommand) (err error) {
	return h.inner.UpdateShopInfo(msg.GetArgs(ctx))
}

type UpdateUserEmailCommand struct {
	UserID dot.ID
	Email  string

	Result struct {
	} `json:"-"`
}

func (h AggregateHandler) HandleUpdateUserEmail(ctx context.Context, msg *UpdateUserEmailCommand) (err error) {
	return h.inner.UpdateUserEmail(msg.GetArgs(ctx))
}

type UpdateUserPhoneCommand struct {
	UserID dot.ID
	Phone  string

	Result struct {
	} `json:"-"`
}

func (h AggregateHandler) HandleUpdateUserPhone(ctx context.Context, msg *UpdateUserPhoneCommand) (err error) {
	return h.inner.UpdateUserPhone(msg.GetArgs(ctx))
}

type UpdateUserRefCommand struct {
	UserID  dot.ID
	RefAff  dot.NullString
	RefSale dot.NullString

	Result *UserRefSaff `json:"-"`
}

func (h AggregateHandler) HandleUpdateUserRef(ctx context.Context, msg *UpdateUserRefCommand) (err error) {
	msg.Result, err = h.inner.UpdateUserRef(msg.GetArgs(ctx))
	return err
}

type UpdateUserReferenceSaleIDCommand struct {
	UserID       dot.ID
	RefSalePhone string

	Result struct {
	} `json:"-"`
}

func (h AggregateHandler) HandleUpdateUserReferenceSaleID(ctx context.Context, msg *UpdateUserReferenceSaleIDCommand) (err error) {
	return h.inner.UpdateUserReferenceSaleID(msg.GetArgs(ctx))
}

type UpdateUserReferenceUserIDCommand struct {
	UserID       dot.ID
	RefUserPhone string

	Result struct {
	} `json:"-"`
}

func (h AggregateHandler) HandleUpdateUserReferenceUserID(ctx context.Context, msg *UpdateUserReferenceUserIDCommand) (err error) {
	return h.inner.UpdateUserReferenceUserID(msg.GetArgs(ctx))
}

type GetAccountByIDQuery struct {
	ID dot.ID

	Result *Account `json:"-"`
}

func (h QueryServiceHandler) HandleGetAccountByID(ctx context.Context, msg *GetAccountByIDQuery) (err error) {
	msg.Result, err = h.inner.GetAccountByID(msg.GetArgs(ctx))
	return err
}

type GetAccountUserQuery struct {
	UserID    dot.ID
	AccountID dot.ID

	Result *AccountUser `json:"-"`
}

func (h QueryServiceHandler) HandleGetAccountUser(ctx context.Context, msg *GetAccountUserQuery) (err error) {
	msg.Result, err = h.inner.GetAccountUser(msg.GetArgs(ctx))
	return err
}

type GetAffiliateByIDQuery struct {
	ID dot.ID

	Result *Affiliate `json:"-"`
}

func (h QueryServiceHandler) HandleGetAffiliateByID(ctx context.Context, msg *GetAffiliateByIDQuery) (err error) {
	msg.Result, err = h.inner.GetAffiliateByID(msg.GetArgs(ctx))
	return err
}

type GetAffiliateWithPermissionQuery struct {
	AffiliateID dot.ID
	UserID      dot.ID

	Result *GetAffiliateWithPermissionResult `json:"-"`
}

func (h QueryServiceHandler) HandleGetAffiliateWithPermission(ctx context.Context, msg *GetAffiliateWithPermissionQuery) (err error) {
	msg.Result, err = h.inner.GetAffiliateWithPermission(msg.GetArgs(ctx))
	return err
}

type GetAffiliatesByIDsQuery struct {
	AffiliateIDs []dot.ID

	Result []*Affiliate `json:"-"`
}

func (h QueryServiceHandler) HandleGetAffiliatesByIDs(ctx context.Context, msg *GetAffiliatesByIDsQuery) (err error) {
	msg.Result, err = h.inner.GetAffiliatesByIDs(msg.GetArgs(ctx))
	return err
}

type GetAffiliatesByOwnerIDQuery struct {
	ID dot.ID

	Result []*Affiliate `json:"-"`
}

func (h QueryServiceHandler) HandleGetAffiliatesByOwnerID(ctx context.Context, msg *GetAffiliatesByOwnerIDQuery) (err error) {
	msg.Result, err = h.inner.GetAffiliatesByOwnerID(msg.GetArgs(ctx))
	return err
}

type GetAllAccountsByUsersQuery struct {
	UserIDs []dot.ID
	Type    account_type.NullAccountType
	Roles   []string

	Result []*AccountUser `json:"-"`
}

func (h QueryServiceHandler) HandleGetAllAccountsByUsers(ctx context.Context, msg *GetAllAccountsByUsersQuery) (err error) {
	msg.Result, err = h.inner.GetAllAccountsByUsers(msg.GetArgs(ctx))
	return err
}

type GetPartnerByIDQuery struct {
	ID dot.ID

	Result *Partner `json:"-"`
}

func (h QueryServiceHandler) HandleGetPartnerByID(ctx context.Context, msg *GetPartnerByIDQuery) (err error) {
	msg.Result, err = h.inner.GetPartnerByID(msg.GetArgs(ctx))
	return err
}

type GetShopByIDQuery struct {
	ID dot.ID

	Result *Shop `json:"-"`
}

func (h QueryServiceHandler) HandleGetShopByID(ctx context.Context, msg *GetShopByIDQuery) (err error) {
	msg.Result, err = h.inner.GetShopByID(msg.GetArgs(ctx))
	return err
}

type GetUserByEmailQuery struct {
	Email string

	Result *User `json:"-"`
}

func (h QueryServiceHandler) HandleGetUserByEmail(ctx context.Context, msg *GetUserByEmailQuery) (err error) {
	msg.Result, err = h.inner.GetUserByEmail(msg.GetArgs(ctx))
	return err
}

type GetUserByIDQuery struct {
	UserID dot.ID

	Result *User `json:"-"`
}

func (h QueryServiceHandler) HandleGetUserByID(ctx context.Context, msg *GetUserByIDQuery) (err error) {
	msg.Result, err = h.inner.GetUserByID(msg.GetArgs(ctx))
	return err
}

type GetUserByPhoneQuery struct {
	Phone string

	Result *User `json:"-"`
}

func (h QueryServiceHandler) HandleGetUserByPhone(ctx context.Context, msg *GetUserByPhoneQuery) (err error) {
	msg.Result, err = h.inner.GetUserByPhone(msg.GetArgs(ctx))
	return err
}

type GetUserByPhoneOrEmailQuery struct {
	Phone string
	Email string

	Result *User `json:"-"`
}

func (h QueryServiceHandler) HandleGetUserByPhoneOrEmail(ctx context.Context, msg *GetUserByPhoneOrEmailQuery) (err error) {
	msg.Result, err = h.inner.GetUserByPhoneOrEmail(msg.GetArgs(ctx))
	return err
}

type GetUserFtRefSaffByIDQuery struct {
	UserID dot.ID

	Result *UserFtRefSaff `json:"-"`
}

func (h QueryServiceHandler) HandleGetUserFtRefSaffByID(ctx context.Context, msg *GetUserFtRefSaffByIDQuery) (err error) {
	msg.Result, err = h.inner.GetUserFtRefSaffByID(msg.GetArgs(ctx))
	return err
}

type GetUserFtRefSaffsQuery struct {
	Name      filter.FullTextSearch
	Phone     string
	Email     string
	RefAff    string
	RefSale   string
	CreatedAt filter.Date
	Paging    meta.Paging

	Result *UserFtRefSaffsResponse `json:"-"`
}

func (h QueryServiceHandler) HandleGetUserFtRefSaffs(ctx context.Context, msg *GetUserFtRefSaffsQuery) (err error) {
	msg.Result, err = h.inner.GetUserFtRefSaffs(msg.GetArgs(ctx))
	return err
}

type GetUsersQuery struct {
	Name      filter.FullTextSearch
	Phone     string
	Email     string
	CreatedAt filter.Date
	Paging    meta.Paging

	Result *UsersResponse `json:"-"`
}

func (h QueryServiceHandler) HandleGetUsers(ctx context.Context, msg *GetUsersQuery) (err error) {
	msg.Result, err = h.inner.GetUsers(msg.GetArgs(ctx))
	return err
}

type GetUsersByAccountQuery struct {
	AccountID dot.ID

	Result []*AccountUser `json:"-"`
}

func (h QueryServiceHandler) HandleGetUsersByAccount(ctx context.Context, msg *GetUsersByAccountQuery) (err error) {
	msg.Result, err = h.inner.GetUsersByAccount(msg.GetArgs(ctx))
	return err
}

type GetUsersByIDsQuery struct {
	IDs []dot.ID

	Result []*User `json:"-"`
}

func (h QueryServiceHandler) HandleGetUsersByIDs(ctx context.Context, msg *GetUsersByIDsQuery) (err error) {
	msg.Result, err = h.inner.GetUsersByIDs(msg.GetArgs(ctx))
	return err
}

type ListAccountUsersQuery struct {
	AccountID dot.ID
	UserID    dot.ID

	Result []*AccountUser `json:"-"`
}

func (h QueryServiceHandler) HandleListAccountUsers(ctx context.Context, msg *ListAccountUsersQuery) (err error) {
	msg.Result, err = h.inner.ListAccountUsers(msg.GetArgs(ctx))
	return err
}

type ListPartnerRelationsBySubjectIDsQuery struct {
	SubjectIDs  []dot.ID
	SubjectType SubjectType

	Result []*PartnerRelation `json:"-"`
}

func (h QueryServiceHandler) HandleListPartnerRelationsBySubjectIDs(ctx context.Context, msg *ListPartnerRelationsBySubjectIDsQuery) (err error) {
	msg.Result, err = h.inner.ListPartnerRelationsBySubjectIDs(msg.GetArgs(ctx))
	return err
}

type ListPartnersForWhiteLabelQuery struct {
	Result []*Partner `json:"-"`
}

func (h QueryServiceHandler) HandleListPartnersForWhiteLabel(ctx context.Context, msg *ListPartnersForWhiteLabelQuery) (err error) {
	msg.Result, err = h.inner.ListPartnersForWhiteLabel(msg.GetArgs(ctx))
	return err
}

type ListShopExtendedsQuery struct {
	Paging               meta.Paging
	Filters              meta.Filters
	Name                 filter.FullTextSearch
	ShopIDs              []dot.ID
	OwnerID              dot.ID
	IncludeWLPartnerShop bool

	Result *ListShopExtendedsResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListShopExtendeds(ctx context.Context, msg *ListShopExtendedsQuery) (err error) {
	msg.Result, err = h.inner.ListShopExtendeds(msg.GetArgs(ctx))
	return err
}

type ListShopsByIDsQuery struct {
	IDs                     []dot.ID
	IsPriorMoneyTransaction bool
	IncludeWLPartnerShop    bool

	Result []*Shop `json:"-"`
}

func (h QueryServiceHandler) HandleListShopsByIDs(ctx context.Context, msg *ListShopsByIDsQuery) (err error) {
	msg.Result, err = h.inner.ListShopsByIDs(msg.GetArgs(ctx))
	return err
}

type ListUsersByIDsAndNameNormQuery struct {
	IDs      []dot.ID
	NameNorm filter.FullTextSearch

	Result []*User `json:"-"`
}

func (h QueryServiceHandler) HandleListUsersByIDsAndNameNorm(ctx context.Context, msg *ListUsersByIDsAndNameNormQuery) (err error) {
	msg.Result, err = h.inner.ListUsersByIDsAndNameNorm(msg.GetArgs(ctx))
	return err
}

type ListUsersByWLPartnerIDQuery struct {
	ID dot.ID

	Result []*User `json:"-"`
}

func (h QueryServiceHandler) HandleListUsersByWLPartnerID(ctx context.Context, msg *ListUsersByWLPartnerIDQuery) (err error) {
	msg.Result, err = h.inner.ListUsersByWLPartnerID(msg.GetArgs(ctx))
	return err
}

// implement interfaces

func (q *BlockUserCommand) command()                   {}
func (q *CreateAccountUserCommand) command()           {}
func (q *CreateAffiliateCommand) command()             {}
func (q *CreateShopCommand) command()                  {}
func (q *DeleteAccountCommand) command()               {}
func (q *DeleteAccountUsersCommand) command()          {}
func (q *DeleteAffiliateCommand) command()             {}
func (q *RegisterSimplifyCommand) command()            {}
func (q *UnblockUserCommand) command()                 {}
func (q *UpdateAccountUserPermissionCommand) command() {}
func (q *UpdateAffiliateBankAccountCommand) command()  {}
func (q *UpdateAffiliateInfoCommand) command()         {}
func (q *UpdateExtensionNumberNormCommand) command()   {}
func (q *UpdateShipFromAddressIDCommand) command()     {}
func (q *UpdateShopInfoCommand) command()              {}
func (q *UpdateUserEmailCommand) command()             {}
func (q *UpdateUserPhoneCommand) command()             {}
func (q *UpdateUserRefCommand) command()               {}
func (q *UpdateUserReferenceSaleIDCommand) command()   {}
func (q *UpdateUserReferenceUserIDCommand) command()   {}

func (q *GetAccountByIDQuery) query()                   {}
func (q *GetAccountUserQuery) query()                   {}
func (q *GetAffiliateByIDQuery) query()                 {}
func (q *GetAffiliateWithPermissionQuery) query()       {}
func (q *GetAffiliatesByIDsQuery) query()               {}
func (q *GetAffiliatesByOwnerIDQuery) query()           {}
func (q *GetAllAccountsByUsersQuery) query()            {}
func (q *GetPartnerByIDQuery) query()                   {}
func (q *GetShopByIDQuery) query()                      {}
func (q *GetUserByEmailQuery) query()                   {}
func (q *GetUserByIDQuery) query()                      {}
func (q *GetUserByPhoneQuery) query()                   {}
func (q *GetUserByPhoneOrEmailQuery) query()            {}
func (q *GetUserFtRefSaffByIDQuery) query()             {}
func (q *GetUserFtRefSaffsQuery) query()                {}
func (q *GetUsersQuery) query()                         {}
func (q *GetUsersByAccountQuery) query()                {}
func (q *GetUsersByIDsQuery) query()                    {}
func (q *ListAccountUsersQuery) query()                 {}
func (q *ListPartnerRelationsBySubjectIDsQuery) query() {}
func (q *ListPartnersForWhiteLabelQuery) query()        {}
func (q *ListShopExtendedsQuery) query()                {}
func (q *ListShopsByIDsQuery) query()                   {}
func (q *ListUsersByIDsAndNameNormQuery) query()        {}
func (q *ListUsersByWLPartnerIDQuery) query()           {}

// implement conversion

func (q *BlockUserCommand) GetArgs(ctx context.Context) (_ context.Context, _ *BlockUserArgs) {
	return ctx,
		&BlockUserArgs{
			UserID:      q.UserID,
			BlockBy:     q.BlockBy,
			BlockReason: q.BlockReason,
		}
}

func (q *BlockUserCommand) SetBlockUserArgs(args *BlockUserArgs) {
	q.UserID = args.UserID
	q.BlockBy = args.BlockBy
	q.BlockReason = args.BlockReason
}

func (q *CreateAccountUserCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateAccountUserArgs) {
	return ctx,
		&CreateAccountUserArgs{
			AccountID:            q.AccountID,
			UserID:               q.UserID,
			Status:               q.Status,
			Permission:           q.Permission,
			FullName:             q.FullName,
			ShortName:            q.ShortName,
			Position:             q.Position,
			InvitationSentAt:     q.InvitationSentAt,
			InvitationSentBy:     q.InvitationSentBy,
			InvitationAcceptedAt: q.InvitationAcceptedAt,
		}
}

func (q *CreateAccountUserCommand) SetCreateAccountUserArgs(args *CreateAccountUserArgs) {
	q.AccountID = args.AccountID
	q.UserID = args.UserID
	q.Status = args.Status
	q.Permission = args.Permission
	q.FullName = args.FullName
	q.ShortName = args.ShortName
	q.Position = args.Position
	q.InvitationSentAt = args.InvitationSentAt
	q.InvitationSentBy = args.InvitationSentBy
	q.InvitationAcceptedAt = args.InvitationAcceptedAt
}

func (q *CreateAffiliateCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateAffiliateArgs) {
	return ctx,
		&CreateAffiliateArgs{
			Name:        q.Name,
			OwnerID:     q.OwnerID,
			Phone:       q.Phone,
			Email:       q.Email,
			IsTest:      q.IsTest,
			BankAccount: q.BankAccount,
		}
}

func (q *CreateAffiliateCommand) SetCreateAffiliateArgs(args *CreateAffiliateArgs) {
	q.Name = args.Name
	q.OwnerID = args.OwnerID
	q.Phone = args.Phone
	q.Email = args.Email
	q.IsTest = args.IsTest
	q.BankAccount = args.BankAccount
}

func (q *CreateShopCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateShopArgs) {
	return ctx,
		&CreateShopArgs{
			ID:                          q.ID,
			Name:                        q.Name,
			OwnerID:                     q.OwnerID,
			AddressID:                   q.AddressID,
			Address:                     q.Address,
			Phone:                       q.Phone,
			BankAccount:                 q.BankAccount,
			WebsiteURL:                  q.WebsiteURL,
			ImageURL:                    q.ImageURL,
			Email:                       q.Email,
			AutoCreateFFM:               q.AutoCreateFFM,
			IsTest:                      q.IsTest,
			URLSlug:                     q.URLSlug,
			CompanyInfo:                 q.CompanyInfo,
			MoneyTransactionRRule:       q.MoneyTransactionRRule,
			SurveyInfo:                  q.SurveyInfo,
			ShippingServicePickStrategy: q.ShippingServicePickStrategy,
		}
}

func (q *CreateShopCommand) SetCreateShopArgs(args *CreateShopArgs) {
	q.ID = args.ID
	q.Name = args.Name
	q.OwnerID = args.OwnerID
	q.AddressID = args.AddressID
	q.Address = args.Address
	q.Phone = args.Phone
	q.BankAccount = args.BankAccount
	q.WebsiteURL = args.WebsiteURL
	q.ImageURL = args.ImageURL
	q.Email = args.Email
	q.AutoCreateFFM = args.AutoCreateFFM
	q.IsTest = args.IsTest
	q.URLSlug = args.URLSlug
	q.CompanyInfo = args.CompanyInfo
	q.MoneyTransactionRRule = args.MoneyTransactionRRule
	q.SurveyInfo = args.SurveyInfo
	q.ShippingServicePickStrategy = args.ShippingServicePickStrategy
}

func (q *DeleteAccountCommand) GetArgs(ctx context.Context) (_ context.Context, _ *DeleteAccountArgs) {
	return ctx,
		&DeleteAccountArgs{
			AccountID: q.AccountID,
			OwnerID:   q.OwnerID,
		}
}

func (q *DeleteAccountCommand) SetDeleteAccountArgs(args *DeleteAccountArgs) {
	q.AccountID = args.AccountID
	q.OwnerID = args.OwnerID
}

func (q *DeleteAccountUsersCommand) GetArgs(ctx context.Context) (_ context.Context, _ *DeleteAccountUsersArgs) {
	return ctx,
		&DeleteAccountUsersArgs{
			AccountID: q.AccountID,
			UserID:    q.UserID,
		}
}

func (q *DeleteAccountUsersCommand) SetDeleteAccountUsersArgs(args *DeleteAccountUsersArgs) {
	q.AccountID = args.AccountID
	q.UserID = args.UserID
}

func (q *DeleteAffiliateCommand) GetArgs(ctx context.Context) (_ context.Context, _ *DeleteAffiliateArgs) {
	return ctx,
		&DeleteAffiliateArgs{
			ID:      q.ID,
			OwnerID: q.OwnerID,
		}
}

func (q *DeleteAffiliateCommand) SetDeleteAffiliateArgs(args *DeleteAffiliateArgs) {
	q.ID = args.ID
	q.OwnerID = args.OwnerID
}

func (q *RegisterSimplifyCommand) GetArgs(ctx context.Context) (_ context.Context, _ *RegisterSimplifyArgs) {
	return ctx,
		&RegisterSimplifyArgs{
			Phone:               q.Phone,
			Password:            q.Password,
			FullName:            q.FullName,
			Email:               q.Email,
			CompanyName:         q.CompanyName,
			IsCreateDefaultShop: q.IsCreateDefaultShop,
		}
}

func (q *RegisterSimplifyCommand) SetRegisterSimplifyArgs(args *RegisterSimplifyArgs) {
	q.Phone = args.Phone
	q.Password = args.Password
	q.FullName = args.FullName
	q.Email = args.Email
	q.CompanyName = args.CompanyName
	q.IsCreateDefaultShop = args.IsCreateDefaultShop
}

func (q *UnblockUserCommand) GetArgs(ctx context.Context) (_ context.Context, userID dot.ID) {
	return ctx,
		q.UserID
}

func (q *UpdateAccountUserPermissionCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateAccountUserPermissionArgs) {
	return ctx,
		&UpdateAccountUserPermissionArgs{
			AccountID:  q.AccountID,
			UserID:     q.UserID,
			Permission: q.Permission,
		}
}

func (q *UpdateAccountUserPermissionCommand) SetUpdateAccountUserPermissionArgs(args *UpdateAccountUserPermissionArgs) {
	q.AccountID = args.AccountID
	q.UserID = args.UserID
	q.Permission = args.Permission
}

func (q *UpdateAffiliateBankAccountCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateAffiliateBankAccountArgs) {
	return ctx,
		&UpdateAffiliateBankAccountArgs{
			ID:          q.ID,
			OwnerID:     q.OwnerID,
			BankAccount: q.BankAccount,
		}
}

func (q *UpdateAffiliateBankAccountCommand) SetUpdateAffiliateBankAccountArgs(args *UpdateAffiliateBankAccountArgs) {
	q.ID = args.ID
	q.OwnerID = args.OwnerID
	q.BankAccount = args.BankAccount
}

func (q *UpdateAffiliateInfoCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateAffiliateInfoArgs) {
	return ctx,
		&UpdateAffiliateInfoArgs{
			ID:      q.ID,
			OwnerID: q.OwnerID,
			Phone:   q.Phone,
			Email:   q.Email,
			Name:    q.Name,
		}
}

func (q *UpdateAffiliateInfoCommand) SetUpdateAffiliateInfoArgs(args *UpdateAffiliateInfoArgs) {
	q.ID = args.ID
	q.OwnerID = args.OwnerID
	q.Phone = args.Phone
	q.Email = args.Email
	q.Name = args.Name
}

func (q *UpdateExtensionNumberNormCommand) GetArgs(ctx context.Context) (_ context.Context, accountID dot.ID, userID dot.ID, extensionNumber string) {
	return ctx,
		q.AccountID,
		q.UserID,
		q.ExtensionNumber
}

func (q *UpdateShipFromAddressIDCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateShipFromAddressArgs) {
	return ctx,
		&UpdateShipFromAddressArgs{
			ID:                q.ID,
			ShipFromAddressID: q.ShipFromAddressID,
		}
}

func (q *UpdateShipFromAddressIDCommand) SetUpdateShipFromAddressArgs(args *UpdateShipFromAddressArgs) {
	q.ID = args.ID
	q.ShipFromAddressID = args.ShipFromAddressID
}

func (q *UpdateShopInfoCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateShopInfoArgs) {
	return ctx,
		&UpdateShopInfoArgs{
			ShopID:                  q.ShopID,
			MoneyTransactionRrule:   q.MoneyTransactionRrule,
			IsPriorMoneyTransaction: q.IsPriorMoneyTransaction,
		}
}

func (q *UpdateShopInfoCommand) SetUpdateShopInfoArgs(args *UpdateShopInfoArgs) {
	q.ShopID = args.ShopID
	q.MoneyTransactionRrule = args.MoneyTransactionRrule
	q.IsPriorMoneyTransaction = args.IsPriorMoneyTransaction
}

func (q *UpdateUserEmailCommand) GetArgs(ctx context.Context) (_ context.Context, userID dot.ID, email string) {
	return ctx,
		q.UserID,
		q.Email
}

func (q *UpdateUserPhoneCommand) GetArgs(ctx context.Context) (_ context.Context, userID dot.ID, phone string) {
	return ctx,
		q.UserID,
		q.Phone
}

func (q *UpdateUserRefCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateUserRefArgs) {
	return ctx,
		&UpdateUserRefArgs{
			UserID:  q.UserID,
			RefAff:  q.RefAff,
			RefSale: q.RefSale,
		}
}

func (q *UpdateUserRefCommand) SetUpdateUserRefArgs(args *UpdateUserRefArgs) {
	q.UserID = args.UserID
	q.RefAff = args.RefAff
	q.RefSale = args.RefSale
}

func (q *UpdateUserReferenceSaleIDCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateUserReferenceSaleIDArgs) {
	return ctx,
		&UpdateUserReferenceSaleIDArgs{
			UserID:       q.UserID,
			RefSalePhone: q.RefSalePhone,
		}
}

func (q *UpdateUserReferenceSaleIDCommand) SetUpdateUserReferenceSaleIDArgs(args *UpdateUserReferenceSaleIDArgs) {
	q.UserID = args.UserID
	q.RefSalePhone = args.RefSalePhone
}

func (q *UpdateUserReferenceUserIDCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateUserReferenceUserIDArgs) {
	return ctx,
		&UpdateUserReferenceUserIDArgs{
			UserID:       q.UserID,
			RefUserPhone: q.RefUserPhone,
		}
}

func (q *UpdateUserReferenceUserIDCommand) SetUpdateUserReferenceUserIDArgs(args *UpdateUserReferenceUserIDArgs) {
	q.UserID = args.UserID
	q.RefUserPhone = args.RefUserPhone
}

func (q *GetAccountByIDQuery) GetArgs(ctx context.Context) (_ context.Context, ID dot.ID) {
	return ctx,
		q.ID
}

func (q *GetAccountUserQuery) GetArgs(ctx context.Context) (_ context.Context, UserID dot.ID, AccountID dot.ID) {
	return ctx,
		q.UserID,
		q.AccountID
}

func (q *GetAffiliateByIDQuery) GetArgs(ctx context.Context) (_ context.Context, ID dot.ID) {
	return ctx,
		q.ID
}

func (q *GetAffiliateWithPermissionQuery) GetArgs(ctx context.Context) (_ context.Context, AffiliateID dot.ID, UserID dot.ID) {
	return ctx,
		q.AffiliateID,
		q.UserID
}

func (q *GetAffiliatesByIDsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetAffiliatesByIDsArgs) {
	return ctx,
		&GetAffiliatesByIDsArgs{
			AffiliateIDs: q.AffiliateIDs,
		}
}

func (q *GetAffiliatesByIDsQuery) SetGetAffiliatesByIDsArgs(args *GetAffiliatesByIDsArgs) {
	q.AffiliateIDs = args.AffiliateIDs
}

func (q *GetAffiliatesByOwnerIDQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetAffiliatesByOwnerIDArgs) {
	return ctx,
		&GetAffiliatesByOwnerIDArgs{
			ID: q.ID,
		}
}

func (q *GetAffiliatesByOwnerIDQuery) SetGetAffiliatesByOwnerIDArgs(args *GetAffiliatesByOwnerIDArgs) {
	q.ID = args.ID
}

func (q *GetAllAccountsByUsersQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetAllAccountUsersArg) {
	return ctx,
		&GetAllAccountUsersArg{
			UserIDs: q.UserIDs,
			Type:    q.Type,
			Roles:   q.Roles,
		}
}

func (q *GetAllAccountsByUsersQuery) SetGetAllAccountUsersArg(args *GetAllAccountUsersArg) {
	q.UserIDs = args.UserIDs
	q.Type = args.Type
	q.Roles = args.Roles
}

func (q *GetPartnerByIDQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetPartnerByIDArgs) {
	return ctx,
		&GetPartnerByIDArgs{
			ID: q.ID,
		}
}

func (q *GetPartnerByIDQuery) SetGetPartnerByIDArgs(args *GetPartnerByIDArgs) {
	q.ID = args.ID
}

func (q *GetShopByIDQuery) GetArgs(ctx context.Context) (_ context.Context, ID dot.ID) {
	return ctx,
		q.ID
}

func (q *GetUserByEmailQuery) GetArgs(ctx context.Context) (_ context.Context, email string) {
	return ctx,
		q.Email
}

func (q *GetUserByIDQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetUserByIDQueryArgs) {
	return ctx,
		&GetUserByIDQueryArgs{
			UserID: q.UserID,
		}
}

func (q *GetUserByIDQuery) SetGetUserByIDQueryArgs(args *GetUserByIDQueryArgs) {
	q.UserID = args.UserID
}

func (q *GetUserByPhoneQuery) GetArgs(ctx context.Context) (_ context.Context, phone string) {
	return ctx,
		q.Phone
}

func (q *GetUserByPhoneOrEmailQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetUserByPhoneOrEmailArgs) {
	return ctx,
		&GetUserByPhoneOrEmailArgs{
			Phone: q.Phone,
			Email: q.Email,
		}
}

func (q *GetUserByPhoneOrEmailQuery) SetGetUserByPhoneOrEmailArgs(args *GetUserByPhoneOrEmailArgs) {
	q.Phone = args.Phone
	q.Email = args.Email
}

func (q *GetUserFtRefSaffByIDQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetUserByIDQueryArgs) {
	return ctx,
		&GetUserByIDQueryArgs{
			UserID: q.UserID,
		}
}

func (q *GetUserFtRefSaffByIDQuery) SetGetUserByIDQueryArgs(args *GetUserByIDQueryArgs) {
	q.UserID = args.UserID
}

func (q *GetUserFtRefSaffsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *ListUserFtRefSaffsArgs) {
	return ctx,
		&ListUserFtRefSaffsArgs{
			Name:      q.Name,
			Phone:     q.Phone,
			Email:     q.Email,
			RefAff:    q.RefAff,
			RefSale:   q.RefSale,
			CreatedAt: q.CreatedAt,
			Paging:    q.Paging,
		}
}

func (q *GetUserFtRefSaffsQuery) SetListUserFtRefSaffsArgs(args *ListUserFtRefSaffsArgs) {
	q.Name = args.Name
	q.Phone = args.Phone
	q.Email = args.Email
	q.RefAff = args.RefAff
	q.RefSale = args.RefSale
	q.CreatedAt = args.CreatedAt
	q.Paging = args.Paging
}

func (q *GetUsersQuery) GetArgs(ctx context.Context) (_ context.Context, _ *ListUsersArgs) {
	return ctx,
		&ListUsersArgs{
			Name:      q.Name,
			Phone:     q.Phone,
			Email:     q.Email,
			CreatedAt: q.CreatedAt,
			Paging:    q.Paging,
		}
}

func (q *GetUsersQuery) SetListUsersArgs(args *ListUsersArgs) {
	q.Name = args.Name
	q.Phone = args.Phone
	q.Email = args.Email
	q.CreatedAt = args.CreatedAt
	q.Paging = args.Paging
}

func (q *GetUsersByAccountQuery) GetArgs(ctx context.Context) (_ context.Context, accountID dot.ID) {
	return ctx,
		q.AccountID
}

func (q *GetUsersByIDsQuery) GetArgs(ctx context.Context) (_ context.Context, IDs []dot.ID) {
	return ctx,
		q.IDs
}

func (q *ListAccountUsersQuery) GetArgs(ctx context.Context) (_ context.Context, _ *ListAccountUsersArgs) {
	return ctx,
		&ListAccountUsersArgs{
			AccountID: q.AccountID,
			UserID:    q.UserID,
		}
}

func (q *ListAccountUsersQuery) SetListAccountUsersArgs(args *ListAccountUsersArgs) {
	q.AccountID = args.AccountID
	q.UserID = args.UserID
}

func (q *ListPartnerRelationsBySubjectIDsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *ListPartnerRelationsBySubjectIDsArgs) {
	return ctx,
		&ListPartnerRelationsBySubjectIDsArgs{
			SubjectIDs:  q.SubjectIDs,
			SubjectType: q.SubjectType,
		}
}

func (q *ListPartnerRelationsBySubjectIDsQuery) SetListPartnerRelationsBySubjectIDsArgs(args *ListPartnerRelationsBySubjectIDsArgs) {
	q.SubjectIDs = args.SubjectIDs
	q.SubjectType = args.SubjectType
}

func (q *ListPartnersForWhiteLabelQuery) GetArgs(ctx context.Context) (_ context.Context, _ *meta.Empty) {
	return ctx,
		&meta.Empty{}
}

func (q *ListPartnersForWhiteLabelQuery) SetEmpty(args *meta.Empty) {
}

func (q *ListShopExtendedsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *ListShopQuery) {
	return ctx,
		&ListShopQuery{
			Paging:               q.Paging,
			Filters:              q.Filters,
			Name:                 q.Name,
			ShopIDs:              q.ShopIDs,
			OwnerID:              q.OwnerID,
			IncludeWLPartnerShop: q.IncludeWLPartnerShop,
		}
}

func (q *ListShopExtendedsQuery) SetListShopQuery(args *ListShopQuery) {
	q.Paging = args.Paging
	q.Filters = args.Filters
	q.Name = args.Name
	q.ShopIDs = args.ShopIDs
	q.OwnerID = args.OwnerID
	q.IncludeWLPartnerShop = args.IncludeWLPartnerShop
}

func (q *ListShopsByIDsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *ListShopsByIDsArgs) {
	return ctx,
		&ListShopsByIDsArgs{
			IDs:                     q.IDs,
			IsPriorMoneyTransaction: q.IsPriorMoneyTransaction,
			IncludeWLPartnerShop:    q.IncludeWLPartnerShop,
		}
}

func (q *ListShopsByIDsQuery) SetListShopsByIDsArgs(args *ListShopsByIDsArgs) {
	q.IDs = args.IDs
	q.IsPriorMoneyTransaction = args.IsPriorMoneyTransaction
	q.IncludeWLPartnerShop = args.IncludeWLPartnerShop
}

func (q *ListUsersByIDsAndNameNormQuery) GetArgs(ctx context.Context) (_ context.Context, _ *ListUsersByIDsAndNameNormArgs) {
	return ctx,
		&ListUsersByIDsAndNameNormArgs{
			IDs:      q.IDs,
			NameNorm: q.NameNorm,
		}
}

func (q *ListUsersByIDsAndNameNormQuery) SetListUsersByIDsAndNameNormArgs(args *ListUsersByIDsAndNameNormArgs) {
	q.IDs = args.IDs
	q.NameNorm = args.NameNorm
}

func (q *ListUsersByWLPartnerIDQuery) GetArgs(ctx context.Context) (_ context.Context, _ *ListUsersByWLPartnerID) {
	return ctx,
		&ListUsersByWLPartnerID{
			ID: q.ID,
		}
}

func (q *ListUsersByWLPartnerIDQuery) SetListUsersByWLPartnerID(args *ListUsersByWLPartnerID) {
	q.ID = args.ID
}

// implement dispatching

type AggregateHandler struct {
	inner Aggregate
}

func NewAggregateHandler(service Aggregate) AggregateHandler { return AggregateHandler{service} }

func (h AggregateHandler) RegisterHandlers(b interface {
	capi.Bus
	AddHandler(handler interface{})
}) CommandBus {
	b.AddHandler(h.HandleBlockUser)
	b.AddHandler(h.HandleCreateAccountUser)
	b.AddHandler(h.HandleCreateAffiliate)
	b.AddHandler(h.HandleCreateShop)
	b.AddHandler(h.HandleDeleteAccount)
	b.AddHandler(h.HandleDeleteAccountUsers)
	b.AddHandler(h.HandleDeleteAffiliate)
	b.AddHandler(h.HandleRegisterSimplify)
	b.AddHandler(h.HandleUnblockUser)
	b.AddHandler(h.HandleUpdateAccountUserPermission)
	b.AddHandler(h.HandleUpdateAffiliateBankAccount)
	b.AddHandler(h.HandleUpdateAffiliateInfo)
	b.AddHandler(h.HandleUpdateExtensionNumberNorm)
	b.AddHandler(h.HandleUpdateShipFromAddressID)
	b.AddHandler(h.HandleUpdateShopInfo)
	b.AddHandler(h.HandleUpdateUserEmail)
	b.AddHandler(h.HandleUpdateUserPhone)
	b.AddHandler(h.HandleUpdateUserRef)
	b.AddHandler(h.HandleUpdateUserReferenceSaleID)
	b.AddHandler(h.HandleUpdateUserReferenceUserID)
	return CommandBus{b}
}

type QueryServiceHandler struct {
	inner QueryService
}

func NewQueryServiceHandler(service QueryService) QueryServiceHandler {
	return QueryServiceHandler{service}
}

func (h QueryServiceHandler) RegisterHandlers(b interface {
	capi.Bus
	AddHandler(handler interface{})
}) QueryBus {
	b.AddHandler(h.HandleGetAccountByID)
	b.AddHandler(h.HandleGetAccountUser)
	b.AddHandler(h.HandleGetAffiliateByID)
	b.AddHandler(h.HandleGetAffiliateWithPermission)
	b.AddHandler(h.HandleGetAffiliatesByIDs)
	b.AddHandler(h.HandleGetAffiliatesByOwnerID)
	b.AddHandler(h.HandleGetAllAccountsByUsers)
	b.AddHandler(h.HandleGetPartnerByID)
	b.AddHandler(h.HandleGetShopByID)
	b.AddHandler(h.HandleGetUserByEmail)
	b.AddHandler(h.HandleGetUserByID)
	b.AddHandler(h.HandleGetUserByPhone)
	b.AddHandler(h.HandleGetUserByPhoneOrEmail)
	b.AddHandler(h.HandleGetUserFtRefSaffByID)
	b.AddHandler(h.HandleGetUserFtRefSaffs)
	b.AddHandler(h.HandleGetUsers)
	b.AddHandler(h.HandleGetUsersByAccount)
	b.AddHandler(h.HandleGetUsersByIDs)
	b.AddHandler(h.HandleListAccountUsers)
	b.AddHandler(h.HandleListPartnerRelationsBySubjectIDs)
	b.AddHandler(h.HandleListPartnersForWhiteLabel)
	b.AddHandler(h.HandleListShopExtendeds)
	b.AddHandler(h.HandleListShopsByIDs)
	b.AddHandler(h.HandleListUsersByIDsAndNameNorm)
	b.AddHandler(h.HandleListUsersByWLPartnerID)
	return QueryBus{b}
}
