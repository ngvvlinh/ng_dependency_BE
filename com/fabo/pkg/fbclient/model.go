package fbclient

import "o.o/backend/com/fabo/pkg/fbclient/model"

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

type ListFeedsRequest struct {
	AccessToken string
	PageID      string
	Pagination  *model.FacebookPagingRequest
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

type GetPostRequest struct {
	AccessToken string
	PostID      string
	PageID      string
}

type GetPostParams struct {
	AccessToken string `url:"access_token"`
	Fields      string `url:"fields"`
	DateFormat  string `url:"date_format"`
}

type ListCommentsRequest struct {
	AccessToken string
	PostID      string
	PageID      string
	Pagination  *model.FacebookPagingRequest
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

type ListConversationsRequest struct {
	AccessToken string
	PageID      string
	Pagination  *model.FacebookPagingRequest
}

type ListConversationsParams struct {
	AccessToken string `url:"access_token"`
	Fields      string `url:"fields"`
	DateFormat  string `url:"date_format"`

	Pagination
}

type GetConversationByUserIDRequest struct {
	AccessToken string
	PageID      string
	UserID      string
}

type GetConversationByUserIDParams struct {
	AccessToken string `url:"access_token"`
	Fields      string `url:"fields"`
	DateFormat  string `url:"date_format"`
	UserID      string `url:"user_id"`
}

type ListMessagesRequest struct {
	AccessToken    string
	ConversationID string
	PageID         string
	Pagination     *model.FacebookPagingRequest
}

type ListMessagesParams struct {
	AccessToken string `url:"access_token"`
	Fields      string `url:"fields"`
	DateFormat  string `url:"date_format"`

	Pagination
}

type GetMessageRequest struct {
	AccessToken string
	MessageID   string
	PageID      string
}

type GetMessageParams struct {
	AccessToken string `url:"access_token"`
	Fields      string `url:"fields"`
	DateFormat  string `url:"date_format"`
}

type GetCommentByIDRequest struct {
	AccessToken string
	CommentID   string
	PageID      string
}

type GetCommentByIDParams struct {
	AccessToken string `url:"access_token"`
	Fields      string `url:"fields"`
	DateFormat  string `url:"date_format"`
}

type CreateSubscribedAppsRequest struct {
	AccessToken string
	Fields      []string
	PageID      string
}

type CreateSubscribedAppsParams struct {
	AccessToken      string `url:"access_token"`
	SubscribedFields string `url:"subscribed_fields"`
}

type SendMessageRequest struct {
	AccessToken     string
	SendMessageArgs *model.SendMessageArgs
	PageID          string
}

type SendMessageParams struct {
	AccessToken string `url:"access_token"`
	Recipient   string `url:"recipient"`
	Message     string `url:"message"`
	Tag         string `url:"tag,omitempty"`
}

type SendCommentRequest struct {
	AccessToken     string
	SendCommentArgs *model.SendCommentArgs
	PageID          string
}

type SendCommentParams struct {
	AccessToken   string `url:"access_token"`
	Message       string `url:"message,omitempty"`
	AttachmentURL string `url:"attachment_url,omitempty"`
}

type LikeCommentRequest struct {
	AccessToken string
	PageID      string
	CommentID   string
}

type LikeCommentParams struct {
	AccessToken string `url:"access_token"`
}

type UnLikeCommentRequest struct {
	AccessToken string
	PageID      string
	CommentID   string
}

type UnLikeCommentParams struct {
	AccessToken string `url:"access_token"`
}

type HideOrUnHideCommentRequest struct {
	AccessToken string
	PageID      string
	CommentID   string
	IsHidden    bool
}

type HideOrUnHideCommentParams struct {
	AccessToken string `url:"access_token"`
	IsHidden    bool   `url:"is_hidden"`
}

type CreatePostRequest struct {
	AccessToken string
	PageID      string
	Content     *model.CreatePostRequest
}

type CreatePostParams struct {
	AccessToken string `url:"access_token"`
	Message     string `url:"message"`
}

type GetProfileRequest struct {
	AccessToken    string
	PSID           string
	PageID         string
	ProfileDefault *model.Profile
}

type GetProfileByPISDParams struct {
	AccessToken string `url:"access_token"`
	Fields      string `url:"fields"`
}

type FbRequest struct {
	Path   string
	Params interface{}
	Resp   interface{}
	PageID string
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
