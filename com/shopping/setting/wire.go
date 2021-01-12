// +build wireinject

package setting

import (
	"github.com/google/wire"
	"o.o/backend/com/shopping/setting/aggregate"
	"o.o/backend/com/shopping/setting/query"
	"o.o/backend/com/shopping/setting/util"
)

var WireSet = wire.NewSet(
	util.NewShopSettingUtil,
	aggregate.NewShopSettingAggregate, aggregate.ShopSettingAggregateMessageBus,
	query.NewShopSettingQuery, query.ShopSettingQueryMessageBus,
)
