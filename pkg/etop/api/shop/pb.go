package shop

import (
	"o.o/api/main/catalog"
	"o.o/api/main/purchaserefund"
	"o.o/api/main/refund"
	"o.o/api/top/int/shop"
	"o.o/api/webserver"
	"o.o/backend/pkg/common/apifw/cmapi"
)

func ConvertShopInfo(args *shop.ShopInfo) *webserver.ShopInfo {
	var shopInfo *webserver.ShopInfo
	if args != nil {
		shopInfo = &webserver.ShopInfo{}
		shopInfo.Phone = args.Phone
		shopInfo.FacebookFanpage = args.FacebookFanpage
		shopInfo.Name = args.Name
		shopInfo.Email = args.Email
		if args.Address != nil {
			shopInfo.Address = &webserver.AddressShopInfo{
				Province:     args.Address.Province,
				District:     args.Address.District,
				Ward:         args.Address.Ward,
				DistrictCode: args.Address.DistrictCode,
				ProvinceCode: args.Address.ProvinceCode,
				WardCode:     args.Address.WardCode,
				Address:      args.Address.Address,
			}
		}

	}
	return shopInfo
}

func ConvertBanner(args *shop.Banner) *webserver.Banner {
	var banner *webserver.Banner
	if args != nil {
		banner = &webserver.Banner{}
		banner.Style = args.Style
		var bannerItems = []*webserver.BannerItem{}
		for _, item := range args.BannerItems {
			bannerItems = append(bannerItems, &webserver.BannerItem{
				Alt:   item.Alt,
				Url:   item.Url,
				Image: item.Image,
			})
		}
		banner.BannerItems = bannerItems
	}
	return banner
}

func ConvertFacebook(args *shop.Facebook) *webserver.Facebook {
	var facebook *webserver.Facebook
	if args != nil {
		facebook = &webserver.Facebook{}
		facebook.FacebookID = args.FacebookID
		facebook.WelcomeMessage = args.WelcomeMessage
	}
	return facebook
}

func ConvertSpecialProduct(args *shop.SpecialProduct) *webserver.SpecialProduct {
	var specialProduct *webserver.SpecialProduct
	if args != nil {
		specialProduct = &webserver.SpecialProduct{}
		specialProduct.Style = args.Style
		specialProduct.ProductIDs = args.ProductIDs
	}
	return specialProduct
}

func ConvertSEOConfig(args *shop.WsSEOConfig) *webserver.WsSEOConfig {
	var wsSEOConfig *webserver.WsSEOConfig
	if args != nil {
		wsSEOConfig = &webserver.WsSEOConfig{}
		wsSEOConfig.Content = args.Content
		wsSEOConfig.Description = args.Description
		wsSEOConfig.Keyword = args.Keyword
	}
	return wsSEOConfig
}

func ConvertComparePrice(args []*shop.ComparePrice) []*webserver.ComparePrice {
	var listComparePrice []*webserver.ComparePrice
	for _, value := range args {
		comparePriceItem := &webserver.ComparePrice{
			VariantID:    value.VariantID,
			ComparePrice: value.ComparePrice,
		}
		listComparePrice = append(listComparePrice, comparePriceItem)
	}
	return listComparePrice
}

func ConvertWsGeneralSEO(args *shop.WsGeneralSEO) *webserver.WsGeneralSEO {
	var wsGeneralSEO *webserver.WsGeneralSEO
	if args != nil {
		wsGeneralSEO = &webserver.WsGeneralSEO{}
		wsGeneralSEO.Title = args.Title
		wsGeneralSEO.SiteContent = args.SiteContent
		wsGeneralSEO.SiteMetaDescription = args.SiteMetaDescription
		wsGeneralSEO.SiteMetaKeyword = args.SiteMetaKeyword
	}
	return wsGeneralSEO
}

func PbPurchaseRefunds(args []*purchaserefund.PurchaseRefund) []*shop.PurchaseRefund {
	var result []*shop.PurchaseRefund
	for _, value := range args {
		result = append(result, PbPurchaseRefund(value))
	}
	return result
}

func PbPurchaseRefund(args *purchaserefund.PurchaseRefund) *shop.PurchaseRefund {
	var result = &shop.PurchaseRefund{
		ID:              args.ID,
		ShopID:          args.ShopID,
		PurchaseOrderID: args.PurchaseOrderID,
		Note:            args.Note,
		Code:            args.Code,
		TotalAdjustment: args.TotalAdjustment,
		AdjustmentLines: args.AdjustmentLines,
		Lines:           PbPurchaseRefundLine(args.Lines),
		CreatedAt:       cmapi.PbTime(args.CreatedAt),
		UpdatedAt:       cmapi.PbTime(args.UpdatedAt),
		CancelledAt:     cmapi.PbTime(args.CancelledAt),
		ConfirmedAt:     cmapi.PbTime(args.ConfirmedAt),
		CreatedBy:       args.CreatedBy,
		UpdatedBy:       args.UpdatedBy,
		CancelReason:    args.CancelReason,
		Status:          args.Status,
		TotalAmount:     args.TotalAmount,
		BasketValue:     args.BasketValue,
	}
	return result
}

func PbPurchaseRefundLine(args []*purchaserefund.PurchaseRefundLine) []*shop.PurchaseRefundLine {
	var result []*shop.PurchaseRefundLine
	for _, v := range args {
		result = append(result, &shop.PurchaseRefundLine{
			VariantID:    v.VariantID,
			ProductID:    v.ProductID,
			Quantity:     v.Quantity,
			Code:         v.Code,
			ImageURL:     v.ImageURL,
			Name:         v.ProductName,
			PaymentPrice: v.PaymentPrice,
			Attributes:   v.Attributes,
			Adjustment:   v.Adjustment,
		})
	}
	return result
}

func PbRefunds(args []*refund.Refund) []*shop.Refund {
	var result []*shop.Refund
	for _, value := range args {
		result = append(result, PbRefund(value))
	}
	return result
}

func PbRefund(args *refund.Refund) *shop.Refund {
	var result = &shop.Refund{
		ID:              args.ID,
		ShopID:          args.ShopID,
		OrderID:         args.OrderID,
		Note:            args.Note,
		Code:            args.Code,
		AdjustmentLines: args.AdjustmentLines,
		TotalAdjustment: args.TotalAdjustment,
		Lines:           PbRefundLine(args.Lines),
		CreatedAt:       cmapi.PbTime(args.CreatedAt),
		UpdatedAt:       cmapi.PbTime(args.UpdatedAt),
		CancelledAt:     cmapi.PbTime(args.CancelledAt),
		ConfirmedAt:     cmapi.PbTime(args.ConfirmedAt),
		CreatedBy:       args.CreatedBy,
		UpdatedBy:       args.UpdatedBy,
		CancelReason:    args.CancelReason,
		Status:          args.Status,
		TotalAmount:     args.TotalAmount,
		BasketValue:     args.BasketValue,
	}
	return result
}

func PbRefundLine(args []*refund.RefundLine) []*shop.RefundLine {
	var result []*shop.RefundLine
	for _, v := range args {
		result = append(result, &shop.RefundLine{
			VariantID:   v.VariantID,
			ProductID:   v.ProductID,
			Quantity:    v.Quantity,
			Code:        v.Code,
			ImageURL:    v.ImageURL,
			Name:        v.ProductName,
			RetailPrice: v.RetailPrice,
			Attributes:  v.Attributes,
			Adjustment:  v.Adjustment,
		})
	}
	return result
}

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
