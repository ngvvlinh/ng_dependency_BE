package shop

import (
	"etop.vn/api/main/catalog"
	pbs3 "etop.vn/backend/pb/etop/etc/status3"
	pbshop "etop.vn/backend/pb/etop/shop"
	"etop.vn/backend/pkg/etop/api/convertpb"
	"etop.vn/backend/pkg/etop/model"
	catalogmodel "etop.vn/backend/pkg/services/catalog/model"
)

func PbShopVariants(items []*catalog.ShopVariant) []*pbshop.ShopVariant {
	res := make([]*pbshop.ShopVariant, len(items))
	for i, item := range items {
		res[i] = PbShopVariant(item)
	}
	return res
}

func PbShopVariant(m *catalog.ShopVariant) *pbshop.ShopVariant {
	res := &pbshop.ShopVariant{
		Id: m.VariantID,
		Info: &pbshop.EtopVariant{
			Id:          0,
			Code:        m.Code,
			Name:        m.Name,
			Description: m.Description,
			ShortDesc:   m.ShortDesc,
			DescHtml:    m.DescHTML,
			ImageUrls:   m.ImageURLs,
			ListPrice:   m.ListPrice,
			CostPrice:   m.CostPrice,
			Attributes:  convertpb.PbAttributes(m.Attributes),
		},
		Code:        m.Code,
		EdCode:      m.Code,
		Name:        m.Name,
		Description: m.Description,
		ShortDesc:   m.ShortDesc,
		DescHtml:    m.DescHTML,
		ImageUrls:   m.ImageURLs,
		Tags:        nil,
		Note:        m.Note,
		Status:      pbs3.Pb(model.Status3(m.Status)),
		ListPrice:   m.ListPrice,
		RetailPrice: coalesceInt32(m.RetailPrice, m.ListPrice),
		CostPrice:   m.CostPrice,
		Attributes:  convertpb.PbAttributes(m.Attributes),
	}
	return res
}

func PbShopProducts(items []*catalogmodel.ShopProduct) []*pbshop.ShopProduct {
	res := make([]*pbshop.ShopProduct, len(items))
	for i, item := range items {
		res[i] = PbShopProduct(item)
	}
	return res
}

func PbShopProduct(m *catalogmodel.ShopProduct) *pbshop.ShopProduct {
	res := &pbshop.ShopProduct{
		Id:            m.ProductID,
		Info:          nil,
		Code:          "",
		EdCode:        "",
		Name:          m.Name,
		Description:   m.Description,
		ShortDesc:     m.ShortDesc,
		DescHtml:      m.DescHTML,
		ImageUrls:     m.ImageURLs,
		Tags:          m.Tags,
		Stags:         nil,
		Note:          "",
		Status:        pbs3.Pb(m.Status),
		IsAvailable:   true,
		ListPrice:     0,
		RetailPrice:   m.RetailPrice,
		CostPrice:     0,
		CollectionIds: m.CollectionIDs,
		Variants:      nil,

		ProductSourceId: m.ShopID, // deprecated
	}
	return res
}

func PbShopProductsWithVariants(items []*catalog.ShopProductWithVariants) []*pbshop.ShopProduct {
	res := make([]*pbshop.ShopProduct, len(items))
	for i, item := range items {
		res[i] = PbShopProductWithVariants(item)
	}
	return res
}

func PbShopProductWithVariants(m *catalog.ShopProductWithVariants) *pbshop.ShopProduct {
	shopID := m.ShopProduct.ShopID
	res := &pbshop.ShopProduct{
		Id: m.ShopProduct.ProductID,
		Info: &pbshop.EtopProduct{
			Id:          m.ShopProduct.ProductID,
			Code:        m.ShopProduct.Code,
			Name:        m.ShopProduct.Name,
			Description: m.ShopProduct.Description,
			ShortDesc:   m.ShopProduct.ShortDesc,
			DescHtml:    m.ShopProduct.DescHTML,
			Unit:        m.ShopProduct.Unit,
			ImageUrls:   m.ShopProduct.ImageURLs,
			ListPrice:   m.ShopProduct.ListPrice,
			CostPrice:   m.ShopProduct.CostPrice,
			CategoryId:  m.ShopProduct.CategoryID,

			// deprecated
			ProductSourceCategoryId: m.ShopProduct.CategoryID,
		},
		Code:            m.ShopProduct.Code,
		EdCode:          m.ShopProduct.Code,
		Name:            m.ShopProduct.Name,
		Description:     m.ShopProduct.Description,
		ShortDesc:       m.ShopProduct.ShortDesc,
		DescHtml:        m.ShopProduct.DescHTML,
		ImageUrls:       m.ShopProduct.ImageURLs,
		Tags:            m.ShopProduct.Tags,
		Stags:           nil,
		Note:            m.Note,
		Status:          pbs3.Pb(model.Status3(m.ShopProduct.Status)),
		IsAvailable:     false,
		ListPrice:       m.ShopProduct.ListPrice,
		RetailPrice:     coalesceInt32(m.ShopProduct.RetailPrice, m.ShopProduct.ListPrice),
		CostPrice:       m.ShopProduct.CostPrice,
		CollectionIds:   m.ShopProduct.CollectionIDs,
		Variants:        PbShopVariants(m.Variants),
		ProductSourceId: shopID, // backward-compatible: use shop_id in place of product_source_id
	}
	return res
}

func coalesce(ss ...string) string {
	for _, s := range ss {
		if s != "" {
			return s
		}
	}
	return ""
}

func coalesceInt32(is ...int32) int32 {
	for _, i := range is {
		if i != 0 {
			return i
		}
	}
	return 0
}

func coalesceStrings(sss ...[]string) []string {
	for _, ss := range sss {
		if len(ss) != 0 {
			return ss
		}
	}
	return nil
}

func merge(sss ...[]string) []string {
	s0 := sss[0]
	for _, ss := range sss[1:] {
		for _, s := range ss {
			if !contain(s0, s) {
				s0 = append(s0, s)
			}
		}
	}
	return s0
}

func contain(ss []string, s string) bool {
	for _, S := range ss {
		if S == s {
			return true
		}
	}
	return false
}
