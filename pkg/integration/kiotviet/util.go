package kiotviet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"go.uber.org/zap/zapcore"
	resty "gopkg.in/resty.v1"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/l"
)

type M map[string]string

type UpdatedQuery struct {
	LastUpdated      time.Time
	Page             int
	PageSize         int
	Desc             bool
	IncludeRemoveIDs bool
}

func (q *UpdatedQuery) Build() url.Values {
	if cm.IsZeroTime(q.LastUpdated) {
		ll.Panic("Zero time", l.Object("q", q))
	}
	if q.Page == 0 {
		ll.Panic("Zero page", l.Object("q", q))
	}
	if q.PageSize == 0 {
		q.PageSize = DefaultPageSize
	}

	fromItem := (q.Page - 1) * q.PageSize

	params := make(url.Values)
	params.Set("lastModifiedFrom", cm.JSONTime(q.LastUpdated))
	params.Set("currentItem", strconv.Itoa(fromItem))
	params.Set("pageSize", strconv.Itoa(q.PageSize))
	params.Set("orderBy", "modifiedDate")
	if q.IncludeRemoveIDs {
		params.Set("includeRemoveIds", "1")
	}
	if q.Desc {
		params.Set("orderDirection", "desc")
	}
	return params
}

func getJWTExpires(tokenStr string) time.Time {
	var claims jwt.StandardClaims
	jwt.ParseWithClaims(tokenStr, &claims, nil)
	if claims.ExpiresAt != 0 {
		return time.Unix(claims.ExpiresAt, 0)
	}
	return time.Time{}
}

var regEx = regexp.MustCompile(`^[0-9A-z-]{10,100}$`)

func validateInput(s string) bool {
	return regEx.MatchString(s)
}

type Response struct {
	*http.Response
	Body []byte
}

func getx(conn *Connection, reqURL string, successV interface{}) (*Response, error) {
	var errV *ErrorResponse
	resp, err := get(conn, reqURL, successV, &errV)
	if err != nil && errV != nil {
		return resp, errV
	}
	return resp, err
}

func postx(conn *Connection, reqURL string, token string, body, successV interface{}) (*Response, error) {
	var errV *ErrorResponse
	resp, err := post(conn, reqURL, body, successV, &errV)
	if err != nil && errV != nil {
		return resp, errV
	}
	return resp, err
}

func get(conn *Connection, reqURL string, successV, errV interface{}) (*Response, error) {
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	if conn != nil {
		req.Header.Add("Authorization", "Bearer "+conn.TokenStr)
		req.Header.Add("Retailer", conn.RetailerID)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	r, err := parseResponse(req, resp, successV, errV)
	logRequest(req, nil, r)
	return r, err
}

func post(conn *Connection, reqURL string, body, successV, errV interface{}) (*Response, error) {
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
		req.Header.Add("Retailer", conn.RetailerID)
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to Kiotviet (%v)", err)
	}
	r, err := parseResponse(req, resp, successV, errV)
	logRequest(req, reqData, r)
	return r, err
}

func parseResponse(req *http.Request, resp *http.Response, successV, errV interface{}) (*Response, error) {
	r := &Response{}
	r.Response = resp
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to Kiotviet (%v)", err)
	}

	r.Body = body
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err = json.Unmarshal(body, errV)
		if err != nil {
			return nil, cm.Error(cm.ExternalServiceError, "Lỗi từ Kiotviet: "+truncLength(string(body), 100), err)
		}
		return r, cm.Error(cm.ExternalServiceError, "Lỗi từ Kiotviet", nil)
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

	msg := "->Kiotviet " + req.Method + " " + req.URL.String()
	args := []zapcore.Field{
		l.Int("status", resp.StatusCode),
		l.String("Retailer", req.Header.Get("Retailer")),
		l.String("RateLimit", resp.Header.Get("Remaining")+"/"+resp.Header.Get("RateLimit")+" ("+t.String()+")"),
		l.String("Authorization", req.Header.Get("Authorization")),
		// l.String("\n->", string(reqBody)),
		// l.String("\n<-", string(resp.Body)),
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		ll.Info(msg, args...)
	} else {
		ll.Error(msg, args...)
	}
}

func truncLength(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}

func coalesceTime(t1, t2 Time) time.Time {
	if !t1.IsZero() {
		return t1.ToTime()
	}
	return t2.ToTime()
}

func wrapRequest(resp *resty.Response, err error) (*resty.Response, error) {
	if err != nil {
		return resp, cm.Error(cm.ExternalServiceError, "Lỗi kết nối với Kiotviet", err)
	}
	code := resp.StatusCode()
	if code < 200 || code >= 300 {
		respErr := resp.Error().(error)
		return resp, cm.Error(
			cm.ExternalServiceError,
			"Lỗi từ Kiotviet: "+respErr.Error(),
			respErr)
	}
	return resp, nil
}
