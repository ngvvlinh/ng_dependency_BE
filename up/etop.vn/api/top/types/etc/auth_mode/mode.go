package auth_mode

func (x AuthMode) MarshalJSON() ([]byte, error) {
	return []byte(`"` + x.String() + `"`), nil
}
