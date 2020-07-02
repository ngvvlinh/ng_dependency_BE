package shop

import (
	"context"

	"o.o/api/main/catalog"
	"o.o/api/top/int/shop"
	api "o.o/api/top/int/shop"
	"o.o/api/webserver"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/capi/dot"
)

func (s *WebServerService) CreateWsWebsite(ctx context.Context, r *api.CreateWsWebsiteRequest) (*api.WsWebsite, error) {
	shopID := s.SS.Shop().ID
	cmd := &webserver.CreateWsWebsiteCommand{
		ShopID:             shopID,
		MainColor:          r.MainColor,
		Banner:             ConvertBanner(r.Banner),
		OutstandingProduct: ConvertSpecialProduct(r.OutstandingProduct),
		NewProduct:         ConvertSpecialProduct(r.NewProduct),
		SEOConfig:          ConvertWsGeneralSEO(r.SEOConfig),
		Facebook:           ConvertFacebook(r.Facebook),
		GoogleAnalyticsID:  r.GoogleAnalyticsID,
		DomainName:         r.DomainName,
		OverStock:          r.OverStock,
		ShopInfo:           ConvertShopInfo(r.ShopInfo),
		Description:        r.Description,
		LogoImage:          r.LogoImage,
		FaviconImage:       r.FaviconImage,
		SiteSubdomain:      r.SiteSubdomain,
	}
	err := s.WebserverAggr.Dispatch(ctx, cmd)
	if err != nil {
		return nil, err
	}
	result := PbWsWebsite(cmd.Result)
	result, err = s.populateWsWebSiteWithProduct(ctx, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *WebServerService) UpdateWsWebsite(ctx context.Context, r *api.UpdateWsWebsiteRequest) (*api.WsWebsite, error) {
	shopID := s.SS.Shop().ID
	cmd := &webserver.UpdateWsWebsiteCommand{
		ShopID:             shopID,
		ID:                 r.ID,
		MainColor:          r.MainColor,
		Banner:             ConvertBanner(r.Banner),
		OutstandingProduct: ConvertSpecialProduct(r.OutstandingProduct),
		NewProduct:         ConvertSpecialProduct(r.NewProduct),
		SEOConfig:          ConvertWsGeneralSEO(r.SEOConfig),
		Facebook:           ConvertFacebook(r.Facebook),
		GoogleAnalyticsID:  r.GoogleAnalyticsID,
		DomainName:         r.DomainName,
		OverStock:          r.OverStock,
		ShopInfo:           ConvertShopInfo(r.ShopInfo),
		Description:        r.Description,
		LogoImage:          r.LogoImage,
		FaviconImage:       r.FaviconImage,
		SiteSubdomain:      r.SiteSubdomain,
	}
	err := s.WebserverAggr.Dispatch(ctx, cmd)
	if err != nil {
		return nil, err
	}
	result := PbWsWebsite(cmd.Result)
	result, err = s.populateWsWebSiteWithProduct(ctx, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *WebServerService) GetWsWebsite(ctx context.Context, r *api.GetWsWebsiteRequest) (*api.WsWebsite, error) {
	shopID := s.SS.Shop().ID
	query := &webserver.GetWsWebsiteByIDQuery{
		ShopID: shopID,
		ID:     0,
		Result: nil,
	}
	err := s.WebserverQuery.Dispatch(ctx, query)
	if err != nil {
		return nil, err
	}
	result := PbWsWebsite(query.Result)
	result, err = s.populateWsWebSiteWithProduct(ctx, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *WebServerService) GetWsWebsites(ctx context.Context, r *api.GetWsWebsitesRequest) (*api.GetWsWebsitesResponse, error) {
	shopID := s.SS.Shop().ID
	paging := cmapi.CMPaging(r.Paging)
	query := &webserver.ListWsWebsitesQuery{
		ShopID:  shopID,
		Paging:  *paging,
		Filters: cmapi.ToFilters(r.Filters),
		Result:  nil,
	}
	err := s.WebserverQuery.Dispatch(ctx, query)
	if err != nil {
		return nil, err
	}
	resp := PbWsWebsites(query.Result.WsWebsites)
	resp, err = s.populateWsWebSitesWithProduct(ctx, resp)
	if err != nil {
		return nil, err
	}
	result := &api.GetWsWebsitesResponse{
		WsWebsites: resp,
		Paging:     cmapi.PbPaging(query.Paging),
	}
	return result, nil
}

func (s *WebServerService) GetWsWebsitesByIDs(ctx context.Context, r *api.GetWsWebsitesByIDsRequest) (*api.GetWsWebsitesByIDsResponse, error) {
	shopID := s.SS.Shop().ID
	query := &webserver.ListWsWebsitesByIDsQuery{
		ShopID: shopID,
		IDs:    r.IDs,
		Result: nil,
	}
	err := s.WebserverQuery.Dispatch(ctx, query)
	if err != nil {
		return nil, err
	}
	resp := PbWsWebsites(query.Result)
	resp, err = s.populateWsWebSitesWithProduct(ctx, resp)
	if err != nil {
		return nil, err
	}
	result := &api.GetWsWebsitesByIDsResponse{
		WsWebsites: resp,
	}
	return result, nil
}

func (s *WebServerService) populateWsWebSiteWithProduct(ctx context.Context, args *api.WsWebsite) (*api.WsWebsite, error) {
	var productIDs []dot.ID
	if args.NewProduct == nil && args.OutstandingProduct == nil {
		return args, nil
	}
	if args.NewProduct != nil {
		productIDs = makeUnduplicatedIDList(productIDs, args.NewProduct.ProductIDs)
	}
	if args.OutstandingProduct != nil {
		productIDs = makeUnduplicatedIDList(productIDs, args.OutstandingProduct.ProductIDs)
	}
	query := &catalog.ListShopProductsWithVariantsByIDsQuery{
		IDs:    productIDs,
		ShopID: args.ShopID,
	}
	err := s.CatalogQuery.Dispatch(ctx, query)
	if err != nil {
		return nil, err
	}
	var mapProducts = make(map[dot.ID]*catalog.ShopProductWithVariants)
	for _, product := range query.Result.Products {
		mapProducts[product.ProductID] = product
	}

	if args.NewProduct != nil {
		for _, v := range args.NewProduct.ProductIDs {
			if mapProducts[v] == nil {
				return nil, cm.Errorf(cm.InvalidArgument, nil, "Sản phẩm id = %v không tồn tại", v)
			}
			args.NewProduct.Products = append(args.NewProduct.Products, PbShopProductWithVariants(mapProducts[v]))
		}
	}
	if args.OutstandingProduct != nil {
		for _, v := range args.OutstandingProduct.ProductIDs {
			if mapProducts[v] == nil {
				return nil, cm.Errorf(cm.InvalidArgument, nil, "Sản phẩm id = %v không tồn tại", v)
			}
			args.OutstandingProduct.Products = append(args.OutstandingProduct.Products, PbShopProductWithVariants(mapProducts[v]))
		}
	}
	return args, nil
}

func (s *WebServerService) populateWsWebSitesWithProduct(ctx context.Context, args []*shop.WsWebsite) ([]*shop.WsWebsite, error) {
	if len(args) == 0 {
		return args, nil
	}
	var productIDs []dot.ID
	for _, wsWebsite := range args {
		if wsWebsite.NewProduct != nil {
			productIDs = makeUnduplicatedIDList(productIDs, wsWebsite.NewProduct.ProductIDs)
		}
		if wsWebsite.OutstandingProduct != nil {
			productIDs = makeUnduplicatedIDList(productIDs, wsWebsite.OutstandingProduct.ProductIDs)
		}
	}
	query := &catalog.ListShopProductsWithVariantsByIDsQuery{
		IDs:    productIDs,
		ShopID: args[0].ShopID,
	}
	err := s.CatalogQuery.Dispatch(ctx, query)
	if err != nil {
		return nil, err
	}
	var mapProducts = make(map[dot.ID]*catalog.ShopProductWithVariants)
	for _, product := range query.Result.Products {
		mapProducts[product.ProductID] = product
	}
	for key, _ := range args {
		if args[key].NewProduct != nil {
			for _, v := range args[key].NewProduct.ProductIDs {
				if mapProducts[v] == nil {
					return nil, cm.Errorf(cm.InvalidArgument, nil, "Sản phẩm id = %v không tồn tại", v)
				}
				args[key].NewProduct.Products = append(args[key].NewProduct.Products, PbShopProductWithVariants(mapProducts[v]))
			}
		}
		if args[key].NewProduct != nil {
			for _, v := range args[key].OutstandingProduct.ProductIDs {
				if mapProducts[v] == nil {
					return nil, cm.Errorf(cm.InvalidArgument, nil, "Sản phẩm id = %v không tồn tại", v)
				}
				args[key].OutstandingProduct.Products = append(args[key].OutstandingProduct.Products, PbShopProductWithVariants(mapProducts[v]))
			}
		}
	}
	return args, nil
}

func makeUnduplicatedIDList(in, out []dot.ID) []dot.ID {
	for _, value := range out {
		isExisted := false
		for _, productID := range in {
			if value == productID {
				isExisted = true
				break
			}
		}
		if !isExisted {
			in = append(in, value)
		}
	}
	return in
}
