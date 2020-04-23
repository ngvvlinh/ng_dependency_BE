package model

import (
	"time"

	"etop.vn/api/top/types/etc/status3"
	"etop.vn/capi/dot"
)

// +sqlgen
type FbPage struct {
	ID                   dot.ID
	ExternalID           string
	FbUserID             dot.ID
	ShopID               dot.ID
	UserID               dot.ID
	ExternalName         string
	ExternalCategory     string
	ExternalCategoryList []*ExternalCategory
	ExternalTasks        []string
	ExternalImageURL     string
	Status               status3.Status
	ConnectionStatus     status3.Status
	CreatedAt            time.Time `sq:"create"`
	UpdatedAt            time.Time `sq:"update"`
	DeletedAt            time.Time
}

type ExternalCategory struct {
	ID   dot.ID `json:"id"`
	Name string `json:"name"`
}

// +sqlgen
type FbPageInternal struct {
	ID        dot.ID
	Token     string
	UpdatedAt time.Time `sq:"update"`
}
