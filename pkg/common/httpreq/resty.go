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
	client := &Resty{}
	if cfg.Client == nil {
		client.Client = *resty.New()
	} else {
		client.Client = *resty.NewWithClient(cfg.Client)
	}
	ll.Watch(func(name string, level zapcore.Level) {
		enabled := level.Enabled(l.V(6))
		client.Client.SetDebug(enabled)
	})
	return client
}
