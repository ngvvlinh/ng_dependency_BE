package address_type

func PbType(s string) AddressType {
	return AddressType(AddressType_value[s])
}

func PbTypeFromInt(s int) AddressType {
	return AddressType(s)
}

func (s *AddressType) ToModel() string {
	if s == nil {
		return ""
	}
	return AddressType_name[int32(*s)]
}
