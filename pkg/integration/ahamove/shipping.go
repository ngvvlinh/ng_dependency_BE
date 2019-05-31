package ahamove

import (
	"context"
	"strconv"
	"strings"
	"time"

	"etop.vn/api/main/identity"

	"etop.vn/backend/pkg/common/bus"

	"github.com/k0kubun/pp"

	"etop.vn/api/main/shipnow"

	"etop.vn/api/main/location"
	ordertypes "etop.vn/api/main/ordering/types"
	"etop.vn/api/main/shipnow/carrier"
	shipnowtypes "etop.vn/api/main/shipnow/types"
	shippingtypes "etop.vn/api/main/shipping/types"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/etop/logic/shipping_provider"
	"etop.vn/backend/pkg/etop/model"
	ahamoveclient "etop.vn/backend/pkg/integration/ahamove/client"
	shipnow_carrier "etop.vn/backend/pkg/services/shipnow-carrier"
)

var _ shipnow_carrier.ShipnowCarrier = &Carrier{}
var _ shipnow_carrier.ShipnowCarrierAccount = &CarrierAccount{}
var identityQuery identity.QueryBus

type Carrier struct {
	client   *ahamoveclient.Client
	location location.QueryBus
}

type CarrierAccount struct {
	client *ahamoveclient.Client
}

func New(cfg ahamoveclient.Config, locationBus location.QueryBus, identityBus identity.QueryBus) (*Carrier, *CarrierAccount) {
	client := ahamoveclient.New(cfg)
	identityQuery = identityBus
	return &Carrier{
		client:   client,
		location: locationBus,
	}, &CarrierAccount{client: client}
}

func (c *Carrier) InitClient(ctx context.Context) error {
	if err := c.client.Ping(); err != nil {
		return cm.Errorf(cm.ExternalServiceError, err, "can not init client")
	}
	return nil
}

func (c *Carrier) CreateExternalShipnow(ctx context.Context, cmd *carrier.CreateExternalShipnowCommand, service *shipnowtypes.ShipnowService) (xshipnow *carrier.ExternalShipnow, _err error) {
	token, err := getToken(ctx, cmd.ShopID)
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Token không được để trống. Vui lòng tạo tài khoản Ahamove")
	}

	deliveryPoints, err := c.PrepareDeliveryPoints(ctx, cmd.PickupAddress, cmd.DeliveryPoints)
	if err != nil {
		return nil, err
	}
	ahamoveCmd := &CreateOrderCommand{
		ServiceID: service.Code,
		Request: &ahamoveclient.CreateOrderRequest{
			Token:          token,
			OrderTime:      0,
			IdleUntil:      0,
			DeliveryPoints: deliveryPoints,
			Remarks:        cmd.ShippingNote,
		},
	}
	err = c.CreateOrder(ctx, ahamoveCmd)
	pp.Println("res :: ", ahamoveCmd.Result)
	if err != nil {
		return nil, err
	}
	res := ahamoveCmd.Result
	feelines := []*shippingtypes.FeeLine{
		{
			ShippingFeeType:     shippingtypes.FeeLineTypeMain,
			Cost:                int32(res.Order.TotalFee),
			ExternalServiceName: "",
			ExternalServiceType: "",
		},
	}
	xshipnow = &carrier.ExternalShipnow{
		ID:        res.OrderId,
		UserID:    res.Order.UserId,
		Duration:  res.Order.Duration,
		Distance:  res.Order.Distance,
		State:     ahamoveclient.OrderState(res.Status).ToCoreState(),
		TotalFee:  res.Order.TotalFee,
		FeeLines:  feelines,
		CreatedAt: time.Now(),
	}

	return xshipnow, nil
}

func (c *Carrier) CancelExternalShipnow(ctx context.Context, cmd *carrier.CancelExternalShipnowCommand) error {
	token, err := getToken(ctx, cmd.ShopID)
	if err != nil {
		return cm.Errorf(cm.InvalidArgument, nil, "Token không được để trống. Vui lòng tạo tài khoản Ahamove")
	}

	ahamoveCmd := &CancelOrderCommand{
		ServiceID: cmd.CarrierServiceCode,
		Request: &ahamoveclient.CancelOrderRequest{
			Token:   token,
			OrderId: cmd.ExternalShipnowID,
			Comment: cmd.CancelReason,
		},
	}
	return c.CancelOrder(ctx, ahamoveCmd)
}

func (c *Carrier) GetShippingServices(ctx context.Context, args shipnow_carrier.GetShippingServiceArgs) ([]*shipnowtypes.ShipnowService, error) {
	deliveryPoints, err := c.PrepareDeliveryPoints(ctx, args.PickupAddress, args.DeliveryPoints)
	if err != nil {
		return nil, err
	}
	cmd := &CalcShippingFeeCommand{
		ArbitraryID: args.ShopID,
		Request: &ahamoveclient.CalcShippingFeeRequest{
			OrderTime:      0,
			IdleUntil:      0,
			DeliveryPoints: deliveryPoints,
		},
	}
	err = c.CalcShippingFee(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return cmd.Result, nil
}

func (c *Carrier) PrepareDeliveryPoints(ctx context.Context, pickupAddress *ordertypes.Address, points []*shipnow.DeliveryPoint) (deliveryPoints []*ahamoveclient.DeliveryPointRequest, _ error) {
	// add pickupAddress
	pickupLocation, err := c.ValidateAndGetAddress(ctx, pickupAddress)
	if err != nil {
		return nil, err
	}
	deliveryPoints = append(deliveryPoints, &ahamoveclient.DeliveryPointRequest{
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
		_point := &ahamoveclient.DeliveryPointRequest{
			Address:        ordertypes.GetFullAddress(point.ShippingAddress, loc),
			Lat:            address.Coordinates.Latitude,
			Lng:            address.Coordinates.Longitude,
			Mobile:         address.Phone,
			Name:           address.FullName,
			COD:            point.CodAmount,
			TrackingNumber: strconv.FormatInt(point.OrderId, 10),
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

func prepareProducts(points []*shipnow.DeliveryPoint) (products []ahamoveclient.Item, _ error) {
	for _, point := range points {
		for _, line := range point.Lines {
			p := ahamoveclient.Item{
				Num:   int(line.Quantity),
				Name:  line.ProductInfo.ProductName,
				Price: int(line.TotalPrice),
			}
			products = append(products, p)
		}
	}
	if len(products) == 0 {
		return products, cm.Errorf(cm.InvalidArgument, nil, "Đơn vận chuyển thiếu thông tin sản phẩm")
	}
	return
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

func (c *Carrier) GetAllShippingServices(ctx context.Context, args shipping_provider.GetShippingServicesArgs) ([]*model.AvailableShippingService, error) {
	return nil, cm.ErrTODO
}

func (c *CarrierAccount) RegisterExternalAccount(ctx context.Context, args *shipnow_carrier.RegisterExternalAccountArgs) (*carrier.RegisterExternalAccountResult, error) {
	cmd := &RegisterAccountCommand{
		Request: &ahamoveclient.RegisterAccountRequest{
			Mobile: args.Phone,
			Name:   args.Name,
		},
	}

	if err := c.RegisterAccount(ctx, cmd); err != nil {
		return nil, err
	}
	res := &carrier.RegisterExternalAccountResult{
		Token: cmd.Result.Token,
	}
	return res, nil
}

func (c *CarrierAccount) GetExternalAccount(ctx context.Context, args *shipnow_carrier.GetExternalAccountArgs) (*carrier.ExternalAccount, error) {
	token, err := getToken(ctx, args.ShopID)
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Token không được để trống. Vui lòng tạo tài khoản Ahamove")
	}

	cmd := &GetAccountCommand{
		Request: &ahamoveclient.GetAccountRequest{
			Token: token,
		},
	}
	if err := c.GetAccount(ctx, cmd); err != nil {
		return nil, err
	}
	account := cmd.Result
	res := &carrier.ExternalAccount{
		ID:       account.ID,
		Name:     account.Name,
		Email:    account.Email,
		Verified: account.Verified,
	}
	return res, nil
}

func getToken(ctx context.Context, shopID int64) (token string, _err error) {
	queryShop := &model.GetShopExtendedQuery{
		ShopID: shopID,
	}
	if err := bus.Dispatch(ctx, queryShop); err != nil {
		return "", err
	}
	user := queryShop.Result.User
	query := &identity.GetExternalAccountAhamoveByPhoneQuery{
		Phone:   user.Phone,
		OwnerID: user.ID,
	}
	if err := identityQuery.Dispatch(ctx, query); err != nil {
		return "", err
	}
	return query.Result.ExternalToken, nil
}
