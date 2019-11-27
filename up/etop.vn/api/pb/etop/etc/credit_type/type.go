package credit_type

func PbType(s string) CreditType {
	return CreditType(CreditType_value[s])
}

func PbTypeFromInt(s int) CreditType {
	return CreditType(s)
}

func (s *CreditType) ToModel() string {
	if s == nil {
		return ""
	}
	return CreditType_name[int(*s)]
}

func (x CreditType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + x.String() + `"`), nil
}
