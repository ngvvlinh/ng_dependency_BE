package model

import (
	"time"

	"o.o/api/top/types/etc/ghn_note_code"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/try_on"
	addressmodel "o.o/backend/com/main/address/model"
	identitysharemodel "o.o/backend/com/main/identity/sharemodel"
	"o.o/capi/dot"
)

// +sqlgen
type Shop struct {
	ID      dot.ID
	Name    string
	OwnerID dot.ID

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

	Status    status3.Status `sql_type:"int2"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Address *addressmodel.Address `sq:"-"`

	RecognizedHosts []string

	// @deprecated use try_on instead
	GhnNoteCode ghn_note_code.GHNNoteCode `sql_type:"enum(ghn_note_code)"`
	TryOn       try_on.TryOnCode          `sql_type:"enum(try_on)"`
	CompanyInfo *identitysharemodel.CompanyInfo
	// MoneyTransactionRRule format:
	// FREQ=DAILY
	// FREQ=WEEKLY;BYDAY=TU,TH,SA
	// FREQ=MONTHLY;BYDAY=-1FR
	MoneyTransactionRRule string `sq:"'money_transaction_rrule'"`
	SurveyInfo            []*SurveyInfo

	Rid dot.ID
}

type ShippingServiceSelectStrategyItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type SurveyInfo struct {
	Key      string `json:"key"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}
