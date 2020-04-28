package sqlstore

import (
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
)

var SortWsCategory = map[string]string{
	"id":         "id",
	"created_at": "created_at",
	"updated_at": "updated_at",
}

var FilterWsCategory = sqlstore.FilterWhitelist{
	Contains: []string{},
	Dates:    []string{"created_at", "updated_at"},
	Equals:   []string{"id", "shop_id"},
	Numbers:  []string{},
	Status:   []string{},
	Arrays:   []string{},
}

var SortWsProduct = map[string]string{
	"id":         "id",
	"created_at": "created_at",
	"updated_at": "updated_at",
}

var FilterWsProduct = sqlstore.FilterWhitelist{
	Contains: []string{},
	Dates:    []string{"created_at", "updated_at"},
	Equals:   []string{"id", "shop_id"},
	Numbers:  []string{},
	Status:   []string{},
	Arrays:   []string{},
}

var SortWsPage = map[string]string{
	"id":         "id",
	"created_at": "created_at",
	"updated_at": "updated_at",
}

var FilterWsPage = sqlstore.FilterWhitelist{
	Contains: []string{},
	Dates:    []string{"created_at", "updated_at"},
	Equals:   []string{"id", "shop_id"},
	Numbers:  []string{},
	Status:   []string{},
	Arrays:   []string{},
}

func (ft WsPageFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

var SortWsWebsite = map[string]string{
	"id":         "id",
	"created_at": "created_at",
	"updated_at": "updated_at",
}

var FilterWsWebsite = sqlstore.FilterWhitelist{
	Contains: []string{},
	Dates:    []string{"created_at", "updated_at"},
	Equals:   []string{"id", "shop_id"},
	Numbers:  []string{},
	Status:   []string{},
	Arrays:   []string{},
}

func (ft WsWebsiteFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}
