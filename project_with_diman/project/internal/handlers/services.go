package handlers

import (
	"Project/internal/handlers/auth"
	"Project/internal/manager"
	"Project/internal/metrics/prometheus"
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"time"
)

func StartEcho(ctx context.Context, port string, project *manager.Project) {

	e := echo.New()
	e.Server.ReadTimeout = 10 * time.Minute  // 60 * time.Second
	e.Server.WriteTimeout = 10 * time.Minute // 60 * time.Second

	e.Use(
		middleware.CORS(),
		middleware.GzipWithConfig(middleware.GzipConfig{
			Level: 6,
		}),
		middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
			LogURI:      true,
			LogStatus:   true,
			LogLatency:  true,
			LogRemoteIP: true,
			LogMethod:   true,
			LogError:    true,
			LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
				project.Logger.Printf("Request: method=%s, uri=%s, status=%d, latency=%v, ip=%s, error=%v",
					v.Method, v.URI, v.Status, v.Latency, v.RemoteIP, v.Error)
				return nil
			},
		}),
	)

	go prometheus.StartPrometheus(ctx)

	auth.HandlerAuth(e, project)
	fmt.Println("starting serve on port %s", port)
	if err := http.ListenAndServe(":"+port, e); err != nil {
	}

}
