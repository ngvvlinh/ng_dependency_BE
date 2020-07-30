package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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
var encoder = schema.NewEncoder()

type Client struct {
	clientID    int
	affiliateID int

	baseUrl string
	token   string
	rclient *httpreq.Resty
}

func New(env string, cfg GHNAccountCfg) *Client {
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	rcfg := httpreq.RestyConfig{Client: client}
	c := &Client{
		token:       cfg.Token,
		affiliateID: cfg.AffiliateID,
		clientID:    cfg.ClientID,
		rclient:     httpreq.NewResty(rcfg),
	}
	switch env {
	case cmenv.PartnerEnvTest, cmenv.PartnerEnvDev:
		c.baseUrl = "https://console.ghn.vn/api/v1/apiv3/"
	case cmenv.PartnerEnvProd:
		c.baseUrl = "https://console.ghn.vn/api/v1/apiv3/"
	default:
		ll.Fatal("ghn: Invalid env")
	}
	return c
}

func (c *Client) ClientID() int {
	return c.clientID
}

func (c *Client) GetAffiliateID() string {
	if c.affiliateID > 0 {
		return strconv.Itoa(c.affiliateID)
	}
	return ""
}

func (c *Client) GetToken() string {
	return c.token
}

func (c *Client) Ping() error {
	req := &FindAvailableServicesRequest{
		FromDistrictID: 1442, // Quận 1, HCM
		ToDistrictID:   1578, // Vĩnh Yên, Vĩnh Phúc
	}
	_, err := c.FindAvailableServices(context.Background(), req)
	return err
}

func (c *Client) CreateTicket(ctx context.Context, req *CreateTicketRequest) (*CreateTicketResponse, error) {
	req.CEmail = "etop@etop.vn" //TODO(Nam) config
	var resp CreateTicketResponse
	fullUrl := c.baseUrl + "shiip/public-api/ticket/create"
	if err := c.sendPostFormRequest(ctx, "post", fullUrl, req, &resp, "Không thể tạo ticket"); err != nil {
		return nil, err
	}
	return &resp, nil
}

// reply in ghn = comment etop
func (c *Client) CreateReply(ctx context.Context, req *CreateTicketReplyRequest) (*CreateTicketReplyResponse, error) {
	req.UserID = "2043053874563" //TODO(Nam) Config
	var resp CreateTicketReplyResponse
	fullUrl := c.baseUrl + "shiip/public-api/ticket/reply"
	if err := c.sendPostFormRequest(ctx, "put", fullUrl, req, &resp, "Không thể tạo ticket"); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*CreateOrderResponse, error) {
	req.Token = c.token
	req.AffiliateID = c.affiliateID

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

func (c *Client) SignIn(ctx context.Context, req *SignInRequest) (*SignInResponse, error) {
	req.Token = c.token
	var resp SignInResponse
	err := c.sendRequest(ctx, "SignIn", req, &resp)
	return &resp, err
}

func (c *Client) SignUp(ctx context.Context, req *SignUpRequest) (*SignInResponse, error) {
	req.Token = c.token
	var resp SignInResponse
	err := c.sendRequest(ctx, "SignUp", req, &resp)
	return &resp, err
}

func (c *Client) RegisterWebhookForClient(ctx context.Context, req *RegisterWebhookForClientRequest) error {
	req.Token = c.token
	// auto turn on all configs
	req.ConfigCOD = true
	req.ConfigReturnData = true
	req.ConfigField = WebhookConfigField{
		CODAmount:            true,
		CurrentWarehouseName: true,
		CustomerID:           true,
		CustomerName:         true,
		CustomerPhone:        true,
		Note:                 true,
		OrderCode:            true,
		ServiceName:          true,
		ShippingOrderCosts:   true,
		Weight:               true,
		ExternalCode:         true,
		ReturnInfo:           true,
	}
	req.ConfigStatus = WebhookConfigStatus{
		ReadyToPick:     true,
		Picking:         true,
		Storing:         true,
		Delivering:      true,
		Delivered:       true,
		WaitingToFinish: true,
		Return:          true,
		Returned:        true,
		Finish:          true,
		LostOrder:       true,
		Cancel:          true,
	}
	err := c.sendRequest(ctx, "SetConfigClient", req, nil)
	return err
}

func (c *Client) sendPostFormRequest(ctx context.Context, method string, fullURL string, req interface{}, resp interface{}, msg string) error {
	values := url.Values{}
	if req != nil {
		if err := encoder.Encode(req, values); err != nil {
			return err
		}
	}
	var formData = make(map[string]string)
	for key := range values {
		formData[key] = values.Get(key)
	}
	var res *resty.Response
	var err error
	switch method {
	case "post":
		res, err = c.rclient.R().SetFormData(formData).Post(fullURL)
	case "put":
		res, err = c.rclient.R().SetFormData(formData).Put(fullURL)
	default:
		return cm.Errorf(cm.ExternalServiceError, err, "Method not found")
	}
	if err != nil {
		return cm.Errorf(cm.ExternalServiceError, err, "Lỗi kết nối với GHN")
	}
	return httpreq.HandleResponse(ctx, res, resp, msg)
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
				return cm.Errorf(cm.ExternalServiceError, nil, "Lỗi không xác định từ Giao Hang Nhanh: null response. Chúng tôi đang liên hệ với Giao Hang Nhanh để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ %v.", wl.X(ctx).CSEmail)
			}
			if err = jsonx.Unmarshal(errResp.Data, resp); err != nil {
				return cm.Errorf(cm.ExternalServiceError, err, "Lỗi không xác định từ Giao Hang Nhanh: %v. Chúng tôi đang liên hệ với Giao Hang Nhanh để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ %v.", err, wl.X(ctx).CSEmail)
			}
		}
		return nil

	case status >= 400:
		var meta map[string]string
		if !httpreq.IsNullJsonRaw(errResp.Data) {
			if err = jsonx.Unmarshal(errResp.Data, &meta); err != nil {
				// The slow path
				var metaX map[string]interface{}
				_ = jsonx.Unmarshal(errResp.Data, &metaX)
				meta = make(map[string]string)
				for k, v := range metaX {
					meta[k] = fmt.Sprint(v)
				}
			}
			errResp.ErrorData = meta
		}

		// Handle "An error occur"
		if strings.Contains(errResp.Msg.String(), "error") {
			return cm.Errorf(cm.ExternalServiceError, &errResp, "Lỗi từ Giao Hang Nhanh: %v. Chúng tôi đang liên hệ với Giao Hang Nhanh để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ %.", errResp.Error(), wl.X(ctx).CSEmail).WithMetaM(meta)
		}

		return cm.Errorf(cm.ExternalServiceError, &errResp, "Lỗi từ Giao Hang Nhanh: %v. Nếu cần thêm thông tin vui lòng liên hệ %v.", errResp.Error(), wl.X(ctx).CSEmail).WithMetaM(meta)

	default:
		return cm.Errorf(cm.ExternalServiceError, nil, "Lỗi không xác định từ Giao Hang Nhanh: Invalid status (%v). Chúng tôi đang liên hệ với Giao Hang Nhanh để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ %v.", status, wl.X(ctx).CSEmail)
	}
}
