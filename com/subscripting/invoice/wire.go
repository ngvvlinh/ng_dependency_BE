// +build wireinject

package invoice

import (
	"github.com/google/wire"
	invoicepm "o.o/backend/com/subscripting/invoice/pm"
)

var WireSet = wire.NewSet(
	invoicepm.New,
	NewInvoiceAggregate, InvoiceAggregateMessageBus,
	NewInvoiceQuery, InvoiceQueryMessageBus,
)
