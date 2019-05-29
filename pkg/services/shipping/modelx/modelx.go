package modelx

import (
	"database/sql"
	"time"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/sql/core"
	"etop.vn/backend/pkg/etop/model"
	shipmodel "etop.vn/backend/pkg/services/shipping/model"
	"etop.vn/backend/pkg/services/shipping/modely"
)

type GetFulfillmentExtendedsQuery struct {
	ShopIDs       []int64 // MixedAccount
	PartnerID     int64
	OrderID       int64
	Status        *model.Status3
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
	ShopID               int64
	PartnerID            int64
	FulfillmentID        int64
	ExternalShippingCode string

	Result *modely.FulfillmentExtended
}

type GetFulfillmentQuery struct {
	ShopID        int64
	PartnerID     int64
	FulfillmentID int64

	ShippingProvider     model.ShippingProvider
	ShippingCode         string
	ExternalShippingCode string

	Result *shipmodel.Fulfillment
}

type GetFulfillmentsQuery struct {
	ShopIDs   []int64 // MixedAccount
	PartnerID int64
	OrderID   int64

	Status                *model.Status3
	ShippingCodes         []string
	ExternalShippingCodes []string
	IDs                   []int64

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
	FromID                int64
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
	ExternalShippingNote     *string
	ExternalShippingSubState *string
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
	FulfillmentIDs []int64
	Status         *model.Status4
	ShopConfirm    *model.Status3
	SyncStatus     *model.Status4
	ShippingState  string
}

type SyncUpdateFulfillmentsCommand struct {
	ShippingSourceID int64
	LastSyncAt       time.Time
	Fulfillments     []*shipmodel.Fulfillment
}

type UpdateFulfillmentsShippingStateCommand struct {
	ShopID        int64
	PartnerID     int64
	IDs           []int64
	ShippingState model.ShippingState

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
	ShippingState            model.ShippingState

	Result struct {
		Updated int
	}
}