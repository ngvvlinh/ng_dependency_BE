package status3

func (x Status) MarshalJSON() ([]byte, error) {
	return []byte(`"` + x.String() + `"`), nil
}
