package shop

import (
	pbcm "etop.vn/backend/pb/common"
	pbs3 "etop.vn/backend/pb/etop/etc/status3"
	pbshop "etop.vn/backend/pb/etop/shop"
	"etop.vn/backend/pkg/etop/api/admin"
	"etop.vn/backend/pkg/etop/api/convertpb"
	"etop.vn/backend/pkg/etop/model"
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
		Id: m.ID,
		// ShortName:         strings.Join([]string{m.Product.Name, m.Name}, " - "),
		Name:              m.GetName(),
		Description:       coalesce(m.Description, m.EdDescription),
		ShortDesc:         coalesce(m.ShortDesc, m.EdShortDesc),
		DescHtml:          coalesce(m.DescHTML, m.EdDescHTML),
		ImageUrls:         coalesceStrings(m.ImageURLs),
		WholesalePrice:    int32(m.WholesalePrice),
		ListPrice:         int32(m.ListPrice),
		RetailPriceMin:    int32(m.RetailPriceMin),
		RetailPriceMax:    int32(m.RetailPriceMax),
		IsAvailable:       m.IsAvailable(),
		QuantityAvailable: int32(m.QuantityAvailable),
		Status:            pbs3.Pb(m.Status),

		Code:   m.Code,
		EdCode: m.EdCode,

		// deprecated
		Sku: m.Code,

		// XId:         m.ExternalID,
		// XBaseId:     m.ExternalBaseID,
		// XAttributes: supplier.PbAttributes(m.VariantExternal.ExternalAttributes),

		SMeta:      pbcm.RawJSONObjectMsg(m.SupplierMeta),
		CostPrice:  int32(m.CostPrice),
		Attributes: convertpb.PbAttributes(m.Attributes),
		UpdatedAt:  pbcm.PbTime(m.Product.UpdatedAt),
		CreatedAt:  pbcm.PbTime(m.Product.CreatedAt),
	}

	if m.VariantExternal != nil {
		res.XAttributes = convertpb.PbAttributes(m.VariantExternal.ExternalAttributes)
	}

	if m.Product != nil {
		res.CategoryId = m.Product.EtopCategoryID
	}

	return res
}

func PbEtopProducts(items []*catalogmodel.ProductFtVariant) []*pbshop.EtopProduct {
	res := make([]*pbshop.EtopProduct, len(items))
	for i, item := range items {
		res[i] = PbEtopProduct(item)
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
		IsAvailable:       m.IsAvailable(),
		QuantityAvailable: int32(m.QuantityAvailable),
		Status:            pbs3.Pb(m.Product.Status),
		Code:              m.Product.Code,
		EdCode:            m.Product.EdCode,
		Unit:              m.Product.Unit,

		// XId:         m.ExternalID,

		Variants:  PbEtopVariants(admin.VExternalExtendedToVExtended(m.Variants)),
		UpdatedAt: pbcm.PbTime(m.Product.UpdatedAt),
		CreatedAt: pbcm.PbTime(m.Product.CreatedAt),
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
		Name:         sv.Name,
		Description:  sv.Description,
		ShortDesc:    sv.ShortDesc,
		DescHtml:     sv.DescHTML,
		ImageUrls:    sv.ImageURLs,
		Tags:         sv.Tags,
		Status:       pbs3.Pb(sv.Status),
		IsAvailable:  m.VariantExtended.IsAvailable(),
		RetailPrice:  int32(sv.RetailPrice),
		CollectionId: sv.CollectionID,

		Note: sv.Note,
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
		Id:            m.ProductID,
		Name:          m.Name,
		Description:   m.Description,
		DescHtml:      m.DescHTML,
		ShortDesc:     m.ShortDesc,
		ImageUrls:     m.ImageURLs,
		Status:        pbs3.Pb(m.Status),
		Tags:          m.Tags,
		CollectionIds: m.CollectionIDs,
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

func PbProductSources(items []*model.ProductSource) []*pbshop.ProductSource {
	result := make([]*pbshop.ProductSource, len(items))
	for i, item := range items {
		result[i] = PbProductSource(item)
	}
	return result
}

func PbProductSource(m *model.ProductSource) *pbshop.ProductSource {
	return &pbshop.ProductSource{
		Id:        m.ID,
		Type:      m.Type,
		Name:      m.Name,
		Status:    pbs3.Pb(m.Status),
		CreatedAt: pbcm.PbTime(m.CreatedAt),
		UpdatedAt: pbcm.PbTime(m.UpdatedAt),
	}
}

func PbProductSourceCategories(items []*model.ProductSourceCategory) []*pbshop.ProductSourceCategory {
	result := make([]*pbshop.ProductSourceCategory, len(items))
	for i, item := range items {
		result[i] = PbProductSourceCategory(item)
	}
	return result
}

func PbProductSourceCategory(m *model.ProductSourceCategory) *pbshop.ProductSourceCategory {
	return &pbshop.ProductSourceCategory{
		Id:                m.ID,
		Name:              m.Name,
		ProductSourceId:   m.ProductSourceID,
		ProductSourceType: m.ProductSourceType,
		ParentId:          m.ParentID,
		ShopId:            m.ShopID,
	}
}
