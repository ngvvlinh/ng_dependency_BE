package identity

import (
	"context"
	"fmt"
	"time"

	"etop.vn/api/external/haravan/identity"
	"etop.vn/api/meta"
	"etop.vn/backend/com/external/haravan/identity/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/sql/cmsql"
	haravanclient "etop.vn/backend/pkg/integration/haravan/client"
	"etop.vn/capi/dot"
)

const (
	PathGetShippingRates = "GetShippingRates"
	PathCreateOrder      = "CreateOrder"
	PathGetOrder         = "GetOrder"
	PathCancelOrder      = "CancelOrder"
	CarrierServiceName   = "TOPSHIP"
)

func BuildGatewayRoute(path string) string {
	return "/haravan/gateway/:shopid/" + path
}

func BuildURLForRegistration(thirdPartyHost string, externalShopID int) func(path string) string {
	return func(path string) string {
		return fmt.Sprintf(
			"%v/haravan/gateway/%v/%v",
			thirdPartyHost, externalShopID, path)
	}
}

var _ identity.Aggregate = &Aggregate{}
var thirdPartyHost string

type Aggregate struct {
	xAccountHaravanStore sqlstore.XAccountHaravanStoreFactory
	haravanClient        *haravanclient.Client
}

func NewAggregate(db *cmsql.Database, thirdParty string, cfg haravanclient.Config) *Aggregate {
	thirdPartyHost = thirdParty
	return &Aggregate{
		xAccountHaravanStore: sqlstore.NewXAccountHaravanStore(db),
		haravanClient:        haravanclient.New(cfg),
	}
}

func (a *Aggregate) MessageBus() identity.CommandBus {
	b := bus.New()
	return identity.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *Aggregate) CreateExternalAccountHaravan(ctx context.Context, args *identity.CreateExternalAccountHaravanArgs) (*identity.ExternalAccountHaravan, error) {
	account, err := a.xAccountHaravanStore(ctx).ShopID(args.ShopID).GetXAccountHaravan()
	if err != nil && cm.ErrorCode(err) != cm.NotFound {
		return nil, err
	}

	if account != nil && account.AccessToken != "" {
		return account, nil
	}

	var id dot.ID
	// code: user for Oauth
	cmd := &haravanclient.GetAccessTokenRequest{
		Subdomain:   args.Subdomain,
		Code:        args.Code,
		RedirectURI: args.RedirectURI,
	}
	tokenResp, err := a.haravanClient.GetAccessToken(ctx, cmd)
	if err != nil {
		return nil, err
	}
	// Get Haravan account
	query := &haravanclient.GetShopRequest{
		Connection: haravanclient.Connection{
			Subdomain: args.Subdomain,
			TokenStr:  tokenResp.AccessToken,
		},
	}
	externalShop, err := a.haravanClient.GetShop(ctx, query)
	if err != nil {
		return nil, err
	}

	if account == nil {
		// create new account
		id = cm.NewID()
		createArgs := &sqlstore.CreateXAccountHaravanArgs{
			ID:             id,
			ShopID:         args.ShopID,
			Subdomain:      args.Subdomain,
			AccessToken:    tokenResp.AccessToken,
			ExternalShopID: externalShop.Id,
		}

		return a.xAccountHaravanStore(ctx).CreateXAccountHaravan(createArgs)
	}

	args2 := &sqlstore.UpdateXAccountHaravanInfoArgs{
		ShopID:         account.ShopID,
		Subdomain:      args.Subdomain,
		AccessToken:    tokenResp.AccessToken,
		ExternalShopID: externalShop.Id,
	}
	return a.xAccountHaravanStore(ctx).UpdateXAccountHaravan(args2)
}

func (a *Aggregate) UpdateExternalAccountHaravanToken(ctx context.Context, args *identity.UpdateExternalAccountHaravanTokenArgs) (*identity.ExternalAccountHaravan, error) {
	cmd := &haravanclient.GetAccessTokenRequest{
		Subdomain:   args.Subdomain,
		Code:        args.Code,
		RedirectURI: args.RedirectURI,
	}
	tokenResp, err := a.haravanClient.GetAccessToken(ctx, cmd)
	if err != nil {
		return nil, err
	}

	cmdUpdate := &sqlstore.UpdateXAccountHaravanInfoArgs{
		ShopID:      args.ShopID,
		Subdomain:   args.Subdomain,
		AccessToken: tokenResp.AccessToken,
	}
	return a.xAccountHaravanStore(ctx).UpdateXAccountHaravan(cmdUpdate)
}

func (a *Aggregate) GetExternalAccountIDHaravan(ctx context.Context, shopID dot.ID) (int, error) {
	account, err := a.xAccountHaravanStore(ctx).ShopID(shopID).GetXAccountHaravan()
	if err != nil {
		return 0, nil
	}
	cmd := &haravanclient.GetShopRequest{
		Connection: haravanclient.Connection{
			Subdomain: account.Subdomain,
			TokenStr:  account.AccessToken,
		},
	}
	externalShop, err := a.haravanClient.GetShop(ctx, cmd)
	if err != nil {
		return 0, nil
	}
	return externalShop.Id, nil
}

func (a *Aggregate) UpdateExternalShopIDAccountHaravan(ctx context.Context, args *identity.UpdateExternalShopIDAccountHaravanArgs) (*identity.ExternalAccountHaravan, error) {
	cmd := &sqlstore.UpdateXShopIDAccountHaravanArgs{
		ShopID:         args.ShopID,
		ExternalShopID: args.ExternalShopID,
	}
	return a.xAccountHaravanStore(ctx).UpdateXShopIDAccountHaravan(cmd)
}

func (a *Aggregate) ConnectCarrierServiceExternalAccountHaravan(ctx context.Context, args *identity.ConnectCarrierServiceExternalAccountHaravanArgs) (*meta.Empty, error) {
	account, err := a.xAccountHaravanStore(ctx).ShopID(args.ShopID).GetXAccountHaravan()
	if err != nil {
		return nil, err
	}

	currentCarrierServiceID, _ := a.GetCurrentCarrierServiceID(ctx, args.ShopID)
	if currentCarrierServiceID != 0 && currentCarrierServiceID == account.ExternalCarrierServiceID {
		// Keep the old connection
		return &meta.Empty{}, nil
	}
	if currentCarrierServiceID != 0 {
		// Delete this connection to create the new one
		cmd := &haravanclient.DeleteConnectedCarrierServiceRequest{
			Connection: haravanclient.Connection{
				Subdomain: account.Subdomain,
				TokenStr:  account.AccessToken,
			},
			CarrierServiceID: currentCarrierServiceID,
		}
		if err := a.haravanClient.DeleteConnectedCarrierService(ctx, cmd); err != nil {
			return nil, err
		}
		// Delete all ConnectedXCarrierService
		updateArgs := &sqlstore.UpdateDeleteConnectedXCarrierSeriveArgs{
			SubDomain: account.Subdomain,
		}
		if err := a.xAccountHaravanStore(ctx).UpdateDeleteConnectedXCarrierService(updateArgs); err != nil {
			return nil, err
		}
	}

	buildURL := BuildURLForRegistration(thirdPartyHost, account.ExternalShopID)
	cmd := &haravanclient.ConnectCarrierServiceRequest{
		Connection: haravanclient.Connection{
			Subdomain: account.Subdomain,
			TokenStr:  account.AccessToken,
		},
		CarrierService: &haravanclient.CarrierService{
			Active:              true,
			TrackingUrl:         "https://www.etop.vn/",
			CreateOrderUrl:      buildURL(PathCreateOrder),
			GetOrderDetailUrl:   buildURL(PathGetOrder),
			GetShippingRatesUrl: buildURL(PathGetShippingRates),
			CancelOrderUrl:      buildURL(PathCancelOrder),
			Name:                CarrierServiceName,
		},
	}
	res, err := a.haravanClient.ConnectCarrierService(ctx, cmd)
	if err != nil {
		return nil, err
	}
	updateArgs := &sqlstore.UpdateXCarrierServiceInfoArgs{
		ShopID:                            args.ShopID,
		ExternalCarrierServiceID:          res.ID,
		ExternalConnectedCarrierServiceAt: time.Now(),
	}

	if _, err := a.xAccountHaravanStore(ctx).UpdateXCarrierServiceInfo(updateArgs); err != nil {
		return nil, err
	}
	return &meta.Empty{}, nil
}

func (a *Aggregate) DeleteConnectedCarrierServiceExternalAccountHaravan(ctx context.Context, args *identity.DeleteConnectedCarrierServiceExternalAccountHaravanArgs) (*meta.Empty, error) {
	account, err := a.xAccountHaravanStore(ctx).ShopID(args.ShopID).GetXAccountHaravan()
	if err != nil {
		return nil, err
	}
	if account.ExternalCarrierServiceID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Shop chưa tạo kết nối nhà vận chuyển với Haravan")
	}

	cmd := &haravanclient.DeleteConnectedCarrierServiceRequest{
		Connection: haravanclient.Connection{
			Subdomain: account.Subdomain,
			TokenStr:  account.AccessToken,
		},
		CarrierServiceID: account.ExternalCarrierServiceID,
	}
	// ignore error
	_ = a.haravanClient.DeleteConnectedCarrierService(ctx, cmd)

	updateArgs := &sqlstore.UpdateDeleteConnectedXCarrierSeriveArgs{
		ShopID: args.ShopID,
	}
	if err := a.xAccountHaravanStore(ctx).UpdateDeleteConnectedXCarrierService(updateArgs); err != nil {
		return nil, err
	}

	return &meta.Empty{}, nil
}

func (a *Aggregate) GetCurrentCarrierServiceID(ctx context.Context, shopID dot.ID) (int, error) {
	account, err := a.xAccountHaravanStore(ctx).ShopID(shopID).GetXAccountHaravan()
	if err != nil {
		return 0, err
	}

	query := &haravanclient.GetCarrierServicesRequest{
		Connection: haravanclient.Connection{
			Subdomain: account.Subdomain,
			TokenStr:  account.AccessToken,
		},
	}
	carrierServices, err := a.haravanClient.GetCarrierServices(ctx, query)
	if err != nil {
		return 0, nil
	}
	var currentCarrierServiceID int
	for _, cs := range carrierServices {
		if cs.CarrierName == CarrierServiceName {
			currentCarrierServiceID = cs.ID
			break
		}
	}
	if currentCarrierServiceID == 0 {
		return 0, cm.Errorf(cm.NotFound, nil, "Shop chưa kết nối vận chuyển với Haravan")
	}
	return currentCarrierServiceID, nil
}
