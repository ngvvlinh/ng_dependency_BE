package affiliate

import (
	"time"
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

type AffiliateCommission struct {
	ID              int64
	AffiliateID     int64
	FromAffiliateID int64
	ProductID       int64
	OrderId         int64
	Value           int32
	Description     string
	Note            string
	Type            string
	Status          int
	ValidAt         time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
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
	Level1DirectCommission   int
	Level1IndirectCommission int
	Level2DirectCommission   int
	Level2IndirectCommission int
	DependOn                 string
	Level1LimitCount         int
	Level1LimitDuration      int64
	MLevel1LimitDuration     *Duration
	LifetimeDuration         int64
	MLifetimeDuration        *Duration
	CreatedAt                time.Time
	UpdatedAt                time.Time
}

type Duration struct {
	Type     string
	Duration int32
}
