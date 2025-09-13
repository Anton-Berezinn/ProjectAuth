package auth

import (
	"Project/internal/manager"
	"Project/internal/models"
	"Project/internal/repository/user/handler"
	j "Project/internal/token/jwt"
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"time"
)

func HandlerAuth(e *echo.Echo, project *manager.Project) {
	e.POST("/login", func(c echo.Context) error {
		return Login(c, project)
	})
	e.POST("/register", func(c echo.Context) error {
		return Register(c, project)
	})
}

func Register(c echo.Context, project *manager.Project) error {
	var req models.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error in server"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*1)
	defer cancel()

	userId, err := handler.RegisterHandler(ctx, &req, project)
	if err != nil {
		if errors.Is(err, handler.ErrDuplicateEmail) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Email already exists"})
		}
		if errors.Is(err, handler.ErrNotValidEmail) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid email"})
		}
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "serve Error"})
	}

	token, err := j.NewJwt(ctx, userId)
	if err != nil {
		project.Logger.Fatalln("Error to create token %s", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error in server"})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})

}

func Login(c echo.Context, project *manager.Project) error {
	var req models.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error in server"})
	}
	fmt.Println(req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*1)
	defer cancel()

	userId, err := handler.LoginHandler(ctx, &req, project)
	if err != nil {
		if errors.Is(err, handler.ErrNotValid) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid email or password"})
		}
		log.Println("Error to Login %s", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "serve Error"})
	}

	token, err := j.NewJwt(ctx, userId)
	if err != nil {
		project.Logger.Fatalln("Error to create token %s", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error in server"})
	}
	return c.JSON(http.StatusOK, map[string]string{"token": token})
}
