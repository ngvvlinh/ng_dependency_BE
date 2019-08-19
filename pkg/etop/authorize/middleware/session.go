package middleware

import (
	"context"
	"net/http"
	"time"

	"etop.vn/api/main/identity"
	identityconvert "etop.vn/backend/com/main/identity/convert"
	identitymodel "etop.vn/backend/com/main/identity/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/auth"
	"etop.vn/backend/pkg/etop/authorize/authkey"
	"etop.vn/backend/pkg/etop/authorize/claims"
	"etop.vn/backend/pkg/etop/authorize/permission"
	"etop.vn/backend/pkg/etop/authorize/tokens"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/common/bus"
	"etop.vn/common/l"
)

var ll = l.New()
var sadminToken string
var identityQS identity.QueryBus

func init() {
	bus.AddHandler("session", StartSession)
}

func Init(token string, identityQuery identity.QueryBus) {
	sadminToken = token
	identityQS = identityQuery
}

// StartSessionQuery ...
type StartSessionQuery struct {
	Token   string
	Context context.Context
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
	User       *model.SignedInUser
	Admin      *model.SignedInUser
	Claim      *claims.Claim
	Partner    *model.Partner
	CtxPartner *model.Partner
	Shop       *model.Shop
	Affiliate  *identitymodel.Affiliate
	model.Permission

	IsOwner      bool
	IsEtopAdmin  bool
	IsSuperAdmin bool
}

func (s *Session) GetUserID() int64 {
	if s.Admin != nil {
		return s.Admin.ID
	}
	if s.User != nil {
		return s.User.ID
	}
	return 0
}

// GetBearerTokenFromCtx ...
func GetBearerTokenFromCtx(ctx context.Context) string {
	authHeader, ok := ctx.Value(authKey{}).(string)
	if !ok {
		return ""
	}
	token, _ := auth.FromHeaderString(authHeader)
	return token
}

func getToken(ctx context.Context, q *StartSessionQuery) string {
	if q.Token != "" {
		return q.Token
	}

	if q.Request != nil {
		token, _ := auth.FromHTTPHeader(q.Request.Header)
		return token
	}
	if q.Context != nil {
		return GetBearerTokenFromCtx(q.Context)
	}
	return ""
}

// StartSession ...
func StartSession(ctx context.Context, q *StartSessionQuery) error {
	if !q.RequireAuth {
		return nil
	}

	token := getToken(ctx, q)
	if token == "" {
		return cm.ErrUnauthenticated
	}

	session := new(Session)
	q.Result = session
	if q.RequireSuperAdmin {
		if token == sadminToken {
			session.IsSuperAdmin = true
			return nil
		}
		ll.Error("Invalid sadmin token")
		return cm.ErrPermissionDenied
	}

	var claim *claims.Claim
	var err error
	var account model.AccountInterface
	if q.RequireAPIKey {
		var expectType model.AccountType
		switch {
		case q.RequireShop:
			expectType = model.TypeShop
		case q.RequirePartner:
			expectType = model.TypePartner
		default:
			ll.Panic("unexpected account type")
		}
		claim, account, err = verifyAPIKey(ctx, token, expectType)
	} else if q.RequireAPIPartnerShopKey {
		claim, account, err = verifyAPIPartnerShopKey(ctx, token)
	} else {
		claim, err = tokens.Store.Validate(token)
	}

	if err != nil {
		return cm.ErrUnauthenticated
	}
	if claim.STokenExpiresAt != nil && claim.STokenExpiresAt.Before(time.Now()) {
		// Invalidate stoken
		claim.SToken = false
		claim.STokenExpiresAt = nil
	}

	if claim.AdminID != 0 {
		query := &model.GetSignedInUserQuery{UserID: claim.AdminID}
		if err := bus.Dispatch(ctx, query); err != nil {
			ll.Error("Invalid AdminID", l.Error(err))
			return nil
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
		return cm.ErrPermissionDenied
	}
	q.Result = session
	return nil
}

func startSessionAPIKey(ctx context.Context, require bool, s *Session) bool {
	return false
}

func startSessionUser(ctx context.Context, require bool, s *Session) bool {
	if require {
		if s.Claim.UserID == 0 {
			return false
		}
		query := &model.GetSignedInUserQuery{UserID: s.Claim.UserID}
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

	var partner *model.Partner
	partnerID := s.Claim.AuthPartnerID
	if partnerID != 0 {
		query := &model.GetPartner{
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

func startSessionPartner(ctx context.Context, require bool, s *Session, account model.AccountInterface) bool {
	if partner, ok := account.(*model.Partner); ok && s.Claim.AccountID == partner.ID {
		s.Partner = partner
		return true
	}
	if require {
		if !model.IsPartnerID(s.Claim.AccountID) {
			return false
		}
		query := &model.GetPartner{
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

func startSessionShop(ctx context.Context, require bool, s *Session, account model.AccountInterface) bool {
	if shop, ok := account.(*model.Shop); ok && s.Claim.AccountID == shop.ID {
		s.Shop = shop
		return true
	}
	if require {
		if !model.IsShopID(s.Claim.AccountID) && !model.IsAccountWhiteList(s.Claim.AccountID) {
			return false
		}

		if s.Claim.UserID != 0 {
			query := &model.GetShopWithPermissionQuery{
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
			query := &model.GetShopQuery{
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

func startSessionAffiliate(ctx context.Context, require bool, s *Session, account model.AccountInterface) bool {
	if affiliate, ok := account.(*identitymodel.Affiliate); ok && s.Claim.AccountID == affiliate.ID {
		s.Affiliate = affiliate
		return true
	}
	if require {
		if !model.IsAffiliateID(s.Claim.AccountID) && !model.IsAccountWhiteList(s.Claim.AccountID) {
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
		query := &model.GetAccountRolesQuery{
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

func verifyAPIKey(ctx context.Context, apikey string, expectType model.AccountType) (*claims.Claim, model.AccountInterface, error) {
	info, ok := authkey.ValidateAuthKeyWithType(authkey.TypeAPIKey, apikey)
	if !ok {
		return nil, nil, cm.Error(cm.Unauthenticated, "api_key không hợp lệ", nil)
	}

	query := &model.GetAccountAuthQuery{
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
	case model.TypePartner:
		partner := query.Result.Account.(*model.Partner)
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

	case model.TypeShop:
		shop := query.Result.Account.(*model.Shop)
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

func verifyAPIPartnerShopKey(ctx context.Context, apikey string) (*claims.Claim, model.AccountInterface, error) {
	_, ok := authkey.ValidateAuthKeyWithType(authkey.TypePartnerShopKey, apikey)
	if !ok {
		return nil, nil, cm.Error(cm.Unauthenticated, "api_key không hợp lệ", nil)
	}

	relationQuery := &model.GetPartnerRelationQuery{
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
