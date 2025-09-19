package prometheus

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

var (
	MyCounterRegister = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "my_custom_counter_register",
		Help: "Пример пользовательского регистрации",
	})
	MyCounterLogin = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "login_custom_counter",
		Help: "Пример пользовательского авторизации",
	})
)

func StartPrometheus(ctx context.Context) {
	reg := prometheus.NewRegistry()

	reg.MustRegister(MyCounterRegister)
	reg.MustRegister(MyCounterLogin)

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
