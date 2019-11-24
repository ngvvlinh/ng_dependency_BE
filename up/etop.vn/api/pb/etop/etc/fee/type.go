package fee

func (x FeeType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + x.String() + `"`), nil
}
