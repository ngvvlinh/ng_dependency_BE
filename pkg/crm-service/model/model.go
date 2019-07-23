package model

import (
	"fmt"
	"strings"
	"time"

	"etop.vn/backend/pkg/common/validate"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenEtopCount(&EtopAcount{})

type EtopAcount struct {
	ID            string
	FullName      string
	Phone         string
	Email         string
	AccountID     string
	AcountName    string
	AccountType   string
	IsOperator    bool
	VtigerAccount string
}

var _ = sqlgenVtigerContact(&VtigerContact{})

// VtigerAccount table vtiger_acount
type VtigerContact struct {
	ID                   string
	Firstname            string
	ContactNo            string
	Phone                string
	Lastname             string
	Mobile               string
	Email                string
	Leadsource           string
	Secondaryemail       string
	AssignedUserID       string
	CreatedAt            time.Time `sq:"create"`
	EtopID               int64
	UpdatedAt            time.Time `sq:"update"`
	Description          string
	Source               string
	UsedShippingProvider string
	OrdersPerDay         string
	Company              string
	City                 string
	State                string
	Website              string
	Lane                 string
	Country              string
	SearchNorm           string
}

func (p *VtigerContact) BeforeInsertOrUpdate() error {
	var b strings.Builder
	fmt.Fprintf(&b, " %v", p.Phone)
	fmt.Fprintf(&b, " %v", p.Email)
	fmt.Fprintf(&b, " %v", p.Lastname)
	fmt.Fprintf(&b, " %v", p.Firstname)
	fmt.Fprintf(&b, " %v", p.Lane)
	fmt.Fprintf(&b, " %v", p.State)
	fmt.Fprintf(&b, " %v", p.City)
	fmt.Fprintf(&b, " %v", p.Company)
	p.SearchNorm = validate.NormalizeSearch(b.String())

	return nil
}

var _ = sqlgenVtigerAccount(&VtigerAccount{})

type VtigerAccount struct {
	ID             string
	UserName       string
	FirstName      string
	RoleID         int64
	Email1         string
	Secondaryemail string
	Status         string
}

var _ = sqlgenVthCallHistory(&VthCallHistory{})

type VthCallHistory struct {
	ID              string
	CdrID           string
	CallID          string
	SipCallID       string
	SdkCallID       string
	Cause           string
	Q850Call        string
	FromExtension   string
	ToExtension     string
	FromNumber      string
	ToNumber        string
	Duration        string
	Direction       string
	TimeStarted     string
	TimeConnected   string
	RecordingPath   string
	RecordingURL    string
	EtopAcountID    int64
	VtigerAccountID string
	SearchNorm      string
}
