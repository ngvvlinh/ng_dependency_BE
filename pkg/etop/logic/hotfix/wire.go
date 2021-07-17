package hotfix

import (
	"github.com/google/wire"
	"o.o/backend/pkg/etop/logic/hotfix/accountuser"
	"o.o/backend/pkg/etop/logic/hotfix/extension"
	"o.o/backend/pkg/etop/logic/hotfix/moneytx"
	"o.o/backend/pkg/etop/logic/hotfix/user"
)

var WireSet = wire.NewSet(
	moneytx.New,
	extension.New,
	user.New,
	accountuser.New,
)
