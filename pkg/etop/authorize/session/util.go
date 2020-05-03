package session

import (
	identitymodel "o.o/backend/com/main/identity/model"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	"o.o/backend/pkg/etop/authorize/claims"
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

func (s *Session) GetClaim() *claims.Claim {
	s.ensureInit()
	return s.claim
}

func (s *Session) GetAdmin() *identitymodelx.SignedInUser {
	s.ensureInit()
	return s.admin
}

func (s *Session) GetUser() *identitymodelx.SignedInUser {
	s.ensureInit()
	return s.user
}

func (s *Session) GetShop() *identitymodel.Shop {
	s.ensureInit()
	return s.shop
}

func (s *Session) GetPartner() *identitymodel.Partner {
	s.ensureInit()
	return s.partner
}

func (s *Session) GetCtxPartner() *identitymodel.Partner {
	s.ensureInit()
	return s.ctxPartner
}

func (s *Session) GetAffiliate() *identitymodel.Affiliate {
	s.ensureInit()
	return s.affiliate
}

func (s *Session) GetPermission() identitymodel.Permission {
	s.ensureInit()
	return s.permission
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
