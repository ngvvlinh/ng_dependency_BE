package shipping_fee_type

func (x ShippingFeeType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + x.String() + `"`), nil
}
