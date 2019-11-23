package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/schema"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/httpreq"
	"etop.vn/common/jsonx"
	"etop.vn/common/l"
	"etop.vn/common/xerrors"
)

var ll = l.New()
var encoder = schema.NewEncoder()

func init() {
	encoder.SetAliasTag("url")
}

type Client struct {
	baseUrl          string
	verifyAccountUrl string
	apiKey           string
	rclient          *httpreq.Resty
	env              string
}

const (
	PathCalcShippingFee = "/order/estimated_fee"
	PathCreateOrder     = "/order/create"
	PathGetOrder        = "/order/detail"
	PathCancelOrder     = "/order/cancel"
	PathRegisterAccount = "/partner/register_account"
	PathGetAccount      = "/user/profile"
	PathGetServices     = "/order/service_types"
)

func New(cfg Config) *Client {
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	rcfg := httpreq.RestyConfig{Client: client}
	c := &Client{
		apiKey:  cfg.ApiKey,
		rclient: httpreq.NewResty(rcfg),
		env:     cfg.Env,
	}
	switch cfg.Env {
	case cm.PartnerEnvTest:
		c.baseUrl = "http://apistg.ahamove.com/v1/"
		c.verifyAccountUrl = "https://ws.ahamove.com/partner/create_ticket_stg"
	case cm.PartnerEnvProd:
		c.baseUrl = "https://api.ahamove.com/v1/"
		c.verifyAccountUrl = "https://ws.ahamove.com/partner/create_ticket"
	default:
		ll.Fatal("ahamove: Invalid ENV")
	}
	return c
}

func (c *Client) Ping() error {
	return c.TestGetServices()
}

var points = []*DeliveryPointRequest{
	{
		Address: "384 Võ Văn Ngân, Bình Thọ, Thủ Đức, Hồ Chí Minh",
		Lat:     10.849550,
		Lng:     106.772640,
		Mobile:  "0999999999",
		Name:    "Tuan pn",
		Remarks: "Gọi shop nếu ko lấy được hàng",
	}, {
		Address: "163 Tô Hiến Thành, P15, Q10, Hồ Chí Minh",
		Lat:     10.779970,
		Lng:     106.669170,
	}, {
		Address: "67 Lê Lợi, Phường Bến Nghé, Quận 1, Tp. Hồ Chí Minh",
		Lat:     10.774960,
		Lng:     106.701690,
	}, {
		Address: "263 Tô Hiến Thành, P15, Q10, Hồ Chí Minh",
		Lat:     10.778680,
		Lng:     106.666610,
	}, {
		Address: "167 Lê Lợi, Phường Bến Nghé, Quận 1, Tp. Hồ Chí Minh",
		Lat:     10.774960,
		Lng:     106.701690,
	}, {
		Address: "363 Tô Hiến Thành, P15, Q10, Hồ Chí Minh",
		Lat:     10.778680,
		Lng:     106.663830,
	}, {
		Address: "267 Lê Lợi, Phường Bến Nghé, Quận 1, Tp. Hồ Chí Minh",
		Lat:     10.774960,
		Lng:     106.701690,
	},
}

func (c *Client) TestCalcShippingFee() error {
	req := &CalcShippingFeeRequest{
		OrderTime:      0,
		IdleUntil:      0,
		DeliveryPoints: points,
		ServiceID:      "SGN-BIKE",
	}
	_, err := c.CalcShippingFee(context.Background(), req)
	return err
}

func (c *Client) TestGetOrder() error {
	req := &GetOrderRequest{
		OrderID: "1JFU54",
	}
	_, err := c.GetOrder(context.Background(), req)
	return err
}

func (c *Client) TestCancelOrder() error {
	req := &CancelOrderRequest{
		OrderId: "1JFU54",
		Comment: "test cancel",
	}
	err := c.CancelOrder(context.Background(), req)
	return err
}

func (c *Client) TestCreateOrder() error {
	req := &CreateOrderRequest{
		OrderTime:      0,
		DeliveryPoints: points,
		ServiceID:      "SGN-BIKE",
		Remarks:        "test order",
		PaymentMethod:  "CASH",
	}
	_, err := c.CreateOrder(context.Background(), req)
	return err
}

func (c *Client) TestGetServices() error {
	req := &GetServicesRequest{
		CityID: "SGN",
	}
	_, err := c.GetServices(context.Background(), req)
	return err
}

func (c *Client) CalcShippingFee(ctx context.Context, req *CalcShippingFeeRequest) (*CalcShippingFeeResponse, error) {
	req.Path = ConvertDeliveryPointsRequestToString(req.DeliveryPoints)
	req.PaymentMethod = PaymentMethodCash
	var resp CalcShippingFeeResponse
	err := c.sendGetRequest(ctx, PathCalcShippingFee, req, &resp,
		"Không thể tính phí giao hàng")
	return &resp, err
}

func (c *Client) GetServices(ctx context.Context, req *GetServicesRequest) ([]*ServiceType, error) {
	var resp []*ServiceType
	err := c.sendGetRequest(ctx, PathGetServices, req, &resp, "Không thể lấy danh sách dịch vụ")
	return resp, err
}

func (c *Client) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*CreateOrderResponse, error) {
	req.Path = ConvertDeliveryPointsRequestToString(req.DeliveryPoints)
	req.PaymentMethod = PaymentMethodCash
	var resp CreateOrderResponse
	err := c.sendGetRequest(ctx, PathCreateOrder, req, &resp, "Không thể tạo đơn hàng")
	return &resp, err
}

func (c *Client) GetOrder(ctx context.Context, req *GetOrderRequest) (*Order, error) {
	var resp Order
	err := c.sendGetRequest(ctx, PathGetOrder, req, &resp, "Không thể lấy thông tin đơn hàng")
	return &resp, err
}

func (c *Client) CancelOrder(ctx context.Context, req *CancelOrderRequest) error {
	err := c.sendGetRequest(ctx, PathCancelOrder, req, nil, "Không thể hủy đơn hàng")
	return err
}

func (c *Client) RegisterAccount(ctx context.Context, req *RegisterAccountRequest) (*RegisterAccountResponse, error) {
	req.ApiKey = c.apiKey
	var resp RegisterAccountResponse
	err := c.sendGetRequest(ctx, PathRegisterAccount, req, &resp, "Không thể tạo tài khoản Ahamove")
	return &resp, err
}

func (c *Client) VerifyAccount(ctx context.Context, req *VerifyAccountRequest) (*VerifyAccountResponse, error) {
	switch c.env {
	case "test":
		req.Description = "[TEST-TOPSHIP] " + req.Description
	case "prod":
		req.Description = "[TOPSHIP] " + req.Description
	}
	// default value
	req.Subject = "VERIFY USER - COD SERVICE"
	req.Type = "ahamove_verify_user"

	queryString := url.Values{}
	err := encoder.Encode(req, queryString)
	if err != nil {
		return nil, cm.Error(cm.Internal, "", err)
	}

	res, err := c.rclient.R().
		SetQueryString(queryString.Encode()).
		Get(c.verifyAccountUrl)
	if err != nil {
		return nil, cm.Error(cm.ExternalServiceError, "Lỗi kết nối với ahamove", err)
	}
	var resp VerifyAccountResponse
	err = handleResponse(res, &resp, "Không thể gửi yêu cầu xác thực tài khoản")
	return &resp, err
}

func (c *Client) GetAccount(ctx context.Context, req *GetAccountRequest) (*Account, error) {
	var resp Account
	err := c.sendGetRequest(ctx, PathGetAccount, req, &resp, "Không thể lấy thông tin tài khoản")
	return &resp, err
}

func (c *Client) sendGetRequest(ctx context.Context, path string, req interface{}, resp interface{}, msg string) error {
	queryString := url.Values{}
	if req != nil {
		err := encoder.Encode(req, queryString)
		if err != nil {
			return cm.Error(cm.Internal, "", err)
		}
	}
	res, err := c.rclient.R().
		SetQueryString(queryString.Encode()).
		Get(buildUrl(c.baseUrl, path))
	if err != nil {
		return cm.Error(cm.ExternalServiceError, "Lỗi kết nối với ahamove", err)
	}
	err = handleResponse(res, resp, msg)
	return err
}

func handleResponse(res *httpreq.RestyResponse, result interface{}, msg string) error {
	status := res.StatusCode()
	var err error
	body := res.Body()
	switch {
	case status >= 200 && status < 300:
		if result != nil {
			if httpreq.IsNullJsonRaw(body) {
				return cm.Error(cm.ExternalServiceError, "Lỗi không xác định từ ahamove: null response. Chúng tôi đang liên hệ với ahamove để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn.", nil)
			}
			if err = jsonx.Unmarshal(body, result); err != nil {
				return cm.Errorf(cm.ExternalServiceError, err, "Lỗi không xác định từ ahamove: %v. Chúng tôi đang liên hệ với ahamove để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn.", err)
			}
		}
		return nil

	case status >= 400:
		var meta map[string]string
		var errJSON xerrors.ErrorJSON
		if !httpreq.IsNullJsonRaw(body) {
			if err = jsonx.Unmarshal(body, &meta); err != nil {
				// The slow path
				var metaX map[string]interface{}
				_ = jsonx.Unmarshal(body, &metaX)
				meta = make(map[string]string)
				for k, v := range metaX {
					meta[k] = fmt.Sprint(v)
				}
				errJSON.Meta = meta
			}
		}

		return cm.Errorf(cm.ExternalServiceError, &errJSON, "Lỗi từ ahamove: %v. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn.", errJSON.Error()).WithMetaM(meta)
	default:
		return cm.Errorf(cm.ExternalServiceError, nil, "Lỗi không xác định từ ahamove: Invalid status (%v). Chúng tôi đang liên hệ với ahamove để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn.", status)
	}
}

func buildUrl(baseUrl, path string) string {
	return baseUrl + path
}
