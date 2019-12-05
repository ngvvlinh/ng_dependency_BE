package shipping

func (x State) MarshalJSON() ([]byte, error) {
	return []byte(`"` + x.String() + `"`), nil
}
