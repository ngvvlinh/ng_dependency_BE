package middleware

import (
	"context"
	"net/http"
	"time"

	"o.o/api/main/identity"
	"o.o/api/top/types/etc/account_type"
	"o.o/api/top/types/etc/status3"
	identityconvert "o.o/backend/com/main/identity/convert"
	identitymodel "o.o/backend/com/main/identity/model"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/authorization/auth"
	"o.o/backend/pkg/common/headers"
	"o.o/backend/pkg/etc/idutil"
	"o.o/backend/pkg/etop/authorize/authkey"
	"o.o/backend/pkg/etop/authorize/claims"
	"o.o/backend/pkg/etop/authorize/permission"
	"o.o/backend/pkg/etop/authorize/tokens"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/capi/dot"
	"o.o/common/l"
)

var ll = l.New()

type SAdminToken string

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

type SessionStarter struct {
	SAdminToken   SAdminToken
	IdentityQuery identity.QueryBus
	TokenStore    tokens.TokenStore

	AccountStore     sqlstore.AccountStoreInterface
	UserStore        sqlstore.UserStoreInterface
	PartnerStore     sqlstore.PartnerStoreInterface
	AccountUserStore sqlstore.AccountUserStoreInterface
	ShopStore        sqlstore.ShopStoreInterface
}

// StartSession ...
func (st *SessionStarter) StartSession(ctx context.Context, q *StartSessionQuery) (newCtx context.Context, _err error) {
	token := getToken(ctx, q)
	return st.StartSessionWithToken(ctx, token, q)
}

func (st *SessionStarter) StartSessionWithToken(ctx context.Context, token string, q *StartSessionQuery) (newCtx context.Context, _err error) {
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
			claim, err := st.TokenStore.Validate(token)
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
		if token == string(st.SAdminToken) {
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
		claim, account, err = st.VerifyAPIKey(ctx, token, expectType)
		if err != nil {
			return ctx, err
		}
		_acc := account.GetAccount()
		if _acc.Type != expectType {
			return ctx, cm.Errorf(cm.PermissionDenied, nil, "")
		}
		wlPartnerID = _acc.ID

	} else if q.RequireAPIPartnerShopKey {
		claim, account, err = st.VerifyAPIPartnerShopKey(ctx, token)
		if err != nil {
			return ctx, err
		}
		wlPartnerID = claim.AuthPartnerID

	} else {
		claim, err = st.TokenStore.Validate(token)
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
		if err := st.UserStore.GetSignedInUser(ctx, query); err != nil {
			ll.Error("Invalid AdminID", l.Error(err))
			return ctx, nil
		}
		session.Admin = query.Result
	}

	session.Claim = claim
	ok := st.StartSessionUser(ctx, q.RequireUser, claim, &session.User) &&
		st.StartSessionPartner(ctx, q.RequirePartner, claim, account, &session.Partner) &&
		st.StartSessionShop(ctx, q.RequireShop, claim, account, &session.Shop, &session.Permission) &&
		st.StartSessionAffiliate(ctx, q.RequireAffiliate, claim, account, &session.Affiliate, &session.Permission) &&
		st.StartSessionEtopAdmin(ctx, q.RequireEtopAdmin, claim, &session.Permission) &&
		st.StartSessionAuthPartner(ctx, q.AuthPartner, claim, &session.CtxPartner)
	if !ok {
		return ctx, cm.ErrPermissionDenied
	}
	if account != nil {
		session.IsOwner = account.GetAccount().OwnerID == claim.UserID
	}
	q.Result = session
	return ctx, nil
}

func (st *SessionStarter) StartSessionUser(ctx context.Context, require bool, claim *claims.Claim, user **identitymodelx.SignedInUser) bool {
	if require {
		if claim.UserID == 0 {
			return false
		}
		query := &identitymodelx.GetSignedInUserQuery{
			UserID: claim.UserID,
		}
		if err := st.UserStore.GetSignedInUser(ctx, query); err != nil {
			ll.Error("Invalid UserID", l.Error(err))
			return false
		}
		if query.Result.Status == status3.N {
			return false
		}
		*user = query.Result
	}
	return true
}

func (st *SessionStarter) StartSessionAuthPartner(ctx context.Context, authOpt permission.AuthOpt, claim *claims.Claim, ctxPartner **identitymodel.Partner) bool {
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
		if err := st.PartnerStore.GetPartner(ctx, query); err != nil {
			ll.Error("Invalid PartnerID", l.Error(err))
			return false
		}
		partner = query.Result.Partner
	}

	*ctxPartner = partner
	return true
}

func (st *SessionStarter) StartSessionPartner(ctx context.Context, require bool, claim *claims.Claim, account identitymodel.AccountInterface, _partner **identitymodel.Partner) bool {
	if partner, ok := account.(*identitymodel.Partner); ok && claim.AccountID == partner.ID {
		*_partner = partner
		return true
	}
	if require {
		if !idutil.IsPartnerID(claim.AccountID) {
			return false
		}
		query := &identitymodelx.GetPartner{
			PartnerID: claim.AccountID,
		}
		if err := st.PartnerStore.GetPartner(ctx, query); err != nil {
			ll.Error("Invalid Name", l.Error(err))
			return false
		}
		*_partner = query.Result.Partner
	}
	return true
}

func (st *SessionStarter) StartSessionShop(ctx context.Context, require bool, claim *claims.Claim, account identitymodel.AccountInterface, _shop **identitymodel.Shop, _permission *identitymodel.Permission) bool {
	if shop, ok := account.(*identitymodel.Shop); ok && claim.AccountID == shop.ID {
		*_shop = shop
		return true
	}
	if require {
		if !idutil.IsShopID(claim.AccountID) {
			return false
		}

		if claim.UserID != 0 {
			query := &identitymodelx.GetShopWithPermissionQuery{
				ShopID: claim.AccountID,
				UserID: claim.UserID,
			}
			if err := st.ShopStore.GetShopWithPermission(ctx, query); err != nil {
				ll.Error("Invalid Name", l.Error(err))
				return false
			}
			if query.Result.Shop.Status == status3.N {
				return false
			}

			*_shop = query.Result.Shop
			*_permission = query.Result.Permission
		} else {
			query := &identitymodelx.GetShopQuery{
				ShopID: claim.AccountID,
			}
			if err := st.ShopStore.GetShop(ctx, query); err != nil {
				ll.Error("Invalid Name", l.Error(err))
				return false
			}
			*_shop = query.Result
		}
	}
	return true
}

func (st *SessionStarter) StartSessionAffiliate(ctx context.Context, require bool, claim *claims.Claim, account identitymodel.AccountInterface, _affiliate **identitymodel.Affiliate, _permission *identitymodel.Permission) bool {
	if affiliate, ok := account.(*identitymodel.Affiliate); ok && claim.AccountID == affiliate.ID {
		*_affiliate = affiliate
		return true
	}
	if require {
		if !idutil.IsAffiliateID(claim.AccountID) {
			return false
		}

		if claim.UserID != 0 {
			query := &identity.GetAffiliateWithPermissionQuery{
				AffiliateID: claim.AccountID,
				UserID:      claim.UserID,
			}
			if err := st.IdentityQuery.Dispatch(ctx, query); err != nil {
				ll.Error("Invalid Name", l.Error(err))
				return false
			}

			*_affiliate = identityconvert.AffiliateDB(query.Result.Affiliate)
			*_permission = identityconvert.PermissionToModel(query.Result.Permission)
		} else {
			query := &identity.GetAffiliateByIDQuery{
				ID: claim.AccountID,
			}
			if err := st.IdentityQuery.Dispatch(ctx, query); err != nil {
				ll.Error("Invalid Name", l.Error(err))
				return false
			}
			*_affiliate = identityconvert.AffiliateDB(query.Result)
		}
	}
	return true
}

func (st *SessionStarter) StartSessionEtopAdmin(ctx context.Context, require bool, claim *claims.Claim, _permission *identitymodel.Permission) bool {
	if require {
		if !idutil.IsEtopAccountID(claim.AccountID) {
			return false
		}
		query := &identitymodelx.GetAccountRolesQuery{
			AccountID: idutil.EtopAccountID,
			UserID:    claim.UserID,
		}
		if err := st.AccountUserStore.GetAccountUserExtended(ctx, query); err != nil {
			ll.Error("Invalid GetAccountRolesQuery", l.Error(err))
			return false
		}
		*_permission = query.Result.AccountUser.Permission
	}
	return true
}

func (st *SessionStarter) VerifyAPIKey(ctx context.Context, apikey string, expectType account_type.AccountType) (*claims.Claim, identitymodel.AccountInterface, error) {
	info, ok := authkey.ValidateAuthKeyWithType(authkey.TypeAPIKey, apikey)
	if !ok {
		return nil, nil, cm.Error(cm.Unauthenticated, "api_key không hợp lệ", nil)
	}

	query := &identitymodelx.GetAccountAuthQuery{
		AuthKey:     apikey,
		AccountType: expectType,
		AccountID:   info.AccountID,
	}
	if err := st.AccountStore.GetAccountAuth(ctx, query); err != nil {
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

func (st *SessionStarter) VerifyAPIPartnerShopKey(ctx context.Context, apikey string) (*claims.Claim, identitymodel.AccountInterface, error) {
	_, ok := authkey.ValidateAuthKeyWithType(authkey.TypePartnerShopKey, apikey)
	if !ok {
		return nil, nil, cm.Error(cm.Unauthenticated, "api_key không hợp lệ", nil)
	}

	relationQuery := &identitymodelx.GetPartnerRelationQuery{
		AuthKey: apikey,
	}
	relationError := st.PartnerStore.GetPartnerRelationQuery(ctx, relationQuery)
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
