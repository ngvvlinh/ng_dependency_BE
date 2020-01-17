package pgevent

import "etop.vn/backend/pkg/etop/model"

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
}}

var TopicMap map[string]TopicDef

func init() {
	TopicMap = make(map[string]TopicDef)
	for _, d := range Topics {
		TopicMap[d.Name] = d
	}
}
