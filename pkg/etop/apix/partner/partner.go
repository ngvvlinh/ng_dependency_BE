package partner

import (
	"context"
	"net/url"
	"regexp"
	"strings"

	"etop.vn/api/main/catalog"
	"etop.vn/api/main/inventory"
	"etop.vn/api/main/location"
	"etop.vn/api/shopping/addressing"
	"etop.vn/api/shopping/customering"
	extpartner "etop.vn/api/top/external/partner"
	pbcm "etop.vn/api/top/types/common"
	"etop.vn/api/top/types/etc/authorize_shop_config"
	"etop.vn/api/top/types/etc/status3"
	identitymodelx "etop.vn/backend/com/main/identity/modelx"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/cmapi"
	"etop.vn/backend/pkg/common/apifw/idemp"
	cmService "etop.vn/backend/pkg/common/apifw/service"
	"etop.vn/backend/pkg/common/apifw/whitelabel"
	"etop.vn/backend/pkg/common/apifw/whitelabel/wl"
	"etop.vn/backend/pkg/common/authorization/auth"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/redis"
	"etop.vn/backend/pkg/common/validate"
	apiconvertpb "etop.vn/backend/pkg/etop/api/convertpb"
	"etop.vn/backend/pkg/etop/apix/convertpb"
	"etop.vn/capi/dot"
	"etop.vn/common/l"
)

var (
	idempgroup        *idemp.RedisGroup
	authStore         auth.Generator
	authURL           string
	locationQuery     location.QueryBus
	customerQuery     *customering.QueryBus
	customerAggregate *customering.CommandBus
	addressQuery      *addressing.QueryBus
	addressAggregate  *addressing.CommandBus
	inventoryQuery    *inventory.QueryBus
	catalogAggregate  *catalog.CommandBus
	catalogQuery      *catalog.QueryBus

	ll = l.New()
)

const PrefixIdempPartnerAPI = "IdempPartnerAPI"

const ttlShopRequest = 15 * 60 // 15 minutes
const msgShopRequest = `Sử dụng mã này để hỏi quyền tạo đơn hàng với tư cách shop (có hiệu lực trong 15 phút)`
const msgShopKey = `Sử dụng mã này để tạo đơn hàng với tư cách shop (có hiệu lực khi shop vẫn tiếp tục sử dụng dịch vụ của đối tác)`

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

	// case 1: if the shop has linked to the partner
	if q.ShopId != 0 && q.ExternalShopId != "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Chỉ cần cung cấp shop_id hoặc external_shop_id")
	}
	if q.ShopId != 0 || q.ExternalShopId != "" {
		if q.ShopId != 0 && (q.Email != "" || q.Phone != "") {
			return cm.Errorf(cm.InvalidArgument, nil, "Nếu cung cấp shop_id thì không cần kèm theo email hoặc phone. Nếu cung cấp email hoặc phone thì không cần kèm theo shop_id.")
		}
		if q.ExternalShopId != "" && !validate.ExternalCode(q.ExternalShopId) {
			return cm.Errorf(cm.InvalidArgument, nil, "Giá trị external_shop_id không hợp lệ")
		}

		relationQuery := &identitymodelx.GetPartnerRelationQuery{
			PartnerID:         partner.ID,
			AccountID:         q.ShopId,
			ExternalAccountID: q.ExternalShopId,
		}
		err := bus.Dispatch(ctx, relationQuery)
		switch {
		case err == nil:
			rel := relationQuery.Result.PartnerRelation
			shop := relationQuery.Result.Shop
			user := relationQuery.Result.User
			info := PartnerShopToken{
				PartnerID:      partner.ID,
				ShopID:         shop.ID,
				ShopName:       shop.Name,
				ShopOwnerEmail: user.Email,
				ShopOwnerPhone: user.Phone,
				ShopOwnerName:  user.FullName,
				ExternalShopID: q.ExternalShopId,

				// client must keep the current email/phone when calling
				// request_login
				RetainCurrentInfo: true,
				Config:            getAuthorizeShopConfig(q.Config),
				RedirectURL:       q.RedirectUrl,
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

	// case 2: if the shop has not linked to the partner
	return generateAuthTokenWithRequestLogin(ctx, q, 0)
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
		ShopID:         0,
		ShopName:       q.ShopName,
		ShopOwnerEmail: q.Email,
		ShopOwnerPhone: q.Phone,
		ShopOwnerName:  q.Name,
		ExternalShopID: q.ExternalShopId,

		RetainCurrentInfo: false,
		RedirectURL:       q.RedirectUrl,
		Config:            getAuthorizeShopConfig(q.Config),
	}
	if shopID != 0 {
		info.ShopID = shopID
		info.RetainCurrentInfo = true
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
	redirectUrl := authURL
	if whitelabelData != nil {
		redirectUrl = whitelabelData.Config.AuthURL
	}
	u, err := url.Parse(redirectUrl)
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

func validateWhiteLabel(ctx context.Context, q *AuthorizeShopEndpoint) (*whitelabel.WL, error) {
	wlPartner := wl.X(ctx)
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
			Name:  "Tên cửa hàng (shop_name)",
			Value: q.ShopName,
		}, {
			Name:  "Mã shop (external_shop_id)",
			Value: q.ExternalShopId,
		},
	}
	if err := cmapi.ValidateEmptyField(fields...); err != nil {
		return nil, err
	}
	return wlPartner, nil
}
