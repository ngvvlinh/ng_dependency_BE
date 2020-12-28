package shop

import (
	"o.o/api/main/catalog"
	"o.o/api/top/int/shop"
	"o.o/backend/pkg/common/apifw/cmapi"
)

func PbBrand(args *catalog.ShopBrand) *shop.Brand {
	return &shop.Brand{
		ShopId:      args.ShopID,
		Id:          args.ID,
		Name:        args.BrandName,
		Description: args.Description,
		CreatedAt:   cmapi.PbTime(args.CreatedAt),
		UpdatedAt:   cmapi.PbTime(args.UpdatedAt),
	}
}

func PbBrands(args []*catalog.ShopBrand) []*shop.Brand {
	var brands []*shop.Brand
	for _, value := range args {
		brands = append(brands, PbBrand(value))
	}
	return brands
}

func PbShopCollection(m *catalog.ShopCollection) *shop.ShopCollection {
	res := &shop.ShopCollection{
		Id:          m.ID,
		ShopId:      m.ShopID,
		Description: m.Description,
		DescHtml:    m.DescHTML,
		Name:        m.Name,
		ShortDesc:   m.ShortDesc,
	}
	return res
}

func PbShopCollections(items []*catalog.ShopCollection) []*shop.ShopCollection {
	res := make([]*shop.ShopCollection, len(items))
	for i, item := range items {
		res[i] = PbShopCollection(item)
	}
	return res
}

func coalesceInt(is ...int) int {
	for _, i := range is {
		if i != 0 {
			return i
		}
	}
	return 0
}
