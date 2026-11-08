package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.deployport.com/specular-runtime/client"
)

const (
	testTelemetryService  = "machinery"
	testTelemetryPrefix   = "ljnyc"
	testTelemetryWriteURL = "http://vm.test:8428/"
	testTelemetryEndpoint = "http://vm.test:8428/api/v1/import/prometheus"
)

// telemetryClock is moved by hand so a recorded duration is exact. A test that
// slept for 300ms could not pin a bucket boundary.
type telemetryClock struct {
	mu sync.Mutex
	at time.Time
}

func newTelemetryClock() *telemetryClock {
	return &telemetryClock{at: time.Date(2026, time.November, 7, 12, 0, 0, 0, time.UTC)}
}

func (c *telemetryClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.at
}

func (c *telemetryClock) advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.at = c.at.Add(d)
}

// telemetryPush is one attempt to reach the store.
type telemetryPush struct {
	Method      string
	URL         string
	ContentType string
	Body        string
	Username    string
	Password    string
	HasAuth     bool
}

// telemetryTransport stands in for the metrics store. It records every attempt,
// including the ones it fails, so a test can tell "not pushed" from "pushed and
// refused".
type telemetryTransport struct {
	mu       sync.Mutex
	pushes   []telemetryPush
	failures int
	status   int
}

func (t *telemetryTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	username, password, hasAuth := r.BasicAuth()
	t.mu.Lock()
	t.pushes = append(t.pushes, telemetryPush{
		Method:      r.Method,
		URL:         r.URL.String(),
		ContentType: r.Header.Get("Content-Type"),
		Body:        string(body),
		Username:    username,
		Password:    password,
		HasAuth:     hasAuth,
	})
	fail := t.failures > 0
	if fail {
		t.failures--
	}
	status := t.status
	if status == 0 {
		status = http.StatusNoContent
	}
	t.mu.Unlock()

	if fail {
		return nil, errors.New("metrics store unreachable")
	}
	return &http.Response{
		StatusCode: status,
		Status:     fmt.Sprintf("%d %s", status, http.StatusText(status)),
		Header:     http.Header{},
		Body:       io.NopCloser(strings.NewReader("")),
		Request:    r,
	}, nil
}

func (t *telemetryTransport) failNext(count int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.failures = count
}

func (t *telemetryTransport) PushList() []telemetryPush {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]telemetryPush, len(t.pushes))
	copy(out, t.pushes)
	return out
}

// Len is how many pushes were attempted, read without copying the bodies.
func (t *telemetryTransport) Len() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.pushes)
}

// flushTelemetry flushes until the interval has carried `want` pushes.
//
// The observation is written after ServeHTTP returns, so a client can finish a
// response before the server has recorded it. One flush can then carry nothing
// while the request lands in the next interval, which is the designed cost of
// swapping the store before the render rather than locking every request.
func flushTelemetry(t *testing.T, svc *Service, transport *telemetryTransport, want int) {
	t.Helper()
	require.Eventually(t, func() bool {
		_ = svc.telemetry.flush()
		return transport.Len() >= want
	}, 2*time.Second, time.Millisecond,
		"no interval carried the %d request(s) the test made", want)
}

// telemetryWarnLog records the messages a Warn produced.
type telemetryWarnLog struct {
	mu    sync.Mutex
	lines []string
}

func (l *telemetryWarnLog) Enabled(context.Context, slog.Level) bool { return true }

func (l *telemetryWarnLog) Handle(_ context.Context, record slog.Record) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if record.Level == slog.LevelWarn {
		l.lines = append(l.lines, record.Message)
	}
	return nil
}

func (l *telemetryWarnLog) WithAttrs([]slog.Attr) slog.Handler { return l }

func (l *telemetryWarnLog) WithGroup(string) slog.Handler { return l }

func (l *telemetryWarnLog) warnings() []string {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]string, len(l.lines))
	copy(out, l.lines)
	return out
}

// newMeasuredService builds a service that measures into transport, on a clock
// the test moves by hand.
func newMeasuredService(
	t *testing.T,
	pk *client.Package,
	clock *telemetryClock,
	transport *telemetryTransport,
	opts ...ServiceOptionFunc,
) (*Service, *httptest.Server) {
	t.Helper()
	configured := []ServiceOptionFunc{
		WithServicePackage(pk),
		WithTelemetry(Telemetry{
			Service:  testTelemetryService,
			WriteURL: testTelemetryWriteURL,
			Username: "deployport",
			Password: "secret",
			Prefix:   testTelemetryPrefix,
			// Long enough that only the test's own flush can push.
			Every:      time.Hour,
			HTTPClient: &http.Client{Transport: transport},
			Now:        clock.Now,
		}),
	}
	configured = append(configured, opts...)
	svc := NewService(configured...)
	server := httptest.NewServer(svc)
	t.Cleanup(server.Close)
	return svc, server
}

// bookCreateTaking answers Book.Create after moving the clock by a fixed
// duration.
func bookCreateTaking(clock *telemetryClock, took time.Duration) ServiceOptionFunc {
	fixed := &atomic.Int64{}
	fixed.Store(int64(took))
	return bookCreateTakingFrom(clock, fixed)
}

// bookCreateTakingFrom answers Book.Create after moving the clock by whatever
// the test set, so one service can serve several durations in a row.
func bookCreateTakingFrom(clock *telemetryClock, took *atomic.Int64) ServiceOptionFunc {
	return WithOperationHandler(OperationHandlerFunc(func(ctx context.Context, op *OperationExecution) (client.Struct, error) {
		clock.advance(time.Duration(took.Load()))
		return &BookCreateOutput{ID: "123"}, nil
	}))
}

func postBookCreate(t *testing.T, server *httptest.Server) int {
	t.Helper()
	response, err := server.Client().Post(
		server.URL+"/Book/Create",
		"application/json",
		strings.NewReader(`{"name":"foo"}`),
	)
	require.NoError(t, err)
	defer response.Body.Close()
	_, _ = io.Copy(io.Discard, response.Body)
	return response.StatusCode
}

// telemetrySamples reads a Prometheus text body into "name{labels}" to value,
// and fails when a line is not a name, a value and an integer millisecond
// timestamp.
func telemetrySamples(t *testing.T, body string) map[string]float64 {
	t.Helper()
	samples := map[string]float64{}
	for _, line := range strings.Split(strings.TrimSpace(body), "\n") {
		if line == "" {
			continue
		}
		fields := strings.SplitN(line, " ", 3)
		require.Len(t, fields, 3, "a sample is a name, a value and a timestamp: %q", line)
		value, err := strconv.ParseFloat(fields[1], 64)
		require.NoError(t, err, "value of %q", line)
		_, err = strconv.ParseInt(fields[2], 10, 64)
		require.NoError(t, err, "a timestamp is integer milliseconds: %q", line)
		samples[fields[0]] = value
	}
	return samples
}

// telemetryKey builds the sample name a test expects, for one of the three
// labels that may travel plus le.
func telemetryKey(name string, operation string, class string, le string) string {
	labels := `service="` + testTelemetryService +
		`",operation="` + operation +
		`",status_class="` + class + `"`
	if le != "" {
		labels += `,le="` + le + `"`
	}
	return name + "{" + labels + "}"
}

func telemetryRequestBase() string {
	return testTelemetryPrefix + "_specular_operation_duration_seconds"
}

func telemetryStreamBase() string {
	return testTelemetryPrefix + "_specular_stream_open_seconds"
}

// One 2xx request that took 300ms: one observation, in the 0.5 bucket, sum 0.3,
// count 1. The histogram count is the request count, and there is no second
// counter to disagree with it.
func TestAnEndpointRecordsItsDuration(t *testing.T) {
	pk, err := testPackage()
	require.NoError(t, err)
	clock := newTelemetryClock()
	transport := &telemetryTransport{}
	svc, server := newMeasuredService(t, pk, clock, transport, bookCreateTaking(clock, 300*time.Millisecond))

	require.Equal(t, http.StatusOK, postBookCreate(t, server))
	flushTelemetry(t, svc, transport, 1)

	pushes := transport.PushList()
	require.Len(t, pushes, 1)
	push := pushes[0]
	require.Equal(t, http.MethodPost, push.Method)
	require.Equal(t, testTelemetryEndpoint, push.URL)
	require.Equal(t, "text/plain", push.ContentType)
	require.True(t, push.HasAuth)
	require.Equal(t, "deployport", push.Username)
	require.Equal(t, "secret", push.Password)

	base := telemetryRequestBase()
	samples := telemetrySamples(t, push.Body)
	require.Equal(t, 0.0, samples[telemetryKey(base+"_bucket", "Book.Create", "2xx", "0.25")])
	require.Equal(t, 1.0, samples[telemetryKey(base+"_bucket", "Book.Create", "2xx", "0.5")])
	require.Equal(t, 1.0, samples[telemetryKey(base+"_count", "Book.Create", "2xx", "")])
	require.InDelta(t, 0.3, samples[telemetryKey(base+"_sum", "Book.Create", "2xx", "")], 1e-9)
	require.NotContains(t, push.Body, "specular_operation_total",
		"the count already is the request count, so a second counter would repeat it")
}

// A duration exactly on a bound belongs to that bound, which is what le means.
func TestABucketBoundaryMatchesPrometheus(t *testing.T) {
	cases := []struct {
		name     string
		took     time.Duration
		bucket   string
		belowOne string
	}{
		{name: "the smallest bound", took: 5 * time.Millisecond, bucket: "0.005"},
		{name: "exactly a quarter second", took: 250 * time.Millisecond, bucket: "0.25", belowOne: "0.1"},
		{name: "just over a quarter second", took: 251 * time.Millisecond, bucket: "0.5", belowOne: "0.25"},
		{name: "exactly half a second", took: 500 * time.Millisecond, bucket: "0.5", belowOne: "0.25"},
		{name: "exactly ten seconds", took: 10 * time.Second, bucket: "10", belowOne: "5"},
		{name: "past the last bound", took: 11 * time.Second, bucket: "+Inf", belowOne: "10"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			pk, err := testPackage()
			require.NoError(t, err)
			clock := newTelemetryClock()
			transport := &telemetryTransport{}
			svc, server := newMeasuredService(t, pk, clock, transport, bookCreateTaking(clock, c.took))

			postBookCreate(t, server)
			flushTelemetry(t, svc, transport, 1)
			pushes := transport.PushList()
			require.Len(t, pushes, 1)

			base := telemetryRequestBase()
			samples := telemetrySamples(t, pushes[0].Body)
			require.Equal(t, 1.0, samples[telemetryKey(base+"_bucket", "Book.Create", "2xx", c.bucket)],
				"the observation belongs to this bound")
			if c.belowOne != "" {
				require.Equal(t, 0.0, samples[telemetryKey(base+"_bucket", "Book.Create", "2xx", c.belowOne)],
					"the bound under it does not count it")
			}
			require.Equal(t, 1.0, samples[telemetryKey(base+"_count", "Book.Create", "2xx", "")])
		})
	}
}

// A streamed operation is measured as an open duration and never enters the
// request histogram, because one long watch would move every percentile.
func TestAStreamedOperationIsNotARequest(t *testing.T) {
	pk, op := streamTestPackage(t)
	clock := newTelemetryClock()
	transport := &telemetryTransport{}
	svc, server := newMeasuredService(t, pk, clock, transport,
		WithStreamHandler(StreamHandlerFunc(func(ctx context.Context, opx *OperationExecution) (<-chan client.StreamEvent[client.Struct], error) {
			clock.advance(30 * time.Second)
			out := make(chan client.StreamEvent[client.Struct], 1)
			out <- client.StreamEvent[client.Struct]{Output: &BookCreateOutput{ID: "a"}}
			close(out)
			return out, nil
		})),
	)

	require.Equal(t, []string{"a"}, drain(t, server, op))
	flushTelemetry(t, svc, transport, 1)

	pushes := transport.PushList()
	require.Len(t, pushes, 1)
	body := pushes[0].Body
	base := telemetryStreamBase()
	samples := telemetrySamples(t, body)
	require.Equal(t, 1.0, samples[telemetryKey(base+"_count", "Book.Create", "2xx", "")])
	require.InDelta(t, 30.0, samples[telemetryKey(base+"_sum", "Book.Create", "2xx", "")], 1e-9)
	require.Equal(t, 0.0, samples[telemetryKey(base+"_bucket", "Book.Create", "2xx", "10")])
	require.Equal(t, 1.0, samples[telemetryKey(base+"_bucket", "Book.Create", "2xx", "+Inf")])

	require.NotContains(t, body, "specular_operation_duration_seconds",
		"a stream must never enter the request histogram")
	require.NotContains(t, body, "specular_operation_total")
}

// A request refused by authentication is still measured, against the operation
// it named.
func TestAnAuthenticationFailureIsRecordedAgainstTheOperation(t *testing.T) {
	pk, err := testPackage()
	require.NoError(t, err)
	clock := newTelemetryClock()
	transport := &telemetryTransport{}
	svc, server := newMeasuredService(t, pk, clock, transport,
		WithAuthentication(func(params *AuthenticationParams) error {
			return &AuthenticationError{Message: "no identity in this request"}
		}),
		bookCreateTaking(clock, 100*time.Millisecond),
	)

	require.Equal(t, http.StatusForbidden, postBookCreate(t, server))
	flushTelemetry(t, svc, transport, 1)

	pushes := transport.PushList()
	require.Len(t, pushes, 1)
	base := telemetryRequestBase()
	samples := telemetrySamples(t, pushes[0].Body)
	require.Equal(t, 1.0, samples[telemetryKey(base+"_count", "Book.Create", "4xx", "")])
	require.NotContains(t, samples, telemetryKey(base+"_count", "Book.Create", "2xx", ""))
}

// A path that names no operation creates no series, so a scanner cannot mint
// series by trying.
func TestAnUnmatchedPathRecordsNothing(t *testing.T) {
	pk, err := testPackage()
	require.NoError(t, err)
	clock := newTelemetryClock()
	transport := &telemetryTransport{}
	svc, server := newMeasuredService(t, pk, clock, transport, bookCreateTaking(clock, 100*time.Millisecond))

	response, err := server.Client().Post(server.URL+"/Nope/Missing", "application/json", nil)
	require.NoError(t, err)
	defer response.Body.Close()
	require.Equal(t, http.StatusNotFound, response.StatusCode)

	require.Equal(t, 0, svc.telemetry.current.Load().count)
	require.NoError(t, svc.telemetry.flush())
	require.Empty(t, transport.PushList())
}

// Without a write URL there is no recorder, so there is no goroutine and no
// push, whatever the interval says.
func TestADisabledTelemetryStartsNothing(t *testing.T) {
	require.Nil(t, newTelemetryRecorder(Telemetry{Service: testTelemetryService, Prefix: testTelemetryPrefix}),
		"a write URL is what turns the recorder on")

	pk, err := testPackage()
	require.NoError(t, err)
	transport := &telemetryTransport{}
	svc := NewService(
		WithServicePackage(pk),
		bookCreateTaking(newTelemetryClock(), 100*time.Millisecond),
		WithTelemetry(Telemetry{
			Service:  testTelemetryService,
			Prefix:   testTelemetryPrefix,
			Username: "deployport",
			Password: "secret",
			// Short on purpose: a started loop would push several times inside
			// the wait below.
			Every:      10 * time.Millisecond,
			HTTPClient: &http.Client{Transport: transport},
		}),
	)
	require.Nil(t, svc.telemetry)

	server := httptest.NewServer(svc)
	defer server.Close()
	require.Equal(t, http.StatusOK, postBookCreate(t, server))
	time.Sleep(100 * time.Millisecond)
	require.Empty(t, transport.PushList())
}

// A flush renders the buckets cumulatively, the way le reads.
func TestAFlushRendersCumulativeBuckets(t *testing.T) {
	pk, err := testPackage()
	require.NoError(t, err)
	clock := newTelemetryClock()
	transport := &telemetryTransport{}
	took := &atomic.Int64{}
	took.Store(int64(10 * time.Millisecond))
	svc, server := newMeasuredService(t, pk, clock, transport, bookCreateTakingFrom(clock, took))

	for _, d := range []time.Duration{10 * time.Millisecond, 100 * time.Millisecond, time.Second} {
		took.Store(int64(d))
		postBookCreate(t, server)
	}
	flushTelemetry(t, svc, transport, 1)

	pushes := transport.PushList()
	require.Len(t, pushes, 1)
	base := telemetryRequestBase()
	samples := telemetrySamples(t, pushes[0].Body)
	for le, want := range map[string]float64{
		"0.005": 0,
		"0.01":  1,
		"0.025": 1,
		"0.05":  1,
		"0.1":   2,
		"0.25":  2,
		"0.5":   2,
		"1":     3,
		"2.5":   3,
		"5":     3,
		"10":    3,
		"+Inf":  3,
	} {
		require.Equal(t, want, samples[telemetryKey(base+"_bucket", "Book.Create", "2xx", le)], "le %s", le)
	}
	require.Equal(t, 3.0, samples[telemetryKey(base+"_count", "Book.Create", "2xx", "")])
	require.InDelta(t, 1.11, samples[telemetryKey(base+"_sum", "Book.Create", "2xx", "")], 1e-9)
}

// A second flush carries only what arrived since the first, and never the
// running total.
func TestAFlushCarriesOnlyTheNewObservations(t *testing.T) {
	pk, err := testPackage()
	require.NoError(t, err)
	clock := newTelemetryClock()
	transport := &telemetryTransport{}
	took := &atomic.Int64{}
	took.Store(int64(100 * time.Millisecond))
	svc, server := newMeasuredService(t, pk, clock, transport, bookCreateTakingFrom(clock, took))
	base := telemetryRequestBase()

	postBookCreate(t, server)
	flushTelemetry(t, svc, transport, 1)
	pushes := transport.PushList()
	require.Len(t, pushes, 1)
	first := telemetrySamples(t, pushes[0].Body)
	require.Equal(t, 1.0, first[telemetryKey(base+"_count", "Book.Create", "2xx", "")])
	require.InDelta(t, 0.1, first[telemetryKey(base+"_sum", "Book.Create", "2xx", "")], 1e-9)
	require.Equal(t, 1.0, first[telemetryKey(base+"_bucket", "Book.Create", "2xx", "0.1")])

	// The timestamp is the millisecond the flush ran, not the request time.
	require.Contains(t, pushes[0].Body, " "+strconv.FormatInt(clock.Now().UnixMilli(), 10)+"\n",
		"every sample carries a millisecond timestamp")

	took.Store(int64(250 * time.Millisecond))
	postBookCreate(t, server)
	flushTelemetry(t, svc, transport, 2)
	pushes = transport.PushList()
	require.Len(t, pushes, 2)
	second := telemetrySamples(t, pushes[1].Body)
	require.Equal(t, 1.0, second[telemetryKey(base+"_count", "Book.Create", "2xx", "")],
		"the second interval counts only its own request")
	require.InDelta(t, 0.25, second[telemetryKey(base+"_sum", "Book.Create", "2xx", "")], 1e-9)
	require.Equal(t, 0.0, second[telemetryKey(base+"_bucket", "Book.Create", "2xx", "0.1")],
		"the first interval must not travel twice")
	require.Equal(t, 1.0, second[telemetryKey(base+"_bucket", "Book.Create", "2xx", "0.25")])
}

// A part has to reach the client while the handler is still running. The
// telemetry wrapper is what the stream loop writes through, and a wrapper that
// hides http.Flusher turns every stream into one buffer that arrives at the
// end, with no error anywhere and nothing to show it.
//
// The parts are sent before the handler waits, and the wait is two seconds,
// well under the five second heartbeat, so a heartbeat cannot do the flush for
// us.
func TestAStreamPartReachesTheClientBeforeTheHandlerEnds(t *testing.T) {
	pk, op := streamTestPackage(t)
	acknowledged := make(chan struct{})
	var once sync.Once
	release := func() { once.Do(func() { close(acknowledged) }) }
	defer release()

	_, server := newMeasuredService(t, pk, newTelemetryClock(), &telemetryTransport{},
		WithStreamHandler(StreamHandlerFunc(func(ctx context.Context, opx *OperationExecution) (<-chan client.StreamEvent[client.Struct], error) {
			out := make(chan client.StreamEvent[client.Struct])
			go func() {
				defer close(out)
				for _, id := range []string{"a", "b"} {
					select {
					case out <- client.StreamEvent[client.Struct]{Output: &BookCreateOutput{ID: id}}:
					case <-ctx.Done():
						return
					}
				}
				select {
				case <-acknowledged:
				case <-ctx.Done():
					return
				}
			}()
			return out, nil
		})),
	)

	transport, err := client.NewHTTPJSONTransport(server.URL, client.WithHTTPJSONTransportClient(server.Client()))
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	events, err := transport.Stream(ctx, &client.Request{Operation: op, Input: &BookCreateInput{}})
	require.NoError(t, err)

	select {
	case ev, ok := <-events:
		require.True(t, ok)
		require.NoError(t, ev.Err)
		require.Equal(t, "a", ev.Output.(*BookCreateOutput).ID)
	case <-time.After(2 * time.Second):
		t.Fatal("the first part did not arrive while the handler was still running")
	}

	release()
	var rest []string
	for ev := range events {
		require.NoError(t, ev.Err)
		rest = append(rest, ev.Output.(*BookCreateOutput).ID)
	}
	require.Equal(t, []string{"b"}, rest)
}

// A store that refuses one push is logged and does not stop the next interval.
func TestAPushFailureDoesNotStopTheNextFlush(t *testing.T) {
	pk, err := testPackage()
	require.NoError(t, err)
	clock := newTelemetryClock()
	transport := &telemetryTransport{}
	warns := &telemetryWarnLog{}
	took := &atomic.Int64{}
	took.Store(int64(100 * time.Millisecond))
	svc, server := newMeasuredService(t, pk, clock, transport,
		bookCreateTakingFrom(clock, took),
		WithServiceLogger(slog.New(warns)),
	)

	transport.failNext(1)
	postBookCreate(t, server)
	require.Eventually(t, func() bool {
		svc.telemetry.flushAndLog()
		return transport.Len() >= 1
	}, 2*time.Second, time.Millisecond, "the push was never attempted")
	require.Len(t, transport.PushList(), 1, "the push was attempted")
	require.Len(t, warns.warnings(), 1)
	require.Contains(t, warns.warnings()[0], "metrics push failed")

	took.Store(int64(250 * time.Millisecond))
	postBookCreate(t, server)
	require.Eventually(t, func() bool {
		svc.telemetry.flushAndLog()
		return transport.Len() >= 2
	}, 2*time.Second, time.Millisecond, "the next interval never fired")
	pushes := transport.PushList()
	require.Len(t, pushes, 2, "the next interval still fires")
	require.Len(t, warns.warnings(), 1, "a healthy push says nothing")
	base := telemetryRequestBase()
	second := telemetrySamples(t, pushes[1].Body)
	require.Equal(t, 1.0, second[telemetryKey(base+"_count", "Book.Create", "2xx", "")])
}

// The cap drops new keys and says so once, so a flood cannot grow the store
// without bound and cannot fill the log either.
func TestTheKeyCapDropsNewKeysWithOneWarn(t *testing.T) {
	warns := &telemetryWarnLog{}
	recorder := newTelemetryRecorder(Telemetry{
		Service:  testTelemetryService,
		WriteURL: testTelemetryWriteURL,
		Prefix:   testTelemetryPrefix,
		Every:    time.Hour,
	})
	require.NotNil(t, recorder)
	recorder.logger = slog.New(warns)

	for i := 0; i < telemetryMaxKeys; i++ {
		key := fmt.Sprintf("Op.%d%s2xx", i, telemetryKeySeparator)
		require.NotNil(t, recorder.seriesFor(key, false), "key %d is inside the cap", i)
	}
	require.Equal(t, telemetryMaxKeys, recorder.current.Load().count)
	require.Empty(t, warns.warnings(), "the cap was not reached yet")

	for i := 0; i < 8; i++ {
		key := fmt.Sprintf("Extra.%d%s2xx", i, telemetryKeySeparator)
		require.Nil(t, recorder.seriesFor(key, false), "key %d is past the cap", i)
	}
	require.Len(t, warns.warnings(), 1, "one line, not one per request")
	require.Equal(t, telemetryMaxKeys, recorder.current.Load().count)

	require.NotNil(t, recorder.seriesFor(fmt.Sprintf("Op.0%s2xx", telemetryKeySeparator), false),
		"a key already known still records")
}

// Only three labels carry data. An identity anywhere in the request must not
// reach the series.
func TestALabelNeverCarriesACustomerIdentifier(t *testing.T) {
	const identity = "acc_9f2/user_7/mach_kxpkzlq4qf"

	pk, err := testPackage()
	require.NoError(t, err)
	clock := newTelemetryClock()
	transport := &telemetryTransport{}

	type identityKey struct{}
	svc, server := newMeasuredService(t, pk, clock, transport,
		WithAuthentication(func(params *AuthenticationParams) error {
			params.Request = params.Request.WithContext(
				context.WithValue(params.Request.Context(), identityKey{}, identity),
			)
			return nil
		}),
		bookCreateTaking(clock, 100*time.Millisecond),
	)

	postBookCreate(t, server)
	flushTelemetry(t, svc, transport, 1)
	pushes := transport.PushList()
	require.Len(t, pushes, 1)
	body := pushes[0].Body
	require.NotContains(t, body, identity)

	allowed := map[string]bool{"service": true, "operation": true, "status_class": true, "le": true}
	namePattern := regexp.MustCompile(`([a-zA-Z_][a-zA-Z0-9_]*)="`)
	lines := 0
	for _, line := range strings.Split(strings.TrimSpace(body), "\n") {
		if line == "" {
			continue
		}
		lines++
		for _, match := range namePattern.FindAllStringSubmatch(line, -1) {
			require.True(t, allowed[match[1]], "label %q may not travel: %q", match[1], line)
		}
	}
	require.Greater(t, lines, 0, "the guard has to read a body to mean anything")
}
