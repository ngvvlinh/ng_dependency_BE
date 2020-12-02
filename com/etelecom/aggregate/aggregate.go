package aggregate

import (
	"o.o/api/etelecom"
	"o.o/api/main/identity"
	"o.o/backend/com/etelecom/convert"
	"o.o/backend/com/etelecom/sqlstore"
	com "o.o/backend/com/main"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi"
)

var scheme = conversion.Build(convert.RegisterConversions)
var _ etelecom.Aggregate = &EtelecomAggregate{}

type EtelecomAggregate struct {
	txDB           cmsql.Transactioner
	eventBus       capi.EventBus
	hotlineStore   sqlstore.HotlineStoreFactory
	extensionStore sqlstore.ExtensionStoreFactory
	identityQuery  identity.QueryBus
}

func NewEtelecomAggregate(dbEtelecom com.EtelecomDB, eventBus capi.EventBus) *EtelecomAggregate {
	return &EtelecomAggregate{
		txDB:           (*cmsql.Database)(dbEtelecom),
		eventBus:       eventBus,
		hotlineStore:   sqlstore.NewHotlineStore(dbEtelecom),
		extensionStore: sqlstore.NewExtensionStore(dbEtelecom),
	}
}

func AggregateMessageBus(a *EtelecomAggregate) etelecom.CommandBus {
	b := bus.New()
	return etelecom.NewAggregateHandler(a).RegisterHandlers(b)
}
