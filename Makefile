collector:
	docker run \
        -p 4317:4317 \
        -e DD_API_KEY=$DD_API_KEY  \
        --hostname $(hostname) \
        -v $(pwd)/collector.yaml:/etc/otelcol-contrib/config.yaml \
        otel/opentelemetry-collector-contrib