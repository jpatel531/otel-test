package middlewares

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

func Logging(logger zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			start := time.Now()
			ctx := req.Context()

			spctx := trace.SpanContextFromContext(ctx)
			traceID := spctx.TraceID().String()
			logger = logger.With().
				Str("method", req.Method).
				Str("remote_addr", req.RemoteAddr).
				Str("url", req.URL.Path).
				Str("user_agent", req.UserAgent()).
				Str("trace_id", traceID).
				Str("span_id", spctx.SpanID().String()).
				Logger()

			ww := middleware.NewWrapResponseWriter(w, req.ProtoMajor)
			defer func() {
				logger.Info().
					Int("bytes", ww.BytesWritten()).
					Dur("duration", time.Now().Sub(start)).
					Int("status_code", ww.Status()).
					Msgf("%s %s", req.Method, req.URL.String())
			}()

			req = req.WithContext(logger.WithContext(ctx))
			next.ServeHTTP(ww, req)
		})
	}
}

func TraceID() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			spctx := trace.SpanContextFromContext(req.Context())
			if spctx.HasTraceID() {
				w.Header().Set("X-Trace-ID", spctx.TraceID().String())
			}

			next.ServeHTTP(w, req)
		})
	}
}
