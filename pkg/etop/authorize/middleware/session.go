package middleware

import (
	"context"
	"net/http"
	"time"

	"o.o/api/main/identity"
	"o.o/api/top/types/etc/account_type"
	identityconvert "o.o/backend/com/main/identity/convert"
	identitymodel "o.o/backend/com/main/identity/model"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/authorization/auth"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/headers"
	"o.o/backend/pkg/etop/authorize/authkey"
	"o.o/backend/pkg/etop/authorize/claims"
	"o.o/backend/pkg/etop/authorize/permission"
	"o.o/backend/pkg/etop/authorize/tokens"
	"o.o/backend/pkg/etop/model"
	"o.o/capi/dot"
	"o.o/common/l"
)

var ll = l.New()
var sadminToken string
var identityQS identity.QueryBus
var tokenStore tokens.TokenStore

type Middleware struct{}
type SAdminToken string

func New(token SAdminToken, _tokenStore tokens.TokenStore, identityQuery identity.QueryBus) Middleware {
	sadminToken = string(token)
	identityQS = identityQuery
	tokenStore = _tokenStore
	return Middleware{}
}

// StartSessionQuery ...
type StartSessionQuery struct {
	Token   string
	Request *http.Request

	RequireAuth              bool
	RequireUser              bool
	RequireAPIKey            bool
	RequireAPIPartnerShopKey bool

	RequirePartner    bool
	RequireShop       bool
	RequireAffiliate  bool
	RequireEtopAdmin  bool
	RequireSuperAdmin bool

	AuthPartner permission.AuthOpt

	Result *Session
}

type Session struct {
	User       *identitymodelx.SignedInUser
	Admin      *identitymodelx.SignedInUser
	Claim      *claims.Claim
	Partner    *identitymodel.Partner
	CtxPartner *identitymodel.Partner
	Shop       *identitymodel.Shop
	Affiliate  *identitymodel.Affiliate
	identitymodel.Permission

	IsOwner      bool
	IsEtopAdmin  bool
	IsSuperAdmin bool
}

func (s *Session) GetUserID() dot.ID {
	if s.Admin != nil {
		return s.Admin.ID
	}
	if s.User != nil {
		return s.User.ID
	}
	return 0
}

func getToken(ctx context.Context, q *StartSessionQuery) string {
	if q.Token != "" {
		return q.Token
	}
	if q.Request != nil {
		token, _ := auth.FromHTTPHeader(q.Request.Header)
		return token
	}
	return headers.GetBearerTokenFromCtx(ctx)
}

// StartSession ...
func StartSession(ctx context.Context, q *StartSessionQuery) (newCtx context.Context, _err error) {
	token := getToken(ctx, q)
	return StartSessionWithToken(ctx, token, q)
}

func StartSessionWithToken(ctx context.Context, token string, q *StartSessionQuery) (newCtx context.Context, _err error) {
	var wlPartnerID dot.ID
	defer func() {
		if wlPartnerID == 0 {
			newCtx = wl.WrapContext(ctx, 0)
		}
	}()

	// TODO: check UserID, ShopID, etc. correctly. Because InitSession now
	// responses token without any credential.
	if !q.RequireAuth {
		if token != "" {
			claim, err := tokenStore.Validate(token)
			if err != nil {
				return ctx, cm.ErrUnauthenticated
			}
			q.Result = &Session{
				Claim: claim,
			}
		}
		return ctx, nil
	}

	if token == "" {
		return ctx, cm.ErrUnauthenticated
	}

	session := new(Session)
	q.Result = session
	if q.RequireSuperAdmin {
		if token == sadminToken {
			session.IsSuperAdmin = true
			return ctx, nil
		}
		ll.Error("Invalid sadmin token")
		return ctx, cm.ErrPermissionDenied
	}

	var claim *claims.Claim
	var err error
	var account identitymodel.AccountInterface
	if q.RequireAPIKey {
		var expectType account_type.AccountType
		switch {
		case q.RequireShop:
			expectType = account_type.Shop
		case q.RequirePartner:
			expectType = account_type.Partner
		default:
			ll.Panic("unexpected account type")
		}
		claim, account, err = VerifyAPIKey(ctx, token, expectType)
		if err != nil {
			return ctx, err
		}
		_acc := account.GetAccount()
		if _acc.Type != expectType {
			return ctx, cm.Errorf(cm.PermissionDenied, nil, "")
		}
		wlPartnerID = _acc.ID

	} else if q.RequireAPIPartnerShopKey {
		claim, account, err = VerifyAPIPartnerShopKey(ctx, token)
		if err != nil {
			return ctx, err
		}
		wlPartnerID = claim.AuthPartnerID

	} else {
		claim, err = tokenStore.Validate(token)
		if err != nil {
			return ctx, cm.ErrUnauthenticated
		}
		if q.AuthPartner != 0 {
			wlPartnerID = claim.AuthPartnerID
		}
	}
	ctx = wl.WrapContext(ctx, wlPartnerID)

	if claim.STokenExpiresAt != nil && claim.STokenExpiresAt.Before(time.Now()) {
		// Invalidate stoken
		claim.SToken = false
		claim.STokenExpiresAt = nil
	}

	if claim.AdminID != 0 {
		query := &identitymodelx.GetSignedInUserQuery{
			UserID: claim.AdminID,
		}
		if err := bus.Dispatch(ctx, query); err != nil {
			ll.Error("Invalid AdminID", l.Error(err))
			return ctx, nil
		}
		session.Admin = query.Result
	}

	session.Claim = claim
	ok := StartSessionUser(ctx, q.RequireUser, claim, &session.User) &&
		StartSessionPartner(ctx, q.RequirePartner, claim, account, &session.Partner) &&
		StartSessionShop(ctx, q.RequireShop, claim, account, &session.Shop, &session.Permission) &&
		StartSessionAffiliate(ctx, q.RequireAffiliate, claim, account, &session.Affiliate, &session.Permission) &&
		StartSessionEtopAdmin(ctx, q.RequireEtopAdmin, claim, &session.Permission) &&
		StartSessionAuthPartner(ctx, q.AuthPartner, claim, &session.CtxPartner)
	if !ok {
		return ctx, cm.ErrPermissionDenied
	}
	if account != nil {
		session.IsOwner = account.GetAccount().OwnerID == claim.UserID
	}
	q.Result = session
	q.Result.Shop.BankAccount = nil
	return ctx, nil
}

func StartSessionUser(ctx context.Context, require bool, claim *claims.Claim, user **identitymodelx.SignedInUser) bool {
	if require {
		if claim.UserID == 0 {
			return false
		}
		query := &identitymodelx.GetSignedInUserQuery{
			UserID: claim.UserID,
		}
		if err := bus.Dispatch(ctx, query); err != nil {
			ll.Error("Invalid UserID", l.Error(err))
			return false
		}
		*user = query.Result
	}
	return true
}

func StartSessionAuthPartner(ctx context.Context, authOpt permission.AuthOpt, claim *claims.Claim, ctxPartner **identitymodel.Partner) bool {
	if authOpt == 0 {
		return true
	}
	if authOpt == permission.Required && claim.AuthPartnerID == 0 {
		return false
	}

	var partner *identitymodel.Partner
	partnerID := claim.AuthPartnerID
	if partnerID != 0 {
		query := &identitymodelx.GetPartner{
			PartnerID: partnerID,
		}
		if err := bus.Dispatch(ctx, query); err != nil {
			ll.Error("Invalid PartnerID", l.Error(err))
			return false
		}
		partner = query.Result.Partner
	}

	*ctxPartner = partner
	return true
}

func StartSessionPartner(ctx context.Context, require bool, claim *claims.Claim, account identitymodel.AccountInterface, _partner **identitymodel.Partner) bool {
	if partner, ok := account.(*identitymodel.Partner); ok && claim.AccountID == partner.ID {
		*_partner = partner
		return true
	}
	if require {
		if !model.IsPartnerID(claim.AccountID) {
			return false
		}
		query := &identitymodelx.GetPartner{
			PartnerID: claim.AccountID,
		}
		if err := bus.Dispatch(ctx, query); err != nil {
			ll.Error("Invalid Name", l.Error(err))
			return false
		}
		*_partner = query.Result.Partner
	}
	return true
}

func StartSessionShop(ctx context.Context, require bool, claim *claims.Claim, account identitymodel.AccountInterface, _shop **identitymodel.Shop, _permission *identitymodel.Permission) bool {
	if shop, ok := account.(*identitymodel.Shop); ok && claim.AccountID == shop.ID {
		*_shop = shop
		return true
	}
	if require {
		if !model.IsShopID(claim.AccountID) {
			return false
		}

		if claim.UserID != 0 {
			query := &identitymodelx.GetShopWithPermissionQuery{
				ShopID: claim.AccountID,
				UserID: claim.UserID,
			}
			if err := bus.Dispatch(ctx, query); err != nil {
				ll.Error("Invalid Name", l.Error(err))
				return false
			}

			*_shop = query.Result.Shop
			*_permission = query.Result.Permission
		} else {
			query := &identitymodelx.GetShopQuery{
				ShopID: claim.AccountID,
			}
			if err := bus.Dispatch(ctx, query); err != nil {
				ll.Error("Invalid Name", l.Error(err))
				return false
			}
			*_shop = query.Result
		}
	}
	return true
}

func StartSessionAffiliate(ctx context.Context, require bool, claim *claims.Claim, account identitymodel.AccountInterface, _affiliate **identitymodel.Affiliate, _permission *identitymodel.Permission) bool {
	if affiliate, ok := account.(*identitymodel.Affiliate); ok && claim.AccountID == affiliate.ID {
		*_affiliate = affiliate
		return true
	}
	if require {
		if !model.IsAffiliateID(claim.AccountID) {
			return false
		}

		if claim.UserID != 0 {
			query := &identity.GetAffiliateWithPermissionQuery{
				AffiliateID: claim.AccountID,
				UserID:      claim.UserID,
			}
			if err := identityQS.Dispatch(ctx, query); err != nil {
				ll.Error("Invalid Name", l.Error(err))
				return false
			}

			*_affiliate = identityconvert.AffiliateDB(query.Result.Affiliate)
			*_permission = identityconvert.PermissionToModel(query.Result.Permission)
		} else {
			query := &identity.GetAffiliateByIDQuery{
				ID: claim.AccountID,
			}
			if err := identityQS.Dispatch(ctx, query); err != nil {
				ll.Error("Invalid Name", l.Error(err))
				return false
			}
			*_affiliate = identityconvert.AffiliateDB(query.Result)
		}
	}
	return true
}

func StartSessionEtopAdmin(ctx context.Context, require bool, claim *claims.Claim, _permission *identitymodel.Permission) bool {
	if require {
		if !model.IsEtopAccountID(claim.AccountID) {
			return false
		}
		query := &identitymodelx.GetAccountRolesQuery{
			AccountID: model.EtopAccountID,
			UserID:    claim.UserID,
		}
		if err := bus.Dispatch(ctx, query); err != nil {
			ll.Error("Invalid GetAccountRolesQuery", l.Error(err))
			return false
		}
		*_permission = query.Result.AccountUser.Permission
	}
	return true
}

func VerifyAPIKey(ctx context.Context, apikey string, expectType account_type.AccountType) (*claims.Claim, identitymodel.AccountInterface, error) {
	info, ok := authkey.ValidateAuthKeyWithType(authkey.TypeAPIKey, apikey)
	if !ok {
		return nil, nil, cm.Error(cm.Unauthenticated, "api_key không hợp lệ", nil)
	}

	query := &identitymodelx.GetAccountAuthQuery{
		AuthKey:     apikey,
		AccountType: expectType,
		AccountID:   info.AccountID,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return nil, nil, cm.MapError(err).
			Map(cm.NotFound, cm.Unauthenticated, "api_key không hợp lệ").
			Throw()
	}

	switch expectType {
	case account_type.Partner,
		account_type.Carrier:
		partner := query.Result.Account.(*identitymodel.Partner)
		claim := &claims.Claim{
			ClaimInfo: claims.ClaimInfo{
				Token:           apikey,
				UserID:          0,
				AdminID:         0,
				AccountID:       partner.ID,
				SToken:          false,
				AccountIDs:      nil,
				STokenExpiresAt: nil,
			},
			LastLoginAt: time.Time{},
		}
		return claim, partner, nil

	case account_type.Shop:
		shop := query.Result.Account.(*identitymodel.Shop)
		claim := &claims.Claim{
			ClaimInfo: claims.ClaimInfo{
				Token:           apikey,
				UserID:          0,
				AdminID:         0,
				AccountID:       shop.ID,
				SToken:          false,
				AccountIDs:      nil,
				STokenExpiresAt: nil,
			},
			LastLoginAt: time.Time{},
		}
		return claim, shop, nil
	}
	panic("unexpected")
}

func VerifyAPIPartnerShopKey(ctx context.Context, apikey string) (*claims.Claim, identitymodel.AccountInterface, error) {
	_, ok := authkey.ValidateAuthKeyWithType(authkey.TypePartnerShopKey, apikey)
	if !ok {
		return nil, nil, cm.Error(cm.Unauthenticated, "api_key không hợp lệ", nil)
	}

	relationQuery := &identitymodelx.GetPartnerRelationQuery{
		AuthKey: apikey,
	}
	relationError := bus.Dispatch(ctx, relationQuery)
	if relationError != nil {
		return nil, nil, cm.MapError(relationError).
			Map(cm.NotFound, cm.PermissionDenied, "").
			Throw()
	}

	partnerID := relationQuery.Result.PartnerID
	shop := relationQuery.Result.Shop

	claim := &claims.Claim{
		ClaimInfo: claims.ClaimInfo{
			Token:           apikey,
			UserID:          0,
			AdminID:         0,
			AccountID:       shop.ID,
			AuthPartnerID:   partnerID,
			SToken:          false,
			AccountIDs:      nil,
			STokenExpiresAt: nil,
		},
		LastLoginAt: time.Time{},
	}
	return claim, shop, nil
}
