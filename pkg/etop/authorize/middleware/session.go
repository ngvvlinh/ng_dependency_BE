package middleware

import (
	"context"
	"net/http"
	"time"

	"etop.vn/api/main/identity"
	"etop.vn/api/top/types/etc/account_type"
	identityconvert "etop.vn/backend/com/main/identity/convert"
	identitymodel "etop.vn/backend/com/main/identity/model"
	identitymodelx "etop.vn/backend/com/main/identity/modelx"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/whitelabel/wl"
	"etop.vn/backend/pkg/common/authorization/auth"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/headers"
	"etop.vn/backend/pkg/etop/authorize/authkey"
	"etop.vn/backend/pkg/etop/authorize/claims"
	"etop.vn/backend/pkg/etop/authorize/permission"
	"etop.vn/backend/pkg/etop/authorize/tokens"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
	"etop.vn/common/l"
)

var ll = l.New()
var sadminToken string
var identityQS identity.QueryBus

func Init(token string, identityQuery identity.QueryBus) {
	sadminToken = token
	identityQS = identityQuery
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
	var wlPartnerID dot.ID
	defer func() {
		if _err == nil && wlPartnerID == 0 {
			newCtx = wl.WrapContext(ctx, 0)
		}
	}()

	// TODO: check UserID, ShopID, etc. correctly. Because InitSession now
	// responses token without any credential.
	if !q.RequireAuth {
		token := getToken(ctx, q)
		if token != "" {
			claim, err := tokens.Store.Validate(token)
			if err != nil {
				return ctx, cm.ErrUnauthenticated
			}
			q.Result = &Session{
				Claim: claim,
			}
		}
		return ctx, nil
	}

	token := getToken(ctx, q)
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
		claim, account, err = verifyAPIKey(ctx, token, expectType)
		if err != nil {
			return ctx, err
		}
		wlPartnerID = account.GetAccount().ID

	} else if q.RequireAPIPartnerShopKey {
		claim, account, err = verifyAPIPartnerShopKey(ctx, token)
		if err != nil {
			return ctx, err
		}
		wlPartnerID = claim.AuthPartnerID

	} else {
		claim, err = tokens.Store.Validate(token)
		if err != nil {
			return ctx, cm.ErrUnauthenticated
		}
	}
	ctx = wl.WrapContext(ctx, wlPartnerID)

	if claim.STokenExpiresAt != nil && claim.STokenExpiresAt.Before(time.Now()) {
		// Invalidate stoken
		claim.SToken = false
		claim.STokenExpiresAt = nil
	}

	if claim.AdminID != 0 {
		query := &identitymodelx.GetSignedInUserQuery{UserID: claim.AdminID}
		if err := bus.Dispatch(ctx, query); err != nil {
			ll.Error("Invalid AdminID", l.Error(err))
			return ctx, nil
		}
		session.Admin = query.Result
	}

	session.Claim = claim
	ok := startSessionUser(ctx, q.RequireUser, session) &&
		startSessionPartner(ctx, q.RequirePartner, session, account) &&
		startSessionShop(ctx, q.RequireShop, session, account) &&
		startSessionAffiliate(ctx, q.RequireAffiliate, session, account) &&
		startSessionEtopAdmin(ctx, q.RequireEtopAdmin, session) &&
		startSessionAuthPartner(ctx, q.AuthPartner, session)
	if !ok {
		return ctx, cm.ErrPermissionDenied
	}
	q.Result = session
	return ctx, nil
}

func startSessionUser(ctx context.Context, require bool, s *Session) bool {
	if require {
		if s.Claim.UserID == 0 {
			return false
		}
		query := &identitymodelx.GetSignedInUserQuery{UserID: s.Claim.UserID}
		if err := bus.Dispatch(ctx, query); err != nil {
			ll.Error("Invalid UserID", l.Error(err))
			return false
		}
		s.User = query.Result
	}
	return true
}

func startSessionAuthPartner(ctx context.Context, authOpt permission.AuthOpt, s *Session) bool {
	if authOpt == 0 {
		return true
	}
	if authOpt == permission.Required && s.Claim.AuthPartnerID == 0 {
		return false
	}

	var partner *identitymodel.Partner
	partnerID := s.Claim.AuthPartnerID
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

	s.CtxPartner = partner
	return true
}

func startSessionPartner(ctx context.Context, require bool, s *Session, account identitymodel.AccountInterface) bool {
	if partner, ok := account.(*identitymodel.Partner); ok && s.Claim.AccountID == partner.ID {
		s.Partner = partner
		return true
	}
	if require {
		if !model.IsPartnerID(s.Claim.AccountID) {
			return false
		}
		query := &identitymodelx.GetPartner{
			PartnerID: s.Claim.AccountID,
		}
		if err := bus.Dispatch(ctx, query); err != nil {
			ll.Error("Invalid Name", l.Error(err))
			return false
		}
		s.Partner = query.Result.Partner
	}
	return true
}

func startSessionShop(ctx context.Context, require bool, s *Session, account identitymodel.AccountInterface) bool {
	if shop, ok := account.(*identitymodel.Shop); ok && s.Claim.AccountID == shop.ID {
		s.Shop = shop
		return true
	}
	if require {
		if !model.IsShopID(s.Claim.AccountID) {
			return false
		}

		if s.Claim.UserID != 0 {
			query := &identitymodelx.GetShopWithPermissionQuery{
				ShopID: s.Claim.AccountID,
				UserID: s.Claim.UserID,
			}
			if err := bus.Dispatch(ctx, query); err != nil {
				ll.Error("Invalid Name", l.Error(err))
				return false
			}

			s.Shop = query.Result.Shop
			s.Permission = query.Result.Permission
			s.IsOwner = s.Shop.OwnerID == s.Claim.UserID
		} else {
			query := &identitymodelx.GetShopQuery{
				ShopID: s.Claim.AccountID,
			}
			if err := bus.Dispatch(ctx, query); err != nil {
				ll.Error("Invalid Name", l.Error(err))
				return false
			}
			s.Shop = query.Result
		}
	}
	return true
}

func startSessionAffiliate(ctx context.Context, require bool, s *Session, account identitymodel.AccountInterface) bool {
	if affiliate, ok := account.(*identitymodel.Affiliate); ok && s.Claim.AccountID == affiliate.ID {
		s.Affiliate = affiliate
		return true
	}
	if require {
		if !model.IsAffiliateID(s.Claim.AccountID) {
			return false
		}

		if s.Claim.UserID != 0 {
			query := &identity.GetAffiliateWithPermissionQuery{
				AffiliateID: s.Claim.AccountID,
				UserID:      s.Claim.UserID,
			}
			if err := identityQS.Dispatch(ctx, query); err != nil {
				ll.Error("Invalid Name", l.Error(err))
				return false
			}

			s.Affiliate = identityconvert.AffiliateDB(query.Result.Affiliate)
			s.Permission = identityconvert.PermissionToModel(query.Result.Permission)
			s.IsOwner = s.Affiliate.OwnerID == s.Claim.UserID
		} else {
			query := &identity.GetAffiliateByIDQuery{
				ID: s.Claim.AccountID,
			}
			if err := identityQS.Dispatch(ctx, query); err != nil {
				ll.Error("Invalid Name", l.Error(err))
				return false
			}
			s.Affiliate = identityconvert.AffiliateDB(query.Result)
		}
	}
	return true
}

func startSessionEtopAdmin(ctx context.Context, require bool, s *Session) bool {
	if require {
		if !model.IsEtopAccountID(s.Claim.AccountID) {
			return false
		}
		query := &identitymodelx.GetAccountRolesQuery{
			AccountID: model.EtopAccountID,
			UserID:    s.Claim.UserID,
		}
		if err := bus.Dispatch(ctx, query); err != nil {
			ll.Error("Invalid GetAccountRolesQuery", l.Error(err))
			return false
		}
		s.Permission = query.Result.AccountUser.Permission
		s.IsEtopAdmin = len(s.Roles) > 0 || len(s.Permissions) > 0
	}
	return true
}

func verifyAPIKey(ctx context.Context, apikey string, expectType account_type.AccountType) (*claims.Claim, identitymodel.AccountInterface, error) {
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
	case account_type.Partner:
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

func verifyAPIPartnerShopKey(ctx context.Context, apikey string) (*claims.Claim, identitymodel.AccountInterface, error) {
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
