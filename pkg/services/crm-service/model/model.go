package model

import (
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
	s := p.Phone + p.Email + p.Lastname + p.Firstname +
		p.Lane + p.State + p.City + p.Company
	p.SearchNorm = validate.NormalizeSearch(s)
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

var _ = sqlgenVthCallHistory(&VhtCallHistory{})

type VhtCallHistory struct {
	CdrID           string
	CallID          string
	SipCallID       string
	SdkCallID       string
	Cause           string
	Q850Cause       string
	FromExtension   string
	ToExtension     string
	FromNumber      string
	ToNumber        string
	Duration        int32
	Direction       int32
	TimeStarted     time.Time
	TimeConnected   time.Time
	TimeEnded       time.Time
	CreatedAt       time.Time `sq:"create"`
	UpdatedAt       time.Time `sq:"update"`
	RecordingPath   string
	RecordingURL    string
	RecordFileSize  int32
	EtopAccountID   int64
	VtigerAccountID string
	SyncStatus      string
	OData           string
	SearchNorm      string
}

func (p *VhtCallHistory) BeforeInsertOrUpdate() error {
	s := p.ToNumber + " " + p.FromNumber
	p.SearchNorm = validate.NormalizeSearch(s)
	return nil
}
