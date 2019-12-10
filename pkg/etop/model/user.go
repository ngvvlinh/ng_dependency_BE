package model

import (
	"time"

	"etop.vn/api/top/types/etc/status3"
	"etop.vn/api/top/types/etc/user_source"
	"etop.vn/capi/dot"
)

// SignedInUser ...
type SignedInUser struct {
	*User
}

// GetSignedInUserQuery ...
type GetSignedInUserQuery struct {
	UserID dot.ID

	Result *SignedInUser
}

type GetUserByIDQuery struct {
	UserID dot.ID

	Result *User
}

type GetUserByEmailQuery struct {
	Email string

	Result *User
}

type GetUsersByIDsQuery struct {
	UserIDs []dot.ID

	Result []*User
}

type GetUserByLoginQuery struct {
	UserID       dot.ID
	PhoneOrEmail string

	Result UserExtended
}

type CreateUserCommand struct {
	UserInner
	Password       string
	Status         status3.Status
	AgreeTOS       bool
	AgreeEmailInfo bool
	IsTest         bool
	Source         user_source.UserSource
	Result         struct {
		User         *User
		UserInternal *UserInternal
	}
}

type SetPasswordCommand struct {
	UserID   dot.ID
	Password string
}

type UpdateUserVerificationCommand struct {
	UserID dot.ID

	EmailVerifiedAt time.Time
	PhoneVerifiedAt time.Time

	EmailVerificationSentAt time.Time
	PhoneVerificationSentAt time.Time
}

type UpdateUserIdentifierCommand struct {
	UserID    dot.ID
	Status    status3.Status // We don't allow update status to 0
	UserInner                // Must be normalized identifier

	Password string

	CreatedAt       time.Time
	PhoneVerifiedAt time.Time // Automatically verify phone if the user register from phone

	Result struct {
		User *User
	}
}

type MergeUserCommand struct {
}
