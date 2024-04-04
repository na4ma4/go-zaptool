package zaptool

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	loggerCallerSkip    = 2
	chiLoggerCallerSkip = 2
)

const (
	HeaderUsername = "X-Logging-Username"
	HeaderNoop     = "X-Logging-Noop"
)

// ErrUnimplemented is returned when a method is unimplemented.
var ErrUnimplemented = errors.New("unimplemented method")

// loggingHandler is the http.Handler implementation for LoggingHandlerTo and its
// friends.
type loggingHandler struct {
	logger  *zap.Logger
	handler http.Handler
	opts    *loggingOptions
}

// ServeHTTP wraps the next handler ServeHTTP.
func (h loggingHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	t := time.Now()
	logger := makeLogger(w)
	url := *req.URL
	req.Header.Del(HeaderNoop)
	h.handler.ServeHTTP(logger, req)
	writeLog(&h, req, url, t, logger.Status(), logger.Size())
}

type loggingResponseWriter interface {
	http.ResponseWriter
	http.Flusher
	http.Pusher
	Status() int
	Size() int
}

func makeLogger(w http.ResponseWriter) loggingResponseWriter {
	return &responseLogger{w: w, status: http.StatusOK, size: 0}
}

// responseLogger is wrapper of http.ResponseWriter that keeps track of its HTTP
// status code and body size.
type responseLogger struct {
	w      http.ResponseWriter
	status int
	size   int
}

//nolint:wrapcheck // wrapping adds nothing.
func (l *responseLogger) Push(target string, opts *http.PushOptions) error {
	p, ok := l.w.(http.Pusher)
	if !ok {
		return ErrUnimplemented
	}

	return p.Push(target, opts)
}

func (l *responseLogger) Header() http.Header {
	return l.w.Header()
}

func (l *responseLogger) Write(b []byte) (int, error) {
	size, err := l.w.Write(b)
	l.size += size

	if err != nil {
		return size, fmt.Errorf("unable to write: %w", err)
	}

	return size, nil
}

func (l *responseLogger) WriteHeader(s int) {
	l.w.WriteHeader(s)
	l.status = s
}

func (l *responseLogger) Status() int {
	return l.status
}

func (l *responseLogger) Size() int {
	return l.size
}

func (l *responseLogger) Flush() {
	f, ok := l.w.(http.Flusher)
	if ok {
		f.Flush()
	}
}

func zapFieldOrSkip(returnField bool, field zapcore.Field) zapcore.Field {
	if returnField {
		return field
	}

	return zap.Skip()
}

// writeLog writes a log entry for req to w in Apache Combined Log Format.
// ts is the timestamp with which the entry should be logged.
// status and size are used to provide the response HTTP status and size.
func writeLog(lh *loggingHandler, req *http.Request, url url.URL, ts time.Time, status, size int) {
	if req.Header.Get(HeaderNoop) != "" {
		return
	}

	// Extract `X-Logging-Username` from request, added by authentication function earlier in process.
	username := "-"
	if req.Header.Get(HeaderUsername) != "" {
		username = sanitizeUsername(req.Header.Get(HeaderUsername))
	}

	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		host = req.RemoteAddr
	}

	uri := req.RequestURI
	if req.ProtoMajor == 2 && req.Method == "CONNECT" {
		uri = req.Host
	}

	if uri == "" {
		uri = url.RequestURI()
	}

	fields := []zapcore.Field{
		zap.Namespace("http"),            // 0
		zap.String("host", host),         // 1
		zap.String("username", username), // 2
		zapFieldOrSkip(lh.opts.includeTimestamp, zap.String("timestamp", ts.Format(time.RFC3339Nano))), // 3
		zap.String("method", req.Method),                                                    // 4
		zap.String("uri", sanitizeURI(uri)),                                                 // 5
		zap.String("proto", req.Proto),                                                      // 6
		zap.Int("status", status),                                                           // 7
		zap.Int("size", size),                                                               // 8
		zap.String("referer", sanitizeURI(req.Referer())),                                   // 9
		zap.String("user-agent", sanitizeUserAgent(req.UserAgent())),                        // 10
		zapFieldOrSkip(lh.opts.includeTiming, zap.Duration("request-time", time.Since(ts))), // 11
		zapFieldOrSkip(lh.opts.includeXForwardedFor, zap.String("forwarded_for", req.Header.Get("X-Forwarded-For"))), // 12
	}

	lh.logger.Info(
		"Request",
		fields...,
	)
}

// LoggingHTTPHandler return a http.Handler that wraps h and logs requests to out using
// a *zap.Logger.
func LoggingHTTPHandler(logger *zap.Logger, httpHandler http.Handler, opts ...loggingOptionsFunc) http.Handler {
	opt := &loggingOptions{
		includeTiming:        true,
		includeTimestamp:     true,
		includeXForwardedFor: false,
	}

	for _, f := range opts {
		f(opt)
	}

	logger = logger.WithOptions(zap.AddCallerSkip(loggerCallerSkip))

	return loggingHandler{
		logger,
		httpHandler,
		opt,
	}
}

func LoggingHTTPHandlerWrapper(logger *zap.Logger, opts ...loggingOptionsFunc) func(next http.Handler) http.Handler {
	opt := &loggingOptions{
		includeTiming:        true,
		includeTimestamp:     true,
		includeXForwardedFor: false,
	}

	for _, f := range opts {
		f(opt)
	}

	logger = logger.WithOptions(zap.AddCallerSkip(chiLoggerCallerSkip))

	return func(next http.Handler) http.Handler {
		return loggingHandler{
			logger,
			next,
			opt,
		}
	}
}
