package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/schema"
	"gopkg.in/resty.v1"

	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/httpreq"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/cmenv"
	"o.o/common/jsonx"
	"o.o/common/l"
)

var ll = l.New()

type Client struct {
	baseUrl     string
	token       string
	username    string
	password    string
	tenantHost  string
	tenantToken string

	rclient *httpreq.Resty
}

type sendRequestArgs struct {
	url   string
	token string
	req   interface{}
	resp  interface{}
}

type Method string

const (
	PostMethod Method = "POST"
	GetMethod  Method = "GET"
)

var encoder = schema.NewEncoder()

func New(env string, cfg VHTAccountCfg) *Client {
	encoder.SetAliasTag("url")
	client := &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	rcfg := httpreq.RestyConfig{Client: client}
	c := &Client{
		token:       cfg.Token,
		username:    cfg.Username,
		password:    cfg.Password,
		tenantHost:  cfg.TenantHost,
		tenantToken: cfg.TenantToken,
		rclient:     httpreq.NewResty(rcfg),
	}

	switch env {
	case cmenv.PartnerEnvTest, cmenv.PartnerEnvDev:
		c.baseUrl = "https://sip.etelecom.vn:8900/api"
		if c.tenantHost == "" {
			c.tenantHost = "https://api-sandbox.etelecom.vn"
		}
	case cmenv.PartnerEnvProd:
		c.baseUrl = "https://sip.etelecom.vn:8900/api"
	default:
		ll.Fatal("Portsip: Invalid env")
	}

	return c
}

func (c *Client) UpdateToken(newToken string) {
	c.token = newToken
}

func (c *Client) Login(ctx context.Context) (*LoginResponse, error) {
	var resp LoginResponse
	err := c.sendPostRequest(ctx, sendRequestArgs{
		url:   URL(c.baseUrl, "/account/credentials/verify"),
		token: c.token,
		req: LoginRequest{
			Name:     c.username,
			Password: c.password,
		},
		resp: &resp,
	})
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) CreateExtension(ctx context.Context, req *CreateExtensionsRequest) (*CreateExtensionResponse, error) {
	var resp CreateExtensionResponse

	err := c.sendPostRequest(ctx, sendRequestArgs{
		url:   URL(c.baseUrl, "/extensions/create"),
		token: c.token,
		req:   req,
		resp:  &resp,
	})
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) GetCallLogs(ctx context.Context, req *GetCallLogsRequest) (*GetCallLogsResponse, error) {
	var resp GetCallLogsResponse

	err := c.sendGetRequest(ctx, sendRequestArgs{
		url:   URL(c.tenantHost, "/portsip-pbx/v1/cdr"),
		token: c.tenantToken,
		req:   req,
		resp:  &resp,
	})
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) sendPostRequest(ctx context.Context, args sendRequestArgs) error {
	return c.sendRequest(ctx, PostMethod, args)
}

func (c *Client) sendGetRequest(ctx context.Context, args sendRequestArgs) error {
	return c.sendRequest(ctx, GetMethod, args)
}

func (c *Client) sendRequest(ctx context.Context, method Method, args sendRequestArgs) error {
	var errResp ErrorResponse
	var res *resty.Response
	var err error

	request := c.rclient.R().
		SetResult(&args.resp).
		SetError(&errResp)
	if args.token != "" {
		request.SetHeader("access_token", args.token)
		request.SetHeader("Authorization", "Bearer "+args.token)
	}

	switch method {
	case PostMethod:
		res, err = request.SetBody(args.req).Post(args.url)
	case GetMethod:
		queryString := url.Values{}
		if args.req != nil {
			err := encoder.Encode(args.req, queryString)
			if err != nil {
				return cm.Error(cm.Internal, "", err)
			}
		}

		res, err = request.SetQueryString(queryString.Encode()).Get(args.url)
	default:
		panic(fmt.Sprintf("unsupported method %v", method))
	}
	if err != nil {
		return cm.Error(cm.ExternalServiceError, "Lỗi kết nối với VHT", err)
	}

	status := res.StatusCode()
	switch {
	case status == 200:
		if err := jsonx.Unmarshal(res.Body(), &args.resp); err != nil {
			return cm.Errorf(cm.ExternalServiceError, err, "Lỗi không xác định từ VHT: %v. Chúng tôi đang liên hệ với VHT để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ %v.", err, wl.X(ctx).CSEmail)
		}
		return nil
	case status >= 400:
		if err = jsonx.Unmarshal(res.Body(), &errResp); err != nil {
			return cm.Errorf(cm.ExternalServiceError, err, "Lỗi không xác định từ VHT: %v. Chúng tôi đang liên hệ với VHT để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ %v.", err, wl.X(ctx).CSEmail)
		}
		return cm.Errorf(cm.ExternalServiceError, &errResp, "Lỗi từ VHT: %v. Nếu cần thêm thông tin vui lòng liên hệ %v.", errResp.Error(), wl.X(ctx).CSEmail)
	default:
		return cm.Errorf(cm.ExternalServiceError, nil, "Lỗi không xác định từ VHT: Invalid status (%v). Chúng tôi đang liên hệ với VHT để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ %v.", status, wl.X(ctx).CSEmail)
	}
}
