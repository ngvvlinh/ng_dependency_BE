package model

import (
	"time"

	"etop.vn/api/top/types/etc/status3"
	"etop.vn/api/top/types/etc/user_source"
	"etop.vn/capi/dot"
)

type UserInner struct {
	FullName  string
	ShortName string
	Email     string
	Phone     string
}

// +sqlgen
type User struct {
	ID dot.ID

	UserInner `sq:"inline"`

	Status status3.Status

	CreatedAt time.Time
	UpdatedAt time.Time

	AgreedTOSAt       time.Time
	AgreedEmailInfoAt time.Time
	EmailVerifiedAt   time.Time
	PhoneVerifiedAt   time.Time

	EmailVerificationSentAt time.Time
	PhoneVerificationSentAt time.Time

	Source    user_source.UserSource
	RefUserID dot.ID
	RefSaleID dot.ID

	Rid dot.ID
}
