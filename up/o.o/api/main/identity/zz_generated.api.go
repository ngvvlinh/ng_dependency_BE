// +build !generator

// Code generated by generator api. DO NOT EDIT.

package identity

import (
	context "context"

	identitytypes "o.o/api/main/identity/types"
	meta "o.o/api/meta"
	account_type "o.o/api/top/types/etc/account_type"
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

type CreateExternalAccountAhamoveCommand struct {
	OwnerID dot.ID
	Phone   string
	Name    string
	Address string

	Result *ExternalAccountAhamove `json:"-"`
}

func (h AggregateHandler) HandleCreateExternalAccountAhamove(ctx context.Context, msg *CreateExternalAccountAhamoveCommand) (err error) {
	msg.Result, err = h.inner.CreateExternalAccountAhamove(msg.GetArgs(ctx))
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

type RequestVerifyExternalAccountAhamoveCommand struct {
	OwnerID dot.ID
	Phone   string

	Result *RequestVerifyExternalAccountAhamoveResult `json:"-"`
}

func (h AggregateHandler) HandleRequestVerifyExternalAccountAhamove(ctx context.Context, msg *RequestVerifyExternalAccountAhamoveCommand) (err error) {
	msg.Result, err = h.inner.RequestVerifyExternalAccountAhamove(msg.GetArgs(ctx))
	return err
}

type UnblockUserCommand struct {
	UserID dot.ID

	Result *User `json:"-"`
}

func (h AggregateHandler) HandleUnblockUser(ctx context.Context, msg *UnblockUserCommand) (err error) {
	msg.Result, err = h.inner.UnblockUser(msg.GetArgs(ctx))
	return err
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

type UpdateExternalAccountAhamoveVerificationCommand struct {
	OwnerID             dot.ID
	Phone               string
	IDCardFrontImg      string
	IDCardBackImg       string
	PortraitImg         string
	WebsiteURL          string
	FanpageURL          string
	CompanyImgs         []string
	BusinessLicenseImgs []string

	Result *ExternalAccountAhamove `json:"-"`
}

func (h AggregateHandler) HandleUpdateExternalAccountAhamoveVerification(ctx context.Context, msg *UpdateExternalAccountAhamoveVerificationCommand) (err error) {
	msg.Result, err = h.inner.UpdateExternalAccountAhamoveVerification(msg.GetArgs(ctx))
	return err
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
	RefAff  string
	RefSale string

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

type UpdateVerifiedExternalAccountAhamoveCommand struct {
	OwnerID dot.ID
	Phone   string

	Result *ExternalAccountAhamove `json:"-"`
}

func (h AggregateHandler) HandleUpdateVerifiedExternalAccountAhamove(ctx context.Context, msg *UpdateVerifiedExternalAccountAhamoveCommand) (err error) {
	msg.Result, err = h.inner.UpdateVerifiedExternalAccountAhamove(msg.GetArgs(ctx))
	return err
}

type GetAccountByIDQuery struct {
	ID dot.ID

	Result *Account `json:"-"`
}

func (h QueryServiceHandler) HandleGetAccountByID(ctx context.Context, msg *GetAccountByIDQuery) (err error) {
	msg.Result, err = h.inner.GetAccountByID(msg.GetArgs(ctx))
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

	Result []*AccountUser `json:"-"`
}

func (h QueryServiceHandler) HandleGetAllAccountsByUsers(ctx context.Context, msg *GetAllAccountsByUsersQuery) (err error) {
	msg.Result, err = h.inner.GetAllAccountsByUsers(msg.GetArgs(ctx))
	return err
}

type GetExternalAccountAhamoveQuery struct {
	OwnerID dot.ID
	Phone   string

	Result *ExternalAccountAhamove `json:"-"`
}

func (h QueryServiceHandler) HandleGetExternalAccountAhamove(ctx context.Context, msg *GetExternalAccountAhamoveQuery) (err error) {
	msg.Result, err = h.inner.GetExternalAccountAhamove(msg.GetArgs(ctx))
	return err
}

type GetExternalAccountAhamoveByExternalIDQuery struct {
	ExternalID string

	Result *ExternalAccountAhamove `json:"-"`
}

func (h QueryServiceHandler) HandleGetExternalAccountAhamoveByExternalID(ctx context.Context, msg *GetExternalAccountAhamoveByExternalIDQuery) (err error) {
	msg.Result, err = h.inner.GetExternalAccountAhamoveByExternalID(msg.GetArgs(ctx))
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
	Name      string
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
	Name      string
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

type ListPartnersForWhiteLabelQuery struct {
	Result []*Partner `json:"-"`
}

func (h QueryServiceHandler) HandleListPartnersForWhiteLabel(ctx context.Context, msg *ListPartnersForWhiteLabelQuery) (err error) {
	msg.Result, err = h.inner.ListPartnersForWhiteLabel(msg.GetArgs(ctx))
	return err
}

type ListShopExtendedsQuery struct {
	Paging  meta.Paging
	Filters meta.Filters
	Name    filter.FullTextSearch
	ShopIDs []dot.ID

	Result *ListShopExtendedsResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListShopExtendeds(ctx context.Context, msg *ListShopExtendedsQuery) (err error) {
	msg.Result, err = h.inner.ListShopExtendeds(msg.GetArgs(ctx))
	return err
}

type ListShopsByIDsQuery struct {
	IDs []dot.ID

	Result []*Shop `json:"-"`
}

func (h QueryServiceHandler) HandleListShopsByIDs(ctx context.Context, msg *ListShopsByIDsQuery) (err error) {
	msg.Result, err = h.inner.ListShopsByIDs(msg.GetArgs(ctx))
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

func (q *BlockUserCommand) command()                                {}
func (q *CreateAffiliateCommand) command()                          {}
func (q *CreateExternalAccountAhamoveCommand) command()             {}
func (q *DeleteAffiliateCommand) command()                          {}
func (q *RequestVerifyExternalAccountAhamoveCommand) command()      {}
func (q *UnblockUserCommand) command()                              {}
func (q *UpdateAffiliateBankAccountCommand) command()               {}
func (q *UpdateAffiliateInfoCommand) command()                      {}
func (q *UpdateExternalAccountAhamoveVerificationCommand) command() {}
func (q *UpdateShipFromAddressIDCommand) command()                  {}
func (q *UpdateUserEmailCommand) command()                          {}
func (q *UpdateUserPhoneCommand) command()                          {}
func (q *UpdateUserRefCommand) command()                            {}
func (q *UpdateUserReferenceSaleIDCommand) command()                {}
func (q *UpdateUserReferenceUserIDCommand) command()                {}
func (q *UpdateVerifiedExternalAccountAhamoveCommand) command()     {}

func (q *GetAccountByIDQuery) query()                        {}
func (q *GetAffiliateByIDQuery) query()                      {}
func (q *GetAffiliateWithPermissionQuery) query()            {}
func (q *GetAffiliatesByIDsQuery) query()                    {}
func (q *GetAffiliatesByOwnerIDQuery) query()                {}
func (q *GetAllAccountsByUsersQuery) query()                 {}
func (q *GetExternalAccountAhamoveQuery) query()             {}
func (q *GetExternalAccountAhamoveByExternalIDQuery) query() {}
func (q *GetPartnerByIDQuery) query()                        {}
func (q *GetShopByIDQuery) query()                           {}
func (q *GetUserByEmailQuery) query()                        {}
func (q *GetUserByIDQuery) query()                           {}
func (q *GetUserByPhoneQuery) query()                        {}
func (q *GetUserByPhoneOrEmailQuery) query()                 {}
func (q *GetUserFtRefSaffByIDQuery) query()                  {}
func (q *GetUserFtRefSaffsQuery) query()                     {}
func (q *GetUsersQuery) query()                              {}
func (q *GetUsersByAccountQuery) query()                     {}
func (q *GetUsersByIDsQuery) query()                         {}
func (q *ListPartnersForWhiteLabelQuery) query()             {}
func (q *ListShopExtendedsQuery) query()                     {}
func (q *ListShopsByIDsQuery) query()                        {}
func (q *ListUsersByWLPartnerIDQuery) query()                {}

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

func (q *CreateExternalAccountAhamoveCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateExternalAccountAhamoveArgs) {
	return ctx,
		&CreateExternalAccountAhamoveArgs{
			OwnerID: q.OwnerID,
			Phone:   q.Phone,
			Name:    q.Name,
			Address: q.Address,
		}
}

func (q *CreateExternalAccountAhamoveCommand) SetCreateExternalAccountAhamoveArgs(args *CreateExternalAccountAhamoveArgs) {
	q.OwnerID = args.OwnerID
	q.Phone = args.Phone
	q.Name = args.Name
	q.Address = args.Address
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

func (q *RequestVerifyExternalAccountAhamoveCommand) GetArgs(ctx context.Context) (_ context.Context, _ *RequestVerifyExternalAccountAhamoveArgs) {
	return ctx,
		&RequestVerifyExternalAccountAhamoveArgs{
			OwnerID: q.OwnerID,
			Phone:   q.Phone,
		}
}

func (q *RequestVerifyExternalAccountAhamoveCommand) SetRequestVerifyExternalAccountAhamoveArgs(args *RequestVerifyExternalAccountAhamoveArgs) {
	q.OwnerID = args.OwnerID
	q.Phone = args.Phone
}

func (q *UnblockUserCommand) GetArgs(ctx context.Context) (_ context.Context, userID dot.ID) {
	return ctx,
		q.UserID
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

func (q *UpdateExternalAccountAhamoveVerificationCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateExternalAccountAhamoveVerificationArgs) {
	return ctx,
		&UpdateExternalAccountAhamoveVerificationArgs{
			OwnerID:             q.OwnerID,
			Phone:               q.Phone,
			IDCardFrontImg:      q.IDCardFrontImg,
			IDCardBackImg:       q.IDCardBackImg,
			PortraitImg:         q.PortraitImg,
			WebsiteURL:          q.WebsiteURL,
			FanpageURL:          q.FanpageURL,
			CompanyImgs:         q.CompanyImgs,
			BusinessLicenseImgs: q.BusinessLicenseImgs,
		}
}

func (q *UpdateExternalAccountAhamoveVerificationCommand) SetUpdateExternalAccountAhamoveVerificationArgs(args *UpdateExternalAccountAhamoveVerificationArgs) {
	q.OwnerID = args.OwnerID
	q.Phone = args.Phone
	q.IDCardFrontImg = args.IDCardFrontImg
	q.IDCardBackImg = args.IDCardBackImg
	q.PortraitImg = args.PortraitImg
	q.WebsiteURL = args.WebsiteURL
	q.FanpageURL = args.FanpageURL
	q.CompanyImgs = args.CompanyImgs
	q.BusinessLicenseImgs = args.BusinessLicenseImgs
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

func (q *UpdateVerifiedExternalAccountAhamoveCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateVerifiedExternalAccountAhamoveArgs) {
	return ctx,
		&UpdateVerifiedExternalAccountAhamoveArgs{
			OwnerID: q.OwnerID,
			Phone:   q.Phone,
		}
}

func (q *UpdateVerifiedExternalAccountAhamoveCommand) SetUpdateVerifiedExternalAccountAhamoveArgs(args *UpdateVerifiedExternalAccountAhamoveArgs) {
	q.OwnerID = args.OwnerID
	q.Phone = args.Phone
}

func (q *GetAccountByIDQuery) GetArgs(ctx context.Context) (_ context.Context, ID dot.ID) {
	return ctx,
		q.ID
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
		}
}

func (q *GetAllAccountsByUsersQuery) SetGetAllAccountUsersArg(args *GetAllAccountUsersArg) {
	q.UserIDs = args.UserIDs
	q.Type = args.Type
}

func (q *GetExternalAccountAhamoveQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetExternalAccountAhamoveArgs) {
	return ctx,
		&GetExternalAccountAhamoveArgs{
			OwnerID: q.OwnerID,
			Phone:   q.Phone,
		}
}

func (q *GetExternalAccountAhamoveQuery) SetGetExternalAccountAhamoveArgs(args *GetExternalAccountAhamoveArgs) {
	q.OwnerID = args.OwnerID
	q.Phone = args.Phone
}

func (q *GetExternalAccountAhamoveByExternalIDQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetExternalAccountAhamoveByExternalIDQueryArgs) {
	return ctx,
		&GetExternalAccountAhamoveByExternalIDQueryArgs{
			ExternalID: q.ExternalID,
		}
}

func (q *GetExternalAccountAhamoveByExternalIDQuery) SetGetExternalAccountAhamoveByExternalIDQueryArgs(args *GetExternalAccountAhamoveByExternalIDQueryArgs) {
	q.ExternalID = args.ExternalID
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

func (q *ListPartnersForWhiteLabelQuery) GetArgs(ctx context.Context) (_ context.Context, _ *meta.Empty) {
	return ctx,
		&meta.Empty{}
}

func (q *ListPartnersForWhiteLabelQuery) SetEmpty(args *meta.Empty) {
}

func (q *ListShopExtendedsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *ListShopQuery) {
	return ctx,
		&ListShopQuery{
			Paging:  q.Paging,
			Filters: q.Filters,
			Name:    q.Name,
			ShopIDs: q.ShopIDs,
		}
}

func (q *ListShopExtendedsQuery) SetListShopQuery(args *ListShopQuery) {
	q.Paging = args.Paging
	q.Filters = args.Filters
	q.Name = args.Name
	q.ShopIDs = args.ShopIDs
}

func (q *ListShopsByIDsQuery) GetArgs(ctx context.Context) (_ context.Context, IDs []dot.ID) {
	return ctx,
		q.IDs
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
	b.AddHandler(h.HandleCreateAffiliate)
	b.AddHandler(h.HandleCreateExternalAccountAhamove)
	b.AddHandler(h.HandleDeleteAffiliate)
	b.AddHandler(h.HandleRequestVerifyExternalAccountAhamove)
	b.AddHandler(h.HandleUnblockUser)
	b.AddHandler(h.HandleUpdateAffiliateBankAccount)
	b.AddHandler(h.HandleUpdateAffiliateInfo)
	b.AddHandler(h.HandleUpdateExternalAccountAhamoveVerification)
	b.AddHandler(h.HandleUpdateShipFromAddressID)
	b.AddHandler(h.HandleUpdateUserEmail)
	b.AddHandler(h.HandleUpdateUserPhone)
	b.AddHandler(h.HandleUpdateUserRef)
	b.AddHandler(h.HandleUpdateUserReferenceSaleID)
	b.AddHandler(h.HandleUpdateUserReferenceUserID)
	b.AddHandler(h.HandleUpdateVerifiedExternalAccountAhamove)
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
	b.AddHandler(h.HandleGetAffiliateByID)
	b.AddHandler(h.HandleGetAffiliateWithPermission)
	b.AddHandler(h.HandleGetAffiliatesByIDs)
	b.AddHandler(h.HandleGetAffiliatesByOwnerID)
	b.AddHandler(h.HandleGetAllAccountsByUsers)
	b.AddHandler(h.HandleGetExternalAccountAhamove)
	b.AddHandler(h.HandleGetExternalAccountAhamoveByExternalID)
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
	b.AddHandler(h.HandleListPartnersForWhiteLabel)
	b.AddHandler(h.HandleListShopExtendeds)
	b.AddHandler(h.HandleListShopsByIDs)
	b.AddHandler(h.HandleListUsersByWLPartnerID)
	return QueryBus{b}
}
