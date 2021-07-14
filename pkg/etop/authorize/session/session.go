package session

import (
	"context"
	"time"

	"o.o/api/top/types/etc/account_type"
	identitymodel "o.o/backend/com/main/identity/model"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/authorize/auth"
	"o.o/backend/pkg/etop/authorize/claims"
	"o.o/backend/pkg/etop/authorize/middleware"
	"o.o/backend/pkg/etop/authorize/permission"
	"o.o/backend/pkg/etop/authorize/tokens"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/capi/dot"
	"o.o/common/l"
)

var ll = l.New()

type session struct {
	init        bool
	ctx         context.Context
	sadminToken string
	validator   tokens.Validator

	auth  *auth.Authorizer
	perm  permission.Decl
	claim claims.Claim

	admin      *identitymodelx.SignedInUser
	user       *identitymodelx.SignedInUser
	shop       *identitymodel.Shop
	partner    *identitymodel.Partner
	ctxPartner *identitymodel.Partner
	affiliate  *identitymodel.Affiliate
	permission identitymodel.Permission

	isSuperAdmin bool
	isAdmin      bool
	isOwner      bool

	st               *middleware.SessionStarter
	AccountUserStore sqlstore.AccountUserStoreInterface
	UserStore        sqlstore.UserStoreInterface
}

func (s *session) ensureInit() {
	if !s.init {
		panic("not init")
	}
}

func (s *session) ensureNotInit() {
	if s.init {
		panic("already init")
	}
}

func (s *session) startSession(ctx context.Context, perm permission.Decl, tokenStr string) (newCtx context.Context, _ error) {
	if s.init {
		panic("already init")
	}
	s.init = true
	s.ctx = ctx
	var wlPartnerID dot.ID
	defer func() {
		if wlPartnerID == 0 {
			newCtx = wl.WrapContext(ctx, 0)
		}
		s.ctx = newCtx
	}()

	if tokenStr == "" && perm.Type == permission.Public {
		return ctx, nil
	}
	if tokenStr == "" && perm.Type != permission.Public {
		return ctx, cm.Errorf(cm.Unauthenticated, nil, "")
	}

	if perm.Type == permission.SuperAdmin {
		if tokenStr != s.sadminToken {
			return ctx, cm.Errorf(cm.Unauthenticated, nil, "").
				Logf("invalid sadmin token")
		}
		s.isSuperAdmin = true
		return ctx, nil
	}

	wlPartnerID, claim, account, err := s.verifyToken(ctx, perm, tokenStr)
	if err != nil {
		// Ignore invalid token for public permission. TopShip App is sending
		// invalid token even for public API.
		if perm.Type == permission.Public {
			return ctx, nil
		}
		return ctx, err
	}
	ctx = wl.WrapContext(ctx, wlPartnerID)
	s.claim = *claim
	s.Permission() // load permission of account

	// handle stoken
	if claim.STokenExpiresAt != nil && claim.STokenExpiresAt.Before(time.Now()) {
		// invalidate stoken
		claim.SToken = false
		claim.STokenExpiresAt = nil
	}

	// handle admin authorized as user
	if claim.AdminID != 0 {
		query := &identitymodelx.GetSignedInUserQuery{
			UserID: claim.AdminID,
		}
		if err = s.UserStore.GetSignedInUser(ctx, query); err != nil {
			ll.Error("Invalid AdminID", l.Error(err))
			return ctx, nil
		}
		s.admin = query.Result
	}

	// verify permission
	for _, action := range perm.Actions {
		if !s.auth.Check(s.permission.Roles, string(action), 0) {
			return ctx, cm.ErrPermissionDenied
		}
	}

	ok := s.st.StartSessionUser(ctx, perm.Type == permission.CurUsr || perm.Auth == permission.User, claim, &s.user) &&
		s.st.StartSessionPartner(ctx, perm.Type == permission.Partner, claim, account, &s.partner) &&
		s.st.StartSessionShop(ctx, perm.Type == permission.Shop, claim, account, &s.shop, &s.permission) &&
		s.st.StartSessionAffiliate(ctx, perm.Type == permission.Affiliate, claim, account, &s.affiliate, &s.permission) &&
		s.st.StartSessionEtopAdmin(ctx, perm.Type == permission.EtopAdmin, claim, &s.permission) &&
		s.st.StartSessionAuthPartner(ctx, perm.AuthPartner, claim, &s.ctxPartner)
	if !ok {
		return ctx, cm.ErrPermissionDenied
	}
	if account != nil {
		s.isOwner = account.GetAccount().OwnerID == claim.UserID
	}

	ctx = bus.NewRootContext(ctx)
	return ctx, nil
}

func (s *session) verifyToken(
	ctx context.Context,
	perm permission.Decl,
	tokenStr string,
) (
	wlPartnerID dot.ID,
	claim *claims.Claim,
	account identitymodel.AccountInterface,
	err error,
) {
	defer func() {
		switch cm.ErrorCode(err) {
		case cm.NotFound:
			err = cm.Errorf(cm.Unauthenticated, err, "")
		}
	}()

	switch perm.Auth {
	case permission.APIKey:
		switch perm.Type {
		case permission.Shop:
			claim, account, err = s.st.VerifyAPIKey(ctx, tokenStr, account_type.Shop)
			if err != nil {
				return
			}
			wlPartnerID = 0 // TODO: api.itopx.vn?
			return

		case permission.Partner:
			claim, account, err = s.st.VerifyAPIKey(ctx, tokenStr, account_type.Partner)
			if err != nil {
				return
			}
			_acc := account.GetAccount()
			if _acc.Type != account_type.Partner {
				err = cm.Errorf(cm.PermissionDenied, nil, "")
				return
			}
			wlPartnerID = _acc.ID
			return

		default:
			ll.Panic("unexpected type", l.Any("type", perm.Type))
			return
		}

	case permission.APIPartnerShopKey:
		claim, account, err = s.st.VerifyAPIPartnerShopKey(ctx, tokenStr)
		if err != nil {
			return
		}
		wlPartnerID = claim.AuthPartnerID
		return

	case permission.APIPartnerCarrierKey:
		claim, account, err = s.st.VerifyAPIKey(ctx, tokenStr, account_type.Carrier)
		if err != nil {
			return
		}
		_acc := account.GetAccount()
		if _acc.Type != account_type.Carrier {
			err = cm.Errorf(cm.PermissionDenied, nil, "")
			return
		}

		wlPartnerID = _acc.ID
		return

	default:
		claim, err = s.validator.Validate(tokenStr)
		if err != nil {
			return
		}
		if perm.AuthPartner.AuthPartner() {
			wlPartnerID = claim.AuthPartnerID
		} else if claim.WLPartnerID != 0 {
			wlPartnerID = claim.WLPartnerID
		}
		return
	}
}

func (s *session) GetRoles() []string {
	return s.Permission().Roles
}
