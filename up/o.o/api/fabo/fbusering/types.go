package fbusering

import (
	"time"

	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
)

type FbExternalUser struct {
	ID           dot.ID
	UserID       dot.ID
	ExternalID   string
	ExternalInfo *FbExternalUserInfo
	Status       status3.Status
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type FbExternalUserInfo struct {
	Name      string
	FirstName string
	LastName  string
	ShortName string
	ImageURL  string
}

type FbExternalUserInternal struct {
	ID        dot.ID
	Token     string
	ExpiresIn int
	UpdatedAt time.Time
}

type FbExternalUserCombined struct {
	FbExternalUser         *FbExternalUser
	FbExternalUserInternal *FbExternalUserInternal
}
