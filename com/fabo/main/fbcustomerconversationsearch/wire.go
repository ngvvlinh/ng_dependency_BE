// +build wireinject

package fbcustomerconversationsearch

import (
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	NewFbSearchServiceQuery,
	FbSearchQueryMessageBus,
)
