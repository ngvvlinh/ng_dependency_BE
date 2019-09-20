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
