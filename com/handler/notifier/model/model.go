package model

import (
	"encoding/json"
	"time"

	"etop.vn/api/top/types/etc/notifier_entity"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
	"etop.vn/common/jsonx"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh
//go:generate rm ../sqlstore/filters.gen.go

var _ = sqlgenNotification(&Notification{})

type NotiEntity = notifier_entity.NotifierEntity

const (
	NotiFulfillment              = notifier_entity.Fulfillment
	NotiMoneyTransactionShipping = notifier_entity.MoneyTransactionShipping
	// OneSignal service ID default
	ExternalServiceOneSignalID = 101
)

type Notification struct {
	ID                dot.ID
	Title             string
	Message           string
	IsRead            bool
	EntityID          dot.ID
	Entity            NotiEntity
	AccountID         dot.ID
	SyncStatus        model.Status3
	ExternalServiceID int
	ExternalNotiID    string
	SendNotification  bool
	SyncedAt          time.Time
	SeenAt            time.Time
	CreatedAt         time.Time `sq:"create"`
	UpdatedAt         time.Time `sq:"update"`
	MetaData          json.RawMessage
}

var _ = sqlgenDevice(&Device{})

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
	Entity   NotiEntity
	EntityID string
	NotiID   string
	ShopID   string
	MetaData json.RawMessage
}

func PrepareNotiData(args NotiDataAddition) json.RawMessage {
	dataRaw, _ := jsonx.Marshal(args)
	return dataRaw
}

type DeviceConfig struct {
	SubcribeAllShop bool     `json:"subcribe_all_shop"`
	SubcribeShopIDs []dot.ID `json:"subcribe_shop_ids"`
	Mute            bool     `json:"mute"`
}
