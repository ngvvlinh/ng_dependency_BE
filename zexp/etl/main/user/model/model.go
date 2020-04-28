package model

import (
	"time"

	"o.o/api/top/types/etc/status3"
	"o.o/capi/dot"
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

	Status status3.Status `sql_type:"int2"`

	CreatedAt time.Time
	UpdatedAt time.Time

	AgreedTOSAt       time.Time
	AgreedEmailInfoAt time.Time
	EmailVerifiedAt   time.Time
	PhoneVerifiedAt   time.Time

	EmailVerificationSentAt time.Time
	PhoneVerificationSentAt time.Time

	Rid dot.ID
}
