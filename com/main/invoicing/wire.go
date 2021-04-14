// +build wireinject

package invoicing

import (
	"github.com/google/wire"
	invoicingpm "o.o/backend/com/main/invoicing/pm"
)

var WireSet = wire.NewSet(
	invoicingpm.New,
	NewInvoiceAggregate, InvoiceAggregateMessageBus,
	NewInvoiceQuery, InvoiceQueryMessageBus,
)
