package model

import (
	"time"

	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
)

// +sqlgen
type FbExternalPage struct {
	ID                   dot.ID
	ShopID               dot.ID
	ExternalID           string
	ExternalName         string
	ExternalTasks        []string
	ExternalCategory     string
	ExternalCategoryList []*ExternalCategory
	ExternalPermissions  []string
	ExternalImageURL     string
	ConnectionStatus     status3.Status
	Status               status3.Status
	CreatedAt            time.Time `sq:"create"`
	UpdatedAt            time.Time `sq:"update"`
	DeletedAt            time.Time
}

type ExternalCategory struct {
	ID   dot.ID `json:"id"`
	Name string `json:"name"`
}

// +sqlgen
type FbExternalPageInternal struct {
	ID         dot.ID
	ExternalID string
	Token      string
	UpdatedAt  time.Time `sq:"update"`
}
