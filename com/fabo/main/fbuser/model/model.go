package model

import (
	"time"

	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
)

// +sqlgen
type FbExternalUser struct {
	ExternalID     string
	ExternalInfo   *FbExternalUserInfo
	Status         status3.Status
	ExternalPageID string
	Tags           []dot.ID
	CreatedAt      time.Time `sq:"create"`
	UpdatedAt      time.Time `sq:"update"`
}

type FbExternalUserInfo struct {
	Name      string `json:"name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	ShortName string `json:"short_name"`
	ImageURL  string `json:"image_url"`
}

// +sqlgen
type FbExternalUserInternal struct {
	ExternalID string
	Token      string
	ExpiresIn  int
	UpdatedAt  time.Time `sq:"update"`
}

// +sqlgen
type FbExternalUserShopCustomer struct {
	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`

	ShopID           dot.ID
	FbExternalUserID string
	CustomerID       dot.ID
	Status           status3.Status
}

// +sqlgen
type FbShopTag struct {
	ID     dot.ID
	Name   string
	Color  string
	ShopID dot.ID

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
}
