package prometheus

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"time"
)

func recordMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(2 * time.Second)
		}
	}()
}

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})
	MyCounterRegister = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "my_custom_counter_register",
		Help: "Пример пользовательского регистрации",
	})
	MyCounterLogin = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "login_custom_counter",
		Help: "Пример пользовательского авторизации",
	})
	RequestDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "request_duration_seconds",
		Help:    "Время обработки запроса",
		Buckets: prometheus.LinearBuckets(0.0001, 2, 10), // 0.5s, 1s, 1.5s, ..., 5s
	})
)

func StartPrometheus(ctx context.Context) {
	recordMetrics()
	reg := prometheus.NewRegistry()

	reg.MustRegister(MyCounterRegister)
	reg.MustRegister(MyCounterLogin)
	reg.MustRegister(opsProcessed)
	reg.MustRegister(RequestDuration)

	handler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})

	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", handler)

		log.Println("Starting Prometheus metrics server on :2112/metrics")
		err := http.ListenAndServe(":2112", mux)
		if err != nil {
			log.Fatalf("Metrics server failed: %v", err)
		}
	}()

}
