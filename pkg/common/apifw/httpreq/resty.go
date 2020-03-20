package httpreq

import (
	"net/http"
	"time"

	"go.uber.org/zap/zapcore"
	"gopkg.in/resty.v1"

	"etop.vn/backend/pkg/common/metrics"
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
	httpClient := &http.Client{} // make a new client
	if cfg.Client != nil {
		*httpClient = *cfg.Client // copy the provided client
	}
	if _, ok := httpClient.Transport.(MeasuredTransport); !ok {
		httpClient.Transport = NewMeasuredTransport(httpClient.Transport)
	}

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

type MeasuredTransport struct {
	http.RoundTripper
}

func NewMeasuredTransport(rt http.RoundTripper) http.RoundTripper {
	if rt == nil {
		rt = http.DefaultTransport
	}
	return MeasuredTransport{RoundTripper: rt}
}

func (m MeasuredTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	t0 := time.Now()
	defer func() {
		d := time.Now().Sub(t0)
		if err != nil {
			metrics.EgressRequest(req.URL, -1, d)
		} else {
			metrics.EgressRequest(req.URL, resp.StatusCode, d)
		}
	}()
	return m.RoundTripper.RoundTrip(req)
}
