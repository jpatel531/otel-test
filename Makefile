collector:
	docker run \
        -p 4317:4317 \
        -e DD_API_KEY=$$DD_API_KEY \
        --hostname=`hostname` \
        -v $(CURDIR)/collector.yaml:/etc/otelcol-contrib/config.yaml \
		-v $(CURDIR)/log:/var/log/app.log \
        otel/opentelemetry-collector-contrib
