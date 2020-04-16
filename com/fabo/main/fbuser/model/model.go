package model

import (
	"time"

	"etop.vn/api/top/types/etc/status3"
	"etop.vn/capi/dot"
)

// +sqlgen
type FbUser struct {
	ID         dot.ID
	ExternalID string
	UserID     dot.ID
	Info       FBUserInfo `json:"info"`
	Status     status3.Status
	CreatedAt  time.Time `sq:"create"`
	UpdatedAt  time.Time `sq:"update"`
}

type FBUserInfo struct {
	Name      string `json:"name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	ShortName string `json:"short_name"`
	ImgURL    string `json:"img_url"`
}

// +sqlgen
type FbUserInternal struct {
	ID        dot.ID
	Token     string
	ExpiresIn int
	UpdatedAt time.Time `sq:"update"`
}
