package metrics

import (
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"

	"etop.vn/backend/pkg/common/sql/sq"
)

var databaseQuery = registerHistogramVec(prometheus.HistogramOpts{
	Namespace: "",
	Subsystem: "",
	Name:      "main_database_query_histogram",
}, "tx", "query")

func DatabaseQuery(entry *sq.LogEntry) {
	databaseQuery.WithLabelValues(
		strconv.FormatBool(entry.IsTx()),
		entry.Query,
	).Observe(entry.Duration.Seconds())
}

var databaseTransaction = registerHistogramVec(prometheus.HistogramOpts{
	Namespace: "",
	Subsystem: "",
	Name:      "main_database_transaction_histogram",
}, "n", "type", "query")

func DatabaseTransaction(entry *sq.LogEntry) {
	var n int
	for _, query := range entry.TxQueries {
		n += len(query.Query) + 1
	}
	var b strings.Builder
	b.Grow(n)
	for _, query := range entry.TxQueries {
		b.WriteString(query.Query)
		b.WriteByte(';')
	}

	flag, _ := entry.Flags.MarshalJSON()
	databaseTransaction.WithLabelValues(
		strconv.Itoa(len(entry.TxQueries)),
		string(flag),
		b.String(),
	).Observe(entry.Duration.Seconds())
}
