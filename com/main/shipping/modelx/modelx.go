package modelx

import (
	"database/sql"
	"time"

	"etop.vn/api/top/types/etc/status3"
	"etop.vn/api/top/types/etc/status4"
	shipmodel "etop.vn/backend/com/main/shipping/model"
	"etop.vn/backend/com/main/shipping/modely"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/sq/core"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
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

	ShippingProvider     model.ShippingProvider
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
	ShippingProviders []model.ShippingProvider

	Result []*shipmodel.Fulfillment
}

type GetFulfillmentsCallbackLogs struct {
	FromID                dot.ID
	Paging                *cm.Paging
	ExcludeShippingStates []model.ShippingState

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
	ShippingState model.ShippingState

	Result struct {
		Updated int
	}
}

type AdminUpdateFulfillmentCommand struct {
	FulfillmentID            dot.ID
	FullName                 string
	Phone                    string
	TotalCODAmount           dot.NullInt
	IsPartialDelivery        bool
	AdminNote                string
	ActualCompensationAmount int
	ShippingState            model.ShippingState

	Result struct {
		Updated int
	}
}
