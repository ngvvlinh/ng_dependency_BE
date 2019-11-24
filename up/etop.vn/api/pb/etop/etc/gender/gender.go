package gender

func PbGender(s string) Gender {
	if s == "" {
		return 0
	}
	return Gender(Gender_value[s])
}

func PbGenderPtr(s *string) Gender {
	if s == nil || *s == "" {
		return 0
	}
	return Gender(Gender_value[*s])
}

func (g Gender) ToModel() string {
	if g == 0 {
		return ""
	}
	return g.String()
}

func (x Gender) MarshalJSON() ([]byte, error) {
	return []byte(`"` + x.String() + `"`), nil
}
