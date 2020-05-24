// +build wireinject

package main

import (
	"github.com/google/wire"
	"o.o/api/main/identity"
	comidentity "o.o/backend/com/main/identity"
	"o.o/backend/pkg/common/sql/cmsql"
)

func NewIdentity(db *cmsql.Database) identity.CommandBus {
	panic(wire.Build(comidentity.WireSet))
}
