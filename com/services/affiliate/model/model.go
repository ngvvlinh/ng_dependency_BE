package model

import (
	"time"

	"etop.vn/capi/dot"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenCommissionSetting(&CommissionSetting{})

type CommissionSetting struct {
	ProductID dot.ID
	AccountID dot.ID
	Amount    int32
	Unit      string
	Type      string
	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
}

var _ = sqlgenProductPromotion(&ProductPromotion{})

type ProductPromotion struct {
	ID          dot.ID
	ProductID   dot.ID
	ShopID      dot.ID
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

var _ = sqlgenSellerCommission(&SellerCommission{})

type SellerCommission struct {
	ID           dot.ID
	SellerID     dot.ID
	FromSellerID dot.ID
	ProductID    dot.ID
	ShopID       dot.ID
	SupplyID     dot.ID
	OrderId      dot.ID
	Amount       int32
	Description  string
	Note         string
	Type         string
	Status       int
	OValue       int32
	OBaseValue   int32
	ValidAt      time.Time
	CreatedAt    time.Time `sq:"create"`
	UpdatedAt    time.Time `sq:"update"`
}

var _ = sqlgenOrderCreatedNotify(&OrderCreatedNotify{})

type OrderCreatedNotify struct {
	ID                       dot.ID
	OrderID                  dot.ID
	ShopUserID               dot.ID
	SellerID                 dot.ID
	ShopID                   dot.ID
	SupplyID                 dot.ID
	ReferralCode             string
	PromotionSnapshotStatus  int32
	PromotionSnapshotErr     string
	CommissionSnapshotStatus int32
	CommissionSnapshotErr    string
	CashbackProcessStatus    int32
	CashbackProcessErr       string
	CommissionProcessStatus  int32
	CommissionProcessErr     string
	PaymentStatus            int32
	Status                   int32
	CompletedAt              time.Time
	CreatedAt                time.Time `sq:"create"`
	UpdatedAt                time.Time `sq:"update"`
}

var _ = sqlgenAffiliateReferralCode(&AffiliateReferralCode{})

type AffiliateReferralCode struct {
	ID          dot.ID
	Code        string
	AffiliateID dot.ID
	UserID      dot.ID
	CreatedAt   time.Time `sq:"create"`
	UpdatedAt   time.Time `sql:"update"`
}

var _ = sqlgenUserReferral(&UserReferral{})

type UserReferral struct {
	UserID           dot.ID
	ReferralID       dot.ID
	ReferralCode     string
	SaleReferralID   dot.ID
	SaleReferralCode string
	ReferralAt       time.Time
	SaleReferralAt   time.Time
	CreatedAt        time.Time `sq:"create"`
	UpdatedAt        time.Time `sq:"update"`
}

var _ = sqlgenSupplyCommissionSetting(&SupplyCommissionSetting{})

type SupplyCommissionSetting struct {
	ShopID                   dot.ID
	ProductID                dot.ID
	Level1DirectCommission   int32
	Level1IndirectCommission int32
	Level2DirectCommission   int32
	Level2IndirectCommission int32
	DependOn                 string
	Level1LimitCount         int32
	Level1LimitDuration      int64
	MLevel1LimitDuration     *DurationJSON
	LifetimeDuration         int64
	MLifetimeDuration        *DurationJSON
	CustomerPolicyGroupID    dot.ID
	Group                    string
	CreatedAt                time.Time `sq:"create"`
	UpdatedAt                time.Time `sq:"update"`
}

type DurationJSON struct {
	Duration int32  `json:"duration"`
	Type     string `json:"type"`
}

var _ = sqlgenOrderPromotion(&OrderPromotion{})

type OrderPromotion struct {
	ID                   dot.ID
	ProductID            dot.ID
	OrderID              dot.ID
	ProductQuantity      int32
	BaseValue            int32
	Amount               int32
	Unit                 string
	Type                 string
	OrderCreatedNotifyID dot.ID
	Description          string
	Src                  string
	CreatedAt            time.Time `sq:"create"`
	UpdatedAt            time.Time `sq:"update"`
}

var _ = sqlgenOrderCommissionSetting(&OrderCommissionSetting{})

type OrderCommissionSetting struct {
	OrderID                  dot.ID
	SupplyID                 dot.ID
	ProductID                dot.ID
	ProductQuantity          int32
	Level1DirectCommission   int32
	Level1IndirectCommission int32
	Level2DirectCommission   int32
	Level2IndirectCommission int32
	DependOn                 string
	Level1LimitCount         int32
	Level1LimitDuration      int64
	LifetimeDuration         int64
	Group                    string
	CustomerPolicyGroupID    dot.ID
	CreatedAt                time.Time `sq:"create"`
	UpdatedAt                time.Time `sq:"update"`
}

var _ = sqlgenShopCashback(&ShopCashback{})

type ShopCashback struct {
	ID                   dot.ID
	ShopID               dot.ID
	OrderID              dot.ID
	Amount               int32
	OrderCreatedNotifyID dot.ID
	Description          string
	Status               int8
	ValidAt              time.Time
	CreatedAt            time.Time `sq:"create"`
	UpdatedAt            time.Time `sq:"update"`
}

var _ = sqlgenShopOrderProductHistory(&ShopOrderProductHistory{})

type ShopOrderProductHistory struct {
	UserID                dot.ID
	ShopID                dot.ID
	OrderID               dot.ID
	SupplyID              dot.ID
	ProductID             dot.ID
	CustomerPolicyGroupID dot.ID
	ProductQuantity       int32
	CreatedAt             time.Time `sq:"create"`
	UpdatedAt             time.Time `sq:"update"`
}

var _ = sqlgenCustomerPolicyGroup(&CustomerPolicyGroup{})

type CustomerPolicyGroup struct {
	ID        dot.ID
	SupplyID  dot.ID
	Name      string
	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
}
