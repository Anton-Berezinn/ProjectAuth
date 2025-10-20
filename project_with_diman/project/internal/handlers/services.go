package handlers

import (
	"Project/internal/handlers/auth"
	"Project/internal/handlers/catalog"
	"Project/internal/manager"
	"Project/internal/metrics/prometheus"
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"time"
	"io"
	"net"
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
	catalog.HandlerCatalogs(e, project)

	// external IP endpoint
	e.GET("/external-ip", func(c echo.Context) error {
		ctxReq, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
		defer cancel()

		// prefer a transport that respects timeouts
		transport := &http.Transport{
			DialContext: (&net.Dialer{Timeout: 3 * time.Second}).DialContext,
		}
		client := &http.Client{Transport: transport, Timeout: 5 * time.Second}

		req, err := http.NewRequestWithContext(ctxReq, http.MethodGet, "https://api.ipify.org?format=json", nil)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "request build failed"})
		}
		resp, err := client.Do(req)
		if err != nil {
			return c.JSON(http.StatusBadGateway, map[string]string{"error": "upstream failed"})
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "read upstream failed"})
		}
		return c.Blob(resp.StatusCode, "application/json", body)
	})
	fmt.Println("starting serve on port %s", port)
	if err := http.ListenAndServe(":"+port, e); err != nil {
	}

}
