package model

import "time"

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenExternalAccountAhamove(&ExternalAccountAhamove{})

type ExternalAccountAhamove struct {
	ID                 int64
	OwnerID            int64
	Phone              string
	Name               string
	ExternalID         string
	ExternalVerified   bool
	ExternalCreatedAt  time.Time
	ExternalToken      string
	CreatedAt          time.Time `sq:"create"`
	UpdatedAt          time.Time `sq:"update"`
	LastSendVerifiedAt time.Time
	ExternalTicketID   string
	IDCardFrontImg     string
	IDCardBackImg      string
	PortraitImg        string
	UploadedAt         time.Time
}
