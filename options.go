package otlp_stream_exporter

import (
	"io"
	"os"
)



type Option interface {
	Apply(*config)
}

type config struct {
	metric, trace io.Writer
}

func defaultConfig() *config {
	return &config{
		metric: os.Stdout,
		trace: os.Stdout,
	}
}

func (cfg *config) Apply(opts ...Option) {
	for _,opt := range opts {
		opt.Apply(cfg)
	}
}

type metricOption struct {
	io.Writer
}
func (o metricOption) Apply(cfg *config) {
	cfg.metric = io.Writer(o)
}
func WithMetricWriter(w io.Writer) metricOption {
	return metricOption{w}
}

type traceOption struct {
	io.Writer
}
func (o traceOption) Apply(cfg *config) {
	cfg.trace = io.Writer(o)
}
func WithTraceWriter(w io.Writer) traceOption {
	return traceOption{w}
}