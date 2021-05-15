package hotfix

import (
	"github.com/google/wire"
	"o.o/backend/pkg/etop/logic/hotfix/extension"
	"o.o/backend/pkg/etop/logic/hotfix/moneytx"
)

var WireSet = wire.NewSet(
	moneytx.New,
	extension.New,
)
