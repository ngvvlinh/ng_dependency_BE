package session

import "o.o/backend/pkg/etop/authorize/tokens"

type HookOption func(h *Hook) error

func OptSecret(secret string) HookOption {
	return func(h *Hook) error {
		h.secret = secret
		return nil
	}
}

type Option func(s *Session) error

func OptSuperAdmin(token string) Option {
	return func(s *Session) error {
		s.SS.sadminToken = token
		return nil
	}
}

func OptValidator(validator tokens.Validator) Option {
	return func(s *Session) error {
		s.SS.validator = validator
		return nil
	}
}
