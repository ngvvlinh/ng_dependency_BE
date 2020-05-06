package fbclient

const (
	AccessToken              = "access_token"
	IDs                      = "ids"
	ClientIDKey              = "client_id"
	ClientSecret             = "client_secret"
	GrantType                = "grant_type"
	Fields                   = "fields"
	Filter                   = "filter"
	Summary                  = "summary"
	Limit                    = "limit"
	InputToken               = "input_token"
	GrantTypeFBExchangeToken = "fb_exchange_token"
	FBExchangeToken          = "fb_exchange_token"
	ClientCredentials        = "client_credentials"
	ExpiresInUserToken       = 5184000 // 60 days

	DateFormat     = "date_format"
	UnixDateFormat = "U"

	PermissionGranted  = "granted"
	PermissionDeclined = "declined"

	DefaultLimitGetPosts         = 100
	DefaultLimitGetComments      = 100
	DefaultLimitGetConversations = 100
	DefaultLimitGetMessages      = 100
	MaximumIDs                   = 1
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

// TODO: Ngoc create enum
type PostAttachmentType string

const (
	Album PostAttachmentType = "album"
)
