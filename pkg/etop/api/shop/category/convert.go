package category

import (
	"o.o/api/main/catalog"
	"o.o/api/top/int/shop"
)

func PbShopCategory(m *catalog.ShopCategory) *shop.ShopCategory {
	res := &shop.ShopCategory{
		Id:       m.ID,
		ShopId:   m.ShopID,
		Status:   0,
		ParentId: m.ParentID,
		Name:     m.Name,
	}
	return res
}

func PbShopCategories(items []*catalog.ShopCategory) []*shop.ShopCategory {
	res := make([]*shop.ShopCategory, len(items))
	for i, item := range items {
		res[i] = PbShopCategory(item)
	}
	return res
}
