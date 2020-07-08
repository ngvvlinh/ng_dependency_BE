package addressing

import (
	"context"
	"time"

	ordertypes "o.o/api/main/ordering/types"
	"o.o/api/meta"
	. "o.o/capi/dot"
	dot "o.o/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateAddress(ctx context.Context, _ *CreateAddressArgs) (*ShopTraderAddress, error)

	UpdateAddress(ctx context.Context, ID dot.ID, ShopID dot.ID, _ *UpdateAddressArgs) (*ShopTraderAddress, error)

	DeleteAddress(ctx context.Context, ID dot.ID, ShopID dot.ID) (deleted int, _ error)

	SetDefaultAddress(ctx context.Context, ID, traderID, ShopID dot.ID) (*meta.UpdatedResponse, error)
}

type QueryService interface {
	GetAddressByID(ctx context.Context, ID dot.ID, ShopID dot.ID) (*ShopTraderAddress, error)

	GetAddressActiveByTraderID(ctx context.Context, traderID, ShopID dot.ID) (*ShopTraderAddress, error)

	GetAddressByTraderID(ctx context.Context, traderID, shopID dot.ID) (*ShopTraderAddress, error)

	ListAddressesByTraderID(ctx context.Context, _ *ListAddressesByTraderIDArgs) (*ShopTraderAddressesResponse, error)

	ListAddressesByTraderIDs(ctx context.Context, _ *ListAddressesByTraderIDsArgs) (*ShopTraderAddressesResponse, error)

	ListAddresses(ctx context.Context, _ *ListAddressesArgs) (*ShopTraderAddressesResponse, error)
}
type ListAddressesArgs struct {
	ShopID   dot.ID
	TraderID dot.ID
	Phone    string

	Paging meta.Paging
}

type ShopTraderAddress struct {
	ID           dot.ID
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
	IsDefault    bool
	Position     string
	Coordinates  *ordertypes.Coordinates
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Deleted      bool
}

type ShopTraderAddressesResponse struct {
	ShopTraderAddresses []*ShopTraderAddress
	Count               int
	Paging              meta.PageInfo
}

// +convert:create=ShopTraderAddress
type CreateAddressArgs struct {
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
	Coordinates  *ordertypes.Coordinates
}

// +convert:update=ShopTraderAddress(ID)
type UpdateAddressArgs struct {
	FullName     NullString
	Phone        NullString
	Email        NullString
	Company      NullString
	Address1     NullString
	Address2     NullString
	DistrictCode NullString
	WardCode     NullString
	Position     NullString
	IsDefault    NullBool
	Coordinates  *ordertypes.Coordinates
}

type ListAddressesByTraderIDArgs struct {
	ShopID   dot.ID
	TraderID dot.ID
	Phone    string

	Paging meta.Paging
}

type ListAddressesByTraderIDsArgs struct {
	ShopID    dot.ID
	TraderIDs []dot.ID
	Phone     string

	Paging meta.Paging

	IncludeDeleted bool
}
