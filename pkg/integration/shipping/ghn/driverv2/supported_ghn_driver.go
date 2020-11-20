package driverv2

import (
	"context"
	"time"
)

type SupportedGHNDriver interface {
	// update setting of shipment price list for merchant
	AddClientContract(ctx context.Context, clientID int) error

	GetPromotionCoupon(*GetPromotionCouponArgs) (coupon string, _ error)
}

type GetPromotionCouponArgs struct {
	FromProvinceCode string
	CurrentTime      time.Time
}
