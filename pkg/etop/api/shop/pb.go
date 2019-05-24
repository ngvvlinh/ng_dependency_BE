package shop

import (
	pbcm "etop.vn/backend/pb/common"
	pbs3 "etop.vn/backend/pb/etop/etc/status3"
	pbshop "etop.vn/backend/pb/etop/shop"
	"etop.vn/backend/pkg/etop/api/convertpb"
	catalogmodel "etop.vn/backend/pkg/services/catalog/model"
)

func PbEtopVariants(items []*catalogmodel.VariantExtended) []*pbshop.EtopVariant {
	if items == nil || len(items) == 0 {
		return nil
	}
	res := make([]*pbshop.EtopVariant, len(items))
	for i, item := range items {
		res[i] = PbEtopVariant(item)
	}
	return res
}

func PbEtopVariant(m *catalogmodel.VariantExtended) *pbshop.EtopVariant {
	res := &pbshop.EtopVariant{
		Id:             m.ID,
		Code:           m.EdCode, // yes, it's EdCode
		ListPrice:      int32(m.ListPrice),
		CostPrice:      int32(m.CostPrice),
		WholesalePrice: 0,
		RetailPriceMin: 0,
		RetailPriceMax: 0,
		Name:           m.GetName(),
		Description:    coalesce(m.Description, m.EdDescription),
		ShortDesc:      coalesce(m.ShortDesc, m.EdShortDesc),
		DescHtml:       coalesce(m.DescHTML, m.EdDescHTML),
		ImageUrls:      coalesceStrings(m.ImageURLs),
	}
	return res
}

func PbEtopProduct(m *catalogmodel.ProductFtVariant) *pbshop.EtopProduct {
	return &pbshop.EtopProduct{
		Id:         m.Product.ID,
		CategoryId: m.Product.EtopCategoryID,

		ProductSourceCategoryId: m.Product.ProductSourceCategoryID,

		Name:              coalesce(m.Product.EdName, m.Product.Name),
		Description:       m.Product.Description,
		ShortDesc:         coalesce(m.EdShortDesc, m.ShortDesc, m.Product.ShortDesc),
		DescHtml:          coalesce(m.EdDescHTML, m.DescHTML, m.Product.DescHTML),
		ImageUrls:         m.Product.ImageURLs,
		QuantityAvailable: 100,
		Code:              m.Product.EdCode, // yes, it's EdCode
	}
}

func PbShopVariants(items []*catalogmodel.ShopVariantExtended) []*pbshop.ShopVariant {
	res := make([]*pbshop.ShopVariant, len(items))
	for i, item := range items {
		res[i] = PbShopVariant(item)
	}
	return res
}

func PbShopVariant(m *catalogmodel.ShopVariantExtended) *pbshop.ShopVariant {
	sv := m.ShopVariant
	res := &pbshop.ShopVariant{
		Id:           sv.VariantID,
		Info:         PbEtopVariant(&m.VariantExtended),
		Code:         m.EdCode, // yes, it's EdCode
		EdCode:       m.EdCode,
		Name:         sv.Name,
		Description:  sv.Description,
		ShortDesc:    sv.ShortDesc,
		DescHtml:     sv.DescHTML,
		ImageUrls:    sv.ImageURLs,
		Tags:         sv.Tags,
		Note:         sv.Note,
		Status:       pbs3.Pb(sv.Status),
		IsAvailable:  m.VariantExtended.IsAvailable(),
		ListPrice:    int32(m.ListPrice),
		RetailPrice:  int32(sv.RetailPrice),
		CostPrice:    int32(m.CostPrice),
		CollectionId: sv.CollectionID,
		Attributes:   convertpb.PbAttributes(m.Attributes),
	}
	res.Info = PbEtopVariant(&m.VariantExtended)
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

func PbShopProductsFtVariant(items []*catalogmodel.ShopProductFtVariant) []*pbshop.ShopProduct {
	res := make([]*pbshop.ShopProduct, len(items))
	for i, item := range items {
		res[i] = PbShopProductFtVariant(item)
	}
	return res
}

func PbShopProductFtVariant(m *catalogmodel.ShopProductFtVariant) *pbshop.ShopProduct {
	res := &pbshop.ShopProduct{
		Id:                m.ShopProduct.ProductID,
		Name:              m.ShopProduct.Name,
		Description:       m.ShopProduct.Description,
		DescHtml:          m.ShopProduct.DescHTML,
		ShortDesc:         m.ShopProduct.ShortDesc,
		ImageUrls:         m.ShopProduct.ImageURLs,
		Status:            pbs3.Pb(m.ShopProduct.Status),
		Tags:              m.Tags,
		CollectionIds:     m.CollectionIDs,
		Variants:          PbShopVariants(m.Variants),
		ProductSourceId:   m.ShopProduct.ProductSourceID,
		ProductSourceName: m.ShopProduct.ProductSourceName,
		ProductSourceType: m.ShopProduct.ProductSourceType,
	}

	res.Info = PbEtopProduct(&catalogmodel.ProductFtVariant{
		ProductExtended: catalogmodel.ProductExtended{
			Product: m.Product,
		},
	})
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
