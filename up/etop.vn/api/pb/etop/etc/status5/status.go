package status5

func (x Status) MarshalJSON() ([]byte, error) {
	return []byte(`"` + x.String() + `"`), nil
}
