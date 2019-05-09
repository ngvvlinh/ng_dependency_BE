package httpreq

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap/zapcore"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/l"
)

var (
	ll     = l.New()
	client = &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
)

type Response struct {
	*http.Response
	Body []byte
}

type AccessTokenResponse struct {
	TokenStr     string `json:"access_token"`
	ExpiresIn    Int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
}

type ErrorResponse struct {
	ErrorStr       string `json:"error"`
	ResponseStatus struct {
		ErrorCode string `json:"errorCode"`
		Message   string `json:"message"`
	} `json:"responseStatus"`
}

func (r *ErrorResponse) Error() string {
	if r.ErrorStr != "" {
		return r.ErrorStr
	}
	if r.ResponseStatus.ErrorCode != "" {
		if r.ResponseStatus.Message == "" {
			return r.ResponseStatus.ErrorCode
		}
		return r.ResponseStatus.Message +
			" (" + r.ResponseStatus.ErrorCode + ")"
	}
	data, _ := json.Marshal(r)
	return string(data)
}

type Connection struct {
	TokenStr  string
	ExpiresAt time.Time
	Headers   map[string]string
}

func GetX(conn *Connection, reqURL string, successV interface{}) (*Response, error) {
	var errV *ErrorResponse
	resp, err := Get(conn, reqURL, successV, &errV)
	if err != nil && errV != nil {
		return resp, errV
	}
	return resp, err
}

func PostX(conn *Connection, reqURL string, token string, body, successV interface{}) (*Response, error) {
	var errV *ErrorResponse
	resp, err := Post(conn, reqURL, body, successV, &errV)
	if err != nil && errV != nil {
		return resp, errV
	}
	return resp, err
}

func Get(conn *Connection, reqURL string, successV, errV interface{}) (*Response, error) {
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	if conn != nil {
		req.Header.Add("Authorization", "Bearer "+conn.TokenStr)
		for k, v := range conn.Headers {
			req.Header.Add(k, v)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	r, err := ParseResponse(resp, successV, errV)
	logRequest(req, nil, r)
	return r, err
}

func Post(conn *Connection, reqURL string, body, successV interface{}, errV interface{}) (*Response, error) {
	var reqBody io.Reader
	var reqData []byte
	if body != nil {
		switch body := body.(type) {
		case []byte:
			reqBody = bytes.NewBuffer(body)
			reqData = body
		case string:
			reqBody = bytes.NewBufferString(body)
			reqData = []byte(body)
		default:
			data, err := json.Marshal(body)
			if err != nil {
				return nil, err
			}
			reqBody = bytes.NewBuffer(data)
			reqData = data
		}
	}
	req, err := http.NewRequest("POST", reqURL, reqBody)
	if err != nil {
		return nil, err
	}
	if conn != nil {
		req.Header.Add("Authorization", "Bearer "+conn.TokenStr)
		for k, v := range conn.Headers {
			req.Header.Add(k, v)
		}
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to external service (%v)", err)
	}
	r, err := ParseResponse(resp, successV, errV)
	logRequest(req, reqData, r)
	return r, err
}

func ParseResponse(resp *http.Response, successV, errV interface{}) (*Response, error) {
	r := &Response{}
	r.Response = resp
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, cm.Error(cm.ExternalServiceError, "", err)
	}

	r.Body = body
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err = json.Unmarshal(body, errV)
		if err != nil {
			return nil, cm.Error(cm.ExternalServiceError, "Error parsing response", err)
		}
		return r, nil
	}

	if successV != nil {
		err = json.Unmarshal(body, successV)
		if err != nil {
			return r, err
		}
	}
	return r, err
}

func logRequest(req *http.Request, reqBody []byte, resp *Response) {
	tReset, _ := strconv.Atoi(resp.Header.Get("Reset"))
	t := time.Unix(int64(tReset), 0)

	msg := "->Request " + req.Method + " " + req.URL.String()
	args := []zapcore.Field{
		l.Int("status", resp.StatusCode),
		l.String("RateLimit", resp.Header.Get("Remaining")+"/"+resp.Header.Get("RateLimit")+" ("+t.String()+")"),
		l.String("Authorization", req.Header.Get("Authorization")),
		l.String("\n->", string(reqBody)),
		l.String("\n<-", string(resp.Body))}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		ll.Info(msg, args...)
	} else {
		ll.Error(msg, args...)
	}
}

func IsNullJsonRaw(data json.RawMessage) bool {
	return len(data) == 0 ||
		len(data) == 4 && string(data) == "null"
}
