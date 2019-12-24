package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const DefaultRoute = "/==/prometrics"

func RegisterHTTPHandler(mux *http.ServeMux) {
	mux.Handle(DefaultRoute, promhttp.Handler())
}

func registerCounterVec(opts prometheus.CounterOpts, labels ...string) *prometheus.CounterVec {
	result := prometheus.NewCounterVec(opts, labels)
	if err := prometheus.Register(result); err != nil {
		panic(err)
	}
	return result
}

func registerHistogramVec(opts prometheus.HistogramOpts, labels ...string) *prometheus.HistogramVec {
	result := prometheus.NewHistogramVec(opts, labels)
	if err := prometheus.Register(result); err != nil {
		panic(err)
	}
	return result
}
