package fbclient

type Pagination struct {
	Limit  string `url:"limit,omitempty"`
	Before string `url:"before,omitempty"`
	After  string `url:"after,omitempty"`
	Offset string `url:"offset,omitempty"`
	Until  string `url:"until,omitempty"`
	Since  string `url:"since,omitempty"`
}

type GetMeParams struct {
	AccessToken string `url:"access_token"`
	Fields      string `url:"fields"`
}

type PingParams struct {
	ClientID     string `url:"client_id"`
	ClientSecret string `url:"client_secret"`
	GrantType    string `url:"grant_type"`
}

type GetAccountsParams struct {
	AccessToken string `url:"access_token"`
	Fields      string `url:"fields"`
	DateFormat  string `url:"date_format"`
}

type GetLongLivedAccessTokenParams struct {
	GrantType       string `url:"grant_type"`
	FBExchangeToken string `url:"fb_exchange_token"`
	ClientID        string `url:"client_id"`
	ClientSecret    string `url:"client_secret"`
}

type DebugTokenParams struct {
	AccessToken string `url:"access_token"`
	InputToken  string `url:"input_token"`
}

type ListFeedsParams struct {
	AccessToken string `url:"access_token"`
	Fields      string `url:"fields"`
	DateFormat  string `url:"date_format"`

	Pagination
}

type ListPublishedPostsParams struct {
	AccessToken string `url:"access_token"`
	Fields      string `url:"fields"`
	DateFormat  string `url:"date_format"`

	Pagination
}

type GetPostParams struct {
	AccessToken string `url:"access_token"`
	Fields      string `url:"fields"`
	DateFormat  string `url:"date_format"`
}

type ListCommentsParams struct {
	AccessToken string `url:"access_token"`
	Fields      string `url:"fields"`
	DateFormat  string `url:"date_format"`

	Pagination
}

type ListCommentByPostIDsParams struct {
	AccessToken string `url:"access_token"`
	IDs         string `url:"ids"`
	Filter      string `url:"filter"`
	Limit       string `url:"limit"`
	Fields      string `url:"fields"`
	DateFormat  string `url:"date_format"`
}

type ListConversationsParams struct {
	AccessToken string `url:"access_token"`
	Fields      string `url:"fields"`
	DateFormat  string `url:"date_format"`

	Pagination
}

type GetConversationByUserIDParams struct {
	AccessToken string `url:"access_token"`
	Fields      string `url:"fields"`
	DateFormat  string `url:"date_format"`
	UserID      string `url:"user_id"`
}

type ListMessagesParams struct {
	AccessToken string `url:"access_token"`
	Fields      string `url:"fields"`
	DateFormat  string `url:"date_format"`

	Pagination
}

type GetMessageParams struct {
	AccessToken string `url:"access_token"`
	Fields      string `url:"fields"`
	DateFormat  string `url:"date_format"`
}

type GetCommentByIDParams struct {
	AccessToken string `url:"access_token"`
	Fields      string `url:"fields"`
	DateFormat  string `url:"date_format"`
}

type CreateSubscribedAppsParams struct {
	AccessToken      string `url:"access_token"`
	SubscribedFields string `url:"subscribed_fields"`
}

type SendMessageParams struct {
	AccessToken string `url:"access_token"`
	Recipient   string `url:"recipient"`
	Message     string `url:"message"`
}

type SendCommentParams struct {
	AccessToken   string `url:"access_token"`
	Message       string `url:"message,omitempty"`
	AttachmentURL string `url:"attachment_url,omitempty"`
}

type CreatePostParams struct {
	AccessToken string `url:"access_token"`
	Message     string `url:"message"`
}

type GetProfileByPISDParams struct {
	AccessToken string `url:"access_token"`
	Fields      string `url:"fields"`
}

// key is businessID
// docs: https://developers.facebook.com/docs/graph-api/overview/rate-limiting/
type XBusinessUseCaseUsageHeader map[string]*XBusinessUseCaseUsageHeaderContents

type XBusinessUseCaseUsageHeaderContents []*XBusinessUseCaseUsageHeaderContent

type XBusinessUseCaseUsageHeaderContent struct {
	Type                        string `json:"type"`
	CallCount                   int    `json:"call_count"`
	TotalCPUTime                int    `json:"total_cputime"`
	TotalTime                   int    `json:"total_time"`
	EstimatedTimeToRegainAccess int    `json:"estimated_time_to_regain_access"` // in minutes
}

func (x XBusinessUseCaseUsageHeader) IsRateLimit(businessID string) bool {
	xBusinessUseCaseUsageHeaderContents, ok := x[businessID]
	if !ok {
		return false
	}

	if xBusinessUseCaseUsageHeaderContents == nil {
		return false
	}

	for _, xBusinessUseCaseUsageHeaderContent := range *xBusinessUseCaseUsageHeaderContents {
		if xBusinessUseCaseUsageHeaderContent.CallCount >= 100 {
			return true
		}
	}

	return false
}

// in minutes
func (x XBusinessUseCaseUsageHeader) GetEstimatedTimeToRegainAccessOfPage(businessID string) int {
	xBusinessUseCaseUsageHeaderContents, ok := x[businessID]
	if !ok {
		return 0
	}

	if xBusinessUseCaseUsageHeaderContents == nil {
		return 0
	}

	for _, xBusinessUseCaseUsageHeaderContent := range *xBusinessUseCaseUsageHeaderContents {
		if xBusinessUseCaseUsageHeaderContent.Type == "pages" {
			return xBusinessUseCaseUsageHeaderContent.EstimatedTimeToRegainAccess
		}
	}
	return 0
}

func (x XBusinessUseCaseUsageHeader) GetEstimatedTimeToRegainAccessOfMessenger(businessID string) int {
	xBusinessUseCaseUsageHeaderContents, ok := x[businessID]
	if !ok {
		return 0
	}

	if xBusinessUseCaseUsageHeaderContents == nil {
		return 0
	}

	for _, xBusinessUseCaseUsageHeaderContent := range *xBusinessUseCaseUsageHeaderContents {
		if xBusinessUseCaseUsageHeaderContent.Type == "messenger" {
			return xBusinessUseCaseUsageHeaderContent.EstimatedTimeToRegainAccess
		}
	}
	return 0
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
