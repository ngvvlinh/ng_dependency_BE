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

func (x Gender) ToModel() string {
	if x == 0 {
		return ""
	}
	return x.String()
}

func (x Gender) MarshalJSON() ([]byte, error) {
	return []byte(`"` + x.String() + `"`), nil
}
