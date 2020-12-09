// +build wireinject

package setting

import (
	"github.com/google/wire"
	"o.o/backend/com/shopping/setting/aggregate"
	"o.o/backend/com/shopping/setting/query"
)

var WireSet = wire.NewSet(
	aggregate.NewShopSettingAggregate, aggregate.ShopSettingAggregateMessageBus,
	query.NewShopSettingQuery, query.ShopSettingQueryMessageBus,
)
