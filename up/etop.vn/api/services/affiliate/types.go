package affiliate

import (
	"time"

	"etop.vn/api/meta"
)

type CommissionSetting struct {
	ProductID int64
	Amount    int32
	Unit      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ProductPromotion struct {
	ID          int64
	ProductID   int64
	Amount      int32
	Unit        string
	Code        string
	Description string
	Note        string
	Type        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type SellerCommission struct {
	ID           int64
	SellerID     int64
	FromSellerID int64
	ProductID    int64
	OrderID      int64
	ShopID       int64
	SupplyID     int64
	Amount       int32
	Description  string
	Note         string
	Type         string
	Status       int
	OValue       int32
	OBaseValue   int32
	ValidAt      time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type AffiliateReferralCode struct {
	ID          int64
	Code        string
	AffiliateID int64
	UserID      int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UserReferral struct {
	UserID           int64
	ReferralID       int64
	ReferralCode     string
	SaleReferralID   int64
	SaleReferralCode string
	ReferralAt       time.Time
	ReferralSaleAt   time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type SupplyCommissionSetting struct {
	ShopID                   int64
	ProductID                int64
	Level1DirectCommission   int32
	Level1IndirectCommission int32
	Level2DirectCommission   int32
	Level2IndirectCommission int32
	DependOn                 string
	Level1LimitCount         int32
	Level1LimitDuration      int64
	MLevel1LimitDuration     *Duration
	LifetimeDuration         int64
	MLifetimeDuration        *Duration
	Group                    string
	CreatedAt                time.Time
	UpdatedAt                time.Time
}

type Duration struct {
	Type     string
	Duration int32
}

type OrderPaymentSuccessEvent struct {
	meta.EventMeta

	OrderID int64
}
