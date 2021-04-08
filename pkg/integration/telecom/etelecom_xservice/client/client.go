package client

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/schema"
	"gopkg.in/resty.v1"

	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/httpreq"
	"o.o/backend/pkg/common/cmenv"
	"o.o/common/l"
)

// This service is built by eTelecom's dev in-house
// Use to connect tenant Portsip, config to get CDR
type Client struct {
	baseUrl string
	// token: api_key which is generated for shop partner account eTop
	token string

	rclient *httpreq.Resty
}

var (
	encoder = schema.NewEncoder()
	ll      = l.New()
)

func New(token string) *Client {
	encoder.SetAliasTag("url")
	client := &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	rcfg := httpreq.RestyConfig{Client: client}
	c := &Client{
		token:   token,
		rclient: httpreq.NewResty(rcfg),
	}
	switch cmenv.Env() {
	case cmenv.EnvDev:
		c.baseUrl = "https://api-dev.etelecom.vn"
	case cmenv.EnvSandbox, cmenv.EnvStag:
		c.baseUrl = "https://api-sandbox.etelecom.vn"
	case cmenv.EnvProd:
		c.baseUrl = "https://api.etelecom.vn"
	default:
		ll.Fatal("eTelecom external service: Invalid env")
	}
	return c
}

type sendRequestArgs struct {
	Path     string
	Req      interface{}
	Resp     interface{}
	Method   httpreq.RequestMethod
	ErrorMsg string
}

func (c *Client) GetCallLogs(ctx context.Context, req *GetCallLogsRequest) (*GetCallLogsResponse, error) {
	var resp GetCallLogsResponse
	args := sendRequestArgs{
		Path:     URL(c.baseUrl, "/portsip-pbx/v1/cdr"),
		Req:      req,
		Resp:     &resp,
		Method:   httpreq.MethodGet,
		ErrorMsg: "Get call logs failed",
	}
	err := c.sendRequest(ctx, args)
	return &resp, err
}

func (c *Client) ConfigTenantCDR(ctx context.Context, req *ConfigTenantCDRRequest) (*ConfigTenantCDRResponse, error) {
	var resp ConfigTenantCDRResponse
	args := sendRequestArgs{
		Path:     URL(c.baseUrl, "/portsip-pbx/v1/accounts"),
		Req:      req,
		Resp:     &resp,
		Method:   httpreq.MethodPost,
		ErrorMsg: "Config tenant CDR failed",
	}
	err := c.sendRequest(ctx, args)
	return &resp, err
}

func (c *Client) sendRequest(ctx context.Context, args sendRequestArgs) error {
	var res *resty.Response
	var err error
	req := c.rclient.R().
		SetHeader("Authorization", "Bearer "+c.token)
	switch args.Method {
	case httpreq.MethodGet:
		queryString := url.Values{}
		if req != nil {
			if err = encoder.Encode(req, queryString); err != nil {
				return cm.Errorf(cm.Internal, err, "")
			}
		}
		res, err = req.SetQueryString(queryString.Encode()).
			Get(args.Path)
	case httpreq.MethodPost:
		res, err = req.SetBody(args.Req).Post(args.Path)
	default:
		return cm.Errorf(cm.Internal, nil, "eTelecom external service does not support request %v method", args.Method)
	}

	if err != nil {
		return cm.Errorf(cm.Internal, err, "Lỗi kết nối với eTelecom external service")
	}
	if err = httpreq.HandleResponse(ctx, res, args.Resp, args.ErrorMsg); err != nil {
		return cm.Errorf(cm.ExternalServiceError, err, "Lỗi từ eTelecom external service: %v", err.Error())
	}
	return nil
}

func URL(baseUrl, path string) string {
	return baseUrl + path
}
