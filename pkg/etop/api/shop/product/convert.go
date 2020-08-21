package product

import (
	"o.o/api/main/catalog"
	"o.o/api/main/inventory"
	"o.o/api/top/int/shop"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/capi/dot"
)

func PbProductsQuantity(shopProducts []*catalog.ShopProductWithVariants, inventoryVariants map[dot.ID]*inventory.InventoryVariant) (res []*shop.ShopProduct) {
	for _, product := range shopProducts {
		productPb := PbProductQuantity(product, inventoryVariants)
		res = append(res, productPb)
	}
	return
}

func PbProductQuantity(shopProduct *catalog.ShopProductWithVariants, inventoryVariants map[dot.ID]*inventory.InventoryVariant) *shop.ShopProduct {
	shopProductPb := PbShopProductWithVariants(shopProduct)
	shopProductPb.Variants = PbVariantsQuantity(shopProduct.Variants, inventoryVariants)
	return shopProductPb
}

func PbVariantsQuantity(shopVariants []*catalog.ShopVariant, inventoryVariants map[dot.ID]*inventory.InventoryVariant) []*shop.ShopVariant {
	var variants []*shop.ShopVariant
	for _, variant := range shopVariants {
		inventoryVariant := inventoryVariants[variant.VariantID]
		valuePb := PbVariantQuantity(variant, inventoryVariant)
		variants = append(variants, valuePb)
	}
	return variants
}

func PbVariantQuantity(shopVariant *catalog.ShopVariant, inventoryVariant *inventory.InventoryVariant) *shop.ShopVariant {
	shopVariantDB := PbShopVariant(shopVariant)
	if inventoryVariant != nil {
		shopVariantDB.InventoryVariant = &shop.InventoryVariantShopVariant{
			QuantityOnHand: inventoryVariant.QuantityOnHand,
			QuantityPicked: inventoryVariant.QuantityPicked,
			Quantity:       inventoryVariant.QuantitySummary,
			CostPrice:      inventoryVariant.CostPrice,
		}
		shopVariantDB.QuantityOnHand = inventoryVariant.QuantityOnHand
		shopVariantDB.QuantityPicked = inventoryVariant.QuantityPicked
		shopVariantDB.Quantity = inventoryVariant.QuantitySummary
		shopVariantDB.CostPrice = inventoryVariant.CostPrice
	}
	return shopVariantDB
}

func PbShopProductsWithVariants(items []*catalog.ShopProductWithVariants) []*shop.ShopProduct {
	res := make([]*shop.ShopProduct, len(items))
	for i, item := range items {
		res[i] = PbShopProductWithVariants(item)
	}
	return res
}

func PbShopProductWithVariants(m *catalog.ShopProductWithVariants) *shop.ShopProduct {
	shopID := m.ShopProduct.ShopID
	metaFields := []*pbcm.MetaField{}
	for _, metaField := range m.MetaFields {
		metaFields = append(metaFields, &pbcm.MetaField{
			Key:   metaField.Key,
			Value: metaField.Value,
		})
	}
	res := &shop.ShopProduct{
		Id: m.ShopProduct.ProductID,
		Info: &shop.EtopProduct{
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
		CategoryId:      m.ShopProduct.CategoryID,
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
		Status:          m.ShopProduct.Status,
		IsAvailable:     false,
		ListPrice:       m.ShopProduct.ListPrice,
		RetailPrice:     coalesceInt(m.ShopProduct.RetailPrice, m.ShopProduct.ListPrice),
		CollectionIds:   m.ShopProduct.CollectionIDs,
		Variants:        PbShopVariants(m.Variants),
		ProductSourceId: shopID, // backward-compatible: use shop_id in place of product_source_id
		CreatedAt:       cmapi.PbTime(m.CreatedAt),
		UpdatedAt:       cmapi.PbTime(m.UpdatedAt),
		ProductType:     m.ProductType.Wrap(),
		MetaFields:      metaFields,
		BrandId:         m.BrandID,
	}
	return res
}

func PbShopVariants(items []*catalog.ShopVariant) []*shop.ShopVariant {
	res := make([]*shop.ShopVariant, len(items))
	for i, item := range items {
		res[i] = PbShopVariant(item)
	}
	return res
}

func PbShopVariant(m *catalog.ShopVariant) *shop.ShopVariant {
	res := &shop.ShopVariant{
		Id: m.VariantID,
		Info: &shop.EtopVariant{
			Id:          0,
			Code:        m.Code,
			Name:        m.Name,
			Description: m.Description,
			ShortDesc:   m.ShortDesc,
			DescHtml:    m.DescHTML,
			ImageUrls:   m.ImageURLs,
			ListPrice:   m.ListPrice,
			Attributes:  m.Attributes,
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
		Status:      m.Status,
		ListPrice:   m.ListPrice,
		RetailPrice: coalesceInt(m.RetailPrice, m.ListPrice),
		Attributes:  m.Attributes,
		ProductId:   m.ProductID,
	}
	return res
}

func PbShopVariantsWithProducts(items []*catalog.ShopVariantWithProduct) []*shop.ShopVariant {
	res := make([]*shop.ShopVariant, len(items))
	for i, item := range items {
		res[i] = PbShopVariantWithProduct(item)
	}
	return res
}

func PbShopVariantWithProduct(m *catalog.ShopVariantWithProduct) *shop.ShopVariant {
	if m == nil {
		return nil
	}
	res := &shop.ShopVariant{
		Id: m.VariantID,
		Info: &shop.EtopVariant{
			Id:          0,
			Code:        m.Code,
			Name:        m.Name,
			Description: m.Description,
			ShortDesc:   m.ShortDesc,
			DescHtml:    m.DescHTML,
			ImageUrls:   m.ImageURLs,
			ListPrice:   m.ListPrice,
			CostPrice:   m.CostPrice,
			Attributes:  m.Attributes,
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
		Status:      m.Status,
		ListPrice:   m.ListPrice,
		RetailPrice: coalesceInt(m.RetailPrice, m.ListPrice),
		Attributes:  m.Attributes,
		ProductId:   m.ShopProduct.ProductID,
	}
	if m.ShopProduct != nil {
		res.Product = &shop.ShopShortProduct{
			Id:   m.ShopProduct.ProductID,
			Name: m.ShopProduct.Name,
		}
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
