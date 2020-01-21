package modelx

import (
	"database/sql"
	"time"

	"etop.vn/api/main/shipnow"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/api/top/types/etc/status4"
	"etop.vn/api/top/types/etc/status5"
	"etop.vn/api/top/types/etc/try_on"
	identitymodel "etop.vn/backend/com/main/identity/model"
	ordermodel "etop.vn/backend/com/main/ordering/model"
	"etop.vn/backend/com/main/ordering/modely"
	shipmodel "etop.vn/backend/com/main/shipping/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/sql/sq/core"
	"etop.vn/capi/dot"
)

type GetOrderQuery struct {
	OrderID       dot.ID
	ExternalID    string
	ShopID        dot.ID
	PartnerID     dot.ID
	Code          string
	TradingShopID dot.ID

	IncludeFulfillment bool

	Result struct {
		Order         *ordermodel.Order
		Fulfillments  []*shipmodel.Fulfillment
		XFulfillments []*Fulfillment
	}
}

type GetOrderExtendedsQuery struct {
	IDs           []dot.ID
	ShopIDs       []dot.ID // MixedAccount
	PartnerID     dot.ID
	Status        status3.NullStatus
	TradingShopID dot.ID
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
	ShopIDs       []dot.ID // MixedAccount
	PartnerID     dot.ID
	TradingShopID dot.ID

	Paging  *cm.Paging
	Filters []cm.Filter
	IDs     []dot.ID
	Status  dot.NullInt

	Result struct {
		Orders []OrderWithFulfillments
		Shops  []*identitymodel.Shop
	}
}

type SimpleGetOrdersByExternalIDsQuery struct {
	SourceID    dot.ID
	SourceType  string
	ExternalIDs []string

	Result struct {
		Orders []*ordermodel.Order
	}
}

type VerifyOrdersByEdCodeQuery struct {
	ShopID           dot.ID
	EdCodes          []string
	OnlyActiveOrders bool // shop_confirm is not -1

	Result struct {
		EdCodes []string
	}
}

type UpdateOrderCommand struct {
	ID        dot.ID
	ShopID    dot.ID
	PartnerID dot.ID

	ShopName        string
	Customer        *ordermodel.OrderCustomer
	CustomerAddress *ordermodel.OrderAddress
	BillingAddress  *ordermodel.OrderAddress
	ShippingAddress *ordermodel.OrderAddress
	ShopShippingFee dot.NullInt
	FeeLines        []ordermodel.OrderFeeLine
	TotalFee        dot.NullInt
	ShopCOD         dot.NullInt
	TotalWeight     int
	OrderNote       string
	ShopNote        string
	ShippingNote    string
	ShopShipping    *ordermodel.OrderShipping
	TryOn           try_on.TryOnCode
	OrderDiscount   dot.NullInt
	Lines           []*ordermodel.OrderLine
	BasketValue     int
	TotalAmount     int
	TotalItems      int
	TotalDiscount   int

	CustomerID dot.ID
}

type UpdateOrdersStatusCommand struct {
	ShopID    dot.ID
	PartnerID dot.ID
	OrderIDs  []dot.ID

	Status        status5.NullStatus
	ShopConfirm   status3.NullStatus
	ConfirmStatus status3.NullStatus
	CancelReason  string

	Result struct {
		Updated int
	}
}

type CreateOrderCommand struct {
	Order *ordermodel.Order
}

type CreateOrdersCommand struct {
	ShopID dot.ID
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
	ShopID        dot.ID
	OrderID       dot.ID
	PaymentStatus status4.NullStatus

	Result struct {
		Updated int
	}
}

type UpdateOrderShippingInfoCommand struct {
	ShopID          dot.ID
	OrderID         dot.ID
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
