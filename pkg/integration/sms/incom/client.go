package incom

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"time"

	"gopkg.in/resty.v1"

	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/httpreq"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/validate"
	"o.o/backend/pkg/integration/sms/util"
	"o.o/common/l"
)

var ll = l.New()

type Config struct {
	APIKey    string `yaml:"api_key"`
	SecretKey string `yaml:"secret_key"`
	BrandName string `yaml:"brand_name"`
}

type Client struct {
	cfg     Config
	rClient *httpreq.Resty
	baseURL string
}

func New(cfg Config) *Client {
	if _, err := validate.ValidateStruct(cfg); err != nil {
		ll.Fatal("invalid incom sms config", l.Error(err))
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	rcfg := httpreq.RestyConfig{Client: client}
	return &Client{
		cfg:     cfg,
		rClient: httpreq.NewResty(rcfg),
		baseURL: "http://210.211.101.107",
	}
}

func (c *Client) SendSMS(ctx context.Context, phone string, content string) (string, error) {
	var resp SMSRes
	sms := []SMS{
		{
			ID:               cm.NewID().String(),
			BrandName:        c.cfg.BrandName,
			Text:             util.ModifyMsgPhone(content),
			To:               phone,
			FeeTypeID:        ByBrandNameID,
			MSGContentTypeID: AccentContentID,
		},
	}
	sendSMSRequest := SendSMSRequest{
		Submission: Submission{
			ApiKey:    c.cfg.APIKey,
			ApiSecret: c.cfg.SecretKey,
			Sms:       sms,
		}}
	if err := c.sendPostRequest(ctx, "/ccsmsunicode/Sms/SMSService.svc/ccsms/json", &sendSMSRequest, &resp); err != nil {
		return "", err
	}
	statusCode := resp.Status
	if _, ok := SMSResultCodeMap[statusCode]; !ok {
		return "", cm.Errorf(cm.ExternalServiceError, nil, "Lỗi từ Incom: nếu cần thêm thông tin vui lòng liên hệ %v.", wl.X(ctx).CSEmail)
	}

	// Chỉ succes khi status code = 0
	if statusCode != "0" {
		return "", cm.Errorf(cm.ExternalServiceError, nil, "Lỗi từ Incom: %v. Chúng tôi đang liên hệ với Incom để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ %v.", SMSResultCodeMap[statusCode], wl.X(ctx).CSEmail)
	}

	return resp.ID, nil
}

func (c *Client) sendPostRequest(ctx context.Context, path string, payload, resp interface{}) error {
	return c.sendRequest(ctx, path, httpreq.MethodPost, payload, resp)
}

func (c *Client) sendRequest(
	ctx context.Context, path string,
	method httpreq.RequestMethod, payload, resp interface{},
) (err error) {
	var (
		bodyResp CommonResponse
		res      *httpreq.RestyResponse
		req      *resty.Request
	)

	req = c.rClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		SetResult(&bodyResp).
		SetError(&bodyResp)

	switch method {
	case httpreq.MethodGet:
		res, err = req.Get(url(c.baseURL, path))
	case httpreq.MethodPost:
		res, err = req.Post(url(c.baseURL, path))
	case httpreq.MethodPut:
		res, err = req.Put(url(c.baseURL, path))
	default:
		return cm.Errorf(cm.Internal, nil, "Incom: unsupported method %v", method)
	}

	status := res.StatusCode()
	switch {
	case status >= 200 && status < 300:
		return c.handleResponseBody(ctx, &bodyResp, resp)
	case status >= 400:
		return cm.Errorf(cm.ExternalServiceError, nil, "Lỗi từ Incom: nếu cần thêm thông tin vui lòng liên hệ %v.", wl.X(ctx).CSEmail)
	default:
		return cm.Errorf(cm.ExternalServiceError, nil, "Lỗi không xác định từ Incom: Invalid status (%v). Chúng tôi đang liên hệ với Incom để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ %v.", status, wl.X(ctx).CSEmail)
	}
	return nil
}

func (c *Client) handleResponseBody(ctx context.Context, bodyResp *CommonResponse, resp interface{}) (err error) {
	if bodyResp == nil || bodyResp.Response == nil || bodyResp.Response.Submission == nil || len(bodyResp.Response.Submission.Sms) < 1 {
		return cm.Errorf(cm.ExternalServiceError, nil, "Lỗi từ Incom: nếu cần thêm thông tin vui lòng liên hệ %v.", wl.X(ctx).CSEmail)
	}
	if err = json.Unmarshal(bodyResp.Response.Submission.Sms[0], &resp); err != nil {
		return err
	}
	return nil
}

func url(baseUrl, path string) string {
	return baseUrl + path
}
