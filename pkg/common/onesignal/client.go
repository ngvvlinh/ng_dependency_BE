package onesignal

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	resty "gopkg.in/resty.v1"

	"etop.vn/backend/pkg/etop/model"
	"github.com/gorilla/schema"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/httpreq"
)

const (
	APICreateNotification = "/notifications"
	APIGetDevices         = "/players"
)

var encoder = schema.NewEncoder()

func init() {
	encoder.SetAliasTag("url")
}

type Client struct {
	apiKey  string
	appID   string
	baseUrl string
	rclient *resty.Client
}

func New(appID string, apiKey string) *Client {
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	c := &Client{
		apiKey:  apiKey,
		appID:   appID,
		rclient: resty.NewWithClient(client).SetDebug(true),
		baseUrl: "https://onesignal.com/api/v1",
	}
	return c
}

func (c *Client) Ping() error {
	req := &GetDevicesRequest{
		AppID:  c.appID,
		Limit:  1,
		Offset: 0,
	}
	_, err := c.GetDevices(context.Background(), req)
	return err
}

func (c *Client) sendPostRequest(ctx context.Context, path string, req interface{}, resp ResponseInterface, msg string) error {
	res, err := c.rclient.R().
		SetBody(req).
		SetHeader("Authorization", fmt.Sprintf("Basic %v", c.apiKey)).
		Post(cm.URL(c.baseUrl, path))
	if err != nil {
		return cm.Error(cm.ExternalServiceError, "Lỗi kết nối với Onesignal", err)
	}
	err = handleResponse(res, resp, msg)
	return err
}

func (c *Client) sendGetRequest(ctx context.Context, path string, req interface{}, resp ResponseInterface, msg string) error {
	queryString := url.Values{}
	if req != nil {
		if err := encoder.Encode(req, queryString); err != nil {
			return cm.Error(cm.Internal, "", err)
		}
	}
	res, err := c.rclient.R().
		SetQueryString(queryString.Encode()).
		SetHeader("Authorization", fmt.Sprintf("Basic %v", c.apiKey)).
		Get(model.URL(c.baseUrl, path))
	if err != nil {
		return cm.Error(cm.ExternalServiceError, "Lỗi kết nối với Onesignal", err)
	}
	err = handleResponse(res, resp, msg)
	return err
}

// Use to Ping
func (c *Client) GetDevices(ctx context.Context, req *GetDevicesRequest) (*GetDevicesResponse, error) {
	endPoint := APIGetDevices
	var resp GetDevicesResponse
	err := c.sendGetRequest(ctx, endPoint, req, &resp, "Không thể lấy danh sách thiết bị.")
	return &resp, err
}

func (c *Client) CreateNotification(ctx context.Context, req *CreateNotificationRequest) (*CreateNotificationResponse, error) {
	req.AppID = c.appID
	var resp CreateNotificationResponse
	err := c.sendPostRequest(ctx, APICreateNotification, req, &resp, "Không thể push notification")
	return &resp, err
}

func handleResponse(res *resty.Response, result ResponseInterface, msg string) error {
	status := res.StatusCode()
	var err error
	body := res.Body()
	switch {
	case status >= 200 && status < 300:
		if result != nil {
			if httpreq.IsNullJsonRaw(body) {
				return cm.Error(cm.ExternalServiceError, "Lỗi không xác định từ Onesignal: null response.", nil)
			}
			if err = json.Unmarshal(body, result); err != nil {
				return cm.Errorf(cm.ExternalServiceError, err, "Lỗi không xác định từ Onesignal: %v.", err)
			}
			r := result.GetCommonResponse()
			if r.ID == "" && r.Errors != nil {
				return cm.Errorf(cm.ExternalServiceError, nil, "Lỗi từ Onesignal: %v (%v)", msg, r.Errors)
			}
		}
		return nil

	case status >= 400:
		var meta map[string]string
		if !httpreq.IsNullJsonRaw(body) {
			if err = json.Unmarshal(body, &meta); err != nil {
				// The slow path
				var metaX map[string]interface{}
				_ = json.Unmarshal(body, &metaX)
				meta = make(map[string]string)
				for k, v := range metaX {
					meta[k] = fmt.Sprint(v)
				}
			}
		}

		return cm.Errorf(cm.ExternalServiceError, nil, "Lỗi từ Onesignal.").WithMetaM(meta)
	default:
		return cm.Errorf(cm.ExternalServiceError, nil, "Lỗi không xác định từ Onesignal: Invalid status (%v).", status)
	}
}