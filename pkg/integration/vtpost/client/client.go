package vtpostclient

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/schema"
	resty "gopkg.in/resty.v1"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/httpreq"
	"etop.vn/backend/pkg/common/l"
	"etop.vn/backend/pkg/etop/model"
)

type ConfigAccount struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type ClientStates struct {
	AccessToken          string
	ExpiresAt            time.Time
	AccessTokenCreatedAt time.Time

	CustomerID     int
	GroupAddressID int
}

type Client interface {
	InitFromSavedStates(states ClientStates)
	GetStatesForSerialization() ClientStates
	LoginAndRefreshToken(ctx context.Context) error
	AutoLoginAndRefreshToken(ctx context.Context) (bool, error)

	Ping() error
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	GetWarehouses(ctx context.Context) (*WarehouseResponse, error)
	CalcShippingFee(ctx context.Context, req *CalcShippingFeeRequest) (*CalcShippingFeeResponse, error)
	CalcShippingFeeAllServices(ctx context.Context, req *CalcShippingFeeAllServicesRequest) ([]*ShippingFeeService, error)
	GetProvinces(ctx context.Context) (*GetProvinceResponse, error)
	GetDistricts(ctx context.Context) (*GetDistrictResponse, error)
	GetDistrictsByProvince(ctx context.Context, req *GetDistrictsByProvinceRequest) (*GetDistrictResponse, error)
	GetWardsByDistrict(ctx context.Context, req *GetWardsByDistrictRequest) (*GetWardsResponse, error)
	CreateOrder(ctx context.Context, req *CreateOrderRequest) (*CreateOrderResponse, error)
	CancelOrder(ctx context.Context, req *CancelOrderRequest) (*CommonResponse, error)
}

type ClientImpl struct {
	baseUrl string
	rclient *resty.Client

	Username string
	Password string

	ClientStates
}

var _ Client = &ClientImpl{}

var ll = l.New()
var encoder = schema.NewEncoder()

func init() {
	encoder.SetAliasTag("url")
}

const (
	PathLogin                  = "/user/Login"
	PathListInventory          = "/user/listInventory"
	PathGetPrice               = "/order/getPrice"
	PathGetPriceAll            = "/order/getPriceAll"
	PathGetProvinces           = "/categories/listProvince"
	PathGetDistricts           = "/categories/listDistrict"
	PathGetDistrictsByProvince = "/categories/listDistrict?provinceId=-1"
	PathGetWardsByDistrict     = "/categories/listWards"
	PathCreateOrder            = "/order/createOrder"
	PathUpdateOrder            = "/order/UpdateOrder"
)

func New(env string, cfg ConfigAccount) *ClientImpl {
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	c := &ClientImpl{
		rclient:  resty.NewWithClient(client).SetDebug(true),
		Username: cfg.Username,
		Password: cfg.Password,
	}
	c.baseUrl = "https://partner.viettelpost.vn/v2/"
	return c
}

func (c *ClientImpl) Ping() error {
	ctx := context.Background()
	req := &CalcShippingFeeRequest{
		SenderProvince:   2,
		SenderDistrict:   35,
		ReceiverProvince: 1,
		ReceiverDistrict: 9,
		ProductType:      "HH",
		OrderService:     OrderServiceCodeSCOD,
		OrderServiceAdd:  "GBH",
		ProductWeight:    2400,
		ProductPrice:     500000,
		MoneyCollection:  500000,
		ProductQuantity:  0,
		NATIONAL_TYPE:    1,
	}

	_, err := c.CalcShippingFee(ctx, req)
	return err
}

func (c *ClientImpl) InitFromSavedStates(states ClientStates) {
	c.ClientStates = states
}

func (c *ClientImpl) GetStatesForSerialization() ClientStates {
	return c.ClientStates
}

func (c *ClientImpl) LoginAndRefreshToken(ctx context.Context) error {
	req := &LoginRequest{
		Username: c.Username,
		Password: c.Password,
	}
	resp, err := c.Login(ctx, req)
	if err != nil {
		return err
	}

	token := resp.Data.Token
	expiresAt := cm.GetJWTExpires(token)
	if expiresAt.IsZero() || expiresAt.Before(time.Now()) {
		// according to ViettelPost documentation, the token expires after 1 month
		expiresAt = time.Now().AddDate(0, 1, 0)
	}

	c.AccessToken = token
	c.ExpiresAt = expiresAt
	c.AccessTokenCreatedAt = time.Now()
	return c.Ping()
}

func (c *ClientImpl) AutoLoginAndRefreshToken(ctx context.Context) (bool, error) {
	if c.ExpiresAt.Sub(time.Now()) < 30*time.Minute {
		return true, c.LoginAndRefreshToken(ctx)
	}
	return false, nil
}

func (c *ClientImpl) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	var resp LoginResponse
	err := c.sendPostRequest(ctx, PathLogin, req, &resp, "Không thể đăng nhập")
	return &resp, err
}

func (c *ClientImpl) GetWarehouses(ctx context.Context) (*WarehouseResponse, error) {
	var resp WarehouseResponse
	err := c.sendGetRequest(ctx, PathListInventory, nil, &resp, "Không thể lấy danh sách kho")
	return &resp, err
}

func (c *ClientImpl) CalcShippingFee(ctx context.Context, req *CalcShippingFeeRequest) (*CalcShippingFeeResponse, error) {
	var _resp _CalcShippingFeeResponse
	// Hard code
	req.ProductType = ProductTypeHH // Hàng Hóa
	req.NATIONAL_TYPE = 1           // Trọng nước
	if req.ProductPrice > 0 {
		// tính phí bảo hiểm
		req.OrderServiceAdd = OrderServiceCodeInsurance
	}
	if err := c.sendPostRequest(ctx, PathGetPrice, req, &_resp,
		"Không thể tính phí giao hàng"); err != nil {
		return nil, err
	}
	var data ShippingFeeData
	if err := json.Unmarshal(_resp.Data, &data); err != nil {

		return nil, cm.Errorf(cm.ExternalServiceError, nil, "Lỗi không xác định từ viettelpost: %v. Chúng tôi đang liên hệ với viettelpost để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn.", err)

	}

	return &CalcShippingFeeResponse{
		CommonResponse: _resp.CommonResponse,
		Data:           &data,
	}, nil
}

func (c *ClientImpl) CalcShippingFeeAllServices(ctx context.Context, req *CalcShippingFeeAllServicesRequest) ([]*ShippingFeeService, error) {
	// Phí trả về là phí căn bản, không bao gồm phí bảo hiểm
	// Cách tính phí bảo hiểm
	// BẢO HIỂM (Mã: GBH)
	// 1% Giá trị khai giá
	// Tối thiểu 15.000VNĐ/bưu gửi.
	var resp []*ShippingFeeService
	req.ProductType = "HH"
	req.Type = 1
	err := c.sendPostRequest(ctx, PathGetPriceAll, req, &resp, "Không thể tính phí giao hàng")
	return resp, err
}

func (c *ClientImpl) GetProvinces(ctx context.Context) (*GetProvinceResponse, error) {
	var resp GetProvinceResponse
	err := c.sendGetRequest(ctx, PathGetProvinces, nil, &resp, "Không thể lấy danh sách tỉnh thành")
	return &resp, err
}

func (c *ClientImpl) GetDistricts(ctx context.Context) (*GetDistrictResponse, error) {
	var resp GetDistrictResponse
	err := c.sendGetRequest(ctx, PathGetDistricts, nil, &resp, "Không thể lấy danh sách quận huyện")
	return &resp, err
}

func (c *ClientImpl) GetDistrictsByProvince(ctx context.Context, req *GetDistrictsByProvinceRequest) (*GetDistrictResponse, error) {
	var resp GetDistrictResponse
	if req.ProvinceID == "" {
		return nil, cm.Error(cm.InvalidArgument, "Missing Province ID", nil)
	}
	err := c.sendGetRequest(ctx, PathGetDistrictsByProvince, req, &resp, "Không thể lấy danh sách quận huyện")
	return &resp, err
}

func (c *ClientImpl) GetWardsByDistrict(ctx context.Context, req *GetWardsByDistrictRequest) (*GetWardsResponse, error) {
	var resp GetWardsResponse
	if req.DistrictID == 0 {
		return nil, cm.Error(cm.InvalidArgument, "Missing District ID", nil)
	}
	err := c.sendGetRequest(ctx, PathGetWardsByDistrict, req, &resp, "Không thể lấy danh sách phường xã")
	return &resp, err
}

func (c *ClientImpl) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*CreateOrderResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	// set default data
	// req.GroupAddressID = c.GroupAddressID
	req.CusID = c.CustomerID
	req.ProductType = ProductTypeHH
	// 3: Thu hộ tiền hàng/ Collect price of goods
	req.Orderpayment = 3

	if req.ProductPrice > 0 {
		// tính phí bảo hiểm
		req.OrderServiceAdd = OrderServiceCodeInsurance
	}

	var resp CreateOrderResponse
	err := c.sendPostRequest(ctx, PathCreateOrder, req, &resp, "Không thể tạo đơn hàng")
	return &resp, err
}

func (c *ClientImpl) CancelOrder(ctx context.Context, req *CancelOrderRequest) (*CommonResponse, error) {
	endPoint := PathUpdateOrder
	var resp CommonResponse
	// Loại trạng thái:/ Status type
	// 1. Duyệt đơn hàng/ Confirm order
	// 2. Duyệt chuyển hoàn/ Confirm return shipping
	// 3. Phát tiếp/ delivery again
	// 4. Hủy đơn hàng/ delivery again
	// 5. Lấy lại đơn hàng (Gửi lại)/ get back order (re-order)
	// 11. Xóa đơn hàng đã hủy(delete canceled order)
	req.Type = 4
	err := c.sendPostRequest(ctx, endPoint, req, &resp, "Không thể hủy đơn hàng")
	return &resp, err
}

func (c *ClientImpl) sendGetRequest(ctx context.Context, path string, params interface{}, resp ResponseInterface, msg string) error {
	queryString := url.Values{}
	if params != nil {
		err := encoder.Encode(params, queryString)
		if err != nil {
			return cm.Error(cm.Internal, "", err)
		}
	}

	res, err := c.rclient.R().
		SetQueryString(queryString.Encode()).
		SetHeader("token", c.AccessToken).
		Get(model.URL(c.baseUrl, path))
	if err != nil {
		return cm.Error(cm.ExternalServiceError, "Lỗi kết nối với VTPost", err)
	}
	err = handleResponse(res, resp, msg)
	return err
}

func (c *ClientImpl) sendPostRequest(ctx context.Context, path string, body interface{}, resp interface{}, msg string) error {
	res, err := c.rclient.R().
		SetBody(body).
		SetHeader("Token", c.AccessToken).
		Post(model.URL(c.baseUrl, path))
	if err != nil {
		return cm.Errorf(cm.ExternalServiceError, err, "Lỗi kết nối với VTPost: %v (%v)", msg, err)
	}
	err = handleResponse(res, resp, msg)
	return err
}

func handleResponse(res *resty.Response, result interface{}, msg string) error {
	status := res.StatusCode()
	var err error
	body := res.Body()
	type ResponseInterface struct {
		CommonResponse
	}

	switch {
	case status >= 200 && status < 300:
		if result != nil {
			if httpreq.IsNullJsonRaw(body) {
				return cm.Error(cm.ExternalServiceError, "Lỗi không xác định từ viettelpost: null response. Chúng tôi đang liên hệ với viettelpost để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn.", nil)
			}
			if err = json.Unmarshal(body, result); err != nil {
				return cm.Errorf(cm.ExternalServiceError, err, "Lỗi không xác định từ viettelpost: %v. Chúng tôi đang liên hệ với viettelpost để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn.", err)
			}
			var responseFormat ResponseInterface
			if err := json.Unmarshal(body, &responseFormat); err == nil {
				if responseFormat.Error {
					return cm.Errorf(cm.ExternalServiceError, nil, "Lỗi từ viettelpost: %v (%v)", msg, responseFormat.Message)
				}
			}

			if err = json.Unmarshal(body, result); err != nil {
				return cm.Errorf(cm.ExternalServiceError, err, "Lỗi không xác định từ viettelpost: %v. Chúng tôi đang liên hệ với viettelpost để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn.", err)
			}
		}
		return nil

	case status >= 400:
		var meta map[string]string
		if !httpreq.IsNullJsonRaw(body) {
			if err = json.Unmarshal(body, &meta); err != nil {
				// The slow path
				var metaX map[string]interface{}
				_ = json.Unmarshal(body, &metaX)
				meta = make(map[string]string)
				for k, v := range metaX {
					meta[k] = fmt.Sprint(v)
				}
			}
		}

		return cm.Errorf(cm.ExternalServiceError, nil, "Lỗi từ viettelpost. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn.").WithMetaM(meta)
	default:
		return cm.Errorf(cm.ExternalServiceError, nil, "Lỗi không xác định từ viettelpost: Invalid status (%v). Chúng tôi đang liên hệ với viettelpost để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn.", status)
	}
}
