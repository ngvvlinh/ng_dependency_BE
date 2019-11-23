package model

import (
	"time"

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
	Status         Status3
	AgreeTOS       bool
	AgreeEmailInfo bool
	IsTest         bool
	IsStub         bool
	Source         UserSource
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
	Status    Status3 // We don't allow update status to 0
	UserInner         // Must be normalized identifier

	Password    string
	Identifying UserIdentifying

	CreatedAt       time.Time
	PhoneVerifiedAt time.Time // Automatically verify phone if the user register from phone

	Result struct {
		User *User
	}
}

type MergeUserCommand struct {
}
