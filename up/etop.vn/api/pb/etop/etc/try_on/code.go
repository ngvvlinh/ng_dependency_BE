package try_on

func (x TryOnCode) MarshalJSON() ([]byte, error) {
	return []byte(x.String()), nil
}
