package telemetry

import (
	"context"

	"github.com/jpatel531/otel-test/config"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"google.golang.org/grpc"
)

func Init(ctx context.Context, cfg config.Config) func(ctx context.Context) {
	if cfg.OtelEndpoint == "" {
		return func(ctx context.Context) {}
	}

	logger := zerolog.Ctx(ctx)

	tp, err := initTracer(ctx, cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("initialising tracer")
	}

	mp, err := initMeter(ctx, cfg.OtelEndpoint)
	if err != nil {
		logger.Fatal().Err(err).Msg("initialising meter")
	}

	return func(ctx context.Context) {
		if err := tp.Shutdown(ctx); err != nil {
			logger.Error().Err(err).Msgf("shutting down tracer provider")
		}

		if err := mp.Shutdown(ctx); err != nil {
			logger.Error().Err(err).Msg("shutting down meter provider")
		}
	}
}

func initTracer(ctx context.Context, cfg config.Config) (*sdktrace.TracerProvider, error) {
	traceClient := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(cfg.OtelEndpoint),
		otlptracegrpc.WithDialOption(grpc.WithBlock()))
	exp, err := otlptrace.New(ctx, traceClient)
	if err != nil {
		return nil, err
	}

	// For the demonstration, use sdktrace.AlwaysSample sampler to sample all traces.
	// In a production application, use sdktrace.ProbabilitySampler with a desired probability.
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(getResource(cfg)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, err
}

func getResource(cfg config.Config) *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(cfg.ServiceName),
		semconv.DeploymentEnvironment(cfg.Environment),
	)
}

func initMeter(ctx context.Context, endpoint string) (*sdkmetric.MeterProvider, error) {
	exp, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(endpoint),
	)
	if err != nil {
		return nil, err
	}

	mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exp)))
	otel.SetMeterProvider(mp)
	return mp, nil
}
