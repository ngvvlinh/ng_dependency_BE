package ahamove

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"time"

	"o.o/api/main/accountshipnow"
	"o.o/api/main/identity"
	"o.o/api/main/location"
	ordertypes "o.o/api/main/ordering/types"
	"o.o/api/main/shipnow"
	"o.o/api/main/shipnow/carrier"
	carriertypes "o.o/api/main/shipnow/carrier/types"
	shipnowtypes "o.o/api/main/shipnow/types"
	shippingtypes "o.o/api/main/shipping/types"
	"o.o/api/top/types/etc/shipping_fee_type"
	shipnowcarriertypes "o.o/backend/com/main/shipnow/carrier/types"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/integration/shipnow/ahamove/client"
)

var _ shipnowcarriertypes.ShipnowCarrier = &Carrier{}

type Carrier struct {
	client              *client.Client
	urlConfig           URLConfig
	location            location.QueryBus
	identityQuery       identity.QueryBus
	accountshipnowQuery accountshipnow.QueryBus
}

type URLConfig struct {
	ThirdPartyHost       string
	PathUserVerification string
}

func New(ahamoveClient *client.Client, urlConfig URLConfig, locationBus location.QueryBus, identityBus identity.QueryBus, accountshipnowQS accountshipnow.QueryBus) *Carrier {
	c := &Carrier{
		client:              ahamoveClient,
		urlConfig:           urlConfig,
		location:            locationBus,
		identityQuery:       identityBus,
		accountshipnowQuery: accountshipnowQS,
	}
	return c
}

func (c *Carrier) Code() carriertypes.ShipnowCarrier {
	return carriertypes.Ahamove
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
		OrderTime:      0,
		IdleUntil:      0,
		DeliveryPoints: deliveryPoints,
		Remarks:        cmd.ShippingNote,
		PromoCode:      cmd.Coupon,
	}
	response, err := c.client.CreateOrder(ctx, request)
	if err != nil {
		return nil, err
	}

	feelines := []*shippingtypes.ShippingFeeLine{
		{
			ShippingFeeType:     shipping_fee_type.Main,
			Cost:                int(response.Order.TotalFee),
			ExternalServiceName: "",
			ExternalServiceType: "",
		},
	}
	xshipnow = &carrier.ExternalShipnow{
		ID:         response.OrderId,
		UserID:     response.Order.UserId,
		Duration:   int(response.Order.Duration),
		Distance:   float32(response.Order.Distance),
		State:      client.OrderState(response.Status).ToCoreState(),
		TotalFee:   int(response.Order.TotalFee),
		FeeLines:   feelines,
		CreatedAt:  time.Now(),
		SharedLink: response.SharedLink,
	}

	return xshipnow, nil
}

func (c *Carrier) CancelExternalShipnow(ctx context.Context, cmd *carrier.CancelExternalShipnowCommand) error {
	request := &client.CancelOrderRequest{
		OrderId: cmd.ExternalShipnowID,
		Comment: cmd.CancelReason,
	}
	return c.client.CancelOrder(ctx, request)
}

func (c *Carrier) GetShipnowServices(ctx context.Context, args shipnowcarriertypes.GetShipnowServiceArgs) ([]*shipnowtypes.ShipnowService, error) {
	deliveryPoints, err := c.PrepareDeliveryPoints(ctx, args.PickupAddress, args.DeliveryPoints)
	if err != nil {
		return nil, err
	}
	request := &client.CalcShippingFeeRequest{
		OrderTime:      0,
		IdleUntil:      0,
		DeliveryPoints: deliveryPoints,
		PromoCode:      args.Coupon,
	}

	arbitraryID := args.ShopID.Int64() + args.ArbitraryID.Int64()
	services, err := c.calcShippingFee(ctx, arbitraryID, request)
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
		return nil, cm.Error(cm.ExternalServiceError, "ahamove: Kh??ng c?? g??i c?????c ph?? h???p", nil)
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
		return nil, cm.Error(cm.ExternalServiceError, "L???i t??? ahamove: Kh??ng th??? l???y th??ng tin g??i c?????c d???ch v???", nil)
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
			COD:            point.CODAmount,
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
		b.WriteString(line.ProductInfo.ProductName + " x " + strconv.Itoa(line.Quantity))
		b.WriteString("\n")
	}
	b.WriteString(point.ShippingNote)
	return b.String()
}

func (c *Carrier) ValidateAndGetAddress(ctx context.Context, in *ordertypes.Address) (*location.LocationQueryResult, error) {
	if in == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "?????a ch??? kh??ng ???????c ????? tr???ng")
	}
	if in.Coordinates == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "?????a ch??? thi???u th??ng tin latitude, longitude")
	}
	query := &location.GetLocationQuery{DistrictCode: in.DistrictCode}
	if err := c.location.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	return query.Result, nil
}

func toShipnowService(sfResp *client.CalcShippingFeeResponse, service *AhamoveShippingService, providerServiceID string) *shipnowtypes.ShipnowService {
	if sfResp == nil {
		return nil
	}

	res := &shipnowtypes.ShipnowService{
		Carrier:     carriertypes.Ahamove,
		Name:        service.Name,
		Code:        providerServiceID,
		Fee:         int(sfResp.TotalPrice),
		Description: service.Description,
	}

	// Avoid fee == 0
	if res.Fee == 0 {
		return nil
	}
	return res
}
