package fbusering

import (
	"time"

	"etop.vn/api/top/types/etc/status3"
	"etop.vn/capi/dot"
)

type FbUser struct {
	ID           dot.ID
	ExternalID   string
	UserID       dot.ID
	ExternalInfo *ExternalFBUserInfo
	Status       status3.Status
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type ExternalFBUserInfo struct {
	Name      string
	FirstName string
	LastName  string
	ShortName string
	ImageURL  string
}

type FbUserInternal struct {
	ID        dot.ID
	Token     string
	ExpiresIn int
	UpdatedAt time.Time
}

type FbUserCombined struct {
	FbUser         *FbUser
	FbUserInternal *FbUserInternal
}
