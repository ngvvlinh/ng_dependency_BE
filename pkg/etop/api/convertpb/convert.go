package convertpb

import (
	"o.o/api/top/int/shop"
	catalogmodel "o.o/backend/com/main/catalog/model"
)

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
