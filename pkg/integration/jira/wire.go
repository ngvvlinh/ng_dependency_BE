package jira

import (
	"github.com/google/wire"
	"o.o/backend/pkg/integration/jira/driver"
)

var WireSet = wire.NewSet(
	driver.New,
)
