package metrics

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
)

type Config struct {
	Enabled bool
	Port    string
}

var mux = http.NewServeMux()

func Register(cfg Config, version string) {
	if !cfg.Enabled {
		slog.Debug("monitoring.Register: monitoring disabled")
		return
	}
	exporter, err := prometheus.New()
	if err != nil {
		slog.Warn(fmt.Sprintf("metric.Register: error creating Prometheus exporter: %v", err), slog.Any("error", err))
		return
	}

	meterProvider := metric.NewMeterProvider(metric.WithReader(exporter))
	otel.SetMeterProvider(meterProvider)

	slog.Info("metric.Register: Prometheus exporter enabled")
	mux.Handle("/metrics", promhttp.Handler())

	go func() {
		srv := &http.Server{
			Addr:              fmt.Sprintf(":%s", cfg.Port),
			Handler:           mux,
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 2 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       30 * time.Second,
			MaxHeaderBytes:    1 << 20, // 1 MB
		}

		if err := srv.ListenAndServe(); err != nil {
			slog.Error(fmt.Sprintf("monitoring.Server: %v", err), slog.Any("error", err))
		}
	}()
}
