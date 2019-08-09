package ahamove

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"time"

	"etop.vn/api/main/identity"
	"etop.vn/api/main/location"
	ordertypes "etop.vn/api/main/ordering/types"
	"etop.vn/api/main/shipnow"
	"etop.vn/api/main/shipnow/carrier"
	shipnowtypes "etop.vn/api/main/shipnow/types"
	shippingtypes "etop.vn/api/main/shipping/types"
	shipnowcarrier "etop.vn/backend/com/main/shipnow-carrier"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/integration/shipnow/ahamove/client"
)

var _ shipnowcarrier.ShipnowCarrier = &Carrier{}
var identityQuery identity.QueryBus

type Carrier struct {
	client   *client.Client
	location location.QueryBus
}

func New(cfg client.Config, locationBus location.QueryBus, identityBus identity.QueryBus, urlConfig URLConfig) (*Carrier, *CarrierAccount) {
	ahamoveClient := client.New(cfg)
	identityQuery = identityBus
	c := &Carrier{
		client:   ahamoveClient,
		location: locationBus,
	}
	ca := &CarrierAccount{
		client:    ahamoveClient,
		urlConfig: urlConfig,
	}
	return c, ca
}

func (c *Carrier) GetServiceName(code string) (serviceName string, ok bool) {
	return decodeShippingServiceName(code)
}

func (c *Carrier) ParseServiceCode(code string) (serviceCode string, _err error) {
	sCode, err := parseServiceCode(code)
	return sCode, err
}

func (c *Carrier) InitClient(ctx context.Context) error {
	if err := c.client.Ping(); err != nil {
		return cm.Errorf(cm.ExternalServiceError, err, "can not init client")
	}
	return nil
}

func (c *Carrier) CreateExternalShipnow(ctx context.Context, cmd *carrier.CreateExternalShipnowCommand, service *shipnowtypes.ShipnowService) (xshipnow *carrier.ExternalShipnow, _err error) {
	queryShop := &identity.GetShopByIDQuery{
		ID: cmd.ShopID,
	}
	if err := identityQuery.Dispatch(ctx, queryShop); err != nil {
		return nil, err
	}
	userID := queryShop.Result.Shop.OwnerID
	if ok, err := isXAccountAhamoveVerified(ctx, userID); err != nil {
		return nil, err
	} else if !ok {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Vui lòng gửi yêu cầu xác thực tài khoản Ahamove trước khi tạo đơn.")
	}

	token, err := getToken(ctx, userID)
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Token không được để trống. Vui lòng tạo tài khoản Ahamove")
	}

	deliveryPoints, err := c.PrepareDeliveryPoints(ctx, cmd.PickupAddress, cmd.DeliveryPoints)
	if err != nil {
		return nil, err
	}

	serviceID, err := parseServiceCode(service.Code)
	if err != nil {
		return nil, err
	}
	request := &client.CreateOrderRequest{
		ServiceID:      serviceID,
		Token:          token,
		OrderTime:      0,
		IdleUntil:      0,
		DeliveryPoints: deliveryPoints,
		Remarks:        cmd.ShippingNote,
	}
	response, err := c.client.CreateOrder(ctx, request)
	if err != nil {
		return nil, err
	}

	feelines := []*shippingtypes.FeeLine{
		{
			ShippingFeeType:     shippingtypes.FeeLineTypeMain,
			Cost:                int32(response.Order.TotalFee),
			ExternalServiceName: "",
			ExternalServiceType: "",
		},
	}
	xshipnow = &carrier.ExternalShipnow{
		ID:         response.OrderId,
		UserID:     response.Order.UserId,
		Duration:   response.Order.Duration,
		Distance:   response.Order.Distance,
		State:      client.OrderState(response.Status).ToCoreState(),
		TotalFee:   response.Order.TotalFee,
		FeeLines:   feelines,
		CreatedAt:  time.Now(),
		SharedLink: response.SharedLink,
	}

	return xshipnow, nil
}

func (c *Carrier) CancelExternalShipnow(ctx context.Context, cmd *carrier.CancelExternalShipnowCommand) error {
	queryShop := &identity.GetShopByIDQuery{
		ID: cmd.ShopID,
	}
	if err := identityQuery.Dispatch(ctx, queryShop); err != nil {
		return err
	}
	userID := queryShop.Result.Shop.OwnerID

	token, err := getToken(ctx, userID)
	if err != nil {
		return cm.Errorf(cm.InvalidArgument, nil, "Token không được để trống. Vui lòng tạo tài khoản Ahamove")
	}

	request := &client.CancelOrderRequest{
		Token:   token,
		OrderId: cmd.ExternalShipnowID,
		Comment: cmd.CancelReason,
	}
	return c.client.CancelOrder(ctx, request)
}

func (c *Carrier) GetShippingServices(ctx context.Context, args shipnowcarrier.GetShippingServiceArgs) ([]*shipnowtypes.ShipnowService, error) {
	deliveryPoints, err := c.PrepareDeliveryPoints(ctx, args.PickupAddress, args.DeliveryPoints)
	if err != nil {
		return nil, err
	}
	request := &client.CalcShippingFeeRequest{
		OrderTime:      0,
		IdleUntil:      0,
		DeliveryPoints: deliveryPoints,
	}
	services, err := c.calcShippingFee(ctx, args.ShopID, request)
	if err != nil {
		return nil, err
	}
	return services, nil
}

// arbitraryID is provided as a seed, for stable randomization
func (c *Carrier) calcShippingFee(
	ctx context.Context, arbitraryID int64, request *client.CalcShippingFeeRequest,
) (res []*shipnowtypes.ShipnowService, _ error) {
	type Result struct {
		Service *AhamoveShippingService
		Result  *client.CalcShippingFeeResponse
		Error   error
	}
	var results []Result
	var wg sync.WaitGroup
	var m sync.Mutex

	services := c.getAvailableServices(ctx, request.DeliveryPoints)
	if len(services) == 0 {
		return nil, cm.Error(cm.ExternalServiceError, "ahamove: Không có gói cước phù hợp", nil)
	}
	wg.Add(len(services))
	for _, s := range services {
		go func(s *AhamoveShippingService) {
			defer wg.Done()
			req := *request // clone the request
			req.ServiceID = s.Code
			resp, err := c.client.CalcShippingFee(ctx, &req)
			m.Lock()
			result := Result{
				s, resp, err,
			}
			results = append(results, result)
			m.Unlock()
		}(s)
	}
	wg.Wait()
	if len(results) == 0 {
		return nil, cm.Error(cm.ExternalServiceError, "Lỗi từ ahamove: Không thể lấy thông tin gói cước dịch vụ", nil)
	}

	generator := newServiceIDGenerator(arbitraryID)
	for _, result := range results {
		providerServiceID, err := generator.generateServiceID(result.Service.Code)
		if err != nil {
			return nil, err
		}
		if result.Error != nil {
			continue
		}
		_r := toShipnowService(result.Result, result.Service, providerServiceID)
		res = append(res, _r)
	}
	return res, nil
}

func (c *Carrier) PrepareDeliveryPoints(ctx context.Context, pickupAddress *ordertypes.Address, points []*shipnow.DeliveryPoint) (deliveryPoints []*client.DeliveryPointRequest, _ error) {
	// add pickupAddress
	pickupLocation, err := c.ValidateAndGetAddress(ctx, pickupAddress)
	if err != nil {
		return nil, err
	}
	deliveryPoints = append(deliveryPoints, &client.DeliveryPointRequest{
		Address:      ordertypes.GetFullAddress(pickupAddress, pickupLocation),
		Lat:          pickupAddress.Coordinates.Latitude,
		Lng:          pickupAddress.Coordinates.Longitude,
		Mobile:       pickupAddress.Phone,
		Name:         pickupAddress.FullName,
		ProvinceCode: pickupAddress.ProvinceCode,
		DistrictCode: pickupAddress.DistrictCode,
		WardCode:     pickupAddress.WardCode,
	})

	for _, point := range points {
		address := point.ShippingAddress
		loc, err := c.ValidateAndGetAddress(ctx, address)
		if err != nil {
			return nil, err
		}
		_point := &client.DeliveryPointRequest{
			Address:        ordertypes.GetFullAddress(point.ShippingAddress, loc),
			Lat:            address.Coordinates.Latitude,
			Lng:            address.Coordinates.Longitude,
			Mobile:         address.Phone,
			Name:           address.FullName,
			COD:            point.CodAmount,
			TrackingNumber: point.OrderCode,
			Remarks:        prepareRemarksForDeliveryPoint(point),
			ProvinceCode:   address.ProvinceCode,
			DistrictCode:   address.DistrictCode,
			WardCode:       address.WardCode,
		}
		deliveryPoints = append(deliveryPoints, _point)
	}
	return deliveryPoints, nil
}

func prepareRemarksForDeliveryPoint(point *shipnow.DeliveryPoint) string {
	var b strings.Builder
	for _, line := range point.Lines {
		b.WriteString(line.ProductInfo.ProductName + " x " + strconv.Itoa(int(line.Quantity)))
		b.WriteString("\n")
	}
	b.WriteString(point.ShippingNote)
	return b.String()
}

func (c *Carrier) ValidateAndGetAddress(ctx context.Context, in *ordertypes.Address) (*location.LocationQueryResult, error) {
	if in.Coordinates == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Địa chỉ thiếu thông tin latitude, longitude")
	}
	query := &location.GetLocationQuery{DistrictCode: in.DistrictCode}
	if err := c.location.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	return query.Result, nil
}

func getToken(ctx context.Context, userID int64) (token string, _err error) {
	queryUser := &identity.GetUserByIDQuery{
		UserID: userID,
	}
	if err := identityQuery.Dispatch(ctx, queryUser); err != nil {
		return "", err
	}
	user := queryUser.Result

	query := &identity.GetExternalAccountAhamoveQuery{
		Phone:   user.Phone,
		OwnerID: user.ID,
	}
	if err := identityQuery.Dispatch(ctx, query); err != nil {
		return "", err
	}
	return query.Result.ExternalToken, nil
}

func isXAccountAhamoveVerified(ctx context.Context, userID int64) (bool, error) {
	queryUser := &identity.GetUserByIDQuery{
		UserID: userID,
	}
	if err := identityQuery.Dispatch(ctx, queryUser); err != nil {
		return false, err
	}
	user := queryUser.Result

	query := &identity.GetExternalAccountAhamoveQuery{
		OwnerID: user.ID,
		Phone:   user.Phone,
	}
	if err := identityQuery.Dispatch(ctx, query); err != nil {
		return false, err
	}
	return query.Result.ExternalVerified, nil
}

func toShipnowService(sfResp *client.CalcShippingFeeResponse, service *AhamoveShippingService, providerServiceID string) *shipnowtypes.ShipnowService {
	if sfResp == nil {
		return nil
	}

	res := &shipnowtypes.ShipnowService{
		Carrier:     carrier.Ahamove,
		Name:        service.Name,
		Code:        providerServiceID,
		Fee:         int32(sfResp.TotalFee),
		Description: service.Description,
	}
	// BIKE/POOL: discount, total_fee, total_pay
	// SAMEDAY: partner_discount, partner_fee, partner_pay
	// Ahamove đang fix, sau này sẽ dùng total_fee hết
	if strings.Contains(service.Code, string(SAMEDAY)) {
		res.Fee = int32(sfResp.PartnerFee)
	}

	// Avoid fee == 0
	if res.Fee == 0 {
		return nil
	}
	return res
}
