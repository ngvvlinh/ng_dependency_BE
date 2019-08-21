package httpreq

import (
	"net/http"

	"go.uber.org/zap/zapcore"
	"gopkg.in/resty.v1"

	"etop.vn/common/l"
)

type RestyConfig struct {
	Client *http.Client
}

type Resty struct {
	resty.Client
}

type RestyResponse = resty.Response

func NewResty(cfg RestyConfig) *Resty {
	var r *resty.Client
	if cfg.Client == nil {
		r = resty.New()
	} else {
		r = resty.NewWithClient(cfg.Client)
	}
	ll.Watch(func(level zapcore.Level) {
		enabled := level.Enabled(l.V(6))
		r.SetDebug(enabled)
	})

	client := &Resty{
		Client: *r,
	}
	return client
}
