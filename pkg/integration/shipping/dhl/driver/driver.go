package driver

import (
	"context"
	"fmt"
	"strings"
	"time"

	"o.o/api/main/location"
	shipping2 "o.o/api/main/shipping"
	shippingstate "o.o/api/top/types/etc/shipping"
	"o.o/api/top/types/etc/shipping_provider"
	"o.o/api/top/types/etc/status4"
	"o.o/api/top/types/etc/status5"
	"o.o/api/top/types/etc/try_on"
	carriertypes "o.o/backend/com/main/shipping/carrier/types"
	shipmodel "o.o/backend/com/main/shipping/model"
	shippingsharemodel "o.o/backend/com/main/shipping/sharemodel"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/randgenerator"
	"o.o/backend/pkg/integration/shipping"
	dhlclient "o.o/backend/pkg/integration/shipping/dhl/client"
	"o.o/capi"
	"o.o/capi/dot"
)

var _ carriertypes.ShipmentCarrier = &DHLDriver{}

type DHLDriver struct {
	client     *dhlclient.Client
	locationQS location.QueryBus
	eventBus   capi.EventBus
}

func New(env string, cfg dhlclient.DHLAccountCfg, locationQS location.QueryBus, eventBus capi.EventBus) *DHLDriver {
	client := dhlclient.New(env, cfg)
	return &DHLDriver{
		client:     client,
		locationQS: locationQS,
		eventBus:   eventBus,
	}
}

func (d *DHLDriver) GetClient() *dhlclient.Client {
	return d.client
}

func (d *DHLDriver) Ping(ctx context.Context) error {
	return nil
}

func (d *DHLDriver) GetAffiliateID() string {
	return ""
}

func (d *DHLDriver) GenerateToken(ctx context.Context) (*carriertypes.GenerateTokenResponse, error) {
	generateAccessTokenResp, err := d.client.GenerateAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	accessTokenResponse := generateAccessTokenResp.AccessTokenResponse
	token := accessTokenResponse.Token.String()
	expiresIn := accessTokenResponse.ExpiresInSeconds.Int()
	expiresAt := time.Now().Add(time.Duration(expiresIn) * time.Second)

	d.client.UpdateToken(token)

	return &carriertypes.GenerateTokenResponse{
		AccessToken: token,
		ExpiresAt:   expiresAt,
		TokenType:   accessTokenResponse.TokenType.String(),
		ExpiresIn:   expiresIn,
	}, nil
}

func (d *DHLDriver) CreateFulfillment(
	ctx context.Context, ffm *shipmodel.Fulfillment,
	args *carriertypes.GetShippingServicesArgs, service *shippingsharemodel.AvailableShippingService,
) (ffmToUpdate *shipmodel.Fulfillment, _ error) {
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
	fromProvince := fromQuery.Result.Province

	addressFrom := dhlclient.ToAddress(ffm.AddressFrom)
	addressTo := dhlclient.ToAddress(ffm.AddressTo)
	// hi???n t???i remark kh??ng hi???n ??? h??? th???ng giao h??ng
	// field remark ????, ch??? n???m ??? tr??n h??? th???ng ????? b??o c??o theo d??i, ch??a ???????c link v???i h??? th???ng giao h??ng
	// ph????ng ??n l?? chia remark v??o address2 v?? address3
	// do address2 v?? address3 ??ang gi???i h???n 60 k?? t??? n??n ??? ????y c???t shippingNote v?? chia v??o address2, address3

	// n???i dung ghi ch?? giao h??ng
	// maximum length of address2 is 60
	addressTo.Address2 = strings.TrimSuffix(ffm.ShippingNote, carriertypes.CarrierNote)
	if len([]rune(addressTo.Address2)) > 60 {
		addressTo.Address2 = string(([]rune(addressTo.Address2))[:60])
	}
	addressTo.Address3 = carriertypes.CarrierNote

	var packageDesc string
	if len(ffm.Lines) > 0 {
		for i, line := range ffm.Lines {
			packageDesc += line.ProductName
			if i != len(ffm.Lines)-1 {
				packageDesc += ". "
			}
		}
	} else {
		if ffm.LinesContent != "" {
			packageDesc = ffm.LinesContent
		}
	}
	// maximum length of packageDesc is 50
	if len([]rune(packageDesc)) > 50 {
		packageDesc = string(([]rune(packageDesc))[:50])
	}

	serviceID, err := d.ParseServiceID(service.ProviderServiceID)
	if err != nil {
		return nil, err
	}

	// get insurance value
	maxValueFreeInsuranceFee := d.GetMaxValueFreeInsuranceFee()
	insuranceValue := args.GetInsuranceAmount(maxValueFreeInsuranceFee)

	now := time.Now()
	var vasServices *dhlclient.ValueAddedServices
	switch ffm.TryOn {
	case try_on.None:
	case try_on.Open, try_on.Try:
		// DHL ch??? support cho xem h??ng kh??ng th??? (OBOX)
		vasServices = &dhlclient.ValueAddedServices{
			ValueAddedService: []*dhlclient.ValueAddedService{{VasCode: dhlclient.VasCodeOBOX}},
		}
	default:

	}

	cmd := &dhlclient.CreateOrdersRequest{
		ManifestRequest: &dhlclient.ManiFestReq{
			Bd: &dhlclient.BdReq{
				PickupAccountID: string(getPickupAccountID(fromProvince.Region)),
				SoldToAccountID: "", // add it in client
				// HandoverMethod
				// Default DHL qua l???y h??ng
				HandoverMethod: dhlclient.HandoverMethodPickup,
				PickupAddress:  addressFrom,
				ShipmentItems: []*dhlclient.ShipmentItemReq{
					{
						ConsigneeAddress: addressTo,
						ReturnAddress:    addressFrom,
						ShipmentID:       ffm.ID.String(),
						ReturnMode:       string(dhlclient.ReturnToNewAddress),
						PackageDesc:      packageDesc,
						TotalWeight:      args.ChargeableWeight,
						TotalWeightUOM:   "g",
						Height:           float64(args.Height),
						Length:           float64(args.Length),
						Width:            float64(args.Width),
						ProductCode:      serviceID,
						// Return Product Code
						// Field optional, default l?? g??i PDO (d???ch v??? chu???n)
						// Tuy nhi??n, TopShip ch??? c?? g??i PDE (giao nhanh)
						// n??n n???u kh??ng truy???n field n??y l??n th?? DHL s??? b??o l???i (do ch??a k??ch ho???t g??i PDO)
						// Workaround: lu??n truy???n returnProductCode = ProductCode l??n
						ReturnProductCode: serviceID,
						CodValue:          args.CODAmount,
						// Total declared value of the shipment (in 2 decimal points). Mandatory for Cross Border shipment, optional for Domestic shipment.
						// For Vietnam Domestic, totalValue must be a multiple of 500.
						TotalValue:         float64(args.BasketValue),
						InsuranceValue:     float64(insuranceValue),
						Currency:           "VND",
						Remarks:            "KH??NG T??? ?? HO??N H??NG. G???i shop n???u giao th???t b???i",
						ValueAddedServices: vasServices,
					},
				},
			},
		},
	}
	r, err := d.client.CreateOrders(ctx, cmd)
	if err != nil {
		return nil, err
	}

	shipmentItem := r.ManifestResponse.Bd.ShipmentItems[0]
	updateFfm := &shipmodel.Fulfillment{
		ID:                        ffm.ID,
		ProviderServiceID:         service.ProviderServiceID,
		Status:                    status5.S, // Now processing
		ShippingCode:              shipmentItem.DeliveryConfirmationNo.String(),
		ExternalShippingName:      service.Name,
		ExternalShippingCode:      shipmentItem.DeliveryConfirmationNo.String(),
		ExternalShippingCreatedAt: now,
		ExternalShippingUpdatedAt: now,
		ShippingCreatedAt:         now,
		ShippingState:             shippingstate.Created,
		SyncStatus:                status4.P,
		SyncStates: &shippingsharemodel.FulfillmentSyncStates{
			SyncAt:    now,
			TrySyncAt: now,
		},
		InsuranceValue: dot.Int(insuranceValue),
	}

	return updateFfm, nil
}

func (d DHLDriver) RefreshFulfillment(ctx context.Context, fulfillment *shipmodel.Fulfillment) (ffmToUpdate *shipmodel.Fulfillment, _ error) {
	return nil, cm.Errorf(cm.ExternalServiceError, nil, "This carrier does not support this method")
}

func (d DHLDriver) UpdateFulfillmentInfo(ctx context.Context, fulfillment *shipmodel.Fulfillment) error {
	return cm.Errorf(cm.ExternalServiceError, nil, "This carrier does not support this method")
}

func (d DHLDriver) UpdateFulfillmentCOD(ctx context.Context, fulfillment *shipmodel.Fulfillment) error {
	return cm.Errorf(cm.ExternalServiceError, nil, "This carrier does not support this method")
}

func (d DHLDriver) CancelFulfillment(ctx context.Context, ffm *shipmodel.Fulfillment) error {
	externalCreatedAt := ffm.ExternalShippingCreatedAt

	// Add 10 minutes to createdAt
	// System behaviour DHL don't allow to cancel fulfillment before 10 minutes from created
	if externalCreatedAt.Add(10 * time.Minute).After(time.Now()) {
		event := &shipping2.DHLFulfillmentCancelledEvent{
			FulfillmentID: ffm.ID,
		}
		return d.eventBus.Publish(ctx, event)
	}

	shipmentID := ffm.ID.String()
	getLocationQuery := &location.GetLocationQuery{
		ProvinceCode: ffm.AddressFrom.ProvinceCode,
	}
	if err := d.locationQS.Dispatch(ctx, getLocationQuery); err != nil {
		return err
	}
	province := getLocationQuery.Result.Province

	cancelFulfillmentReq := &dhlclient.CancelOrderRequest{
		DeleteShipmentReq: &dhlclient.DeleteShipmentReq{
			Bd: &dhlclient.BdDeleteShipmentReq{
				PickupAccountId: string(getPickupAccountID(province.Region)),
				SoldToAccountId: "", // add it in client
				ShipmentItems: []*dhlclient.ShipmentItemDeleteShipmentReq{
					{
						ShipmentID: shipmentID,
					},
				},
			},
		}}
	if _, err := d.client.CancelOrder(ctx, cancelFulfillmentReq); err != nil {
		return err
	}
	return nil
}

func getPickupAccountID(VietNamRegion location.VietnamRegion) dhlclient.PickupAccountID {
	switch VietNamRegion {
	case location.North:
		return dhlclient.PickupAccountIDNorth
	case location.Middle:
		return dhlclient.PickupAccountIDMiddle
	case location.South:
		return dhlclient.PickupAccountIDSouth
	default:
		panic(fmt.Sprintf("unsupported region %v", VietNamRegion))
	}
}

func (d DHLDriver) GetShippingServices(ctx context.Context, args *carriertypes.GetShippingServicesArgs) ([]*shippingsharemodel.AvailableShippingService, error) {
	if args.FromWardCode == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "DHL: ?????a ch??? g???i h??ng - ph?????ng/x?? kh??ng ???????c ????? tr???ng!")
	}
	if args.ToWardCode == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "DHL: ?????a ch??? nh???n h??ng - ph?????ng/x?? kh??ng ???????c ????? tr???ng!")
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
		providerServiceID, err := GenerateServiceID(generator, service.ServiceCode.String())
		if err != nil {
			return nil, err
		}
		result = append(result, &shippingsharemodel.AvailableShippingService{
			Name:              service.Name.String(),
			Provider:          shipping_provider.DHL,
			ProviderServiceID: providerServiceID,
		})
	}

	result = shipping.CalcServicesTime(shipping_provider.DHL, fromDistrict, toDistrict, result)
	return result, nil
}

func GenerateServiceID(generator *randgenerator.RandGenerator, serviceName string) (string, error) {
	if serviceName == "" {
		return "", cm.Errorf(cm.Internal, nil, "Service Name can not be empty").WithMeta("DHL", "func GenerateServiceID")
	}

	code := generator.RandomAlphabet32(8)
	switch serviceName {
	case string(dhlclient.OrderServiceCodeSPD):
		code[1] = 'T'
	case string(dhlclient.OrderServiceCodePDE):
		code[1] = 'N'
	case string(dhlclient.OrderServiceCodePDO):
		code[1] = 'C'
	default:
		return "", cm.Errorf(cm.Internal, nil, "DHL invalid service name")
	}

	return string(code), nil
}

func DecodeShippingServiceName(code string) (name string, ok bool) {
	if len(code) < 6 {
		return "", false
	}
	switch {
	case code[1] == 'T': // T???i ??u
		return "T???i ??u", true
	case code[1] == 'N': // Nhanh
		return "Nhanh", true
	case code[1] == 'C': // Chu???n
		return "Chu???n", true
	}
	return "", false
}

func (d DHLDriver) GetServiceName(code string) (serviceName string, ok bool) {
	return DecodeShippingServiceName(code)
}

func (d DHLDriver) ParseServiceID(code string) (serviceID string, err error) {
	if code == "" {
		err = cm.Errorf(cm.InvalidArgument, nil, "Missing service id")
		return
	}
	if len(code) <= 3 {
		err = cm.Errorf(cm.InvalidArgument, nil, "DHL invalid service id (code = %v)", code)
		return
	}

	switch code[1] {
	case 'T': // T???i ??u
		serviceID = string(dhlclient.OrderServiceCodeSPD)
	case 'N': // Nhanh
		serviceID = string(dhlclient.OrderServiceCodePDE)
	case 'C': // Chu???n
		serviceID = string(dhlclient.OrderServiceCodePDO)
	default:
		return "", cm.Errorf(cm.Internal, nil, fmt.Sprintf("DHL invalid service id (code = %v)", code))
	}
	return serviceID, nil
}

func (d DHLDriver) GetMaxValueFreeInsuranceFee() int {
	return 0
}

func (d DHLDriver) SignIn(ctx context.Context, args *carriertypes.SignInArgs) (*carriertypes.AccountResponse, error) {
	return nil, cm.Errorf(cm.ExternalServiceError, nil, "This carrier does not support this method")
}

func (d DHLDriver) SignUp(ctx context.Context, args *carriertypes.SignUpArgs) (*carriertypes.AccountResponse, error) {
	return nil, cm.Errorf(cm.ExternalServiceError, nil, "This carrier does not support this method")
}
