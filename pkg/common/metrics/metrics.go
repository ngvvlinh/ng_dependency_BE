package metrics

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	cm "etop.vn/backend/pkg/common"
)

const DefaultRoute = "/==prometrics=="

func init() {
	prometheus.MustRegister(countRequests)
}

func RegisterHTTPHandler(mux *http.ServeMux) {
	mux.Handle(DefaultRoute, promhttp.Handler())
}

var countRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
	Namespace:   "",
	Subsystem:   "",
	Name:        "main_api_requests_total",
	Help:        "",
	ConstLabels: nil,
}, []string{"name", "code"})

func CountRequest(name string, err error) {
	code := cm.ErrorCode(err)
	countRequests.WithLabelValues(name, strconv.Itoa(int(code))).Inc()
}
