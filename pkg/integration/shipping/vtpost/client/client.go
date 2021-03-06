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
	"o.o/backend/pkg/etop/model"
	"o.o/common/jsonx"
	"o.o/common/l"
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
	GetUserName() string
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
	rclient *httpreq.Resty

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

	rcfg := httpreq.RestyConfig{Client: client}
	c := &ClientImpl{
		rclient:  httpreq.NewResty(rcfg),
		Username: cfg.Username,
		Password: cfg.Password,
	}
	c.baseUrl = "https://partner.viettelpost.vn/v2"
	return c
}

func NewClientWithToken(env string, token string) *ClientImpl {
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	rcfg := httpreq.RestyConfig{Client: client}
	c := &ClientImpl{
		rclient: httpreq.NewResty(rcfg),
		ClientStates: ClientStates{
			AccessToken: token,
		},
	}
	c.baseUrl = "https://partner.viettelpost.vn/v2"
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

func (c *ClientImpl) GetUserName() string {
	return c.Username
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
	c.CustomerID = resp.Data.UserId
	return c.Ping()
}

func (c *ClientImpl) AutoLoginAndRefreshToken(ctx context.Context) (bool, error) {
	if time.Until(c.ExpiresAt) < 30*time.Minute {
		return true, c.LoginAndRefreshToken(ctx)
	}
	return false, nil
}

func (c *ClientImpl) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	var resp LoginResponse
	err := c.sendPostRequest(ctx, PathLogin, req, &resp, "Kh??ng th??? ????ng nh???p")
	return &resp, err
}

func (c *ClientImpl) GetWarehouses(ctx context.Context) (*WarehouseResponse, error) {
	var resp WarehouseResponse
	err := c.sendGetRequest(ctx, PathListInventory, nil, &resp, "Kh??ng th??? l???y danh s??ch kho")
	return &resp, err
}

func (c *ClientImpl) CalcShippingFee(ctx context.Context, req *CalcShippingFeeRequest) (*CalcShippingFeeResponse, error) {
	var _resp _CalcShippingFeeResponse
	// Hard code
	req.ProductType = ProductTypeHH // H??ng H??a
	req.NATIONAL_TYPE = 1           // Tr???ng n?????c
	if req.ProductPrice > 0 {
		// t??nh ph?? b???o hi???m
		req.OrderServiceAdd = OrderServiceCodeInsurance
	}
	if err := c.sendPostRequest(ctx, PathGetPrice, req, &_resp,
		"Kh??ng th??? t??nh ph?? giao h??ng"); err != nil {
		return nil, err
	}
	var data ShippingFeeData
	if err := jsonx.Unmarshal(_resp.Data, &data); err != nil {

		return nil, cm.Errorf(cm.ExternalServiceError, nil, "L???i kh??ng x??c ?????nh t??? viettelpost: %v. Ch??ng t??i ??ang li??n h??? v???i viettelpost ????? x??? l??. Xin l???i qu?? kh??ch v?? s??? b???t ti???n n??y. N???u c???n th??m th??ng tin vui l??ng li??n h??? %v.", err, wl.X(ctx).CSEmail)

	}

	return &CalcShippingFeeResponse{
		CommonResponse: _resp.CommonResponse,
		Data:           &data,
	}, nil
}

func (c *ClientImpl) CalcShippingFeeAllServices(ctx context.Context, req *CalcShippingFeeAllServicesRequest) ([]*ShippingFeeService, error) {
	// Ph?? tr??? v??? l?? ph?? c??n b???n, kh??ng bao g???m ph?? b???o hi???m
	// C??ch t??nh ph?? b???o hi???m
	// B???O HI???M (M??: GBH)
	// 1% Gi?? tr??? khai gi??
	// T???i thi???u 15.000VN??/b??u g???i.
	var resp []*ShippingFeeService
	req.ProductType = "HH"
	req.Type = 1
	err := c.sendPostRequest(ctx, PathGetPriceAll, req, &resp, "Kh??ng th??? t??nh ph?? giao h??ng")
	return resp, err
}

func (c *ClientImpl) GetProvinces(ctx context.Context) (*GetProvinceResponse, error) {
	var resp GetProvinceResponse
	err := c.sendGetRequest(ctx, PathGetProvinces, nil, &resp, "Kh??ng th??? l???y danh s??ch t???nh th??nh")
	return &resp, err
}

func (c *ClientImpl) GetDistricts(ctx context.Context) (*GetDistrictResponse, error) {
	var resp GetDistrictResponse
	err := c.sendGetRequest(ctx, PathGetDistricts, nil, &resp, "Kh??ng th??? l???y danh s??ch qu???n huy???n")
	return &resp, err
}

func (c *ClientImpl) GetDistrictsByProvince(ctx context.Context, req *GetDistrictsByProvinceRequest) (*GetDistrictResponse, error) {
	var resp GetDistrictResponse
	if req.ProvinceID == "" {
		return nil, cm.Error(cm.InvalidArgument, "Missing Province ID", nil)
	}
	err := c.sendGetRequest(ctx, PathGetDistrictsByProvince, req, &resp, "Kh??ng th??? l???y danh s??ch qu???n huy???n")
	return &resp, err
}

func (c *ClientImpl) GetWardsByDistrict(ctx context.Context, req *GetWardsByDistrictRequest) (*GetWardsResponse, error) {
	var resp GetWardsResponse
	if req.DistrictID == 0 {
		return nil, cm.Error(cm.InvalidArgument, "Missing District ID", nil)
	}
	err := c.sendGetRequest(ctx, PathGetWardsByDistrict, req, &resp, "Kh??ng th??? l???y danh s??ch ph?????ng x??")
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
	// 3: Thu h??? ti???n h??ng/ Collect price of goods
	req.Orderpayment = 3

	if req.ProductPrice > 0 {
		// t??nh ph?? b???o hi???m
		req.OrderServiceAdd = OrderServiceCodeInsurance
	}

	var resp CreateOrderResponse
	err := c.sendPostRequest(ctx, PathCreateOrder, req, &resp, "Kh??ng th??? t???o ????n h??ng")
	return &resp, err
}

func (c *ClientImpl) CancelOrder(ctx context.Context, req *CancelOrderRequest) (*CommonResponse, error) {
	endPoint := PathUpdateOrder
	var resp CommonResponse
	// Lo???i tr???ng th??i:/ Status type
	// 1. Duy???t ????n h??ng/ Confirm order
	// 2. Duy???t chuy???n ho??n/ Confirm return shipping
	// 3. Ph??t ti???p/ delivery again
	// 4. H???y ????n h??ng/ delivery again
	// 5. L???y l???i ????n h??ng (G???i l???i)/ get back order (re-order)
	// 11. X??a ????n h??ng ???? h???y(delete canceled order)
	req.Type = 4
	err := c.sendPostRequest(ctx, endPoint, req, &resp, "Kh??ng th??? h???y ????n h??ng")
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
		return cm.Error(cm.ExternalServiceError, "L???i k???t n???i v???i VTPost", err)
	}
	err = handleResponse(ctx, res, resp, msg)
	return err
}

func (c *ClientImpl) sendPostRequest(ctx context.Context, path string, body interface{}, resp interface{}, msg string) error {
	res, err := c.rclient.R().
		SetBody(body).
		SetHeader("Token", c.AccessToken).
		Post(model.URL(c.baseUrl, path))
	if err != nil {
		return cm.Errorf(cm.ExternalServiceError, err, "L???i k???t n???i v???i VTPost: %v (%v)", msg, err)
	}
	err = handleResponse(ctx, res, resp, msg)
	return err
}

func handleResponse(ctx context.Context, res *httpreq.RestyResponse, result interface{}, msg string) error {
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
				return cm.Errorf(cm.ExternalServiceError, nil, "L???i kh??ng x??c ?????nh t??? viettelpost: null response. Ch??ng t??i ??ang li??n h??? v???i viettelpost ????? x??? l??. Xin l???i qu?? kh??ch v?? s??? b???t ti???n n??y. N???u c???n th??m th??ng tin vui l??ng li??n h??? %v.", wl.X(ctx).CSEmail)
			}
			if err = jsonx.Unmarshal(body, result); err != nil {
				return cm.Errorf(cm.ExternalServiceError, err, "L???i kh??ng x??c ?????nh t??? viettelpost: %v. Ch??ng t??i ??ang li??n h??? v???i viettelpost ????? x??? l??. Xin l???i qu?? kh??ch v?? s??? b???t ti???n n??y. N???u c???n th??m th??ng tin vui l??ng li??n h??? %v.", err, wl.X(ctx).CSEmail)
			}
			var responseFormat ResponseInterface
			if err := jsonx.Unmarshal(body, &responseFormat); err == nil {
				if responseFormat.Error {
					return cm.Errorf(cm.ExternalServiceError, nil, "L???i t??? viettelpost: %v (%v)", msg, responseFormat.Message)
				}
			}

			if err = jsonx.Unmarshal(body, result); err != nil {
				return cm.Errorf(cm.ExternalServiceError, err, "L???i kh??ng x??c ?????nh t??? viettelpost: %v. Ch??ng t??i ??ang li??n h??? v???i viettelpost ????? x??? l??. Xin l???i qu?? kh??ch v?? s??? b???t ti???n n??y. N???u c???n th??m th??ng tin vui l??ng li??n h??? %v.", err, wl.X(ctx).CSEmail)
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

		return cm.Errorf(cm.ExternalServiceError, nil, "L???i t??? viettelpost. N???u c???n th??m th??ng tin vui l??ng li??n h??? %v.", wl.X(ctx).CSEmail).WithMetaM(meta)
	default:
		return cm.Errorf(cm.ExternalServiceError, nil, "L???i kh??ng x??c ?????nh t??? viettelpost: Invalid status (%v). Ch??ng t??i ??ang li??n h??? v???i viettelpost ????? x??? l??. Xin l???i qu?? kh??ch v?? s??? b???t ti???n n??y. N???u c???n th??m th??ng tin vui l??ng li??n h??? %v", status, wl.X(ctx).CSEmail)
	}
}
