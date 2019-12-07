package model

import (
	"time"

	"etop.vn/api/top/types/etc/status4"

	"etop.vn/capi/dot"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenCommissionSetting(&CommissionSetting{})

type CommissionSetting struct {
	ProductID dot.ID
	AccountID dot.ID
	Amount    int
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
	Amount      int
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
	Amount       int
	Description  string
	Note         string
	Type         string
	Status       status4.Status
	OValue       int
	OBaseValue   int
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
	PromotionSnapshotStatus  int
	PromotionSnapshotErr     string
	CommissionSnapshotStatus int
	CommissionSnapshotErr    string
	CashbackProcessStatus    int
	CashbackProcessErr       string
	CommissionProcessStatus  int
	CommissionProcessErr     string
	PaymentStatus            int
	Status                   int
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
	Level1DirectCommission   int
	Level1IndirectCommission int
	Level2DirectCommission   int
	Level2IndirectCommission int
	DependOn                 string
	Level1LimitCount         int
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
	Duration int    `json:"duration"`
	Type     string `json:"type"`
}

var _ = sqlgenOrderPromotion(&OrderPromotion{})

type OrderPromotion struct {
	ID                   dot.ID
	ProductID            dot.ID
	OrderID              dot.ID
	ProductQuantity      int
	BaseValue            int
	Amount               int
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
	ProductQuantity          int
	Level1DirectCommission   int
	Level1IndirectCommission int
	Level2DirectCommission   int
	Level2IndirectCommission int
	DependOn                 string
	Level1LimitCount         int
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
	Amount               int
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
	ProductQuantity       int
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
