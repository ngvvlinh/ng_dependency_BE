package model

import (
	"time"

	"etop.vn/api/top/types/etc/status3"
	"etop.vn/capi/dot"
)

// +sqlgen
type FbPage struct {
	ID           dot.ID
	ExternalID   string
	ShopID       dot.ID
	UserID       dot.ID
	Name         string
	Category     string
	CategoryList []Category
	Tasks        []string
	Status       status3.Status
	CreatedAt    time.Time `sq:"create"`
	UpdatedAt    time.Time `sq:"update"`
	DeletedAt    time.Time
}

type Category struct {
	ID   dot.ID `json:"id"`
	Name string `json:"name"`
}

// +sqlgen
type FbPageInternal struct {
	ID        dot.ID
	Token     string
	ExpiresIn int
	UpdateAt  time.Time `sq:"update"`
}
