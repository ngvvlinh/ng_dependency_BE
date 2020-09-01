package session

import (
	"context"

	"o.o/backend/pkg/etop/authorize/auth"
	"o.o/backend/pkg/etop/authorize/permission"
)

// The pointer-Session implement Sessioner
var _ Sessioner = &Session{}

type Sessioner interface {
	StartSession(ctx context.Context, perm permission.Decl, tokenStr string) (context.Context, error)
	GetSession() *Session
}

// Session is designed so that it can be cloned with a simple assign. All
// pointer and non-pointer methods are intentional.
//
//     ss := session.New()
//     s2 := ss            // clone
type Session struct {
	SS session

	linkedSessions []Sessioner
}

// New returns a non-pointer Session, for embedding in other structs.
func New(auth *auth.Authorizer, opts ...Option) Session {
	s := Session{}
	s.SS.auth = auth
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
	newCtx, err := s.SS.startSession(ctx, perm, tokenStr)
	if err != nil {
		return newCtx, err
	}
	linkSessions(s, s.linkedSessions)
	return newCtx, nil
}

func linkSessions(top *Session, sessions []Sessioner) {
	for _, childSession := range sessions {
		child := childSession.GetSession()
		child.SS = top.SS
		linkSessions(top, child.linkedSessions)
	}
}

func (s *Session) GetSession() *Session {
	return s
}

// Link connects child sessions into a tree, linked with the top session. When
// the top session is initialized, all its children are copied. (It only happens
// once, as the session is only initialized once)
func (s *Session) Link(ss ...Sessioner) {
	if s.SS.init {
		panic("already init")
	}
	if s.linkedSessions != nil {
		panic("already linked")
	}
	s.linkedSessions = ss
}
