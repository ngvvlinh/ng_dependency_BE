package model

import (
	"encoding/json"
	"time"

	"etop.vn/api/top/types/etc/account_type"

	"etop.vn/api/top/types/etc/status3"

	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenExternalAccountAhamove(&ExternalAccountAhamove{})

type ExternalAccountAhamove struct {
	ID                  dot.ID
	OwnerID             dot.ID
	Phone               string
	Name                string
	ExternalID          string
	ExternalVerified    bool
	ExternalCreatedAt   time.Time
	ExternalToken       string
	CreatedAt           time.Time `sq:"create"`
	UpdatedAt           time.Time `sq:"update"`
	LastSendVerifiedAt  time.Time
	ExternalTicketID    string
	IDCardFrontImg      string
	IDCardBackImg       string
	PortraitImg         string
	WebsiteURL          string
	FanpageURL          string
	CompanyImgs         []string
	BusinessLicenseImgs []string

	ExternalDataVerified json.RawMessage

	UploadedAt time.Time
}

var _ = sqlgenSale(&Affiliate{})

type Affiliate struct {
	ID          dot.ID
	OwnerID     dot.ID
	Name        string
	Phone       string
	Email       string
	IsTest      int
	Status      status3.Status
	CreatedAt   time.Time `sq:"create"`
	UpdatedAt   time.Time `sq:"update"`
	DeletedAt   time.Time
	BankAccount *model.BankAccount
}

var _ model.AccountInterface = &Affiliate{}

func (s *Affiliate) GetAccount() *model.Account {
	return &model.Account{
		ID:      s.ID,
		OwnerID: s.OwnerID,
		Name:    s.Name,
		Type:    account_type.Affiliate,
	}
}
