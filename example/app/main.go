package main

import (
	"context"
	"log"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/label"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	otlpjson "github.com/madvikinggod/otlp-stream-exporter"
)

func setupOTLP(ctx context.Context, driver otlp.ProtocolDriver) func() {
	exporter, _ := otlp.NewExporter(ctx, driver)

	tracerProvider := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter))
	otel.SetTracerProvider(tracerProvider)

	return func() {
		tracerProvider.Shutdown(ctx)
	}
}

func main() {
	traceFile, err := os.Create("trace.json")
	if err != nil {
		log.Fatalf("Could not open trace files %v", err)
	}
	driver := otlpjson.NewDriver(
		otlpjson.WithTraceWriter(traceFile),
	)

	ctx := context.Background()
	done := setupOTLP(ctx, driver)
	defer done()

	tracer := otel.Tracer("exampleTracer")
	ctx2, trace := tracer.Start(ctx, "outerTrace")
	defer trace.End()
	_, trace2 := tracer.Start(ctx2, "innerTrace")
	defer trace2.End()
	trace2.SetAttributes(label.String("inner-trace-label", "foos"))
}
