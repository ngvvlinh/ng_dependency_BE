package imgroup

import (
	"context"
	"net/http"
	"time"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/httpreq"
	"etop.vn/common/l"
)

var ll = l.New()
var httpClient = httpreq.NewResty(httpreq.RestyConfig{
	Client: &http.Client{
		Timeout: 10 * time.Second,
	},
})

const requestURL = "http://itopx.webitop.com/sms"

type Config struct {
	APIKey string `yaml:"api_key"`
}

type Client struct {
	cfg Config
}

func New(cfg Config) *Client {
	if cfg.APIKey == "" {
		panic("missing api key")
	}
	return &Client{cfg: cfg}
}

type requestBody struct {
	Phone string `json:"phone"`
	SMS   string `json:"sms"`
}

type responseBody struct {
	Success bool `json:"success"`
}

func (c *Client) SendSMS(ctx context.Context, phone string, content string) (smsID string, _ error) {
	body := requestBody{
		Phone: phone,
		SMS:   content,
	}
	var result responseBody
	resp, err := httpClient.NewRequest().
		SetHeader("Authorization", c.cfg.APIKey).
		SetBody(&body).
		SetResult(&result).
		Post(requestURL)
	if err != nil {
		ll.Error("send sms", l.Error(err))
		return "", cm.Errorf(cm.ExternalServiceError, err, "imgroup sms: %v", err)
	}
	if !result.Success {
		ll.Error("send sms", l.String("resp", resp.String()))
		return "", cm.Errorf(cm.ExternalServiceError, err, "imgroup sms: %v", err)
	}
	return "", nil
}
