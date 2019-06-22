package ahamove

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"etop.vn/backend/cmd/etop-server/config"

	"etop.vn/api/main/identity"

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
var thirdPartyHost string

type Carrier struct {
	client   *ahamoveclient.Client
	location location.QueryBus
}

type CarrierAccount struct {
	client *ahamoveclient.Client
}

func New(cfg ahamoveclient.Config, locationBus location.QueryBus, identityBus identity.QueryBus, thirdParty string) (*Carrier, *CarrierAccount) {
	client := ahamoveclient.New(cfg)
	identityQuery = identityBus
	thirdPartyHost = thirdParty
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
		ID:         res.OrderId,
		UserID:     res.Order.UserId,
		Duration:   res.Order.Duration,
		Distance:   res.Order.Distance,
		State:      ahamoveclient.OrderState(res.Status).ToCoreState(),
		TotalFee:   res.Order.TotalFee,
		FeeLines:   feelines,
		CreatedAt:  time.Now(),
		SharedLink: res.SharedLink,
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
			Mobile:  args.Phone,
			Name:    args.Name,
			Address: args.Address,
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
	token, err := getToken(ctx, args.OwnerID)
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
	createAt := time.Unix(int64(account.CreateTime), 0)

	res := &carrier.ExternalAccount{
		ID:        account.ID,
		Name:      account.Name,
		Email:     account.Email,
		Verified:  account.Verified,
		CreatedAt: createAt,
	}
	return res, nil
}

func (c *CarrierAccount) VerifyExternalAccount(ctx context.Context, args *shipnow_carrier.VerifyExternalAccountArgs) (*carrier.VerifyExternalAccountResult, error) {
	token, err := getToken(ctx, args.OwnerID)
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Token không được để trống. Vui lòng tạo tài khoản Ahamove")
	}

	description, err := getDescriptionForVerification(ctx, args.OwnerID)
	if err != nil {
		return nil, err
	}

	cmd := &VerifyAccountCommand{
		Request: &ahamoveclient.VerifyAccountRequest{
			Token:       token,
			Description: description,
		},
	}
	if err := c.VerifyAccount(ctx, cmd); err != nil {
		return nil, err
	}
	res := &carrier.VerifyExternalAccountResult{
		TicketID: strconv.Itoa(cmd.Result.Ticket.ID),
	}
	return res, nil
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

func prepareAhamovePhotoUrl(ahamoveAccount *identity.ExternalAccountAhamove, uri string, typeImg string) string {
	ext := filepath.Ext(uri)
	filename := strings.TrimSuffix(filepath.Base(uri), ext)

	newName := ""
	switch typeImg {
	case "front":
		newName = fmt.Sprintf("user_id_front_%v_%v", ahamoveAccount.ExternalID, ahamoveAccount.ExternalCreatedAt.Unix())
	case "back":
		newName = fmt.Sprintf("user_id_back_%v_%v", ahamoveAccount.ExternalID, ahamoveAccount.ExternalCreatedAt.Unix())
	case "portrait":
		newName = fmt.Sprintf("user_portrait_%v_%v", ahamoveAccount.ExternalID, ahamoveAccount.ExternalCreatedAt.Unix())
	}

	// expected result:
	// https://3rd.d.etop.vn/ahamove/user_verification/BdVzaWz6ssamNKrRV7W8/user_id_front_84909090999_1444118656.jpg
	return fmt.Sprintf("%v%v/%v/%v%v", thirdPartyHost, config.PathAhamoveUserVerification, filename, newName, ext)
}

// description format: <user._id>, <user.name>, <photo_urls>
// photo_url format: <topship_domain>/upload/ahamove/user_verification/user_id_front<user.id>_<user.create_time>.jpg

func getDescriptionForVerification(ctx context.Context, userID int64) (des string, _err error) {
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
	account := query.Result

	front := prepareAhamovePhotoUrl(account, account.IDCardFrontImg, "front")
	back := prepareAhamovePhotoUrl(account, account.IDCardBackImg, "back")
	portrait := prepareAhamovePhotoUrl(account, account.PortraitImg, "portrait")

	des = fmt.Sprintf("%v, %v, %v, %v, %v", account.ExternalID, account.Name, front, back, portrait)
	return des, nil
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
	account := query.Result
	if !account.ExternalVerified {
		return false, nil
	}
	return true, nil
}
