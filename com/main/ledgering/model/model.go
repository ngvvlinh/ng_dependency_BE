package model

import (
	"time"

	"etop.vn/backend/pkg/etop/model"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenShopLedger(&ShopLedger{})

type ShopLedger struct {
	ID          int64
	ShopID      int64
	Name        string
	BankAccount *model.BankAccount
	Note        string
	Type        string
	Status      int32
	CreatedBy   int64
	CreatedAt   time.Time `sq:"create"`
	UpdatedAt   time.Time `sq:"update"`
	DeletedAt   time.Time
}
