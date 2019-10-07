package addressing

import (
	"context"
	"time"

	"etop.vn/api/meta"

	ordertypes "etop.vn/api/main/ordering/types"
	. "etop.vn/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateAddress(ctx context.Context, _ *CreateAddressArgs) (*ShopTraderAddress, error)

	UpdateAddress(ctx context.Context, ID int64, ShopID int64, _ *UpdateAddressArgs) (*ShopTraderAddress, error)

	DeleteAddress(ctx context.Context, ID int64, ShopID int64) (deleted int, _ error)

	SetDefaultAddress(ctx context.Context, ID, traderID, ShopID int64) (*meta.UpdatedResponse, error)
}

type QueryService interface {
	GetAddressByID(ctx context.Context, ID int64, ShopID int64) (*ShopTraderAddress, error)

	GetAddressActiveByTraderID(ctx context.Context, traderID, ShopID int64) (*ShopTraderAddress, error)

	GetAddressByTraderID(ctx context.Context, traderID, shopID int64) (*ShopTraderAddress, error)

	ListAddressesByTraderID(ctx context.Context, ShopID int64, TraderID int64) ([]*ShopTraderAddress, error)
}

type ShopTraderAddress struct {
	ID           int64
	ShopID       int64
	TraderID     int64
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

type CreateAddressArgs struct {
	ShopID       int64
	TraderID     int64
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
