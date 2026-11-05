// Package telemetry wires the process to the collector. It is the only place
// that knows about exporters, so the rest of the service depends on the
// OpenTelemetry API alone.
package telemetry

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"os"
	"runtime"

	"github.com/grafana/pyroscope-go"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	runtimeinst "go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// profileLogger reports failed profile uploads and drops the rest. Uploads run
// asynchronously, so a silent profiler is indistinguishable from a working one;
// its debug stream, on the other hand, buries every other line in the log.
type profileLogger struct{}

func (profileLogger) Debugf(string, ...any) {}
func (profileLogger) Infof(string, ...any)  {}

func (profileLogger) Errorf(format string, args ...any) {
	log.Printf("pyroscope: "+format, args...)
}

// Setup installs traces, metrics, logs and continuous profiling for service.
// Endpoint and protocol come from the standard OTEL_EXPORTER_OTLP_* variables,
// which the exporters read themselves, so the deployment configures telemetry
// without a code change. The returned function flushes every pipeline.
func Setup(ctx context.Context, service string) (func(context.Context) error, error) {
	traceExp, err := otlptracegrpc.New(ctx)
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(traceExp))
	otel.SetTracerProvider(tp)

	metricExp, err := otlpmetricgrpc.New(ctx)
	if err != nil {
		return nil, err
	}
	mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExp)))
	otel.SetMeterProvider(mp)

	logExp, err := otlploggrpc.New(ctx)
	if err != nil {
		return nil, err
	}
	lp := sdklog.NewLoggerProvider(sdklog.WithProcessor(sdklog.NewBatchProcessor(logExp)))
	slog.SetDefault(slog.New(otelslog.NewHandler(service, otelslog.WithLoggerProvider(lp))))

	// Goroutine and heap counters come from here; the container's own memory is
	// reported separately by the cluster collector.
	if err := runtimeinst.Start(runtimeinst.WithMeterProvider(mp)); err != nil {
		return nil, err
	}

	// Profiles go straight to Pyroscope, which speaks its own protocol rather
	// than OTLP.
	runtime.SetMutexProfileFraction(5)
	runtime.SetBlockProfileRate(5)
	prof, err := pyroscope.Start(pyroscope.Config{
		ApplicationName: service,
		ServerAddress:   os.Getenv("PYROSCOPE_SERVER_ADDRESS"),
		ProfileTypes:    pyroscope.DefaultProfileTypes,
		Logger:          profileLogger{},
	})
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) error {
		return errors.Join(
			prof.Stop(),
			tp.Shutdown(ctx),
			mp.Shutdown(ctx),
			lp.Shutdown(ctx),
		)
	}, nil
}
