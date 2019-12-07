package modelx

import (
	"time"

	"etop.vn/api/top/types/etc/status3"
	txmodel "etop.vn/backend/com/main/moneytx/model"
	"etop.vn/backend/com/main/moneytx/modely"
	shipmodel "etop.vn/backend/com/main/shipping/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
)

// use for import csv file
type CreateMoneyTransactions struct {
	ShopIDMapFfms map[dot.ID][]*shipmodel.Fulfillment
	ShopIDMap     map[dot.ID]*model.Shop

	Result struct {
		Created int
	}
}

type CreateMoneyTransaction struct {
	Shop           *model.Shop
	FulFillmentIDs []dot.ID
	TotalCOD       int
	TotalAmount    int
	TotalOrders    int

	Result *modely.MoneyTransactionExtended
}

type GetMoneyTransaction struct {
	ID     dot.ID
	ShopID dot.ID

	Result *modely.MoneyTransactionExtended
}

type GetMoneyTransactions struct {
	ShopID                             dot.ID
	IDs                                []dot.ID
	Paging                             *cm.Paging
	MoneyTransactionShippingExternalID dot.ID
	IncludeFulfillments                bool
	Filters                            []cm.Filter

	Result struct {
		MoneyTransactions []*modely.MoneyTransactionExtended
		Total             int
	}
}

type UpdateMoneyTransaction struct {
	ID            dot.ID
	ShopID        dot.ID
	Note          string
	InvoiceNumber string
	BankAccount   *model.BankAccount

	Result *modely.MoneyTransactionExtended
}

type AddFfmsMoneyTransaction struct {
	FulfillmentIDs     []dot.ID
	MoneyTransactionID dot.ID
	ShopID             dot.ID

	Result *modely.MoneyTransactionExtended
}

type RemoveFfmsMoneyTransaction struct {
	FulfillmentIDs     []dot.ID
	ShopID             dot.ID
	MoneyTransactionID dot.ID

	Result *modely.MoneyTransactionExtended
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
	BankAccount    *model.BankAccount
	Note           string
	InvoiceNumber  string

	Result *txmodel.MoneyTransactionShippingExternalExtended
}

type UpdateMoneyTransactionShippingExternal struct {
	ID            dot.ID
	BankAccount   *model.BankAccount
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

type ConfirmMoneyTransactionShippingExternal struct {
	ID dot.ID

	Result struct {
		Updated int
	}
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
		Total                             int
	}
}

type CreateMoneyTransactionShippingEtop struct {
	MoneyTransactionShippingIDs []dot.ID
	BankAccount                 *model.BankAccount
	Note                        string
	InvoiceNumber               string

	Result *modely.MoneyTransactionShippingEtopExtended
}

type GetMoneyTransactionShippingEtop struct {
	ID dot.ID

	Result *modely.MoneyTransactionShippingEtopExtended
}

type GetMoneyTransactionShippingEtops struct {
	IDs     []dot.ID
	Paging  *cm.Paging
	Status  status3.NullStatus
	Filters []cm.Filter

	Result struct {
		Total                         int
		MoneyTransactionShippingEtops []*modely.MoneyTransactionShippingEtopExtended
	}
}
type UpdateMoneyTransactionShippingEtop struct {
	ID            dot.ID
	Adds          []dot.ID
	Deletes       []dot.ID
	ReplaceAll    []dot.ID
	BankAccount   *model.BankAccount
	Note          string
	InvoiceNumber string

	Result *modely.MoneyTransactionShippingEtopExtended
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
