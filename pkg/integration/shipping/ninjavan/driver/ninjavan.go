package driver

import (
	"context"
	"fmt"
	"time"

	"o.o/api/main/identity"
	"o.o/api/main/location"
	"o.o/api/main/shippingcode"
	shippingstate "o.o/api/top/types/etc/shipping"
	"o.o/api/top/types/etc/shipping_provider"
	"o.o/api/top/types/etc/status4"
	"o.o/api/top/types/etc/status5"
	carriertypes "o.o/backend/com/main/shipping/carrier/types"
	shipmodel "o.o/backend/com/main/shipping/model"
	shippingsharemodel "o.o/backend/com/main/shipping/sharemodel"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/randgenerator"
	etopmodel "o.o/backend/pkg/etop/model"
	"o.o/backend/pkg/integration/shipping"
	ninjavanclient "o.o/backend/pkg/integration/shipping/ninjavan/client"
	"o.o/capi/dot"
)

var _ carriertypes.ShipmentCarrier = &NinjaVanDriver{}

type NinjaVanDriver struct {
	client         *ninjavanclient.Client
	locationQS     location.QueryBus
	identityQS     identity.QueryBus
	shippingcodeQS shippingcode.QueryBus
}

func New(
	env string, cfg ninjavanclient.NinjaVanCfg,
	locationQS location.QueryBus, identityQS identity.QueryBus,
	shippingcodeQS shippingcode.QueryBus,
) *NinjaVanDriver {
	client := ninjavanclient.New(env, cfg)
	return &NinjaVanDriver{
		client:         client,
		locationQS:     locationQS,
		identityQS:     identityQS,
		shippingcodeQS: shippingcodeQS,
	}
}

func (d *NinjaVanDriver) Ping(ctx context.Context) error {
	return nil
}

func (d *NinjaVanDriver) GetAffiliateID() string {
	return ""
}

func (d *NinjaVanDriver) GenerateToken(ctx context.Context) (*carriertypes.GenerateTokenResponse, error) {
	accessTokenResponse, err := d.client.GenerateOAuthAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	expiresAt := time.Unix(int64(accessTokenResponse.Expires.Int()), 0)
	d.client.UpdateToken(accessTokenResponse.AccessToken.String())
	return &carriertypes.GenerateTokenResponse{
		AccessToken: accessTokenResponse.AccessToken.String(),
		ExpiresAt:   expiresAt,
		TokenType:   accessTokenResponse.TokenType.String(),
		ExpiresIn:   accessTokenResponse.ExpiresIn.Int(),
	}, nil
}

func (d *NinjaVanDriver) CreateFulfillment(ctx context.Context, ffm *shipmodel.Fulfillment, args *carriertypes.GetShippingServicesArgs, service *shippingsharemodel.AvailableShippingService) (ffmToUpdate *shipmodel.Fulfillment, _ error) {
	if args.FromWardCode == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Y??u c???u nh???p ph?????ng/x?? c???a ?????a ch??? l???y h??ng")
	}
	if args.ToWardCode == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Y??u c???u nh???p ph?????ng/x?? c???a ?????a ch??? nh???n h??ng")
	}

	fromQuery := &location.GetLocationQuery{
		DistrictCode: args.FromDistrictCode,
		WardCode:     args.FromWardCode,
	}
	toQuery := &location.GetLocationQuery{
		DistrictCode: args.ToDistrictCode,
		WardCode:     args.ToWardCode,
	}
	if err := d.locationQS.DispatchAll(ctx, fromQuery, toQuery); err != nil {
		return nil, err
	}

	serviceID, err := d.ParseServiceID(service.ProviderServiceID)
	if err != nil {
		return nil, err
	}

	getShopQuery := &identity.GetShopByIDQuery{ID: ffm.ShopID}
	if err := d.identityQS.Dispatch(ctx, getShopQuery); err != nil {
		return nil, err
	}
	shopName := getShopQuery.Result.Name

	// get insurance value
	maxValueFreeInsuranceFee := d.GetMaxValueFreeInsuranceFee()
	insuranceValue := args.GetInsuranceAmount(maxValueFreeInsuranceFee)
	addressFrom := ninjavanclient.ToAddress(ffm.AddressFrom)
	addressTo := ninjavanclient.ToAddress(ffm.AddressTo)
	now := time.Now()

	// pickupDate m???c ?????nh l???y ng??y t???o ????n
	pickupDate := now.Format(ninjavanclient.LayoutISO)
	// deliveryStartDate = now + 3 days
	deliveryStartDate := now.Add(ninjavanclient.ThreeDays).Format(ninjavanclient.LayoutISO)

	generateShippingCodeQuery := &shippingcode.GenerateShippingCodeQuery{}
	if err := d.shippingcodeQS.Dispatch(ctx, generateShippingCodeQuery); err != nil {
		return nil, err
	}
	shippingCode := generateShippingCodeQuery.Result

	cmd := &ninjavanclient.CreateOrderRequest{
		ServiceType:             string(ninjavanclient.ServiceTypeMarketPlace),
		ServiceLevel:            serviceID,
		RequestedTrackingNumber: shippingCode,
		Reference: &ninjavanclient.Reference{
			MerchantOrderNumber: shippingCode,
		},
		From: addressFrom,
		To:   addressTo,
		ParcelJob: &ninjavanclient.ParcelJob{
			IsPickupRequired: true,
			PickupDate:       pickupDate,
			PickupTimeslot: &ninjavanclient.Timeslot{
				StartTime: "09:00",
				EndTime:   "22:00",
				TimeZone:  ninjavanclient.TimeZoneHCM,
			},
			PickupInstructions: ffm.ShippingNote,
			PickupAddress:      addressFrom,
			DeliveryStartDate:  deliveryStartDate,
			DeliveryTimeslot: &ninjavanclient.Timeslot{
				StartTime: "09:00",
				EndTime:   "22:00",
				TimeZone:  ninjavanclient.TimeZoneHCM,
			},
			DeliveryInstructions: ffm.ShippingNote,
			CashOnDelivery:       float64(ffm.TotalCODAmount),
			InsuredValue:         float64(insuranceValue),
			Dimensions: &ninjavanclient.Dimensions{
				Weight: float64(args.ChargeableWeight) / 1000, // g -> kg
				Length: cm.CoalesceFloat(float64(args.Length), 10),
				Width:  cm.CoalesceFloat(float64(args.Width), 10),
				Height: cm.CoalesceFloat(float64(args.Height), 10),
			},
		},
		Marketplace: &ninjavanclient.Marketplace{
			SellerID:          ffm.ShopID.String(),
			SellerCompanyName: shopName,
		},
	}
	r, err := d.client.CreateOrder(ctx, cmd)
	if err != nil {
		return nil, err
	}

	expectedDeliveryAt, err := time.Parse(ninjavanclient.LayoutISO, r.ParcelJob.DeliveryStartDate)
	if err != nil {
		return nil, cm.Errorf(cm.ExternalServiceError, nil, fmt.Sprintf("Ninja Van: error when parse deliveryStartDate = %s", r.ParcelJob.DeliveryStartDate))
	}

	updateFfm := &shipmodel.Fulfillment{
		ID:                        ffm.ID,
		ProviderServiceID:         service.ProviderServiceID,
		Status:                    status5.S, // Now processing
		ShippingCode:              r.TrackingNumber.String(),
		ExternalShippingName:      service.Name,
		ExternalShippingCode:      r.TrackingNumber.String(),
		ExternalShippingCreatedAt: now,
		ExternalShippingUpdatedAt: now,
		ShippingCreatedAt:         now,
		ShippingState:             shippingstate.Created,
		SyncStatus:                status4.P,
		SyncStates: &shippingsharemodel.FulfillmentSyncStates{
			SyncAt:    now,
			TrySyncAt: now,
		},
		ExpectedDeliveryAt: expectedDeliveryAt,
		InsuranceValue:     dot.Int(insuranceValue),
	}

	return updateFfm, nil
}

func (d *NinjaVanDriver) RefreshFulfillment(ctx context.Context, fulfillment *shipmodel.Fulfillment) (ffmToUpdate *shipmodel.Fulfillment, _ error) {
	return nil, cm.Errorf(cm.ExternalServiceError, nil, "This carrier does not support this method")
}

func (d *NinjaVanDriver) UpdateFulfillmentInfo(ctx context.Context, fulfillment *shipmodel.Fulfillment) error {
	return cm.Errorf(cm.ExternalServiceError, nil, "This carrier does not support this method")
}

func (d *NinjaVanDriver) UpdateFulfillmentCOD(ctx context.Context, fulfillment *shipmodel.Fulfillment) error {
	return cm.Errorf(cm.ExternalServiceError, nil, "This carrier does not support this method")
}

func (d *NinjaVanDriver) CancelFulfillment(ctx context.Context, ffm *shipmodel.Fulfillment) error {
	code := ffm.ExternalShippingCode
	if _, err := d.client.CancelOrder(ctx, code); err != nil {
		return err
	}
	return nil
}

func (d *NinjaVanDriver) GetShippingServices(ctx context.Context, args *carriertypes.GetShippingServicesArgs) ([]*shippingsharemodel.AvailableShippingService, error) {
	if args.FromWardCode == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "NinjaVan: ?????a ch??? g???i h??ng - ph?????ng/x?? kh??ng ???????c ????? tr???ng!")
	}
	if args.ToWardCode == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "NinjaVan: ?????a ch??? nh???n h??ng - ph?????ng/x?? kh??ng ???????c ????? tr???ng!")
	}
	fromQuery := &location.GetLocationQuery{
		DistrictCode: args.FromDistrictCode,
		WardCode:     args.FromWardCode,
	}
	toQuery := &location.GetLocationQuery{
		DistrictCode: args.ToDistrictCode,
		WardCode:     args.ToWardCode,
	}
	if err := d.locationQS.DispatchAll(ctx, fromQuery, toQuery); err != nil {
		return nil, err
	}
	_, fromDistrict, _ := fromQuery.Result.Province, fromQuery.Result.District, fromQuery.Result.Ward
	_, toDistrict, _ := toQuery.Result.Province, toQuery.Result.District, toQuery.Result.Ward

	arbitraryID := args.AccountID.Int64() + args.ArbitraryID.Int64()

	var result []*shippingsharemodel.AvailableShippingService
	availableServices := d.client.FindAvailableServices(ctx)
	generator := randgenerator.NewGenerator(arbitraryID)
	for _, service := range availableServices.AvailableServices {
		providerServiceID, err := GenerateServiceID(generator, service.Name.String())
		if err != nil {
			return nil, err
		}
		result = append(result, &shippingsharemodel.AvailableShippingService{
			Name:              service.Name.String(),
			Provider:          shipping_provider.NinjaVan,
			ProviderServiceID: providerServiceID,
		})
	}

	result = shipping.CalcServicesTime(shipping_provider.NinjaVan, fromDistrict, toDistrict, result)
	return result, nil
}

func GenerateServiceID(generator *randgenerator.RandGenerator, serviceName string) (string, error) {
	if serviceName == "" {
		return "", cm.Errorf(cm.Internal, nil, "Service Name can not be empty").WithMeta("GHN", "func GenerateServiceID")
	}

	code := generator.RandomAlphabet32(8)
	switch serviceName {
	case string(ninjavanclient.ServiceLevelStandard):
		code[1] = 'C'
	default:
		return "", cm.Errorf(cm.Internal, nil, "Ninja Van invalid service name")
	}

	return string(code), nil
}

func DecodeShippingServiceName(code string) (name string, ok bool) {
	if len(code) < 6 {
		return "", false
	}
	switch {
	case code[1] == 'C': // Chu???n
		return etopmodel.ShippingServiceNameStandard, true
	}
	return "", false
}

func (d *NinjaVanDriver) GetServiceName(code string) (serviceName string, ok bool) {
	return DecodeShippingServiceName(code)
}

func (d *NinjaVanDriver) ParseServiceID(code string) (serviceID string, err error) {
	if code == "" {
		err = cm.Errorf(cm.InvalidArgument, nil, "Missing service id")
		return
	}
	if len(code) <= 3 {
		err = cm.Errorf(cm.InvalidArgument, nil, "Ninja Van invalid service id (code = %v)", code)
		return
	}

	switch code[1] {
	case 'C':
		serviceID = string(ninjavanclient.ServiceLevelStandard)
	default:
		return "", cm.Errorf(cm.Internal, nil, fmt.Sprintf("Ninja Van invalid service id (code = %v)", code))
	}
	return serviceID, nil
}

func (d *NinjaVanDriver) GetMaxValueFreeInsuranceFee() int {
	return 0
}

func (d *NinjaVanDriver) SignIn(ctx context.Context, args *carriertypes.SignInArgs) (*carriertypes.AccountResponse, error) {
	return nil, cm.Errorf(cm.ExternalServiceError, nil, "This carrier does not support this method")
}

func (d *NinjaVanDriver) SignUp(ctx context.Context, args *carriertypes.SignUpArgs) (*carriertypes.AccountResponse, error) {
	return nil, cm.Errorf(cm.ExternalServiceError, nil, "This carrier does not support this method")
}
