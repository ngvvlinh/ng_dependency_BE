// +build !generator

// Code generated by generator api. DO NOT EDIT.

package addressing

import (
	context "context"

	orderingtypes "o.o/api/main/ordering/types"
	meta "o.o/api/meta"
	capi "o.o/capi"
	dot "o.o/capi/dot"
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

type CreateAddressCommand struct {
	ShopID       dot.ID
	PartnerID    dot.ID
	TraderID     dot.ID
	FullName     string
	Phone        string
	Email        string
	Company      string
	Address1     string
	Address2     string
	DistrictCode string
	WardCode     string
	Position     string
	IsDefault    bool
	Coordinates  *orderingtypes.Coordinates

	Result *ShopTraderAddress `json:"-"`
}

func (h AggregateHandler) HandleCreateAddress(ctx context.Context, msg *CreateAddressCommand) (err error) {
	msg.Result, err = h.inner.CreateAddress(msg.GetArgs(ctx))
	return err
}

type DeleteAddressCommand struct {
	ID     dot.ID
	ShopID dot.ID

	Result int `json:"-"`
}

func (h AggregateHandler) HandleDeleteAddress(ctx context.Context, msg *DeleteAddressCommand) (err error) {
	msg.Result, err = h.inner.DeleteAddress(msg.GetArgs(ctx))
	return err
}

type SetDefaultAddressCommand struct {
	ID       dot.ID
	TraderID dot.ID
	ShopID   dot.ID

	Result *meta.UpdatedResponse `json:"-"`
}

func (h AggregateHandler) HandleSetDefaultAddress(ctx context.Context, msg *SetDefaultAddressCommand) (err error) {
	msg.Result, err = h.inner.SetDefaultAddress(msg.GetArgs(ctx))
	return err
}

type UpdateAddressCommand struct {
	ID           dot.ID
	ShopID       dot.ID
	FullName     dot.NullString
	Phone        dot.NullString
	Email        dot.NullString
	Company      dot.NullString
	Address1     dot.NullString
	Address2     dot.NullString
	DistrictCode dot.NullString
	WardCode     dot.NullString
	Position     dot.NullString
	IsDefault    dot.NullBool
	Coordinates  *orderingtypes.Coordinates

	Result *ShopTraderAddress `json:"-"`
}

func (h AggregateHandler) HandleUpdateAddress(ctx context.Context, msg *UpdateAddressCommand) (err error) {
	msg.Result, err = h.inner.UpdateAddress(msg.GetArgs(ctx))
	return err
}

type GetAddressActiveByTraderIDQuery struct {
	TraderID dot.ID
	ShopID   dot.ID

	Result *ShopTraderAddress `json:"-"`
}

func (h QueryServiceHandler) HandleGetAddressActiveByTraderID(ctx context.Context, msg *GetAddressActiveByTraderIDQuery) (err error) {
	msg.Result, err = h.inner.GetAddressActiveByTraderID(msg.GetArgs(ctx))
	return err
}

type GetAddressByIDQuery struct {
	ID     dot.ID
	ShopID dot.ID

	Result *ShopTraderAddress `json:"-"`
}

func (h QueryServiceHandler) HandleGetAddressByID(ctx context.Context, msg *GetAddressByIDQuery) (err error) {
	msg.Result, err = h.inner.GetAddressByID(msg.GetArgs(ctx))
	return err
}

type GetAddressByTraderIDQuery struct {
	TraderID dot.ID
	ShopID   dot.ID

	Result *ShopTraderAddress `json:"-"`
}

func (h QueryServiceHandler) HandleGetAddressByTraderID(ctx context.Context, msg *GetAddressByTraderIDQuery) (err error) {
	msg.Result, err = h.inner.GetAddressByTraderID(msg.GetArgs(ctx))
	return err
}

type ListAddressesQuery struct {
	ShopID   dot.ID
	TraderID dot.ID
	Phone    string
	Paging   meta.Paging

	Result *ShopTraderAddressesResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListAddresses(ctx context.Context, msg *ListAddressesQuery) (err error) {
	msg.Result, err = h.inner.ListAddresses(msg.GetArgs(ctx))
	return err
}

type ListAddressesByTraderIDQuery struct {
	ShopID   dot.ID
	TraderID dot.ID
	Phone    string
	Paging   meta.Paging

	Result *ShopTraderAddressesResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListAddressesByTraderID(ctx context.Context, msg *ListAddressesByTraderIDQuery) (err error) {
	msg.Result, err = h.inner.ListAddressesByTraderID(msg.GetArgs(ctx))
	return err
}

type ListAddressesByTraderIDsQuery struct {
	ShopID         dot.ID
	TraderIDs      []dot.ID
	Phone          string
	Paging         meta.Paging
	IncludeDeleted bool

	Result *ShopTraderAddressesResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListAddressesByTraderIDs(ctx context.Context, msg *ListAddressesByTraderIDsQuery) (err error) {
	msg.Result, err = h.inner.ListAddressesByTraderIDs(msg.GetArgs(ctx))
	return err
}

// implement interfaces

func (q *CreateAddressCommand) command()     {}
func (q *DeleteAddressCommand) command()     {}
func (q *SetDefaultAddressCommand) command() {}
func (q *UpdateAddressCommand) command()     {}

func (q *GetAddressActiveByTraderIDQuery) query() {}
func (q *GetAddressByIDQuery) query()             {}
func (q *GetAddressByTraderIDQuery) query()       {}
func (q *ListAddressesQuery) query()              {}
func (q *ListAddressesByTraderIDQuery) query()    {}
func (q *ListAddressesByTraderIDsQuery) query()   {}

// implement conversion

func (q *CreateAddressCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateAddressArgs) {
	return ctx,
		&CreateAddressArgs{
			ShopID:       q.ShopID,
			PartnerID:    q.PartnerID,
			TraderID:     q.TraderID,
			FullName:     q.FullName,
			Phone:        q.Phone,
			Email:        q.Email,
			Company:      q.Company,
			Address1:     q.Address1,
			Address2:     q.Address2,
			DistrictCode: q.DistrictCode,
			WardCode:     q.WardCode,
			Position:     q.Position,
			IsDefault:    q.IsDefault,
			Coordinates:  q.Coordinates,
		}
}

func (q *CreateAddressCommand) SetCreateAddressArgs(args *CreateAddressArgs) {
	q.ShopID = args.ShopID
	q.PartnerID = args.PartnerID
	q.TraderID = args.TraderID
	q.FullName = args.FullName
	q.Phone = args.Phone
	q.Email = args.Email
	q.Company = args.Company
	q.Address1 = args.Address1
	q.Address2 = args.Address2
	q.DistrictCode = args.DistrictCode
	q.WardCode = args.WardCode
	q.Position = args.Position
	q.IsDefault = args.IsDefault
	q.Coordinates = args.Coordinates
}

func (q *DeleteAddressCommand) GetArgs(ctx context.Context) (_ context.Context, ID dot.ID, ShopID dot.ID) {
	return ctx,
		q.ID,
		q.ShopID
}

func (q *SetDefaultAddressCommand) GetArgs(ctx context.Context) (_ context.Context, ID dot.ID, traderID dot.ID, ShopID dot.ID) {
	return ctx,
		q.ID,
		q.TraderID,
		q.ShopID
}

func (q *UpdateAddressCommand) GetArgs(ctx context.Context) (_ context.Context, ID dot.ID, ShopID dot.ID, _ *UpdateAddressArgs) {
	return ctx,
		q.ID,
		q.ShopID,
		&UpdateAddressArgs{
			FullName:     q.FullName,
			Phone:        q.Phone,
			Email:        q.Email,
			Company:      q.Company,
			Address1:     q.Address1,
			Address2:     q.Address2,
			DistrictCode: q.DistrictCode,
			WardCode:     q.WardCode,
			Position:     q.Position,
			IsDefault:    q.IsDefault,
			Coordinates:  q.Coordinates,
		}
}

func (q *UpdateAddressCommand) SetUpdateAddressArgs(args *UpdateAddressArgs) {
	q.FullName = args.FullName
	q.Phone = args.Phone
	q.Email = args.Email
	q.Company = args.Company
	q.Address1 = args.Address1
	q.Address2 = args.Address2
	q.DistrictCode = args.DistrictCode
	q.WardCode = args.WardCode
	q.Position = args.Position
	q.IsDefault = args.IsDefault
	q.Coordinates = args.Coordinates
}

func (q *GetAddressActiveByTraderIDQuery) GetArgs(ctx context.Context) (_ context.Context, traderID dot.ID, ShopID dot.ID) {
	return ctx,
		q.TraderID,
		q.ShopID
}

func (q *GetAddressByIDQuery) GetArgs(ctx context.Context) (_ context.Context, ID dot.ID, ShopID dot.ID) {
	return ctx,
		q.ID,
		q.ShopID
}

func (q *GetAddressByTraderIDQuery) GetArgs(ctx context.Context) (_ context.Context, traderID dot.ID, shopID dot.ID) {
	return ctx,
		q.TraderID,
		q.ShopID
}

func (q *ListAddressesQuery) GetArgs(ctx context.Context) (_ context.Context, _ *ListAddressesArgs) {
	return ctx,
		&ListAddressesArgs{
			ShopID:   q.ShopID,
			TraderID: q.TraderID,
			Phone:    q.Phone,
			Paging:   q.Paging,
		}
}

func (q *ListAddressesQuery) SetListAddressesArgs(args *ListAddressesArgs) {
	q.ShopID = args.ShopID
	q.TraderID = args.TraderID
	q.Phone = args.Phone
	q.Paging = args.Paging
}

func (q *ListAddressesByTraderIDQuery) GetArgs(ctx context.Context) (_ context.Context, _ *ListAddressesByTraderIDArgs) {
	return ctx,
		&ListAddressesByTraderIDArgs{
			ShopID:   q.ShopID,
			TraderID: q.TraderID,
			Phone:    q.Phone,
			Paging:   q.Paging,
		}
}

func (q *ListAddressesByTraderIDQuery) SetListAddressesByTraderIDArgs(args *ListAddressesByTraderIDArgs) {
	q.ShopID = args.ShopID
	q.TraderID = args.TraderID
	q.Phone = args.Phone
	q.Paging = args.Paging
}

func (q *ListAddressesByTraderIDsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *ListAddressesByTraderIDsArgs) {
	return ctx,
		&ListAddressesByTraderIDsArgs{
			ShopID:         q.ShopID,
			TraderIDs:      q.TraderIDs,
			Phone:          q.Phone,
			Paging:         q.Paging,
			IncludeDeleted: q.IncludeDeleted,
		}
}

func (q *ListAddressesByTraderIDsQuery) SetListAddressesByTraderIDsArgs(args *ListAddressesByTraderIDsArgs) {
	q.ShopID = args.ShopID
	q.TraderIDs = args.TraderIDs
	q.Phone = args.Phone
	q.Paging = args.Paging
	q.IncludeDeleted = args.IncludeDeleted
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
	b.AddHandler(h.HandleCreateAddress)
	b.AddHandler(h.HandleDeleteAddress)
	b.AddHandler(h.HandleSetDefaultAddress)
	b.AddHandler(h.HandleUpdateAddress)
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
	b.AddHandler(h.HandleGetAddressActiveByTraderID)
	b.AddHandler(h.HandleGetAddressByID)
	b.AddHandler(h.HandleGetAddressByTraderID)
	b.AddHandler(h.HandleListAddresses)
	b.AddHandler(h.HandleListAddressesByTraderID)
	b.AddHandler(h.HandleListAddressesByTraderIDs)
	return QueryBus{b}
}
