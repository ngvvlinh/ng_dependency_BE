package handler

import (
	"o.o/backend/com/eventhandler"
	"o.o/backend/pkg/common/mq"
	"o.o/backend/pkg/etc/dbdecl"
)

func Topics() []eventhandler.TopicDef {
	return []eventhandler.TopicDef{
		{
			Name:       "fulfillment",
			Partitions: 64,
			DBName:     dbdecl.DBMain,
		},
		{
			Name:       "order",
			Partitions: 64,
			DBName:     dbdecl.DBMain,
		},
		{
			Name:       "money_transaction_shipping",
			Partitions: 8,
			DBName:     dbdecl.DBMain,
		},
		{
			Name:       "notification",
			Partitions: 64,
			DBName:     dbdecl.DBNotifier,
		},
		{
			Name:       "shop_product",
			Partitions: 64,
			DBName:     dbdecl.DBMain,
		},
		{
			Name:       "shop_collection",
			Partitions: 64,
			DBName:     dbdecl.DBMain,
		}, {
			Name:       "shop_product_collection",
			Partitions: 64,
			DBName:     dbdecl.DBMain,
		},
		{
			Name:       "shop_variant",
			Partitions: 64,
			DBName:     dbdecl.DBMain,
		},
		{
			Name:       "shop_customer",
			Partitions: 64,
			DBName:     dbdecl.DBMain,
		},
		{
			Name:       "shop_customer_group",
			Partitions: 64,
			DBName:     dbdecl.DBMain,
		},
		{
			Name:       "shop_customer_group_customer",
			Partitions: 64,
			DBName:     dbdecl.DBMain,
		},
		{
			Name:       "inventory_variant",
			Partitions: 64,
			DBName:     dbdecl.DBMain,
		},
		{
			Name:       "shop_trader_address",
			Partitions: 64,
			DBName:     dbdecl.DBMain,
		},
		{
			Name:       "shipnow_fulfillment",
			Partitions: 64,
			DBName:     dbdecl.DBMain,
		},
	}
}

var mapTopics = eventhandler.MapTopics(Topics())

func GetTopics(topics map[string]mq.EventHandler) []eventhandler.TopicDef {
	var result []eventhandler.TopicDef
	for name := range topics {
		result = append(result, mapTopics[name])
	}
	return result
}
