// Database model

package model

import (
	"time"

	"etop.vn/api/top/types/etc/try_on"

	"etop.vn/api/top/types/etc/status5"

	"etop.vn/api/top/types/etc/status4"

	"etop.vn/api/top/types/etc/status3"

	ordermodel "etop.vn/backend/com/main/ordering/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
	"etop.vn/common/xerrors"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenShipnowFulfillment(&ShipnowFulfillment{})

type Carrier string

const (
	Ahamove Carrier = "ahamove"
)

func (c Carrier) ToString() string {
	return string(c)
}

type ShipnowFulfillment struct {
	ID dot.ID

	ShopID    dot.ID
	PartnerID dot.ID
	OrderIDs  []dot.ID

	PickupAddress *ordermodel.OrderAddress

	Carrier Carrier

	ShippingServiceCode        string
	ShippingServiceFee         int
	ShippingServiceName        string
	ShippingServiceDescription string

	ChargeableWeight int
	GrossWeight      int
	BasketValue      int
	CODAmount        int
	ShippingNote     string
	RequestPickupAt  time.Time

	DeliveryPoints []*DeliveryPoint
	CancelReason   string

	Status            status5.Status
	ConfirmStatus     status3.Status
	ShippingStatus    status5.Status
	EtopPaymentStatus status4.Status

	ShippingState        string
	ShippingCode         string
	FeeLines             []*model.ShippingFeeLine
	CarrierFeeLines      []*model.ShippingFeeLine
	TotalFee             int
	ShippingCreatedAt    time.Time
	ShippingPickingAt    time.Time
	ShippingDeliveringAt time.Time
	ShippingDeliveredAt  time.Time
	ShippingCancelledAt  time.Time

	SyncStatus          status4.Status
	SyncStates          *model.FulfillmentSyncStates
	CreatedAt           time.Time `sq:"create"`
	UpdatedAt           time.Time `sq:"update"`
	CODEtopTransferedAt time.Time
	ShippingSharedLink  string
}

func (m *ShipnowFulfillment) Validate() error {
	var errs xerrors.Errors
	if m.ChargeableWeight <= 0 && m.GrossWeight <= 0 {
		err := cm.Errorf(cm.InvalidArgument, nil, "tổng khối lượng tính phí không hợp lệ")
		errs = append(errs, err)
	}

	// TODO
	return errs.ToError()
}

type DeliveryPoint struct {
	ShippingAddress *ordermodel.OrderAddress `json:"shipping_address"`
	Items           []*ordermodel.OrderLine  `json:"items"`

	OrderID          dot.ID           `json:"order_id"`
	OrderCode        string           `json:"order_code"`
	GrossWeight      int              `json:"gross_weight"`
	ChargeableWeight int              `json:"chargeable_weight"`
	Length           int              `json:"lenght"`
	Width            int              `json:"width"`
	Height           int              `json:"height"`
	BasketValue      int              `json:"basket_value"`
	CODAmount        int              `json:"cod_amount"`
	TryOn            try_on.TryOnCode `json:"try_on"`
	ShippingNote     string           `json:"shipping_note"`
}
