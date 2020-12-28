package _all

import (
	"o.o/api/top/int/shop"
	"o.o/api/webserver"
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
