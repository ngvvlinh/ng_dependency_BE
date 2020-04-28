package partner

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"o.o/api/main/catalog"
	"o.o/api/main/connectioning"
	"o.o/api/main/inventory"
	"o.o/api/main/location"
	"o.o/api/main/shipping"
	"o.o/api/shopping/addressing"
	"o.o/api/shopping/customering"
	extpartner "o.o/api/top/external/partner"
	pbcm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/authorize_shop_config"
	"o.o/api/top/types/etc/status3"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/common/apifw/idemp"
	cmService "o.o/backend/pkg/common/apifw/service"
	"o.o/backend/pkg/common/apifw/whitelabel"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/authorization/auth"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/common/validate"
	apiconvertpb "o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/apix/convertpb"
	"o.o/capi/dot"
	"o.o/common/l"
)

var (
	idempgroup          *idemp.RedisGroup
	authStore           auth.Generator
	authURL             string
	locationQuery       location.QueryBus
	customerQuery       *customering.QueryBus
	customerAggregate   *customering.CommandBus
	addressQuery        *addressing.QueryBus
	addressAggregate    *addressing.CommandBus
	inventoryQuery      *inventory.QueryBus
	catalogAggregate    *catalog.CommandBus
	catalogQuery        *catalog.QueryBus
	connectionQuery     connectioning.QueryBus
	connectionAggregate connectioning.CommandBus
	shippingAggregate   shipping.CommandBus

	ll = l.New()
)

const PrefixIdempPartnerAPI = "IdempPartnerAPI"

const ttlShopRequest = 15 * 60 // 15 minutes
const msgShopRequest = `Sử dụng mã này để hỏi quyền tạo đơn hàng với tư cách shop (có hiệu lực trong 15 phút)`
const msgShopKey = `Sử dụng mã này để tạo đơn hàng với tư cách shop (có hiệu lực khi shop vẫn tiếp tục sử dụng dịch vụ của đối tác)`
const msgUserKey = `Sử dụng mã này để truy cập hệ thống với tư cách user (có hiệu lực khi user vẫn tiếp tục sử dụng dịch vụ của đối tác)`

func init() {
	bus.AddHandlers("apix",
		miscService.VersionInfo,
		miscService.CurrentAccount,
		shopService.CurrentShop,
		shopService.AuthorizeShop,
	)
}

type MiscService struct{}
type ShopService struct{}
type ShipmentConnectionService struct{}
type ShipmentService struct{}
type WebhookService struct{}
type HistoryService struct{}

type ShippingService struct{}
type CustomerService struct{}
type CustomerAddressService struct{}
type CustomerGroupService struct{}
type CustomerGroupRelationshipService struct{}
type InventoryService struct{}
type OrderService struct{}
type FulfillmentService struct{}
type ProductService struct{}
type ProductCollectionService struct{}
type ProductCollectionRelationshipService struct{}
type VariantService struct{}

var miscService = &MiscService{}
var shopService = &ShopService{}
var webhookService = &WebhookService{}
var historyService = &HistoryService{}
var shippingService = &ShippingService{}
var customerService = &CustomerService{}
var customerAddressService = &CustomerAddressService{}
var customerGroupService = &CustomerGroupService{}
var customerGroupRelationshipService = &CustomerGroupRelationshipService{}
var inventoryService = &InventoryService{}
var orderService = &OrderService{}
var fulfillmentService = &FulfillmentService{}
var productService = &ProductService{}
var productCollectionService = &ProductCollectionService{}
var productCollectionRelationshipService = &ProductCollectionRelationshipService{}
var variantService = &VariantService{}
var shipmentConnectionService = &ShipmentConnectionService{}
var shipmentService = &ShipmentService{}

func Init(
	sd cmService.Shutdowner,
	rd redis.Store,
	s auth.Generator,
	_authURL string,
	locationQ location.QueryBus,
	customerQ *customering.QueryBus,
	customerA *customering.CommandBus,
	addressQ *addressing.QueryBus,
	addressA *addressing.CommandBus,
	inventoryQ *inventory.QueryBus,
	catalogQ *catalog.QueryBus,
	catalogA *catalog.CommandBus,
	connectionQ connectioning.QueryBus,
	connectionA connectioning.CommandBus,
	shippingAggr shipping.CommandBus,
) {
	if _authURL == "" {
		ll.Panic("no auth_url")
	}
	if _, err := url.Parse(_authURL); err != nil {
		ll.Panic("invalid auth_url", l.String("url", _authURL))
	}

	authStore = s
	authURL = _authURL

	idempgroup = idemp.NewRedisGroup(rd, PrefixIdempPartnerAPI, 0)
	sd.Register(idempgroup.Shutdown)

	locationQuery = locationQ
	customerQuery = customerQ
	customerAggregate = customerA
	addressQuery = addressQ
	addressAggregate = addressA
	inventoryQuery = inventoryQ
	catalogAggregate = catalogA
	catalogQuery = catalogQ
	connectionQuery = connectionQ
	connectionAggregate = connectionA
	shippingAggregate = shippingAggr
}

func (s *MiscService) VersionInfo(ctx context.Context, q *VersionInfoEndpoint) error {
	q.Result = &pbcm.VersionInfoResponse{
		Service: "partner",
		Version: "1.0.0",
	}
	return nil
}

func (s *MiscService) CurrentAccount(ctx context.Context, q *CurrentAccountEndpoint) error {
	if q.Context.Partner == nil {
		return cm.Errorf(cm.Internal, nil, "")
	}
	q.Result = convertpb.PbPartner(q.Context.Partner)
	if wl.X(ctx).IsWhiteLabel() {
		q.Result.Meta = map[string]string{
			"wl_name": wl.X(ctx).Name,
			"wl_key":  wl.X(ctx).Key,
			"wl_host": wl.X(ctx).Host,
		}
	}
	return nil
}

func (s *ShopService) CurrentShop(ctx context.Context, q *CurrentShopEndpoint) error {
	if q.Context.Shop == nil {
		return cm.Errorf(cm.Internal, nil, "")
	}
	q.Result = apiconvertpb.PbPublicAccountInfo(q.Context.Shop)
	return nil
}

func getAuthorizeShopConfig(configs []authorize_shop_config.AuthorizeShopConfig) string {
	var res []string
	for _, c := range configs {
		res = append(res, c.String())
	}
	return strings.Join(res, ",")
}

func (s *ShopService) AuthorizeShop(ctx context.Context, q *AuthorizeShopEndpoint) error {
	partner := q.Context.Partner
	if q.RedirectUrl != "" {
		if err := validateRedirectURL(partner.RedirectURLs, q.RedirectUrl, true); err != nil {
			return err
		}
	}

	whiteLabelData, err := validateWhiteLabel(ctx, q)
	if err != nil {
		return err
	}

	// always verify email and phone
	var emailNorm, phoneNorm string
	if q.Email != "" {
		email, ok := validate.NormalizeEmail(q.Email)
		if !ok {
			return cm.Errorf(cm.InvalidArgument, nil, `Giá trị email không hợp lệ`)
		}
		emailNorm = email.String()
	}
	if q.Phone != "" {
		phone, ok := validate.NormalizePhone(q.Phone)
		if !ok {
			return cm.Errorf(cm.InvalidArgument, nil, `Giá trị phone không hợp lệ`)
		}
		phoneNorm = phone.String()
	}
	_, _ = emailNorm, phoneNorm // temporary ignore error here

	// verify external code
	if q.ExternalShopID != "" && !validate.ExternalCode(q.ExternalShopID) {
		return cm.Errorf(cm.InvalidArgument, nil, "Giá trị external_shop_id không hợp lệ")
	}
	if q.ExternalUserID != "" && !validate.ExternalCode(q.ExternalUserID) {
		return cm.Errorf(cm.InvalidArgument, nil, "Giá trị external_user_id không hợp lệ")
	}

	// case 1: if the shop has linked to the partner
	if q.ShopId != 0 && q.ExternalShopID != "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Chỉ cần cung cấp shop_id hoặc external_shop_id")
	}

	if q.ExternalUserID != "" {
		if !wl.X(ctx).IsWhiteLabel() || whiteLabelData == nil {
			// only whitelabel partner can provide this field,
			// we fake the error message here
			return cm.Errorf(cm.InvalidArgument, nil, "external_user_id is removed")
		}
		if q.ExtraToken != "" || q.ExternalShopID == "" {
			return handleWLAuthorizeShopByExternalUserID(ctx, q)
		}
	}

	if q.ShopId != 0 || q.ExternalShopID != "" {
		if q.ShopId != 0 && (q.Email != "" || q.Phone != "") {
			return cm.Errorf(cm.InvalidArgument, nil, "Nếu cung cấp shop_id thì không cần kèm theo email hoặc phone. Nếu cung cấp email hoặc phone thì không cần kèm theo shop_id.")
		}

		relationQuery := &identitymodelx.GetPartnerRelationQuery{
			PartnerID:         partner.ID,
			AccountID:         q.ShopId,
			ExternalAccountID: q.ExternalShopID,
		}
		err := bus.Dispatch(ctx, relationQuery)
		switch {
		case err == nil:
			rel := relationQuery.Result.PartnerRelation
			shop := relationQuery.Result.Shop
			user := relationQuery.Result.User
			if q.ExternalUserID != "" {
				relationUserQuery := &identitymodelx.GetPartnerRelationQuery{
					PartnerID:      partner.ID,
					ExternalUserID: q.ExternalUserID,
				}
				if err2 := bus.Dispatch(ctx, relationUserQuery); err2 == nil {
					userID := relationUserQuery.Result.SubjectID
					if userID != 0 && userID != user.ID {
						return cm.Errorf(cm.FailedPrecondition, nil, "external_shop_id (id=%v) does not belong to external_user_id (id=%v).", q.ExternalShopID, q.ExternalUserID)
					}
				}
			}

			info := PartnerShopToken{
				PartnerID:      partner.ID,
				ShopID:         shop.ID,
				ShopName:       shop.Name,
				ShopOwnerEmail: user.Email,
				ShopOwnerPhone: user.Phone,
				ShopOwnerName:  user.FullName,
				ExternalShopID: q.ExternalShopID,
				ExternalUserID: q.ExternalUserID,

				// client must keep the current email/phone when calling
				// request_login
				RetainCurrentInfo: true,
				Config:            getAuthorizeShopConfig(q.Config),
				RedirectURL:       q.RedirectUrl,
				AuthType:          AuthTypeShopKey,
			}
			if whiteLabelData != nil {
				// Trường hợp whiteLabel
				// Đã có sẵn tài khoản sẽ redirect về trang whiteLabel
				info.RedirectURL = whiteLabelData.Config.RootURL
			}
			token, err := generateAuthToken(info)
			if err != nil {
				return err
			}

			if rel.Status == status3.P && rel.DeletedAt.IsZero() &&
				shop.Status == status3.P && shop.DeletedAt.IsZero() &&
				user.Status == status3.P {
				// TODO: handle config: "shipment"
				q.Result = &extpartner.AuthorizeShopResponse{
					Code:      "ok",
					Msg:       msgShopKey,
					Type:      "shop_key",
					AuthToken: rel.AuthKey,
					ExpiresIn: -1,
					AuthUrl:   generateAuthURL(whiteLabelData, authURL, token.TokenStr),
				}
				return nil
			}
			if shop.Status != status3.P || !shop.DeletedAt.IsZero() ||
				user.Status != status3.P {
				return cm.Errorf(cm.AccountClosed, nil, "")
			}
			if rel.Status != status3.P || !rel.DeletedAt.IsZero() {
				q.Result = &extpartner.AuthorizeShopResponse{
					Code:      "ok",
					Msg:       msgShopRequest,
					Type:      "shop_request",
					AuthToken: token.TokenStr,
					ExpiresIn: token.ExpiresIn,
				}
				if q.RedirectUrl != "" {
					q.Result.AuthUrl = generateAuthURL(whiteLabelData, authURL, q.Result.AuthToken)
				}
				return nil
			}

		case cm.ErrorCode(err) == cm.NotFound:
			if q.ShopId != 0 {
				return cm.Errorf(cm.PermissionDenied, nil, "").
					WithMeta("reason", "Chỉ có thể sử dụng shop_id nếu shop đã từng đăng nhập qua hệ thống của đối tác")
			}
			return generateAuthTokenWithRequestLogin(ctx, q, 0)

		default:
			return cm.Errorf(cm.Internal, err, "")
		}

		// prevent unexpected condition
		return cm.Errorf(cm.Internal, nil, "").
			WithMeta("reason", "unexpected condition")
	}

	// case 2: if the user/shop has not linked to the partner
	return generateAuthTokenWithRequestLogin(ctx, q, 0)
}

func handleWLAuthorizeShopByExternalUserID(ctx context.Context, q *AuthorizeShopEndpoint) error {
	if !wl.X(ctx).IsWhiteLabel() {
		return cm.Errorf(cm.FailedPrecondition, nil, "Only whitelabel partner can use this function")
	}

	wlPartner := wl.X(ctx)
	redirectURL := ""
	switch {
	case q.ExternalUserID != "" && q.ExtraToken != "":
		// TH invitation
		// check if invitation
		if err := checkExtraTokenInvitation(q.ExtraToken); err != nil {
			return err
		}
		redirectURL = fmt.Sprintf("%v?t=%v", wlPartner.InviteUserURLByEmail, q.ExtraToken)
	case q.ExternalUserID != "" && q.ExternalShopID == "":
		// TH login bằng external_user_id
		redirectURL = wlPartner.Config.RootURL
	default:
		return cm.Errorf(cm.Internal, nil, "").WithMeta("reason", "unexpected condition")
	}

	// try retrieving the partner relation
	relationQuery := &identitymodelx.GetPartnerRelationQuery{
		PartnerID:      wlPartner.ID,
		ExternalUserID: q.ExternalUserID,
	}
	err := bus.Dispatch(ctx, relationQuery)
	switch cm.ErrorCode(err) {
	case cm.NoError:
		// generate user_key and redirect to whiteLabelData.Config.RootURL
		rel := relationQuery.Result.PartnerRelation
		user := relationQuery.Result.User
		info := PartnerShopToken{
			PartnerID:      wlPartner.ID,
			ShopOwnerEmail: user.Email,
			ShopOwnerPhone: user.Phone,
			ShopOwnerName:  user.FullName,
			ExternalUserID: q.ExternalUserID,
			ExtraToken:     q.ExtraToken,
			// client must keep the current email/phone when calling
			// request_login
			RetainCurrentInfo: true,
			Config:            getAuthorizeShopConfig(q.Config),
			RedirectURL:       redirectURL,
			AuthType:          AuthTypeUserKey,
		}
		token, err := generateAuthToken(info)
		if err != nil {
			return err
		}

		if user.Status != status3.P {
			return cm.Errorf(cm.AccountClosed, nil, "")
		}
		// TODO: handle config: "shipment"
		q.Result = &extpartner.AuthorizeShopResponse{
			Code:      "ok",
			Msg:       msgUserKey,
			Type:      "user_key",
			AuthToken: rel.AuthKey,
			ExpiresIn: -1,
			AuthUrl:   generateAuthURL(wlPartner, authURL, token.TokenStr),
		}
		return nil

	case cm.NotFound:
		// chỉ cho phép pass bước này trong trường hợp invitation
		if _err := checkExtraTokenInvitation(q.ExtraToken); _err != nil {
			return cm.Errorf(cm.InvalidArgument, err, "Request không hợp lệ.")
		}

		// continue to case 2 (the user/shop has not linked to the partner)
		return generateAuthTokenWithRequestLogin(ctx, q, 0)
	default:
		return err
	}
}

func generateAuthTokenWithRequestLogin(ctx context.Context, q *AuthorizeShopEndpoint, shopID dot.ID) error {
	whiteLabelData, err := validateWhiteLabel(ctx, q)
	if err != nil {
		return err
	}

	info := PartnerShopToken{
		PartnerID: q.Context.Partner.ID,

		// leave this field empty because we don't want to expose our account
		// information to partner
		ShopID:            0,
		ShopName:          q.ShopName,
		ShopOwnerEmail:    q.Email,
		ShopOwnerPhone:    q.Phone,
		ShopOwnerName:     q.Name,
		ExternalShopID:    q.ExternalShopID,
		ExternalUserID:    q.ExternalUserID,
		ExtraToken:        q.ExtraToken,
		RetainCurrentInfo: false,
		RedirectURL:       q.RedirectUrl,
		Config:            getAuthorizeShopConfig(q.Config),
	}
	if shopID != 0 {
		info.ShopID = shopID
		info.RetainCurrentInfo = true
	}
	if whiteLabelData != nil && q.ExternalUserID != "" && q.ExtraToken != "" {
		if err := checkExtraTokenInvitation(q.ExtraToken); err != nil {
			return err
		}
		info.RedirectURL = fmt.Sprintf("%v?t=%v", whiteLabelData.InviteUserURLByEmail, q.ExtraToken)
	}

	token, err := generateAuthToken(info)
	if err != nil {
		return err
	}

	q.Result = &extpartner.AuthorizeShopResponse{
		Code:      "ok",
		Msg:       msgShopRequest,
		Type:      "shop_request",
		AuthToken: token.TokenStr,
		ExpiresIn: token.ExpiresIn,
		Meta: map[string]string{
			"recaptcha_token": "1",
		},
	}
	if q.RedirectUrl != "" {
		q.Result.AuthUrl = generateAuthURL(whiteLabelData, authURL, q.Result.AuthToken)
	}
	return nil
}

func generateAuthToken(info PartnerShopToken) (*auth.Token, error) {
	if info.PartnerID == 0 {
		return nil, cm.Errorf(cm.Internal, nil, "Missing PartnerID")
	}

	tokStr := "request:" + auth.RandomToken(auth.DefaultTokenLength)
	tok := &auth.Token{
		TokenStr: tokStr,
		Usage:    auth.UsagePartnerIntegration,
		UserID:   0,
		Value:    info,
	}
	_, err := authStore.GenerateWithValue(tok, ttlShopRequest)
	if err != nil {
		return nil, cm.Errorf(cm.Internal, err, "")
	}
	return tok, nil
}

func generateAuthURL(whitelabelData *whitelabel.WL, authURL string, token string) string {
	if whitelabelData != nil {
		authURL = whitelabelData.Config.AuthURL
	}
	u, err := url.Parse(authURL)
	if err != nil {
		ll.Panic("invalid auth_url", l.Error(err))
	}
	query := u.Query()
	query.Set("token", token)
	u.RawQuery = query.Encode()
	return u.String()
}

var reLoopback = regexp.MustCompile(`^127\.0\.0\.[0-9]{3}$`)

func validateRedirectURL(redirectURLs []string, redirectURL string, skipCheckIfNoURL bool) error {
	rURL, err := url.Parse(redirectURL)
	if err != nil {
		return cm.Errorf(cm.InvalidArgument, nil, "Địa chỉ url không hợp lệ")
	}

	// allow localhost for testing
	if rURL.Host == "localhost" || reLoopback.MatchString(rURL.Host) {
		return nil
	}

	if skipCheckIfNoURL && len(redirectURLs) == 0 {
		return nil
	}
	for _, registerURL := range redirectURLs {
		if redirectURL == registerURL {
			return nil
		}
	}
	return cm.Errorf(cm.InvalidArgument, nil, "Địa chỉ url cần được đăng ký trước")
}

func validateWhiteLabel(ctx context.Context, q *AuthorizeShopEndpoint) (wlPartner *whitelabel.WL, err error) {
	wlPartner = wl.X(ctx)
	if !wlPartner.IsWhiteLabel() {
		return nil, nil
	}

	// validate data
	fields := []cmapi.Field{
		{
			Name:  "Số điện thoại (phone)",
			Value: q.Phone,
		}, {
			Name:  "Email",
			Value: q.Email,
		}, {
			Name:  "Tên (name)",
			Value: q.Name,
		}, {
			Name:  "redirect_url",
			Value: q.RedirectUrl,
		}, {
			Name:  "Mã người dùng (external_user_id)",
			Value: q.ExternalUserID,
		},
		// NOTE(qv): external_shop_id may be empty
	}

	if q.ExternalShopID != "" {
		validateFields := []cmapi.Field{
			{
				Name:  "Tên cửa hàng (shop_name)",
				Value: q.ShopName,
			},
		}
		fields = append(fields, validateFields...)
	}

	if err := cmapi.ValidateEmptyField(fields...); err != nil {
		return nil, err
	}
	return wlPartner, nil
}

func checkExtraTokenInvitation(extraToken string) error {
	if strings.HasPrefix(extraToken, auth.UsageInviteUser+":") {
		return nil
	}
	return cm.Errorf(cm.InvalidArgument, nil, "extra_token does not valid")
}
