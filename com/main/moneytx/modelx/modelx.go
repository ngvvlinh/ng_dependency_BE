package modelx

import (
	"time"

	"etop.vn/api/top/types/etc/status3"
	identitymodel "etop.vn/backend/com/main/identity/model"
	identitysharemodel "etop.vn/backend/com/main/identity/sharemodel"
	txmodel "etop.vn/backend/com/main/moneytx/model"
	"etop.vn/backend/com/main/moneytx/txmodely"
	shipmodel "etop.vn/backend/com/main/shipping/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/capi/dot"
)

// use for import csv file
type CreateMoneyTransactions struct {
	ShopIDMapFfms map[dot.ID][]*shipmodel.Fulfillment
	ShopIDMap     map[dot.ID]*identitymodel.Shop

	Result struct {
		Created int
	}
}

type CreateMoneyTransaction struct {
	Shop           *identitymodel.Shop
	FulFillmentIDs []dot.ID
	TotalCOD       int
	TotalAmount    int
	TotalOrders    int

	Result *txmodely.MoneyTransactionExtended
}

type GetMoneyTransaction struct {
	ID     dot.ID
	ShopID dot.ID

	Result *txmodely.MoneyTransactionExtended
}

type GetMoneyTransactions struct {
	ShopID                             dot.ID
	IDs                                []dot.ID
	Paging                             *cm.Paging
	MoneyTransactionShippingExternalID dot.ID
	IncludeFulfillments                bool
	Filters                            []cm.Filter

	Result struct {
		MoneyTransactions []*txmodely.MoneyTransactionExtended
	}
}

type GetMoneyTxsByMoneyTxShippingEtopID struct {
	MoneyTxShippingEtopID dot.ID

	Result struct {
		MoneyTransactions []*txmodel.MoneyTransactionShipping
	}
}

type UpdateMoneyTransaction struct {
	ID            dot.ID
	ShopID        dot.ID
	Note          string
	InvoiceNumber string
	BankAccount   *identitysharemodel.BankAccount

	Result *txmodely.MoneyTransactionExtended
}

type AddFfmsMoneyTransaction struct {
	FulfillmentIDs     []dot.ID
	MoneyTransactionID dot.ID
	ShopID             dot.ID

	Result *txmodely.MoneyTransactionExtended
}

type RemoveFfmsMoneyTransaction struct {
	FulfillmentIDs     []dot.ID
	ShopID             dot.ID
	MoneyTransactionID dot.ID

	Result *txmodely.MoneyTransactionExtended
}

type ConfirmMoneyTransaction struct {
	MoneyTransactionID dot.ID
	ShopID             dot.ID
	TotalCOD           int
	TotalAmount        int
	TotalOrders        int

	Result struct {
		Updated int
	}
}

type DeleteMoneyTransaction struct {
	MoneyTransactionID dot.ID
	ShopID             dot.ID

	Result struct {
		Deleted int
	}
}

type CreateMoneyTransactionShippingExternal struct {
	Provider       string
	ExternalPaidAt time.Time
	Lines          []*txmodel.MoneyTransactionShippingExternalLine
	BankAccount    *identitysharemodel.BankAccount
	Note           string
	InvoiceNumber  string

	Result *txmodel.MoneyTransactionShippingExternalExtended
}

type UpdateMoneyTransactionShippingExternal struct {
	ID            dot.ID
	BankAccount   *identitysharemodel.BankAccount
	Note          string
	InvoiceNumber string

	Result *txmodel.MoneyTransactionShippingExternalExtended
}

type CreateMoneyTransactionShippingExternalLine struct {
	ExternalCode                       string
	ExternalTotalCOD                   int
	ExternalCreatedAt                  time.Time
	ExternalClosedAt                   time.Time
	EtopFulfillmentIdRaw               string
	ExternalCustomer                   string
	ExternalAddress                    string
	MoneyTransactionShippingExternalID dot.ID
	ExternalTotalShippingFee           int

	Result *txmodel.MoneyTransactionShippingExternalLine
}

type RemoveMoneyTransactionShippingExternalLines struct {
	MoneyTransactionShippingExternalID dot.ID
	LineIDs                            []dot.ID

	Result *txmodel.MoneyTransactionShippingExternalExtended
}

type ConfirmMoneyTransactionShippingExternals struct {
	IDs []dot.ID

	Result struct {
		Updated int
	}
}

type DeleteMoneyTransactionShippingExternal struct {
	ID dot.ID

	Result struct {
		Deleted int
	}
}

type GetMoneyTransactionShippingExternal struct {
	ID dot.ID

	Result *txmodel.MoneyTransactionShippingExternalExtended
}

type GetMoneyTransactionShippingExternals struct {
	IDs     []dot.ID
	Paging  *cm.Paging
	Filters []cm.Filter

	Result struct {
		MoneyTransactionShippingExternals []*txmodel.MoneyTransactionShippingExternalExtended
	}
}

type CreateMoneyTransactionShippingEtop struct {
	MoneyTransactionShippingIDs []dot.ID
	BankAccount                 *identitysharemodel.BankAccount
	Note                        string
	InvoiceNumber               string

	Result *txmodely.MoneyTransactionShippingEtopExtended
}

type GetMoneyTransactionShippingEtop struct {
	ID dot.ID

	Result *txmodely.MoneyTransactionShippingEtopExtended
}

type GetMoneyTransactionShippingEtops struct {
	IDs     []dot.ID
	Paging  *cm.Paging
	Status  status3.NullStatus
	Filters []cm.Filter

	Result struct {
		MoneyTransactionShippingEtops []*txmodely.MoneyTransactionShippingEtopExtended
	}
}
type UpdateMoneyTransactionShippingEtop struct {
	ID            dot.ID
	Adds          []dot.ID
	Deletes       []dot.ID
	ReplaceAll    []dot.ID
	BankAccount   *identitysharemodel.BankAccount
	Note          string
	InvoiceNumber string

	Result *txmodely.MoneyTransactionShippingEtopExtended
}

type ConfirmMoneyTransactionShippingEtop struct {
	ID          dot.ID
	TotalCOD    int
	TotalAmount int
	TotalOrders int

	Result struct {
		Updated int
	}
}

type DeleteMoneyTransactionShippingEtop struct {
	ID dot.ID

	Result struct {
		Deleted int
	}
}
