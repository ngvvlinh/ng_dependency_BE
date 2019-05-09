package model

import (
	"time"

	cm "etop.vn/backend/pkg/common"
)

type SimpleGetOrdersByExternalIDsQuery struct {
	SourceID    int64
	SourceType  string
	ExternalIDs []string

	Result struct {
		Orders []*Order
	}
}

type GetOrderQuery struct {
	OrderID    int64
	ExternalID string
	ShopID     int64
	SupplierID int64
	PartnerID  int64
	Code       string

	// If true, don't filter order lines from other suppliers
	AllSuppliers       bool
	IncludeFulfillment bool

	Result struct {
		Order        *Order
		Fulfillments []*Fulfillment
	}
}

type GetOrdersQuery struct {
	ShopIDs    []int64 // MixedAccount
	SupplierID int64
	PartnerID  int64

	// If true, don't filter order lines from other suppliers
	AllSuppliers bool

	Paging  *cm.Paging
	Filters []cm.Filter
	IDs     []int64
	Status  *int

	Result struct {
		Orders []*Order
		Shops  []*Shop
		Total  int
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

type GetFulfillmentQuery struct {
	ShopID        int64
	SupplierID    int64
	PartnerID     int64
	FulfillmentID int64

	ShippingProvider     ShippingProvider
	ShippingCode         string
	ExternalShippingCode string

	Result *Fulfillment
}

type GetFulfillmentsQuery struct {
	ShopIDs    []int64 // MixedAccount
	SupplierID int64
	PartnerID  int64
	OrderID    int64

	Status                *Status3
	ShippingCodes         []string
	ExternalShippingCodes []string
	IDs                   []int64

	Paging  *cm.Paging
	Filters []cm.Filter

	Result struct {
		Fulfillments []*Fulfillment
		Total        int
	}
}

type GetUnCompleteFulfillmentsQuery struct {
	ShippingProviders []ShippingProvider

	Result []*Fulfillment
}

type GetFulfillmentsCallbackLogs struct {
	FromID                int64
	Paging                *cm.Paging
	ExcludeShippingStates []ShippingState

	Result struct {
		Fulfillments []*Fulfillment
	}
}

type CreateFulfillmentsCommand struct {
	Fulfillments []*Fulfillment

	Result struct {
		Fulfillments []*Fulfillment
	}
}

type UpdateFulfillmentCommand struct {
	Fulfillment              *Fulfillment
	ExternalShippingNote     *string
	ExternalShippingSubState *string
}

type UpdateFulfillmentsCommand struct {
	Fulfillments []*Fulfillment

	Result struct {
		Updated int64
	}
}

type UpdateFulfillmentsWithoutTransactionCommand struct {
	Fulfillments []*Fulfillment

	Result struct {
		Updated int
		Error   int
	}
}

type UpdateFulfillmentsStatusCommand struct {
	FulfillmentIDs  []int64
	Status          *Status4
	ShopConfirm     *Status3
	SupplierConfirm *Status3
	SyncStatus      *Status4
	ShippingState   string
}

type SyncUpdateFulfillmentsCommand struct {
	ShippingSourceID int64
	LastSyncAt       time.Time
	Fulfillments     []*Fulfillment
}

func SumOrderLineDiscount(lines []*OrderLine) int {
	sum := 0
	for _, line := range lines {
		sum += line.TotalDiscount
	}
	return sum
}

type UpdateFulfillmentsShippingStateCommand struct {
	ShopID        int64
	PartnerID     int64
	IDs           []int64
	ShippingState ShippingState

	Result struct {
		Updated int
	}
}

type UpdateOrderPaymentStatusCommand struct {
	ShopID  int64
	OrderID int64
	Status  *Status3

	Result struct {
		Updated int
	}
}

type AdminUpdateFulfillmentCommand struct {
	FulfillmentID            int64
	FullName                 string
	Phone                    string
	TotalCODAmount           *int
	IsPartialDelivery        bool
	AdminNote                string
	ActualCompensationAmount int
	ShippingState            ShippingState

	Result struct {
		Updated int
	}
}
