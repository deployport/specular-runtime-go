package server

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"go.deployport.com/specular-runtime/client"
)

const (
	// telemetryDefaultEvery is the push interval when Telemetry.Every is zero.
	telemetryDefaultEvery = 60 * time.Second
	// telemetryMaxKeys caps the distinct series one interval may mint, so a
	// flood of generated operation names cannot grow the store without bound.
	telemetryMaxKeys = 1024
	// The push uses the same tight budget as the platform metrics clients: a
	// slow store is treated as unavailability, never as a reason to wait.
	telemetryDialTimeout   = 2 * time.Second
	telemetryHeaderTimeout = 5 * time.Second
	telemetryPushTimeout   = 5 * time.Second
	// telemetryKeySeparator splits the operation from the status class inside a
	// series key. Neither value can carry it.
	telemetryKeySeparator = "\x1f"
)

// telemetryBuckets are the histogram bounds in seconds. A specular endpoint
// slower than ten seconds is a fault, not a bucket to distinguish.
var telemetryBuckets = []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10}

// Telemetry configures the per-operation measurement. Every field is optional
// except Service; the recorder is OFF unless WriteURL is set, and then it
// starts no goroutine and pushes nothing.
type Telemetry struct {
	// Service is the series label, e.g. "machinery". Each service passes its
	// own DEPLOYPORT_SERVICE value.
	Service string
	// WriteURL is the VictoriaMetrics base URL, no path. The recorder appends
	// /api/v1/import/prometheus, the route the machinery and via metrics
	// clients already use.
	WriteURL string
	// Username and Password are sent as basic auth when Username is set.
	Username string
	Password string
	// Prefix namespaces the series, as the platform does today: one VM host
	// backs dev and both production regions, so names read "<prefix>_<name>".
	Prefix string
	// Every is the push interval. Zero means 60 seconds.
	Every time.Duration
	// HTTPClient and Now are for tests.
	HTTPClient *http.Client
	Now        func() time.Time
}

// WithTelemetry turns the measurement on. Every operation the service serves is
// then timed from one place, keyed by operation and status class.
func WithTelemetry(t Telemetry) ServiceOptionFunc {
	return func(s *Service) {
		s.telemetry = newTelemetryRecorder(t)
	}
}

// telemetrySeries is one (operation, status class) pair. Every field is atomic
// so an observation never waits, and a flush reads a value mid-flight without
// blocking the writers.
type telemetrySeries struct {
	// stream tells which histogram this series belongs to. An operation is
	// either streamed or not for the life of the process.
	stream bool
	// buckets counts observations that landed in each bound, EXCLUSIVE. The
	// render accumulates them, so a bucket is cumulative on the wire and a
	// single observation is stored once.
	buckets []atomic.Int64
	// sum is the total duration in nanoseconds.
	sum   atomic.Int64
	count atomic.Int64
	// total counts the requests, whatever the duration. It is the counter
	// beside the histogram, so a failing endpoint is one query and needs no
	// bucket.
	total atomic.Int64
}

func newTelemetrySeries(stream bool) *telemetrySeries {
	return &telemetrySeries{
		stream:  stream,
		buckets: make([]atomic.Int64, len(telemetryBuckets)),
	}
}

// observe records one request. The first bucket whose bound reaches the
// duration takes it, which is what Prometheus means by le: a duration exactly
// on a bound belongs to that bound, not the next one.
func (s *telemetrySeries) observe(took time.Duration) {
	s.count.Add(1)
	s.sum.Add(took.Nanoseconds())
	seconds := took.Seconds()
	for i, bound := range telemetryBuckets {
		if seconds <= bound {
			s.buckets[i].Add(1)
			return
		}
	}
	// Longer than the last bound: it counts in +Inf only, and the render takes
	// +Inf straight from count so the two can never disagree.
}

// telemetryStore is one interval of observations. The recorder swaps the whole
// store, so a flush renders a map nothing new is written to.
type telemetryStore struct {
	series sync.Map // series key -> *telemetrySeries
	// count is the number of keys in series. It is guarded by the recorder's
	// insertMu, which is the only writer.
	count int
}

// telemetryRecorder renders the observations and pushes them to the store.
type telemetryRecorder struct {
	service  string
	prefix   string
	endpoint string
	username string
	password string
	every    time.Duration
	client   *http.Client
	now      func() time.Time
	logger   *slog.Logger

	// current is the interval being written. A swap is one pointer store, so
	// the request path is never blocked by a flush.
	current atomic.Pointer[telemetryStore]

	// insertMu guards a new key. An observation of a known key never takes it.
	insertMu sync.Mutex
	// capped makes the key cap warn once for the life of the process, so a
	// flood produces one line and not one per request.
	capped atomic.Bool
}

func newTelemetryRecorder(t Telemetry) *telemetryRecorder {
	if t.WriteURL == "" {
		return nil
	}
	httpClient := t.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Transport: newTelemetryTransport()}
	}
	every := t.Every
	if every <= 0 {
		every = telemetryDefaultEvery
	}
	now := t.Now
	if now == nil {
		now = time.Now
	}
	r := &telemetryRecorder{
		service:  t.Service,
		prefix:   t.Prefix,
		endpoint: strings.TrimRight(t.WriteURL, "/") + "/api/v1/import/prometheus",
		username: t.Username,
		password: t.Password,
		every:    every,
		client:   httpClient,
		now:      now,
		logger:   createDefaultLogger(),
	}
	r.current.Store(&telemetryStore{})
	return r
}

// telemetryTransport bounds every leg of a push, so a store that is restarting
// cannot tie up a connection.
func newTelemetryTransport() *http.Transport {
	return &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   telemetryDialTimeout,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          2,
		MaxIdleConnsPerHost:   2,
		IdleConnTimeout:       90 * time.Second,
		ResponseHeaderTimeout: telemetryHeaderTimeout,
	}
}

// seriesName puts the configured prefix in front of a base metric name, as the
// platform does, because one store is shared by dev and both regions.
func (r *telemetryRecorder) seriesName(name string) string {
	if r.prefix == "" {
		return name
	}
	return r.prefix + "_" + name
}

// run pushes one interval at a time until the process ends. A service lives as
// long as its process, so the recorder needs no stop.
func (r *telemetryRecorder) run() {
	ticker := time.NewTicker(r.every)
	defer ticker.Stop()
	for range ticker.C {
		r.flushAndLog()
	}
}

// flushAndLog pushes the interval and reports a failure without failing
// anything else. A metrics outage is not a request outage.
func (r *telemetryRecorder) flushAndLog() {
	if err := r.flush(); err != nil {
		r.logger.Warn("metrics push failed", "err", err, "endpoint", r.endpoint)
	}
}

// flush renders everything observed since the previous flush and pushes it.
//
// The store is swapped FIRST, so observations that arrive during the render and
// the push land in the next interval and never wait. The visible cost is that
// an observation written into the old store after its render is lost, which is
// one request in a minute against a lock on every request.
func (r *telemetryRecorder) flush() error {
	previous := r.current.Swap(&telemetryStore{})
	body := r.render(previous, r.now().UnixMilli())
	if len(body) == 0 {
		// Nothing was served in the interval. An empty push says nothing and
		// costs a round trip.
		return nil
	}
	return r.push(body)
}

// observe records one served request against its operation.
func (r *telemetryRecorder) observe(op *client.Operation, status int, took time.Duration) {
	if op == nil {
		return
	}
	streamed := op.IsStreamed()
	operation := op.Resource().Name() + "." + op.Name()
	series := r.seriesFor(operation+telemetryKeySeparator+statusClassOf(status), streamed)
	if series == nil {
		return
	}
	series.observe(took)
	if !streamed {
		series.total.Add(1)
	}
}

// seriesFor returns the series for a key, minting it when the cap permits. The
// fast path is lock free; only the first observation of a key takes the lock.
func (r *telemetryRecorder) seriesFor(key string, streamed bool) *telemetrySeries {
	store := r.current.Load()
	if v, ok := store.series.Load(key); ok {
		return v.(*telemetrySeries)
	}
	r.insertMu.Lock()
	defer r.insertMu.Unlock()
	// The interval may have been swapped while the lock was pending, so resolve
	// the store again before adding to it.
	store = r.current.Load()
	if v, ok := store.series.Load(key); ok {
		return v.(*telemetrySeries)
	}
	if store.count >= telemetryMaxKeys {
		if r.capped.CompareAndSwap(false, true) {
			r.logger.Warn(
				"metrics series cap reached, new keys are dropped",
				"service", r.service,
				"cap", telemetryMaxKeys,
			)
		}
		return nil
	}
	series := newTelemetrySeries(streamed)
	store.series.Store(key, series)
	store.count++
	return series
}

// push sends one rendered interval. The body is dropped on the floor on
// failure: see flushAndLog.
func (r *telemetryRecorder) push(body []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), telemetryPushTimeout)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, r.endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "text/plain")
	if r.username != "" {
		request.SetBasicAuth(r.username, r.password)
	}
	response, err := r.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	// Read a little so the connection can be reused, but never the whole body
	// of an unexpected answer.
	_, _ = io.Copy(io.Discard, io.LimitReader(response.Body, 4<<10))
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("metrics store answered %s", response.Status)
	}
	return nil
}

// render turns one interval into the Prometheus text the store imports, with a
// millisecond timestamp on every sample. Keys are sorted so the same interval
// always renders the same bytes.
func (r *telemetryRecorder) render(store *telemetryStore, timestampMillis int64) []byte {
	keys := make([]string, 0, 64)
	store.series.Range(func(k, _ any) bool {
		keys = append(keys, k.(string))
		return true
	})
	if len(keys) == 0 {
		return nil
	}
	sort.Strings(keys)
	var out bytes.Buffer
	for _, key := range keys {
		operation, class, ok := strings.Cut(key, telemetryKeySeparator)
		if !ok {
			continue
		}
		value, ok := store.series.Load(key)
		if !ok {
			continue
		}
		r.renderSeries(&out, value.(*telemetrySeries), operation, class, timestampMillis)
	}
	return out.Bytes()
}

func (r *telemetryRecorder) renderSeries(
	out *bytes.Buffer,
	series *telemetrySeries,
	operation string,
	class string,
	timestampMillis int64,
) {
	base := "specular_operation_duration_seconds"
	if series.stream {
		base = "specular_stream_open_seconds"
	}
	base = r.seriesName(base)
	labels := telemetryLabels(r.service, operation, class)

	count := series.count.Load()
	var cumulative int64
	for i, bound := range telemetryBuckets {
		cumulative += series.buckets[i].Load()
		writeTelemetrySample(
			out,
			base+"_bucket",
			labels+`,le="`+formatTelemetryBound(bound)+`"`,
			strconv.FormatInt(cumulative, 10),
			timestampMillis,
		)
	}
	// +Inf is count itself, so it cannot disagree with the observations that
	// were longer than the last bound.
	writeTelemetrySample(out, base+"_bucket", labels+`,le="+Inf"`, strconv.FormatInt(count, 10), timestampMillis)
	writeTelemetrySample(out, base+"_sum", labels, formatTelemetrySeconds(series.sum.Load()), timestampMillis)
	writeTelemetrySample(out, base+"_count", labels, strconv.FormatInt(count, 10), timestampMillis)
	if !series.stream {
		writeTelemetrySample(
			out,
			r.seriesName("specular_operation_total"),
			labels,
			strconv.FormatInt(series.total.Load(), 10),
			timestampMillis,
		)
	}
}

// telemetryLabels is the whole label set. Only the operation and the status
// class travel: an account, a user, a machine id or a DRN would turn one series
// per endpoint into one series per customer, and a metrics store is not a place
// to keep customer identifiers.
func telemetryLabels(service string, operation string, class string) string {
	return `service="` + escapeTelemetryLabel(service) +
		`",operation="` + escapeTelemetryLabel(operation) +
		`",status_class="` + escapeTelemetryLabel(class) + `"`
}

func escapeTelemetryLabel(value string) string {
	if !strings.ContainsAny(value, "\\\"\n") {
		return value
	}
	replacer := strings.NewReplacer(`\`, `\\`, `"`, `\"`, "\n", `\n`)
	return replacer.Replace(value)
}

func writeTelemetrySample(out *bytes.Buffer, name string, labels string, value string, timestampMillis int64) {
	out.WriteString(name)
	out.WriteByte('{')
	out.WriteString(labels)
	out.WriteString("} ")
	out.WriteString(value)
	out.WriteByte(' ')
	out.WriteString(strconv.FormatInt(timestampMillis, 10))
	out.WriteByte('\n')
}

func formatTelemetryBound(bound float64) string {
	return strconv.FormatFloat(bound, 'f', -1, 64)
}

func formatTelemetrySeconds(nanoseconds int64) string {
	return strconv.FormatFloat(float64(nanoseconds)/float64(time.Second), 'f', -1, 64)
}

// statusClassOf reduces a status code to the class a dashboard groups by.
func statusClassOf(status int) string {
	switch {
	case status >= 200 && status < 300:
		return "2xx"
	case status >= 300 && status < 400:
		return "3xx"
	case status >= 400 && status < 500:
		return "4xx"
	case status >= 500 && status < 600:
		return "5xx"
	default:
		return "other"
	}
}

// statusResponseWriter records the status the handler answered and forwards
// everything else, so one measurement covers every endpoint.
type statusResponseWriter struct {
	http.ResponseWriter
	status int
	wrote  bool
}

func newStatusResponseWriter(w http.ResponseWriter) *statusResponseWriter {
	return &statusResponseWriter{ResponseWriter: w, status: http.StatusOK}
}

func (w *statusResponseWriter) WriteHeader(status int) {
	if !w.wrote {
		w.status = status
		w.wrote = true
	}
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusResponseWriter) Write(p []byte) (int, error) {
	if !w.wrote {
		w.status = http.StatusOK
		w.wrote = true
	}
	return w.ResponseWriter.Write(p)
}

// Flush forwards to the writer underneath. It exists so a streamed operation
// still reaches the client one part at a time: a wrapper that hides the
// http.Flusher turns every stream into a buffer that arrives at the end.
func (w *statusResponseWriter) Flush() {
	if !w.wrote {
		w.status = http.StatusOK
		w.wrote = true
	}
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// Unwrap returns the writer underneath so http.ResponseController still reaches
// it. Without this a stream cannot clear the server write deadline.
func (w *statusResponseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}
