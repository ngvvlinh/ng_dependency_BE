package session

import (
	identitymodel "o.o/backend/com/main/identity/model"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	"o.o/backend/pkg/etop/authorize/claims"
	"o.o/backend/pkg/etop/authorize/middleware"
	"o.o/backend/pkg/etop/authorize/permission"
	"o.o/backend/pkg/etop/authorize/tokens"
)

type Option func(s *Session) error

func OptSuperAdmin(token string) Option {
	return func(s *Session) error {
		s.sadminToken = token
		return nil
	}
}

func OptValidator(validator tokens.Validator) Option {
	return func(s *Session) error {
		s.validator = validator
		return nil
	}
}

func (s *Session) Claim() *claims.Claim {
	s.ensureInit()
	return s.claim
}

func (s *Session) Admin() *identitymodelx.SignedInUser {
	s.ensureInit()
	return s.admin
}

func (s *Session) User() *identitymodelx.SignedInUser {
	s.ensureInit()
	if s.user != nil {
		return s.user
	}
	middleware.StartSessionUser(s.ctx, true, s.claim, &s.user)
	return s.user
}

func (s *Session) Shop() *identitymodel.Shop {
	s.ensureInit()
	return s.shop
}

func (s *Session) Partner() *identitymodel.Partner {
	s.ensureInit()
	return s.partner
}

func (s *Session) CtxPartner() *identitymodel.Partner {
	s.ensureInit()
	return s.ctxPartner
}

func (s *Session) Affiliate() *identitymodel.Affiliate {
	s.ensureInit()
	return s.affiliate
}

func (s *Session) Permission() identitymodel.Permission {
	s.ensureInit()
	return s.permission
}

func (s *Session) PermissionDecl() permission.Decl {
	s.ensureInit()
	return s.perm
}

func (s *Session) IsSuperAdmin() bool {
	s.ensureInit()
	return s.isSuperAdmin
}

func (s *Session) IsAdmin() bool {
	s.ensureInit()
	return s.isAdmin
}

func (s *Session) IsOwner() bool {
	s.ensureInit()
	return s.isOwner
}
