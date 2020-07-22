package health

import (
	"fmt"
	"net/http"

	"go.uber.org/atomic"

	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/cmenv"
	"o.o/backend/pkg/common/redis"
	"o.o/common/l"
)

const DefaultRoute = "/healthcheck"

var ll = l.New()
var llDeploy = ll.WithChannel("deploy")

// New returns new health.Service
func New(redis redis.Store) *Service {
	return &Service{redis: redis}
}

type Service struct {
	ready atomic.Bool

	// redis may be nil
	redis redis.Store
}

// RegisterHTTPHandler registers health service.
func (s *Service) RegisterHTTPHandler(mux *http.ServeMux) {
	mux.Handle(DefaultRoute, s)
}

// MarkReady marks the service ready
func (s *Service) MarkReady() {
	s.ready.Store(true)

	msg := fmt.Sprintf("✨ %v on %v started ✨\n%v", cmenv.ServiceName(), cmenv.Env(), cm.CommitMessage())
	ll.SendMessage(msg)

	// redis may be nil, ignore the deploy channel
	if s.redis == nil {
		return
	}

	// also send message to the deploy channel, ignore duplicated messages
	key := fmt.Sprintf("deploy:%v+%v", cmenv.ServiceName(), cmenv.Env())
	lastMsg, err := s.redis.GetString(key)
	if err != nil && err != redis.ErrNil {
		ll.Panic("redis error", l.Error(err))
	}
	if msg != lastMsg {
		if err = s.redis.SetString(key, msg); err != nil {
			ll.Panic("redis error", l.Error(err))
		}
		llDeploy.SendMessagef("✨ %v on %v started ✨\n%v", cmenv.ServiceName(), cmenv.Env(), cm.CommitMessage())
	}
}

// Shutdown marks the service not ready
func (s *Service) Shutdown() {
	s.ready.Store(false)

	ll.SendMessagef("✨ %v on %v stopped ✨\n%v", cmenv.ServiceName(), cmenv.Env(), cm.CommitMessage())
}

// ServeHTTP implements http healthcheck
func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status := 400
	if s.ready.Load() {
		status = 200
	}
	w.WriteHeader(status)
}
