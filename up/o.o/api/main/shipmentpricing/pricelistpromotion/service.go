package pricelistpromotion

import (
	"context"
	"time"

	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreatePriceListPromotion(context.Context, *CreatePriceListPromotionArgs) (*ShipmentPriceListPromotion, error)

	UpdatePriceListPromotion(context.Context, *UpdatePriceListPromotionArgs) error

	DeletePriceListPromotion(ctx context.Context, ID dot.ID) error
}

type QueryService interface {
	GetPriceListPromotion(ctx context.Context, ID dot.ID) (*ShipmentPriceListPromotion, error)

	ListPriceListPromotion(context.Context, *ListPriceListPromotionArgs) ([]*ShipmentPriceListPromotion, error)

	GetValidPriceListPromotion(context.Context, *GetValidPriceListPromotionArgs) (*ShipmentPriceListPromotion, error)
}

// +convert:create=ShipmentPriceListPromotion
type CreatePriceListPromotionArgs struct {
	PriceListID   dot.ID
	Name          string
	Description   string
	ConnectionID  dot.ID
	DateFrom      time.Time
	DateTo        time.Time
	AppliedRules  *AppliedRules
	PriorityPoint int
}

// +convert:update=ShipmentPriceListPromotion
type UpdatePriceListPromotionArgs struct {
	ID            dot.ID
	Name          string
	Description   string
	DateFrom      time.Time
	DateTo        time.Time
	AppliedRules  *AppliedRules
	PriorityPoint int
	Status        status3.NullStatus
}

type ListPriceListPromotionArgs struct {
	ConnectionID dot.ID
}

type GetValidPriceListPromotionArgs struct {
	ShopID           dot.ID
	FromProvinceCode string
	ConnectionID     dot.ID
}
