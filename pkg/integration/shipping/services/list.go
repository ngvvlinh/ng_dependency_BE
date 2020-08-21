package services

import (
	"o.o/api/top/types/etc/shipping_provider"
)

type ShipmentService struct {
	ServiceID string
	Name      string
}

type MapShipmentServices map[shipping_provider.ShippingProvider][]ShipmentService

func (m MapShipmentServices) ByCarrier(carrier shipping_provider.ShippingProvider) []ShipmentService {
	res, ok := m[carrier]
	if !ok {
		return []ShipmentService{}
	}
	return res
}
