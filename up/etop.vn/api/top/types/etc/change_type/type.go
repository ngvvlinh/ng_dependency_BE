package change_type

func (x ChangeType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + x.String() + `"`), nil
}
