// Database model

package model

import (
	"time"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/etop/model"
	ordermodel "etop.vn/backend/pkg/services/ordering/model"
)

//go:generate ../../../../scripts/derive.sh

var _ = sqlgenShipnowFulfillment(&ShipnowFulfillment{})

type ShipnowFulfillment struct {
	ID int64

	ShopID    int64
	PartnerID int64
	OrderIDs  []int64

	PickupAddress *model.Address

	Carrier model.ShippingProvider

	ShippingServiceCode string
	ShippingServiceFee  int32

	ChargeableWeight int32
	GrossWeight      int32
	BasketValue      int32
	CODAmount        int32
	ShippingNote     string
	RequestPickupAt  time.Time

	DeliveryPoints []*DeliveryPoint
}

func (m *ShipnowFulfillment) Validate() error {
	var errs cm.Errors
	if m.ChargeableWeight <= 0 {
		err := cm.Errorf(cm.InvalidArgument, nil, "tổng khối lượng tính phí không hợp lệ")
		errs = append(errs, err)
	}
	if m.GrossWeight <= 0 {
		err := cm.Errorf(cm.InvalidArgument, nil, "tổng khối lượng không hợp lệ")
		errs = append(errs, err)
	}

	// TODO
	return errs.ToError()
}

type DeliveryPoint struct {
	ShippingAddress *model.Address
	Items           []*ordermodel.OrderLine

	GrossWeight      int32
	ChargeableWeight int32
	Length           int32
	Width            int32
	Height           int32
	BasketValue      int32
	CODAmount        int32
	TryOn            model.TryOn
	ShippingNote     string
}
