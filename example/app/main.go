package main

import (
	"context"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/semconv"
	"log"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	otlpjson "github.com/madvikinggod/otlp-stream-exporter"
)

func setupOTLP(ctx context.Context, driver otlp.ProtocolDriver) func() {
	exporter, _ := otlp.NewExporter(ctx, driver)

	res := resource.NewWithAttributes(semconv.ServiceNameKey.String("Example Service"))

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSyncer(exporter),
		sdktrace.WithResource(res),
	)
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
	time.Sleep(time.Second)
	defer trace.End()
	_, trace2 := tracer.Start(ctx2, "innerTrace")
	defer trace2.End()
	time.Sleep(time.Second * 2)
	trace2.SetAttributes(attribute.String("inner-trace-label", "foos"))

	tracer2 := otel.Tracer("other Tracer")
	_, trace3 := tracer2.Start(ctx, "other Trace")
	time.Sleep(time.Second)
	trace3.End()

}
