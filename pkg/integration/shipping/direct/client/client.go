package client

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"github.com/gorilla/schema"

	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/httpreq"
)

var encoder = schema.NewEncoder()

func init() {
	encoder.SetAliasTag("json")
}

type Client struct {
	connName               string
	trackingURL            string
	createFulfillmentURL   string
	getFulfillmentURL      string
	getShippingServicesURL string
	cancelFulfillmentURL   string
	signInURL              string
	signUpURL              string

	affiliateID string
	token       string
	rclient     *httpreq.Resty
}

func New(cfg PartnerAccountCfg) (*Client, error) {
	conn := cfg.Connection
	if err := validateConnection(conn); err != nil {
		return nil, cm.Errorf(cm.FailedPrecondition, err, err.Error()).WithMetap("connection", conn)
	}
	// if cfg.Token == "" {
	// 	return nil, cm.Errorf(cm.FailedPrecondition, nil, "Không thể khởi tạo direct shipment cho connection '%v'. Token không được để trống.", conn.Name)
	// }

	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	rcfg := httpreq.RestyConfig{Client: client}
	driverCfg := conn.DriverConfig
	return &Client{
		connName:               conn.Name,
		trackingURL:            driverCfg.TrackingURL,
		createFulfillmentURL:   driverCfg.CreateFulfillmentURL,
		getFulfillmentURL:      driverCfg.GetFulfillmentURL,
		getShippingServicesURL: driverCfg.GetShippingServicesURL,
		cancelFulfillmentURL:   driverCfg.CancelFulfillmentURL,
		signInURL:              driverCfg.SignInURL,
		signUpURL:              driverCfg.SignUpURL,
		affiliateID:            cfg.AffiliateID,
		token:                  cfg.Token,
		rclient:                httpreq.NewResty(rcfg),
	}, nil
}

func (c *Client) Ping() error {
	return cm.ErrTODO
}

func (c *Client) GetShippingServices(ctx context.Context, req *GetShippingServicesRequest) ([]*ShippingService, error) {
	var resp GetShippingServiceResponse
	if err := c.sendPostRequest(ctx, c.getShippingServicesURL, req, &resp, "Không thể lấy gói vận chuyển"); err != nil {
		return nil, err
	}
	for _, s := range resp.Data {
		if err := s.validate(); err != nil {
			return nil, cm.Errorf(cm.ExternalServiceError, err, "Gói dịch vụ trả về không hợp lệ: %v", err.Error())
		}
	}
	return resp.Data, nil
}

func (c *Client) CreateFulfillment(ctx context.Context, req *CreateFulfillmentRequest) (*Fulfillment, error) {
	var resp CreateFulfillmentResponse
	req.AffiliateID = c.affiliateID
	if err := c.sendPostRequest(ctx, c.createFulfillmentURL, req, &resp, "Không thể tạo đơn giao hàng"); err != nil {
		return nil, err
	}
	if err := resp.Data.validate(); err != nil {
		return nil, cm.Errorf(cm.ExternalServiceError, err, "Đơn giao hàng không hợp lệ: %v", err.Error())
	}
	return resp.Data, nil
}

func (c *Client) GetFulfillment(ctx context.Context, req *GetFulfillmentRequest) (*Fulfillment, error) {
	var resp GetFulfillmentResponse
	if err := c.sendPostRequest(ctx, c.getFulfillmentURL, req, &resp, "Không thể lấy thông tin đơn giao hàng"); err != nil {
		return nil, err
	}
	if err := resp.Data.validate(); err != nil {
		return nil, cm.Errorf(cm.ExternalServiceError, err, "Thông tin đơn giao hàng không hợp lệ: %v", err.Error())
	}
	return resp.Data, nil
}

func (c *Client) CancelFulfillment(ctx context.Context, req *CancelFulfillmentRequest) error {
	var resp CommonResponse
	if err := c.sendPostRequest(ctx, c.cancelFulfillmentURL, req, &resp, "Không thể hủy đơn giao hàng"); err != nil {
		return err
	}
	return nil
}

func (c *Client) sendPostRequest(ctx context.Context, path string, req interface{}, resp ResponseInterface, msg string) error {
	res, err := c.rclient.R().
		SetBody(req).
		SetHeader("Authorization", "Bearer "+c.token).
		SetHeader("Ref-Account-ID", c.affiliateID).
		Post(path)
	if err != nil {
		return cm.Errorf(cm.Internal, err, "Lỗi kết nối với %v", c.connName)
	}
	if err := httpreq.HandleResponse(ctx, res, resp, msg); err != nil {
		return cm.Errorf(cm.ExternalServiceError, err, "Lỗi từ %v: %v", c.connName, err.Error())
	}

	cr := resp.GetCommonResponse()
	if !cr.Success {
		return cm.Errorf(cm.ExternalServiceError, nil, "Lỗi từ %v: %v (%v)", c.connName, msg, cr.Message)
	}
	return nil
}

func (c *Client) SignIn(ctx context.Context, req *SignInRequest) (*AccountData, error) {
	if c.signInURL == "" {
		return nil, cm.Errorf(cm.ExternalServiceError, nil, "Nhà vận chuyển không hỗ trợ đăng nhập tài khoản")
	}
	var resp SignInResponse
	if err := c.sendPostRequest(ctx, c.signInURL, req, &resp, "Không thể đăng nhập tài khoản"); err != nil {
		return nil, err
	}
	if err := resp.Data.validate(); err != nil {
		return nil, cm.Errorf(cm.ExternalServiceError, err, "Thông tin tài khoản không hợp lệ: %v", err.Error())
	}
	return resp.Data, nil
}

func (c *Client) SignUp(ctx context.Context, req *SignUpRequest) (*AccountData, error) {
	if c.signInURL == "" {
		return nil, cm.Errorf(cm.ExternalServiceError, nil, "Nhà vận chuyển không hỗ trợ đăng ký tài khoản")
	}
	var resp SignUpResponse
	if err := c.sendPostRequest(ctx, c.signUpURL, req, &resp, "Không thể tạo tài khoản"); err != nil {
		return nil, err
	}
	if err := resp.Data.validate(); err != nil {
		return nil, cm.Errorf(cm.ExternalServiceError, err, "Thông tin tài khoản không hợp lệ: %v", err.Error())
	}
	return resp.Data, nil
}
