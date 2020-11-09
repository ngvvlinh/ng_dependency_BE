package client

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"time"

	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/httpreq"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/cmenv"
	"o.o/common/l"
)

const authorization = "Authorization"

var ll = l.New()

type Method string

const (
	POST Method = "POST"
)

type Client struct {
	baseUrl string
	token   string

	rclient *httpreq.Resty
}

func New(env string, cfg *SuiteCRMCfg) *Client {
	client := httpreq.RestyConfig{Client: &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}}
	c := &Client{
		token:   cfg.Token,
		rclient: httpreq.NewResty(client),
	}
	switch env {
	case cmenv.PartnerEnvTest, cmenv.PartnerEnvDev, cmenv.PartnerEnvProd:
		c.baseUrl = "https://suitecrm-uat.vnpost.vn/Api"
	default:
		ll.Fatal("suite crm: Invalid env")
	}
	return c
}

func (c *Client) GetToken() string {
	return c.token
}

// create ticket
func (c *Client) InsertCase(ctx context.Context, req *InsertCaseRequest) (*InsertCaseResponse, error) {
	var resp InsertCaseResponse
	if err := c.sendRequest(ctx, "/insertCase", req, &resp); err != nil {
		return nil, err
	}

	if resp.Error != "" {
		return nil, ErrorResp(ctx, resp.Error)
	}
	return &resp, nil
}

func ErrorResp(ctx context.Context, errMsg string) error {
	return cm.Errorf(cm.ExternalServiceError, nil, "Lỗi từ SuiteCRM: %v. Nếu cần thêm thông tin vui lòng liên hệ %v", errMsg, wl.X(ctx).CSEmail)
}

func (c *Client) sendRequest(ctx context.Context, path string, req, resp interface{}) error {
	res, err := c.rclient.R().
		SetHeader(authorization, "Bearer "+c.token).
		SetBody(req).
		Post(URL(c.baseUrl, path))
	if err != nil {
		return cm.Error(cm.ExternalServiceError, "Lỗi kết nối với SuiteCRM", err)
	}

	status := res.StatusCode()
	switch {
	case status >= 200 && status < 300:
		if err := json.Unmarshal(res.Body(), &resp); err != nil {
			return cm.Errorf(cm.ExternalServiceError, err, "Lỗi không xác định từ SuiteCRM: %v. Chúng tôi đang liên hệ với SuiteCRM để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ %v.", err, wl.X(ctx).CSEmail)
		}
		return nil
	default:
		return cm.Errorf(cm.ExternalServiceError, nil, "Lỗi không xác định từ SuiteCRM: Invalid status (%v). Chúng tôi đang liên hệ với SuiteCRM để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ %v.", status, wl.X(ctx).CSEmail)
	}
}
