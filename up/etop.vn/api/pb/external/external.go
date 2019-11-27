package external

func (m *Order) HasChanged() bool {
	return m.Status != nil ||
		m.ConfirmStatus != nil ||
		m.FulfillmentShippingStatus != nil ||
		m.EtopPaymentStatus != nil ||
		m.BasketValue.Valid ||
		m.TotalAmount.Valid ||
		m.Shipping != nil ||
		m.CustomerAddress != nil || m.ShippingAddress != nil
}

func (m *Fulfillment) HasChanged() bool {
	return m.Status != nil ||
		m.ShippingState != nil ||
		m.EtopPaymentStatus != nil ||
		m.ActualShippingServiceFee.Valid ||
		m.CodAmount.Valid ||
		m.ActualCodAmount.Valid ||
		m.ShippingNote.Valid ||
		m.ChargeableWeight.Valid
}
