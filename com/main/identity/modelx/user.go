package modelx

import (
	"time"

	"etop.vn/api/top/types/etc/status3"
	"etop.vn/api/top/types/etc/user_source"
	identitymodel "etop.vn/backend/com/main/identity/model"
	"etop.vn/capi/dot"
)

// SignedInUser ...
type SignedInUser struct {
	*identitymodel.User
}

// GetSignedInUserQuery ...
type GetSignedInUserQuery struct {
	UserID      dot.ID
	WLPartnerID dot.ID

	Result *SignedInUser
}

type GetUserByIDQuery struct {
	UserID      dot.ID
	WLPartnerID dot.ID

	Result *identitymodel.User
}

type GetUserByEmailOrPhoneQuery struct {
	Email       string
	Phone       string
	WLPartnerID dot.ID

	Result *identitymodel.User
}

type GetUsersByIDsQuery struct {
	UserIDs []dot.ID

	Result []*identitymodel.User
}

type GetUserByLoginQuery struct {
	UserID       dot.ID
	PhoneOrEmail string
	WLPartnerID  dot.ID

	Result identitymodel.UserExtended
}

type CreateUserCommand struct {
	identitymodel.UserInner
	Password       string
	Status         status3.Status
	AgreeTOS       bool
	AgreeEmailInfo bool
	IsTest         bool
	Source         user_source.UserSource
	WLPartnerID    dot.ID
	Result         struct {
		User         *identitymodel.User
		UserInternal *identitymodel.UserInternal
	}
}

type SetPasswordCommand struct {
	UserID   dot.ID
	Password string
}

type UpdateUserVerificationCommand struct {
	UserID      dot.ID
	WLPartnerID dot.ID

	EmailVerifiedAt time.Time
	PhoneVerifiedAt time.Time

	EmailVerificationSentAt time.Time
	PhoneVerificationSentAt time.Time
}

type UpdateUserIdentifierCommand struct {
	UserID                  dot.ID
	WLPartnerID             dot.ID
	Status                  status3.Status // We don't allow update status to 0
	identitymodel.UserInner                // Must be normalized identifier

	Password string

	CreatedAt       time.Time
	PhoneVerifiedAt time.Time // Automatically verify phone if the user register from phone

	Result struct {
		User *identitymodel.User
	}
}

type MergeUserCommand struct {
}
