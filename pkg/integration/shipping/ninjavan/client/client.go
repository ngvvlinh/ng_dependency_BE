package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

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
	baseUrl   string
	token     string
	clientID  string
	secretKey string

	rclient *httpreq.Resty
}

type sendRequestArgs struct {
	path string
	req  interface{}
	resp interface{}
}

type Method string

const (
	PostMethod   Method = "POST"
	DeleteMethod Method = "DELETE"
)

type ServiceType string
type ServiceLevel string

const (
	ServiceTypeParcel        ServiceType = "Parcel"
	ServiceTypeDocument      ServiceType = "Document"
	ServiceTypeReturn        ServiceType = "Return"
	ServiceTypeMarketPlace   ServiceType = "Marketplace"
	ServiceTypeBulky         ServiceType = "Bulky"
	ServiceTypeInternational ServiceType = "International"

	ServiceLevelStandard ServiceLevel = "Standard"
	ServiceLevelExpress  ServiceLevel = "Express"
	ServiceLevelSameDay  ServiceLevel = "Sameday"
	ServiceLevelNextDay  ServiceLevel = "Nextday"
)

func New(env string, cfg NinjaVanCfg) *Client {
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	rcfg := httpreq.RestyConfig{Client: client}
	c := &Client{
		token:     cfg.Token,
		clientID:  cfg.ClientID,
		secretKey: cfg.SecretKey,
		rclient:   httpreq.NewResty(rcfg),
	}
	switch env {
	case cmenv.PartnerEnvTest, cmenv.PartnerEnvDev:
		c.baseUrl = "https://api-sandbox.ninjavan.co/SG"
	case cmenv.PartnerEnvProd:
		c.baseUrl = "https://api.ninjavan.co/VN"
	default:
		ll.Fatal("NinjaVan: Invalid env")
	}
	return c
}

func (c *Client) UpdateToken(newToken string) {
	c.token = newToken
}

func (c *Client) GenerateOAuthAccessToken(ctx context.Context) (*GenerateAccessTokenResponse, error) {
	req := &GenerateAccessTokenRequest{
		ClientID:     c.clientID,
		ClientSecret: c.secretKey,
		GrantType:    ClientCredentials,
	}

	var resp GenerateAccessTokenResponse
	err := c.sendPostRequest(ctx, sendRequestArgs{
		path: "/2.0/oauth/access_token",
		req:  req,
		resp: &resp,
	})
	return &resp, err
}

// NJV does not support api get available services
// NJV only support Standard service
func (c *Client) FindAvailableServices(ctx context.Context) *FindAvailableServicesResponse {
	return &FindAvailableServicesResponse{
		AvailableServices: []*AvailableService{
			{
				Name: String(ServiceLevelStandard),
			},
		},
	}
}

func (c *Client) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*CreateOrderResponse, error) {
	var resp CreateOrderResponse
	err := c.sendPostRequest(ctx, sendRequestArgs{
		path: "/4.1/orders",
		req:  req,
		resp: &resp,
	})
	return &resp, err
}

func (c *Client) CancelOrder(ctx context.Context, trackingNo string) (*CancelOrderResponse, error) {
	var resp CancelOrderResponse
	path := fmt.Sprintf("/2.2/orders/%s", trackingNo)
	err := c.sendDeleteRequest(ctx, sendRequestArgs{
		path: path,
		resp: &resp,
	})
	return &resp, err
}

func (c *Client) sendPostRequest(ctx context.Context, args sendRequestArgs) error {
	return c.sendRequest(ctx, PostMethod, args)
}

func (c *Client) sendDeleteRequest(ctx context.Context, args sendRequestArgs) error {
	return c.sendRequest(ctx, DeleteMethod, args)
}

func (c *Client) sendRequest(ctx context.Context, method Method, args sendRequestArgs) error {
	var errResp ErrorResponse
	var res *resty.Response
	var err error

	request := c.rclient.R().
		SetBody(args.req).
		SetResult(&args.resp).
		SetError(&errResp)
	if c.token != "" {
		request.SetHeader("Authorization", "Bearer "+c.token)
	}

	switch method {
	case PostMethod:
		res, err = request.Post(URL(c.baseUrl, args.path))
	case DeleteMethod:
		res, err = request.Delete(URL(c.baseUrl, args.path))
	default:
		panic(fmt.Sprintf("unsupported method %v", method))
	}
	if err != nil {
		return cm.Error(cm.ExternalServiceError, "Lỗi kết nối với Ninja Van", err)
	}

	status := res.StatusCode()
	switch {
	case status == 200:
		return nil
	case status >= 400:
		if err = jsonx.Unmarshal(res.Body(), &errResp); err != nil {
			return cm.Errorf(cm.ExternalServiceError, err, "Lỗi không xác định từ Ninja Van: %v. Chúng tôi đang liên hệ với Giao Hang Nhanh để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ %v.", err, wl.X(ctx).CSEmail)
		}
		return cm.Errorf(cm.ExternalServiceError, &errResp, "Lỗi từ Ninja Van: %v. Nếu cần thêm thông tin vui lòng liên hệ %v.", errResp.Error(), wl.X(ctx).CSEmail)
	default:
		return cm.Errorf(cm.ExternalServiceError, nil, "Lỗi không xác định từ Ninja Van: Invalid status (%v). Chúng tôi đang liên hệ với Ninja Van để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ %v.", status, wl.X(ctx).CSEmail)
	}
}
