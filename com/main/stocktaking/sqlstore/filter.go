package sqlstore

import "etop.vn/backend/pkg/etop/sqlstore"

var SortShopStocktake = map[string]string{
	"id":         "",
	"created_at": "",
}

var FilterStocktake = sqlstore.FilterWhitelist{
	Arrays: []string{"product_ids"},
	Equals: []string{"code"},
	Status: []string{"status"},
	Dates:  []string{"created_at"},
}
