package model

import (
	"time"

	"etop.vn/backend/pkg/common/validate"
	"etop.vn/capi/dot"
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
	EtopUserID           dot.ID
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
	VtigerCreatedAt      time.Time
	VtigerUpdatedAt      time.Time
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
	RoleID         dot.ID
	Email1         string
	Secondaryemail string
	Status         string
}
