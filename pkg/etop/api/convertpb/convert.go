package convertpb

import (
	"etop.vn/api/main/catalog"
	pbshop "etop.vn/backend/pb/etop/shop"
	catalogmodel "etop.vn/backend/pkg/services/catalog/model"
)

func PbAttributes(as catalog.Attributes) []*pbshop.Attribute {
	attrs := make([]*pbshop.Attribute, len(as))
	for i, a := range as {
		attrs[i] = &pbshop.Attribute{
			Name:  a.Name,
			Value: a.Value,
		}
	}
	return attrs
}

func AttributesTomodel(items []*pbshop.Attribute) []catalogmodel.ProductAttribute {
	result := make([]catalogmodel.ProductAttribute, 0, len(items))
	for _, item := range items {
		if item.Name == "" {
			continue
		}
		result = append(result, item.ToModel())
	}
	return result
}

func PbCategories(cs []*catalogmodel.ShopCategory) []*pbshop.Category {
	res := make([]*pbshop.Category, len(cs))
	for i, c := range cs {
		res[i] = PbCategory(c)
	}
	return res
}

func PbCategory(m *catalogmodel.ShopCategory) *pbshop.Category {
	return &pbshop.Category{
		Id:       m.ID,
		Name:     m.Name,
		ParentId: m.ParentID,
		ShopId:   m.ShopID,

		// deprecated: 2019.07.24+14
		ProductSourceId: m.ShopID,
	}
}
