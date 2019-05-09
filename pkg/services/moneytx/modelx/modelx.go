package modelx

import (
	"time"

	"etop.vn/backend/pkg/services/moneytx/modely"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/etop/model"
	txmodel "etop.vn/backend/pkg/services/moneytx/model"
)

// use for import csv file
type CreateMoneyTransactions struct {
	ShopIDMapFfms map[int64][]*model.Fulfillment
	ShopIDMap     map[int64]*model.Shop

	Result struct {
		Created int
	}
}

type CreateMoneyTransaction struct {
	Shop           *model.Shop
	FulFillmentIDs []int64
	TotalCOD       int
	TotalAmount    int
	TotalOrders    int

	Result *modely.MoneyTransactionExtended
}

type GetMoneyTransaction struct {
	ID     int64
	ShopID int64

	Result *modely.MoneyTransactionExtended
}

type GetMoneyTransactions struct {
	ShopID                             int64
	IDs                                []int64
	Paging                             *cm.Paging
	MoneyTransactionShippingExternalID int64
	IncludeFulfillments                bool
	Filters                            []cm.Filter

	Result struct {
		MoneyTransactions []*modely.MoneyTransactionExtended
		Total             int
	}
}

type UpdateMoneyTransaction struct {
	ID            int64
	ShopID        int64
	Note          string
	InvoiceNumber string
	BankAccount   *model.BankAccount

	Result *modely.MoneyTransactionExtended
}

type AddFfmsMoneyTransaction struct {
	FulfillmentIDs     []int64
	MoneyTransactionID int64
	ShopID             int64

	Result *modely.MoneyTransactionExtended
}

type RemoveFfmsMoneyTransaction struct {
	FulfillmentIDs     []int64
	ShopID             int64
	MoneyTransactionID int64

	Result *modely.MoneyTransactionExtended
}

type ConfirmMoneyTransaction struct {
	MoneyTransactionID int64
	ShopID             int64
	TotalCOD           int
	TotalAmount        int
	TotalOrders        int

	Result struct {
		Updated int
	}
}

type DeleteMoneyTransaction struct {
	MoneyTransactionID int64
	ShopID             int64

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
	ID            int64
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
	MoneyTransactionShippingExternalID int64
	ExternalTotalShippingFee           int

	Result *txmodel.MoneyTransactionShippingExternalLine
}

type RemoveMoneyTransactionShippingExternalLines struct {
	MoneyTransactionShippingExternalID int64
	LineIDs                            []int64

	Result *txmodel.MoneyTransactionShippingExternalExtended
}

type ConfirmMoneyTransactionShippingExternal struct {
	ID int64

	Result struct {
		Updated int
	}
}

type ConfirmMoneyTransactionShippingExternals struct {
	IDs []int64

	Result struct {
		Updated int
	}
}

type DeleteMoneyTransactionShippingExternal struct {
	ID int64

	Result struct {
		Deleted int
	}
}

type GetMoneyTransactionShippingExternal struct {
	ID int64

	Result *txmodel.MoneyTransactionShippingExternalExtended
}

type GetMoneyTransactionShippingExternals struct {
	IDs     []int64
	Paging  *cm.Paging
	Filters []cm.Filter

	Result struct {
		MoneyTransactionShippingExternals []*txmodel.MoneyTransactionShippingExternalExtended
		Total                             int
	}
}

type CreateMoneyTransactionShippingEtop struct {
	MoneyTransactionShippingIDs []int64
	BankAccount                 *model.BankAccount
	Note                        string
	InvoiceNumber               string

	Result *modely.MoneyTransactionShippingEtopExtended
}

type GetMoneyTransactionShippingEtop struct {
	ID int64

	Result *modely.MoneyTransactionShippingEtopExtended
}

type GetMoneyTransactionShippingEtops struct {
	IDs     []int64
	Paging  *cm.Paging
	Status  *model.Status3
	Filters []cm.Filter

	Result struct {
		Total                         int
		MoneyTransactionShippingEtops []*modely.MoneyTransactionShippingEtopExtended
	}
}
type UpdateMoneyTransactionShippingEtop struct {
	ID            int64
	Adds          []int64
	Deletes       []int64
	ReplaceAll    []int64
	BankAccount   *model.BankAccount
	Note          string
	InvoiceNumber string

	Result *modely.MoneyTransactionShippingEtopExtended
}

type ConfirmMoneyTransactionShippingEtop struct {
	ID          int64
	TotalCOD    int
	TotalAmount int
	TotalOrders int

	Result struct {
		Updated int
	}
}

type DeleteMoneyTransactionShippingEtop struct {
	ID int64

	Result struct {
		Deleted int
	}
}
