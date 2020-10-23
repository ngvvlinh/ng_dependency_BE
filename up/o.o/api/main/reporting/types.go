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

// notion: https://www.notion.so/Report-T-i-ch-nh-B-o-c-o-k-t-qu-ho-t-ng-kinh-doanh-e5e8c48ef3d9497db95d8696d9a5e462
type ReportIncomeStatement struct {
	Revenue         int // Doanh thu bán hàng (1)
	Discounts       int // Giảm trừ Doanh thu (2)
	NetRevenue      int // Doanh thu thuần (3=1-2)
	CostPrice       int // Giá vốn hàng bán (4)
	GrossProfit     int // Lợi nhuận gộp về bán hàng (5=3-4)
	Expenses        int // Chi phí (6 = 6.1 + 6.2 + 6.3)
	ShippingFee     int // Phí giao hàng (6.1)
	Discards        int // Xuất hủy hàng hóa  (6.2)
	Others          int // Khác (6.3)
	ProfitStatement int // Lợi nhuận từ hoạt động kinh doanh (7=5-6)
	OtherIncomes    int // Thu nhập khác (8)
	OtherExpenses   int // Chi phí khác (9)
	NetProfit       int // Lợi nhuận thuần (10=(7+8)-9)
}
