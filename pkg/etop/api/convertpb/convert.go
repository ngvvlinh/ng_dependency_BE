package convertpb

import (
	"etop.vn/backend/pb/common"
	pbadmin "etop.vn/backend/pb/etop/admin"
	"etop.vn/backend/pb/etop/etc/status3"
	pbshop "etop.vn/backend/pb/etop/shop"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/etop/model"
	catalogmodel "etop.vn/backend/pkg/services/catalog/model"
)

func PbAttributes(as []catalogmodel.ProductAttribute) []*pbshop.Attribute {
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

func PbCategories(cs []*model.ProductSourceCategoryExtended) []*pbshop.Category {
	res := make([]*pbshop.Category, len(cs))
	for i, c := range cs {
		res[i] = PbCategory(c)
	}
	return res
}

func PbCategory(m *model.ProductSourceCategoryExtended) *pbshop.Category {
	return &pbshop.Category{
		Id:                m.ProductSourceCategory.ID,
		Name:              m.Name,
		ProductSourceId:   m.ProductSourceCategory.ProductSourceID,
		ProductSourceType: m.ProductSourceCategory.ProductSourceType,
		ParentId:          m.ParentID,
		ShopId:            m.ShopID,
		XId:               m.ExternalID,
		XName:             m.ExternalName,
		XParentId:         m.ExternalParentID,
	}
}

func PbVariants(items []*catalogmodel.VariantExtended) []*pbadmin.Variant {
	res := make([]*pbadmin.Variant, len(items))
	for i, item := range items {
		res[i] = PbVariant(item)
	}
	return res
}

func PbVariant(m *catalogmodel.VariantExtended) *pbadmin.Variant {
	// units := make([]*Unit, len(m.ExternalUnits))
	// for i, u := range m.ExternalUnits {
	// 	units[i] = &Unit{
	// 		Code:     u.Code,
	// 		Name:     u.Name,
	// 		FullName: u.FullName,
	// 		Unit:     u.Unit,
	// 		UnitConv: float32(u.UnitConv),
	// 		Price:    int32(u.Price),
	// 	}
	// }

	return &pbadmin.Variant{
		Id:         m.ID,
		CategoryId: m.Product.EtopCategoryID,

		// SCategoryId:  m.ProductSourceCategoryExtendedID, // unused

		// SMeta: common.RawJSONObjectMsg(m.SupplierMeta),
		SName:        m.EdName,
		SShortDesc:   m.EdShortDesc,
		SDescription: m.EdDescription,
		SDescHtml:    m.EdDescHTML,
		ImageUrls:    m.ImageURLs,

		XId:          m.VariantExternal.ExternalID,
		XBaseId:      m.VariantExternal.ExternalProductID,
		XCategoryId:  m.VariantExternal.ExternalCategoryID,
		XCode:        m.VariantExternal.ExternalCode,
		XName:        m.VariantExternal.ExternalName,
		XDescription: m.VariantExternal.ExternalDescription,
		XImageUrls:   m.VariantExternal.ExternalImageURLs,
		XUnit:        m.VariantExternal.ExternalUnit,
		XUnitConv:    float32(m.VariantExternal.ExternalUnitConv),
		XPrice:       int32(m.VariantExternal.ExternalPrice),
		// XUnits:      units,
		XAttributes: PbAttributes(m.VariantExternal.ExternalAttributes),

		// XCreatedAt: etop.PbTime(m.ExternalCreatedAt),
		XUpdatedAt: common.PbTime(m.VariantExternal.ExternalUpdatedAt),
		// XSyncAt:    common.PbTime(m.LastSyncAt),
		UpdatedAt: common.PbTime(m.UpdatedAt),
		CreatedAt: common.PbTime(m.CreatedAt),
		SStatus:   status3.Pb(m.EdStatus),
		XStatus:   status3.Pb(m.VariantExternal.ExternalStatus),
		EStatus:   status3.Pb(m.EtopStatus),
		Status:    status3.Pb(m.Status),

		QuantityAvailable: int32(m.QuantityAvailable),
		QuantityOnHand:    int32(m.QuantityOnHand),
		QuantityReserved:  int32(m.QuantityReserved),

		SWholesalePrice: int32(m.EdWholesalePrice),
		SListPrice:      int32(m.EdListPrice),
		SRetailPriceMin: int32(m.EdRetailPriceMin),
		SRetailPriceMax: int32(m.EdRetailPriceMax),

		WholesalePrice_0: int32(m.WholesalePrice0),
		WholesalePrice:   int32(m.WholesalePrice),
		ListPrice:        int32(m.ListPrice),
		RetailPriceMin:   int32(m.RetailPriceMin),
		RetailPriceMax:   int32(m.RetailPriceMax),
	}
}

func PbProducts(items []*catalogmodel.ProductFtVariant) []*pbadmin.Product {
	res := make([]*pbadmin.Product, len(items))
	for i, item := range items {
		res[i] = PbProduct(item)
	}
	return res
}

func PbProduct(m *catalogmodel.ProductFtVariant) *pbadmin.Product {
	return &pbadmin.Product{
		Id:         m.Product.ID,
		CategoryId: m.EtopCategoryID,

		SName:        cm.Coalesce(m.EdName, m.Product.Name),
		SDescription: cm.Coalesce(m.EdDescription, m.Product.Description),
		SShortDesc:   cm.Coalesce(m.EdShortDesc, m.Product.ShortDesc),
		SDescHtml:    cm.Coalesce(m.EdDescHTML, m.Product.DescHTML),
		STags:        m.EdTags,
		ImageUrls:    m.Product.ImageURLs,

		// XId:          m.Product.ExternalID,
		XCategoryId:  m.ExternalCategoryID,
		XCode:        m.ExternalCode,
		XName:        m.ExternalName,
		XDescription: m.ExternalDescription,
		XImageUrls:   m.ExternalImageURLs,
		XUnit:        m.ExternalUnit,
		// XUnits:      units,
		// XAttributes: PbAttributes(m.ExternalAttributes),

		XCreatedAt: common.PbTime(m.ExternalCreatedAt),
		XUpdatedAt: common.PbTime(m.ExternalUpdatedAt),
		// XSyncAt:    common.PbTime(m.LastSyncAt),
		UpdatedAt: common.PbTime(m.Product.UpdatedAt),
		CreatedAt: common.PbTime(m.Product.CreatedAt),
		// SStatus:    status3.Pb(m.SupplierStatus),
		XStatus: status3.Pb(m.ExternalStatus),
		// EStatus:    status3.Pb(m.EtopStatus),
		Status: status3.Pb(m.Product.Status),

		QuantityAvailable: int32(m.QuantityAvailable),
		QuantityOnHand:    int32(m.QuantityOnHand),
		QuantityReserved:  int32(m.QuantityReserved),

		Variants: PbVariantFromExternals(m.Variants, m.Product),
	}
}

func PbVariantFromExternals(items []*catalogmodel.VariantExternalExtended, p *catalogmodel.Product) []*pbadmin.Variant {
	res := make([]*pbadmin.Variant, len(items))
	for i, item := range items {
		res[i] = PbVariantFromExternal(item, p)
	}
	return res
}

func PbVariantFromExternal(m *catalogmodel.VariantExternalExtended, p *catalogmodel.Product) *pbadmin.Variant {
	return &pbadmin.Variant{
		Id:         m.Variant.ID,
		CategoryId: p.EtopCategoryID,

		// SCategoryId: p.ProductSourceCategoryID, // unused

		// SMeta:        common.RawJSONObjectMsg(m.SupplierMeta),
		SName:        m.Variant.GetName(),
		SDescription: cm.Coalesce(m.EdDescription, m.Variant.Description),
		SShortDesc:   cm.Coalesce(m.EdShortDesc, m.Variant.ShortDesc),
		SDescHtml:    cm.Coalesce(m.EdDescHTML, m.Variant.DescHTML),
		ImageUrls:    m.ImageURLs,

		XId: m.ExternalID,
		// XBaseId:      m.ExternalBaseID,
		XCategoryId:  m.ExternalCategoryID,
		XCode:        m.ExternalCode,
		XName:        m.ExternalName,
		XDescription: m.ExternalDescription,
		XImageUrls:   m.ExternalImageURLs,
		XUnit:        m.ExternalUnit,
		XUnitConv:    float32(m.ExternalUnitConv),
		XPrice:       int32(m.ExternalPrice),
		// XUnits:      units,
		XAttributes: PbAttributes(m.ExternalAttributes),

		// XCreatedAt: common.PbTime(m.ExternalCreatedAt),
		XUpdatedAt: common.PbTime(m.ExternalUpdatedAt),
		// XSyncAt:    common.PbTime(m.LastSyncAt),
		UpdatedAt: common.PbTime(m.UpdatedAt),
		CreatedAt: common.PbTime(m.CreatedAt),
		SStatus:   status3.Pb(m.EdStatus),
		XStatus:   status3.Pb(m.ExternalStatus),
		EStatus:   status3.Pb(m.EtopStatus),
		Status:    status3.Pb(m.Status),

		QuantityAvailable: int32(m.QuantityAvailable),
		QuantityOnHand:    int32(m.QuantityOnHand),
		QuantityReserved:  int32(m.QuantityReserved),

		SWholesalePrice: int32(m.EdWholesalePrice),
		SListPrice:      int32(m.EdListPrice),
		SRetailPriceMin: int32(m.EdRetailPriceMin),
		SRetailPriceMax: int32(m.EdRetailPriceMax),

		WholesalePrice_0: int32(m.WholesalePrice0),
		WholesalePrice:   int32(m.WholesalePrice),
		ListPrice:        int32(m.ListPrice),
		RetailPriceMin:   int32(m.RetailPriceMin),
		RetailPriceMax:   int32(m.RetailPriceMax),
		Attributes:       PbAttributes(m.Attributes),
	}
}
