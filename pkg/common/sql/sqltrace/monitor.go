package sqltrace

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"

	"o.o/backend/pkg/common/apifw/idemp"
	"o.o/backend/pkg/common/cmenv"
	"o.o/backend/pkg/common/metrics"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/common/l"
)

var ll = l.New()
var monitor *Monitor

const DefaultRoute = "/==/sqltrace"

func Init() {
	monitor = newMonitor()
	monitor.Start()
}

func RegisterHTTPHandler(mux *http.ServeMux) {
	if monitor == nil {
		Init()
	}
	mux.Handle(DefaultRoute, monitor)
}

func Trace(entry *sq.LogEntry) {
	if monitor != nil {
		monitor.Trace(entry)
	}
}

type Monitor struct {
	ch      chan *sq.LogEntry
	signal  chan int
	running bool
	mu      sync.RWMutex

	mapTraceItem map[string]*TraceItem
}

type TraceItem struct {
	Count int
	Last  TraceEntry
	Group []TraceEntry
}

type TraceEntry struct {
	Count       int
	Fingerprint string
	DBID        string
	Query       string
	Time        time.Time
	Duration    time.Duration
	Error       error
	OrigError   error
	Args        sq.LogArgs
}

func newMonitor() *Monitor {
	return &Monitor{
		ch:           make(chan *sq.LogEntry, 1024),
		signal:       make(chan int),
		mapTraceItem: map[string]*TraceItem{},
	}
}

func (m *Monitor) Start() {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.running {
		panic("already running")
	}

	m.running = true
	go m.run()
}

func (m *Monitor) run() {
	ll.Debug("sqltrace start")

	for {
		select {
		case <-m.signal:
			ll.Debug("sqltrace stop")
			m.mu.Lock()
			m.running = false
			m.mu.Unlock()
			return

		case entry := <-m.ch:
			m.handleEntry(entry)
		}
	}
}

func (m *Monitor) handleEntry(entry *sq.LogEntry) {
	if entry.IsQuery() {
		fingerprint, _ := Fingerprint(entry.Query)
		metrics.DatabaseQuery(fingerprint, entry)

		m.mu.Lock()
		defer m.mu.Unlock()
		item := m.mapTraceItem[fingerprint]
		if item == nil {
			item = &TraceItem{
				Group: make([]TraceEntry, 0, 10),
			}
			m.mapTraceItem[fingerprint] = item
		}
		savedEntry := TraceEntry{
			Count:       1,
			Fingerprint: fingerprint,
			DBID:        entry.DBID,
			Query:       entry.Query,
			Time:        entry.Time,
			Duration:    entry.Duration,
			Error:       entry.Error,
			OrigError:   entry.OrigError,
			Args:        entry.Args,
		}
		item.Group = appendTen(item.Group, &savedEntry)
		item.Last = savedEntry
		item.Count++
		return
	}

	t := entry.Type()
	if t == sq.TypeCommit || t == sq.TypeRollback {
		fingerprints := make([]string, len(entry.TxQueries))
		for i, query := range entry.TxQueries {
			fp, _ := Fingerprint(query.Query)
			fingerprints[i] = fp
		}
		metrics.DatabaseTransaction(fingerprints, entry)
	}
}

func appendTen(s []TraceEntry, entry *TraceEntry) []TraceEntry {
	for i, it := range s {
		if entry.Query == it.Query {
			entry.Count = it.Count + 1
			s[i] = *entry
			return s
		}
	}
	if len(s) == cap(s) {
		return s
	}
	s = append(s, *entry)
	return s
}

func (m *Monitor) Trace(entry *sq.LogEntry) {
	select {
	case m.ch <- entry:
		// no-op
	default:
		// no-blocking
	}
}

var idempgroup = idemp.NewGroup()

const timeLayout = "2006-01-02 15:04:05"

func (m *Monitor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var timeout time.Duration
	if cmenv.IsProd() {
		timeout = time.Minute
	}
	resp, _, _ := idempgroup.Do("key", timeout, func() (interface{}, error) {
		return m.serveHTTP(r), nil
	})
	_, _ = w.Write(resp.([]byte))
}

func (m *Monitor) serveHTTP(r *http.Request) []byte {
	m.mu.RLock()
	defer m.mu.RUnlock()

	items := make([]*TraceItem, 0, len(m.mapTraceItem))
	for _, item := range m.mapTraceItem {
		items = append(items, item)
	}
	// sort by count, then time
	sort.Slice(items, func(i, j int) bool {
		if items[i].Count != items[j].Count {
			return items[i].Count < items[j].Count
		}
		return items[i].Last.Time.Before(items[j].Last.Time)
	})

	b := &bytes.Buffer{}

	// list all queries
	w(b, "Refreshed %v\n\n", time.Now().UTC().Format(timeLayout))
	w(b, "[ ] % 5v\t%v\t% 19v\t%42v\t%v\t%v\t%v\n",
		"COUNT", "DURATION", "TIME", "FINGERPRINT", "DATABASE", "ERROR", "QUERY")
	for i := len(items) - 1; i >= 0; i-- {
		item := items[i]
		printEntry(b, -1, item.Last)
		for i, entry := range item.Group {
			if entry.Query != item.Last.Query {
				printEntry(b, i, entry)
			}
		}
	}
	return b.Bytes()
}

func printEntry(b io.Writer, i int, entry TraceEntry) {
	a := "*"
	if i >= 0 {
		a = strconv.Itoa(i)
	}
	w(b, "[%v] % 5v\t%v\t%v\t%v\t%v\t%v\t%v\n",
		a, entry.Count,
		entry.Duration, entry.Time.UTC().Format(timeLayout),
		entry.Fingerprint, entry.DBID, entry.OrigError, entry.Query)
}

func w(w io.Writer, format string, args ...interface{}) {
	_, _ = fmt.Fprintf(w, format, args...)
}
