package addressing

import (
	"context"
	"time"

	ordertypes "etop.vn/api/main/ordering/types"
	"etop.vn/api/meta"
	. "etop.vn/capi/dot"
	dot "etop.vn/capi/dot"
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

	ListAddressesByTraderID(ctx context.Context, ShopID dot.ID, TraderID dot.ID) ([]*ShopTraderAddress, error)
}

type ShopTraderAddress struct {
	ID           dot.ID
	ShopID       dot.ID
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
	Coordinates  *ordertypes.Coordinates
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// +convert:create=ShopTraderAddress
type CreateAddressArgs struct {
	ShopID       dot.ID
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
	IsDefault    NullBool
	Coordinates  *ordertypes.Coordinates
}
