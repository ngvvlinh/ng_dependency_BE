package client

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	model2 "etop.vn/backend/com/supporting/crm/vht/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/common/jsonx"
)

var (
	MethodGET                = "GET"
	URLGetHistory            = "https://acd-api.vht.com.vn/rest/cdrs"
	URLGetHistoryBySdkCallID = "http://acd-api.vht.com.vn/rest/softphones/cdrs"
)

func (c *Client) GetHistories(queryDTO *VHTHistoryQueryDTO) (*VhCallHistoriesResponse, error) {
	query, err := c.MakeQueryVht(queryDTO)
	if err != nil {
		return nil, err
	}
	var vhtResp VhCallHistoriesResponse
	err = c.Request(URLGetHistory, MethodGET, query, &vhtResp)
	if err != nil {
		return nil, err
	}
	return &vhtResp, nil
}

func (c *Client) GetHistoryBySDKCallID(sdkCallID string) (*VhtCallHistory, error) {
	u, err := url.Parse("")
	if err != nil {
		return nil, err
	}
	query := u.Query()
	query.Set("sdk_call_id", sdkCallID)
	u.RawQuery = query.Encode()
	var vhtResp VhtCallHistory
	err = c.Request(URLGetHistoryBySdkCallID, MethodGET, u.String(), &vhtResp)
	if err != nil {
		return nil, err
	}
	return &vhtResp, nil
}

func (c *Client) Request(baseUrl string, method string, query string, vhtResp interface{}) error {
	requestURL := baseUrl + query
	request, err := http.NewRequest(method, requestURL, nil)
	if err != nil {
		return err
	}
	request.SetBasicAuth(c.UserName, c.PassWord)
	request.Header.Add("Content-Type", "application/json")

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = jsonx.Unmarshal(respBytes, vhtResp)
	if err != nil {
		return err
	}

	return nil
}

func ConvertToModel(callHistory *VhtCallHistory) *model2.VhtCallHistory {
	return &model2.VhtCallHistory{
		CdrID:           callHistory.CdrID,
		CallID:          callHistory.CallID,
		SipCallID:       callHistory.SipCallID,
		SdkCallID:       callHistory.SdkCallID,
		Cause:           callHistory.Cause,
		Q850Cause:       callHistory.Q850Cause,
		FromExtension:   callHistory.FromExtension,
		ToExtension:     callHistory.ToExtension,
		FromNumber:      callHistory.FromNumber,
		ToNumber:        callHistory.ToNumber,
		Duration:        callHistory.Duration,
		Direction:       callHistory.Direction,
		TimeStarted:     time.Unix(callHistory.TimeStarted, 0),
		TimeConnected:   time.Unix(callHistory.TimeConnected, 0),
		TimeEnded:       time.Unix(callHistory.TimeEnd, 0),
		RecordingPath:   callHistory.RecordingPath,
		RecordingURL:    callHistory.RecordingUrl,
		RecordFileSize:  callHistory.RecordFileSize,
		EtopAccountID:   callHistory.EtopAccountID,
		VtigerAccountID: callHistory.VtigerAccountID,
		SyncStatus:      "Done",
	}
}

func (c *Client) MakeQueryVht(queryDTO *VHTHistoryQueryDTO) (string, error) {
	u, err := url.Parse("")
	if err != nil {
		return "", err
	}
	query := u.Query()
	if queryDTO.Page != 0 {
		query.Set("page", strconv.FormatInt(int64(queryDTO.Page), 10))
	}
	if queryDTO.Limit != 0 {
		query.Set("limit", strconv.FormatInt(int64(queryDTO.Limit), 10))
	}
	if queryDTO.Direction != 0 {
		query.Set("direction", strconv.FormatInt(int64(queryDTO.Direction), 10))
	}
	if queryDTO.Extension != 0 {
		query.Set("extension", strconv.FormatInt(int64(queryDTO.Extension), 10))
	}
	if queryDTO.DateStarted != 0 {
		query.Set("date_started", strconv.FormatInt(int64(queryDTO.DateStarted), 10))
	}
	if queryDTO.DateEnded != 0 {
		query.Set("date_ended", strconv.FormatInt(int64(queryDTO.DateEnded), 10))
	}
	if queryDTO.SortBy != "" {
		query.Set("sort_by", queryDTO.SortBy)
	}
	if queryDTO.SortType != "" {
		query.Set("sort_type", queryDTO.SortType)
	}
	if queryDTO.State != "" {
		query.Set("state", queryDTO.State)
	}
	if queryDTO.FromNumber != "" {
		query.Set("from_number", queryDTO.FromNumber)
	}
	if queryDTO.ToNumber != "" {
		query.Set("to_number", queryDTO.ToNumber)
	}
	u.RawQuery = query.Encode()
	return u.String(), nil
}

func (c *Client) PingServerVht() error {
	queryDTO := VHTHistoryQueryDTO{
		Page:  1,
		Limit: 1,
	}
	query, err := c.MakeQueryVht(&queryDTO)
	if err != nil {
		return err
	}
	var vhtResp VhCallHistoriesResponse
	err = c.Request(URLGetHistory, "GET", query, &vhtResp)
	if err != nil {
		return err
	}
	if len(vhtResp.Items) == 0 {
		return cm.Error(cm.ExternalServiceError, "Can't connect to Vht server", nil)
	}
	return nil

}
