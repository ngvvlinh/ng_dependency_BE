package client

import (
	"net/http"
	"time"

	"etop.vn/capi/dot"
)

// Client define header and url of request
type Client struct {
	httpClient *http.Client
	UserName   string
	PassWord   string
}

// NewClient create new Client
func NewClient(userName string, passWord string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 20 * time.Second},
		UserName:   userName,
		PassWord:   passWord,
	}
}

// query dto
type VHTHistoryQueryDTO struct {
	Page        int
	Limit       int
	SortBy      string
	SortType    string //: 'ASC' | 'DESC';
	State       string
	Direction   int
	Extension   int
	FromNumber  string
	ToNumber    string
	DateStarted int64
	DateEnded   int64
}

type VhtCallHistory struct {
	CdrID           string `json:"cdr_id"`
	CallID          string `json:"call_id"`
	SipCallID       string `json:"sip_call_id"`
	Cause           string `json:"cause"`
	SdkCallID       string `json:"sdk_call_id"`
	Q850Cause       string `json:"q850_cause"`
	FromExtension   string `json:"from_extension"`
	ToExtension     string `json:"to_extension"`
	FromNumber      string `json:"from_number"`
	ToNumber        string `json:"to_number"`
	Duration        int    `json:"duration"`
	Direction       int    `json:"direction"`
	TimeStarted     int64  `json:"time_started"`
	TimeConnected   int64  `json:"time_connected"`
	TimeEnd         int64  `json:"time_ended"`
	RecordingPath   string `json:"recording_path"`
	RecordingUrl    string `json:"recording_url"`
	RecordFileSize  int    `json:"record_file_size"`
	EtopAccountID   dot.ID `json:"etop_account_"id`
	VtigerAccountID string `json:"vtiger_account_id"`
}

// Response
type VhCallHistoriesResponse struct {
	Total       int               `json:"total"`
	Currentpage int               `json:"currentpage"`
	Limit       int               `json:"limit"`
	Items       []*VhtCallHistory `json:"items"`
}
