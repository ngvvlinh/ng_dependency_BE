package external

func (m *Order) HasChanged() bool {
	return m.Status != nil ||
		m.ConfirmStatus != nil ||
		m.FulfillmentShippingStatus != nil ||
		m.EtopPaymentStatus != nil ||
		m.BasketValue != nil ||
		m.TotalAmount != nil ||
		m.Shipping != nil ||
		m.CustomerAddress != nil || m.ShippingAddress != nil
}

func (m *Fulfillment) HasChanged() bool {
	return m.Status != nil ||
		m.ShippingState != nil ||
		m.EtopPaymentStatus != nil ||
		m.ActualShippingServiceFee != nil ||
		m.CodAmount != nil ||
		m.ActualCodAmount != nil ||
		m.ShippingNote != nil ||
		m.ChargeableWeight != nil
}
