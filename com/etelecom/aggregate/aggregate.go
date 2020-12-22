package aggregate

import (
	"o.o/api/etelecom"
	"o.o/api/main/connectioning"
	"o.o/api/main/contact"
	"o.o/api/main/identity"
	"o.o/backend/com/etelecom/convert"
	telecomprovider "o.o/backend/com/etelecom/provider"
	"o.o/backend/com/etelecom/sqlstore"
	com "o.o/backend/com/main"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi"
	"o.o/common/l"
)

var scheme = conversion.Build(convert.RegisterConversions)
var _ etelecom.Aggregate = &EtelecomAggregate{}

var ll = l.New()

type EtelecomAggregate struct {
	txDB            cmsql.Transactioner
	eventBus        capi.EventBus
	hotlineStore    sqlstore.HotlineStoreFactory
	extensionStore  sqlstore.ExtensionStoreFactory
	callLogStore    sqlstore.CallLogStoreFactory
	contactQuery    contact.QueryBus
	identityQuery   identity.QueryBus
	telecomManager  *telecomprovider.TelecomManager
	connectionQuery connectioning.QueryBus
}

func NewEtelecomAggregate(
	dbEtelecom com.EtelecomDB, eventBus capi.EventBus,
	contactQS contact.QueryBus, telecomManager *telecomprovider.TelecomManager,
	connectionQ connectioning.QueryBus,
	identityQ identity.QueryBus,
) *EtelecomAggregate {
	return &EtelecomAggregate{
		txDB:            (*cmsql.Database)(dbEtelecom),
		eventBus:        eventBus,
		contactQuery:    contactQS,
		hotlineStore:    sqlstore.NewHotlineStore(dbEtelecom),
		extensionStore:  sqlstore.NewExtensionStore(dbEtelecom),
		callLogStore:    sqlstore.NewCallLogStore(dbEtelecom),
		telecomManager:  telecomManager,
		connectionQuery: connectionQ,
		identityQuery:   identityQ,
	}
}

func AggregateMessageBus(a *EtelecomAggregate) etelecom.CommandBus {
	b := bus.New()
	return etelecom.NewAggregateHandler(a).RegisterHandlers(b)
}
