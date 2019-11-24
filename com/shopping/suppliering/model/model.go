package model

import (
	"time"

	"etop.vn/capi/dot"
)

//go:generate $ETOPDIR/backend/scripts/derive.sh

var _ = sqlgenShopSupplier(&ShopSupplier{})

type ShopSupplier struct {
	ID                dot.ID
	ShopID            dot.ID
	FullName          string
	Phone             string
	Email             string
	Code              string
	CodeNorm          int32
	CompanyName       string
	TaxNumber         string
	HeadquaterAddress string
	Note              string
	FullNameNorm      string
	PhoneNorm         string
	Status            int32
	CreatedAt         time.Time `sq:"create"`
	UpdatedAt         time.Time `sq:"update"`
	DeletedAt         time.Time
}
