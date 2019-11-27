package product_type

import "etop.vn/api/main/catalog"

func (s *ProductType) ToProductType() catalog.ProductType {
	if s == nil || *s == 0 {
		return ""
	}
	return catalog.ProductType(ProductType_name[int(*s)])
}

func PbProductType(s string) *ProductType {
	res := ProductType(ProductType_value[s])
	return &res
}

func (x ProductType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + x.String() + `"`), nil
}
