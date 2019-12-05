package convertpb

import (
	"etop.vn/api/main/catalog"
	"etop.vn/api/top/int/shop"
	catalogmodel "etop.vn/backend/com/main/catalog/model"
)

func PbAttributesToDomain(as []*shop.Attribute) []*catalog.Attribute {
	attrs := make([]*catalog.Attribute, len(as))
	for i, a := range as {
		attrs[i] = &catalog.Attribute{
			Name:  a.Name,
			Value: a.Value,
		}
	}
	return attrs
}

func PbAttributes(as catalog.Attributes) []*shop.Attribute {
	attrs := make([]*shop.Attribute, len(as))
	for i, a := range as {
		attrs[i] = &shop.Attribute{
			Name:  a.Name,
			Value: a.Value,
		}
	}
	return attrs
}

func AttributesToModel(items []*shop.Attribute) []*catalogmodel.ProductAttribute {
	result := make([]*catalogmodel.ProductAttribute, 0, len(items))
	for _, item := range items {
		if item.Name == "" {
			continue
		}
		result = append(result, &catalogmodel.ProductAttribute{
			Name:  item.Name,
			Value: item.Value,
		})
	}
	return result
}

func PbCategories(cs []*catalogmodel.ShopCategory) []*shop.Category {
	res := make([]*shop.Category, len(cs))
	for i, c := range cs {
		res[i] = PbCategory(c)
	}
	return res
}

func PbCategory(m *catalogmodel.ShopCategory) *shop.Category {
	return &shop.Category{
		Id:       m.ID,
		Name:     m.Name,
		ParentId: m.ParentID,
		ShopId:   m.ShopID,

		// deprecated: 2019.07.24+14
		ProductSourceId: m.ShopID,
	}
}
