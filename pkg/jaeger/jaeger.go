package jaeger

import (
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

func New(service, host string) (opentracing.Tracer, io.Closer, error) {
	headerCfg := jaeger.HeadersConfig{
		TraceContextHeaderName:   "uber-trace-id",
		TraceBaggageHeaderPrefix: "uberctx-",
	}
	cfg := &config.Configuration{
		ServiceName: service,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           false,
			LocalAgentHostPort: host,
		},
		Headers: &headerCfg,
	}

	textPropagator := jaeger.NewTextMapPropagator(&headerCfg, *jaeger.NewNullMetrics())
	return cfg.NewTracer(
		config.ZipkinSharedRPCSpan(true),
		config.Injector(opentracing.TextMap, textPropagator),
		config.Extractor(opentracing.TextMap, textPropagator),
	)
}
