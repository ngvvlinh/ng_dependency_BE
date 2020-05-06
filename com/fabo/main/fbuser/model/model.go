package model

import (
	"time"

	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
)

// +sqlgen
type FbExternalUser struct {
	ID           dot.ID
	UserID       dot.ID
	ExternalID   string
	ExternalInfo *FbExternalUserInfo
	Status       status3.Status
	CreatedAt    time.Time `sq:"create"`
	UpdatedAt    time.Time `sq:"update"`
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
	ID        dot.ID
	Token     string
	ExpiresIn int
	UpdatedAt time.Time `sq:"update"`
}
