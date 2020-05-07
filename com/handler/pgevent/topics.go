package pgevent

import "o.o/backend/pkg/etop/model"

type TopicDef struct {
	Name       string
	Partitions int
	DBName     model.DBName
}

var Topics = []TopicDef{{
	Name:       "fulfillment",
	Partitions: 64,
	DBName:     model.DBMain,
}, {
	Name:       "order",
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
}, {
	Name:       "shop_product",
	Partitions: 64,
	DBName:     model.DBMain,
}, {
	Name:       "shop_collection",
	Partitions: 64,
	DBName:     model.DBMain,
}, {
	Name:       "shop_product_collection",
	Partitions: 64,
	DBName:     model.DBMain,
}, {
	Name:       "shop_variant",
	Partitions: 64,
	DBName:     model.DBMain,
}, {
	Name:       "shop_customer",
	Partitions: 64,
	DBName:     model.DBMain,
}, {
	Name:       "shop_customer_group",
	Partitions: 64,
	DBName:     model.DBMain,
}, {
	Name:       "shop_customer_group_customer",
	Partitions: 64,
	DBName:     model.DBMain,
}, {
	Name:       "inventory_variant",
	Partitions: 64,
	DBName:     model.DBMain,
}, {
	Name:       "shop_trader_address",
	Partitions: 64,
	DBName:     model.DBMain,
}, {
	Name:       "fb_external_conversation",
	Partitions: 64,
	DBName:     model.DBMain,
}, {
	Name:       "fb_external_message",
	Partitions: 64,
	DBName:     model.DBMain,
}, {
	Name:       "fb_external_comment",
	Partitions: 64,
	DBName:     model.DBMain,
}, {
	Name:       "fb_external_conversation_fabo",
	Partitions: 64,
	DBName:     model.DBMain,
}, {
	Name:       "fb_external_message_fabo",
	Partitions: 64,
	DBName:     model.DBMain,
}, {
	Name:       "fb_external_comment_fabo",
	Partitions: 64,
	DBName:     model.DBMain,
}}

var TopicMap map[string]TopicDef

func init() {
	TopicMap = make(map[string]TopicDef)
	for _, d := range Topics {
		TopicMap[d.Name] = d
	}
}
