package otlp_stream_exporter

import (
	"context"
	"io"

	colmetricspb "go.opentelemetry.io/proto/otlp/collector/metrics/v1"
	coltracepb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	"github.com/gogo/protobuf/jsonpb"
	"go.opentelemetry.io/otel/exporters/otlp"
	metricsdk "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/trace"
	"github.com/madvikinggod/otlp-stream-exporter/internal/transform"
)

type driver struct {
	marshler *jsonpb.Marshaler
	metric io.Writer
	trace io.Writer
}

var _ otlp.ProtocolDriver = (*driver)(nil)

func NewDriver(opts ...Option) otlp.ProtocolDriver {
	cfg := defaultConfig()
	cfg.Apply(opts...)
	return &driver{
		marshler: &jsonpb.Marshaler{},
		metric: cfg.metric,
		trace: cfg.trace,
	}
}

func (d *driver) Start(ctx context.Context) error {
	return nil
}

func (d *driver) Stop(ctx context.Context) error {
	return nil
}

func (d *driver) ExportMetrics(ctx context.Context, cps metricsdk.CheckpointSet, selector metricsdk.ExportKindSelector) error {
	rms, err := transform.CheckpointSet(ctx, selector, cps, 1)
	if err != nil {
		return err
	}
	if len(rms) == 0 {
		return nil
	}
	pbRequest := &colmetricspb.ExportMetricsServiceRequest{
		ResourceMetrics: rms,
	}

	return d.marshler.Marshal(d.metric, pbRequest)
}

func (d *driver) ExportTraces(ctx context.Context, ss []*trace.SpanSnapshot) error {
	protoSpans := transform.SpanData(ss)
	if len(protoSpans) == 0 {
		return nil
	}
	pbRequest := &coltracepb.ExportTraceServiceRequest{
		ResourceSpans: protoSpans,
	}
	return d.marshler.Marshal(d.trace, pbRequest)
}
