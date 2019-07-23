package vtiger

import (
	"net/http"
	"time"
)

// Client define header and url of request
type Client struct {
	httpClient *http.Client
	Token      string
	BaseURL    string
}

// NewClient create new Client
func NewClient(token, baseURL string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		Token:      token,
		BaseURL:    baseURL,
	}
}

// SessionInfo token information
type SessionInfo struct {
	UserInfo           UserSessionInfo    `json:"user"`
	AccountSessionInfo AccountSessionInfo `json:"account"`
}

// UserSessionInfo properties user in token information
type UserSessionInfo struct {
	ID        string `json:"id"`
	FullName  string `json:"full_name"`
	ShortName string `json:"short_name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}

// AccountSessionInfo properties account in token information
type AccountSessionInfo struct {
	AccountID   string `json:"id"`
	AccountName string `json:"name"`
	AccountType string `json:"type"`
	IsOperator  bool   `json:"is_operator"`
	Vtiger      string
}

// VtigerServiceTokenResponse get token response from vtiger
type VtigerServiceTokenResponse struct {
	Success bool `json:"success"`
	Result  *VtigerTokenResult
}

// VtigerTokenResult properti result of response get token
type VtigerTokenResult struct {
	Token      string `json:"token"`
	ServerTime int32  `json:"serverTime"`
	ExpireTime int32  `json:"expireTime"`
}

// VtigerSessionResponse session response from vtiger
type VtigerSessionResponse struct {
	Success bool `json:"success"`
	Result  *VtigerSessionResult
}

// VtigerSessionResult propoerty Result in session response
type VtigerSessionResult struct {
	SessionName   string `json:"sessionName"`
	UserID        string `json:"userId"`
	Version       string `json:"version"`
	VtigerVersion string `json:"vtigerVersion"`
}

// BodySessionRequest information for get session
type BodySessionRequest struct {
	Operation string `json:"operation"`
	Username  string `json:"username"`
	AccessKey string `json:"accessKey"`
}

// VtigerClient vtiger object contain information
type VtigerClient struct {
	BaseURL     string
	httpClient  *http.Client
	SessionInfo string
}

// NewVigerClient create VtigerClient
func NewVigerClient(sessionInfo string, baseUrl string) *VtigerClient {
	return &VtigerClient{
		SessionInfo: sessionInfo,
		BaseURL:     baseUrl,
		httpClient:  &http.Client{Timeout: 10 * time.Second},
	}
}

// VtigerResponse contain
type VtigerResponse struct {
	Success string              `json:"success"`
	Result  []map[string]string `json:"result"`
}

// VtigerConfig information vtiger's config
type VtigerConfig struct {
	VtigerService   string `yaml:"vtiger_service"`
	VtigerUsername  string `yaml:"vtiger_username"`
	VtigerAccesskey string `yaml:"vtiger_accesskey"`
}
