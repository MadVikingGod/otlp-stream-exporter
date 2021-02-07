module github.com/madvikinggod/otlp-stream-exporter

go 1.15

replace (
	go.opentelemetry.io/otel/exporters/otlp/internal/transform v0.16.0 => ./internal/transform
	go.opentelemetry.io/proto/otlp v0.16.0 => ./internal/otlp
)

require (
	go.opentelemetry.io/otel v0.16.0
	go.opentelemetry.io/otel/exporters/otlp v0.16.0
	go.opentelemetry.io/otel/sdk v0.16.0
	go.opentelemetry.io/proto/otlp v0.16.0
)
