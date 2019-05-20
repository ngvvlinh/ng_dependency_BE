package model

type SimpleGetOrdersByExternalIDsQuery struct {
	SourceID    int64
	SourceType  string
	ExternalIDs []string

	Result struct {
		Orders []*Order
	}
}

type VerifyOrdersByEdCodeQuery struct {
	ShopID           int64
	EdCodes          []string
	OnlyActiveOrders bool // shop_confirm is not -1

	Result struct {
		EdCodes []string
	}
}

type UpdateOrderCommand struct {
	ID        int64
	ShopID    int64
	PartnerID int64

	ShopName        string
	Customer        *OrderCustomer
	CustomerAddress *OrderAddress
	BillingAddress  *OrderAddress
	ShippingAddress *OrderAddress
	ShopShippingFee *int
	FeeLines        []OrderFeeLine
	TotalFee        *int
	ShopCOD         *int
	TotalWeight     int
	OrderNote       string
	ShopNote        string
	ShippingNote    string
	ShopShipping    *OrderShipping
	TryOn           TryOn
	OrderDiscount   *int
	Lines           []*OrderLine
	BasketValue     int
	TotalAmount     int
	TotalItems      int
	TotalDiscount   int
}

type UpdateOrdersStatusCommand struct {
	ShopID    int64
	PartnerID int64
	OrderIDs  []int64

	Status        *Status4
	ShopConfirm   *Status3
	ConfirmStatus *Status3
	CancelReason  string

	Result struct {
		Updated int
	}
}

type UpdateOrderLinesStatusCommand struct {
	SupplierID int64
	Updates    []UpdateOrderLinesStatus

	Result struct {
		Updated int
	}
}

type UpdateOrderLinesStatus struct {
	OrderID         int64
	ProductIDs      []int64
	SupplierConfirm *Status3
	CancelReason    string
}

type CreateOrderCommand struct {
	Order *Order
}

type CreateOrdersCommand struct {
	ShopID int64
	Orders []*Order

	Result struct {
		Errors []error
	}
}

func SumOrderLineDiscount(lines []*OrderLine) int {
	sum := 0
	for _, line := range lines {
		sum += line.TotalDiscount
	}
	return sum
}

type UpdateOrderPaymentStatusCommand struct {
	ShopID  int64
	OrderID int64
	Status  *Status3

	Result struct {
		Updated int
	}
}
