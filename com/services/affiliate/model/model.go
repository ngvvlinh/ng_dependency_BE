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

var _ = sqlgenCommission(&Commission{})

type Commission struct {
	ID          int64
	AffiliateID int64
	Value       int32
	Unit        string
	description string
	note        string
	OrderID     int64
	Status      int
	Type        string
	CreatedAt   time.Time `sq:"create"`
	UpdatedAt   time.Time `sq:"update"`
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
