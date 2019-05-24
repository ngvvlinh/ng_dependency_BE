package model

import "time"

// SignedInUser ...
type SignedInUser struct {
	*User
}

// GetSignedInUserQuery ...
type GetSignedInUserQuery struct {
	UserID int64

	Result *SignedInUser
}

type GetUserByIDQuery struct {
	UserID int64

	Result *User
}

type GetUserByLoginQuery struct {
	UserID       int64
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
	UserID   int64
	Password string
}

type UpdateUserVerificationCommand struct {
	UserID int64

	EmailVerifiedAt time.Time
	PhoneVerifiedAt time.Time

	EmailVerificationSentAt time.Time
	PhoneVerificationSentAt time.Time
}

type UpdateUserIdentifierCommand struct {
	UserID    int64
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
