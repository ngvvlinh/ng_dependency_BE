package session

import "o.o/backend/pkg/etop/authorize/tokens"

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
