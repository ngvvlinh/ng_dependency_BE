package handler

import (
	"o.o/backend/com/eventhandler"
	"o.o/backend/pkg/common/mq"
	"o.o/backend/pkg/etop/model"
)

func Topics() []eventhandler.TopicDef {
	return []eventhandler.TopicDef{
		{
			Name:       "fulfillment",
			Partitions: 64,
			DBName:     model.DBMain,
		},
		{
			Name:       "order",
			Partitions: 64,
			DBName:     model.DBMain,
		},
		{
			Name:       "money_transaction_shipping",
			Partitions: 8,
			DBName:     model.DBMain,
		},
		{
			Name:       "notification",
			Partitions: 64,
			DBName:     model.DBNotifier,
		},
		{
			Name:       "shop_product",
			Partitions: 64,
			DBName:     model.DBMain,
		},
		{
			Name:       "shop_collection",
			Partitions: 64,
			DBName:     model.DBMain,
		}, {
			Name:       "shop_product_collection",
			Partitions: 64,
			DBName:     model.DBMain,
		},
		{
			Name:       "shop_variant",
			Partitions: 64,
			DBName:     model.DBMain,
		},
		{
			Name:       "shop_customer",
			Partitions: 64,
			DBName:     model.DBMain,
		},
		{
			Name:       "shop_customer_group",
			Partitions: 64,
			DBName:     model.DBMain,
		},
		{
			Name:       "shop_customer_group_customer",
			Partitions: 64,
			DBName:     model.DBMain,
		},
		{
			Name:       "inventory_variant",
			Partitions: 64,
			DBName:     model.DBMain,
		},
		{
			Name:       "shop_trader_address",
			Partitions: 64,
			DBName:     model.DBMain,
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
