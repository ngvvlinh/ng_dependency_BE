package address_type

func PbType(s string) AddressType {
	return AddressType(AddressType_value[s])
}

func PbTypeFromInt(s int) AddressType {
	return AddressType(s)
}

func (x *AddressType) ToModel() string {
	if x == nil {
		return ""
	}
	return AddressType_name[int(*x)]
}

func (x AddressType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + x.String() + `"`), nil
}
