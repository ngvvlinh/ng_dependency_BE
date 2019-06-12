package shop

import (
	"etop.vn/api/main/catalog"
	pbcm "etop.vn/backend/pb/common"
	pbs3 "etop.vn/backend/pb/etop/etc/status3"
	pbshop "etop.vn/backend/pb/etop/shop"
	"etop.vn/backend/pkg/etop/api/convertpb"
	"etop.vn/backend/pkg/etop/model"
	catalogmodel "etop.vn/backend/pkg/services/catalog/model"
)

func PbEtopVariant(m *catalog.Variant) *pbshop.EtopVariant {
	res := &pbshop.EtopVariant{
		Id:          m.ID,
		Code:        m.Code,
		Name:        m.Name,
		Description: m.Description,
		ShortDesc:   m.ShortDesc,
		DescHtml:    m.DescHTML,
		ImageUrls:   coalesceStrings(m.ImageURLs),
		ListPrice:   int32(m.ListPrice),
		CostPrice:   int32(m.CostPrice),
		Attributes:  convertpb.PbAttributes(m.Attributes),
	}
	return res
}

func PbEtopProduct(m *catalog.Product) *pbshop.EtopProduct {
	return &pbshop.EtopProduct{
		Id:          m.ID,
		Code:        m.Code,
		Name:        m.Name,
		Description: m.Description,
		ShortDesc:   m.ShortDesc,
		DescHtml:    m.DescHTML,
		ImageUrls:   m.ImageURLs,
		ListPrice:   0,
		CostPrice:   0,

		CategoryId: m.ProductSourceCategoryID,
		// @deprecated
		ProductSourceCategoryId: m.ProductSourceCategoryID,
	}
}

func PbShopVariants(items []*catalog.ShopVariantExtended) []*pbshop.ShopVariant {
	res := make([]*pbshop.ShopVariant, len(items))
	for i, item := range items {
		res[i] = PbShopVariant(item)
	}
	return res
}

func PbShopVariant(m *catalog.ShopVariantExtended) *pbshop.ShopVariant {
	sv := m.ShopVariant
	res := &pbshop.ShopVariant{
		Id:           sv.VariantID,
		Info:         PbEtopVariant(m.Variant),
		Code:         m.ShopVariant.Code,
		EdCode:       m.ShopVariant.Code,
		Name:         sv.Name,
		Description:  sv.Description,
		ShortDesc:    sv.ShortDesc,
		DescHtml:     sv.DescHTML,
		ImageUrls:    sv.ImageURLs,
		Tags:         nil,
		Note:         sv.Note,
		Status:       pbs3.Pb(model.Status3(sv.Status)),
		ListPrice:    int32(m.ShopVariant.ListPrice),
		RetailPrice:  int32(sv.RetailPrice),
		CostPrice:    int32(m.ShopVariant.CostPrice),
		CollectionId: sv.CollectionID,
		Attributes:   convertpb.PbAttributes(m.Attributes),
	}
	res.Info = PbEtopVariant(m.Variant)
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
		Id:                m.ProductID,
		Info:              nil,
		Code:              "",
		EdCode:            "",
		Name:              m.Name,
		Description:       m.Description,
		ShortDesc:         m.ShortDesc,
		DescHtml:          m.DescHTML,
		ImageUrls:         m.ImageURLs,
		Tags:              m.Tags,
		Status:            pbs3.Pb(m.Status),
		IsAvailable:       true,
		CollectionIds:     m.CollectionIDs,
		Variants:          nil,
		ProductSourceId:   0,
		ProductSourceType: "",
		ProductSourceName: "",
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
	res := &pbshop.ShopProduct{
		Id:              m.ShopProduct.ProductID,
		Name:            m.ShopProduct.Name,
		Description:     m.ShopProduct.Description,
		DescHtml:        m.ShopProduct.DescHTML,
		ShortDesc:       m.ShopProduct.ShortDesc,
		ImageUrls:       m.ShopProduct.ImageURLs,
		Status:          pbs3.Pb(model.Status3(m.ShopProduct.Status)),
		Tags:            m.ShopProduct.Tags,
		CollectionIds:   m.ShopProduct.CollectionIDs,
		Variants:        PbShopVariants(m.Variants),
		ProductSourceId: m.Product.ProductSourceID,
	}

	res.Info = PbEtopProduct(m.Product)
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

func PbProductSources(items []*catalogmodel.ProductSource) []*pbshop.ProductSource {
	result := make([]*pbshop.ProductSource, len(items))
	for i, item := range items {
		result[i] = PbProductSource(item)
	}
	return result
}

func PbProductSource(m *catalogmodel.ProductSource) *pbshop.ProductSource {
	return &pbshop.ProductSource{
		Id:        m.ID,
		Type:      m.Type,
		Name:      m.Name,
		Status:    pbs3.Pb(m.Status),
		CreatedAt: pbcm.PbTime(m.CreatedAt),
		UpdatedAt: pbcm.PbTime(m.UpdatedAt),
	}
}
