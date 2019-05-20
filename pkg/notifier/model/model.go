package model

import (
	"encoding/json"
	"time"

	"etop.vn/backend/pkg/etop/model"
)

//go:generate ../../../scripts/derive.sh

var _ = sqlgenNotification(&Notification{})

type NotiEntity string

const (
	NotiFulfillment              NotiEntity = "fulfillment"
	NotiMoneyTransactionShipping NotiEntity = "money_transaction_shipping"
	// OneSignal service ID default
	ExternalServiceOneSignalID = 101
)

type Notification struct {
	ID                int64
	Title             string
	Message           string
	IsRead            bool
	EntityID          int64
	Entity            NotiEntity
	AccountID         int64
	SyncStatus        model.Status3
	ExternalServiceID int
	ExternalNotiID    string
	SendNotification  bool
	SyncedAt          time.Time
	SeenAt            time.Time
	CreatedAt         time.Time `sq:"create"`
	UpdatedAt         time.Time `sq:"update"`
}

var _ = sqlgenDevice(&Device{})

type Device struct {
	ID                int64
	DeviceID          string
	DeviceName        string
	ExternalDeviceID  string
	ExternalServiceID int
	AccountID         int64
	CreatedAt         time.Time `sq:"create"`
	UpdatedAt         time.Time `sq:"update"`
}

type NotiDataAddition struct {
	Entity   NotiEntity
	EntityID string
	NotiID   string
}

func PrepareNotiData(args NotiDataAddition) json.RawMessage {
	dataRaw, _ := json.Marshal(args)
	return dataRaw
}
