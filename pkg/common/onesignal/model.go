package onesignal

import "encoding/json"

type CommonParams struct {
	AppID string `json:"app_id"`
}

type MultipleContentLanguages struct {
	EN string `json:"en"`
	VI string `json:"vi"`
}

type CreateNotificationRequest struct {
	CommonParams
	// Limit of 2,000 entries per REST API call
	IncludePlayerIDs []string                 `json:"include_player_ids"`
	Headings         MultipleContentLanguages `json:"headings"`
	Contents         MultipleContentLanguages `json:"contents"`
	Subtitle         MultipleContentLanguages `json:"subtitle"`
	Data             json.RawMessage          `json:"data"`
	WebURL           string                   `json:"web_url"`
}

type CreateNotificationResponse struct {
	CommonResponse
	// ID         string `json:"id"`
	Recipients int `json:"recipients"`
}

type CommonResponse struct {
	ID     string      `json:"id"`
	Errors interface{} `json:"errors"`
}

type ResponseInterface interface {
	GetCommonResponse() CommonResponse
}

func (c CommonResponse) GetCommonResponse() CommonResponse {
	return c
}

type GetDevicesRequest struct {
	AppID  string `url:"app_id"`
	Limit  int    `url:"limit"`
	Offset int    `url:"offset"`
}

type GetDevicesResponse struct {
	CommonResponse
	TotalCount int `json:"total_count"`
	Offset     int `json:"offset"`
	Limit      int `json:"limit"`
}
