package modelx

import (
	"database/sql"
	"time"

	"o.o/api/top/types/etc/shipping"
	"o.o/api/top/types/etc/shipping_provider"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/status4"
	shipmodel "o.o/backend/com/main/shipping/model"
	"o.o/backend/com/main/shipping/modely"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/sql/sq/core"
	"o.o/capi/dot"
)

type GetFulfillmentExtendedsQuery struct {
	IDs           []dot.ID
	ShopIDs       []dot.ID // MixedAccount
	PartnerID     dot.ID
	OrderID       dot.ID
	Status        status3.NullStatus
	ShippingCodes []string
	DateFrom      time.Time
	DateTo        time.Time
	ConnectionIDs []dot.ID

	Paging  *cm.Paging
	Filters []cm.Filter

	// When use this option, remember to always call Rows.Close()
	ResultAsRows bool

	Result struct {
		Fulfillments []*modely.FulfillmentExtended
		Total        int
		// only for ResultAsRows
		Rows *sql.Rows
		Opts core.Opts
	}
}

type GetFulfillmentExtendedQuery struct {
	ShopID               dot.ID
	PartnerID            dot.ID
	FulfillmentID        dot.ID
	ExternalShippingCode string

	Result *modely.FulfillmentExtended
}

type GetFulfillmentQuery struct {
	ShopID        dot.ID
	PartnerID     dot.ID
	FulfillmentID dot.ID

	ShippingProvider     shipping_provider.ShippingProvider
	ShippingCode         string
	ExternalShippingCode string

	Result *shipmodel.Fulfillment
}

type GetFulfillmentsQuery struct {
	ShopIDs   []dot.ID // MixedAccount
	PartnerID dot.ID
	OrderID   dot.ID

	Status                status3.NullStatus
	ShippingCodes         []string
	ExternalShippingCodes []string
	IDs                   []dot.ID

	Paging  *cm.Paging
	Filters []cm.Filter

	Result struct {
		Fulfillments []*shipmodel.Fulfillment
		Total        int
	}
}

type GetUnCompleteFulfillmentsQuery struct {
	ShippingProviders []shipping_provider.ShippingProvider

	Result []*shipmodel.Fulfillment
}

type GetFulfillmentsCallbackLogs struct {
	FromID                dot.ID
	Paging                *cm.Paging
	ExcludeShippingStates []shipping.State

	Result struct {
		Fulfillments []*shipmodel.Fulfillment
	}
}

type CreateFulfillmentsCommand struct {
	Fulfillments []*shipmodel.Fulfillment

	Result struct {
		Fulfillments []*shipmodel.Fulfillment
	}
}

type UpdateFulfillmentCommand struct {
	Fulfillment              *shipmodel.Fulfillment
	ExternalShippingNote     dot.NullString
	ExternalShippingSubState dot.NullString
}

type UpdateFulfillmentsCommand struct {
	Fulfillments []*shipmodel.Fulfillment

	Result struct {
		Updated int64
	}
}

type UpdateFulfillmentsWithoutTransactionCommand struct {
	Fulfillments []*shipmodel.Fulfillment

	Result struct {
		Updated int
		Error   int
	}
}

type UpdateFulfillmentsStatusCommand struct {
	FulfillmentIDs []dot.ID
	Status         status4.NullStatus
	ShopConfirm    status3.NullStatus
	SyncStatus     status4.NullStatus
	ShippingState  string
}

type SyncUpdateFulfillmentsCommand struct {
	ShippingSourceID dot.ID
	LastSyncAt       time.Time
	Fulfillments     []*shipmodel.Fulfillment
}

type UpdateFulfillmentsShippingStateCommand struct {
	ShopID        dot.ID
	PartnerID     dot.ID
	IDs           []dot.ID
	ShippingState shipping.State

	Result struct {
		Updated int
	}
}
