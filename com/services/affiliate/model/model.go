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

var _ = sqlgenSellerCommission(&SellerCommission{})

type SellerCommission struct {
	ID           int64
	SellerID     int64
	FromSellerID int64
	ProductID    int64
	ShopID       int64
	SupplyID     int64
	OrderId      int64
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
	ID                       int64
	OrderID                  int64
	ShopUserID               int64
	SellerID                 int64
	ShopID                   int64
	SupplyID                 int64
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

var _ = sqlgenSupplyCommissionSetting(&SupplyCommissionSetting{})

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
	MLevel1LimitDuration     *DurationJSON
	LifetimeDuration         int64
	MLifetimeDuration        *DurationJSON
	CustomerPolicyGroupID    int64
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
	ID                   int64
	ProductID            int64
	OrderID              int64
	ProductQuantity      int32
	BaseValue            int32
	Amount               int32
	Unit                 string
	Type                 string
	OrderCreatedNotifyID int64
	Description          string
	Src                  string
	CreatedAt            time.Time `sq:"create"`
	UpdatedAt            time.Time `sq:"update"`
}

var _ = sqlgenOrderCommissionSetting(&OrderCommissionSetting{})

type OrderCommissionSetting struct {
	OrderID                  int64
	SupplyID                 int64
	ProductID                int64
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
	CustomerPolicyGroupID    int64
	CreatedAt                time.Time `sq:"create"`
	UpdatedAt                time.Time `sq:"update"`
}

var _ = sqlgenShopCashback(&ShopCashback{})

type ShopCashback struct {
	ID                   int64
	ShopID               int64
	OrderID              int64
	Amount               int32
	OrderCreatedNotifyID int64
	Description          string
	Status               int8
	ValidAt              time.Time
	CreatedAt            time.Time `sq:"create"`
	UpdatedAt            time.Time `sq:"update"`
}

var _ = sqlgenShopOrderProductHistory(&ShopOrderProductHistory{})

type ShopOrderProductHistory struct {
	UserID                int64
	ShopID                int64
	OrderID               int64
	SupplyID              int64
	ProductID             int64
	CustomerPolicyGroupID int64
	ProductQuantity       int32
	CreatedAt             time.Time `sq:"create"`
	UpdatedAt             time.Time `sq:"update"`
}

var _ = sqlgenCustomerPolicyGroup(&CustomerPolicyGroup{})

type CustomerPolicyGroup struct {
	ID        int64
	SupplyID  int64
	Name      string
	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
}
