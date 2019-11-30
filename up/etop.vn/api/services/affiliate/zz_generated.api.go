// +build !generator

// Code generated by generator api. DO NOT EDIT.

package affiliate

import (
	context "context"

	meta "etop.vn/api/meta"
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

type CreateAffiliateReferralCodeCommand struct {
	AffiliateAccountID dot.ID
	Code               string

	Result *AffiliateReferralCode `json:"-"`
}

func (h AggregateHandler) HandleCreateAffiliateReferralCode(ctx context.Context, msg *CreateAffiliateReferralCodeCommand) (err error) {
	msg.Result, err = h.inner.CreateAffiliateReferralCode(msg.GetArgs(ctx))
	return err
}

type CreateOrUpdateCommissionSettingCommand struct {
	ProductID   dot.ID
	AccountID   dot.ID
	Amount      int
	Unit        string
	Type        string
	Description string
	Note        string

	Result *CommissionSetting `json:"-"`
}

func (h AggregateHandler) HandleCreateOrUpdateCommissionSetting(ctx context.Context, msg *CreateOrUpdateCommissionSettingCommand) (err error) {
	msg.Result, err = h.inner.CreateOrUpdateCommissionSetting(msg.GetArgs(ctx))
	return err
}

type CreateOrUpdateSupplyCommissionSettingCommand struct {
	ShopID                   dot.ID
	ProductID                dot.ID
	Level1DirectCommission   int
	Level1IndirectCommission int
	Level2DirectCommission   int
	Level2IndirectCommission int
	DependOn                 string
	Level1LimitCount         int
	Level1LimitDuration      int
	Level1LimitDurationType  string
	LifetimeDuration         int
	LifetimeDurationType     string
	Group                    string

	Result *SupplyCommissionSetting `json:"-"`
}

func (h AggregateHandler) HandleCreateOrUpdateSupplyCommissionSetting(ctx context.Context, msg *CreateOrUpdateSupplyCommissionSettingCommand) (err error) {
	msg.Result, err = h.inner.CreateOrUpdateSupplyCommissionSetting(msg.GetArgs(ctx))
	return err
}

type CreateOrUpdateUserReferralCommand struct {
	UserID           dot.ID
	ReferralCode     string
	SaleReferralCode string

	Result *UserReferral `json:"-"`
}

func (h AggregateHandler) HandleCreateOrUpdateUserReferral(ctx context.Context, msg *CreateOrUpdateUserReferralCommand) (err error) {
	msg.Result, err = h.inner.CreateOrUpdateUserReferral(msg.GetArgs(ctx))
	return err
}

type CreateProductPromotionCommand struct {
	ShopID      dot.ID
	ProductID   dot.ID
	Amount      int
	Code        string
	Description string
	Unit        string
	Note        string
	Type        string

	Result *ProductPromotion `json:"-"`
}

func (h AggregateHandler) HandleCreateProductPromotion(ctx context.Context, msg *CreateProductPromotionCommand) (err error) {
	msg.Result, err = h.inner.CreateProductPromotion(msg.GetArgs(ctx))
	return err
}

type OnTradingOrderCreatedCommand struct {
	OrderID      dot.ID
	ReferralCode string

	Result struct {
	} `json:"-"`
}

func (h AggregateHandler) HandleOnTradingOrderCreated(ctx context.Context, msg *OnTradingOrderCreatedCommand) (err error) {
	return h.inner.OnTradingOrderCreated(msg.GetArgs(ctx))
}

type OrderPaymentSuccessCommand struct {
	EventMeta meta.EventMeta
	OrderID   dot.ID

	Result struct {
	} `json:"-"`
}

func (h AggregateHandler) HandleOrderPaymentSuccess(ctx context.Context, msg *OrderPaymentSuccessCommand) (err error) {
	return h.inner.OrderPaymentSuccess(msg.GetArgs(ctx))
}

type TradingOrderCreatingCommand struct {
	ProductIDs   []dot.ID
	ReferralCode string
	UserID       dot.ID

	Result struct {
	} `json:"-"`
}

func (h AggregateHandler) HandleTradingOrderCreating(ctx context.Context, msg *TradingOrderCreatingCommand) (err error) {
	return h.inner.TradingOrderCreating(msg.GetArgs(ctx))
}

type UpdateProductPromotionCommand struct {
	ID          dot.ID
	Amount      int
	Unit        string
	Code        string
	Description string
	Note        string
	Type        string

	Result *ProductPromotion `json:"-"`
}

func (h AggregateHandler) HandleUpdateProductPromotion(ctx context.Context, msg *UpdateProductPromotionCommand) (err error) {
	msg.Result, err = h.inner.UpdateProductPromotion(msg.GetArgs(ctx))
	return err
}

type GetAffiliateAccountReferralByCodeQuery struct {
	Code string

	Result *AffiliateReferralCode `json:"-"`
}

func (h QueryServiceHandler) HandleGetAffiliateAccountReferralByCode(ctx context.Context, msg *GetAffiliateAccountReferralByCodeQuery) (err error) {
	msg.Result, err = h.inner.GetAffiliateAccountReferralByCode(msg.GetArgs(ctx))
	return err
}

type GetAffiliateAccountReferralCodesQuery struct {
	AffiliateAccountID dot.ID

	Result []*AffiliateReferralCode `json:"-"`
}

func (h QueryServiceHandler) HandleGetAffiliateAccountReferralCodes(ctx context.Context, msg *GetAffiliateAccountReferralCodesQuery) (err error) {
	msg.Result, err = h.inner.GetAffiliateAccountReferralCodes(msg.GetArgs(ctx))
	return err
}

type GetCommissionByProductIDQuery struct {
	AccountID dot.ID
	ProductID dot.ID

	Result *CommissionSetting `json:"-"`
}

func (h QueryServiceHandler) HandleGetCommissionByProductID(ctx context.Context, msg *GetCommissionByProductIDQuery) (err error) {
	msg.Result, err = h.inner.GetCommissionByProductID(msg.GetArgs(ctx))
	return err
}

type GetCommissionByProductIDsQuery struct {
	AccountID  dot.ID
	ProductIDs []dot.ID

	Result []*CommissionSetting `json:"-"`
}

func (h QueryServiceHandler) HandleGetCommissionByProductIDs(ctx context.Context, msg *GetCommissionByProductIDsQuery) (err error) {
	msg.Result, err = h.inner.GetCommissionByProductIDs(msg.GetArgs(ctx))
	return err
}

type GetReferralsByReferralIDQuery struct {
	ID dot.ID

	Result []*UserReferral `json:"-"`
}

func (h QueryServiceHandler) HandleGetReferralsByReferralID(ctx context.Context, msg *GetReferralsByReferralIDQuery) (err error) {
	msg.Result, err = h.inner.GetReferralsByReferralID(msg.GetArgs(ctx))
	return err
}

type GetSellerCommissionsQuery struct {
	SellerID dot.ID
	Paging   meta.Paging
	Filters  meta.Filters

	Result []*SellerCommission `json:"-"`
}

func (h QueryServiceHandler) HandleGetSellerCommissions(ctx context.Context, msg *GetSellerCommissionsQuery) (err error) {
	msg.Result, err = h.inner.GetSellerCommissions(msg.GetArgs(ctx))
	return err
}

type GetShopProductPromotionQuery struct {
	ShopID    dot.ID
	ProductID dot.ID

	Result *ProductPromotion `json:"-"`
}

func (h QueryServiceHandler) HandleGetShopProductPromotion(ctx context.Context, msg *GetShopProductPromotionQuery) (err error) {
	msg.Result, err = h.inner.GetShopProductPromotion(msg.GetArgs(ctx))
	return err
}

type GetShopProductPromotionByProductIDsQuery struct {
	ShopID     dot.ID
	ProductIDs []dot.ID

	Result []*ProductPromotion `json:"-"`
}

func (h QueryServiceHandler) HandleGetShopProductPromotionByProductIDs(ctx context.Context, msg *GetShopProductPromotionByProductIDsQuery) (err error) {
	msg.Result, err = h.inner.GetShopProductPromotionByProductIDs(msg.GetArgs(ctx))
	return err
}

type GetSupplyCommissionSettingsByProductIDsQuery struct {
	ShopID     dot.ID
	ProductIDs []dot.ID

	Result []*SupplyCommissionSetting `json:"-"`
}

func (h QueryServiceHandler) HandleGetSupplyCommissionSettingsByProductIDs(ctx context.Context, msg *GetSupplyCommissionSettingsByProductIDsQuery) (err error) {
	msg.Result, err = h.inner.GetSupplyCommissionSettingsByProductIDs(msg.GetArgs(ctx))
	return err
}

type ListShopProductPromotionsQuery struct {
	ShopID  dot.ID
	Paging  meta.Paging
	Filters meta.Filters

	Result *ListShopProductPromotionsResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListShopProductPromotions(ctx context.Context, msg *ListShopProductPromotionsQuery) (err error) {
	msg.Result, err = h.inner.ListShopProductPromotions(msg.GetArgs(ctx))
	return err
}

// implement interfaces

func (q *CreateAffiliateReferralCodeCommand) command()           {}
func (q *CreateOrUpdateCommissionSettingCommand) command()       {}
func (q *CreateOrUpdateSupplyCommissionSettingCommand) command() {}
func (q *CreateOrUpdateUserReferralCommand) command()            {}
func (q *CreateProductPromotionCommand) command()                {}
func (q *OnTradingOrderCreatedCommand) command()                 {}
func (q *OrderPaymentSuccessCommand) command()                   {}
func (q *TradingOrderCreatingCommand) command()                  {}
func (q *UpdateProductPromotionCommand) command()                {}
func (q *GetAffiliateAccountReferralByCodeQuery) query()         {}
func (q *GetAffiliateAccountReferralCodesQuery) query()          {}
func (q *GetCommissionByProductIDQuery) query()                  {}
func (q *GetCommissionByProductIDsQuery) query()                 {}
func (q *GetReferralsByReferralIDQuery) query()                  {}
func (q *GetSellerCommissionsQuery) query()                      {}
func (q *GetShopProductPromotionQuery) query()                   {}
func (q *GetShopProductPromotionByProductIDsQuery) query()       {}
func (q *GetSupplyCommissionSettingsByProductIDsQuery) query()   {}
func (q *ListShopProductPromotionsQuery) query()                 {}

// implement conversion

func (q *CreateAffiliateReferralCodeCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateReferralCodeArgs) {
	return ctx,
		&CreateReferralCodeArgs{
			AffiliateAccountID: q.AffiliateAccountID,
			Code:               q.Code,
		}
}

func (q *CreateAffiliateReferralCodeCommand) SetCreateReferralCodeArgs(args *CreateReferralCodeArgs) {
	q.AffiliateAccountID = args.AffiliateAccountID
	q.Code = args.Code
}

func (q *CreateOrUpdateCommissionSettingCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateCommissionSettingArgs) {
	return ctx,
		&CreateCommissionSettingArgs{
			ProductID:   q.ProductID,
			AccountID:   q.AccountID,
			Amount:      q.Amount,
			Unit:        q.Unit,
			Type:        q.Type,
			Description: q.Description,
			Note:        q.Note,
		}
}

func (q *CreateOrUpdateCommissionSettingCommand) SetCreateCommissionSettingArgs(args *CreateCommissionSettingArgs) {
	q.ProductID = args.ProductID
	q.AccountID = args.AccountID
	q.Amount = args.Amount
	q.Unit = args.Unit
	q.Type = args.Type
	q.Description = args.Description
	q.Note = args.Note
}

func (q *CreateOrUpdateSupplyCommissionSettingCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateOrUpdateSupplyCommissionSettingArgs) {
	return ctx,
		&CreateOrUpdateSupplyCommissionSettingArgs{
			ShopID:                   q.ShopID,
			ProductID:                q.ProductID,
			Level1DirectCommission:   q.Level1DirectCommission,
			Level1IndirectCommission: q.Level1IndirectCommission,
			Level2DirectCommission:   q.Level2DirectCommission,
			Level2IndirectCommission: q.Level2IndirectCommission,
			DependOn:                 q.DependOn,
			Level1LimitCount:         q.Level1LimitCount,
			Level1LimitDuration:      q.Level1LimitDuration,
			Level1LimitDurationType:  q.Level1LimitDurationType,
			LifetimeDuration:         q.LifetimeDuration,
			LifetimeDurationType:     q.LifetimeDurationType,
			Group:                    q.Group,
		}
}

func (q *CreateOrUpdateSupplyCommissionSettingCommand) SetCreateOrUpdateSupplyCommissionSettingArgs(args *CreateOrUpdateSupplyCommissionSettingArgs) {
	q.ShopID = args.ShopID
	q.ProductID = args.ProductID
	q.Level1DirectCommission = args.Level1DirectCommission
	q.Level1IndirectCommission = args.Level1IndirectCommission
	q.Level2DirectCommission = args.Level2DirectCommission
	q.Level2IndirectCommission = args.Level2IndirectCommission
	q.DependOn = args.DependOn
	q.Level1LimitCount = args.Level1LimitCount
	q.Level1LimitDuration = args.Level1LimitDuration
	q.Level1LimitDurationType = args.Level1LimitDurationType
	q.LifetimeDuration = args.LifetimeDuration
	q.LifetimeDurationType = args.LifetimeDurationType
	q.Group = args.Group
}

func (q *CreateOrUpdateUserReferralCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateOrUpdateReferralArgs) {
	return ctx,
		&CreateOrUpdateReferralArgs{
			UserID:           q.UserID,
			ReferralCode:     q.ReferralCode,
			SaleReferralCode: q.SaleReferralCode,
		}
}

func (q *CreateOrUpdateUserReferralCommand) SetCreateOrUpdateReferralArgs(args *CreateOrUpdateReferralArgs) {
	q.UserID = args.UserID
	q.ReferralCode = args.ReferralCode
	q.SaleReferralCode = args.SaleReferralCode
}

func (q *CreateProductPromotionCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateProductPromotionArgs) {
	return ctx,
		&CreateProductPromotionArgs{
			ShopID:      q.ShopID,
			ProductID:   q.ProductID,
			Amount:      q.Amount,
			Code:        q.Code,
			Description: q.Description,
			Unit:        q.Unit,
			Note:        q.Note,
			Type:        q.Type,
		}
}

func (q *CreateProductPromotionCommand) SetCreateProductPromotionArgs(args *CreateProductPromotionArgs) {
	q.ShopID = args.ShopID
	q.ProductID = args.ProductID
	q.Amount = args.Amount
	q.Code = args.Code
	q.Description = args.Description
	q.Unit = args.Unit
	q.Note = args.Note
	q.Type = args.Type
}

func (q *OnTradingOrderCreatedCommand) GetArgs(ctx context.Context) (_ context.Context, _ *OnTradingOrderCreatedArgs) {
	return ctx,
		&OnTradingOrderCreatedArgs{
			OrderID:      q.OrderID,
			ReferralCode: q.ReferralCode,
		}
}

func (q *OnTradingOrderCreatedCommand) SetOnTradingOrderCreatedArgs(args *OnTradingOrderCreatedArgs) {
	q.OrderID = args.OrderID
	q.ReferralCode = args.ReferralCode
}

func (q *OrderPaymentSuccessCommand) GetArgs(ctx context.Context) (_ context.Context, _ *OrderPaymentSuccessEvent) {
	return ctx,
		&OrderPaymentSuccessEvent{
			EventMeta: q.EventMeta,
			OrderID:   q.OrderID,
		}
}

func (q *OrderPaymentSuccessCommand) SetOrderPaymentSuccessEvent(args *OrderPaymentSuccessEvent) {
	q.EventMeta = args.EventMeta
	q.OrderID = args.OrderID
}

func (q *TradingOrderCreatingCommand) GetArgs(ctx context.Context) (_ context.Context, _ *TradingOrderCreating) {
	return ctx,
		&TradingOrderCreating{
			ProductIDs:   q.ProductIDs,
			ReferralCode: q.ReferralCode,
			UserID:       q.UserID,
		}
}

func (q *TradingOrderCreatingCommand) SetTradingOrderCreating(args *TradingOrderCreating) {
	q.ProductIDs = args.ProductIDs
	q.ReferralCode = args.ReferralCode
	q.UserID = args.UserID
}

func (q *UpdateProductPromotionCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateProductPromotionArgs) {
	return ctx,
		&UpdateProductPromotionArgs{
			ID:          q.ID,
			Amount:      q.Amount,
			Unit:        q.Unit,
			Code:        q.Code,
			Description: q.Description,
			Note:        q.Note,
			Type:        q.Type,
		}
}

func (q *UpdateProductPromotionCommand) SetUpdateProductPromotionArgs(args *UpdateProductPromotionArgs) {
	q.ID = args.ID
	q.Amount = args.Amount
	q.Unit = args.Unit
	q.Code = args.Code
	q.Description = args.Description
	q.Note = args.Note
	q.Type = args.Type
}

func (q *GetAffiliateAccountReferralByCodeQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetAffiliateAccountReferralByCodeArgs) {
	return ctx,
		&GetAffiliateAccountReferralByCodeArgs{
			Code: q.Code,
		}
}

func (q *GetAffiliateAccountReferralByCodeQuery) SetGetAffiliateAccountReferralByCodeArgs(args *GetAffiliateAccountReferralByCodeArgs) {
	q.Code = args.Code
}

func (q *GetAffiliateAccountReferralCodesQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetAffiliateAccountReferralCodesArgs) {
	return ctx,
		&GetAffiliateAccountReferralCodesArgs{
			AffiliateAccountID: q.AffiliateAccountID,
		}
}

func (q *GetAffiliateAccountReferralCodesQuery) SetGetAffiliateAccountReferralCodesArgs(args *GetAffiliateAccountReferralCodesArgs) {
	q.AffiliateAccountID = args.AffiliateAccountID
}

func (q *GetCommissionByProductIDQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetCommissionByProductIDArgs) {
	return ctx,
		&GetCommissionByProductIDArgs{
			AccountID: q.AccountID,
			ProductID: q.ProductID,
		}
}

func (q *GetCommissionByProductIDQuery) SetGetCommissionByProductIDArgs(args *GetCommissionByProductIDArgs) {
	q.AccountID = args.AccountID
	q.ProductID = args.ProductID
}

func (q *GetCommissionByProductIDsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetCommissionByProductIDsArgs) {
	return ctx,
		&GetCommissionByProductIDsArgs{
			AccountID:  q.AccountID,
			ProductIDs: q.ProductIDs,
		}
}

func (q *GetCommissionByProductIDsQuery) SetGetCommissionByProductIDsArgs(args *GetCommissionByProductIDsArgs) {
	q.AccountID = args.AccountID
	q.ProductIDs = args.ProductIDs
}

func (q *GetReferralsByReferralIDQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetReferralsByReferralIDArgs) {
	return ctx,
		&GetReferralsByReferralIDArgs{
			ID: q.ID,
		}
}

func (q *GetReferralsByReferralIDQuery) SetGetReferralsByReferralIDArgs(args *GetReferralsByReferralIDArgs) {
	q.ID = args.ID
}

func (q *GetSellerCommissionsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetSellerCommissionsArgs) {
	return ctx,
		&GetSellerCommissionsArgs{
			SellerID: q.SellerID,
			Paging:   q.Paging,
			Filters:  q.Filters,
		}
}

func (q *GetSellerCommissionsQuery) SetGetSellerCommissionsArgs(args *GetSellerCommissionsArgs) {
	q.SellerID = args.SellerID
	q.Paging = args.Paging
	q.Filters = args.Filters
}

func (q *GetShopProductPromotionQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetProductPromotionArgs) {
	return ctx,
		&GetProductPromotionArgs{
			ShopID:    q.ShopID,
			ProductID: q.ProductID,
		}
}

func (q *GetShopProductPromotionQuery) SetGetProductPromotionArgs(args *GetProductPromotionArgs) {
	q.ShopID = args.ShopID
	q.ProductID = args.ProductID
}

func (q *GetShopProductPromotionByProductIDsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetShopProductPromotionByProductIDs) {
	return ctx,
		&GetShopProductPromotionByProductIDs{
			ShopID:     q.ShopID,
			ProductIDs: q.ProductIDs,
		}
}

func (q *GetShopProductPromotionByProductIDsQuery) SetGetShopProductPromotionByProductIDs(args *GetShopProductPromotionByProductIDs) {
	q.ShopID = args.ShopID
	q.ProductIDs = args.ProductIDs
}

func (q *GetSupplyCommissionSettingsByProductIDsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetSupplyCommissionSettingsByProductIDsArgs) {
	return ctx,
		&GetSupplyCommissionSettingsByProductIDsArgs{
			ShopID:     q.ShopID,
			ProductIDs: q.ProductIDs,
		}
}

func (q *GetSupplyCommissionSettingsByProductIDsQuery) SetGetSupplyCommissionSettingsByProductIDsArgs(args *GetSupplyCommissionSettingsByProductIDsArgs) {
	q.ShopID = args.ShopID
	q.ProductIDs = args.ProductIDs
}

func (q *ListShopProductPromotionsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *ListShopProductPromotionsArgs) {
	return ctx,
		&ListShopProductPromotionsArgs{
			ShopID:  q.ShopID,
			Paging:  q.Paging,
			Filters: q.Filters,
		}
}

func (q *ListShopProductPromotionsQuery) SetListShopProductPromotionsArgs(args *ListShopProductPromotionsArgs) {
	q.ShopID = args.ShopID
	q.Paging = args.Paging
	q.Filters = args.Filters
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
	b.AddHandler(h.HandleCreateAffiliateReferralCode)
	b.AddHandler(h.HandleCreateOrUpdateCommissionSetting)
	b.AddHandler(h.HandleCreateOrUpdateSupplyCommissionSetting)
	b.AddHandler(h.HandleCreateOrUpdateUserReferral)
	b.AddHandler(h.HandleCreateProductPromotion)
	b.AddHandler(h.HandleOnTradingOrderCreated)
	b.AddHandler(h.HandleOrderPaymentSuccess)
	b.AddHandler(h.HandleTradingOrderCreating)
	b.AddHandler(h.HandleUpdateProductPromotion)
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
	b.AddHandler(h.HandleGetAffiliateAccountReferralByCode)
	b.AddHandler(h.HandleGetAffiliateAccountReferralCodes)
	b.AddHandler(h.HandleGetCommissionByProductID)
	b.AddHandler(h.HandleGetCommissionByProductIDs)
	b.AddHandler(h.HandleGetReferralsByReferralID)
	b.AddHandler(h.HandleGetSellerCommissions)
	b.AddHandler(h.HandleGetShopProductPromotion)
	b.AddHandler(h.HandleGetShopProductPromotionByProductIDs)
	b.AddHandler(h.HandleGetSupplyCommissionSettingsByProductIDs)
	b.AddHandler(h.HandleListShopProductPromotions)
	return QueryBus{b}
}