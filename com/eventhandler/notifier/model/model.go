package model

import (
	"encoding/json"
	"time"

	"o.o/api/top/types/etc/notifier_entity"
	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
	"o.o/common/jsonx"
)

type NotiEntity = notifier_entity.NotifierEntity

const (
	NotiFulfillment              = notifier_entity.Fulfillment
	NotiMoneyTransactionShipping = notifier_entity.MoneyTransactionShipping
	// OneSignal service ID default
	ExternalServiceOneSignalID = 101
)

// +sqlgen
type Notification struct {
	ID                dot.ID
	Title             string
	Message           string
	IsRead            bool
	EntityID          dot.ID
	Entity            NotiEntity
	AccountID         dot.ID
	SyncStatus        status3.Status
	ExternalServiceID int
	ExternalNotiID    string
	SendNotification  bool
	SyncedAt          time.Time
	SeenAt            time.Time
	CreatedAt         time.Time `sq:"create"`
	UpdatedAt         time.Time `sq:"update"`
	MetaData          json.RawMessage
}

// +sqlgen
type Device struct {
	ID dot.ID
	// DeviceID: deprecated
	DeviceID          string
	DeviceName        string
	ExternalDeviceID  string
	ExternalServiceID int
	// Name: deprecated
	AccountID     dot.ID
	UserID        dot.ID
	CreatedAt     time.Time `sq:"create"`
	UpdatedAt     time.Time `sq:"update"`
	DeactivatedAt time.Time
	Config        *DeviceConfig
}

type NotiDataAddition struct {
	Entity   NotiEntity      `json:"entity"`
	EntityID string          `json:"entity_id"`
	NotiID   string          `json:"noti_id"`
	ShopID   string          `json:"shop_id"`
	MetaData json.RawMessage `json:"meta_data"`
}

func PrepareNotiData(args *NotiDataAddition) json.RawMessage {
	dataRaw, _ := jsonx.Marshal(args)
	return dataRaw
}

type DeviceConfig struct {
	SubcribeAllShop bool     `json:"subcribe_all_shop"`
	SubcribeShopIDs []dot.ID `json:"subcribe_shop_ids"`
	Mute            bool     `json:"mute"`
}