package model

import (
	"time"

	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
)

// +sqlgen
type FbUser struct {
	ID           dot.ID
	ExternalID   string
	UserID       dot.ID
	ExternalInfo *ExternalFBUserInfo
	Status       status3.Status
	CreatedAt    time.Time `sq:"create"`
	UpdatedAt    time.Time `sq:"update"`
}

type ExternalFBUserInfo struct {
	Name      string `json:"name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	ShortName string `json:"short_name"`
	ImageURL  string `json:"image_url"`
}

// +sqlgen
type FbUserInternal struct {
	ID        dot.ID
	Token     string
	ExpiresIn int
	UpdatedAt time.Time `sq:"update"`
}
