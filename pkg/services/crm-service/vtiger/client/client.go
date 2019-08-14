package client

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	cm "etop.vn/backend/pkg/common"
	"github.com/gorilla/schema"
)

var schemaEncoder = schema.NewEncoder()

func init() {
	schemaEncoder.SetAliasTag("json")
}

// VtigerClient vtiger object contain information
type VtigerClient struct {
	ServiceURL  string
	httpClient  *http.Client
	SessionInfo Session
}

type Session struct {
	VtigerSession *VtigerSessionResult
	ExpriredTime  int64
}

// NewVigerClient create VtigerClient
func NewVigerClient(serviceURL string) *VtigerClient {
	return &VtigerClient{
		SessionInfo: Session{},
		ServiceURL:  serviceURL,
		httpClient:  &http.Client{Timeout: 10 * time.Second},
	}
}

// GetSessionKey which use in excute vtiger API
func (v *VtigerClient) GetSessionKey(vtigerService string, vtigerUser string, vtigerAccessKey string) (*VtigerSessionResult, error) {

	emptySession := Session{}
	now := time.Now()
	sec := now.Unix()

	if v.SessionInfo != emptySession && (v.SessionInfo.ExpriredTime-sec) > 60 {
		return v.SessionInfo.VtigerSession, nil
	}

	requestValues := make(url.Values)
	requestValues.Set("operation", "getchallenge")
	requestValues.Set("username", vtigerUser)

	var responseVtigerMap VtigerServiceTokenResponse

	err := v.SendGet(requestValues, &responseVtigerMap)
	if err != nil {
		return nil, err
	}

	// get session
	accessKey := calcMD5Hash(responseVtigerMap.Result.Token + vtigerAccessKey)
	bodySessionRequest := BodySessionRequest{
		Operation: "login",
		Username:  vtigerUser,
		AccessKey: accessKey,
	}
	requestBody := url.Values{}
	if err = schemaEncoder.Encode(bodySessionRequest, requestBody); err != nil {
		panic(err)
	}

	var vtigerSessionResponse VtigerSessionResponse
	err = v.SendPost(requestBody, &vtigerSessionResponse)
	if err != nil {
		return nil, err
	}
	if !vtigerSessionResponse.Success {
		return nil, cm.Errorf(cm.Unknown, nil, "unknown error while sending request to vtiger")
	}

	v.SessionInfo.ExpriredTime = responseVtigerMap.Result.ExpireTime
	v.SessionInfo.VtigerSession = vtigerSessionResponse.Result

	return vtigerSessionResponse.Result, nil
}

// Get sends Get request to vtiger
func (v *VtigerClient) RequestGet(values url.Values) (*VtigerResponse, error) {
	var resp VtigerResponse
	err := v.SendGet(values, &resp)
	if err != nil {
		return nil, err
	}
	if !resp.Success {
		return nil, cm.Errorf(cm.Unknown, nil, "unknown error while sending request to vtiger")
	}
	return &resp, nil
}

func (v *VtigerClient) SendGet(values url.Values, respBody interface{}) error {
	requestURL := mustParseURL(v.ServiceURL)
	requestURL.RawQuery = values.Encode()
	request, err := http.NewRequest("GET", requestURL.String(), nil)
	if err != nil {
		return err
	}
	resp, err := v.httpClient.Do(request)
	if err != nil {
		return err
	}
	bodyResponse, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(bodyResponse, respBody)
}

func (v *VtigerClient) SendPost(body url.Values, respBody interface{}) error {
	request, err := http.NewRequest("POST", v.ServiceURL, strings.NewReader(body.Encode()))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := v.httpClient.Do(request)
	if err != nil {
		return err
	}

	respBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(respBytes, respBody)

	return err
}

// calcMD5Hash calculate MD5 hash of a string
func calcMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func mustParseURL(s string) *url.URL {
	reqUrl, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	return reqUrl
}
