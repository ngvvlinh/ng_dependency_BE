package client

import cc "etop.vn/backend/pkg/common/config"

type Connection struct {
	Subdomain string
	TokenStr  string
}

type Config struct {
	APIKey string `yaml:"api_key"`
	Secret string `yaml:"secret"`
}

func (c *Config) MustLoadEnv(prefix ...string) {
	p := "ET_HARAVAN"
	if len(prefix) > 0 {
		p = prefix[0]
	}
	cc.EnvMap{
		p + "_API_KEY": &c.APIKey,
		p + "_SECRET":  &c.Secret,
	}.MustLoad()
}

type GetShopRequest struct {
	Connection `json:"-"`
}

type GetShopResponse struct {
	Shop *Shop `json:"shop"`
}

type Shop struct {
	Address1                        string  `json:"address1"`
	City                            string  `json:"city"`
	Country                         string  `json:"country"`
	CreatedAt                       string  `json:"created_at"`
	CustomerEmail                   string  `json:"customer_email"`
	Domain                          string  `json:"domain"`
	Email                           string  `json:"email"`
	Id                              int     `json:"id"`
	Latitude                        float32 `json:"latitude"`
	Longitude                       float32 `json:"longitude"`
	Name                            string  `json:"name"`
	Phone                           string  `json:"phone"`
	PrimaryLocale                   string  `json:"primary_locale"`
	PrimaryLocationId               string  `json:"primary_location_id"`
	Province                        string  `json:"province"`
	Source                          string  `json:"source"`
	Zip                             string  `json:"zip"`
	CountryCode                     string  `json:"country_code"`
	CountryName                     string  `json:"country_name"`
	Currency                        string  `json:"currency"`
	Timezone                        string  `json:"timezone"`
	IanaTimezone                    string  `json:"iana_timezone"`
	ShopOwner                       string  `json:"shop_owner"`
	MoneyFormat                     string  `json:"money_format"`
	MoneyWithCurrencyFormat         string  `json:"money_with_currency_format"`
	ProvinceCode                    string  `json:"province_code"`
	TaxesIncluded                   string  `json:"taxes_included"`
	TaxShipping                     bool    `json:"tax_shipping"`
	CountyTaxes                     string  `json:"county_taxes"`
	PlanDisplayName                 string  `json:"plan_display_name"`
	PlanName                        string  `json:"plan_name"`
	MyharavanDomain                 string  `json:"myharavan_domain"`
	GoogleAppsDomain                string  `json:"google_apps_domain"`
	GoogleAppsLoginEnabled          string  `json:"google_apps_login_enabled"`
	MoneyInEmailsFormat             string  `json:"money_in_emails_format"`
	MoneyWithCurrencyInEmailsFormat string  `json:" money_with_currency_in_emails_format"`
	EligibleForPayments             string  `json:"eligible_for_payments"`
	RequiresExtraPaymentsAgreement  string  `json:" requires_extra_payments_agreement"`
	PasswordEnabled                 bool    `json:"password_enabled"`
	HasStorefront                   bool    `json:"has_storefront"`
}

type ConnectCarrierServiceRequest struct {
	Connection     `json:"-"`
	CarrierService *CarrierService `json:"carrier_service"`
}

type ConnectCarrierServiceResponse struct {
	CarrierService *CarrierService `json:"carrier_service"`
}

type CarrierService struct {
	ID                  int    `json:"id"`
	Active              bool   `json:"active"`                 // : true
	TrackingUrl         string `json:"tracking_url"`           // : Link xem chi tiết vận đơn
	CreateOrderUrl      string `json:"create_order_url"`       // : Link xem chi tiết vận đơn
	GetOrderDetailUrl   string `json:"get_order_detail_url"`   // : Link xem chi tiết vận đơn
	GetShippingRatesUrl string `json:"get_shipping_rates_url"` // : Link xem chi tiết vận đơn
	CancelOrderUrl      string `json:"cancel_order_url"`       // : Link xem chi tiết vận đơn
	Name                string `json:"name"`                   // : "demo Express"
	CarrierServiceType  string `json:"carrier_service_type"`   // : "api"
	ServiceDiscovery    bool   `json:"service_discovery"`      // : false
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
}

type GetAccessTokenRequest struct {
	Subdomain   string `json:"-"`
	Code        string `json:"code"`
	RedirectURI string `json:"redirect_uri"`
}

type GetAccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
}

type DeleteConnectedCarrierServiceRequest struct {
	Connection       `json:"-"`
	CarrierServiceID int
}
