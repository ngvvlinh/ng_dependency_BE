package client

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gopkg.in/resty.v1"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/httpreq"
	"etop.vn/common/l"
)

var ll = l.New()

type Client struct {
	clientID int

	baseUrl string
	token   string
	rclient *resty.Client
}

func New(env string, clientID int, token string) *Client {
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	c := &Client{
		token:    token,
		clientID: clientID,
		rclient:  resty.NewWithClient(client).SetDebug(true),
	}
	switch env {
	case cm.PartnerEnvTest:
		c.baseUrl = "https://console.ghn.vn/api/v1/apiv3/"
	case cm.PartnerEnvProd:
		c.baseUrl = "https://console.ghn.vn/api/v1/apiv3/"
	default:
		ll.Fatal("ghn: Invalid env")
	}
	return c
}

func (c *Client) ClientID() int {
	return c.clientID
}

func (c *Client) Ping() error {
	req := &FindAvailableServicesRequest{
		FromDistrictID: 1442, // Quận 1, HCM
		ToDistrictID:   1578, // Vĩnh Yên, Vĩnh Phúc
	}
	_, err := c.FindAvailableServices(context.Background(), req)
	return err
}

func (c *Client) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*CreateOrderResponse, error) {
	req.Token = c.token
	var resp CreateOrderResponse
	err := c.sendRequest(ctx, "CreateOrder", req, &resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrorMessage != "" {
		return nil, cm.Errorf(cm.ExternalServiceError, nil,
			"Lỗi từ Giao Hàng Nhanh: "+string(resp.ErrorMessage))
	}
	return &resp, nil
}

func (c *Client) FindAvailableServices(ctx context.Context, req *FindAvailableServicesRequest) (*FindAvailableServicesResponse, error) {
	req.Token = c.token
	var resp []*AvailableService
	err := c.sendRequest(ctx, "FindAvailableServices", req, &resp)
	if err != nil {
		return nil, err
	}
	return &FindAvailableServicesResponse{
		AvailableServices: resp,
	}, nil
}

func (c *Client) CalculateFee(ctx context.Context, req *CalculateFeeRequest) (*CalculateFeeResponse, error) {
	req.Token = c.token
	var resp CalculateFeeResponse
	err := c.sendRequest(ctx, "CalculateFee", req, &resp)
	return &resp, err
}

func (c *Client) GetOrderInfo(ctx context.Context, req *OrderCodeRequest) (*Order, error) {
	req.Token = c.token
	var resp Order
	err := c.sendRequest(ctx, "OrderInfo", req, &resp)
	return &resp, err
}

func (c *Client) GetOrderLogs(ctx context.Context, req *OrderLogsRequest) (*OrderLogsResponse, error) {
	req.Token = c.token
	var resp OrderLogsResponse
	if req.FromTime == 0 {
		req.FromTime = 1
	}
	if req.ToTime == 0 {
		now := time.Now().UnixNano() / int64(time.Millisecond)
		req.ToTime = now
	}
	if req.Condition == nil {
		req.Condition = &OrderLogsCondition{}
	}
	req.Condition.CustomerID = c.clientID
	err := c.sendRequest(ctx, "GetOrderLogs", req, &resp)
	return &resp, err
}

func (c *Client) CancelOrder(ctx context.Context, req *OrderCodeRequest) error {
	req.Token = c.token
	return c.sendRequest(ctx, "CancelOrder", req, nil)
}

func (c *Client) ReturnOrder(ctx context.Context, req *OrderCodeRequest) error {
	req.Token = c.token
	return c.sendRequest(ctx, "ReturnOrder", req, nil)
}

func (c *Client) sendRequest(ctx context.Context, path string, req, resp interface{}) error {
	var errResp ErrorResponse
	res, err := c.rclient.R().
		SetBody(req).
		SetResult(&errResp).
		SetError(&errResp).
		Post(URL(c.baseUrl, path))
	if err != nil {
		return cm.Error(cm.ExternalServiceError, "Lỗi kết nối với GHN", err)
	}

	status := res.StatusCode()
	switch {
	case status >= 200 && status < 300:
		if resp != nil {
			if httpreq.IsNullJsonRaw(errResp.Data) {
				return cm.Error(cm.ExternalServiceError, "Lỗi không xác định từ Giao Hang Nhanh: null response. Chúng tôi đang liên hệ với Giao Hang Nhanh để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn.", nil)
			}
			if err = json.Unmarshal(errResp.Data, resp); err != nil {
				return cm.Errorf(cm.ExternalServiceError, err, "Lỗi không xác định từ Giao Hang Nhanh: %v. Chúng tôi đang liên hệ với Giao Hang Nhanh để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn.", err)
			}
		}
		return nil

	case status >= 400:
		var meta map[string]string
		if !httpreq.IsNullJsonRaw(errResp.Data) {
			if err = json.Unmarshal(errResp.Data, &meta); err != nil {
				// The slow path
				var metaX map[string]interface{}
				_ = json.Unmarshal(errResp.Data, &metaX)
				meta = make(map[string]string)
				for k, v := range metaX {
					meta[k] = fmt.Sprint(v)
				}
			}
			errResp.ErrorData = meta
		}

		// Handle "An error occur"
		if strings.Contains(errResp.Msg.String(), "error") {
			return cm.Errorf(cm.ExternalServiceError, &errResp, "Lỗi từ Giao Hang Nhanh: %v. Chúng tôi đang liên hệ với Giao Hang Nhanh để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn.", errResp.Error()).WithMetaM(meta)
		}

		return cm.Errorf(cm.ExternalServiceError, &errResp, "Lỗi từ Giao Hang Nhanh: %v. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn.", errResp.Error()).WithMetaM(meta)

	default:
		return cm.Errorf(cm.ExternalServiceError, nil, "Lỗi không xác định từ Giao Hang Nhanh: Invalid status (%v). Chúng tôi đang liên hệ với Giao Hang Nhanh để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn.", status)
	}
}
