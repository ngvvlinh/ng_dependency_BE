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

type ShopProduct struct {
}
