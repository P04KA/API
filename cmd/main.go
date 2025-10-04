package main

import (
	"net/http"
	"os"

	"log/slog"

	"github.com/P04KA/API/internal/app"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "example_counter",
		Help: "An example counter metric",
	})

	// Регистрируем счетчик в реестре метрик
	prometheus.MustRegister(counter)

	// Инкрементируем счетчик
	counter.Inc()

	// Запускаем HTTP-сервер для экспорта метрик
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8082", nil)

	if err := app.Run(); err != nil {

		slog.Error("run app", slog.Any("error", err))
		os.Exit(1)
	}
}
