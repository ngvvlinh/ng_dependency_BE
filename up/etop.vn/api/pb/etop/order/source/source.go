package source

func (x Source) MarshalJSON() ([]byte, error) {
	return []byte(`"` + x.String() + `"`), nil
}
