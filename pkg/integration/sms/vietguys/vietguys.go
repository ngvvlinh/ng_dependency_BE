package vietguys

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/validate"
	"o.o/common/l"
)

var ll = l.New()
var client = &http.Client{
	Timeout: 10 * time.Second,
}

const requestURL = "https://cloudsms.vietguys.biz:4438/api/"

type Config struct {
	BrandName string `yaml:"brand_name"  valid:"required"`
	Username  string `yaml:"username"    valid:"required"`
	APIKey    string `yaml:"api_key"     valid:"required"`
}

type Client struct {
	cfg Config
}

func New(cfg Config) *Client {
	if _, err := validate.ValidateStruct(cfg); err != nil {
		ll.Fatal("invalid vietguys sms config", l.Error(err))
	}
	return &Client{cfg}
}

func (c *Client) SendSMS(ctx context.Context, phone string, content string) (smsID string, _ error) {
	u, err := url.Parse(requestURL)
	if err != nil {
		panic(err)
	}
	q := u.Query()
	q.Set("u", c.cfg.Username)
	q.Set("pwd", c.cfg.APIKey)
	q.Set("from", c.cfg.BrandName)
	q.Set("phone", phone)
	q.Set("sms", content)
	q.Set("bid", cm.NewID().String())
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return "", cm.Errorf(cm.Internal, err, "không thể gửi tin nhắn")
	}
	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return "", cm.Errorf(cm.ExternalServiceError, err, "không thể gửi tin nhắn")
	}
	if resp.StatusCode != 200 {
		return "", cm.Errorf(cm.ExternalServiceError, nil, "không thể gửi tin nhắn")
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", cm.Errorf(cm.ExternalServiceError, err, "không thể gửi tin nhắn")
	}
	if len(data) <= 3 {
		return "", cm.Errorf(cm.ExternalServiceError, nil, "không thể gửi tin nhắn (mã lỗi %v)", string(data))
	}
	return string(data), nil
}
