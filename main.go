package main

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/jpatel531/otel-test/middlewares"

	"github.com/go-chi/chi/v5"
	"github.com/jpatel531/otel-test/config"
	"github.com/jpatel531/otel-test/telemetry"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	ctx := logger.WithContext(context.Background())

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal().Err(err).Msg("loading config")
	}
	logger = enrichLogger(cfg, logger)

	shutdown := telemetry.Init(ctx, cfg)
	defer shutdown(ctx)

	r := chi.NewRouter()
	r.Use(otelhttp.NewMiddleware(""))
	r.Use(middlewares.TraceID())
	r.Use(middlewares.Logging(logger))

	r.Get("/", rootHandler())

	server := &http.Server{Handler: r, Addr: cfg.Addr}

	logger.Info().Msgf("listening and serving on %s", cfg.Addr)
	if err := server.ListenAndServe(); err != nil {
		logger.Fatal().Err(err).Msg("listening")
	}
}

func rootHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		errMsg := req.URL.Query().Get("error")
		if errMsg != "" {
			span := trace.SpanFromContext(ctx)
			span.SetStatus(codes.Error, "root failed")
			span.RecordError(errors.New(errMsg))
			return
		}

		_, _ = w.Write([]byte("ok"))
	}
}

func enrichLogger(cfg config.Config, logger zerolog.Logger) zerolog.Logger {
	hostname, _ := os.Hostname()
	return logger.With().
		Str("environment", cfg.Environment).
		Str("service", cfg.ServiceName).
		Str("hostname", hostname).
		Logger()
}
