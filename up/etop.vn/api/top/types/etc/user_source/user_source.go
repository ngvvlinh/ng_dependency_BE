package user_source

func (x UserSource) MarshalJSON() ([]byte, error) {
	return []byte(`"` + x.String() + `"`), nil
}
