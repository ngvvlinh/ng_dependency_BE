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
	baseUrl  string
	token    string
	username string
	password string

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

func New(cfg PortsipAccountCfg) *Client {
	encoder.SetAliasTag("url")
	client := &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	rcfg := httpreq.RestyConfig{Client: client}
	c := &Client{
		token:    cfg.Token,
		username: cfg.Username,
		password: cfg.Password,
		rclient:  httpreq.NewResty(rcfg),
	}

	switch cmenv.Env() {
	case cmenv.EnvDev:
		c.baseUrl = "https://sipdev.dinodata.vn:8900/api"
	case cmenv.EnvSandbox, cmenv.EnvStag:
		c.baseUrl = "https://sip.dinodata.vn:8900/api"
	case cmenv.EnvProd:
		c.baseUrl = "https://sip.dinodata.vn:8900/api"
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

// This feature is available to admin user only
// Doc: https://www.portsip.com/pbx-rest-api/pbx/account.html#add_account
func (c *Client) CreateTenant(ctx context.Context, req *CreateTenantRequest) (*CreateTenantResponse, error) {
	var resp CreateTenantResponse
	err := c.sendPostRequest(ctx, sendRequestArgs{
		url:   URL(c.baseUrl, "/account/create"),
		token: c.token,
		req:   req,
		resp:  &resp,
	})
	return &resp, err
}

// This feature is available to admin user only
func (c *Client) UpdateTrunkProvider(ctx context.Context, req *UpdateTrunkProviderRequest) error {
	err := c.sendPostRequest(ctx, sendRequestArgs{
		url:   URL(c.baseUrl, "/providers/update"),
		token: c.token,
		req:   req,
		resp:  nil,
	})
	return err
}

func (c *Client) GetTrunkProvider(ctx context.Context, req *GetTrunkProviderRequest) (*TrunkProvider, error) {
	var resp TrunkProvider
	err := c.sendGetRequest(ctx, sendRequestArgs{
		url:   URL(c.baseUrl, "/providers/show"),
		token: c.token,
		req:   req,
		resp:  &resp,
	})
	return &resp, err
}

func (c *Client) ListOutboundRules(ctx context.Context, req *CommonListRequest) (*ListOutboundRulesResponse, error) {
	var resp ListOutboundRulesResponse
	err := c.sendGetRequest(ctx, sendRequestArgs{
		url:   URL(c.baseUrl, "/outbound_rules/list"),
		token: c.token,
		req:   req,
		resp:  &resp,
	})
	return &resp, err
}

func (c *Client) CreateOutboundRule(ctx context.Context, req *CreateOutboundRuleRequest) (*CreateOutboundRuleResponse, error) {
	var resp CreateOutboundRuleResponse
	err := c.sendPostRequest(ctx, sendRequestArgs{
		url:   URL(c.baseUrl, "/outbound_rules/create"),
		token: c.token,
		req:   req,
		resp:  &resp,
	})
	return &resp, err
}

func (c *Client) GetExtensionGroups(ctx context.Context, req *CommonListRequest) (*GetExtensionGroupsResponse, error) {
	var resp GetExtensionGroupsResponse
	err := c.sendGetRequest(ctx, sendRequestArgs{
		url:   URL(c.baseUrl, "/extensions/group/list"),
		token: c.token,
		req:   req,
		resp:  &resp,
	})
	return &resp, err
}

func (c *Client) DestroyCallSesssion(ctx context.Context, req *DestroyCallSessionRequest) error {
	err := c.sendPostRequest(ctx, sendRequestArgs{
		url:   URL(c.baseUrl, "/call_sessions/destroy"),
		token: c.token,
		req:   req,
		resp:  nil,
	})
	return err
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
			err = encoder.Encode(args.req, queryString)
			if err != nil {
				return cm.Error(cm.Internal, "", err)
			}
		}

		res, err = request.SetQueryString(queryString.Encode()).Get(args.url)
	default:
		panic(fmt.Sprintf("unsupported method %v", method))
	}
	if err != nil {
		return cm.Error(cm.ExternalServiceError, "Lỗi kết nối với Portsip", err)
	}

	status := res.StatusCode()
	switch {
	case status == 200:
		if args.resp == nil {
			return nil
		}
		if err = jsonx.Unmarshal(res.Body(), &args.resp); err != nil {
			return cm.Errorf(cm.ExternalServiceError, err, "Lỗi không xác định từ Portsip: %v. Chúng tôi đang liên hệ với Portsip để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ %v.", err, wl.X(ctx).CSEmail)
		}
		return nil
	case status >= 400:
		if err = jsonx.Unmarshal(res.Body(), &errResp); err != nil {
			return cm.Errorf(cm.ExternalServiceError, err, "Lỗi không xác định từ Portsip: %v. Chúng tôi đang liên hệ với Portsip để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ %v.", err, wl.X(ctx).CSEmail)
		}

		code := cm.ExternalServiceError
		if errResp.ErrCode.String() == string(NameOrDomainIncorrect) {
			code = cm.PortsipNameOrDomainIncorrect
		}
		return cm.Errorf(code, &errResp, "Lỗi từ Portsip: %v. Nếu cần thêm thông tin vui lòng liên hệ %v.", errResp.Error(), wl.X(ctx).CSEmail)
	default:
		return cm.Errorf(cm.ExternalServiceError, nil, "Lỗi không xác định từ Portsip: Invalid status (%v). Chúng tôi đang liên hệ với Portsip để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ %v.", status, wl.X(ctx).CSEmail)
	}
}
