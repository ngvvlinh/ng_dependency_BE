package credit

import (
	"time"

	"o.o/api/main/identity"
	"o.o/api/top/types/etc/credit_type"
	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
)

// +gen:event:topic=event/credit

type Credit struct {
	ID              dot.ID
	Amount          int
	ShopID          dot.ID
	Type            credit_type.CreditType
	Status          status3.Status
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       time.Time
	PaidAt          time.Time
	Classify        credit_type.CreditClassify
	BankStatementID dot.ID
}

type CreditExtended struct {
	*Credit
	Shop *identity.Shop
}

type CreditConfirmedEvent struct {
	CreditID dot.ID
	ShopID   dot.ID
}

type CreditCreatedEvent struct {
	CreditID dot.ID
	ShopID   dot.ID
}
