package model

import "time"

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenCommissionSetting(&CommissionSetting{})

type CommissionSetting struct {
	ProductID int64
	AccountID int64
	Amount    int32
	Unit      string
	Type      string
	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
}

var _ = sqlgenProductPromotion(&ProductPromotion{})

type ProductPromotion struct {
	ID          int64
	ProductID   int64
	ShopID      int64
	Amount      int32
	Unit        string
	Code        string
	Description string
	Note        string
	Type        string
	Status      int
	CreatedAt   time.Time `sq:"create"`
	UpdatedAt   time.Time `sq:"update"`
}

var _ = sqlgenAffiliateCommission(&AffiliateCommission{})

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
	CreatedAt       time.Time `sq:"create"`
	UpdatedAt       time.Time `sq:"update"`
}

var _ = sqlgenOrderCreatedNotify(&OrderCreatedNotify{})

type OrderCreatedNotify struct {
	ID           int64
	OrderID      int64
	ReferralCode string
	Status       int
	CompletedAt  time.Time
	CreatedAt    time.Time `sq:"create"`
	UpdatedAt    time.Time `sq:"update"`
}

var _ = sqlgenAffiliateReferralCode(&AffiliateReferralCode{})

type AffiliateReferralCode struct {
	ID          int64
	Code        string
	AffiliateID int64
	UserID      int64
	CreatedAt   time.Time `sq:"create"`
	UpdatedAt   time.Time `sql:"update"`
}

var _ = sqlgenUserReferral(&UserReferral{})

type UserReferral struct {
	UserID           int64
	ReferralID       int64
	ReferralCode     string
	SaleReferralID   int64
	SaleReferralCode string
	ReferralAt       time.Time
	SaleReferralAt   time.Time
	CreatedAt        time.Time `sq:"create"`
	UpdatedAt        time.Time `sq:"update"`
}
