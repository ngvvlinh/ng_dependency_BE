// Code generated by gen-cmd-query. DO NOT EDIT.

package identity

import (
	context "context"

	meta "etop.vn/api/meta"
)

type Command interface{ command() }
type Query interface{ query() }
type CommandBus struct{ bus meta.Bus }
type QueryBus struct{ bus meta.Bus }

func (c CommandBus) Dispatch(ctx context.Context, msg Command) error {
	return c.bus.Dispatch(ctx, msg)
}
func (c QueryBus) Dispatch(ctx context.Context, msg Query) error {
	return c.bus.Dispatch(ctx, msg)
}
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
	Name    string
	OwnerID int64
	Phone   string
	Email   string
	IsTest  bool

	Result *Affiliate `json:"-"`
}

type CreateExternalAccountAhamoveCommand struct {
	OwnerID int64
	Phone   string
	Name    string
	Address string

	Result *ExternalAccountAhamove `json:"-"`
}

type DeleteAffiliateCommand struct {
	ID      int64
	OwnerID int64

	Result struct {
	} `json:"-"`
}

type RequestVerifyExternalAccountAhamoveCommand struct {
	OwnerID int64
	Phone   string

	Result *RequestVerifyExternalAccountAhamoveResult `json:"-"`
}

type UpdateAffiliateCommand struct {
	ID      int64
	OwnerID int64
	Phone   string
	Email   string
	Name    string

	Result *Affiliate `json:"-"`
}

type UpdateExternalAccountAhamoveVerificationCommand struct {
	OwnerID             int64
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

type UpdateUserReferenceSaleIDCommand struct {
	UserID       int64
	RefSalePhone string

	Result struct {
	} `json:"-"`
}

type UpdateUserReferenceUserIDCommand struct {
	UserID       int64
	RefUserPhone string

	Result struct {
	} `json:"-"`
}

type UpdateVerifiedExternalAccountAhamoveCommand struct {
	OwnerID int64
	Phone   string

	Result *ExternalAccountAhamove `json:"-"`
}

type GetAffiliateByIDQuery struct {
	ID int64

	Result *Affiliate `json:"-"`
}

type GetAffiliateWithPermissionQuery struct {
	AffiliateID int64
	UserID      int64

	Result *GetAffiliateWithPermissionResult `json:"-"`
}

type GetExternalAccountAhamoveQuery struct {
	OwnerID int64
	Phone   string

	Result *ExternalAccountAhamove `json:"-"`
}

type GetExternalAccountAhamoveByExternalIDQuery struct {
	ExternalID string

	Result *ExternalAccountAhamove `json:"-"`
}

type GetShopByIDQuery struct {
	ID int64

	Result *Shop `json:"-"`
}

type GetUserByIDQuery struct {
	UserID int64

	Result *User `json:"-"`
}

type GetUserByPhoneQuery struct {
	Phone string

	Result *User `json:"-"`
}

// implement interfaces

func (q *CreateAffiliateCommand) command()                          {}
func (q *CreateExternalAccountAhamoveCommand) command()             {}
func (q *DeleteAffiliateCommand) command()                          {}
func (q *RequestVerifyExternalAccountAhamoveCommand) command()      {}
func (q *UpdateAffiliateCommand) command()                          {}
func (q *UpdateExternalAccountAhamoveVerificationCommand) command() {}
func (q *UpdateUserReferenceSaleIDCommand) command()                {}
func (q *UpdateUserReferenceUserIDCommand) command()                {}
func (q *UpdateVerifiedExternalAccountAhamoveCommand) command()     {}
func (q *GetAffiliateByIDQuery) query()                             {}
func (q *GetAffiliateWithPermissionQuery) query()                   {}
func (q *GetExternalAccountAhamoveQuery) query()                    {}
func (q *GetExternalAccountAhamoveByExternalIDQuery) query()        {}
func (q *GetShopByIDQuery) query()                                  {}
func (q *GetUserByIDQuery) query()                                  {}
func (q *GetUserByPhoneQuery) query()                               {}

// implement conversion

func (q *CreateAffiliateCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateAffiliateArgs) {
	return ctx,
		&CreateAffiliateArgs{
			Name:    q.Name,
			OwnerID: q.OwnerID,
			Phone:   q.Phone,
			Email:   q.Email,
			IsTest:  q.IsTest,
		}
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

func (q *DeleteAffiliateCommand) GetArgs(ctx context.Context) (_ context.Context, _ *DeleteAffiliateArgs) {
	return ctx,
		&DeleteAffiliateArgs{
			ID:      q.ID,
			OwnerID: q.OwnerID,
		}
}

func (q *RequestVerifyExternalAccountAhamoveCommand) GetArgs(ctx context.Context) (_ context.Context, _ *RequestVerifyExternalAccountAhamoveArgs) {
	return ctx,
		&RequestVerifyExternalAccountAhamoveArgs{
			OwnerID: q.OwnerID,
			Phone:   q.Phone,
		}
}

func (q *UpdateAffiliateCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateAffiliateArgs) {
	return ctx,
		&UpdateAffiliateArgs{
			ID:      q.ID,
			OwnerID: q.OwnerID,
			Phone:   q.Phone,
			Email:   q.Email,
			Name:    q.Name,
		}
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

func (q *UpdateUserReferenceSaleIDCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateUserReferenceSaleIDArgs) {
	return ctx,
		&UpdateUserReferenceSaleIDArgs{
			UserID:       q.UserID,
			RefSalePhone: q.RefSalePhone,
		}
}

func (q *UpdateUserReferenceUserIDCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateUserReferenceUserIDArgs) {
	return ctx,
		&UpdateUserReferenceUserIDArgs{
			UserID:       q.UserID,
			RefUserPhone: q.RefUserPhone,
		}
}

func (q *UpdateVerifiedExternalAccountAhamoveCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateVerifiedExternalAccountAhamoveArgs) {
	return ctx,
		&UpdateVerifiedExternalAccountAhamoveArgs{
			OwnerID: q.OwnerID,
			Phone:   q.Phone,
		}
}

func (q *GetAffiliateByIDQuery) GetArgs(ctx context.Context) (_ context.Context, ID int64) {
	return ctx,
		q.ID
}

func (q *GetAffiliateWithPermissionQuery) GetArgs(ctx context.Context) (_ context.Context, AffiliateID int64, UserID int64) {
	return ctx,
		q.AffiliateID,
		q.UserID
}

func (q *GetExternalAccountAhamoveQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetExternalAccountAhamoveArgs) {
	return ctx,
		&GetExternalAccountAhamoveArgs{
			OwnerID: q.OwnerID,
			Phone:   q.Phone,
		}
}

func (q *GetExternalAccountAhamoveByExternalIDQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetExternalAccountAhamoveByExternalIDQueryArgs) {
	return ctx,
		&GetExternalAccountAhamoveByExternalIDQueryArgs{
			ExternalID: q.ExternalID,
		}
}

func (q *GetShopByIDQuery) GetArgs(ctx context.Context) (_ context.Context, ID int64) {
	return ctx,
		q.ID
}

func (q *GetUserByIDQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetUserByIDQueryArgs) {
	return ctx,
		&GetUserByIDQueryArgs{
			UserID: q.UserID,
		}
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
	meta.Bus
	AddHandler(handler interface{})
}) CommandBus {
	b.AddHandler(h.HandleCreateAffiliate)
	b.AddHandler(h.HandleCreateExternalAccountAhamove)
	b.AddHandler(h.HandleDeleteAffiliate)
	b.AddHandler(h.HandleRequestVerifyExternalAccountAhamove)
	b.AddHandler(h.HandleUpdateAffiliate)
	b.AddHandler(h.HandleUpdateExternalAccountAhamoveVerification)
	b.AddHandler(h.HandleUpdateUserReferenceSaleID)
	b.AddHandler(h.HandleUpdateUserReferenceUserID)
	b.AddHandler(h.HandleUpdateVerifiedExternalAccountAhamove)
	return CommandBus{b}
}

func (h AggregateHandler) HandleCreateAffiliate(ctx context.Context, msg *CreateAffiliateCommand) error {
	result, err := h.inner.CreateAffiliate(msg.GetArgs(ctx))
	msg.Result = result
	return err
}

func (h AggregateHandler) HandleCreateExternalAccountAhamove(ctx context.Context, msg *CreateExternalAccountAhamoveCommand) error {
	result, err := h.inner.CreateExternalAccountAhamove(msg.GetArgs(ctx))
	msg.Result = result
	return err
}

func (h AggregateHandler) HandleDeleteAffiliate(ctx context.Context, msg *DeleteAffiliateCommand) error {
	return h.inner.DeleteAffiliate(msg.GetArgs(ctx))
}

func (h AggregateHandler) HandleRequestVerifyExternalAccountAhamove(ctx context.Context, msg *RequestVerifyExternalAccountAhamoveCommand) error {
	result, err := h.inner.RequestVerifyExternalAccountAhamove(msg.GetArgs(ctx))
	msg.Result = result
	return err
}

func (h AggregateHandler) HandleUpdateAffiliate(ctx context.Context, msg *UpdateAffiliateCommand) error {
	result, err := h.inner.UpdateAffiliate(msg.GetArgs(ctx))
	msg.Result = result
	return err
}

func (h AggregateHandler) HandleUpdateExternalAccountAhamoveVerification(ctx context.Context, msg *UpdateExternalAccountAhamoveVerificationCommand) error {
	result, err := h.inner.UpdateExternalAccountAhamoveVerification(msg.GetArgs(ctx))
	msg.Result = result
	return err
}

func (h AggregateHandler) HandleUpdateUserReferenceSaleID(ctx context.Context, msg *UpdateUserReferenceSaleIDCommand) error {
	return h.inner.UpdateUserReferenceSaleID(msg.GetArgs(ctx))
}

func (h AggregateHandler) HandleUpdateUserReferenceUserID(ctx context.Context, msg *UpdateUserReferenceUserIDCommand) error {
	return h.inner.UpdateUserReferenceUserID(msg.GetArgs(ctx))
}

func (h AggregateHandler) HandleUpdateVerifiedExternalAccountAhamove(ctx context.Context, msg *UpdateVerifiedExternalAccountAhamoveCommand) error {
	result, err := h.inner.UpdateVerifiedExternalAccountAhamove(msg.GetArgs(ctx))
	msg.Result = result
	return err
}

type QueryServiceHandler struct {
	inner QueryService
}

func NewQueryServiceHandler(service QueryService) QueryServiceHandler {
	return QueryServiceHandler{service}
}

func (h QueryServiceHandler) RegisterHandlers(b interface {
	meta.Bus
	AddHandler(handler interface{})
}) QueryBus {
	b.AddHandler(h.HandleGetAffiliateByID)
	b.AddHandler(h.HandleGetAffiliateWithPermission)
	b.AddHandler(h.HandleGetExternalAccountAhamove)
	b.AddHandler(h.HandleGetExternalAccountAhamoveByExternalID)
	b.AddHandler(h.HandleGetShopByID)
	b.AddHandler(h.HandleGetUserByID)
	b.AddHandler(h.HandleGetUserByPhone)
	return QueryBus{b}
}

func (h QueryServiceHandler) HandleGetAffiliateByID(ctx context.Context, msg *GetAffiliateByIDQuery) error {
	result, err := h.inner.GetAffiliateByID(msg.GetArgs(ctx))
	msg.Result = result
	return err
}

func (h QueryServiceHandler) HandleGetAffiliateWithPermission(ctx context.Context, msg *GetAffiliateWithPermissionQuery) error {
	result, err := h.inner.GetAffiliateWithPermission(msg.GetArgs(ctx))
	msg.Result = result
	return err
}

func (h QueryServiceHandler) HandleGetExternalAccountAhamove(ctx context.Context, msg *GetExternalAccountAhamoveQuery) error {
	result, err := h.inner.GetExternalAccountAhamove(msg.GetArgs(ctx))
	msg.Result = result
	return err
}

func (h QueryServiceHandler) HandleGetExternalAccountAhamoveByExternalID(ctx context.Context, msg *GetExternalAccountAhamoveByExternalIDQuery) error {
	result, err := h.inner.GetExternalAccountAhamoveByExternalID(msg.GetArgs(ctx))
	msg.Result = result
	return err
}

func (h QueryServiceHandler) HandleGetShopByID(ctx context.Context, msg *GetShopByIDQuery) error {
	result, err := h.inner.GetShopByID(msg.GetArgs(ctx))
	msg.Result = result
	return err
}

func (h QueryServiceHandler) HandleGetUserByID(ctx context.Context, msg *GetUserByIDQuery) error {
	result, err := h.inner.GetUserByID(msg.GetArgs(ctx))
	msg.Result = result
	return err
}

func (h QueryServiceHandler) HandleGetUserByPhone(ctx context.Context, msg *GetUserByPhoneQuery) error {
	result, err := h.inner.GetUserByPhone(msg.GetArgs(ctx))
	msg.Result = result
	return err
}
