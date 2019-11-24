package shipping_provider

func (x ShippingProvider) MarshalJSON() ([]byte, error) {
	return []byte(`"` + x.String() + `"`), nil
}
