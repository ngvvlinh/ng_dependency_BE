package session

import (
	"context"

	"o.o/backend/pkg/etop/authorize/permission"
)

// The pointer-Session implement Sessioner
var _ Sessioner = &Session{}

type Sessioner interface {
	StartSession(ctx context.Context, perm permission.Decl, tokenStr string) (context.Context, error)
	GetSession() *Session
}

// Session is designed so that it can be cloned with a simple assign. All
// pointer and non-pointer methods are intention.
//
//     ss := session.New()
//     s2 := ss            // clone
type Session struct {
	SS session
}

// New returns a non-pointer Session, for embedding in other structs.
func New(opts ...Option) Session {
	s := Session{}
	return s.MustWith(opts...)
}

func (s Session) With(opts ...Option) (Session, error) {
	s.SS.ensureNotInit()
	if len(opts) == 0 {
		return s, nil
	}
	for _, opt := range opts {
		err := opt(&s)
		if err != nil {
			return s, err
		}
	}
	return s, nil
}

func (s Session) MustWith(opts ...Option) Session {
	res, err := s.With(opts...)
	if err != nil {
		panic(err)
	}
	return res
}

// StartSession runs on a pointer-Session to fill its details.
func (s *Session) StartSession(ctx context.Context, perm permission.Decl, tokenStr string) (newCtx context.Context, _ error) {
	return s.SS.startSession(ctx, perm, tokenStr)
}

func (s *Session) GetSession() *Session {
	return s
}
