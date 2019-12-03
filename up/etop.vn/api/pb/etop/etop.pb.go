package etop

import (
	common "etop.vn/api/pb/common"
	address_type "etop.vn/api/pb/etop/etc/address_type"
	status3 "etop.vn/api/pb/etop/etc/status3"
	try_on "etop.vn/api/pb/etop/etc/try_on"
	user_source "etop.vn/api/pb/etop/etc/user_source"
	"etop.vn/capi/dot"
	"etop.vn/common/jsonx"
)

// Indicates whether given account is **etop**, **shop**, **partner** or **sale**.
type AccountType int

const (
	AccountType_unknown   AccountType = 0
	AccountType_partner   AccountType = 21
	AccountType_shop      AccountType = 33
	AccountType_affiliate AccountType = 35
	AccountType_etop      AccountType = 101
)

var AccountType_name = map[int]string{
	0:   "unknown",
	21:  "partner",
	33:  "shop",
	35:  "affiliate",
	101: "etop",
}

var AccountType_value = map[string]int{
	"unknown":   0,
	"partner":   21,
	"shop":      33,
	"affiliate": 35,
	"etop":      101,
}

func (x AccountType) Enum() *AccountType {
	p := new(AccountType)
	*p = x
	return p
}

func (x AccountType) String() string {
	return jsonx.EnumName(AccountType_name, int(x))
}

func (x *AccountType) UnmarshalJSON(data []byte) error {
	value, err := jsonx.UnmarshalJSONEnum(AccountType_value, data, "AccountType")
	if err != nil {
		return err
	}
	*x = AccountType(value)
	return nil
}

type GetInvitationsRequest struct {
	Paging  *common.Paging   `json:"paging"`
	Filters []*common.Filter `json:"filters"`
}

func (m *GetInvitationsRequest) Reset()         { *m = GetInvitationsRequest{} }
func (m *GetInvitationsRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetInvitationByTokenRequest struct {
	Token string `json:"token"`
}

func (m *GetInvitationByTokenRequest) Reset()         { *m = GetInvitationByTokenRequest{} }
func (m *GetInvitationByTokenRequest) String() string { return jsonx.MustMarshalToString(m) }

type Invitation struct {
	Id         dot.ID         `json:"id"`
	ShopId     dot.ID         `json:"shop_id"`
	Email      string         `json:"email"`
	Roles      []string       `json:"roles"`
	Token      string         `json:"token"`
	Status     status3.Status `json:"status"`
	InvitedBy  dot.ID         `json:"invited_by"`
	AcceptedAt dot.Time       `json:"accepted_at"`
	DeclinedAt dot.Time       `json:"declined_at"`
	ExpiredAt  dot.Time       `json:"expired_at"`
	CreatedAt  dot.Time       `json:"created_at"`
	UpdatedAt  dot.Time       `json:"updated_at"`
}

func (m *Invitation) Reset()         { *m = Invitation{} }
func (m *Invitation) String() string { return jsonx.MustMarshalToString(m) }

type AcceptInvitationRequest struct {
	Token string `json:"token"`
}

func (m *AcceptInvitationRequest) Reset()         { *m = AcceptInvitationRequest{} }
func (m *AcceptInvitationRequest) String() string { return jsonx.MustMarshalToString(m) }

type RejectInvitationRequest struct {
	Token string `json:"token"`
}

func (m *RejectInvitationRequest) Reset()         { *m = RejectInvitationRequest{} }
func (m *RejectInvitationRequest) String() string { return jsonx.MustMarshalToString(m) }

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
}

func (m *User) Reset()         { *m = User{} }
func (m *User) String() string { return jsonx.MustMarshalToString(m) }

type IDsRequest struct {
	// @required
	Ids   []dot.ID      `json:"ids"`
	Mixed *MixedAccount `json:"mixed"`
}

func (m *IDsRequest) Reset()         { *m = IDsRequest{} }
func (m *IDsRequest) String() string { return jsonx.MustMarshalToString(m) }

type MixedAccount struct {
	Ids      []dot.ID `json:"ids"`
	All      bool     `json:"all"`
	AllShops bool     `json:"all_shops"`
}

func (m *MixedAccount) Reset()         { *m = MixedAccount{} }
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

func (m *Partner) Reset()         { *m = Partner{} }
func (m *Partner) String() string { return jsonx.MustMarshalToString(m) }

type PublicAccountInfo struct {
	Id       dot.ID      `json:"id"`
	Name     string      `json:"name"`
	Type     AccountType `json:"type"`
	ImageUrl string      `json:"image_url"`
	Website  string      `json:"website"`
}

func (m *PublicAccountInfo) Reset()         { *m = PublicAccountInfo{} }
func (m *PublicAccountInfo) String() string { return jsonx.MustMarshalToString(m) }

type PublicAuthorizedPartnerInfo struct {
	Id          dot.ID      `json:"id"`
	Name        string      `json:"name"`
	Type        AccountType `json:"type"`
	ImageUrl    string      `json:"image_url"`
	Website     string      `json:"website"`
	RedirectUrl string      `json:"redirect_url"`
}

func (m *PublicAuthorizedPartnerInfo) Reset()         { *m = PublicAuthorizedPartnerInfo{} }
func (m *PublicAuthorizedPartnerInfo) String() string { return jsonx.MustMarshalToString(m) }

type PublicAccountsResponse struct {
	Accounts []*PublicAccountInfo `json:"accounts"`
}

func (m *PublicAccountsResponse) Reset()         { *m = PublicAccountsResponse{} }
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
	// @deprecated use try_on instead
	GhnNoteCode string           `json:"ghn_note_code"`
	TryOn       try_on.TryOnCode `json:"try_on"`
	OwnerId     dot.ID           `json:"owner_id"`
	User        *User            `json:"user"`
	CompanyInfo *CompanyInfo     `json:"company_info"`
	// referrence: https://icalendar.org/rrule-tool.html
	MoneyTransactionRrule         string                               `json:"money_transaction_rrule"`
	SurveyInfo                    []*SurveyInfo                        `json:"survey_info"`
	ShippingServiceSelectStrategy []*ShippingServiceSelectStrategyItem `json:"shipping_service_select_strategy"`
	Code                          string                               `json:"code"`
}

func (m *Shop) Reset()         { *m = Shop{} }
func (m *Shop) String() string { return jsonx.MustMarshalToString(m) }

type ShippingServiceSelectStrategyItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (m *ShippingServiceSelectStrategyItem) Reset()         { *m = ShippingServiceSelectStrategyItem{} }
func (m *ShippingServiceSelectStrategyItem) String() string { return jsonx.MustMarshalToString(m) }

type SurveyInfo struct {
	Key      string `json:"key"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

func (m *SurveyInfo) Reset()         { *m = SurveyInfo{} }
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
}

func (m *CreateUserRequest) Reset()         { *m = CreateUserRequest{} }
func (m *CreateUserRequest) String() string { return jsonx.MustMarshalToString(m) }

type RegisterResponse struct {
	// @required
	User *User `json:"user"`
}

func (m *RegisterResponse) Reset()         { *m = RegisterResponse{} }
func (m *RegisterResponse) String() string { return jsonx.MustMarshalToString(m) }

// Exactly one of phone or email must be provided.
type ResetPasswordRequest struct {
	// Phone number to send reset password instruction.
	Phone string `json:"phone"`
	// Email address to send reset password instruction.
	Email string `json:"email"`
}

func (m *ResetPasswordRequest) Reset()         { *m = ResetPasswordRequest{} }
func (m *ResetPasswordRequest) String() string { return jsonx.MustMarshalToString(m) }

// Exactly one of current_password or reset_password_token must be provided.
type ChangePasswordRequest struct {
	// @required
	CurrentPassword string `json:"current_password"`
	// @required
	NewPassword string `json:"new_password"`
	// @required
	ConfirmPassword string `json:"confirm_password"`
}

func (m *ChangePasswordRequest) Reset()         { *m = ChangePasswordRequest{} }
func (m *ChangePasswordRequest) String() string { return jsonx.MustMarshalToString(m) }

// Exactly one of email or phone must be provided.
type ChangePasswordUsingTokenRequest struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
	// @required
	ResetPasswordToken string `json:"reset_password_token"`
	// @required
	NewPassword string `json:"new_password"`
	// @required
	ConfirmPassword string `json:"confirm_password"`
}

func (m *ChangePasswordUsingTokenRequest) Reset()         { *m = ChangePasswordUsingTokenRequest{} }
func (m *ChangePasswordUsingTokenRequest) String() string { return jsonx.MustMarshalToString(m) }

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

func (m *Permission) Reset()         { *m = Permission{} }
func (m *Permission) String() string { return jsonx.MustMarshalToString(m) }

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
	AccountType AccountType `json:"account_type"`
	// Not implemented.
	AccountKey string `json:"account_key"`
}

func (m *LoginRequest) Reset()         { *m = LoginRequest{} }
func (m *LoginRequest) String() string { return jsonx.MustMarshalToString(m) }

// Represents an account associated with the current user. It has extra fields to represents relation with the user.
type LoginAccount struct {
	ExportedFields []string `json:"exported_fields"`
	// @required
	Id dot.ID `json:"id"`
	// @required
	Name string `json:"name"`
	// @required
	Type AccountType `json:"type"`
	// Associated token for the account. It's returned when calling Login or
	// SwitchAccount with regenerate_tokens set to true.
	AccessToken string `json:"access_token"`
	// The same as access_token.
	ExpiresIn   int              `json:"expires_in"`
	ImageUrl    string           `json:"image_url"`
	UrlSlug     string           `json:"url_slug"`
	UserAccount *UserAccountInfo `json:"user_account"`
}

func (m *LoginAccount) Reset()         { *m = LoginAccount{} }
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

func (m *LoginResponse) Reset()         { *m = LoginResponse{} }
func (m *LoginResponse) String() string { return jsonx.MustMarshalToString(m) }

type SwitchAccountRequest struct {
	// @required
	AccountId dot.ID `json:"account_id"`
	// This field should only be used after creating new accounts. If it is set,
	// account_id can be left empty.
	RegenerateTokens bool `json:"regenerate_tokens"`
}

func (m *SwitchAccountRequest) Reset()         { *m = SwitchAccountRequest{} }
func (m *SwitchAccountRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpgradeAccessTokenRequest struct {
	// @required
	Stoken   string `json:"stoken"`
	Password string `json:"password"`
}

func (m *UpgradeAccessTokenRequest) Reset()         { *m = UpgradeAccessTokenRequest{} }
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

func (m *AccessTokenResponse) Reset()         { *m = AccessTokenResponse{} }
func (m *AccessTokenResponse) String() string { return jsonx.MustMarshalToString(m) }

type SendSTokenEmailRequest struct {
	// @required
	Email string `json:"email"`
	// @required
	AccountId dot.ID `json:"account_id"`
}

func (m *SendSTokenEmailRequest) Reset()         { *m = SendSTokenEmailRequest{} }
func (m *SendSTokenEmailRequest) String() string { return jsonx.MustMarshalToString(m) }

type SendEmailVerificationRequest struct {
	// @required
	Email string `json:"email"`
}

func (m *SendEmailVerificationRequest) Reset()         { *m = SendEmailVerificationRequest{} }
func (m *SendEmailVerificationRequest) String() string { return jsonx.MustMarshalToString(m) }

type SendPhoneVerificationRequest struct {
	// @required
	Phone string `json:"phone"`
}

func (m *SendPhoneVerificationRequest) Reset()         { *m = SendPhoneVerificationRequest{} }
func (m *SendPhoneVerificationRequest) String() string { return jsonx.MustMarshalToString(m) }

type VerifyEmailUsingTokenRequest struct {
	// @required
	VerificationToken string `json:"verification_token"`
}

func (m *VerifyEmailUsingTokenRequest) Reset()         { *m = VerifyEmailUsingTokenRequest{} }
func (m *VerifyEmailUsingTokenRequest) String() string { return jsonx.MustMarshalToString(m) }

type VerifyPhoneUsingTokenRequest struct {
	// @required
	VerificationToken string `json:"verification_token"`
}

func (m *VerifyPhoneUsingTokenRequest) Reset()         { *m = VerifyPhoneUsingTokenRequest{} }
func (m *VerifyPhoneUsingTokenRequest) String() string { return jsonx.MustMarshalToString(m) }

type InviteUserToAccountRequest struct {
	// @required account to manage, must be owner of the account
	AccountId dot.ID `json:"account_id"`
	// @required phone or email
	InviteeIdentifier string      `json:"invitee_identifier"`
	Permission        *Permission `json:"permission"`
}

func (m *InviteUserToAccountRequest) Reset()         { *m = InviteUserToAccountRequest{} }
func (m *InviteUserToAccountRequest) String() string { return jsonx.MustMarshalToString(m) }

type AnswerInvitationRequest struct {
	AccountId dot.ID          `json:"account_id"`
	Response  *status3.Status `json:"response"`
}

func (m *AnswerInvitationRequest) Reset()         { *m = AnswerInvitationRequest{} }
func (m *AnswerInvitationRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetUsersInCurrentAccountsRequest struct {
	Paging  *common.Paging   `json:"paging"`
	Filters []*common.Filter `json:"filters"`
	Mixed   *MixedAccount    `json:"mixed"`
}

func (m *GetUsersInCurrentAccountsRequest) Reset()         { *m = GetUsersInCurrentAccountsRequest{} }
func (m *GetUsersInCurrentAccountsRequest) String() string { return jsonx.MustMarshalToString(m) }

type PublicUserInfo struct {
	// @required
	Id dot.ID `json:"id"`
	// @required
	FullName string `json:"full_name"`
	// @required
	ShortName string `json:"short_name"`
}

func (m *PublicUserInfo) Reset()         { *m = PublicUserInfo{} }
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
	AccountType AccountType `json:"account_type"`
	Position    string      `json:"position"`
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

func (m *UserAccountInfo) Reset()         { *m = UserAccountInfo{} }
func (m *UserAccountInfo) String() string { return jsonx.MustMarshalToString(m) }

// Prepsents users in current account
type ProtectedUsersResponse struct {
	Paging *common.PageInfo   `json:"paging"`
	Users  []*UserAccountInfo `json:"users"`
}

func (m *ProtectedUsersResponse) Reset()         { *m = ProtectedUsersResponse{} }
func (m *ProtectedUsersResponse) String() string { return jsonx.MustMarshalToString(m) }

type LeaveAccountRequest struct {
}

func (m *LeaveAccountRequest) Reset()         { *m = LeaveAccountRequest{} }
func (m *LeaveAccountRequest) String() string { return jsonx.MustMarshalToString(m) }

type RemoveUserFromCurrentAccountRequest struct {
}

func (m *RemoveUserFromCurrentAccountRequest) Reset()         { *m = RemoveUserFromCurrentAccountRequest{} }
func (m *RemoveUserFromCurrentAccountRequest) String() string { return jsonx.MustMarshalToString(m) }

type Province struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Region string `json:"region"`
}

func (m *Province) Reset()         { *m = Province{} }
func (m *Province) String() string { return jsonx.MustMarshalToString(m) }

type GetProvincesResponse struct {
	Provinces []*Province `json:"provinces"`
}

func (m *GetProvincesResponse) Reset()         { *m = GetProvincesResponse{} }
func (m *GetProvincesResponse) String() string { return jsonx.MustMarshalToString(m) }

type District struct {
	Code         string `json:"code"`
	ProvinceCode string `json:"province_code"`
	Name         string `json:"name"`
	IsFreeship   bool   `json:"is_freeship"`
}

func (m *District) Reset()         { *m = District{} }
func (m *District) String() string { return jsonx.MustMarshalToString(m) }

type GetDistrictsResponse struct {
	Districts []*District `json:"districts"`
}

func (m *GetDistrictsResponse) Reset()         { *m = GetDistrictsResponse{} }
func (m *GetDistrictsResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetDistrictsByProvinceRequest struct {
	ProvinceCode string `json:"province_code"`
	ProvinceName string `json:"province_name"`
}

func (m *GetDistrictsByProvinceRequest) Reset()         { *m = GetDistrictsByProvinceRequest{} }
func (m *GetDistrictsByProvinceRequest) String() string { return jsonx.MustMarshalToString(m) }

type Ward struct {
	Code         string `json:"code"`
	DistrictCode string `json:"district_code"`
	Name         string `json:"name"`
}

func (m *Ward) Reset()         { *m = Ward{} }
func (m *Ward) String() string { return jsonx.MustMarshalToString(m) }

type GetWardsResponse struct {
	Wards []*Ward `json:"wards"`
}

func (m *GetWardsResponse) Reset()         { *m = GetWardsResponse{} }
func (m *GetWardsResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetWardsByDistrictRequest struct {
	DistrictCode string `json:"district_code"`
	DistrictName string `json:"district_name"`
}

func (m *GetWardsByDistrictRequest) Reset()         { *m = GetWardsByDistrictRequest{} }
func (m *GetWardsByDistrictRequest) String() string { return jsonx.MustMarshalToString(m) }

type ParseLocationRequest struct {
	ProvinceName string `json:"province_name"`
	DistrictName string `json:"district_name"`
	WardName     string `json:"ward_name"`
}

func (m *ParseLocationRequest) Reset()         { *m = ParseLocationRequest{} }
func (m *ParseLocationRequest) String() string { return jsonx.MustMarshalToString(m) }

type ParseLocationResponse struct {
	Province *Province `json:"province"`
	District *District `json:"district"`
	Ward     *Ward     `json:"ward"`
}

func (m *ParseLocationResponse) Reset()         { *m = ParseLocationResponse{} }
func (m *ParseLocationResponse) String() string { return jsonx.MustMarshalToString(m) }

type Bank struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func (m *Bank) Reset()         { *m = Bank{} }
func (m *Bank) String() string { return jsonx.MustMarshalToString(m) }

type GetBanksResponse struct {
	Banks []*Bank `json:"banks"`
}

func (m *GetBanksResponse) Reset()         { *m = GetBanksResponse{} }
func (m *GetBanksResponse) String() string { return jsonx.MustMarshalToString(m) }

type BankProvince struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	BankCode string `json:"bank_code"`
}

func (m *BankProvince) Reset()         { *m = BankProvince{} }
func (m *BankProvince) String() string { return jsonx.MustMarshalToString(m) }

type GetBankProvincesResponse struct {
	Provinces []*BankProvince `json:"provinces"`
}

func (m *GetBankProvincesResponse) Reset()         { *m = GetBankProvincesResponse{} }
func (m *GetBankProvincesResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetProvincesByBankResquest struct {
	BankCode string `json:"bank_code"`
	BankName string `json:"bank_name"`
}

func (m *GetProvincesByBankResquest) Reset()         { *m = GetProvincesByBankResquest{} }
func (m *GetProvincesByBankResquest) String() string { return jsonx.MustMarshalToString(m) }

type BankBranch struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	BankCode     string `json:"bank_code"`
	ProvinceCode string `json:"province_code"`
}

func (m *BankBranch) Reset()         { *m = BankBranch{} }
func (m *BankBranch) String() string { return jsonx.MustMarshalToString(m) }

type GetBranchesByBankProvinceResponse struct {
	Branches []*BankBranch `json:"branches"`
}

func (m *GetBranchesByBankProvinceResponse) Reset()         { *m = GetBranchesByBankProvinceResponse{} }
func (m *GetBranchesByBankProvinceResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetBranchesByBankProvinceResquest struct {
	BankCode     string `json:"bank_code"`
	BankName     string `json:"bank_name"`
	ProvinceCode string `json:"province_code"`
	ProvinceName string `json:"province_name"`
}

func (m *GetBranchesByBankProvinceResquest) Reset()         { *m = GetBranchesByBankProvinceResquest{} }
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

func (m *Address) Reset()         { *m = Address{} }
func (m *Address) String() string { return jsonx.MustMarshalToString(m) }

type Coordinates struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

func (m *Coordinates) Reset()         { *m = Coordinates{} }
func (m *Coordinates) String() string { return jsonx.MustMarshalToString(m) }

type BankAccount struct {
	Name          string `json:"name"`
	Province      string `json:"province"`
	Branch        string `json:"branch"`
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
}

func (m *BankAccount) Reset()         { *m = BankAccount{} }
func (m *BankAccount) String() string { return jsonx.MustMarshalToString(m) }

type ContactPerson struct {
	Name     string `json:"name"`
	Position string `json:"position"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

func (m *ContactPerson) Reset()         { *m = ContactPerson{} }
func (m *ContactPerson) String() string { return jsonx.MustMarshalToString(m) }

type CompanyInfo struct {
	Name                string         `json:"name"`
	TaxCode             string         `json:"tax_code"`
	Address             string         `json:"address"`
	Website             string         `json:"website"`
	LegalRepresentative *ContactPerson `json:"legal_representative"`
}

func (m *CompanyInfo) Reset()         { *m = CompanyInfo{} }
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

func (m *CreateAddressRequest) Reset()         { *m = CreateAddressRequest{} }
func (m *CreateAddressRequest) String() string { return jsonx.MustMarshalToString(m) }

type AddressNote struct {
	Note       string `json:"note"`
	OpenTime   string `json:"open_time"`
	LunchBreak string `json:"lunch_break"`
	Other      string `json:"other"`
}

func (m *AddressNote) Reset()         { *m = AddressNote{} }
func (m *AddressNote) String() string { return jsonx.MustMarshalToString(m) }

type GetAddressResponse struct {
	Addresses []*Address `json:"addresses"`
}

func (m *GetAddressResponse) Reset()         { *m = GetAddressResponse{} }
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

func (m *UpdateAddressRequest) Reset()         { *m = UpdateAddressRequest{} }
func (m *UpdateAddressRequest) String() string { return jsonx.MustMarshalToString(m) }

type SetDefaultAddressRequest struct {
	Id   dot.ID                   `json:"id"`
	Type address_type.AddressType `json:"type"`
}

func (m *SetDefaultAddressRequest) Reset()         { *m = SetDefaultAddressRequest{} }
func (m *SetDefaultAddressRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateURLSlugRequest struct {
	AccountId dot.ID         `json:"account_id"`
	UrlSlug   dot.NullString `json:"url_slug"`
}

func (m *UpdateURLSlugRequest) Reset()         { *m = UpdateURLSlugRequest{} }
func (m *UpdateURLSlugRequest) String() string { return jsonx.MustMarshalToString(m) }

type HistoryResponse struct {
	Paging *common.PageInfo      `json:"paging"`
	Data   *common.RawJSONObject `json:"data"`
}

func (m *HistoryResponse) Reset()         { *m = HistoryResponse{} }
func (m *HistoryResponse) String() string { return jsonx.MustMarshalToString(m) }

type Credit struct {
	Id        dot.ID         `json:"id"`
	Amount    int            `json:"amount"`
	ShopId    dot.ID         `json:"shop_id"`
	Type      string         `json:"type"`
	Shop      *Shop          `json:"shop"`
	CreatedAt dot.Time       `json:"created_at"`
	UpdatedAt dot.Time       `json:"updated_at"`
	PaidAt    dot.Time       `json:"paid_at"`
	Status    status3.Status `json:"status"`
}

func (m *Credit) Reset()         { *m = Credit{} }
func (m *Credit) String() string { return jsonx.MustMarshalToString(m) }

type CreditsResponse struct {
	Paging  *common.PageInfo `json:"paging"`
	Credits []*Credit        `json:"credits"`
}

func (m *CreditsResponse) Reset()         { *m = CreditsResponse{} }
func (m *CreditsResponse) String() string { return jsonx.MustMarshalToString(m) }

type UpdatePermissionRequest struct {
	Items []*UpdatePermissionItem `json:"items"`
}

func (m *UpdatePermissionRequest) Reset()         { *m = UpdatePermissionRequest{} }
func (m *UpdatePermissionRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdatePermissionItem struct {
	Type       string   `json:"type"`
	Key        string   `json:"key"`
	Grants     []string `json:"grants"`
	Revokes    []string `json:"revokes"`
	ReplaceAll []string `json:"replace_all"`
	RevokeAll  bool     `json:"revoke_all"`
}

func (m *UpdatePermissionItem) Reset()         { *m = UpdatePermissionItem{} }
func (m *UpdatePermissionItem) String() string { return jsonx.MustMarshalToString(m) }

type UpdatePermissionResponse struct {
	Msg string `json:"msg"`
}

func (m *UpdatePermissionResponse) Reset()         { *m = UpdatePermissionResponse{} }
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

func (m *Device) Reset()         { *m = Device{} }
func (m *Device) String() string { return jsonx.MustMarshalToString(m) }

type Notification struct {
	Id               dot.ID         `json:"id"`
	Title            string         `json:"title"`
	Message          string         `json:"message"`
	IsRead           bool           `json:"is_read"`
	Entity           string         `json:"entity"`
	EntityId         dot.ID         `json:"entity_id"`
	AccountId        dot.ID         `json:"account_id"`
	SendNotification bool           `json:"send_notification"`
	SeenAt           dot.Time       `json:"seen_at"`
	CreatedAt        dot.Time       `json:"created_at"`
	UpdatedAt        dot.Time       `json:"updated_at"`
	SyncStatus       status3.Status `json:"sync_status"`
}

func (m *Notification) Reset()         { *m = Notification{} }
func (m *Notification) String() string { return jsonx.MustMarshalToString(m) }

type CreateDeviceRequest struct {
	DeviceId         string `json:"device_id"`
	DeviceName       string `json:"device_name"`
	ExternalDeviceId string `json:"external_device_id"`
}

func (m *CreateDeviceRequest) Reset()         { *m = CreateDeviceRequest{} }
func (m *CreateDeviceRequest) String() string { return jsonx.MustMarshalToString(m) }

type DeleteDeviceRequest struct {
	DeviceId         string `json:"device_id"`
	ExternalDeviceId string `json:"external_device_id"`
}

func (m *DeleteDeviceRequest) Reset()         { *m = DeleteDeviceRequest{} }
func (m *DeleteDeviceRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetNotificationsRequest struct {
	Paging *common.Paging `json:"paging"`
}

func (m *GetNotificationsRequest) Reset()         { *m = GetNotificationsRequest{} }
func (m *GetNotificationsRequest) String() string { return jsonx.MustMarshalToString(m) }

type NotificationsResponse struct {
	Paging        *common.PageInfo `json:"paging"`
	Notifications []*Notification  `json:"notifications"`
}

func (m *NotificationsResponse) Reset()         { *m = NotificationsResponse{} }
func (m *NotificationsResponse) String() string { return jsonx.MustMarshalToString(m) }

type UpdateNotificationsRequest struct {
	Ids    []dot.ID `json:"ids"`
	IsRead bool     `json:"is_read"`
}

func (m *UpdateNotificationsRequest) Reset()         { *m = UpdateNotificationsRequest{} }
func (m *UpdateNotificationsRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateReferenceUserRequest struct {
	// @required
	Phone string `json:"phone"`
}

func (m *UpdateReferenceUserRequest) Reset()         { *m = UpdateReferenceUserRequest{} }
func (m *UpdateReferenceUserRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateReferenceSaleRequest struct {
	// @required
	Phone string `json:"phone"`
}

func (m *UpdateReferenceSaleRequest) Reset()         { *m = UpdateReferenceSaleRequest{} }
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

func (m *Affiliate) Reset()         { *m = Affiliate{} }
func (m *Affiliate) String() string { return jsonx.MustMarshalToString(m) }

type GetUserByPhoneRequest struct {
	// @required
	Phone          string `json:"phone"`
	RecaptchaToken string `json:"recaptcha_token"`
}

func (m *GetUserByPhoneRequest) Reset()         { *m = GetUserByPhoneRequest{} }
func (m *GetUserByPhoneRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetUserByPhoneResponse struct {
	// @required
	Exists bool `json:"exists"`
}

func (m *GetUserByPhoneResponse) Reset()         { *m = GetUserByPhoneResponse{} }
func (m *GetUserByPhoneResponse) String() string { return jsonx.MustMarshalToString(m) }

type SendPhoneVerifyRequest struct {
	// @required
	Phone          string `json:"phone"`
	RecaptchaToken string `json:"recaptcha_token"`
}

func (m *SendPhoneVerifyRequest) Reset()         { *m = SendPhoneVerifyRequest{} }
func (m *SendPhoneVerifyRequest) String() string { return jsonx.MustMarshalToString(m) }
