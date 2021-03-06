package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/schema"

	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/httpreq"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/cmenv"
	"o.o/backend/pkg/common/validate"
	"o.o/backend/pkg/etop/model"
	"o.o/common/jsonx"
	"o.o/common/l"
)

var ll = l.New()
var encoder = schema.NewEncoder()

func init() {
	encoder.SetAliasTag("url")
}

type Client struct {
	baseUrl string
	token   string

	affiliateID string
	b2ctoken    string
	rclient     *httpreq.Resty
}

const (
	PathCalcShippingFee = "/services/shipment/fee"
	PathCreateOrder     = "/services/shipment/order"
	PathGetOrder        = "/services/shipment/v2"
	PathCancelOrder     = "/services/shipment/cancel"
	PathSignIn          = "/services/shops/token"
	PathSignUp          = "/services/shops/add"
)

func New(env string, cfg GhtkAccount) *Client {
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
		b2ctoken:    cfg.B2CToken,
		rclient:     httpreq.NewResty(rcfg),
	}
	switch env {
	case cmenv.PartnerEnvTest, cmenv.PartnerEnvDev:
		c.baseUrl = "https://dev.ghtk.vn"
	case cmenv.PartnerEnvProd:
		c.baseUrl = "https://services.giaohangtietkiem.vn"
	default:
		ll.Fatal("ghtk: Invalid ENV")
	}
	return c
}

func (c *Client) Ping() error {
	req := &CalcShippingFeeRequest{
		Weight:          1000,
		Value:           3000000,
		PickingProvince: "Hà Nội",
		PickingDistrict: "Quận Hai Bà Trưng",
		Address:         "163 Tô Hiến Thành",
		Province:        "Hồ Chí Minh",
		District:        "Quận 10",
		Transport:       "road",
	}
	_, err := c.CalcShippingFee(context.Background(), req)
	return err
}

func (c *Client) GetAffiliateID() string {
	return c.affiliateID
}

func (c *Client) TestCreateOrder() error {
	req := &CreateOrderRequest{
		Products: []*ProductRequest{
			{
				Name:     "bút",
				Weight:   0.1,
				Quantity: 1,
			},
			{
				Name:     "tẩy",
				Weight:   0.2,
				Quantity: 1,
			},
		},
		Order: &OrderRequest{
			ID:           "1234.5678",
			PickName:     "HCM-nội thành",
			PickAddress:  "590 CMT8 P.11",
			PickProvince: "TP. Hồ Chí Minh",
			PickDistrict: "Quận 3",
			PickTel:      "0911222333",
			Tel:          "0911222333",
			Name:         "GHTK - HCM - Noi Thanh",
			Address:      "123 nguyễn chí thanh",
			Province:     "TP. Hồ Chí Minh",
			District:     "Quận 1",
			IsFreeship:   1,
			PickMoney:    47000,
			Note:         "Khối lượng tính cước tối đa: 1.00 kg",
			Value:        3000000,
		},
	}
	_, err := c.CreateOrder(context.Background(), req)
	return err
}

func (c *Client) TestGetOrder() error {
	labelID := "S1858017.SG11.10I.236273443"
	_, err := c.GetOrder(context.Background(), labelID, "")
	return err
}

func (c *Client) TestCancelOrder() error {
	labelID := "S1858017.298467358"
	_, err := c.CancelOrder(context.Background(), labelID, "")
	return err
}

func (c *Client) CalcShippingFee(ctx context.Context, req *CalcShippingFeeRequest) (*CalcShippingFeeResponse, error) {
	req.PickingProvince = validate.NormalizeSearchSimple(req.PickingProvince)
	req.PickingDistrict = validate.NormalizeSearchSimple(req.PickingDistrict)
	req.PickingWard = validate.NormalizeSearchSimple(req.PickingWard)
	req.PickingAddress = validate.NormalizeSearchSimple(req.PickingAddress)
	req.Province = validate.NormalizeSearchSimple(req.Province)
	req.District = validate.NormalizeSearchSimple(req.District)
	req.Ward = validate.NormalizeSearchSimple(req.Ward)
	req.Address = validate.NormalizeSearchSimple(req.Address)
	var resp CalcShippingFeeResponse
	err := c.sendGetRequest(ctx, PathCalcShippingFee, req, &resp,
		"Không thể tính phí giao hàng")
	return &resp, err
}

func (c *Client) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*CreateOrderResponse, error) {
	// Setup freeship: COD sẽ chỉ thu người nhận hàng số tiền bằng pick_money
	req.Order.IsFreeship = 1
	req.Order.Hamlet = "Khác"
	if err := req.Validate(); err != nil {
		return nil, err
	}
	type Resp struct {
		Err           error
		OrderResponse *CreateOrderResponse
	}
	var ch = make(chan Resp, 1)
	go func() {
		var resp CreateOrderResponse
		err := c.sendPostRequest(ctx, PathCreateOrder, req, &resp, "Không thể tạo đơn hàng")
		ch <- Resp{
			Err:           err,
			OrderResponse: &resp,
		}
	}()
	tick := time.NewTicker(10 * time.Second)
	defer tick.Stop()
	go func() {
		for _ = range tick.C {
			newCtx, ctxCancel := context.WithTimeout(ctx, 10*time.Second)
			defer ctxCancel()
			orderResp, err := c.GetOrder(newCtx, "", req.Order.ID)
			if err == nil && orderResp != nil {
				order := orderResp.Order
				resp := &CreateOrderResponse{
					CommonResponse: orderResp.CommonResponse,
					Order: OrderResponse{
						PartnerID:            order.PartnerID,
						Label:                order.LabelID,
						Fee:                  order.ShipMoney,
						InsuranceFee:         order.Insurance,
						EstimatedPickTime:    order.PickDate,
						EstimatedDeliverTime: order.DeliverDate,
						StatusID:             order.Status,
					},
				}
				ch <- Resp{
					Err:           nil,
					OrderResponse: resp,
				}
			}
		}
	}()
	resp := <-ch
	return resp.OrderResponse, resp.Err
}

func (c *Client) GetOrder(ctx context.Context, labelID, orderPartnerID string) (*GetOrderResponse, error) {
	endPoint := PathGetOrder
	if labelID != "" {
		endPoint += "/" + labelID
	} else if orderPartnerID != "" {
		endPoint += "/partner_id:" + orderPartnerID
	}
	var resp GetOrderResponse
	err := c.sendGetRequest(ctx, endPoint, nil, &resp, "Không thể lấy thông tin đơn hàng")
	return &resp, err
}

func (c *Client) CancelOrder(ctx context.Context, labelID, orderPartnerID string) (*CommonResponse, error) {
	endPoint := PathCancelOrder
	if labelID != "" {
		endPoint += "/" + labelID
	} else if orderPartnerID != "" {
		endPoint += "/" + orderPartnerID
	}
	var resp CommonResponse
	err := c.sendPostRequest(ctx, endPoint, nil, &resp, "Không thể hủy đơn hàng")
	return &resp, err
}

func (c *Client) SignIn(ctx context.Context, req *SignInRequest) (*SignInResponse, error) {
	var resp SignInResponse
	err := c.sendPostRequest(ctx, PathSignIn, req, &resp, "Không thể đăng nhập tài khoản")
	return &resp, err
}

func (c *Client) SignUp(ctx context.Context, req *SignUpRequest) (*SignUpResponse, error) {
	var resp SignUpResponse
	err := c.sendPostRequest(ctx, PathSignUp, req, &resp, "Không thể tạo tài khoản mới")
	return &resp, err
}

func (c *Client) sendGetRequest(ctx context.Context, path string, req interface{}, resp ResponseInterface, msg string) error {
	queryString := url.Values{}
	if req != nil {
		err := encoder.Encode(req, queryString)
		if err != nil {
			return cm.Error(cm.Internal, "", err)
		}
	}

	res, err := c.rclient.R().
		SetQueryString(queryString.Encode()).
		SetHeader("token", c.token).
		Get(model.URL(c.baseUrl, path))
	if err != nil {
		return cm.Error(cm.ExternalServiceError, "Lỗi kết nối với GHTK", err)
	}
	err = handleResponse(ctx, res, resp, msg)
	return err
}

func (c *Client) sendPostRequest(ctx context.Context, path string, req interface{}, resp ResponseInterface, msg string) error {
	_req := c.rclient.R().
		SetBody(req).
		SetHeader("token", c.token)
	if c.b2ctoken != "" {
		_req.SetHeader("X-Refer-Token", c.b2ctoken)
	}

	res, err := _req.Post(model.URL(c.baseUrl, path))
	if err != nil {
		return cm.Error(cm.ExternalServiceError, "Lỗi kết nối với GHTK", err)
	}
	err = handleResponse(ctx, res, resp, msg)
	return err
}

func handleResponse(ctx context.Context, res *httpreq.RestyResponse, result ResponseInterface, msg string) error {
	status := res.StatusCode()
	var err error
	body := res.Body()
	switch {
	case status >= 200 && status < 300:
		if result != nil {
			if httpreq.IsNullJsonRaw(body) {
				return cm.Errorf(cm.ExternalServiceError, nil, "Lỗi không xác định từ Giaohangtietkiem: null response. Chúng tôi đang liên hệ với Giaohangtietkiem để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ %v.", wl.X(ctx).CSEmail)
			}
			if err = jsonx.Unmarshal(body, result); err != nil {
				return cm.Errorf(cm.ExternalServiceError, err, "Lỗi không xác định từ Giaohangtietkiem: %v. Chúng tôi đang liên hệ với Giaohangtietkiem để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ %v.", err, wl.X(ctx).CSEmail)
			}
			cr := result.GetCommonResponse()
			if !cr.Success {
				return cm.Errorf(cm.ExternalServiceError, nil, "Lỗi từ Giaohangtietkiem: %v (%v)", msg, cr.Message).
					WithMeta("messsage", cr.Message)
			}
		}
		return nil

	case status >= 400:
		var meta map[string]string
		if !httpreq.IsNullJsonRaw(body) {
			if err = jsonx.Unmarshal(body, &meta); err != nil {
				// The slow path
				var metaX map[string]interface{}
				_ = jsonx.Unmarshal(body, &metaX)
				meta = make(map[string]string)
				for k, v := range metaX {
					meta[k] = fmt.Sprint(v)
				}
			}
		}

		return cm.Errorf(cm.ExternalServiceError, nil, "Lỗi từ Giaohangtietkiem. Nếu cần thêm thông tin vui lòng liên hệ %v.", wl.X(ctx).CSEmail).WithMetaM(meta)
	default:
		return cm.Errorf(cm.ExternalServiceError, nil, "Lỗi không xác định từ Giaohangtietkiem: Invalid status (%v). Chúng tôi đang liên hệ với Giaohangtietkiem để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ %v.", status, wl.X(ctx).CSEmail)
	}
}
