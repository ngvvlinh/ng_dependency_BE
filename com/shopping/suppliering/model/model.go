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
	CodeNorm          int
	CompanyName       string
	TaxNumber         string
	HeadquaterAddress string
	Note              string
	FullNameNorm      string
	PhoneNorm         string
	Status            int
	CreatedAt         time.Time `sq:"create"`
	UpdatedAt         time.Time `sq:"update"`
	DeletedAt         time.Time
}
