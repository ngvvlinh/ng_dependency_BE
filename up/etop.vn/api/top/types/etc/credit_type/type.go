package credit_type

func PbType(s string) CreditType {
	return CreditType(CreditType_value[s])
}

func PbTypeFromInt(s int) CreditType {
	return CreditType(s)
}

func (x *CreditType) ToModel() string {
	if x == nil {
		return ""
	}
	return CreditType_name[int(*x)]
}

func (x CreditType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + x.String() + `"`), nil
}
