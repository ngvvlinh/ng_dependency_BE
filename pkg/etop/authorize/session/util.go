package session

import (
	identitymodel "o.o/backend/com/main/identity/model"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/authorize/claims"
	"o.o/backend/pkg/etop/authorize/middleware"
	"o.o/backend/pkg/etop/authorize/permission"
)

func (s *session) Claim() claims.Claim {
	s.ensureInit()
	return s.claim
}

func (s *session) Admin() *identitymodelx.SignedInUser {
	s.ensureInit()
	return s.admin
}

func (s *session) User() *identitymodelx.SignedInUser {
	s.ensureInit()
	if s.user != nil {
		return s.user
	}
	middleware.StartSessionUser(s.ctx, true, &s.claim, &s.user)
	return s.user
}

func (s *session) Shop() *identitymodel.Shop {
	s.ensureInit()
	return s.shop
}

func (s *session) Partner() *identitymodel.Partner {
	s.ensureInit()
	return s.partner
}

func (s *session) CtxPartner() *identitymodel.Partner {
	s.ensureInit()
	return s.ctxPartner
}

func (s *session) Affiliate() *identitymodel.Affiliate {
	s.ensureInit()
	return s.affiliate
}

func (s *session) Permission() identitymodel.Permission {
	s.ensureInit()
	if s.claim.UserID == 0 || s.claim.AccountID == 0 {
		s.permission = identitymodel.Permission{Permissions: []string{}}
	}
	if s.permission.Permissions != nil {
		return s.permission
	}
	accQuery := &identitymodelx.GetAccountRolesQuery{
		AccountID: s.claim.AccountID,
		UserID:    s.claim.UserID,
	}
	if err := bus.Dispatch(s.ctx, accQuery); err != nil {
		panic(err)
	}
	s.permission = accQuery.Result.AccountUser.Permission
	if s.permission.Permissions == nil {
		s.permission.Permissions = []string{}
	}
	return s.permission
}

func (s *session) PermissionDecl() permission.Decl {
	s.ensureInit()
	return s.perm
}

func (s *session) IsSuperAdmin() bool {
	s.ensureInit()
	return s.isSuperAdmin
}

func (s *session) IsAdmin() bool {
	s.ensureInit()
	return s.isAdmin
}

func (s *session) IsOwner() bool {
	s.ensureInit()
	return s.isOwner
}
