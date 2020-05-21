package model

import (
	"time"

	"o.o/api/top/types/etc/status3"
)

// +sqlgen
type FbExternalUser struct {
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
	ExternalID string
	Token      string
	ExpiresIn  int
	UpdatedAt  time.Time `sq:"update"`
}
