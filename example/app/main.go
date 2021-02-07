package main

import (
	"context"
	"log"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/metric"
	metricsdk "go.opentelemetry.io/otel/sdk/export/metric"
	push "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	otlpjson "github.com/madvikinggod/otlp-stream-exporter"
)

func setupOTLP(ctx context.Context, driver otlp.ProtocolDriver) func() {
	exporter, _ := otlp.NewExporter(ctx, driver)

	tracerProvider := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter))
	processor := processor.New(simple.NewWithInexpensiveDistribution(), metricsdk.StatelessExportKindSelector())
	pusher := push.New(processor, push.WithPusher(exporter))
	pusher.Start(ctx)
	metricProvider := pusher.MeterProvider()

	otel.SetMeterProvider(metricProvider)
	otel.SetTracerProvider(tracerProvider)
	return func() {
		pusher.Stop(ctx)
		tracerProvider.Shutdown(ctx)
	}
}

func main() {
	metricFile, err := os.Create("metric.json")
	if err != nil {
		log.Fatalf("Could not open metric files %v", err)
	}
	traceFile, err := os.Create("trace.json")
	if err != nil {
		log.Fatalf("Could not open trace files %v", err)
	}
	driver := otlpjson.NewDriver(
		otlpjson.WithMetricWriter(metricFile),
		otlpjson.WithTraceWriter(traceFile),
		)

	ctx := context.Background()
	done := setupOTLP(ctx, driver)
	defer done()

	meter := metric.Must(otel.Meter("exampleMetrics"))
	counter1 := meter.NewInt64Counter("count1")
	counter1.Add(ctx,4, label.String("counter-label", "values"))
	counter2 := meter.NewInt64UpDownCounter("count2")
	counter2.Add(ctx, -4)

	tracer := otel.Tracer("exampleTracer")
	ctx2, trace := tracer.Start(ctx, "outerTrace")
	defer trace.End()
	_,trace2 := tracer.Start(ctx2, "innerTrace")
	defer trace2.End()
	trace2.SetAttributes(label.String("inner-trace-label", "foos"))
}
