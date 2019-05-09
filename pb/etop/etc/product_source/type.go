package product_source

func PbType(s string) ProductSourceType {
	return ProductSourceType(ProductSourceType_value[s])
}

func PbTypeFromInt(s int) ProductSourceType {
	return ProductSourceType(s)
}

func (s *ProductSourceType) ToModel() string {
	if s == nil {
		return ""
	}
	return ProductSourceType_name[int32(*s)]
}
