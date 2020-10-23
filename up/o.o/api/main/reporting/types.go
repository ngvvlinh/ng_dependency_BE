package reporting

import "time"

type ReportOrder struct {
	OrderCode     string
	CreatedAt     time.Time
	TotalItems    int
	TotalFee      int
	TotalDiscount int
	TotalAmount   int
	Revenue       int
}
