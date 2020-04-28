package convert

import (
	"time"

	"github.com/microcosm-cc/bluemonday"

	"o.o/api/webserver"
	cm "o.o/backend/pkg/common"
)

// +gen:convert: o.o/backend/com/web/webserver/model  -> o.o/api/webserver
// +gen:convert:  o.o/api/webserver

func createOrUpdateWsCategory(in *webserver.CreateOrUpdateWsCategoryArgs, out *webserver.WsCategory) {
	if out == nil {
		out = &webserver.WsCategory{}
	}
	apply_webserver_CreateOrUpdateWsCategoryArgs_webserver_WsCategory(in, out)
	out.UpdatedAt = time.Now()
	if in.ID == 0 {
		out.ID = cm.NewID()
	}
}

func createOrUpdateWsProduct(in *webserver.CreateOrUpdateWsProductArgs, out *webserver.WsProduct) {
	if out == nil {
		out = &webserver.WsProduct{}
	}
	apply_webserver_CreateOrUpdateWsProductArgs_webserver_WsProduct(in, out)
	out.UpdatedAt = time.Now()
	if in.ID == 0 {
		out.ID = cm.NewID()
	}
	if in.DescHTML.Valid == true {
		p := bluemonday.UGCPolicy()
		var descHTML = p.Sanitize(in.DescHTML.String)
		out.DescHTML = descHTML
	}
}

func updateWsPage(in *webserver.UpdateWsPageArgs, out *webserver.WsPage) {
	if out == nil {
		out = &webserver.WsPage{}
	}
	apply_webserver_UpdateWsPageArgs_webserver_WsPage(in, out)
	out.UpdatedAt = time.Now()
	if in.DescHTML.Valid == true {
		p := bluemonday.UGCPolicy()
		var descHTML = p.Sanitize(in.DescHTML.String)
		out.DescHTML = descHTML
	}
}

func createWsPage(in *webserver.CreateWsPageArgs, out *webserver.WsPage) {
	if out == nil {
		out = &webserver.WsPage{}
	}
	apply_webserver_CreateWsPageArgs_webserver_WsPage(in, out)
	out.ID = cm.NewID()
}

func createWsWebsite(in *webserver.CreateWsWebsiteArgs, out *webserver.WsWebsite) {
	if out == nil {
		out = &webserver.WsWebsite{}
	}
	apply_webserver_CreateWsWebsiteArgs_webserver_WsWebsite(in, out)
	out.UpdatedAt = time.Now()
	out.ID = cm.NewID()
	if in.Description != "" {
		p := bluemonday.UGCPolicy()
		out.Description = p.Sanitize(in.Description)
	}
}

func updateWsWebsite(in *webserver.UpdateWsWebsiteArgs, out *webserver.WsWebsite) {
	if out == nil {
		out = &webserver.WsWebsite{}
	}
	if in.OutstandingProduct == nil {
		in.OutstandingProduct = out.OutstandingProduct
	}
	if in.Banner == nil {
		in.Banner = out.Banner
	}
	if in.SEOConfig == nil {
		in.SEOConfig = out.SEOConfig
	}
	if in.Facebook == nil {
		in.Facebook = out.Facebook
	}
	if in.ShopInfo == nil {
		in.ShopInfo = out.ShopInfo
	}
	apply_webserver_UpdateWsWebsiteArgs_webserver_WsWebsite(in, out)
	out.ID = cm.NewID()
	if in.Description.Valid == true {
		p := bluemonday.UGCPolicy()
		var description = p.Sanitize(in.Description.String)
		out.Description = description
	}

}
