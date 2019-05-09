package pgevent

import "etop.vn/backend/pkg/etop/model"

type TopicDef struct {
	Name       string
	Partitions int
	DBName     model.DBName
}

var Topics = []TopicDef{{
	Name:       "account",
	Partitions: 8,
	DBName:     model.DBMain,
}, {
	Name:       "address",
	Partitions: 8,
	DBName:     model.DBMain,
}, {
	Name:       "etop_category",
	Partitions: 8,
	DBName:     model.DBMain,
}, {
	Name:       "fulfillment",
	Partitions: 64,
	DBName:     model.DBMain,
}, {
	Name:       "order_external",
	Partitions: 64,
	DBName:     model.DBMain,
}, {
	Name:       "order_source",
	Partitions: 64,
	DBName:     model.DBMain,
}, {
	Name:       "order",
	Partitions: 64,
	DBName:     model.DBMain,
}, {
	Name:       "product_brand",
	Partitions: 8,
	DBName:     model.DBMain,
}, {
	Name:       "product_external",
	Partitions: 64,
	DBName:     model.DBMain,
}, {
	Name:       "product_source_category_external",
	Partitions: 8,
	DBName:     model.DBMain,
}, {
	Name:       "product_source_category",
	Partitions: 8,
	DBName:     model.DBMain,
}, {
	Name:       "product_source",
	Partitions: 8,
	DBName:     model.DBMain,
}, {
	Name:       "product",
	Partitions: 64,
	DBName:     model.DBMain,
}, {
	Name:       "shop_collection",
	Partitions: 8,
	DBName:     model.DBMain,
}, {
	Name:       "shop_product",
	Partitions: 64,
	DBName:     model.DBMain,
}, {
	Name:       "shop",
	Partitions: 8,
	DBName:     model.DBMain,
}, {
	Name:       "supplier",
	Partitions: 8,
	DBName:     model.DBMain,
}, {
	Name:       "user",
	Partitions: 8,
	DBName:     model.DBMain,
}, {
	Name:       "variant_external",
	Partitions: 64,
	DBName:     model.DBMain,
}, {
	Name:       "variant",
	Partitions: 64,
	DBName:     model.DBMain,
}, {
	Name:       "money_transaction_shipping",
	Partitions: 8,
	DBName:     model.DBMain,
}, {
	Name:       "notification",
	Partitions: 64,
	DBName:     model.DBNotifier,
}}

var TopicMap map[string]TopicDef

func init() {
	TopicMap = make(map[string]TopicDef)
	for _, d := range Topics {
		TopicMap[d.Name] = d
	}
}
