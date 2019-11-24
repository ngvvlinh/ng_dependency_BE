package affiliate

import (
	"time"

	"etop.vn/api/meta"
	"etop.vn/capi/dot"
)

type CommissionSetting struct {
	ProductID dot.ID
	Amount    int32
	Unit      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ProductPromotion struct {
	ID          dot.ID
	ProductID   dot.ID
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
	ID           dot.ID
	SellerID     dot.ID
	FromSellerID dot.ID
	ProductID    dot.ID
	OrderID      dot.ID
	ShopID       dot.ID
	SupplyID     dot.ID
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
	ID          dot.ID
	Code        string
	AffiliateID dot.ID
	UserID      dot.ID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UserReferral struct {
	UserID           dot.ID
	ReferralID       dot.ID
	ReferralCode     string
	SaleReferralID   dot.ID
	SaleReferralCode string
	ReferralAt       time.Time
	ReferralSaleAt   time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

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

	OrderID dot.ID
}
