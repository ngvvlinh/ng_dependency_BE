package model

import (
	"time"

	"etop.vn/api/top/types/etc/ghn_note_code"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/api/top/types/etc/try_on"
	addressmodel "etop.vn/backend/com/main/address/model"
	identitysharemodel "etop.vn/backend/com/main/identity/sharemodel"
	"etop.vn/capi/dot"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenShop(&Shop{})

type Shop struct {
	ID      dot.ID
	Name    string
	OwnerID dot.ID
	IsTest  int

	AddressID         dot.ID
	ShipToAddressID   dot.ID
	ShipFromAddressID dot.ID
	Phone             string
	BankAccount       *identitysharemodel.BankAccount
	WebsiteURL        string
	ImageURL          string
	Email             string
	Code              string
	AutoCreateFFM     bool

	OrderSourceID dot.ID

	Status    status3.Status
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time

	Address *addressmodel.Address `sq:"-"`

	RecognizedHosts []string

	// @deprecated use try_on instead
	GhnNoteCode ghn_note_code.GHNNoteCode
	TryOn       try_on.TryOnCode
	CompanyInfo *identitysharemodel.CompanyInfo
	// MoneyTransactionRRule format:
	// FREQ=DAILY
	// FREQ=WEEKLY;BYDAY=TU,TH,SA
	// FREQ=MONTHLY;BYDAY=-1FR
	MoneyTransactionRRule         string `sq:"'money_transaction_rrule'"`
	SurveyInfo                    []*SurveyInfo
	ShippingServiceSelectStrategy []*ShippingServiceSelectStrategyItem

	InventoryOverstock dot.NullBool

	Rid dot.ID
}

type ShippingServiceSelectStrategyItem struct {
	Key   string
	Value string
}

type SurveyInfo struct {
	Key      string `json:"key"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}
