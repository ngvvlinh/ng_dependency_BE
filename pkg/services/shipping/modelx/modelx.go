package modelx

import (
	"database/sql"
	"time"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/sql/core"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/services/shipping/modely"
)

type GetFulfillmentExtendedsQuery struct {
	ShopIDs       []int64 // MixedAccount
	SupplierID    int64
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
	SupplierID           int64
	PartnerID            int64
	FulfillmentID        int64
	ExternalShippingCode string

	Result *modely.FulfillmentExtended
}
