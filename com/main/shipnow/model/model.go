// Database model

package model

import (
	"time"

	ordermodel "etop.vn/backend/com/main/ordering/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/etop/model"
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
	ID int64

	ShopID    int64
	PartnerID int64
	OrderIDs  []int64

	PickupAddress *ordermodel.OrderAddress

	Carrier Carrier

	ShippingServiceCode        string
	ShippingServiceFee         int32
	ShippingServiceName        string
	ShippingServiceDescription string

	ChargeableWeight int32
	GrossWeight      int32
	BasketValue      int32
	CODAmount        int32
	ShippingNote     string
	RequestPickupAt  time.Time

	DeliveryPoints []*DeliveryPoint
	CancelReason   string

	Status            model.Status5
	ConfirmStatus     model.Status3
	ShippingStatus    model.Status5
	EtopPaymentStatus model.Status4

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

	SyncStatus          model.Status4
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

	OrderID          int64       `json:"order_id"`
	OrderCode        string      `json:"order_code"`
	GrossWeight      int32       `json:"gross_weight"`
	ChargeableWeight int32       `json:"chargeable_weight"`
	Length           int32       `json:"lenght"`
	Width            int32       `json:"width"`
	Height           int32       `json:"height"`
	BasketValue      int32       `json:"basket_value"`
	CODAmount        int32       `json:"cod_amount"`
	TryOn            model.TryOn `json:"try_on"`
	ShippingNote     string      `json:"shipping_note"`
}
