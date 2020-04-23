package fbclient

const (
	AccessTokenKey           = "access_token"
	ClientIDKey              = "client_id"
	ClientSecretKey          = "client_secret"
	GrantType                = "grant_type"
	Fields                   = "fields"
	InputToken               = "input_token"
	GrantTypeFBExchangeToken = "fb_exchange_token"
	FBExchangeToken          = "fb_exchange_token"
	ClientCredentials        = "client_credentials"
	ExpiresInUserToken       = 5184000 // 60 days
)

type FacebookRole int

const (
	UNKNOWN    FacebookRole = 0
	ADMIN      FacebookRole = 1
	ADVERTISER FacebookRole = 2
	ANALYST    FacebookRole = 3
	EDITOR     FacebookRole = 4
	MODERATOR  FacebookRole = 5
)
