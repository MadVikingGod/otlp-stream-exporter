module example

go 1.14

replace github.com/madvikinggod/otlp-stream-exporter v0.1.0 => ../

replace go.opentelemetry.io/proto/otlp v0.16.0 => ../internal/otlp

require (
	github.com/golang/protobuf v1.4.3
	github.com/madvikinggod/otlp-stream-exporter v0.1.0
	go.opentelemetry.io/otel v0.16.0
	go.opentelemetry.io/otel/exporters/otlp v0.16.0
	go.opentelemetry.io/otel/sdk v0.16.0
	go.opentelemetry.io/proto/otlp v0.16.0
	google.golang.org/grpc v1.35.0
)
