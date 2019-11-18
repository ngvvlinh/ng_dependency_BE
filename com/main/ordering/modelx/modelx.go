package modelx

import (
	"database/sql"
	"time"

	"etop.vn/api/main/shipnow"
	ordermodel "etop.vn/backend/com/main/ordering/model"
	"etop.vn/backend/com/main/ordering/modely"
	shipmodel "etop.vn/backend/com/main/shipping/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/sq/core"
	"etop.vn/backend/pkg/etop/model"
)

type GetOrderQuery struct {
	OrderID       int64
	ExternalID    string
	ShopID        int64
	PartnerID     int64
	Code          string
	TradingShopID int64

	IncludeFulfillment bool

	Result struct {
		Order         *ordermodel.Order
		Fulfillments  []*shipmodel.Fulfillment
		XFulfillments []*Fulfillment
	}
}

type GetOrderExtendedsQuery struct {
	IDs           []int64
	ShopIDs       []int64 // MixedAccount
	PartnerID     int64
	Status        *model.Status3
	TradingShopID int64
	DateFrom      time.Time
	DateTo        time.Time

	Paging  *cm.Paging
	Filters []cm.Filter

	// When use this option, remember to always call Rows.Close()
	ResultAsRows bool

	Result struct {
		Orders []*modely.OrderExtended
		Total  int

		// only for ResultAsRows
		Rows *sql.Rows
		Opts core.Opts
	}
}

type OrderWithFulfillments struct {
	*ordermodel.Order

	Fulfillments []*Fulfillment
}

type GetOrdersQuery struct {
	ShopIDs       []int64 // MixedAccount
	PartnerID     int64
	TradingShopID int64

	Paging  *cm.Paging
	Filters []cm.Filter
	IDs     []int64
	Status  *int

	Result struct {
		Orders []OrderWithFulfillments
		Shops  []*model.Shop
		Total  int
	}
}

type SimpleGetOrdersByExternalIDsQuery struct {
	SourceID    int64
	SourceType  string
	ExternalIDs []string

	Result struct {
		Orders []*ordermodel.Order
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
	Customer        *ordermodel.OrderCustomer
	CustomerAddress *ordermodel.OrderAddress
	BillingAddress  *ordermodel.OrderAddress
	ShippingAddress *ordermodel.OrderAddress
	ShopShippingFee *int
	FeeLines        []ordermodel.OrderFeeLine
	TotalFee        *int
	ShopCOD         *int
	TotalWeight     int
	OrderNote       string
	ShopNote        string
	ShippingNote    string
	ShopShipping    *ordermodel.OrderShipping
	TryOn           model.TryOn
	OrderDiscount   *int
	Lines           []*ordermodel.OrderLine
	BasketValue     int
	TotalAmount     int
	TotalItems      int
	TotalDiscount   int

	CustomerID int64
}

type UpdateOrdersStatusCommand struct {
	ShopID    int64
	PartnerID int64
	OrderIDs  []int64

	Status        *model.Status4
	ShopConfirm   *model.Status3
	ConfirmStatus *model.Status3
	CancelReason  string

	Result struct {
		Updated int
	}
}

type CreateOrderCommand struct {
	Order *ordermodel.Order
}

type CreateOrdersCommand struct {
	ShopID int64
	Orders []*ordermodel.Order

	Result struct {
		Errors []error
	}
}

func SumOrderLineDiscount(lines []*ordermodel.OrderLine) int {
	sum := 0
	for _, line := range lines {
		sum += line.TotalDiscount
	}
	return sum
}

type UpdateOrderPaymentStatusCommand struct {
	ShopID  int64
	OrderID int64
	Status  *model.Status3

	Result struct {
		Updated int
	}
}

type UpdateOrderShippingInfoCommand struct {
	ShopID          int64
	OrderID         int64
	ShippingAddress *ordermodel.OrderAddress
	Shipping        *ordermodel.OrderShipping

	Result struct {
		Updated int
	}
}

type Fulfillment struct {
	Shipment *shipmodel.Fulfillment
	Shipnow  *shipnow.ShipnowFulfillment
}
