package etop

import (
	shoptypes "o.o/api/top/int/shop/types"
	common "o.o/api/top/types/common"
	"o.o/api/top/types/etc/account_type"
	address_type "o.o/api/top/types/etc/address_type"
	"o.o/api/top/types/etc/authentication_method"
	"o.o/api/top/types/etc/credit_type"
	"o.o/api/top/types/etc/ghn_note_code"
	"o.o/api/top/types/etc/notifier_entity"
	status3 "o.o/api/top/types/etc/status3"
	try_on "o.o/api/top/types/etc/try_on"
	user_source "o.o/api/top/types/etc/user_source"
	"o.o/capi/dot"
	"o.o/capi/filter"
	"o.o/common/jsonx"
)

type GetTicketLabelsResponse struct {
	TicketLabels []*shoptypes.TicketLabel `json:"ticket_labels"`
}

func (m *GetTicketLabelsResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetTicketLabelsRequest struct {
	Tree bool `json:"tree"`
}

func (m *GetTicketLabelsRequest) String() string { return jsonx.MustMarshalToString(m) }

type Authorization struct {
	UserId  dot.ID   `json:"user_id"`
	Name    string   `json:"name"`
	Email   string   `json:"email"`
	Roles   []string `json:"roles"`
	Actions []string `json:"actions"`
}

func (m *Authorization) String() string { return jsonx.MustMarshalToString(m) }

type UpdateAuthorizationRequest struct {
	UserId dot.ID   `json:"user_id"`
	Roles  []string `json:"roles"`
}

func (m *UpdateAuthorizationRequest) String() string { return jsonx.MustMarshalToString(m) }

type AuthorizationsResponse struct {
	Authorizations []*Authorization `json:"authorizations"`
}

func (m *AuthorizationsResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetInvitationByTokenRequest struct {
	Token string `json:"token"`
}

func (m *GetInvitationByTokenRequest) String() string { return jsonx.MustMarshalToString(m) }

type Invitation struct {
	Id            dot.ID         `json:"id"`
	UserId        dot.ID         `json:"user_id"`
	ShopId        dot.ID         `json:"shop_id"`
	Phone         string         `json:"phone"`
	Email         string         `json:"email"`
	FullName      string         `json:"full_name"`
	ShortName     string         `json:"short_name"`
	ShopShort     *ShopShort     `json:"shop"`
	InvitedByUser string         `json:"invited_by_user"`
	Roles         []string       `json:"roles"`
	Token         string         `json:"token"`
	Status        status3.Status `json:"status"`
	InvitedBy     dot.ID         `json:"invited_by"`
	AcceptedAt    dot.Time       `json:"accepted_at"`
	DeclinedAt    dot.Time       `json:"declined_at"`
	ExpiresAt     dot.Time       `json:"expires_at"`
	CreatedAt     dot.Time       `json:"created_at"`
	UpdatedAt     dot.Time       `json:"updated_at"`
	InvitationURL string         `json:"invitation_url"`
}

func (m *Invitation) String() string { return jsonx.MustMarshalToString(m) }

type InvitationsResponse struct {
	Invitations []*Invitation    `json:"invitations"`
	Paging      *common.PageInfo `json:"paging"`
}

func (m *InvitationsResponse) String() string { return jsonx.MustMarshalToString(m) }

type UpdateAccountUserPermissionRequest struct {
	UserID dot.ID   `json:"user_id"`
	Roles  []string `json:"roles"`
}

func (m *UpdateAccountUserPermissionRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateRelationshipRequest struct {
	UserID    dot.ID         `json:"user_id"`
	FullName  dot.NullString `json:"full_name"`
	ShortName dot.NullString `json:"short_name"`
	Position  dot.NullString `json:"position"`
}

func (m *UpdateRelationshipRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetInvitationsRequest struct {
	Paging  *common.Paging   `json:"paging"`
	Filters []*common.Filter `json:"filters"`
}

func (m *GetInvitationsRequest) String() string { return jsonx.MustMarshalToString(m) }

type DeleteInvitationRequest struct {
	Token string `json:"token"`
}

func (m *DeleteInvitationRequest) String() string { return jsonx.MustMarshalToString(m) }

type UserRelationshipLeaveAccountRequest struct {
	AccountID dot.ID `json:"account_id"`
}

func (m *UserRelationshipLeaveAccountRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreateInvitationRequest struct {
	// @Required
	Email string `json:"email"`
	// @Required
	Phone string `json:"phone"`
	// @Optional
	FullName  string `json:"full_name"`
	ShortName string `json:"short_name"`
	Position  string `json:"position"`
	// @Required
	Roles []string `json:"roles"`
}

func (m *CreateInvitationRequest) String() string { return jsonx.MustMarshalToString(m) }

type ResendInvitationRequest struct {
	// @Required
	Email string `json:"email"`
	// @Required
	Phone string `json:"phone"`
}

func (m *ResendInvitationRequest) String() string { return jsonx.MustMarshalToString(m) }

type AcceptInvitationRequest struct {
	Token string `json:"token"`
}

func (m *AcceptInvitationRequest) String() string { return jsonx.MustMarshalToString(m) }

type RejectInvitationRequest struct {
	Token string `json:"token"`
}

func (m *RejectInvitationRequest) String() string { return jsonx.MustMarshalToString(m) }

type Relationship struct {
	UserID      dot.ID   `json:"user_id"`
	AccountID   dot.ID   `json:"account_id"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
	FullName    string   `json:"full_name"`
	ShortName   string   `json:"short_name"`
	Email       string   `json:"email"`
	Phone       string   `json:"phone"`
	Position    string   `json:"position"`
	Deleted     bool     `json:"deleted"`
}

func (m *Relationship) String() string { return jsonx.MustMarshalToString(m) }

type RelationshipsResponse struct {
	Relationships []*Relationship  `json:"relationships"`
	Paging        *common.PageInfo `json:"paging"`
}

func (m *RelationshipsResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetRelationshipsRequest struct {
	Paging  *common.Paging                `json:"paging"`
	Filters []*common.Filter              `json:"filters"`
	Filter  *FilterGetAccountUsersRequest `json:"filter"`
}

func (m *GetRelationshipsRequest) String() string { return jsonx.MustMarshalToString(m) }

type FilterGetAccountUsersRequest struct {
	Name            filter.FullTextSearch `json:"name"`
	Phone           filter.FullTextSearch `json:"phone"`
	ExtensionNumber filter.FullTextSearch `json:"extension_number"`
}

func (m *FilterGetAccountUsersRequest) String() string { return jsonx.MustMarshalToString(m) }

type RemoveUserRequest struct {
	UserID dot.ID `json:"user_id"`
}

func (m *RemoveUserRequest) String() string { return jsonx.MustMarshalToString(m) }

// Represents a user in eTop system. The user may or may not have associated accounts.
type User struct {
	// @required
	Id dot.ID `json:"id"`
	// @required
	FullName string `json:"full_name"`
	// @required
	ShortName string `json:"short_name"`
	// @required
	Phone string `json:"phone"`
	// @required
	Email string `json:"email"`
	// @required
	CreatedAt dot.Time `json:"created_at"`
	// @required
	UpdatedAt               dot.Time               `json:"updated_at"`
	EmailVerifiedAt         dot.Time               `json:"email_verified_at"`
	PhoneVerifiedAt         dot.Time               `json:"phone_verified_at"`
	EmailVerificationSentAt dot.Time               `json:"email_verification_sent_at"`
	PhoneVerificationSentAt dot.Time               `json:"phone_verification_sent_at"`
	Source                  user_source.UserSource `json:"source"`
	TotalShop               int                    `json:"total_shop"`
	IsBlocked               bool                   `json:"is_blocked"`
	BlockReason             string                 `json:"block_reason"`
	BlockedAt               dot.Time               `json:"blocked_at"`

	RefSale string `json:"ref_sale"`
	RefAff  string `json:"ref_aff"`
}

func (m *User) String() string { return jsonx.MustMarshalToString(m) }

type IDsRequest struct {
	// @required
	Ids   []dot.ID      `json:"ids"`
	Mixed *MixedAccount `json:"mixed"`
}

func (m *IDsRequest) String() string { return jsonx.MustMarshalToString(m) }

type MixedAccount struct {
	Ids      []dot.ID `json:"ids"`
	All      bool     `json:"all"`
	AllShops bool     `json:"all_shops"`
}

func (m *MixedAccount) String() string { return jsonx.MustMarshalToString(m) }

type Partner struct {
	Id             dot.ID           `json:"id"`
	Name           string           `json:"name"`
	PublicName     string           `json:"public_name"`
	Status         status3.Status   `json:"status"`
	IsTest         bool             `json:"is_test"`
	ContactPersons []*ContactPerson `json:"contact_persons"`
	Phone          string           `json:"phone"`
	WebsiteUrl     string           `json:"website_url"`
	ImageUrl       string           `json:"image_url"`
	Email          string           `json:"email"`
	OwnerId        dot.ID           `json:"owner_id"`
	User           *User            `json:"user"`
}

func (m *Partner) String() string { return jsonx.MustMarshalToString(m) }

type PublicAccountInfo struct {
	Id       dot.ID                   `json:"id"`
	Name     string                   `json:"name"`
	Type     account_type.AccountType `json:"type"`
	ImageUrl string                   `json:"image_url"`
	Website  string                   `json:"website"`
}

func (m *PublicAccountInfo) String() string { return jsonx.MustMarshalToString(m) }

type PublicAuthorizedPartnerInfo struct {
	Id          dot.ID                   `json:"id"`
	Name        string                   `json:"name"`
	Type        account_type.AccountType `json:"type"`
	ImageUrl    string                   `json:"image_url"`
	Website     string                   `json:"website"`
	RedirectUrl string                   `json:"redirect_url"`
}

func (m *PublicAuthorizedPartnerInfo) String() string { return jsonx.MustMarshalToString(m) }

type PublicAccountsResponse struct {
	Accounts []*PublicAccountInfo `json:"accounts"`
}

func (m *PublicAccountsResponse) String() string { return jsonx.MustMarshalToString(m) }

type Shop struct {
	ExportedFields     []string       `json:"exported_fields"`
	InventoryOverstock bool           `json:"inventory_overstock"`
	Id                 dot.ID         `json:"id"`
	Name               string         `json:"name"`
	Status             status3.Status `json:"status"`
	IsTest             bool           `json:"is_test"`
	Address            *Address       `json:"address"`
	Phone              string         `json:"phone"`
	BankAccount        *BankAccount   `json:"bank_account"`
	AutoCreateFfm      bool           `json:"auto_create_ffm"`
	WebsiteUrl         string         `json:"website_url"`
	ImageUrl           string         `json:"image_url"`
	Email              string         `json:"email"`
	ProductSourceId    dot.ID         `json:"product_source_id"`
	ShipToAddressId    dot.ID         `json:"ship_to_address_id"`
	ShipFromAddressId  dot.ID         `json:"ship_from_address_id"`
	IsBlocked          bool           `json:"is_blocked"`
	BlockReason        string         `json:"block_reason"`
	// @deprecated use try_on instead
	GhnNoteCode ghn_note_code.GHNNoteCode `json:"ghn_note_code,omitempty"`
	TryOn       try_on.TryOnCode          `json:"try_on"`
	OwnerId     dot.ID                    `json:"owner_id"`
	User        *User                     `json:"user"`
	CompanyInfo *CompanyInfo              `json:"company_info"`
	// referrence: https://icalendar.org/rrule-tool.html
	MoneyTransactionRrule         string                               `json:"money_transaction_rrule"`
	SurveyInfo                    []*SurveyInfo                        `json:"survey_info"`
	ShippingServiceSelectStrategy []*ShippingServiceSelectStrategyItem `json:"shipping_service_select_strategy"`
	Code                          string                               `json:"code"`
	IsPriorMoneyTransaction       dot.NullBool                         `json:"is_prior_money_transaction"`

	CreatedAt dot.Time `json:"created_at"`
	UpdatedAt dot.Time `json:"updated_at"`

	MoneyTransactionCount int `json:"money_transaction_count"`
}

func (m *Shop) String() string { return jsonx.MustMarshalToString(m) }

type ShopShort struct {
	ID       dot.ID `json:"id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	ImageUrl string `json:"image_url"`
}

func (m *ShopShort) String() string { return jsonx.MustMarshalToString(m) }

type ShippingServiceSelectStrategyItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (m *ShippingServiceSelectStrategyItem) String() string { return jsonx.MustMarshalToString(m) }

type SurveyInfo struct {
	Key      string `json:"key"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

func (m *SurveyInfo) String() string { return jsonx.MustMarshalToString(m) }

type CreateUserRequest struct {
	// @required
	FullName string `json:"full_name"`
	// Can be automatically deduce from full_name.
	ShortName string `json:"short_name"`
	// @required
	Phone string `json:"phone"`
	// It's not required if the user provides register_token
	Email string `json:"email"`
	// @required
	Password string `json:"password"`
	// @required
	AgreeTos bool `json:"agree_tos"`
	// @required
	AgreeEmailInfo dot.NullBool `json:"agree_email_info"`
	// This field must be set if the user uses generated password to register.
	// Automatically set phone_verified if it's sent within a specific time.
	RegisterToken string                 `json:"register_token"`
	Source        user_source.UserSource `json:"source"`

	AutoAcceptInvitation bool `json:"auto_accept_invitation"`

	RefSale string `json:"ref_sale"`
	RefAff  string `json:"ref_aff"`
}

func (m *CreateUserRequest) String() string { return jsonx.MustMarshalToString(m) }

type RequestRegisterSimplifyRequest struct {
	// @required
	Phone string `json:"phone"`
}

func (m *RequestRegisterSimplifyRequest) String() string { return jsonx.MustMarshalToString(m) }

type RegisterSimplifyRequest struct {
	// @required
	Phone string `json:"phone"`
	// @required
	OTP string `json:"otp"`
}

func (m *RegisterSimplifyRequest) String() string { return jsonx.MustMarshalToString(m) }

type RegisterResponse struct {
	// @required
	User *User `json:"user"`
}

func (m *RegisterResponse) String() string { return jsonx.MustMarshalToString(m) }

// Exactly one of phone or email must be provided.
type ResetPasswordRequest struct {
	// Email address to send reset password instruction.
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	RecaptchaToken string `json:"recaptcha_token"`
}

func (m *ResetPasswordRequest) String() string { return jsonx.MustMarshalToString(m) }
func (m *ResetPasswordRequest) RequireCaptcha() bool {
	return m.Phone != ""
}

// Exactly one of current_password or reset_password_token must be provided.
type ChangePasswordRequest struct {
	// @required
	CurrentPassword string `json:"current_password"`
	// @required
	NewPassword string `json:"new_password"`
	// @required
	ConfirmPassword string `json:"confirm_password"`
}

func (m *ChangePasswordRequest) String() string { return jsonx.MustMarshalToString(m) }

type ChangeUserRefAffRequest struct {
	RefAff dot.NullString `json:"ref_aff"`
}

func (c *ChangeUserRefAffRequest) String() string { return jsonx.MustMarshalToString(c) }

// Exactly one of email or phone must be provided.
type ChangePasswordUsingTokenRequest struct {
	// @required
	ResetPasswordToken string `json:"reset_password_token"`
	// @required
	NewPassword string `json:"new_password"`
	// @required
	ConfirmPassword string `json:"confirm_password"`
}

func (m *ChangePasswordUsingTokenRequest) String() string { return jsonx.MustMarshalToString(m) }

type ChangePasswordForPhoneUsingTokenRequest struct {
	// @required
	NewPassword string `json:"new_password"`
	// @required
	ConfirmPassword string `json:"confirm_password"`
}

func (m *ChangePasswordForPhoneUsingTokenRequest) Reset() {
	*m = ChangePasswordForPhoneUsingTokenRequest{}
}
func (m *ChangePasswordForPhoneUsingTokenRequest) String() string {
	return jsonx.MustMarshalToString(m)
}

// Represents permission of the current user relation with an account.
type Permission struct {
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}

func (p *Permission) GetRoles() []string {
	if p != nil {
		return p.Roles
	}
	return nil
}

func (p *Permission) GetPermissions() []string {
	if p != nil {
		return p.Permissions
	}
	return nil
}

func (p *Permission) String() string { return jsonx.MustMarshalToString(p) }

type LoginRequest struct {
	// @required Phone or email
	Login string `json:"login"`
	// @required
	Password string `json:"password"`
	// Automatically switch to this account if available.
	//
	// It's *ignored* if the user *does not* have permission to this account.
	AccountId dot.ID `json:"account_id"`
	// Automatically switch to the only account of this account type if available.
	//
	// It's *ignored* if the user *does not* have any account of this account type, or if the user *has more than one* account of this account type.
	AccountType account_type.AccountType `json:"account_type"`
	// Not implemented.
	AccountKey string `json:"account_key"`
}

func (m *LoginRequest) String() string { return jsonx.MustMarshalToString(m) }

// Represents an account associated with the current user. It has extra fields to represents relation with the user.
type LoginAccount struct {
	ExportedFields []string `json:"exported_fields"`
	// @required
	Id dot.ID `json:"id"`
	// @required
	Name string `json:"name"`
	// @required
	Type account_type.AccountType `json:"type"`
	// Associated token for the account. It's returned when calling Login or
	// SwitchAccount with regenerate_tokens set to true.
	AccessToken string `json:"access_token"`
	// The same as access_token.
	ExpiresIn   int              `json:"expires_in"`
	ImageUrl    string           `json:"image_url"`
	UrlSlug     string           `json:"url_slug"`
	UserAccount *UserAccountInfo `json:"user_account"`
}

func (m *LoginAccount) String() string { return jsonx.MustMarshalToString(m) }

type LoginResponse struct {
	// @required
	AccessToken string `json:"access_token"`
	// @required
	ExpiresIn int `json:"expires_in"`
	// @required
	User      *User         `json:"user"`
	Account   *LoginAccount `json:"account"`
	Shop      *Shop         `json:"shop"`
	Affiliate *Affiliate    `json:"affiliate"`
	// @required
	AvailableAccounts  []*LoginAccount    `json:"available_accounts"`
	InvitationAccounts []*UserAccountInfo `json:"invitation_accounts"`
	Stoken             bool               `json:"stoken"`
	StokenExpiresAt    dot.Time           `json:"stoken_expires_at"`
}

func (m *LoginResponse) String() string { return jsonx.MustMarshalToString(m) }

type SwitchAccountRequest struct {
	// @required
	AccountId dot.ID `json:"account_id"`
	// This field should only be used after creating new accounts. If it is set,
	// account_id can be left empty.
	RegenerateTokens bool `json:"regenerate_tokens"`
}

func (m *SwitchAccountRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateUserEmailRequest struct {
	Email string `json:"email"`

	FirstCode  string `json:"first_code"`
	SecondCode string `json:"second_code"`

	AuthenticationMethod authentication_method.AuthenticationMethod `json:"authentication_method"`
}

func (m *UpdateUserEmailRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateUserEmailResponse struct {
	Msg string `json:"msg"`
}

func (m *UpdateUserEmailResponse) String() string { return jsonx.MustMarshalToString(m) }

type UpdateUserPhoneRequest struct {
	Phone string `json:"phone"`

	FirstCode  string `json:"first_code"`
	SecondCode string `json:"second_code"`

	AuthenticationMethod authentication_method.AuthenticationMethod `json:"authentication_method"`
}

func (m *UpdateUserPhoneRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateUserPhoneResponse struct {
	Msg string `json:"msg"`
}

func (m *UpdateUserPhoneResponse) String() string { return jsonx.MustMarshalToString(m) }

type UpgradeAccessTokenRequest struct {
	// @required
	Stoken   string `json:"stoken"`
	Password string `json:"password"`
}

func (m *UpgradeAccessTokenRequest) String() string { return jsonx.MustMarshalToString(m) }

type AccessTokenResponse struct {
	// @required
	AccessToken string `json:"access_token"`
	// @required
	ExpiresIn int `json:"expires_in"`
	// @required
	User            *User         `json:"user"`
	Account         *LoginAccount `json:"account"`
	Shop            *Shop         `json:"shop"`
	Affiliate       *Affiliate    `json:"affiliate"`
	Stoken          bool          `json:"stoken"`
	StokenExpiresAt dot.Time      `json:"stoken_expires_at"`
}

func (m *AccessTokenResponse) String() string { return jsonx.MustMarshalToString(m) }

type SendSTokenEmailRequest struct {
	// @required
	Email string `json:"email"`
	// @required
	AccountId dot.ID `json:"account_id"`
}

func (m *SendSTokenEmailRequest) String() string { return jsonx.MustMarshalToString(m) }

type SendEmailVerificationRequest struct {
	// @required
	Email string `json:"email"`
}

func (m *SendEmailVerificationRequest) String() string { return jsonx.MustMarshalToString(m) }

type SendEmailVerificationUsingOTPRequest struct {
	// @required
	Email string `json:"email"`
}

func (m *SendEmailVerificationUsingOTPRequest) String() string { return jsonx.MustMarshalToString(m) }

type SendPhoneVerificationRequest struct {
	// @required
	Phone          string `json:"phone"`
	RecaptchaToken string `json:"recaptcha_token"`
}

func (m *SendPhoneVerificationRequest) String() string { return jsonx.MustMarshalToString(m) }

type VerifyEmailUsingTokenRequest struct {
	// @required
	VerificationToken string `json:"verification_token"`
}

func (m *VerifyEmailUsingTokenRequest) String() string { return jsonx.MustMarshalToString(m) }

type VerifyEmailUsingOTPRequest struct {
	// @required
	VerificationToken string `json:"verification_token"`
}

func (m *VerifyEmailUsingOTPRequest) String() string { return jsonx.MustMarshalToString(m) }

type VerifyPhoneUsingTokenRequest struct {
	// @required
	VerificationToken string `json:"verification_token"`
}

func (m *VerifyPhoneUsingTokenRequest) String() string { return jsonx.MustMarshalToString(m) }

type VerifyPhoneResetPasswordUsingTokenRequest struct {
	// @required
	VerificationToken string `json:"verification_token"`
}

func (m *VerifyPhoneResetPasswordUsingTokenRequest) String() string {
	return jsonx.MustMarshalToString(m)
}

type VerifyPhoneResetPasswordUsingTokenResponse struct {
	ResetPasswordToken string `json:"reset_password_token"`
}

func (m *VerifyPhoneResetPasswordUsingTokenResponse) String() string {
	return jsonx.MustMarshalToString(m)
}

type InviteUserToAccountRequest struct {
	// @required account to manage, must be owner of the account
	AccountId dot.ID `json:"account_id"`
	// @required phone or email
	InviteeIdentifier string      `json:"invitee_identifier"`
	Permission        *Permission `json:"permission"`
}

func (m *InviteUserToAccountRequest) String() string { return jsonx.MustMarshalToString(m) }

type AnswerInvitationRequest struct {
	AccountId dot.ID             `json:"account_id"`
	Response  status3.NullStatus `json:"response"`
}

func (m *AnswerInvitationRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetUsersInCurrentAccountsRequest struct {
	Paging  *common.Paging   `json:"paging"`
	Filters []*common.Filter `json:"filters"`
	Mixed   *MixedAccount    `json:"mixed"`
}

func (m *GetUsersInCurrentAccountsRequest) String() string { return jsonx.MustMarshalToString(m) }

type PublicUserInfo struct {
	// @required
	Id dot.ID `json:"id"`
	// @required
	FullName string `json:"full_name"`
	// @required
	ShortName string `json:"short_name"`
}

func (m *PublicUserInfo) String() string { return jsonx.MustMarshalToString(m) }

// Presents user information inside an account
type UserAccountInfo struct {
	// @required
	UserId dot.ID `json:"user_id"`
	// @required
	UserFullName string `json:"user_full_name"`
	// @required
	UserShortName string `json:"user_short_name"`
	// @required
	AccountId dot.ID `json:"account_id"`
	// @required
	AccountName string `json:"account_name"`
	// @required
	AccountType account_type.AccountType `json:"account_type"`
	Position    string                   `json:"position"`
	// @required
	Permission           *Permission    `json:"permission"`
	Status               status3.Status `json:"status"`
	ResponseStatus       status3.Status `json:"response_status"`
	InvitationSentBy     dot.ID         `json:"invitation_sent_by"`
	InvitationSentAt     dot.Time       `json:"invitation_sent_at"`
	InvitationAcceptedAt dot.Time       `json:"invitation_accepted_at"`
	InvitationRejectedAt dot.Time       `json:"invitation_rejected_at"`
	DisabledAt           dot.Time       `json:"disabled_at"`
	DisabledBy           dot.ID         `json:"disabled_by"`
	DisableReason        string         `json:"disable_reason"`
}

func (m *UserAccountInfo) String() string { return jsonx.MustMarshalToString(m) }

// Prepsents users in current account
type ProtectedUsersResponse struct {
	Paging *common.PageInfo   `json:"paging"`
	Users  []*UserAccountInfo `json:"users"`
}

func (m *ProtectedUsersResponse) String() string { return jsonx.MustMarshalToString(m) }

type LeaveAccountRequest struct {
}

func (m *LeaveAccountRequest) String() string { return jsonx.MustMarshalToString(m) }

type RemoveUserFromCurrentAccountRequest struct {
}

func (m *RemoveUserFromCurrentAccountRequest) String() string { return jsonx.MustMarshalToString(m) }

type Province struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Region string `json:"region"`
}

func (m *Province) String() string { return jsonx.MustMarshalToString(m) }

type GetProvincesResponse struct {
	Provinces []*Province `json:"provinces"`
}

func (m *GetProvincesResponse) String() string { return jsonx.MustMarshalToString(m) }

type District struct {
	Code         string `json:"code"`
	ProvinceCode string `json:"province_code"`
	Name         string `json:"name"`
	IsFreeship   bool   `json:"is_freeship"`
}

func (m *District) String() string { return jsonx.MustMarshalToString(m) }

type GetDistrictsResponse struct {
	Districts []*District `json:"districts"`
}

func (m *GetDistrictsResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetDistrictsByProvinceRequest struct {
	ProvinceCode string `json:"province_code"`
	ProvinceName string `json:"province_name"`
}

func (m *GetDistrictsByProvinceRequest) String() string { return jsonx.MustMarshalToString(m) }

type Ward struct {
	Code         string `json:"code"`
	DistrictCode string `json:"district_code"`
	Name         string `json:"name"`
}

func (m *Ward) String() string { return jsonx.MustMarshalToString(m) }

type GetWardsResponse struct {
	Wards []*Ward `json:"wards"`
}

func (m *GetWardsResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetWardsByDistrictRequest struct {
	DistrictCode string `json:"district_code"`
	DistrictName string `json:"district_name"`
}

func (m *GetWardsByDistrictRequest) String() string { return jsonx.MustMarshalToString(m) }

type ParseLocationRequest struct {
	ProvinceName string `json:"province_name"`
	DistrictName string `json:"district_name"`
	WardName     string `json:"ward_name"`
}

func (m *ParseLocationRequest) String() string { return jsonx.MustMarshalToString(m) }

type ParseLocationResponse struct {
	Province *Province `json:"province"`
	District *District `json:"district"`
	Ward     *Ward     `json:"ward"`
}

func (m *ParseLocationResponse) String() string { return jsonx.MustMarshalToString(m) }

type ResetPasswordResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	// @required
	AccessToken string `json:"access_token"`
	// @required
	ExpiresIn int `json:"expires_in"`
}

func (m *ResetPasswordResponse) String() string { return jsonx.MustMarshalToString(m) }

type Bank struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func (m *Bank) String() string { return jsonx.MustMarshalToString(m) }

type GetBanksResponse struct {
	Banks []*Bank `json:"banks"`
}

func (m *GetBanksResponse) String() string { return jsonx.MustMarshalToString(m) }

type BankProvince struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	BankCode string `json:"bank_code"`
}

func (m *BankProvince) String() string { return jsonx.MustMarshalToString(m) }

type GetBankProvincesResponse struct {
	Provinces []*BankProvince `json:"provinces"`
}

func (m *GetBankProvincesResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetProvincesByBankResquest struct {
	BankCode string `json:"bank_code"`
	BankName string `json:"bank_name"`
}

func (m *GetProvincesByBankResquest) String() string { return jsonx.MustMarshalToString(m) }

type BankBranch struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	BankCode     string `json:"bank_code"`
	ProvinceCode string `json:"province_code"`
}

func (m *BankBranch) String() string { return jsonx.MustMarshalToString(m) }

type GetBranchesByBankProvinceResponse struct {
	Branches []*BankBranch `json:"branches"`
}

func (m *GetBranchesByBankProvinceResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetBranchesByBankProvinceResquest struct {
	BankCode     string `json:"bank_code"`
	BankName     string `json:"bank_name"`
	ProvinceCode string `json:"province_code"`
	ProvinceName string `json:"province_name"`
}

func (m *GetBranchesByBankProvinceResquest) String() string { return jsonx.MustMarshalToString(m) }

type Address struct {
	ExportedFields []string `json:"exported_fields"`
	Id             dot.ID   `json:"id"`
	Province       string   `json:"province"`
	ProvinceCode   string   `json:"province_code"`
	District       string   `json:"district"`
	DistrictCode   string   `json:"district_code"`
	Ward           string   `json:"ward"`
	WardCode       string   `json:"ward_code"`
	Address1       string   `json:"address1"`
	Address2       string   `json:"address2"`
	Zip            string   `json:"zip"`
	Country        string   `json:"country"`
	FullName       string   `json:"full_name"`
	// deprecated: use full_name instead
	FirstName string `json:"first_name"`
	// deprecated: use full_name instead
	LastName    string                   `json:"last_name"`
	Phone       string                   `json:"phone"`
	Email       string                   `json:"email"`
	Position    string                   `json:"position"`
	Type        address_type.AddressType `json:"type"`
	Notes       *AddressNote             `json:"notes"`
	Coordinates *Coordinates             `json:"coordinates"`
}

func (m *Address) String() string { return jsonx.MustMarshalToString(m) }

type Coordinates struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

func (m *Coordinates) String() string { return jsonx.MustMarshalToString(m) }

type BankAccount struct {
	Name          string `json:"name"`
	Province      string `json:"province"`
	Branch        string `json:"branch"`
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
}

func (m *BankAccount) String() string { return jsonx.MustMarshalToString(m) }

type ContactPerson struct {
	Name     string `json:"name"`
	Position string `json:"position"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

func (m *ContactPerson) String() string { return jsonx.MustMarshalToString(m) }

type CompanyInfo struct {
	Name                string         `json:"name"`
	TaxCode             string         `json:"tax_code"`
	Address             string         `json:"address"`
	Website             string         `json:"website"`
	LegalRepresentative *ContactPerson `json:"legal_representative"`
}

func (m *CompanyInfo) String() string { return jsonx.MustMarshalToString(m) }

type CreateAddressRequest struct {
	// @required
	Province     string `json:"province"`
	ProvinceCode string `json:"province_code"`
	// @required
	District     string `json:"district"`
	DistrictCode string `json:"district_code"`
	// @required
	Ward        string                   `json:"ward"`
	WardCode    string                   `json:"ward_code"`
	Address1    string                   `json:"address1"`
	Address2    string                   `json:"address2"`
	Zip         string                   `json:"zip"`
	Country     string                   `json:"country"`
	FullName    string                   `json:"full_name"`
	FirstName   string                   `json:"first_name"`
	LastName    string                   `json:"last_name"`
	Phone       string                   `json:"phone"`
	Email       string                   `json:"email"`
	Position    string                   `json:"position"`
	Type        address_type.AddressType `json:"type"`
	Notes       *AddressNote             `json:"notes"`
	Coordinates *Coordinates             `json:"coordinates"`
}

func (m *CreateAddressRequest) String() string { return jsonx.MustMarshalToString(m) }

type AddressNote struct {
	Note       string `json:"note"`
	OpenTime   string `json:"open_time"`
	LunchBreak string `json:"lunch_break"`
	Other      string `json:"other"`
}

func (m *AddressNote) String() string { return jsonx.MustMarshalToString(m) }

type GetAddressResponse struct {
	Addresses []*Address `json:"addresses"`
}

func (m *GetAddressResponse) String() string { return jsonx.MustMarshalToString(m) }

type UpdateAddressRequest struct {
	Id dot.ID `json:"id"`
	// @required
	Province     string `json:"province"`
	ProvinceCode string `json:"province_code"`
	// @required
	District     string `json:"district"`
	DistrictCode string `json:"district_code"`
	// @required
	Ward        string                   `json:"ward"`
	WardCode    string                   `json:"ward_code"`
	Address1    string                   `json:"address1"`
	Address2    string                   `json:"address2"`
	Zip         string                   `json:"zip"`
	Country     string                   `json:"country"`
	FullName    string                   `json:"full_name"`
	FirstName   string                   `json:"first_name"`
	LastName    string                   `json:"last_name"`
	Phone       string                   `json:"phone"`
	Email       string                   `json:"email"`
	Position    string                   `json:"position"`
	Type        address_type.AddressType `json:"type"`
	Notes       *AddressNote             `json:"notes"`
	Coordinates *Coordinates             `json:"coordinates"`
}

func (m *UpdateAddressRequest) String() string { return jsonx.MustMarshalToString(m) }

type SetDefaultAddressRequest struct {
	Id   dot.ID                   `json:"id"`
	Type address_type.AddressType `json:"type"`
}

func (m *SetDefaultAddressRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateURLSlugRequest struct {
	AccountId dot.ID         `json:"account_id"`
	UrlSlug   dot.NullString `json:"url_slug"`
}

func (m *UpdateURLSlugRequest) String() string { return jsonx.MustMarshalToString(m) }

type HistoryResponse struct {
	Paging *common.PageInfo      `json:"paging"`
	Data   *common.RawJSONObject `json:"data"`
}

func (m *HistoryResponse) String() string { return jsonx.MustMarshalToString(m) }

type Credit struct {
	ID        dot.ID                     `json:"id"`
	Amount    int                        `json:"amount"`
	ShopID    dot.ID                     `json:"shop_id"`
	Type      credit_type.CreditType     `json:"type"`
	Shop      *Shop                      `json:"shop"`
	CreatedAt dot.Time                   `json:"created_at"`
	UpdatedAt dot.Time                   `json:"updated_at"`
	PaidAt    dot.Time                   `json:"paid_at"`
	Status    status3.Status             `json:"status"`
	Classify  credit_type.CreditClassify `json:"classify"`
}

func (m *Credit) String() string { return jsonx.MustMarshalToString(m) }

type CreditsResponse struct {
	Paging  *common.PageInfo `json:"paging"`
	Credits []*Credit        `json:"credits"`
}

func (m *CreditsResponse) String() string { return jsonx.MustMarshalToString(m) }

type UpdatePermissionRequest struct {
	Items []*UpdatePermissionItem `json:"items"`
}

func (m *UpdatePermissionRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdatePermissionItem struct {
	Type       string   `json:"type"`
	Key        string   `json:"key"`
	Grants     []string `json:"grants"`
	Revokes    []string `json:"revokes"`
	ReplaceAll []string `json:"replace_all"`
	RevokeAll  bool     `json:"revoke_all"`
}

func (m *UpdatePermissionItem) String() string { return jsonx.MustMarshalToString(m) }

type UpdatePermissionResponse struct {
	Msg string `json:"msg"`
}

func (m *UpdatePermissionResponse) String() string { return jsonx.MustMarshalToString(m) }

type Device struct {
	Id                dot.ID   `json:"id"`
	DeviceId          string   `json:"device_id"`
	DeviceName        string   `json:"device_name"`
	ExternalDeviceId  string   `json:"external_device_id"`
	ExternalServiceId int      `json:"external_service_id"`
	AccountId         dot.ID   `json:"account_id"`
	CreatedAt         dot.Time `json:"created_at"`
	UpdatedAt         dot.Time `json:"updated_at"`
}

func (m *Device) String() string { return jsonx.MustMarshalToString(m) }

type NotificationMeta struct {
	ConversationID      string `json:"conversation_id,omitempty"`
	FbExternalPostID    string `json:"fb_external_post_id,omitempty"`
	FbExternalCommentID string `json:"fb_external_comment_id,omitempty"`
}

type Notification struct {
	Id               dot.ID            `json:"id"`
	Title            string            `json:"title"`
	Message          string            `json:"message"`
	IsRead           bool              `json:"is_read"`
	Entity           string            `json:"entity"`
	EntityId         dot.ID            `json:"entity_id"`
	AccountId        dot.ID            `json:"account_id"`
	SendNotification bool              `json:"send_notification"`
	SeenAt           dot.Time          `json:"seen_at"`
	CreatedAt        dot.Time          `json:"created_at"`
	UpdatedAt        dot.Time          `json:"updated_at"`
	SyncStatus       status3.Status    `json:"sync_status"`
	MetaData         *NotificationMeta `json:"meta_data"`
}

func (m *Notification) String() string { return jsonx.MustMarshalToString(m) }

type CreateDeviceRequest struct {
	DeviceId         string `json:"device_id"`
	DeviceName       string `json:"device_name"`
	ExternalDeviceId string `json:"external_device_id"`
}

func (m *CreateDeviceRequest) String() string { return jsonx.MustMarshalToString(m) }

type DeleteDeviceRequest struct {
	DeviceId         string `json:"device_id"`
	ExternalDeviceId string `json:"external_device_id"`
}

func (m *DeleteDeviceRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetNotificationFilter struct {
	Entity notifier_entity.NotifierEntity `json:"entity"`
}

type GetNotificationsRequest struct {
	Paging *common.Paging         `json:"paging"`
	Filter *GetNotificationFilter `json:"filter"`
}

func (m *GetNotificationsRequest) String() string { return jsonx.MustMarshalToString(m) }

type NotificationsResponse struct {
	Paging        *common.PageInfo `json:"paging"`
	Notifications []*Notification  `json:"notifications"`
}

func (m *NotificationsResponse) String() string { return jsonx.MustMarshalToString(m) }

type UpdateNotificationsRequest struct {
	Ids    []dot.ID `json:"ids"`
	IsRead bool     `json:"is_read"`
}

func (m *UpdateNotificationsRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateReferenceUserRequest struct {
	// @required
	Phone string `json:"phone"`
}

func (m *UpdateReferenceUserRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateReferenceSaleRequest struct {
	// @required
	Phone string `json:"phone"`
}

func (m *UpdateReferenceSaleRequest) String() string { return jsonx.MustMarshalToString(m) }

type Affiliate struct {
	Id          dot.ID         `json:"id"`
	Name        string         `json:"name"`
	Status      status3.Status `json:"status"`
	IsTest      bool           `json:"is_test"`
	Phone       string         `json:"phone"`
	Email       string         `json:"email"`
	BankAccount *BankAccount   `json:"bank_account"`
}

func (m *Affiliate) String() string { return jsonx.MustMarshalToString(m) }

type GetUserByPhoneRequest struct {
	// @required
	Phone          string `json:"phone"`
	RecaptchaToken string `json:"recaptcha_token"`
}

func (m *GetUserByPhoneRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetUserByPhoneResponse struct {
	// @required
	Exists bool `json:"exists"`
}

func (m *GetUserByPhoneResponse) String() string { return jsonx.MustMarshalToString(m) }

type SendPhoneVerifyRequest struct {
	// @required
	Phone          string `json:"phone"`
	RecaptchaToken string `json:"recaptcha_token"`
}

func (m *SendPhoneVerifyRequest) String() string { return jsonx.MustMarshalToString(m) }

func (m *CreateUserRequest) Censor() {
	if m.Password != "" {
		m.Password = "..."
	}
	if m.RegisterToken != "" {
		m.RegisterToken = "..."
	}
}

func (m *LoginRequest) Censor() {
	if m.Password != "" {
		m.Password = "..."
	}
}

func (m *ChangePasswordRequest) Censor() {
	if m.CurrentPassword != "" {
		m.CurrentPassword = "..."
	}
	if m.NewPassword != "" {
		m.NewPassword = "..."
	}
	if m.ConfirmPassword != "" {
		m.ConfirmPassword = "..."
	}
}

func (m *ChangePasswordUsingTokenRequest) Censor() {
	if m.ResetPasswordToken != "" {
		m.ResetPasswordToken = "..."
	}
	if m.NewPassword != "" {
		m.NewPassword = "..."
	}
	if m.ConfirmPassword != "" {
		m.ConfirmPassword = "..."
	}
}

type EcomSessionInfoResponse struct {
	AllowAccess bool `json:"allow_access"`
}

func (m *EcomSessionInfoResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetBankProvincesRequest struct {
	All      bool   `json:"all"`
	BankCode string `json:"bank_code"`
	BankName string `json:"bank_name"`
}

func (m *GetBankProvincesRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetBankProvinceResponse struct {
	Provinces []*BankProvince `json:"provinces"`
}

func (m *GetBankProvinceResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetBankBranchesRequest struct {
	All      bool   `json:"all"`
	BankCode string `json:"bank_code"`
	BankName string `json:"bank_name"`
}

func (m *GetBankBranchesRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetBankBranchesResponse struct {
	Branches []*BankBranch `json:"branches"`
}

func (m *GetBankBranchesResponse) String() string { return jsonx.MustMarshalToString(m) }

type NotifyTopic struct {
	Topic  string `json:"topic"`
	Enable bool   `json:"enable"`
}

type GetNotifySettingResponse struct {
	UserID dot.ID         `json:"user_id"`
	Topics []*NotifyTopic `json:"topics"`
}

func (m *GetNotifySettingResponse) String() string { return jsonx.MustMarshalToString(m) }

type UpdateNotifyTopicRequest struct {
	Topic string `json:"topic"`
}

func (m *UpdateNotifyTopicRequest) String() string { return jsonx.MustMarshalToString(m) }

type WebphoneRequestLoginRequest struct {
	Phone string `json:"phone"`
}

func (m *WebphoneRequestLoginRequest) String() string { return jsonx.MustMarshalToString(m) }

type WebphoneRequestLoginResponse struct {
	SecretKey string `json:"secret_key"`
}

func (m *WebphoneRequestLoginResponse) String() string { return jsonx.MustMarshalToString(m) }

type WebphoneLoginRequest struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

func (m *WebphoneLoginRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetAuthCodeURLResponse struct {
	AuthURL string `json:"auth_url"`
}

func (m *GetAuthCodeURLResponse) String() string { return jsonx.MustMarshalToString(m) }

type VerifyTokenUsingCodeRequest struct {
	Code string `json:"code"`
}

func (m *VerifyTokenUsingCodeRequest) String() string { return jsonx.MustMarshalToString(m) }
