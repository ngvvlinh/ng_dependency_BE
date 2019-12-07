package integration

import (
	etop "etop.vn/api/top/int/etop"
	"etop.vn/api/top/types/etc/account_type"
	status3 "etop.vn/api/top/types/etc/status3"
	"etop.vn/capi/dot"
	"etop.vn/common/jsonx"
)

type InitRequest struct {
	AuthToken string `json:"auth_token"`
}

func (m *InitRequest) Reset()         { *m = InitRequest{} }
func (m *InitRequest) String() string { return jsonx.MustMarshalToString(m) }

type LoginResponse struct {
	AccessToken       string                     `json:"access_token"`
	ExpiresIn         int                        `json:"expires_in"`
	User              *PartnerUserLogin          `json:"user"`
	Account           *PartnerShopLoginAccount   `json:"account"`
	Shop              *PartnerShopInfo           `json:"shop"`
	AvailableAccounts []*PartnerShopLoginAccount `json:"available_accounts"`
	AuthPartner       *etop.PublicAccountInfo    `json:"auth_partner"`
	Actions           []*Action                  `json:"actions"`
	RedirectUrl       string                     `json:"redirect_url"`
}

func (m *LoginResponse) Reset()         { *m = LoginResponse{} }
func (m *LoginResponse) String() string { return jsonx.MustMarshalToString(m) }

type Action struct {
	Name  string            `json:"name"`
	Label string            `json:"label"`
	Msg   string            `json:"msg"`
	Meta  map[string]string `json:"meta" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *Action) Reset()         { *m = Action{} }
func (m *Action) String() string { return jsonx.MustMarshalToString(m) }

type PartnerUserLogin struct {
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
}

func (m *PartnerUserLogin) Reset()         { *m = PartnerUserLogin{} }
func (m *PartnerUserLogin) String() string { return jsonx.MustMarshalToString(m) }

type PartnerShopLoginAccount struct {
	// @required
	Id         dot.ID `json:"id"`
	ExternalId string `json:"external_id"`
	// @required
	Name string `json:"name"`
	// @required
	Type account_type.AccountType `json:"type"`
	// Associated token for the account. It's returned when calling Login or
	// SwitchAccount with regenerate_tokens set to true.
	AccessToken string `json:"access_token"`
	// The same as access_token.
	ExpiresIn int    `json:"expires_in"`
	ImageUrl  string `json:"image_url"`
}

func (m *PartnerShopLoginAccount) Reset()         { *m = PartnerShopLoginAccount{} }
func (m *PartnerShopLoginAccount) String() string { return jsonx.MustMarshalToString(m) }

type PartnerShopInfo struct {
	Id      dot.ID         `json:"id"`
	Name    string         `json:"name"`
	Status  status3.Status `json:"status"`
	IsTest  bool           `json:"is_test"`
	Address *etop.Address  `json:"address"`
	Phone   string         `json:"phone"`
	//    optional string website_url = 14 [(gogoproto.nullable)=false];
	ImageUrl string `json:"image_url"`
	Email    string `json:"email"`
	//    optional dot.ID product_source_id = 17 [(gogoproto.nullable)=false];
	//    optional dot.ID ship_to_address_id = 18 [(gogoproto.nullable)=false];
	ShipFromAddressId dot.ID `json:"ship_from_address_id"`
}

func (m *PartnerShopInfo) Reset()         { *m = PartnerShopInfo{} }
func (m *PartnerShopInfo) String() string { return jsonx.MustMarshalToString(m) }

type RequestLoginRequest struct {
	// @required Phone or email
	Login string `json:"login"`
	// @required
	RecaptchaToken string `json:"recaptcha_token"`
}

func (m *RequestLoginRequest) Reset()         { *m = RequestLoginRequest{} }
func (m *RequestLoginRequest) String() string { return jsonx.MustMarshalToString(m) }

type RequestLoginResponse struct {
	// @required
	Code string `json:"code"`
	// @required
	Msg string `json:"msg"`
	// @required
	Actions []*Action `json:"actions"`
}

func (m *RequestLoginResponse) Reset()         { *m = RequestLoginResponse{} }
func (m *RequestLoginResponse) String() string { return jsonx.MustMarshalToString(m) }

type LoginUsingTokenRequest struct {
	Login            string `json:"login"`
	VerificationCode string `json:"verification_code"`
}

func (m *LoginUsingTokenRequest) Reset()         { *m = LoginUsingTokenRequest{} }
func (m *LoginUsingTokenRequest) String() string { return jsonx.MustMarshalToString(m) }

type RegisterRequest struct {
	FullName       string       `json:"full_name"`
	Phone          string       `json:"phone"`
	Email          string       `json:"email"`
	AgreeTos       bool         `json:"agree_tos"`
	AgreeEmailInfo dot.NullBool `json:"agree_email_info"`
}

func (m *RegisterRequest) Reset()         { *m = RegisterRequest{} }
func (m *RegisterRequest) String() string { return jsonx.MustMarshalToString(m) }

type RegisterResponse struct {
	User        *etop.User `json:"user"`
	AccessToken string     `json:"access_token"`
	ExpiresIn   int        `json:"expires_in"`
}

func (m *RegisterResponse) Reset()         { *m = RegisterResponse{} }
func (m *RegisterResponse) String() string { return jsonx.MustMarshalToString(m) }

type GrantAccessRequest struct {
	ShopId dot.ID `json:"shop_id"`
}

func (m *GrantAccessRequest) Reset()         { *m = GrantAccessRequest{} }
func (m *GrantAccessRequest) String() string { return jsonx.MustMarshalToString(m) }

type GrantAccessResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func (m *GrantAccessResponse) Reset()         { *m = GrantAccessResponse{} }
func (m *GrantAccessResponse) String() string { return jsonx.MustMarshalToString(m) }
