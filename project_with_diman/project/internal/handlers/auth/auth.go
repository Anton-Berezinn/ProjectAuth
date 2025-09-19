package auth

import (
	"Project/internal/manager"
	"Project/internal/metrics/prometheus"
	"Project/internal/models"
	"Project/internal/repository/user/handler"
	"Project/internal/token/jwt"
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strconv"
	"time"
)

func HandlerAuth(e *echo.Echo, project *manager.Project) {
	e.POST("/login", func(c echo.Context) error {
		return Login(c, project)
	})
	e.POST("/register", func(c echo.Context) error {
		return Register(c, project)
	})
	e.GET("/f", func(c echo.Context) error { return CheckVersion(c, project) })
}

func CheckVersion(c echo.Context, project *manager.Project) error {
	t := c.Request().Header.Get(echo.HeaderAuthorization)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	j, err := jwt.DecodeJwt(ctx, t)
	if err != nil {
		fmt.Println(err)
	}
	version, err := project.Cache.GetData(ctx, strconv.Itoa(j.UserId))
	if err != nil {
		fmt.Println(err)
	}

	if strconv.Itoa(j.Version) != version {
		return c.JSON(401, "you need to login first")
	}

	return c.JSON(http.StatusOK, map[string]string{})
}

func Register(c echo.Context, project *manager.Project) error {
	defer prometheus.MyCounterRegister.Inc()

	var req models.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error in server"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*1)
	defer cancel()

	token, err := handler.RegisterHandler(ctx, &req, project)
	if err != nil {
		if errors.Is(err, handler.ErrDuplicateEmail) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Email already exists"})
		}
		if errors.Is(err, handler.ErrNotValidEmail) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid email"})
		}
		if errors.Is(err, handler.ErrNotValidName) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid name"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "serve Error"})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})

}

func Login(c echo.Context, project *manager.Project) error {
	defer prometheus.MyCounterLogin.Inc()
	var req models.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error in server"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*1)
	defer cancel()

	token, err := handler.LoginHandler(ctx, &req, project)
	if err != nil {
		if errors.Is(err, handler.ErrNotValid) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid email or password"})
		}
		log.Println("Error to Login %s", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "serve Error"})
	}
	return c.JSON(http.StatusOK, map[string]string{"token": token})
}
