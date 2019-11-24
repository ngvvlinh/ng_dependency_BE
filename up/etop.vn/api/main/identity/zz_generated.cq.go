// +build !generator

// Code generated by generator cq. DO NOT EDIT.

package identity

import (
	context "context"

	capi "etop.vn/capi"
	dot "etop.vn/capi/dot"
)

type Command interface{ command() }
type Query interface{ query() }
type CommandBus struct{ bus capi.Bus }
type QueryBus struct{ bus capi.Bus }

func NewCommandBus(bus capi.Bus) CommandBus                          { return CommandBus{bus} }
func NewQueryBus(bus capi.Bus) QueryBus                              { return QueryBus{bus} }
func (c CommandBus) Dispatch(ctx context.Context, msg Command) error { return c.bus.Dispatch(ctx, msg) }
func (c QueryBus) Dispatch(ctx context.Context, msg Query) error     { return c.bus.Dispatch(ctx, msg) }
func (c CommandBus) DispatchAll(ctx context.Context, msgs ...Command) error {
	for _, msg := range msgs {
		if err := c.bus.Dispatch(ctx, msg); err != nil {
			return err
		}
	}
	return nil
}
func (c QueryBus) DispatchAll(ctx context.Context, msgs ...Query) error {
	for _, msg := range msgs {
		if err := c.bus.Dispatch(ctx, msg); err != nil {
			return err
		}
	}
	return nil
}

type CreateAffiliateCommand struct {
	Name        string
	OwnerID     dot.ID
	Phone       string
	Email       string
	IsTest      bool
	BankAccount *BankAccount

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

type UpdateAffiliateBankAccountCommand struct {
	ID          dot.ID
	OwnerID     dot.ID
	BankAccount *BankAccount

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

// implement interfaces

func (q *CreateAffiliateCommand) command()                          {}
func (q *CreateExternalAccountAhamoveCommand) command()             {}
func (q *DeleteAffiliateCommand) command()                          {}
func (q *RequestVerifyExternalAccountAhamoveCommand) command()      {}
func (q *UpdateAffiliateBankAccountCommand) command()               {}
func (q *UpdateAffiliateInfoCommand) command()                      {}
func (q *UpdateExternalAccountAhamoveVerificationCommand) command() {}
func (q *UpdateUserReferenceSaleIDCommand) command()                {}
func (q *UpdateUserReferenceUserIDCommand) command()                {}
func (q *UpdateVerifiedExternalAccountAhamoveCommand) command()     {}
func (q *GetAffiliateByIDQuery) query()                             {}
func (q *GetAffiliateWithPermissionQuery) query()                   {}
func (q *GetAffiliatesByIDsQuery) query()                           {}
func (q *GetAffiliatesByOwnerIDQuery) query()                       {}
func (q *GetExternalAccountAhamoveQuery) query()                    {}
func (q *GetExternalAccountAhamoveByExternalIDQuery) query()        {}
func (q *GetShopByIDQuery) query()                                  {}
func (q *GetUserByEmailQuery) query()                               {}
func (q *GetUserByIDQuery) query()                                  {}
func (q *GetUserByPhoneQuery) query()                               {}

// implement conversion

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

// implement dispatching

type AggregateHandler struct {
	inner Aggregate
}

func NewAggregateHandler(service Aggregate) AggregateHandler { return AggregateHandler{service} }

func (h AggregateHandler) RegisterHandlers(b interface {
	capi.Bus
	AddHandler(handler interface{})
}) CommandBus {
	b.AddHandler(h.HandleCreateAffiliate)
	b.AddHandler(h.HandleCreateExternalAccountAhamove)
	b.AddHandler(h.HandleDeleteAffiliate)
	b.AddHandler(h.HandleRequestVerifyExternalAccountAhamove)
	b.AddHandler(h.HandleUpdateAffiliateBankAccount)
	b.AddHandler(h.HandleUpdateAffiliateInfo)
	b.AddHandler(h.HandleUpdateExternalAccountAhamoveVerification)
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
	b.AddHandler(h.HandleGetAffiliateByID)
	b.AddHandler(h.HandleGetAffiliateWithPermission)
	b.AddHandler(h.HandleGetAffiliatesByIDs)
	b.AddHandler(h.HandleGetAffiliatesByOwnerID)
	b.AddHandler(h.HandleGetExternalAccountAhamove)
	b.AddHandler(h.HandleGetExternalAccountAhamoveByExternalID)
	b.AddHandler(h.HandleGetShopByID)
	b.AddHandler(h.HandleGetUserByEmail)
	b.AddHandler(h.HandleGetUserByID)
	b.AddHandler(h.HandleGetUserByPhone)
	return QueryBus{b}
}
