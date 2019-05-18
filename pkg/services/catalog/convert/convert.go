package convert

import (
	catalogtypes "etop.vn/api/main/catalog/v1/types"
	catalogmodel "etop.vn/backend/pkg/services/catalog/model"
)

func AttributeToModel(in *catalogtypes.Attribute) (out *catalogmodel.ProductAttribute) {
	if in == nil {
		return nil
	}
	return &catalogmodel.ProductAttribute{
		Name:  in.Name,
		Value: in.Value,
	}
}

func AttributesToModel(ins []*catalogtypes.Attribute) (outs []catalogmodel.ProductAttribute) {
	for _, in := range ins {
		outs = append(outs, *AttributeToModel(in))
	}
	return outs
}

func Attribute(in catalogmodel.ProductAttribute) (out *catalogtypes.Attribute) {
	return &catalogtypes.Attribute{
		Name:  in.Name,
		Value: in.Value,
	}
}

func Attributes(ins []catalogmodel.ProductAttribute) (outs []*catalogtypes.Attribute) {
	for _, in := range ins {
		outs = append(outs, Attribute(in))
	}
	return outs
}
