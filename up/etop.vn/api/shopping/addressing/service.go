package addressing

import (
	"context"
	"time"

	ordertypes "etop.vn/api/main/ordering/types"
	. "etop.vn/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateAddress(ctx context.Context, _ *CreateAddressArgs) (*ShopTraderAddress, error)

	UpdateAddress(ctx context.Context, ID int64, ShopID int64, _ *UpdateAddressArgs) (*ShopTraderAddress, error)

	DeleteAddress(ctx context.Context, ID int64, ShopID int64) (deleted int, _ error)
}

type QueryService interface {
	GetAddressByID(ctx context.Context, ID int64, ShopID int64) (*ShopTraderAddress, error)

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
	Coordinates  *ordertypes.Coordinates
}
